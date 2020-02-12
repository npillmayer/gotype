package terex

func (env *Environment) Eval(list *GCons) *GCons {
	r := env.eval(Elem(list))
	T().Debugf("Eval => %s", r.String())
	return r.AsList()
}

func (env *Environment) eval(el Element) Element {
	T().Debugf("eval of %v", el)
	if el.IsAtom() {
		if el.AsAtom().Type() == ConsType {
			sublist := el.AsAtom().Data.(*GCons)
			mapped := env.evalList(sublist)
			return mapped
		}
		return env.evalAtom(el.AsAtom())
	}
	list := el.AsList()
	l := env.evalList(list)
	return l
}

func (env *Environment) evalList(list *GCons) Element {
	if list == nil || list.Car == NilAtom {
		return Elem(list)
	}
	op := list.Car
	if op.typ != OperatorType {
		verblist := list.Map(env.eval)
		return Elem(verblist)
	}
	T().Infof("-------- op=%s -----------", op.String())
	operator := op.Data.(Operator)
	args := list.Cdr.Map(env.eval)
	T().Infof("-------- call ------------")
	ev := operator.Call(Elem(args), env)
	T().Infof("  eval result = %v", ev)
	T().Infof("--------------------------")
	return ev
}

func (env *Environment) evalAtom(atom Atom) Element {
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
