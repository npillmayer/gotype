package earley

import (
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/iteratable"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
)

// TokenAt returns the input token at position pos.
func (p *Parser) TokenAt(pos uint64) interface{} {
	if pos >= 0 && pos < uint64(len(p.tokens)) {
		return p.tokens[pos+1] // tokens start at index 1
	}
	return nil
}

// --- Derivation listener ---------------------------------------------------

// Listener is a type for walking a parse tree/forest.
type Listener interface {
	Reduce(sym *lr.Symbol, rule int, rhs []*RuleNode, span lr.Span, level int) interface{}
	Terminal(tokenValue int, token interface{}, span lr.Span, level int) interface{}
}

// RuleNode represents a node occuring during a parse tree/forest walk.
type RuleNode struct {
	sym    *lr.Symbol
	Extent lr.Span     // span of intput symbols this rule reduced
	Value  interface{} // user defined value
}

// Symbol returns the grammar symbol a RuleNode refers to.
// It is either a terminal or the LHS of a reduced rule.
func (rnode *RuleNode) Symbol() *lr.Symbol {
	return rnode.sym
}

// --- Tree Walker -----------------------------------------------------------

// WalkDerivation walks the grammar items which occured during the parse.
// It uses a listener, which gets called for every terminal and for every
// non-terminal reduction.
func (p *Parser) WalkDerivation(listener Listener) *RuleNode {
	T().Debugf("=== Walk ===============================")
	var root *RuleNode
	S := p.states[p.sc]
	S.IterateOnce()
	for S.Next() {
		item := S.Item().(lr.Item)
		if item.PeekSymbol() == nil && item.Rule().LHS == p.ga.Grammar().Rule(0).LHS {
			root = p.walk(item, p.sc, listener, 0)
		}
	}
	T().Debugf("========================================")
	T().Debugf("TOKENS: %d", len(p.tokens))
	for i, t := range p.tokens {
		T().Debugf("        [%d]=%v", i, t)
	}
	return root
}

