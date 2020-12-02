package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/npillmayer/gotype/syntax/terex/terexlang"
	"github.com/npillmayer/gotype/syntax/terex/termr"
	"github.com/pterm/pterm"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
)

/*
BSD License

Copyright (c) 2019–21, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */

var stdEnv *terex.Environment

// We provide a simple expression grammar as a default for AST generation
// and rewriting experiments.
//
//  S'     ➞ Expr #eof
//  Expr   ➞ Expr SumOp Term  |  Term
//  Term   ➞ Term ProdOp Factor  |   Factor
//  Factor ➞ number  |   ( Expr )
//  SumOp  ➞ +  |  -
//  ProdOp ➞ *  |  /
//
func makeExprGrammar() *lr.LRAnalysis {
	level := T().GetTraceLevel()
	T().SetTraceLevel(tracing.LevelError)
	b := lr.NewGrammarBuilder("G")
	b.LHS("Expr").N("Expr").N("SumOp").N("Term").End()
	b.LHS("Expr").N("Term").End()
	b.LHS("Term").N("Term").N("ProdOp").N("Factor").End()
	b.LHS("Term").N("Factor").End()
	b.LHS("Factor").T("number", scanner.Int).End()
	b.LHS("Factor").T("(", '(').N("Expr").T(")", ')').End()
	b.LHS("SumOp").T("+", '+').End()
	b.LHS("SumOp").T("-", '-').End()
	b.LHS("ProdOp").T("*", '*').End()
	b.LHS("ProdOp").T("/", '/').End()
	g, err := b.Grammar()
	if err != nil {
		panic(fmt.Errorf("error creating grammar: %s", err.Error()))
	}
	T().SetTraceLevel(level)
	return lr.Analysis(g)
}

// main() starts an interactive CLI ("T.REPL"), where users may enter TeREx
// s-expressions. T.REPL will evaluate the s-expr and print out the result.
// T.REPL is intended as a sandbox for experiments during the early phase of
// parser/compiler/interpreter development, with a focus on term rewriting.
// It will allow users to test rewriting-expression for parser runs or for
// AST walking (AST = abstract syntax tree).
//
// Please refer to modules "terex" and "terexlang".
//
func main() {
	// set up logging
	initDisplay()
	gtrace.SyntaxTracer = gologadapter.New()
	tlevel := flag.String("trace", "Info", "Trace level [Debug|Info|Error]")
	initf := flag.String("init", "", "Initial load")
	flag.Parse()
	T().SetTraceLevel(tracing.LevelInfo)   // will set the correct level later
	pterm.Info.Println("Welcome to TREPL") // colored welcome message
	T().Infof("Trace level is %s", *tlevel)
	//
	// set up grammar and symbol-environment
	ga := makeExprGrammar()
	T().SetTraceLevel(traceLevel(*tlevel)) // now set the user supplied level
	ga.Grammar().Dump()                    // only visible in debug mode
	input := strings.Join(flag.Args(), " ")
	input = strings.TrimSpace(input)
	T().Infof("Input argument is \"%s\"", input)
	env := initSymbols(ga)
	//
	// set up REPL
	repl, err := readline.New("trepl> ")
	if err != nil {
		T().Errorf(err.Error())
		os.Exit(3)
	}
	intp := &Intp{
		lastInput: input,
		GA:        ga,
		repl:      repl,
		env:       env,
	}
	if input != "" {
		intp.tree, intp.tretr, err = Parse(ga, input)
		if err != nil {
			T().Errorf("%v", err)
			os.Exit(2)
		}
	}
	//
	// load an init file and start receiving commands / s-expressions
	T().Infof("Quit with <ctrl>D") // inform user how to stop the CLI
	intp.loadInitFile(*initf)      // init file name provided by flag
	intp.REPL()                    // go into interactive mode
}

