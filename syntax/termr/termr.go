package termr

import (
	"bytes"
	"fmt"
)

// https://hanshuebner.github.io/lmman/fd-con.xml
// https://www.tutorialspoint.com/lisp/lisp_basic_syntax.htm

type AtomType int

const (
	AnyType AtomType = iota
	NumType
	StringType
	BoolType
	ConsType
	OperatorType
	TokenType
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

type Operator interface {
	function(*GCons) *GCons
}

type carNode struct {
	atom  Atom
	child *GCons
}

var nullCar carNode = carNode{atom: Atom{}}

func (cnode carNode) Type() AtomType {
	if cnode.child == nil {
		return cnode.atom.typ
	}
	return ConsType
}

func atomize(thing interface{}) carNode {
	if thing == nil {
		return nullCar
	}
	carnode := nullCar
	switch c := thing.(type) {
	case *GCons:
		carnode.child = c
	case int, int32, int64, uint, uint32, uint64:
		carnode.atom.Data = thing
		carnode.atom.typ = NumType
	case string, []byte, []rune:
		carnode.atom.Data = thing
		carnode.atom.typ = StringType
	case bool:
		carnode.atom.Data = thing
		carnode.atom.typ = BoolType
	case Operator:
		carnode.atom.Data = thing
		carnode.atom.typ = OperatorType
	default:
		carnode.atom.Data = thing
		carnode.atom.typ = UserType
	}
	return carnode
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
		cons.car = atomize(e)
		if first == nil {
			first = cons
		} else {
			last.cdr = cons
		}
		last = cons
	}
	return first
}

func Cons(car interface{}, cdr *GCons) *GCons {
	if car == nil {
		return cdr
	}
	carnode := atomize(car)
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

/*
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
An equivalent expression would be

(let ((tem foo))
  (and (consp tem)
       (consp (car tem))
       (eq (caar tem) 'a)
       (consp (cdar tem))
       (progn (setq x (cadar tem)) t)
       (null (cddar tem))
       (consp (cdr tem))
	   (setq c (cddr tem))))

but list-match-p is faster.
*/
func (l *GCons) Match(other *GCons, matcher func(Atom, Atom) bool) bool {
	if l == nil {
		return other == nil
	}
	if other == nil {
		return false
	}
	match := true
	if t := l.car.Type(); t != ConsType {
		if other.car.Type() == t {
			match = match && matcher(l.car.atom, other.car.atom)
		} else {
			match = false
		}
	} else {
		match = match && l.Car().Match(other.Car(), matcher)
	}
	if match {
		match = match && l.Cdr().Match(other.Cdr(), matcher)
	}
	return match
}
