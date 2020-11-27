package terexlang

import (
	"io/ioutil"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/npillmayer/gotype/syntax/terex/termr"
)

func TestScanner(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	lex, _ := Lexer()
	input := "a + '(1.2 world !) #var nil ;"
	scan, err := lex.Scanner(input)
	if err != nil {
		t.Errorf(err.Error())
	}
	scan.SetErrorHandler(func(e error) {
		t.Error(e)
	})
	done := false
	for !done {
		tokval, token, _, _ := scan.NextToken(nil)
		if tokval == -1 {
			done = true
		} else {
			t.Logf("token = %v with value = %d", token, tokval)
		}
	}
}

func TestAssignability(t *testing.T) {
	var e interface{} = &sExprTermR{name: "Hello"}
	switch x := e.(type) {
	case termr.TermR:
		t.Logf("sExprTermR %v accepted as termr.TermR", x)
		switch o := x.Operator().(type) {
		case terex.Operator:
			t.Logf("sExprTermR.Operator() %v accepted as terex.Operator", o)
		default:
			t.Errorf("Expected %v to implement terex.Operator interface", o)
		}
	default:
		t.Errorf("Expected terexlang.sExprTermR to implement termr.TermR interface")
	}
}

// func TestMatchAnyOp(t *testing.T) {
// 	gtrace.SyntaxTracer = gotestingadapter.New()
// 	teardown := gotestingadapter.RedirectTracing(t)
// 	defer teardown()
// 	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
// 	initTokens()
// 	initDefaultPatterns()
// 	l1 := terex.List(makeASTTermR("S", "start").Operator(), 1)
// 	t.Logf("l1 = %s, pattern = %s", l1.ListString(), SingleTokenArg.ListString())
// 	if !SingleTokenArg.Match(l1, terex.GlobalEnvironment) {
// 		t.Errorf("Expected l1 to match pattern (Op any)")
// 	}
// }

func TestMatchAnything(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	initRewriters()
	l := terex.List(1, 2, 3)
	if !termr.Anything().Match(l, terex.GlobalEnvironment) {
		t.Errorf("Expected !Anything to match (1 2 3)")
	}
}

func TestParse(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	terex.InitGlobalEnvironment()
	input := `a`
	//input := "(Hello 'World 1)"
	parser := createParser()
	scan, _ := lexer.Scanner(input)
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	accept, err := parser.Parse(scan, nil)
	t.Logf("accept=%v, input=%s", accept, input)
	if err != nil {
		t.Error(err)
	}
	if !accept {
		t.Errorf("No accept. Not a valid TeREx expression")
	}
	parsetree := parser.ParseForest()
	tmpfile, err := ioutil.TempFile(".", "eval-parsetree-*.dot")
	if err != nil {
		t.Error("cannot open tmp file for graphviz output")
	}
	sppf.ToGraphViz(parsetree, tmpfile)
	T().Infof("Exported parse tree to %s", tmpfile.Name())
}

func TestAST(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	terex.InitGlobalEnvironment()
	input := `a`
	//input := `(+ '(1 "Hi" 3) 4)`
	parsetree, retr, err := Parse(input)
	if err != nil {
		t.Error(err)
	}
	if parsetree == nil || retr == nil {
		t.Errorf("parse tree or  token retriever is nil")
	}
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	T().Infof("####################################################")
	ab := newASTBuilder()
	env := ab.AST(parsetree, retr)
	if env == nil {
		t.Errorf("Cannot create AST from parsetree")
	}
	ast := env.AST
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	T().Infof("AST: %s", ast.ListString())
	T().Infof("####################################################")
	// terseAst := terex.GlobalEnvironment.Quote(ast)
	// T().Infof("reduced AST: %s", terseAst.ListString())
}

func TestQuoteAST(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	terex.InitGlobalEnvironment()
	//input := `(Hello 'World (+ 1 2) "string")`
	input := `a`
	tree, retr, err := Parse(input)
	if err != nil {
		t.Errorf(err.Error())
	}
	ast, env, err := AST(tree, retr)
	//t.Logf("\n\n" + debugString(terex.Elem(ast.Car)))
	t.Logf("\n\n" + debugString(terex.Elem(ast)))
	env.Def("a", terex.Atomize(7))
	q, err := QuoteAST(terex.Elem(ast.Car), env)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("\n\n" + debugString(q))
	t.Fail()
}

func debugString(e terex.Element) string {
	if e.IsAtom() {
		return e.String()
	}
	return e.AsList().IndentedListString()
}

/*
func TestQuote(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	terex.InitGlobalEnvironment()
	input := "(Hello 'World 1)"
	result, _ := Quote(input)
	if result.Length() != 2 {
		t.Errorf("Expected resulting list to be of length 2, is %d", result.Length())
	} else {
		t.Logf("AST=%s", result.Car.Data.(*terex.GCons).ListString())
		if result.Car.Data.(*terex.GCons).Length() != 3 {
			t.Errorf("Expected AST to be of length 3, is %d", result.Cadr().Length())
		}
	}
}
*/

func NoTestEval(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelError)
	terex.InitGlobalEnvironment()
	sym := terex.GlobalEnvironment.FindSymbol("+", false)
	if sym == nil {
		t.Error("Expected to find operator '+' in global environment")
	}
	input := "(+ 1 2 3)"
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	result, env := Eval(input)
	if result.Length() != 2 {
		t.Errorf("Expected resulting list to be of length 2, is %d", result.Length())
	} else {
		ast := env.AST.Cadr()
		t.Logf("AST=%s", result.Car.Data.(*terex.GCons).ListString())
		if ast.Length() != 3 {
			t.Errorf("Expected AST to be of length 3, is %d", ast.Length())
		}
	}
	t.Fail()
}
