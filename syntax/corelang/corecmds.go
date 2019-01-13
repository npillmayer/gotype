package corelang

/*
---------------------------------------------------------------------------

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

---------------------------------------------------------------------------

Internal commands for our core DSL, borrowing from MetaFont/MetaPost.

*/

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/npillmayer/gotype/core/arithmetic"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	"github.com/npillmayer/gotype/syntax/variables/varparse"
	dec "github.com/shopspring/decimal"
)

// Declaration of 'whatever', used to instantiate anonymous whatever-variables.
var WhateverDeclaration *variables.PMMPVarDecl

// --- Custom Error Listener -------------------------------------------------

// We create our own type of error listener for the ANTLR parser
type TracingErrorListener struct {
	*antlr.DefaultErrorListener // use default as base class
}

// Our error listener prints an error to the trace.
func (c *TracingErrorListener) SyntaxErrorf(r antlr.Recognizer, sym interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	//
	at := fmt.Sprintf("%s:%s", strconv.Itoa(line), strconv.Itoa(column))
	T().P("line", at).Errorf("%.44s", msg)
}

// === Handling Variables and Constants ======================================

/*
Construct a valid variable reference string from parts on the stack.

Collect fragments of a variable reference, e.g. "x[k+1]r".
Subscripts should be found on the expression stack and inserted as
numeric constants, i.e. resulting in "x[5]r" (if k=4).

Parameter t is the text of the variable ref literal, e.g. "x[k+1]r".
It is split by the parser into:

  . "x" -> TAG x
  . subscript { k+1 }
  . "r" -> TAG r

*/
func CollectVarRefParts(rt *runtime.Runtime, t string, children []antlr.Tree) string {
	var vname bytes.Buffer
	for _, ch := range children {
		T().Debugf("collecting var ref part: %s", GetCtxText(ch))
		if IsTerminal(ch) { // just copy string parts to output
			T().Debugf("adding suffix verbatim: %s", GetCtxText(ch))
			vname.WriteString(ch.(antlr.ParseTree).GetText())
		} else { // non-terminal is a subscript-expression
			subscript, ok := rt.ExprStack.Pop() // take subscript from stack
			if !ok {
				T().P("var", t).Errorf("expected subscript on expression stack")
				T().P("var", t).Errorf("substituting 0 instead")
				vname.WriteString("[0]")
			} else {
				c, isconst := subscript.XPolyn.IsConstant()
				if !isconst { // we cannot handle unknown subscripts
					T().P("var", t).Errorf("subscript must be known numeric")
					T().P("var", t).Errorf("substituting 0 for %s",
						rt.ExprStack.TraceString(subscript))
					vname.WriteString("[0]")
				} else {
					vname.WriteString("[")
					vname.WriteString(c.String())
					vname.WriteString("]")
				}
			}
		}
	}
	varname := vname.String()
	T().P("var", varname).Debugf("collected parts")
	return varname
}

/*
Get or create a variable reference. To get the canonical representation of
the variable reference, we parse it and construct a small AST. This AST
is fed into GetVarRefFromVarSyntax(). The resulting variable reference
struct is used to find the memory location of the variable reference.

Example:

	vref := MakeCanonicalAndResolve(rt, "a2r", true)
	// now vref.String() gives something like:
	//      "<var a[2].r=<nil> w/ <decl a[].r/numeric>>"

If a variable has been undeclared and is now created, the top-most scope
and memory-frame will hold the newly created variable.
*/
func MakeCanonicalAndResolve(rt *runtime.Runtime, v string, create bool) (
	*variables.PMMPVarRef, error) {
	//
	var vref *variables.PMMPVarRef
	resolve, err := varparse.ParseVariableName(v)
	if err != nil {
		return nil, err
	}
	vref = resolve.VariableReference(rt.ScopeTree)
	if vref != nil {
		vref, _ = FindVariableReferenceInMemory(rt, vref, create)
	}
	return vref, nil
}

// Allocate a variable in a memory frame. Existing variable references in
// this memory frame will be overwritten !
// Clients should probably first call FindVariableReferenceInMemory(vref).
func AllocateVariableInMemory(vref *variables.PMMPVarRef,
	mf *runtime.DynamicMemoryFrame) *variables.PMMPVarRef {
	//
	mf.Symbols().InsertSymbol(vref)
	T().P("var", vref.GetFullName()).Debugf("allocating variable in %s", mf.GetName())
	return vref
}

