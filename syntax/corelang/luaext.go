package corelang

import (
	"fmt"

	"github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtlocate"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	"github.com/shopspring/decimal"
	lua "github.com/yuin/gopher-lua"
)

// ---------------------------------------------------------------------------

/*
Type Scripting is an opaque data type to provide access to the
scripting sub-system.

DSLs built on top of this language core may be scripted with Lua.
Lua scripts may be called as hooks or as functions on primary level.
Lua functions are preceded by an '@'.

Example:

    a2r = 7 + @inlua(x0)

This will delegate to the Lua scripting subsystem, putting the value of x0 onto
the Lua stack, and then call Lua-function inlua(...) on it.
*/
type Scripting struct {
	*lua.LState
	hooks   map[string]lua.LGFunction
	runtime *runtime.Runtime
}

// Create a new scripting subsystem. Scripting sub-systems are not thread safe.
func NewScripting(rt *runtime.Runtime) *Scripting {
	luastate := lua.NewState()
	if luastate == nil {
		T.Error("failed to create Lua scripting subsystem")
		return nil
	}
	T.Info("Loading initial Lua scripts")
	hostlang := gtlocate.FileResource("hostlang", "lua")
	T.Infof("- %s", hostlang)
	luastate.DoFile(hostlang)
	luastate.DoString("hostlang = HostLang")
	scr := &Scripting{luastate, nil, rt}
	scr.registerDSLRuntimeEnvType()
	scr.registerVarRefType()
	scr.registerPairType()
	scr.hooks = make(map[string]lua.LGFunction)
	scr.RegisterHook("trace", trace)
	T.Info("Scripting initialized")
	return scr
}

// --- Scripting Arguments ---------------------------------------------------

// We shield the Lua implementation from our scripting API.
// Parameters and return values are provided as generic values.
type scriptingArgs []interface{}

func asScriptingArgs(values ...interface{}) scriptingArgs {
	args := make(scriptingArgs, len(values))
	for i, a := range values {
		args[i] = a
	}
	return args
}

// Convert a Go type argument list into Lua types
func forLua(values scriptingArgs) []lua.LValue {
	args := make([]lua.LValue, len(values))
	for i, a := range values {
		if a == nil {
			args[i] = lua.LNil
		} else {
			if d, ok := a.(decimal.Decimal); ok {
				f, _ := d.Float64()
				args[i] = lua.LNumber(f)
			} else if s, isstr := a.(string); isstr {
				args[i] = lua.LString(s)
			} else if n, isnum := a.(int); isnum {
				args[i] = lua.LNumber(n)
			} else if f, isfloat := a.(float64); isfloat {
				args[i] = lua.LNumber(f)
			} else {
				// TODO implement UserData types
				T.Error("type not implemented for scripting")
			}
		}
	}
	return args
}

// --- Helpers for decoding Lua LValues --------------------------------------

func isUserData(lv lua.LValue) (*lua.LUserData, bool) {
	if lv.Type() == lua.LTUserData {
		udata := lv.(*lua.LUserData)
		return udata, true
	}
	return nil, false
}

func isVariable(lv lua.LValue) (*runtime.ExprNode, *variables.PMMPVarRef, bool) {
	if udata, ok := isUserData(lv); ok {
		if v, isvref := udata.Value.(*LuaVarRef); isvref {
			var e *runtime.ExprNode
			if v.vref.IsPair() {
				x := v.vref.XPart()
				y := v.vref.YPart()
				e = runtime.NewPairVarExpression(x, y)
			} else {
				e = runtime.NewNumericVarExpression(v.vref)
			}
			return e, v.vref, true
		}
	}
	return nil, nil, false
}

func isNumericConstant(lv lua.LValue) (decimal.Decimal, bool) {
	if lv.Type() == lua.LTNumber {
		f := float64(lv.(lua.LNumber))
		n := decimal.NewFromFloat(f)
		return n, true
	}
	return arithmetic.ConstZero, false
}

