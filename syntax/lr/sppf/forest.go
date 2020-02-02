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

“We could [create] separate instances of the items for different substring matches,
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
the resulting parser will be of at least unbounded polynomial order in worst case.”
(Notation slightly modified by me to conform to notations elsewhere in my
parser packages).

Despite the shortcomings of the forest described by Grune & Jacobs, I won't
implement Scott's improvements. For practical use, the worst case spatial complexity
seems never to materialize. However, after gaining more insights in the future when
using the SPPF for more complex real word scenarios I will be prepared to reconsider.


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
// Symbol nodes will be searched by span (x…y), followed by a check of A
// (for RHS-nodes this will be modified slightly, see remarks below).
// Searching is implemented as a tree of height 2, with the edges labeled by input position
// and the leafs being sets of nodes. The tree is implemented by a map of maps of sets.
// We introduce a small helper type for it.
type searchTree map[uint64]map[uint64]*iteratable.Set // methods below

// Forest is a data structure for a "shared packed parse forest" (SPPF).
// A packed parse forest re-uses existing parse tree nodes between different
// parse trees. For a conventional non-ambiguous parse, a parse forest consists
// of a single tree. Ambiguous grammars, on the other hand, may result in parse
// runs where more than one parse tree is created. To save space these parse
// trees will share common nodes.
type Forest struct {
	symbolNodes searchTree                      // search tree for [A (x…y)] (type SymbolNode)
	rhsNodes    searchTree                      // search tree for RHSs, see type rhsNode
	orEdges     map[*SymbolNode]*iteratable.Set // or-edges from symbols to RHSs, indexed by symbol
	andEdges    map[*rhsNode]*iteratable.Set    // and-edges
	parent      map[*SymbolNode]*SymbolNode     // parent-edges
}

// NewForest returns an empty forest.
func NewForest() *Forest {
	return &Forest{
		symbolNodes: make(map[uint64]map[uint64]*iteratable.Set),
		rhsNodes:    make(map[uint64]map[uint64]*iteratable.Set),
		orEdges:     make(map[*SymbolNode]*iteratable.Set),
		andEdges:    make(map[*rhsNode]*iteratable.Set),
		parent:      make(map[*SymbolNode]*SymbolNode),
	}
}

// --- Exported Functions ----------------------------------------------------

// AddReduction adds a node for a reduced grammar rule into the forest.
// The extent of the reduction is derived from the RHS-nodes.
//
// If the RHS is void, nil is returned. For this case, clients should use
// AddEpsilonReduction instead.
func (f *Forest) AddReduction(sym *lr.Symbol, rule int, rhs []*SymbolNode) *SymbolNode {
	if len(rhs) == 0 {
		return nil
	}
	start := rhs[0].Extent.From()
	end := rhs[len(rhs)-1].Extent.To()
	rhsnode := f.addRHSNode(rule, rhs, rhs[0].Extent.From())
	f.addOrEdge(sym, rhsnode, start, end)
	for seq, d := range rhs {
		f.addAndEdge(rhsnode, d.Symbol, uint(seq), start, end)
		f.parent[d] = f.findSymNode(sym, start, end)
	}
	return f.findSymNode(sym, start, end)
}

// AddEpsilonReduction adds a node for a reduced ε-production.
func (f *Forest) AddEpsilonReduction(sym *lr.Symbol, rule int, pos uint64) *SymbolNode {
	rhsnode := f.addRHSNode(rule, []*SymbolNode{}, pos)
	f.addOrEdge(sym, rhsnode, pos, pos)
	return f.findSymNode(sym, pos, pos)
}

// AddTerminal adds a node for a recognized terminal into the forest.
func (f *Forest) AddTerminal(t *lr.Symbol, pos uint64) *SymbolNode {
	return f.addSymNode(t, pos, pos+1)
}

// --- Nodes -----------------------------------------------------------------

// SymbolNode represents a node in the parse forest, referencing a
// grammar symbol which has been reduced (Earley: completed).
type SymbolNode struct { // this is [A (x…y)]
	Symbol *lr.Symbol // A
	Extent lr.Span    // (x…y), i.e., positions in the input covered by this symbol
}

func makeSym(symbol *lr.Symbol) *SymbolNode {
	return &SymbolNode{Symbol: symbol}
}

// Use as makeSym(A).spanning(x, y), resulting in [A (x…y)]
func (sn *SymbolNode) spanning(from, to uint64) *SymbolNode {
	sn.Extent = lr.Span{from, to}
	return sn
}

func (sn *SymbolNode) String() string {
	return fmt.Sprintf("(%s %s)", sn.Symbol, sn.Extent.String())
}

// FindSymNode finds a (shared) node for a symbol node in the forest.
func (f *Forest) findSymNode(sym *lr.Symbol, start, end uint64) *SymbolNode {
	return f.symbolNodes.findSymbol(start, end, sym)
}

