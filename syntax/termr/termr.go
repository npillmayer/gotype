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
)

type Atom struct {
	typ  AtomType
	Data interface{}
}

var NullAtom Atom = Atom{}

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
	default:
		atom.Data = c
		atom.typ = UserType
	}
	return atom
}

type Operator interface {
	Call(*GCons) *GCons
}

type carNode struct {
	atom  Atom
	child *GCons
}

var nullCar carNode = carNode{atom: Atom{}}

func makeCar(thing interface{}) carNode {
	car := nullCar
	if cons, ok := thing.(*GCons); ok {
		car.child = cons
	} else {
		car.atom = atomize(thing)
	}
	return car
}

func (cnode carNode) Type() AtomType {
	if cnode.child == nil {
		return cnode.atom.typ
	}
	return ConsType
}

func (car carNode) String() string {
	if car.child != nil {
		return "(list)"
	}
	switch car.atom.typ {
	case NumType:
		return fmt.Sprintf("%d", car.atom.Data)
	case BoolType:
		return fmt.Sprintf("%v", car.atom.Data)
	case StringType:
		return fmt.Sprintf("\"%s\"", car.atom.Data)
	}
	return fmt.Sprintf("%v", car.atom.Data)
}

func (car carNode) ListString() string {
	if car.child != nil {
		return car.child.ListString()
	}
	return car.String()
}

type GCons struct {
	car carNode
	cdr *GCons
}

func (c GCons) String() string {
	var cdrstring string
	if c.cdr == nil {
		cdrstring = "∖"
	} else {
		cdrstring = "→"
	}
	return fmt.Sprintf("(%s,%s)", c.car, cdrstring)
}

func (c *GCons) ListString() string {
	if c == nil {
		return "()"
	}
	var b bytes.Buffer
	b.WriteString("(")
	first := true
	for c != nil {
		if first {
			first = false
		} else {
			b.WriteString(" ")
		}
		b.WriteString(c.car.ListString())
		c = c.cdr
	}
	b.WriteString(")")
	return b.String()
}

func (l *GCons) Atom() Atom {
	if l == nil {
		return NullAtom
	}
	return l.car.atom
}

func List(elements ...interface{}) *GCons {
	if len(elements) == 0 {
		return nil
	}
	last := &GCons{}
	var first *GCons
	for _, e := range elements {
		cons := &GCons{}
		cons.car = makeCar(e)
		if first == nil {
			first = cons
		} else {
			last.cdr = cons
		}
		last = cons
	}
	return first
}

func Cons(thing interface{}, cdr *GCons) *GCons {
	if thing == nil {
		return cdr
	}
	carnode := makeCar(thing)
	return &GCons{car: carnode, cdr: cdr}
}

func (l *GCons) Car() *GCons {
	if l == nil {
		return nil
	}
	if l.car.Type() == ConsType {
		return l.car.child
	}
	return &GCons{car: l.car}
}

func (l *GCons) Cdr() *GCons {
	if l == nil {
		return nil
	}
	return l.cdr
}

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
