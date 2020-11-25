package terex

/*
BSD License

Copyright (c) 2019–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

/*
https://www.tutorialspoint.com/lisp/lisp_discussion.htm :
Lisp is the second-oldest high-level programming language after Fortran and has
changed a great deal since its early days, and a number of dialects have existed
over its history. Today, the most widely known general-purpose Lisp dialects are
Common Lisp and Scheme. Lisp was invented by John McCarthy in 1958 while he was at
the Massachusetts Institute of Technology (MIT).
*/

// Clojure Script: http://cljs.github.io/api/
//                 https://funcool.github.io/clojurescript-unraveled/

// https://hanshuebner.github.io/lmman/fd-con.xml
// https://www.tutorialspoint.com/lisp/lisp_basic_syntax.htm

// Properties: https://www.tutorialspoint.com/lisp/lisp_symbols.htm

// Atom is a type for atomic values (in lists).
// A cons will consist of an atom and a cdr.
type Atom struct {
	typ  AtomType
	Data interface{}
}

// AtomType is a type specifier for an atom.
type AtomType int

//go:generate stringer -type AtomType
const (
	NoType AtomType = iota
	ConsType
	VarType
	NumType
	StringType
	BoolType
	OperatorType
	TokenType
	EnvironmentType
	UserType
	AnyType
	ErrorType
)

// NilAtom is a zero value for atoms.
var NilAtom Atom = Atom{} // NIL

// Type returns an atom's type.
func (a Atom) Type() AtomType {
	return a.typ
}

// Atomize creates an Atom from an untyped value.
func Atomize(thing interface{}) Atom {
	if thing == nil {
		return NilAtom
	}
	if a, ok := thing.(Atom); ok {
		return a
	}
	atom := Atom{Data: thing}
	switch c := thing.(type) {
	case *GCons:
		atom.typ = ConsType
	case AtomType:
		atom.typ = c
		atom.Data = nil
		T().Debugf("atomize(%s) = %v", thing, atom)
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		f, err := toFloat(c)
		if err != nil {
			return ErrorAtom(err.Error())
		}
		atom.typ = NumType
		atom.Data = f
	case string, []byte, []rune:
		atom.typ = StringType
	case bool:
		atom.typ = BoolType
	case Operator:
		atom.typ = OperatorType
	case *Symbol:
		atom.typ = VarType
	case *Token:
		atom.typ = TokenType
	case *Environment:
		atom.typ = EnvironmentType
	case error:
		atom.typ = ErrorType
	default:
		atom.typ = UserType
	}
	return atom
}

// ErrorAtom returns an error message wrapped in an Atom.
func ErrorAtom(msg string) Atom {
	return Atomize(errors.New(msg))
}

// IsAtom returns t.
func (a Atom) IsAtom() Atom {
	return Atomize(true)
}

func (a Atom) String() string {
	if a == NilAtom {
		return "NIL"
	}
	if a.typ == ConsType {
		if a.Data == nil {
			return "()"
		}
		return "(list)"
	}
	switch a.typ {
	case NumType:
		return fmt.Sprintf("%g", a.Data)
	case BoolType:
		return fmt.Sprintf("%v", a.Data)
	case StringType:
		return fmt.Sprintf("\"%s\"", a.Data)
	case TokenType:
		if a.Data == nil {
			return ":any"
		}
		t := a.Data.(*Token)
		return fmt.Sprintf(":%s", t.String())
	case OperatorType:
		if a.Data == nil {
			return "Op:any"
		}
		o := a.Data.(Operator)
		return fmt.Sprintf("#%s", o.String())
	}
	return fmt.Sprintf("%s[%v]", a.typ, a.Data)
}

// ListString returns an Atom's string representation within a list. Will usually be called
// indirectly with GCons.ListString().
func (a Atom) ListString() string {
	if a.typ == ConsType {
		if a.Data == nil {
			return "NIL"
		}
		return a.Data.(*GCons).ListString()
	}
	return a.String()
}

