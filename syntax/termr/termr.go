package termr

import (
	"bytes"
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

// https://hanshuebner.github.io/lmman/fd-con.xml
// https://www.tutorialspoint.com/lisp/lisp_basic_syntax.htm

// TODO Properties: https://www.tutorialspoint.com/lisp/lisp_symbols.htm

// Atom is a type for atomic values (in lists).
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
)

// NullAtom is a zero value for atoms.
var NullAtom Atom = Atom{}

// Type returns an atom's type.
func (a Atom) Type() AtomType {
	return a.typ
}

func atomize(thing interface{}) Atom {
	if thing == nil {
		return NullAtom
	}
	atom := NullAtom
	switch c := thing.(type) {
	case *GCons:
	case int, int32, int64, uint, uint32, uint64:
		atom.Data = c
		atom.typ = NumType
	case string, []byte, []rune:
		atom.Data = c
		atom.typ = StringType
	case bool:
		atom.Data = c
		atom.typ = BoolType
	case Operator:
		atom.Data = c
		atom.typ = OperatorType
	case *Symbol:
		atom.Data = c
		atom.typ = SymbolType
	case *Environment:
		atom.Data = c
		atom.typ = EnvironmentType
	case *Token:
		atom.Data = c
		atom.typ = TokenType
	default:
		atom.Data = c
		atom.typ = UserType
	}
	return atom
}

// Operator is an interface to be implemented by every node being able to
// operate on an argument list.
type Operator interface {
	Call(*GCons) *GCons
}

// Node is a type for a list node.
// A Cons will consist of a Node and a Cdr.
type Node struct {
	atom  Atom
	child *GCons
}

var nullNode Node = Node{atom: Atom{}}

func makeNode(thing interface{}) Node {
	node := nullNode
	if cons, ok := thing.(*GCons); ok {
		node.child = cons
	} else {
		node.atom = atomize(thing)
	}
	return node
}

// Type returns the atom type of a node.
func (node Node) Type() AtomType {
	if node.child == nil {
		return node.atom.typ
	}
	return ConsType
}

func (node Node) String() string {
	if node.child != nil {
		return "(list)"
	}
	switch node.atom.typ {
	case NumType:
		return fmt.Sprintf("%d", node.atom.Data)
	case BoolType:
		return fmt.Sprintf("%v", node.atom.Data)
	case StringType:
		return fmt.Sprintf("\"%s\"", node.atom.Data)
	case TokenType:
		t := node.atom.Data.(*Token)
		return fmt.Sprintf("%s", t.String())
	}
	return fmt.Sprintf("%v", node.atom.Data)
}

// ListString returns a Node within a list representation. Will usually be called
// indirectly with GCons.ListString().
func (node Node) ListString() string {
	if node.child != nil {
		return node.child.ListString()
	}
	return node.String()
}

// GCons is a type for a list cons.
type GCons struct {
	car Node
	cdr *GCons
}

func (l GCons) String() string {
	var cdrstring string
	if l.cdr == nil {
		cdrstring = "∖"
	} else {
		cdrstring = "→"
	}
	return fmt.Sprintf("(%s,%s)", l.car, cdrstring)
}

// ListString returns a string representing a list (or cons).
func (l *GCons) ListString() string {
	if l == nil {
		return "()"
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
		b.WriteString(l.car.ListString())
		l = l.cdr
	}
	b.WriteString(")")
	return b.String()
}

// Atom returns the Car atom of a list or cons.
func (l *GCons) Atom() Atom {
	if l == nil {
		return NullAtom
	}
	return l.car.atom
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
		cons.car = makeNode(e)
		if first == nil {
			first = cons
		} else {
			last.cdr = cons
		}
		last = cons
	}
	return first
}

// Cons creates a cons from a given node (Car) value.
func Cons(car Node, cdr *GCons) *GCons {
	if car == nullNode {
		return cdr
	}
	return &GCons{car: car, cdr: cdr}
}

// Car returns the Car of a list/cons.
func (l *GCons) Car() Node {
	if l == nil {
		return nullNode
	}
	return l.car
}

// Cdr returns the Cdr of a list/node.
func (l *GCons) Cdr() *GCons {
	if l == nil {
		return nil
	}
	return l.cdr
}

// Cddr returns Cdr(Cdr(...)) of a list/node.
func (l *GCons) Cddr() *GCons {
	if l == nil || l.cdr == nil {
		return nil
	}
	return l.cdr.cdr
}

// Length returns the length of a list.
func (l *GCons) Length() int {
	if l == nil {
		return 0
	}
	length := 0
	for l != nil {
		length++
		l = l.Cdr()
	}
	return length
}

func (l *GCons) copyCons() *GCons {
	node := l.car
	return Cons(node, nil)
}

// First returns the frist n elements of a list.
func (l *GCons) First(n int) *GCons {
	if l == nil || n <= 0 {
		return nil
	}
	f := l.copyCons()
	start := f
	l = l.cdr
	for n--; n > 0; n-- {
		if l == nil {
			break
		}
		f.cdr = l.copyCons()
		f, l = f.cdr, l.cdr
	}
	return start
}
