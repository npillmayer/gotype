package arithmetic

import (
	"bytes"
	"fmt"

	"github.com/emirpasic/gods/maps"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treemap"
	numeric "github.com/shopspring/decimal"
)

/*
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
3. Neither the name of Tom Everett nor the names of its contributors
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

 * Objects and interfaces for solving systems of linear equations.
 *
 * Inspired by Donald E. Knuth's MetaFont, John Hobby's MetaPost and by
 * a Lua project by John D. Ramsdell: http://luaforge.net/projects/lineqpp/

*/

/* We use an interface to resolve "real" variable names. Within the LEQ
 * variables are encoded by their serial ID, which is used as their position
 * within polynomias. Example: variable "n[3].a" with ID=4711 will become x.4711
 * internally. The resolver maps x.4711 => "n[3].a", i.e., IDs to names.
 */
type VariableResolver interface {
	GetVariableName(int) string             // get real-life name of x.i
	SetVariableSolved(int, numeric.Decimal) // message: x.i is solved
	IsCapsule(int) bool                     // x.i has gone out of scope
}

// === System of linear equations =======================================

type LinEqSolver struct {
	variables   *hashmap.Map     // stores pointers to Symbols // TODO remove
	dependents  *treemap.Map     // dependent variable at position i has dependencies[i]
	solved      *treemap.Map     // map x.i => numeric
	varresolver VariableResolver // to resolve variable names from term positions
}

/* Create a new sytem of linear equations.
 */
func CreateLinEqSolver() *LinEqSolver {
	leq := LinEqSolver{
		dependents: treemap.NewWithIntComparator(),
		solved:     treemap.NewWithIntComparator(),
	}
	return &leq
}

/* Set a variable resolver. Within the LEQ variables are
 * encoded by their serial ID, which is used as their position within
 * polynomias. Example: variable "n[3].a" with ID=4711 will become x.4711
 * internally. The resolver maps x.4711 => "n[3].a".
 */
func (leq *LinEqSolver) SetVariableResolver(resolver VariableResolver) {
	leq.varresolver = resolver
}

/* Collect all currently solved variables from a system of linear equations.
 * Solved variables are returned as a map: i(var) -> numeric, where i(var) is an
 * integer representing the position of variable var.
 */
func (leq *LinEqSolver) getSolvedVars() maps.Map {
	setOfSolved := treemap.NewWithIntComparator() // return value
	it := leq.solved.Iterator()
	for it.Next() { // for every x.i = p[x.i = c] => put [x.i = c] into new set
		setOfSolved.Put(it.Key().(int), it.Value().(Polynomial).GetCoeffForTerm(0))
	}
	return setOfSolved
}

/* Add a new equation 0 = p (p is Polynomial) to a system of linear equations.
 * Immediately starts to solve the -- possibly incomplete -- system, as
 * far as possible.
 */
func (leq *LinEqSolver) AddEq(p Polynomial) *LinEqSolver {
	leq.addEq(p, false)
	leq.Dump(leq.varresolver)
	return leq
}

/* Add a list of linear equations to the LEQ.
 */
func (leq *LinEqSolver) AddEqs(plist []Polynomial) *LinEqSolver {
	l := len(plist)
	if l == 0 {
		T.Error("given empty list of equations")
	} else {
		for i, p := range plist {
			T.Debugf("adding equation %d/%d: 0 = %s", i+1, l, p)
			leq.addEq(p, i+1 < l)
		}
	}
	leq.Dump(leq.varresolver)
	return leq
}

/* If parameter cont is true, expect another equation immediately after this
 * one. This is necessary to suppress harvesting of capsules.
 */
