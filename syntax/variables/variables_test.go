package variables

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/syntax"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type TestingErrorListener struct {
	*antlr.DefaultErrorListener // use default as base class
}

var parserr string

/* Our error listener prints an error and set a global error flag.
 */
func (c *TestingErrorListener) SyntaxError(r antlr.Recognizer, sym interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	//
	at := fmt.Sprintf("%s:%s", strconv.Itoa(line), strconv.Itoa(column))
	T.P("line", at).Errorf("%.44s", msg)
	parserr = msg
}

func newErrL() antlr.ErrorListener {
	return &TestingErrorListener{}
}

func checkErr(t *testing.T) {
	if parserr != "" {
		parserr = ""
		T.Error(parserr)
		t.Fail()
	}
}

func TestVarDecl1(t *testing.T) {
	T.SetLevel(logrus.InfoLevel)
	symtab := syntax.NewSymbolTable(NewPMMPVarDecl)
	symtab.DefineSymbol("x")
}

func TestVarDecl2(t *testing.T) {
	symtab := syntax.NewSymbolTable(NewPMMPVarDecl)
	x, _ := symtab.DefineSymbol("x")
	CreatePMMPVarDecl("r", ComplexSuffix, x.(*PMMPVarDecl))
}

func TestVarDecl3(t *testing.T) {
	symtab := syntax.NewSymbolTable(NewPMMPVarDecl)
	x, _ := symtab.DefineSymbol("x")
	var v *PMMPVarDecl = x.(*PMMPVarDecl)
	CreatePMMPVarDecl("r", ComplexSuffix, v)
	arr := CreatePMMPVarDecl("<array>", ComplexArray, x.(*PMMPVarDecl))
	CreatePMMPVarDecl("a", ComplexSuffix, arr)
	//var b *bytes.Buffer
	//b = v.ShowVariable(b)
	//fmt.Printf("## showvariable %s;\n%s\n", v.BaseTag.GetName(), b.String())
}

func TestVarRef1(t *testing.T) {
	x := CreatePMMPVarDecl("x", NumericType, nil)
	var v *PMMPVarRef = CreatePMMPVarRef(x, 1, nil)
	t.Logf("var ref: %v\n", v)
}

func TestVarRef2(t *testing.T) {
	x := CreatePMMPVarDecl("x", NumericType, nil)
	r := CreatePMMPVarDecl("r", ComplexSuffix, x)
	var v *PMMPVarRef = CreatePMMPVarRef(r, 1, nil)
	t.Logf("var ref: %v\n", v)
}

func TestVarRef3(t *testing.T) {
	x := CreatePMMPVarDecl("x", NumericType, nil)
	arr := CreatePMMPVarDecl("<[]>", ComplexArray, x)
	subs := []decimal.Decimal{decimal.New(7, 0)}
	var v *PMMPVarRef = CreatePMMPVarRef(arr, 7, subs)
	t.Logf("var ref: %v\n", v)
}

func TestVarRefNameCache(t *testing.T) {
	x := CreatePMMPVarDecl("x", NumericType, nil)
	var v *PMMPVarRef = CreatePMMPVarRef(x, 1, nil)
	_ = v.GetName()
	if len(v.cachedName) == 0 {
		t.Fail()
	}
}

func TestVarRefParser1(t *testing.T) {
	ParseVariableFromString("x@", newErrL())
	checkErr(t)
}

func TestVarRefParser2(t *testing.T) {
	ParseVariableFromString("x1@", newErrL())
	checkErr(t)
}

func TestVarRefParser3(t *testing.T) {
	ParseVariableFromString("x.a@", newErrL())
	checkErr(t)
}

func TestVarRefParser4(t *testing.T) {
	ParseVariableFromString("xyz18abc@", newErrL())
	checkErr(t)
}

func TestVarRefParser5(t *testing.T) {
	ParseVariableFromString("xyz18.abc@", newErrL())
	checkErr(t)
}

func TestVarRefParser6(t *testing.T) {
	ParseVariableFromString("x1a[2]b@", newErrL())
	checkErr(t)
}
