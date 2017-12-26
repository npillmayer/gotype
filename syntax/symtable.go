package syntax

import (
	"fmt"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
)

var T tracing.Trace = tracing.SyntaxTracer

// --- Symbols ------------------------------------------------

/* Symbols (constants, variables, etc.) to be stored.
 */

var serialID int = 1 // may not start with 0 !

type Symbol interface {
	GetName() string
	GetID() int
}

type Typable interface {
	GetType() int
	SetType(int) int
}

type TreeNode interface {
	GetSibling() TreeNode
	SetSibling(TreeNode)
	GetFirstChild() TreeNode
	SetFirstChild(TreeNode)
}

type StdSymbol struct {
	Name     string
	Id       int
	Symtype  int
	Sibling  TreeNode
	Children TreeNode
}

// Pre-defined types, if you want to use them
const (
	Undefined int = iota
	IntegerType
	PairType
	PathType
	PenType
	ColorType
	StringType
)

func NewStdSymbol(nm string) Symbol {
	serialID += 1
	var sym = &StdSymbol{
		Name: nm,
		Id:   serialID,
	}
	return sym
}

func (s *StdSymbol) String() string {
	return fmt.Sprintf("<std-symbol '%s'[%d]:%d>", s.Name, s.Id, s.Symtype)
}

func (s *StdSymbol) GetName() string {
	return s.Name
}

func (s *StdSymbol) GetID() int {
	return s.Id
}

func (s *StdSymbol) GetType() int {
	return s.Symtype
}

func (s *StdSymbol) SetType(tp int) int {
	t := s.Symtype
	s.Symtype = tp
	return t
}

func (s *StdSymbol) GetFirstChild() TreeNode {
	return s.Children
}

func (s *StdSymbol) SetFirstChild(ch TreeNode) {
	s.Children = ch
}

func (s *StdSymbol) GetSibling() TreeNode {
	return s.Sibling
}

func (s *StdSymbol) SetSibling(sibling TreeNode) {
	s.Sibling = sibling
}

func (s *StdSymbol) AppendChild(ch TreeNode) TreeNode {
	//T.Debug("---> append child %v to %v\n", ch, s)
	if s.Children == nil {
		T.Debugf("appending first child: %s", ch)
		s.Children = ch
		s.GetFirstChild()
	} else {
		next := s.Children
		//T.Debug("skipping first child: %s\n", next)
		for ; next.GetSibling() != nil; next = next.GetSibling() {
			//T.Debug("skipping existing child: %s\n", next)
			// do nothing
		}
		//fmt.Debug("---> next = %v\n", next)
		next.SetSibling(ch)
		T.Debug("appending child: %s\n", next.GetSibling())
	}
	return s
}

//func (s *StdSymbol)

// check interface assignability
var _ Symbol = &StdSymbol{}
var _ Typable = &StdSymbol{}
var _ TreeNode = &StdSymbol{}

// --- Symbol tables ------------------------------------------

/* Symbol tables to store symbols (map-like semantics).
 */
type SymbolTable struct {
	//sync.Mutex
	Table        map[string]Symbol
	createSymbol func(string) Symbol
}

/* Create an empty symbol table.
 */
func NewSymbolTable(symcreator func(string) Symbol) *SymbolTable {
	if symcreator == nil {
		symcreator = NewStdSymbol
	}
	var symtab = SymbolTable{
		//sync.Mutex{},
		Table:        make(map[string]Symbol),
		createSymbol: symcreator,
	}
	return &symtab
}

func (t *SymbolTable) GetSymbolCreator() func(string) Symbol {
	return t.createSymbol
}

/* Check for a symbol in the symbol table.
 * Returns a pointer to a symbol or nil.
 */
func (t *SymbolTable) ResolveSymbol(symname string) Symbol {
	//t.Lock()
	sym := t.Table[symname] // get symbol by name
	//t.Unlock()
	return sym
}

/* Find a symbol in the table, insert a new one if not found.
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

/* Create a new symbol to store into the symbol table.
 * The symbol's name may not be empty
 * Overwrites existing symbol with this name, if any.
 * Returns the new symbol and the previously stored symbol (or nil).
 */
func (t *SymbolTable) DefineSymbol(symname string) (Symbol, Symbol) {
	if len(symname) == 0 {
		return nil, nil
	}
	//T.Debug("inserting 'undefined' symbol: \"%s\"\n", symname)
	sym := t.createSymbol(symname)
	old := t.InsertSymbol(sym)
	return sym, old
}

func (t *SymbolTable) InsertSymbol(sym Symbol) Symbol {
	old := t.ResolveSymbol(sym.GetName())
	//t.Lock()
	t.Table[sym.GetName()] = sym
	//t.Unlock()
	return old
}

func (t *SymbolTable) Each(mapper func(string, Symbol)) {
	for k, v := range t.Table {
		mapper(k, v)
	}
}

// --- Scopes -----------------------------------------------------------

type Scope struct {
	Name   string
	Parent *Scope
	symtab *SymbolTable
}

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

/* Find symbol. Returns the symbol (or nil) and a scope. The scope is
 * the scope of a scope-tree-path the symbol was found in.
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
	T.P("scope", newsc.Name).Debug("pushing new scope")
	return newsc
}

/* Pop a scope.
 */
func (scst *ScopeTree) PopScope() *Scope {
	if scst.ScopeTOS == nil {
		panic("attempt to pop scope from empty stack")
	}
	sc := scst.ScopeTOS
	T.Debugf("popping scope [%s]", sc.Name)
	scst.ScopeTOS = scst.ScopeTOS.Parent
	return sc
}

// ====== obsolete =======
/* Return the current scope.
func (l *ParseListener) scope() *syntax.Scope {
	if l.scopeTOS == nil {
		panic("attempt to access scope from empty stack")
	}
	return l.scopeTOS
}

/* Get the outermost scope, containing global symbols.
func (l *ParseListener) globals() *syntax.Scope {
	if l.scopeStack == nil {
		panic("attempt to access global scope from empty stack")
	}
	return l.scopeStack
}

/* Push a scope. A scope is constructed, including a symbol table
 * for variable declarations.
func (l *ParseListener) pushNewScope(nm string) *syntax.Scope {
	scp := l.scopeTOS
	newsc := syntax.NewScope(nm, scp, NewPMMPVarDecl)
	if scp == nil { // the new scope is the global scope
		l.scopeStack = newsc // make new scope anchor
	}
	l.scopeTOS = newsc // new scope now TOS
	T.Debug("@@@ pushing new scope [%s]\n", newsc.Name)
	return newsc
}

/* Pop a scope.
func (l *ParseListener) popScope() *syntax.Scope {
	if l.scopeTOS == nil {
		panic("attempt to pop scope from empty stack")
	}
	sc := l.scopeTOS
	T.Debug("@@@ popping scope [%s]\n", sc.Name)
	l.scopeTOS = l.scopeTOS.Parent
	return sc
}
*/
