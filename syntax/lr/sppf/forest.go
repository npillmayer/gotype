package sppf

/*
Code for SPPFs are rare, mostly found in academic papers. One of them
is "SPPF-Style Parsing From Earley Recognisers" by Elizabeth Scott
(https://reader.elsevier.com/reader/sd/pii/S1571066108001497?token=11DA10F6D3F3B941B251F0A0FB5CEFAE7CE954EEEFF5066D98928CFB16E01B8840108BA73F9D1DE7644B0CFD11F9DBC6).
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

type symNode struct {
	symbol *lr.Symbol
	span   lr.Span // positions in the input covered by this symbol
}

type rhsNode struct {
	rule int     // rule number this RHS is from
	span lr.Span // positions in the input covered by this RHS
}

func makeSym(symbol *lr.Symbol) symNode {
	return symNode{symbol: symbol}
}

// Use as makeSym(A).spanning(x, y), resulting in [A (x…y)]
func (sn *symNode) spanning(from, to uint64) *symNode {
	sn.span = lr.Span{from, to}
	return sn
}

func makeRHS(rule int) rhsNode {
	return rhsNode{rule: rule}
}

// Use as makeRHS(δ).spanning(x, y), resulting in [δ (x…y)]
func (rhs *rhsNode) spanning(from, to uint64) *rhsNode {
	rhs.span = lr.Span{from, to}
	return rhs
}

// addSymNode adds a symbol node to the forest. Returns a reference to a symNode,
// which may already have been in the SPPF beforehand.
func (f *Forest) addSymNode(sym *lr.Symbol, start, end uint64) *symNode {
	if pos, found := f.findSymNode(sym.Value, start, end); found {
		return pos
	}
	sn := makesym(sym).spanning(start, end)
	f.symNodes[sym.Value] = sn
	return len(f.symNodes) - 1
}

func (f *Forest) findSymNode(ID int, start, end uint64) (*symNode, bool) {
	f.symNodes[start][end].FirstMatch()
	for i, sn := range f.symNodes {
		if sn.symbol.Value == ID && sn.span.From() == start && sn.span.To() == end {
			return i, true
		}
	}
	return -1, false
}

// for now, just brute force. 1st optimization should be: search from back to front.
// If necessary, introduce a map.
func (f *Forest) findRHSNode(rule int, start, end uint64) (int, bool) {
	for i, rhs := range f.rhsNodes {
		if rhs.rule == rule && rhs.span.From() == start && rhs.span.To() == end {
			return i, true
		}
	}
	return -1, false
}

type orEdge struct {
	fromSym int // index of symNode
	toRHS   int // index of RHSNode
}

// // collection of edges. Made a struct to make it hashable.
// type orEdges struct {
// 	e []orEdge
// }

// nullOrEdge denotes an or-edge that is not present in a graph.
var nullOrEdge = orEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e orEdge) isNull() bool {
	return e == nullOrEdge
}

func (f *Forest) addOrEdge(sym *lr.Symbol, rule int, start, end uint64) {
	sn := f.addSymNode(sym, start, end)
	rhs := f.addRHSNode(rule, start, end)
	f.insertOrEdge(sn, rhs)
}

func (f *Forest) insertOrEdge(snpos, rhspos int) *orEdge {
	if e := f.findOrEdge(snpos, rhspos); e != nil {
		return e
	}
	e := &orEdge{fromSym: snpos, toRHS: rhspos}
	if _, ok := f.orEdges[snpos]; !ok {
		f.orEdges[snpos] = iteratable.NewSet(0)
	}
	f.orEdges[snpos].Add(e)
	return e
}

func (f *Forest) findOrEdge(snpos, rhspos int) *orEdge {
	if edges, ok := f.orEdges[snpos]; ok {
		// TODO introduce a function FirstMatch(predicate) for efficiency
		es := edges.Copy().Subset(func(el interface{}) bool {
			e := el.(*orEdge)
			return e.toRHS == rhspos
		})
		return asOrEdge(es.First())
	}
	return nil
}

func asOrEdge(edge interface{}) *orEdge {
	if e, ok := edge.(*orEdge); ok {
		return e
	}
	panic("orEdge cannot be converted from interface{}")
}

type andEdge struct {
	from     int // index of rhsNode
	toSym    int // index of symNode
	sequence int // sequence number 0…n, used for ordering children
}

// collection of edges. Made a struct to make it hashable.
type andEdges struct {
	edges []andEdge
}

// nullAndEdge denotes an and-edge that is not present in a graph.
var nullAndEdge = andEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e andEdge) isNull() bool {
	return e == nullAndEdge
}

func (f *Forest) addAndEdge(sym *lr.Symbol, rhs int, toSym *lr.Symbol, start, end uint64, seq int) {
	var from int
	if sym != nil {
		from = f.addSymNode(sym, start, end)
	} else {
		from = f.addRHSNode(rhs, start, end)
	}
	to := f.addSymNode(toSym, start, end)
	f.insertAndEdge(from, to, seq)
}

func (f *Forest) insertAndEdge(pos, snpos int, sequence int) *andEdge {
	if e := f.findAndEdge(pos, snpos); e != nil {
		return e
	}
	e := andEdge{from: pos, toSym: snpos, sequence: sequence}
	if edges, ok := f.andEdges[from]; ok {
		edges.edges = append(edges.edges, e)
		return &edges.edges[len(edges.edges)-1]
	}
	f.andEdges[from] = make(map[int]andEdges{from: e})
	return &f.andEdges[from][0]
}

func (f *Forest) findAndEdge(snfrom, snto int) *andEdge {
	if edges, ok := f.andEdges[snfrom]; ok {
		// TODO introduce a function FirstMatch(predicate) for efficiency
		es := edges.Copy().Subset(func(el interface{}) bool {
			e := el.(*andEdge)
			return e.toSym == snto
		})
		return asAndEdge(es.First())
	}
	return nil
}

func asAndEdge(edge interface{}) *andEdge {
	if e, ok := edge.(*andEdge); ok {
		return e
	}
	panic("andEdge cannot be converted from interface{}")
}

type edge interface {
	From() int
	To() int
}

func (f *Forest) endpoints(e edge) (int, int) {
	f, t := e.From(), e.To()
	if f < 0 {
		f = -f
	}
	return f, t
}

func (f *Forest) orEndpoints(e edge) (*lr.Symbol, *lr.Rule) {
	return nil, nil // TODO
}

// --------------------------------------------------------------------------------

/*
func (g *fbGraph) StartOfEdge(edge wEdge) *feasibleBreakpoint {
	if from, ok := g.nodes[edge.from]; ok {
		return from
	}
	return nil
}
*/

