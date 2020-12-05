/*
Package varparse implements a parser to get variable references from strings.

The implementation is tightly coupled to the ANTLR V4 parser generator.


BSD License

Copyright (c) 2017, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer nor the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/
package varparse

import (
	"fmt"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	"github.com/npillmayer/gotype/syntax/variables/grammar"
	"github.com/npillmayer/schuko/tracing"
	dec "github.com/shopspring/decimal"
)

// We're tracing to the InterpreterTracer
func T() tracing.Trace {
	return tracing.InterpreterTracer
}

// === Variable Parser =======================================================

/*
We use a small ANTLR V4 sub-grammar for parsing variable references.
We'll attach a listener to ANTLR's AST walker.
*/
type varParseListener struct {
	*grammar.BasePMMPVarListener // build on top of ANTLR's base 'class'
	scopeTree                    *runtime.ScopeTree
	def                          *variables.PMMPVarDecl
	ref                          *variables.PMMPVarRef
	suffixes                     *sll.List
	fullname                     string // full variable name
}

// Construct a new variable parse listener.
func newVarParseListener() *varParseListener {
	pl := &varParseListener{} // no need to initialize base class
	return pl
}

// Parse a variable from a string. This function will try to find an existing
// declaration and extend it accordingly. If no declaration can be found, this
// function will construct a new one from the variable reference.
//
// To find a variable in memory we first construct the "canonical" form.
// Variables live in symbol tables and are resolved by name, so the exact
// spelling of the variable's name is important. The canonical form uses
// brackets for subscripts and dots for suffixes.
//
// Examples:
//
//    x ⟹ x
//    x3 ⟹ x[3]
//    z4r ⟹ z[4].r
//    hello.world ⟹ hello.world
//
func parseVariableFromString(vstr string, err antlr.ErrorListener) antlr.RuleContext {
	// We let ANTLR to the heavy lifting. This may change in the future,
	// as it would be fairly straightforward to implement this by hand.
	input := antlr.NewInputStream(vstr + "@")
	T().Debugf("parsing variable ref = %s", vstr)
	varlexer := grammar.NewPMMPVarLexer(input)
	stream := antlr.NewCommonTokenStream(varlexer, 0)
	// TODO: make the parser re-usable....!
	// TODO: We can re-use the parser, but not the lexer (see ANTLR book, chapter 10).
	varParser := grammar.NewPMMPVarParser(stream)
	if err == nil {
		err = antlr.NewDiagnosticErrorListener(true)
	} else {
		T().Debugf("setting error listener")
	}
	varParser.RemoveErrorListeners()
	varParser.AddErrorListener(err)
	varlexer.RemoveErrorListeners()
	varlexer.AddErrorListener(err)
	varParser.BuildParseTrees = true
	tree := varParser.Variable()
	sexpr := antlr.TreesStringTree(tree, nil, varParser)
	T().Debugf("### VAR = %s", sexpr)
	return tree
}

// After having ANTLR create us a parse tree for a variable identifier, we
// will construct a variable reference from the parse tree. This var ref
// has an initial serial which is unique. This may not be what you want:
// usually you will try to find an existing incarnation (with a lower serial)
// in the memory (see method FindVariableReferenceInMemory).
//
// We walk the ANTLR parse tree using a listener (varParseListener).
//
func getVarRefFromVarSyntax(vtree antlr.RuleContext, scopes *runtime.ScopeTree) (
	*variables.PMMPVarRef, string) {
	//
	listener := newVarParseListener()
	listener.scopeTree = scopes                        // where variables declarations have to be found
	listener.suffixes = sll.New()                      // list to collect suffixes
	antlr.ParseTreeWalkerDefault.Walk(listener, vtree) // fills listener.ref
	return listener.ref, listener.def.GetFullName()
}

// Helper struct to collect suffixes
type varsuffix struct {
	text   string
	number bool
}

