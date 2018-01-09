package gallery

import (
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/syntax/variables"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// Helper
func createInterpreter(s string) *GalleryInterpreter {
	intp := NewGalleryInterpreter(false)
	input := antlr.NewInputStream(s)
	intp.ASTListener.LazyCreateParser(input)
	return intp
}

// ---------------------------------------------------------------------------

func TestParseVariable(t *testing.T) {
	T.SetLevel(logrus.ErrorLevel)
	intp := createInterpreter("hello") // variable reference
	tree := intp.ASTListener.statemParser.Variable()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### variable = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	if sym, _ := intp.runtime.ScopeTree.Globals().ResolveSymbol("hello"); sym == nil {
		t.Fail()
	}
}

func TestParseDecimal(t *testing.T) {
	intp := createInterpreter("3.14") // numeric token atom
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
	intp := createInterpreter("3/14") // numeric token atom
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
	intp := createInterpreter("(1,0.5)") // literal pair
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
	intp := createInterpreter("show a, b")
	tree := intp.ASTListener.statemParser.Command()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### show = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
}

func TestParseAssignment(t *testing.T) {
	intp := createInterpreter("a := 1") // assign a numeric variable
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
	intp := createInterpreter("3a") // scale a variable
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
	coeff := tos.GetXPolyn().GetCoeffForTerm(sym.GetID())
	if !coeff.Equal(decimal.New(3, 0)) {
		t.Fail()
	}
}

func TestParseScalarmulop2(t *testing.T) {
	intp := createInterpreter("-3/2(4,8)") // scale a literal pair
	tree := intp.ASTListener.statemParser.Primary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### primary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.GetYPolyn().GetConstantValue().Equal(decimal.New(-12, 0)) {
		t.Fail()
	}
}

func TestParseExprgroup(t *testing.T) {
	intp := createInterpreter("begingroup 1 endgroup") // group returning 1
	tree := intp.ASTListener.statemParser.Atom()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### atom = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.GetXPolyn().GetConstantValue().Equal(arithmetic.ConstOne) {
		t.Fail()
	}
}

func TestParseMathfunc(t *testing.T) {
	intp := createInterpreter("floor 3.14") // should yield 3
	tree := intp.ASTListener.statemParser.Primary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### primary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.GetXPolyn().GetConstantValue().Equal(decimal.New(3, 0)) {
		t.Fail()
	}
}

func TestParsePairpart(t *testing.T) {
	intp := createInterpreter("ypart (3,1)") // should yield 1
	tree := intp.ASTListener.statemParser.Primary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### primary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.GetXPolyn().GetConstantValue().Equal(arithmetic.ConstOne) {
		t.Fail()
	}
}

func TestParseNumericSecondary(t *testing.T) {
	intp := createInterpreter("3 * 2") // should yield 6
	tree := intp.ASTListener.statemParser.Secondary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### secondary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.GetXPolyn().GetConstantValue().Equal(decimal.New(6, 0)) {
		t.Fail()
	}
}

func TestParsePairSecondary(t *testing.T) {
	intp := createInterpreter("(3,6)/2") // should yield (1.5,3)
	tree := intp.ASTListener.statemParser.Secondary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### secondary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.GetYPolyn().GetConstantValue().Equal(decimal.New(3, 0)) {
		t.Fail()
	}
}

func TestParseNumericTertiary(t *testing.T) {
	intp := createInterpreter("3 - 2") // should yield 1
	tree := intp.ASTListener.statemParser.Tertiary()
	sexpr := antlr.TreesStringTree(tree, nil, intp.ASTListener.statemParser)
	T.Debugf("### tertiary = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(intp.ASTListener, tree)
	tos := intp.runtime.ExprStack.Top()
	T.Debugf("TOS = %v", tos)
	if !tos.GetXPolyn().GetConstantValue().Equal(arithmetic.ConstOne) {
		t.Fail()
	}
}

func TestParseEquation(t *testing.T) {
	intp := createInterpreter("2a = b = 2") // minimal equation
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
