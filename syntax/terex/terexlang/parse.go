package terexlang

import (
	"fmt"
	"sync"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/npillmayer/gotype/syntax/termr"
)

// --- Grammar ---------------------------------------------------------------

// Atom       ::=  '^' Atom   // currently un-ambiguated by QuoteOrAtom
// Atom       ::=  ident
// Atom       ::=  string
// Atom       ::=  int
// Atom       ::=  Op
// Op         ::=  '+' | '-' | '*' | '/' | '='
// Atom       ::=  List
// List       ::=  '(' Sequence ')'
// Sequence   ::=  Sequence Atom
// Sequence   ::=  Sequence
//
// Comments starting with // will be filtered by the scanner.
// ',' counts as whitespace and is discarded as well.
func makeTermRGrammar() (*lr.LRAnalysis, error) {
	b := lr.NewGrammarBuilder("TermR")
	b.LHS("QuoteOrAtom").N("Quote").End()
	b.LHS("QuoteOrAtom").N("Atom").End()
	b.LHS("Quote").T(Token("'")).N("Atom").End()
	b.LHS("Atom").T(Token("ID")).End()
	b.LHS("Atom").T(Token("STRING")).End()
	b.LHS("Atom").T(Token("NUM")).End()
	b.LHS("Atom").N("Op").End()
	b.LHS("Op").T(Token("+")).End()
	b.LHS("Op").T(Token("-")).End()
	b.LHS("Op").T(Token("*")).End()
	b.LHS("Op").T(Token("/")).End()
	b.LHS("Op").T(Token("=")).End()
	b.LHS("Atom").N("List").End()
	b.LHS("List").T(Token("(")).N("Sequence").T(Token(")")).End()
	b.LHS("Sequence").N("Sequence").N("QuoteOrAtom").End()
	b.LHS("Sequence").N("QuoteOrAtom").End()
	g, err := b.Grammar()
	if err != nil {
		return nil, err
	}
	return lr.Analysis(g), nil
}

var grammar *lr.LRAnalysis
var lexer *scanner.LMAdapter
var astBuilder *termr.ASTBuilder
var startOnce sync.Once // monitors one-time creation of grammar, lexer and AST-builder

func createParser() *earley.Parser {
	startOnce.Do(func() {
		var err error
		T().Infof("Creating lexer")
		if lexer, err = Lexer(); err != nil { // MUST be called before grammar builing !
			panic("Cannot create lexer")
		}
		T().Infof("Creating grammar")
		if grammar, err = makeTermRGrammar(); err != nil {
			panic("Cannot create global grammar")
		}
		initDefaultPatterns()
		astBuilder = termr.NewASTBuilder(grammar.Grammar())
		astBuilder.AddTermR(atomOp)
		astBuilder.AddTermR(opOp)
		astBuilder.AddTermR(quoteOp)
		astBuilder.AddTermR(seqOp)
		astBuilder.AddTermR(listOp)
	})
	return earley.NewParser(grammar, earley.GenerateTree(true))
}

func parse(input string, source string) (*sppf.Forest, error) {
	parser := createParser()
	scan, err := lexer.Scanner(input)
	if err != nil {
		return nil, err
	}
	//gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	accept, err := parser.Parse(scan, nil)
	//T().Infof("accept=%v, input=%s", accept, input)
	if err != nil {
		return nil, err
	} else if !accept {
		return nil, fmt.Errorf("Not a valid TeREx expression")
	}
	return parser.ParseForest(), nil
}

// Eval evaluates an s-expr (given in textual form).
// It will parse the string, create an internal S-expr structure and evaluate it,
// using the symbols in env,
func Eval(sexpr string, env *terex.Environment) *terex.GCons {
	ast := Quote(sexpr, env)
	r := env.Eval(ast)
	T().Infof("eval(AST) = %s", r.ListString())
	return r
}

// Quote qotes an s-expr (given in textual form).
// It will parse the string, create an internal S-expr structure and quote it,
// using the symbols in env,
func Quote(sexpr string, env *terex.Environment) *terex.GCons {
	parsetree, err := parse(sexpr, "eval")
	if err != nil {
		//env.lastError = err
		T().Errorf("Eval parsing error: %v", err)
		return nil
	}
	if parsetree == nil {
		T().Errorf("Empty eval() parse tree")
		return nil
	}
	// tmpfile, err := ioutil.TempFile(".", "eval-parsetree-*.dot")
	// if err != nil {
	// 	panic("cannot open tmp file")
	// }
	// sppf.ToGraphViz(parsetree, tmpfile)
	// T().Errorf("Exported parse tree to %s", tmpfile.Name())
	ast, _ := astBuilder.AST(parsetree)
	if ast == nil {
		T().Errorf("Cannot create AST from parsetree")
		return nil
	}
	T().Infof("AST: %s", ast.ListString())
	r := env.Quote(ast)
	T().Infof("quote(AST) = %s", r.ListString())
	return r
}

// --- S-expr AST builder listener -------------------------------------------

