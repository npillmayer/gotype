package dss

import (
	"io/ioutil"
	"testing"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
)

func traceOn() {
	T.SetLevel(tracing.LevelDebug)
}

func TestRoot(t *testing.T) {
	r := NewRoot("G")
	if r.Name != "G" {
		t.Fail()
	}
}

func TestCreateStack(t *testing.T) {
	r := NewRoot("G")
	s := NewStack(r)
	if s.Pop() != nil {
		t.Fail()
	}
}

func TestPush(t *testing.T) {
	r := NewRoot("G")
	s := NewStack(r)
	A := pseudosym("A")
	s.Push(1, A)
	if s.Size() != 1 {
		t.Fail()
	}
}

func TestPush2(t *testing.T) {
	A := pseudosym("A")
	r := NewRoot("G")
	s1 := NewStack(r)
	s2 := NewStack(r)
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
	A, B := pseudosym("A"), pseudosym("B")
	r := NewRoot("G")
	s1 := NewStack(r)
	s1.Push(1, A)
	s2 := NewStack(r)
	s2.Push(2, B).Push(3, B)
	if s1.Size() != 1 || s2.Size() != 2 {
		t.Fail()
	}
	if s1.tos == s2.tos {
		t.Fail()
	}
}

func TestPush4(t *testing.T) {
	A, B := pseudosym("A"), pseudosym("B")
	r := NewRoot("G")
	s1 := NewStack(r)
	s2 := NewStack(r)
	s1.Push(1, A)
	if s1.Size() != 1 || s2.Size() != 0 {
		t.Fail()
	}
	s2.Push(1, A).Push(2, B)
	if s2.Size() != 2 {
		t.Fail()
	}
}

func TestPush5(t *testing.T) {
	A, B, C := pseudosym("A"), pseudosym("B"), pseudosym("C")
	r := NewRoot("G")
	s1 := NewStack(r)
	s1.Push(1, C).Push(3, A)
	s2 := NewStack(r)
	s2.Push(2, B).Push(3, A)
	if s1.Size() != 2 || s2.Size() != 2 {
		T.Error("unexpected stack sizes: %d, %d", s1.Size(), s2.Size())
		t.Fail()
	}
	if s1.tos != s2.tos {
		t.Fail()
	}
}

func TestUnzip1(t *testing.T) {
	A, B, C := pseudosym("A"), pseudosym("B"), pseudosym("C")
	r := NewRoot("G")
	s1 := NewStack(r)
	s1.Push(1, C).Push(3, A)
	s2 := NewStack(r)
	s2.Push(2, B).Push(3, A).Push(4, C)
	if s1.Size() != 2 || s2.Size() != 3 {
		T.Error("unexpected stack sizes: %d, %d", s1.Size(), s2.Size())
		t.Fail()
	}
	handle := []lr.Symbol{B, A, C}
	path := s2.findHandleBranch(handle)
	T.Debugf("path = %v", path)
	if path == nil || len(path) != 3 || path[1] == nil {
		T.Errorf("path not found or incorrect: ", path)
		t.Fail()
	}
}

func TestSplitOff1(t *testing.T) {
	A, B, C := pseudosym("A"), pseudosym("B"), pseudosym("C")
	r := NewRoot("G")
	s2 := NewStack(r)
	s2.Push(2, B).Push(3, A).Push(4, C)
	handle := []lr.Symbol{A, C}
	path := s2.findHandleBranch(handle)
	T.Debugf("path = %v", path)
	s3 := s2.splitOff(path)
	T.Debugf("TOS of s2 = %v", s2.tos)
	T.Debugf("TOS of s3 = %v", s3.tos)
	if s2.tos == s3.tos {
		T.Error("stacks not split")
		t.Fail()
	}
	if len(r.stacks) != 2 {
		t.Fail()
	}
}

func TestSplitOff2(t *testing.T) {
	traceOn()
	A, B, C := pseudosym("A"), pseudosym("B"), pseudosym("C")
	r := NewRoot("G")
	s1 := NewStack(r)
	s1.Push(1, C).Push(3, A)
	s2 := NewStack(r)
	s2.Push(2, B).Push(3, A).Push(4, C)
	handle := []lr.Symbol{A, C}
	path := s2.findHandleBranch(handle)
	T.Debugf("path = %v", path)
	s3 := s2.splitOff(path)
	T.Debugf("TOS of s2 = %v", s2.tos)
	T.Debugf("TOS of s3 = %v", s3.tos)
	if s2.tos == s3.tos {
		T.Error("stacks not split")
		t.Fail()
	}
	//PrintDSS(r)
	if len(r.stacks) != 3 {
		t.Fail()
	}
	_ = s3.Fork()
	s3.Pop()
	s3.Pop()
	tmp, _ := ioutil.TempFile("", "stack_")
	T.Infof("writing 3rd DOT to %s", tmp.Name())
	DSS2Dot(r, tmp)
	if len(r.stacks) != 5 {
		t.Fail()
	}
}
