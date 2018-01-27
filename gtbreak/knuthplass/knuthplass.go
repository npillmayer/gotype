/*
Line breaking algorithm described by D.E. Knuth and M.F. Plass (early draft).

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
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	sg "gonum.org/v1/gonum/graph/simple"
)

const (
	BoxType int8 = iota
	GlueType
	KernType
	DiscretionaryType
)

const InfinityDemerits int32 = 10000
const InfinityMerits int32 = -10000

func beadtype(t int8) string {
	switch t {
	case BoxType:
		return "box"
	case GlueType:
		return "glue"
	case KernType:
		return "kern"
	case DiscretionaryType:
		return "discr"
	}
	return "<unknown>"
}

func demerits(d int32) int32 {
	if d > InfinityDemerits {
		d = InfinityDemerits
	}
	return d
}

// --- Beads -----------------------------------------------------------------

type Bead interface {
	BeadType() int8
	Text() string
	Width() WSS                    // width, w-shrink, w+stretch
	Dimens() (int64, int64, int64) // width, height, depth
}

// --- Beading ---------------------------------------------------------------

// at = nil : start
type Beading interface {
	GetCursor(at BeadingCursor) BeadingCursor
}

type BeadingCursor interface {
	GetBead() Bead
	Advance() bool
	ID() int64
}

// ---------------------------------------------------------------------------

type LL func(int64) int64 // get line length at line n

type WSS struct {
	W   int64
	Min int64
	Max int64
}

// --- Horizon (active Nodes) ------------------------------------------------

type activeFeasibleBreakpoints struct {
	*hashset.Set
	values  []interface{}
	iterinx int
}

func newActiveFeasibleBreakpoints() *activeFeasibleBreakpoints {
	set := hashset.New()
	h := &activeFeasibleBreakpoints{set, nil, -1}
	return h
}

func (h *activeFeasibleBreakpoints) first() *FeasibleBreakpoint {
	var f *FeasibleBreakpoint
	fmt.Printf("horizon: there are %d active FBs\n", h.Size())
	if h.Size() > 0 {
		h.values = h.Values()
		f = h.values[0].(*FeasibleBreakpoint)
		h.iterinx = 1
	}
	return f
}

func (h *activeFeasibleBreakpoints) next() *FeasibleBreakpoint {
	var f *FeasibleBreakpoint
	if h.values != nil && h.iterinx < len(h.values) {
		f = h.values[h.iterinx].(*FeasibleBreakpoint)
		h.iterinx++
	}
	return f
}

func (h *activeFeasibleBreakpoints) append(f *FeasibleBreakpoint) {
	h.Add(f)
}

func (h *activeFeasibleBreakpoints) remove(f *FeasibleBreakpoint) {
	h.Remove(f)
}

// ---------------------------------------------------------------------------

func w(wss WSS) (int64, int64, int64) {
	return wss.W, wss.Min, wss.Max
}

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

type FeasibleBreakpoint struct {
	graph.Node
	LineNum   int64
	marker    BeadingCursor
	fragment  WSS   // sum of widths up to last bead iteration
	totalcost int32 // sum of costs for segments up to this breakpoint
}

func (h *FeasibleBreakpoint) String() string {
	if h.marker == nil {
		return "<para-start>"
	}
	return fmt.Sprintf("<break l.%d @ %v>", h.LineNum, h.marker.GetBead())
}

// internal
func (kp *KnuthPlassLinebreaker) newFeasibleBreakpoint(id int64) *FeasibleBreakpoint {
	node := sg.Node(id)
	fb := &FeasibleBreakpoint{node, 1, nil, WSS{}, 0}
	kp.AddNode(fb)
	return fb
}

func (kp *KnuthPlassLinebreaker) breakpointAfterBead(marker BeadingCursor) *FeasibleBreakpoint {
	fb := kp.newFeasibleBreakpoint(marker.ID())
	fb.marker = marker
	return fb
}

func (kp *KnuthPlassLinebreaker) findBreakpointAtMarker(marker BeadingCursor) *FeasibleBreakpoint {
	var fb *FeasibleBreakpoint
	if marker != nil {
		if kp.Has(sg.Node(marker.ID())) {
			node := kp.Node(marker.ID())
			fmt.Printf("found existing feas. bp = %v\n", node)
			if node != nil {
				fb = node.(*FeasibleBreakpoint)
			}
		}
	}
	return fb
}

type feasibleSegment struct { // an edge between nodes of feasible breakpoints
	graph.WeightedEdge
	Totals WSS
}

func (fseg *feasibleSegment) String() string {
	if fseg == nil {
		return "<--no edge-->"
	} else {
		return fmt.Sprintf("<-%v--(%.1f)--%v->", fseg.From(), fseg.Weight(), fseg.To())
	}
}

func (kp *KnuthPlassLinebreaker) feasibleLineBetween(from, to *FeasibleBreakpoint,
	cost int32, totals WSS, prune bool) *feasibleSegment {
	//
	predecs := kp.To(to)
	var seg *feasibleSegment
	var p *FeasibleBreakpoint
	mintotal := InfinityDemerits
	if prune && len(predecs) > 0 {
		if len(predecs) > 1 {
			fmt.Printf("breakpoint %v has %d predecessors", to, len(predecs))
			panic("breakpoint (with pruning) has more than one predecessor")
		}
		p = predecs[0].(*FeasibleBreakpoint)
		w, _ := kp.Weight(p, to)
		mintotal = p.totalcost + int32(w)
	}
	if from.totalcost+cost < mintotal {
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

/*
Calculate the cost of a breakpoint. A breakpoint may result either in being
infeasible (demerits >= infinity) or having a positive (demerits) or negative
(merits) cost/benefit.

TODO: the cost/badness function should be deleted to a strategy delegate (or
a chain of delegate => Holkner). Maybe refactor: skeleton algorithm in
central package, & different strategy delegates in sub-packages. K&P is a
special delegate with the TeX strategy formalized.

Question: Is the box-glue-model part of the central algorithm? Or is it
already a strategy (component) ?
*/
func (fb *FeasibleBreakpoint) calculateCostTo(bead Bead, linelen int64) (int32, bool) {
	w, min, max := w(bead.Width())
	fb.fragment.W += w
	fb.fragment.Min += min
	fb.fragment.Max += max
	stillreachable := true
	var d int32 = InfinityDemerits
	if fb.fragment.W <= linelen {
		if fb.fragment.Max >= linelen {
			d = int32(linelen * 100 / (fb.fragment.Max - fb.fragment.W))
			//fmt.Printf("w < l: demerits = %d\n", d)
		}
	} else if fb.fragment.W >= linelen {
		if fb.fragment.Min <= linelen { // compressible enough?
			d = int32(linelen * 100 / (fb.fragment.W - fb.fragment.Min))
			//fmt.Printf("w > l: demerits = %d\n", d)
		} else { // will not fit any more
			stillreachable = true
		}
	}
	return demerits(d), stillreachable
}