// We use pterm for moderately fancy output.
func initDisplay() {
	pterm.EnableDebugMessages()
	pterm.Info.Prefix = pterm.Prefix{
		Text:  "  >>",
		Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
	}
	pterm.Error.Prefix = pterm.Prefix{
		Text:  "  Error",
		Style: pterm.NewStyle(pterm.BgRed, pterm.FgBlack),
	}
}

// Pre-load some symbols:
// G    = grammar information for demo expression grammar (see above)
// tree = a print-like tree command
//
func initSymbols(ga *lr.LRAnalysis) *terex.Environment {
	terex.InitGlobalEnvironment()
	stdEnv = terexlang.LoadStandardLanguage()
	env := terex.NewEnvironment("trepl", stdEnv)
	env.Def("#ns#", terex.Elem(env)) // TODO put this into "terex.NewEnvironment"
	// G is expression grammar (analyzed)
	env.Def("G", terex.Elem(ga))
	makeTreeOps(env)
	return env
}

// Intp is our interpreter object
type Intp struct {
	lastInput string
	lastValue interface{}
	GA        *lr.LRAnalysis
	repl      *readline.Instance
	tree      *sppf.Forest
	ast       *terex.GCons
	env       *terex.Environment
	tretr     termr.TokenRetriever
}

func (intp *Intp) loadInitFile(filename string) {
	if filename == "" {
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		T().Errorf("Unable to open init file: %s", filename)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineno := 1
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		line := scanner.Text()
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		_, err := intp.Eval(line)
		if err != nil {
			T().Errorf("Error line %d: "+err.Error(), lineno)
		}
		lineno++
	}
	if err := scanner.Err(); err != nil {
		T().Errorf("Error while reading init file: " + err.Error())
	}
}

// REPL starts interactive mode.
func (intp *Intp) REPL() {
	intp.env.Def("a", terex.Elem(7)) // pre-set for debugging purposes
	intp.env.Def("b", terex.Elem(terex.Cons(terex.Atomize(1), terex.Cons(terex.Atomize(2), nil))))
	intp.env.Def("c", terex.Elem("c")) // pre-set for debugging purposes
	for {
		line, err := intp.repl.Readline()
		if err != nil { // io.EOF
			break
		}
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		//println(line)
		// args := strings.Split(line, " ")
		// cmd := args[0]
		// err, quit := intp.Execute(cmd, args)
		quit, err := intp.Eval(line)
		if err != nil {
			//T().Errorf(err.Error())
			//pterm.Error.Println(err.Error())
			continue
		}
		if quit {
			break
		}
	}
	println("Good bye!")
}

// Eval evaluates a TeREx s-expr, given on a line by itself.
//
func (intp *Intp) Eval(line string) (bool, error) {
	T().Infof("----------------------- Parse & AST ------------------------------")
	level := T().GetTraceLevel()
	T().SetTraceLevel(tracing.LevelError)
	tree, retr, err := terexlang.Parse(line)
	T().SetTraceLevel(level)
	if err != nil {
		//T().Errorf(err.Error())
		return false, err
	}
	T().SetTraceLevel(tracing.LevelError)
	ast, env, err := terexlang.AST(tree, retr)
	T().SetTraceLevel(level)
	//T().Infof("\n\n" + ast.IndentedListString() + "\n\n")
	terex.Elem(ast).Dump(level)
	T().Infof("----------------------- Quoted AST -------------------------------")
	//q, err := terexlang.QuoteAST(terex.Elem(ast.Car), env)
	first := terex.Elem(ast).First()
	// first.Dump(tracing.LevelInfo)
	q, err := terexlang.QuoteAST(first, env)
	T().SetTraceLevel(level)
	if err != nil {
		//T().Errorf(err.Error())
		return false, err
	}
	//T().Infof("\n\n" + q.String() + "\n\n")
	q.Dump(level)
	//T().Infof(env.Dump())
	T().Infof("-------------------------- Output --------------------------------")
	//T().Infof(intp.env.Dump())
	result := terex.Eval(q, intp.env)
	intp.printResult(result, intp.env)
	intp.env.Error(nil)
	return false, nil
}

