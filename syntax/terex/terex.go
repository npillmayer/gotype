package terex

import (
	"bytes"
	"errors"
	"fmt"
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
	ConsType AtomType = iota
	SymbolType
	NumType
	StringType
	BoolType
	OperatorType
	TokenType
	EnvironmentType
	UserType
	AnyType
	AnyList
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
	atom := Atom{Data: thing}
	switch c := thing.(type) {
	case *GCons:
		atom.typ = ConsType
	case AtomType:
		atom.typ = c
	case int, int32, int64, uint, uint32, uint64:
		atom.typ = NumType
	case string, []byte, []rune:
		atom.typ = StringType
	case bool:
		atom.typ = BoolType
	case Operator:
		atom.typ = OperatorType
	case *Symbol:
		atom.typ = SymbolType
	case *Environment:
		atom.typ = EnvironmentType
	case *Token:
		atom.typ = TokenType
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
	if a.typ == ConsType {
		return "(list)"
	}
	switch a.typ {
	case NumType:
		return fmt.Sprintf("%d", a.Data)
	case BoolType:
		return fmt.Sprintf("%v", a.Data)
	case StringType:
		return fmt.Sprintf("\"%s\"", a.Data)
	case TokenType:
		t := a.Data.(*Token)
		return fmt.Sprintf("T(%s)", t.String())
	}
	return fmt.Sprintf("%v", a.Data)
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
		return "NIL"
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

// IsAtom returns false, i.e. NIL.
func (l *GCons) IsAtom() Atom {
	return NilAtom
}

// List makes a list from given elements.
func List(elements ...interface{}) *GCons {
	if len(elements) == 0 {
		return nil
	}
	last := &GCons{}
	var first *GCons
	for _, e := range elements {
		cons := &GCons{}
		cons.Car = Atomize(e)
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
	if car == NilAtom {
		return cdr
	}
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

// Cadr returns Cdr(Car(...)) of a list/node.
func (l *GCons) Cadr() *GCons {
	if l == nil || l.Car.typ != ConsType || l.Car.Data == nil {
		return nil
	}
	return l.Car.Data.(*GCons).Cdr
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

// --- Internal Operations ---------------------------------------------------

// Operator is an interface to be implemented by every node being able to
// operate on an argument list.
type Operator interface {
	Name() string            // returns the string representation of this operator
	Call(*GCons) interface{} // returns *GCons or Node
}

func _Add(args *GCons) interface{} {
	sum := 0
	for args != nil {
		if args.Car.Type() != NumType {
			return ErrorAtom
		}
		sum += args.Car.Data.(int)
	}
	return Atomize(sum)
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
	T().Debugf("Match: %s vs %s", l, other)
	if l != nil && l.Car.Type() == AnyList {
		return true
	}
	if l == nil {
		return other == nil
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
// 	if car.Type() == SymbolType {
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
	if atom.Type() == SymbolType {
		return bindSymbol(atom, otherAtom, env)
	}
	typeMatches, doMatchData := typeMatch(atom.typ, otherAtom.typ)
	if !typeMatches {
		return false
	}
	if doMatchData {
		return dataMatch(atom.Data, otherAtom.Data, atom.typ)
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
	if t1 == t2 {
		return true, true
	}
	if t1 == AnyType {
		return true, false
	}
	return false, true
}

func dataMatch(d1 interface{}, d2 interface{}, t AtomType) bool {
	if t == OperatorType && d1 != nil {
		if op, ok := d1.(Operator); ok {
			if op.Name() == "Any" {
				return true
			}
		}
	}
	if t == TokenType && d2 != nil {
		tok1, _ := d1.(*Token)
		if tok2, ok := d2.(*Token); ok {
			if tok1.Value == tok2.Value { // only tokval must match
				return true
			}
		}
	}
	T().Errorf("dataMatch()")
	return d1 == d2
}
