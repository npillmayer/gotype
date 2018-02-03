package lr

import (
	"fmt"
	"os"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
	"github.com/npillmayer/gotype/syntax/runtime"
)

// === Table Generation ======================================================

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

func (iset *itemSet) forDot() string {
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

func (iset *itemSet) Dump() {
	items := iset.Values()
	//T.Debug("--- item set ------------")
	for k, i := range items {
		T.Debugf("item %2d = %v", k, i)
	}
	//T.Debug("-------------------------")
}

var _ *itemSet = newItemSet() // verify assignability

func (g *Grammar) closure(i *item, A Symbol) *itemSet {
	iset := newItemSet()
	iset.Add(i)
	if A == nil {
		A = i.peekSymbol() // get symbol after dot
	}
	if A != nil {
		iset = g.closureSet(iset)
		T.Debugf("closure(%v) = %v", i, iset)
		return iset
	}
	return iset
}

// https://www.cs.bgu.ac.il/~comp151/wiki.files/ps6.html#sec-2-7-3
func (g *Grammar) closureSet(iset *itemSet) *itemSet {
	cset := newItemSet()
	cset.union(iset)
	for _, x := range iset.Values() {
		i := x.(*item)
		if A := i.peekSymbol(); A != nil {
			// iterate through all rules
			// is LHS = A ?
			// create item A ::= * RHS  ? How to proceed with eps-rules?
			if !A.IsTerminal() {
				iiset := g.findNonTermRules(A)
				//T.Debugf("found %d items for closure", iiset.Size())
				cset.union(iiset)
			}
		}
	}
	return cset
}

func (g *Grammar) gotoSet(closure *itemSet, A Symbol) (*itemSet, Symbol) {
	// for every item in closure C
	// if item in C:  N -> ... * A ...
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

func (g *Grammar) gotoSetClosure(i *itemSet, A Symbol) (*itemSet, Symbol) {
	gotoset, _ := g.gotoSet(i, A)
	//T.Infof("gotoset  = %v", gotoset)
	gclosure := g.closureSet(gotoset)
	//T.Infof("gclosure = %v", gclosure)
	T.Debugf("goto(%s) --%s--> %s", i, A, gclosure)
	return gclosure, A
}

var cfsmIds int = 1

type cfsmState struct {
	id    int
	items *itemSet
}

type cfsmEdge struct {
	from  *cfsmState
	to    *cfsmState
	label Symbol
}

func (s *cfsmState) Dump() {
	T.Debugf("--- state %03d -----------", s.id)
	s.items.Dump()
	T.Debug("-------------------------")
}

func (s *cfsmState) isErrorState() bool {
	return s.items.Size() == 0
}

func state(iset *itemSet) *cfsmState {
	s := &cfsmState{id: cfsmIds}
	cfsmIds++
	if iset == nil {
		s.items = newItemSet()
	} else {
		s.items = iset
	}
	return s
}

func edge(from, to *cfsmState, label Symbol) *cfsmEdge {
	return &cfsmEdge{
		from:  from,
		to:    to,
		label: label,
	}
}

func stateComparator(s1, s2 interface{}) int {
	c1 := s1.(*cfsmState)
	c2 := s2.(*cfsmState)
	return utils.IntComparator(c1.id, c2.id)
}

/*
func edgeComparator(e1, e2 interface{}) int {
    c1 := e1.(*cfsmEdge)
    c2 := e2.(*cfsmEdge)
    return utils.IntComparator(c1.from.id, c2.id)
}
*/

func (c *cfsm) addState(iset *itemSet) *cfsmState {
	s := c.findStateByItems(iset)
	if s == nil {
		s = state(iset)
	}
	c.states.Add(s)
	return s
}

func (c *cfsm) findStateByItems(iset *itemSet) *cfsmState {
	it := c.states.Iterator()
	for it.Next() {
		s := it.Value().(*cfsmState)
		if s.items.equals(iset) {
			return s
		}
	}
	return nil
}

func (c *cfsm) addEdge(s0, s1 *cfsmState, sym Symbol) *cfsmEdge {
	e := edge(s0, s1, sym)
	c.edges.Add(e)
	return e
}

type cfsm struct {
	states *treeset.Set
	edges  *arraylist.List
}

func emptyCFSM() *cfsm {
	c := &cfsm{}
	c.states = treeset.NewWith(stateComparator)
	c.edges = arraylist.New()
	return c
}

func (g *Grammar) buildCFSM() *cfsm {
	T.Debug("=== build CFSM ==================================================")
	r0 := g.rules[0]
	closure0 := g.closure(r0.startItem())
	cfsm := emptyCFSM()
	s0 := cfsm.addState(closure0)
	s0.Dump()
	S := treeset.NewWith(stateComparator)
	S.Add(s0)
	for S.Size() > 0 {
		s := S.Values()[0].(*cfsmState)
		S.Remove(s)
		g.symbols.Each(func(name string, sym runtime.Symbol) {
			T.Debugf("sym %s = %v", name, sym)
			A := sym.(Symbol)
			gotoset, _ := g.gotoSetClosure(s.items, A)
			snew := cfsm.findStateByItems(gotoset)
			if snew == nil {
				snew = cfsm.addState(gotoset)
				if !snew.isErrorState() {
					S.Add(snew)
				}
			}
			if !snew.isErrorState() {
				cfsm.addEdge(s, snew, A)
			}
			snew.Dump()
		})
		T.Debug("-----------------------------------------------------------------")
	}
	return cfsm
}

func cfsm2dot(c *cfsm) {
	f, err := os.Create("/tmp/cfsm.dot")
	if err != nil {
		panic(fmt.Sprintf("file open error: %v", err.Error()))
	}
	defer f.Close()
	f.WriteString(`digraph {
node [shape=record];

`)
	for _, x := range c.states.Values() {
		s := x.(*cfsmState)
		f.WriteString(fmt.Sprintf("s%03d [label=\"{%03d | %s}\"]\n", s.id, s.id, s.items.forDot()))
	}
	it := c.edges.Iterator()
	for it.Next() {
		x := it.Value()
		edge := x.(*cfsmEdge)
		f.WriteString(fmt.Sprintf("s%03d -> s%03d [label=\"%s\"]\n", edge.from.id, edge.to.id, edge.label))
	}
	f.WriteString("}\n")
}
