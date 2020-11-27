package terexlang

import (
	"fmt"
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
		args = args.Map(terex.Eval, env)
		if args.Length() != 2 {
			return terex.Elem(terex.ErrorAtom("Wrong number of arguments for def"))
		}
		if args.Car.Type() != terex.VarType {
			return terex.Elem(terex.ErrorAtom("Missing symbol for def"))
		}
		sym := env.Intern(args.Car.Data.(*terex.Symbol).Name, true)
		sym.Value = args.Cdar()
		return terex.Elem(terex.Atomize(sym))
	})
	env.Defn("print", func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		args = args.Map(terex.Eval, env)
		args.Map(func(e terex.Element, env *terex.Environment) terex.Element {
			T().Infof(e.String())
			return terex.Elem(nil)
		}, env)
		last := args.Last()
		return terex.Elem(terex.Atomize(last))
	})
}

func makeParserOps(env *terex.Environment) {
	env.Defn("parse", func(e terex.Element, env *terex.Environment) terex.Element {
		// (parse G "input") => tree
		args := e.AsList()
		args = args.Map(terex.Eval, env)
		if args.Length() != 2 {
			return terex.Elem(terex.ErrorAtom("Wrong number of arguments for parse"))
		}
		if args.Car.Type() != terex.UserType || args.Cdar().Type() != terex.StringType {
			return terex.Elem(terex.ErrorAtom("Wrong argument type for parse"))
		}
		if args.Cdar().Data == nil {
			return terex.Elem(terex.ErrorAtom("Cannot parse nil-input"))
		}
		input := args.Cdar().Data.(string)
		g := args.Car.Data.(*lr.LRAnalysis)
		tree, retr, err := parseAny(g, input)
		if err != nil {
			return terex.Elem(terex.ErrorAtom(err.Error()))
		}
		result := &parsetree{
			tree: tree,
			retr: retr,
		}
		return terex.Elem(terex.Atomize(result))
	})
	ast := makeOp(terex.UserType, 1, func(args *terex.GCons, env *terex.Environment) terex.Element {
		// (ast T) => AST
		if tree, ok := args.Car.Data.(*parsetree); ok {
			ast, env, err := AST(tree.tree, tree.retr)
			//t.Logf("\n\n" + debugString(terex.Elem(ast.Car)))
			T().Infof("\n\n" + terex.Elem(ast).String())
			if err != nil {
				T().Errorf(err.Error())
				return terex.Elem(terex.ErrorAtom(err.Error()))
			}
			q, err := QuoteAST(terex.Elem(ast.Car), env)
			if err != nil {
				T().Errorf(err.Error())
				return terex.Elem(terex.ErrorAtom(err.Error()))
			}
			return q
		}
		return terex.Elem(terex.ErrorAtom("Argument for ast is not a parse tree"))
	})
	env.Defn("ast", ast)
}

type parsetree struct {
	tree *sppf.Forest
	retr termr.TokenRetriever
}

// --- Helpers ---------------------------------------------------------------

func makeOp(t terex.AtomType, a int, op func(*terex.GCons, *terex.Environment) terex.Element) func(
	e terex.Element, env *terex.Environment) terex.Element {
	arity := a
	typ := t
	call := func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		args = args.Map(terex.Eval, env)
		if arity > 0 && args.Length() != arity {
			return terex.Elem(terex.ErrorAtom("Wrong number of arguments"))
		} else if arity < 0 && args.Length() < -arity {
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