func isPair(lv lua.LValue) (*LuaPair, bool) {
	if udata, ok := isUserData(lv); ok {
		if lpr, ispair := udata.Value.(*LuaPair); ispair {
			return lpr, true
		}
	}
	return nil, false
}

func isLiteralPair(lv lua.LValue) (*runtime.ExprNode, []*variables.PMMPVarRef, bool) {
	var vars []*variables.PMMPVarRef = make([]*variables.PMMPVarRef, 2)
	if lpr, ok := isPair(lv); ok {
		var x, y arithmetic.Polynomial
		if xc, xisconst := isNumericConstant(lpr.X); xisconst {
			x = arithmetic.NewConstantPolynomial(xc)
		} else if xv, v, xisvar := isVariable(lpr.X); xisvar {
			x = xv.XPolyn
			vars[0] = v
		} else {
			T.Error("illegal x-part for pair returned from Lua, assuming 0")
			x = arithmetic.NewConstantPolynomial(arithmetic.ConstZero)
		}
		if yc, yisconst := isNumericConstant(lpr.Y); yisconst {
			y = arithmetic.NewConstantPolynomial(yc)
		} else if yv, v, yisvar := isVariable(lpr.Y); yisvar {
			y = yv.XPolyn
			vars[1] = v
		} else {
			T.Error("illegal y-part for pair returned from Lua, assuming 0")
			y = arithmetic.NewConstantPolynomial(arithmetic.ConstZero)
		}
		return runtime.NewPairExpression(x, y), vars, true
	}
	return nil, vars, false
}

func isString(lv lua.LValue) (string, bool) {
	if lv.Type() == lua.LTString {
		s := lua.LVAsString(lv)
		return s, true
	}
	return "", false
}

// ---------------------------------------------------------------------------

// Type to return values from Lua scripts.
// Single values are accessed with an iterator.
//
// see ScriptingReturnValueIterator
type ScriptingReturnValues struct {
	values []lua.LValue
}

// Iterator type for scripting return values.
// Return values from Lua are wrapped into an opaque type
// ScriptingReturnValues and accessed using this iterator type.
//
// see ScriptingReturnValues.Iterator()
type ScriptingReturnValueIterator struct {
	values *ScriptingReturnValues
	inx    int
}

// Create an iterator for scripting arguments / return values.
func (r *ScriptingReturnValues) Iterator() *ScriptingReturnValueIterator {
	return &ScriptingReturnValueIterator{r, -1}
}

// Is there a next scripting argument?
// Advances the iterator's cursor.
func (it *ScriptingReturnValueIterator) Next() bool {
	T.Debugf("%d return values to iterate", len(it.values.values))
	it.inx++
	if it.inx < len(it.values.values) {
		return true
	} else {
		return false
	}
}

// Get the value of the scripting argument under the iterator's cursor.
// Returns the value and a type (see package 'variables' for the
// definition of variable types).
func (it *ScriptingReturnValueIterator) Value() (interface{}, int) {
	if it.inx < len(it.values.values) {
		a := it.values.values[it.inx]
		if a == nil {
			return nil, variables.Undefined
		} else if v, ok := isNumericConstant(a); ok {
			return v, variables.NumericType
		} else {
			T.Errorf("not yet implemented: type for %v", a)
		}
	}
	return nil, variables.Undefined
}

