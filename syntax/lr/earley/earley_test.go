package earley

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
)

// We use a small unambiguous expression grammar for testing.
// It is slightly adapted from
//
//      http://loup-vaillant.fr/tutorials/earley-parsing/recogniser
//
// This way we will be able to follow the examples there.
//
//     Sum     = Sum     '+' Product
//             | Product
//     Product = Product '*' Factor
//             | Factor
//     Factor  = '(' Sum ')'
//             | number
//
// 'number' is a terminal symbol recognizing Go integers.
//
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
		t.Error(err)
	}
	ga := lr.Analysis(g)
	if ga == nil {
		t.Errorf("Could not analyze grammar")
	}
	return ga
}

func makeParser(t *testing.T, no int, input string) (*Parser, scanner.Tokenizer) {
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	reader := strings.NewReader(input)
	scanner := scanner.GoTokenizer(fmt.Sprintf("test #%d", no), reader)
	ga := makeGrammar(t)
	return NewParser(ga), scanner
}

var inputStrings = []string{
	"1", "1+2", "1*2", "1+2*3", "1*(2+3)", "1+2+3+4", "1*2+3*4",
}

// --- the tests -------------------------------------------------------------

func TestParser1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	for n, input := range inputStrings {
		T().Infof("=== '%s' ========================", input)
		parser, scanner := makeParser(t, 1, input)
		gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
		accept, err := parser.Parse(scanner, nil)
		if err != nil {
			t.Error(err)
		}
		if !accept {
			t.Errorf("Valid input string #%d not accepted: '%s'", n+1, input)
		}
	}
}

func TestTree1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	input := "1+2*3"
	parser, scanner := makeParser(t, 1, input)
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	accept, err := parser.Parse(scanner, nil)
	if err != nil {
		t.Error(err)
	}
	if !accept {
		t.Errorf("Valid input string not accepted: '%s'", input)
	}
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	v := parser.TreeWalk(&ExprListener{0, t})
	value, ok := v.Value.(int)
	if !ok || value != 7 {
		t.Errorf("Expected %s to be 7, is %d", input, value)
	}
}

// ---------------------------------------------------------------------------

type ExprListener struct {
	total int
	t     *testing.T
}

func (el *ExprListener) Reduce(lhs *lr.Symbol, children []*RuleNode, extent span,
	level int) interface{} {
	//
	switch lhs.Name {
	case "Sum":
		return el.ReduceSum(lhs, children, level)
	case "Product":
		return el.ReduceProduct(lhs, children, level)
	case "Factor":
		return el.ReduceFactor(lhs, children, level)
	}
	el.t.Logf("%sReduce of grammar symbol: %v", indent(level), lhs)
	return children[0].Value
}

func (el *ExprListener) ReduceSum(lhs *lr.Symbol, children []*RuleNode, level int) interface{} {
	v := children[0].Value // Product
	if len(children) > 1 {
		v = children[0].Value.(int) + children[2].Value.(int) // Sum + Product
	}
	el.t.Logf("%sSUM %v\n", indent(level), v)
	return v
}

func (el *ExprListener) ReduceProduct(lhs *lr.Symbol, children []*RuleNode, level int) interface{} {
	v := children[0].Value // Factor
	if len(children) > 1 {
		v = children[0].Value.(int) * children[2].Value.(int) // Product * Factor
	}
	el.t.Logf("%sPRODUCT %v\n", indent(level), v)
	return v
}

func (el *ExprListener) ReduceFactor(lhs *lr.Symbol, children []*RuleNode, level int) interface{} {
	v := children[0].Value // number
	if len(children) > 1 {
		v = children[1].Value // ( Sum )
	}
	el.t.Logf("%sFACTOR %v\n", indent(level), v)
	return v
}

func (el *ExprListener) Terminal(tokval int, token interface{}, extent span, level int) interface{} {
	el.t.Logf("%sToken %s|%d\n", indent(level), scanner.Lexeme(token), tokval)
	if tokval == scanner.Int {
		n, _ := strconv.Atoi(scanner.Lexeme(token))
		return n
	}
	return tokval
}

func indent(level int) string {
	in := ""
	for level > 0 {
		in = in + ". "
		level--
	}
	return in
}
