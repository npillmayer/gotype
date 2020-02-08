package terexlang

import (
	"fmt"
	"strings"
	"sync"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/npillmayer/gotype/syntax/termr"
)

// T traces to the global syntax tracer
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

// --- Grammar ---------------------------------------------------------------

// Atom       ::=  '^' Atom   // currently un-ambiguated by QuoteOrAtom
// Atom       ::=  ident
// Atom       ::=  string
// Atom       ::=  int
// Atom       ::=  float      // TODO unify this to number
// Atom       ::=  List
// List       ::=  '(' Sequence ')'
// Sequence   ::=  Sequence Atom
// Sequence   ::=  Sequence
func makeTermRGrammar() (*lr.LRAnalysis, error) {
	b := lr.NewGrammarBuilder("TermR")
	b.LHS("QuoteOrAtom").N("Quote").End()
	b.LHS("QuoteOrAtom").N("Atom").End()
	b.LHS("Quote").T("quote", '^').N("Atom").End()
	b.LHS("Atom").T("ident", scanner.Ident).End()
	b.LHS("Atom").T("string", scanner.String).End()
	b.LHS("Atom").T("int", scanner.Int).End()
	b.LHS("Atom").T("float", scanner.Float).End()
	b.LHS("Atom").N("List").End()
	b.LHS("List").T("(", '(').N("Sequence").T(")", ')').End()
	b.LHS("Sequence").N("Sequence").N("QuoteOrAtom").End()
	b.LHS("Sequence").N("QuoteOrAtom").End()
	g, err := b.Grammar()
	if err != nil {
		return nil, err
	}
	return lr.Analysis(g), nil
}

var grammar *lr.LRAnalysis
var astBuilder *termr.ASTBuilder
var startOnce sync.Once // monitors one-time creation of grammar and AST-builder

func createParser() *earley.Parser {
	startOnce.Do(func() {
		var err error
		if grammar, err = makeTermRGrammar(); err != nil {
			panic("Cannot create global grammar")
		}
		initDefaultPatterns()
		astBuilder = termr.NewASTBuilder(grammar.Grammar())
		astBuilder.AddOperator(atomOp)
		astBuilder.AddOperator(quoteOp)
		astBuilder.AddOperator(listOp)
	})
	return earley.NewParser(grammar, earley.GenerateTree(true))
}

func parse(sexpr string, source string) (*sppf.Forest, error) {
	parser := createParser()
	r := strings.NewReader(sexpr)
	// TODO create a (lexmachine?) tokenizer
	scan := scanner.GoTokenizer(source, r, scanner.SkipComments(true))
	accept, err := parser.Parse(scan, nil)
	if err != nil {
		return nil, err
	}
	if !accept {
		return nil, fmt.Errorf("Not a valid expression")
	}
	return parser.ParseForest(), nil
}

// Eval evaluates an s-expr (given in textual form).
// It will parse the string, create an internal S-expr structure and evaluate it,
// using the symbols in env,
func Eval(sexpr string, env *terex.Environment) *terex.GCons {
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
	return ast
}

// --- S-expr AST builder listener -------------------------------------------

var atomOp *sExprOp  // for Atom -> ... productions
var quoteOp *sExprOp // for Quote -> ... productions
var listOp *sExprOp  // for List -> ... productions

// --- AST operator helpers --------------------------------------------------

type sExprOp struct {
	name  string
	rules []termr.RewriteRule
}

func makeASTOp(name string) *sExprOp {
	op := &sExprOp{
		name:  name,
		rules: make([]termr.RewriteRule, 0, 1),
	}
	return op
}

func (op *sExprOp) Name() string {
	return op.name
}

func (op *sExprOp) String() string {
	return op.name + ":Op"
}

func (op *sExprOp) Rewrite(l *terex.GCons, env *terex.Environment) *terex.GCons {
	T().Errorf("%s:Op.Rewrite[%s] called, %d rules", op.Name(), l.ListString(), len(op.rules))
	for _, rule := range op.rules {
		T().Errorf("match: trying %s %% %s ?", rule.Pattern.ListString(), l.ListString())
		if rule.Pattern.Match(l, env) {
			T().Infof("Op %s has a match", op.Name())
			v := rule.Rewrite(l, env)
			T().Infof("Op %s rewrite -> %s", op.Name(), v.ListString())
			//return rule.Rewrite(l, env)
			return v
		}
	}
	return l
}

func (op *sExprOp) Descend(sppf.RuleCtxt) bool {
	return true
}

// Call is part of interface Operator.
func (op *sExprOp) Call(term *terex.GCons) interface{} {
	return op.Rewrite(term, terex.GlobalEnvironment)
}

var _ terex.Operator = &sExprOp{}

func (op *sExprOp) Rule(pattern *terex.GCons, rw termr.Rewriter) *sExprOp {
	r := termr.RewriteRule{
		Pattern: pattern,
		Rewrite: rw,
	}
	op.rules = append(op.rules, r)
	return op
}

// Anything is a pattern matching any s-expr.
var Anything *terex.GCons

// AnyToken is a pattern matching any arg of TokenType
var AnyToken *terex.GCons

// SingleTokenArg is a pattern matching an operator with a single arg of TokenType.
var SingleTokenArg *terex.GCons

// AnyOp is a pattern matching any node with OperatorType
var AnyOp = makeASTOp("!AnyOp")

func initDefaultPatterns() {
	Anything = terex.Cons(terex.Atomize(terex.AnyList), nil)
	T().Errorf("Anything=%s", Anything.ListString())
	arg := terex.Atomize(terex.AnyType)
	AnyToken = terex.Cons(arg, nil)
	SingleTokenArg = terex.Cons(terex.Atomize(AnyOp), AnyToken)
	//p := Cons(makeNode(AnyOp), Anything)
	atomOp = makeASTOp("Atom").Rule(Anything, func(l *terex.GCons, env *terex.Environment) *terex.GCons {
		return l.Cdr
	})
	// .Rule(SingleTokenArg, func(l *terex.GCons, env *terex.Environment) *terex.GCons {
	// 	return l.Cdr()
	// })
	p := terex.Cons(terex.Atomize(AnyOp),
		terex.Cons(terex.Atomize(&terex.Token{Name: "^", Value: '^'}), AnyToken))
	T().Errorf("PATTERN Quote = %s", p.ListString())
	quoteOp = makeASTOp("Quote").Rule(p, func(l *terex.GCons, env *terex.Environment) *terex.GCons {
		return terex.Cons(l.Car, l.Cddr())
	})
	//p = Cons(makeNode(AnyOp), Cons(makeNode(&Token{"'('", '('}), Cons(R, nil)))
	listOp = makeASTOp("List").Rule(Anything, func(l *terex.GCons, env *terex.Environment) *terex.GCons {
		if l.Length() <= 3 { // ( )
			return nil
		}
		content := l.Cddr()                      // strip (
		content = content.FirstN(l.Length() - 1) // strip )
		return terex.Cons(l.Car, content)        // (List:Op ...)
	})
}
