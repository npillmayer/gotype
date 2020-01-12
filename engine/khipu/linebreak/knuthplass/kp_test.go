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

func setupKhipu(t *testing.T, paragraph string) (*khipu.Khipu, linebreak.Cursor) {
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	kh := khipu.KnotEncode(strings.NewReader(paragraph), nil, regs)
	if kh == nil {
		t.Errorf("nokKhipu to test; input is %s", paragraph)
	}
	kh.AppendKnot(khipu.Penalty(linebreak.InfinityMerits))
	t.Logf("khipu=%s", kh.String())
	cursor := linebreak.NewFixedWidthCursor(khipu.NewCursor(kh), 10*dimen.BP)
	return kh, cursor
}

func TestKP1(t *testing.T) {
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelError)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	kh, cursor := setupKhipu(t, "The quick ")
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParshape(10 * 10 * dimen.BP)
	n, breaks, err := FindBreakpoints(cursor, parshape, true, nil)
	t.Logf("%d linebreaking-variants found, error = %v", n, err)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < n; i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	t.Fail()
}

func testKP2(t *testing.T) {
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelError)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	_, cursor := setupKhipu(t, "The quick brown fox jumps over the lazy dog!")
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParshape(10 * 10 * dimen.BP)
	n, breaks, err := FindBreakpoints(cursor, parshape, true, nil)
	t.Logf("%d linebreaking-variants found, error = %v", n, err)
	if err != nil || n != 2 || len(breaks[3]) != 2 {
		t.Fail()
	}
}
