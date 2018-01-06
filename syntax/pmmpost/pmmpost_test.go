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
