package runtime

import (
	"fmt"
)

/*
----------------------------------------------------------------------

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
3. Neither the name of Norbert Pillmayer or the names of its contributors
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

 * Symbol table for variables. Symbol tables are attached to scopes.
 * Scopes are organized in a tree.

*/

// --- Symbols ------------------------------------------------

/* Symbols (constants, variables, etc.) to be stored.
 */

// Every symbol has a serial ID.
var serialID int = 1 // may not start with 0 !

// The most basic symbol needs a name to be referenced under.
type Symbol interface {
	GetName() string
	GetID() int
}

// Some symbols may be typed.
type Typed interface {
	GetType() int
}

// Some symbols may be typable.
type Typable interface {
	Typed
	SetType(int) int
}

// Some symbols are non-atomic.
type TreeNode interface {
	GetSibling() TreeNode
	SetSibling(TreeNode)
	GetFirstChild() TreeNode
	SetFirstChild(TreeNode)
}

// Our goto-symbol implements all of the interfaces above.
type StdSymbol struct {
	Name     string
	Id       int
	Symtype  int
	Sibling  TreeNode
	Children TreeNode
}

// Pre-defined symbol types, if you want to use them.
const (
	Undefined int = iota
	IntegerType
	PairType
	PathType
	PenType
	ColorType
	StringType
)

/* Create a new standard symbol, with a new ID.
 */
func NewStdSymbol(nm string) Symbol {
	serialID += 1
	var sym = &StdSymbol{
		Name: nm,
		Id:   serialID,
	}
	return sym
}

/* A debug Stringer for symbols.
 */
func (s *StdSymbol) String() string {
	return fmt.Sprintf("<std-symbol '%s'[%d]:%d>", s.Name, s.Id, s.Symtype)
}

/* Interface Symbol: get the name of the symbol.
 */
func (s *StdSymbol) GetName() string {
	return s.Name
}

/* Interface Symbol: get the ID of the symbol.
 */
func (s *StdSymbol) GetID() int {
	return s.Id
}

/* Interface Typeable: get the symbols type.
 */
func (s *StdSymbol) GetType() int {
	return s.Symtype
}

/* Interface Typeable: set the symbols type.
 */
func (s *StdSymbol) SetType(tp int) int {
	t := s.Symtype
	s.Symtype = tp
	return t
}

/* Complex symbols: get the first sub-symbol.
 */
func (s *StdSymbol) GetFirstChild() TreeNode {
	return s.Children
}

/* Complex symbols: set the first sub-symbol.
 */
func (s *StdSymbol) SetFirstChild(ch TreeNode) {
	s.Children = ch
}

/* Complex symbols: get the next adjacent symbol.
 */
func (s *StdSymbol) GetSibling() TreeNode {
	return s.Sibling
}

/* Complex symbols: set the next adjacent symbol.
 */
func (s *StdSymbol) SetSibling(sibling TreeNode) {
	s.Sibling = sibling
}

/* Complex symbols: add a sub-symbol.
 */
func (s *StdSymbol) AppendChild(ch TreeNode) TreeNode {
	//T().Debug("---> append child %v to %v\n", ch, s)
	if s.Children == nil {
		T().Debugf("appending first child: %s", ch)
		s.Children = ch
		s.GetFirstChild()
	} else {
		next := s.Children
		for ; next.GetSibling() != nil; next = next.GetSibling() {
			// do nothing
		}
		next.SetSibling(ch)
		T().Debugf("appending child: %s\n", next.GetSibling())
	}
	return s
}

// check interface assignability
var _ Symbol = &StdSymbol{}
var _ Typable = &StdSymbol{}
var _ TreeNode = &StdSymbol{}

// === Symbol Tables =========================================================

/* Symbol tables to store symbols (map-like semantics).
 */
type SymbolTable struct {
	Table        map[string]Symbol
	createSymbol func(string) Symbol
}

/* Create an empty symbol table. Clients may pass a function to create new
 * symbols, thus enabling the symbol table to create missing symbols on the
 * fly. If none is given and the demand for a new symbol arises, a StdSymbol
 * is created.
 */
func NewSymbolTable(symcreator func(string) Symbol) *SymbolTable {
	if symcreator == nil {
		symcreator = NewStdSymbol
	}
	var symtab = SymbolTable{
		Table:        make(map[string]Symbol),
		createSymbol: symcreator,
	}
	return &symtab
}

/* Retrieve the symbol creation function.
 */
func (t *SymbolTable) GetSymbolCreator() func(string) Symbol {
	return t.createSymbol
}

/* Check for a symbol in the symbol table.
 * Returns a symbol or nil.
 */
func (t *SymbolTable) ResolveSymbol(symname string) Symbol {
	//t.Lock()
	sym := t.Table[symname] // get symbol by name
	//t.Unlock()
	return sym
}

