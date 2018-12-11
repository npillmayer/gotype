package khipus

import (
	"strings"
	"testing"

	p "github.com/npillmayer/gotype/gtcore/parameters"
)

func TestDimen(t *testing.T) {
	if p.BP.String() != "65536sp" {
		t.Error("a big point BP should be 65536 scaled points SP")
	}
}

func TestKhipu(t *testing.T) {
	khipu := NewKhipu()
	khipu.AppendKnot(NewKnot(KTKern)).AppendKnot(NewKnot(KTGlue))
	khipu.AppendKnot(NewTextBox("Hello"))
	t.Logf("khipu = %s\n", khipu.String())
	if khipu.Length() != 3 {
		t.Errorf("Length of khipu should be 3")
	}
}

func TestBreaking1(t *testing.T) {
	regs := p.NewTypesettingRegisters()
	regs.Push(p.P_MINHYPHENLENGTH, 3)
	khipu := KnotEncode(strings.NewReader("Hello world!"), nil, regs)
	if khipu.Length() != 9 {
		t.Errorf("khipu length is %d, should be 9", khipu.Length())
	}
}

/*
func TestUAX14(t *testing.T) {
	UAX14LineWrap("Title(\"Héllô\") 世界", nil)
}
*/
