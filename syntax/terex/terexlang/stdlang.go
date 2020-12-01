package terexlang

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/npillmayer/gotype/syntax/terex/termr"
)

// LoadStandardLanguage returns an environment with all the standard
// symbols and operators loaded.
//
func LoadStandardLanguage() *terex.Environment {
	env := terex.NewEnvironment("stdlang", nil)
	makeLispOps(env)
	makeArithmOps(env)
	makeNamespaceOps(env)
	makeParserOps(env)
	return env
}

func makeArithmOps(env *terex.Environment) {
	plus := makeOp(terex.NumType, -1, func(args *terex.GCons, env *terex.Environment) terex.Element {
		result := args.Reduce(func(a terex.Atom, b terex.Atom) terex.Atom {
			return terex.Atomize(a.Data.(float64) + b.Data.(float64))
		}, terex.Atomize(0), env)
		return terex.Elem(result)
	})
	env.Defn("+", plus)
	minus := makeOp(terex.NumType, -1, func(args *terex.GCons, env *terex.Environment) terex.Element {
		if args.Length() == 0 {
			return terex.Elem(terex.Atomize(0))
		}
		args.Car.Data = -args.Car.Data.(float64)
		result := args.Reduce(func(a terex.Atom, b terex.Atom) terex.Atom {
			return terex.Atomize(a.Data.(float64) - b.Data.(float64))
		}, terex.Atomize(0), env)
		args.Car.Data = -args.Car.Data.(float64)
		return terex.Elem(result)
	})
	env.Defn("-", minus)
	times := makeOp(terex.NumType, -1, func(args *terex.GCons, env *terex.Environment) terex.Element {
		if args.Length() == 0 {
			return terex.Elem(terex.Atomize(0))
		}
		result := args.Reduce(func(a terex.Atom, b terex.Atom) terex.Atom {
			return terex.Atomize(a.Data.(float64) * b.Data.(float64))
		}, terex.Atomize(1), env)
		return terex.Elem(result)
	})
	env.Defn("*", times)
	divide := makeOp(terex.NumType, -1, func(args *terex.GCons, env *terex.Environment) terex.Element {
		if args.Length() == 0 {
			return terex.Elem(terex.Atomize(0))
		} else if args.Length() == 1 {
			return terex.Elem(args.Car)
		}
		result := args.Cdr.Reduce(func(a terex.Atom, b terex.Atom) terex.Atom {
			if a.Type() == terex.ErrorType {
				return a // propagate error
			} else if b.Data.(float64) == 0 {
				T().Debugf("Division by zero")
				return terex.ErrorAtom("Division by zero")
			}
			return terex.Atomize(a.Data.(float64) / b.Data.(float64))
		}, args.Car, env)
		return terex.Elem(result)
	})
	env.Defn("/", divide)
}

