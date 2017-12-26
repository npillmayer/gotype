package pmmpost

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/syntax"
	"github.com/shopspring/decimal"
)

func TestVarDecl1(t *testing.T) {
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
	var b *bytes.Buffer
	b = v.ShowVariable(b)
	fmt.Printf("## showvariable %s;\n%s\n", v.BaseTag.GetName(), b.String())
}

func TestVarRef1(t *testing.T) {
	x := CreatePMMPVarDecl("x", Numeric, nil)
	var v *PMMPVarRef = CreatePMMPVarRef(x, 1, nil)
	t.Logf("var ref: %v\n", v)
}

func TestVarRef2(t *testing.T) {
	x := CreatePMMPVarDecl("x", Numeric, nil)
	r := CreatePMMPVarDecl("r", ComplexSuffix, x)
	var v *PMMPVarRef = CreatePMMPVarRef(r, 1, nil)
	t.Logf("var ref: %v\n", v)
}

func TestVarRef3(t *testing.T) {
	x := CreatePMMPVarDecl("x", Numeric, nil)
	arr := CreatePMMPVarDecl("<[]>", ComplexArray, x)
	subs := []decimal.Decimal{decimal.New(7, 0)}
	var v *PMMPVarRef = CreatePMMPVarRef(arr, 7, subs)
	t.Logf("var ref: %v\n", v)
}

func TestVarRefNameCache(t *testing.T) {
	x := CreatePMMPVarDecl("x", Numeric, nil)
	var v *PMMPVarRef = CreatePMMPVarRef(x, 1, nil)
	_ = v.GetName()
	if len(v.cachedName) == 0 {
		t.Fail()
	}
}

func TestVarRefParser1(t *testing.T) {
	p := NewParseListener()
	p.ParseVarFromString("x@")
}

func TestVarRefParser2(t *testing.T) {
	p := NewParseListener()
	p.ParseVarFromString("x1@")
}

func TestVarRefParser3(t *testing.T) {
	p := NewParseListener()
	p.ParseVarFromString("x.a@")
}

func TestVarRefParser4(t *testing.T) {
	p := NewParseListener()
	p.ParseVarFromString("xyz18abc@")
}

func TestVarRefParser5(t *testing.T) {
	p := NewParseListener()
	p.ParseVarFromString("xyz18.abc@")
}

func TestVarRefParser6(t *testing.T) {
	p := NewParseListener()
	p.ParseVarFromString("x1a[2]b@")
}

func TestStatement1(t *testing.T) {
	log.Println("--------------------------------------------------")
	input := antlr.NewInputStream("numeric x;")
	p := NewParseListener()
	p.ParseStatements(input)
}

func TestStatement2(t *testing.T) {
	log.Println("--------------------------------------------------")
	input := antlr.NewInputStream("numeric x, y, a, b;")
	p := NewParseListener()
	p.ParseStatements(input)
}

func TestInterprTest1(t *testing.T) {
	log.Println("--------------------------------------------------")
	input := antlr.NewInputStream("numeric a;")
	p := NewParseListener()
	p.ParseStatements(input)
	p.Summary()
}

func TestInterprTest2(t *testing.T) {
	log.Println("--------------------------------------------------")
	input := antlr.NewInputStream("numeric a; begingroup save a; pair a; endgroup;")
	p := NewParseListener()
	p.ParseStatements(input)
	p.Summary()
}
