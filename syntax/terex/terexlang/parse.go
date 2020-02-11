package terexlang

import (
	"fmt"
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

// b.LHS("Quote").T("quote", '^').N("Atom").End()
// b.LHS("Atom").T("ident", scanner.Ident).End()
// b.LHS("Atom").T("string", scanner.String).End()
// b.LHS("Atom").T("int", scanner.Int).End()
// b.LHS("Atom").T("float", scanner.Float).End()

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
	b.LHS("Quote").T(Token("'")).N("Atom").End()
	b.LHS("Atom").T(Token("ID")).End()
	b.LHS("Atom").T(Token("STRING")).End()
	b.LHS("Atom").T(Token("NUM")).End()
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
var startOnce sync.Once // monitors one-time creation of grammar and AST-builder

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
		astBuilder.AddTermR(quoteOp)
		astBuilder.AddTermR(listOp)
	})
	return earley.NewParser(grammar, earley.GenerateTree(true))
}

func parse(input string, source string) (*sppf.Forest, error) {
	parser := createParser()
	//r := strings.NewReader(sexpr)
	// TODO create a (lexmachine?) tokenizer
	//scan := scanner.GoTokenizer(source, r, scanner.SkipComments(true))
	lexer, _ := Lexer()
	scan, err := lexer.Scanner(input)
	if err != nil {
		return nil, err
	}
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelInfo)
	accept, err := parser.Parse(scan, nil)
	T().Errorf("accept=%v, input=%s", accept, input)
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
	r := env.Quote(ast)
	T().Infof("quote(AST) = %s", r.ListString())
	return r
}

// --- S-expr AST builder listener -------------------------------------------

var atomOp *sExprTermR  // for Atom -> ... productions
var quoteOp *sExprTermR // for Quote -> ... productions
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

func (op *sExprTermR) Rewrite(l *terex.GCons, env *terex.Environment) *terex.GCons {
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
	//arg := terex.Atomize(terex.AnyType)
	//AnyToken = terex.Cons(arg, nil)
	SingleTokenArg = terex.Cons(terex.Atomize(terex.OperatorType), termr.AnySymbol())
	//p := Cons(makeNode(AnyOp), Anything)
	atomOp = makeASTTermR("Atom", "").Rule(termr.Anything(), func(l *terex.GCons, env *terex.Environment) *terex.GCons {
		return l.Cdr
	})
	// .Rule(SingleTokenArg, func(l *terex.GCons, env *terex.Environment) *terex.GCons {
	// 	return l.Cdr()
	// })
	_, tokval := Token("'")
	p := terex.Cons(terex.Atomize(terex.OperatorType),
		terex.Cons(terex.Atomize(&terex.Token{Name: "'", Value: tokval}), termr.AnySymbol()))
	//T().Debugf("PATTERN Quote = %s", p.ListString())
	quoteOp = makeASTTermR("Quote", "quote").Rule(p, func(l *terex.GCons, env *terex.Environment) *terex.GCons {
		return terex.Cons(l.Car, l.Cddr())
	})
	//p = Cons(makeNode(AnyOp), Cons(makeNode(&Token{"'('", '('}), Cons(R, nil)))
	listOp = makeASTTermR("List", "list").Rule(termr.Anything(), func(l *terex.GCons, env *terex.Environment) *terex.GCons {
		if l.Length() <= 3 { // ( )
			return nil
		}
		content := l.Cddr()                            // strip (
		content = content.FirstN(content.Length() - 1) // strip )
		return terex.Cons(l.Car, content)              // (List:Op ...)
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
	return operator.Quote(term, env)
}
