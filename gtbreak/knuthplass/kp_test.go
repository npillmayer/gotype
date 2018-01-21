package knuthplass

import (
	"fmt"
	"testing"
)

var llequal LL = func(line int64) int64 {
	return 30
}

func TestKPInit(t *testing.T) {
	text := "Because of these two problems, one major principle of Go's design is that the meaning of the program "
	text += "should be as transparent as possible. Your should be able to, as best possible, go to any source code "
	text += "and - after spending a few minutes studying it - immediately be able to start hacking away."

	kp := NewKPLinebreaker()
	kp.beading = newPseudoBeading(text)
	kp.lineLength = llequal
	cursor := kp.beading.GetCursor(nil)
	bead := cursor.GetBead()
	fmt.Printf("bead = [%5s] '%s'\n", beadtype(bead.BeadType()), bead.Text())
	cursor.Advance()
	bead = cursor.GetBead()
	fmt.Printf("bead = [%5s] '%s'\n", beadtype(bead.BeadType()), bead.Text())
	cursor.Advance()
	bead = cursor.GetBead()
	fmt.Printf("bead = [%5s] '%s'\n", beadtype(bead.BeadType()), bead.Text())
	cursor.Advance()
	bead = cursor.GetBead()
	fmt.Printf("bead = [%5s] '%s'\n", beadtype(bead.BeadType()), bead.Text())
}

func TestKPCost(t *testing.T) {
	text := "Because of these two problems, one major principle of Go's design is that the meaning of the program "
	kp := NewKPLinebreaker()
	kp.beading = newPseudoBeading(text)
	kp.lineLength = llequal
	cursor := kp.beading.GetCursor(nil)
	cursor.Advance()
	cursor.Advance()
	cursor.Advance()
	cursor.Advance()
	cursor.Advance()
	cursor.Advance()
	cursor.Advance()
	cursor.Advance()
	cursor.Advance()
	bead := cursor.GetBead()
	fmt.Printf("bead = [%5s] '%s'\n", beadtype(bead.BeadType()), bead.Text())
	f1 := kp.breakpointAfterBead(kp.beading.GetCursor(cursor))
	fmt.Printf("feasible breakpoint %v\n", f1)
	f1.fragment = WSS{20, 18, 25}
	cost, _ := f1.calculateCostTo(bead, 30)
	fmt.Printf("cost = %d, frag = %v\n", cost, f1.fragment)
}

func TestKP1(t *testing.T) {
	fmt.Println("------------------------------------------------------------")
	text := "Because of these two problems, one major principle of Go's design is that the meaning of the program "
	kp := NewKPLinebreaker()
	kp.beading = newPseudoBeading(text)
	kp.lineLength = llequal
	kp.FindBreakpoints(false)
}

func TestKP2(t *testing.T) {
	fmt.Println("------------------------------------------------------------")
	text := "Because of these two problems, one major principle of Go's design is that the meaning of the program "
	kp := NewKPLinebreaker()
	kp.beading = newPseudoBeading(text)
	kp.lineLength = llequal
	kp.FindBreakpoints(true)
	kp.MarshalToDotFile("test2")
}