/*
Given a variable reference, locate an incarnation in a memory frame. The
frame is determined by the variable's declaring scope: search for the top
frame linked to the scope.

Variable references live in memory frames. Memory frames correspond to
scopes. To find a variable reference -- i.e. a living variable with a possible
value -- we have to proceed as follows:

(1) find the variable declaration in a scope, beginning at the top

(2) find the most recent memory frame pointing to this scope

(3) find a variable reference with the correct name in the memory frame

(4) if no reference/incarnation exists, create one

Parameter doAlloc: should step (4) be performed ?
*/
func FindVariableReferenceInMemory(rt *runtime.Runtime, vref *variables.PMMPVarRef, doAlloc bool) (
	*variables.PMMPVarRef, *runtime.DynamicMemoryFrame) {
	//
	if vref.Decl == nil {
		T().P("var", vref.GetFullName()).Errorf("attempt to store variable without decl. in memory")
		return vref, nil
	}
	var sym *variables.PMMPVarRef
	var memframe *runtime.DynamicMemoryFrame
	tagname := vref.Decl.BaseTag.GetName()
	tag, scope := rt.ScopeTree.Current().ResolveSymbol(tagname)
	if tag != nil { // found tag declaration in scope
		memframe = rt.MemFrameStack.FindMemoryFrameWithScope(scope)
		varname := vref.GetName()
		T().P("var", varname).Debugf("var in ? %s", memframe)
		s := memframe.Symbols().ResolveSymbol(varname)
		if s == nil { // no variable ref incarnation => create one
			T().P("var", varname).Debugf("not found in memory")
			if doAlloc {
				sym = AllocateVariableInMemory(vref, memframe)
			}
		} else { // already present, return this one
			T().P("var", varname).Debugf("variable already present in memory")
			sym = s.(*variables.PMMPVarRef)
		}
	} else {
		// this should never happen: we could neither find nor construct a var decl
		panic(fmt.Sprintf("declaration for %s mysteriously vanished...", tagname))
	}
	return sym, memframe
}

// The expression stack knows nothing about the interpreter's symbols, except
// the few properties of interface Symbol. The expression stack deals with
// polynomials and serial IDs of variables.
//
// Push a variable (numeric or pair type) onto the expression stack.
//
func PushVariable(rt *runtime.Runtime, vref *variables.PMMPVarRef, asLValue bool) {
	if vref.IsPair() {
		if vref.IsKnown() && !asLValue {
			PushConstant(rt, vref) // put constant on expression stack
		} else {
			xpart, ypart := vref.XPart(), vref.YPart()
			rt.ExprStack.PushVariable(xpart, ypart)
		}
	} else {
		if vref.IsKnown() && !asLValue {
			PushConstant(rt, vref) // put constant on expression stack
		} else {
			rt.ExprStack.PushVariable(vref, nil)
		}
	}
}

// Push a constant (numeric or pair type) onto the expression stack.
func PushConstant(rt *runtime.Runtime, vref *variables.PMMPVarRef) {
	switch vref.Decl.Type() {
	case variables.NumericType:
		rt.ExprStack.PushConstant(vref.Value.(dec.Decimal))
	case variables.PairType:
		x := vref.XPart().GetValue()
		y := vref.YPart().GetValue()
		pair := arithmetic.MakePair(x.(dec.Decimal), y.(dec.Decimal))
		rt.ExprStack.PushPairConstant(pair)
	case variables.PathType:
		rt.ExprStack.PushOtherConstant(vref.Value)
	}
	/*
		if vref.IsPair() {
			x := vref.XPart().GetValue()
			y := vref.YPart().GetValue()
			pair := arithmetic.MakePair(x.(dec.Decimal), y.(dec.Decimal))
			rt.ExprStack.PushPairConstant(pair)
		} else {
			rt.ExprStack.PushConstant(vref.Value.(dec.Decimal))
		}
	*/
}

// The expression stack knows nothing about the interpreter's symbols, except
// the few properties of interface Symbol. The expression stack deals with
// polynomials and serial IDs of variables. To get back from IDs to
// variable references, we ask the expression stack for a Symbol (from an
// ID). If the variable is of type pair, the Symbol may be a pair part (x-part
// or y-part). Parts point to their parent symbol, thus giving us the
// variable reference.
func GetVariableFromExpression(rt *runtime.Runtime, e *runtime.ExprNode) *variables.PMMPVarRef {
	var v *variables.PMMPVarRef
	if sym := rt.ExprStack.GetVariable(e); sym != nil {
		var part *variables.PairPartRef
		var ok bool
		if part, ok = sym.(*variables.PairPartRef); ok {
			sym = part.Pairvar
		}
		v = sym.(*variables.PMMPVarRef)
		T().P("var", v.GetName()).Debugf("variable of type %s", variables.TypeString(v.Type()))
	}
	return v
}

