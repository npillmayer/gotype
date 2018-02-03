package lr

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/runtime"
)

// --- Rules -----------------------------------------------------------------

type Symbol interface {
	//fmt.Stringer
	IsTerminal() bool
	Token() int
	GetID() int
}

type Rule struct {
	lhs   []Symbol
	rhs   []Symbol
	items map[int]*item
}

func newRule() *Rule {
	r := &Rule{}
	r.lhs = make([]Symbol, 0, 5)
	r.rhs = make([]Symbol, 0, 5)
	r.items = make(map[int]*item)
	return r
}

func (r *Rule) startItem() (*item, Symbol) {
	if r.isEps() {
		return nil, nil
	} else {
		var i *item
		if i = r.items[0]; i == nil {
			i = &item{rule: r, dot: 0}
			r.items[0] = i
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

type Grammar struct {
	name    string
	rules   []*Rule
	symbols *runtime.SymbolTable
	epsilon Symbol
}

func NewLRGrammar(gname string) *Grammar {
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
	epsilonType  = -1
	NonTermType  = 100
	TerminalType = 1000
)

type lrSymbol struct {
	*runtime.StdSymbol
}

func (lrsym *lrSymbol) String() string {
	return lrsym.GetName()
}

func (lrsym *lrSymbol) IsTerminal() bool {
	return lrsym.GetType() >= TerminalType
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

type GrammarBuilder struct {
	g *Grammar
}

func NewGrammarBuilder(gname string) *GrammarBuilder {
	g := NewLRGrammar(gname)
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

func (gb *GrammarBuilder) LHS(s string) *RuleBuilder {
	rb := gb.newRuleBuilder()
	sym, _ := rb.gb.g.symbols.ResolveOrDefineSymbol(s)
	lrs := sym.(*lrSymbol)
	rb.rule.lhs = append(rb.rule.lhs, lrs)
	return rb
}

func (gb *GrammarBuilder) Grammar() *Grammar {
	return gb.g
}

type RuleBuilder struct {
	gb   *GrammarBuilder
	rule *Rule
}

func (rb *RuleBuilder) N(s string) *RuleBuilder {
	sym, _ := rb.gb.g.symbols.ResolveOrDefineSymbol(s)
	lrs := sym.(*lrSymbol)
	rb.rule.rhs = append(rb.rule.rhs, lrs)
	return rb
}

func (rb *RuleBuilder) T(s string, tokval int) *RuleBuilder {
	sym, _ := rb.gb.g.symbols.ResolveOrDefineSymbol(s)
	lrs := sym.(*lrSymbol)
	lrs.Symtype = TerminalType + tokval
	rb.rule.rhs = append(rb.rule.rhs, lrs)
	return rb
}

func (rb *RuleBuilder) Epsilon() *Rule {
	rb.gb.appendRule(rb.rule)
	T.Debugf("appending epsilon-rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}

func (rb *RuleBuilder) EOF() *Rule {
	rb.T("#eof", 0)
	return rb.End()
}

func (rb *RuleBuilder) End() *Rule {
	rb.gb.appendRule(rb.rule)
	T.Debugf("appending rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}
