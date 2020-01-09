package knuthplass

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/core/parameters"
	"github.com/npillmayer/gotype/engine/khipu"
	"github.com/npillmayer/gotype/engine/khipu/linebreak"
)

func init() {
	gtrace.CoreTracer = gotestingadapter.New()
}

func TestKP1(t *testing.T) {
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelError)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	kh := khipu.KnotEncode(strings.NewReader("The quick brown fox jumps over the lazy dog!"), nil, regs)
	if kh.Length() < 5 {
		t.Errorf("Length of khipu is too short; is %d", kh.Length())
	}
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	n, breaks := FindBreakpoints(khipu.NewCursor(kh), linebreak.RectangularParshape(20), nil, true)
	t.Logf("%d breakpoints found: %v", n, breaks)
	j := 0
	for i := 0; i < n; i++ {
		t.Logf(": %s", kh.Text(j, breaks[i].Position()))
		j = breaks[i].Position()
	}
	t.Fail()
}
