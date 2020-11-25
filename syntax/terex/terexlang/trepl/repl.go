package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
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
	T().Infof("Welcome to TermRL")
	T().Infof("Trace level is %s", *tlevel)
	ga := makeExprGrammar()
	T().SetTraceLevel(traceLevel(*tlevel))
	ga.Grammar().Dump()
	input := strings.Join(flag.Args(), " ")
	input = strings.TrimSpace(input)
	T().Infof("Input argument is \"%s\"", input)
	repl, err := readline.New("trepl> ")
	if err != nil {
		T().Errorf(err.Error())
		os.Exit(3)
	}
	intp := &Intp{
		lastInput: input,
		GA:        ga,
		repl:      repl,
	}
	if input != "" {
		intp.tree, intp.tretr, err = Parse(ga, input)
		if err != nil {
			T().Errorf("%v", err)
			os.Exit(2)
		}
	}
	T().Infof("Quit with <ctrl>D")
	stdEnv = terexlang.LoadStandardLanguage()
	intp.REPL()
}

func initDisplay() {
	pterm.EnableDebugMessages()
	// Customize default error.
	pterm.Error.Prefix = pterm.Prefix{
		Text:  "  Error",
		Style: pterm.NewStyle(pterm.BgLightRed, pterm.FgBlack),
	}
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
			T().Errorf(err.Error())
			pterm.Error.Println(err.Error())
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
	terex.InitGlobalEnvironment()
	tree, retr, err := terexlang.Parse(line)
	T().SetTraceLevel(level)
	if err != nil {
		//T().Errorf(err.Error())
		return err, false
	}
	T().SetTraceLevel(tracing.LevelError)
	ast, env, err := terexlang.AST(tree, retr)
	T().SetTraceLevel(level)
	T().Infof("\n\n" + ast.IndentedListString() + "\n\n")
	T().Infof("------------------------------------------------------------------")
	q, err := terexlang.QuoteAST(ast.Tee(), env)
	T().SetTraceLevel(level)
	if err != nil {
		//T().Errorf(err.Error())
		return err, false
	}
	T().Infof("\n\n" + q.IndentedListString() + "\n\n")
	T().Infof(env.Dump())
	T().Infof("==================================================================")
	result := terex.Eval(terex.Elem(q), stdEnv)
	if result.AsAtom().Type() == terex.ErrorType {
		return fmt.Errorf(result.AsAtom().Data.(error).Error()), false
	}
	T().Infof(result.String())
	return nil, false
}

func (intp *Intp) Execute(cmd string, args []string) (error, bool) {
	switch cmd {
	case "quit", "bye":
		return nil, true
	case "parse":
		intp.ast = nil
		intp.lastInput = strings.TrimSpace(strings.Join(args[1:], " "))
		var err error
		intp.tree, intp.tretr, err = Parse(intp.GA, intp.lastInput)
		return err, false
	case "dot":
		if intp.tree == nil {
			return errors.New("No parse tree present"), false
		}
		return ExportTree(intp.tree), false
	case "ast":
		if intp.ast == nil {
			if intp.tree == nil {
				return errors.New("No parse tree present"), false
			}
			astbuild := termr.NewASTBuilder(intp.GA.Grammar())
			intp.env = astbuild.AST(intp.tree, intp.tretr)
		}
		out := intp.ast.ListString()
		println(out)
	}
	return nil, false
}

func ExportTree(tree *sppf.Forest) error {
	tmpfile, err := ioutil.TempFile(".", "tree-*.dot")
	if err != nil {
		T().Errorf("Cannot create tmp-fiile in local directory")
		return err
	}
	defer tmpfile.Close()
	sppf.ToGraphViz(tree, tmpfile)
	T().Infof("Exported parse tree to %s", tmpfile.Name())
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
