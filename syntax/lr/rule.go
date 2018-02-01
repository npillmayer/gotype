package lr

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/runtime"
)

type Symbol interface {
	fmt.Stringer
	IsTerminal() bool
}

type Rule struct {
	lhs []Symbol
	rhs []Symbol
}

func newRule() *Rule {
	r := &Rule{}
	r.lhs = make([]Symbol, 0, 5)
	r.rhs = make([]Symbol, 0, 5)
	return r
}

func (r *Rule) String() string {
	s := fmt.Sprintf("%v ::= %v", r.lhs, r.rhs)
	return s
}

type Item struct {
	rule *Rule
	dot  int
}

type Grammar struct {
	name    string
	rules   []*Rule
	symbols *runtime.SymbolTable
}

// ===========================================================================

const (
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
	return lrsym.GetType() > TerminalType
}

func newLRSymbol(s string) runtime.Symbol {
	sym := runtime.NewStdSymbol(s)
	st := sym.(*runtime.StdSymbol)
	st.Symtype = NonTermType
	return &lrSymbol{st}
}

func NewLRGrammar(gname string) *Grammar {
	g := &Grammar{}
	g.name = gname
	g.rules = make([]*Rule, 0, 30)
	g.symbols = runtime.NewSymbolTable(newLRSymbol)
	return g
}

func (g *Grammar) Builder() *GrammarBuilder {
	gb := &GrammarBuilder{g}
	return gb
}

// === Builder ===============================================================

type GrammarBuilder struct {
	g *Grammar
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
	T.Debugf("appending rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}

func (rb *RuleBuilder) End() *Rule {
	rb.gb.appendRule(rb.rule)
	T.Debugf("appending rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}

// === Table Generation ======================================================

func (g *Grammar) closure(r *lrItem)