// Get the value of the scripting argument under the iterator's cursor.
// Returns the value wrapped in an expression node (or nil).
// If variables are part of the expression(s), they are returned in a
// separate array.
func (it *ScriptingReturnValueIterator) ValueAsExprNode() (*runtime.ExprNode, []*variables.PMMPVarRef) {
	if it.inx < len(it.values.values) {
		lv := it.values.values[it.inx]
		T.Debugf("iterator: value as expression = %v", lv)
		if lua.LVIsFalse(lv) {
			T.Error("'nil' return value from Lua, substituting numeric 0")
			p := arithmetic.NewConstantPolynomial(arithmetic.ConstZero)
			return runtime.NewNumericExpression(p), nil
		} else if _, ok := isString(lv); ok {
			T.Error("cannot process string return value, substituting numeric 0")
			p := arithmetic.NewConstantPolynomial(arithmetic.ConstZero)
			return runtime.NewNumericExpression(p), nil
		} else if e, v, ok := isVariable(lv); ok {
			T.Debugf("return values iterator: variable = %v", e)
			vars := make([]*variables.PMMPVarRef, 1)
			vars[0] = v
			return e, vars
		} else if c, ok := isNumericConstant(lv); ok {
			p := arithmetic.NewConstantPolynomial(c)
			return runtime.NewNumericExpression(p), nil
		} else if e, vxy, ok := isLiteralPair(lv); ok {
			T.Debugf("return values iterator: literal pair = %v", e)
			return e, vxy
		} else {
			T.Error("unknown return value from Lua, substituting numeric 0")
			p := arithmetic.NewConstantPolynomial(arithmetic.ConstZero)
			return runtime.NewNumericExpression(p), nil
		}
	}
	return nil, nil
}

// ---------------------------------------------------------------------------

/*
Register a hook function for a key given as string parameter. The
hook function must accept a single argument: the Lua state, and return
a single int: the number of return values on the Lua stack.

Hooks may be called from the Lua side by name, or from the Go side by
CallHook(...).

Example:

    scripting := NewScripting()
    scripting.RegisterHook("stars", func(L *lua.LState) int {
        L.Push(lua.LString("* * * * *")) // push result
        return 1                         // return value count
    })

In Lua:

    print(stars())    -- prints "* * * * *" to stdout

*/
func (lscript *Scripting) RegisterHook(name string, f lua.LGFunction) {
	lscript.Register(name, f)
	lscript.hooks[name] = f
}

/*
Call a registered hook from Go. Arguments may be passed (Go data types) in
a variable argument list. Return values are converted back from Lua types to
Go types.

see RegisterHook()
*/
func (lscript *Scripting) CallHook(hook string, arguments ...interface{}) (
	*ScriptingReturnValues, error) {
	//
	T.P("lua", hook).Debug("call hook")
	if lscript.hooks[hook] == nil {
		T.P("lua", hook).Error("hook not registered")
	}
	r, err := lscript.Call("", hook, arguments...)
	return r, err
}

/*
Call a Lua function, possibly qualified by a table prefix.
*/
func (lscript *Scripting) Call(table string, function string, arguments ...interface{}) (
	*ScriptingReturnValues, error) {
	//
	var lfunc lua.LValue
	if table == "" {
		lfunc = lscript.GetGlobal(function)
	} else {
		t := lscript.GetGlobal(table)
		lfunc = lscript.GetField(t, function)
	}
	args := asScriptingArgs(arguments...)
	err := lscript.CallByParam(lua.P{
		Fn:      lfunc,
		NRet:    1,
		Protect: true,
	}, forLua(args)...)
	var r *ScriptingReturnValues
	if err == nil {
		T.P("lua", function).Debugf("%d return values on the stack", lscript.GetTop())
		r = lscript.returnFromScripting(lscript.GetTop()) // return all values on the stack
	}
	return r, err
}

// return n values (off the stack) from a Lua call.
func (lscript *Scripting) returnFromScripting(n int) *ScriptingReturnValues {
	values := make([]lua.LValue, n)
	i, j := n, -1
	for ; i > 0; i-- {
		lv := lscript.Get(j)
		values[-(j + 1)] = lv
		T.Debugf("return value %d = %v", -(j + 1), lv)
		j--
	}
	lscript.Pop(n)
	rv := &ScriptingReturnValues{values: values}
	return rv
}

