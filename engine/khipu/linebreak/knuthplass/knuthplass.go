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
	"github.com/emirpasic/gods/stacks/arraystack"
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
	//lineno    int            // line number (line count this break creates)
	mark  khipu.Mark           // location of this breakpoint
	books map[int]*bookkeeping // bookkeeping per linecount
}

/*
type segment struct {
	cost      int32
	linecount int
}
*/

/*
func newSegment(cost int32, linecount int) segment {
	return segment{
		cost:      cost,
		linecount: linecount,
	}
}
*/

/*
func (s segment) extract() (cost int32, linecount int) {
	cost = s.cost
	linecount = s.linecount
	return
}
*/

type bookkeeping struct {
	fragment  linebreak.WSS // sum of widths up to last knot
	totalcost int32         // sum of costs for segments up to this breakpoint
}

type provisionalMark int // provisional mark from an integer position

func (m provisionalMark) Position() int    { return int(m) }
func (m provisionalMark) Knot() khipu.Knot { return nil }

func (fb *feasibleBreakpoint) String() string {
	if fb.mark == nil {
		return "<para-start>"
	}
	return fmt.Sprintf("<break p.%d @ %v>", fb.mark.Position(), fb.mark.Knot())
}

func (fb *feasibleBreakpoint) UpdateBook(linecnt int, fragment linebreak.WSS, total int32) {
	if fb.books == nil {
		fb.books = make(map[int]*bookkeeping)
	}
	fb.books[linecnt] = &bookkeeping{
		fragment:  fragment,
		totalcost: total,
	}
}

func (fb *feasibleBreakpoint) UpdateSegment(linecnt int, diff linebreak.WSS) {
	if fb.books == nil {
		fb.books = make(map[int]*bookkeeping)
	}
	segment := linebreak.WSS{}
	total := int32(0)
	book, ok := fb.books[linecnt]
	if ok {
		segment = book.fragment
		total = book.totalcost
	}
	fb.books[linecnt] = &bookkeeping{
		fragment:  segment.Add(diff),
		totalcost: total,
	}
}

func (fb *feasibleBreakpoint) ClearBook(linecnt int) {
	if fb.books == nil {
		return
	}
	delete(fb.books, linecnt)
}

func (fb *feasibleBreakpoint) Book(linecnt int) (*bookkeeping, bool) {
	if fb.books == nil {
		return &bookkeeping{}, false
	}
	b, ok := fb.books[linecnt]
	return b, ok
}

