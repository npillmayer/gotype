package dagstack

import (
	"io/ioutil"
	"testing"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
)

func traceOn() {
	T.SetLevel(tracing.LevelDebug)
}

func TestRoot(t *testing.T) {
	r := DAGStackRoot("G")
	if r.name != "G" {
		t.Fail()
	}
}

func TestCreateStack(t *testing.T) {
	r := DAGStackRoot("G")
	s := NewDAGStack(r)
	if s.Pop() != nil {
		t.Fail()
	}
}

func TestPush(t *testing.T) {
	r := DAGStackRoot("G")
	s := NewDAGStack(r)
	A := pseudosym("A")
	s.Push(1, A)
	if s.Size() != 1 {
		t.Fail()
	}
}

func TestPush2(t *testing.T) {
	A := pseudosym("A")
	r := DAGStackRoot("G")
	s1 := NewDAGStack(r)
	s2 := NewDAGStack(r)
	s1.Push(1, A)
	if s1.Size() != 1 || s2.Size() != 0 {
		t.Fail()
	}
	s2.Push(1, A)
	if s2.Size() != 1 {
		t.Fail()
	}
	if s1.tos != s2.tos {
		t.Fail()
	}
}

func TestPush3(t *testing.T) {
	traceOn()
	A, B := pseudosym("A"), pseudosym("B")
	r := DAGStackRoot("G")
	s1 := NewDAGStack(r)
	s1.Push(1, A)
	s2 := NewDAGStack(r)
	s2.Push(2, B)
	s2.Push(3, B)
	if s1.Size() != 1 || s2.Size() != 2 {
		t.Fail()
	}
	if s1.tos == s2.tos {
		t.Fail()
	}
}

func TestPush4(t *testing.T) {
	A, B := pseudosym("A"), pseudosym("B")
	r := DAGStackRoot("G")
	s1 := NewDAGStack(r)
	s2 := NewDAGStack(r)
	s1.Push(1, A)
	if s1.Size() != 1 || s2.Size() != 0 {
		t.Fail()
	}
	s2.Push(1, A)
	s2.Push(2, B)
	if s2.Size() != 2 {
		t.Fail()
	}
}

func TestPush5(t *testing.T) {
	traceOn()
	A, B := pseudosym("A"), pseudosym("B")
	r := DAGStackRoot("G")
	s1 := NewDAGStack(r)
	s1.Push(1, A)
	s1.Push(3, A)
	s2 := NewDAGStack(r)
	s2.Push(2, B)
	s2.Push(3, A)
	if s1.Size() != 2 || s2.Size() != 2 {
		t.Fail()
	}
	if s1.tos != s2.tos {
		t.Fail()
	}
	tmp, _ := ioutil.TempFile("", "stack_")
	T.Infof("writing DOT to %s", tmp.Name())
	Stack2Dot(r, tmp)
}
