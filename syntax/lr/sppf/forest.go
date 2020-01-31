package sppf

/*
Code for SPPFs are rare, mostly found in academic papers. One of them
is "SPPF-Style Parsing From Earley Recognisers" by Elizabeth Scott
(https://www.sciencedirect.com/science/article/pii/S1571066108001497).
It describes a binarised variant of an SPPF, which we will not follow.
A more accessible discussion of parse SPPFs may be found in
"Parsing Techniques" by  Dick Grune and Ceriel J.H. Jacobs
(https://dickgrune.com/Books/PTAPG_2nd_Edition/), Section 3.7.3.
Scott explains the downside of this simpler approach:
"We could [create] separate instances of the items for different substring matches,
so if [B→δ●,k], [B→σ●,k'] ∈ Ei where k≠k' then we create two copies of [D→τB●μ,h], one
pointing to each of the two items. In the above example we would create two items [S→SS●,0]
in E3, one in which the second S points to [S→b●,2] and the other in which the second S
points to [S→SS●,1]. This would cause correct derivations to be generated, but it also
effectively embeds all the derivation trees in the construction and, as reported by Johnson,
the size cannot be bounded by O(n^p) for any fixed integer p.
[...]
Grune has described a parser which exploits an Unger style parser to construct the
derivations of a string from the sets produced by Earley’s recogniser. However,
as noted by Grune, in the case where the number of derivations is exponential
the resulting parser will be of at least unbounded polynomial order in worst case."
(Notation slightly modified by me to conform to notations elsewhere in these
parser packages).

Despite the shortcomings of the forest described by Grune & Jacobs, I won't
implement Scott's improvements. For practical use, the worst case spatial complexity
seems never to materialize. However, gaining more insights in the future when using the SPPF
for more complex real word scenarios I'm prepared to reconsider.


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

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/iteratable"
)

// Forest implements a Shared Packed Parse Forest (SPPF).
// A packed parse forest re-uses existing parse tree nodes between different
// parse trees. For a conventional non-ambiguous parse, a parse forest consists
// of a single tree. Ambiguous grammars, on the other hand, may result in parse
// runs where more than one parse tree is created. To save space these parse
// trees will share common nodes.
//
// Our task is to store nodes representing recognition of a substring of the input, i.e.,
// [A→δ,(x…y)], where A is a grammar symbol and δ denotes the right hand side (RHS)
// of a production. (x…y) is the position interval within the input covered by A.
//
// We split up these nodes in two parts: A symbol node for A, and a RHS-node for δ.
// Symbol nodes fan out via or-edges to RHS-nodes. RHS-nodes fan out to the symbols
// of the RHS in order of their appearance in the corresponding grammar rule.
// If a tree segment is unambiguous, our node
// [A→δ,(x…y)] would be split into [A (x…y)]⟶[δ (x…y)], i.e. connected by an
// or-edge without alternatives (fan-out of A is 1).
// For ambiguous parses, subtrees can be shared if [δ (x…y)] is already found in
// the forest, meaning that there is another derivation of this input span present.
//
// How can we quickly identify nodes [A (x…y)] or [δ (x…y)] to find out if they
// are already present in the forest, and thus can be re-used?
// Nodes will be searched by span (x…y), followed by a check of either A or δ.
// This is implemented as a tree of height 2, with the edges labelled by input position
// and the leafs being sets of nodes. The tree is implemented by a map of maps of sets.
// We introduce a small helper type for it.
//
// Currently we carry the span in both types of nodes. This is redundant, and we
// will omit one of them in the future. For now, we use it to keep some operations
// simpler.

type searchTree map[uint64]map[uint64]*iteratable.Set // methods below

// Forest is a data structure for a "shared packed parse forest" (SPPF).
type Forest struct {
	symNodes searchTree                   // search tree for [A (x…y)]
	rhsNodes searchTree                   // search tree for [δ (x…y)]
	orEdges  map[*symNode]*iteratable.Set // or-edges from symbols to RHSs, indexed by symbol
	andEdges map[*rhsNode]*iteratable.Set // and-edges
}

// NewForest returns an empty forest.
func NewForest() *Forest {
	return &Forest{
		symNodes: make(map[uint64]map[uint64]*iteratable.Set),
		rhsNodes: make(map[uint64]map[uint64]*iteratable.Set),
		orEdges:  make(map[*symNode]*iteratable.Set),
		andEdges: make(map[*rhsNode]*iteratable.Set),
	}
}

// --- Nodes -----------------------------------------------------------------

// Nodes [A (x…y)] in the parse forest.
type symNode struct {
	symbol *lr.Symbol // A
	span   lr.Span    // positions in the input covered by this symbol
}

// Nodes [δ (x…y)] in the parse forest.
type rhsNode struct {
	rule int     // rule number this RHS δ is from
	span lr.Span // positions in the input covered by this RHS
}

func makeSym(symbol *lr.Symbol) *symNode {
	return &symNode{symbol: symbol}
}

// Use as makeSym(A).spanning(x, y), resulting in [A (x…y)]
func (sn *symNode) spanning(from, to uint64) *symNode {
	sn.span = lr.Span{from, to}
	return sn
}

// FindSymNode finds a (shared) node for a symbol node in the forest.
func (f *Forest) findSymNode(sym *lr.Symbol, start, end uint64) *symNode {
	return f.symNodes.findSymbol(start, end, sym)
}

// addSymNode adds a symbol node to the forest. Returns a reference to a symNode,
// which may already have been in the SPPF beforehand.
func (f *Forest) addSymNode(sym *lr.Symbol, start, end uint64) *symNode {
	sn := f.findSymNode(sym, start, end)
	if sn == nil {
		sn := makeSym(sym).spanning(start, end)
		f.symNodes.Add(start, end, sn)
	}
	return sn
}

func makeRHS(rule int) *rhsNode {
	return &rhsNode{rule: rule}
}

// Use as makeRHS(δ).spanning(x, y), resulting in [δ (x…y)]
func (rhs *rhsNode) spanning(from, to uint64) *rhsNode {
	rhs.span = lr.Span{from, to}
	return rhs
}

// FindRHSNode finds a (shared) node for a right hand side in the forest.
func (f *Forest) findRHSNode(rule int, start, end uint64) *rhsNode {
	return f.rhsNodes.findRHS(start, end, rule)
}

// addRHSNode adds a symbol node to the forest. Returns a reference to a rhsNode,
// which may already have been in the SPPF beforehand.
func (f *Forest) addRHSNode(rule int, start, end uint64) *rhsNode {
	rhs := f.findRHSNode(rule, start, end)
	if rhs == nil {
		rhs := makeRHS(rule).spanning(start, end)
		f.symNodes.Add(start, end, rhs)
	}
	return rhs
}

// --- Edges -----------------------------------------------------------------

// orEdges are ambiguity forks in the parse forest.
type orEdge struct {
	fromSym *symNode
	toRHS   *rhsNode
}

// addOrEdge inserts an edge between a symbol and a RHS.
// If start or end are not already contained in the forest, they are added.
//
// If the edge already exists, nothing is done.
func (f *Forest) addOrEdge(sym *lr.Symbol, rule int, start, end uint64) {
	sn := f.addSymNode(sym, start, end)
	rhs := f.addRHSNode(rule, start, end)
	if e := f.findOrEdge(sn, rhs); e.isVoid() {
		e = orEdge{sn, rhs}
		if _, ok := f.orEdges[sn]; !ok {
			f.orEdges[sn] = iteratable.NewSet(0)
		}
		f.orEdges[sn].Add(rhs)
	}
}

// findOrEdge finds an or-edge starting from a symbol and pointing to an
// RHS-node. If none is found, nullOrEdge is returned.
func (f *Forest) findOrEdge(sn *symNode, rhs *rhsNode) orEdge {
	if edges := f.orEdges[sn]; edges != nil {
		v := edges.FirstMatch(func(el interface{}) bool {
			e := el.(orEdge)
			return e.fromSym == sn && e.toRHS == rhs
		})
		return v.(orEdge)
	}
	return nullOrEdge
}

// nullOrEdge denotes an or-edge that is not present in a graph.
var nullOrEdge = orEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e orEdge) isVoid() bool {
	return e == nullOrEdge
}

// An andEdge connects a RHS to the symbols it consists of.
type andEdge struct {
	fromRHS  *rhsNode // RHS node starts the edge
	toSym    *symNode // symbol node this edge points to
	sequence uint     // sequence number 0…n, used for ordering children
}

// addAndEdge inserts an edge between a RHS and a symbol, labelled with a seqence
// number. If start or end are not already contained in the forest, they are
// added. Note that it cannot happen that two edges between identical nodes
// exist for different sequence numbers. The function panics if such a condition
// is found.
//
// If the edge already exists, nothing is done.
func (f *Forest) addAndEdge(rule int, sym *lr.Symbol, seq uint, start, end uint64) {
	rhs := f.addRHSNode(rule, start, end)
	sn := f.addSymNode(sym, start, end)
	if e := f.findAndEdge(rhs, sn); e.isVoid() {
		e = andEdge{rhs, sn, seq}
		if _, ok := f.andEdges[rhs]; !ok {
			f.andEdges[rhs] = iteratable.NewSet(0)
		}
		f.andEdges[rhs].Add(sn)
	} else if e.sequence != seq {
		panic(fmt.Sprintf("new edge with sequence=%d replaces sequence=%d", seq, e.sequence))
	}
}

// findAndEdge finds an and-edge starting from an RHS node and pointing to a
// symbol-node. If none is found, nullAndEdge is returned.
func (f *Forest) findAndEdge(rhs *rhsNode, sn *symNode) andEdge {
	if edges := f.andEdges[rhs]; edges != nil {
		v := edges.FirstMatch(func(el interface{}) bool {
			e := el.(andEdge)
			return e.fromRHS == rhs && e.toSym == sn
		})
		return v.(andEdge)
	}
	return nullAndEdge
}

// nullAndEdge denotes an and-edge that is not present in a graph.
var nullAndEdge = andEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e andEdge) isVoid() bool {
	return e == nullAndEdge
}

// --- Exported Functions ----------------------------------------------------

// --- searchTree -----------------------------------------------------------------

func (t searchTree) find(from, to uint64, predicate func(el interface{}) bool) interface{} {
	if t1, ok := t[from]; ok {
		if t2, ok := t1[to]; ok {
			return t2.FirstMatch(predicate)
		}
	}
	return nil
}

func (t searchTree) findSymbol(from, to uint64, sym *lr.Symbol) *symNode {
	node := t.find(from, to, func(el interface{}) bool {
		s := el.(*symNode)
		return s.symbol == sym
	})
	return node.(*symNode)
}

func (t searchTree) findRHS(from, to uint64, rule int) *rhsNode {
	node := t.find(from, to, func(el interface{}) bool {
		rhs := el.(*rhsNode)
		return rhs.rule == rule
	})
	return node.(*rhsNode)
}

func (t searchTree) Add(from, to uint64, item interface{}) {
	if t1, ok := t[from]; !ok {
		t[from] = make(map[uint64]*iteratable.Set)
		t[from][to] = iteratable.NewSet(0)
	} else if _, ok := t1[to]; !ok {
		t[from][to] = iteratable.NewSet(0)
	}
	t[from][to].Add(item)
}
