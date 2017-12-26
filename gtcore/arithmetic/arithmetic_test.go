package arithmetic

import (
	"testing"

	numeric "github.com/shopspring/decimal"
)

func TestNumericCompare(t *testing.T) {
	var x numeric.Decimal = numeric.New(1, -2)
	var y numeric.Decimal = numeric.New(1, -1)
	if x.LessThan(y) {
		t.Logf(" %s < %s", x.String(), y.String())
	} else {
		t.Errorf(" %s < %s", x.String(), y.String())
	}
}

/*
func TestNumericOrigin(t *testing.T) {
	x := MakePair(numeric.New(1, -10), numeric.New(0, 0))
	if !x.Equal(Origin) {
		t.Fail()
	}
	if !x.Zero() {
		t.Fail()
	}
}
*/