// addSymNode adds a symbol node to the forest. Returns a reference to a SymbolNode,
// which may already have been in the SPPF beforehand.
func (f *Forest) addSymNode(sym *lr.Symbol, start, end uint64) *SymbolNode {
	sn := f.findSymNode(sym, start, end)
	if sn == nil {
		sn = makeSym(sym).spanning(start, end)
		f.symbolNodes.Add(start, end, sn)
	}
	return sn
}

/*
RHS-Nodes

We are handling ambiguity by inserting multiple RHS nodes per symbol. Refer
to “Parsing Techniques” by  Dick Grune and Ceriel J.H. Jacobs
(https://dickgrune.com/Books/PTAPG_2nd_Edition/), Section 3.7.3.1
Combining Duplicate Subtrees:

“We […] combine all the duplicate subtrees in the forest […] by having only
one copy of a node labeled with a non-terminal A and spanning a given substring
of the input. If A produces that substring in more than one way, more than one
or-arrow will emanate from the OR-node labeled A, each pointing to an AND-node
labeled with a rule number. In this way the AND-OR-tree turns into a directed
acyclic graph […].

It is important to note that two OR-nodes (which represent right-hand sides of
rules) can only be combined if all members of the one node are the same as the
corresponding members of the other node.”

The last remark of Grune & Jacobs leads us to the question of identity of RHS-Nodes.
It is not enough to have [δ (x…y)] as a unique label; we need identity of every sub-symbol
(including its span) as well.

To avoid iterating repeatedly over the children we use a signature-function Σ to
encode the following information:

	let RHS = [δ1 (x1…y1)] [δ2 (x2…y2)] … [δn (xn…yn)]
	then Σ(RHS) := Σ(δ1, x1, δ2, x2, … , δn, xn)

and for a reduced ε-production RHS=[ε (x)]

	Σ(RHS) := Σ(x)

Thus instead of storing [δ (x…y)] as RHS-nodes, we store [δ (x) Σ] as unique
RHS-nodes.
*/

// Nodes [δ (x) Σ] in the parse forest.
type rhsNode struct {
	rule  int    // rule of which this RHS δ is from
	start uint64 // start position in the input
	sigma int32  // signature Σ of RHS children symbol nodes
}

func makeRHS(rule int) *rhsNode {
	return &rhsNode{rule: rule}
}

// Use as makeRHS(δ).identified(x, Σ), resulting in [δ (x) Σ]
func (rhs *rhsNode) identified(start uint64, signature int32) *rhsNode {
	rhs.start = start
	rhs.sigma = signature
	return rhs
}

// rhsSignature hashes over the symbols of a RHS, given a slice of symbols and
// a start position. The latter is used only in cases where RHS=ε.
// To randomize input positions, we map them to an array o of offsets.
var o = [...]int32{107, 401, 353, 223, 811, 569, 619, 173, 433, 757}

func rhsSignature(rhs []*SymbolNode, start uint64) int32 {
	const largePrime = int32(143743)
	if len(rhs) == 0 { // ε
		return o[start%uint64(len(o))]
	}
	h := int32(817)
	for _, symnode := range rhs {
		h *= int32(symnode.Symbol.Value) * o[symnode.Extent.From()%uint64(len(o))]
		h %= largePrime
	}
	return 0
}

// FindRHSNode finds a (shared) node for a right hand side in the forest.
func (f *Forest) findRHSNode(rule int, rhs []*SymbolNode, start uint64) *rhsNode {
	signature := rhsSignature(rhs, start)
	return f.rhsNodes.findRHS(start, rule, signature)
}

// addRHSNode adds a symbol node to the forest. Returns a reference to a rhsNode,
// which may already have been in the SPPF beforehand.
func (f *Forest) addRHSNode(rule int, rhs []*SymbolNode, start uint64) *rhsNode {
	node := f.findRHSNode(rule, rhs, start)
	if node == nil {
		signature := rhsSignature(rhs, start)
		node = makeRHS(rule).identified(start, signature)
		f.rhsNodes.Add(start, uint64(rule), node)
	}
	return node
}

// --- Edges -----------------------------------------------------------------

// orEdges are ambiguity forks in the parse forest.
type orEdge struct {
	fromSym *SymbolNode
	toRHS   *rhsNode
}

// addOrEdge inserts an edge between a symbol and a RHS.
// If start or end are not already contained in the forest, they are added.
//
// If the edge already exists, nothing is done.
func (f *Forest) addOrEdge(sym *lr.Symbol, rhs *rhsNode, start, end uint64) {
	T().Debugf("Add OR-edge %v ----> %v", sym, rhs.rule)
	sn := f.addSymNode(sym, start, end)
	if e := f.findOrEdge(sn, rhs); e.isNull() {
		e = orEdge{sn, rhs}
		if _, ok := f.orEdges[sn]; !ok {
			f.orEdges[sn] = iteratable.NewSet(0)
		}
		f.orEdges[sn].Add(rhs)
	}
}

