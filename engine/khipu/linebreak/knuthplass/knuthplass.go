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
	"github.com/npillmayer/gotype/engine/khipu"
	"github.com/npillmayer/gotype/engine/khipu/linebreak"
	"gonum.org/v1/gonum/graph/encoding/dot"
)

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
	//fmt.Printf("horizon: there are %d active FBs\n", h.Size())
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

type linebreaker struct {
	*fbGraph
	//*sg.WeightedDirectedGraph
	//beading  Beading
	parshape linebreak.Parshaper
	measure  linebreak.GlyphMeasure
	horizon  *activeFeasibleBreakpoints
}

func newLinebreaker(parshape linebreak.Parshaper, measure linebreak.GlyphMeasure) *linebreaker {
	kp := &linebreaker{}
	kp.fbGraph = newFBGraph()
	kp.parshape = parshape
	kp.measure = measure
	kp.horizon = newActiveFeasibleBreakpoints()
	return kp
}

type feasibleBreakpoint struct {
	//graph.Node
	lineno    int            // line number
	mark      linebreak.Mark // location of this breakpoint
	fragment  linebreak.WSS  // sum of widths up to last knot
	totalcost int32          // sum of costs for segments up to this breakpoint
}

type provisionalMark int // provisional mark from an integer position

func (m provisionalMark) Position() int    { return int(m) }
func (m provisionalMark) Knot() khipu.Knot { return nil }

func (fb *feasibleBreakpoint) String() string {
	if fb.mark == nil {
		return "<para-start>"
	}
	return fmt.Sprintf("<break l.%d @ %v>", fb.lineno, fb.mark.Knot())
}

// internal constructor
/*
func (kp *linebreaker) newFeasibleBreakpoint(position int) *feasibleBreakpoint {
	mark := provisionalMark(position)
	fb := &feasibleBreakpoint{
		lineno:    1,
		mark:      mark,
		fragment:  linebreak.WSS{},
		totalcost: 0,
	}
	kp.Add(fb)
	return fb
}
*/

// newBreakpointAfterKnot creates a breakpoint at the given cursor position.
func (kp *linebreaker) newBreakpointAfterKnot(mark linebreak.Mark) *feasibleBreakpoint {
	fb := &feasibleBreakpoint{
		lineno:    1,
		mark:      mark,
		fragment:  linebreak.WSS{},
		totalcost: 0,
	}
	kp.Add(fb)
	return fb
}

func (kp *linebreaker) findBreakpointAtMark(mark linebreak.Mark) *feasibleBreakpoint {
	if mark == nil {
		return nil
	}
	if fb := kp.Breakpoint(mark.Position()); fb != nil {
		//node := kp.Node(int64(mark.Position()))
		//fmt.Printf("found existing feas. bp = %v\n", node)
		//if node != nil {
		//	fb = node.(*feasibleBreakpoint)
		//}
		return fb
	}
	return nil
}

type feasibleSegment struct { // an edge between nodes of feasible breakpoints
	wEdge                // edge with cost
	totals linebreak.WSS // the cumulative widths of this segment
}

