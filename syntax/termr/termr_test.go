package termr

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
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
	t.Logf(terex.GlobalEnvironment.Dump())
	if env.FindSymbol("E", true) == nil {
		t.Errorf("Expected to have symbol E in environment")
	}
}

func TestAST1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
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
	expected := `(:a :+ :a :#eof)`
	if ast.Cdr == nil {
		t.Errorf("AST is empty")
	} else if ast.ListString() != expected {
		t.Errorf("AST should be %s, is %s", expected, ast.ListString())
	}
}

func TestAST2(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
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
	builder.AddTermR(makeOp("E"))
	ast, _ := builder.AST(parser.ParseForest())
	expected := `((#E (#E :a) :+ :a) :#eof)`
	if ast.Cdr == nil {
		t.Errorf("AST is empty")
	} else if ast.ListString() != expected {
		t.Errorf("AST should be %s, is %s", expected, ast.ListString())
	}
}

// ---------------------------------------------------------------------------

type testOp struct {
	name string
}

func (op *testOp) Rewrite(list *terex.GCons, env *terex.Environment) *terex.GCons {
	T().Debugf(env.Dump())
	return list
}

func (op *testOp) Descend(sppf.RuleCtxt) bool {
	return true
}

func (op *testOp) Name() string {
	return op.name
}

func (op *testOp) String() string {
	return op.name
}

func (op *testOp) Operator() terex.Operator {
	return op
}

func (op *testOp) Call(el terex.Element, env *terex.Environment) terex.Element {
	return terex.Elem(nil)
}

func (op *testOp) Quote(el terex.Element, env *terex.Environment) terex.Element {
	return el
}

func makeOp(name string) *testOp {
	return &testOp{
		name: name,
	}
}

var _ terex.Operator = makeOp("Hello")
