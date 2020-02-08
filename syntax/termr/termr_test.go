package termr

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
)

func TestEnvSym(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	b := lr.NewGrammarBuilder("TermR")
	b.LHS("E").N("E").T("+", '+').T("var", scanner.Ident).End()
	b.LHS("E").T("var", scanner.Ident).End()
	G, _ := b.Grammar()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	env, err := EnvironmentForGrammarSymbol("E", G)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf(env.Dump())
	t.Logf(globalEnvironment.Dump())
	t.Fail()
}

func TestAST1(t *testing.T) {
	//gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer = gologadapter.New()
	//teardown := gotestingadapter.RedirectTracing(t)
	//defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	b := lr.NewGrammarBuilder("TermR")
	b.LHS("E").N("E").T("+", '+').T("a", scanner.Ident).End()
	b.LHS("E").T("a", scanner.Ident).End()
	G, _ := b.Grammar()
	ga := lr.Analysis(G)
	parser := earley.NewParser(ga, earley.GenerateTree(true))
	input := strings.NewReader("a+a")
	scanner := scanner.GoTokenizer("TestAST", input)
	acc, err := parser.Parse(scanner, nil)
	if !acc || err != nil {
		t.Errorf("parser could not parse input")
	}
	// tmpfile, _ := ioutil.TempFile(".", "tree-*.dot")
	// sppf.ToGraphViz(parser.ParseForest(), tmpfile)
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	builder := NewASTBuilder(G)
	ast, _ := builder.AST(parser.ParseForest())
	expected := `(a + a)`
	if ast.cdr == nil {
		t.Errorf("AST is empty")
	} else if ast.ListString() != expected {
		t.Errorf("AST should be %s, is %s", expected, ast.ListString())
	}
}

func TestAST2(t *testing.T) {
	//gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer = gologadapter.New()
	//teardown := gotestingadapter.RedirectTracing(t)
	//defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	b := lr.NewGrammarBuilder("TermR")
	b.LHS("E").N("E").T("+", '+').T("a", scanner.Ident).End()
	b.LHS("E").T("a", scanner.Ident).End()
	G, _ := b.Grammar()
	ga := lr.Analysis(G)
	parser := earley.NewParser(ga, earley.GenerateTree(true))
	input := strings.NewReader("a+a")
	scanner := scanner.GoTokenizer("TestAST", input)
	acc, err := parser.Parse(scanner, nil)
	if !acc || err != nil {
		t.Errorf("parser could not parse input")
	}
	// tmpfile, _ := ioutil.TempFile(".", "tree-*.dot")
	// sppf.ToGraphViz(parser.ParseForest(), tmpfile)
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	builder := NewASTBuilder(G)
	builder.AddOperator(makeOp("E"))
	ast, _ := builder.AST(parser.ParseForest())
	expected := `(a + a)`
	if ast.cdr == nil {
		t.Errorf("AST is empty")
	} else if ast.ListString() != expected {
		t.Errorf("AST should be %s, is %s", expected, ast.ListString())
	}
}