func (fseg *feasibleSegment) String() string {
	if fseg == nil {
		return "+--no edge-->"
	}
	return fmt.Sprintf("+-%d--(%d)--%d->", fseg.from, fseg.cost, fseg.to)
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
func (kp *linebreaker) feasibleLineBetween(from, to *feasibleBreakpoint,
	cost int32, totals linebreak.WSS, prune bool) *feasibleSegment {
	//
	predecessors := kp.To(to.mark.Position())
	var seg *feasibleSegment  // return value
	var p *feasibleBreakpoint // predecessor breakpoint position
	mintotal := linebreak.InfinityDemerits
	if prune && len(predecessors) > 0 { // then get existing demerits
		if len(predecessors) > 1 {
			fmt.Printf("breakpoint %v has %d predecessors", to, len(predecessors))
			panic("breakpoint (with pruning) has more than one predecessor")
		}
		p = predecessors[0]
		c := kp.Edge(p.mark.Position(), to.mark.Position()).cost
		//w, _ := kp.Weight(p, to)
		mintotal = p.totalcost + c
	}
	if from.totalcost+cost < mintotal { // new line is cheaper
		segment := kp.newWEdge(from, to, cost, totals)
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
func (fb *feasibleBreakpoint) calculateCostTo(knot khipu.Knot, parshape linebreak.Parshaper) (int32, bool) {
	fb.fragment.SetFromKnot(knot)
	stillreachable := true
	var d = linebreak.InfinityDemerits
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
	return linebreak.CapDemerits(d), stillreachable
}

//func (kp *linebreaker) FindBreakpoints(input *khipu.Khipu, prune bool) (int, []*feasibleBreakpoint) {

func FindBreakpoints(input *khipu.Khipu, cursor linebreak.Cursor, parshape linebreak.Parshaper,
	measure linebreak.GlyphMeasure, prune bool) (int, []linebreak.Mark) {
	//
	kp := newLinebreaker(parshape, measure)
	fb := kp.newFeasibleBreakpoint(-1) // without bead
	kp.horizon.append(fb)              // begin of paragraph is first "active node"
	//cursor := input.Iterator()
	//var knot khipu.Knot // will hold the current knot of the input khipu
	for cursor.Next() { // loop over input
		knot := cursor.Knot() // will shift horizon one knot further
		//fmt.Printf("next knot is %s\n", khipu.KnotString(knot))
		fb = kp.horizon.first() // loop over feasible breakpoints of horizon
		for fb != nil {         // while there are active breakpoints in horizon n
			//fmt.Printf("fb in horizon = %v\n", fb)
			cost, stillreachable := fb.calculateCostTo(knot, parshape) // TODO linelength dynamic
			if cost < linebreak.InfinityDemerits {                     // is breakpoint and is feasible
				//fmt.Printf("create feasible breakpoint at %v\n", cursor)
				newfb := kp.findBreakpointAtMark(cursor.Mark())
				if newfb == nil { // breakpoint not yet existent
					newfb = kp.newBreakpointAfterKnot(cursor.Mark()) // create a feasible breakpoint
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
	return kp.collectFeasibleBreakpoints(cursor.Mark())
}

func (kp *linebreaker) collectFeasibleBreakpoints(last linebreak.Mark) (int, []linebreak.Mark) {
	//var optimalBreaks []*feasibleBreakpoint
	//fmt.Printf("collecting breakpoints, backwards from #%d\n", last)
	fb := kp.findBreakpointAtMark(last)
	if fb == nil {
		fb = kp.newBreakpointAfterKnot(last)
	}
	stack := arraystack.New() // for reversing the breakpoint order, TODO this is overkill
	var breakpoints = make([]linebreak.Mark, 20)
	for {
		stack.Push(fb.mark)
		predecessors := kp.To(fb)
		if predecessors == nil || len(predecessors) == 0 { // at start node
			break
		}
		//fmt.Printf("node #%d has %d predecessors\n", fb.ID(), len(predecessors))
		fb = predecessors[0] // with pruning, there should only be one
	}
	if stack.Size() > 0 {
		fmt.Printf("optimal paragraph breaking uses %d breakpoints\n", stack.Size())
		p, ok := stack.Pop()
		for ok {
			breakpoints = append(breakpoints, p.(linebreak.Mark))
			p, ok = stack.Pop()
		}
	}
	return len(breakpoints), breakpoints
}

func (kp *linebreaker) MarshalToDotFile(id string) {
	dot, err := dot.Marshal(kp, "linebreaks", "", "", false)
	if err != nil {
		fmt.Printf("mashalling error: %v", err.Error())
	} else {
		ioutil.WriteFile(fmt.Sprintf("./kp_graph_%s.dot", id), dot, 0644)
	}
}
