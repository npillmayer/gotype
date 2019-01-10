package pmmpost

/*
----------------------------------------------------------------------

BSD License
Copyright (c) 2017, Norbert Pillmayer <norbert@pillmayer.com>

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

----------------------------------------------------------------------

*/

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/backend/gfx"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/corelang"
	"github.com/npillmayer/gotype/syntax/pmmpost/listener"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
)

// === Interpreter ===========================================================

/*
We use AST-driven interpretation to execute the input program. Input
is more or less a list of statements and function definitions.
We will annotate the AST with scope-information, holding the symbols of
dynamic scopes. Scopes stem from either:

- function definitions (macros in MetaFont: def and vardef)

- compound statements, i.e. groups (begingroup ... endgroup)

The interpreter relies on the scopes and definitions constructed earlier.
It manages a memory frame stack to track the calling sequence of functions
and groups.

So the overall picture looks like this:

1. ANTLR V4 constructs an AST for us.

2. We use a listener to walk the AST and execute the statements.

Metafont, and therefore PMMPost, is a dynamically scoped language. This means,
functions can access local variables from calling functions or groups.
Nevertheless we will find the definition of all variables (which are explicitly
defined) in a scope definition. This is mainly for type checking reasons and
due to the complex structure of MetaFont variable identifiers.

PMMPostInterpreter

This is an umbrella object to hold together the various tools needed to
execute steps 1 to 3 from above. It will orchestrate and instrument the
tools and execute them in the correct order. Also, this object will hold
and respect parameters we pass to the interpreter, so we can alter the
behaviour in certain aspects.
*/
type PMMPostInterpreter struct {
	ASTListener *listener.PMMPostParseListener // parse / AST listener
	runtime     *runtime.Runtime               // runtime environment
	scripting   *corelang.Scripting            // Lua scripting subsystem
}

// Create a new Interpreter for the "Poor Man's MetaPost". This is the top-level
// object for this package.
//
// The interpreter manages an AST listener, a scope tree, a memory frame stack,
// and stacks for numeric/pair and path expressions.
//
// Loads builtin symbols (variables and types) if argument is true.
//
func NewPMMPostInterpreter(loadBuiltins bool) *PMMPostInterpreter {
	T = tracing.InterpreterTracer
	intp := &PMMPostInterpreter{}
	intp.runtime = runtime.NewRuntimeEnvironment(variables.NewPMMPVarDecl)
	intp.scripting = corelang.NewScripting(intp.runtime) // scripting subsystem
	if loadBuiltins {
		corelang.LoadBuiltinSymbols(intp.runtime, intp.scripting) // load syms into global scope
		intp.loadAdditionalBuiltinSymbols(intp.runtime, intp.scripting)
	}
	intp.ASTListener = listener.NewParseListener(intp.runtime, intp.scripting) // listener for ANTLR
	return intp
}

// Set an output routine for drawing. Default is nil.
func (intp *PMMPostInterpreter) SetOutputRoutine(o gfx.OutputRoutine) {
	if intp.ASTListener != nil {
		intp.ASTListener.SetOutputRoutine(o)
	}
}

// Parse and interpret a statement list.
//
func (intp *PMMPostInterpreter) ParseStatements(input []byte) {
	// func (intp *PMMPostInterpreter) ParseStatements(input antlr.CharStream) {
	inputStream := antlr.NewInputStream(string(input))
	intp.ASTListener.ParseStatements(inputStream)
}

// Load additionl builtins for this language (added to the core symbols)
func (intp *PMMPostInterpreter) loadAdditionalBuiltinSymbols(rt *runtime.Runtime,
	scripting *corelang.Scripting) {
	//
	//scripting.RegisterHook("TODO: z", ping)
}

// Value returns the current value of a variable.
// Return values are: canonical name of the variable & value of the
// variable as string.
func (intp *PMMPostInterpreter) Value(variable string) (string, string) {
	variable = strings.TrimSpace(variable)
	r, _ := utf8.DecodeRuneInString(variable)
	if !unicode.IsLetter(r) {
		return variable, "<illegal variable name>"
	}
	vtree := variables.ParseVariableFromString(variable, nil)
	if vtree == nil {
		return variable, "<illegal variable name>"
	}
	vref := variables.GetVarRefFromVarSyntax(vtree, intp.runtime.ScopeTree)
	if vref == nil {
		return variable, "<undefined>"
	}
	v := corelang.MakeCanonicalAndResolve(intp.runtime, variable, false)
	if v == nil {
		t := variables.TypeString(vref.GetType())
		return vref.GetFullName(), fmt.Sprintf("<unset %s>", t)
	}
	return v.GetFullName(), v.ValueString()
}