// A variable which goes out of scope becomes a "capsule". We send a message
// to the expression stack to forget the Symbol(s) for the ID(s) of a
// variable. Variables are of type numeric or pair.
func EncapsulateVariable(rt *runtime.Runtime, v *variables.PMMPVarRef) {
	rt.ExprStack.EncapsuleVariable(v.GetID())
	if v.IsPair() {
		var ypart *variables.PairPartRef = v.GetFirstChild().(*variables.PairPartRef)
		rt.ExprStack.EncapsuleVariable(ypart.GetID())
	}
}

// Make all variables in a memory frame "capsules".
//
// When a memory frame is popped from the stack, the local variables living
// in the frame have to be made "capsules". This is necessary, because they
// may still be relevant to the LEQ-solver. The LEQ will finally decide
// when to abondon the "zombie" variable.
func EncapsulateVarsInMemory(rt *runtime.Runtime, mf *runtime.DynamicMemoryFrame) {
	mf.Symbols().Each(func(name string, sym runtime.Symbol) {
		vref := sym.(*variables.PMMPVarRef)
		T().P("var", vref.GetFullName()).Debugf("encapsule")
		rt.ExprStack.EncapsuleVariable(vref.GetID()) // vref is now capsule
	})
}

// Load builtin symbols into a scope (usually the global scope).
// Additionally loads initial Lua definitions.
func LoadBuiltinSymbols(rt *runtime.Runtime, scripting *Scripting) {
	originDef := Declare(rt, "origin", variables.PairType)
	origin := arithmetic.MakePair(arithmetic.ConstZero, arithmetic.ConstZero)
	_ = Variable(rt, originDef, origin, nil, true)
	upDef := Declare(rt, "up", variables.PairType)
	up := arithmetic.MakePair(arithmetic.ConstZero, arithmetic.ConstOne)
	_ = Variable(rt, upDef, up, nil, true)
	downDef := Declare(rt, "down", variables.PairType)
	down := arithmetic.MakePair(arithmetic.ConstZero, arithmetic.MinusOne)
	_ = Variable(rt, downDef, down, nil, true)
	rightDef := Declare(rt, "right", variables.PairType)
	right := arithmetic.MakePair(arithmetic.ConstOne, arithmetic.ConstZero)
	_ = Variable(rt, rightDef, right, nil, true)
	leftDef := Declare(rt, "left", variables.PairType)
	left := arithmetic.MakePair(arithmetic.MinusOne, arithmetic.ConstZero)
	_ = Variable(rt, leftDef, left, nil, true)
	_ = Declare(rt, "p", variables.PairType)
	_ = Declare(rt, "q", variables.PairType)
	_ = Declare(rt, "P", variables.PathType)
	w := Declare(rt, "_whtvr", variables.NumericType) // 'whatever' variables
	// make whatever[]
	WhateverDeclaration = variables.CreatePMMPVarDecl("<array>", variables.ComplexArray, w)
}

// === Commands ==============================================================

/*
Variable assignment.

   assignment : lvalue ASSIGN numtertiary

(1) Retract lvalue from the resolver's table (make a capsule)

(3) Unset the value of lvalue

(3) Re-incarnate lvalue (get a new ID for it)

(4) If type is numeric or pair: Create equation on expression stack,
else assign a path value to a path variable.
*/
func Assign(rt *runtime.Runtime, lvalue *variables.PMMPVarRef, e *runtime.ExprNode) {
	varname := lvalue.GetName()
	oldserial := lvalue.GetID()
	T().P("var", varname).Debugf("assignment of lvalue #%d", oldserial)
	EncapsulateVariable(rt, lvalue)
	vref, mf := FindVariableReferenceInMemory(rt, lvalue, false)
	vref.SetValue(nil) // now lvalue is unset / unsolved
	T().P("var", varname).Debugf("unset in %v", mf)
	vref.Reincarnate()
	T().P("var", vref.GetName()).Debugf("new lvalue incarnation #%d", vref.GetID())
	if vref.Type() == variables.PathType {
		vref.SetValue(e.Other)
	} else { // create linear equation
		PushVariable(rt, vref, false) // push LHS on stack
		rt.ExprStack.Push(e)          // push RHS on stack
		rt.ExprStack.EquateTOS2OS()   // construct equation
	}
}