// ---------------------------------------------------------------------------

// GCons is a type for a list cons.
type GCons struct {
	Car Atom
	Cdr *GCons
}

func (l GCons) String() string {
	var cdrstring string
	if l.Cdr == nil {
		cdrstring = "∖"
	} else {
		cdrstring = "→"
	}
	return fmt.Sprintf("(%s,%s)", l.Car, cdrstring)
}

// ListString returns a string representing a list (or cons).
func (l *GCons) ListString() string {
	if l == nil {
		return "<NIL>"
	}
	var b bytes.Buffer
	b.WriteString("(")
	first := true
	for l != nil {
		if first {
			first = false
		} else {
			b.WriteString(" ")
		}
		b.WriteString(l.Car.ListString())
		l = l.Cdr
	}
	b.WriteString(")")
	return b.String()
}

// IndentedListString returns a string representing a list (or cons).
func (l *GCons) IndentedListString() string {
	var bf bytes.Buffer
	bf = l.indLString(bf, 0)
	return bf.String()
}

func (l *GCons) indLString(bf bytes.Buffer, ind int) bytes.Buffer {
	if l == nil {
		bf.WriteString("<NIL>")
	}
	bf.WriteString("(")
	first := true
	for l != nil {
		if first {
			first = false
		} else {
			bf.WriteString("\n")
			bf.WriteString(indentation[:(ind+1)*3])
		}
		if l.Car.typ == ConsType {
			bf = l.Car.Data.(*GCons).indLString(bf, ind+1)
		} else {
			bf.WriteString(l.Car.String())
		}
		l = l.Cdr
	}
	bf.WriteString(")")
	return bf
}

var indentation = "                                                                   "

// IsAtom returns false, i.e. NIL.
func (l *GCons) IsAtom() Atom {
	return NilAtom
}

// IsLeaf returns true if this node does have neither a Cdr nor
// a left child.
func (l *GCons) IsLeaf() bool {
	return l.Cdr == nil && (l.Car.typ != ConsType || l.Car.Data == nil)
}

// QuotedList makes a list from given elements, quoting them.
func QuotedList(things ...interface{}) *GCons {
	return makeList(true, things)
}

// List makes a list from given elements.
func List(things ...interface{}) *GCons {
	return makeList(false, things)
}

func makeList(quoted bool, things []interface{}) *GCons {
	if len(things) == 0 {
		return nil
	}
	last := &GCons{}
	var first *GCons
	for _, e := range things {
		cons := &GCons{}
		if quoted {
			cons.Car = Atomize(e)
		} else if sym, ok := e.(*Symbol); ok {
			cons.Car = Atomize(sym.Value)
		} else {
			cons.Car = Atomize(e)
		}
		if first == nil {
			first = cons
		} else {
			last.Cdr = cons
		}
		last = cons
	}
	return first
}

// Cons creates a cons from a given Atom and a Cdr.
func Cons(car Atom, cdr *GCons) *GCons {
	// if car == NilAtom {
	// 	return cdr
	// }
	return &GCons{Car: car, Cdr: cdr}
}

// First returns the Car atom of a list/cons.
func (l *GCons) First() Atom {
	if l == nil {
		return NilAtom
	}
	return l.Car
}

// Rest returns the Cdr of a list/node.
func (l *GCons) Rest() *GCons {
	if l == nil {
		return nil
	}
	return l.Cdr
}

// Tee returns the Car as a list, if it is of sublist-type, nil otherwise.
func (l *GCons) Tee() *GCons {
	if l == nil || l.Car.typ != ConsType || l.Car.Data == nil {
		return nil
	}
	return l.Car.Data.(*GCons)
}

// Cadr returns Cdr(Car(...)) of a list/node.
func (l *GCons) Cadr() *GCons {
	if l == nil || l.Car.typ != ConsType || l.Car.Data == nil {
		return nil
	}
	return l.Car.Data.(*GCons).Cdr
}

