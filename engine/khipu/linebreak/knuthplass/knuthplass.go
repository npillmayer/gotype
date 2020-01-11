package knuthplass

/*
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

import (
	"fmt"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/khipu"
	"github.com/npillmayer/gotype/engine/khipu/linebreak"
)

// linebreaker is an internal entity for K&P-linebreaking.
type linebreaker struct {
	*fbGraph
	//*sg.WeightedDirectedGraph
	//beading  Beading
	parshape linebreak.Parshape
	horizon  *activeFeasibleBreakpoints
}

func newLinebreaker(parshape linebreak.Parshape) *linebreaker {
	kp := &linebreaker{}
	kp.fbGraph = newFBGraph()
	kp.parshape = parshape
	kp.horizon = newActiveFeasibleBreakpoints()
	return kp
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

// --- Breakpoints -----------------------------------------------------------

// A feasible breakpoint is uniquely identified by a text position (mark).
// A break position may be selectable for different line counts, and we
// retain all of them. Different line-count paths usually will have different costs.
// We will hold some bookkeeping information to reflect active segments.
type feasibleBreakpoint struct {
	mark  khipu.Mark           // location of this breakpoint
	books map[int]*bookkeeping // bookkeeping per linecount
}

type bookkeeping struct {
	segment      linebreak.WSS // sum of widths from this breakpoint up to current knot
	totalcost    int32         // sum of costs for segment up to this breakpoint
	startDiscard linebreak.WSS // sum of discardable space at start of segment / line
	breakDiscard linebreak.WSS // sum of discardable space while lookinf for next breakpoint
	hasContent   bool          // does this segment contain non-discardable item?
}

type provisionalMark int // provisional mark from an integer position

func (m provisionalMark) Position() int    { return int(m) }
func (m provisionalMark) Knot() khipu.Knot { return khipu.Penalty(-10000) }

func (fb *feasibleBreakpoint) String() string {
	if fb.mark == nil {
		return "<para-start>"
	}
	return fmt.Sprintf("<brk %d/%v>", fb.mark.Position(), fb.mark.Knot())
}

// func (fb *feasibleBreakpoint) initBook(linecnt int) {
// 	if fb.books == nil {
// 		fb.books = make(map[int]*bookkeeping)
// 	}
// 	fb.books[linecnt] = &bookkeeping{}
// }

func (fb *feasibleBreakpoint) UpdateSegment(linecnt int, diff linebreak.WSS) {
	if fb.books == nil {
		fb.books = make(map[int]*bookkeeping)
	}
	segment := linebreak.WSS{}
	total := int32(0)
	book, ok := fb.books[linecnt]
	if ok {
		segment = book.segment
		total = book.totalcost
	}
	fb.books[linecnt] = &bookkeeping{
		segment:   segment.Add(diff),
		totalcost: total,
	}
}

func (fb *feasibleBreakpoint) UpdateSegmentBookkeeping(mark khipu.Mark) {
	wss := linebreak.WSS{}.SetFromKnot(mark.Knot()) // get dimensions of knot
	for _, book := range fb.books {
		book.segment = book.segment.Add(wss)
		if book.hasContent {
			if mark.Knot().IsDiscardable() {
				book.breakDiscard = book.breakDiscard.Add(wss)
			} else {
				book.breakDiscard = linebreak.WSS{}
			}
		} else {
			if mark.Knot().IsDiscardable() {
				book.startDiscard = book.startDiscard.Add(wss)
			} else {
				book.hasContent = true
			}
		}
		T().Debugf("extending segment to %v", book.segment)
	}
}

// func (fb *feasibleBreakpoint) ClearBook(linecnt int) {
// 	if fb.books == nil {
// 		return
// 	}
// 	delete(fb.books, linecnt)
// }

func (fb *feasibleBreakpoint) Book(linecnt int) *bookkeeping {
	b, ok := fb.books[linecnt]
	if !ok {
		fb.books[linecnt] = &bookkeeping{}
	}
	return b
}

// newBreakpointAtMark creates a breakpoint at the given cursor position.
func (kp *linebreaker) newBreakpointAtMark(mark khipu.Mark) *feasibleBreakpoint {
	fb := &feasibleBreakpoint{
		mark:  mark,
		books: make(map[int]*bookkeeping),
	}
	kp.Add(fb)
	return fb
}

func (kp *linebreaker) findBreakpointAtMark(mark khipu.Mark) *feasibleBreakpoint {
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

func (kp *linebreaker) findPredecessorsWithLinecount(fb *feasibleBreakpoint, linecnt int) []*feasibleBreakpoint {
	var predecessors []*feasibleBreakpoint
	edges := kp.EdgesTo(fb).WithLabel(linecnt)
	for _, edge := range edges {
		if edge.isNull() {
			panic("this should not happen!") // TODO remove this
		}
		if edge.linecount == linecnt {
			predecessors = append(predecessors, kp.EdgeFrom(edge))
		}
	}
	return predecessors
}

// --- Segments ---------------------------------------------------------

func (kp *linebreaker) newFeasibleLine(fb *feasibleBreakpoint, mark khipu.Mark, cost int32,
	linecnt int, prune bool) *feasibleBreakpoint {
	//
	newfb := kp.findBreakpointAtMark(mark)
	if newfb == nil { // breakpoint not yet existent
		newfb = kp.newBreakpointAtMark(mark) // create a new feasible breakpoint
	}
	targettotal := fb.Book(linecnt-1).totalcost + cost // total cost of new line
	if !prune || kp.isCheapestSurvivor(newfb, targettotal, linecnt, true) {
		kp.AddEdge(fb, newfb, cost, linecnt)
		newfb.books[linecnt] = &bookkeeping{
			totalcost: targettotal,
			// segment:      linebreak.WSS{},
			// startDiscard: linebreak.WSS{},
			// breakDiscard: linebreak.WSS{},
		}
		T().Debugf("new line established %v --> %v", fb, newfb)
	}
	return newfb
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
//
func (kp *linebreaker) isCheapestSurvivor(fb *feasibleBreakpoint, totalcost int32,
	linecnt int, deleteOthers bool) bool {
	//
	// func (kp *linebreaker) feasibleLineBetween(from, to *feasibleBreakpoint,
	// 	cost int32, totals linebreak.WSS, prune bool) *feasibleSegment {
	//
	//predecessors := kp.To(to)
	//var seg *feasibleSegment  // return value
	var predecessor *feasibleBreakpoint    // predecessor breakpoint position
	mintotal := linebreak.InfinityDemerits // pre-set
	//
	// CONNECT A BREAKPOINT AND KEEP THE LEAST COST (if pruning)
	// target BP may or may not have an edge, but is part of the graph
	//if prune { // then look for existing demertics
	if pp := kp.findPredecessorsWithLinecount(fb, linecnt); pp != nil {
		if len(pp) > 1 {
			panic("breakpoint (with pruning) has more than one predecessor[line]")
		}
		predecessor = pp[0]
		// isolate total cost of predecessor for segment to fb
		predCost := kp.Edge(predecessor, fb, linecnt).cost
		mintotal = predecessor.books[linecnt].totalcost + predCost
	}
	//}
	// book, ok := from.Book(linecnt)
	// if !ok {
	// 	panic("totalcost of horizon node must be set") // TODO remove
	// }
	// if book.totalcost+cost < mintotal { // new line is cheaper
	if totalcost < mintotal { // new line is cheaper
		//edge := newWEdge(from, to, cost, linecnt)
		//seg = &feasibleSegment{edge, totals}
		// kp.AddSegment(from, to, cost, linecnt)
		// totalcost := book.totalcost + cost
		// wss := linebreak.WSS{}.SetFromKnot(to.mark.Knot()) // dimensions of next knot
		// segment := book.segment.Add(wss)
		// to.UpdateBook(linecnt, segment, totalcost)
		if deleteOthers && predecessor != nil {
			kp.RemoveEdge(predecessor, fb, linecnt)
		}
		return true
	}
	return false
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
func (fb *feasibleBreakpoint) calculateCostsTo(knot khipu.Knot, parshape linebreak.Parshape) (
	map[int]int32, bool) {
	//
	T().Infof("### calculateCostsTo(%v)", knot)
	var costs = make(map[int]int32) // linecount => cost, i.e. costs for different line targets
	//var wss = linebreak.WSS{}.SetFromKnot(knot) // dimensions of next knot
	//T().Debugf(" width of %v is %.2f bp", knot, knot.W().Points())
	cannotReachIt := 0
	//for linecnt, bookkeeping := range fb.books {
	for linecnt := range fb.books {
		T().Debugf(" ## checking cost at linecnt=%d", linecnt)
		d := linebreak.InfinityDemerits             // pre-set result variable
		linelen := parshape.LineLength(linecnt + 1) // length of line to fit into
		// wss := linebreak.WSS{}.SetFromKnot(knot)    // dimensions of next knot, should be 0
		// segwss := fb.Book(linecnt).segment.Add(wss) // widths of segment including knot
		segwss := fb.Book(linecnt).segment
		//segwss := fb.calculateSegmentWidth(knot, linecnt)
		T().Debugf("    +---%.2f--->    | %.2f", segwss.W.Points(), linelen.Points())
		if segwss.W <= linelen { // natural width less than line-length
			if segwss.Max >= linelen { // segment can stretch enough
				d = calculateDemerits(segwss, linelen-segwss.W, 0)
			} else { // segment is just too short
				// try with tolerance - misnomer, used otherweise in TeX
				tolerance := 3 // TODO from typesetting parameters; 1 = rigid
				stretchedwss := segwss.Copy()
				stretchedwss.Max = dimen.Dimen(tolerance) * (segwss.Max - segwss.W)
				if stretchedwss.Max >= linelen { // now segment can stretch enough
					d = calculateDemerits(stretchedwss, linelen-segwss.W, tolerance)
				}
			}
		} else { // natural width larger than line-length
			if segwss.Min <= linelen { // segment can shrink enough
				d = calculateDemerits(segwss, segwss.W-linelen, 0)
			} else { // segment will not fit any more
				// TeX has no tolerance for shrinking. Good?
				// TODO introduce overfull-hbox break here? d slightly smaller than infinity?
				cannotReachIt++
			}
		}
		costs[linecnt] = linebreak.CapDemerits(d)
		T().Debugf(" ## cost for line %d would be %s", linecnt+1, demeritsString(costs[linecnt]))
	}
	stillreachable := (cannotReachIt < len(fb.books))
	T().Debugf("### costs to %v is %v, reachable is %v", knot, costs, stillreachable)
	return costs, stillreachable
}

// func (fb *feasibleBreakpoint) calculateSegmentWidth(knot khipu.Knot, linecnt int) linebreak.WSS {
// 	var segwss linebreak.WSS
// 	wss := linebreak.WSS{}.SetFromKnot(knot) // dimensions of next knot
// 	segwss = fb.Book(linecnt).segment.Add(wss)
// 	return segwss
// }

func calculateDemerits(segwss linebreak.WSS, stretch dimen.Dimen, tolerance int) int32 {
	tolerancepenalty := 1000 // TODO from typesetting parameters
	tolerancepenalty *= tolerance - 1
	return 200 // TODO
}

func demeritsString(d int32) string {
	if d >= linebreak.InfinityDemerits {
		return "\u221e"
	} else if d <= linebreak.InfinityMerits {
		return "-\u221e"
	}
	return fmt.Sprintf("%d", d)
}

// --- Main API ---------------------------------------------------------

// FindBreakpoints is the main client API.
func FindBreakpoints(cursor linebreak.Cursor, parshape linebreak.Parshape, prune bool) (
	int, map[int][]khipu.Mark) {
	//
	kp := newLinebreaker(parshape)
	fb := kp.newBreakpointAtMark(provisionalMark(-1)) // start of paragraph
	fb.books[0] = &bookkeeping{}
	kp.horizon.append(fb) // this is the first "active node"
	var last khipu.Mark   // will hold last position within input khipu
	for cursor.Next() {   // loop over input
		last = cursor.Mark()
		T().Infof("_____________________________________________")
		T().Infof("_______________ %d/%s ___________________", last.Position(), last.Knot())
		fb = kp.horizon.first() // loop over feasible breakpoints of horizon
		if fb == nil {
			T().Errorf("no more active breakpoints, but input available")
			T().Errorf("this should probably have produces an overfull hbox")
		}
		for fb != nil { // while there are active breakpoints in horizon n
			stillreachable := true
			T().Infof("                %d/%v  (in horizon)", fb.mark.Position(), fb.mark.Knot())
			fb.UpdateSegmentBookkeeping(cursor.Mark())
			// Breakpoints are only allowed at penalties
			if cursor.Mark().Knot().Type() == khipu.KTPenalty { // TODO discretionaries
				var costs map[int]int32 // cost per linecnt-alternative
				costs, stillreachable = fb.calculateCostsTo(cursor.Knot(), parshape)
				for linecnt, cost := range costs {
					if cost < linebreak.InfinityDemerits { // new breakpoint is feasible
						newfb := kp.newFeasibleLine(fb, cursor.Mark(), cost, linecnt+1, prune)
						kp.horizon.append(newfb) // make new fb member of horizon n+1
					}
				}
			}
			if !stillreachable {
				kp.horizon.remove(fb) // no longer valid in horizon
			}
			fb = kp.horizon.next()
		}
	}
	return kp.collectFeasibleBreakpoints(last)
}

func (kp *linebreaker) collectFeasibleBreakpoints(last khipu.Mark) (int, map[int][]khipu.Mark) {
	// Collecting breakpoints, backwards from last
	fb := kp.findBreakpointAtMark(last)
	if fb == nil {
		panic("last breakpoint not found") // khipu didn't end with penalty -10000
		// for now panic, for debugging purposes
		//fb = kp.newBreakpointAtMark(last)
	}
	//stack := arraystack.New() // for reversing the breakpoint order, TODO this is overkill
	var breakpoints = make(map[int][]khipu.Mark) // list of breakpoints per linecount-variant
	for linecnt := range fb.books {
		l := linecnt
		//stack.Push(fb.mark)
		lines := make([]khipu.Mark, 0, 20)
		lines = append(lines, fb.mark)
		predecessors := kp.To(fb)
		for len(predecessors) > 0 { // while not at start node
			l-- // searching for predecessor with linecount-1
			var pred *feasibleBreakpoint
			for _, pred = range predecessors {
				if pred.books[l] != nil {
					lines = append(lines, pred.mark)
					fb = pred
				}
			}
			if pred == nil {
				panic("BREAKPOINT without correct predecessor")
			}
			predecessors = kp.To(pred)
		}
		T().Errorf("reversing the breakpoint list for line %d: %v", linecnt, lines)
		// Reverse the breakpoints-list
		for i := len(lines)/2 - 1; i >= 0; i-- {
			opp := len(lines) - 1 - i
			lines[i], lines[opp] = lines[opp], lines[i]
		}
		breakpoints[linecnt] = lines
	}
	return len(breakpoints), breakpoints
}