// findOrEdge finds an or-edge starting from a symbol and pointing to an
// RHS-node. If none is found, nullOrEdge is returned.
func (f *Forest) findOrEdge(sn *SymbolNode, rhs *rhsNode) orEdge {
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
func (e orEdge) isNull() bool {
	return e == nullOrEdge
}

// An andEdge connects a RHS to the symbols it consists of.
type andEdge struct {
	fromRHS  *rhsNode    // RHS node starts the edge
	toSym    *SymbolNode // symbol node this edge points to
	sequence uint        // sequence number 0…n, used for ordering children
}

// addAndEdge inserts an edge between a RHS and a symbol, labeled with a seqence
// number. If start or end are not already contained in the forest, they are
// added. Note that it cannot happen that two edges between identical nodes
// exist for different sequence numbers. The function panics if such a condition
// is found.
//
// If the edge already exists, nothing is done.
func (f *Forest) addAndEdge(rhs *rhsNode, sym *lr.Symbol, seq uint, start, end uint64) {
	T().Debugf("Add AND-edge %v --(%d)--> %v", rhs.rule, seq, sym)
	sn := f.addSymNode(sym, start, end)
	if e := f.findAndEdge(rhs, sn); e.isNull() {
		e = andEdge{rhs, sn, seq}
		if _, ok := f.andEdges[rhs]; !ok {
			f.andEdges[rhs] = iteratable.NewSet(0)
		}
		f.andEdges[rhs].Add(e)
	} else if e.sequence != seq {
		panic(fmt.Sprintf("new edge with sequence=%d replaces sequence=%d", seq, e.sequence))
	}
}

// findAndEdge finds an and-edge starting from an RHS node and pointing to a
// symbol-node. If none is found, nullAndEdge is returned.
func (f *Forest) findAndEdge(rhs *rhsNode, sn *SymbolNode) andEdge {
	if edges := f.andEdges[rhs]; edges != nil {
		v := edges.FirstMatch(func(el interface{}) bool {
			e := el.(andEdge)
			return e.fromRHS == rhs && e.toSym == sn
		})
		if v == nil {
			return nullAndEdge
		}
		return v.(andEdge)
	}
	return nullAndEdge
}

// nullAndEdge denotes an and-edge that is not present in a graph.
var nullAndEdge = andEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e andEdge) isNull() bool {
	return e == nullAndEdge
}

// --- searchTree -----------------------------------------------------------------

// searchTree models a tree of height 3 with the edges labeled with p1 and p2.
// Semantics of (p1, p2) differ between symbol-nodes and RHS-nodes:
// For symbols, (p1, p2) = (start, end) input positions,
// for RHS, (p1, p2) = (start, rule).
// The leaf is a parse-forest node, either a symbol node or an RHS-node.
//
// find() searches the leaf at (p1, p2) and tests all nodes there for a given criteria.
// The search criteria is given as a predicate-function, returning true for a match.
func (t searchTree) find(p1, p2 uint64, predicate func(el interface{}) bool) interface{} {
	if t1, ok := t[p1]; ok {
		if t2, ok := t1[p2]; ok {
			return t2.FirstMatch(predicate)
		}
	}
	return nil
}

// find a symbol-node for (start, end, symbol).
func (t searchTree) findSymbol(from, to uint64, sym *lr.Symbol) *SymbolNode {
	node := t.find(from, to, func(el interface{}) bool {
		s := el.(*SymbolNode)
		return s.Symbol == sym
	})
	if node == nil {
		return nil
	}
	return node.(*SymbolNode)
}

// find an RHS-node for (position, rule-no, signature).
func (t searchTree) findRHS(start uint64, rule int, signature int32) *rhsNode {
	node := t.find(start, uint64(rule), func(el interface{}) bool {
		rhs := el.(*rhsNode)
		return rhs.sigma == signature
	})
	if node == nil {
		return nil
	}
	return node.(*rhsNode)
}

// Add adds an item as a leaf of the searchTree-path (p1, p2).
// Semantics of (p1, p2) differ between symbol-nodes and RHS-nodes:
// For symbols, (p1, p2) = (start, end) input positions,
// for RHS, (p1, p2) = (start, rule).
func (t searchTree) Add(p1, p2 uint64, item interface{}) {
	if t1, ok := t[p1]; !ok {
		t[p1] = make(map[uint64]*iteratable.Set)
		t[p1][p2] = iteratable.NewSet(0)
	} else if _, ok := t1[p2]; !ok {
		t[p1][p2] = iteratable.NewSet(0)
	}
	t[p1][p2].Add(item)
}
