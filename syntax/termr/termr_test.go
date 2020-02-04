package termr

import (
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
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

func TestEval1(t *testing.T) {
	gtrace.SyntaxTracer = gologadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	input := "(Hello ^World 1)"
	env := &Environment{}
	env.Eval(input)
	t.Fail()
}
