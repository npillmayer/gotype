package lr

import (
	"fmt"
	"text/scanner"

	"github.com/npillmayer/gotype/syntax/runtime"
)

// http://scottmcpeak.com/elkhound/

// --- Rules -----------------------------------------------------------------

// REMARK: the left hand side of rules is an array, although they consist
// of just one symbol. This is due to the fact that I intended to have
// context sensitive rules, too (not for LR-parsing of course).

// Symbol is an interface for symbols of grammars, i.e. terminals and non-terminals.
type Symbol interface {
	IsTerminal() bool // is this symbol a terminal?
	Token() int       // token value of terminals
	GetID() int       // ID of non-terminals (and possibly terminals)
}

func symvalue(A Symbol) int {
	if A.IsTerminal() {
		return A.Token()
	}
	return A.GetID()
}

// A Rule is a type for rules of a grammar. Rules cannot be shared between grammars.
type Rule struct {
	no    int           // order number of this rule within a grammar
	lhs   []Symbol      // symbols of left hand side
	rhs   []Symbol      // symbols of right hand side
	items map[int]*item // Earley items for this rule, int = dot position
}

func newRule() *Rule {
	r := &Rule{}
	r.lhs = make([]Symbol, 0, 5)
	r.rhs = make([]Symbol, 0, 5)
	r.items = make(map[int]*item)
	return r
}

/*
Get the right hand side of a rule as a shallow copy. Clients should treat
it as read-only.
*/
func (r *Rule) GetRHS() []Symbol {
	dup := make([]Symbol, len(r.rhs))
	copy(dup, r.rhs)
	return dup
}

/*
Get the LHS symbol of a rule.
*/
func (r *Rule) GetLHSSymbol() Symbol {
	return r.lhs[0]
}