// newBreakpointAtMark creates a breakpoint at the given cursor position.
func (kp *linebreaker) newBreakpointAtMark(mark khipu.Mark) *feasibleBreakpoint {
	fb := &feasibleBreakpoint{
		mark: mark,
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

/*
type feasibleSegment struct { // an edge between nodes of feasible breakpoints
	wEdge                // edge with cost
	totals linebreak.WSS // the cumulative widths of this segment
}
*/

/*
func (fseg *feasibleSegment) String() string {
	if fseg == nil {
		return "+--no edge-->"
	}
	return fmt.Sprintf("+-%d--(%d)--%d->", fseg.from, fseg.cost, fseg.to)
}
*/

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
func (kp *linebreaker) feasibleLineBetween(from, to *feasibleBreakpoint,
	cost int32, linecnt int, prune bool) bool {
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
	if prune { // then look for existing demertics
		if pp := kp.findPredecessorsWithLinecount(to, linecnt); pp != nil {
			if len(pp) > 1 {
				panic("breakpoint (with pruning) has more than one predecessor[line]")
			}
			predecessor = pp[0]
			// isolate total cost of predecessor for segment to 'to'
			predCost := kp.Edge(predecessor, to, linecnt).cost
			mintotal = predecessor.books[linecnt].totalcost + predCost
		}
	}
	book, ok := from.Book(linecnt)
	if !ok {
		panic("totalcost of horizon node must be set") // TODO remove
	}
	if book.totalcost+cost < mintotal { // new line is cheaper
		//edge := newWEdge(from, to, cost, linecnt)
		//seg = &feasibleSegment{edge, totals}
		kp.AddSegment(from, to, cost, linecnt)
		totalcost := book.totalcost + cost
		wss := linebreak.WSS{}.SetFromKnot(to.mark.Knot()) // dimensions of next knot
		fragment := book.fragment.Add(wss)
		to.UpdateBook(linecnt, fragment, totalcost)
		if prune && predecessor != nil {
			kp.RemoveEdge(predecessor, to, linecnt)
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
	T().Debugf(" width of %v is %.2f", knot, knot.W)
	cannotReachIt := 0
	//for linecnt, bookkeeping := range fb.books {
	for linecnt := range fb.books {
		T().Debugf(" ## checking cost at linecnt=%d", linecnt)
		d := linebreak.InfinityDemerits                   // pre-set result variable
		linelen := parshape.LineLength(linecnt + 1)       // length of line to fit into
		segwss := fb.calculateSegmentWidth(knot, linecnt) // widths of segment including knot
		T().Debugf("    +---%.2f--->    | %.2f", segwss.W.Points(), linelen.Points())
		if segwss.W <= linelen { // natural width less than line-length
			if segwss.Max >= linelen { // segment can stretch enough
				d = calculateDemerits(segwss, linelen-segwss.W, 0)
			} else { // segment is just too short
				// try with tolerance
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

func (fb *feasibleBreakpoint) calculateSegmentWidth(knot khipu.Knot, linecnt int) linebreak.WSS {
	var segwss linebreak.WSS
	bookkeeping, ok := fb.Book(linecnt)
	if !ok {
		panic(fmt.Errorf("bookkeeping for line %d MUST be present", linecnt))
	}
	wss := linebreak.WSS{}.SetFromKnot(knot) // dimensions of next knot
	segwss = bookkeeping.fragment.Add(wss)
	return segwss
}

func calculateDemerits(segwss linebreak.WSS, stretch dimen.Dimen, tolerance int) int32 {
	tolerancepenalty := 1000 // TODO from typesetting parameters
	tolerancepenalty *= tolerance
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

//func (kp *linebreaker) FindBreakpoints(input *khipu.Khipu, prune bool) (int, []*feasibleBreakpoint) {

// FindBreakpoints is the main client API.
func FindBreakpoints(cursor linebreak.Cursor, parshape linebreak.Parshape, prune bool) (
	int, []khipu.Mark) {
	//
	kp := newLinebreaker(parshape)
	fb := kp.newBreakpointAtMark(provisionalMark(-1)) // start of paragraph
	fb.UpdateBook(0, linebreak.WSS{}, int32(0))
	kp.horizon.append(fb) // this is the first "active node"
	//cursor := input.Iterator()
	//var knot khipu.Knot // will hold the current knot of the input khipu
	var last khipu.Mark // will hold last position within input khipu
	for cursor.Next() { // loop over input
		last = cursor.Mark()
		//fmt.Printf("next knot is %s\n", khipu.KnotString(knot))
		fb = kp.horizon.first() // loop over feasible breakpoints of horizon
		if fb == nil {
			T().Errorf("no more active breakpoints, but input available")
			T().Errorf("this should probably have produces an overfull hbox")
		}
		for fb != nil { // while there are active breakpoints in horizon n
			//fmt.Printf("fb in horizon = %v\n", fb)
			costs, stillreachable := fb.calculateCostsTo(cursor.Knot(), parshape)
			// TODO only do this for valid breakpoints, i.e. discr & penalties
			// suppress for boxes and kern
			for linecnt, cost := range costs {
				if cost < linebreak.InfinityDemerits { // new breakpoint is feasible
					//fmt.Printf("create feasible breakpoint at %v\n", cursor)
					newfb := kp.findBreakpointAtMark(cursor.Mark())
					if newfb == nil { // breakpoint not yet existent
						newfb = kp.newBreakpointAtMark(cursor.Mark()) // create a feasible breakpoint
					}
					ok := kp.feasibleLineBetween(fb, newfb, cost, linecnt, prune)
					if ok {
						T().Debugf("new line established to %v", newfb)
					}
					//fmt.Printf("feasible line = %v\n", edge)
					kp.horizon.append(newfb) // make new fb member of horizon n+1
				} else {
					wss := linebreak.WSS{}.SetFromKnot(cursor.Knot()) // get dimensions of knot
					fb.UpdateSegment(linecnt, wss)
					book, _ := fb.Book(linecnt)
					T().Debugf("extending segment to %v", book.fragment)
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

func (kp *linebreaker) collectFeasibleBreakpoints(last khipu.Mark) (int, []khipu.Mark) {
	//var optimalBreaks []*feasibleBreakpoint
	//fmt.Printf("collecting breakpoints, backwards from #%d\n", last)
	fb := kp.findBreakpointAtMark(last)
	if fb == nil {
		fb = kp.newBreakpointAtMark(last)
	}
	stack := arraystack.New() // for reversing the breakpoint order, TODO this is overkill
	var breakpoints = make([]khipu.Mark, 0, 20)
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
		T().Infof("optimal paragraph breaking uses %d breakpoints\n", stack.Size())
		p, ok := stack.Pop()
		for ok {
			breakpoints = append(breakpoints, p.(khipu.Mark))
			p, ok = stack.Pop()
		}
	}
	return len(breakpoints), breakpoints
}
