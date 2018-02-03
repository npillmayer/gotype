package lr

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/runtime"
)

/*
Grammar analysis (FIRST and FOLLOW sets).
*/
type GrammarAnalysis struct {
	g          *Grammar
	derivesEps map[Symbol]bool
	firstSets  *symSetMap
	followSets *symSetMap
}

// Create an analyser for a grammar.
// The analyser immediately starts its work and computes
// FIRST and FOLLOW.
func NewGrammarAnalysis(g *Grammar) *GrammarAnalysis {
	ga := &GrammarAnalysis{}
	ga.g = g
	ga.derivesEps = make(map[Symbol]bool)
	ga.firstSets = newSymSetMap()
	ga.followSets = newSymSetMap()
	return ga
}

// Return the FIRST set for a non-terminal.
// Returns a list of tokens.
func (ga *GrammarAnalysis) First(sym Symbol) []int {
	return ga.firstSets.getSetFor(sym).syms
}

// Return the FOLLOW set for a non-terminal.
// Returns a list of tokens.
func (ga *GrammarAnalysis) Follow(sym Symbol) []int {
	return ga.followSets.getSetFor(sym).syms
}

// ---------------------------------------------------------------------------

func (ga *GrammarAnalysis) markEps() {
	changed := true
	for changed {
		changed = false
		for _, r := range ga.g.rules {
			rhsDerivesEps := true
			for _, A := range r.rhs {
				rhsDerivesEps = rhsDerivesEps && ga.derivesEps[A]
			}
			//T.Debugf("rhs eps?  %v   -> %v", r, rhsDerivesEps)
			if rhsDerivesEps && !ga.derivesEps[r.lhs[0]] {
				changed = true
				ga.derivesEps[r.lhs[0]] = true
			}
		}
	}
}

// --- Sets of terminals and non-terminals -----------------------------------

type symSetMap struct {
	sets map[Symbol]*symSet
}

func newSymSetMap() *symSetMap {
	m := &symSetMap{}
	m.sets = make(map[Symbol]*symSet)
	return m
}

func (m *symSetMap) getSetFor(forN Symbol) *symSet {
	symset := m.sets[forN]
	if symset == nil {
		symset = newSymbolSet()
		m.sets[forN] = symset
	}
	return symset
}

func (m *symSetMap) addSymFor(forN Symbol, A Symbol) {
	symset := m.sets[forN]
	if symset == nil {
		symset = newSymbolSet()
		m.sets[forN] = symset
	}
	symset.add(A)
}

func symvalue(A Symbol) int {
	if A.IsTerminal() {
		return A.Token()
	}
	return A.GetID()
}

type symSet struct { // enable hashing on symSet
	syms []int
}

func newSymbolSet() *symSet {
	symset := &symSet{}
	symset.syms = make([]int, 0, 3)
	return symset
}

func (symset *symSet) clone() *symSet {
	c := newSymbolSet()
	copy(symset.syms, c.syms)
	return c
}

func (symset *symSet) contains(symval int) bool {
	for _, sym := range symset.syms {
		if sym == symval {
			return true
		}
	}
	return false
}

func (symset *symSet) containsEps() bool {
	for _, sym := range symset.syms {
		if sym == epsilonType {
			return true
		}
	}
	return false
}

func (symset *symSet) withoutEps() *symSet {
	if len(symset.syms) == 0 {
		return symset.clone()
	}
	woeps := newSymbolSet()
	for _, s := range symset.syms {
		if s != epsilonType {
			woeps.syms = append(woeps.syms, s)
		}
	}
	return woeps
}

func (symset *symSet) add(A Symbol) *symSet {
	if A == nil {
		symset.syms = append(symset.syms, epsilonType)
	} else {
		symset.syms = append(symset.syms, symvalue(A))
	}
	return symset
}

