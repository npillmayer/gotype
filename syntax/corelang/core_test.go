package corelang

import (
	"fmt"
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	lua "github.com/yuin/gopher-lua"
)

func TestLuaInit(t *testing.T) {
	T.SetLevel(tracing.LevelDebug)
	L := NewScripting(nil)
	if L == nil {
		t.Fail()
	}
}

func TestLuaEval1(t *testing.T) {
	l := NewScripting(nil)
	_, err := l.Eval(`print("lua")`)
	if err != nil {
		t.Fail()
	}
}

func TestLuaEval2(t *testing.T) {
	l := NewScripting(nil)
	l.Register("ping", ping)
	_, err := l.Eval(`print(ping())`)
	if err != nil {
		t.Fail()
	}
}

func TestLuaArgs1(t *testing.T) {
	a := asScriptingArgs(1, "hello")
	T.Debugf("args = %v", a)
	if len(a) != 2 {
		t.Fail()
	}
}

func TestLuaArgs2(t *testing.T) {
	a := asScriptingArgs(1, "hello")
	T.Debugf("args  = %v", a)
	lv := forLua(a)
	T.Debugf("largs = %v", lv)
	if len(lv) != 2 {
		t.Fail()
	}
}

func TestLuaHook(t *testing.T) {
	T.SetLevel(tracing.LevelDebug)
	l := NewScripting(nil)
	l.RegisterHook("hello", ping)
	r, _ := l.CallHook("hello")
	T.Debugf("r = %v", r)
	if r[0] != "ok" {
		t.Fail()
	}
}

func TestLuaHook2(t *testing.T) {
	scripting := NewScripting(nil)
	scripting.RegisterHook("stars", func(L *lua.LState) int {
		L.Push(lua.LString("* * * * *")) // push result
		return 1                         // return value count
	})
	scripting.Eval("print(stars())") // Lua: print(stars())
}

func ExampleScripting_hook() {
	// Call hook from Go
	scripting := NewScripting(nil)
	scripting.RegisterHook("echo", func(L *lua.LState) int { // register closure
		lv := L.Get(-1)                                      // get top of stack
		msg := fmt.Sprintf("echo: %s !", lua.LVAsString(lv)) // process
		L.Push(lua.LString(msg))                             // push result
		return 1                                             // return value count
	})
	r, _ := scripting.CallHook("echo", "hello world") // Lua: echo("hello world")
	fmt.Println(r)
	// Output: [echo: hello world !]
}

func TestUserDataPair(t *testing.T) {
	l := NewScripting(nil)
	l.registerPairType()
	T.Debug("----------------------------------------")
	l.Eval(`p = pair.new{1, 1}`)
	l.Eval(`p:x(4)`)
	l.Eval(`print(p:x())`)
}

func TestUserDataVarRef(t *testing.T) {
	l := NewScripting(runtime.NewRuntimeEnvironment(variables.NewPMMPVarDecl))
	l.registerDSLRuntimeEnvType()
	l.registerVarRefType()
	T.Debug("----------------------------------------")
	l.Eval(`
        a = varref.refer_to("a2r")
        a:value(7)
        print(a:value())
    `)
}

func TestUserDataRuntime(t *testing.T) {
	l := NewScripting(runtime.NewRuntimeEnvironment(variables.NewPMMPVarDecl))
	l.registerDSLRuntimeEnvType()
	l.registerVarRefType()
	T.Debug("----------------------------------------")
	l.Eval(`rt = runtime.current`)
	l.Eval(`x = rt.connect_variable("x")`)
	l.Eval(`print(x)`)
	l.Eval(`print(x:value())`)
	l.Eval(`print(x:isknown())`)
	l.Eval(`x:value(3.14)`)
	l.Eval(`print("x="..x:value())`)
	l.Eval(`print(x:isknown())`)
}
