package lr

import (
	"fmt"
	"io"
	"os"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
	"github.com/npillmayer/gotype/syntax/lr/sparse"
)

// TODO: Improve documentation...
// https://stackoverflow.com/questions/12968048/what-is-the-closure-of-a-left-recursive-lr0-item-with-epsilon-transitions
// = optimization
// https://www.cs.bgu.ac.il/~comp151/wiki.files/ps6.html#sec-2-7-3

// === Items and Item Sets ===================================================

// A set of Earley items ( A -> B *C D ).
type itemSet struct {
	*hashset.Set
}

func newItemSet() *itemSet {
	s := hashset.New()
	iset := &itemSet{s}
	return iset
}

func (iset *itemSet) union(iset2 *itemSet) {
	iset.Add(iset2.Values()...)
}

func (iset *itemSet) equals(iset2 *itemSet) bool {
	if iset.Size() == iset2.Size() {
		for _, i := range iset.Values() {
			if !iset2.Contains(i) {
				return false
			}
		}
		return true
	}
	return false
}

// iset2 ist input and output (destructive)
func (iset *itemSet) difference(iset2 *itemSet) *itemSet {
	for _, i := range iset2.Values() {
		if iset.Contains(i) {
			iset2.Remove(i)
		}
	}
	return iset2
}

func (iset *itemSet) String() string {
	items := iset.Values()
	if len(items) == 0 {
		return "{ }"
	} else {
		s := "{\n"
		for _, i := range items {
			s = s + fmt.Sprintf("\t%v\n", i)
		}
		s += "}"
		return s
	}
}

// Prepare an item set for export to Graphviz.
func (iset *itemSet) forGraphviz() string {
	items := iset.Values()
	if len(items) == 0 {
		return "err\\n"
	} else {
		s := ""
		for _, i := range items {
			s = s + fmt.Sprintf("%v\\l", i)
		}
		return s
	}
}

// Debugging helper
func (iset *itemSet) Dump() {
	items := iset.Values()
	//T.Debug("--- item set ------------")
	for k, i := range items {
		T.Debugf("item %2d = %v", k, i)
	}
	//T.Debug("-------------------------")
}

var _ *itemSet = newItemSet() // verify assignability

// === Closure and Goto-Set Operations =======================================

// Refer to "Crafting A Compiler" by Charles N. Fisher & Richard J. LeBlanc, Jr.
// Section 6.2.1 LR(0) Parsing

// Compute the closure of an Earley item.
func (ga *GrammarAnalysis) closure(i *item, A Symbol) *itemSet {
	iset := newItemSet()
	iset.Add(i)
	if A == nil {
		A = i.peekSymbol() // get symbol after dot
	}
	if A != nil {
		//T.Debugf("pre closure(%v) = %v", i, iset)
		iset = ga.closureSet(iset)
		//T.Debugf("    closure(%v) = %v", i, iset)
		return iset
	}
	return iset
}

// Compute the closure of an Earley item.
// https://www.cs.bgu.ac.il/~comp151/wiki.files/ps6.html#sec-2-7-3
func (ga *GrammarAnalysis) closureSet(iset *itemSet) *itemSet {
	cset := newItemSet()   // will be our closure result set
	cset.union(iset)       // add start item to closure
	tmpset := newItemSet() // this will collect derived items for the next iteration
	for !iset.Empty() {
		for _, x := range iset.Values() {
			i := x.(*item)                                        // LHS -> X *A Y
			if A := i.peekSymbol(); A != nil && !A.IsTerminal() { // if A is non-term
				iiset := ga.g.findNonTermRules(A, false) // without eps-productions
				cset.union(iiset)                        // add { A -> *... } to closure
				tmpset.union(iiset)                      // prepare for next iteration
				if ga.derivesEps[A] {                    // then we have to look past A
					j, B := i.advance() // move dot 1 position
					if B != nil {       // if not at end of RHS
						cset.Add(j)   // add to closure
						tmpset.Add(j) // prepare for next iteration
					}
				}
			}
		}
		tmpset, iset = iset, iset.difference(tmpset) // swap; this is correct!
		tmpset.Clear()
	}
	return cset
}

