package corelang

import (
	"fmt"
	"testing"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	lua "github.com/yuin/gopher-lua"
)

func TestLuaInit(t *testing.T) {
	T.SetLevel(tracing.LevelDebug)
	L := NewScripting()
	if L == nil {
		t.Fail()
	}
}

func TestLuaEval1(t *testing.T) {
	l := NewScripting()
	n, err := l.Eval(`print("lua")`)
	T.Debugf("Lua eval return value = %d", n)
	if n != 0 || err != nil {
		t.Fail()
	}
}

func TestLuaEval2(t *testing.T) {
	l := NewScripting()
	l.Register("ping", ping)
	n, err := l.Eval(`print(ping())`)
	if n != 0 || err != nil {
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
	l := NewScripting()
	l.RegisterHook("hello", ping)
	r, _ := l.CallHook("hello")
	T.Debugf("r = %v", r)
	if r[0] != "ok" {
		t.Fail()
	}
}

func TestLuaHook2(t *testing.T) {
	scripting := NewScripting()
	scripting.RegisterHook("stars", func(L *lua.LState) int {
		L.Push(lua.LString("* * * * *")) // push result
		return 1                         // return value count
	})
	scripting.Eval("print(stars())") // Lua: print(stars())
}

func ExampleScripting_hook() {
	// Call hook from Go
	scripting := NewScripting()
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
	l := NewScripting()
	l.registerPairType()
	T.Debug("----------------------------------------")
	l.Eval(`p = pair.new{1, 1}`)
	l.Eval(`p:x(4)`)
	l.Eval(`print(p:x())`)
}
