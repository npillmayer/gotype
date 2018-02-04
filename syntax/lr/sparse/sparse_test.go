package sparse

import "testing"

func TestSparse1(t *testing.T) {
	m := NewIntMatrix(10, 11, -1)
	m.Set(5, 5, 77)
	if m.Value(5, 5) != 77 {
		t.Fail()
	}
	if m.Value(0, 0) != -1 {
		t.Fail()
	}
	if m.ValueCount() != 1 {
		t.Fail()
	}
}