func (leq *LinEqSolver) addEq(p Polynomial, cont bool) *LinEqSolver {
	p = p.Zap()
	T.Infof("# 0 = %s", leq.PolynString(p))
	p = leq.substituteSolved(p)
	if _, off := p.isOff(); !off { //  :-))  no pun intended
		i, _ := p.maxCoeff(leq.dependents) // start with max (free) coefficient of p
		if i == 0 {
			panic("I think this is an impossible error: seeing equation 0 = c")
		}
		p = leq.activateEquationTowards(i, p) // now  x.i = -1/a * p(...).
		// now loop: substitute variables in dependent equations until no more changes
		j, q, setOfSolved := i, p, treemap.NewWithIntComparator() // set up iteration
		for true {                                                // repeat until set of solved vars is empty
			// substitute p for a.j*x.j in all p.n
			setOfSolved = leq.substituteDependencies(j, q, setOfSolved)
			if !setOfSolved.Empty() { // this may result in some variables becoming known : x.i = { c }
				j = setOfSolved.Keys()[0].(int) // get first solved var and re-iterate
				p1, _ := setOfSolved.Get(j)     // we know this polyn is of constant type { c }
				q = p1.(Polynomial)             // set has returned interface{} -> cast
				leq.solved.Put(j, q)            // don't lose solution...
				setOfSolved.Remove(j)           // ...but remove from working set of solved vars
			} else {
				break // no more changes in set of dependent variables => stop
			}
		}
		leq.checkIfSolved(i, p, leq.solved, leq.dependents) // p may solve x.i
	}
	if !cont {
		leq.harvestCapsules()
	}
	return leq
}

/* In a (new) equation, substitute all variables which are already known.
 */
func (leq *LinEqSolver) substituteSolved(p Polynomial) Polynomial {
	it := leq.solved.Iterator()
	T.Debug("---------- subst solved -----------")
	for it.Next() { // iterate over all solved x.i = c
		i := it.Key().(int)
		c := it.Value().(Polynomial).GetConstantValue()
		coeff := p.GetCoeffForTerm(i)
		if !Zero(coeff) {
			coeff = coeff.Mul(c)
			pc := p.GetConstantValue()
			p.SetTerm(0, pc.Add(coeff))
			p.Terms.Remove(i)
			T.Debugf("x.%d = %s  => %s", i, c.String(), p.String())
		}
	}
	T.Debug("-----------------------------------")
	return p
}

/* Transform an equation 0 = p(a x.i) to make x.i the dependent variable, i.e.
 * x.i = -1/a * p(...).
 */
func (leq *LinEqSolver) activateEquationTowards(i int, p Polynomial) Polynomial {
	coeff := p.GetCoeffForTerm(i)
	p.Terms.Remove(i)                                          // remove term x.i from RHS(p)
	pp := NewConstantPolynomial(numeric.New(-1, 0).Div(coeff)) // -1/coeff
	p = p.Multiply(pp, true).Zap()
	T.Infof("## x.%d = %s", i, p.String())
	return p
}

/* Given an equation x.v = p, substitute p for x.w in Polynomials of every
 * dependent variable. This may solve dependent equations, thus the
 * method returns a set of solved variables (and their solutions).
 * The set of dependent variables of the LEQ may shrink.
 *
 * A special case is, with x.v = p1 given, that there may be another dependent
 * equations x.v = p2. Then we select a free x.j in p2 and make it the dependent
 * variable, resulting in x.j = p2' - x.v, and we proceed substituting x.v with
 * p1, resulting in x.j = p2' - p1.
 */
