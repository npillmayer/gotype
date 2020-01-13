package knuthplass

/*
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
	"errors"
	"fmt"
	"io"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/khipu"
	"github.com/npillmayer/gotype/engine/khipu/linebreak"
)

// linebreaker is an internal entity for K&P-linebreaking.
type linebreaker struct {
	*fbGraph
	horizon  *activeFeasibleBreakpoints // horizon of possible linebreaks
	params   *linebreak.Parameters      // typesetting parameters relevant for line-breaking
	parshape linebreak.Parshape         // target shape of the paragraph
	root     *feasibleBreakpoint        // "break" at start of paragraph
	end      *feasibleBreakpoint        // "break" at end of paragraph
}

func newLinebreaker(parshape linebreak.Parshape, params *linebreak.Parameters) *linebreaker {
	kp := &linebreaker{}
	kp.fbGraph = newFBGraph()
	kp.horizon = newActiveFeasibleBreakpoints()
	kp.parshape = parshape
	if params == nil {
		params = NewKPDefaultParameters()
	}
	kp.params = params
	return kp
}

// NewKPDefaultParameters creates line-breaking parameters similar to
// (but not identical) to TeX's.
func NewKPDefaultParameters() *linebreak.Parameters {
	return &linebreak.Parameters{
		Tolerance:            1000,
		PreTolerance:         500,
		LinePenalty:          20,
		DoubleHyphenDemerits: 200,
		FinalHyphenDemerits:  500,
		EmergencyStretch:     dimen.Dimen(dimen.BP * 20),
	}
}

func setupLinebreaker(cursor linebreak.Cursor, parshape linebreak.Parshape,
	params *linebreak.Parameters) (*linebreaker, error) {
	if parshape == nil {
		return nil, errors.New("Cannot shape a paragraph without a Parshape")
	}
	kp := newLinebreaker(parshape, params)
	fb := kp.newBreakpointAtMark(provisionalMark(-1)) // start of paragraph
	fb.books[0] = &bookkeeping{}
	kp.root = fb       // remember the start breakpoint as root of the graph
	kp.horizon.Add(fb) // this is the first 'active node' of horizon
	return kp, nil
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
	if fb.mark == nil || fb.mark.Position() < 0 {
		return "<para-start>"
	}
	return fmt.Sprintf("<brk %d/%v>", fb.mark.Position(), fb.mark.Knot())
}

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
	return kp.Breakpoint(mark.Position()) // may be nil
}

func (kp *linebreaker) findPredecessorsWithLinecount(fb *feasibleBreakpoint, linecnt int) []*feasibleBreakpoint {
	var predecessors []*feasibleBreakpoint
	edges := kp.EdgesTo(fb).WithLabel(linecnt)
	for _, edge := range edges {
		if edge.isNull() {
			panic("this should not happen!") // TODO remove this
		}
		if edge.linecount == linecnt {
			predecessors = append(predecessors, kp.StartOfEdge(edge))
		}
	}
	return predecessors
}

// --- Segments ---------------------------------------------------------

// newFeasibleLine possibly creates a segment between two given breakpoints.
//
// The segment is constructed and compared to
// any existing segments (for the same line-count). If its cost is cheaper
// than the exising one, the new segment replaces the old one
// (just one segment between the two breakpoints can exist with pruning).
func (kp *linebreaker) newFeasibleLine(fb *feasibleBreakpoint, mark khipu.Mark, cost int32,
	linecnt int) *feasibleBreakpoint {
	//
	newfb := kp.findBreakpointAtMark(mark)
	if newfb == nil { // breakpoint not yet existent => create one
		newfb = kp.newBreakpointAtMark(mark)
	}
	targettotal := fb.Book(linecnt-1).totalcost + cost // total cost of new line
	if kp.isCheapestSurvivor(newfb, targettotal, linecnt) {
		kp.AddEdge(fb, newfb, cost, targettotal, linecnt)
		newfb.books[linecnt] = &bookkeeping{totalcost: targettotal}
		T().Debugf("new line established %v ---%d---> %v", fb, cost, newfb)
	}
	return newfb
}

// isCheapestSurvivor calculates the total cost for a new segment, and compares it
// to all existing segments. If the new segment would be cheaper,
// the others will die.
func (kp *linebreaker) isCheapestSurvivor(fb *feasibleBreakpoint, totalcost int32,
	linecnt int) bool {
	//
	var predecessor *feasibleBreakpoint          // predecessor breakpoint position
	mintotal := linebreak.InfinityDemerits * 100 // pre-set to hyper-infinity
	//
	// Calculate an edge from fb to a new hypothetical breakpoint.
	// If the total cost for the new edge would be cheaper than every existing
	// edge, and deleteOthers is set, remove the more expensive edges.
	if pp := kp.findPredecessorsWithLinecount(fb, linecnt); pp != nil {
		if len(pp) > 1 {
			panic("breakpoint (with pruning) has more than one predecessor[line]")
		}
		predecessor = pp[0]
		// isolate total cost of predecessor for segment to fb
		predCost := kp.Edge(predecessor, fb, linecnt).cost
		mintotal = predecessor.books[linecnt].totalcost + predCost
	}
	//T().Debugf("mintotal=%d, totalcost=%d", mintotal, totalcost)
	if totalcost < mintotal { // new line is cheaper
		if predecessor != nil {
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
func (fb *feasibleBreakpoint) calculateCostsTo(penalty khipu.Penalty, parshape linebreak.Parshape,
	params *linebreak.Parameters) (map[int]int32, bool) {
	//
	T().Infof("### calculateCostsTo(%v)", penalty)
	var costs = make(map[int]int32) // linecount => cost, i.e. costs for different line targets
	cannotReachIt := 0
	for linecnt := range fb.books {
		T().Debugf(" ## checking cost at linecnt=%d", linecnt)
		d := linebreak.InfinityDemerits             // pre-set result variable
		linelen := parshape.LineLength(linecnt + 1) // length of line to fit into
		segwss := fb.segmentWidth(linecnt)
		T().Debugf("    +---%.2f--->    | %.2f", segwss.W.Points(), linelen.Points())
		if segwss.W <= linelen { // natural width less than line-length
			if segwss.Max >= linelen { // segment can stretch enough
				d = calculateDemerits(segwss, linelen-segwss.W, penalty, params)
			} else { // segment is just too short
				if params.EmergencyStretch > 0 {
					emStretch := segwss.Max + params.EmergencyStretch
					if emStretch >= linelen { // now segment can stretch enough
						d = calculateDemerits(segwss, linelen-segwss.W, penalty, params)
					}
				}
			}
		} else { // natural width larger than line-length
			if segwss.Min <= linelen { // segment can shrink enough
				d = calculateDemerits(segwss, segwss.W-linelen, penalty, params)
			} else { // segment will not fit any more
				cannotReachIt++
			}
		}
		costs[linecnt] = linebreak.CapDemerits(d)
		T().Debugf(" ## cost for line %d would be %s", linecnt+1, demeritsString(costs[linecnt]))
	}
	stillreachable := (cannotReachIt < len(fb.books))
	T().Debugf("### costs to %v is %v, reachable is %v", penalty, costs, stillreachable)
	return costs, stillreachable
}

// segmentWidth returns the widths of a segment at fb, subtracting discardable
// items at the start of the segment and at the end (= possible breakpoint).
func (fb *feasibleBreakpoint) segmentWidth(linecnt int) linebreak.WSS {
	segw := fb.Book(linecnt).segment
	segw = segw.Subtract(fb.Book(linecnt).startDiscard)
	segw = segw.Subtract(fb.Book(linecnt).breakDiscard)
	return segw
}

// Currently we try to replicated logic of TeX.
func calculateDemerits(segwss linebreak.WSS, stretch dimen.Dimen, penalty khipu.Penalty,
	params *linebreak.Parameters) int32 {
	//
	var d int32
	p := linebreak.CapDemerits(penalty.Demerits())
	//p2 := p * p
	p2 := p // seems to work better for now; related to segmenter behaviour
	badness := int32(abs(stretch) / max(1+abs(segwss.Max-segwss.W), 1) * 100)
	b := (params.LinePenalty + badness)
	b2 := b * b
	if p > 0 {
		d = b2 + p2
	} else if p <= linebreak.InfinityMerits {
		d = b2
	} else {
		d = b2 - p2
	}
	d = linebreak.CapDemerits(d)
	T().Debugf("Calculating demerits for p=%d, b=%d: d=%d", p, b, d)
	return d
}

func demeritsString(d int32) string {
	if d >= linebreak.InfinityDemerits {
		return "\u221e"
	} else if d <= linebreak.InfinityMerits {
		return "-\u221e"
	}
	return fmt.Sprintf("%d", d)
}

// penaltyAt iterates over all penalties, starting at the current cursor mark, and
// collects penalties, searching for the most significant one.
// Will return
//
//        -10000, if present
//        max(p1, p2, ..., pn) otherwise
//
// Returns the most significant penalty. Advances the cursor over all adjacent penalties.
// After this, the cursor mark may not reflect the position of the significant penalty.
func penaltyAt(cursor linebreak.Cursor) (khipu.Penalty, khipu.Mark) {
	if cursor.Knot().Type() != khipu.KTPenalty {
		return khipu.Penalty(linebreak.InfinityDemerits), cursor.Mark()
	}
	penalty := cursor.Knot().(khipu.Penalty)
	ignore := false // final penalty found, ignore all other penalties
	knot, ok := cursor.Peek()
	for ok {
		if knot.Type() == khipu.KTPenalty {
			cursor.Next() // advance to next penalty
			if ignore {
				break // just skip over adjacent penalties
			}
			p := knot.(khipu.Penalty)
			if p.Demerits() <= linebreak.InfinityMerits { // -10000 must break (like in TeX)
				penalty = p
				ignore = true
			} else if p.Demerits() > penalty.Demerits() {
				penalty = p
			}
			knot, ok = cursor.Peek() // now check next knot
		} else {
			ok = false
		}
	}
	p := khipu.Penalty(linebreak.CapDemerits(penalty.Demerits()))
	return p, cursor.Mark()
}

// --- Main API ---------------------------------------------------------

// BreakParagraph determines optimal linebreaks for a paragraph, depending on
// a given set of linebreaking parameters and the desired shape of the paragraph.
//
// Paragraphs may be broken with different line counts. Only one of these will be
// optimal, but Break
func BreakParagraph(cursor linebreak.Cursor, parshape linebreak.Parshape,
	params *linebreak.Parameters) ([]khipu.Mark, error) {
	//
	variants, breakpoints, err := FindBreakpoints(cursor, parshape, params, nil)
	if err != nil {
		return nil, err
	}
	if len(breakpoints) == 0 {
		return nil, fmt.Errorf("No breakpoints could be found for paragraph")
	}
	best := variants[0] // slice is sorted by increasing totalcost, first one is best
	return breakpoints[best], err
}

// FindBreakpoints is the main client API.
func FindBreakpoints(cursor linebreak.Cursor, parshape linebreak.Parshape, params *linebreak.Parameters,
	dotfile io.Writer) ([]int, map[int][]khipu.Mark, error) {
	//
	kp, err := setupLinebreaker(cursor, parshape, params)
	if err != nil {
		return nil, nil, err
	}
	err = kp.constructBreakpointGraph(cursor, parshape, params)
	if err != nil {
		T().Errorf(err.Error())
		return nil, nil, err
	}
	variants, breaks := kp.collectFeasibleBreakpoints(kp.end)
	c := khipu.NewCursor(cursor.Khipu())
	if dotfile != nil {
		kp.toGraphViz(c, breaks, dotfile)
	}
	return variants, breaks, nil
}

func (kp *linebreaker) constructBreakpointGraph(cursor linebreak.Cursor, parshape linebreak.Parshape,
	params *linebreak.Parameters) error {
	//
	var last khipu.Mark        // will hold last position within input khipu
	var fb *feasibleBreakpoint // will hold feasible breakpoint from horizon
	for cursor.Next() {        // outer loop over input knots
		last = cursor.Mark() // we will need the last knot at the end of the loop
		T().Debugf("_______________ %d/%s ___________________", last.Position(), last.Knot())
		if fb = kp.horizon.first(); fb == nil {
			panic("no more active breakpoints, but input available") // TODO remove after debugging
		}
		// --- main loop over active breakpoints in horizon ------------
		for fb != nil { // loop over active feasible breakpoints of horizon
			T().Debugf("                %d/%v  (in horizon)", fb.mark.Position(), fb.mark.Knot())
			fb.UpdateSegmentBookkeeping(cursor.Mark())
			// Breakpoints are allowed at penalties only
			if cursor.Mark().Knot().Type() == khipu.KTPenalty { // TODO discretionaries
				stillreachable := true
				var penalty khipu.Penalty
				penalty, last = penaltyAt(cursor) // find correct p, if more than one
				var costs map[int]int32           // we want cost per linecnt-alternative
				costs, stillreachable = fb.calculateCostsTo(penalty, parshape, kp.params)
				if stillreachable { // yes, position may have been reached in this iteration
					for linecnt, cost := range costs { // check for every linecount alternative
						if penalty.Demerits() <= linebreak.InfinityMerits && // forced break
							cost > kp.params.Tolerance {
							T().Infof("Underfull box at line %d, cost=%v", linecnt+1, cost)
							newfb := kp.newFeasibleLine(fb, cursor.Mark(), cost, linecnt+1)
							kp.horizon.Add(newfb) // make forced break member of horizon n+1
						} else if cost < kp.params.Tolerance { // happy case: new breakpoint is feasible
							newfb := kp.newFeasibleLine(fb, cursor.Mark(), cost, linecnt+1)
							kp.horizon.Add(newfb) // make new breakpoint member of horizon n+1
						}
					}
				} else { // no longer reachable => check against draining of horizon
					if kp.horizon.Size() <= 1 { // oops, low on options
						for linecnt := range costs {
							T().Infof("Overfull box at line %d, cost=10000", linecnt+1)
							newfb := kp.newFeasibleLine(fb, cursor.Mark(), linebreak.InfinityDemerits, linecnt+1)
							kp.horizon.Add(newfb) // make new fb member of horizon n+1
							if newfb.mark.Position() == fb.mark.Position() {
								panic("THIS SHOULD NOT HAPPEN ?!?")
							}
						}
					}
					kp.horizon.Remove(fb) // no longer valid in horizon
				}
			}
			fb = kp.horizon.next()
		} // --- end of main loop over horizon ----------------------
	} // end of outer loop over input knots
	fb = kp.findBreakpointAtMark(last)
	if fb == nil {
		// for now panic, for debugging purposes
		panic("last breakpoint not found") // khipu didn't end with penalty -10000
		// TODO add fb(-10000) and connect to last horizon
		// in this situation, input is drained but horizon is not ?!
	}
	kp.end = fb // remember last breakpoint of paragraph
	return nil
}

// Collecting breakpoints, backwards from last
func (kp *linebreaker) collectFeasibleBreakpoints(last *feasibleBreakpoint) (
	[]int, map[int][]khipu.Mark) {
	breakpoints := make(map[int][]khipu.Mark) // list of breakpoints per linecount-variant
	costDict := make(map[int]int32)           // list of total-costs per linecount-variant
	lineVariants := make([]int, 0, 3)         // will become sorted list of linecount-variants
	for linecnt, book := range last.books {
		costDict[linecnt] = book.totalcost
		i := 0
		for j, lv := range lineVariants {
			if book.totalcost < costDict[lv] {
				i = j
				break
			}
		}
		insert(lineVariants, i, linecnt)
		l := linecnt
		lines := make([]khipu.Mark, 0, 20)
		lines = append(lines, last.mark)
		predecessors := kp.findPredecessorsWithLinecount(last, l)
		for len(predecessors) > 0 { // while not at start node
			l-- // searching for predecessor with linecount-1
			if len(predecessors) > 1 {
				panic("THERE SHOULD ONLY BE ONE PREDECESSOR") // TODO remove after debugging
			}
			pred := predecessors[0]
			lines = append(lines, pred.mark)
			predecessors = kp.findPredecessorsWithLinecount(pred, l)
		}
		T().Errorf("reversing the breakpoint list for line %d: %v", linecnt, lines)
		// Now reverse the breakpoints-list
		for i := len(lines)/2 - 1; i >= 0; i-- {
			opp := len(lines) - 1 - i
			lines[i], lines[opp] = lines[opp], lines[i]
		}
		breakpoints[linecnt] = lines
	}
	for l := range costDict {
		lineVariants = append(lineVariants, l)
	}
	return lineVariants, breakpoints
}

// --- Helpers ----------------------------------------------------------

func abs(n dimen.Dimen) dimen.Dimen {
	if n < 0 {
		return -n
	}
	return n
}

func max(n, m dimen.Dimen) dimen.Dimen {
	if n > m {
		return n
	}
	return m
}

func insert(s []int, i, n int) {
	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = n
}
