package knuthplass

import (
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
)

func init() {
	gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
}

func TestGraph1(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	g := newLinebreaker(nil)
	g.newBreakpointAtMark(provisionalMark(1))
	if g.Breakpoint(1) == nil {
		t.Errorf("Expected to find breakpoint at %d in graph, is nil", 1)
	}
}
