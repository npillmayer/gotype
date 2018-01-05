package gtcore

import (
	"testing"
)

func TestDimen(t *testing.T) {
	if BP.String() != "65536sp" {
		t.Fail()
	}
}

func TestNodeFactory(t *testing.T) {
	kern := NewNode(NTKern).(*Kern)
	kern.Width = BP
	t.Logf("kern = %s\n", kern.String())
	if kern.W() != BP {
		t.Fail()
	}
	glue := NewNode(NTGlue).(*Glue)
	glue.Width, glue.MaxWidth, glue.MinWidth = BP, BP, BP
	t.Logf("glue = %s\n", glue.String())
	if glue.W() != BP {
		t.Fail()
	}
}

func TestNodelist(t *testing.T) {
	nodelist := NewNodelist()
	nodelist.AppendNode(NewNode(NTKern)).AppendNode(NewNode(NTGlue))
	t.Logf("nodelist = %s\n", nodelist.String())
	if nodelist.Length() != 2 {
		t.Fail()
	}
}
