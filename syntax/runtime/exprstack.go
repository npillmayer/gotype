package runtime

import (
	"fmt"
	"math"

	"github.com/emirpasic/gods/stacks/linkedliststack"
	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
	dec "github.com/shopspring/decimal"
)

/*
----------------------------------------------------------------------

BSD License
Copyright (c) 2017, Norbert Pillmayer

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. Neither the name of Norbert Pillmayer or the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----------------------------------------------------------------------

 * This module implements a stack of expressions. It is used for
 * expression evaluation during a parser walk of an expression AST.
 * Expressions can be of type numeric or of type pair.
 *
 * Complexity arises from the fact that we handle not only known
 * quantities, but unknown ones, too. Unknown variables will be handled
 * as terms in linear polynomials. Expressions on the stack are always
 * represented by linear polynomials, containing solved and unsolved variables.
 *
 * The expression stack is connected to a system of linear equations (LEQ).
 * If an equation is constructed from 2 polynomials, it is put into the LEQ.
 * The LEQ operates on generic identifiers and knows nothing of the
 * 'real life' symbols we use in the parser. The expression stack is
 * a bridge between both worlds: It holds a table (VariableResolver) to
 * map LEQ-internal variables to real-life symbols. The variable resolver
 * will receive a message from the LEQ whenever an equation gets solved,
 * i.e. variables become known.

*/

// Some symbols are lvalues, i.e. can be assigned a value
type Assignable interface {
	GetValue() interface{}
	SetValue(interface{})
	IsKnown() bool
}

// === Expressions ===========================================================

/* Interface for expressions: will contain linear polynomials, possibly
 * containing variables with unknown value. Expressions are either of type
 * pair or type numeric. Numeric expressions are modelled as pair values,
 * with the y-part set to 0.
 */
type Expression interface {
	fmt.Stringer
	GetXPolyn() arithm.Polynomial // x-part
	SetXPolyn(p arithm.Polynomial)
	GetYPolyn() arithm.Polynomial // y-part
	SetYPolyn(p arithm.Polynomial)
	IsPair() bool
	IsValid() bool
}

// The type we will push onto the stack
type ExprNode struct {
	xpart  arithm.Polynomial
	ypart  arithm.Polynomial
	ispair bool
}

/* Interface Expression.
 */
func (e *ExprNode) String() string {
	if e.ispair {
		return fmt.Sprintf("(%s,%s)", e.xpart.String(), e.ypart.String())
	} else {
		return e.xpart.String()
	}
}

/* Create a new expression node given a polynomial.
 */
func NewNumericExpression(p arithm.Polynomial) *ExprNode {
	return &ExprNode{xpart: p, ispair: false}
}

/* Create a new pair expression node. Arguments are x-part and y-part
 * for the pair. If no y-part is supplied, the type of the expression will
 * still be type pair - although an invalid one.
 */
func NewPairExpression(xp arithm.Polynomial, yp arithm.Polynomial) *ExprNode {
	return &ExprNode{xpart: xp, ypart: yp, ispair: true}
}

/* Interface Expression.
 */
func (e *ExprNode) GetXPolyn() arithm.Polynomial {
	return e.xpart
}

/* Interface Expression.
 */
func (e *ExprNode) SetXPolyn(p arithm.Polynomial) {
	e.xpart = p
}

/* Interface Expression.
 */
func (e *ExprNode) GetYPolyn() arithm.Polynomial {
	return e.ypart
}

/* Interface Expression.
 */
func (e *ExprNode) SetYPolyn(p arithm.Polynomial) {
	e.ypart = p
}

/* Interface Expression.
 */
func (e *ExprNode) IsPair() bool {
	return e.ispair
}

/* Interface Expression.
 */
func (e *ExprNode) IsValid() bool {
	if e.IsPair() {
		return e.GetXPolyn().Terms != nil && e.GetYPolyn().Terms != nil
	} else {
		return e.GetXPolyn().Terms != nil
	}
}

var _ Expression = &ExprNode{}

