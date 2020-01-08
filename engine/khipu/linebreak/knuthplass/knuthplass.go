/*
Package knuthplass implements (in an early draft) a
line breaking algorithm described by D.E. Knuth and M.F. Plass.

Definite source of information is of course

   Computers & Typesetting, Vol. A & C.
   http://www-cs-faculty.stanford.edu/~knuth/abcde.html

A good summary may be found in

   http://defoe.sourceforge.net/folio/knuth-plass.html

Horizon hält nicht nur die aktiven Knoten für die "optimalen"
Umbrüche, sondern auch für Varianten.
Was ist "optimal"? Jeder Durchlauf bekommt ein Strategy-Objekt mit, das
in einer Kette von Delegates einen Beitrag zur Bestimmung leistet (eher
nach Holkner). Diese Delegates bestimmen die Demerits.

Beim Aufbrechen der Zeilen ist der Folio-Artikel ziemich aufschlussreich:
wieweit muss man eigentlich glue usw explizit machen? Wie weit nimmt einem
Unicode das ab?

Folio kennt folgende Objekte:
- Breakpoint (possible und feasible)
- Segment (f. Breakpoint, sub-path bis zurück zu paragraph start)
- Line (zwischen 2 f. Breakpoints, mit key (lineno, break1, break2))

*/
package knuthplass

import (
	"fmt"
	"io/ioutil"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/khipu"
	"github.com/npillmayer/gotype/engine/khipu/linebreak"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	sg "gonum.org/v1/gonum/graph/simple"
)

// --- Beads -----------------------------------------------------------------

type Bead interface {
	BeadType() int8
	Text() string
	Width() linebreak.WSS          // width, w-shrink, w+stretch
	Dimens() (int64, int64, int64) // width, height, depth
}

// --- Horizon (active Nodes) ------------------------------------------------

type activeFeasibleBreakpoints struct {
	*hashset.Set               // a set of feasible breakpoints
	values       []interface{} // holds breakpoints during iteration
	iterinx      int           // current iteration index
}

// constructor
func newActiveFeasibleBreakpoints() *activeFeasibleBreakpoints {
	set := hashset.New()
	afb := &activeFeasibleBreakpoints{set, nil, -1}
	return afb
}

// first starts iteration over the feasible breakpoints of the current horizon.
func (h *activeFeasibleBreakpoints) first() *feasibleBreakpoint {
	var fb *feasibleBreakpoint // return value
	fmt.Printf("horizon: there are %d active FBs\n", h.Size())
	if h.Size() > 0 {
		h.values = h.Values() // get set values as list, unsorted
		fb = h.values[0].(*feasibleBreakpoint)
		h.iterinx = 1
	}
	return fb
}

// next gets the next feasible breakpoints of the current horizon.
func (h *activeFeasibleBreakpoints) next() *feasibleBreakpoint {
	var fb *feasibleBreakpoint
	if h.values != nil && h.iterinx < len(h.values) {
		fb = h.values[h.iterinx].(*feasibleBreakpoint)
		h.iterinx++
	}
	return fb
}

func (h *activeFeasibleBreakpoints) append(fb *feasibleBreakpoint) {
	h.Add(fb)
}

func (h *activeFeasibleBreakpoints) remove(fb *feasibleBreakpoint) {
	h.Remove(fb)
}

// ---------------------------------------------------------------------------

type KnuthPlassLinebreaker struct {
	*sg.WeightedDirectedGraph
	beading    Beading
	lineLength LL
	horizon    *activeFeasibleBreakpoints
}

func NewKPLinebreaker() *KnuthPlassLinebreaker {
	kp := &KnuthPlassLinebreaker{sg.NewWeightedDirectedGraph(
		float64(InfinityDemerits), float64(InfinityDemerits)), nil, nil,
		newActiveFeasibleBreakpoints()}
	return kp
}

type feasibleBreakpoint struct {
	graph.Node
	LineNum   int64
	marker    BeadingCursor
	fragment  linebreak.WSS // sum of widths up to last bead iteration
	totalcost int32         // sum of costs for segments up to this breakpoint
}