func (ga *GrammarAnalysis) gotoSet(closure *itemSet, A Symbol) (*itemSet, Symbol) {
	// for every item in closure C
	// if item in C:  N -> ... *A ...
	//     advance N -> ... A * ...
	gotoset := newItemSet()
	for _, x := range closure.Values() {
		i := x.(*item)
		if i.peekSymbol() == A {
			ii, _ := i.advance()
			T.Debugf("goto(%s) -%s-> %s", i, A, ii)
			gotoset.Add(ii)
		}
	}
	//gotoset.Dump()
	return gotoset, A
}

func (ga *GrammarAnalysis) gotoSetClosure(i *itemSet, A Symbol) (*itemSet, Symbol) {
	gotoset, _ := ga.gotoSet(i, A)
	//T.Infof("gotoset  = %v", gotoset)
	gclosure := ga.closureSet(gotoset)
	//T.Infof("gclosure = %v", gclosure)
	T.Debugf("goto(%s) --%s--> %s", i, A, gclosure)
	return gclosure, A
}

// === CFSM Construction =====================================================

// CFSM state
type CFSMState struct {
	ID     int      // serial ID of this state
	items  *itemSet // configuration items within this state
	Accept bool     // is this an accepting state?
}

// CFSM edge between 2 states, directed and labeled with a terminal
type cfsmEdge struct {
	from  *CFSMState
	to    *CFSMState
	label Symbol
}

// Debugging helper
func (s *CFSMState) Dump() {
	T.Debugf("--- state %03d -----------", s.ID)
	s.items.Dump()
	T.Debugf("-------------------------")
}

func (s *CFSMState) isErrorState() bool {
	return s.items.Size() == 0
}

// Create a state from an item set
func state(id int, iset *itemSet) *CFSMState {
	s := &CFSMState{ID: id}
	if iset == nil {
		s.items = newItemSet()
	} else {
		s.items = iset
	}
	return s
}

func (s *CFSMState) allItems() []interface{} {
	vals := s.items.Values()
	return vals
}

func (s *CFSMState) String() string {
	return fmt.Sprintf("(state %d | [%d])", s.ID, s.items.Size())
}

func (s *CFSMState) containsCompletedStartRule() bool {
	for _, x := range s.items.Values() {
		i := x.(*item)
		if i.rule.no == 0 && i.peekSymbol() == nil {
			return true
		}
	}
	return false
}

// Create an edge
func edge(from, to *CFSMState, label Symbol) *cfsmEdge {
	return &cfsmEdge{
		from:  from,
		to:    to,
		label: label,
	}
}

// We need this for the set of states. It sorts states by serial ID.
func stateComparator(s1, s2 interface{}) int {
	c1 := s1.(*CFSMState)
	c2 := s2.(*CFSMState)
	return utils.IntComparator(c1.ID, c2.ID)
}

// Add a state to the CFSM. Checks first if state is present.
func (c *CFSM) addState(iset *itemSet) *CFSMState {
	s := c.findStateByItems(iset)
	if s == nil {
		s = state(c.cfsmIds, iset)
		c.cfsmIds++
	}
	c.states.Add(s)
	return s
}

// Find a CFSM state by the contained item set.
func (c *CFSM) findStateByItems(iset *itemSet) *CFSMState {
	it := c.states.Iterator()
	for it.Next() {
		s := it.Value().(*CFSMState)
		if s.items.equals(iset) {
			return s
		}
	}
	return nil
}

func (c *CFSM) addEdge(s0, s1 *CFSMState, sym Symbol) *cfsmEdge {
	e := edge(s0, s1, sym)
	c.edges.Add(e)
	return e
}

func (c *CFSM) allEdges(s *CFSMState) []*cfsmEdge {
	it := c.edges.Iterator()
	r := make([]*cfsmEdge, 0, 2)
	for it.Next() {
		e := it.Value().(*cfsmEdge)
		if e.from == s {
			r = append(r, e)
		}
	}
	return r
}

// LR(0) state diagram for a grammar, i.e. the characteristic finite
// state automata CFSM. Will be constructed by a LRTableGenerator.
// Clients normally do not use it directly. Nevertheless, there are some methods
// defined on it, e.g, for debugging purposes, or even to
// compute your own tables from it.
type CFSM struct {
	g       *Grammar        // this CFSM is for Grammar g
	states  *treeset.Set    // all the states
	edges   *arraylist.List // all the edges between states
	S0      *CFSMState      // start state
	cfsmIds int             // serial IDs for CFSM states
}

