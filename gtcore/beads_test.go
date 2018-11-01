package gtcore

import (
	"testing"

	p "github.com/npillmayer/gotype/gtcore/parameters"
)

func TestDimen(t *testing.T) {
	if p.BP.String() != "65536sp" {
		t.Fail()
	}
}

func TestBeadFactory(t *testing.T) {
	kern := NewBead(NTKern).(*Kern)
	kern.Width = p.BP
	t.Logf("kern = %s\n", kern.String())
	if kern.W() != p.BP {
		t.Fail()
	}
	glue := NewBead(NTGlue).(*Glue)
	glue.Width, glue.MaxWidth, glue.MinWidth = p.BP, p.BP, p.BP
	t.Logf("glue = %s\n", glue.String())
	if glue.W() != p.BP {
		t.Fail()
	}
}

func TestBeadChain(t *testing.T) {
	beading := NewBeadChain()
	beading.AppendBead(NewBead(NTKern)).AppendBead(NewBead(NTGlue))
	t.Logf("beading = %s\n", beading.String())
	if beading.Length() != 2 {
		t.Fail()
	}
}
