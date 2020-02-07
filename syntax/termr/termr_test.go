package termr

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/sppf"

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

func TestMatch1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	l1 := List("a", 1, 2)
	l2 := List("a", 1, 2)
	if !l1.Match(l2, globalEnvironment) {
		t.Errorf("l1 and l2 expected to match, don't")
	}
	l3 := List("a", 1, 2)
	l4 := List("a", "b", 2)
	if l3.Match(l4, globalEnvironment) {
		t.Errorf("l3 and l4 not expected to match, do")
	}
	l51 := List("b")
	l5 := List("a", l51, 2)
	l61 := List("b")
	l6 := List("a", l61, 2)
	t.Logf("l5 = %s, l6 = %s", l5.ListString(), l6.ListString())
	if !l5.Match(l6, globalEnvironment) {
		t.Errorf("l5 and l6 expected to match, don't")
	}
	t.Logf(globalEnvironment.Dump())
	t.Fail()
}

func TestMatch2(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	sym := globalEnvironment.Intern("S", false)
	l1 := List("a", sym, 2)
	l2 := List("a", 1, 2)
	t.Logf("l1 = %s, l2 = %s", l1.ListString(), l2.ListString())
	if !l1.Match(l2, globalEnvironment) {
		t.Errorf("l1 and l2 expected to match, don't")
	}
	t.Logf(globalEnvironment.Dump())
	t.Fail()
}

func TestString(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	l2 := List("+", 5, 7, 12)
	l := Cons(1, nil)
	l.cdr = Cons(2, nil)
	x := &GCons{Node{NullAtom, l2}, nil}
	l.cdr.cdr = x
	x.cdr = Cons("Hello", nil)
	s := l.ListString()
	t.Logf("l = %s", s)
	if s != `(1 2 ("+" 5 7 12) "Hello")` {
		t.Errorf(`Expected list to be (1 2 ("+" 5 7 12) "Hello"), is "%s"`, s)
	}
}

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

// TODO
func TestEval1(t *testing.T) {
	gtrace.SyntaxTracer = gologadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	input := "(Hello ^World 1)"
	globalEnvironment.Eval(input)
	t.Errorf("TODO: TestEval1 to be cleaned up !")
}

func TestAST(t *testing.T) {
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

func echoRewrite(list *GCons, env *Environment) *GCons {
	T().Debugf(env.Dump())
	return list
}

func alwaysDescend(sppf.RuleCtxt) bool {
	return true
}

func makeOp(name string) *ASTOperator {
	return &ASTOperator{
		Name:    name,
		Rewrite: echoRewrite,
		Descend: alwaysDescend,
	}
}
