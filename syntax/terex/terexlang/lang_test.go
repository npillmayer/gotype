package terexlang

import (
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/syntax/terex"
)

func TestMatch3(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	initDefaultPatterns()
	l1 := terex.List(makeASTOp("S"), 1)
	t.Logf("l1 = %s, pattern = %s", l1.ListString(), SingleTokenArg.ListString())
	if !SingleTokenArg.Match(l1, terex.GlobalEnvironment) {
		t.Errorf("Expected l1 to match pattern (Op any)")
	}
}

func TestMatch4(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	initDefaultPatterns()
	l := terex.List(1, 2, 3)
	if !Anything.Match(l, terex.GlobalEnvironment) {
		t.Errorf("Expected !Anything to match (1 2 3)")
	}
}

// TODO
func TestEval1(t *testing.T) {
	gtrace.SyntaxTracer = gologadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	input := "(Hello ^World 1)"
	Eval(input, terex.GlobalEnvironment)
	t.Errorf("TODO: TestEval1 to be cleaned up !")
}
