package terex

import "fmt"

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

func (env *Environment) Eval(list *GCons) *GCons {
	r := Eval(Elem(list), env)
	T().Debugf("Eval => %s", r.String())
	return r.AsList()
}

func Eval(el Element, env *Environment) Element {
	T().Debugf("eval of %v", el)
	if el.IsAtom() {
		if el.AsAtom().Type() == ConsType {
			sublist := el.AsAtom().Data.(*GCons)
			mapped := evalList(sublist, env)
			return mapped
		}
		return evalAtom(el.AsAtom(), env)
	}
	list := el.AsList()
	l := evalList(list, env)
	return l
}

func evalList(list *GCons, env *Environment) Element {
	if list == nil || list.Car == NilAtom {
		return Elem(list)
	}
	carAtom, err := resolve(list.Car, env, true)
	if err != nil {
		return Elem(ErrorAtom(err.Error()))
	}
	if carAtom.Type() != OperatorType { // Resolver gave us this
		verblist := list.Map(Eval, env) // ⇒ accept it
		return Elem(verblist)           // and return non-operated list
	}
	T().Debugf("-------- op=%s -----------", carAtom.String())
	operator := carAtom.Data.(Operator)
	T().Debugf("OP = %s", operator)
	//args := list.Cdr.Map(Eval, env)
	args := Elem(list.Cdr)
	T().Debugf("-------- call ------------")
	ev := operator.Call(args, env)
	T().Debugf("  eval result = %v", ev)
	T().Debugf("--------------------------")
	return ev
}

func evalAtom(atom Atom, env *Environment) Element {
	resolved, _ := resolve(atom, env, false)
	T().Debugf("resolved -> %v", resolved)
	return Elem(resolved)
}

func resolve(atom Atom, env *Environment, asOp bool) (Atom, error) {
	if env.Resolver == nil {
		return DefaultSymbolResolver{}.Resolve(atom, env, asOp)
	}
	return env.Resolver.Resolve(atom, env, asOp)
}

type DefaultSymbolResolver struct {
	// TODO options
}

func (dsr DefaultSymbolResolver) Resolve(atom Atom, env *Environment, asOp bool) (Atom, error) {
	if atom.Type() == OperatorType {
		return atom, nil // shortcut, not resolved in env
	}
	if atom.Type() == VarType {
		atomSym := atom.Data.(*Symbol)
		sym := env.FindSymbol(atomSym.Name, true)
		if sym == nil {
			T().Errorf("Unable to resolve symbol '%s' in environment", atomSym.Name)
			err := fmt.Errorf("Unable to resolve symbol '%s' in environment", atomSym.Name)
			env.Error(err)
			return atom, err
		}
		valueAtom := sym.Value
		if asOp && valueAtom.Type() != OperatorType {
			env.lastError = fmt.Errorf("Symbol '%s' cannot be resolved as operator", sym.Name)
			T().Errorf("Symbol '%s' cannot be resolved as operator", sym.Name)
			err := fmt.Errorf("Symbol '%s' cannot be resolved as operator", sym.Name)
			env.Error(err)
			return NilAtom, err
		}
		return valueAtom, nil
	}
	if asOp { // atom is neither a symbol nor an operator, but operator expected
		env.lastError = fmt.Errorf("Atom '%s' cannot be cast to operator ", atom)
		T().Errorf("Atom '%s' cannot be cast to operator ", atom)
		err := fmt.Errorf("Atom '%s' cannot be cast to operator ", atom)
		env.Error(err)
		return NilAtom, err
	}
	return atom, nil
}

var _ SymbolResolver = &DefaultSymbolResolver{}

// --- Quote -----------------------------------------------------------------

// Quote traverses an s-expr and returns it as pure list/tree data.
// It gets rid of #list:op and #quote:op nodes.
//
// If the environment contains a symbol's value, quoting will replace the symbol
// by its value. For example, if the s-expr contains a symbol 'str' with a value
// of "this is a string", the resulting data structure will contain the string,
// not the name of the symbol. If you do not have use for this kind of substitution,
// simply call Quote(…) for the global environment.
//
/*
func (env *Environment) Quote(list *GCons) *GCons {
	r := env.quote(Elem(list))
	T().Debugf("Quote => %s", r)
	return r.AsList()
}

func (env *Environment) quote(el Element) Element {
	if el.IsAtom() {
		//T().Errorf("quote of atom %v, type=%s", el, el.AsAtom().Type().String())
		if el.AsAtom().Type() == ConsType {
			//T().Infof("---- atom/list --------------")
			//T().Errorf("QUOTE el.Cdr=%s", el.AsAtom().Data.(*GCons).ListString())
			sublist := el.AsAtom().Data.(*GCons)
			mapped := env.quoteList(sublist).AsList()
			//T().Infof("MAPPED = %s", mapped.ListString())
			//T().Infof("-----------------------------")
			return Elem(mapped)
			//return Elem(Cons(Atomize(mapped), nil))
		}
		return el
	}
	//T().Infof("==== as list ================")
	list := el.AsList()
	l := env.quoteList(list)
	//T().Infof("=============================")
	return l
	//return env.quoteList(list)
	// T().Errorf("quote of list %v", list.ListString())
	// if list == nil || list.Car == NilAtom {
	// 	return Elem(list)
	// }
	// op := list.Car
	// if op.typ != OperatorType {
	// 	verblist := list.Map(env.quote)
	// 	return Elem(verblist)
	// }
	// T().Errorf("quote-OP = %s", op)
	// operator := op.Data.(Operator)
	// args := list.Cdr.Map(env.quote)
	// return operator.Quote(Elem(args), env)
}

func (env *Environment) quoteAtom(atom Atom) Element {
	return Elem(atom) // TODO
}

func (env *Environment) quoteList(list *GCons) Element {
	//
	//T().Errorf("quote of list %v", list.ListString())
	if list == nil || list.Car == NilAtom {
		return Elem(list)
	}
	op := list.Car
	if op.typ != OperatorType {
		T().Infof("------- VerbList -----------------------------")
		T().Infof("   > Map(quote(...))")
		verblist := list.Map(env.quote)
		T().Infof("----------------------------------------------")
		return Elem(verblist)
	}
	//T().Errorf("quote-OP = %s", op)
	operator := op.Data.(Operator)
	args := list.Cdr.Map(env.quote)
	//args := list.Cdr
	T().Infof("-------- Op = %s -----------------------------", operator.String())
	T().Infof("     args=%s", args.ListString())
	T().Infof("   > quote(args...)")
	quoted := operator.Quote(Elem(args), env)
	T().Infof("     quoted=%s", quoted.String())
	T().Infof("----------------------------------------------")
	return quoted
}
*/