func makeLispOps(env *terex.Environment) {
	env.Defn("def", func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		//args = args.Map(terex.EvalAtom, env)
		if args.Length() != 2 {
			return ErrorPacker("Wrong number of arguments for def", env)
		}
		if args.Car.Type() != terex.VarType {
			return ErrorPacker("Missing symbol for def", env)
		}
		arg := terex.EvalAtom(terex.Elem(args.Cdar()), env)
		sym := env.Intern(args.Car.Data.(*terex.Symbol).Name, true)
		sym.Value = arg
		return terex.Elem(terex.Atomize(sym))
	})
	env.Defn("print", func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		//args = args.Map(terex.EvalAtom, env)
		args = args.Map(resolve, env)
		args.Map(func(e terex.Element, env *terex.Environment) terex.Element {
			// T().Infof(e.String())
			print(e)
			return terex.Elem(nil)
		}, env)
		last := args.Last()
		if last == nil {
			return terex.Elem(terex.NilAtom)
		}
		return terex.Elem(last.Car)
	})
	env.Defn("list", func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		args = args.Map(terex.EvalAtom, env)
		return terex.Elem(args)
	})
	env.Defn("quote", func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		if args.Length() != 1 {
			return ErrorPacker("Wrong number of arguments for quote", env)
		}
		last := args.Last()
		if last == nil {
			return terex.Elem(terex.NilAtom)
		}
		//terex.Elem(last).First().Dump(tracing.LevelError)
		return terex.Elem(last).First()
	})
	env.Defn("assoc", func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		if args.Length() != 3 {
			return ErrorPacker("Wrong number of arguments for assoc", env)
		}
		if args.Car.Type() != terex.VarType {
			return ErrorPacker("Missing dict symbol for assoc", env)
		}
		k := terex.Eval(terex.Elem(args.Cdar()), env)
		if k.AsAtom().Type() != terex.StringType {
			return ErrorPacker("Missing key (string) for assoc", env)
		}
		key := k.AsAtom().Data.(string)
		value := terex.Eval(terex.Elem(args.Last().Car), env)
		sym := env.Intern(args.Car.Data.(*terex.Symbol).Name, true)
		if sym.Value.IsNil() {
			sym.Value = terex.Elem(make(Dict))
		}
		d := sym.AtomValue().(Dict)
		d[key] = value.AsAtom()
		return terex.Elem(sym)
	})
	env.Defn("get", func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		if args.Length() != 2 {
			return ErrorPacker("Wrong number of arguments for get", env)
		}
		if args.Car.Type() != terex.VarType {
			return ErrorPacker("Missing dict symbol for get", env)
		}
		k := terex.Eval(terex.Elem(args.Cdar()), env)
		if k.AsAtom().Type() != terex.StringType {
			return ErrorPacker("Missing key (string) for get", env)
		}
		key := k.AsAtom().Data.(string)
		sym := env.FindSymbol(args.Car.Data.(*terex.Symbol).Name, true)
		if sym == nil {
			return ErrorPacker("Unable to resolve symbol as dict", env)
		}
		d := sym.AtomValue().(Dict)
		T().Errorf("DICT GET = %v", terex.Elem(d[key]).AsList().ListString())
		//return terex.Elem(d[key]) // TODO
		return terex.Elem(terex.Atomize(d[key]))
	})
}

type Dict map[string]terex.Atom

func (d Dict) getMap() map[string]terex.Atom {
	return map[string]terex.Atom(d)
}

func makeNamespaceOps(env *terex.Environment) {
	env.Def("#ns#", terex.Elem(env))
	env.Def("t", terex.Elem(true))
	env.Def("nil", terex.Elem(nil))
}

