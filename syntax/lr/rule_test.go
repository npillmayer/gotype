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
	gb := NewLRGrammar("G")
	b := gb.Builder()
	b.LHS("S").Epsilon()
	if len(gb.rules) != 1 {
		t.Fail()
	}
}
func TestClosure1(t *testing.T) {
	gb := NewLRGrammar("G")
	b := gb.Builder()
	r1 := b.LHS("S").N("E").EOF()
	r2 := b.LHS("E").N("E").T("+", 1).N("E").End()
	if len(gb.rules) != 2 {
		t.Fail()
	}
	item1 := &Item{r1, 0}
	item2 := &Item{r2, 1}
	T.Debug(item1)
	T.Debug(item2)
}