func (leq *LinEqSolver) substituteDependencies(v int, p Polynomial, setOfSolved *treemap.Map) *treemap.Map {
	//setOfSolved := treemap.NewWithIntComparator() // set of solved variables = return value
	setOfDeps := treemap.NewWithIntComparator() // temp. set of dependents variables
	it := leq.dependents.Iterator()
	T.Debug("------------ subst dep ------------")
	for it.Next() { // iterate over all dependent x.w = p.w ( c ... { a x.v } ... )
		w := it.Key().(int)
		pw := it.Value().(Polynomial)
		T.Debugf("# (1) p(x.%d) in x.%d = %s", v, w, pw.String())
		if w == v { // x.w = x.v, i.e. equations with identical LHS
			j, _ := pw.maxCoeff(leq.dependents) // start with max (free) coefficient of p.w
			if j == 0 {
				panic("I think this is an impossible error: seeing equation 0 = c")
			}
			pp := NewConstantPolynomial(ConstZero)  // construct LHS as pp
			pp.SetTerm(w, ConstOne)                 // now LHS is { 0 + 1 x.w }
			pw = pw.Subtract(pp, false)             // move to RHS
			pw = leq.activateEquationTowards(j, pw) // now  x.j = -1/a * p(... x.v ...).
			w = j                                   // ride the new horse
		}
		T.Debugf("# (2) p(x.%d) = %s  in  x.%d = %s", v, p, w, pw.String())
		pw = pw.substitute(v, p).Zap()
		T.Debugf("# => RHS = %s", pw.String())
		aw := pw.GetCoeffForTerm(w) // check for case x.w = a x.w + [...] + c
		if !Zero(aw) {              // found x.w in p.w => solve for x.w
			aw = aw.Sub(ConstOne)
			aw = numeric.Zero.Sub(ConstOne.Div(aw)) // sc = -1/(a-1)
			T.Debugf("# => -1/(a-1) = %s", aw.String())
			t := NewConstantPolynomial(aw)
			pw.Terms.Remove(w)               // now solve for x.w => remove a*x.w from RHS p
			pw = pw.Multiply(t, false).Zap() // and Multiply RHS by -1/(a-1)
			T.Debugf("# => x.%d = %s", w, pw.String())
			T.Infof("## x.%d = %s", w, pw.String())
		}
		leq.checkIfSolved(w, pw, setOfSolved, setOfDeps) // pw may now solve x.w
	}
	leq.dependents = setOfDeps
	T.Debug("-----------------------------------")
	return setOfSolved
}

/* For x.i = p, check if p is constant and solves x.i, then put x.i either
 * into the set of solved variables or the set of dependent variables.
 */
func (leq *LinEqSolver) checkIfSolved(i int, p Polynomial, setOfSolved maps.Map, setOfDependents maps.Map) {
	if rhs, isconst := p.IsConstant(); isconst { // RHS is const and solves x.i
		if _, found := setOfSolved.Get(i); !found { // if not already in set
			rhs = Round(rhs)            // round to epsilon
			p.SetTerm(0, rhs)           // replace const coeff by rounded value
			if leq.varresolver == nil { // now print it out
				T.Infof("#### x.%d = %s", i, rhs.String()) // print internal solution
			} else {
				varname := leq.varresolver.GetVariableName(i)
				T.Infof("#### x.%d = %s", i, rhs.String()) // print internal solution
				T.P("var", varname).Infof("#### %s = %s", varname, rhs.String())
			}
			setOfSolved.Put(i, p) // move x.i to set of solved variables
			if leq.varresolver != nil {
				leq.varresolver.SetVariableSolved(i, rhs) // notify variable solver
			}
		}
	} else {
		setOfDependents.Put(i, p) // move x.i to set of dependent variables
	}
}

/* Helper: variable as string. Uses VariableResolver, if present.
 */
func (leq *LinEqSolver) VarString(i int) string {
	if leq.varresolver == nil {
		return fmt.Sprintf("x.%d", i)
	} else {
		return leq.varresolver.GetVariableName(i)
	}
}

/* Helper: Polynomial as string. Uses VariableResolver, if present.
 */
func (leq *LinEqSolver) PolynString(p Polynomial) string {
	if leq.varresolver != nil {
		return p.TraceString(leq.varresolver)
	} else {
		return p.String()
	}
}

// === Capsules ==============================================================

/* 'Capsule' is a MetaFont terminus for variables in the LEQ, which have
 * fallen out of scope. This may happen on "endgroup", if a variable has
 * been "save"d, or with assignments, where the old incarnation of an
 * lvalue may still be entangled in the LEQ. Capsules may still be relevant
 * in the LEQ, but are of no further interest to the user.
 *
 * The typical case in MetaFont ist the use of "whatever", e.g. in the
 * equation z0 = whatever[z1,z2] (z0 is somewhere on the straight line
 * trough z1 and z2). "whatever" is defined as "begingroup save ?; ? endgroup".
 * The variable ? falls out of scope, but is still relevant for solving the
 * equations for z0 (the above command produces 2 equations).
 */

