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

func TestKP1(t *testing.T) {
	gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelError)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	khipu := khipu.KnotEncode(strings.NewReader("Hello world!"), nil, regs)
	if khipu.Length() < 5 {
		t.Errorf("Length of khipu is too short; is %d", khipu.Length())
	}
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	FindBreakpoints(linebreak.KhipuCursor(khipu), linebreak.RectangularParshape(20), nil, true)
}
