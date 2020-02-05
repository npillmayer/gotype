package termr

import (
	"text/scanner"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
)

// ASTBuilder is a parse tree listener for building ASTs.
// AST is a homogenous abstract syntax tree.
type ASTBuilder struct {
	G                *lr.Grammar  // input grammar the parse forest stems from
	forest           *sppf.Forest // input parse forest
	ast              *GCons       // root of the AST to construct
	last             *GCons       // current last node to append conses
	dispatchEnter    func(string) ruleABEnter
	dispatchExit     func(string) ruleABExit
	conflictStrategy sppf.Pruner
	Error            func(error)
	stack            []*GCons
}

type ruleABEnter func(sym *lr.Symbol, rhs []*sppf.RuleNode) bool
type ruleABExit func(sym *lr.Symbol, rhs []*sppf.RuleNode) interface{}

// NewASTBuilder creates an AST builder from a parse forest/tree.
func NewASTBuilder(g *lr.Grammar) *ASTBuilder {
	ab := &ASTBuilder{
		G:             g,
		ast:           &GCons{carNode{NullAtom, nil}, nil}, // AST anchor
		stack:         make([]*GCons, 0, 256),
		dispatchEnter: nullABEnter,
		dispatchExit:  nullABExit,
	}
	ab.last = ab.ast
	ab.stack = append(ab.stack, ab.ast) // push as stopper
	return ab
}

func nullABEnter(string) ruleABEnter {
	return nil
}

func nullABExit(string) ruleABExit {
	return nil
}

// AST creates an abstract syntax tree from a parse tree/forest.
func (ab *ASTBuilder) AST(parseTree *sppf.Forest) (*GCons, interface{}) {
	if parseTree == nil {
		return nil, nil
	}
	ab.forest = parseTree
	cursor := ab.forest.SetCursor(nil, nil) // TODO set Pruner
	value := cursor.TopDown(ab, sppf.LtoR, sppf.Break)
	return ab.ast, value
}

func (ab *ASTBuilder) append(cons *GCons) {
	if cons == nil {
		return
	}
	ab.last.cdr = cons
	ab.last = cons
}

func (ab *ASTBuilder) up() {
	if ab.stack[len(ab.stack)-1].car.atom == NullAtom {
		T().Errorf("up() called with empty stack")
		panic("empty stack")
	}
	ab.stack = ab.stack[:len(ab.stack)-1] // pop last
}

func (ab *ASTBuilder) down(car carNode) {
	ab.stack = append(ab.stack, ab.last)
	child := &GCons{car, nil}
	ab.last.cdr = &GCons{carNode{NullAtom, child}, nil}
	ab.last = child
}

// --- sppf.Listener interface -----------------------------------------------

// EnterRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) EnterRule(sym *lr.Symbol, rhs []*sppf.RuleNode, span lr.Span, level int) bool {
	if r := ab.dispatchEnter(sym.Name); r != nil {
		return r(sym, rhs)
	}
	T().Debugf("enter grammar symbol: %v", sym)
	return true
}

// ExitRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) ExitRule(sym *lr.Symbol, rhs []*sppf.RuleNode, span lr.Span, level int) interface{} {
	if r := ab.dispatchExit(sym.Name); r != nil {
		return r(sym, rhs)
	}
	if len(rhs) > 0 {
		return rhs[0].Value
	}
	T().Debugf("exit grammar symbol: %v", sym)
	return nil
}

// Terminal is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Terminal(tokval int, token interface{}, span lr.Span, level int) interface{} {
	if tokval == scanner.EOF { // filter out artificial EOF
		return nil
	}
	t := ab.G.Terminal(tokval).Name
	car := atomize(t)
	car.atom.typ = TokenType
	T().Debugf("cons(terminal=%d) = %s", tokval, car)
	ab.append(&GCons{car, nil})
	return nil
}

// Conflict is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Conflict(sym *lr.Symbol, rule int, span lr.Span, level int) (int, error) {
	panic("Conflict of AST building not yet implemented")
	//return 0, nil
}
