package main

import (
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

func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

var stdEnv *terex.Environment

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

func main() {
	initDisplay()
	gtrace.SyntaxTracer = gologadapter.New()
	//grammarName := flag.String("grammar", "G", "Name of the grammar to load")
	tlevel := flag.String("trace", "I", "Trace level")
	flag.Parse()
	T().SetTraceLevel(tracing.LevelInfo)
	pterm.Info.Println("Welcome to TREPL")
	T().Infof("Trace level is %s", *tlevel)
	ga := makeExprGrammar()
	T().SetTraceLevel(traceLevel(*tlevel))
	ga.Grammar().Dump()
	input := strings.Join(flag.Args(), " ")
	input = strings.TrimSpace(input)
	T().Infof("Input argument is \"%s\"", input)
	env := initSymbols(ga)
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
	T().Infof("Quit with <ctrl>D")
	intp.REPL()
}

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

func initSymbols(ga *lr.LRAnalysis) *terex.Environment {
	terex.InitGlobalEnvironment()
	stdEnv = terexlang.LoadStandardLanguage()
	env := terex.NewEnvironment("trepl", stdEnv)
	// G is expression grammar (analyzed)
	env.Def("#ns#", terex.Atomize(env)) // TODO put this into "terex.NewEnvironment"
	env.Def("G", terex.Atomize(ga))
	makeTreeOps(env)
	return env
}

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

func (intp *Intp) REPL() {
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
		err, quit := intp.Eval(line)
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

func (intp *Intp) Eval(line string) (error, bool) {
	level := T().GetTraceLevel()
	T().SetTraceLevel(tracing.LevelError)
	tree, retr, err := terexlang.Parse(line)
	T().SetTraceLevel(level)
	if err != nil {
		//T().Errorf(err.Error())
		return err, false
	}
	T().SetTraceLevel(tracing.LevelError)
	ast, env, err := terexlang.AST(tree, retr)
	T().SetTraceLevel(level)
	// T().Infof("\n\n" + ast.IndentedListString() + "\n\n")
	// T().Infof("------------------------------------------------------------------")
	q, err := terexlang.QuoteAST(terex.Elem(ast.Car), env)
	T().SetTraceLevel(level)
	if err != nil {
		//T().Errorf(err.Error())
		return err, false
	}
	T().Infof("\n\n" + q.String() + "\n\n")
	//T().Infof(env.Dump())
	T().Infof("-------------------------- Output --------------------------------")
	//T().Infof(intp.env.Dump())
	result := terex.Eval(q, intp.env)
	intp.printResult(result, intp.env)
	intp.env.Error(nil)
	return nil, false
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
	pterm.Info.Println(result.String())
	return nil
}

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
	switch l {
	case "D":
		return tracing.LevelDebug
	case "I":
		return tracing.LevelInfo
	case "E":
		return tracing.LevelError
	}
	return tracing.LevelDebug
}