// === Expression Stack ======================================================

/* Type ExprStack. This implements a stack of numeric or pair expressions.
 * Various mathematical operations may be performed on the stack values.
 *
 * The expression stack is connected to a system of linear equations (LEQ).
 * If an equation is constructed from 2 polynomials, it is put into the LEQ.
 * The LEQ operates on generic identifiers and knows nothing of the
 * 'real life' symbols we use in the parser. The expression stack is
 * a bridge between both worlds: It holds a table (VariableResolver) to
 * map LEQ-internal variables to real-life symbols. The variable resolver
 * will receive a message from the LEQ whenever an equation gets solved,
 * i.e. variables become known.
 *
 * The connection between symbols and LEQ-variables is by symbol-ID:
 * symbol "a" with ID=7 will be x.7 in LEQ.
 */

type ExprStack struct {
	stack    *linkedliststack.Stack // a stack of expressions
	leq      *arithm.LinEqSolver    // we need a system of linear equations
	resolver map[int]Symbol         // used to resolve variable names from IDs
}

/* Create a new expression stack. It is fully initialized and empty.
 */
func NewExprStack() *ExprStack {
	est := &ExprStack{
		stack:    linkedliststack.New(), // stack of interface{}
		leq:      arithm.CreateLinEqSolver(),
		resolver: make(map[int]Symbol),
	}
	est.leq.SetVariableResolver(est)
	return est
}

/* Give notice of a new variable used in expressions / polynomials.
 * This will put the variable's symbol into the variable resolver's table.
 *
 * Example: symbol "a"|ID=7  =>  resolver table[i] = "a"
 *
 * If a variable ID is not known by the resolver, it is assumed to be
 * a "capsule", which is MetaFont's notation for a variable, which has
 * fallen out of scope.
 */
func (es *ExprStack) AnnounceVariable(v Symbol) {
	T.P("var", v.GetName()).Debugf("announcing id=%d", v.GetID())
	es.resolver[v.GetID()] = v
}

/* Return the name of a variable, given its ID. Will return the string
 * "?nnnn" for capsules.
 *
 * Interface VariableResolver.
 */
func (es *ExprStack) GetVariableName(id int) string {
	v, ok := es.resolver[id]
	if !ok {
		return fmt.Sprintf("?%04d", id)
	}
	return v.GetName()
}

/* Is a variable (index) a capsule, i.e., has it gone out of scope?
 * The terminus stems from MetaFont (with "whatever" being a prominent
 * example for a capsule).
 *
 * Interface VariableResolver.
 */
func (es *ExprStack) IsCapsule(id int) bool {
	_, found := es.resolver[id]
	return !found
}

/* Set the value of a variable. If the LEQ solves a variable and it becomes
 * known, the LEQ will send us this message.
 *
 * Interface VariableResolver.
 */
func (es *ExprStack) SetVariableSolved(id int, val dec.Decimal) {
	v, ok := es.resolver[id]
	if ok { // yes, we know about this variable
		a, isAssignable := v.(Assignable)
		if isAssignable { // is it an lvalue ?
			a.SetValue(val)
		}
	}
}

/* Drop the name of a variable from the variable resolver. The variable itself
 * is not dropped, but rather lives on as an anonymous quantity (i.e., a
 * capsule) as long as it is part of an equation.
 */
func (es *ExprStack) EncapsuleVariable(id int) {
	delete(es.resolver, id)
}

/* Stack functionality. Will return an invalid expression if stack is empty.
 */
func (es *ExprStack) Top() Expression {
	tos, ok := es.stack.Peek()
	if !ok {
		tos = &ExprNode{}
	}
	return tos.(Expression)
}

/* Stack functionality.
 */
func (es *ExprStack) Pop() (Expression, bool) {
	tos, ok := es.stack.Pop()
	return tos.(Expression), ok
}

/* Convenience method: return TOS as a numeric constant.
 */