/* Find a symbol in the table, insert a new one if not found.
 * The symcreator function is used for creating non-existent symbols on
 * the fly. Returns the symbol and a flag, signalling wether the symbol
 * already has been present.
 */
func (t *SymbolTable) ResolveOrDefineSymbol(symname string) (Symbol, bool) {
	if len(symname) == 0 {
		return nil, false
	}
	found := true
	sym := t.ResolveSymbol(symname)
	if sym == nil { // if not already there, insert it
		sym, _ = t.DefineSymbol(symname)
		found = false
	}
	return sym, found
}

/*
Create a new symbol to store into the symbol table.
The symbol's name may not be empty
Overwrites existing symbol with this name, if any.
Returns the new symbol and the previously stored symbol (or nil).
*/
func (t *SymbolTable) DefineSymbol(symname string) (Symbol, Symbol) {
	if len(symname) == 0 {
		return nil, nil
	}
	sym := t.createSymbol(symname)
	old := t.InsertSymbol(sym)
	return sym, old
}

// Insert a pre-created symbol.
func (t *SymbolTable) InsertSymbol(sym Symbol) Symbol {
	old := t.ResolveSymbol(sym.GetName())
	t.Table[sym.GetName()] = sym
	return old
}

// Count the symbols in a symbol table.
func (t *SymbolTable) Size() int {
	return len(t.Table)
}

// Iterate over each symbol in the table, executing a mapper function.
func (t *SymbolTable) Each(mapper func(string, Symbol)) {
	for k, v := range t.Table {
		mapper(k, v)
	}
}

// === Scopes ================================================================

// A named scope, which may contain symbol definitions. Scopes link back to a
// parent scope, forming a tree.
type Scope struct {
	Name   string
	Parent *Scope
	symtab *SymbolTable
}

/* Create a new scope.
 */
func NewScope(nm string, parent *Scope, symcreator func(string) Symbol) *Scope {
	sc := &Scope{
		Name:   nm,
		Parent: parent,
		symtab: NewSymbolTable(symcreator),
	}
	return sc
}

/* Prettyfied Stringer.
 */
func (s *Scope) String() string {
	return fmt.Sprintf("<scope %s>", s.Name)
}

/* Return the name of a scope.
 */
func (s *Scope) GetName() string {
	return s.Name
}

/* Return the symbol table of a scope.
 */
func (s *Scope) Symbols() *SymbolTable {
	return s.symtab
}

/* Define a symbol in the scope. Returns the new symbol and the previously
 * stored symbol under this key, if any.
 */
func (s *Scope) DefineSymbol(symname string) (Symbol, Symbol) {
	return s.symtab.DefineSymbol(symname)
}

/* Find a symbol. Returns the symbol (or nil) and a scope. The scope is
 * the scope (of a scope-tree-path) the symbol was found in.
 */
func (s *Scope) ResolveSymbol(symname string) (Symbol, *Scope) {
	sym := s.symtab.ResolveSymbol(symname)
	if sym != nil {
		return sym, s
	}
	for s.Parent != nil {
		s = s.Parent
		sym, _ = s.ResolveSymbol(symname)
		if sym != nil {
			return sym, s
		}
	}
	return sym, nil
}

// ---------------------------------------------------------------------------

/* Scope tree. Can be treated as a stack during static analysis, thus
 * building a tree from scopes which are pushed an popped to/from the stack.
 */
type ScopeTree struct {
	ScopeBase *Scope
	ScopeTOS  *Scope
}

/* Get the current scope of a stack (TOS).
 */
func (scst *ScopeTree) Current() *Scope {
	if scst.ScopeTOS == nil {
		panic("attempt to access scope from empty stack")
	}
	return scst.ScopeTOS
}

/* Get the outermost scope, containing global symbols.
 */
func (scst *ScopeTree) Globals() *Scope {
	if scst.ScopeBase == nil {
		panic("attempt to access global scope from empty stack")
	}
	return scst.ScopeBase
}

/* Push a scope. A scope is constructed, including a symbol table
 * for variable declarations.
 */
func (scst *ScopeTree) PushNewScope(nm string, symcreator func(string) Symbol) *Scope {
	scp := scst.ScopeTOS
	newsc := NewScope(nm, scp, symcreator)
	if scp == nil { // the new scope is the global scope
		scst.ScopeBase = newsc // make new scope anchor
	}
	scst.ScopeTOS = newsc // new scope now TOS
	T().P("scope", newsc.Name).Debugf("pushing new scope")
	return newsc
}

/* Pop a scope.
 */
func (scst *ScopeTree) PopScope() *Scope {
	if scst.ScopeTOS == nil {
		panic("attempt to pop scope from empty stack")
	}
	sc := scst.ScopeTOS
	T().Debugf("popping scope [%s]", sc.Name)
	scst.ScopeTOS = scst.ScopeTOS.Parent
	return sc
}
