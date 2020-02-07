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
var startOnce sync.Once

func createParser() *earley.Parser {
	startOnce.Do(func() {
		var err error
		if grammar, err = makeTermRGrammar(); err != nil {
			panic("Cannot create global grammar")
		}
	})
	return earley.NewParser(grammar, earley.GenerateTree(true))
}

func parse(prog string, source string) (*sppf.Forest, error) {
	parser := createParser()
	r := strings.NewReader(prog)
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

// Symbol is a type for language symbols (stored in the Environment).
type Symbol struct {
	Name  string
	props properties
	value Node
}

// newSymbol creates a new symbol for a given initial value.
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

// Eval evaluates an S-expr (given in textual form).
// It will parse the string, create an internal S-expr structure and evaluate it,
// using the symbols in env,
func (env *Environment) Eval(prog string) *GCons {
	parsetree, err := parse(prog, "eval")
	if err != nil {
		env.lastError = err
		T().Errorf("Eval parsing error: %v", err)
		return nil
	}
	if parsetree == nil {
		T().Errorf("Empty eval() parse tree")
		return nil
	}
	// tmpfile, err := ioutil.TempFile(".", "parsetree.*.dot")
	// if err != nil {
	// 	panic("cannot open tmp file")
	// }
	// sppf.ToGraphViz(parsetree, tmpfile)
	// T().Errorf("Exported parse tree to %s", tmpfile.Name())
	astbuilder := NewASTBuilder(grammar.Grammar())
	ast, _ := astbuilder.AST(parsetree)
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

// func (env *Environment) listFromParseTree(tree *sppf.Forest) *GCons {
// 	listener := newListener()
// 	if listener != nil {
// 		// remove this if
// 	}
// 	// TODO
// 	return nil
// }

// ---------------------------------------------------------------------------

/* TODO remove this

type Listener struct {
	list          *GCons
	dispatchEnter map[string]ruleEnterOp
	dispatchExit  map[string]ruleOp
}
type ruleEnterOp func(*lr.Symbol, []*sppf.RuleNode, sppf.RuleCtxt) bool
type ruleOp func(*lr.Symbol, []*sppf.RuleNode, sppf.RuleCtxt) interface{}

func newListener() *Listener {
	l := &Listener{}
	l.dispatchEnter = map[string]ruleEnterOp{
		"QuoteOrAtom": l.EnterQuoteOrAtom,
	}
	l.dispatchExit = map[string]ruleOp{
		"List":     l.ExitList,
		"Sequence": l.ExitSequence,
		"Atom":     l.ExitAtom,
	}
	return l
}

func (l *Listener) EnterRule(lhs *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) bool {
	if r, ok := l.dispatchEnter[lhs.Name]; ok {
		return r(lhs, rhs, ctxt)
	}
	T().Debugf("%senter grammar symbol: %v", indent(ctxt.Level), lhs)
	return true
}

func (l *Listener) EnterQuoteOrAtom(lhs *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) bool {
	if len(rhs) <= 2 {
		T().Debugf("%sempty Q-LIST", indent(ctxt.Level))
		return false
	}
	values := make([]interface{}, len(rhs)-2)
	for i, r := range rhs[1 : len(rhs)-1] {
		values[i-1] = r.Value
	}
	T().Debugf("%senter Q-LIST (%v)\n", indent(ctxt.Level), values)
	return true
}

func (l *Listener) ExitRule(lhs *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) interface{} {
	if r, ok := l.dispatchExit[lhs.Name]; ok {
		return r(lhs, rhs, ctxt)
	}
	T().Debugf("%sReduce of grammar symbol: %v", indent(ctxt.Level), lhs)
	return rhs[0].Value
}

func (l *Listener) ExitList(lhs *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) interface{} {
	if len(rhs) <= 2 {
		T().Debugf("%sempty LIST", indent(ctxt.Level))
		return nil
	}
	values := make([]interface{}, len(rhs)-2)
	for i, r := range rhs[1 : len(rhs)-1] {
		values[i-1] = r.Value
	}
	T().Debugf("%sLIST (%v)\n", indent(ctxt.Level), values)
	return List(values)
}

func (l *Listener) ExitSequence(lhs *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) interface{} {
	T().Debugf("%sSEQUENCE\n", indent(ctxt.Level))
	return nil
}

func (l *Listener) ExitAtom(lhs *lr.Symbol, rhs []*sppf.RuleNode, ctxt sppf.RuleCtxt) interface{} {
	T().Debugf("%sATOM\n", indent(ctxt.Level))
	return nil
}

func (l *Listener) Terminal(tokval int, token interface{}, ctxt sppf.RuleCtxt) interface{} {
	// switch tokval {
	// case scanner.Ident:
	// case scanner.String:
	// case scanner.Int:
	// case scanner.Float:
	// default:
	// }
	a := atomize(token)
	T().Debugf("new Atom(%v) of type %d", token, tokval)
	return a
}

func (l *Listener) Conflict(*lr.Symbol, sppf.RuleCtxt) (int, error) {
	panic("AST should never be ambiguous")
}

func (l *Listener) MakeAttrs(*lr.Symbol) interface{} {
	return nil
}
*/

func indent(level int) string {
	in := ""
	for level > 0 {
		in = in + ". "
		level--
	}
	return in
}