func (es *ExprStack) PopAsNumeric() (dec.Decimal, bool) {
	tos, ok := es.stack.Pop()
	if ok && tos.(Expression).IsValid() {
		p := tos.(Expression).GetXPolyn()
		if c, isconst := p.IsConstant(); isconst {
			return c, true
		}
	}
	return arithm.ConstZero, false
}

/* Convenience method: return TOS as a pair constant.
 * If TOS is not a known pair, returns (0,0) and false.
 */
func (es *ExprStack) PopAsPair() (arithm.Pair, bool) {
	tos, ok := es.Pop()
	if ok && tos.(Expression).IsValid() && tos.IsPair() {
		xpart, isxconst := tos.(Expression).GetXPolyn().IsConstant()
		ypart, isyconst := tos.(Expression).GetYPolyn().IsConstant()
		if isxconst && isyconst {
			return arithm.MakePair(xpart, ypart), true
		} else {
			T.Debugf("TOS is not constant: %s", tos)
		}
	} else {
		T.Debugf("TOS is not a pair: %s", tos)
	}
	return arithm.MakePair(arithm.ConstZero, arithm.ConstZero), false
}

/* Stack functionality.
 */
func (es *ExprStack) Push(e Expression) *ExprStack {
	es.stack.Push(e)
	return es
}

/* Push a numeric constant onto the stack. It will be wrapped into a
 * polynomial p = c. For pair constants use PushPairConstant(c).
 */
func (es *ExprStack) PushConstant(c dec.Decimal) *ExprStack {
	constant := arithm.NewConstantPolynomial(c)
	T.Debugf("pushing constant = %s", c.String())
	//return es.Push(&ExprNode{constant})
	return es.Push(NewNumericExpression(constant))
}

/* Push a pair constant onto the stack. It will be wrapped into a
 * polynomial p = c. For numeric constants use PushConstant(c).
 */
func (es *ExprStack) PushPairConstant(pc arithm.Pair) *ExprStack {
	xpart := arithm.NewConstantPolynomial(pc.XPart())
	ypart := arithm.NewConstantPolynomial(pc.YPart())
	e := NewPairExpression(xpart, ypart)
	T.Debugf("pushing pair constant = %s", e.String())
	return es.Push(e)
}

/* Push a variable onto the stack. The ID of the variable must be > 0 !
 * It will be wrapped into a polynomial p = 0 + 1 * v.
 * If the variable is of type pair we will push a pair expression.
 */
func (es *ExprStack) PushVariable(v Symbol, w Symbol) *ExprStack {
	es.AnnounceVariable(v)
	p := arithm.NewConstantPolynomial(arithm.ConstZero)
	p = p.SetTerm(v.GetID(), arithm.ConstOne) // p = 0 + 1*v
	if w != nil {
		es.AnnounceVariable(w)
		py := arithm.NewConstantPolynomial(arithm.ConstZero)
		py = py.SetTerm(w.GetID(), arithm.ConstOne) // py = 0 + 1*w
		e := NewPairExpression(p, py)
		symname := fmt.Sprintf("(%s,%s)", v.GetName(), w.GetName())
		T.P("var", symname).Debugf("pushing %s", e.String())
		return es.Push(e)
	} else {
		T.P("var", v.GetName()).Debugf("pushing p = %s", p.String())
		return es.Push(NewNumericExpression(p))
	}
}

func (es *ExprStack) PushPairVariable(xpart Symbol, xconst dec.Decimal, ypart Symbol,
	yconst dec.Decimal) *ExprStack {
	//
	es.AnnounceVariable(xpart)
	es.AnnounceVariable(ypart)
	px := arithm.NewConstantPolynomial(xconst)
	if !xconst.Equal(arithm.ConstZero) {
		px = px.SetTerm(xpart.GetID(), arithm.ConstOne) // px = xconst + 1*xpart
	}
	py := arithm.NewConstantPolynomial(yconst)
	if !yconst.Equal(arithm.ConstZero) {
		py = py.SetTerm(ypart.GetID(), arithm.ConstOne) // py = yconst + 1*ypart
	}
	e := NewPairExpression(px, py)
	symname := fmt.Sprintf("(%s,%s)", xpart.GetName(), ypart.GetName())
	T.P("var", symname).Debugf("pushing %s", e.String())
	return es.Push(e)
}

