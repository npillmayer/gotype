package lr

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/lr/iteratable"
)

// --- Symbol ----------------------------------------------------------------

// Symbols' values must adhere to these ranges.
const (
	EpsilonType = 0
	EOFType     = -1    // pseudo terminal token for end of input
	NonTermType = -1000 // IDs of non-terminals MUST be in { -2 … -999 }
)

// Serial no. for lrSymbol IDs
var lrSymbolIDSerial = NonTermType - 1

// Symbol is a symbol type used for grammars and grammar builders.
type Symbol struct {
	Name  string // visual representation, if any
	Value int    // ID or token value
}

func (lrsym *Symbol) String() string {
	//return fmt.Sprintf("<%s|%d>", lrsym.Name, lrsym.Value)
	return fmt.Sprintf("%s", lrsym.Name)
}

// IsTerminal returns true if this symbol represents a terminal.
func (lrsym *Symbol) IsTerminal() bool {
	return lrsym.Value > NonTermType
}

// Token is just an alias for lrsym.Value.
func (lrsym *Symbol) Token() int {
	return lrsym.Value
}

func newSymbol(s string) *Symbol {
	lrSymbolIDSerial--
	return &Symbol{
		Name:  s,
		Value: lrSymbolIDSerial,
	}
}

// --- Rules -----------------------------------------------------------------

// A Rule is a type for rules of a grammar. Rules cannot be shared between grammars.
type Rule struct {
	Serial int       // order number of this rule within a grammar
	LHS    *Symbol   // symbols of left hand side
	rhs    []*Symbol // symbols of right hand side
}

func newRule() *Rule {
	r := &Rule{}
	r.rhs = make([]*Symbol, 0, 5)
	return r
}

func (r *Rule) String() string {
	return fmt.Sprintf("%v ➞ %v", r.LHS, r.rhs)
}

// RHS gets the right hand side of a rule as a shallow copy. Clients should treat
// it as read-only.
func (r *Rule) RHS() []*Symbol {
	dup := make([]*Symbol, len(r.rhs))
	copy(dup, r.rhs)
	return dup
}

// IsEps returns true if this an epsilon-rule.
func (r *Rule) IsEps() bool {
	return len(r.rhs) == 0
}

// Check if two RHS are identical. If parameter prefix is true, this function
// returns true when handle is a prefix of r, even if handle is of length 0.
func (r *Rule) eqRHS(handle []*Symbol, prefix bool) bool {
	if r.rhs == nil && handle == nil {
		return true
	}
	if r.rhs == nil || handle == nil {
		return false
	}
	if !prefix && (len(r.rhs) != len(handle)) {
		return false
	}
	for i := range r.rhs {
		if i >= len(handle) {
			return true
		}
		if r.rhs[i] != handle[i] {
			return false
		}
	}
	return true
}

// --- Grammar ---------------------------------------------------------------

// Grammar is a type for a grammar. Usually created using a GrammarBuilder.
type Grammar struct {
	Name         string          // a grammar has a name, for documentation only
	rules        []*Rule         // grammar productions, first one is start rule
	Epsilon      *Symbol         // a special symbol representing epsilon
	EOF          *Symbol         // a special symbol representing end of input
	nonterminals map[int]*Symbol // all non-terminals
	terminals    map[int]*Symbol // all terminals
	//terminalsByToken map[int]*Symbol // terminals, indexed by token value
}

func newLRGrammar(gname string) *Grammar {
	g := &Grammar{}
	g.Name = gname
	g.rules = make([]*Rule, 0, 30)
	g.terminals = make(map[int]*Symbol)
	g.nonterminals = make(map[int]*Symbol)
	//g.terminalsByToken = make(map[int]Symbol)
	g.Epsilon = &Symbol{Name: "_eps", Value: EpsilonType}
	g.EOF = &Symbol{Name: "_eof", Value: EOFType}
	return g
}

// Size returns the number of rules in the grammar.
func (g *Grammar) Size() int {
	return len(g.rules)
}

// Rule gets a grammar rule.
func (g *Grammar) Rule(no int) *Rule {
	if no < 0 || no >= len(g.rules) {
		return nil
	}
	return g.rules[no]
}

