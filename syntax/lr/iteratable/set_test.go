package iteratable

import (
	"strings"
	"testing"
)

type T interface {
	Value() int
}

func TestSetCreation(t *testing.T) {
	S := NewSet(-1)
	if S.items == nil {
		t.Errorf("set creation failed, S.items == nil")
	}
	R := NewSet(-1)
	R.Add("1")
	if R.items == nil {
		t.Errorf("set creation failed, R.items == nil")
	}
	if R.Size() != 1 {
		t.Errorf("size of R expected to be 1, is %d", R.Size())
	}
	R.Add("2")
	if R.Size() != 2 {
		t.Errorf("size of R expected to be 2, is %d", R.Size())
	}
}

func TestSetEquals(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("5")
	R := NewSet(-1)
	R.Add("1")
	R.Add("5")
	if !S.Equals(R) {
		t.Errorf("S and R expected to be equal")
	}
	R.Add("9")
	if S.Equals(R) {
		t.Errorf("S and modified R expected to be not equal")
	}
}

func TestSetAddAndRemove(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("1")
	if S.Size() != 2 {
		t.Errorf("size of S expected to be 2, is %d", S.Size())
	}
	S.Remove("2")
	t.Logf("S.items=%v", S.items)
	if S.Size() != 1 {
		t.Errorf("size of S now expected to be 1, is %d", S.Size())
	}
	if !S.Contains("1") {
		t.Errorf("S expected to contain '1', does not")
	}
}

func TestSetUnion(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("1")
	R := NewSet(-1)
	R.Add("1")
	R.Add("3")
	S.Union(R)
	t.Logf("S.items=%v", S.items)
	if R.Size() != 2 {
		t.Errorf("R should contain 2 elements, has %d", R.Size())
	}
	if S.Size() != 3 {
		t.Errorf("size of S expected to be 3, is %d", S.Size())
	}
	if !S.Contains("1") {
		t.Errorf("S expected to contain '1', does not")
	}
}

func TestSetIntersection(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("1")
	t.Logf("size of S is %d", S.Size())
	R := NewSet(-1)
	R.Add("1")
	R.Add("3")
	t.Logf("size of R is %d", R.Size())
	S.Intersection(R)
	t.Logf("S.items=%v", S.items)
	if S.Size() != 1 {
		t.Errorf("size of S expected to be 1, is %d", S.Size())
	}
	if !S.Contains("1") {
		t.Errorf("S expected to contain '1', does not")
	}
}

func TestSubset(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("1")
	S.Add("5")
	t.Logf("S.items=%v", S.items)
	S.Subset(func(el interface{}) bool {
		return el.(string) == "5"
	})
	t.Logf("S.items=%v", S.items)
	if S.Size() != 1 {
		t.Errorf("size of S expected to be 1, is %d", S.Size())
	}
	if !S.Contains("5") {
		t.Errorf("S expected to contain '5', does not")
	}
}

func TestSetIteration1(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("1")
	S.Add("5")
	t.Logf("S.items=%v", S.items)
	S.IterateOnce()
	out := ""
	for S.Next() {
		el := S.Item()
		t.Logf("el=%v", el)
		out = out + el.(string)
	}
	if out != "125" {
		t.Errorf("output after iteration should be '125', is %s", out)
	}
}

func TestSetIteration2(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("1")
	S.Add("5")
	t.Logf("S.items=%v", S.items)
	out := ""
	i := 0
	S.IterateOnce()
	for S.Next() {
		if i == 1 {
			S.Add("7")
		}
		el := S.Item()
		t.Logf("el=%v", el)
		out = out + el.(string)
		i++
	}
	if out != "1257" {
		t.Errorf("output after iteration should be '1257', is %s", out)
	}
}

func TestSetIteration3(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("1")
	S.Add("5")
	t.Logf("S.items=%v", S.items)
	out := ""
	S.Exhaust()
	for S.Next() {
		el := S.Take()
		t.Logf("el=%v", el)
		out = out + el.(string)
	}
	if out != "521" {
		t.Errorf("output after iteration should be '521', is %s", out)
	}
}

func TestSetIteration4(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("3")
	t.Logf("S.items=%v", S.items)
	out := ""
	i := 0
	S.Exhaust() // now elements-list is consulted backwards
	for S.Next() {
		if i == 1 { // 3 has been consumed, now is { 1, 2@ }
			t.Logf("S=%v", S.items)
			S.Add("7")
			S.Add("9") // now should be { 9, 7, 1, 2@ }
			t.Logf("S=%v", S.items)
			S.Remove("1") // now should be { 9, 7, 2@ }
			t.Logf("S=%v", S.items)
		}
		i++
		el := S.Take()
		t.Logf("took el=%v", el)
		out = out + el.(string)
	}
	if out != "3279" {
		t.Errorf("output after iteration should be '3279', is %s", out)
	}
}

func TestSetIteration5(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("2")
	S.Add("3")
	t.Logf("S.items=%v", S.items)
	i := 0
	S.Exhaust()
	for S.Next() && !S.Stagnates() {
		i++
		if i > 10 {
			t.Errorf("Stagnation of S not recognized")
			break
		}
		el := S.Take()
		S.Add(el)
	}
	t.Logf("S stagnated, has %d elements, detected after %d iterations", S.Size(), i)
	if S.Size() != 3 || i != 3 {
		t.Error("Detection of stagnation did not work correctly")
	}
}

func TestSetIteration6(t *testing.T) {
	S := NewSet(-1)
	S.Add("1")
	S.Add("-1")
	S.Add("4")
	S.Add("-2")
	t.Logf("S.items=%v", S.items)
	S.IterateOnce()
	out := ""
	for S.Next() {
		el := S.Item().(string)
		t.Logf("el=%v", el)
		if strings.HasPrefix(el, "-") {
			S.Remove(el)
			t.Logf("S.items=%v", S.items)
		} else {
			out = out + el
		}
	}
	if out != "14" {
		t.Errorf("output after iteration should be '14', is %s", out)
	}
}

type x struct {
	A int
	B int
}

func TestItemSet(t *testing.T) {
	S := NewSet(-1)
	S.Add(x{1, 3})
	S.Add(x{2, 3})
	S.Add(x{2, 4})
	S.Add(x{1, 3})
	if S.Size() != 3 {
		t.Logf("S=%v", S)
		t.Errorf("S expected to be of size 3, is %d", S.Size())
	}
}

func TestSorting(t *testing.T) {
	S := NewSet(-1)
	S.Add("2")
	S.Add("4")
	S.Add("1")
	S.Add("3")
	t.Logf("S.items=%v", S.items)
	S.Sort(func(x, y interface{}) bool {
		return x.(string) < y.(string)
	})
	t.Logf("S.items=%v", S.items)
	if S.items[0] != "1" || S.items[3] != "4" {
		t.Errorf("Set should be lexicographically sorted by now, isn't")
	}
}