/* If an expression is a simple variable reference, return the symbol /
 * variable reference. The variable must have been previously announced
 * (see PushVariable(v)).
 */
func (es *ExprStack) GetVariable(e Expression) Symbol {
	if e.IsValid() {
		v, ok := e.GetXPolyn().IsVariable()
		if ok {
			return es.resolver[v]
		}
	}
	return nil
}

/* Stack functionality.
 */
func (es *ExprStack) IsEmpty() bool {
	return es.stack.Empty()
}

/* Stack functionality.
 */
func (es *ExprStack) Size() int {
	return es.stack.Size()
}

/* Internal helper: dump expression stack. This is printed to the trace
 * with level=DEBUG.
 */
func (es *ExprStack) Dump() {
	T.P("size", es.Size()).Debug("Expression Stack, TOS first:")
	it := es.stack.Iterator()
	for it.Next() {
		e := it.Value().(Expression)
		T.P("#", it.Index()).Debugf("    %s", e.GetXPolyn().TraceString(es))
	}
}

/* Print a summary of LEQ and stack contents.
 */
func (es *ExprStack) Summary() {
	es.leq.Dump(es)
	es.Dump()
}

/* Check: is this a valid expression? Will reject un-initialized expressions.
 */
func (es *ExprStack) isValid(e Expression) bool {
	return e.GetXPolyn().Terms != nil
}

/* Pretty print an expression.
 */
func (es *ExprStack) TraceString(e Expression) string {
	if e.IsValid() {
		return e.GetXPolyn().TraceString(es)
	} else {
		return "<empty>"
	}
}

/* Predicate: is an expression of a certain type?
 */
func (es *ExprStack) CheckTypeMatch(e1 Expression, e2 Expression) bool {
	match := false
	if e1.IsValid() && e2.IsValid() {
		if e1.IsPair() == e2.IsPair() {
			match = true
		}
	}
	return match
}

/* Check the operands on the stack for an arithmetic operation.
 * Currently will panic if operands are invalid or not enough operands (n) on
 * stack.
 */
func (es *ExprStack) CheckOperands(n int, op string) {
	if n <= 0 {
		return
	}
	if es.Size() < n {
		panic(fmt.Sprintf("attempt to %s %d operand(s), but %d on stack", op, n, es.Size()))
	}
	if !es.isValid(es.Top()) {
		panic(fmt.Sprintf("TOS operand is invalid for <%s>", op))
	}
}

// Check interface assignabiliy
var _ arithm.VariableResolver = &ExprStack{}

// === Arithmetic Operations =================================================

/* Length of a pair (i.e., distance from origin). Argument must be a known
 * pair.
 */
func (es *ExprStack) LengthTOS() {
	es.CheckOperands(1, "get length of")
	e, _ := es.Pop()
	cx, isconstx := e.GetXPolyn().IsConstant()
	cy, isconsty := e.GetYPolyn().IsConstant()
	if !e.IsPair() || !isconstx || !isconsty {
		T.P("op", "length").Errorf("argument must be known pair")
		panic("not implemented: length(<unknown>)")
	} else {
		T.P("op", "length").Debugf("length of (%s,%s)", cx, cy)
		x, _ := cx.Float64()
		y, _ := cy.Float64()
		l := math.Sqrt(math.Pow(x, 2.0) + math.Pow(y, 2.0))
		es.PushConstant(dec.NewFromFloat(l))
	}
}

/* Add TOS and 2ndOS. Allowed for known and unknown terms.
 */