func (kp *KnuthPlassLinebreaker) FindBreakpoints(prune bool) (int, []*FeasibleBreakpoint) {
	fb := kp.newFeasibleBreakpoint(-1) // without bead
	kp.horizon.append(fb)              // begin of paragraph is first "active node"
	var cursor BeadingCursor = kp.beading.GetCursor(nil)
	var bead Bead          // will hold the current bead of kp.beading
	for cursor.Advance() { // loop over beading
		bead = cursor.GetBead() // will shift horizon one bead further
		fmt.Printf("next bead is %s\n", bead)
		fb = kp.horizon.first() // loop over feas. breakpoints of horizon
		for fb != nil {         // while there are active f. breakpoints in horizon n
			//fmt.Printf("fb in horizon = %v\n", fb)
			cost, stillreachable := fb.calculateCostTo(bead, 30)
			if cost < InfinityDemerits { // is breakpoint and is feasible
				fmt.Printf("create feasible breakpoint at %v\n", cursor)
				newfb := kp.findBreakpointAtMarker(cursor)
				if newfb == nil { // feas. breakpoint not yet existent
					newfb = kp.breakpointAfterBead(cursor) // create a f. breakpoint
				}
				edge := kp.feasibleLineBetween(fb, newfb, cost, fb.fragment, prune)
				fmt.Printf("feasible line = %v\n", edge)
				kp.horizon.append(newfb) // make new fb member of horizon n+1
			} else if !stillreachable {
				kp.horizon.remove(fb) // no longer valid in horizon
			}
			fb = kp.horizon.next()
		}
	}
	return kp.collectFeasibleBreakpoints(cursor.ID())
}

func (kp *KnuthPlassLinebreaker) collectFeasibleBreakpoints(last int64) (int, []*FeasibleBreakpoint) {
	//var optimalBreaks []*FeasibleBreakpoint
	fmt.Printf("collecting breakpoints, backwards from #%d\n", last)
	fb := kp.Node(last)
	stack := arraystack.New()
	var breakpoints []*FeasibleBreakpoint = make([]*FeasibleBreakpoint, 20)
	for fb != nil {
		stack.Push(fb)
		predecs := kp.To(fb)
		if predecs == nil || len(predecs) == 0 {
			fb = nil
		} else {
			fmt.Printf("node #%d has %d predecessors\n", fb.ID(), len(predecs))
			fb = predecs[0]
		}
	}
	if stack.Size() > 0 {
		fmt.Printf("optimal paragraph breaking uses %d breakpoints\n", stack.Size())
		p, ok := stack.Pop()
		for ok {
			breakpoints = append(breakpoints, p.(*FeasibleBreakpoint))
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
