package lr

import (
	"testing"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
)

func traceOn() {
	T.SetLevel(tracing.LevelDebug)
}

func TestBuilder1(t *testing.T) {
	traceOn()
	gb := NewLRGrammar("G")
	b := gb.Builder()
	b.LHS("S").N("A").End()
	if len(gb.rules) != 1 {
		t.Fail()
	}
}
func TestBuilder2(t *testing.T) {
	traceOn()
	gb := NewLRGrammar("G")
	b := gb.Builder()
	b.LHS("S").Epsilon()
	if len(gb.rules) != 1 {
		t.Fail()
	}
}
