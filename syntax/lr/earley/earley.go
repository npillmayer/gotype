/*
Package earley will some day implement an Earley-Parser.

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
)

// T traces to the global syntax tracer.
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

// Scanner is a scanner-interface the parser relies on to receive the next input token.
type Scanner interface {
	//MoveTo(position uint64) // is this necessary?
	//NextToken(expected []int) (tokval int, token interface{}, start, len uint64)
	NextToken(expected []int) (tokval int, token interface{}, start, len uint64)
}

var any []int = nil

// Parser is an Earley-parser type. Create and initialize one with earley.NewParser(...)
type Parser struct {
	GA      *lr.LRAnalysis    // the analyzed grammar we operate on
	scanner Scanner           // scanner deliveres tokens
	states  []*iteratable.Set // list of states, each a set of Earley-items
	SC      int               // state counter
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

// Parse startes a new parse, given a start state and a scanner tokenizing the input.
// The parser must have been initialized.
//
// The parser returns true if the input string has been accepted.
func (p *Parser) Parse(scan Scanner) (bool, error) {
	if scan == nil {
		return false, fmt.Errorf("Earley-parser needs a valid scanner, is void")
	}
	p.scanner = scan
	startItem, _ := lr.StartItem(p.GA.Grammar().Rule(0))
	p.states[0] = iteratable.NewSet(0)
	p.states[0].Add(startItem)

	return false, nil
}

// tokval, token, start, len := p.scanner.NextToken(any)
// if tokval != scanner.EOF {
// 	T().Debugf("Scanner delivered '%v' @ %d (%d)", token, start, len)
// }

// http://citeseerx.ist.psu.edu/viewdoc/download
// http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.12.4254&rep=rep1&type=pdf
// From "Practical Earley Parsing" by John Aycock and R. Nigel Horspool, 2002:
//
// Earley parsers operate by constructing a sequence of sets, sometimes called
// Earley sets. Given an input
//       x1x2...xn,
// the parser builds n+1 sets: an initial set S0 and one set Si foreach input
// symbol xi. Elements of these sets are referred to as (Earley) items, which
// consist of three parts: a grammar rule, a position in the right-hand side
// of the rule indicating how much of that rule has been seen and a pointer
// to an earlier Earley set. Typically Earley items are written as
//       [A→α•β, j]
// where the position in the rule’s right-hand side is denoted bya dot (•)and
// j is a pointer to set Sj.
//
// In terms of implementation, the Earley sets are built in increasing order as
// the input is read. Also, each set is typically represented as a list of items,
// as suggested by Earley[...]. This list representation of a set is particularly
// convenient, because the list of items acts as a ‘work queue’ when building the set:
// items are examined in order, applying Scanner, Predictor and Completer as
// necessary; items added to the set are appended onto the end of the list.

// Scanner:
// If [A→…•a…, j] is in Si and a=xi+1, add [A→…a•…, j] to Si+1
func (p *Parser) scan(tokval int, pos uint64) {
	S := p.states[p.SC]
	S1 := p.states[p.SC+1]
	S.Iterate(false) // iterate over each item once
	for S.Next() {
		item := S.Item().(lr.Item)
		if item.PeekSymbol().Value == tokval {
			S1.Add(item.Advance())
		}
	}
}

// Predictor:
// If [A→…•B…, j] is in Si, add [B→•α, i] to Si for all rules B→α.
// If B is nullable, also add [A→…B•…, j] to Si.
func (p *Parser) predict(tokval int, pos uint64) {
	S := p.states[p.SC]
	S.Iterate(false) // iterate over each item once
	for S.Next() {
		item := S.Item().(lr.Item)
		B := item.PeekSymbol()
		rulesForB := p.GA.Grammar().FindNonTermRules(B, true)
		rulesForB.Each(func(e interface{}) {
			i := e.(lr.Item)
			i.Origin = pos
		})
		if p.GA.DerivesEpsilon(B) {
			S.Add(item.Advance())
		}
	}
}

// Completer:
// If [A→…•, j] is in Si, add [B→…A•…, k] to Si for all items [B→…•A…, k] in Sj.
func (p *Parser) complete(tokval int, pos uint64) {
	S := p.states[p.SC]
	S.Iterate(false) // iterate over each item once
	for S.Next() {
		item := S.Item().(lr.Item)
		if item.PeekSymbol() == nil { // dot is behind RHS
			A, j := item.Rule().LHS, item.Origin
			Sj := p.states[j]
			R := Sj.Copy().Subset(func(e interface{}) bool {
				jtem := e.(lr.Item)
				return jtem.PeekSymbol() == A
			})
			R.Each(func(e interface{}) {
				jtem := e.(lr.Item)
				if jadv := jtem.Advance(); jadv != lr.NullItem {
					S.Add(jadv)
				}
			})
		}
	}
}
