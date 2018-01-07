package runtime

import (
	"testing"
)

func TestNewSymTab(t *testing.T) {
	symtab := NewSymbolTable(nil)
	if symtab == nil {
		t.Error("no symbol table created")
	}
}

func TestNewSymbol(t *testing.T) {
	symtab := NewSymbolTable(nil)
	sym, _ := symtab.DefineSymbol("new-sym")
	if sym == nil {
		t.Error("no symbol created for table")
	}
}

func TestTwoSymbolsDistinctId(t *testing.T) {
	symtab := NewSymbolTable(nil)
	sym1, _ := symtab.DefineSymbol("new-sym1")
	sym2, _ := symtab.DefineSymbol("new-sym2")
	if sym1 == sym2 {
		t.Error("2 symbols with equal name")
	}
}

func TestResolveSymbol(t *testing.T) {
	symtab := NewSymbolTable(nil)
	sym, _ := symtab.DefineSymbol("new-sym")
	if s := symtab.ResolveSymbol(sym.GetName()); s == nil {
		t.Error("cannot find stored symbol in table")
	}
}

func TestResolveOrDefineSymbol(t *testing.T) {
	symtab := NewSymbolTable(nil)
	sym, _ := symtab.DefineSymbol("new-sym")
	if _, found := symtab.ResolveOrDefineSymbol(sym.GetName()); !found {
		t.Error("cannot find stored symbol in table")
	}
}

func TestDefineSymbol(t *testing.T) {
	symtab := NewSymbolTable(nil)
	sym, _ := symtab.DefineSymbol("new-sym")
	if _, old := symtab.DefineSymbol("new-sym"); old != sym {
		t.Error("symbol should have been replaced")
	}
}

func TestScopeUpsearch(t *testing.T) {
	scopep := NewScope("parent", nil, nil)
	scope := NewScope("current", scopep, nil)
	scopep.DefineSymbol("new-sym")
	if sym, _ := scope.ResolveSymbol("new-sym"); sym != nil {
		t.Logf("found symbol '%s' in parent scope, ok\n", sym.GetName())
	} else {
		t.Fail()
	}
}

func TestAddChild(t *testing.T) {
	sym := NewStdSymbol("new-sym").(*StdSymbol)
	ch1 := NewStdSymbol("child-sym1").(*StdSymbol)
	ch2 := NewStdSymbol("child-sym2").(*StdSymbol)
	sym.AppendChild(ch1)
	sym.AppendChild(ch2)
	if sym.GetFirstChild().(*StdSymbol).GetName() != "child-sym1" {
		t.Fail()
	}
}