func (fb *feasibleBreakpoint) String() string {
	if fb.marker == nil {
		return "<para-start>"
	}
	return fmt.Sprintf("<break l.%d @ %v>", fb.LineNum, fb.marker.GetBead())
}

// internal
func (kp *KnuthPlassLinebreaker) newFeasibleBreakpoint(id int64) *feasibleBreakpoint {
	node := sg.Node(id)
	fb := &feasibleBreakpoint{node, 1, nil, WSS{}, 0}
	kp.AddNode(fb)
	return fb
}

// newBreakpointAfterKnot creates a breakpoint at the given cursor position.
func (kp *KnuthPlassLinebreaker) newBreakpointAfterKnot(marker BeadingCursor) *feasibleBreakpoint {
	fb := kp.newFeasibleBreakpoint(marker.ID())
	fb.marker = marker
	return fb
}

func (kp *KnuthPlassLinebreaker) findBreakpointAtMarker(marker BeadingCursor) *feasibleBreakpoint {
	if marker == nil {
		return nil
	}
	var fb *feasibleBreakpoint
	if kp.Has(sg.Node(marker.ID())) {
		node := kp.Node(marker.ID())
		fmt.Printf("found existing feas. bp = %v\n", node)
		if node != nil {
			fb = node.(*feasibleBreakpoint)
		}
	}
	return fb
}

type feasibleSegment struct { // an edge between nodes of feasible breakpoints
	graph.WeightedEdge
	Totals WSS // the cumulative widths of this segment
}

func (fseg *feasibleSegment) String() string {
	if fseg == nil {
		return "<--no edge-->"
	}
	return fmt.Sprintf("<-%v--(%.1f)--%v->", fseg.From(), fseg.Weight(), fseg.To())
}

// feasibleLineBetween possibly creates a segment between two given breakpoints.
//
// If parameter prune is true, the segment is constructed and compared to
// any existing segments. If its demerits are better than the exising one, the
// new segment replaces the old one.
// (just one segment between the two breakpoints can exist).
//
// If prune is false, the segment is constructed and added to the existing ones,
// if any (more than one segment between the two breakpoints can co-exist).
func (kp *KnuthPlassLinebreaker) feasibleLineBetween(from, to *feasibleBreakpoint,
	cost int32, totals WSS, prune bool) *feasibleSegment {
	//
	predecs := kp.To(to)
	var seg *feasibleSegment  // return value
	var p *feasibleBreakpoint // predecessor breakpoint
	mintotal := InfinityDemerits
	if prune && len(predecs) > 0 { // then get existing demerits
		if len(predecs) > 1 {
			fmt.Printf("breakpoint %v has %d predecessors", to, len(predecs))
			panic("breakpoint (with pruning) has more than one predecessor")
		}
		p = predecs[0].(*feasibleBreakpoint)
		w, _ := kp.Weight(p, to)
		mintotal = p.totalcost + int32(w)
	}
	if from.totalcost+cost < mintotal { // new line is cheaper
		edge := kp.NewWeightedEdge(from, to, float64(cost))
		seg = &feasibleSegment{edge, totals}
		kp.SetWeightedEdge(seg)
		if prune && p != nil {
			kp.RemoveEdge(kp.WeightedEdge(p, to))
		}
	}
	return seg
}

// === Algorithms ============================================================

