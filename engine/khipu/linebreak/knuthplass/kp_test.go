package knuthplass

import (
	"io"
	"io/ioutil"
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

var graphviz = false // globally switches GraphViz output on/off

func TestGraph1(t *testing.T) {
	gtrace.CoreTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParshape(10 * 10 * dimen.BP)
	g := newLinebreaker(parshape, nil)
	g.newBreakpointAtMark(provisionalMark(1))
	if g.Breakpoint(1) == nil {
		t.Errorf("Expected to find breakpoint at %d in graph, is nil", 1)
	}
}

func setupKPTest(t *testing.T, paragraph string) (*khipu.Khipu, linebreak.Cursor, io.Writer) {
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	kh := khipu.KnotEncode(strings.NewReader(paragraph), nil, regs)
	if kh == nil {
		t.Errorf("no Khipu to test; input is %s", paragraph)
	}
	kh.AppendKnot(khipu.Penalty(linebreak.InfinityMerits))
	gtrace.CoreTracer.Infof("input khipu=%s", kh.String())
	cursor := linebreak.NewFixedWidthCursor(khipu.NewCursor(kh), 10*dimen.BP)
	var dotfile io.Writer
	var err error
	if graphviz {
		dotfile, err = ioutil.TempFile(".", "knuthplass-*.dot")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	return kh, cursor, dotfile
}

func TestKPUnderfull(t *testing.T) {
	gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	kh, cursor, dotfile := setupKPTest(t, " ")
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParshape(10 * 10 * dimen.BP)
	n, breaks, err := FindBreakpoints(cursor, parshape, true, nil, dotfile)
	t.Logf("%d linebreaking-variants for empty line found, error = %v", n, err)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < n; i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	if err != nil || n != 1 || len(breaks[1]) != 2 {
		t.Fail()
	}
}

func TestKPExactFit(t *testing.T) {
	gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	kh, cursor, dotfile := setupKPTest(t, "The quick.")
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParshape(10 * 10 * dimen.BP)
	n, breaks, err := FindBreakpoints(cursor, parshape, true, nil, dotfile)
	t.Logf("%d linebreaking-variants found, error = %v", n, err)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < n; i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	if err != nil || n != 1 || len(breaks[1]) != 2 {
		t.Fail()
	}
}

func TestKPOverfull(t *testing.T) {
	gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	kh, cursor, dotfile := setupKPTest(t, "The quick brown fox.")
	params := NewKPDefaultParameters()
	params.EmergencyStretch = dimen.Dimen(0)
	params.Tolerance = 400
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParshape(10 * 10 * dimen.BP)
	n, breaks, err := FindBreakpoints(cursor, parshape, true, params, dotfile)
	t.Logf("%d linebreaking-variants found, error = %v", n, err)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < n; i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	if err != nil || n != 1 || len(breaks[2]) != 3 {
		t.Fail()
	}
}
