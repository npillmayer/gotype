/*
Package earley will some day implement an Earley-Parser.

	TODO error handling: provide an error function, like scanner.Scanner
	TODO semantic actions

There already exists a solution in

	https://github.com/jakub-m/gearley

which is based on the nice Earley-introduction from

	http://loup-vaillant.fr/tutorials/earley-parsing/

(which boasts an implementation in Lua and OcaML) */
package earley

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/iteratable"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
)

// T traces to the global syntax tracer.
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

// Parser is an Earley-parser type. Create and initialize one with earley.NewParser(...)
type Parser struct {
	GA      *lr.LRAnalysis    // the analyzed grammar we operate on
	scanner scanner.Tokenizer // scanner deliveres tokens
	states  []*iteratable.Set // list of states, each a set of Earley-items
	SC      uint64            // state counter
}

// NewParser creates and initializes an Earley parser.
func NewParser(ga *lr.LRAnalysis) *Parser {
	return &Parser{
		GA:      ga,
		scanner: nil,
		states:  make([]*iteratable.Set, 0, 512),
		SC:      0,
	}
}

// The parser consumes input symbols until the token value is EOF.
type inputSymbol struct {
	tokval int         // token value
	lexeme interface{} // visual representation of the symbol, if any
	span   span        // position and extent in the input stream
}

// Parse startes a new parse, given a scanner tokenizing the input.
// The parser must have been initialized with an analyzed grammar.
// It returns true if the input string has been accepted.
func (p *Parser) Parse(scan scanner.Tokenizer) (bool, error) {
	if p.scanner = scan; scan == nil {
		return false, fmt.Errorf("Earley-parser needs a valid scanner, is void")
	}
	startItem, _ := lr.StartItem(p.GA.Grammar().Rule(0))
	p.states[0] = iteratable.NewSet(0) // S0
	p.states[0].Add(startItem)         // S0 = { [S′→•S, 0] }
	tokval, token, start, len := p.scanner.NextToken(scanner.AnyToken)
	for tokval != scanner.EOF { // outer loop
		T().Debugf("Scanner read '%v|%d' @ %d (%d)", token, tokval, start, len)
		x := inputSymbol{tokval, token, span{start, start + len - 1}}
		i := p.setupNextState()
		p.innerLoop(i, x)
		tokval, token, start, len = p.scanner.NextToken(scanner.AnyToken)
	}
	return false, nil
}

// http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.12.4254&rep=rep1&type=pdf
// From "Practical Earley Parsing" by John Aycock and R. Nigel Horspool, 2002:
//
// Earley parsers operate by constructing a sequence of sets, sometimes called
// Earley sets. Given an input
//       x1 x2 … xn,
// the parser builds n+1 sets: an initial set S0 and one set Si foreach input
// symbol xi. Elements of these sets are referred to as (Earley) items, which
// consist of three parts: a grammar rule, a position in the right-hand side
// of the rule indicating how much of that rule has been seen and a pointer
// to an earlier Earley set. Typically Earley items are written as
//       [A→α•β, j]
// where the position in the rule’s right-hand side is denoted bya dot (•) and
// j is a pointer to set Sj.
//
// […] In terms of implementation, the Earley sets are built in increasing order as
// the input is read. Also, each set is typically represented as a list of items,
// as suggested by Earley[…]. This list representation of a set is particularly
// convenient, because the list of items acts as a ‘work queue’ when building the set:
// items are examined in order, applying Scanner, Predictor and Completer as
// necessary; items added to the set are appended onto the end of the list.

func (p *Parser) setupNextState() uint64 {
	if p.SC != 0 { // first one has already been created beforehand
		p.states = append(p.states, iteratable.NewSet(0))
	}
	i := p.SC
	p.SC++
	return i // ready to operate on set Si
}

// The inner loop iterates over Si, applying Scanner, Predictor and Completer.
// The variable for Si is called S and Si+1 is called S1.
func (p *Parser) innerLoop(i uint64, x inputSymbol) {
	S := p.states[i]
	S1 := p.states[i+1]
	S.IterateOnce()
	for S.Next() {
		item := S.Item().(lr.Item)
		p.scan(S, S1, item, x.tokval) // may add items to S1
		p.predict(S, S1, item, i)     // may add items to S
		p.complete(S, S1, item)       // may add items to S
	}
}

// Scanner:
// If [A→…•a…, j] is in Si and a=xi+1, add [A→…a•…, j] to Si+1
func (p *Parser) scan(S, S1 *iteratable.Set, item lr.Item, tokval int) {
	if item.PeekSymbol().Value == tokval {
		S1.Add(item.Advance())
	}
}

// Predictor:
// If [A→…•B…, j] is in Si, add [B→•α, i] to Si for all rules B→α.
// If B is nullable, also add [A→…B•…, j] to Si.
func (p *Parser) predict(S, S1 *iteratable.Set, item lr.Item, i uint64) {
	B := item.PeekSymbol()
	startitemsForB := p.GA.Grammar().FindNonTermRules(B, true)
	startitemsForB.Each(func(e interface{}) { // e is a start item
		startitem := e.(lr.Item)
		startitem.Origin = i
	})
	if p.GA.DerivesEpsilon(B) { // B is nullable?
		S.Add(item.Advance())
	}
}

// Completer:
// If [A→…•, j] is in Si, add [B→…A•…, k] to Si for all items [B→…•A…, k] in Sj.
func (p *Parser) complete(S, S1 *iteratable.Set, item lr.Item) {
	if item.PeekSymbol() == nil { // dot is behind RHS
		A, j := item.Rule().LHS, item.Origin
		Sj := p.states[j]
		R := Sj.Copy().Subset(func(e interface{}) bool { // find all [B→…•A…, k]
			jtem := e.(lr.Item)
			return jtem.PeekSymbol() == A
		})
		R.Each(func(e interface{}) { // now add [B→…A•…, k]
			jtem := e.(lr.Item)
			if jadv := jtem.Advance(); jadv != lr.NullItem {
				S.Add(jadv)
			}
		})
	}
}

// ----------------------------------------------------------------------

// span is a small type for capturing a length of input token.
type span [2]uint64 // start and end positions in the input string

func (s span) from() uint64 {
	return s[0]
}

func (s span) to() uint64 {
	return s[1]
}

func (s span) isNull() bool {
	return s == span{}
}