/* Remove all equations which are dependent on a capsule, but only if the
 * capsule is a loner. If a capsule occurs in at least 2 equations, it
 * is still relevant for solving the LEQ.
 */
func (leq *LinEqSolver) harvestCapsules() {
	var counts map[int]int = make(map[int]int)
	it := leq.dependents.Iterator()
	for it.Next() { // iterate over all dependent x.w = p.w ( c ... { a x.v } ... )
		w := it.Key().(int)
		pw := it.Value().(Polynomial)
		leq.checkAndCountCapsule(w, counts) // check LHS variable
		pit := pw.Terms.Iterator()          // for all terms in polynomial
		for pit.Next() {
			i := pit.Key().(int) // get every term.i
			if i > 0 {           // omit constant term
				leq.checkAndCountCapsule(i, counts)
			}
		}
	}
	itsolv := leq.solved.Iterator() // count solved capsules
	for itsolv.Next() {
		j := itsolv.Key().(int)
		leq.checkAndCountCapsule(j, counts)
	}
	for pos, count := range counts { // now remove capsules with count == 1
		if count == 1 { // only remove loners
			T.P("capsule", pos).Debug("capsule removed")
			leq.retractVariable(pos)
		}
	}
}

/* Helper for counting capsule references. Updates the count for a capsule.
 */
func (leq *LinEqSolver) checkAndCountCapsule(i int, counts map[int]int) {
	if leq.varresolver != nil && leq.varresolver.IsCapsule(i) {
		counts[i] += 1
		//T.P("capsule", i).Debugf("capsule counted, #=%d", counts[i])
	}
}

/* If a capsule is removed, all equations containing the capsule must
 * be deleted from the LEQ.
 *
 * TODO: The whole procedure for removing capsules is awfully inefficient:
 * lots of set iterations (some nested loops) and set creations. But for
 * my use cases the number of simultaneous equations is small, therefore
 * I'll clean this up sometime later... :-)
 */
func (leq *LinEqSolver) retractVariable(i int) {
	if _, ok := leq.solved.Get(i); ok {
		T.Debugf("unsolve %s", leq.VarString(i))
		leq.solved.Remove(i)
	}
	leq.dependents.Remove(i)              // possibly remove from dependents
	eqs := treemap.NewWithIntComparator() // set of equation indices, i.e. int
	it := leq.dependents.Iterator()
	for it.Next() { // iterate over all dependent x.j = p.i ( c ... { a x.i } ... )
		j := it.Key().(int)
		p := it.Value().(Polynomial)
		if a := p.GetCoeffForTerm(i); !numeric.Zero.Equal(a) { // yes, x.i in p
			eqs.Put(j, p) // mark for deletion, as it is invalid now
		}
	}
	it = eqs.Iterator()
	for it.Next() { // iterate over marked equations
		leq.dependents.Remove(it.Key().(int))
	}
}

// === (Linear) Polynomials =============================================

/* Type for linear Polynomials  c + a.1 x.1 + a.2 x.2 + ... a.n x.n .
 * We store the coefficients only. Index 0 is the constant term.
 * We store the scales/coeff in a TreeMap (sorted map). Coefficients are of
 * type numeric.Decimal.
 */
type Polynomial struct {
	Terms *treemap.Map
}

/* Create a Polynomial consisting of just a constant term.
 */
func NewConstantPolynomial(c numeric.Decimal) Polynomial {
	m := treemap.NewWithIntComparator()
	//p := Polynomial{m, false}
	p := Polynomial{m}
	p.Terms.Put(0, c) // initialize with constant term (at position 0)
	return p.Zap()
}

/* Set the coefficient for a term a.i within a Polynomial.
 * for i=0: constant term. If this is a Pair Polynomial and i = 0, then
 * the constant term will be set to (scale,scale).
 */
func (p Polynomial) SetTerm(i int, scale numeric.Decimal) Polynomial {
	p.Terms.Put(i, scale)
	return p
}

