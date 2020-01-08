package knuthplass

import (
	"fmt"
	"sort"

	"github.com/npillmayer/gotype/engine/khipu/linebreak"

	"github.com/emirpasic/gods/sets/hashset"
)

// fbGraph is a directed weighted graph of feasible breakpoints.
// Its implementation is inspired by the great Gonum packages. However, Gonum
// is in some respects too restrictive and in others too much of an overkill.
// Moreover, it panics in certain error conditions.
// We taylor it heavily to fit our specific needs.
//
// Gonum-license:
// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of Gonum source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
type fbGraph struct {
	nodes   map[int]*feasibleBreakpoint
	from    map[int]map[int]wEdge
	to      map[int]map[int]wEdge
	nodeIDs *hashset.Set
}

// newFBGraph returns a fbGraph with the specified self and absent
// edge weight values.
func newFBGraph() *fbGraph {
	return &fbGraph{
		nodes:   make(map[int]*feasibleBreakpoint),
		from:    make(map[int]map[int]wEdge),
		to:      make(map[int]map[int]wEdge),
		nodeIDs: hashset.New(),
	}
}

type wEdge struct {
	from, to int   // this is an edge between from and to
	cost     int32 // demerits for breaking here
}

// nullEdge denotes an edge that is not present in a graph.
var nullEdge = wEdge{}

// isNull checks if an edge is null, i.e. non-existent.
func (e wEdge) isNull() bool {
	return e == nullEdge
}

// newWEdge returns a new weighted edge from one breakpoint to another.
func (g *fbGraph) newWEdge(from, to *feasibleBreakpoint, cost int32, wss linebreak.WSS) wEdge {
	return wEdge{
		from: from.mark.Position(),
		to:   to.mark.Position(),
		cost: cost,
	}
}

// Add adds a feasible breakpoint to the graph.
// It returns an error if the breakpoint is already present.
func (g *fbGraph) Add(fb *feasibleBreakpoint) error {
	if _, exists := g.nodes[fb.mark.Position()]; exists {
		return fmt.Errorf("Breakpoint at position %d already known", fb.mark.Position())
	}
	g.nodes[fb.mark.Position()] = fb
	g.nodeIDs.Add(fb.mark.Position())
	return nil
}

// Edge returns the edge (from,to), if such an edge exists,
// otherwise it returns nullEdge.
// The to-node must be directly reachable from the from-node.
func (g *fbGraph) Edge(from, to int) wEdge {
	edge, ok := g.from[from][to]
	if !ok {
		return nullEdge
	}
	return edge
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

// RemoveEdge removes the edge between two breakpoints.
// The breakpoints are not deleted from the graph.
// If the edge does not exist, this is a no-op.
func (g *fbGraph) RemoveEdge(from, to int) {
	if _, ok := g.nodes[from]; !ok {
		return
	}
	if _, ok := g.nodes[to]; !ok {
		return
	}
	delete(g.from[from], to)
	delete(g.to[to], from)
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

// AddEdge adds a weighted edge from one node to another. If the nodes do not exist, they are added
// and are set to the endpoints of the edge otherwise.
func (g *fbGraph) AddSegment(from, to *feasibleBreakpoint, cost int32, totals linebreak.WSS) {
	edge := newWEdge()
	if from. == e.to {
		return
	}
	if g.Breakpoint(from) == nil {
		g.Add(from)
	} else {
		g.nodes[fid] = from
	}
	if _, ok := g.nodes[tid]; !ok {
		g.AddNode(to)
	} else {
		g.nodes[tid] = to
	}

	if fm, ok := g.from[fid]; ok {
		fm[tid] = e
	} else {
		g.from[fid] = map[int]graph.wEdge{tid: e}
	}
	if tm, ok := g.to[tid]; ok {
		tm[fid] = e
	} else {
		g.to[tid] = map[int]graph.wEdge{fid: e}
	}
}

var noBreakpoints []*feasibleBreakpoint

// To returns all breakpoints in g that can reach directly to a breakpoint given by
// a position. The returned breakpoints are sorted by position.
func (g *fbGraph) To(position int) []*feasibleBreakpoint {
	if _, ok := g.to[position]; !ok || len(g.to[position]) == 0 {
		return noBreakpoints
	}
	breakpoints := make([]*feasibleBreakpoint, len(g.to[position]))
	i := 0
	for pos := range g.to[position] {
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

// Cost returns the cost for the edge between between two breakpoints.
// If from and to are the same node or if there is no edge (from,to),
// in infinite cost ist returned.
// Cost returns true if an edge (from,to) exists, false otherwise.
func (g *fbGraph) Cost(from, to int) (int32, bool) {
	if from == to {
		return linebreak.InfinityDemerits, false
	}
	if edges, ok := g.from[from]; ok {
		if e, ok := edges[to]; ok {
			return e.cost, true
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