// Calculate the cost of a breakpoint. A breakpoint may result either in being
// infeasible (demerits >= infinity) or having a positive (demerits) or negative
// (merits) cost/benefit.
//
// TODO: the cost/badness function should be deleted to a strategy delegate (or
// a chain of delegate => Holkner). Maybe refactor: skeleton algorithm in
// central package, & different strategy delegates in sub-packages. K&P is a
// special delegate with the TeX strategy formalized.
//
// Question: Is the box-glue-model part of the central algorithm? Or is it
// already a strategy (component) ?
func (fb *feasibleBreakpoint) calculateCostTo(knot khipu.Knot, linelen dimen.Dimen) (int32, bool) {
	fb.fragment.SetFromKnot(knot)
	stillreachable := true
	var d = InfinityDemerits
	if fb.fragment.W <= linelen {
		if fb.fragment.Max >= linelen { // TODO really?
			d = int32(linelen * 100 / (fb.fragment.Max - fb.fragment.W))
			//fmt.Printf("w < l: demerits = %d\n", d)
		}
	} else if fb.fragment.W >= linelen {
		if fb.fragment.Min <= linelen { // compressible enough?
			d = int32(linelen * 100 / (fb.fragment.W - fb.fragment.Min))
			//fmt.Printf("w > l: demerits = %d\n", d)
		} else { // will not fit any more
			stillreachable = false
		}
	}
	return demerits(d), stillreachable
}

func (kp *KnuthPlassLinebreaker) FindBreakpoints(input *khipu.Khipu, prune bool) (int, []*feasibleBreakpoint) {
	fb := kp.newFeasibleBreakpoint(-1) // without bead
	kp.horizon.append(fb)              // begin of paragraph is first "active node"
	cursor := input.Iterator()
	//var knot khipu.Knot // will hold the current knot of the input khipu
	for cursor.Next() { // loop over input
		knot := cursor.Knot() // will shift horizon one knot further
		//fmt.Printf("next knot is %s\n", khipu.KnotString(knot))
		fb = kp.horizon.first() // loop over feasible breakpoints of horizon
		for fb != nil {         // while there are active breakpoints in horizon n
			//fmt.Printf("fb in horizon = %v\n", fb)
			cost, stillreachable := fb.calculateCostTo(bead, 30) // TODO linelength dynamic
			if cost < InfinityDemerits {                         // is breakpoint and is feasible
				//fmt.Printf("create feasible breakpoint at %v\n", cursor)
				newfb := kp.findBreakpointAtMarker(cursor)
				if newfb == nil { // feas. breakpoint not yet existent
					newfb = kp.newBreakpointAfterKnot(cursor) // create a feasible breakpoint
				}
				edge := kp.feasibleLineBetween(fb, newfb, cost, fb.fragment, prune)
				//fmt.Printf("feasible line = %v\n", edge)
				kp.horizon.append(newfb) // make new fb member of horizon n+1
			} else if !stillreachable {
				kp.horizon.remove(fb) // no longer valid in horizon
			}
			fb = kp.horizon.next()
		}
	}
	return kp.collectFeasibleBreakpoints(cursor.ID())
}

func (kp *KnuthPlassLinebreaker) collectFeasibleBreakpoints(last int64) (int, []*feasibleBreakpoint) {
	//var optimalBreaks []*feasibleBreakpoint
	fmt.Printf("collecting breakpoints, backwards from #%d\n", last)
	fb := kp.Node(last)
	stack := arraystack.New() // for reversing the breakpoint order, TODO this is overkill
	var breakpoints = make([]*feasibleBreakpoint, 20)
	for {
		stack.Push(fb)
		predecs := kp.To(fb)
		if predecs == nil || len(predecs) == 0 { // at start node
			break
		}
		fmt.Printf("node #%d has %d predecessors\n", fb.ID(), len(predecs))
		fb = predecs[0] // with pruning, there should only be one
	}
	if stack.Size() > 0 {
		fmt.Printf("optimal paragraph breaking uses %d breakpoints\n", stack.Size())
		p, ok := stack.Pop()
		for ok {
			breakpoints = append(breakpoints, p.(*feasibleBreakpoint))
			p, ok = stack.Pop()
		}
	}
	return len(breakpoints), breakpoints
}

func (kp *KnuthPlassLinebreaker) MarshalToDotFile(id string) {
	dot, err := dot.Marshal(kp, "linebreaks", "", "", false)
	if err != nil {
		fmt.Printf("mashalling error: %v", err.Error())
	} else {
		ioutil.WriteFile(fmt.Sprintf("./kp_graph_%s.dot", id), dot, 0644)
	}
}
