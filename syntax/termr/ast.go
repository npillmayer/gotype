package termr

import (
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
)

// ASTBuilder is a parse tree listener for building ASTs.
// AST is a homogenous abstract syntax tree.
type ASTBuilder struct {
	G      *lr.Grammar  // input grammar the parse forest stems from
	forest *sppf.Forest // input parse forest
	ast    *GCons       // root of the AST to construct
	last   *GCons       // current last node to append conses
	// dispatchEnter    func(string) ruleABEnter
	// dispatchExit     func(string) ruleABExit
	operators        map[string]*ASTOperator
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
		ast:       &GCons{carNode{NullAtom, nil}, nil}, // AST anchor
		stack:     make([]*GCons, 0, 256),
		operators: make(map[string]*ASTOperator),
		// dispatchEnter: nullABEnter,
		// dispatchExit:  nullABExit,
	}
	ab.last = ab.ast
	ab.stack = append(ab.stack, ab.ast) // push as stopper
	return ab
}

// func nullABEnter(string) ruleABEnter {
// 	return nil
// }

// func nullABExit(string) ruleABExit {
// 	return nil
// }

type ASTOperator struct {
	Name    string
	Rewrite func(*GCons, *Environment) *GCons
	Descend func(sppf.RuleCtxt) bool
}

func (op *ASTOperator) String() string {
	return op.Name + ":Op"
}

func (op *ASTOperator) Call(term *GCons) *GCons {
	return op.Rewrite(term, globalEnvironment)
}

func (ab *ASTBuilder) AddOperator(op *ASTOperator) {
	if op != nil {
		ab.operators[op.Name] = op
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

// func (ab *ASTBuilder) appendList(list *GCons) {
// 	if list == nil {
// 		return
// 	}
// 	T().Debugf("APPEND LIST %s", list.ListString())
// 	//ab.last.cdr = &GCons{carNode{NullAtom, list}, nil}
// 	ab.last.cdr = &GCons{list.Car().car, nil}
// 	ab.last = ab.last.cdr
// }

// func (ab *ASTBuilder) append(car carNode) {
// 	if car == nullCar {
// 		return
// 	}
// 	ab.last.cdr = &GCons{car, nil}
// 	ab.last = ab.last.cdr
// }

// func (ab *ASTBuilder) up() {
// 	if ab.stack[len(ab.stack)-1].car.atom == NullAtom {
// 		T().Errorf("up() called with empty stack")
// 		panic("empty stack")
// 	}
// 	ab.last = ab.stack[len(ab.stack)-1]
// 	ab.stack = ab.stack[:len(ab.stack)-1] // pop last
// }

// func (ab *ASTBuilder) down(car carNode) {
// 	ab.stack = append(ab.stack, ab.last)
// 	child := &GCons{car, nil}
// 	ab.last.cdr = &GCons{carNode{NullAtom, child}, nil}
// 	ab.last = child
// }

// --- sppf.Listener interface -----------------------------------------------

// EnterRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) EnterRule(sym *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) bool {
	// if r := ab.dispatchEnter(sym.Name); r != nil {
	// 	return r(sym, rhs)
	// }
	if op, ok := ab.operators[sym.Name]; ok {
		if !op.Descend(ctxt) {
			return false
		}
		T().Debugf("enter operator symbol: %v", sym)
		ab.stack = append(ab.stack, &GCons{makeCar(op), nil})
	} else {
		T().Debugf("enter grammar symbol: %v", sym)
	}
	return true
}

// ExitRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) ExitRule(sym *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) interface{} {
	// if r := ab.dispatchExit(sym.Name); r != nil {
	// 	return r(sym, rhs)
	// }
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
				sym.car.atom.Data = r.Value // must be a carNode
			}
			switch v := r.Value.(type) {
			case carNode:
				end = appendCar(end, v)
			case *GCons:
				end = appendTee(end, v)
			default:
				panic("Unknown value type of RHS symbol")
			}
		}
		T().Debugf("Rewrite of %s", rhsList.ListString())
		rhsList = op.Rewrite(rhsList, env)
		//ab.stack[len(ab.stack)-1] = rhsList
		ab.stack = ab.stack[:len(ab.stack)-1]
		T().Debugf("%s returns %s", sym.Name, rhsList.ListString())
		return rhsList
	}
	var list *GCons
	var end *GCons
	for _, r := range rhs {
		switch v := r.Value.(type) {
		case carNode:
			end = appendCar(end, v)
		case *GCons:
			end = appendTee(end, v)
		default:
			panic("Unknown value type of RHS symbol")
		}
		if list == nil {
			list = end
		}
	}
	T().Debugf("%s returns %s", sym.Name, list.ListString())
	T().Debugf("exit grammar symbol: %v", sym)
	return list
}

// Terminal is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Terminal(tokval int, token interface{}, ctxt sppf.RuleCtxt) interface{} {
	//t := ab.G.Terminal(tokval).Name
	car := makeCar(tokval)
	car.atom.typ = TokenType
	T().Debugf("cons(terminal=%s) = %v", ab.G.Terminal(tokval).Name, car)
	return car
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

func appendCar(cons *GCons, car carNode) *GCons {
	if cons == nil {
		return &GCons{car, nil}
	}
	cons.cdr = &GCons{car, nil}
	return cons.cdr
}

func appendList(cons *GCons, list *GCons) *GCons {
	if cons == nil {
		cons = list
	} else {
		cons.cdr = list
	}
	for cons.cdr != nil {
		cons = cons.cdr
	}
	return cons
}

func appendTee(cons *GCons, list *GCons) *GCons {
	tee := &GCons{nullCar, nil}
	tee.car.child = list
	if cons == nil {
		cons = tee
	} else {
		cons.cdr = tee
	}
	return tee
}
