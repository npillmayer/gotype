package termr

import (
	"errors"
	"fmt"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/iteratable"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
)

// ASTBuilder is a parse tree listener for building ASTs.
// AST is a homogenous abstract syntax tree.
type ASTBuilder struct {
	G                *lr.Grammar  // input grammar the parse forest stems from
	forest           *sppf.Forest // input parse forest
	ast              *terex.GCons // root of the AST to construct
	last             *terex.GCons // current last node to append conses
	rewriters        map[string]TermR
	conflictStrategy sppf.Pruner
	Error            func(error)
	stack            []*terex.GCons
}

type ruleABEnter func(sym *lr.Symbol, rhs []*sppf.RuleNode) bool
type ruleABExit func(sym *lr.Symbol, rhs []*sppf.RuleNode) interface{}

// NewASTBuilder creates an AST builder from a parse forest/tree.
func NewASTBuilder(g *lr.Grammar) *ASTBuilder {
	ab := &ASTBuilder{
		G:         g,
		ast:       &terex.GCons{Car: terex.NilAtom, Cdr: nil}, // AST anchor
		stack:     make([]*terex.GCons, 0, 256),
		rewriters: make(map[string]TermR),
	}
	ab.last = ab.ast
	ab.stack = append(ab.stack, ab.ast) // push as stopper
	return ab
}

// TermR is a type for a rewriter for AST creation and transformation.
type TermR interface {
	Name() string                                           // printable name
	Rewrite(*terex.GCons, *terex.Environment) terex.Element // term rewriting
	Descend(sppf.RuleCtxt) bool                             // predicate wether to descend to children nodes
	Operator() terex.Operator                               // operator to place as sub-tree node
}

// AddTermR adds an AST rewriter for a grammar symbol to the builder.
func (ab *ASTBuilder) AddTermR(op TermR) {
	if op != nil {
		ab.rewriters[op.Name()] = op
	}
}

// AST creates an abstract syntax tree from a parse tree/forest.
func (ab *ASTBuilder) AST(parseTree *sppf.Forest) (*terex.GCons, interface{}) {
	if parseTree == nil {
		return nil, nil
	}
	ab.forest = parseTree
	cursor := ab.forest.SetCursor(nil, nil) // TODO set Pruner
	value := cursor.TopDown(ab, sppf.LtoR, sppf.Break)
	T().Infof("AST creation return value = %v", value)
	if value != nil {
		ab.ast = value.(terex.Element).AsList()
		T().Infof("AST = %s", ab.ast.ListString())
	}
	return ab.ast, value
}

// --- sppf.Listener interface -----------------------------------------------

// EnterRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) EnterRule(sym *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) bool {
	if rew, ok := ab.rewriters[sym.Name]; ok {
		if !rew.Descend(ctxt) {
			return false
		}
		T().Errorf("enter operator symbol: %v", sym)
		opSymListStart := terex.Cons(terex.Atomize(rew.Operator()), nil)
		ab.stack = append(ab.stack, opSymListStart) // put '(op ... ' on stack
	} else {
		T().Debugf("enter grammar symbol: %v", sym)
	}
	return true
}

// ExitRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) ExitRule(sym *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) interface{} {
	if op, ok := ab.rewriters[sym.Name]; ok {
		env, err := EnvironmentForGrammarSymbol(sym.Name, ab.G)
		if err != nil && ab.Error != nil {
			ab.Error(err)
		}
		rhsList := ab.stack[len(ab.stack)-1]
		end := rhsList
		T().Infof("iterating over %d RHS elements", len(rhs))
		for _, r := range rhs {
			T().Infof("r = %v", r)
			// TODO set value of RHS vars in Env
			rhssym := env.Intern(r.Symbol().Name, true)
			T().Infof("sym = %v", sym)
			if !r.Symbol().IsTerminal() {
				rhssym.Value.Data = r.Value
			}
			rhsList, end = growRHSList(rhsList, end, r, env)
		}
		T().Infof("%s: Rewrite of %s", sym.Name, rhsList.ListString())
		rewritten := op.Rewrite(rhsList, env) // returns an terex.Element
		ab.stack = ab.stack[:len(ab.stack)-1] // pop initial '(op ...'
		T().Infof("%s returns %s", sym.Name, rewritten.String())
		return rewritten
	}
	//var list, end *terex.GCons
	var list, end *terex.GCons
	for _, r := range rhs {
		list, end = growRHSList(list, end, r, terex.GlobalEnvironment)
	}
	rew := noOpRewrite(list)
	T().Infof("%s returns %s", sym.Name, rew.String())
	T().Infof("exit grammar symbol: %v", sym)
	return rew
}

