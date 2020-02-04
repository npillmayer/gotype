package gospel

import (
	"fmt"
	"strings"
	"sync"

	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
)

func makeGospelGrammar() (*lr.LRAnalysis, error) {
	b := lr.NewGrammarBuilder("Gospel")
	b.LHS("List").T("(", '(').N("Sequence").T(")", ')').End()
	b.LHS("QuotedList").T("quote", '^').T("(", '(').N("Sequence").T(")", ')').End()
	b.LHS("Sequence").N("Sequence").N("Atom").End()
	b.LHS("Sequence").N("Atom").End()
	b.LHS("Atom").T("quote", '^').T("ident", scanner.Ident).End()
	b.LHS("Atom").T("ident", scanner.Ident).End()
	b.LHS("Atom").T("string", scanner.String).End()
	b.LHS("Atom").T("int", scanner.Int).End()
	b.LHS("Atom").T("float", scanner.Float).End()
	b.LHS("Atom").N("List").End()
	b.LHS("Atom").N("QuotedList").End()
	g, err := b.Grammar()
	if err != nil {
		return nil, err
	}
	return lr.Analysis(g), nil
}

var grammar *lr.LRAnalysis
var startOnce sync.Once

func createParser() *earley.Parser {
	var ga *lr.LRAnalysis
	startOnce.Do(func() {
		var err error
		if ga, err = makeGospelGrammar(); err != nil {
			panic("Cannot create global grammar")
		}
	})
	return earley.NewParser(ga, earley.GenerateTree(true))
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

type Symbol struct {
	Name string
	atom Atom
}

type Environment struct {
	name      string
	parent    *Environment
	dict      map[string]*Symbol
	lastError error
}

var globalEnvironment *Environment = &Environment{
	name: "global",
	dict: make(map[string]*Symbol),
}

func NewEnvironment(name string, parent *Environment) *Environment {
	if parent == nil {
		parent = globalEnvironment
	}
	return &Environment{
		name:   "global",
		parent: parent,
		dict:   make(map[string]*Symbol),
	}
}

func (env *Environment) Eval(prog string) *GCons {
	parsetree, err := parse(prog, "eval")
	if err != nil {
		env.lastError = err
		T().Errorf("Eval parsing error: %v", err)
		return nil
	}
	if parsetree == nil {
		T().Errorf("Empty eval() parse tree")
	}
	return nil
}

func (env *Environment) Intern(name string, inherit bool) *Symbol {
	var sym *Symbol
	var ok bool
	for env != nil {
		sym, ok = env.dict[name]
		if ok {
			return sym
		}
		if inherit {
			env = env.parent
		}
	}
	sym = &Symbol{Name: name, atom: NullAtom}
	env.dict[name] = sym
	return sym
}

func (env *Environment) listFromParseTree(tree *sppf.Forest) *GCons {
	listener := newListener()
	if listener != nil {
		// remove this if
	}
	// TODO
	return nil
}

type Listener struct {
	list          *GCons
	dispatchEnter map[string]ruleEnterOp
	dispatchExit  map[string]ruleOp
}
type ruleEnterOp func(*lr.Symbol, []*sppf.RuleNode, lr.Span, int) bool
type ruleOp func(*lr.Symbol, []*sppf.RuleNode, lr.Span, int) interface{}

func newListener() *Listener {
	l := &Listener{}
	l.dispatchEnter = map[string]ruleEnterOp{
		"QuotedList": l.EnterQuotedList,
	}
	l.dispatchExit = map[string]ruleOp{
		"List":     l.ExitList,
		"Sequence": l.ExitSequence,
		"Atom":     l.ExitAtom,
	}
	return l
}

func (l *Listener) EnterRule(lhs *lr.Symbol, rhs []*sppf.RuleNode, span lr.Span, level int) bool {
	if r, ok := l.dispatchEnter[lhs.Name]; ok {
		return r(lhs, rhs, span, level)
	}
	T().Debugf("%senter grammar symbol: %v", indent(level), lhs)
	return true
}

func (l *Listener) EnterQuotedList(lhs *lr.Symbol, rhs []*sppf.RuleNode, span lr.Span, level int) bool {
	if len(rhs) <= 2 {
		T().Debugf("%sempty Q-LIST", indent(level))
		return false
	}
	values := make([]interface{}, len(rhs)-2)
	for i, r := range rhs[1 : len(rhs)-1] {
		values[i-1] = r.Value
	}
	T().Debugf("%senter Q-LIST (%v)\n", indent(level), values)
	return true
}

func (l *Listener) ExitRule(lhs *lr.Symbol, rhs []*sppf.RuleNode, span lr.Span, level int) interface{} {
	if r, ok := l.dispatchExit[lhs.Name]; ok {
		return r(lhs, rhs, span, level)
	}
	T().Debugf("%sReduce of grammar symbol: %v", indent(level), lhs)
	return rhs[0].Value
}

func (l *Listener) ExitList(lhs *lr.Symbol, rhs []*sppf.RuleNode, span lr.Span, level int) interface{} {
	if len(rhs) <= 2 {
		T().Debugf("%sempty LIST", indent(level))
		return nil
	}
	values := make([]interface{}, len(rhs)-2)
	for i, r := range rhs[1 : len(rhs)-1] {
		values[i-1] = r.Value
	}
	T().Debugf("%sLIST (%v)\n", indent(level), values)
	return List(values)
}

func (l *Listener) ExitSequence(lhs *lr.Symbol, rhs []*sppf.RuleNode, span lr.Span, level int) interface{} {
	T().Debugf("%sSEQUENCE\n", indent(level))
	return nil
}

func (l *Listener) ExitAtom(lhs *lr.Symbol, rhs []*sppf.RuleNode, span lr.Span, level int) interface{} {
	T().Debugf("%sATOM\n", indent(level))
	return nil
}

func (l *Listener) Terminal(tokval int, token interface{}, span lr.Span, level int) interface{} {
	/* TODO check type
	switch tokval {
	case scanner.Ident:
	case scanner.String:
	case scanner.Int:
	case scanner.Float:
	default:
	}
	*/
	a := atomize(token)
	T().Debugf("new Atom(%v) of type %d", token, tokval)
	return a
}

func (l *Listener) Conflict(*lr.Symbol, int, lr.Span, int) (int, error) {
	panic("AST should never be ambiguous")
}

func indent(level int) string {
	in := ""
	for level > 0 {
		in = in + ". "
		level--
	}
	return in
}
