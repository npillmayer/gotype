package uax29

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/gtcore/uax/segment"
)

func Test0(t *testing.T) {
	TC.SetLevel(tracing.LevelError)
}

func ExampleWordBreaker() {
	onWords := NewWordBreaker()
	segmenter := segment.NewSegmenter(onWords)
	segmenter.Init(strings.NewReader("Hello World🇩🇪!"))
	for segmenter.Next() {
		fmt.Printf("'%s'\n", segmenter.Text())
	}
	// Output: 'Hello'
	// ' '
	// 'World'
	// '🇩🇪'
	// '!'
}

func TestWordBreakTestFile(t *testing.T) {
	//TC.SetLevel(tracing.LevelDebug)
	TC.SetLevel(tracing.LevelError)
	SetupUAX29Classes()
	onWordBreak := NewWordBreaker()
	seg := segment.NewSegmenter(onWordBreak)
	gopath := os.Getenv("GOPATH")
	f, err := os.Open(gopath + "/etc/WordBreakTest.txt")
	if err != nil {
		t.Errorf("ERROR loading " + gopath + "/etc/WordBreakTest.txt\n")
	}
	defer f.Close()
	failcnt, i, from, to := 0, 0, 1, 1900
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		line = strings.TrimSpace(line)
		if line[0] == '#' { // ignore comment lines
			continue
		}
		i++
		if i >= from {
			parts := strings.Split(line, "#")
			testInput, comment := parts[0], parts[1]
			TC.Infoln(comment)
			in, out := breakTestInput(testInput)
			if !executeSingleTest(t, seg, i, in, out) {
				failcnt++
			}
		}
		if i >= to {
			break
		}
	}
	if err := scan.Err(); err != nil {
		TC.Errorf("reading input:", err)
	}
	t.Logf("%d TEST CASES OUT of %d FAILED", failcnt, i-from+1)
}

func breakTestInput(ti string) (string, []string) {
	//fmt.Printf("breaking up %s\n", ti)
	sc := bufio.NewScanner(strings.NewReader(ti))
	sc.Split(bufio.ScanWords)
	out := make([]string, 0, 5)
	inp := bytes.NewBuffer(make([]byte, 0, 20))
	run := bytes.NewBuffer(make([]byte, 0, 20))
	for sc.Scan() {
		token := sc.Text()
		if token == "÷" {
			if run.Len() > 0 {
				out = append(out, run.String())
				run.Reset()
			}
		} else if token == "×" {
			// do nothing
		} else {
			n, _ := strconv.ParseUint(token, 16, 64)
			run.WriteRune(rune(n))
			inp.WriteRune(rune(n))
		}
	}
	//fmt.Printf("input = '%s'\n", inp.String())
	//fmt.Printf("output = %#v\n", out)
	return inp.String(), out
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