// Terminal returns the terminal symbol for a given token value, if it
// is defined in the grammar.
func (g *Grammar) Terminal(tokval int) *Symbol {
	if t, ok := g.terminals[tokval]; ok {
		return t
	}
	return nil
}

// SymbolByName gets a symbol for a given name, if found in the grammar.
func (g *Grammar) SymbolByName(name string) *Symbol {
	var found *Symbol
	g.EachSymbol(func(sym *Symbol) interface{} {
		if sym.Name == name {
			found = sym
		}
		return nil
	})
	return found
}

// FindNonTermRules returns a set of Earley-items, where each item stems from
// a rule with a given LHS and the dot is at position 0.
func (g *Grammar) FindNonTermRules(sym *Symbol, includeEpsRules bool) *iteratable.Set {
	iset := iteratable.NewSet(0)
	for _, r := range g.rules {
		if r.LHS == sym {
			if !r.IsEps() || includeEpsRules {
				item, _ := StartItem(r)
				iset.Add(item)
			}
		}
	}
	return iset
}

func (g *Grammar) matchesRHS(lhs *Symbol, handle []*Symbol) (*Rule, int) {
	//func (g *Grammar) matchesRHS(lhs Symbol, handle []Symbol, prefix bool) (*Rule, int) {
	for i, r := range g.rules {
		if r.LHS == lhs {
			if r.eqRHS(handle, false) {
				return r, i
			}
		}
	}
	return nil, -1
}

// EachNonTerminal iterates over all non-terminal symbols of the grammar.
// Return values of the mapper function for all non-terminals are returned as an
// array.
func (g *Grammar) EachNonTerminal(mapper func(sym *Symbol) interface{}) []interface{} {
	var r = make([]interface{}, 0, len(g.nonterminals))
	for _, A := range g.nonterminals {
		r = append(r, mapper(A))
	}
	return r
}

// EachTerminal iterates over all terminals of the grammar.
// Return values of the mapper function for all terminals are returned as an array.
func (g *Grammar) EachTerminal(mapper func(sym *Symbol) interface{}) []interface{} {
	var r = make([]interface{}, 0, len(g.terminals))
	for _, B := range g.terminals {
		r = append(r, mapper(B))
	}
	return r
}

// EachSymbol iterates over all symbols of the grammar.
// Return values of the mapper function are returned as an array.
func (g *Grammar) EachSymbol(mapper func(sym *Symbol) interface{}) []interface{} {
	var r = make([]interface{}, 0, len(g.terminals)+len(g.nonterminals))
	for _, A := range g.nonterminals {
		r = append(r, mapper(A))
	}
	for _, B := range g.terminals {
		r = append(r, mapper(B))
	}
	return r
}

func (g *Grammar) resolveOrDefineNonTerminal(s string) *Symbol {
	for _, nt := range g.nonterminals {
		if nt.Name == s {
			return nt
		}
	}
	lrsym := newSymbol(s)
	g.nonterminals[lrsym.Value] = lrsym
	return lrsym
}

func (g *Grammar) resolveOrDefineTerminal(s string, tokval int) *Symbol {
	for _, t := range g.terminals {
		if t.Name == s {
			return t
		}
	}
	t := &Symbol{Name: s, Value: tokval}
	g.terminals[tokval] = t
	// if g.terminalsByToken[tokval] != nil {
	// 	T().Errorf("duplicate terminal symbol for token value = %d", tokval)
	// 	// proceed with fingers crossed
	// }
	//g.terminalsByToken[tokval] = lrsym
	return t
}

// Dump is a debugging helper: dump symbols and rules to stdout.
func (g *Grammar) Dump() {
	T().Debugf("--- %s --------------------------------------------\n", g.Name)
	T().Debugf("epsilon  = %d\n", g.Epsilon.Value)
	T().Debugf("EOF      = %d\n", g.EOF.Value)
	for _, A := range g.nonterminals {
		T().Debugf("N  %v = %d\n", A, A.Value)
	}
	for _, A := range g.terminals {
		T().Debugf("T  %v = %d\n", A, A.Token())
	}
	for _, r := range g.rules {
		T().Debugf("%3d: %s\n", r.Serial, r.String())
	}
	T().Debugf("-------------------------------------------------------")
}
