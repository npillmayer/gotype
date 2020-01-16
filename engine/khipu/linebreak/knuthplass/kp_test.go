package knuthplass

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
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
	parshape := linebreak.RectangularParShape(10 * 10 * dimen.BP)
	g := newLinebreaker(parshape, nil)
	g.newBreakpointAtMark(provisionalMark(1))
	if g.Breakpoint(1) == nil {
		t.Errorf("Expected to find breakpoint at %d in graph, is nil", 1)
	}
}

func setupKPTest(t *testing.T, paragraph string, hyphens bool) (*khipu.Khipu, linebreak.Cursor, io.Writer) {
	regs := parameters.NewTypesettingRegisters()
	if hyphens {
		regs.Push(parameters.P_MINHYPHENLENGTH, 3) // allow hyphenation
	} else {
		regs.Push(parameters.P_MINHYPHENLENGTH, 100) // inhibit hyphenation
	}
	kh := khipu.KnotEncode(strings.NewReader(paragraph), nil, regs)
	if kh == nil {
		t.Errorf("no Khipu to test; input is %s", paragraph)
	}
	kh.AppendKnot(khipu.Penalty(linebreak.InfinityMerits))
	gtrace.CoreTracer.Infof("input khipu=%s", kh.String())
	cursor := linebreak.NewFixedWidthCursor(khipu.NewCursor(kh), 10*dimen.BP, 0)
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
	kh, cursor, dotfile := setupKPTest(t, " ", false)
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParShape(10 * 10 * dimen.BP)
	v, breaks, err := FindBreakpoints(cursor, parshape, nil, dotfile)
	t.Logf("%d linebreaking-variants for empty line found, error = %v", len(v), err)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < len(v); i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	if err != nil || len(v) != 1 || len(breaks[1]) != 2 {
		t.Fail()
	}
}

func TestKPExactFit(t *testing.T) {
	gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	kh, cursor, dotfile := setupKPTest(t, "The quick.", false)
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParShape(10 * 10 * dimen.BP)
	v, breaks, err := FindBreakpoints(cursor, parshape, nil, dotfile)
	t.Logf("%d linebreaking-variants found, error = %v", len(v), err)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < len(v); i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	if err != nil || len(v) != 1 || len(breaks[1]) != 2 {
		t.Fail()
	}
}

func TestKPOverfull(t *testing.T) {
	gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	kh, cursor, dotfile := setupKPTest(t, "The quick brown fox.", false)
	params := NewKPDefaultParameters()
	params.EmergencyStretch = dimen.Dimen(0)
	params.Tolerance = 400
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	parshape := linebreak.RectangularParShape(10 * 10 * dimen.BP)
	v, breaks, err := FindBreakpoints(cursor, parshape, params, dotfile)
	t.Logf("%d linebreaking-variants found, error = %v", len(v), err)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < len(v); i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	if err != nil || len(v) != 1 || len(breaks[2]) != 3 {
		t.Fail()
	}
}

var princess = `In olden times when wishing still helped one, there lived a king whose daughters were all beautiful; and the youngest was so beautiful that the sun itself, which has seen so much, was astonished whenever it shone in her face. Close by the king's castle lay a great dark forest, and under an old lime-tree in the forest was a well, and when the day was very warm, the king's child went out into the forest and sat down by the side of the cool fountain; and when she was bored she took a golden ball, and threw it up on high and caught it; and this ball was her favorite plaything.`
var king = `In olden times when wishing still helped one, there lived a king`

func TestKPParaKing(t *testing.T) {
	//gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer = gologadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	// teardown := gotestingadapter.RedirectTracing(t)
	// defer teardown()
	kh, _, dotfile := setupKPTest(t, king, false)
	cursor := linebreak.NewFixedWidthCursor(khipu.NewCursor(kh), 10*dimen.BP, 3)
	params := NewKPDefaultParameters()
	parshape := linebreak.RectangularParShape(45 * 10 * dimen.BP)
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	v, breaks, err := FindBreakpoints(cursor, parshape, params, dotfile)
	t.Logf("%d linebreaking-variants found, error = %v", len(v), err)
	for linecnt, breakpoints := range breaks {
		t.Logf("# Paragraph with %d lines: %v", linecnt, breakpoints)
		j := 0
		for i := 1; i < len(v); i++ {
			t.Logf(": %s", kh.Text(j, breakpoints[i].Position()))
			j = breakpoints[i].Position()
		}
	}
	if err != nil || len(v) != 1 || len(breaks[2]) != 3 {
		t.Fail()
	}
}

func TestKPParaPrincess(t *testing.T) {
	//gtrace.CoreTracer = gotestingadapter.New()
	gtrace.CoreTracer = gologadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	// teardown := gotestingadapter.RedirectTracing(t)
	// defer teardown()
	kh, _, _ := setupKPTest(t, princess, false)
	// change to cursor with flexible interword-spacing
	cursor := linebreak.NewFixedWidthCursor(khipu.NewCursor(kh), 10*dimen.BP, 2)
	params := NewKPDefaultParameters()
	parshape := linebreak.RectangularParShape(45 * 10 * dimen.BP)
	//gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	breakpoints, err := BreakParagraph(cursor, parshape, params)
	//v, breaks, err := FindBreakpoints(cursor, parshape, params, dotfile)
	//t.Logf("%d linebreaking-variants found, error = %v", len(v), err)
	t.Logf("# Paragraph with %d lines: %v", len(breakpoints)-1, breakpoints)
	t.Logf("    |---------+---------+---------+---------+-----|")
	j := 0
	for i := 1; i < len(breakpoints); i++ {
		//	t.Logf("%3d: %s", i, kh.Text(j, breakpoints[i].Position()))
		text := kh.Text(j, breakpoints[i].Position())
		t.Logf("%3d: %-45s|", i, justify(text, 45, i%2 == 0))
		j = breakpoints[i].Position()
	}
	if err != nil {
		t.Fail()
	}
	t.Fail()
}

// crude implementation just for testing
func justify(text string, l int, even bool) string {
	t := strings.Trim(text, " \t\n")
	d := l - len(t)
	if d == 0 {
		return t // fit
	} else if d < 0 { // overfull box
		return text + "\u25ae"
	}
	s := strings.Fields(text)
	if len(s) == 1 {
		return text
	}
	var b bytes.Buffer
	W := 0 // length of all words
	for _, w := range s {
		W += len(w)
	}
	d = l - W // amount of WS to distribute
	ws := d / (len(s) - 1)
	r := d - ws*(len(s)-1) + 1
	b.WriteString(s[0])
	if even {
		for j := 1; j < r; j++ {
			for i := 0; i < ws+1; i++ {
				b.WriteString(" ")
			}
			b.WriteString(s[j])
		}
		for j := r; j < len(s); j++ {
			for i := 0; i < ws; i++ {
				b.WriteString(" ")
			}
			b.WriteString(s[j])
		}
	} else {
		for j := 1; j <= len(s)-r; j++ {
			for i := 0; i < ws; i++ {
				b.WriteString(" ")
			}
			b.WriteString(s[j])
		}
		for j := len(s) - r + 1; j < len(s); j++ {
			for i := 0; i < ws+1; i++ {
				b.WriteString(" ")
			}
			b.WriteString(s[j])
		}
	}
	return b.String()
}