/*
Evaluate a Lua statement, given as string. Return arguments (from the Lua stack)
are packed into an opaque data structure. The second return value is a possible
error condition. The Lua command(s) must be syntactically correct and complete
statements (no expressions etc. accepted).

Eval accepts a variable number of untyped arguments. These are put on the
Lua stack before the statement is executed.
*/
func (lscript *Scripting) Eval(luacmd string, arguments ...interface{}) (*ScriptingReturnValues, error) {
	args := asScriptingArgs(arguments...)
	largs := forLua(args)
	for _, larg := range largs {
		lscript.Push(larg)
	}
	var r *ScriptingReturnValues
	err := lscript.DoString(luacmd)
	if err != nil {
		T.P("script", "lua").Errorf("scripting error: %s", err.Error())
	} else {
		if err == nil {
			T.P("lua", "eval").Debugf("%d return values on the stack", lscript.GetTop())
			r = lscript.returnFromScripting(lscript.GetTop()) // return all values on the stack
		}
	}
	return r, err
}

// For testing purposes
func ping(L *lua.LState) int {
	L.Push(lua.LString("ok")) // push return value on stack
	return 1                  // number of results
}

// For testing purposes
func echo(L *lua.LState) int {
	lv := L.Get(-1)
	msg := fmt.Sprintf("echo: %s !", lua.LVAsString(lv))
	L.Push(lua.LString(msg))
	return 1
}

// Tracer for the Lua scripting sub-system.
// Will be registered under 'trace(level, message)'.
func trace(L *lua.LState) int {
	s := L.CheckString(-1)
	level := L.CheckInt(-2)
	if level == 0 {
		S.Debug(s)
	} else {
		S.Info(s)
	}
	return 0
}

// vardef z(suffixes)
/*
func lua_z(L *lua.LState) int {
	suffixes := L.Get(-1)
	T.Debugf("lua_z(%v)", suffixes)
	L.Push(suffixes)
	hostlang := L.GetGlobal("hostlang")
	fmt.Printf("hostlang     = %v\n", hostlang)
	z := L.GetField(hostlang, "z")
	fmt.Printf("hostlang.z() = %v\n", z)
	err := L.CallByParam(lua.P{
		Fn:      z,
		NRet:    1,    // number of returned values
		Protect: true, // return err or panic
	}, suffixes)
	if err != nil {
		T.Errorf("CallByParam error: %v", err.Error())
	}
	return 1
}
*/

// === User Data Type: Pair ==================================================

const luaPairTypeName = "pair"

/*
Lua UserData type for pairs.

Example (Lua):

  p = pair.new{2, 5}
  print(p:x())          -- get x-part
  p:y(3.14)             -- set y-part

*/
type LuaPair struct {
	X lua.LValue
	Y lua.LValue
}

var pairMethods = map[string]lua.LGFunction{
	"x": pairGetSetX,
	"y": pairGetSetY,
}

// Registers pair type to given Scripting.
func (lscript *Scripting) registerPairType() {
	mt := lscript.NewTypeMetatable(luaPairTypeName)
	lscript.SetGlobal("pair", mt)
	// static attributes
	lscript.SetField(mt, "new", lscript.NewFunction(newPair))
	// methods
	lscript.SetField(mt, "__index", lscript.SetFuncs(lscript.NewTable(), pairMethods))
}

// Constructor
func newPair(L *lua.LState) int {
	T.P("lua", "udata").Debug("creating pair() UserData")
	var pair *LuaPair
	xy := L.CheckTable(1)
	if xy.MaxN() != 2 {
		T.P("lua", "udata").Errorf("wrong number of arguments to pair:new(): %d", xy.MaxN())
		pair = &LuaPair{lua.LNumber(0), lua.LNumber(0)}
	} else {
		x := xy.RawGet(lua.LNumber(1))
		y := xy.RawGet(lua.LNumber(2))
		pair = &LuaPair{x, y}
	}
	ud := L.NewUserData()
	ud.Value = pair
	L.SetMetatable(ud, L.GetTypeMetatable(luaPairTypeName))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Pair and returns
// this *Pair.
func checkPair(L *lua.LState) *LuaPair {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*LuaPair); ok {
		return v
	}
	L.ArgError(1, "pair expected")
	return nil
}

