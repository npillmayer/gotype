package terex

import (
	"bytes"
	"fmt"
)

// Symbol is a type for language symbols (stored in the Environment).
// A symbol can change its value. A value may be any atom or s-expr type.
// A value of nil means the symbol is not yet bound.
type Symbol struct {
	Name  string
	props properties
	Value Atom
}

// newSymbol creates a new symbol for a given initial value (which may be nil).
func newSymbol(name string, thing interface{}) *Symbol {
	return &Symbol{
		Name:  name,
		props: makeProps(),
		Value: Atomize(thing),
	}
}

func (sym Symbol) String() string {
	return fmt.Sprintf("%s:%s(%s)", sym.Name, sym.Value.Type().String(), sym.Value.String())
}

// IsAtom returns true if a symbol represents an atom (not a cons).
func (sym Symbol) IsAtom() bool {
	return sym.Value.Type() != ConsType
}

// Get returns a property value for a given key.
func (sym *Symbol) Get(key string) Atom {
	value := sym.props.Get(key, sym)
	if value == nil {
		return NilAtom
	}
	return Atomize(value)
}

type properties []keyvalue
type keyvalue struct {
	Key   string
	Value interface{}
}

func makeProps() properties {
	return make([]keyvalue, 0, 2)
}

func (p properties) Get(key string, sym *Symbol) interface{} {
	if key == "type" {
		return sym.Value.Type()
	}
	if key == "name" {
		return sym.Name
	}
	for _, kv := range p {
		if kv.Key == key {
			return kv.Value
		}
	}
	return nil
}

func (p properties) Set(key string, value interface{}) properties {
	found := false
	for _, kv := range p {
		if kv.Key == key {
			kv.Value = value
			found = true
		}
	}
	if !found {
		return append(p, keyvalue{key, value})
	}
	return p
}

// Token represents a grammar terminal, and a corresponding input token, respectively.
type Token struct {
	Name  string
	Value int
}

func (t Token) String() string {
	return t.Name
}

// --- Environments ----------------------------------------------------------

// Environment is a type for a symbol environment.
type Environment struct {
	name      string
	parent    *Environment
	dict      map[string]*Symbol
	lastError error
}

// GlobalEnvironment is the base environment all other environments stem from.
var GlobalEnvironment *Environment = &Environment{
	name: "#global",
	dict: make(map[string]*Symbol),
}

func initGlobalEnvironment() {
	defun("+", _Add)
}

func defun(opname string, funcBody Mapper) {
	opsym := GlobalEnvironment.Intern(opname, false)
	opsym.Value = Atomize(&internalOp{sym: opsym, Execute: funcBody})
	T().Debugf("new interal op %s = %v", opsym.Name, opsym.Value)
}

type internalOp struct {
	sym     *Symbol
	Execute Mapper
}

func (iop *internalOp) Call(el Element) Element {
	return iop.Execute(el)
}

func (iop *internalOp) String() string {
	if iop.sym != nil {
		return iop.sym.Name
	}
	return "internal"
}

var _ Operator = &internalOp{}

// NewEnvironment creates a new environment.
func NewEnvironment(name string, parent *Environment) *Environment {
	if parent == nil {
		parent = GlobalEnvironment
	}
	return &Environment{
		name:   name,
		parent: parent,
		dict:   make(map[string]*Symbol),
	}
}

// FindSymbol checks wether a symbol is defined in env and returns it, if found.
// Otherwise nil is returned.
func (env *Environment) FindSymbol(name string, inherit bool) *Symbol {
	var sym *Symbol
	var ok bool
	for env != nil {
		sym, ok = env.dict[name]
		if ok {
			return sym
		}
		if inherit {
			env = env.parent
		} else {
			break
		}
	}
	return nil
}

// Intern interns a symbol name as a symbol, returning a reference to that symbol.
// If the symbol already exists, the existing symbol is returned.
// Parameter inherit dictates wether ancestor environments should be searched, too,
// to detect the symbol.
func (env *Environment) Intern(name string, inherit bool) *Symbol {
	sym := env.FindSymbol(name, inherit)
	if sym == nil {
		sym = &Symbol{Name: name}
		env.dict[name] = sym
	}
	return sym
}

func (env *Environment) String() string {
	return env.name
}

// Dump is a debugging helper, listing all known symbols in env.
func (env *Environment) Dump() string {
	var b bytes.Buffer
	b.WriteString(env.String())
	b.WriteString(" {\n")
	for k, v := range env.dict {
		b.WriteString(fmt.Sprintf("    %s = %v\n", k, v))
	}
	b.WriteString("}\n")
	return b.String()
}
