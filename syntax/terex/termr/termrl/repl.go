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
	"github.com/npillmayer/gotype/syntax/termr"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
)

func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

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
	gtrace.SyntaxTracer = gologadapter.New()
	//grammarName := flag.String("grammar", "G", "Name of the grammar to load")
	tlevel := flag.String("trace", "I", "Trace level")
	flag.Parse()
	T().Infof("Welcome to TermRL")
	ga := makeExprGrammar()
	T().SetTraceLevel(traceLevel(*tlevel))
	ga.Grammar().Dump()
	input := strings.Join(flag.Args(), " ")
	input = strings.TrimSpace(input)
	T().Infof("Input argument is \"%s\"", input)
	repl, err := readline.New("termrl> ")
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
		intp.tree, err = Parse(ga, input)
		if err != nil {
			T().Errorf("%v", err)
			os.Exit(2)
		}
	}
	intp.REPL()
}

type Intp struct {
	lastInput string
	lastValue interface{}
	GA        *lr.LRAnalysis
	repl      *readline.Instance
	tree      *sppf.Forest
	ast       *terex.GCons
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
		args := strings.Split(line, " ")
		cmd := args[0]
		err, quit := intp.Execute(cmd, args)
		if err != nil {
			T().Errorf(err.Error())
			continue
		}
		if quit {
			break
		}
	}
	println("Good bye!")
}

func (intp *Intp) Execute(cmd string, args []string) (error, bool) {
	switch cmd {
	case "quit", "bye":
		return nil, true
	case "parse":
		intp.ast = nil
		intp.lastInput = strings.TrimSpace(strings.Join(args[1:], " "))
		var err error
		intp.tree, err = Parse(intp.GA, intp.lastInput)
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
			intp.ast, intp.lastValue = astbuild.AST(intp.tree)
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

func Parse(ga *lr.LRAnalysis, input string) (*sppf.Forest, error) {
	level := T().GetTraceLevel()
	T().SetTraceLevel(tracing.LevelError)
	parser := earley.NewParser(ga, earley.GenerateTree(true))
	reader := strings.NewReader(input)
	scanner := scanner.GoTokenizer(ga.Grammar().Name, reader)
	acc, err := parser.Parse(scanner, nil)
	if !acc || err != nil {
		return nil, fmt.Errorf("could not parse input: %v", err)
	}
	T().SetTraceLevel(level)
	T().Infof("Successfully parsed input")
	return parser.ParseForest(), nil
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
