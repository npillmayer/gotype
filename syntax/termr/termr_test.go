package termr

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/syntax/lr/earley"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
)

func TestList1(t *testing.T) {
	l1 := List(1, 4, 6, 8, 12)
	if l1.Length() != 5 {
		t.Errorf("length of list expected to be 5, is %d", l1.Length())
	}
	if l1.Car().Atom().Data != 1 {
		t.Errorf("element #1 expected to be 1, is %v", l1.Car().Atom().Data)
	}
}

func TestString(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	l2 := List("+", 5, 7, 12)
	l := Cons(1, nil)
	l.cdr = Cons(2, nil)
	x := &GCons{carNode{NullAtom, l2}, nil}
	l.cdr.cdr = x
	x.cdr = Cons("Hello", nil)
	s := l.ListString()
	t.Logf("l = %s", s)
	if s != `(1 2 ("+" 5 7 12) "Hello")` {
		t.Errorf(`Expected list to be (1 2 ("+" 5 7 12) "Hello"), is "%s"`, s)
	}
}

// TODO
func TestEval1(t *testing.T) {
	gtrace.SyntaxTracer = gologadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	input := "(Hello ^World 1)"
	env := &Environment{}
	env.Eval(input)
	t.Fail()
}

func TestAST(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
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
	builder := NewASTBuilder(G)
	ast, _ := builder.AST(parser.ParseForest())
	expected := `(a + a)`
	if ast.cdr == nil {
		t.Errorf("AST is empty")
	} else if ast.cdr.ListString() != expected {
		t.Errorf("AST should be %s, is %s", expected, ast.cdr.ListString())
	}
}
