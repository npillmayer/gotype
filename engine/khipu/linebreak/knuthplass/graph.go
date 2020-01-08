package knuthplass

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
// Gonum-license:
// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of Gonum source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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
	edgesTo map[int]map[int]map[int]wEdge // edge to, from, with linecount
}

// newFBGraph returns a fbGraph with the specified self and absent
// edge weight values.
func newFBGraph() *fbGraph {
	return &fbGraph{
		nodes: make(map[int]*feasibleBreakpoint),
		//from:    make(map[int]map[int][]wEdge),
		edgesTo: make(map[int]map[int]map[int]wEdge),
	}
}

type wEdge struct {
	from, to  int // this is an edge between two text-positions
	cost      int32
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
func newWEdge(from, to *feasibleBreakpoint, cost int32, linecnt int) wEdge {
	return wEdge{
		from:      from.mark.Position(),
		to:        to.mark.Position(),
		cost:      cost,
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

func (g *fbGraph) EdgeFrom(edge wEdge) *feasibleBreakpoint {
	if from, ok := g.nodes[edge.from]; ok {
		return from
	}
	return nil
}

// Edges returns all the edges in the graph.
/*
func (g *fbGraph) Edges() graph.Edges {
	var edges []graph.Edge
	for _, u := range g.nodes {
		for _, e := range g.from[u.ID()] {
			edges = append(edges, e)
		}
	}
	if len(edges) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedEdges(edges)
}
*/

// From returns all nodes in g that can be reached directly from n.
/*
func (g *fbGraph) From(id int) graph.Nodes {
	if _, ok := g.from[id]; !ok {
		return graph.Empty
	}

	from := make([]graph.Node, len(g.from[id]))
	i := 0
	for vid := range g.from[id] {
		from[i] = g.nodes[vid]
		i++
	}
	if len(from) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedNodes(from)
}
*/

// HasEdgeBetween returns whether an edge exists between nodes x and y without
// considering direction.
/*
func (g *fbGraph) HasEdgeBetween(xid, yid int) bool {
	if _, ok := g.from[xid][yid]; ok {
		return true
	}
	_, ok := g.from[yid][xid]
	return ok
}
*/

// hasEdge returns true if an edge (from,to) exists in the graph.
/*
func (g *fbGraph) hasEdge(from, to int) bool {
	if _, ok := g.from[from][to]; !ok {
		return false
	}
	return true
}
*/

// NewNode returns a new unique Node to be added to g. The Node's ID does
// not become valid in g until the Node is added to g.
/*
func (g *fbGraph) NewNode() graph.Node {
	if len(g.nodes) == 0 {
		return Node(0)
	}
	if int(len(g.nodes)) == uid.Max {
		panic("simple: cannot allocate node: no slot")
	}
	return Node(g.nodeIDs.NewID())
}
*/

// Breakpoint returns the feasible breakpoint at the given position if it exists in the graph,
// and nil otherwise.
func (g *fbGraph) Breakpoint(position int) *feasibleBreakpoint {
	return g.nodes[position]
}

// Nodes returns all the nodes in the graph.
/*
func (g *fbGraph) Nodes() graph.Nodes {
	if len(g.nodes) == 0 {
		return graph.Empty
	}
	nodes := make([]graph.Node, len(g.nodes))
	i := 0
	for _, n := range g.nodes {
		nodes[i] = n
		i++
	}
	return iterator.NewOrderedNodes(nodes)
}
*/

// RemoveEdge removes the edge between two breakpoints for a linecount.
// The breakpoints are not deleted from the graph.
// If the edge does not exist, this is a no-op.
func (g *fbGraph) RemoveEdge(from, to *feasibleBreakpoint, linecnt int) {
	if _, ok := g.nodes[from.mark.Position()]; !ok {
		return
	}
	if _, ok := g.nodes[to.mark.Position()]; !ok {
		return
	}
	//delete(g.from[from.mark.Position()], to.mark.Position())
	if edgesFrom, ok := g.edgesTo[to.mark.Position()]; ok {
		if edges, ok := edgesFrom[from.mark.Position()]; ok {
			delete(edges, linecnt)
			if len(edges) == 0 {
				delete(edgesFrom, to.mark.Position())
			}
		}
	}
}

// RemoveNode removes the node with the given ID from the graph, as well as any edges attached
// to it. If the node is not in the graph it is a no-op.
/*
func (g *fbGraph) RemoveNode(id int) {
	if _, ok := g.nodes[id]; !ok {
		return
	}
	delete(g.nodes, id)

	for from := range g.from[id] {
		delete(g.to[from], id)
	}
	delete(g.from, id)

	for to := range g.to[id] {
		delete(g.from[to], id)
	}
	delete(g.to, id)

	g.nodeIDs.Release(id)
}
*/

// AddSegment adds a weighted edge from one node to another. Endpoints which are
// not yet contained in the graph are added.
// Does nothing if from=to.
func (g *fbGraph) AddSegment(from, to *feasibleBreakpoint, cost int32, linecnt int) {
	if from.mark.Position() == to.mark.Position() {
		return
	}
	// if g.Breakpoint(from.mark.Position()) == nil {
	// 	g.Add(from)
	// }
	if g.Breakpoint(to.mark.Position()) == nil {
		g.Add(to)
	}
	if g.Edge(from, to, linecnt).isNull() {
		edge := newWEdge(from, to, cost, linecnt)
		// if f, ok := g.from[from.mark.Position()]; ok {
		// 	f[to.mark.Position()] = edge
		// } else {
		// 	edges := []wEdge{edge}
		// 	g.from[from.mark.Position()] = map[int][]wEdge{to.mark.Position(): edges}
		// }
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

// wEdge returns the weighted edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
/*
func (g *fbGraph) wEdge(uid, vid int) wEdge {
	edge, ok := g.from[uid][vid]
	if !ok {
		return nil
	}
	return edge
}
*/

// wEdges returns all the weighted edges in the graph.
/*
func (g *fbGraph) wEdges() graph.wEdges {
	var edges []graph.wEdge
	for _, u := range g.nodes {
		for _, e := range g.from[u.ID()] {
			edges = append(edges, e)
		}
	}
	if len(edges) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedwEdges(edges)
}
*/
