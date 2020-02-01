package sppf

import "github.com/npillmayer/gotype/syntax/lr"

/*
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

/*
Traversing a parse forest resulting from an ambiguous grammar in practice
mainly come in two variants:

- The choice of parse tree is irrelevant. An example for this is determining
  the result of an arithmetic expression, ignoring associativity: 1+2+3=5,
  independently of associativity.

- The client has additional knowledge of how to prune the parse forest and
  select one tree. Using again the example from above: selecting the tree
  representing left-associativity would result in the pruning of 1+(2+3),
  leaving (1+2)+3 as the un-ambiguous parse-tree.

Many heavily grammars can be invented (and often are for tests and for
research), but real-world clients most of the time have a clear understanding
of how strings for a given grammar are supposed to be structured. It is more
a problem of communicating this easily to a parser.

The focus of this package is therefore, to enable the user to prune ambiguous
parse trees, without making silent descicions which cannot be influenced by
the user. That means: having sensible defaults, but provide options for the
advanced user.
*/

// Listener is a type for walking a parse tree/forest.
type Listener interface {
	Reduce(*lr.Symbol, []*RuleNode, lr.Span, int) interface{}
	Terminal(int, interface{}, lr.Span, int) interface{}
	Conflict(*lr.Symbol, int, lr.Span, int) (int, error)
}

// RuleNode represents a node occuring during a parse tree/forest walk.
type RuleNode struct {
	sym    *lr.Symbol
	Extent lr.Span
	Value  interface{}
}

// Symbol returns the symbol a RuleNode refers to.
// It is either a terminal or the LHS of a reduced rule.
func (rnode *RuleNode) Symbol() *lr.Symbol {
	return rnode.sym
}

func (rnode *RuleNode) HasConflict() bool {
	return false // TODO
}

func TreeWalk(listener *Listener) {

}
