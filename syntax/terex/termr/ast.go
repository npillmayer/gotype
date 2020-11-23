package termr

/*
BSD License

Copyright (c) 2019–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/timtadh/lexmachine"
)

// ASTBuilder is a parse tree listener for building ASTs. It will construct an
// abstract syntax tree (AST) while walking the parse tree. As we may support
// ambiguous grammars, the parse tree may happen to be a parse forest, containing
// more than one derivation for a sentence.
type ASTBuilder struct {
	G                *lr.Grammar        // input grammar the parse forest stems from
	Env              *terex.Environment // environment for symbols of the AST to create
	forest           *sppf.Forest       // input parse tree/forest
	conflictStrategy sppf.Pruner        // how to deal with parse-forest ambiguities
	toks             TokenRetriever     // retriever for parse tree leafs
	rewriters        map[string]TermR   // term rewriters to apply
	Error            func(error)        // user supplied handler for errors
	stack            []*terex.GCons     // used for recursive operator walking
}

// ErrorHandler is an interface to process errors occuring during parsing.
type ErrorHandler interface {
	Error(error) bool
}

// NewASTBuilder creates an AST builder from a parse forest/tree.
// It is initialized with the grammar, which has been used to create the parse tree.
//
// Clients will first create an ASTBuilder, then initialize it with all the
// term rewriters and variables/symbols necessary, and finally call ASTBuilder.AST(…).
func NewASTBuilder(g *lr.Grammar) *ASTBuilder {
	ab := &ASTBuilder{
		G:         g,
		Env:       terex.NewEnvironment("AST "+g.Name, terex.GlobalEnvironment),
		stack:     make([]*terex.GCons, 0, 256),
		rewriters: make(map[string]TermR),
	}
	return ab
}

// TermR is a type for a rewriter for AST creation and transformation.
type TermR interface {
	Name() string                                           // printable name
	Rewrite(*terex.GCons, *terex.Environment) terex.Element // term rewriting
	Descend(sppf.RuleCtxt) bool                             // predicate wether to descend to children nodes
	Operator() terex.Operator                               // operator to place as sub-tree node
}

// AddTermR adds an AST rewriter for a non-terminal grammar symbol to the builder.
func (ab *ASTBuilder) AddTermR(op TermR) {
	if op != nil {
		ab.rewriters[op.Name()] = op
	}
}

// AST creates an abstract syntax tree from a parse tree/forest.
// The type of ASTs we create is a homogenous abstract syntax tree.
func (ab *ASTBuilder) AST(parseTree *sppf.Forest, tokRetr TokenRetriever) *terex.Environment {
	if parseTree == nil || tokRetr == nil {
		return nil
	}
	ab.forest = parseTree
	ab.toks = tokRetr
	cursor := ab.forest.SetCursor(nil, nil) // TODO set Pruner
	value := cursor.TopDown(ab, sppf.LtoR, sppf.Break)
	T().Infof("AST creation return value = %v", value)
	if value != nil {
		ab.Env.AST = value.(terex.Element).AsList()
		T().Infof("AST = %s", ab.Env.AST.ListString())
	}
	return ab.Env
}

// TokenRetriever is a type for getting tokens at an input position.
type TokenRetriever func(uint64) interface{}

// --- sppf.Listener interface -----------------------------------------------

// EnterRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) EnterRule(sym *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) bool {
	if rew, ok := ab.rewriters[sym.Name]; ok {
		if !rew.Descend(ctxt) {
			return false
		}
		T().Debugf("-------> enter operator symbol: %v", sym)
		opSymListStart := terex.Cons(terex.Atomize(rew.Operator()), nil)
		ab.stack = append(ab.stack, opSymListStart) // put '(op ... ' on stack
	} else {
		T().Debugf("-------> enter grammar symbol: %v", sym)
	}
	return true
}

// ExitRule is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) ExitRule(sym *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) interface{} {
	T().Debugf("<------- exit symbol: %v, now rewriting", sym)
	if op, ok := ab.rewriters[sym.Name]; ok {
		//env, err := ab.EnvironmentForGrammarSymbol(sym.Name)
		env, err := ab.EnvironmentForGrammarRule(sym.Name, rhs)
		if err != nil && ab.Error != nil {
			ab.Error(err)
		}
		rhsList := ab.stack[len(ab.stack)-1] // operator is TOS ⇒ first element of RHS list
		//end := rhsList
		for _, r := range rhs { // collect all the value of RHS symbols
			// T().Infof("r = %v", r)
			// // TODO set value of RHS vars in Env
			// rhssym := env.Intern(r.Symbol().Name, true)
			// T().Infof("sym = %v", sym)
			// if !r.Symbol().IsTerminal() {
			// 	rhssym.Value.Data = r.Value
			// }
			//rhsList, end = growRHSList(rhsList, end, r, env)
			rhsList = appendRHSResult(rhsList, r)
		}
		T().Infof("%s: Rewrite of %s", sym.Name, rhsList.ListString())
		rewritten := op.Rewrite(rhsList, env) // returns a terex.Element
		ab.stack = ab.stack[:len(ab.stack)-1] // pop initial '(op ...'
		T().Infof("%s returns %s", sym.Name, rewritten.String())
		return rewritten
	}
	//var list, end *terex.GCons
	var list *terex.GCons
	T().Infof("%s will rewrite |rhs| = %d symbols", sym.Name, len(rhs))
	for _, r := range rhs {
		//list, end = growRHSList(list, end, r, terex.GlobalEnvironment)
		list = appendRHSResult(list, r)
	}
	rew := noOpRewrite(list)
	T().Infof("%s returns %s", sym.Name, rew.String())
	T().Infof("done with rewriting grammar symbol: %v", sym)
	return rew
}

func noOpRewrite(list *terex.GCons) terex.Element {
	T().Debugf("no-op rewrite of %v", list)
	if list != nil && list.Length() == 1 {
		return terex.Elem(list.Car)
	}
	return terex.Elem(list)
}

func appendRHSResult(list *terex.GCons, r *sppf.RuleNode) *terex.GCons {
	if _, ok := r.Value.(terex.Element); !ok {
		T().Errorf("r.Value=%v", r.Value)
		panic("RHS symbol is not of type Element")
	}
	e := r.Value.(terex.Element) // value of rule-node r is either atom or list
	// T().Debugf("RHS-appending value = %s", e.String())
	if e.IsNil() {
		return list
	}
	if e.IsAtom() {
		list = list.Append(terex.Cons(e.AsAtom(), nil))
		return list
	}
	l := e.AsList()
	if l.Car.Type() == terex.OperatorType {
		list = list.Branch(l)
	} else {
		list = list.Append(l)
	}
	return list
}

//func growRHSList(start, end *terex.GCons, r *sppf.RuleNode, env *terex.Environment) (*terex.GCons, *terex.GCons) {
/* func growRHSList(start, end *terex.GCons, r *sppf.RuleNode, env *terex.Environment) (*terex.GCons, *terex.GCons) {
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
} */