func noOpRewrite(list *terex.GCons) terex.Element {
	if list != nil && list.Length() == 1 {
		return terex.Elem(list.Car)
	}
	return terex.Elem(list)
}

//func growRHSList(start, end *terex.GCons, r *sppf.RuleNode, env *terex.Environment) (*terex.GCons, *terex.GCons) {
func growRHSList(start, end *terex.GCons, r *sppf.RuleNode, env *terex.Environment) (*terex.GCons, *terex.GCons) {
	if _, ok := r.Value.(terex.Element); !ok {
		T().Errorf("r.Value=%v", r.Value)
		panic("RHS symbol is not of type Element")
	}
	e := r.Value.(terex.Element) // value of rule-node r is either atom or list
	if e.IsNil() {
		return start, end
	}
	if e.IsAtom() {
		end = (appendAtom(end, e.AsAtom()))
		if start == nil {
			start = end
		}
	} else {
		l := e.AsList()
		if l.Car.Type() == terex.OperatorType {
			//T().Infof("%s: tee appending %v", sym, v.ListString())
			end = appendTee(end, l)
			if start == nil {
				start = end
			}
		} else { // append l at end of current list
			var concat *terex.GCons
			//T().Infof("%s: inline appending %v", sym, v.ListString())
			concat, end = appendList(end, l)
			if start == nil {
				start = concat
			}
		}
	}
	return start, end
}

// Terminal is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Terminal(tokval int, token interface{}, ctxt sppf.RuleCtxt) interface{} {
	//t := ab.G.Terminal(tokval).Name
	terminal := ab.G.Terminal(tokval)
	atom := terex.Atomize(&terex.Token{Name: terminal.Name, Value: tokval})
	T().Debugf("cons(terminal=%s) = %v", ab.G.Terminal(tokval).Name, atom)
	return terex.Elem(atom)
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

func appendAtom(cons *terex.GCons, atom terex.Atom) *terex.GCons {
	if atom == terex.NilAtom {
		return cons
	}
	if cons == nil {
		return terex.Cons(atom, nil)
	}
	cons.Cdr = terex.Cons(atom, nil)
	return cons.Cdr
}

func appendList(cons *terex.GCons, list *terex.GCons) (*terex.GCons, *terex.GCons) {
	start := cons
	if cons == nil {
		cons = list
		start = list
	} else {
		cons.Cdr = list
	}
	for cons.Cdr != nil {
		cons = cons.Cdr
	}
	T().Debugf("appendList: new list is %s", start.ListString())
	return start, cons
}

func appendTee(cons *terex.GCons, list *terex.GCons) *terex.GCons {
	tee := terex.Cons(terex.Atomize(list), nil)
	if cons == nil {
		cons = tee
	} else {
		cons.Cdr = tee
	}
	return tee
}

// EnvironmentForGrammarSymbol creates a new environment, suitable for the
// grammar symbols at a given tree node of a parse-tree or AST.
//
// Given a grammar production
//
//     A -> B C D
//
// it will create an environment #A for A, with pre-interned (but empty) symbols
// for A, B, C and D. If any of the right-hand-side symbols are terminals, they will
// be created as nodes with an appropriate atom type.
//
func EnvironmentForGrammarSymbol(symname string, G *lr.Grammar) (*terex.Environment, error) {
	if G == nil {
		return terex.GlobalEnvironment, errors.New("Grammar is null")
	}
	envname := "#" + symname
	if env := terex.GlobalEnvironment.FindSymbol(envname, false); env != nil {
		if env.Value.Type() != terex.EnvironmentType {
			panic(fmt.Errorf("Internal error, environment misconstructed: %s", envname))
		}
		return env.Value.Data.(*terex.Environment), nil
	}
	gsym := G.SymbolByName(symname)
	if gsym == nil || gsym.IsTerminal() {
		return terex.GlobalEnvironment, fmt.Errorf("Non-terminal not found in grammar: %s", symname)
	}
	env := terex.NewEnvironment(envname, nil)
	rhsSyms := iteratable.NewSet(0)
	rules := G.FindNonTermRules(gsym, false)
	rules.IterateOnce()
	for rules.Next() {
		rule := rules.Item().(lr.Item).Rule()
		for _, s := range rule.RHS() {
			rhsSyms.Add(s)
		}
	}
	rhsSyms.IterateOnce()
	for rhsSyms.Next() {
		gsym := rhsSyms.Item().(*lr.Symbol)
		sym := env.Intern(gsym.Name, false)
		if gsym.IsTerminal() {
			sym.Value = terex.Atomize(gsym.Value)
		}
	}
	// e := globalEnvironment.Intern(envname, false)
	// e.atom.typ = EnvironmentType
	// e.atom.Data = env
	return env, nil
}
