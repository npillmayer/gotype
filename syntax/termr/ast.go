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
	operators        map[string]ASTOperator
	conflictStrategy sppf.Pruner
	Error            func(error)
	stack            []*GCons
}

type ruleABEnter func(sym *lr.Symbol, rhs []*sppf.RuleNode) bool
type ruleABExit func(sym *lr.Symbol, rhs []*sppf.RuleNode) interface{}

// NewASTBuilder creates an AST builder from a parse forest/tree.
func NewASTBuilder(g *lr.Grammar) *ASTBuilder {
	ab := &ASTBuilder{
		G:         g,
		ast:       &GCons{Node{NullAtom, nil}, nil}, // AST anchor
		stack:     make([]*GCons, 0, 256),
		operators: make(map[string]ASTOperator),
	}
	ab.last = ab.ast
	ab.stack = append(ab.stack, ab.ast) // push as stopper
	return ab
}

// ASTOperator is a type for an operator for AST creation and transformation
// (rewriting).
type ASTOperator interface {
	Name() string
	Rewrite(*GCons, *Environment) *GCons
	Descend(sppf.RuleCtxt) bool
}

// AddOperator adds an AST operator for a grammar symbol to the builder.
func (ab *ASTBuilder) AddOperator(op ASTOperator) {
	if op != nil {
		ab.operators[op.Name()] = op
	}
}

// AST creates an abstract syntax tree from a parse tree/forest.
func (ab *ASTBuilder) AST(parseTree *sppf.Forest) (*GCons, interface{}) {
	if parseTree == nil {
		return nil, nil
	}
	ab.forest = parseTree
	cursor := ab.forest.SetCursor(nil, nil) // TODO set Pruner
	value := cursor.TopDown(ab, sppf.LtoR, sppf.Break)
	T().Infof("AST creation return value = %v", value)
	if value != nil {
		ab.ast = value.(*GCons)
		T().Infof("AST = %s", ab.ast.ListString())
	}
	return ab.ast, value
}

// --- sppf.Listener interface -----------------------------------------------

// EnterRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) EnterRule(sym *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) bool {
	if op, ok := ab.operators[sym.Name]; ok {
		if !op.Descend(ctxt) {
			return false
		}
		T().Debugf("enter operator symbol: %v", sym)
		ab.stack = append(ab.stack, &GCons{makeNode(op), nil})
	} else {
		T().Debugf("enter grammar symbol: %v", sym)
	}
	return true
}

// ExitRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) ExitRule(sym *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) interface{} {
	if op, ok := ab.operators[sym.Name]; ok {
		env, err := EnvironmentForGrammarSymbol(sym.Name, ab.G)
		if err != nil && ab.Error != nil {
			ab.Error(err)
		}
		rhsList := ab.stack[len(ab.stack)-1]
		end := rhsList
		//T().Debugf("iterating over %d RHS elements", len(rhs))
		for _, r := range rhs {
			sym := env.Intern(r.Symbol().Name, true)
			//T().Debugf("sym = %v", sym)
			if !r.Symbol().IsTerminal() {
				sym.value.atom.Data = r.Value // value must be a Node
			}
			rhsList, end = growRHSList(rhsList, end, r, env)
			// switch v := r.Value.(type) { // TODO same logic as below (factor out)
			// case Node:
			// 	end = appendNode(end, v)
			// case *GCons:
			// 	end = appendTee(end, v)
			// default:
			// 	panic("Unknown value type of RHS symbol")
			// }
		}
		T().Infof("Rewrite of %s", rhsList.ListString())
		rhsList = op.Rewrite(rhsList, env)
		ab.stack = ab.stack[:len(ab.stack)-1]
		T().Debugf("%s returns %s", sym.Name, rhsList.ListString())
		return rhsList
	}
	var list, end *GCons
	for _, r := range rhs {
		list, end = growRHSList(list, end, r, globalEnvironment)
		// switch v := r.Value.(type) {
		// case Node:
		// 	end = appendNode(end, v)
		// 	if list == nil {
		// 		list = end
		// 	}
		// case *GCons:
		// 	if v.car.Type() == OperatorType {
		// 		//T().Infof("%s: tee appending %v", sym, v.ListString())
		// 		end = appendTee(end, v)
		// 		if list == nil {
		// 			list = end
		// 		}
		// 	} else {
		// 		var l *GCons
		// 		//T().Infof("%s: inline appending %v", sym, v.ListString())
		// 		l, end = appendList(end, v)
		// 		if list == nil {
		// 			list = l
		// 		}
		// 	}
		// default:
		// 	panic("Unknown value type of RHS symbol")
		// }
	}
	//T().Infof("List of length %d: %s", list.Length(), list.ListString())
	if list.Length() == 1 && list.car.Type() == ConsType {
		//T().Infof("Inner list of length %d: %s", list.car.child.Length(), list.car.child.ListString())
		list = list.car.child // unwrap
	}
	T().Infof("%s returns %s", sym.Name, list.ListString())
	T().Debugf("exit grammar symbol: %v", sym)
	return list
}

func growRHSList(start, end *GCons, r *sppf.RuleNode, env *Environment) (*GCons, *GCons) {
	switch v := r.Value.(type) {
	case Node:
		end = appendNode(end, v)
		if start == nil {
			start = end
		}
	case *GCons:
		if v.car.Type() == OperatorType {
			//T().Infof("%s: tee appending %v", sym, v.ListString())
			end = appendTee(end, v)
			if start == nil {
				start = end
			}
		} else {
			var l *GCons
			//T().Infof("%s: inline appending %v", sym, v.ListString())
			l, end = appendList(end, v)
			if start == nil {
				start = l
			}
		}
	default:
		panic("Unknown value type of RHS symbol")
	}
	return start, end
}

// Terminal is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Terminal(tokval int, token interface{}, ctxt sppf.RuleCtxt) interface{} {
	//t := ab.G.Terminal(tokval).Name
	terminal := ab.G.Terminal(tokval)
	node := makeNode(&Token{terminal.Name, tokval})
	T().Debugf("cons(terminal=%s) = %v", ab.G.Terminal(tokval).Name, node)
	return node
}

// Conflict is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Conflict(sym *lr.Symbol, ctxt sppf.RuleCtxt) (int, error) {
	panic("Conflict of AST building not yet implemented")
	//return 0, nil
}

// MakeAttrs is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) MakeAttrs(*lr.Symbol) interface{} {
	return nil // TODO
}

// ---------------------------------------------------------------------------

func appendNode(cons *GCons, node Node) *GCons {
	if cons == nil {
		return &GCons{node, nil}
	}
	cons.cdr = &GCons{node, nil}
	return cons.cdr
}

func appendList(cons *GCons, list *GCons) (*GCons, *GCons) {
	start := cons
	if cons == nil {
		cons = list
		start = list
	} else {
		cons.cdr = list
	}
	for cons.cdr != nil {
		cons = cons.cdr
	}
	T().Debugf("appendList: new list is %s", start.ListString())
	return start, cons
}

func appendTee(cons *GCons, list *GCons) *GCons {
	tee := &GCons{nullNode, nil}
	tee.car.child = list
	if cons == nil {
		cons = tee
	} else {
		cons.cdr = tee
	}
	return tee
}
