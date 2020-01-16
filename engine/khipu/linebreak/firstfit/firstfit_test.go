package firstfit

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/core/parameters"
	"github.com/npillmayer/gotype/engine/khipu"
	"github.com/npillmayer/gotype/engine/khipu/linebreak"
)

var graphviz = false // global switch for GraphViz DOT output

func TestBuffer(t *testing.T) {
	//gtrace.CoreTracer = gologadapter.New()
	gtrace.CoreTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	lb := newTestLinebreaker(t, "Hello World!")
	k, ok := lb.peek()
	if !ok || k.Type() != khipu.KTTextBox {
		t.Logf("lb.pos=%d, lb.mark=%d", lb.pos, lb.mark)
		t.Errorf("expected the first knot to be TextBox('Hello'), is %v", k)
	}
	if knot := lb.next(); k != knot {
		t.Errorf("first knot is %v, re-read knot is %v", k, knot)
	}
}

func TestBacktrack(t *testing.T) {
	//gtrace.CoreTracer = gologadapter.New()
	gtrace.CoreTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	lb := newTestLinebreaker(t, "the quick brown fox jumps over the lazy dog.")
	k := lb.next()
	lb.checkpoint()
	lb.next()
	lb.next()
	lb.next()
	knot := lb.backtrack()
	if k != knot {
		t.Errorf("remembered start knot is %v, backtracked knot is %v", k, knot)
	}
	k = lb.next()
	lb.checkpoint()
	knot = lb.backtrack()
	if k != knot {
		t.Errorf("remembered knot is %v, backtracked knot is %v", k, knot)
	}
}

func newTestLinebreaker(t *testing.T, text string) *linebreaker {
	_, cursor, _ := setupFFTest(t, text, false)
	parshape := linebreak.RectangularParShape(10 * 10 * dimen.BP)
	lb, err := newLinebreaker(cursor, parshape, nil)
	if err != nil {
		t.Error(err)
	}
	return lb
}

func setupFFTest(t *testing.T, paragraph string, hyphens bool) (*khipu.Khipu, linebreak.Cursor, io.Writer) {
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
	//kh.AppendKnot(khipu.Penalty(linebreak.InfinityMerits))
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
