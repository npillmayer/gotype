package lr

import (
	"golang.org/x/tools/container/intsets"
)

// LRAnalysis is an object for grammar analysis (compute FIRST and FOLLOW sets).
type LRAnalysis struct {
	g          *Grammar
	derivesEps map[*Symbol]bool
	firstSets  symSetMap
	followSets symSetMap
}

// Analysis creates an analyser for a grammar.
// The analyser immediately starts its work and computes FIRST and FOLLOW sets.
func Analysis(g *Grammar) *LRAnalysis {
	ga := makeAnalysis(g)
	ga.analyse()
	return ga
}

// create a fully initialized grammar analysis object
func makeAnalysis(g *Grammar) *LRAnalysis {
	ga := &LRAnalysis{}
	ga.g = g
	ga.derivesEps = make(map[*Symbol]bool)
	ga.firstSets = newSymSetMap()
	ga.followSets = newSymSetMap()
	return ga
}

func (ga *LRAnalysis) analyse() {
	ga.markEps()
	ga.initFirstSets()
	ga.initFollowSets()
}

// Grammar retunrns the grammer this analyser operates on.
func (ga *LRAnalysis) Grammar() *Grammar {
	return ga.g
}

// markEps marks all non-terminals from which epsilon can be derived.
func (ga *LRAnalysis) markEps() {
	changed := true
	for changed {
		changed = false
		for _, r := range ga.g.rules {
			rhsDerivesEps := true
			for _, A := range r.rhs {
				rhsDerivesEps = rhsDerivesEps && ga.derivesEps[A]
			}
			//T.Debugf("rhs eps?  %v   -> %v", r, rhsDerivesEps)
			if rhsDerivesEps && !ga.derivesEps[r.LHS] {
				changed = true
				ga.derivesEps[r.LHS] = true
			}
		}
	}
}

// --- Sets of terminals and non-terminals -----------------------------------

type symSetMap map[*Symbol]*intsets.Sparse

func newSymSetMap() symSetMap {
	return make(map[*Symbol]*intsets.Sparse)
}

func (m symSetMap) SetFor(forN *Symbol) *intsets.Sparse {
	symset, found := m[forN]
	if !found {
		symset = &intsets.Sparse{}
		m[forN] = symset
	}
	return symset
}

func (m symSetMap) addSymFor(forN *Symbol, A *Symbol) {
	m.SetFor(forN).Insert(A.Value)
}

func withoutEps(symset *intsets.Sparse) *intsets.Sparse {
	woeps := &intsets.Sparse{}
	woeps.Copy(symset)
	woeps.Remove(EpsilonType)
	return woeps
}

// --- FIRST and FOLLOW ------------------------------------------------------

// First returns the FIRST set for a non-terminal.
// Returns a set of token values.
func (ga *LRAnalysis) First(sym *Symbol) *intsets.Sparse {
	return ga.firstSets.SetFor(sym)
}

// Follow returns the FOLLOW set for a non-terminal.
// Returns a set of token values.
func (ga *LRAnalysis) Follow(sym *Symbol) *intsets.Sparse {
	return ga.followSets.SetFor(sym)
}

// DerivesEpsilon returns true if there are rules in the grammar which let
// a given symbol sym be derived to epsilon.
func (ga *LRAnalysis) DerivesEpsilon(sym *Symbol) bool {
	return ga.derivesEps[sym]
}

func (ga *LRAnalysis) computeFirst(syms []*Symbol) *intsets.Sparse {
	if len(syms) == 0 {
		epsset := &intsets.Sparse{}
		epsset.Insert(ga.g.Epsilon.Value)
		return epsset
	}
	first := ga.firstSets.SetFor(syms[0])
	result := withoutEps(first)
	//T.Infof(". c_first  : first(%v) = %v", syms, first)
	k := 1
	for ; k < len(syms); k++ {
		if first.Has(EpsilonType) { // prev one did contain epsilon
			//T.Infof("  . c_first: first(%v) = %v", syms[k-1:], first)
			first = ga.firstSets.SetFor(syms[k])
			result.UnionWith(first)
			result.Remove(EpsilonType)
			//T.Infof("  . c_first: first'(%v) = %v", syms[k:], first)
		} else {
			break
		}
	}
	if k == len(syms) && ga.firstSets.SetFor(syms[k-1]).Has(EpsilonType) {
		result.Insert(EpsilonType)
	}
	return result
}

