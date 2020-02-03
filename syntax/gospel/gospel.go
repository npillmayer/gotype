package gospel

type AtomType int

const (
	NilType AtomType = iota
	NumType
	StringType
	BoolType
	ConsType
	OperatorType
	UserType
)

type atom struct {
	typ  AtomType
	data interface{}
}

type Operator interface {
	function(*GCons) *GCons
}

type carNode struct {
	atom  atom
	child *GCons
}

func (cnode carNode) Type() AtomType {
	if cnode.child == nil {
		return cnode.atom.typ
	}
	return ConsType
}

type GCons struct {
	car carNode
	cdr *GCons
}

func List(elements ...interface{}) *GCons {
	return nil // TODO
}

func Cons(car interface{}, cdr *GCons) *GCons {
	if car == nil {
		return cdr
	}
	carnode := carNode{}
	switch c := car.(type) {
	case *GCons:
		carnode.child = c
	case int, int32, int64, uint, uint32, uint64:
		carnode.atom.data = car
		carnode.atom.typ = NumType
	case string, []byte, []rune:
		carnode.atom.data = car
		carnode.atom.typ = StringType
	case bool:
		carnode.atom.data = car
		carnode.atom.typ = BoolType
	case Operator:
		carnode.atom.data = car
		carnode.atom.typ = OperatorType
	default:
		carnode.atom.data = car
		carnode.atom.typ = UserType
	}
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