// Save a tag within a group. The tag will be restored at the end of the
// group. Save-commands within global scope will be ignored.
// This method simply creates a var decl for the tag in the current scope.
func Save(rt *runtime.Runtime, tag string) {
	sym, scope := rt.ScopeTree.Current().ResolveSymbol(tag)
	if sym != nil { // already found in scope stack
		T().P("tag", tag).Debugf("save: found tag in scope %s",
			scope.GetName())
	}
	T().Debugf("declaring %s in current scope", tag)
	rt.ScopeTree.Current().DefineSymbol(tag)
}

// Declare a tag to be of type tp.
//
// If the tag is not declared, insert a new symbol in global scope. If a
// declaration already exists, erase all variables and re-enter a declaration
// (MetaFont semantics). If the tag has been "saved" in the current or in an outer
// scope, make this tag a new undefined symbol.
//
func Declare(rt *runtime.Runtime, tag string, tp variables.VariableType) *variables.PMMPVarDecl {
	sym, scope := rt.ScopeTree.Current().ResolveSymbol(tag)
	if sym != nil { // already found in scope stack
		T().P("tag", tag).Debugf("declare: found tag in scope %s", scope.GetName())
		T().P("decl", tag).Debugf("variable already declared - re-declaring")
		// Erase all existing variables and re-define symbol
		sym, _ = scope.DefineSymbol(tag)
		sym.(*variables.PMMPVarDecl).SetType(int(tp))
	} else { // enter new symbol in global scope
		scope = rt.ScopeTree.Globals()
		sym, _ = scope.DefineSymbol(tag)
		sym.(*variables.PMMPVarDecl).SetType(int(tp))
	}
	T().P("decl", sym.GetName()).Debugf("declared symbol in %s", scope.GetName())
	return sym.(*variables.PMMPVarDecl)
}

// Create a variable reference. Parameters are the declaration for the variable,
// a value and a flag, indicating if this variable should go to global memory.
// The subscripts parameter is a slice of array-subscripts, if the variable
// declaration is of array (complex) type.
func Variable(rt *runtime.Runtime, decl *variables.PMMPVarDecl, value interface{},
	subscripts []dec.Decimal, global bool) *variables.PMMPVarRef {
	//
	var v *variables.PMMPVarRef
	if decl.Type() == variables.NumericType {
		v = variables.CreatePMMPVarRef(decl, value, subscripts)
	} else {
		v = variables.CreatePMMPPairTypeVarRef(decl, value, subscripts)
	}
	if global {
		rt.MemFrameStack.Globals().Symbols().InsertSymbol(v)
	} else {
		rt.MemFrameStack.Current().Symbols().InsertSymbol(v)
	}
	return v
}

// Create a whatever anonymous variable. In MetaFont this is a macro, but
// it is a frequent use case, so we put it in the core.
func Whatever(rt *runtime.Runtime) *variables.PMMPVarRef {
	var vref *variables.PMMPVarRef
	sym, _ := rt.ScopeTree.Globals().ResolveSymbol("_whtvr")
	if sym == nil {
		T().Errorf("'whatever'-variable not correctly initialized")
	} else {
		//func CreatePMMPVarRef(*PMMPVarDecl, value, indices []dec.Decimal) *PMMPVarRef {
		inx := make([]dec.Decimal, 1)
		whateverCounter++
		inx[0] = dec.New(whateverCounter, 0)
		vref = variables.CreatePMMPVarRef(WhateverDeclaration, nil, inx)
	}
	return vref
}

// Counter for 'whatever' anonymous variables.
var whateverCounter int64

// Apply a (math or scripting) function, given by name, to a known/constant argument.
// Internal math functions are floor(), ceil() and sqrt(). Other function names
// will be delegated to the scripting subsystem (Lua).
//
// Lua functions will return just one value (of type numeric, pair or path).
//
func CallFunc(val interface{}, fun string, scripting *Scripting) (*runtime.ExprNode, []*variables.PMMPVarRef) {
	n := arithmetic.ConstZero
	if strings.HasPrefix(fun, "@") {
		fun = strings.TrimLeft(fun, "@")
		T().P("func", fun).Debugf("calling Lua scripting subsytem")
		r, err := scripting.CallHook(fun, val)
		if err == nil {
			it := r.Iterator() // iterator over return values
			if it.Next() {     // go to first return value
				e, vars := it.ValueAsExprNode() // unpack first return value only
				return e, vars
			}
		} else {
			T().P("func", fun).Errorf("scripting error: %v", err.Error())
		}
	} else {
		switch fun {
		case "floor":
			n = val.(dec.Decimal)
			n = n.Floor()
		case "ceil":
			n = val.(dec.Decimal)
			n = n.Ceil()
		case "sqrt":
			T().P("func", fun).Errorf("function not yet implemented")
		default:
			T().P("func", fun).Errorf("function not implemented")
		}
	}
	p := arithmetic.NewConstantPolynomial(n)
	e := runtime.NewNumericExpression(p)
	return e, nil
}