// Edges returns all known edges of a graph. if includePruned is true,
// deleted edges will be included.
/*
func (g *fbGraph) Edges(includePrunded bool) []wEdge {
	//edgesTo: map[int]map[int]map[int]wEdge
	var edges []wEdge
	for _, from := range g.edgesTo {
		for _, edgesDict := range from {
			for _, e := range edgesDict {
				// edge for line l
				edges = append(edges, e)
			}
		}
	}
	return edges
}
*/

// Breakpoint returns the feasible breakpoint at the given position if it exists in the graph,
// and nil otherwise.
/*
func (g *fbGraph) Breakpoint(position int) *feasibleBreakpoint {
	return g.nodes[position]
}
*/

// RemoveEdge removes the edge between two breakpoints for a linecount.
// The breakpoints are not deleted from the graph.
// If the edge does not exist, this is a no-op.
//
// Deleted edges are conserved and may be collected with g.Edges(true).
/*
func (g *fbGraph) RemoveEdge(from, to *feasibleBreakpoint, linecnt int) {
	if _, ok := g.nodes[from.mark.Position()]; !ok {
		return
	}
	if _, ok := g.nodes[to.mark.Position()]; !ok {
		return
	}
	if edgesFrom, ok := g.edgesTo[to.mark.Position()]; ok {
		if edges, ok := edgesFrom[from.mark.Position()]; ok {
			if e, ok := edges[linecnt]; ok { // edge exists => move it to prunedEdges
				if t, ok := g.prunedEdges[to.mark.Position()]; ok { // 'to' dict exists
					edges := t[from.mark.Position()]
					if edges == nil {
						edges = make(map[int]wEdge)
						t[from.mark.Position()] = edges
					}
					edges[linecnt] = e
				} else {
					edges := map[int]wEdge{linecnt: e}
					g.prunedEdges[to.mark.Position()] = map[int]map[int]wEdge{from.mark.Position(): edges}
					//g.prunedEdges[to.mark.Position()][from.mark.Position()][linecnt] = e
				}
			}
			delete(edges, linecnt)
			if len(edges) == 0 {
				delete(edgesFrom, to.mark.Position())
			}
		}
	}
}
*/