// Cdar returns Car(Cdr(...)) of a list/node.
func (l *GCons) Cdar() Atom {
	if l == nil || l.Cdr == nil {
		return NilAtom
	}
	return l.Cdr.Car
}

// Cddr returns Cdr(Cdr(...)) of a list/node.
func (l *GCons) Cddr() *GCons {
	if l == nil || l.Cdr == nil {
		return nil
	}
	return l.Cdr.Cdr
}

// Length returns the length of a list.
func (l *GCons) Length() int {
	if l == nil {
		return 0
	}
	length := 0
	for l != nil {
		length++
		l = l.Cdr
	}
	return length
}

func (l *GCons) copyCons() *GCons {
	if l == nil {
		return nil
	}
	node := l.Car
	return Cons(node, nil)
}

// FirstN returns the frist n elements of a list.
func (l *GCons) FirstN(n int) *GCons {
	if l == nil || n <= 0 {
		return nil
	}
	f := l.copyCons()
	start := f
	l = l.Cdr
	for n--; n > 0; n-- {
		if l == nil {
			break
		}
		f.Cdr = l.copyCons()
		f, l = f.Cdr, l.Cdr
	}
	return start
}

// Last returns the last element of a list or nil.
func (l *GCons) Last() *GCons {
	if l == nil {
		return nil
	}
	for l.Cdr != nil {
		l = l.Cdr
	}
	return l
}

// Concat appends a list or element at the end of the copy of a list.
func (l *GCons) Concat(other *GCons) *GCons {
	if l == nil {
		return other
	}
	infinity := 999999
	cc := l.FirstN(infinity) // make a copy
	cc.Last().Cdr = other
	return cc
}

// Append destructively appends a list to a list.
func (l *GCons) Append(other *GCons) *GCons {
	if l == nil {
		return other
	}
	l.Last().Cdr = other
	return l
}

// Branch destructively appends a list as a sublist to l.
func (l *GCons) Branch(other *GCons) *GCons {
	tee := Cons(Atomize(other), nil)
	if l == nil {
		l = tee
	} else {
		l.Last().Cdr = tee
	}
	return l
}

// Map applies a mapping-function to every element of a list.
func (l *GCons) Map(mapper Mapper, env *Environment) *GCons {
	return _Map(mapper, Elem(l), env).AsList()
}

// Reduce applies a mapping-function to every element of a list.
func (l *GCons) Reduce(f func(Atom, Atom) Atom, initial Atom, env *Environment) Atom {
	if l.Length() == 0 {
		return initial
	}
	result := f(initial, l.Car)
	if l.Length() > 1 {
		rest := l.Cdr
		for rest != nil {
			result = f(result, rest.Car)
			rest = rest.Cdr
		}
	}
	//T().Debugf("_Map result = %s", result.ListString())
	return result
}

// --- Internal Operations ---------------------------------------------------

// A Mapper takes an atom or list and maps it to an atom or list
type Mapper func(Element, *Environment) Element

type Element struct {
	thing interface{}
}

func Elem(thing interface{}) Element {
	return Element{thing: thing}
}

func DumpElement(el Element) {
	if el.IsAtom() {
		T().Debugf("Dump element: %v", el)
	} else {
		switch e := el.thing.(type) {
		case Element:
			T().Debugf("Dump element: recursive element:")
			DumpElement(e)
		case *GCons:
			T().Debugf("Dump element: list = %s", e.ListString())
		default:
			T().Debugf("Dump element: unknown = %v", e)
		}
	}
}

func (el Element) IsAtom() bool {
	if _, ok := el.thing.(Atom); ok {
		return true
	}
	return false
}

func (el Element) IsNil() bool {
	if t, ok := el.thing.(*GCons); ok {
		if t == nil {
			return true
		}
	}
	return el.thing == nil
}

func (el Element) IsError() bool {
	return el.AsAtom().typ == ErrorType
}

func (el Element) AsAtom() Atom {
	if el.IsAtom() {
		return el.thing.(Atom)
	}
	return Atomize(el.thing.(*GCons))
}

