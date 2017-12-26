package syntax

import (
	"log"
	"testing"

	"github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/shopspring/decimal"
)

type X struct { // helper for quick construction of polynomials
	i int
	c float64
}

func polyn(c float64, tms ...X) arithmetic.Polynomial { // construct a polynomial
	p := arithmetic.NewConstantPolynomial(decimal.NewFromFloat(c))
	for _, t := range tms {
		p.SetTerm(t.i, decimal.NewFromFloat(t.c))
	}
	return p
}

func TestStackCreate(t *testing.T) {
	p := polyn(5, X{1, 1}, X{2, 2})
	log.Printf("p = %s\n", p.String())
}

func TestStackVar1(t *testing.T) {
	est := NewExprStack()
	est.AnnounceVariable(NewStdSymbol("sym1"))
	p := polyn(5, X{1, 1}, X{2, 2})
	log.Printf("p = %s\n", p.TraceString(est))
}

func TestStackVar2(t *testing.T) {
	est := NewExprStack()
	est.PushConstant(decimal.New(4711, 0))
	log.Printf("TOS = %s\n", est.Top().GetPolynomial().TraceString(est))
}

func TestStackVar3(t *testing.T) {
	est := NewExprStack()
	v := NewStdSymbol("a")
	est.PushVariable(v)
	log.Printf("TOS = %s\n", est.Top().GetPolynomial().TraceString(est))
}

func TestStackAdd(t *testing.T) {
	est := NewExprStack()
	v1 := NewStdSymbol("a")
	est.PushVariable(v1)
	v2 := NewStdSymbol("b")
	est.PushVariable(v2)
	est.AddTOS2OS()
	log.Printf("TOS = %s\n", est.Top().GetPolynomial().TraceString(est))
	if !est.Top().GetPolynomial().GetConstantValue().Equal(decimal.Zero) {
		t.Fail()
	}
	if est.Top().GetPolynomial().Terms.Size() != 3 {
		t.Fail()
	}
}

func TestStackSubtract(t *testing.T) {
	est := NewExprStack()
	v1 := NewStdSymbol("a")
	est.PushVariable(v1)
	a := v1.GetID()
	p := polyn(2.0, X{a, 3.0})
	est.Push(&ExprSimplifier{p}) // push p = 3a + 2
	est.SubtractTOS2OS()         // should result in p = -2a - 2
	log.Printf("TOS = %s\n", est.Top().GetPolynomial().TraceString(est))
	if !est.Top().GetPolynomial().GetConstantValue().Equal(decimal.New(-2, 0)) {
		t.Fail()
	}
	if est.Top().GetPolynomial().Terms.Size() != 2 {
		t.Fail()
	}
}

func TestStackMultiply(t *testing.T) {
	est := NewExprStack()
	v1 := NewStdSymbol("a")
	est.PushVariable(v1)
	p := polyn(2.0)              // constant
	est.Push(&ExprSimplifier{p}) // push p = 2
	est.MultiplyTOS2OS()
	log.Printf("TOS = %s\n", est.Top().GetPolynomial().TraceString(est))
	if !est.Top().GetPolynomial().GetConstantValue().Equal(decimal.Zero) {
		t.Fail()
	}
	if est.Top().GetPolynomial().Terms.Size() != 2 {
		t.Fail()
	}
}

func TestStackDivide(t *testing.T) {
	est := NewExprStack()
	v1 := NewStdSymbol("a")
	est.PushVariable(v1)
	p := polyn(2.0)              // constant
	est.Push(&ExprSimplifier{p}) // push p = 2
	est.DivideTOS2OS()
	log.Printf("TOS = %s\n", est.Top().GetPolynomial().TraceString(est))
	if !est.Top().GetPolynomial().GetConstantValue().Equal(decimal.Zero) {
		t.Fail()
	}
	if est.Top().GetPolynomial().Terms.Size() != 2 {
		t.Fail()
	}
}

func TestStackEquation(t *testing.T) {
	est := NewExprStack()
	v1 := NewStdSymbol("a")
	est.PushVariable(v1)
	v2 := NewStdSymbol("b")
	est.PushVariable(v2)
	est.Dump()
	est.EquateTOS2OS()
	if !est.IsEmpty() {
		t.Fail()
	}
	est.leq.Dump(est)
}