func makeParserOps(env *terex.Environment) {
	env.Defn("parse", func(e terex.Element, env *terex.Environment) terex.Element {
		// (parse G "input") => tree
		args := e.AsList()
		args = args.Map(terex.EvalAtom, env)
		if args.Length() != 2 {
			return ErrorPacker("Wrong number of arguments for parse", env)
		}
		if args.Car.Type() != terex.UserType || args.Cdar().Type() != terex.StringType {
			return ErrorPacker("Wrong argument type for parse", env)
		}
		if args.Cdar().Data == nil {
			return ErrorPacker("Cannot parse nil-input", env)
		}
		input := args.Cdar().Data.(string)
		g := args.Car.Data.(*lr.LRAnalysis)
		tree, retr, err := parseAny(g, input)
		if err != nil {
			return ErrorPacker(err.Error(), env)
		}
		result := &parsetree{
			G:    g,
			tree: tree,
			retr: retr,
		}
		return terex.Elem(terex.Atomize(result))
	})
	ast := makeOp(terex.UserType, 2, func(args *terex.GCons, env *terex.Environment) terex.Element {
		// (ast <rew> <T>) => AST
		T().Errorf("AST ARGS = %s", args.ListString())
		if rews, ok1 := args.Car.Data.([]*Rewriter); ok1 {
			if tree, ok2 := args.Cdar().Data.(*parsetree); ok2 {
				ab := termr.NewASTBuilder(tree.G.Grammar())
				for _, r := range rews {
					ab.AddTermR(r)
				}
				env := ab.AST(tree.tree, tree.retr)
				if env == nil {
					return ErrorPacker("Error while creating AST", env)
				}
				T().Infof("\n\n" + env.AST.IndentedListString())
				return terex.Elem(env.AST)
			}
		}
		return ErrorPacker("Arguments for ast should be rewriters and parse tree", env)
	})
	env.Defn("ast", ast)
	rew := makeOp(terex.UserType, 1, func(args *terex.GCons, env *terex.Environment) terex.Element {
		// (rewriters dict) => UType
		//ab := termr.NewASTBuilder(&lr.Grammar{Name: "G"})
		//T().Errorf("AB created")
		if dict, ok := args.Car.Data.(Dict); ok {
			T().Errorf("dict found")
			rewriters := createASTRewritersFromDict(dict)
			return terex.Elem(terex.Atomize(rewriters))
		}
		return ErrorPacker("Expected argument to be Dict", env)
	})
	env.Defn("rewriters", rew)
	env.Defn("match?", func(e terex.Element, env *terex.Environment) terex.Element {
		// (match? <rewriters> key cons) => cons
		args := e.AsList()
		args = args.Map(terex.EvalAtom, env)
		if args.Length() != 3 {
			return ErrorPacker("Wrong number of arguments for match", env)
		}
		if args.Car.Type() != terex.UserType {
			return ErrorPacker("Expected a rewriter-set for match", env)
		}
		if args.Cdar().Type() != terex.StringType && args.Cdar().Type() != terex.VarType {
			return ErrorPacker("Expected a string-like argument as key for match", env)
		}
		input := terex.Elem(args.Cddar())
		key := ""
		if args.Cdar().Type() == terex.StringType {
			key = args.Cdar().Data.(string)
		} else if args.Cdar().Type() == terex.VarType {
			key = args.Cdar().Data.(*terex.Symbol).Name
		}
		T().Errorf("MATCH for key = '%s'", key)
		e = terex.Elem(nil)
		if rews, ok := args.Car.Data.([]*Rewriter); ok {
			for _, rew := range rews {
				if rew.String() == key {
					m := rew.Match(input, env)
					if !m.IsNil() {
						e = m
						break
					}
				}
			}
		}
		return terex.Elem(e)
	})
}

type parsetree struct {
	G    *lr.LRAnalysis
	tree *sppf.Forest
	retr termr.TokenRetriever
}

// --- Helpers ---------------------------------------------------------------

func resolve(el terex.Element, env *terex.Environment) terex.Element {
	if el.IsNil() {
		return el
	}
	if !el.IsAtom() {
		return el
	}
	atom := el.AsAtom()
	resolved, err := terex.DefaultSymbolResolver{}.Resolve(atom, env, false)
	if err != nil {
		return ErrorPacker(err.Error(), env)
	}
	return resolved
}

func makeOp(t terex.AtomType, a int, op func(*terex.GCons, *terex.Environment) terex.Element) func(
	e terex.Element, env *terex.Environment) terex.Element {
	arity := a
	typ := t
	call := func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		args = args.Map(terex.EvalAtom, env)
		if arity > 0 && args.Length() != arity {
			env.Error(errors.New("Wrong number of arguments"))
			return terex.Elem(terex.ErrorAtom("Wrong number of arguments"))
		} else if arity < 0 && args.Length() < -arity {
			env.Error(errors.New("Wrong number of arguments"))
			return terex.Elem(terex.ErrorAtom("Wrong number of arguments"))
		}
		if typ != terex.NoType {
			errAtom := terex.ErrorAtom("")
			errAtom.Data = nil
			errAtom = args.Reduce(func(e terex.Atom, a terex.Atom) terex.Atom {
				var err error
				if e.Data != nil {
					err = e.Data.(error)
				}
				_, err = cast(a, typ, env, err)
				e.Data = err
				return e
			}, errAtom, env)
			if errAtom.Data != nil {
				return terex.Elem(errAtom)
			}
		}
		return op(args, env)
	}
	return call
}

func parseAny(ga *lr.LRAnalysis, input string) (*sppf.Forest, termr.TokenRetriever, error) {
	level := T().GetTraceLevel()
	T().SetTraceLevel(tracing.LevelError)
	parser := earley.NewParser(ga, earley.GenerateTree(true))
	reader := strings.NewReader(input)
	scanner := scanner.GoTokenizer(ga.Grammar().Name, reader)
	acc, err := parser.Parse(scanner, nil)
	if !acc || err != nil {
		return nil, nil, fmt.Errorf("Could not parse input: %v", err)
	}
	T().SetTraceLevel(level)
	T().Infof("Successfully parsed input")
	tokretr := func(pos uint64) interface{} {
		return parser.TokenAt(pos)
	}
	return parser.ParseForest(), tokretr, nil
}