func (symset *symSet) union(set2 *symSet) (*symSet, bool) {
	changed := false
	if set2 != nil && len(set2.syms) > 0 {
		l := len(symset.syms)
		for _, sym := range set2.syms {
			if !symset.contains(sym) {
				symset.syms = append(symset.syms, sym)
			}
		}
		changed = l < len(symset.syms) // more elements than before?
	}
	return symset, changed
}

func (symset *symSet) String() string {
	s := fmt.Sprintf("%v", symset.syms)
	return s
}

// --- FIRST and FOLLOW ------------------------------------------------------

func (ga *GrammarAnalysis) computeFirst(syms []Symbol) *symSet {
	if len(syms) == 0 {
		epsset := newSymbolSet()
		epsset.add(ga.g.epsilon)
		return epsset
	}
	first := ga.firstSets.getSetFor(syms[0])
	var result *symSet = newSymbolSet()
	result.union(first)
	result = result.withoutEps()
	var k int = 1
	//T.Infof(". c_first  : first(%v) = %v", syms, first)
	for ; k < len(syms); k++ {
		if first.containsEps() { // prev one did contain epsilon
			//T.Infof("  . c_first: first(%v) = %v", syms[k-1:], first)
			first = ga.firstSets.getSetFor(syms[k])
			result.union(first)
			result = result.withoutEps()
			//T.Infof("  . c_first: first'(%v) = %v", syms[k:], first)
		} else {
			break
		}
	}
	if k == len(syms) && ga.firstSets.getSetFor(syms[k-1]).containsEps() {
		result.add(ga.g.epsilon)
	}
	return result
}

func (ga *GrammarAnalysis) initFirstSets() {
	ga.g.symbols.Each(func(name string, sym runtime.Symbol) {
		A := sym.(Symbol)
		if !A.IsTerminal() { // for all non-terminals A
			if ga.derivesEps[A] {
				ga.firstSets.addSymFor(A, ga.g.epsilon)
			}
		}
	})
	ga.g.symbols.Each(func(name string, sym runtime.Symbol) {
		B := sym.(Symbol)
		if B.IsTerminal() { // for all terminals B
			ga.firstSets.addSymFor(B, B)
			for _, r := range ga.g.rules {
				A := r.lhs[0]
				if !r.isEps() && symvalue(r.rhs[0]) == symvalue(B) {
					// if A -> t(B) ...
					ga.firstSets.addSymFor(A, B)
					//T.Infof("adding term = %v to first(%v)", B, A)
				}
			}
		}
	})
	for changed := true; changed; {
		changed = false
		for _, r := range ga.g.rules {
			A := r.lhs[0]
			first := ga.computeFirst(r.rhs)
			//T.Infof("c_first(%v) = %v", r, first)
			_, ch := ga.firstSets.getSetFor(A).union(first)
			//T.Infof("(ch)anged = %v", ch)
			//T.Infof("new first(%v) = %v", A, ga.firstSets.getSetFor(A))
			changed = changed || ch
		}
	}
}

func (ga *GrammarAnalysis) initFollowSets() {
	ga.followSets.addSymFor(ga.g.rules[0].lhs[0], ga.g.epsilon) // start symbol
	for changed := true; changed; {
		changed = false
		for _, r := range ga.g.rules {
			//T.Infof("rule  %v", r)
			A := r.lhs[0]             // look for A -> ... B y
			for k, B := range r.rhs { // look for non-terms in RHS of r
				if !B.IsTerminal() {
					y := r.rhs[k+1:]
					//T.Infof("      %v in RHS(%v),  y = %v", B, A, y)
					//T.Infof("      y = %v", y)
					yfirst := ga.computeFirst(y)
					_, ch := ga.followSets.getSetFor(B).union(yfirst.withoutEps())
					if yfirst.containsEps() {
						followA := ga.followSets.getSetFor(A)
						_, ch = ga.followSets.getSetFor(B).union(followA)
					}
					changed = changed || ch
				}
			}
		}
	}
}
