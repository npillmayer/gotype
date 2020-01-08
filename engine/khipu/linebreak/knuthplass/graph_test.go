package knuthplass

import "testing"

func TestGraph1(t *testing.T) {
	g := newLinebreaker(nil, nil)
	g.Add(g.newBreakpointAfterKnot(provisionalMark(1)))
	if g.Breakpoint(1) == nil {
		t.Errorf("Expected to find breakpoint at %d in graph, is nil", 1)
	}
}