/*
Walk backwards over the items of Earley states.

A good overview of how to construct a parse forest from Earley-items may be found in
"Parsing Techniques" by  Dick Grune and Ceriel J.H. Jacobs
(https://dickgrune.com/Books/PTAPG_2nd_Edition/), Section 7.2.1.2.

Even more practical, a great tutorial by Loup Vaillant
(http://loup-vaillant.fr/tutorials/earley-parsing/parser)
provides a very approachable summary of how to create a parse forest from an Earley-parse.
Here is a relevant excerpt:

Imagine we have an item like this ('a', 'b', and 'c' are symbols, and 'i' is an integer):

    Foo -> a b c •  (i)

The fact that this item even exist means the following items also exist somewhere:

    Foo ->   a   b • c  (i)
    Foo ->   a • b   c  (i)
    Foo -> • a   b   c  (i)

We know the parse was successful: the recogniser said so, by showing us this item:

    === 9 ===
    Sum -> Sum [+-] Product • (0)

There is a dot at the end, so this is a completed item. It starts at (0) (the beginning),
and stops at (9) (the very end). There's only one way Earley's algorithm could possibly
produce such an item: the whole input is a Sum. In our current example this means we can
find those items:

	Sum ->   Sum   [+-] • Product (0)
	Sum ->   Sum • [+-]   Product (0)
	Sum -> • Sum   [+-]   Product (0)

But that's not the end of it. To advance an item one step, you need two things:
an un-advanced version of the item (which we have here), and a completed something:
either a completed state, or a successful scan. This has several implications:

   * There is another completed Sum somewhere. It starts at (0), and
     finishes at… well… let's say (x).
   * There is a successful scan between (x) and (x+1). Meaning, the input at x matches [+-].
   * There is a completed Product somewhere. It starts at (x+1), and finishes at… wait
     a minute this is the last one! it's got to finish wherever the overall Sum finishes!
     That would be the end of the input, or (9).

The problem now is to search for those states, and determine the value of (x).
Given how Earley items are stored in the state sets, we need to start at the end.
*/
func (p *Parser) walk(item lr.Item, pos uint64, listener Listener, level int) *RuleNode {
	rhs := reverse(item.Rule().RHS()) // we iterate over RHS symbols of item
	l := len(rhs)
	T().Debugf("Walk from item=%s (%d…%d)", item, item.Origin, pos)
	extent := lr.Span{item.Origin, pos}
	ruleNodes := make([]*RuleNode, len(rhs)) // we will collect children nodes
	end := pos
	for n, B := range rhs {
		T().Debugf("Next symbol in rev(RHS) is %s", B)
		if B.IsTerminal() { // collect a terminal node
			T().Infof("Tree node    %d: %s", pos-1, B)
			value := listener.Terminal(B.Value, p.tokens[pos], lr.Span{pos - 1, pos}, level+1)
			ruleNodes[l-n-1] = &RuleNode{
				sym:    B,
				Extent: lr.Span{pos - 1, pos},
				Value:  value,
			}
			pos--
			continue
		}
		// for each symbol B, find an item [B→…A•, k] which has completed it
		S := p.states[pos]
		cleanupState(S)
		T().Debugf("Looking for item which completed %s", B)
		dumpState(p.states, pos)
		T().Debugf("---------------------------------------------")
		R := S.Copy().Subset(func(el interface{}) bool {
			jtem := el.(lr.Item)
			return itemCompletes(jtem, B)
		}) // now R contains all items [B→…A•, k]
		T().Debugf("R=%s", itemSetString(R))
		switch R.Size() {
		case 0: // cannot happen
			panic("predecessor for item missing")
		case 1: // non-ambiguous
			child := R.First().(lr.Item)
			ruleNodes[l-n-1] = p.walk(child, pos, listener, level+1)
			pos = child.Origin // k
		default: // ambiguous: resolve by longest rule first, then by rule number
			/* TODO
			Earley parse-tree generation: // ambiguous: resolve by longest rule first,
			then by rule number

			Let clients decide this via option. The implemented default works for many cases,
			but automatically prefers a right-derivation. This may not be what clients want,
			e.g. for left-associative operators.
			*/
			var longest lr.Item
			R.IterateOnce()
			for R.Next() {
				rule := R.Item().(lr.Item)
				// interval(longest) < interval(item) ?
				// avoid looping with parent-rule = child-rule
				if item.Origin <= rule.Origin && !(item.Origin == rule.Origin && pos == end) {
					if longest.Rule() == nil || len(rule.Prefix()) > len(longest.Prefix()) {
						longest = rule
					} else if len(rule.Prefix()) == len(longest.Prefix()) {
						if rule.Origin < longest.Origin {
							longest = rule
						} else if rule.Origin == longest.Origin && rule.Rule().Serial < longest.Rule().Serial {
							longest = rule
						}
					}
				}
			}
			T().Debugf("Selected rule %s", longest)
			ruleNodes[l-n-1] = p.walk(longest, pos, listener, level+1)
			pos = longest.Origin // k
		}
	}
	value := listener.Reduce(item.Rule().LHS, item.Rule().Serial, ruleNodes, extent, level)
	node := &RuleNode{
		sym:    item.Rule().LHS,
		Extent: extent,
		Value:  value,
	}
	T().Infof("Tree node    %d|-----%s-----|%d", extent.From(), item.Rule().LHS.Name, extent.To())
	return node
}

// Does item complete a rule with LHS B ?
func itemCompletes(item lr.Item, B *lr.Symbol) bool {
	return item.PeekSymbol() == nil &&
		item.Rule().LHS.Value == B.Value
}

// Throw away non-completing items, as they are not needed for parse tree construction.
func cleanupState(S *iteratable.Set) {
	S.IterateOnce()
	for S.Next() {
		item := S.Item().(lr.Item)
		if item.PeekSymbol() != nil {
			S.Remove(item)
		}
	}
}

