package terexlang

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
	"fmt"
	"strconv"
	"sync"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/npillmayer/gotype/syntax/terex/termr"
	"github.com/timtadh/lexmachine"
)

// --- Grammar ---------------------------------------------------------------

// Atom       ::=  '\'' Atom   // currently un-ambiguated by QuoteOrAtom
// Atom       ::=  ident             // a
// Atom       ::=  string            // "abc"
// Atom       ::=  number            // 123.45
// Atom       ::=  variable          // $a
// Atom       ::=  List
// List       ::=  '(' Sequence ')'
// Sequence   ::=  Sequence Atom
// Sequence   ::=  Atom
//
// Comments starting with ';' will be filtered by the scanner.
//
func makeTeRExGrammar() (*lr.LRAnalysis, error) {
	b := lr.NewGrammarBuilder("TeREx S-Expr")
	b.LHS("QuoteOrAtom").N("Quote").End()
	b.LHS("QuoteOrAtom").N("Atom").End()
	b.LHS("Quote").T(Token("'")).N("Atom").End()
	b.LHS("Atom").T(Token("ID")).End()
	b.LHS("Atom").T(Token("STRING")).End()
	b.LHS("Atom").T(Token("NUM")).End()
	b.LHS("Atom").T(Token("VAR")).End()
	b.LHS("Atom").N("List").End()
	b.LHS("List").T(Token("(")).N("Sequence").T(Token(")")).End()
	b.LHS("Sequence").N("QuoteOrAtom").N("Sequence").End()
	b.LHS("Sequence").Epsilon()
	g, err := b.Grammar()
	if err != nil {
		return nil, err
	}
	return lr.Analysis(g), nil
}

var grammar *lr.LRAnalysis
var lexer *scanner.LMAdapter

var startOnce sync.Once // monitors one-time creation of grammar and lexer

func createParser() *earley.Parser {
	startOnce.Do(func() {
		var err error
		T().Infof("Creating lexer")
		if lexer, err = Lexer(); err != nil { // MUST be called before grammar builing !
			panic("Cannot create lexer")
		}
		T().Infof("Creating grammar")
		if grammar, err = makeTeRExGrammar(); err != nil {
			panic("Cannot create global grammar")
		}
		initRewriters()
	})
	return earley.NewParser(grammar, earley.GenerateTree(true))
}

// NewASTBuilder returns a new AST builder for the TeREx language
func newASTBuilder() *termr.ASTBuilder {
	ab := termr.NewASTBuilder(grammar.Grammar())
	//ab.AddTermR(opOp)
	ab.AddTermR(quoteOp)
	ab.AddTermR(seqOp)
	ab.AddTermR(listOp)
	atomOp := makeASTTermR("Atom", "atom")
	atomOp.rewrite = func(l *terex.GCons, env *terex.Environment) terex.Element {
		// rewrite: (:atom x)  => x
		e := terex.Elem(l.Cdar())
		// in (atom x), check if x is terminal
		T().Infof("atomOp.rewrite: l.Cdr = %v", e)
		if l.Cdr.IsLeaf() {
			e = convertTerminalToken(e, ab.Env)
		}
		T().Infof("atomOp.rewrite => %s", e.String())
		return e
	}
	ab.AddTermR(atomOp)
	return ab
}

// Parse parses an input string, given in TeREx language format. It returns the
// parse forest and a TokenReceiver, or an error in case of failure.
//
// Clients may use a terex.ASTBuilder to create an abstract syntax tree
// from the parse forest.
//
func Parse(input string) (*sppf.Forest, termr.TokenRetriever, error) {
	parser := createParser()
	scan, err := lexer.Scanner(input)
	if err != nil {
		return nil, nil, err
	}
	accept, err := parser.Parse(scan, nil)
	if err != nil {
		return nil, nil, err
	} else if !accept {
		return nil, nil, fmt.Errorf("Not a valid TeREx expression")
	}
	return parser.ParseForest(), earleyTokenReceiver(parser), nil
}

func earleyTokenReceiver(parser *earley.Parser) termr.TokenRetriever {
	return func(pos uint64) interface{} {
		return parser.TokenAt(pos)
	}
}

