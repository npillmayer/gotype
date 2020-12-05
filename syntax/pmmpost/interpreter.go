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
	"github.com/npillmayer/gotype/core/arithmetic"
	"github.com/npillmayer/gotype/core/path"
	"github.com/npillmayer/gotype/syntax/corelang"
	"github.com/npillmayer/gotype/syntax/pmmpost/listener"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	"github.com/npillmayer/gotype/syntax/variables/varparse"
	"github.com/npillmayer/schuko/tracing"
	"github.com/shopspring/decimal"
)

// === Interpreter ===========================================================

// This is an umbrella object to hold together the various tools needed to
// parse and execute PMMPost statements. It will orchestrate and instrument the
// tools and execute them in the correct order. Also, this object will hold
// and respect parameters we pass to the interpreter, so we can alter the
// behaviour in certain aspects.
//
type PMMPostInterpreter struct {
	listener  *listener.PMMPostParseListener // parse-listener
	runtime   *runtime.Runtime               // runtime environment
	scripting *corelang.Scripting            // Lua scripting subsystem
}

// Create a new Interpreter for the "Poor Man's MetaPost". This is the top-level
// object for this package.
//
// The interpreter manages an AST listener, a scope tree, a memory frame stack,
// and stacks for numeric/pair and path expressions.
//
// Loads builtin symbols (variables and types) if argument is true.
//
func NewPMMPostInterpreter(loadBuiltins bool, callback func(*gfx.Picture)) *PMMPostInterpreter {
	T = tracing.InterpreterTracer
	intp := &PMMPostInterpreter{}
	intp.runtime = runtime.NewRuntimeEnvironment(variables.NewPMMPVarDecl)
	intp.scripting = corelang.NewScripting(intp.runtime) // scripting subsystem
	if loadBuiltins {
		corelang.LoadBuiltinSymbols(intp.runtime, intp.scripting) // load syms into global scope
		intp.loadAdditionalBuiltinSymbols(intp.runtime, intp.scripting)
	}
	intp.listener = listener.NewParseListener(intp.runtime, intp.scripting,
		callback) // listener for ANTLR
	return intp
}

// Parse and interpret a statement list.
//
func (intp *PMMPostInterpreter) ParseStatements(input []byte) []error {
	inputStream := antlr.NewInputStream(string(input))
	errs := intp.listener.ParseStatements(inputStream)
	//T.Infof("ParseStmt: ERR = %v", errs)
	return errs
}

// Load additionl builtins for this language (added to the core symbols)
func (intp *PMMPostInterpreter) loadAdditionalBuiltinSymbols(rt *runtime.Runtime,
	scripting *corelang.Scripting) {
	//
	//scripting.RegisterHook("TODO: z", ping)
}

// ----------------------------------------------------------------------

// VarValue is a wrapper type for variable values. It hides the complexities
// of the various types of variables from clients.
//
// This is part of an API for the interpreter. The PMMPost interpreter can
// be used as a standalone batch drawing CLI, but its foremost purpose
// is to be integrated in applications.
type VarValue struct {
	VariableFullName string      // full name of the variable this value is from
	val              interface{} // variable's value
}

// IsSet checks if a variable's value is set.
func (vval VarValue) IsSet() bool {
	return vval.val != nil
}

// String returns a variable's value as a string.
func (vval VarValue) String() string {
	if vval.val == nil {
		return ""
	}
	return fmt.Sprintf("%v", vval.val)
}

// Type returns the type of a variable's value.
func (vval VarValue) Type() variables.VariableType {
	if vval.IsSet() {
		switch vval.val.(type) {
		case arithmetic.Pair:
			return variables.PairType
		case *path.Path:
			return variables.PathType
		case decimal.Decimal:
			return variables.NumericType
		default:
			return variables.Undefined
		}
	}
	return variables.Undefined
}

func (vval VarValue) AsNumeric() decimal.Decimal {
	return vval.val.(decimal.Decimal)
}

func (vval VarValue) AsPair() arithmetic.Pair {
	return vval.val.(arithmetic.Pair)
}

func (vval VarValue) AsPath() *path.Path {
	return vval.val.(*path.Path)
}

// Value returns the current value of a variable.
// Return values are: canonical name of the variable & value of the
// variable as string.
func (intp *PMMPostInterpreter) ValueOf(variable string) (VarValue, error) {
	variable = strings.TrimSpace(variable)
	r, _ := utf8.DecodeRuneInString(variable)
	if !unicode.IsLetter(r) {
		return VarValue{
			VariableFullName: variable,
			val:              nil,
		}, fmt.Errorf("Illegal variable name: %s", variable)
	}
	resolve, err := varparse.ParseVariableName(variable)
	if err != nil {
		return VarValue{
			VariableFullName: variable,
			val:              nil,
		}, fmt.Errorf("Illegal variable name: %s", variable)
	}
	vref := resolve.VariableReference(intp.runtime.ScopeTree)
	if vref == nil {
		return VarValue{
			VariableFullName: resolve.VariableName(),
			val:              nil,
		}, fmt.Errorf("Variable is undefined")
	}
	v, _ := corelang.FindVariableReferenceInMemory(intp.runtime, vref, false)
	if v == nil {
		return VarValue{
			VariableFullName: vref.GetFullName(),
			val:              nil,
		}, nil
	}
	return VarValue{
		VariableFullName: v.GetFullName(),
		val:              v.GetValue(),
	}, nil
}

// ValueString returns the current value of a variable.
// Return values are: canonical name of the variable & value of the
// variable as string.
func (intp *PMMPostInterpreter) ValueString(variable string) (string, string) {
	variable = strings.TrimSpace(variable)
	r, _ := utf8.DecodeRuneInString(variable)
	if !unicode.IsLetter(r) {
		return variable, "<illegal variable name>"
	}
	resolve, err := varparse.ParseVariableName(variable)
	if err != nil {
		return variable, "<illegal variable name>"
	}
	vref := resolve.VariableReference(intp.runtime.ScopeTree)
	if vref == nil {
		return variable, "<undefined>"
	}
	v, _ := corelang.FindVariableReferenceInMemory(intp.runtime, vref, false)
	if v == nil {
		t := variables.TypeString(vref.Type())
		return vref.GetFullName(), fmt.Sprintf("<unset %s>", t)
	}
	return v.GetFullName(), v.ValueString()
}
