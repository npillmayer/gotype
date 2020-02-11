package fp_test

import (
	"testing"

	"github.com/npillmayer/gotype/syntax/terex/fp"
)

func TestN(t *testing.T) {
	a := make([]int64, 0, 10)
	for n, N := fp.N().First(); n < 10; n = N.Next() {
		t.Logf("n=%d", n)
		a = append(a, n)
	}
	t.Logf("a=%v", a)
	if a[0] != 0 || a[9] != 9 || len(a) != 10 {
		t.Errorf("Generating 10 n in N failed")
	}
}

func TestR(t *testing.T) {
	a := make([]float64, 0, 10)
	for x, R := fp.R().First(); x < 10; x = R.Next() {
		t.Logf("x=%.1f", x)
		a = append(a, x)
	}
	t.Logf("a=%v", a)
	if a[0] != 0.0 || a[9] != 9.0 || len(a) != 10 {
		t.Errorf("Generating 10 x in R failed")
	}
}

func TestF(t *testing.T) {
	seq := fp.N().Where(fp.EvenN()).Map(fp.SquareN())
	t.Logf("seq=%v", seq)
	n, N := fp.N().Where(fp.EvenN()).Map(fp.SquareN()).First()
	t.Logf("n=%d, N.n=%d", n, N.N())
	n = N.Next()
	t.Logf("n=%d, N.n=%d", n, N.N())
	n = N.Next()
	t.Logf("n=%d, N.n=%d", n, N.N())
	n = N.Next()
	t.Logf("n=%d, N.n=%d", n, N.N())
	t.Fail()
}

func TestIntFilter(t *testing.T) {
	a := make([]int64, 0, 10)
	for n, N := fp.N().Where(fp.EvenN()).First(); n < 20; n = N.Next() {
		//t.Logf("n=%d", n)
		a = append(a, n)
	}
	if a[0] != 0 || a[9] != 18 || len(a) != 10 {
		t.Errorf("Generating 10 even(n) failed")
	}
	t.Logf("a=%v", a)
}

func TestIntMapper(t *testing.T) {
	a := make([]int64, 0, 10)
	i := 0
	for n, N := fp.N().Map(fp.SquareN()).First(); i < 10; n = N.Next() {
		i++
		//t.Logf("n=%d", n)
		a = append(a, n)
	}
	if a[0] != 0 || a[9] != 81 || len(a) != 10 {
		t.Errorf("Generating 10 squares failed")
	}
	t.Logf("a=%v", a)
}
