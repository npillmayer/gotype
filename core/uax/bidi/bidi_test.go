package bidi

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"golang.org/x/text/unicode/bidi"
)

func TestScanner(t *testing.T) {
	gtrace.CoreTracer = gologadapter.New()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	input := "Hell\u0302o 吾輩は World!"
	reader := strings.NewReader(input)
	sc := NewScanner(reader)
	cnt := 0
	for {
		cnt++
		tokval, token, pos, _ := sc.NextToken(scanner.AnyToken)
		t.Logf("token '%s' at %d = %s", token, pos, ClassString(bidi.Class(tokval)))
		if tokval == scanner.EOF {
			break
		}
	}
	if cnt != 9 {
		t.Errorf("Expected scanner to return 9 tokens, was %d", cnt)
	}
}

func TestWeakTypes(t *testing.T) {
	gtrace.CoreTracer = gologadapter.New()
	T().SetTraceLevel(tracing.LevelDebug)
	gtrace.SyntaxTracer = gologadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	//input := "hell\u0302o"
	input := "123,45"
	scan := NewScanner(strings.NewReader(input), Testing(true))
	accept, err := Parse(scan)
	if err != nil {
		t.Error(err)
	}
	if !accept {
		t.Errorf("Input not recognized as a valid Bidi run")
	}
}
