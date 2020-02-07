package termr

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/iteratable"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
)

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
var astBuilder *ASTBuilder
var startOnce sync.Once // monitors one-time creation of grammar and AST-builder

func createParser() *earley.Parser {
	startOnce.Do(func() {
		var err error
		if grammar, err = makeTermRGrammar(); err != nil {
			panic("Cannot create global grammar")
		}
		initDefaultPatterns()
		astBuilder = NewASTBuilder(grammar.Grammar())
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

// --- Symbols ---------------------------------------------------------------

// Symbol is a type for language symbols (stored in the Environment).
// A symbol can change its value. A value may be any atom or s-expr type.
// A value of nil means the symbol is not yet bound.
type Symbol struct {
	Name  string
	props properties
	value Node
}

// newSymbol creates a new symbol for a given initial value (which may be nil).
func newSymbol(name string, thing interface{}) *Symbol {
	return &Symbol{
		Name:  name,
		props: makeProps(),
		value: makeNode(thing),
	}
}

func (sym Symbol) String() string {
	return fmt.Sprintf("%s:%s(%s)", sym.Name, sym.value.Type().String(), sym.value.String())
}

// IsAtom returns true if a symbol represents an atom (not a cons).
func (sym Symbol) IsAtom() bool {
	return sym.value.Type() != ConsType
}

// Get returns a property value for a given key.
func (sym *Symbol) Get(key string) Atom {
	value := sym.props.Get(key, sym)
	if value == nil {
		return NullAtom
	}
	return atomize(value)
}

type properties []keyvalue
type keyvalue struct {
	Key   string
	Value interface{}
}

func makeProps() properties {
	return make([]keyvalue, 0, 2)
}

func (p properties) Get(key string, sym *Symbol) interface{} {
	if key == "type" {
		return sym.value.Type()
	}
	if key == "name" {
		return sym.Name
	}
	for _, kv := range p {
		if kv.Key == key {
			return kv.Value
		}
	}
	return nil
}

func (p properties) Set(key string, value interface{}) properties {
	found := false
	for _, kv := range p {
		if kv.Key == key {
			kv.Value = value
			found = true
		}
	}
	if !found {
		return append(p, keyvalue{key, value})
	}
	return p
}

// Token represents a grammar terminal, and a corresponding input token, respectively.
type Token struct {
	Name  string
	value int
}

func (t Token) String() string {
	return t.Name
}

// --- Environments ----------------------------------------------------------

// Environment is a type for a symbol environment.
type Environment struct {
	name      string
	parent    *Environment
	dict      map[string]*Symbol
	lastError error
}

var globalEnvironment *Environment = &Environment{
	name: "#global",
	dict: make(map[string]*Symbol),
}

// NewEnvironment creates a new environment.
func NewEnvironment(name string, parent *Environment) *Environment {
	if parent == nil {
		parent = globalEnvironment
	}
	return &Environment{
		name:   name,
		parent: parent,
		dict:   make(map[string]*Symbol),
	}
}

// Eval evaluates an s-expr (given in textual form).
// It will parse the string, create an internal S-expr structure and evaluate it,
// using the symbols in env,
func (env *Environment) Eval(sexpr string) *GCons {
	parsetree, err := parse(sexpr, "eval")
	if err != nil {
		env.lastError = err
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
	//
	//astbuilder := NewASTBuilder(grammar.Grammar())
	ast, _ := astBuilder.AST(parsetree)
	if ast == nil {
		T().Errorf("Cannot create AST from parsetree")
		return nil
	}
	T().Infof("AST: %s", ast.ListString())
	return ast
	// cursor := parsetree.SetCursor(nil, nil)
	// value := cursor.TopDown(newListener(), sppf.LtoR, sppf.Break)
	// T().Debugf("return value of top-down traversal: %v", value)
}

// FindSymbol checks wether a symbol is defined in env and returns it, if found.
// Otherwise nil is returned.
func (env *Environment) FindSymbol(name string, inherit bool) *Symbol {
	var sym *Symbol
	var ok bool
	for env != nil {
		sym, ok = env.dict[name]
		if ok {
			return sym
		}
		if inherit {
			env = env.parent
		} else {
			break
		}
	}
	return nil
}

// Intern interns a symbol name as a symbol, returning a reference to that symbol.
// If the symbol already exists, the existing symbol is returned.
// Parameter inherit dictates wether ancestor environments should be searched, too,
// to detect the symbol.
func (env *Environment) Intern(name string, inherit bool) *Symbol {
	sym := env.FindSymbol(name, inherit)
	if sym == nil {
		sym = &Symbol{Name: name}
		env.dict[name] = sym
	}
	return sym
}

func (env *Environment) String() string {
	return env.name
}

// Dump is a debugging helper, listing all known symbols in env.
func (env *Environment) Dump() string {
	var b bytes.Buffer
	b.WriteString(env.String())
	b.WriteString(" {\n")
	for k, v := range env.dict {
		b.WriteString(fmt.Sprintf("    %s = %v\n", k, v))
	}
	b.WriteString("}\n")
	return b.String()
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
func EnvironmentForGrammarSymbol(symname string, G *lr.Grammar) (*Environment, error) {
	if G == nil {
		return globalEnvironment, errors.New("Grammar is null")
	}
	envname := "#" + symname
	if env := globalEnvironment.FindSymbol(envname, false); env != nil {
		if env.value.Type() != EnvironmentType {
			panic(fmt.Errorf("Internal error, environment misconstructed: %s", envname))
		}
		return env.value.atom.Data.(*Environment), nil
	}
	gsym := G.SymbolByName(symname)
	if gsym == nil || gsym.IsTerminal() {
		return globalEnvironment, fmt.Errorf("Non-terminal not found in grammar: %s", symname)
	}
	env := NewEnvironment(envname, nil)
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
			sym.value = makeNode(gsym.Value)
		}
		// else {
		// 	sym.atom.typ = SymbolType
		// }
	}
	// e := globalEnvironment.Intern(envname, false)
	// e.atom.typ = EnvironmentType
	// e.atom.Data = env
	return env, nil
}

// --- S-expr AST builder listener -------------------------------------------

var atomOp *sExprOp  // for Atom -> ... productions
var quoteOp *sExprOp // for Quote -> ... productions
var listOp *sExprOp  // for List -> ... productions

// --- AST operator helpers --------------------------------------------------

type sExprOp struct {
	name  string
	rules []RewriteRule
}

func makeASTOp(name string) *sExprOp {
	op := &sExprOp{
		name:  name,
		rules: make([]RewriteRule, 0, 1),
	}
	return op
}

func (op *sExprOp) Name() string {
	return op.name
}

func (op *sExprOp) String() string {
	return op.name + ":Op"
}

func (op *sExprOp) Rewrite(l *GCons, env *Environment) *GCons {
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
func (op *sExprOp) Call(term *GCons) *GCons {
	return op.Rewrite(term, globalEnvironment)
}

func (op *sExprOp) Rule(pattern *GCons, rw Rewriter) *sExprOp {
	r := RewriteRule{
		Pattern: pattern,
		Rewrite: rw,
	}
	op.rules = append(op.rules, r)
	return op
}

// Anything is a pattern matching any s-expr.
var Anything *GCons

// AnyToken is a pattern matching any arg of TokenType
var AnyToken *GCons

// SingleTokenArg is a pattern matching an operator with a single arg of TokenType.
var SingleTokenArg *GCons

// AnyOp is a pattern matching any node with OperatorType
var AnyOp = makeASTOp("!AnyOp")

func initDefaultPatterns() {
	Anything = Cons(makeNode(1), nil)
	T().Errorf("Anything=%s", Anything.ListString())
	Anything.car.atom.typ = AnyList
	arg := makeNode(nil)
	arg.atom.typ = AnyType
	AnyToken = Cons(arg, nil)
	SingleTokenArg = Cons(makeNode(AnyOp), AnyToken)
	//p := Cons(makeNode(AnyOp), Anything)
	atomOp = makeASTOp("Atom").Rule(Anything, func(l *GCons, env *Environment) *GCons {
		return l.Cdr()
	})
	// .Rule(SingleTokenArg, func(l *GCons, env *Environment) *GCons {
	// 	return l.Cdr()
	// })
	p := Cons(makeNode(AnyOp), Cons(makeNode(&Token{"^", '^'}), AnyToken))
	T().Errorf("PATTERN Quote = %s", p.ListString())
	quoteOp = makeASTOp("Quote").Rule(p, func(l *GCons, env *Environment) *GCons {
		return Cons(l.car, l.Cddr())
	})
	//p = Cons(makeNode(AnyOp), Cons(makeNode(&Token{"'('", '('}), Cons(R, nil)))
	listOp = makeASTOp("List").Rule(Anything, func(l *GCons, env *Environment) *GCons {
		if l.Length() <= 3 { // ( )
			return nil
		}
		content := l.Cddr()                     // strip (
		content = content.First(l.Length() - 1) // strip )
		return Cons(l.car, content)             // (List:Op ...)
	})
}