// Getter and setter for pair#x
func pairGetSetX(L *lua.LState) int {
	p := checkPair(L)
	if L.GetTop() == 2 { // setter
		p.X = L.CheckAny(2)
		return 0
	} else { // getter
		L.Push(p.X)
		return 1
	}
}

// Getter and setter for pair#y
func pairGetSetY(L *lua.LState) int {
	p := checkPair(L)
	if L.GetTop() == 2 { // setter
		p.Y = L.CheckAny(2)
		return 0
	} else { // getter
		L.Push(p.Y)
		return 1
	}
}

// === User Data Type: Variable ==============================================

const luaVarRefTypeName = "varref"

/*
Lua UserData type for variables. Variables reference DSL-variables in
the DSL's runtime environment (MetaFont-like variables of type numeric, pair,
etc.) A variable may be known or unknown.

Example (Lua):

    a = hostlang.numeric("a")   -- connect to a tag of the host language
    a[2].r = 3.14               -- assign a numeric value to an 'a'-variable
    print(a[2].r:value())       -- prints "3.14"

Variable a[2].r  (or short: a2r) is now set/known in the host-language (DSL):

    DSL> show a;
    ## show a;                                    tag=a
    a : numeric
    a[] : numeric
    a[].r : numeric
    ## a[2].r = 3.14

Lua's notation for (sub-)tables lends itself nicely for a congruency to
MetaFont-style variable notations. However, it is not possible to use the
DSL shorthand notation ("a2r") for variable names in Lua.

In Lua, there are two member-functions defined for type varref: value() and
isknown(). value() is a getter/setter for the variable value. isknown()
returns a boolean value.

Example (Lua):

    a = hostlang.numeric("a")   -- connect to a tag of the host language
    print(a:isknown())          -- prints "false" if not yet defined in the DSL
    a:value(3.14)               -- must use this notation for 'a' base tag
    print(a:isknown())          -- prints "true"

Variables of this kind are 'live'-objects, i.e. they are always synchronous
between the two languages.
*/
type LuaVarRef struct {
	vref *variables.PMMPVarRef
}

// Stringer for variable references. Used for varref.__tostring(...).
// Will give a debug representation of the DSL-connected variable.
func (lvref *LuaVarRef) String() string {
	if lvref.vref == nil {
		return "<undefined variable>"
	} else {
		return lvref.vref.String()
	}
}

// Metatable functions for type varref
var varRefMethods = map[string]lua.LGFunction{
	"value":   varRefGetSetValue,
	"isknown": varRefIsKnown,
}

// Registers varRef type to given Scripting.
func (lscript *Scripting) registerVarRefType() {
	mt := lscript.NewTypeMetatable(luaVarRefTypeName)
	lscript.SetGlobal("varref", mt)
	lscript.SetField(mt, "refer_to", lscript.NewFunction(referToVar))
	lscript.SetField(mt, "__index", lscript.SetFuncs(lscript.NewTable(), varRefMethods))
	lscript.SetField(mt, "__tostring", lscript.NewFunction(varRef2String))
}

/*
Lua constructor: Construct a variable reference from a variable name.
Variable names are (complex) MetaFont-style variables.

Examples:

    "a", "z[2]", "x[3.14].r", "hello.world[3]"

Essentially performs a call to MakeCanonicalAndResolve(...).
*/
func referToVar(L *lua.LState) int {
	varname := L.CheckString(1)
	if lrt := getGlobalDSLRuntimeEnv(L); lrt != nil {
		vref := MakeCanonicalAndResolve(lrt.rt, varname, true)
		T.Debugf("var. ref. = %v", vref)
		vudata := newVarRefUserData(L, vref)
		L.Push(vudata)
		return 1
	}
	return 0
}

// Create a LuaVarRef UserData wrapper for a variable reference.
// Sets the correct metatable for the variable.
func newVarRefUserData(L *lua.LState, vref *variables.PMMPVarRef) *lua.LUserData {
	vudata := L.NewUserData()
	lvref := &LuaVarRef{vref}
	vudata.Value = lvref
	L.SetMetatable(vudata, L.GetTypeMetatable(luaVarRefTypeName))
	return vudata
}

