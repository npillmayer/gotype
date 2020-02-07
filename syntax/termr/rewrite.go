package termr

// Rewriter is a function
//
//     list × env ↦ list
//
// i.e., a term rewriting function.
type Rewriter func(l *GCons, env *Environment) *GCons

// RewriteRule is a type representing a rule for term rewriting.
// It contains a pattern and a rewriting-function. The pattern will be applied
// to nodes in an AST, and if it matches the rewriter will be called on the redex.
type RewriteRule struct {
	Pattern *GCons
	Rewrite Rewriter
}

// --- Matching --------------------------------------------------------------

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

List l is the pattern, other is the argument to be matched against the pattern.
*/
func (l *GCons) Match(other *GCons, env *Environment) bool {
	T().Debugf("Match: %s vs %s", l, other)
	if l != nil && l.car.Type() == AnyList {
		return true
	}
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
	T().Debugf("Match Car: %s vs %s", car, otherNode)
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
	T().Debugf("Match Atom: %v vs %v", atom, otherAtom)
	if atom == NullAtom {
		return otherAtom == NullAtom
	}
	if otherAtom == NullAtom {
		return false
	}
	typeMatches, doMatchData := typeMatch(atom.typ, otherAtom.typ)
	if !typeMatches {
		return false
	}
	if doMatchData {
		return dataMatch(atom.Data, otherAtom.Data, atom.typ)
	}
	return true
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

// typeMatch returns (typesAreMatching, mustMatchValue)
func typeMatch(t1 AtomType, t2 AtomType) (bool, bool) {
	if t1 == t2 {
		return true, true
	}
	if t1 == AnyType {
		return true, false
	}
	return false, true
}

func dataMatch(d1 interface{}, d2 interface{}, t AtomType) bool {
	if t == OperatorType && d1 != nil {
		if op, ok := d1.(*sExprOp); ok {
			if op.Name() == AnyOp.name {
				return true
			}
		}
	}
	if t == TokenType && d2 != nil {
		tok1, _ := d1.(*Token)
		if tok2, ok := d2.(*Token); ok {
			if tok1.value == tok2.value { // only tokval must match
				return true
			}
		}
	}
	T().Errorf("dataMatch()")
	return d1 == d2
}
