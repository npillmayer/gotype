/*
Package glr implements a small-scale GLR(1)-parser.
It is mainly intended for Markdown parsing, but may be of use for
other purposes, too.

The parser relies on tables built from an lr.Grammar object (see
package documentation there).

Note: The API is still very much in flux! Currently it is something like

	scanner := parser.NewStdScanner(strings.NewReader("some input text"))
	p := glr.Create(grammar, gotoTable, actionTable)
	p.Parse(startState, scanner)

----------------------------------------------------------------------

BSD License

Copyright (c) 2017–2020, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software or the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package glr

import (
	"fmt"
	"io"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/dss"
	"github.com/npillmayer/gotype/syntax/lr/sparse"
)

// T traces to the SyntaxTracer
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

// A Parser type. Create and initialize one with parser.Create(...)
type Parser struct {
	G         *lr.Grammar
	dss       *dss.DSSRoot      // stack
	gotoT     *sparse.IntMatrix // GOTO table
	actionT   *sparse.IntMatrix // ACTION table
	accepting []int             // slice of accepting states
}

// Create and initialize a parser object, providing information from an
// lr.LRTableGenerator. Clients have to provide a link to the grammar and the
// parser tables.
func Create(g *lr.Grammar, gotoTable *sparse.IntMatrix, actionTable *sparse.IntMatrix) *Parser {
	parser := &Parser{
		G:       g,
		gotoT:   gotoTable,
		actionT: actionTable,
	}
	return parser
}

// Parse startes a new parse, given a start state and a scanner tokenizing the input.
// The parser must have been initialized.
func (p *Parser) Parse(S *lr.CFSMState, scan Scanner) (bool, error) {
	if p.G == nil || p.gotoT == nil {
		T().Errorf("GLR parser not initialized")
		return false, fmt.Errorf("GLR parser not initialized")
	}
	p.dss = dss.NewRoot("G", -1)    // forget existing one, if present
	start := dss.NewStack(p.dss)    // create first stack instance in DSS
	start.Push(S.ID, p.G.Epsilon()) // push the start state onto the stack
	accepting := false
	done := false
	for !done {
		tval, token := scan.NextToken(nil)
		if token == nil {
			tval = scanner.EOF
		}
		T().Debugf("got token %v from scanner", token)
		activeStacks := p.dss.ActiveStacks()
		T().P("glr", "parse").Debugf("currently %d active stack(s)", len(activeStacks))
		for _, stack := range activeStacks {
			p.reducesAndShiftsForToken(stack, tval)
		}
		if p.checkAccepted() {
			T().Errorf("ACCEPT")
			accepting, done = true, true
		}
		if tval == scanner.EOF {
			done = true
		}
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~")
	}
	return accepting, nil
}

// With a new lookahead (tokval): execute all possible reduces and shifts,
// cascading. The general outline is as follows:
//
//   1. do until no more reduces:
//      1.a if action(s) =
//           | shift: store stack and params in set S
//           | reduce: do reduce and store stack and params in set R
//           | conflict: shift/reduce or reduce/reduce
//              | do reduce(s) and store stack(s) in S or R respectively
//      1.b iterate again with R
//   2. shifts are now collected in S => execute
func (p *Parser) reducesAndShiftsForToken(stack *dss.Stack, tokval int) {
	var heads [2]*dss.Stack
	var actions [2]int32
	S := newStackSet() // will collect shift actions
	R := newStackSet() // re-consider stack/action for reduce
	R = R.add(stack)   // start with this active stack
	for !R.empty() {
		heads[0] = R.get()
		stateID, sym := heads[0].Peek()
		T().P("dss", "TOS").Debugf("state = %d, symbol = %v", stateID, sym)
		actions[0], actions[1] = p.actionT.Values(stateID, tokval)
		if actions[0] == p.actionT.NullValue() {
			T().Infof("no entry in ACTION table found, parser dies")
			heads[0].Die()
		} else {
			headcnt := 1
			T().Debugf("action 1 = %d, action 2 = %d", actions[0], actions[1])
			conflict := actions[1] != p.actionT.NullValue()
			if conflict { // shift/reduce or reduce/reduce conflict
				T().Infof("conflict, forking stack")
				heads[1] = stack.Fork() // must happen before action 1 !
				headcnt = 2
			}
			for i := 0; i < headcnt; i++ {
				if actions[i] >= 0 { // reduce action
					stacks := p.reduce(stateID, p.G.Rule(int(actions[i])), heads[i])
					R = R.add(stacks...)
				} else { // shift action
					S = S.add(heads[i])
				}
			}
		}
		T().Infof("%d shift operations in S", len(S))
		for !S.empty() {
			heads[0] = S.get()
			p.shift(stateID, tokval, heads[0])
		}
	}
}

func (p *Parser) shift(stateID int, tokval int, stack *dss.Stack) []*dss.Stack {
	nextstate := p.gotoT.Value(stateID, tokval)
	T().Infof("shifting %v to %d", tokenString(tokval), nextstate)
	terminal := p.G.GetTerminalSymbolFor(tokval)
	head := stack.Push(int(nextstate), terminal)
	return []*dss.Stack{head}
}

func (p *Parser) reduce(stateID int, rule *lr.Rule, stack *dss.Stack) []*dss.Stack {
	T().Infof("reduce %v", rule)
	handle := rule.GetRHS()
	heads := stack.Reduce(handle)
	if heads != nil {
		T().Debugf("reduce resulted in %d stacks", len(heads))
		lhs := rule.GetLHSSymbol()
		for i, head := range heads {
			state, _ := head.Peek()
			T().Debugf("state on stack#%d is %d", i, state)
			nextstate := p.gotoT.Value(state, lhs.GetID())
			newhead := head.Push(int(nextstate), lhs)
			T().Debugf("new head = %v", newhead)
		}
	}
	return heads
}

func (p *Parser) checkAccepted() bool {
	for _, stack := range p.dss.ActiveStacks() {
		state, _ := stack.Peek()
		for _, accstate := range p.accepting {
			if state == accstate {
				return true
			}
		}
	}
	return false
}

// ---------------------------------------------------------------------------

// helper: set of stacks
type stackSet []*dss.Stack

// TODO use sync.pool
func newStackSet() stackSet {
	s := make([]*dss.Stack, 0, 5)
	return stackSet(s)
}

// add a stack to the set
func (sset stackSet) add(stack ...*dss.Stack) stackSet {
	return append(sset, stack...)
}

// get a stack from the set
func (sset *stackSet) get() *dss.Stack {
	l := len(*sset)
	if l == 0 {
		return nil
	}
	s := (*sset)[l-1]
	(*sset)[l-1] = nil
	*sset = (*sset)[:l-1]
	return s
}

// is this set empty?
func (sset stackSet) empty() bool {
	return len(sset) == 0
}

// make a stack set empty
func (sset *stackSet) clear() {
	for k := len(*sset) - 1; k >= 0; k-- {
		(*sset)[k] = nil
	}
	*sset = (*sset)[:0]
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

// Scanner is an interface the parser relies on.
type Scanner interface {
	MoveTo(position uint64)
	NextToken(expected []int) (tokval int, token interface{})
}

func tokenString(tok int) string {
	return scanner.TokenString(rune(tok))
}

func (token *Token) String() string {
	return fmt.Sprintf("(%s:%d|\"%s\")", tokenString(token.Value), token.Value,
		string(token.Lexeme))
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
