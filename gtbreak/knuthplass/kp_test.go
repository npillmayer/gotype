package knuthplass

import (
	"testing"
	"fmt"
)

var llequal LL = func(line int64) int64 {
	return 50
}

func TestKPInit(t *testing.T) {
	text := "Because of these two problems, one major principle of Go's design is that the meaning of the program "
	text += "should be as transparent as possible. Your should be able to, as best possible, go to any source code "
	text += "and - after spending a few minutes studying it - immediately be able to start hacking away."

	kp := &KnuthPlassLinebreaker{}
	kp.beading = NewPseudoBeading(text)
	kp.lineLength = llequal
	fmt.Printf("bead = '%s'\n", kp.beading.GetBead().Text())
	kp.beading.Advance()
	fmt.Printf("bead = '%s'\n", kp.beading.GetBead().Text())
}