// Call the Lua script for a vardef variable. The parameters of the call will
// be the suffixes of the variable, i.e. all subscripts and tags after the base
// tag.
//
// Return values are wrapped into an expression node. If variables are part of
// the returned expressions, they are packed into an array of variable
// references.
//
func CallVardef(vref *variables.PMMPVarRef, scripting *Scripting) (*runtime.ExprNode, []*variables.PMMPVarRef) {
	basetag := vref.Decl.BaseTag
	suffixes := vref.GetSuffixesString()
	T().Debugf("vardef call %s(%s)", basetag.GetName(), suffixes)
	r, err := scripting.Call("vardef", basetag.GetName(), suffixes)
	if err == nil {
		it := r.Iterator() // iterator over return values
		if it.Next() {     // go to first return value
			e, vars := it.ValueAsExprNode() // unpack first return value only
			return e, vars
		}
	} else {
		T().P("vardef", basetag.GetName()).Errorf("scripting error: %v", err.Error())
	}
	return nil, nil
}

// MetaFont begingroup command: push a new scope and memory frame.
// Clients may supply a name for the group, otherwise it will be set
// to "group".
func Begingroup(rt *runtime.Runtime, name string) (*runtime.Scope, *runtime.DynamicMemoryFrame) {
	if name == "" {
		name = "group"
	}
	groupscope := rt.ScopeTree.PushNewScope(name, variables.NewPMMPVarDecl)
	groupmf := rt.MemFrameStack.PushNewMemoryFrame(name, groupscope)
	return groupscope, groupmf
}

//MetaFont endgroup command: pop scope and memory frame of group.
func Endgroup(rt *runtime.Runtime) {
	mf := PopScopeAndMemory(rt)
	EncapsulateVarsInMemory(rt, mf)
}

// Decrease grouping level.
// We pop the topmost scope and topmost memory frame. This happens after
// a group is left.
//
// Returns the previously topmost memory frame.
//
func PopScopeAndMemory(rt *runtime.Runtime) *runtime.DynamicMemoryFrame {
	hidden := rt.ScopeTree.PopScope()
	hidden.Name = "(hidden)"
	mf := rt.MemFrameStack.PopMemoryFrame()
	if mf.GetScope() != hidden {
		T().P("mem", mf.GetName()).Errorf("groups out of sync?")
	}
	return mf
}

// === Show Commands =========================================================

// Show all declarations and references for a tag.
func Showvariable(rt *runtime.Runtime, tag string) string {
	sym, scope := rt.ScopeTree.Current().ResolveSymbol(tag)
	if sym == nil {
		return fmt.Sprintf("%s : tag\n", tag)
	} else {
		v := sym.(*variables.PMMPVarDecl)
		var b *bytes.Buffer
		b = v.ShowDeclarations(b)
		if mf := rt.MemFrameStack.FindMemoryFrameWithScope(scope); mf != nil {
			for _, v := range mf.Symbols().Table {
				vref := v.(*variables.PMMPVarRef)
				if vref.Decl.BaseTag == sym {
					s := fmt.Sprintf("%s = %s\n", vref.GetFullName(), vref.ValueString())
					b.WriteString(s)
				}
			}
		}
		return b.String()
	}
}

// === Utilities =============================================================

/*
TODO complete this.
Return scaled points for high level units (cm, mm, pt, in, ...)
*/
func Unit2numeric(u string) dec.Decimal {
	switch u {
	case "in":
		return dec.NewFromFloat(0.01388888)
	}
	return arithmetic.ConstOne
}

// Scale a numeric value by a unit.
func ScaleDimension(dimen dec.Decimal, unit string) dec.Decimal {
	u := Unit2numeric(unit)
	return dimen.Mul(u)
}

// --- Utilities -------------------------------------------------------------

// I do not always understand ANTLR V4's Go runtime typing and tree semantics
// (rather poorly documented), so I introduce some helpers. Some of these are
// probably unnecessary for a better versed ANTLR Go user...
func IsTerminal(node antlr.Tree) bool {
	_, ok := node.GetPayload().(antlr.RuleNode)
	//fmt.Printf("node is terminal: %v\n", !ok)
	return !ok
}

// I do not always understand ANTLR V4's Go runtime typing and tree semantics
// (rather poorly documented), so I introduce some helpers. Some of these are
// probably unnecessary for a better versed ANTLR Go user...
func GetCtxText(ctx antlr.Tree) string {
	t := ctx.(antlr.ParseTree).GetText()
	return t
}
