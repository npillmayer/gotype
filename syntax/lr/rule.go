package lr

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/runtime"
)

// --- Rules -----------------------------------------------------------------

// An interface for symbols of grammars, i.e. terminals and non-terminals.
type Symbol interface {
	IsTerminal() bool // is this symbol a terminal?
	Token() int       // token value of terminals
	GetID() int       // ID of non-terminals (and possibly terminals)
}

// A type for rules of a grammar
type Rule struct {
	lhs   []Symbol      // symbols of left hand side
	rhs   []Symbol      // symbols of right hand side
	items map[int]*item // Earley items for this rule
}

func newRule() *Rule {
	r := &Rule{}
	r.lhs = make([]Symbol, 0, 5)
	r.rhs = make([]Symbol, 0, 5)
	r.items = make(map[int]*item)
	return r
}

// Return an item for a rule wth the dot at the start of RHS
func (r *Rule) startItem() (*item, Symbol) {
	if r.isEps() {
		return nil, nil
	} else {
		var i *item
		if i = r.items[0]; i == nil {
			i = &item{rule: r, dot: 0}
			r.items[0] = i // store this item
		}
		return i, r.rhs[0]
	}
}

func (r *Rule) findOrCreateItem(dot int) (*item, Symbol) {
	if dot > len(r.rhs) {
		return nil, nil
	}
	var i *item
	if i = r.items[dot]; i == nil {
		i = &item{rule: r, dot: dot}
		r.items[dot] = i
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
	name    string
	rules   []*Rule
	symbols *runtime.SymbolTable
	epsilon Symbol
}

func newLRGrammar(gname string) *Grammar {
	g := &Grammar{}
	g.name = gname
	g.rules = make([]*Rule, 0, 30)
	g.symbols = runtime.NewSymbolTable(newLRSymbol)
	eps := newLRSymbol("_eps").(*lrSymbol)
	eps.SetType(epsilonType)
	g.epsilon = eps
	return g
}

func (g *Grammar) findNonTermRules(sym Symbol) *itemSet {
	iset := newItemSet()
	for _, r := range g.rules {
		if r.lhs[0] == sym {
			i, _ := r.startItem()
			iset.Add(i)
		}
	}
	return iset
}

// Debugging helper: dump symbols and rules to stdout
func (g *Grammar) Dump() {
	fmt.Printf("--- %s --------------------------------------------\n", g.name)
	fmt.Printf("epsilon  = %d\n", g.epsilon.GetID())
	g.symbols.Each(func(name string, sym runtime.Symbol) {
		A := sym.(Symbol)
		if A.IsTerminal() {
			fmt.Printf("T  %5s = %d\n", name, A.Token())
		} else {
			fmt.Printf("N  %5s = %d\n", name, A.GetID())
		}
	})
	for k, r := range g.rules {
		fmt.Printf("%3d: %s\n", k, r.String())
	}
	fmt.Println("-------------------------------------------------------")
}

// ---------------------------------------------------------------------------

const (
	NonTermType  = 100
	TerminalType = 1000
	epsilonType  = 999
)

type lrSymbol struct {
	*runtime.StdSymbol
}

func (lrsym *lrSymbol) String() string {
	return lrsym.GetName()
}

func (lrsym *lrSymbol) IsTerminal() bool {
	return lrsym.GetType() >= epsilonType
}

func (lrsym *lrSymbol) Token() int {
	return lrsym.GetType()
}

func newLRSymbol(s string) runtime.Symbol {
	sym := runtime.NewStdSymbol(s)
	st := sym.(*runtime.StdSymbol)
	st.Symtype = NonTermType
	return &lrSymbol{st}
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

This results in the following grammar:  b.Grammar().Dump() :

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
	gb.g.rules = append(gb.g.rules, r)
}

// Start a rule given the left hand side symbol (non-terminal).
func (gb *GrammarBuilder) LHS(s string) *RuleBuilder {
	rb := gb.newRuleBuilder()
	sym, _ := rb.gb.g.symbols.ResolveOrDefineSymbol(s)
	lrs := sym.(*lrSymbol)
	rb.rule.lhs = append(rb.rule.lhs, lrs)
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
func (rb *RuleBuilder) N(s string) *RuleBuilder {
	sym, _ := rb.gb.g.symbols.ResolveOrDefineSymbol(s)
	lrs := sym.(*lrSymbol)
	rb.rule.rhs = append(rb.rule.rhs, lrs)
	return rb
}

// Append a terminal to the builder.
func (rb *RuleBuilder) T(s string, tokval int) *RuleBuilder {
	sym, _ := rb.gb.g.symbols.ResolveOrDefineSymbol(s)
	lrs := sym.(*lrSymbol)
	lrs.Symtype = TerminalType + tokval
	rb.rule.rhs = append(rb.rule.rhs, lrs)
	return rb
}

// Set epsilon as the RHS of a production.
// This must be called directly after rb.LHS(...).
func (rb *RuleBuilder) Epsilon() *Rule {
	rb.gb.appendRule(rb.rule)
	T.Debugf("appending epsilon-rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}

// Append EOF as a (terminal) symbol to a rule.
// This completes the rule (no other builder calls should be made
// for this rule).
func (rb *RuleBuilder) EOF() *Rule {
	rb.T("#eof", 0)
	return rb.End()
}

// End a rule.
// This completes the rule (no other builder calls should be made
// for this rule).
func (rb *RuleBuilder) End() *Rule {
	rb.gb.appendRule(rb.rule)
	T.Debugf("appending rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}