func (es *ExprStack) AddTOS2OS() {
	var e, e1, e2 Expression
	es.CheckOperands(2, "add")
	e2, _ = es.Pop()
	e1, _ = es.Pop()
	if e1.IsPair() {
		if !e2.IsPair() {
			T.Error("type mismatch: <pair> + <numeric>")
			panic("not implemented: <pair> + <numeric>")
		}
		px := e1.GetXPolyn().Add(e2.GetXPolyn(), false)
		py := e1.GetYPolyn().Add(e2.GetYPolyn(), false)
		e = NewPairExpression(px, py)
	} else {
		if e2.IsPair() {
			T.Error("type mismatch: <numeric> + <pair>")
			panic("not implemented: <numeric> + <pair>")
		}
		px := e1.GetXPolyn().Add(e2.GetXPolyn(), false)
		e = NewNumericExpression(px)
	}
	//es.Push(&ExprNode{p})
	es.Push(e)
	T.P("op", "ADD").Debugf("result %s", e.String())
}

/* Subtract TOS from 2ndOS. Allowed for known and unknown terms.
 */
func (es *ExprStack) SubtractTOS2OS() {
	var e, e1, e2 Expression
	es.CheckOperands(2, "subtract")
	e2, _ = es.Pop()
	e1, _ = es.Pop()
	if e1.IsPair() {
		if !e2.IsPair() {
			T.Error("type mismatch: <pair> - <numeric>")
			panic("not implemented: <pair> - <numeric>")
		}
		px := e1.GetXPolyn().Subtract(e2.GetXPolyn(), false)
		py := e1.GetYPolyn().Subtract(e2.GetYPolyn(), false)
		e = NewPairExpression(px, py)
	} else {
		if e2.IsPair() {
			T.Error("type mismatch: <numeric> - <pair>")
			panic("not implemented: <numeric> - <pair>")
		}
		px := e1.GetXPolyn().Subtract(e2.GetXPolyn(), false)
		e = NewNumericExpression(px)
	}
	//es.Push(&ExprNode{p})
	es.Push(e)
	T.P("op", "SUB").Debugf("result %s", e.String())
}

/* Multiply TOS and 2ndOS. One multiplicant must be a known numeric constant.
 */
func (es *ExprStack) MultiplyTOS2OS() {
	var e, e1, e2 Expression
	es.CheckOperands(2, "multiply")
	e2, _ = es.Pop()
	e1, _ = es.Pop()
	if e2.IsPair() {
		e1, e2 = e2, e1
	}
	if e1.IsPair() {
		if e2.IsPair() {
			T.Errorf("one multiplicant must be a known numeric")
			panic("not implemented: <pair> * <pair>")
		} else {
			n := e2.GetXPolyn()
			px := e1.GetXPolyn().Multiply(n, false)
			py := e1.GetYPolyn().Multiply(n, false)
			e = NewPairExpression(px, py)
		}
	} else {
		px := e1.GetXPolyn().Multiply(e2.GetXPolyn(), false)
		e = NewNumericExpression(px)
	}
	es.Push(e)
	T.P("op", "MUL").Debugf("result = %s", e.String())
}

/* Divide 2ndOS by TOS. Divisor must be numeric non-0 constant.
 */
func (es *ExprStack) DivideTOS2OS() {
	var e, e1, e2 Expression
	es.CheckOperands(2, "divide")
	e2, _ = es.Pop()
	e1, _ = es.Pop()
	if e2.IsPair() {
		T.Errorf("divisor must be a known non-zero numeric")
		panic("not implemented: division by <pair>")
	}
	if e1.IsPair() {
		n := e2.GetXPolyn()
		nn := n.CopyPolynomial()
		px := e1.GetXPolyn().Divide(n, false) // n will be destroyed
		py := e1.GetYPolyn().Divide(nn, false)
		e = NewPairExpression(px, py)
	} else {
		px := e1.GetXPolyn().Divide(e2.GetXPolyn(), false)
		e = NewNumericExpression(px)
	}
	es.Push(e)
	T.P("op", "DIV").Debugf("result = %s", e.String())
}

/* Numeric interpolation operation. Either n must be known or a and b.
 *
 * n[a,b] => a - na + nb.
 */
