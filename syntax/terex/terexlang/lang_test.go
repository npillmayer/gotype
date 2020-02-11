package terexlang

import (
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/npillmayer/gotype/syntax/termr"
)

func TestAssignability(t *testing.T) {
	var e interface{} = &sExprTermR{name: "Hello"}
	switch x := e.(type) {
	case termr.TermR:
		t.Logf("sExprOp %v accepted as termr.TermR", x)
		switch o := x.Operator().(type) {
		case terex.Operator:
			t.Logf("sExprOp.Operator() %v accepted as terex.Operator", o)
		default:
			t.Errorf("Expected %v to implement terex.Operator interface", o)
		}
	default:
		t.Errorf("Expected terexlang.sExprOp to implement termr.TermR interface")
	}
}

func TestMatchAnyOp(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	initDefaultPatterns()
	l1 := terex.List(makeASTTermR("S", "start").Operator(), 1)
	t.Logf("l1 = %s, pattern = %s", l1.ListString(), SingleTokenArg.ListString())
	if !SingleTokenArg.Match(l1, terex.GlobalEnvironment) {
		t.Errorf("Expected l1 to match pattern (Op any)")
	}
}

func TestMatchAnything(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	initDefaultPatterns()
	l := terex.List(1, 2, 3)
	if !termr.Anything().Match(l, terex.GlobalEnvironment) {
		t.Errorf("Expected !Anything to match (1 2 3)")
	}
}

func TestParse(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	terex.InitGlobalEnvironment()
	input := "(Hello 'World 1)"
	// lexer, err := Lexer()
	// if err != nil {
	// 	t.Error(err)
	// }
	// scan, err := lexer.Scanner(input)
	// if err != nil {
	// 	t.Error(err)
	// }
	//gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	lexer, _ := Lexer()
	parser := createParser()
	scan, _ := lexer.Scanner(input)
	accept, err := parser.Parse(scan, nil)
	t.Logf("accept=%v, input=%s", accept, input)
	if err != nil {
		t.Error(err)
	}
	if !accept {
		t.Errorf("No accept. Not a valid TeREx expression")
	}
	t.Fail()
}

func TestQuote(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	terex.InitGlobalEnvironment()
	input := "(Hello 'World 1)"
	//lexer, _ := Lexer()
	parser := createParser()
	scan, err := lexer.Scanner(input)
	if err != nil {
		t.Error(err)
	}
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	accept, err := parser.Parse(scan, nil)
	T().Errorf("accept=%v, input=%s", accept, input)
	// result := Eval(input, terex.GlobalEnvironment)
	// if result.Length() != 2 {
	// 	t.Errorf("Expected resulting list to be of length 2, is %d", result.Length())
	// } else {
	// 	t.Logf("AST=%s", result.Car.Data.(*terex.GCons).ListString())
	// 	if result.Car.Data.(*terex.GCons).Length() != 3 {
	// 		t.Errorf("Expected AST to be of length 3, is %d", result.Cadr().Length())
	// 	}
	// }
	t.Fail()
}

func TestEval(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	terex.InitGlobalEnvironment()
	input := "(+ 1 2 3)"
	result := Eval(input, terex.GlobalEnvironment)
	if result.Length() != 2 {
		t.Errorf("Expected resulting list to be of length 2, is %d", result.Length())
	}
	t.Logf("AST=%s", result.Car.Data.(*terex.GCons).ListString())
	if result.Car.Data.(*terex.GCons).Length() != 3 {
		t.Errorf("Expected AST to be of length 3, is %d", result.Cadr().Length())
	}
}
