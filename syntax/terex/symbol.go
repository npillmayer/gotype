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
	"sync"
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
	AST       *GCons
}

// GlobalEnvironment is the base environment all other environments stem from.
var GlobalEnvironment *Environment = NewEnvironment("#global", nil)

var initOnce sync.Once // monitors one-time initialization of global environment

// InitGlobalEnvironment initializes the global environment. It is guarded against
// multiple execution. Without calling this, the "native" operators will not be
// found in the symbol table.
func InitGlobalEnvironment() {
	initOnce.Do(func() {
		Defun("+", _Add, nil, GlobalEnvironment)
		Defun("quote", _Quote, nil, GlobalEnvironment)
		Defun("list", _ErrorMapper(errors.New("list used as function call")), _Identity,
			GlobalEnvironment)
	})
}

// Defun defines a new operator and stores its symbol in the given environment.
// funcBody is the operator function, called during eval().
func Defun(opname string, funcBody Mapper, quoter Mapper, env *Environment) {
	opsym := GlobalEnvironment.Intern(opname, false)
	opsym.Value = Atomize(&internalOp{sym: opsym, caller: funcBody, quoter: quoter})
	T().Debugf("new interal op %s = %v", opsym.Name, opsym.Value)
}

type internalOp struct {
	sym    *Symbol
	caller Mapper
	quoter Mapper
}

func (iop *internalOp) Call(el Element, env *Environment) Element {
	// TODO is env needed for internal ops?
	T().Errorf("######## iop=%s #################", iop.String())
	if iop.caller == nil {
		return el
	}
	return iop.caller(el)
}

// Quote
// TODO The whole quote-thing is unnecessary. Currently we use it to get rid
// of the #list:op quoting, but this should be replaced by a term rewrite.
// #list is used as a sentinel to stop sequences from flowing in parent nodes.
// It is useful until the first AST is complete. Afterwards, instead of quoting,
// we should rewrite AST nodes of type #list.
//
func (iop *internalOp) Quote(el Element, env *Environment) Element {
	// TODO is env needed for internal ops?
	if iop.quoter == nil {
		if el.IsAtom() {
			return Elem(Cons(Atomize(iop), Cons(el.AsAtom(), nil)))
		}
		return Elem(Cons(Atomize(iop), el.AsList()))
	}
	return iop.quoter(el)
}

func (iop *internalOp) String() string {
	if iop.sym != nil {
		return iop.sym.Name
	}
	return "internal"
}

// NewEnvironment creates a new environment.
func NewEnvironment(name string, parent *Environment) *Environment {
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
