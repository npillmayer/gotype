package main

import (
	"fmt"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/syntax/mpost/parser"
)

type TreeShapeListener struct {
	*parser.BasearithmListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (this *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

func main() {
	fmt.Println("This is a test, MPOST syntax, draft")
	input, _ := antlr.NewFileStream(os.Args[1])
	lexer := parser.NewarithmLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewarithmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Equation()
	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}
