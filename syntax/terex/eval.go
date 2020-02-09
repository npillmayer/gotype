package terex

func (env *Environment) Eval(list *GCons) *GCons {
	r := env.eval(elem(list))
	T().Debugf("Eval => %s", r)
	return r.AsList()
}

func (env *Environment) eval(el Element) Element {
	T().Debugf("eval of %v", el)
	if el.IsAtom() {
		return el
	}
	list := el.AsList()
	if list == nil || list.Car == NilAtom {
		return elem(list)
	}
	op := list.Car
	if op.typ != OperatorType {
		return elem(list)
	}
	operator := op.Data.(Operator)
	args := list.Cdr.Map(env.eval)
	return operator.Call(elem(args))
}