// create an empty (initial) CFSM automata.
func emptyCFSM(g *Grammar) *CFSM {
	c := &CFSM{g: g}
	c.states = treeset.NewWith(stateComparator)
	c.edges = arraylist.New()
	return c
}

/*
Generator object to construct LR parser tables.
Clients usually create a Grammar G, then a GrammarAnalysis-object for G,
and then a table generator. LRTableGenerator.CreateTables() constructs
the CFSM and parser tables for an LR-parser recognizing grammar G.
*/
type LRTableGenerator struct {
	g           *Grammar
	ga          *GrammarAnalysis
	dfa         *CFSM
	gototable   *sparse.IntMatrix
	actiontable *sparse.IntMatrix
}

/*
Create a new LRTableGenerator for a (previously analysed) grammar.
*/
func NewLRTableGenerator(ga *GrammarAnalysis) *LRTableGenerator {
	lrgen := &LRTableGenerator{}
	lrgen.g = ga.Grammar()
	lrgen.ga = ga
	return lrgen
}

/*
Return the characteristic finite state machine (CFSM) for a grammar.
Usually clients call lrgen.CreateTables() beforehand, but it is possible
to call lrgen.CFSM() directly. The CFSM will be created, if it has not
been constructed previously.
*/
func (lrgen *LRTableGenerator) CFSM() *CFSM {
	if lrgen.dfa == nil {
		lrgen.dfa = lrgen.buildCFSM()
	}
	return lrgen.dfa
}

/*
Return the GOTO table for LR-parsing a grammar. The tables have to be
built by calling CreateTables() previously (or a separate call to
BuildGotoTable(...).)
*/
func (lrgen *LRTableGenerator) GotoTable() *sparse.IntMatrix {
	if lrgen.gototable == nil {
		T.P("lr", "gen").Errorf("tables not yet initialized")
	}
	return lrgen.gototable
}

/*
Return the ACTION table for LR-parsing a grammar. The tables have to be
built by calling CreateTables() previously (or a separate call to
BuildSLR1ActionTable(...).)
*/
func (lrgen *LRTableGenerator) ActionTable() *sparse.IntMatrix {
	if lrgen.actiontable == nil {
		T.P("lr", "gen").Errorf("tables not yet initialized")
	}
	return lrgen.actiontable
}

/*
Create the necessary data structures for an SLR parser.
*/
func (lrgen *LRTableGenerator) CreateTables() {
	lrgen.dfa = lrgen.buildCFSM()
	lrgen.gototable = lrgen.BuildGotoTable()
	lrgen.actiontable = lrgen.BuildSLR1ActionTable()
}

// Return all states of the CFSM which represent an accept action.
// Clients have to call CreateTables() first.
func (lrgen *LRTableGenerator) AcceptingStates() []int {
	if lrgen.dfa == nil {
		T.Errorf("tables not yet generated; call CreateTables() first")
		return nil
	}
	acc := make([]int, 0, 3)
	for _, x := range lrgen.dfa.states.Values() {
		state := x.(*CFSMState)
		if state.Accept {
			acc = append(acc, state.ID)
		}
	}
	return acc
}

// Construct the characteristic finite state machine CFSM for a grammar.
func (lrgen *LRTableGenerator) buildCFSM() *CFSM {
	T.Debugf("=== build CFSM ==================================================")
	G := lrgen.g
	cfsm := emptyCFSM(G)
	closure0 := lrgen.ga.closure(G.rules[0].startItem())
	cfsm.S0 = cfsm.addState(closure0)
	cfsm.S0.Dump()
	S := treeset.NewWith(stateComparator)
	S.Add(cfsm.S0)
	for S.Size() > 0 {
		s := S.Values()[0].(*CFSMState)
		S.Remove(s)
		G.EachSymbol(func(A Symbol) interface{} {
			T.Debugf("checking goto-set for symbol = %v", A)
			gotoset, _ := lrgen.ga.gotoSetClosure(s.items, A)
			snew := cfsm.findStateByItems(gotoset)
			if snew == nil {
				snew = cfsm.addState(gotoset)
				if !snew.isErrorState() {
					S.Add(snew)
					if snew.containsCompletedStartRule() {
						snew.Accept = true
					}
				}
			}
			if !snew.isErrorState() {
				cfsm.addEdge(s, snew, A)
			}
			snew.Dump()
			return nil
		})
		T.Debugf("-----------------------------------------------------------------")
	}
	return cfsm
}