// --- Tree building listener -------------------------------------------

// TreeBuilder is a DerivationListener which is able to create a parse tree/forest
// from the Earley-states. Users may create one and call it themselves, but the more
// common usage pattern is by setting the GenerateTree-option for a parser and
// retrieving the parse-tree/forest from parser.Forest.
type TreeBuilder struct {
	forest  *sppf.Forest
	grammar *lr.Grammar
}

// NewTreeBuilder creates a TreeBuilder given an input grammar. This should obviously
// be the same grammar as the one used for parsing, but this is not checked.
// The TreeBuilder uses the grammar to
func NewTreeBuilder(g *lr.Grammar) *TreeBuilder {
	return &TreeBuilder{
		forest:  sppf.NewForest(),
		grammar: g,
	}
}

// Forest returns the parse forest after walking the derivation tree.
func (tb *TreeBuilder) Forest() *sppf.Forest {
	return tb.forest
}

// Reduce is a listener method, called for Earley-completions.
func (tb *TreeBuilder) Reduce(sym *lr.Symbol, rule int, rhs []*RuleNode, span lr.Span, level int) interface{} {
	if len(rhs) == 0 {
		return tb.forest.AddEpsilonReduction(sym, rule, span.From())
	}
	treenodes := make([]*sppf.SymbolNode, len(rhs))
	for i, r := range rhs {
		treenodes[i] = r.Value.(*sppf.SymbolNode)
	}
	return tb.forest.AddReduction(sym, rule, treenodes)
}

// Terminal is a listener method, called when matching input tokens.
func (tb *TreeBuilder) Terminal(tokval int, token interface{}, span lr.Span, level int) interface{} {
	// TODO
	t := tb.grammar.Terminal(tokval)
	return tb.forest.AddTerminal(t, span.From())
}

var _ Listener = &TreeBuilder{}

// --- Helpers ----------------------------------------------------------

// Reverse the symbols of a RHS of a rule (i.e., a handle)
// Creates a new slice.
func reverse(syms []*lr.Symbol) []*lr.Symbol {
	r := append([]*lr.Symbol(nil), syms...) // make copy first
	for i := len(syms)/2 - 1; i >= 0; i-- {
		opp := len(syms) - 1 - i
		r[i], r[opp] = r[opp], r[i]
	}
	return r
}

/*
From http://loup-vaillant.fr/tutorials/earley-parsing/parser
The author states:

	A completed item only stores its beginning and its rule. Its end is implicit:
	it's the Earley set it is stored on. We can reverse that. Instead of having this:

		=== 9 ===
		Product -> Factor (2)

	We could have the beginning be implicit, and store the end. Like that:

		=== 2 ===
		Product -> Factor (9)

	It is basically the same thing, but now we can perform searches from the beginning.

Unfortunately, there is a complication when searching from the beginning: We would
need to check for terminals in the input against successors of non-terminals
of the completion-sets. The reason is, that the completion sets contain dead-ends,
i.e. completions which did not produce a valid shift on the lookahead.
When searching from the back, we never see these dead-ends. When searching from the
beginning, we need backtracking to identify them.
*/

// reverseStates reverses the states after a successful parse, following the idea
// of http://loup-vaillant.fr/tutorials/earley-parsing/parser
// However, currently it seems not very useful.
func (p *Parser) reverseStates() []*iteratable.Set {
	l := len(p.states)
	reversed := make([]*iteratable.Set, l)
	for n, S := range p.states {
		reversed[n] = iteratable.NewSet(0)
		R := S.Subset(func(el interface{}) bool {
			item := el.(lr.Item)
			return item.PeekSymbol() == nil
		}) // now R contains only completion-items
		R.Each(func(el interface{}) {
			item := el.(lr.Item)
			o := item.Origin // misuse as span.To
			item.Origin = uint64(n)
			reversed[o].Add(item)
		})
	}
	for n := range reversed {
		dumpState(reversed, uint64(n))
	}
	return reversed
}
