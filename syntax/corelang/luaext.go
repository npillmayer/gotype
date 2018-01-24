package corelang

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	"github.com/shopspring/decimal"
	lua "github.com/yuin/gopher-lua"
)

// ---------------------------------------------------------------------------

/*
Type Scripting is an opaque data type to provide access to the
scripting subsystem.

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
	hooks map[string]lua.LGFunction
}

// Create a new scripting subsystem. Not thread safe.
func NewScripting() *Scripting {
	luastate := lua.NewState()
	if luastate == nil {
		T.Error("failed to create Lua scripting subsystem")
		return nil
	}
	scr := &Scripting{luastate, nil}
	scr.hooks = make(map[string]lua.LGFunction)
	T.Info("Scripting initialized")
	return scr
}

// --- Scripting Arguments ---------------------------------------------------

// We shield the Lua implementation from our scripting API.
// Parameters and return values are provided as generic values.
type ScriptingArgs []interface{}

func asScriptingArgs(values ...interface{}) ScriptingArgs {
	args := make(ScriptingArgs, len(values))
	for i, a := range values {
		args[i] = a
	}
	return args
}

// Convert a Go type argument list into Lua types
func forLua(values ScriptingArgs) []lua.LValue {
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

// Convert and wrap a Lua return value into ScriptingArgs
func returnScriptingArgs(lv lua.LValue) ScriptingArgs {
	a := make(ScriptingArgs, 1)
	if lua.LVIsFalse(lv) {
		a[0] = nil
	} else {
		switch lv.Type() {
		case lua.LTNumber:
			a[0] = float64(lv.(lua.LNumber))
			a[0] = decimal.NewFromFloat(a[0].(float64))
		case lua.LTString:
			a[0] = lua.LVAsString(lv)
		case lua.LTUserData:
			//a = append(a, lua.LTUserData)
			T.Error("not yet implemented: user data, ignored")
		default:
			T.Error("unexpected scripting value type, ignored")
		}
	}
	return a
}

type ScriptingArgsIterator struct {
	args ScriptingArgs
	inx  int
}

func (scrargs ScriptingArgs) Iterator() *ScriptingArgsIterator {
	return &ScriptingArgsIterator{scrargs, -1}
}

func (it *ScriptingArgsIterator) Next() bool {
	it.inx++
	if it.inx < len(it.args) {
		return true
	} else {
		return false
	}
}

func (it *ScriptingArgsIterator) Value() (interface{}, int) {
	if it.inx < len(it.args) {
		a := it.args[it.inx]
		if a == nil {
			return nil, variables.Undefined
		} else if v, ok := a.(decimal.Decimal); ok {
			return v, variables.NumericType
		} else {
			T.Errorf("not yet implemented: type for %v", a)
		}
	}
	return nil, variables.Undefined
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
func (lscript *Scripting) CallHook(hook string, arguments ...interface{}) (ScriptingArgs, error) {
	args := asScriptingArgs(arguments...)
	err := lscript.CallByParam(lua.P{
		Fn:      lscript.GetGlobal(hook),
		NRet:    1,
		Protect: true,
	}, forLua(args)...)
	var r ScriptingArgs
	if err == nil {
		r = returnScriptingArgs(lscript.Get(-1)) // wrap returned value
		lscript.Pop(1)                           // remove received value from stack
	}
	return r, err
}

/*
Evaluate a Lua statement, given as string. Returns the number of return
arguments on the stack and a possible error condition. The statement must
be a syntactically correct and complete statement (no expressions etc.
accepted).

Eval accepts a variable number of untyped arguments. These are put on the
Lua stack before the statement is executed.
*/
func (lscript *Scripting) Eval(luacmd string, arguments ...interface{}) (int, error) {
	args := asScriptingArgs(arguments...)
	largs := forLua(args)
	for _, larg := range largs {
		lscript.Push(larg)
	}
	var lv lua.LValue
	err := lscript.DoString(luacmd)
	if err != nil {
		T.Errorf("scripting error: %s", err.Error())
	} else {
		lv = lscript.Get(-1)
	}
	T.P("lua", "eval").Debugf("return value on Lua stack = %s", lua.LVAsString(lv))
	return lscript.GetTop(), err
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

// === User Data Type: Pair ==================================================

const luaPairTypeName = "pair"

/*
Lua UserData type for pairs.

Example:

  p = pair.new{2, 5}
  print(p:x())
  p:y(3.14)

*/
type LuaPair struct {
	X lua.LNumber
	Y lua.LNumber
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
	var pair *LuaPair
	xy := L.CheckTable(1)
	if xy.MaxN() != 2 {
		T.P("lua", "udata").Errorf("wrong number of arguments to pair:new(): %d", xy.MaxN())
		pair = &LuaPair{lua.LNumber(0), lua.LNumber(0)}
	} else {
		x := xy.RawGet(lua.LNumber(1))
		y := xy.RawGet(lua.LNumber(2))
		pair = &LuaPair{x.(lua.LNumber), y.(lua.LNumber)}
	}
	//pair := &LuaPair{L.CheckNumber(1), L.CheckNumber(2)}
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
		p.X = L.CheckNumber(2)
		return 0
	} else { // getter
		L.Push(lua.LNumber(p.X))
		return 1
	}
}

// Getter and setter for pair#y
func pairGetSetY(L *lua.LState) int {
	p := checkPair(L)
	if L.GetTop() == 2 { // setter
		p.Y = L.CheckNumber(2)
		return 0
	} else { // getter
		L.Push(lua.LNumber(p.Y))
		return 1
	}
}

// === User Data Type: Variable ==============================================

const luaNumericTypeName = "numeric"

/*
Lua UserData type for numeric variables. A numeric variable may be known
or unknown.
*/
type LuaNumeric struct {
	e *runtime.ExprNode
}

var numericMethods = map[string]lua.LGFunction{
	"value":   numericGetSetValue,
	"isknown": numericIsKnown,
}

// Registers numeric type to given Scripting.
func (lscript *Scripting) registerNumericType() {
	mt := lscript.NewTypeMetatable(luaNumericTypeName)
	lscript.SetGlobal("numeric", mt)
	// static attributes
	lscript.SetField(mt, "new", lscript.NewFunction(newNumeric))
	// methods
	lscript.SetField(mt, "__index", lscript.SetFuncs(lscript.NewTable(), numericMethods))
}

// Constructor
func newNumeric(L *lua.LState) int {
	var n *LuaNumeric
	if L.GetTop() >= 2 {
		n = newLuaNumeric(L.CheckNumber(1))
	} else {
		n = newLuaNumeric(nil)
	}
	ud := L.NewUserData()
	ud.Value = n
	L.SetMetatable(ud, L.GetTypeMetatable(luaNumericTypeName))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Numeric and returns
// this *Numeric.
func checkNumeric(L *lua.LState) *LuaNumeric {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*LuaNumeric); ok {
		return v
	}
	L.ArgError(1, "numeric expected")
	return nil
}