var atomOp *sExprTermR  // for Atom -> ... productions
var opOp *sExprTermR    // for Op -> ... productions
var quoteOp *sExprTermR // for Quote -> ... productions
var seqOp *sExprTermR   // for Sequence -> ... productions
var listOp *sExprTermR  // for List -> ... productions

type sExprTermR struct {
	name   string
	opname string
	rules  []termr.RewriteRule
	call   func(terex.Element, *terex.Environment) terex.Element
	quote  func(terex.Element, *terex.Environment) terex.Element
}

func makeASTTermR(name string, opname string) *sExprTermR {
	termr := &sExprTermR{
		name:   name,
		opname: opname,
		rules:  make([]termr.RewriteRule, 0, 1),
	}
	return termr
}

func (op *sExprTermR) Name() string {
	return op.name
}

func (op *sExprTermR) Operator() terex.Operator {
	return &globalOpInEnv{op.opname}
}

func (op *sExprTermR) Rewrite(l *terex.GCons, env *terex.Environment) terex.Element {
	T().Debugf("%s:Op.Rewrite[%s] called, %d rules", op.Name(), l.ListString(), len(op.rules))
	for _, rule := range op.rules {
		T().Infof("match: trying %s %% %s ?", rule.Pattern.ListString(), l.ListString())
		if rule.Pattern.Match(l, env) {
			T().Infof("Op %s has a match", op.Name())
			v := rule.Rewrite(l, env)
			T().Infof("Op %s rewrite -> %s", op.Name(), v.String())
			//return rule.Rewrite(l, env)
			return v
		}
	}
	return terex.Elem(l)
}

func (op *sExprTermR) Descend(sppf.RuleCtxt) bool {
	return true
}

func (op *sExprTermR) Rule(pattern *terex.GCons, rw termr.Rewriter) *sExprTermR {
	r := termr.RewriteRule{
		Pattern: pattern,
		Rewrite: rw,
	}
	op.rules = append(op.rules, r)
	return op
}

// SingleTokenArg is a pattern matching an operator with a single arg of TokenType.
var SingleTokenArg *terex.GCons

func initDefaultPatterns() {
	SingleTokenArg = terex.Cons(terex.Atomize(terex.OperatorType), termr.AnySymbol())
	atomOp = makeASTTermR("Atom", "").Rule(termr.Anything(), func(l *terex.GCons, env *terex.Environment) terex.Element {
		// l=(atom x),  => x     // should be standard behaviour, thus unneccessary
		return terex.Elem(l.Cdar()) // an atom simply unwraps its arg
	})
	opOp = makeASTTermR("Op", "").Rule(termr.Anything(), func(l *terex.GCons, env *terex.Environment) terex.Element {
		if l.Length() <= 1 || l.Cdar().Type() != terex.TokenType {
			return terex.Elem(l)
		}
		// (op "x") => x:op
		tname := l.Cdar().Data.(*terex.Token).String()
		T().Errorf("FOUND TOKEN=%s", tname)
		if sym := terex.GlobalEnvironment.FindSymbol(tname, true); sym != nil {
			T().Errorf("FOUND SYM=%v", sym)
			if sym.Value.Type() == terex.OperatorType {
				op := terex.Atomize(&globalOpInEnv{tname})
				return terex.Elem(op)
			}
		}
		return terex.Elem(l)
	})
	_, tokval := Token("'")
	p := terex.Cons(terex.Atomize(terex.OperatorType),
		terex.Cons(terex.Atomize(&terex.Token{Name: "'", Value: tokval}), termr.AnySymbol()))
	quoteOp = makeASTTermR("Quote", "quote").Rule(p, func(l *terex.GCons, env *terex.Environment) terex.Element {
		return terex.Elem(terex.Cons(l.Car, l.Cddr()))
	})
	seqOp = makeASTTermR("Sequence", "!seq").Rule(termr.Anything(), func(l *terex.GCons, env *terex.Environment) terex.Element {
		if l.Cdar().Type() == terex.ConsType {
			seq := l.Cdr.Tee().Concat(l.Cddr())
			return terex.Elem(seq)
		} else if l.Cddr() == nil {
			return terex.Elem(l.Cdar())
		}
		return terex.Elem(l.Cdr)
	})
	listOp = makeASTTermR("List", "list").Rule(termr.Anything(), func(l *terex.GCons, env *terex.Environment) terex.Element {
		// list '(' x y ... ')'  => (#list x y ...)
		if l.Length() <= 3 { // ( )
			return terex.Elem(nil)
		}
		content := l.Cddr()                            // strip (
		content = content.FirstN(content.Length() - 1) // strip )
		return terex.Elem(terex.Cons(l.Car, content))  // (List:Op ...)
	})
}

// ---------------------------------------------------------------------------

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
		T().Errorf("Cannot quote-call parsing operation %s", op.opname)
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
	if op.String() == "list" {
		T().Infof("========================================================")
		T().Infof("=====     LIST         =================================")
		T().Infof("========================================================")
	}
	return operator.Quote(term, env)
}