package terexlang

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/terex"
)

// LoadStandardLanguage returns an environment with all the standard
// symbols and operators loaded.
//
func LoadStandardLanguage() *terex.Environment {
	env := terex.NewEnvironment("stdlang", nil)
	makeArithmOps(env)
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

// --- Helpers ---------------------------------------------------------------

func makeOp(t terex.AtomType, a int, op func(*terex.GCons, *terex.Environment) terex.Element) func(
	e terex.Element, env *terex.Environment) terex.Element {
	arity := a
	typ := t
	call := func(e terex.Element, env *terex.Environment) terex.Element {
		args := e.AsList()
		args = args.Map(terex.Eval, env)
		if arity > 0 && args.Length() != arity {
			return terex.Elem(terex.ErrorAtom("wrong number of arguments"))
		} else if arity < 0 && args.Length() < -arity {
			return terex.Elem(terex.ErrorAtom("wrong number of arguments"))
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
