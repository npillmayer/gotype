package arithmetic

import (
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestInit0(t *testing.T) {
	tracing.EquationsTracer = gologadapter.New()
	T = tracing.EquationsTracer
}

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

type res map[int]decimal.Decimal // a variable resolver for testing purposes

func NewResolver() res {
	var r res
	r = make(map[int]decimal.Decimal)
	return r
}

func (r res) GetVariableName(n int) string { // get real-life name of x.i
	return string(rune(n + 96)) // 'a', 'b', ...
}

func (r res) SetVariableSolved(n int, v decimal.Decimal) { // message: x.i is solved
	//T.P("msg", "SOLVED").Infof("%s = %s", r.GetVariableName(n), v.String())
	r[n] = v // remember the value to assert test conditions
}

func (r res) IsCapsule(int) bool { // x.i has gone out of scope
	return false // no capsules
}

// --- Tests -----------------------------------------------------------------

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

func TestSubst1(t *testing.T) {
	p := polyn(1, X{3, 3})
	q := polyn(2, X{1, 3}, X{4, 4}, X{5, 5})
	t.Logf("   x.%d = %s\n", 1, p.String())
	t.Logf("   x.%d = %s\n", 2, q.String())
	var k int
	k, q = subst(1, p, 2, q)
	t.Logf("=> x.%d = %s\n", k, q.String())
	if !q.GetCoeffForTerm(3).Equal(decimal.New(9, 0)) {
		t.Fail()
	}
}

func TestSubst2(t *testing.T) {
	p := polyn(1, X{3, 3})
	q := polyn(2, X{1, 3}, X{3, 4}, X{5, 5})
	t.Logf("   x.%d = %s\n", 1, p.String())
	t.Logf("   x.%d = %s\n", 2, q.String())
	var k int
	k, q = subst(1, p, 2, q)
	t.Logf("=> x.%d = %s\n", k, q.String())
	if !q.GetCoeffForTerm(3).Equal(decimal.New(13, 0)) {
		t.Fail()
	}
}

func TestSubst3(t *testing.T) {
	p := polyn(1, X{2, 1})
	q := polyn(2, X{1, 1}, X{4, 4}, X{5, 5})
	t.Logf("   x.%d = %s\n", 1, p.String())
	t.Logf("   x.%d = %s\n", 2, q.String())
	var k int
	k, q = subst(1, p, 2, q)
	t.Logf("=> x.%d = %s\n", k, q.String())
	if k != 0 {
		t.Fail()
	}
}

func TestNewSolver(t *testing.T) {
	leq := CreateLinEqSolver()
	if leq == nil {
		t.Error("cannot create solver")
		t.Fail()
	}
}

func TestLinEqAddPolyn(t *testing.T) {
	leq := CreateLinEqSolver()
	r := NewResolver()
	leq.SetVariableResolver(r)
	p := polyn(1, X{1, 2})
	leq.AddEq(p)
	if _, found := r[1]; !found {
		t.Error("a still unsolved")
		t.Fail()
	}
}

func TestLinEqAddPolyn2(t *testing.T) {
	leq := CreateLinEqSolver()
	r := NewResolver()
	leq.SetVariableResolver(r)
	p := polyn(6, X{1, -1}, X{2, -1})
	leq.AddEq(p)
	q := polyn(2, X{1, 3}, X{2, -1})
	leq.AddEq(q)
	if _, found := r[1]; !found {
		t.Error("a still unsolved")
		t.Fail()
	}
}

func TestLEQ1(t *testing.T) {
	leq := CreateLinEqSolver()
	r := NewResolver()
	leq.SetVariableResolver(r)
	p1 := polyn(100, X{1, -2})           // 2a=100   =>  0=100-2a
	p2 := polyn(100, X{2, -1}, X{3, -1}) // 100=b+c  =>  0=100-b-c
	leq.AddEq(p1)
	leq.AddEq(p2)
	if _, found := r[1]; !found {
		t.Error("a still unsolved")
		t.Fail()
	}
}

func TestLEQ2(t *testing.T) {
	leq := CreateLinEqSolver()
	r := NewResolver()
	leq.SetVariableResolver(r)
	p1 := polyn(100, X{2, -1}, X{3, -1})        // b+c=100 =>  0=100-b-c
	p2 := polyn(0, X{1, 2}, X{2, -1}, X{3, -1}) // 2a=b+c  =>  0=2a-b-c
	leq.AddEq(p1)
	leq.AddEq(p2)
	if _, found := r[1]; !found {
		t.Error("a still unsolved")
		t.Fail()
	}
}

func TestLEQ3(t *testing.T) {
	leq := CreateLinEqSolver()
	r := NewResolver()
	leq.SetVariableResolver(r)
	p1 := polyn(100, X{1, -1}) // a = 100
	p2 := polyn(99, X{1, -2})  // 2a = 99
	leq.AddEq(p1)
	//leq.AddEq(p2)
	assert.Panics(t, func() { leq.AddEq(p2) }, "equation should be off by -101")
}

func TestLEQ4(t *testing.T) {
	//T.SetLevel(logrus.DebugLevel)
	leq := CreateLinEqSolver()
	r := NewResolver()
	leq.SetVariableResolver(r)
	p1 := polyn(100, X{1, -1})                          // a=100
	p2 := polyn(0, X{1, 2}, X{2, -1}, X{3, 1}, X{4, 4}) // 2a=b-c-e
	p3 := polyn(0, X{2, 1}, X{3, -1})                   // b=c
	leq.AddEq(p1)
	leq.AddEq(p2)
	leq.AddEq(p3) // eliminates b and c from p2 => d solved
	if _, found := r[4]; !found {
		t.Error("d still unsolved")
		t.Fail()
	}
}

func TestLEQ5(t *testing.T) {
	//T.SetLevel(logrus.DebugLevel)
	leq := CreateLinEqSolver()
	//leq.showdependencies = true
	r := NewResolver()
	leq.SetVariableResolver(r)
	p1 := polyn(0, X{2, -1}, X{3, 1}) // b=c
	p2 := polyn(0, X{3, -1}, X{4, 1}) // c=d
	p3 := polyn(0, X{4, -1}, X{2, 1}) // d=b
	leq.AddEq(p1)
	leq.AddEq(p2)
	leq.AddEq(p3)
	p4 := polyn(0, X{1, -1}, X{2, 1}, X{3, 1}, X{4, 1}) // a=b+c+d
	leq.AddEq(p4)                                       // now a=3d (or b or c)
	a, _ := leq.dependents.Get(1)
	p := a.(Polynomial)
	if termlength(p) != 2 { // a = 0 + 3d
		t.Fail()
	}
}

// Example for solving linear equations. We use a variable resolver, which
// maps a numeric value of 0..<n> to lowercase letters 'a'..'z'.
func TestExampleLinEqSolver_usage(t *testing.T) {
	//func TestExampleLinEqSolver_usage() {
	leq := CreateLinEqSolver()
	r := NewResolver() // clients have to provide their own
	leq.SetVariableResolver(r)
	p := polyn(6, X{1, -1}, X{2, -1})
	leq.AddEq(p)
	q := polyn(2, X{1, 3}, X{2, -1})
	leq.AddEq(q)
}