func (es *ExprStack) Interpolate() {
	es.CheckOperands(3, "interpolate")
	var n, a, b Expression
	b, _ = es.Pop()
	a, _ = es.Pop()
	n, _ = es.Pop()
	if a.IsPair() {
		es.InterpolatePair(n, a, b)
	} else {
		// second operand will be destroyed, n must be first !
		p1 := n.GetXPolyn().Multiply(a.GetXPolyn(), false)
		p2 := n.GetXPolyn().Multiply(b.GetXPolyn(), false)
		p := a.GetXPolyn().Subtract(p1, false)
		p = p.Add(p2, false)
		e := NewNumericExpression(p)
		es.Push(e)
		T.P("op", "INTERP").Debugf("result = %s", p.String())
	}
}

/* Pair interpolation operation. Either n must be known or z1 and z2.
 *
 * n[z1,z2] => z1 - n*z1 + n*z2.
 */
func (es *ExprStack) InterpolatePair(n Expression, z1 Expression, z2 Expression) {
	// second operand will be destroyed, n must be first !
	px1 := n.GetXPolyn().Multiply(z1.GetXPolyn(), false)
	px2 := n.GetXPolyn().Multiply(z2.GetXPolyn(), false)
	px := z1.GetXPolyn().Subtract(px1, false)
	px = px.Add(px2, false)
	py1 := n.GetXPolyn().Multiply(z1.GetYPolyn(), false)
	py2 := n.GetXPolyn().Multiply(z2.GetYPolyn(), false)
	py := z1.GetYPolyn().Subtract(py1, false)
	py = py.Add(py2, false)
	e := NewPairExpression(px, py)
	es.Push(e)
	T.P("op", "INTERP").Debugf("result = %s", e.String())
}

/* Rotate a pair around origin for TOS degrees, counterclockwise.
 * TOS must be a known numeric constant.
 */
func (es *ExprStack) Rotate2OSbyTOS() {
	es.CheckOperands(2, "rotate")
	e, _ := es.Pop()
	c, _ := e.GetXPolyn().IsConstant()
	angle, _ := c.Mul(arithm.Deg2Rad).Float64()
	T.Debugf("rotating by %s° = %f rad", c, angle)
	sin := arithm.NewConstantPolynomial(dec.NewFromFloat(math.Sin(angle)))
	cos := arithm.NewConstantPolynomial(dec.NewFromFloat(math.Cos(angle)))
	//T.Debugf("sin %s° = %s, cos %s° = %s", c, sin, c, cos)
	e, _ = es.Pop()
	if e.IsPair() {
		var ysin, ycos, xpart, ypart, tmp arithm.Polynomial
		tmp = sin.CopyPolynomial()
		ysin = e.GetYPolyn().Multiply(tmp, false)
		tmp = cos.CopyPolynomial()
		ycos = e.GetYPolyn().Multiply(tmp, false)
		xpart = e.GetXPolyn().Multiply(cos, false).Subtract(ysin, false)
		ypart = e.GetXPolyn().Multiply(sin, false).Subtract(ycos, false)
		e = NewPairExpression(xpart, ypart)
		es.Push(e)
	} else {
		T.P("op", "rotate").Errorf("not implemented: rotate <non-pair>")
	}
}

/* Create an equation of the polynomials of TOS and 2ndOS.
 * Introduces the equation to the solver's linear equation system.
 *
 * If the polynomials are of type pair polynomial, then there will be 2
 * equations, one for the x-part and one for the y-part. LEQ will only handle
 * numeric linear equations.
 */
func (es *ExprStack) EquateTOS2OS() {
	es.SubtractTOS2OS() // now 0 = p1 - p2
	e, _ := es.Pop()    // e is interpreted as an equation, one side 0
	if e.IsPair() {
		var eqs = []arithm.Polynomial{
			e.GetXPolyn(),
			e.GetYPolyn(),
		}
		es.leq.AddEqs(eqs)
	} else {
		es.leq.AddEq(e.GetXPolyn())
	}
}
