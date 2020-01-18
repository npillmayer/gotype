package bidi

import (
	"strings"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/parser"
)

func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

type BiDiGrammar struct {
	scanner parser.Scanner
	parser  *parser.Parser
	lrgen   *lr.LRTableGenerator
}

func NewGrammar() *BiDiGrammar {
	bidi := &BiDiGrammar{}
	b := lr.NewGrammarBuilder("UAX#9")
	b.LHS("R").N("AL").N("NSM").EOF()
	b.LHS("AL").T("a", scanner.Ident).End()
	b.LHS("NSM").T("b", scanner.Ident).End()
	//b.LHS("AL").N("AL").N("NSM").End()
	//b.LHS("R").N("R").N("NSM").End()
	//b.LHS("L").N("L").N("NSM").End()
	g := b.Grammar()
	ga := lr.NewGrammarAnalysis(g)
	lrgen := lr.NewLRTableGenerator(ga)
	lrgen.CreateTables()
	parser := parser.Create(g, lrgen.GotoTable(), lrgen.ActionTable(), lrgen.AcceptingStates())
	bidi.parser = parser
	bidi.lrgen = lrgen
	return bidi
}

func (bidi *BiDiGrammar) Parse(input string) {
	r := strings.NewReader(input)
	scanner := parser.NewStdScanner(r)
	T().Debugf("===========================================================")
	bidi.parser.Parse(bidi.lrgen.CFSM().S0, scanner)
}
