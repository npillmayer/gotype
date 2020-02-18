package termr

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

import (
	"github.com/npillmayer/gotype/syntax/terex"
	"github.com/npillmayer/gotype/syntax/terex/fp"
)

// Rewriter is a function
//
//     list × env ↦ list
//
// i.e., a term rewriting function.
type Rewriter func(l *terex.GCons, env *terex.Environment) terex.Element

// RewriteRule is a type representing a rule for term rewriting.
// It contains a pattern and a rewriting-function. The pattern will be applied
// to nodes in an AST, and if it matches the rewriter will be called on the redex.
type RewriteRule struct {
	Pattern *terex.GCons
	Rewrite Rewriter
}

// env is an environment which holds an abstract syntax tree.
func Rewrite(env *terex.Environment, rule RewriteRule) error {
	return nil
}

// ---------------------------------------------------------------------------

// Anything is a pattern matching any s-expr.
func Anything() *terex.GCons {
	return terex.Cons(terex.Atomize(terex.ConsType), nil)
}

// AnySymbol is a pattern matching any single symbol or token.
func AnySymbol() *terex.GCons {
	return terex.Cons(terex.Atomize(terex.AnyType), nil)
}

type NodeEnv struct {
	node   *terex.GCons
	parent *terex.GCons
	env    *terex.Environment
}

func (seq fp.TreeSeq) Match(pattern *terex.GCons, env *terex.Environment) NodeEnv {
	node := seq.node.Node
	if pattern.Match(seq.node.Node, env) {
		return NodeEnv{
			node:   node,
			parent: seq.node.Parent(),
			env:    env,
		}
	}
	return NodeEnv{}
}

// EnvironmentForOperator creates an environment for an operator. The intention is for
// op to be the head of a TeREx list. The environment will have the operatore pre-stored
// as a symbol.
//
// If the parent environment is not given, it will be set to the global environment.
//
// Will return nil, if op is nil, a new environment otherwise.
func EnvironmentForOperator(op terex.Operator, parent *terex.Environment) *terex.Environment {
	if op == nil {
		return nil
	}
	if parent == nil {
		parent = terex.GlobalEnvironment
	}
	env := terex.NewEnvironment("#"+op.String(), parent)
	sym := env.Intern(op.String(), false)
	sym.Value = terex.Atomize(op)
	return env
}