/*
Listener callback, receiving a complete variable reference.

   variable : tag (suffix | subscript)* MARKER

A variable has been referenced. We will have to find the declaration of
this variable and push a variable reference onto the expression stack.
A variable reference looks something like this: "x2a.b" or "y[3.14]".

Some complications may arise:

- No declaration for the variable's tag can be found: we will create
a declaration for the tag in global scope with type numeric.

- The declaration is incomplete, i.e. the tag is declared, but not the
suffix(es). We will extend the declaration appropriately.

- We will have to create a subscript vector for the var reference.
we'll collect them (together with suffixes) in a list.

Example:
Parser reads "x2a", thus

   tag="x" + subscript="2" + suffix="a"

We create (if not yet known)

   vl.def="x[].a"

and

   vl.ref="x[2].a"

The MARKER will be ignored.
*/
func (vl *varParseListener) ExitVariable(ctx *grammar.VariableContext) {
	tag := ctx.Tag().GetText()
	T().P("tag", tag).Debugf("looking for declaration for tag")
	sym, scope := vl.scopeTree.Current().ResolveSymbol(tag)
	if sym != nil {
		vl.def = sym.(*variables.PMMPVarDecl) // scopes are assumed to create these
		T().P("decl", vl.def.GetFullName()).Debugf("found %v in scope %s", vl.def, scope.GetName())
	} else { // variable declaration for tag not found => create it
		sym, _ = vl.scopeTree.Globals().DefineSymbol(tag)
		vl.def = sym.(*variables.PMMPVarDecl)      // scopes are assumed to create these
		vl.def.SetType(int(variables.NumericType)) // un-declared variables default to type numeric
		T().P("decl", vl.def.GetName()).Debugf("created %v in global scope", vl.def)
	} // now def declaration of <tag> is in vl.def
	// produce declarations for suffixes, if necessary
	it := vl.suffixes.Iterator()
	subscrCount := 0
	for it.Next() {
		i, vs := it.Index(), it.Value().(varsuffix)
		T().P("decl", vl.def.GetFullName()).Debugf("appending suffix #%d: %s", i, vs)
		if vs.number { // subscript
			vl.def = variables.CreatePMMPVarDecl("<array>", variables.ComplexArray, vl.def)
			subscrCount += 1
		} else { // tag suffix
			vl.def = variables.CreatePMMPVarDecl(vs.text, variables.ComplexSuffix, vl.def)
		}
	}
	T().P("decl", vl.def.GetFullName()).Debugf("full declared type: %v", vl.def)
	// now create variable ref and push onto expression stack
	var subscripts []dec.Decimal = make([]dec.Decimal, subscrCount, subscrCount+1)
	it = vl.suffixes.Iterator()
	for it.Next() { // extract subscripts -> array
		_, vs2 := it.Index(), it.Value().(varsuffix)
		if vs2.number { // subscript
			d, _ := dec.NewFromString(vs2.text)
			subscripts = append(subscripts, d)
		}
	}
	vl.ref = variables.CreatePMMPVarRef(vl.def, nil, subscripts)
	T().P("var", vl.ref.GetName()).Debugf("var ref %v", vl.ref)
}

// Variable parsing: Collect a suffix.
func (vl *varParseListener) ExitSuffix(ctx *grammar.SuffixContext) {
	tag := ctx.TAG().GetText()
	T().Debugf("suffix tag: %s", tag)
	vl.suffixes.Add(varsuffix{tag, false})
}

// Variable parsing: Collect a numeric subscript.
func (vl *varParseListener) ExitSubscript(ctx *grammar.SubscriptContext) {
	d := ctx.DECIMAL().GetText()
	T().Debugf("subscript: %s", d)
	vl.suffixes.Add(varsuffix{d, true})
}

// ----------------------------------------------------------------------

// VariableResolver will find a variable's declaration in a scope tree,
// if present. It will then construct a valid variable reference for
// that declaration.
type VariableResolver interface {
	VariableName() string // full name of variable
	VariableReference(*runtime.ScopeTree) *variables.PMMPVarRef
}

// ParseVariableName will parse a string as a variable's name and
// return a VariableResolver.
func ParseVariableName(v string) (VariableResolver, error) {
	el := &varErrorListener{}
	vtree := parseVariableFromString(v, el)
	if el.err != nil {
		return nil, el.err
	}
	return &varResolver{vtree, "", nil}, el.err
}

type varResolver struct {
	ctx      antlr.RuleContext
	fullname string // variable full name
	varref   *variables.PMMPVarRef
}

func (r *varResolver) VariableReference(sc *runtime.ScopeTree) *variables.PMMPVarRef {
	if r == nil || r.ctx == nil || sc == nil {
		return nil
	}
	if r.varref == nil {
		r.varref, r.fullname = getVarRefFromVarSyntax(r.ctx, sc)
	}
	return r.varref
}

func (r *varResolver) VariableName() string {
	if r.varref == nil {
		r.VariableReference(nil) // provisional call will set fullname
	}
	return r.fullname
}

type varErrorListener struct {
	*antlr.DefaultErrorListener // use default as base class
	err                         error
}

func (el *varErrorListener) SyntaxError(r antlr.Recognizer, sym interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	//
	el.err = fmt.Errorf("[%s|%s] %.44s", strconv.Itoa(line), strconv.Itoa(column), msg)
}