func (intp *Intp) printResult(result terex.Element, env *terex.Environment) error {
	if result.IsNil() {
		if env.LastError() != nil {
			pterm.Error.Println(stdEnv.LastError())
			return env.LastError()
		}
		//T().Infof("result: nil")
		pterm.Info.Println("nil")
		return nil
	}
	if result.AsAtom().Type() == terex.ErrorType {
		pterm.Error.Println(result.AsAtom().Data.(error).Error())
		return fmt.Errorf(result.AsAtom().Data.(error).Error())
	}
	//T().Infof(result.String())
	if env.LastError() != nil {
		pterm.Error.Println(env.LastError())
		return env.LastError()
	}
	result.Dump(tracing.LevelInfo)
	pterm.Info.Println(result.String())
	return nil
}

// Parse parses input for a given experimental grammar and returns a parse forest.
//
func Parse(ga *lr.LRAnalysis, input string) (*sppf.Forest, termr.TokenRetriever, error) {
	level := T().GetTraceLevel()
	T().SetTraceLevel(tracing.LevelError)
	parser := earley.NewParser(ga, earley.GenerateTree(true))
	reader := strings.NewReader(input)
	scanner := scanner.GoTokenizer(ga.Grammar().Name, reader)
	acc, err := parser.Parse(scanner, nil)
	if !acc || err != nil {
		return nil, nil, fmt.Errorf("could not parse input: %v", err)
	}
	T().SetTraceLevel(level)
	T().Infof("Successfully parsed input")
	tokretr := func(pos uint64) interface{} {
		return parser.TokenAt(pos)
	}
	return parser.ParseForest(), tokretr, nil
}

func makeTreeOps(env *terex.Environment) {
	// tree is a helper command to display an AST (abstract syntax tree) as a tree
	// on a terminal
	env.Defn("tree", func(e terex.Element, env *terex.Environment) terex.Element {
		// (tree T) => print tree representation of T
		T().Debugf("   e = %v", e.String())
		args := e.AsList()
		T().Debugf("args = %v", args.ListString())
		if args.Length() != 1 {
			return terexlang.ErrorPacker("Can only print tree for one symbol at a time", env)
		}
		T().Debugf("arg = %v", args.Car.String())
		first := args.Car
		label := "tree"
		T().Debugf("Atom = %v", first.Type())
		if args.Car.Type() == terex.VarType {
			label = first.Data.(*terex.Symbol).Name
		}
		arg := terex.Eval(terex.Elem(first), env)
		pterm.Println(label)
		root := indentedListFrom(arg, env)
		pterm.DefaultTree.WithRoot(root).Render()
		return terex.Elem(args.Car)
	})
}

func indentedListFrom(e terex.Element, env *terex.Environment) pterm.TreeNode {
	ll := leveledElem(e.AsList(), pterm.LeveledList{}, 0)
	T().Debugf("|ll| = %d, ll = %v", len(ll), ll)
	root := pterm.NewTreeFromLeveledList(ll)
	return root
}

func leveledElem(list *terex.GCons, ll pterm.LeveledList, level int) pterm.LeveledList {
	if list == nil {
		ll = append(ll, pterm.LeveledListItem{
			Level: level,
			Text:  "nil",
		})
		return ll
	}
	first := true
	for list != nil {
		if first {
			first = false // TODO modify level
		}
		if list.Car.Type() == terex.ConsType {
			ll = leveledElem(list.Car.Data.(*terex.GCons), ll, level+1)
		} else {
			ll = append(ll, pterm.LeveledListItem{
				Level: level,
				Text:  list.Car.String(),
			})
		}
		list = list.Cdr
	}
	return ll
}

func traceLevel(l string) tracing.TraceLevel {
	return tracing.TraceLevelFromString(l)
}