func (ga *LRAnalysis) initFirstSets() {
	ga.g.EachNonTerminal(func(A *Symbol) interface{} {
		if ga.derivesEps[A] { // if A derives epsilon => epsilon in FIRST(A)
			ga.firstSets.addSymFor(A, ga.g.Epsilon)
		}
		return nil
	})
	ga.g.EachTerminal(func(B *Symbol) interface{} {
		ga.firstSets.addSymFor(B, B)
		for _, r := range ga.g.rules {
			A := r.LHS
			if !r.IsEps() && r.rhs[0].Value == B.Value {
				// if A -> t(B) ...
				ga.firstSets.addSymFor(A, B)
				//T.Infof("adding term = %v to first(%v)", B, A)
			}
		}
		return nil
	})
	for changed := true; changed; {
		changed = false
		for _, r := range ga.g.rules {
			A := r.LHS
			first := ga.computeFirst(r.rhs)
			//T.Infof("c_first(%v) = %v", r, first)
			//ch := ga.firstSets.SetFor(A).UnionWith(first) // must w/around bug in UnionWith()
			l := ga.firstSets.SetFor(A).Len()
			ga.firstSets.SetFor(A).UnionWith(first)
			ch := l < ga.firstSets.SetFor(A).Len()
			//T.Infof("(ch)anged = %v", ch)
			//T.Infof("new first(%v) = %v", A, ga.firstSets.getSetFor(A))
			changed = changed || ch
		}
	}
}

// 1) FOLLOW(S) = { $ }   // where S is the starting Non-Terminal
//    better: FOLLOW(S) = { Є }
//
// 2) If A -> xBy is a production, where p, B and q are any grammar symbols,
//    then everything in FIRST(y)  except Є is in FOLLOW(B).
//
// 3) If A->xB is a production, then everything in FOLLOW(A) is in FOLLOW(B).
//
// 4) If A->xBy is a production and FIRST(y) contains Є,
//    then FOLLOW(B) contains { FIRST(y) – Є } U FOLLOW(A)
//
func (ga *LRAnalysis) initFollowSets() {
	ga.followSets.addSymFor(ga.g.rules[0].LHS, ga.g.Epsilon) // start symbol
	changed := true
	for changed {
		changed = false
		for _, r := range ga.g.rules {
			T().Debugf("rule  %v", r)
			A := r.LHS                // look for A -> ... B y
			for k, B := range r.rhs { // look for non-terms in RHS of r
				if !B.IsTerminal() {
					y := r.rhs[k+1:]
					T().Debugf("      %v in RHS(%v),  y = %v", B, A, y)
					T().Debugf("      y = %v", y)
					yfirst := ga.computeFirst(y)
					//ch := ga.followSets.SetFor(B).UnionWith(withoutEps(yfirst)) // bug in intesets.UnionWith()
					l := ga.followSets.SetFor(B).Len()
					ga.followSets.SetFor(B).UnionWith(withoutEps(yfirst))
					ch := l < ga.followSets.SetFor(B).Len()
					if yfirst.Has(EpsilonType) {
						followA := ga.followSets.SetFor(A)
						l = ga.followSets.SetFor(B).Len()
						ga.followSets.SetFor(B).UnionWith(followA)
						ch = l < ga.followSets.SetFor(B).Len()
					}
					changed = changed || ch
				}
			}
		}
	}
}

// CleanUp cleans a context free grammar by removing unproductive non-terminals
// and rules, and by removing unreachable non-terminals.
//
// If the returned error is of kind ProductionsRemovedError, sub-errors are
// giving details about the removals. Check with errors.Is(…).
func (ga *LRAnalysis) CleanUp() error {
	// TODO Grune et al, 2.9.5
	return nil
}