// AST creates an abstract syntax tree from a parse tree/forest.
//
// Returns a homogenous AST, a TeREx-environment and an error status.
func AST(parsetree *sppf.Forest, tokRetr termr.TokenRetriever) (*terex.GCons,
	*terex.Environment, error) {
	ab := newASTBuilder()
	env := ab.AST(parsetree, tokRetr)
	if env == nil {
		T().Errorf("Cannot create AST from parsetree")
		return nil, nil, fmt.Errorf("Error while creating AST")
	}
	ast := env.AST
	T().Infof("AST: %s", env.AST.ListString())
	return ast, env, nil
}

// ---------------------------------------------------------------------------

// Eval evaluates an s-expr (given in textual form).
// It will parse the string, create an internal S-expr structure and evaluate it,
// using the symbols in env,
/* func Eval(sexpr string) (*terex.GCons, *terex.Environment) {
	quoted, env, _ := Quote(sexpr)
	if quoted == nil {
		return nil, env
	}
	evald := env.Eval(quoted.Tee())
	T().Infof("eval(AST) = %s", evald.ListString())
	return evald, env
} */

// QuoteAST returns an AST, which should be the result of parsing an s-expr, as
// pure data.
//
// If the environment contains any symbol's value, quoting will replace the symbol
// by its value. For example, if the s-expr contains a symbol 'str' with a value
// of "this is a string", the resulting data structure will contain the string,
// not the name of the symbol. If you do not have use for this kind of substitution,
// simply call Quote(…) for the global environment.
//
func QuoteAST(ast terex.Element, env *terex.Environment) (terex.Element, error) {
	// ast *terex.GCons
	if env == nil {
		env = terex.GlobalEnvironment
	}
	quEnv := terex.NewEnvironment("quoting", env)
	quEnv.Defn("list", listOp.call)
	//quEnv.Defn("quote", quoteOp.call)
	quEnv.Resolver = symbolPreservingResolver{}
	q := terex.Eval(ast, quEnv)
	return q, quEnv.LastError()
}

// Quote quotes an s-expr (given in textual form), thus returning it as data.
// It will parse the string, create an internal S-expr structure (AST) and quote it,
// using the symbols in env,
/* func Quote(sexpr string) (*terex.GCons, *terex.Environment, error) {
	parsetree, tokRetr, err := Parse(sexpr)
	if err != nil {
		T().Errorf("Parsing error: %v", err)
		return nil, nil, err
	}
	if parsetree == nil {
		T().Errorf("Empty parse tree, no AST to create")
		return nil, nil, fmt.Errorf("Empty parse tree, no AST to create")
	}
	astBuilder := newASTBuilder()
	env := astBuilder.AST(parsetree, tokRetr)
	if env == nil {
		T().Errorf("Error creating AST from parsetree")
		return nil, nil, env.LastError()
	}
	T().Infof("AST: %s", env.AST.ListString())
	//r := env.Quote(env.AST)
	//T().Infof("quote(AST) = %s", r.ListString())
	//return r, env, nil
	return nil, env, nil
} */

// --- S-expr AST-builder listener -------------------------------------------

var quoteOp *sExprTermR // for Quote -> ... productions
var seqOp *sExprTermR   // for Sequence -> ... productions
var listOp *sExprTermR  // for List -> ... productions

type sExprTermR struct {
	name    string
	opname  string
	rewrite func(*terex.GCons, *terex.Environment) terex.Element
	call    func(terex.Element, *terex.Environment) terex.Element
	//isquote bool
	// call   func(terex.Element, *terex.Environment) terex.Element
}

var _ terex.Operator = &sExprTermR{}
var _ termr.TermR = &sExprTermR{}

//func makeASTTermR(name string, opname string, quoter bool) *sExprTermR {
func makeASTTermR(name string, opname string) *sExprTermR {
	termr := &sExprTermR{
		name:   name,
		opname: opname,
		//isquote: quoter,
		//rules:  make([]termr.RewriteRule, 0, 1),
	}
	return termr
}

func (trew *sExprTermR) String() string {
	return trew.name
}

func (trew *sExprTermR) Operator() terex.Operator {
	return trew
	//return listOp
	//return &globalOpInEnv{op.opname}
}

