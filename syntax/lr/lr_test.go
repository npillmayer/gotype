package lr

import (
	"io/ioutil"
	"log"
	"testing"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"golang.org/x/tools/container/intsets"
)

var graphviz bool

func TestBuilder1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	b.LHS("S").N("A").End()
	g, err := b.Grammar()
	if len(g.rules) != 2 {
		t.Errorf("Expected G to have 2 rules, has %d", len(g.rules))
	}
	sym := g.rules[0].LHS
	if err != nil || sym == nil || sym.Value > NonTermType {
		t.Errorf("Symbol S is broken")
	}
}

func TestBuilder2(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	b.LHS("S").Epsilon()
	g, err := b.Grammar()
	if err != nil || len(g.rules) != 2 {
		t.Errorf("Expected G to have 2 rules, has %d", len(g.rules))
	}
}

func TestClosure1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	r1 := b.LHS("S").N("E").End()
	r2 := b.LHS("E").N("E").T("+", 1).N("E").End()
	item1, _ := StartItem(r1)
	item2, _ := StartItem(r2)
	t.Logf("item2=%v", item1)
	t.Logf("item2=%v", item2)
	if item1.dot != 0 {
		t.Errorf("Start items must have dot at position 0")
	}
	item2 = item2.Advance()
	t.Logf("item2=%v", item2)
	if item2.dot != 1 {
		t.Errorf("Item expected to have dot at position 1, is %d", item2.dot)
	}
}

func TestSet1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	set1 := &intsets.Sparse{}
	set2 := &intsets.Sparse{}
	set1.Insert(1)
	set1.Insert(2)
	set1.Insert(9)
	var set3 = &intsets.Sparse{}
	set3.Copy(set1)
	t.Logf("set 1 = %s", set1)
	set2.Insert(1)
	set2.Insert(2)
	t.Logf("set 2 = %s", set2)
	if set1.UnionWith(set2) {
		t.Errorf("set 1 not expected to change, is now %s", set1)
	}
	if set1.Equals(set3) {
		t.Logf("however, set 1 is unchanged")
		t.Logf("intsets.UnionWith() contains a bug")
	}
}

func TestClosure2(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	r0 := b.LHS("S").N("E").End()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g, _ := b.Grammar()
	g.Dump()
	ga := Analysis(g)
	closure0 := ga.closure(StartItem(r0))
	Dump(closure0)
}

func TestItemSetEquality(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	r0 := b.LHS("S").N("E").End()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g, _ := b.Grammar()
	ga := Analysis(g)
	//b.Grammar().Dump()
	closure0 := ga.closure(StartItem(r0))
	//closure0.Dump()
	closure1 := ga.closure(StartItem(r0))
	if !closure0.Equals(closure1) {
		t.Errorf("closure0 expected to equal closure1")
	}
}

func TestClosure4(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	r0 := b.LHS("S").N("E").End()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g, _ := b.Grammar()
	ga := Analysis(g)
	//b.Grammar().Dump()
	i, A := StartItem(r0)
	closure0 := ga.closure(i, A)
	//closure0.Dump()
	ga.gotoSet(closure0, A)
}

func TestStateRetrieval(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	r0 := b.LHS("S").N("E").End()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g, _ := b.Grammar()
	ga := Analysis(g)
	cfsm := emptyCFSM(g)
	closure0 := ga.closure(StartItem(r0))
	s0 := cfsm.addState(closure0)
	//s0.Dump()
	s1 := cfsm.addState(closure0)
	if s0.ID != s1.ID {
		t.Errorf("State 0 and state 1 expected to have the same ID")
	}
}

func TestBuildCFSM(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	b.LHS("S").N("E").End()
	b.LHS("E").N("E").T("+", 1).T("(", 2).N("E").T(")", 3).End()
	b.LHS("E").T("a", 4).End()
	g, _ := b.Grammar()
	ga := Analysis(g)
	lrgen := NewTableGenerator(ga)
	c := lrgen.buildCFSM()
	c.CFSM2GraphViz("/tmp/cfsm-" + "G1" + ".dot")
}

func TestDerivesEps(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	b.LHS("S").N("E").End()
	b.LHS("E").N("T").T("a", 1).End()
	b.LHS("E").N("F").End()
	b.LHS("F").Epsilon()
	g, _ := b.Grammar()
	ga := Analysis(g)
	//g.Dump()
	ga.markEps()
	cnt := 0
	g.EachSymbol(func(A *Symbol) interface{} {
		T().Debugf("%v => eps  : %v", A, ga.derivesEps[A])
		if ga.derivesEps[A] {
			cnt++
		}
		return nil
	})
	if cnt != 3 {
		t.Errorf("S, E and F should derive epsilon")
	}
}

func TestFirstSet(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	b.LHS("S").N("E").End()
	b.LHS("E").N("T").End()
	b.LHS("T").T("a", 1).End()
	b.LHS("T").Epsilon()
	g, _ := b.Grammar()
	//g.Dump()
	//ga := Analysis(g)
	ga := makeAnalysis(g)
	ga.markEps()
	ga.initFirstSets()
	for key, value := range ga.firstSets {
		T().Debugf("key = %v     value = %v", key, value)
	}
}