func (el Element) AsList() *GCons {
	if el.IsNil() {
		return nil
	}
	if el.IsAtom() {
		return Cons(el.thing.(Atom), nil)
	}
	if el.thing == nil {
		return nil
	}
	return el.thing.(*GCons)
}

func (el Element) String() string {
	if el.IsAtom() {
		return el.AsAtom().String()
	}
	if el.thing == nil {
		return "nil"
	}
	return el.AsList().ListString()
}

func _First(args Element) Element {
	if args.IsAtom() {
		return args
	}
	return Elem(args.AsList().Car)
}

func _Rest(args Element) Element {
	return Elem(args.AsList().Cdr)
}

func _Identity(args Element) Element {
	return args
}

func _Add(args Element) Element {
	T().Infof("_Add args=%s", args.String())
	if args.IsAtom() {
		if a := args.AsAtom(); a.typ == NumType {
			return Elem(a)
		}
	}
	sum := 0.0
	arglist := args.AsList()
	for arglist != nil {
		T().Infof("     arg=%v", arglist.Car)
		if arglist.Car.Type() == NumType {
			sum += arglist.Car.Data.(float64)
		} else if arglist.Car.Type() == TokenType {
			v := arglist.Car.Data.(*Token).Value
			f, err := toFloat(v)
			if err != nil {
				return Elem(ErrorAtom)
			}
			sum += f
		} else {
			return Elem(ErrorAtom)
		}
		arglist = arglist.Cdr
	}
	return Elem(Atomize(sum))
}

func _Inc(args Element) Element {
	if args.IsAtom() {
		if a := args.AsAtom(); a.typ == NumType {
			return Elem(Atomize(a.Data.(float64) + 1))
		}
	}
	return Elem(ErrorAtom)
}

// func _Quote(op Element, args Element) Element {
// 	if args.IsAtom() {
// 		return args
// 	}
// 	if op.IsAtom() {
// 		qargs := GlobalEnvironment.quoteList(args.AsList())
// 		return Elem(Cons(op.AsAtom(), qargs.AsList()))
// 	}
// 	panic(fmt.Errorf("_Quote called with op=list %s", op))
// }

func _Eval(args Element, env *Environment) Element {
	if args.IsAtom() {
		return args
	}
	return evalList(args.AsList(), env)
}

func _ErrorMapper(err error) Mapper {
	return func(Element, *Environment) Element {
		return Elem(ErrorAtom(err.Error()))
	}
}

func _Map(mapper Mapper, args Element, env *Environment) Element {
	arglist := args.AsList()
	if arglist == nil {
		return Elem(nil)
	}
	r := mapper(Elem(arglist.Car), env)
	T().Debugf("Map: mapping(%s) = %s", arglist.Car, r)
	if arglist.Cdr == nil {
		return r
	}
	result := Cons(r.AsAtom(), nil)
	iter := result
	cons := arglist.Cdr
	for cons != nil {
		el := mapper(Elem(cons.Car), env)
		T().Debugf("Map: mapping(%s) = %s", cons.Car, el)
		if el.IsError() {
			return el
		}
		iter.Cdr = Cons(el.AsAtom(), nil)
		iter = iter.Cdr
		cons = cons.Cdr
	}
	//T().Debugf("_Map result = %s", result.ListString())
	return Elem(result)
}

// --- Matching --------------------------------------------------------------

