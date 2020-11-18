package lr

import (
	"fmt"
	"text/scanner"
)

// A GrammarBuilder is used to construct a Grammar.
//
//     b := NewGrammarBuilder("G")
//     b.LHS("S").N("A").T("a", 1).End()  // S  ->  A a
//     b.LHS("A").N("B").N("D").End()     // A  ->  B D
//     b.LHS("B").T("b", 2).End()         // B  ->  b
//     b.LHS("B").Epsilon()               // B  ->
//     b.LHS("D").T("d", 3).End()         // D  ->  d
//     b.LHS("D").Epsilon()               // D  ->
//
// This results in the following grammar:  b.Grammar().Dump()
//
//   0: [S']::= [S]
//   1: [S] ::= [A a]
//   2: [A] ::= [B D]
//   3: [B] ::= [b]
//   4: [B] ::= []
//   5: [D] ::= [d]
//   6: [D] ::= []
//
// A call to b.Grammar() returns the (completed) grammar.
//
type GrammarBuilder struct {
	g                  *Grammar      // the grammar to build
	initial            *Rule         // the top-level rule we will wrap around the user's first rule
	tokenizerHook      TokenizerHook // tokenizer hook in the resuting grammar
	tokenValueSequence int           // internal sequence for terminal token values
}

// NewGrammarBuilder gets a new grammar builder, given the name of the grammar to build.
func NewGrammarBuilder(gname string) *GrammarBuilder {
	g := newLRGrammar(gname)
	gb := &GrammarBuilder{g: g, initial: newRule()}
	sym := g.resolveOrDefineNonTerminal("S'")
	gb.initial.LHS = sym                        // LHS of wrapper rule S' -> S #eof
	gb.g.rules = append(gb.g.rules, gb.initial) // RHS to be added later
	gb.tokenValueSequence = NonTermType
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
	r.Serial = rno
	gb.g.rules = append(gb.g.rules, r)
}

// LHS starts a rule given the left hand side symbol (non-terminal).
func (gb *GrammarBuilder) LHS(s string) *RuleBuilder {
	rb := gb.newRuleBuilder()
	sym := rb.gb.g.resolveOrDefineNonTerminal(s)
	rb.rule.LHS = sym
	return rb
}

// Grammar returns the (completed) grammar.
func (gb *GrammarBuilder) Grammar() (*Grammar, error) {
	if len(gb.g.rules) <= 1 {
		T().Errorf("Grammar does not contain any rules")
		return nil, fmt.Errorf("Grammar does not contain any rules")
	}
	gb.initial.rhs = append(gb.initial.rhs, gb.g.rules[1].LHS)
	eof := gb.g.resolveOrDefineTerminal("#eof", scanner.EOF)
	gb.initial.rhs = append(gb.initial.rhs, eof)
	return gb.g, nil
}

// SetTokenizerHook sets a tokenizer hook, which will be called by the grammar
// to produce terminal tokens.
func (gb *GrammarBuilder) SetTokenizerHook(hook TokenizerHook) {
	gb.tokenizerHook = hook
}

// TokenizerHook is an interface for a hook function, which will produce properties
// for a valid terminal token for this grammar.
type TokenizerHook interface {
	NewToken(string) (string, int)
}

// RuleBuilder is a builder type for rules.
type RuleBuilder struct {
	gb   *GrammarBuilder
	rule *Rule
}

// N appends a non-terminal to the builder.
// The internal symbol created for the non-terminal will have an ID
// less than -1000.
func (rb *RuleBuilder) N(s string) *RuleBuilder {
	sym := rb.gb.g.resolveOrDefineNonTerminal(s)
	rb.rule.rhs = append(rb.rule.rhs, sym)
	return rb
}

// T appends a terminal to the builder.
// The symbol created for the terminal must not have a token value
// <= -1000 and not have value 0 or -1.
// This is due to the convention of the stdlib-package text/parser, which
// uses token values > 0 for single-rune tokens and token values < 0 for
// common language elements like identifiers, strings, numbers, etc.
// (it is assumed that no symbol set will require more than 1000 of such
// language elements). The method call will panic if this restriction is
// violated.
func (rb *RuleBuilder) T(s string, tokval int) *RuleBuilder {
	if tokval <= NonTermType {
		T().Errorf("illegal token value parameter (%d), must be > %d", tokval, NonTermType)
		panic(fmt.Sprintf("illegal token value parameter (%d)", tokval))
	}
	sym := rb.gb.g.resolveOrDefineTerminal(s, tokval)
	lrs := sym
	rb.rule.rhs = append(rb.rule.rhs, lrs)
	return rb
}

// L appends a terminal/lexeme to the builder.
// This will create a symbol for a terminal, with a token value > -1000.
// This is due to the convention of the stdlib-package text/parser, which
// uses token values > 0 for single-rune tokens and token values < 0 for
// common language elements like identifiers, strings, numbers, etc.
// (it is assumed that no symbol set will require more than 1000 of such
// language elements). The method call will panic if this restriction is
// violated.
//
// The token value will either be generated from an internal sequence, or –
// if a tokenizer-hook is set – by the hook.
func (rb *RuleBuilder) L(s string) *RuleBuilder {
	tokval := 0
	if rb.gb.tokenizerHook != nil {
		s, tokval = rb.gb.tokenizerHook.NewToken(s)
	} else {
		rb.gb.tokenValueSequence++
		tokval = rb.gb.tokenValueSequence
	}
	sym := rb.gb.g.resolveOrDefineTerminal(s, tokval)
	lrs := sym
	rb.rule.rhs = append(rb.rule.rhs, lrs)
	return rb
}

// AppendSymbol appends your own symbol objects to the builder to extend the RHS of a rule.
// Clients will have to make sure no different 2 symbols have the same ID
// and no symbol ID equals a token value of a non-terminal. This restriction
// is necessary to help produce correct GOTO tables for LR-parsing.
func (rb *RuleBuilder) AppendSymbol(sym *Symbol) *RuleBuilder {
	rb.rule.rhs = append(rb.rule.rhs, sym)
	return rb
}

// Epsilon sets epsilon as the RHS of a production.
// This must be called directly after rb.LHS(...).
// It closes the rule, thus no call to End() or EOF() must follow.
func (rb *RuleBuilder) Epsilon() *Rule {
	rb.gb.appendRule(rb.rule)
	T().Debugf("appending epsilon-rule:  %v", rb.rule)
	r := rb.rule
	rb.rule = nil
	return r
}

// EOF appends EOF as a (terminal) symbol to a rule.
// This is usually not called by clients, but rather internally by the grammar
// builder. If you know what you're doing, be careful.
//
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