func cast(atom terex.Atom, typ terex.AtomType, env *terex.Environment, err error) (terex.Atom, error) {
	if err != nil {
		return atom, err
	}
	if atom.Type() != typ {
		// TODO try to cast
		T().Errorf("%s cannot be cast to %s", atom, typ.String())
		err := fmt.Errorf("%s cannot be cast to %s", atom, typ.String())
		env.Error(err)
		return atom, err
	}
	return atom, nil
}

func ErrorPacker(emsg string, env *terex.Environment) terex.Element {
	T().Errorf(emsg)
	env.Error(errors.New(emsg))
	return terex.Elem(terex.ErrorAtom(emsg))
}

func print(e terex.Element) {
	if e.IsNil() {
		T().Infof("nil")
		return
	}
	if e.IsAtom() {
		if e.AsAtom().Type() == terex.UserType {
			if G, ok := e.AsAtom().Data.(*lr.LRAnalysis); ok {
				level := T().GetTraceLevel()
				T().SetTraceLevel(tracing.LevelDebug)
				G.Grammar().Dump()
				T().SetTraceLevel(level)
				return
			}
			if t, ok := e.AsAtom().Data.(*parsetree); ok {
				err := showTree(t.tree)
				if err != nil {
					T().Errorf(err.Error())
				}
				return
			}
			if d, ok := e.AsAtom().Data.(Dict); ok {
				T().Infof("Dict {")
				for k, v := range d {
					T().Infof("\t%s => %v", k, terex.Elem(v))
				}
				T().Infof("}")
				return
			}
			if rew, ok := e.AsAtom().Data.([]*Rewriter); ok {
				T().Infof("Rewriter {")
				for _, r := range rew {
					T().Infof("\t%s:", r.String())
					for i, rule := range r.rules {
						T().Infof("\t   [%d] %s => %v", i, rule.pattern, rule.target)
					}
				}
				T().Infof("}")
				return
			}
		}
		// if e.AsAtom().Type() == terex.ConsType {
		// 	T().Infof(e.AsList().ListString())
		// 	return
		// }
		if e.AsAtom().Type() == terex.EnvironmentType {
			if env, ok := e.AsAtom().Data.(*terex.Environment); ok {
				T().Infof(env.Dump())
				return
			}
		}
		T().Infof(e.String())
		return
	}
	//T().Infof(e.AsList().ListString())
	T().Infof(e.String())
}

func showTree(tree *sppf.Forest) error {
	tmpfile, err := ioutil.TempFile("", "parsetree-*.dot")
	if err != nil {
		T().Errorf("Cannot create temporary file in TMP directory")
		return err
	}
	defer tmpfile.Close()
	sppf.ToGraphViz(tree, tmpfile)
	T().Infof("Exported parse tree to %s", tmpfile.Name())
	svg := strings.Replace(tmpfile.Name(), ".dot", ".svg", 1)
	cmd := exec.Command("dot", "-T", "svg", "-o", svg, tmpfile.Name())
	// T().Infof("Executing %v", cmd.String())
	err = cmd.Run()
	if err != nil {
		return err
	}
	T().Infof("Opening SVG output %s", svg)
	cmd = exec.Command("open", svg)
	err = cmd.Run()
	return err
}

// --- Rewriter --------------------------------------------------------------

func createASTRewritersFromDict(dict Dict) []*Rewriter {
	var arr []*Rewriter
	for k, v := range dict.getMap() {
		T().Errorf("ADD REW[%s] = %v", k, v)
		rew := NewRewriter(k)
		switch v.Type() {
		case terex.ConsType:
			rew.SetRewriters(v.Data.(*terex.GCons))
			arr = append(arr, rew)
		default:
			rew.SetRewriters(terex.Cons(terex.Atomize(terex.Cons(v, nil)), nil))
			arr = append(arr, rew)
		}
	}
	return arr
}

