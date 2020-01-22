package sppf

/*

Some code in this file is loosely beased on ideas from the great
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
)

// fbGraph is a directed weighted graph of feasible breakpoints.
// Its implementation is inspired by the great Gonum packages. However, Gonum
// is in some respects too restrictive and in others too much of an overkill
// for our needs. Moreover, it panics in certain error conditions.
// We taylor it heavily to fit our specific needs.
//
// A node in the graph refers to a numeric position within an input text.
// The text is represented by a khipu (see package khipu), which is something
// like a TeX hlist, i.e. a list of boxes, glue, penalties etc.  These kind
// of list-items are called knots. Positions are indices of knots.
//
// A single position can be reached optimally by exactly one segment (path
// through breakpoints). However, for reasons explained by Knuth/Plass it is
// advantageous in some situations to permit for more than one segment, if
// they result in different line-counts. This allows in effect to defer the
// optimality-decision until the target for the last line is known.
//
// We implement this feed-forward by labelling the edges with a tuple
// (cost, linecount). A breakpoint may be reached by more than one edge
// if either cost or linecount differ.
type fbGraph struct {
	nodes map[int]*feasibleBreakpoint
	//from    map[int]map[int][]wEdge
	edgesTo     map[int]map[int]map[int]wEdge // edge to, from, with linecount
	prunedEdges map[int]map[int]map[int]wEdge // we conserve deleted edges
}

// newFBGraph returns a fbGraph with the specified self and absent
// edge weight values.
func newFBGraph() *fbGraph {
	return &fbGraph{
		nodes:       make(map[int]*feasibleBreakpoint),
		edgesTo:     make(map[int]map[int]map[int]wEdge),
		prunedEdges: make(map[int]map[int]map[int]wEdge),
	}
}

type wEdge struct {
	from, to  int // this is an edge between two text-positions
	cost      int32
	total     int32
	linecount int
}

// nullEdge denotes an edge that is not present in a graph.
var nullEdge = wEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e wEdge) isNull() bool {
	return e == nullEdge
}

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

// Add adds a feasible breakpoint to the graph.
// It returns an error if the breakpoint is already present.
func (g *fbGraph) Add(fb *feasibleBreakpoint) error {
	T().Debugf("Added new breakpoint at %d/%v", fb.mark.Position(), fb.mark.Knot())
	if _, exists := g.nodes[fb.mark.Position()]; exists {
		return fmt.Errorf("Breakpoint at position %d already known", fb.mark.Position())
	}
	g.nodes[fb.mark.Position()] = fb
	return nil
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
