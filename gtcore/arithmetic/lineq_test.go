package arithmetic

import (
	"testing"

	"github.com/shopspring/decimal"
)

type X struct { // helper for quick construction of polynomials
	i int
	c float64
}

func polyn(c float64, tms ...X) Polynomial { // construct a polynomial
	p := NewConstantPolynomial(decimal.NewFromFloat(c))
	for _, t := range tms {
		p.SetTerm(t.i, decimal.NewFromFloat(t.c))
	}
	return p
}

/*
func mkpair(a int64, b int64) Pair {
	return Pair{decimal.New(a, 0), decimal.New(b, 0)}
}

func pairpolyn(pc Pair, tms ...X) Polynomial {
	p := NewConstantPairPolynomial(pc)
	for _, t := range tms {
		p.SetTerm(t.i, decimal.NewFromFloat(t.c))
	}
	return p
}
*/

func TestPolynSimple1(t *testing.T) {
	p := NewConstantPolynomial(decimal.Zero)
	if p.Terms.Size() != 1 {
		t.Fail()
	}
}

func TestPolynSimple2(t *testing.T) {
	p := NewConstantPolynomial(decimal.NewFromFloat(0.5))
	p.SetTerm(1, decimal.NewFromFloat(3))
	if p.Terms.Size() != 2 {
		t.Fail()
	}
}

func TestPolynConstant(t *testing.T) {
	p := NewConstantPolynomial(decimal.NewFromFloat(0.5))
	_, isconst := p.IsConstant()
	if !isconst {
		t.Error("did not recognize constant polynomial as constant")
	}
	p.SetTerm(1, decimal.NewFromFloat(2))
	_, isconst = p.IsConstant()
	if isconst {
		t.Error("did falsely recognize non-constant polynomial as constant")
	}
}

func TestZapPolyn(t *testing.T) {
	p := NewConstantPolynomial(decimal.NewFromFloat(0.5))
	p.SetTerm(1, decimal.NewFromFloat(0.0000000005))
	p.Zap()
}

func TestPolynAdd(t *testing.T) {
	p := polyn(5, X{1, 1}, X{2, 2})
	t.Logf("# p  = %s\n", p.String())
	p2 := polyn(4, X{1, 6}, X{5, 4})
	t.Logf("# p2 = %s\n", p2.String())
	pr := p.Add(p2, false)
	t.Logf("# pr = %s\n", pr.String())
	pr.Zap()
	if !pr.GetCoeffForTerm(1).Equal(decimal.New(7, 0)) {
		t.Fail()
	}
}

func TestPolynMul(t *testing.T) {
	p := polyn(6, X{1, 4}, X{2, 2})
	t.Logf("T p  = %s\n", p.String())
	p2 := NewConstantPolynomial(decimal.New(2, 0))
	t.Logf("T p2 = %s\n", p2.String())
	pr := p.Multiply(p2, false)
	t.Logf("T pr = %s\n", pr.String())
	pr.Zap()
	if !pr.GetCoeffForTerm(1).Equal(decimal.New(8, 0)) {
		t.Fail()
	}
}

func TestPolynDiv(t *testing.T) {
	p := polyn(6, X{1, 4}, X{2, 2})
	t.Logf("T p  = %s\n", p.String())
	p2 := NewConstantPolynomial(decimal.New(2, 0))
	t.Logf("T p2 = %s\n", p2.String())
	pr := p.Divide(p2, false)
	t.Logf("T pr = %s\n", pr.String())
	pr.Zap()
	p2 = NewConstantPolynomial(decimal.Zero)
	t.Logf("T p2 = %s\n", p2.String())
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			_, ok = r.(error)
			if ok {
				t.Error("did not detect division by 0\n")
				t.Fail()
			} else {
				t.Logf("test (division by 0) passed: %v\n", r)
			}
		}
	}()
	pr = p.Divide(p2, false)
}

func TestPolynSubst(t *testing.T) {
	p := polyn(1, X{1, 10}, X{2, 20})
	t.Logf("T p  = %s\n", p.String())
	p2 := polyn(2, X{3, 30}, X{4, 40})
	t.Logf("T p2 = %s\n", p2.String())
	p = p.substitute(1, p2)
	t.Logf("T -> p = %s\n", p.String())
	if !p.GetCoeffForTerm(3).Equal(decimal.New(300, 0)) {
		t.Fail()
	}
}

func TestPolynMaxCoeff(t *testing.T) {
	p := polyn(1, X{1, 8}, X{2, 2}, X{3, -2})
	t.Logf("T p  = %s\n", p.String())
	i, c := p.maxCoeff(nil)
	if i != 1 || !c.Equal(decimal.New(8, 0)) {
		t.Fail()
	}
	t.Logf("T ->max coeff @%d is %s, ok\n", i, c.String())
}

func TestNewSolver(t *testing.T) {
	leq := CreateLinEqSolver()
	if leq == nil {
		t.Error("cannot create solver")
	}
}

func TestLinEqAddPolyn(t *testing.T) {
	leq := CreateLinEqSolver()
	p := polyn(1, X{1, 2})
	leq.AddEq(p)
}

func TestLinEqAddPolyn2(t *testing.T) {
	leq := CreateLinEqSolver()
	p := polyn(6, X{1, -1}, X{2, -1})
	leq.AddEq(p)
	q := polyn(2, X{1, 3}, X{2, -1})
	leq.AddEq(q)
}

/*
func TestPolynPairMul(t *testing.T) {
	p1 := pairpolyn(mkpair(1, 1), X{1, 4}, X{2, 2})
	t.Logf("  p1  = %s\n", p1.String())
	p2 := pairpolyn(mkpair(3, 4), X{2, 4}, X{5, 2})
	t.Logf("+ p2  = %s\n", p2.String())
	p := p1.AddPair(p2, false)
	t.Logf("= p   = %s\n", p.String())
	if !decimal.New(4, 0).Equal(p.GetConstantPair().x) { // c(p).x = 4 ?
		t.Fail()
	}
	q := polyn(7)
	p = p.MultiplyPair(q, false)
	t.Logf("* 7   = %s\n", p.String())
	if !decimal.New(28, 0).Equal(p.GetConstantPair().x) { // c(p).x = 28 ?
		t.Fail()
	}
}
*/
