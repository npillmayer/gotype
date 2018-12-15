package uax14_test

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/uax/segment"
	"github.com/npillmayer/gotype/core/uax/uax14"
	"github.com/npillmayer/gotype/core/uax/ucd"
)

//var TC tracing.Trace = tracing.CoreTracer
var TC tracing.Trace = gologadapter.New()

func Test0(t *testing.T) {
	t.Log("setting core tracer")
	TC.SetTraceLevel(tracing.LevelDebug)
	tracing.CoreTracer = TC
	TC.Infof("Core tracer alive")
	TC.SetTraceLevel(tracing.LevelError)
}

func TestWordBreakTestFile(t *testing.T) {
	linewrap := uax14.NewLineWrap()
	seg := segment.NewSegmenter(linewrap)
	tf := ucd.OpenTestFile("./LineBreakTest.txt", t)
	defer tf.Close()
	//failcnt, i, from, to := 0, 0, 6263, 6263
	failcnt, i, from, to := 0, 0, 1, 8000
	for tf.Scan() {
		i++
		if i >= from {
			TC.Infof(tf.Comment())
			in, out := ucd.BreakTestInput(tf.Text())
			if !executeSingleTest(t, seg, i, in, out) {
				failcnt++
			}
		}
		if i >= to {
			break
		}
	}
	if err := tf.Err(); err != nil {
		TC.Errorf("reading input:", err)
	}
	t.Logf("%d TEST CASES OUT of %d FAILED", failcnt, i-from+1)
}

func executeSingleTest(t *testing.T, seg *segment.Segmenter, tno int, in string, out []string) bool {
	seg.Init(strings.NewReader(in))
	i := 0
	ok := true
	for seg.Next() {
		if len(out) <= i {
			t.Errorf("test #%d: number of segments too large: %d > %d", tno, i+1, len(out))
		} else if out[i] != seg.Text() {
			t.Errorf("test #%d: '%+q' should be '%+q'", tno, seg.Bytes(), out[i])
			ok = false
		}
		i++
	}
	return ok
}
