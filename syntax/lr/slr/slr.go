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
	"io"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/sparse"
)

// T traces to the global SyntaxTracer.
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

// Parser is an SLR(1)-parser type. Create and initialize one with slr.NewParser(...)
type Parser struct {
	G *lr.Grammar
	//stack   *stack            // parser stack
	stack   []stackitem       // parser stack
	gotoT   *sparse.IntMatrix // GOTO table
	actionT *sparse.IntMatrix // ACTION table
}

// We store pairs of state-IDs and symbol-IDs on the parse stack.
type stackitem struct {
	ID    int
	symID int
}

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
type Scanner interface {
	MoveTo(position uint64)
	NextToken(expected []int) (tokval int, token interface{})
}

// Parse startes a new parse, given a start state and a scanner tokenizing the input.
// The parser must have been initialized.
//
// The parser returns true if the input string has been accepted.
func (p *Parser) Parse(S *lr.CFSMState, scan Scanner) (bool, error) {
	T().Debugf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	if p.G == nil || p.gotoT == nil {
		T().Errorf("SLR(1)-parser not initialized")
		return false, fmt.Errorf("SLR(1)-parser not initialized")
	}
	var accepting bool
	p.stack = append(p.stack, stackitem{S.ID, 0}) // push S
	// http://www.cse.unt.edu/~sweany/CSCE3650/HANDOUTS/LRParseAlg.pdf
	tokval, token := scan.NextToken(nil)
	done := false
	for !done {
		if token == nil {
			tokval = scanner.EOF
		}
		T().Debugf("got token %s/%d from scanner", token, tokval)
		state := p.stack[len(p.stack)-1] // TOS
		action := p.actionT.Value(state.ID, tokval)
		T().Debugf("action(%d,%d)=%s", state.ID, tokval, valstring(action, p.actionT))
		if action == p.actionT.NullValue() {
			return false, fmt.Errorf("Syntax error at %d/%v", tokval, token)
		}
		if action == lr.AcceptAction {
			T().Infof("ACCEPT")
			accepting = true
			done = true
		} else if action == lr.ShiftAction {
			nextstate := int(p.gotoT.Value(state.ID, tokval))
			T().Debugf("shifting, next state = %d", nextstate)
			p.stack = append(p.stack, stackitem{nextstate, tokval}) // push
			tokval, token = scan.NextToken(nil)
		} else if action > 0 { // reduce action
			rule := p.G.Rule(int(action))
			nextstate := p.reduce(state.ID, rule)
			T().Debugf("next state = %d", nextstate)
			p.stack = append(p.stack, stackitem{nextstate, rule.GetLHSSymbol().GetID()}) // push
		} else { // no action found
			done = true
		}
	}
	return accepting, nil
}

func (p *Parser) reduce(stateID int, rule *lr.Rule) int {
	T().Infof("reduce %v", rule)
	handle := reverse(rule.GetRHS())
	for _, sym := range handle {
		p.stack = p.stack[:len(p.stack)-1] // pop TOS
		tos := p.stack[len(p.stack)-1]
		if tos.symID != sym.GetID() {
			// if tos := p.stack.Pop(); sym.GetID() != tos.ID {
			T().Errorf("Expected %v on top of stack, got %d", sym, tos.symID)
			// }
		}
	}
	lhs := rule.GetLHSSymbol()
	state := p.stack[len(p.stack)-1] // TOS
	nextstate := p.gotoT.Value(state.ID, lhs.GetID())
	return int(nextstate)
}

// --- Scanner ----------------------------------------------------------

// A Token type, if you want to use it. Tokens of this type are returned
// by StdScanner.
//
// Clients may provide their own token data type.
type Token struct {
	Value  int
	Lexeme []byte
}

// StdScanner provides a default scanner implementation, but clients are free (and
// even encouraged) to provide their own. This implementation is based on
// stdlib's text/scanner.
type StdScanner struct {
	reader io.Reader // will be io.ReaderAt in the future
	scan   scanner.Scanner
}

// NewStdScanner creates a new default scanner from a Reader.
func NewStdScanner(r io.Reader) *StdScanner {
	s := &StdScanner{reader: r}
	s.scan.Init(r)
	s.scan.Filename = "Go symbols"
	return s
}

// MoveTo is not functional for default scanners.
// Default scanners allow sequential processing only.
func (s *StdScanner) MoveTo(position uint64) {
	T().Errorf("MoveTo() not yet supported by parser.StdScanner")
}

// NextToken gets the next token scanned from the input source. Returns the token
// value and a user-defined token type.
//
// Clients may provide an array of token values, one of which is expected
// at the current parse position. For the default scanner, as of now this is
// unused. In the future it will help with error-repair.
func (s *StdScanner) NextToken(expected []int) (int, interface{}) {
	tokval := int(s.scan.Scan())
	token := &Token{Value: tokval, Lexeme: []byte(s.scan.TokenText())}
	T().P("token", tokenString(tokval)).Debugf("scanned token at %s = \"%s\"",
		s.scan.Position, s.scan.TokenText())
	return tokval, token
}

func tokenString(tok int) string {
	return scanner.TokenString(rune(tok))
}

// --- Helpers ----------------------------------------------------------

// reverse the symbols of a RHS of a rule (i.e., a handle)
func reverse(syms []lr.Symbol) []lr.Symbol {
	r := append([]lr.Symbol(nil), syms...) // make copy first
	for i := len(syms)/2 - 1; i >= 0; i-- {
		opp := len(syms) - 1 - i
		syms[i], syms[opp] = syms[opp], syms[i]
	}
	return r
}

// valstring is a short helper to stringify an action table entry.
func valstring(v int32, m *sparse.IntMatrix) string {
	if v == m.NullValue() {
		return "<none>"
	}
	return fmt.Sprintf("%d", v)
}
