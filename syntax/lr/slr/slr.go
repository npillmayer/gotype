/*
Package slr provides an SLR(1)-parser. Clients have to use the tools
of package lr to prepare the necessary parse tables. The SLR parser
utilizes these tables to create a right derivation for a given input,
provided through a scanner interface.

This parser is intended for small to moderate grammars, e.g. for configuration
input or small domain-specific languages. It is *not* intended for full-fledged
programming languages (there are superb other tools around for these kinds of
usages, usually creating LALR(1)-parsers, which are able to recognize a super-set
of SLR-languages).

The main focus for this implementation is adaptability and on-the-fly usage.
Clients are able to construct the parse tables from a grammar and use the
parser directly, without a code-generation or compile step. If you want, you
can create a grammar from user input and use a parser for it in a couple of
lines of code.

Package slr can only handle SLR(1) grammars. All SLR-grammars are deterministic
(but not vice versa). For parsing ambiguous grammars, see package glr.

Usage

Clients construct a grammar, usually by using a grammar builder:

	b := lr.NewGrammarBuilder("Signed Variables Grammar")
	b.LHS("Var").N("Sign").T("a", scanner.Ident).End()  // Var  --> Sign Id
	b.LHS("Sign").T("+", '+').End()                     // Sign --> +
	b.LHS("Sign").T("-", '-').End()                     // Sign --> -
	b.LHS("Sign").Epsilon()                             // Sign -->
	g, err := b.Grammar()

This grammar is subjected to grammar analysis and table generation.

	ga := lr.NewGrammarAnalysis(g)
	lrgen := lr.NewTableGenerator(ga)
	lrgen.CreateTables()
	if lrgen.HasConflicts { ... }  // cannot use an SLR parser

Finally parse some input:

	p := slr.NewParser(g, lrgen.GotoTable(), lrgen.ActionTable())
	scanner := slr.NewStdScanner(string.NewReader("+a")
	accepted, err := p.Parse(lrgen.CFSM().S0, scanner)

Clients may instrument the grammar with semantic operations or let the
parser create a parse tree. See the examples below.

Warning

This is a very early implementation. Currently you should use it for study purposes
only. The API may change significantly without prior notice.


BSD License

Copyright (c) 2017–20, Norbert Pillmayer

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
package slr

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/npillmayer/gotype/syntax/lr/sparse"
)

// T traces to the global SyntaxTracer.
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

// Parser is an SLR(1)-parser type. Create and initialize one with slr.NewParser(...)
type Parser struct {
	G       *lr.Grammar
	stack   []stackitem       // parser stack
	gotoT   *sparse.IntMatrix // GOTO table
	actionT *sparse.IntMatrix // ACTION table
}

// We store pairs of state-IDs and symbol-IDs on the parse stack.
type stackitem struct {
	stateID int  // ID of a CFSM state
	symID   int  // ID of a grammar symbol (terminal or non-terminal)
	span    span // input span over which this symbol reaches
}

// span is a small type for capturing a length of input token run. For every
// terminal and non-terminal, a parse tree/forest will track which input positions
// this symbol covers.
// Some useful operations on spans are to be found further down in this file.
type span [2]uint64 // start and end positions in the input string

// NewParser creates an SLR(1) parser.
func NewParser(g *lr.Grammar, gotoTable *sparse.IntMatrix, actionTable *sparse.IntMatrix) *Parser {
	parser := &Parser{
		G:       g,
		stack:   make([]stackitem, 0, 512),
		gotoT:   gotoTable,
		actionT: actionTable,
	}
	return parser
}

// Scanner is a scanner-interface the parser relies on to receive the next input token.
// type Scanner interface {
// 	MoveTo(position uint64)
// 	NextToken(expected []int) (tokval int, token interface{}, start, len uint64)
// }

// Parse startes a new parse, given a start state and a scanner tokenizing the input.
// The parser must have been initialized.
//
// The parser returns true if the input string has been accepted.
func (p *Parser) Parse(S *lr.CFSMState, scan scanner.Tokenizer) (bool, error) {
	T().Debugf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	if p.G == nil || p.gotoT == nil {
		T().Errorf("SLR(1)-parser not initialized")
		return false, fmt.Errorf("SLR(1)-parser not initialized")
	}
	var accepting bool
	p.stack = append(p.stack, stackitem{S.ID, 0, span{0, 0}}) // push S
	// http://www.cse.unt.edu/~sweany/CSCE3650/HANDOUTS/LRParseAlg.pdf
	tokval, token, pos, length := scan.NextToken(nil)
	done := false
	for !done {
		if token == nil {
			tokval = scanner.EOF
		}
		T().Debugf("got token %s/%d from scanner", token, tokval)
		state := p.stack[len(p.stack)-1] // TOS
		action := p.actionT.Value(state.stateID, tokval)
		T().Debugf("action(%d,%d)=%s", state.stateID, tokval, valstring(action, p.actionT))
		if action == p.actionT.NullValue() {
			return false, fmt.Errorf("Syntax error at %d/%v", tokval, token)
		}
		if action == lr.AcceptAction {
			accepting, done = true, true
			// TODO patch start symbol to have span(0,pos)
		} else if action == lr.ShiftAction {
			nextstate := int(p.gotoT.Value(state.stateID, tokval))
			T().Debugf("shifting, next state = %d", nextstate)
			p.stack = append(p.stack, // push a terminal state onto stack
				stackitem{nextstate, tokval, span{pos, pos + length - 1}})
			tokval, token, pos, length = scan.NextToken(nil)
		} else if action > 0 { // reduce action
			rule := p.G.Rule(int(action))
			nextstate, handlespan := p.reduce(state.stateID, rule)
			if handlespan.isNull() { // resulted from an epsilon production
				handlespan = span{pos - 1, pos - 1} // epsilon was just before lookahead
			}
			T().Debugf("reduced to next state = %d", nextstate)
			p.stack = append(p.stack, // push a non-terminal state onto stack
				stackitem{nextstate, rule.LHS.Value, handlespan})
		} else { // no action found
			done = true
		}
	}
	return accepting, nil
}

// reduce performs a reduce action for a rule
//
//    LHS --> X1 ... Xn   (with X being terminals or non-terminals)
//
// Symbols X1 to Xn should be represented on the stack as states
//
//    [TOS]  Sn(Xn, span_n) ... S1(X1, span1)  ...
//
// TODO: perform semantic action on reduce, either by calling a user-provided
// function from the grammar, or by constructing a node in a parse tree/forest.
func (p *Parser) reduce(stateID int, rule *lr.Rule) (int, span) {
	T().Infof("reduce %v", rule)
	handle := reverse(rule.RHS())
	var handlespan span
	for _, sym := range handle {
		p.stack = p.stack[:len(p.stack)-1] // pop TOS
		tos := p.stack[len(p.stack)-1]
		if tos.symID != sym.Value {
			T().Errorf("Expected %v on top of stack, got %d", sym, tos.symID)
		}
		handlespan = handlespan.extendFrom(tos.span)
	}
	lhs := rule.LHS
	// TODO: now perform sematic action or parse-tree build
	state := p.stack[len(p.stack)-1] // TOS
	nextstate := p.gotoT.Value(state.stateID, lhs.Value)
	return int(nextstate), handlespan
}

// --- Scanner ----------------------------------------------------------

// A Token type, if you want to use it. Tokens of this type are returned
// by StdScanner.
//
// Clients may provide their own token data type.
/*
type Token struct {
	Value  int
	Lexeme []byte
}
*/