// Check if two RHS are identical. If parameter prefix is true, this function
// returns true when handle is a prefix of r, even if handle is of length 0.
func (r *Rule) eqRHS(handle []Symbol, prefix bool) bool {
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

// Return an item for a rule wth the dot at the start of RHS
func (r *Rule) startItem() (*item, Symbol) {
	var i *item
	if i = r.items[0]; i == nil {
		i = &item{rule: r, dot: 0}
		r.items[0] = i // store this item for dot-position = 0
	}
	if r.isEps() {
		return i, nil
	} else {
		return i, r.rhs[0]
	}
}

func (r *Rule) findOrCreateItem(dot int) (*item, Symbol) {
	if dot > len(r.rhs) {
		return nil, nil
	}
	var i *item
	if i = r.items[dot]; i == nil { // not found => create item for dot position
		i = &item{rule: r, dot: dot}
		r.items[dot] = i // remember the new item
	}
	return i, i.peekSymbol()
}

// Debugging helper: string representation of a grammar rule
func (r *Rule) String() string {
	s := fmt.Sprintf("%v ::= %v", r.lhs, r.rhs)
	return s
}

func (r *Rule) isEps() bool {
	return len(r.rhs) == 0
}

type item struct {
	rule *Rule
	dot  int
}

func (i *item) String() string {
	s := fmt.Sprintf("%v ::= %v @ %v", i.rule.lhs, i.rule.rhs[0:i.dot], i.rule.rhs[i.dot:])
	return s
}

func (i *item) peekSymbol() Symbol {
	if i.dot >= len(i.rule.rhs) {
		return nil
	}
	return i.rule.rhs[i.dot]
}

// returns a slice, so the result should probably be considered read-only.
func (i *item) getPrefix() []Symbol {
	return i.rule.rhs[:i.dot]
}

func (i *item) advance() (*item, Symbol) {
	if i.dot >= len(i.rule.rhs) {
		return nil, nil
	} else {
		ii, A := i.rule.findOrCreateItem(i.dot + 1)
		return ii, A
	}
}

// --- Grammar ---------------------------------------------------------------

// Type for a grammar. Usually created using a GrammarBuilder.
type Grammar struct {
	Name             string         // a grammar has a name, for documentation only
	rules            []*Rule        // grammar productions, first one is start rule
	epsilon          Symbol         // a special symbol representing epsilon
	nonterminals     map[int]Symbol // all non-terminals
	terminals        map[int]Symbol // all terminals
	terminalsByToken map[int]Symbol // terminals, indexed by token value
}

func newLRGrammar(gname string) *Grammar {
	g := &Grammar{}
	g.Name = gname
	g.rules = make([]*Rule, 0, 30)
	g.terminals = make(map[int]Symbol)
	g.nonterminals = make(map[int]Symbol)
	g.terminalsByToken = make(map[int]Symbol)
	eps := newLRSymbol("_eps").(*lrSymbol)
	eps.SetType(epsilonType)
	eps.Id = 0
	g.epsilon = eps
	return g
}

// Get the symbol representing epsilon for this grammar. Treat as read-only!
func (g *Grammar) Epsilon() Symbol {
	return g.epsilon
}

// Get the symbol representing epsilon for this grammar. Treat as read-only!
func (g *Grammar) Rule(no int) *Rule {
	if no < 0 || no >= len(g.rules) {
		return nil
	}
	return g.rules[no]
}

func (g *Grammar) GetTerminalSymbolFor(tokenvalue int) Symbol {
	return g.terminalsByToken[tokenvalue]
}

func (g *Grammar) findNonTermRules(sym Symbol, includeEpsRules bool) *itemSet {
	iset := newItemSet()
	for _, r := range g.rules {
		if r.lhs[0] == sym {
			if !r.isEps() || includeEpsRules {
				i, _ := r.startItem()
				if i == nil {
					T().Errorf("inconsistency? start-item == NIL")
				}
				iset.Add(i)
			}
		}
	}
	return iset
}

func (g *Grammar) matchesRHS(handle []Symbol, prefix bool) (*Rule, int) {
	for i, r := range g.rules {
		if r.eqRHS(handle, prefix) {
			return r, i
		}
	}
	return nil, -1
}

/*
Iterate over all non-terminal symbols of the grammar.
Return values of the mapper function for all non-terminals are returned as an
array.
*/
func (g *Grammar) EachNonTerminal(mapper func(sym Symbol) interface{}) []interface{} {
	var r []interface{} = make([]interface{}, 0, len(g.nonterminals))
	for _, A := range g.nonterminals {
		r = append(r, mapper(A))
	}
	return r
}

/*
Iterate over all terminals of the grammar.
Return values of the mapper function for all terminals are returned as an array.
*/
func (g *Grammar) EachTerminal(mapper func(sym Symbol) interface{}) []interface{} {
	var r []interface{} = make([]interface{}, 0, len(g.terminals))
	for _, B := range g.terminals {
		r = append(r, mapper(B))
	}
	return r
}

/*
Iterate over all symbols of the grammar.
Return values of the mapper function are returned as an array.
*/
func (g *Grammar) EachSymbol(mapper func(sym Symbol) interface{}) []interface{} {
	var r []interface{} = make([]interface{}, 0, len(g.terminals)+len(g.nonterminals))
	for _, A := range g.nonterminals {
		r = append(r, mapper(A))
	}
	for _, B := range g.terminals {
		r = append(r, mapper(B))
	}
	return r
}

func (g *Grammar) resolveOrDefineNonTerminal(s string) Symbol {
	for _, nt := range g.nonterminals {
		if lrsym, ok := nt.(*lrSymbol); ok {
			if lrsym.GetName() == s {
				return lrsym
			}
		}
	}
	lrsym := newLRSymbol(s)
	g.nonterminals[lrsym.GetID()] = lrsym
	return lrsym
}

func (g *Grammar) resolveOrDefineTerminal(s string, tokval int) Symbol {
	for _, nt := range g.terminals {
		if lrsym, ok := nt.(*lrSymbol); ok {
			if lrsym.GetName() == s {
				return lrsym
			}
		}
	}
	lrsym := newLRSymbol(s).(*lrSymbol)
	lrsym.Symtype = terminalType
	lrsym.tokval = tokval
	g.terminals[lrsym.GetID()] = lrsym
	if g.terminalsByToken[tokval] != nil {
		T().Errorf("duplicate terminal symbol for token value = %d", tokval)
		// proceed with fingers crossed
	}
	g.terminalsByToken[tokval] = lrsym
	return lrsym
}

// Debugging helper: dump symbols and rules to stdout.
func (g *Grammar) Dump() {
	fmt.Printf("--- %s --------------------------------------------\n", g.Name)
	fmt.Printf("epsilon  = %d\n", g.epsilon.GetID())
	for _, A := range g.nonterminals {
		fmt.Printf("N  %v = %d\n", A, A.GetID())
	}
	for _, A := range g.terminals {
		fmt.Printf("T  %v = %d\n", A, A.Token())
	}
	for _, r := range g.rules {
		fmt.Printf("%3d: %s\n", r.no, r.String())
	}
	fmt.Println("-------------------------------------------------------")
}

// ---------------------------------------------------------------------------

const (
	epsilonType  = 0
	terminalType = 1
	nonTermType  = -1000 // IDs of non-terminals MUST be below this !
)

// Serial no. for lrSymbol IDs
var lrSymbolIDSerial int = nonTermType - 1

// An internal symbol type used by the grammar builder.
// Clients may supply their own symbol type, but should be able to
// rely on a standard implementation, if they do not bring one themselves.
type lrSymbol struct {
	*runtime.StdSymbol     // TODO: create our own (do not couple to StdSymbol)
	tokval             int // for terminals
}

func (lrsym *lrSymbol) String() string {
	return lrsym.GetName()
}

func (lrsym *lrSymbol) IsTerminal() bool {
	return lrsym.GetType() == terminalType
}

func (lrsym *lrSymbol) Token() int {
	return lrsym.tokval
}

func newLRSymbol(s string) Symbol {
	stdsym := runtime.NewStdSymbol(s)
	sym := &lrSymbol{stdsym.(*runtime.StdSymbol), 0}
	sym.Symtype = nonTermType
	sym.Id = lrSymbolIDSerial // -1001, -1002, -1003, ...
	lrSymbolIDSerial--
	return sym
}

// === Builder ===============================================================

/*
Grammars are constructed using a GrammarBuilder.

    b := NewGrammarBuilder("G")
    b.LHS("S").N("A").T("a", 1).EOF()  // S  ->  A a EOF
    b.LHS("A").N("B").N("D").End()     // A  ->  B D
    b.LHS("B").T("b", 2).End()         // B  ->  b
    b.LHS("B").Epsilon()               // B  ->
    b.LHS("D").T("d", 3).End()         // D  ->  d
    b.LHS("D").Epsilon()               // D  ->

This results in the following grammar:  b.Grammar().Dump()

  0: [S] ::= [A a #eof]
  1: [A] ::= [B D]
  2: [B] ::= [b]
  3: [B] ::= []
  4: [D] ::= [d]
  5: [D] ::= []

A call to b.Grammar() returns the (completed) grammar.
*/
type GrammarBuilder struct {
	g *Grammar
}

// Get a new grammar builder, given the name of the grammar to build.
func NewGrammarBuilder(gname string) *GrammarBuilder {
	g := newLRGrammar(gname)
	gb := &GrammarBuilder{g}
	return gb
}

func (gb *GrammarBuilder) newRuleBuilder() *RuleBuilder {
	rb := &RuleBuilder{}
	rb.gb = gb
	rb.rule = newRule()
	return rb
}

func (gb *GrammarBuilder) appendRule(r *Rule) {
	rno := len(gb.g.rules)
	r.no = rno
	gb.g.rules = append(gb.g.rules, r)
}

// Start a rule given the left hand side symbol (non-terminal).
func (gb *GrammarBuilder) LHS(s string) *RuleBuilder {
	rb := gb.newRuleBuilder()
	sym := rb.gb.g.resolveOrDefineNonTerminal(s)
	rb.rule.lhs = append(rb.rule.lhs, sym)
	return rb
}

// Return the (completed) grammar.
func (gb *GrammarBuilder) Grammar() *Grammar {
	return gb.g
}

// A builder type for rules
type RuleBuilder struct {
	gb   *GrammarBuilder
	rule *Rule
}

// Append a non-terminal to the builder.
// The internal symbol created for the non-terminal will have an ID
// less than -1000.
func (rb *RuleBuilder) N(s string) *RuleBuilder {
	sym := rb.gb.g.resolveOrDefineNonTerminal(s)
	rb.rule.rhs = append(rb.rule.rhs, sym)
	return rb
}

/*
Append a terminal to the builder.
The symbol created for the terminal must not have a token value
<= -1000 and not have value 0 or -1.
This is due to the convention of the stdlib-package text/parser, which
uses token values > 0 for single-rune tokens and token values < 0 for
common language elements like identifiers, strings, numbers, etc.
(it is assumed that no symbol set will require more than 1000 of such
language elements). The method call will panic if this restriction is
violated.
*/
func (rb *RuleBuilder) T(s string, tokval int) *RuleBuilder {
	if tokval <= nonTermType {
		T().Errorf("illegal token value parameter (%d), must be > %d", tokval, nonTermType)
		panic(fmt.Sprintf("illegal token value parameter (%d)", tokval))
	}
	sym := rb.gb.g.resolveOrDefineTerminal(s, tokval)
	lrs := sym.(*lrSymbol)
	rb.rule.rhs = append(rb.rule.rhs, lrs)
	return rb
}

// Append your own symbol objects to the builder to extend the RHS of a rule.
// Clients will have to make sure no different 2 symbols have the same ID
// and no symbol ID equals a token value of a non-terminal. This restriction
// is necessary to help produce correct GOTO tables for LR-parsing.
func (rb *RuleBuilder) AppendSymbol(sym Symbol) *RuleBuilder {
	rb.rule.rhs = append(rb.rule.rhs, sym)
	return rb
}

// Set epsilon as the RHS of a production.
// This must be called directly after rb.LHS(...).
func (rb *RuleBuilder) Epsilon() *Rule {
	rb.gb.appendRule(rb.rule)
	T().Debugf("appending epsilon-rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}

// Append EOF as a (terminal) symbol to a rule.
// This completes the rule (no other builder calls should be made
// for this rule).
func (rb *RuleBuilder) EOF() *Rule {
	rb.T("#eof", scanner.EOF)
	return rb.End()
}

// End a rule.
// This completes the rule (no other builder calls should be made
// for this rule).
func (rb *RuleBuilder) End() *Rule {
	rb.gb.appendRule(rb.rule)
	T().Debugf("appending rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}
