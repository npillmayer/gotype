package termr

type RewriteRule struct {
	Pattern *GCons
	Rewrite func(*GCons, *Environment) *GCons
}

/*
Match an s-expr to a pattern.

From https://hanshuebner.github.io/lmman/fd-con.xml:

list-match-p object pattern

object is evaluated and matched against pattern; the value is t if it matches, nil otherwise.
pattern is made with backquotes (Aids for Defining Macros); whereas normally a backquote
expression says how to construct list structure out of constant and variable parts, in
this context it says how to match list structure against constants and variables. Constant
parts of the backquote expression must match exactly; variables preceded by commas can
match anything but set the variable to what was matched. (Some of the variables may be
set even if there is no match.) If a variable appears more than once, it must match
the same thing (equal list structures) each time. ,ignore can be used to match anything
and ignore it. For example, `(x (,y) . ,z) is a pattern that matches a list of length
at least two whose first element is x and whose second element is a list of length one;
if a list matches, the caadr of the list is stored into the value of y and the cddr of
the list is stored into z. Variables set during the matching remain set after the
list-match-p returns; in effect, list-match-p expands into code which can setq the
variables. If the match fails, some or all of the variables may already have been set.

Example:

    (list-match-p foo
              `((a ,x) ,ignore . ,c))

is t if foo's value is a list of two or more elements, the first of which is a list
of two elements; and in that case it sets x to (cadar foo) and c to (cddr foo).
*/
func (l *GCons) Match(other *GCons, env *Environment) bool {
	if l == nil {
		return other == nil
	}
	if other == nil {
		return false
	}
	if !matchCar(l.car, other.car, env) {
		return false
	}
	return l.cdr.Match(other.cdr, env)
}

func matchCar(car Node, otherNode Node, env *Environment) bool {
	if car == nullNode {
		return otherNode == nullNode
	}
	if car.Type() == SymbolType {
		return bindSymbol(car, otherNode, env)
	}
	if car.Type() == ConsType {
		if otherNode.Type() != ConsType {
			return false
		}
		return car.child.Match(otherNode.child, env)
	}
	return matchAtom(car.atom, otherNode.atom)
}

func matchAtom(atom Atom, otherAtom Atom) bool {
	if atom == NullAtom {
		return otherAtom == NullAtom
	}
	if otherAtom == NullAtom {
		return false
	}
	if atom.typ != otherAtom.typ {
		return false
	}
	return atom.Data == otherAtom.Data
}

func bindSymbol(symcar Node, otherNode Node, env *Environment) bool {
	sym := symcar.atom.Data.(*Symbol)
	T().Debugf("binding symbol %s to %s", sym.String(), otherNode.String())
	if sym.value == nullNode {
		//symcar.atom.Data = otherNode
		sym.value = otherNode // bind it
		T().Debugf("bound symbol %s", sym.String())
		return true
	}
	// valuecar := symcar.atom.Data.(Node)
	// return matchCar(valuecar, otherNode, env)
	return true // TODO match sym.value to other
}