// Terminal is part of sppf.Listener interface.
// Not intended for direct client use.
func (ab *ASTBuilder) Terminal(tokval int, token interface{}, ctxt sppf.RuleCtxt) interface{} {
	// T().Errorf("TERMINAL TOKEN FOR %d ------------------", tokval)
	// T().Errorf("SPAN = %d", ctxt.Span.Len())
	if ctxt.Span.Len() == 0 { // RHS is epsilon
		return terex.Elem(nil)
	}
	terminal := ab.G.Terminal(tokval)
	tokpos := ctxt.Span.From()
	t := ab.toks(tokpos) // opaque token type
	atom := terex.Atomize(&terex.Token{Name: terminal.Name, TokType: tokval, Token: t})
	if t != nil {
		s := ""
		if tt, ok := t.(*lexmachine.Token); ok {
			s = string(tt.Lexeme)
		} else if tt, ok := t.(fmt.Stringer); ok {
			s = tt.String()
		} else {
			s = t.(string)
		}
		T().Debugf("CONS(terminal=%s) = %v @%d [%s]", terminal.Name, atom, tokpos, s)
	} else {
		T().Debugf("CONS(terminal=%s) = %v @%d", terminal.Name, atom, tokpos)
	}
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

// EnvironmentForGrammarSymbol creates an empty new environment, suitable for the
// grammar symbols at a given tree node of a parse-tree or AST.
func (ab *ASTBuilder) environmentForGrammarSymbol(symname string) (*terex.Environment, error) {
	if ab.G == nil {
		return terex.GlobalEnvironment, errors.New("Grammar is null")
	}
	envname := "#" + symname
	//if env := terex.GlobalEnvironment.FindSymbol(envname, false); env != nil {
	/* 	if env := ab.Env.FindSymbol(envname, false); env != nil {
		if env.Value.Type() != terex.EnvironmentType {
			panic(fmt.Errorf("Internal error, environment misconstructed: %s", envname))
		}
		return env.Value.Data.(*terex.Environment), nil
	} */
	env := terex.NewEnvironment(envname, ab.Env)
	gsym := ab.G.SymbolByName(symname)
	if gsym == nil || gsym.IsTerminal() {
		//return terex.GlobalEnvironment, fmt.Errorf("Non-terminal not found in grammar: %s", symname)
		return env, fmt.Errorf("Non-terminal not found in grammar: %s", symname)
	}
	/* 	rhsSyms := iteratable.NewSet(0)
	   	rules := ab.G.FindNonTermRules(gsym, false)
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
	   	} */
	// e := globalEnvironment.Intern(envname, false)
	// e.atom.typ = EnvironmentType
	// e.atom.Data = env
	return env, nil
}

// EnvironmentForGrammarRule creates a new environment, suitable for the
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
// For recurring non-terminal RHS symbols, as in
//
//     A -> B + B
//
// sequence number suffixes will be created. The grammar rule above will produce an
// environment containing
//
//     B, B.1   for the 1st B
//     B.2      for the 2nd B
//
func (ab *ASTBuilder) EnvironmentForGrammarRule(symname string, rhs []*sppf.RuleNode) (*terex.Environment, error) {
	env, err := ab.environmentForGrammarSymbol(symname)
	if err != nil {
		return env, err
	}
	symtrack := make(map[string]int)
	for _, r := range rhs { // set value of RHS vars in Env
		T().Debugf("RHS: r = %v", r)
		symname := r.Symbol().Name
		count := symtrack[symname]
		envsym := env.Intern(symname+"."+strconv.Itoa(count), true)
		setSymbolValue(envsym, r)
		if count == 0 {
			//envsym := env.Intern(r.Symbol().Name, true)
			envsym = env.Intern(symname, true)
			setSymbolValue(envsym, r)
		}
		symtrack[symname] = count + 1
		T().Debugf("env sym = %v", envsym)
	}
	return env, nil
}

// --- Helpers ---------------------------------------------------------------

func setSymbolValue(envsym *terex.Symbol, r *sppf.RuleNode) {
	if r.Symbol().IsTerminal() {
		envsym.Value.Data = r.Value
	} else {
		envsym.Value = terex.Atomize(r.Value)
	}
}

/*
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
	//T().Debugf("appendList: new list is %s", start.ListString())
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
*/