// Export an CFSM to the Graphviz Dot format, given a filename.
func (cfsm *CFSM) CFSM2GraphViz(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("file open error: %v", err.Error()))
	}
	defer f.Close()
	f.WriteString(`digraph {
node [shape=Mrecord style=filled];

`)
	for _, x := range cfsm.states.Values() {
		s := x.(*CFSMState)
		f.WriteString(fmt.Sprintf("s%03d [fillcolor=%s label=\"{%03d | %s}\"]\n",
			s.ID, nodecolor(s), s.ID, s.items.forGraphviz()))
	}
	it := cfsm.edges.Iterator()
	for it.Next() {
		x := it.Value()
		edge := x.(*cfsmEdge)
		f.WriteString(fmt.Sprintf("s%03d -> s%03d [label=\"%s\"]\n", edge.from.ID, edge.to.ID, edge.label))
	}
	f.WriteString("}\n")
}

func nodecolor(state *CFSMState) string {
	if state.Accept {
		return "lightgray"
	}
	return "white"
}

// ===========================================================================

/*
Build the GOTO table. This is normally not called directly, but rather
via CreateTables().
*/
func (lrgen *LRTableGenerator) BuildGotoTable() *sparse.IntMatrix {
	statescnt := lrgen.dfa.states.Size()
	maxtok := 0
	lrgen.g.EachSymbol(func(A Symbol) interface{} {
		if A.Token() > maxtok { // find maximum token value
			maxtok = A.Token()
		}
		return nil
	})
	T.Infof("GOTO table of size %d x %d", statescnt, maxtok)
	gototable := sparse.NewIntMatrix(statescnt, maxtok, sparse.DefaultNullValue)
	states := lrgen.dfa.states.Iterator()
	for states.Next() {
		state := states.Value().(*CFSMState)
		edges := lrgen.dfa.allEdges(state)
		for _, e := range edges {
			//T.Debugf("edge %s --%v--> %v", state, e.label, e.to)
			//T.Debugf("GOTO (%d , %d ) = %d", state.ID, symvalue(e.label), e.to.ID)
			gototable.Set(state.ID, symvalue(e.label), int32(e.to.ID))
		}
	}
	return gototable
}

/*
Export the GOTO-table in HTML-format.
*/
func GotoTableAsHTML(lrgen *LRTableGenerator, w io.Writer) {
	if lrgen.gototable == nil {
		T.Errorf("GOTO table not yet created, cannot export to HTML")
		return
	}
	parserTableAsHTML(lrgen, "GOTO", lrgen.gototable, w)
}

/*
Export the SLR(1) ACTION-table in HTML-format.
*/
func ActionTableAsHTML(lrgen *LRTableGenerator, w io.Writer) {
	if lrgen.actiontable == nil {
		T.Errorf("ACTION table not yet created, cannot export to HTML")
		return
	}
	parserTableAsHTML(lrgen, "ACTION", lrgen.actiontable, w)
}

func parserTableAsHTML(lrgen *LRTableGenerator, tname string, table *sparse.IntMatrix, w io.Writer) {
	var symvec []Symbol = make([]Symbol, len(lrgen.g.terminals)+len(lrgen.g.nonterminals))
	io.WriteString(w, "<html><body>\n")
	io.WriteString(w, "<img src=\"cfsm.png\"/><p>")
	io.WriteString(w, fmt.Sprintf("%s table of size = %d<p>", tname, table.ValueCount()))
	io.WriteString(w, "<table border=1 cellspacing=0 cellpadding=5>\n")
	io.WriteString(w, "<tr bgcolor=#cccccc><td></td>\n")
	j := 0
	lrgen.g.EachSymbol(func(A Symbol) interface{} {
		io.WriteString(w, fmt.Sprintf("<td>%s</td>", A))
		symvec[j] = A
		j++
		return nil
	})
	io.WriteString(w, "</tr>\n")
	states := lrgen.dfa.states.Iterator()
	var td string // table cell
	for states.Next() {
		state := states.Value().(*CFSMState)
		io.WriteString(w, fmt.Sprintf("<tr><td>state %d</td>\n", state.ID))
		for _, A := range symvec {
			v1, v2 := table.Values(state.ID, symvalue(A))
			if v1 == table.NullValue() {
				td = "&nbsp;"
			} else if v2 == table.NullValue() {
				td = fmt.Sprintf("%d", v1)
			} else {
				td = fmt.Sprintf("%d/%d", v1, v2)
			}
			io.WriteString(w, "<td>")
			io.WriteString(w, td)
			io.WriteString(w, "</td>\n")
		}
		io.WriteString(w, "</tr>\n")
	}
	io.WriteString(w, "</table></body></html>\n")
}

