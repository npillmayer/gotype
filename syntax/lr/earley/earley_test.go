package earley

import (
	"fmt"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"

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
	b.LHS("Factor").T("number", scanner.Int).End()
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
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	reader := strings.NewReader(input)
	scanner := scanner.GoTokenizer(fmt.Sprintf("test #%d", no), reader)
	ga := makeGrammar(t)
	return NewParser(ga), scanner
}

var inputStrings = []string{
	"1", "1+2", "1*2", "1+2*3", "1*(2+3)",
}

func TestParser1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	for n, input := range inputStrings {
		T().Infof("=== '%s' ========================", input)
		parser, scanner := makeParser(t, 1, input)
		gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
		accept, err := parser.Parse(scanner)
		if err != nil {
			t.Error(err)
		}
		if !accept {
			t.Errorf("Valid input string #%d not accepted: '%s'", n+1, input)
		}
	}
}
