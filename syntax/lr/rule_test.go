package lr

import (
	"io/ioutil"
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"
)

func traceOn() {
	T().SetTraceLevel(tracing.LevelDebug)
}

func TestBuilder1(t *testing.T) {
	b := NewGrammarBuilder("G")
	b.LHS("S").N("A").End()
	if len(b.Grammar().rules) != 1 {
		t.Fail()
	}
	sym := b.Grammar().rules[0].lhs[0]
	if sym == nil || sym.GetID() > nonTermType {
		t.Fail()
	}
}

func TestBuilder2(t *testing.T) {
	b := NewGrammarBuilder("G")
	b.LHS("S").Epsilon()
	if len(b.Grammar().rules) != 1 {
		t.Fail()
	}
}

func TestItems1(t *testing.T) {
	b := NewGrammarBuilder("G")
	r := b.LHS("S").N("E").EOF()
	i1, _ := r.startItem()
	i2, _ := r.startItem() // items are cashed to not get duplicates
	if i1 != i2 {
		t.Fail()
	}
}

func TestClosure1(t *testing.T) {
	b := NewGrammarBuilder("G")
	r1 := b.LHS("S").N("E").EOF()
	r2 := b.LHS("E").N("E").T("+", 1).N("E").End()
	if len(b.Grammar().rules) != 2 {
		t.Fail()
	}
	item1, _ := r1.startItem()
	item2, _ := r2.startItem()
	T().Debugf("item1=%v", item1)
	T().Debugf("item2=%v", item2)
	if item1.dot != 0 {
		t.Fail()
	}
	item2, _ = item2.advance()
	T().Debugf("item2=%v", item2)
	if item2.dot != 1 {
		t.Fail()
	}
}

func TestClosure2(t *testing.T) {
	b := NewGrammarBuilder("G")
	r0 := b.LHS("S").N("E").EOF()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g := b.Grammar()
	ga := NewGrammarAnalysis(g)
	//b.Grammar().Dump()
	closure0 := ga.closure(r0.startItem())
	closure0.Dump()
}

func TestItemSetEquality(t *testing.T) {
	b := NewGrammarBuilder("G")
	r0 := b.LHS("S").N("E").EOF()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g := b.Grammar()
	ga := NewGrammarAnalysis(g)
	//b.Grammar().Dump()
	closure0 := ga.closure(r0.startItem())
	//closure0.Dump()
	closure1 := ga.closure(r0.startItem())
	if !closure0.equals(closure1) {
		t.Fail()
	}
}

func TestClosure4(t *testing.T) {
	b := NewGrammarBuilder("G")
	r0 := b.LHS("S").N("E").EOF()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g := b.Grammar()
	ga := NewGrammarAnalysis(g)
	//b.Grammar().Dump()
	i, A := r0.startItem()
	closure0 := ga.closure(i, A)
	//closure0.Dump()
	ga.gotoSet(closure0, A)
}

func TestStateRetrieval(t *testing.T) {
	b := NewGrammarBuilder("G")
	r0 := b.LHS("S").N("E").EOF()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g := b.Grammar()
	ga := NewGrammarAnalysis(g)
	cfsm := emptyCFSM(g)
	closure0 := ga.closure(r0.startItem())
	s0 := cfsm.addState(closure0)
	//s0.Dump()
	s1 := cfsm.addState(closure0)
	if s0.ID != s1.ID {
		t.Fail()
	}
}

func TestBuildCFSM(t *testing.T) {
	b := NewGrammarBuilder("G")
	b.LHS("S").N("E").EOF()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g := b.Grammar()
	ga := makeGrammarAnalysis(g)
	lrgen := NewLRTableGenerator(ga)
	c := lrgen.buildCFSM()
	c.CFSM2GraphViz("/tmp/cfsm-" + "G1" + ".dot")
}