// AddEdge adds a weighted edge from one node to another. Endpoints which are
// not yet contained in the graph are added.
// Does nothing if from=to.
/*
func (g *fbGraph) AddEdge(from, to *feasibleBreakpoint, cost int32, total int32, linecnt int) {
	if from.mark.Position() == to.mark.Position() {
		return
	}
	if g.Breakpoint(to.mark.Position()) == nil {
		g.Add(to)
	}
	if g.Edge(from, to, linecnt).isNull() {
		edge := newWEdge(from, to, cost, total, linecnt)
		if t, ok := g.edgesTo[to.mark.Position()]; ok {
			edges := t[from.mark.Position()]
			if edges == nil {
				edges = make(map[int]wEdge)
				t[from.mark.Position()] = edges
			}
			edges[linecnt] = edge
		} else {
			edges := map[int]wEdge{linecnt: edge}
			g.edgesTo[to.mark.Position()] = map[int]map[int]wEdge{from.mark.Position(): edges}
		}
	}
}
*/

/*
var noBreakpoints []*feasibleBreakpoint

func (g *fbGraph) EdgesTo(fb *feasibleBreakpoint) edgesDict {
	return edgesDict{
		g:         g,
		edgesFrom: g.edgesTo[fb.mark.Position()],
	}
}

type edgesDict struct {
	g         *fbGraph
	edgesFrom map[int]map[int]wEdge
}

func (dict edgesDict) SelectFrom(fb *feasibleBreakpoint, linecnt int) wEdge {
	if edges, ok := dict.edgesFrom[fb.mark.Position()]; ok {
		if edge, ok := edges[linecnt]; ok {
			return edge
		}
	}
	return nullEdge
}

func (dict edgesDict) WithLabel(linecnt int) []wEdge {
	var r []wEdge
	for _, edges := range dict.edgesFrom {
		if edge, ok := edges[linecnt]; ok {
			r = append(r, edge)
		}
	}
	return r
}
*/

// To returns all breakpoints in g that can reach directly to a breakpoint given by
// a position. The returned breakpoints are sorted by position.
/*
func (g *fbGraph) To(fb *feasibleBreakpoint) []*feasibleBreakpoint {
	position := fb.mark.Position()
	if _, ok := g.edgesTo[position]; !ok || len(g.edgesTo[position]) == 0 {
		return noBreakpoints
	}
	breakpoints := make([]*feasibleBreakpoint, len(g.edgesTo[position]))
	i := 0
	for pos := range g.edgesTo[position] {
		breakpoints[i] = g.nodes[pos]
		i++
	}
	sort.Sort(breakpointSorter{breakpoints})
	return breakpoints
}

type breakpointSorter struct {
	breakpoints []*feasibleBreakpoint
}

func (s breakpointSorter) Len() int {
	return len(s.breakpoints)
}

func (s breakpointSorter) Less(i, j int) bool {
	fb1 := s.breakpoints[i]
	fb2 := s.breakpoints[j]
	return fb1.mark.Position() < fb2.mark.Position()
}

func (s breakpointSorter) Swap(i, j int) {
	s.breakpoints[i], s.breakpoints[j] = s.breakpoints[j], s.breakpoints[i]
}
*/

// Cost returns the cost for the edge between between two breakpoints,
// valid for a linecount.
// If from and to are the same node or if there is no edge (from,to),
// a pseudo-label with infinite cost is returned.
// Cost returns true if an edge (from,to) exists, false otherwise.
/*
func (g *fbGraph) Cost(from, to *feasibleBreakpoint, linecnt int) (int32, bool) {
	if from == to {
		return linebreak.InfinityDemerits, false
	}
	if edgesFrom, ok := g.edgesTo[to.mark.Position()]; ok {
		if edges, ok := edgesFrom[from.mark.Position()]; ok {
			if edge, ok := edges[linecnt]; ok {
				return edge.cost, true
			}
		}
	}
	return linebreak.InfinityDemerits, false
}
*/