/* Helper: for an equation [ 0 = p ] check if p is constant and != 0.
 */
func (p Polynomial) isOff() (numeric.Decimal, bool) {
	if coeff, isconst := p.IsConstant(); isconst {
		//coeff := p.getCoeffForTerm(0)
		if !Zero(coeff) {
			panic(fmt.Sprintf("equation off by %s", coeff.String()))
		}
		return coeff, true
	}
	return numeric.Zero, false
}

/* Find coefficient of maximum absolute value.
 * If parameter 'dependents' is given, first search for a.i * x.i, with
 * x.i not in dependents (i.e., we're looking for free variables only:
 * find free variable x.i in p, with abs(a.i) is max in p).
 * If no free variable can be found, find max(dependent(a.j)).
 *
 * TODO: will using a non-free var overwrite existing equations?
 */
func (p Polynomial) maxCoeff(dependents maps.Map) (int, numeric.Decimal) {
	it := p.Terms.Iterator()
	var maxp int                            // variable position of max coeff
	var maxc numeric.Decimal = numeric.Zero // max coeff
	var coeff numeric.Decimal               // result coeff
	for it.Next() {
		i := it.Key().(int)
		var isdep = false
		if dependents != nil {
			_, isdep = dependents.Get(i) // could be better de-coupled by providing predicate func
		}
		if i == 0 || isdep {
			continue
		}
		c := p.GetCoeffForTerm(i)
		if c.Abs().GreaterThan(maxc) {
			maxc, maxp, coeff = c.Abs(), i, c
		}
	}
	if maxp == 0 && dependents != nil { // no free variable found
		maxp, coeff = p.maxCoeff(nil) //
	}
	//log.Printf("max-coeff = %s @ %d\n", coeff, maxp)
	return maxp, coeff
}

/* Substitute variable i within p with Polynomial p2.
 * If p does not contain a term.i, p is unchanged
 * This routine is detructive!
 */
func (p Polynomial) substitute(i int, p2 Polynomial) Polynomial {
	scale_i := p2.GetCoeffForTerm(i)
	if !Zero(scale_i) {
		panic(fmt.Sprintf("cyclic call to substitute term #%d: %s", i, p2.String()))
	}
	scale_i = p.GetCoeffForTerm(i)
	if !Zero(scale_i) { // variable i exists in p
		//log.Printf("# found x.%d scaled %s\n", i, scale_i.String())
		p.Terms.Remove(i)
		//log.Printf("# p/%d = %s\n", i, p)
		pp := p2.Multiply(NewConstantPolynomial(scale_i), true)
		//log.Printf("# p2 * %s = %s\n", scale_i, pp)
		p = p.Add(pp, true).Zap()
		//log.Printf("# p + p2 = %s\n", p)
	}
	return p
}

/* Helper: make a copy of a numeric Polynomial.
 */
func (p Polynomial) CopyPolynomial() Polynomial {
	p1 := NewConstantPolynomial(numeric.Zero) // will become our return value
	it := p.Terms.Iterator()
	for it.Next() { // copy all terms of p into p1
		pos := it.Key().(int)
		scale := it.Value().(numeric.Decimal)
		p1.SetTerm(pos, scale)
	}
	return p1
}

/* Internal method: add or subtract 2 polynomials. The high level methods
 * are based on this one.
 * Flag doAdd signals addition or subtraction.
 */
func (p Polynomial) addOrSub(p2 Polynomial, doAdd bool, destructive bool) Polynomial {
	p1 := p.CopyPolynomial() // will become our return value
	it2 := p2.Terms.Iterator()
	for it2.Next() { // inspect all terms of p2
		pos2 := it2.Key().(int)
		scale2 := it2.Value().(numeric.Decimal)
		if !Zero(scale2) {
			scale1 := p1.GetCoeffForTerm(pos2)
			if doAdd {
				scale1 = scale1.Add(scale2) // if present, add a1 + a2
			} else {
				scale1 = scale1.Sub(scale2) // if present, subtract a1 - a2
			}
			p1.SetTerm(pos2, scale1) // we operate on the copy p1
		}
	}
	if destructive {
		p.Terms = p1.Terms
	}
	return p1
}

