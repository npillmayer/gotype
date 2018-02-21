package sparse

import (
	"testing"
)

func TestSparse1(t *testing.T) {
	m := NewIntMatrix(10, 11, -1)
	m.Set(5, 5, 77)
	if m.Value(5, 5) != 77 {
		t.Errorf("m(5,5) should be [77, ], is [%d, ]", m.Value(5, 5))
		t.Fail()
	}
	if m.Value(0, 0) != -1 {
		t.Error("m(0,0) should be empty")
		t.Fail()
	}
	if m.ValueCount() != 1 {
		t.Fail()
	}
}

func TestSparse2(t *testing.T) {
	m := NewIntMatrix(10, 10, -1)
	m.Add(5, 5, 77)
	m.Add(5, 5, 88)
	if m.Value(5, 5) != 77 {
		t.Errorf("m(5,5) should be [77, ], is [%d, ]", m.Value(5, 5))
		t.Fail()
	}
	if a, b := m.Values(5, 5); a != 77 || b != 88 {
		t.Errorf("m(5,5) should be [77,88], is [%d,%d]", a, b)
		t.Fail()
	}
	if m.Value(0, 0) != -1 {
		t.Error("m(0,0) should be empty")
		t.Fail()
	}
	if m.ValueCount() != 1 {
		t.Fail()
	}
}