/*
Match an s-expr to a pattern.

From https://hanshuebner.github.io/lmman/fd-con.xml:

list-match-p object pattern

object is evaluated and matched against pattern; the value is t if it matches, nil otherwise.
pattern is made with backquotes (Aids for Defining Macros); whereas normally a backquote
expression says how to construct list structure out of constant and variable parts, in
this context it says how to match list structure against constants and variables. Constant
parts of the backquote expression must match exactly; variables preceded by commas can
match anything but set the variable to what was matched. (Some of the variables may be
set even if there is no match.) If a variable appears more than once, it must match
the same thing (equal list structures) each time. ,ignore can be used to match anything
and ignore it. For example, `(x (,y) . ,z) is a pattern that matches a list of length
at least two whose first element is x and whose second element is a list of length one;
if a list matches, the caadr of the list is stored into the value of y and the cddr of
the list is stored into z. Variables set during the matching remain set after the
list-match-p returns; in effect, list-match-p expands into code which can setq the
variables. If the match fails, some or all of the variables may already have been set.

Example:

    (list-match-p foo
              `((a ,x) ,ignore . ,c))

is t if foo's value is a list of two or more elements, the first of which is a list
of two elements; and in that case it sets x to (cadar foo) and c to (cddr foo).

List l is the pattern, other is the argument to be matched against the pattern.
*/
func (l *GCons) Match(other *GCons, env *Environment) bool {
	T().Debugf("Match: %s vs %s", l.ListString(), other.ListString())
	if l == nil {
		T().Debugf("l=nil")
	}
	if l == nil {
		return other == nil
	}
	T().Debugf("l.type=%s, l.data=%v", l.Car.Type(), l.Car.Data)
	if l != nil && l.Car.Type() == ConsType && l.Car.Data == nil {
		return true
	}
	if other == nil {
		return false
	}
	if !matchAtom(l.Car, other.Car, env) {
		return false
	}
	return l.Cdr.Match(other.Cdr, env)
}

// func matchCar(car Node, otherNode Node, env *Environment) bool {
// 	T().Debugf("Match Car: %s vs %s", car, otherNode)
// 	if car == nullNode {
// 		return otherNode == nullNode
// 	}
// 	if car.Type() == VarType {
// 		return bindSymbol(car, otherNode, env)
// 	}
// 	if car.Type() == ConsType {
// 		if otherNode.Type() != ConsType {
// 			return false
// 		}
// 		return car.child.Match(otherNode.child, env)
// 	}
// 	return matchAtom(car.atom, otherNode.atom)
// }

func matchAtom(atom Atom, otherAtom Atom, env *Environment) bool {
	T().Debugf("Match Atom: %v vs %v", atom, otherAtom)
	if atom == NilAtom {
		return otherAtom == NilAtom
	}
	if otherAtom == NilAtom {
		return false
	}
	if atom.Type() == VarType {
		return bindSymbol(atom, otherAtom, env)
	}
	typeMatches, doMatchData := typeMatch(atom.typ, otherAtom.typ)
	if !typeMatches {
		return false
	}
	if doMatchData {
		return dataMatch(atom.Data, otherAtom.Data, atom.typ, env)
	}
	return true
}

func bindSymbol(symatom Atom, value Atom, env *Environment) bool {
	sym, ok := symatom.Data.(*Symbol)
	if !ok {
		return false
	}
	T().Debugf("binding symbol %s to %s", sym.String(), value.String())
	if sym.Value == NilAtom {
		sym.Value = value // bind it
		T().Debugf("bound symbol %s", sym.String())
		return true
	}
	return matchAtom(sym.Value, value, env)
}

// typeMatch returns (typesAreMatching, mustMatchValue)
func typeMatch(t1 AtomType, t2 AtomType) (bool, bool) {
	if t1 == AnyType {
		return true, false
	}
	if t1 == t2 {
		return true, true
	}
	T().Debugf("No type match: %s vs %s", t1, t2)
	return false, true
}

func dataMatch(d1 interface{}, d2 interface{}, t AtomType, env *Environment) bool {
	if d1 == nil {
		return true
	}
	if t == TokenType && d2 != nil {
		tok1, _ := d1.(*Token)
		if tok2, ok := d2.(*Token); ok {
			if tok1.TokType == tok2.TokType { // only tokval must match
				return true
			}
		}
	}
	if t == ConsType {
		return d1.(*GCons).Match(d2.(*GCons), env)
	}
	return d1 == d2
}

// ----------------------------------------------------------------------

var floatType = reflect.TypeOf(float64(0))

func toFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}