// StdScanner provides a default scanner implementation, but clients are free (and
// even encouraged) to provide their own. This implementation is based on
// stdlib's text/scanner.
/*
type StdScanner struct {
	reader io.Reader // will be io.ReaderAt in the future
	scan   scanner.Tokenizer
}
*/

// NewStdScanner creates a new default scanner from a Reader.
/*
func NewStdScanner(r io.Reader) *StdScanner {
	s := &StdScanner{reader: r}
	s.scan.Init(r)
	s.scan.Filename = "Go symbols"
	return s
}
*/

// MoveTo is not functional for default scanners.
// Default scanners allow sequential processing only.
/*
func (s *StdScanner) MoveTo(position uint64) {
	T().Errorf("MoveTo() not yet supported by parser.StdScanner")
	panic("MoveTo() not yet supported by parser.StdScanner")
}
*/

// NextToken gets the next token scanned from the input source. Returns the token
// value and a user-defined token type.
//
// Clients may provide an array of token values, one of which is expected
// at the current parse position. For the default scanner, as of now this is
// unused. In the future it will help with error-repair.
/*
func (s *StdScanner) NextToken(expected []int) (int, interface{}, uint64, uint64) {
	tokval := int(s.scan.Scan())
	token := &Token{Value: tokval, Lexeme: []byte(s.scan.TokenText())}
	T().P("token", tokenString(tokval)).Debugf("scanned token at %s = \"%s\"",
		s.scan.Position, s.scan.TokenText())
	return tokval, token, uint64(s.scan.Position.Offset),
		uint64(s.scan.Pos().Offset - s.scan.Position.Offset)
}

func tokenString(tok int) string {
	return scanner.TokenString(rune(tok))
}
*/

// --- spans ----------------------------------------

func (s span) from() uint64 {
	return s[0]
}

func (s span) to() uint64 {
	return s[1]
}

func (s span) isNull() bool {
	return s == span{}
}

func (s span) extendFrom(other span) span {
	if other.from() < s.from() {
		s[0] = other[0]
	}
	if other.to() > s.to() {
		s[1] = other[1]
	}
	return s
}

// --- Helpers ----------------------------------------------------------

// reverse the symbols of a RHS of a rule (i.e., a handle)
func reverse(syms []*lr.Symbol) []*lr.Symbol {
	r := append([]*lr.Symbol(nil), syms...) // make copy first
	for i := len(syms)/2 - 1; i >= 0; i-- {
		opp := len(syms) - 1 - i
		r[i], r[opp] = r[opp], r[i]
	}
	return r
}

// valstring is a short helper to stringify an action table entry.
func valstring(v int32, m *sparse.IntMatrix) string {
	if v == m.NullValue() {
		return "<none>"
	} else if v == lr.AcceptAction {
		return "<accept>"
	} else if v == lr.ShiftAction {
		return "<shift>"
	}
	return fmt.Sprintf("%d", v)
}