type rule struct {
	pattern terex.Element
	target  terex.Element
}

func (r rule) match(e terex.Element, env *terex.Environment) terex.Element {
	T().Errorf("MATCH? %s", r.pattern.AsList().ListString())
	if r.pattern.IsAtom() {
		a := r.pattern.AsAtom()
		T().Errorf("MATCH? Atom = %s", a.String())
		if a.Match(e.AsAtom(), env) {
			T().Errorf("MATCHED !")
			return terex.Element(r.target)
		}
	} else {
		l := r.pattern.AsList()
		T().Errorf("MATCH? List = %s", l.ListString())
		if l.Match(e.AsList(), env) {
			T().Errorf("MATCHED !")
			return terex.Element(r.target)
		}
	}
	return terex.Elem(nil)
}

type Rewriter struct {
	name  string
	Op    terex.Operator
	rules []rule
}

/*
	String() string                                         // printable name
	Rewrite(*terex.GCons, *terex.Environment) terex.Element // term rewriting
	Descend(sppf.RuleCtxt) bool                             // predicate wether to descend to children nodes
	Operator() terex.Operator                               // operator to place as sub-tree node
*/

func NewRewriter(symname string) *Rewriter {
	rew := &Rewriter{
		//name:  "rewrite(" + symname + ")",
		name:  symname,
		Op:    chameleonOp(symname),
		rules: make([]rule, 0, 3),
	}
	return rew
}

func (rew *Rewriter) SetRewriters(rlist *terex.GCons) {
	cons := rlist
	T().Errorf("REW LIST len = %d, list = %s", rlist.Length(), rlist.ListString())
	for cons != nil {
		if !cons.Car.IsNil() {
			atom := cons.Car
			if atom.IsNil() {
				continue
			}
			T().Errorf("CAR = %v", terex.Elem(atom).AsList().ListString())
			switch atom.Type() {
			case terex.ConsType:
				l := atom.Data.(*terex.GCons)
				pattern := terex.Elem(l.Car)
				target := terex.Elem(nil)
				if cons.Cdr != nil && !cons.Cdar().IsNil() {
					T().Errorf("CDR = %v", terex.Elem(cons.Cdar()).AsList().ListString())
					target = terex.Elem(cons.Cdar())
				}
				T().Errorf("ADD RULE %s => %v", pattern, target.String())
				rew.rules = append(rew.rules, rule{
					pattern: pattern,
					target:  target,
				})
			}
		}
		cons = cons.Cdr
	}
}

func (rew Rewriter) String() string {
	return rew.name
}

func (rew *Rewriter) Rewrite(l *terex.GCons, env *terex.Environment) terex.Element {
	T().Errorf("REWRITE")
	if l.Length() > 1 {
		m := rew.Match(terex.Elem(l.Cdr), env)
		if !m.IsNil() {
			return m
		}
	}
	return terex.Elem(l)
}

func (rew Rewriter) Descend(sppf.RuleCtxt) bool {
	return true
}
func (rew Rewriter) Operator() terex.Operator {
	return rew.Op
}

func (rew *Rewriter) Match(e terex.Element, env *terex.Environment) terex.Element {
	for i, r := range rew.rules {
		t := r.match(e, env) // match target
		if !t.AsAtom().IsNil() {
			T().Errorf("OUTPUT TARGET[%d] = %s", i, t.AsList().ListString())
			return t
		}
		//rules: make([]rule, 3, 3),
	}
	return terex.Elem(nil)
}

var _ termr.TermR = &Rewriter{}

// --- Chameleon operator ----------------------------------------------------

type chameleonOp string

func (chop chameleonOp) String() string {
	return string(chop)
}
func (chop chameleonOp) Call(e terex.Element, env *terex.Environment) terex.Element {
	opsym := env.FindSymbol(chop.String(), true)
	if opsym == nil {
		T().Errorf("Cannot find operator %v", chop)
		return e
	}
	operator, ok := opsym.AtomValue().(terex.Operator)
	if !ok {
		T().Errorf("Cannot call operator %s", chop)
		return e
	}
	return operator.Call(e, env)
}

var _ terex.Operator = chameleonOp("")
