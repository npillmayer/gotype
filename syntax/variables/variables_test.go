package variables_test

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	"github.com/npillmayer/gotype/syntax/variables/varparse"
	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/gologadapter"
	"github.com/shopspring/decimal"
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
	log.Printf("line %s: %.44s", at, msg)
	parserr = msg
}

func newErrL() antlr.ErrorListener {
	return &TestingErrorListener{}
}

func checkErr(t *testing.T) {
	if parserr != "" {
		t.Errorf(parserr)
	}
}

// Init the global tracers.
func TestInit0(t *testing.T) {
	tracing.InterpreterTracer = gologadapter.New()
	tracing.SyntaxTracer = gologadapter.New()
}

func TestVarDecl1(t *testing.T) {
	symtab := runtime.NewSymbolTable(variables.NewPMMPVarDecl)
	symtab.DefineSymbol("x")
}

func TestVarDecl2(t *testing.T) {
	symtab := runtime.NewSymbolTable(variables.NewPMMPVarDecl)
	x, _ := symtab.DefineSymbol("x")
	variables.CreatePMMPVarDecl("r", variables.ComplexSuffix, x.(*variables.PMMPVarDecl))
}

func TestVarDecl3(t *testing.T) {
	symtab := runtime.NewSymbolTable(variables.NewPMMPVarDecl)
	x, _ := symtab.DefineSymbol("x")
	var v *variables.PMMPVarDecl = x.(*variables.PMMPVarDecl)
	variables.CreatePMMPVarDecl("r", variables.ComplexSuffix, v)
	arr := variables.CreatePMMPVarDecl("<array>", variables.ComplexArray, x.(*variables.PMMPVarDecl))
	variables.CreatePMMPVarDecl("a", variables.ComplexSuffix, arr)
	//var b *bytes.Buffer
	//b = v.ShowVariable(b)
	//fmt.Printf("## showvariable %s;\n%s\n", v.BaseTag.GetName(), b.String())
}

func TestVarRef1(t *testing.T) {
	x := variables.CreatePMMPVarDecl("x", variables.NumericType, nil)
	var v *variables.PMMPVarRef = variables.CreatePMMPVarRef(x, 1, nil)
	t.Logf("var ref: %v\n", v)
}

func TestVarRef2(t *testing.T) {
	x := variables.CreatePMMPVarDecl("x", variables.NumericType, nil)
	r := variables.CreatePMMPVarDecl("r", variables.ComplexSuffix, x)
	var v *variables.PMMPVarRef = variables.CreatePMMPVarRef(r, 1, nil)
	t.Logf("var ref: %v\n", v)
}

func TestVarRef3(t *testing.T) {
	x := variables.CreatePMMPVarDecl("x", variables.NumericType, nil)
	arr := variables.CreatePMMPVarDecl("<[]>", variables.ComplexArray, x)
	subs := []decimal.Decimal{decimal.New(7, 0)}
	var v *variables.PMMPVarRef = variables.CreatePMMPVarRef(arr, 7, subs)
	t.Logf("var ref: %v\n", v)
}

func TestVarRefParser1(t *testing.T) {
	listener.ParseVariableFromString("x@", newErrL())
	checkErr(t)
}

func TestVarRefParser2(t *testing.T) {
	listener.ParseVariableFromString("x1@", newErrL())
	checkErr(t)
}

func TestVarRefParser3(t *testing.T) {
	listener.ParseVariableFromString("x.a@", newErrL())
	checkErr(t)
}

func TestVarRefParser4(t *testing.T) {
	listener.ParseVariableFromString("xyz18abc@", newErrL())
	checkErr(t)
}

func TestVarRefParser5(t *testing.T) {
	listener.ParseVariableFromString("xyz18.abc@", newErrL())
	checkErr(t)
}

func TestVarRefParser6(t *testing.T) {
	listener.ParseVariableFromString("x1a[2]b@", newErrL())
	checkErr(t)
}