func (trew *sExprTermR) Rewrite(l *terex.GCons, env *terex.Environment) terex.Element {
	T().Debugf("%s:trew.Rewrite[%s] called", trew.String(), l.ListString())
	e := trew.rewrite(l, env)
	// T().Debugf("%s:Op.Rewrite[%s] called, %d rules", op.Name(), l.ListString(), len(op.rules))
	// for _, rule := range op.rules {
	// 	T().Infof("match: trying %s %% %s ?", rule.Pattern.ListString(), l.ListString())
	// 	if rule.Pattern.Match(l, env) {
	// 		T().Infof("Op %s has a match", op.Name())
	// 		//T().Debugf("-> pre rewrite: %s", l.ListString())
	// 		v := rule.Rewrite(l, env)
	// 		//T().Debugf("<- post rewrite:")
	// 		terex.DumpElement(v)
	// 		T().Infof("Op %s rewrite -> %s", op.Name(), v.String())
	// 		//return rule.Rewrite(l, env)
	// 		return v
	// 	}
	// }
	return e
	//return terex.Elem(e)
}

func (trew *sExprTermR) Descend(sppf.RuleCtxt) bool {
	return true
}

func (trew *sExprTermR) Call(e terex.Element, env *terex.Environment) terex.Element {
	opsym := env.FindSymbol(trew.opname, true)
	if opsym == nil {
		T().Errorf("Cannot find parsing operation %s", trew.opname)
		return e
	}
	operator, ok := opsym.Value.Data.(terex.Operator)
	if !ok {
		T().Errorf("Cannot call parsing operation %s", trew.opname)
		return e
	}
	return operator.Call(e, env)
}

// func (trew *sExprTermR) IsQuoter() bool {
// 	return trew.isquote
// }

// func (op *sExprTermR) Rule(pattern *terex.GCons, rw termr.Rewriter) *sExprTermR {
// 	r := termr.RewriteRule{
// 		Pattern: pattern,
// 		Rewrite: rw,
// 	}
// 	op.rules = append(op.rules, r)
// 	return op
// }

// SingleTokenArg is a pattern matching an operator with a single arg of TokenType.
// var SingleTokenArg *terex.GCons

func initRewriters() {
	// SingleTokenArg = terex.Cons(terex.Atomize(terex.OperatorType), termr.AnySymbol())
	// opOp = makeASTTermR("Op", "").Rule(termr.Anything(), func(l *terex.GCons, env *terex.Environment) terex.Element {
	// 	if l.Length() <= 1 || l.Cdar().Type() != terex.TokenType {
	// 		return terex.Elem(l)
	// 	}
	// 	// (op "x") => x:op
	// 	tname := l.Cdar().Data.(*terex.Token).String()
	// 	if sym := terex.GlobalEnvironment.FindSymbol(tname, true); sym != nil {
	// 		if sym.Value.Type() == terex.OperatorType {
	// 			op := terex.Atomize(&globalOpInEnv{tname})
	// 			return terex.Elem(op)
	// 		}
	// 	}
	// 	return terex.Elem(l)
	// })
	// _, tokval := Token("'")
	// p := terex.Cons(terex.Atomize(terex.OperatorType),
	// 	terex.Cons(terex.Atomize(&terex.Token{Name: "'", TokType: tokval}), termr.AnySymbol()))
	quoteOp = makeASTTermR("Quote", "list")
	quoteOp.rewrite = func(l *terex.GCons, env *terex.Environment) terex.Element {
		// (:quote ' ⟨atom⟩) =>  (:list quote ⟨atom⟩)
		q := env.Intern("quote", false)
		qu := terex.Cons(terex.Atomize(q), l.Cddr())
		return terex.Elem(terex.Cons(l.Car, qu))
	}
	// quoteOp.call = func(e terex.Element, env *terex.Environment) terex.Element {
	// 	T().Debugf("Un-QUOTE of %v", e)
	// 	// :quote(atom) =>  atom
	// 	return e
	// }
	seqOp = makeASTTermR("Sequence", "seq")
	seqOp.rewrite = func(l *terex.GCons, env *terex.Environment) terex.Element {
		switch l.Length() {
		case 0:
			return terex.Elem(nil)
		case 1:
			return terex.Elem(nil)
		case 2:
			if l.Cdar().Type() == terex.ConsType {
				return terex.Elem(l.Cdr.Tee())
			}
			return terex.Elem(l.Cdar())
		}
		return terex.Elem(l.Cdr)
	}
	listOp = makeASTTermR("List", "list")
	listOp.rewrite = func(l *terex.GCons, env *terex.Environment) terex.Element {
		// list '(' x y ... ')'  => (:list x y ...)
		// if l.Length() <= 3 { // ( )
		// 	return terex.Elem(l.Car)
		// }
		content := l.Cddr()                            // strip '('
		content = content.FirstN(content.Length() - 1) // strip ')'
		T().Debugf("List content = %v", content)
		return terex.Elem(terex.Cons(l.Car, content)) // (List:Op ...)
	}
	listOp.call = func(e terex.Element, env *terex.Environment) terex.Element {
		// (:list a b c) =>  (a b c)
		T().Debugf("Un-LIST of %v", e)
		list := e.AsList()
		if list.Length() == 0 { //  () => nil  [empty list is nil]
			return terex.Elem(nil)
		}
		//list := args.Map(terex.Eval, env) // eval arguments
		e = terex.Eval(terex.Elem(list), env)
		return e
	}
}