/* Add two Polynomials. Returns a new Polynomial, except when the
 * 'destructive'-flag is set (then p is altered).
 */
func (p Polynomial) Add(p2 Polynomial, destructive bool) Polynomial {
	/*
		if p.ispair {
			return p.AddPair(p2, destructive)
		} else {
			return p.addOrSub(p2, true, destructive)
		}
	*/
	return p.addOrSub(p2, true, destructive)
}

/* Subtract two Polynomials. Returns a new Polynomial, except when the
 * 'destructive'-flag is set (then p is altered).
 */
func (p Polynomial) Subtract(p2 Polynomial, destructive bool) Polynomial {
	/*
		if p.ispair {
			return p.SubtractPair(p2, destructive)
		} else {
			return p.addOrSub(p2, false, destructive)
		}
	*/
	return p.addOrSub(p2, false, destructive)
}

/* Multiply two Polynomials. One of both must be a constant.
 * p2 will be destroyed.
 */
func (p Polynomial) Multiply(p2 Polynomial, destructive bool) Polynomial {
	/*
		if p.ispair {
			return p.MultiplyPair(p2, destructive)
		} else {
	*/
	p1 := p.CopyPolynomial()      // will become our return value
	c, isconst := p2.IsConstant() // is p2 constant?
	if !isconst {
		c, isconst = p1.IsConstant() // is p1 constant?
		if !isconst {
			panic("neither muliplier nor multiplicand constant")
		}
		p1 = p2 // swap to operate on p2
	}
	it := p1.Terms.Iterator()
	for it.Next() { // multiply all coefficients by c
		pos := it.Key().(int)
		scale := it.Value().(numeric.Decimal)
		p1.SetTerm(pos, scale.Mul(c))
	}
	if destructive {
		p.Terms = p1.Terms
	}
	p1 = p1.Zap()
	return p1
}

/* Divide Polynomial by a numeric (not 0).
 * p2 will be destroyed.
 */
func (p Polynomial) Divide(p2 Polynomial, destructive bool) Polynomial {
	/*
		if p.ispair {
			return p.DividePair(p2, destructive)
		} else {
	*/
	c, isconst := p2.IsConstant() // is p2 constant?
	if !isconst || Zero(c) {
		panic(fmt.Sprintf("illegal divisor: %s", p2.String()))
	} else {
		p2.Terms.Remove(0)
		p2.Terms.Put(0, ConstOne.Div(c)) // now p2 = 1/c
	}
	return p.Multiply(p2, destructive)
	//}
}

/* Eliminate all terms with coefficient=0.
 */
func (p Polynomial) Zap() Polynomial {
	positions := p.Terms.Keys()     // all non-Zero terms of p
	for _, pos := range positions { // inspect terms
		//if !(p.ispair && pos == 0) {
		if scale, _ := p.Terms.Get(pos); Zero(scale.(numeric.Decimal)) {
			p.Terms.Remove(pos) // may lose constant term c
		}
		//}
	}
	if _, ok := p.Terms.Get(0); !ok {
		p.Terms.Put(0, numeric.Zero) // set p = 0: re-introduce c
	}
	//T.Debugf("# Zapped: %s", p.String())
	return p
}

/* Is a Polynomial a constant, i.e. p = { c }? Returns the constant and a flag.
 * If p is a Pair Polynomial, this method will return xpart(c).
 */
func (p Polynomial) IsConstant() (numeric.Decimal, bool) {
	/*
		if p.ispair {
			return p.GetConstantPair().x, p.Terms.Size() == 1
		} else {
			return p.GetCoeffForTerm(0), p.Terms.Size() == 1
		}
	*/
	return p.GetCoeffForTerm(0), p.Terms.Size() == 1
}

/* Is a Polynomial a variable?, i.e. a single term with coefficient = 1.
 * Returns the position of the term and a flag.
 */
