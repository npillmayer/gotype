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
	"fmt"
	"sync"
)

// Symbol is a type for language symbols (stored in the Environment).
// A symbol can change its value. A value may be any atom or s-expr type.
// A value of nil means the symbol is not yet bound.
type Symbol struct {
	Name  string
	props properties
	Value Element
}

// newSymbol creates a new symbol for a given initial value (which may be nil).
func newSymbol(name string, thing interface{}) *Symbol {
	return &Symbol{
		Name:  name,
		props: makeProps(),
		Value: Elem(thing),
	}
}

var nilSymbol = newSymbol("<nil symbol>", nil)

func (sym Symbol) String() string {
	return fmt.Sprintf("%s:%s(%s)", sym.Name, sym.ValueType().String(), sym.Value.String())
}

// ValueType returns the type of a symbols value
func (sym Symbol) ValueType() AtomType {
	if sym.Value.IsNil() {
		return NoType
	}
	if sym.Value.IsAtom() {
		return sym.Value.thing.(Atom).Type()
	}
	return ConsType
}

func (sym Symbol) AtomValue() interface{} {
	if sym.Value.IsNil() {
		return nil
	}
	if sym.ValueType() == ConsType {
		return sym.Value.AsList()
	}
	return sym.Value.AsAtom().Data
}

// IsNil returns true if the symbol's value is NilAtom
func (sym Symbol) IsNil() bool {
	return sym.Value.IsNil()
}

// IsAtom returns true if a symbol represents an atom (not a cons).
/* func (sym Symbol) IsAtom() bool {
	return sym.Value.IsAtom()
} */

// IsOperatorType returns true if a symbol represents an atom (not a cons).
func (sym Symbol) IsOperatorType() bool {
	if !sym.Value.IsAtom() {
		return false
	}
	atom := sym.Value.AsAtom()
	return atom.Type() == OperatorType
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
		return sym.ValueType()
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
	Name    string
	TokType int
	Token   interface{}
	Value   interface{}
}

func (t Token) String() string {
	if t.Value == nil {
		return t.Name
	}
	switch v := t.Value.(type) {
	case float64:
		return fmt.Sprintf("%s[%g]", t.Name, v)
	case string:
		return fmt.Sprintf("%s[\"%s\"]", t.Name, v)
	}
	return t.Name + "[?]"
}

// --- Environments ----------------------------------------------------------

// Environment is a type for a symbol environment.
type Environment struct {
	name      string
	parent    *Environment
	dict      map[string]*Symbol
	lastError error
	Resolver  SymbolResolver
	AST       *GCons
}

type SymbolResolver interface {
	Resolve(Atom, *Environment, bool) (Element, error)
}

// GlobalEnvironment is the base environment all other environments stem from.
var GlobalEnvironment *Environment = NewEnvironment("#global", nil)

var initOnce sync.Once // monitors one-time initialization of global environment

// InitGlobalEnvironment initializes the global environment. It is guarded against
// multiple execution. Without calling this, the "native" operators will not be
// found in the symbol table.
func InitGlobalEnvironment() {
	initOnce.Do(func() {
		// GlobalEnvironment.Defn("quote", func(e Element, env *Environment) Element {
		// 	return e
		// })
		GlobalEnvironment.Defn("list", func(e Element, env *Environment) Element {
			// (:list a b c) =>  a.call(b c)
			list := e.AsList()
			if list.Length() == 0 { //  () => nil  [empty list is nil]
				return Elem(nil)
			}
			e = Eval(Elem(list), env)
			return e
		})
	})
}

// NewEnvironment creates a new environment.
func NewEnvironment(name string, parent *Environment) *Environment {
	return &Environment{
		name:   name,
		parent: parent,
		dict:   make(map[string]*Symbol),
	}
}

// Defn defines a new operator and stores its symbol in the given environment.
// funcBody is the operator function, called during eval().
func (env *Environment) Defn(opname string, funcBody Mapper) *Symbol {
	opsym := env.Intern(opname, false)
	opsym.Value = Elem(&internalOp{sym: opsym, call: funcBody})
	T().Debugf("new interal op %s = %v", opsym.Name, opsym.Value)
	return opsym
}

// Def defines a symbol and stores a value for it, if any.
func (env *Environment) Def(symname string, value Element) *Symbol {
	sym := env.Intern(symname, false)
	sym.Value = value
	T().Debugf("new interal sym %s = %v", sym.Name, sym.Value)
	return sym
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

// Error sets an error occuring in this environment.
func (env *Environment) Error(e error) {
	env.lastError = e
}

// LastError returns the last error occuring in this environment.
func (env *Environment) LastError() error {
	return env.lastError
}

// Dump is a debugging helper, listing all known symbols in env.
func (env *Environment) Dump() string {
	var b bytes.Buffer
	b = env.dumpEnv(b)
	return b.String()
}

func (env *Environment) dumpEnv(b bytes.Buffer) bytes.Buffer {
	b.WriteString(env.String())
	b.WriteString(" {\n")
	for k, v := range env.dict {
		b.WriteString(fmt.Sprintf("    %s = %v\n", k, v))
	}
	b.WriteString("}\n")
	if env.parent != nil {
		b = env.parent.dumpEnv(b)
	}
	return b
}

// Operator is an interface to be implemented by every operator-symbol, i.e., one
// being able to operate on an argument list.
type Operator interface {
	String() string                     // returns the string representation of this operator
	Call(Element, *Environment) Element // takes and returns *GCons or Node
	//IsQuoter() bool                     // does this operator prevent eval of arguments?
	//Quote(Element, *Environment) Element // takes and returns *GCons or Node
}

// Internal operators implement the Operator interface.
type internalOp struct {
	sym  *Symbol
	call Mapper
	//quoter Mapper
}

func (iop *internalOp) Call(el Element, env *Environment) Element {
	if iop.call == nil {
		return el
	}
	return iop.call(el, env)
}

// func (iop *internalOp) Quote(el Element, env *Environment) Element {
// 	// TODO is env needed for internal ops?
// 	if iop.call == nil {
// 		if el.IsAtom() {
// 			return Elem(Cons(Atomize(iop), Cons(el.AsAtom(), nil)))
// 		}
// 		return Elem(Cons(Atomize(iop), el.AsList()))
// 	}
// 	return iop.call(el, env)
// }

func (iop *internalOp) String() string {
	if iop.sym != nil {
		return iop.sym.Name
	}
	return "internal"
}
