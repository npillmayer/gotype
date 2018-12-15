package parameters

import "testing"

func TestRegistersCreate(t *testing.T) {
	regs := NewTypesettingRegisters()
	if regs.base[P_LANGUAGE] != "en_EN" {
		t.Fail()
	}
	if lang := regs.S(P_LANGUAGE); lang != "en_EN" {
		t.Fail()
	}
}

func TestRegistersGrouping(t *testing.T) {
	regs := NewTypesettingRegisters()
	regs.Begingroup()
	if regs.grouplevel != 1 {
		t.Fail()
	}
	regs.Push(P_LANGUAGE, "de_DE")
	if regs.groups == nil {
		t.Fail()
	}
	if regs.S(P_LANGUAGE) != "de_DE" {
		t.Fail()
	}
}

func TestRegistersGrouping2(t *testing.T) {
	regs := NewTypesettingRegisters()
	regs.Begingroup()
	regs.Push(P_LANGUAGE, "de_DE")
	regs.Endgroup()
	if regs.S(P_LANGUAGE) != "en_EN" {
		t.Fail()
	}
}
