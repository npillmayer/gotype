package khipus

import (
	"strings"
	"testing"

	p "github.com/npillmayer/gotype/gtcore/parameters"
)

func TestDimen(t *testing.T) {
	if p.BP.String() != "65536sp" {
		t.Fail()
	}
}

func TestKhipu(t *testing.T) {
	khipu := NewKhipu()
	khipu.AppendKnot(NewKnot(KTKern)).AppendKnot(NewKnot(KTGlue))
	khipu.AppendKnot(NewWordBox("Hello"))
	t.Logf("khipu = %s\n", khipu.String())
	if khipu.Length() != 3 {
		t.Fail()
	}
}

func TestBreaking1(t *testing.T) {
	regs := p.NewTypesettingRegisters()
	regs.Push(p.P_MINHYPHENLENGTH, 3)
	KnotEncode(strings.NewReader("Hello world!"), nil, regs)
}

/*
func TestGraphemeIterator1(t *testing.T) {
	_ = graphemes("Hello")
}

func TestGraphemeIterator2(t *testing.T) {
	iterateOverGraphemes("Héllô 世界")
}

func TestWordIterator1(t *testing.T) {
	iterateOverWords(strings.NewReader("Héllô 世界"))
}

func TestUAX14(t *testing.T) {
	UAX14LineWrap("Title(\"Héllô\") 世界", nil)
}
*/