// ---------------------------------------------------------------------------

/*
type globalOpInEnv struct {
	opname string
}

func (op globalOpInEnv) String() string {
	return op.opname
}

// Call is part of interface Operator.
func (op globalOpInEnv) Call(term terex.Element, env *terex.Environment) terex.Element {
	opsym := terex.GlobalEnvironment.FindSymbol(op.opname, true)
	if opsym == nil {
		T().Errorf("Cannot find parsing operation %s", op.opname)
		return term
	}
	operator, ok := opsym.Value.Data.(terex.Operator)
	if !ok {
		T().Errorf("Cannot call parsing operation %s", op.opname)
		return term
	}
	return operator.Call(term, env)
}

// Quote is part of interface Operator
func (op globalOpInEnv) Quote(term terex.Element, env *terex.Environment) terex.Element {
	opsym := terex.GlobalEnvironment.FindSymbol(op.opname, true)
	if opsym == nil {
		T().Errorf("Cannot find parsing operation %s", op.opname)
		return term
	}
	operator, ok := opsym.Value.Data.(terex.Operator)
	if !ok {
		T().Errorf("Cannot quote-call parsing operation %s", op.opname)
		return term
	}
	return operator.Call(term, env)
}
*/

func convertTerminalToken(el terex.Element, env *terex.Environment) terex.Element {
	if !el.IsAtom() {
		return el
	}
	atom := el.AsAtom()
	if atom.Type() != terex.TokenType {
		return el
	}
	t := atom.Data.(*terex.Token)
	token := t.Token.(*lexmachine.Token)
	T().Infof("Convert terminal token: '%v'", string(token.Lexeme))
	switch token.Type {
	case tokenIds["NUM"]:
		if f, err := strconv.ParseFloat(string(token.Lexeme), 64); err == nil {
			T().Debugf("   t.Value=%g", f)
			t.Value = f
		} else {
			T().Errorf("   %s", err.Error())
			return terex.Elem(terex.Atomize(err))
		}
	case tokenIds["STRING"]:
		if (len(token.Lexeme)) <= 2 {
			t.Value = ""
		} else { // trim off "…"
			//runes := []rune(string(token.Lexeme))  // unnecessary
			t.Value = string(token.Lexeme[1 : len(token.Lexeme)-1])
		}
	case tokenIds["VAR"]:
		panic("VAR type tokens not yet implemented")
		fallthrough
	case tokenIds["ID"]:
		s := string(token.Lexeme)
		//sym := terex.GlobalEnvironment.Intern(s, true)
		sym := env.Intern(s, true)
		if sym != nil {
			return terex.Elem(terex.Atomize(sym))
		}
	default:
	}
	return el
}

type symbolPreservingResolver struct{}

func (r symbolPreservingResolver) Resolve(atom terex.Atom, env *terex.Environment, asOp bool) (
	terex.Atom, error) {
	if atom.Type() == terex.TokenType {
		t := atom.Data.(*terex.Token)
		token := t.Token.(*lexmachine.Token)
		T().Debugf("Resolve terminal token: '%v'", string(token.Lexeme))
		switch token.Type {
		case tokenIds["NUM"]:
			return terex.Atomize(t.Value.(float64)), nil
		case tokenIds["STRING"]:
			return terex.Atomize(t.Value.(string)), nil
		}
	}
	return atom, nil
}

var _ terex.SymbolResolver = symbolPreservingResolver{}
