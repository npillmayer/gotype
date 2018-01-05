package arithmetic

import (
	"testing"

	dec "github.com/shopspring/decimal"
)

var two dec.Decimal
var three dec.Decimal
var six dec.Decimal
var seven dec.Decimal

func init() {
	two = dec.New(2, 0)
	three = dec.New(3, 0)
	six = dec.New(6, 0)
	seven = dec.New(7, 0)
}

func TestNumericCompare(t *testing.T) {
	var x dec.Decimal = dec.New(1, -2)
	var y dec.Decimal = dec.New(1, -1)
	if x.LessThan(y) {
		t.Logf(" %s < %s", x.String(), y.String())
	} else {
		t.Errorf(" %s < %s", x.String(), y.String())
	}
}

func TestMatrix1(t *testing.T) {
	a := newAffineTransform()
	if !a.get(1, 1).Equal(ConstZero) {
		t.Error("expected 0 at M(1,1)")
		t.Fail()
	}
	a.set(1, 1, ConstOne)
	if !a.get(1, 1).Equal(ConstOne) {
		t.Error("expected 1 at M(1,1)")
		t.Fail()
	}
}

func TestMatrixRow(t *testing.T) {
	a := Identity()
	row := a.row(0)
	if !row[0].Equal(ConstOne) || !row[1].Equal(ConstZero) {
		t.Errorf("expected (1,0,0), is %v\n", row)
		t.Fail()
	}
}

func TestMatrix2(t *testing.T) {
	a := Identity()
	if !a.get(1, 1).Equal(ConstOne) {
		t.Error("expected 1 at M(1,1)")
		t.Fail()
	}
	p := MakePair(two, three)
	p2 := a.Transform(p)
	if !p2.XPart().Equal(two) || !p2.YPart().Equal(three) {
		t.Errorf("expected %s as result, is %s\n", p, p2)
		t.Fail()
	}
}

func TestTranslation(t *testing.T) {
	a := Translation(MakePair(two, six))
	p := MakePair(ConstZero, ConstZero)
	p = a.Transform(p)
	if !p.XPart().Equal(two) || !p.YPart().Equal(six) {
		t.Errorf("expected (2,6) as result, is %s\n", p)
		t.Fail()
	}
}

func TestRotation(t *testing.T) {
	a := Rotation(dec.New(90, 0).Mul(Deg2Rad))
	p := MakePair(ConstOne, ConstZero)
	p = a.Transform(p)
	if !Round(p.XPart()).Equal(ConstZero) || !Round(p.YPart()).Equal(ConstOne) {
		t.Errorf("expected (0,1) as result, is %s\n", p)
		t.Fail()
	}
}

func TestCombine(t *testing.T) {
	a := Rotation(dec.New(90, 0).Mul(Deg2Rad))
	b := Translation(MakePair(two, six))
	c := a.Combine(b)
	p := MakePair(ConstOne, ConstZero)
	p = c.Transform(p)
	if !Round(p.XPart()).Equal(two) || !Round(p.YPart()).Equal(seven) {
		t.Errorf("expected (2,7) as result, is %s\n", p)
		t.Fail()
	}
}
