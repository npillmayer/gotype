package knuthplass

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/core/dimen"
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
	//kh := khipu.KnotEncode(strings.NewReader("The quick brown fox jumps over the lazy dog!"), nil, regs)
	kh := khipu.KnotEncode(strings.NewReader("The quick brown "), nil, regs)
	if kh.Length() < 5 {
		t.Errorf("Length of khipu is too short; is %d", kh.Length())
	}
	kh.AppendKnot(khipu.Penalty(linebreak.InfinityMerits))
	t.Logf("khipu=%s", kh.String())
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	var cursor linebreak.Cursor
	cursor = khipu.NewCursor(kh)
	cursor = linebreak.NewFixedWidthCursor(cursor, 10*dimen.BP)
	parshape := linebreak.RectangularParshape(10 * 10 * dimen.BP)
	n, breaks := FindBreakpoints(cursor, parshape, true)
	t.Logf("%d linebreaking-variants found", n)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < len(breakpoints); i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	t.Fail()
}