// Checks whether the first lua argument is a *LUserData with *VarRef and returns
// this *VarRef.
func checkVarRef(L *lua.LState) *LuaVarRef {
	udata := L.CheckUserData(1)
	if v, ok := udata.Value.(*LuaVarRef); ok {
		return v
	}
	L.ArgError(1, "varref expected")
	return nil
}

// Function for varref metatable:
// getter and setter for variable's value()
func varRefGetSetValue(L *lua.LState) int {
	v := checkVarRef(L)
	if L.GetTop() == 2 { // setter
		f := L.CheckNumber(2)
		d := decimal.NewFromFloat(float64(f))
		v.vref.SetValue(d)
		return 0
	} else { // getter
		val := v.vref.GetValue()
		if n, ok := val.(decimal.Decimal); ok {
			lf, _ := n.Float64()
			L.Push(lua.LNumber(lf))
		} else {
			L.Push(lua.LNil)
		}
		return 1
	}
	return 0
}

// Function for varref metatable
func varRefIsKnown(L *lua.LState) int {
	v := checkVarRef(L)
	L.Push(lua.LBool(v.vref.IsKnown()))
	return 1
}

// Function for varref metatable
func varRef2String(L *lua.LState) int {
	v := checkVarRef(L)
	s := v.String()
	L.Push(lua.LString(s))
	return 1
}

// === User Data Type: Runtime ===============================================

// Global Lua type: runtime. This connects the Lua scripting sub-sytem
// to the DSL's runtime environment
const luaDSLRuntimeTypeName = "runtime"

/*
Lua UserData type for the DSL's interpreter runtime environment.
The scripting sub-system has access to variables of the DSL (and therefore
access to scopes and memory-frames of the runtime environment).

Example (Lua):

    rt = runtime.current                -- find the host-DSL runtime environment
    x = rt.connect_variable("x")        -- create a varref (UserData) for tag 'x'
    print(x)                            -- print a representation for 'x'
    print(x:value())                    -- print the value of 'x'

This will support other host-DSL commands in the future.
*/
type DSLRuntimeEnv struct {
	rt *runtime.Runtime
}

var runtimeMethods = map[string]lua.LGFunction{
	"connect_variable": runtimeConnectVar,
}

// Registers runtime type to given Scripting.
func (lscript *Scripting) registerDSLRuntimeEnvType() {
	grt := &DSLRuntimeEnv{
		rt: lscript.runtime,
	}
	udata := lscript.NewUserData()
	udata.Value = grt
	mt := lscript.NewTypeMetatable(luaDSLRuntimeTypeName)
	lscript.SetMetatable(udata, mt)
	lscript.SetGlobal("runtime", mt)
	lscript.SetField(mt, "current", udata)
	lscript.SetField(mt, "__index", lscript.SetFuncs(lscript.NewTable(), runtimeMethods))
}

func runtimeConnectVar(L *lua.LState) int {
	varname := L.CheckString(1)
	if lrt := getGlobalDSLRuntimeEnv(L); lrt != nil {
		vref := MakeCanonicalAndResolve(lrt.rt, varname, true)
		T.Debugf("variable reference = %v", vref)
		vudata := newVarRefUserData(L, vref)
		L.Push(vudata)
		return 1
	}
	return 0
}

func getGlobalDSLRuntimeEnv(L *lua.LState) *DSLRuntimeEnv {
	mt := L.GetTypeMetatable(luaDSLRuntimeTypeName)
	if mt != nil {
		udata := L.GetField(mt, "current").(*lua.LUserData)
		lrt := udata.Value.(*DSLRuntimeEnv)
		return lrt
	} else {
		T.P("script", "lua").Error("host language runtime env. not found")
		T.Error("Did you pre-load UserData-type 'runtime'?")
		return nil
	}
}
