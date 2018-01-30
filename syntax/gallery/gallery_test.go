package gallery

import (
	"fmt"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/variables"
	"github.com/shopspring/decimal"
)

var syntrace tracing.Trace = tracing.SyntaxTracer

// Helper
func createInterpreter(s string, builtins bool) *GalleryInterpreter {
	intp := NewGalleryInterpreter(builtins)
	input := antlr.NewInputStream(s)
	intp.ASTListener.LazyCreateParser(input)
	return intp
}

// ---------------------------------------------------------------------------

func TestParseVariable(t *testing.T) {
	T.SetLevel(tracing.LevelError)
	intp := createInterpreter("hello", false) // variable reference
	tree := intp.ASTListener.statemParser.Variable()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### variable = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	if sym, _ := intp.runtime.ScopeTree.Globals().ResolveSymbol("hello"); sym == nil {
		t.Fail()
	}
}

func TestParseDecimal(t *testing.T) {
	intp := createInterpreter("3.14", false) // numeric token atom
	tree := intp.ASTListener.statemParser.Atom()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### atom = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	_, ok := intp.runtime.ExprStack.PopAsNumeric()
	if !ok {
		t.Fail()
	}
}

func TestParseNumtokenatom(t *testing.T) {
	intp := createInterpreter("3/14", false) // numeric token atom
	tree := intp.ASTListener.statemParser.Numtokenatom()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### numtokenatom = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	_, ok := intp.runtime.ExprStack.PopAsNumeric()
	if !ok {
		t.Fail()
	}
}

func TestParsePair(t *testing.T) {
	intp := createInterpreter("(1,0.5)", false) // literal pair
	tree := intp.ASTListener.statemParser.Atom()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### atom = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	_, ok := intp.runtime.ExprStack.PopAsPair()
	if !ok {
		t.Fail()
	}
}

func TestParseShowcmd(t *testing.T) {
	intp := createInterpreter("show a, b", false)
	tree := intp.ASTListener.statemParser.Command()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### show = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
}

func TestParseAssignment(t *testing.T) {
	intp := createInterpreter("a := 1", false) // assign a numeric variable
	tree := intp.ASTListener.statemParser.Assignment()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### assignment = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	sym := intp.runtime.MemFrameStack.Globals().Symbols().ResolveSymbol("a")
	if sym == nil {
		t.Fail()
	}
	var a *variables.PMMPVarRef = sym.(*variables.PMMPVarRef)
	if !a.GetValue().(decimal.Decimal).Equal(arithmetic.ConstOne) {
		t.Fail()
	}
}

func TestParseScalarmulop1(t *testing.T) {
	intp := createInterpreter("3a", false) // scale a variable
	tree := intp.ASTListener.statemParser.Primary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### primary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	sym := intp.runtime.MemFrameStack.Globals().Symbols().ResolveSymbol("a")
	if sym == nil {
		t.Fail()
	}
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	coeff := tos.XPolyn.GetCoeffForTerm(sym.GetID())
	if !coeff.Equal(decimal.New(3, 0)) {
		t.Fail()
	}
}

func TestParseScalarmulop2(t *testing.T) {
	intp := createInterpreter("-3/2(4,8)", false) // scale a literal pair
	tree := intp.ASTListener.statemParser.Primary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### primary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.YPolyn.GetConstantValue().Equal(decimal.New(-12, 0)) {
		t.Fail()
	}
}

func TestParseExprgroup(t *testing.T) {
	intp := createInterpreter("begingroup 1 endgroup", false) // group returning 1
	tree := intp.ASTListener.statemParser.Atom()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### atom = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.XPolyn.GetConstantValue().Equal(arithmetic.ConstOne) {
		t.Fail()
	}
}

func TestParseMathfunc(t *testing.T) {
	intp := createInterpreter("floor 3.14", false) // should yield 3
	tree := intp.ASTListener.statemParser.Primary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### primary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.XPolyn.GetConstantValue().Equal(decimal.New(3, 0)) {
		t.Fail()
	}
}

func TestParsePairpart(t *testing.T) {
	intp := createInterpreter("ypart (3,1)", false) // should yield 1
	tree := intp.ASTListener.statemParser.Primary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### primary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.XPolyn.GetConstantValue().Equal(arithmetic.ConstOne) {
		t.Fail()
	}
}

func TestParseNumericSecondary(t *testing.T) {
	intp := createInterpreter("3 * 2", false) // should yield 6
	tree := intp.ASTListener.statemParser.Secondary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### secondary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.XPolyn.GetConstantValue().Equal(decimal.New(6, 0)) {
		t.Fail()
	}
}

func TestParsePairSecondary(t *testing.T) {
	intp := createInterpreter("(3,6)/2", false) // should yield (1.5,3)
	tree := intp.ASTListener.statemParser.Secondary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### secondary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.YPolyn.GetConstantValue().Equal(decimal.New(3, 0)) {
		t.Fail()
	}
}

func TestParseNumericTertiary(t *testing.T) {
	intp := createInterpreter("3 - 2", false) // should yield 1
	tree := intp.ASTListener.statemParser.Tertiary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### tertiary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.XPolyn.GetConstantValue().Equal(arithmetic.ConstOne) {
		t.Fail()
	}
}

func TestParseEquation(t *testing.T) {
	fmt.Println("-------------------------------------")
	intp := createInterpreter("2a = b = 2", false) // minimal equation
	tree := intp.ASTListener.statemParser.Equation()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### equation = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	sym := intp.runtime.MemFrameStack.Globals().Symbols().ResolveSymbol("a")
	if sym == nil {
		t.Fail()
	}
	var a *variables.PMMPVarRef = sym.(*variables.PMMPVarRef)
	if !a.GetValue().(decimal.Decimal).Equal(arithmetic.ConstOne) {
		t.Fail()
	}
}

func TestParseEquation2(t *testing.T) {
	//tracing.SyntaxTracer.SetLevel(tracing.LevelDebug)
	//config.Initialize()
	fmt.Println("-------------------------------------")
	intp := createInterpreter("p = a * (1,2); ypart p=4;", true)
	tree := intp.ASTListener.statemParser.Statementlist()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### equation = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	sym := intp.runtime.MemFrameStack.Globals().Symbols().ResolveSymbol("p")
	if sym == nil {
		t.Fail()
	}
	//intp.runtime.ExprStack.Summary()
	var p *variables.PMMPVarRef = sym.(*variables.PMMPVarRef)
	if !p.XPart().IsKnown() {
		t.Fail()
	}
}

func TestPathBuilding1(t *testing.T) {
	//T.SetLevel(tracing.LevelDebug)
	intp := createInterpreter("(0,0) -- (1,2) -- (4,2) -- (3,0) -- cycle", false)
	tree := intp.ASTListener.statemParser.Path()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### path = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
}

func TestPathBuilding2(t *testing.T) {
	T.SetLevel(tracing.LevelInfo)
	syntrace.SetLevel(tracing.LevelInfo)
	intp := createInterpreter("(0,0) -- (1,2) shifted (4,3) rotated 45 -- cycle", false)
	//intp := createInterpreter("(0,0) -- (1,2) shifted (4,3) -- cycle")
	tree := intp.ASTListener.statemParser.Path()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### path = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
}
