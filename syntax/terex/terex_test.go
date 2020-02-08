package terex

import (
	"testing"

	"github.com/npillmayer/gotype/syntax/lr/sppf"

	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
)

func TestList1(t *testing.T) {
	l1 := List(1, 4, 6, 8, 12)
	if l1.Length() != 5 {
		t.Errorf("length of list expected to be 5, is %d", l1.Length())
	}
	if l1.Car().atom.Data != 1 {
		t.Errorf("element #1 expected to be 1, is %v", l1.Car().atom.Data)
	}
}
func TestFirst(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	l := List(1, 2, 3, 4, 5)
	f := l.First(3)
	if f.Length() != 3 || f.car.atom.Data != 1 {
		t.Errorf("Expected f to be (1 2 3), is %s", f.ListString())
	}
	l = List(1)
	f = l.First(1)
	if f.Length() != 1 || f.car.atom.Data != 1 {
		t.Errorf("Expected f to be (1), is %s", f.ListString())
	}
	l = List(1)
	f = l.First(5)
	if f.Length() != 1 || f.car.atom.Data != 1 {
		t.Errorf("Expected f to be (1), is %s", f.ListString())
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
	S := globalEnvironment.Intern("S", false)
	l1 := List("a", S, 2)
	l2 := List("a", 1, 2)
	t.Logf("l1 = %s, l2 = %s", l1.ListString(), l2.ListString())
	if !l1.Match(l2, globalEnvironment) {
		t.Errorf("l1 and l2 expected to match, don't")
	}
	t.Logf(globalEnvironment.Dump())
	if S.value.atom.typ != NumType { // S expected to be bound to 1
		t.Errorf("Symbol a expected to be of number type now, is %s", S.value.atom.typ.String())
	}
}

func TestMatch3(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	initDefaultPatterns()
	l1 := List(makeASTOp("S"), 1)
	t.Logf("l1 = %s, pattern = %s", l1.ListString(), SingleTokenArg.ListString())
	if !SingleTokenArg.Match(l1, globalEnvironment) {
		t.Errorf("Expected l1 to match pattern (Op any)")
	}
}

func TestMatch4(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	initDefaultPatterns()
	l := List(1, 2, 3)
	if !Anything.Match(l, globalEnvironment) {
		t.Errorf("Expected !Anything to match (1 2 3)")
	}
}

func TestString(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	l2 := List("+", 5, 7, 12)
	l := Cons(makeNode(1), nil)
	l.cdr = Cons(makeNode(2), nil)
	x := &GCons{Node{NullAtom, l2}, nil}
	l.cdr.cdr = x
	x.cdr = Cons(makeNode("Hello"), nil)
	s := l.ListString()
	t.Logf("l = %s", s)
	if s != `(1 2 ("+" 5 7 12) "Hello")` {
		t.Errorf(`Expected list to be (1 2 ("+" 5 7 12) "Hello"), is "%s"`, s)
	}
}

// TODO
func TestEval1(t *testing.T) {
	gtrace.SyntaxTracer = gologadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	input := "(Hello ^World 1)"
	globalEnvironment.Eval(input)
	t.Errorf("TODO: TestEval1 to be cleaned up !")
}

// ---------------------------------------------------------------------------

type testOp struct {
	name string
}

func (op *testOp) Rewrite(list *GCons, env *Environment) *GCons {
	T().Debugf(env.Dump())
	return list
}

func (op *testOp) Descend(sppf.RuleCtxt) bool {
	return true
}

func (op *testOp) Name() string {
	return op.name
}

func makeOp(name string) ASTOperator {
	return &testOp{
		name: name,
	}
}
