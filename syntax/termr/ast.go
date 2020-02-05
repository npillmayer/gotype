package termr

import (
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
	dispatchExit     func(string) ruleAB
	conflictStrategy sppf.Pruner
	Error            func(error)
}

type ruleABEnter func(sym *lr.Symbol, rhs []*sppf.RuleNode) bool
type ruleAB func(sym *lr.Symbol, rhs []*sppf.RuleNode) interface{}

// NewASTBuilder creates an AST builder from a parse forest/tree.
func NewASTBuilder(g *lr.Grammar, parseTree *sppf.Forest) *ASTBuilder {
	ab := &ASTBuilder{
		G:      g,
		forest: parseTree,
	}
	return ab
}

func (ab *ASTBuilder) append(cons *GCons) {
	if cons == nil {
		return
	}
	if ab.last != nil {
		ab.last.cdr = cons
	}
	ab.last = cons
	if ab.ast == nil {
		ab.ast = ab.last
	}
}

func (ab *ASTBuilder) up() {
}

func (ab *ASTBuilder) down(car *carNode) {
	//
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
	// var values *iteratable.Set
	// for _, node := range rhs {
	// 	switch v := node.Value.(type) {
	// 	case *iteratable.Set:
	// 		if values == nil {
	// 			values = v
	// 		} else {
	// 			values = values.Union(v)
	// 		}
	// 	default:
	// 		if values == nil {
	// 			values = iteratable.NewSet(0)
	// 		}
	// 		values.Add(v)
	// 	}
	// }
	T().Debugf("exit grammar symbol: %v", sym)
	//return values
	return nil
}

// Terminal is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Terminal(tokval int, token interface{}, span lr.Span, level int) interface{} {
	car := atomize(token)
	return nil
}

// Conflict is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Conflict(sym *lr.Symbol, rule int, span lr.Span, level int) (int, error) {
	return 0, nil
}