func TestDerivesEps(t *testing.T) {
	b := NewGrammarBuilder("G")
	b.LHS("S").N("E").EOF()
	b.LHS("E").N("T").T("a", 1).End()
	b.LHS("E").N("F").End()
	b.LHS("F").Epsilon()
	g := b.Grammar()
	//g.Dump()
	ga := makeGrammarAnalysis(g)
	ga.markEps()
	cnt := 0
	g.EachSymbol(func(A Symbol) interface{} {
		T().Debugf("%v => eps  : %v", A, ga.derivesEps[A])
		if ga.derivesEps[A] {
			cnt++
		}
		return nil
	})
	if cnt != 2 {
		t.Fail() // E and F should => eps
	}
}

func TestFirstSet(t *testing.T) {
	b := NewGrammarBuilder("G")
	b.LHS("S").N("E").EOF()
	b.LHS("E").N("T").End()
	b.LHS("T").T("a", 1).End()
	b.LHS("T").Epsilon()
	g := b.Grammar()
	//g.Dump()
	ga := makeGrammarAnalysis(g)
	ga.markEps()
	ga.initFirstSets()
	for key, value := range ga.firstSets.sets {
		T().Debugf("key = %v     value = %v", key, value)
	}
}

func TestFollowSet(t *testing.T) {
	b := NewGrammarBuilder("G")
	b.LHS("S").N("A").T("a", 1).EOF()
	b.LHS("A").N("B").N("D").End()
	b.LHS("B").T("b", 2).End()
	b.LHS("B").Epsilon()
	b.LHS("D").T("d", 3).End()
	b.LHS("D").Epsilon()
	g := b.Grammar()
	//g.Dump()
	ga := makeGrammarAnalysis(g)
	ga.markEps()
	/*
		ga.g.symbols.Each(func(name string, sym runtime.Symbol) {
			A := sym.(Symbol)
			if !A.IsTerminal() {
				if ga.derivesEps[A] {
					T().Debugf("%v  => epsilon", A)
				}
			}
		})
	*/
	ga.initFirstSets()
	ga.Grammar().EachNonTerminal( // supply a mapper function
		func(A Symbol) interface{} {
			T().Debugf("FIRST(%v) = %v", A, ga.First(A))
			return nil
		})
	T().Debugf("-------")
	ga.initFollowSets()
	ga.Grammar().EachNonTerminal(
		func(A Symbol) interface{} {
			T().Debugf("FOLLOW(%v) = %v", A, ga.Follow(A))
			return nil
		})
}

func TestGotoTable(t *testing.T) {
	b := NewGrammarBuilder("G")
	b.LHS("S").N("A").EOF()
	b.LHS("A").T("a", 1).End()
	b.LHS("A").Epsilon()
	g := b.Grammar()
	g.Dump()
	ga := NewGrammarAnalysis(g)
	lrgen := NewLRTableGenerator(ga)
	lrgen.dfa = lrgen.buildCFSM()
	lrgen.gototable = lrgen.BuildGotoTable()
	lrgen.CFSM().CFSM2GraphViz("/tmp/cfsm.dot")
	tmp, _ := ioutil.TempFile("", "lr_")
	T().Infof("writing HTML to %s", tmp.Name())
	GotoTableAsHTML(lrgen, tmp)
}

func TestActionTable(t *testing.T) {
	b := NewGrammarBuilder("G")
	b.LHS("S").N("A").EOF()
	b.LHS("A").T("a", 1).End()
	b.LHS("A").Epsilon()
	g := b.Grammar()
	g.Dump()
	ga := NewGrammarAnalysis(g)
	lrgen := NewLRTableGenerator(ga)
	lrgen.dfa = lrgen.buildCFSM()
	traceOn()
	T().Debugf("\n---------- Action 0 -----------------------------------")
	lrgen.actiontable = lrgen.BuildLR0ActionTable()
	T().Debugf("\n---------- Action 1 -----------------------------------")
	lrgen.actiontable = lrgen.BuildSLR1ActionTable()
	/*
		tmp, _ := ioutil.TempFile("", "lr_")
		T().Infof("writing HTML to %s", tmp.Name())
		ActionTableAsHTML(lrgen, tmp)
	*/
}

// func TestScannerSimple1(t *testing.T) {
// 	scanner := NewScanner(nil)
// 	tokval, token := scanner.NextToken(nil)
// 	T().Infof("scanned: %d / %v", tokval, token)
// }
