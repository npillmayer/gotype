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
	k := khipu.KnotEncode(strings.NewReader("Hello world!"), nil, regs)
	if k.Length() < 5 {
		t.Errorf("Length of khipu is too short; is %d", k.Length())
	}
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	FindBreakpoints(khipu.NewCursor(k), linebreak.RectangularParshape(20), nil, true)
}
