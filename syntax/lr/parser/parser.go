/*
Package parser implements a simple GLR-parser.
It is mainly intended for Markdown parsing, but may be of use for
other purposes, too.

----------------------------------------------------------------------

BSD License

Copyright (c) 2017–2018, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer or the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----------------------------------------------------------------------
*/
package parser

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/dss"
	"github.com/npillmayer/gotype/syntax/lr/sparse"
)

// Configurable trace
var T tracing.Trace = tracing.SyntaxTracer

/*
--------- scan -----------
https://golang.org/pkg/text/scanner/

https://github.com/timtadh/lexmachine
lexer framework
http://hackthology.com/writing-a-lexer-in-go-with-lexmachine.html

Buffered Reader
https://github.com/SteelSeries/bufrr

https://github.com/opennota/re2dfa
regular expression to DFA

--------- parse -----------
http://www.cse.unt.edu/~sweany/CSCE3650/HANDOUTS/LRParseAlg.pdf
*/

var testscanner scanner.Scanner

func init() {
	const src = `a b a a`
	testscanner.Init(strings.NewReader(src))
	testscanner.Filename = "example"
}

type Token struct {
	Value  int
	Lexeme string
}

func TokenString(tok int) string {
	return scanner.TokenString(rune(tok))
}

func (token *Token) String() string {
	return fmt.Sprintf("(%s:%d|\"%s\")", TokenString(token.Value), token.Value, token.Lexeme)
}

type Scanner struct {
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{}
}

func (s *Scanner) NextToken(expected []int) (int, interface{}) {
	tok := testscanner.Scan()
	token := &Token{Value: int(tok), Lexeme: testscanner.TokenText()}
	T.P("token", TokenString(int(tok))).Debugf("scanned token at %s = \"%s\"",
		testscanner.Position, testscanner.TokenText())
	return int(tok), token
}

type Parser struct {
	G       *lr.Grammar
	dss     *dss.DSSRoot
	gotoT   *sparse.IntMatrix
	actionT *sparse.IntMatrix
}

func Create(g *lr.Grammar, gotoTable *sparse.IntMatrix, actionTable *sparse.IntMatrix) *Parser {
	parser := &Parser{G: g, gotoT: gotoTable, actionT: actionTable}
	return parser
}

/*
Start a new parse, given a start state and a scanner tokenizing the input.
The parser must have been initialized.
*/
func (p *Parser) Parse(S *lr.CFSMState, scanner *Scanner) {
	if p.G == nil {
		T.Error("parser not initialized")
		return
	}
	p.dss = dss.NewRoot("G", -1)    // forget existing one, if present
	start := dss.NewStack(p.dss)    // create first stack instance in DSS
	start.Push(S.ID, p.G.Epsilon()) // push the start state onto the stack
	for {
		tval, token := scanner.NextToken(nil)
		if token == nil {
			break // EOF
		}
		T.Debugf("got token %v from scanner", token)
		activeStacks := p.dss.ActiveStacks()
		T.P("lr", "parse").Debugf("currently %d active stacks", len(activeStacks))
		for _, stack := range activeStacks {
			stateID, sym := stack.Peek()
			T.P("dss", "TOS").Debugf("state = %d, symbol = %v", stateID, sym)
			action1, action2 := p.actionT.Values(stateID, tval)
			if action1 == p.actionT.NullValue() {
				T.Info("no entry in ACTION table found")
			} else {
				if action2 != 0 {
					T.Info("conflict, resolving with stack.fork()")
				}
				if action1 < 0 {
					T.Infof("reduce %d", action1)
				} else {
					T.Info("shift")
				}
				if action2 < 0 {
					T.Infof("reduce %d", action2)
				} else if action2 > 0 {
					T.Info("shift")
				}
			}
		}
		break
	}
}
