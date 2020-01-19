package parser

import (
	"strings"
	"testing"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/dss"
)

// https://github.com/golang/go/wiki/Performance

func traceOn() {
	T().SetTraceLevel(tracing.LevelDebug)
}

/*
  0: S ::= [A a #eof]
  1: A ::= [B]
  2: A ::= [C]
  3: B ::= [+]
  4: B ::= []
  5: C ::= [-]
  6: C ::= []
*/
func makeGrammar1() *lr.Grammar { // un-ambiguous
	b := lr.NewGrammarBuilder("G1")
	b.LHS("S").N("A").T("a", scanner.Ident).EOF()
	b.LHS("A").N("B").End()
	b.LHS("A").N("C").End()
	b.LHS("B").T("+", '+').End()
	b.LHS("B").Epsilon()
	b.LHS("C").T("-", '-').End()
	b.LHS("C").Epsilon()
	return b.Grammar()
}

/*
https://cs.au.dk/~amoeller/papers/ambiguity/ambiguity.pdf  -> Example 4

  0: S' ::= [S #eof]
  1: S  ::= [A -]
  2: S  ::= [+ B]
  3: A  ::= [+ a]
  4: B  ::= [a -]
*/
func makeGrammar2() *lr.Grammar { // ambiguous
	b := lr.NewGrammarBuilder("G2")
	b.LHS("S'").N("S").EOF()
	b.LHS("S").N("A").T("-", '-').End()
	b.LHS("S").T("+", '+').N("B").End()
	b.LHS("A").T("+", '+').T("a", scanner.Ident).End()
	b.LHS("B").T("a", scanner.Ident).T("-", '-').End()
	return b.Grammar()
}

// --- Test Cases ------------------------------------------------------------

func TestStackSet1(t *testing.T) {
	set1 := newStackSet()
	r := dss.NewRoot("G", -999)
	s1 := dss.NewStack(r)
	set1 = set1.add(s1)
	T().Debugf("set = %v", set1)
	if set1.get() != s1 {
		t.Fail()
	}
}

func TestParser1(t *testing.T) {
	g := makeGrammar1()
	ga := lr.NewGrammarAnalysis(g)
	lrgen := lr.NewLRTableGenerator(ga)
	lrgen.CreateTables()
	parser := Create(g, lrgen.GotoTable(), lrgen.ActionTable(), lrgen.AcceptingStates())
	r := strings.NewReader("a")
	scanner := NewStdScanner(r)
	parser.Parse(lrgen.CFSM().S0, scanner)
}

func TestParser2(t *testing.T) {
	g := makeGrammar1()
	//g.Dump()
	ga := lr.NewGrammarAnalysis(g)
	lrgen := lr.NewLRTableGenerator(ga)
	lrgen.CreateTables()
	parser := Create(g, lrgen.GotoTable(), lrgen.ActionTable(), lrgen.AcceptingStates())
	r := strings.NewReader("+a")
	scanner := NewStdScanner(r)
	parser.Parse(lrgen.CFSM().S0, scanner)
	/*
	   lrgen.CFSM().CFSM2GraphViz("/tmp/cfsm-" + "G1" + ".dot")
	   tmp, _ := ioutil.TempFile("", "lr_")
	   T().Infof("writing HTML to %s", tmp.Name())
	   lr.ActionTableAsHTML(lrgen, tmp)
	   lr.GotoTableAsHTML(lrgen, tmp)
	   tmp.Close()
	*/
}

func TestParser3(t *testing.T) {
	g := makeGrammar2()
	g.Dump()
	ga := lr.NewGrammarAnalysis(g)
	lrgen := lr.NewLRTableGenerator(ga)
	lrgen.CreateTables()
	parser := Create(g, lrgen.GotoTable(), lrgen.ActionTable(), lrgen.AcceptingStates())
	r := strings.NewReader("+a-")
	scanner := NewStdScanner(r)
	//traceOn()
	parser.Parse(lrgen.CFSM().S0, scanner)
}

func TestParser4(t *testing.T) {
	gtrace.SyntaxTracer = gologadapter.New()
	b := lr.NewGrammarBuilder("G4")
	b.LHS("S").N("A").EOF()
	b.LHS("A").T("a", scanner.Ident).End()
	b.LHS("A").Epsilon()
	g := b.Grammar()
	g.Dump()
	ga := lr.NewGrammarAnalysis(g)
	lrgen := lr.NewLRTableGenerator(ga)
	lrgen.CreateTables()
	lrgen.CFSM().CFSM2GraphViz("cfsm-" + "G4" + ".dot")
	// parser := Create(g, lrgen.GotoTable(), lrgen.ActionTable(), lrgen.AcceptingStates())
	// r := strings.NewReader("+ +")
	// scanner := NewStdScanner(r)
	// traceOn()
	// parser.Parse(lrgen.CFSM().S0, scanner)
}

func TestParser5(t *testing.T) {
	gtrace.SyntaxTracer = gologadapter.New()
	b := lr.NewGrammarBuilder("G5")
	b.LHS("S").N("E").EOF()
	b.LHS("E").N("E").T("+", '+').N("T").End()
	b.LHS("E").N("T").End()
	b.LHS("T").T("a", scanner.Ident).End()
	b.LHS("T").T("(", '(').N("E").T(")", ')').End()
	g := b.Grammar()
	g.Dump()
	ga := lr.NewGrammarAnalysis(g)
	lrgen := lr.NewLRTableGenerator(ga)
	lrgen.CreateTables()
	lrgen.CFSM().CFSM2GraphViz("cfsm-" + "G5" + ".dot")
}
