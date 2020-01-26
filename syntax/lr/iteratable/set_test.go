package iteratable

import "testing"

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
	S.Iterate(Once)
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
	S.Iterate(Once)
	out := ""
	i := 0
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
	S.Iterate(Exhaustively)
	out := ""
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
	S.Add("1")
	S.Add("5")
	t.Logf("S.items=%v", S.items)
	S.Iterate(Exhaustively)
	out := ""
	i := 0
	for S.Next() {
		if i == 1 {
			S.Add("7")
			S.Add("9")
			S.Remove("1")
		}
		if i == 20 {
		}
		i++
		el := S.Take()
		t.Logf("el=%v", el)
		out = out + el.(string)
	}
	if out != "5297" {
		t.Errorf("output after iteration should be '5297', is %s", out)
	}
}
