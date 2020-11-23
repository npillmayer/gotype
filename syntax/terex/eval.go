package terex

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
	r := eval(Elem(list), env)
	T().Debugf("Eval => %s", r.String())
	return r.AsList()
}

func eval(el Element, env *Environment) Element {
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
	op := list.Car
	if op.typ != OperatorType {
		verblist := list.Map(eval, env)
		return Elem(verblist)
	}
	T().Debugf("-------- op=%s -----------", op.String())
	operator := op.Data.(Operator)
	args := list.Cdr.Map(eval, env)
	T().Debugf("-------- call ------------")
	ev := operator.Call(Elem(args), env)
	T().Debugf("  eval result = %v", ev)
	T().Debugf("--------------------------")
	return ev
}

func evalAtom(atom Atom, env *Environment) Element {
	if atom.Type() == TokenType {
		if sym := env.FindSymbol("!TokenEvaluator", true); sym != nil {
			if sym.Value.typ == OperatorType { // then use it
				teval := sym.Value.Data.(Operator)
				result := teval.Call(Elem(atom), env)
				return result
			}
		}
	}
	return Elem(atom)
}

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
