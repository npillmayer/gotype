package parser

import (
	"testing"
	"text/scanner"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
)

func traceOn() {
	T.SetLevel(tracing.LevelDebug)
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
	parser := Create(g, lrgen.GotoTable(), lrgen.ActionTable())
	scanner := NewScanner(nil)
	traceOn()
	parser.Parse(lrgen.CFSM().S0, scanner)
}
