package parser

import (
	"strings"
	"testing"
	"text/scanner"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/dss"
)

func traceOn() {
	T.SetLevel(tracing.LevelDebug)
}

func TestStackSet1(t *testing.T) {
	set1 := newStackSet()
	r := dss.NewRoot("G", -999)
	s1 := dss.NewStack(r)
	set1 = set1.add(s1)
	T.Debugf("set = %v", set1)
	if set1.get() != s1 {
		t.Fail()
	}
}

func TestParser1(t *testing.T) {
	b := lr.NewGrammarBuilder("G")
	b.LHS("S").N("A").T("a", scanner.Ident).EOF()
	b.LHS("A").N("B").End()
	b.LHS("A").N("C").End()
	b.LHS("B").T("+", '+').End()
	b.LHS("B").Epsilon()
	b.LHS("C").T("-", '-').End()
	b.LHS("C").Epsilon()
	g := b.Grammar()
	g.Dump()
	ga := lr.NewGrammarAnalysis(g)
	lrgen := lr.NewLRTableGenerator(ga)
	lrgen.CreateTables()
	parser := Create(g, lrgen.GotoTable(), lrgen.ActionTable(), lrgen.AcceptingStates())
	r := strings.NewReader("+a")
	scanner := NewStdScanner(r)
	//traceOn()
	parser.Parse(lrgen.CFSM().S0, scanner)
	/*
		lrgen.CFSM().CFSM2GraphViz("/tmp/cfsm-" + "G1" + ".dot")
		tmp, _ := ioutil.TempFile("", "lr_")
		T.Infof("writing HTML to %s", tmp.Name())
		lr.ActionTableAsHTML(lrgen, tmp)
		lr.GotoTableAsHTML(lrgen, tmp)
		tmp.Close()
	*/
}