func (p Polynomial) IsVariable() (int, bool) {
	if p.Terms.Size() == 2 { // ok: p = a*x.i + c
		if Zero(p.GetCoeffForTerm(0)) { // if c == 0
			positions := p.Terms.Keys() // all non-Zero Terms of p, ordered
			pos := positions[1].(int)
			a := p.GetCoeffForTerm(pos)
			if One(a) { // if a.i = 0
				return pos, true
			}
		}
	}
	return -77777, false
}

/* Get the constant term of a polynomial.
 */
func (p Polynomial) GetConstantValue() numeric.Decimal {
	return p.GetCoeffForTerm(0)
}

/* Get the coefficient for term no. i.
 * Example: p = x + 3x.2  => coeff(2) = 3
 */
func (p Polynomial) GetCoeffForTerm(i int) numeric.Decimal {
	var sc interface{}
	var found bool
	sc, found = p.Terms.Get(i)
	if found {
		return sc.(numeric.Decimal)
	} else {
		return numeric.Zero
	}
}

// === Utilities =============================================================

/* Helper: dump all known equations.
 */
func (leq *LinEqSolver) Dump(resolv VariableResolver) {
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println("Dependents:                                                        LEQ")
	it := leq.dependents.Iterator()
	for it.Next() { // for every x.i = p[x.i]
		k := it.Key().(int)
		p := it.Value().(Polynomial)
		fmt.Printf("\t%s = %s\n", TraceStringVar(k, resolv), p.TraceString(resolv))
	}
	fmt.Println("Solved:")
	it = leq.solved.Iterator()
	for it.Next() { // for every x.i = { c }
		k := it.Key().(int)
		p := it.Value().(Polynomial)
		fmt.Printf("\t%s = %s\n", TraceStringVar(k, resolv), p.GetConstantValue().String())
	}
	fmt.Println("----------------------------------------------------------------------")
}

/* Create a string representation for a Polynomial.
 * Uses internal variable representations x.<n> where n corresponds to
 * the variable's real life ID.
 */
func (p Polynomial) String() string {
	return p.TraceString(nil)
}

/* Create a string representation for a Polynomial. Uses a variable name
 * resolver to print 'real' variable identifiers. If no resolver is
 * present, variables are printed in a generic form: +/- a.i x.i, where i is
 * the position of the term. Coefficients are rounded to the 3rd place.
 */
func (p Polynomial) TraceString(resolv VariableResolver) string {
	var buffer bytes.Buffer
	it := p.Terms.Iterator()
	var indent bool = false // no space before first term (usually constant)
	for it.Next() {
		pos := it.Key().(int)
		if pos == 0 { // constant term
			/*
				if p.ispair {
					pc := it.Value().(Pair)
				} else {
					pc := it.Value().(numeric.Decimal).Round(3)
				}
			*/
			pc := it.Value().(numeric.Decimal)
			if resolv == nil {
				buffer.WriteString(fmt.Sprintf("{ %s } ", pc.Round(3).String()))
			} else {
				if !Zero(pc) {
					buffer.WriteString(pc.Round(3).String())
					indent = true
				}
			}
		} else { // variable term
			scale := it.Value().(numeric.Decimal)
			if resolv == nil {
				buffer.WriteString(fmt.Sprintf("{ %s x.%d } ", scale.Round(3).String(), pos))
			} else {
				if indent {
					if scale.LessThan(numeric.Zero) {
						buffer.WriteString(" - ")
					} else if scale.GreaterThan(numeric.Zero) {
						buffer.WriteString(" + ")
					}
				} else {
					indent = true
					if scale.LessThan(numeric.Zero) {
						buffer.WriteString("-")
					}
				}
				if !scale.Abs().Equal(ConstOne) {
					buffer.WriteString(scale.Abs().Round(3).String())
				}
				buffer.WriteString(resolv.GetVariableName(pos))
			}
		}
	}
	return buffer.String()
}

/* Helper for tracing output. Parameter resolv may be nil.
 */
func TraceStringVar(i int, resolv VariableResolver) string {
	if resolv == nil {
		return fmt.Sprintf("x.%d", i)
	} else {
		return resolv.GetVariableName(i)
	}
}