// ===========================================================================

/*
Build the LR(0) Action table. This method is not called by CreateTables(),
as we normally use an SLR(1) parser and therefore an action table with
lookahead included. This method is provided as an add-on.
*/
func (lrgen *LRTableGenerator) BuildLR0ActionTable() *sparse.IntMatrix {
	statescnt := lrgen.dfa.states.Size()
	T.Infof("ACTION.0 table of size %d x 1", statescnt)
	actions := sparse.NewIntMatrix(statescnt, 1, sparse.DefaultNullValue)
	return lrgen.buildActionTable(actions, false)
}

/*
Build the SLR(1) Action table. This method is normally not called by
clients, but rather via CreateTables(). It builds an action table including
lookahead (using the FOLLOW-set created by the grammar analyzer).
*/
func (lrgen *LRTableGenerator) BuildSLR1ActionTable() *sparse.IntMatrix {
	statescnt := lrgen.dfa.states.Size()
	maxtok := 0
	lrgen.g.EachSymbol(func(A Symbol) interface{} {
		if A.Token() > maxtok { // find maximum token value
			maxtok = A.Token()
		}
		return nil
	})
	T.Infof("ACTION.1 table of size %d x %d", statescnt, maxtok)
	actions := sparse.NewIntMatrix(statescnt, maxtok, sparse.DefaultNullValue)
	return lrgen.buildActionTable(actions, true)
}

/*
For building an ACTION table we iterate over all the states of the CFSM.
An inner loop iterates over alle the Earley items within a CFSM-state.
If an item has a non-terminal immediately after the dot, we produce a shift
entry. If an item's dot is behind the complete (non-epsilon) RHS of a rule,
then
- for the LR(0) case: we produce a reduce-entry for the rule
- for the SLR case: we produce a reduce-entry for for the rule for each
  terminal from FOLLOW(LHS).

The table is returned as a sparse matrix, where every entry may consist of up
to 2 entries, thus allowing for shift/reduce- or reduce/reduce-conflicts.

Shift entries are represented as -1.  Reduce entries are encoded as the
ordinal no. of the grammar rule to reduce. 0 means reducing the start rule,
i.e., accept.
*/
func (lrgen *LRTableGenerator) buildActionTable(actions *sparse.IntMatrix, slr1 bool) *sparse.IntMatrix {
	states := lrgen.dfa.states.Iterator()
	for states.Next() {
		state := states.Value().(*CFSMState)
		for _, v := range state.items.Values() {
			T.Debugf("item in s%d = %v", state.ID, v)
			i, _ := v.(*item)
			A := i.peekSymbol()
			prefix := i.getPrefix()
			T.Debugf("symbol at dot = %v, prefix = %v", A, prefix)
			if A != nil && A.IsTerminal() { // create a shift entry
				if slr1 {
					T.Debugf("    creating shift action entry --%v-->", A)
					a := actions.Value(state.ID, A.Token())
					if a != -1 { // already shift present ?
						actions.Add(state.ID, A.Token(), -1)
					}
				} else {
					T.Debugf("    creating shift action entry")
					actions.Add(state.ID, 1, -1) // general shift (no lookahead)
				}
			}
			if len(prefix) > 0 && A == nil { // we are at the end of a non-eps rule
				rule, inx := lrgen.g.matchesRHS(prefix, false) // find the rule
				if inx >= 0 {                                  // found => create a reduce entry
					if slr1 {
						lookaheads := lrgen.ga.Follow(rule.lhs[0])
						T.Debugf("    Follow(%v) = %v", rule.lhs[0], lookaheads)
						for _, la := range lookaheads {
							actions.Add(state.ID, la, int32(inx))  // reduce rule[inx]
							if rule.no == 0 && la == epsilonType { // start rule reduced
								T.Debugf("    accepting") // TODO
							}
							T.Debugf("    creating reduce_%d action entry @ %v for %v", inx, la, rule)
						}
					} else {
						T.Debugf("    creating reduce_%d action entry for %v", inx, rule)
						actions.Add(state.ID, 1, int32(inx)) // reduce rule[inx]
					}
				}
			}
		}
	}
	return actions
}