func TestFollowSet(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	b.LHS("S").N("A").T("a", 1).End()
	b.LHS("A").N("B").N("D").End()
	b.LHS("B").T("b", 2).End()
	b.LHS("B").Epsilon()
	b.LHS("D").T("d", 3).End()
	b.LHS("D").Epsilon()
	g, _ := b.Grammar()
	//g.Dump()
	ga := Analysis(g)
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
		func(A *Symbol) interface{} {
			T().Debugf("FIRST(%v) = %v", A, ga.First(A))
			return nil
		})
	T().Debugf("-------")
	ga.initFollowSets()
	ga.Grammar().EachNonTerminal(
		func(A *Symbol) interface{} {
			T().Debugf("FOLLOW(%v) = %v", A, ga.Follow(A))
			return nil
		})
}

func TestGotoTable(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	b.LHS("S").N("A").EOF()
	b.LHS("A").T("a", 1).End()
	b.LHS("A").Epsilon()
	g, _ := b.Grammar()
	g.Dump()
	ga := Analysis(g)
	lrgen := NewTableGenerator(ga)
	lrgen.dfa = lrgen.buildCFSM()
	lrgen.gototable = lrgen.BuildGotoTable()
	if graphviz {
		lrgen.CFSM().CFSM2GraphViz("/tmp/cfsm.dot")
		tmp, _ := ioutil.TempFile("", "lr_")
		T().Infof("writing HTML to %s", tmp.Name())
		GotoTableAsHTML(lrgen, tmp)
	}
}

func TestActionTable(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G")
	b.LHS("S").N("A").EOF()
	b.LHS("A").T("a", 1).End()
	b.LHS("A").Epsilon()
	g, _ := b.Grammar()
	g.Dump()
	ga := Analysis(g)
	lrgen := NewTableGenerator(ga)
	lrgen.dfa = lrgen.buildCFSM()
	T().Debugf("\n---------- Action 0 -----------------------------------")
	lrgen.actiontable, _ = lrgen.BuildLR0ActionTable()
	T().Debugf("\n---------- Action 1 -----------------------------------")
	lrgen.actiontable, _ = lrgen.BuildSLR1ActionTable()
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

func TestCraftingG2(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G2")
	b.LHS("S").N("A").EOF()
	b.LHS("A").T("a", scanner.Ident).End()
	b.LHS("A").Epsilon()
	g, _ := b.Grammar()
	g.Dump()
	ga := Analysis(g)
	lrgen := NewTableGenerator(ga)
	lrgen.CreateTables()
	if graphviz {
		tmpfile, err := ioutil.TempFile(".", "G2_*.html")
		if err != nil {
			log.Fatal(err)
		}
		GotoTableAsHTML(lrgen, tmpfile)
	}
}

func TestTerminals1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G6")
	b.LHS("S").T("(", '(').N("A").T(")", ')').EOF()
	b.LHS("A").T("a", scanner.Ident).End()
	g, _ := b.Grammar()
	g.Dump()
	ga := Analysis(g)
	lrgen := NewTableGenerator(ga)
	lrgen.CreateTables()
	if graphviz {
		tmpfile, err := ioutil.TempFile(".", "G6_*.html")
		if err != nil {
			log.Fatal(err)
		}
		GotoTableAsHTML(lrgen, tmpfile)
		lrgen.dfa = lrgen.buildCFSM()
		lrgen.CFSM().CFSM2GraphViz("./G6_cfsm.dot")
	}
}

func TestExercise1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("E6")
	b.LHS("S").N("A").T("a", scanner.Ident).End()
	b.LHS("A").N("B").N("D").End()
	b.LHS("B").T("b", scanner.Ident-1).End()
	b.LHS("B").Epsilon()
	b.LHS("D").T("d", scanner.Ident-2).End()
	b.LHS("D").Epsilon()
	g, _ := b.Grammar()
	g.Dump()
	ga := Analysis(g)
	for sym, firstSet := range ga.firstSets {
		t.Logf("First(%s) = %v", sym, firstSet)
	}
	for sym, followSet := range ga.followSets {
		t.Logf("Follow(%s) = %v", sym, followSet)
		if sym.Name == "S" { // start symbol
			if followSet.Len() != 1 || !followSet.Has(EOFType) {
				t.Errorf("Expected Follow(S) to be {#eof}")
			}
		}
	}
	lrgen := NewTableGenerator(ga)
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	lrgen.CreateTables()
	// tmpfile, err := ioutil.TempFile(".", "E6_*.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//GotoTableAsHTML(lrgen, tmpfile)
	//lrgen.dfa = lrgen.buildCFSM()
	if graphviz {
		lrgen.CFSM().CFSM2GraphViz("./E6_cfsm.dot")
	}
}

func TestGrammar7(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := NewGrammarBuilder("G7")
	b.LHS("S'").N("S").EOF()
	b.LHS("S").N("A").T("a", scanner.Ident).End()
	b.LHS("A").T("+", '+').End()
	b.LHS("A").T("-", '-').End()
	b.LHS("A").Epsilon()
	g, _ := b.Grammar()
	g.Dump()
	ga := Analysis(g)
	lrgen := NewTableGenerator(ga)
	lrgen.CreateTables()
	if graphviz {
		tmpfile, err := ioutil.TempFile(".", "G7_action_*.html")
		if err != nil {
			log.Fatal(err)
		}
		//GotoTableAsHTML(lrgen, tmpfile)
		ActionTableAsHTML(lrgen, tmpfile)
		//lrgen.dfa = lrgen.buildCFSM()
		//lrgen.CFSM().CFSM2GraphViz("./G7_cfsm.dot")
	}
}
