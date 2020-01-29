/*
Package earley will some day implement an Earley-Parser.

	TODO error handling: provide an error function, like scanner.Scanner
	TODO semantic actions

There already exists a solution in

	https://github.com/jakub-m/gearley (unfortunately not usable for me)

which is based on the nice Earley-introduction from

	http://loup-vaillant.fr/tutorials/earley-parsing/

(which boasts an implementation in Lua and OcaML)


BSD License

Copyright (c) 2019–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */
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

// The parser may be controlled in some ways by setting mode flags.
// Set them with Parser.Mode
//
// If StoreTokens is not set, listeners will not be able to access the tokens.
const (
	StoreTokens  uint = 1 << 1 // store all input tokens, defaults to true
	GenerateTree uint = 1 << 2 // if parse was successful, generate a parse forest (default false)
)

// Parser is an Earley-parser type. Create and initialize one with earley.NewParser(...)
type Parser struct {
	GA      *lr.LRAnalysis    // the analyzed grammar we operate on
	scanner scanner.Tokenizer // scanner deliveres tokens
	states  []*iteratable.Set // list of states, each a set of Earley-items
	tokens  []interface{}     // we remember all input tokens, if requested
	SC      uint64            // state counter
	Mode    uint              // flags controlling some behaviour of the parser
	mode    uint              // private copy of Mode during parse
}

// NewParser creates and initializes an Earley parser.
func NewParser(ga *lr.LRAnalysis) *Parser {
	return &Parser{
		GA:      ga,
		scanner: nil,
		states:  make([]*iteratable.Set, 1, 512),
		tokens:  make([]interface{}, 1, 512),
		SC:      0,
		Mode:    StoreTokens,
		mode:    StoreTokens,
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
	p.mode = p.Mode // do not let client change during parse
	startItem, _ := lr.StartItem(p.GA.Grammar().Rule(0))
	p.states[0] = iteratable.NewSet(0) // S0
	p.states[0].Add(startItem)         // S0 = { [S′→•S, 0] }
	tokval, token, start, len := p.scanner.NextToken(scanner.AnyToken)
	for { // outer loop over Si per input token xi
		T().Debugf("Scanner read '%v|%d' @ %d", token, tokval, start)
		x := inputSymbol{tokval, token, span{start, start + len - 1}}
		i := p.setupNextState(token)
		p.innerLoop(i, x)
		if tokval == scanner.EOF {
			break
		}
		tokval, token, start, len = p.scanner.NextToken(scanner.AnyToken)
	}
	return p.checkAccept(), nil
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

// Invariant: we're in set Si and prepare Si+1
func (p *Parser) setupNextState(token interface{}) uint64 {
	// first one has already been created before outer loop
	p.states = append(p.states, iteratable.NewSet(0))
	if p.hasmode(StoreTokens) {
		p.tokens = append(p.tokens, token)
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
	dumpState(p.states, i)
}

// Scanner:
// If [A→…•a…, j] is in Si and a=xi+1, add [A→…a•…, j] to Si+1
func (p *Parser) scan(S, S1 *iteratable.Set, item lr.Item, tokval int) {
	if a := item.PeekSymbol(); a != nil {
		if a.Value == tokval {
			S1.Add(item.Advance())
		}
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
		S.Add(startitem)
	})
	//T().Debugf("start items from B=%s: %v", B, itemSetString(startitemsForB))
	if p.GA.DerivesEpsilon(B) { // B is nullable?
		//T().Debugf("%s is nullable", B)
		S.Add(item.Advance())
	}
}

// Completer:
// If [A→…•, j] is in Si, add [B→…A•…, k] to Si for all items [B→…•A…, k] in Sj.
func (p *Parser) complete(S, S1 *iteratable.Set, item lr.Item) {
	if item.PeekSymbol() == nil { // dot is behind RHS
		A, j := item.Rule().LHS, item.Origin
		//T().Debugf("Completing rule for %s: %s", A, item)
		Sj := p.states[j]
		//R := Sj.Copy()
		//T().Debugf("   search predecessors: %s", itemSetString(R))
		R := Sj.Copy().Subset(func(e interface{}) bool { // find all [B→…•A…, k]
			jtem := e.(lr.Item)
			// if jtem.PeekSymbol() == A {
			// 	T().Debugf("    => %s", jtem)
			// }
			return jtem.PeekSymbol() == A
		})
		//T().Debugf("   found predecessors: %s", itemSetString(R))
		R.Each(func(e interface{}) { // now add [B→…A•…, k]
			jtem := e.(lr.Item)
			if jadv := jtem.Advance(); jadv != lr.NullItem {
				S.Add(jadv)
			}
		})
	}
}

// checkAccepts searches the final state for items with a dot after #eof
// and a LHS of the start rule.
// It returns true if an accepting item has been found, indicating that the
// input has been recognized.
func (p *Parser) checkAccept() bool {
	dumpState(p.states, p.SC)
	S := p.states[p.SC] // last state should contain accept item
	S.IterateOnce()
	acc := false
	for S.Next() {
		item := S.Item().(lr.Item)
		if item.PeekSymbol() == nil && item.Rule().LHS == p.GA.Grammar().Rule(0).LHS {
			T().Debugf("ACCEPT: %s", item)
			acc = true
		}
	}
	return acc
}

func (p *Parser) hasmode(m uint) bool {
	return p.mode&m > 0
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
