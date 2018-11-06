package gtcore

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

func TestKnotFactory(t *testing.T) {
	kern := NewKnot(NTKern).(*Kern)
	kern.Width = p.BP
	t.Logf("kern = %s\n", kern.String())
	if kern.W() != p.BP {
		t.Fail()
	}
	glue := NewKnot(NTGlue).(*Glue)
	glue.Width, glue.MaxWidth, glue.MinWidth = p.BP, p.BP, p.BP
	t.Logf("glue = %s\n", glue.String())
	if glue.W() != p.BP {
		t.Fail()
	}
}

func TestKhipu(t *testing.T) {
	khipu := NewKhipu()
	khipu.AppendKnot(NewKnot(NTKern)).AppendKnot(NewKnot(NTGlue))
	t.Logf("khipu = %s\n", khipu.String())
	if khipu.Length() != 2 {
		t.Fail()
	}
}

func TestGraphemeIterator1(t *testing.T) {
	_ = graphemes("Hello")
}

func TestGraphemeIterator2(t *testing.T) {
	iterateOverGraphemes("Héllô 世界")
}

func TestWordIterator1(t *testing.T) {
	iterateOverWords(strings.NewReader("Héllô 世界"))
}
