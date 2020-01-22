package sppf

/*

Some code in this file is very loosely based on ideas from the great
Gonum project.

Gonum-License

Copyright ©2014 The Gonum Authors. All rights reserved.
Use of Gonum source code is governed by a BSD-style
license that can be found in the LICENSES file.

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
	"sort"

	"github.com/npillmayer/gotype/engine/khipu/linebreak"
	"github.com/npillmayer/gotype/syntax/lr"
)

// Forest implements a Shared Packed Parse Forest (SPPF).
// A packed parse forest re-uses existing parse tree nodes between different
// parse trees. For a conventional non-ambiguous parse, a parse forest consists
// of a single tree. Ambiguous grammars, on the other hand, may result in parse
// runs where more than one parse tree is created. To save space these parse
// trees will share common nodes.
//
// In our implementation there are two kinds of nodes: Symbol nodes and RHS-nodes.
// Symbol nodes reflect LHSs of rules, which are the results of a reduce action.
// RHS-nodes reflect RHSs of rules, which have been used to reduce to a Symbol node.
type Forest struct {
	symNodes map[int]SymNode
	rhsNodes map[int]RHSNode
	orEdges  map[int]orEdges  // or-edges from symbols to RHSs, indexed by symbol
	andEdges map[int]andEdges // and-edges
}

// NewForest returns a Forest with the specified self and absent
// edge weight values.
func NewForest() *Forest {
	return &Forest{
		symNodes: make(map[int]SymNode),
		rhsNodes: make(map[int]RHSNode),
		orEdges:  make(map[int]orEdge),
		andEdges: make(map[int]andEdge),
	}
}

type span [2]uint64

func (s *span) from() uint64 {
	return s[0]
}

func (s *span) to() uint64 {
	return s[1]
}

func (sn symNode) spanning(from, to uint64) symNode {
	return symNode{
		symbols: sn.symbol,
		span:    span{from, to},
	}
}

func (rhs rhsNode) spanning(from, to uint64) rhsNode {
	return rhsNode{
		rule: rhs.rule,
		span: span{from, to},
	}
}

func sym(symbol lr.Symbol) symNode {
	return symNode{symbol: symbol}
}

func rhs(rule int) rhsNode {
	return rhsNode{rule: rule}
}

type symNode struct {
	symbol lr.Symbol
	span   uint64 // positions in the input covered by this symbol
}

type rhsNode struct {
	rule int    // rule number this RHS is from
	span uint64 // positions in the input covered by this symbol
}

// addSymNode adds a symbol node to the forest. Returns a position in the nodes-table
// (which may already have been occupied beforehand,
// as it may already have been present in the SPPF.)
func (f *Forest) addSymNode(sym lr.Symbol, start, end uint64) int {
	if pos, found := f.findSymNode(sym.ID, start, end); found {
		return pos
	}
	sn = sym(sym).spanning(start, end)
	f.symNodes = append(f.symNodes, sn)
	return len(f.symNodes) - 1
}

// for now, just brute force. 1st optimization should be: search from back to front.
// If necessary, introduce a map.
func (f *Forest) findSymNode(ID int, start, end uint64) (int, bool) {
	for i, sn := range f.symNodes {
		if sn.symbol.GetID() == ID && sn.span.from() == start && sn.span.to() == end {
			return i, true
		}
	}
	return -1, false
}

// for now, just brute force. 1st optimization should be: search from back to front.
// If necessary, introduce a map.
func (f *Forest) findRHSNode(rule int, start, end uint64) (int, bool) {
	for i, rhs := range f.RHSNodes {
		if rhs.rule == rule && sn.span.from() == start && sn.span.to() == end {
			return i, true
		}
	}
	return -1, false
}

type orEdge struct {
	fromSym int // index of symNode
	toRHS   int // index of RHSNode
}

// collection of edges. Made a struct to make it hashable.
type orEdges struct {
	edges []orEdge
}

// We introduce a small space optimization: and-edges may occur between
// sym-->sym or rhs-->sym. If the origin is a RHS, we store the negative positional
// index of the rhsNode, if it is a symbol, we simply store the index positive.
// This helps avoid storing an index and a <none>.
type andEdge struct {
	from     int // index of symNode or -index of rhsNode
	toSym    int // index of symNode
	sequence int // sequence number 0..n, used for ordering children
}

// collection of edges. Made a struct to make it hashable.
type andEdges struct {
	edges []andEdge
}

// nullOrEdge denotes an or-edge that is not present in a graph.
var nullOrEdge = orEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e orEdge) isNull() bool {
	return e == nullOrEdge
}

// nullAndEdge denotes an and-edge that is not present in a graph.
var nullAndOrEdge = orEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e andEdge) isNull() bool {
	return e == nullAndEdge
}

func (f *Forest) addOrEdge(sym lr.Symbol, rule int, start, end uint64) {
	sn := f.addSymNode(sym, start, end)
	rhs := f.addRHSNode(rule, start, end)
	insertOrEdge(sn, rhs)
}

func (f *Forest) insertOrEdge(snpos, rhspos int) *orEdge {
	if e := f.findOrEdge(snpos, rhspos); e != nil {
		return e
	}
	e := orEdge{fromSym: snpos, toRHS: rhspos}
	if edges, ok := f.orEdges[snpos]; ok {
		edges.edges = append(edges.edges, e)
		return &edges.edges[len(edges.edges)-1]
	}
	f.orEdges[snpos] = make(map[int]orEdges{snpos: e})
	return &f.orEdges[snpos][0]
}

func (f *Forest) findOrEdge(snpos, rhspos int) *orEdge {
	if edges, ok := f.orEdges[snpos]; ok {
		for i, e := range edges.edges {
			if e.toRHS == rhspos {
				return &edges.edges[i]
			}
		}
	}
	return nil
}

// --------------------------------------------------------------------------------

// newWEdge returns a new weighted edge from one breakpoint to another,
// given two breakpoints and a label-key.
// It is not yet inserted into a graph.
func newWEdge(from, to *feasibleBreakpoint, cost int32, total int32, linecnt int) wEdge {
	if from.books[linecnt-1] == nil {
		panic(fmt.Errorf("startpoint of new line %d seems to have incorrent books: %v", linecnt, from))
	}
	if to.books[linecnt] == nil {
		panic(fmt.Errorf("endpoint of new line %d seems to have incorrent books: %v", linecnt, to))
	}
	return wEdge{
		from:      from.mark.Position(),
		to:        to.mark.Position(),
		cost:      cost,
		total:     total,
		linecount: linecnt,
	}
}

// Edge returns the edge (from,to), if such an edge exists,
// otherwise it returns nullEdge.
// The to-node must be directly reachable from the from-node.
func (g *fbGraph) Edge(from, to *feasibleBreakpoint, linecnt int) wEdge {
	edges, ok := g.edgesTo[to.mark.Position()][from.mark.Position()]
	if !ok {
		return nullEdge
	}
	if edge, ok := edges[linecnt]; ok {
		return edge
	}
	return nullEdge
}

func (g *fbGraph) StartOfEdge(edge wEdge) *feasibleBreakpoint {
	if from, ok := g.nodes[edge.from]; ok {
		return from
	}
	return nil
}

// Edges returns all known edges of a graph. if includePruned is true,
// deleted edges will be included.
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
	if includePrunded {
		for _, from := range g.prunedEdges {
			for _, edgesDict := range from {
				for _, e := range edgesDict {
					// edge for line l
					edges = append(edges, e)
				}
			}
		}
	}
	return edges
}

// Breakpoint returns the feasible breakpoint at the given position if it exists in the graph,
// and nil otherwise.
func (g *fbGraph) Breakpoint(position int) *feasibleBreakpoint {
	return g.nodes[position]
}

// RemoveEdge removes the edge between two breakpoints for a linecount.
// The breakpoints are not deleted from the graph.
// If the edge does not exist, this is a no-op.
//
// Deleted edges are conserved and may be collected with g.Edges(true).
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

// AddEdge adds a weighted edge from one node to another. Endpoints which are
// not yet contained in the graph are added.
// Does nothing if from=to.
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

// To returns all breakpoints in g that can reach directly to a breakpoint given by
// a position. The returned breakpoints are sorted by position.
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

// Cost returns the cost for the edge between between two breakpoints,
// valid for a linecount.
// If from and to are the same node or if there is no edge (from,to),
// a pseudo-label with infinite cost is returned.
// Cost returns true if an edge (from,to) exists, false otherwise.
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
