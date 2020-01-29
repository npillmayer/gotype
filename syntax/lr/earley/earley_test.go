package earley

import (
	"fmt"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
)

func makeGrammar(t *testing.T) *lr.LRAnalysis {
	b := lr.NewGrammarBuilder("Expressions")
	b.LHS("Sum").N("Sum").T("+", '+').N("Product").End()
	b.LHS("Sum").N("Product").End()
	b.LHS("Product").N("Product").T("*", '*').N("Factor").End()
	b.LHS("Product").N("Factor").End()
	b.LHS("Factor").T("(", '(').N("Sum").T(")", ')').End()
	b.LHS("Factor").T("1", '1').End()
	g, err := b.Grammar()
	if err != nil {
		T().Errorf(err.Error())
		t.Error(err)
	}
	ga := lr.Analysis(g)
	if ga == nil {
		t.Errorf("Could not analyze grammar")
	}
	return ga
}

func makeParser(t *testing.T, no int, input string) (*Parser, scanner.Tokenizer) {
	reader := strings.NewReader(input)
	scanner := scanner.GoTokenizer(fmt.Sprintf("test #%d", no), reader)
	ga := makeGrammar(t)
	return NewParser(ga), scanner
}

func TestParser1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	parser, scanner := makeParser(t, 1, "1")
	parser.Parse(scanner)
}
