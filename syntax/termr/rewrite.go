package termr

import "github.com/npillmayer/gotype/syntax/terex"

// Rewriter is a function
//
//     list × env ↦ list
//
// i.e., a term rewriting function.
type Rewriter func(l *terex.GCons, env *terex.Environment) *terex.GCons

// RewriteRule is a type representing a rule for term rewriting.
// It contains a pattern and a rewriting-function. The pattern will be applied
// to nodes in an AST, and if it matches the rewriter will be called on the redex.
type RewriteRule struct {
	Pattern *terex.GCons
	Rewrite Rewriter
}

// Anything is a pattern matching any s-expr.
func Anything() *terex.GCons {
	return terex.Cons(terex.Atomize(terex.ConsType), nil)
}

// AnySymbol is a pattern matching any single symbol or token.
func AnySymbol() *terex.GCons {
	return terex.Cons(terex.Atomize(terex.AnyType), nil)
}
