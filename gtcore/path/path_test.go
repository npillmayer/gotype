package path

import (
	"fmt"
	"math"
	"testing"

	"github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtcore/config"
	"github.com/shopspring/decimal"
)

var zero = arithmetic.ConstZero
var one = arithmetic.ConstOne
var two = decimal.New(2, 0)

func testpath() (*Path, SplineControls) {
	path, controls := Nullpath().Knot(P(1, 1)).Curve().Knot(P(2, 2)).
		Curve().Knot(P(3, 1)).End()
	return path.(*Path), controls
}

func TestStart(t *testing.T) {
	config.Initialize()
}

func TestSliceEnlargement(t *testing.T) {
	arr := make([]complex128, 0)
	arr = extendC(arr, 3, 2+1i)
	c := arr[3]
	if c != 2+1i {
		t.Fail()
	}
}

func TestCreatePath(t *testing.T) {
	path, _ := testpath()
	fmt.Printf("path = %s\n", path)
	if path.N() != 3 {
		t.Fail()
	}
}

func TestPadding(t *testing.T) {
	path, _ := testpath()
	path.cycle = true
	if path.Z(1) != path.Z(path.N()+1) {
		t.Fail()
	}
}

func TestSetTension(t *testing.T) {
	path, _ := Nullpath().Knot(P(1, 1)).TensionCurve(one, two).Cycle()
	if path.PostTension(0) < 0.99 {
		t.Fail()
	}
	if path.PreTension(1) < 1.99 {
		t.Fail()
	}
}

func TestDir(t *testing.T) {
	path, _ := Nullpath().DirKnot(P(1, 1), P(1, 0)).End()
	fmt.Printf("dir(0) = %g\n", path.PostDir(0))
	if angle(path.PostDir(0)) > 0.01 {
		t.Fail()
	}
}

func TestCurl(t *testing.T) {
	path, _ := Nullpath().Knot(P(1, 1)).Line().Cycle()
	fmt.Printf("curl(0) = %g\n", path.PostCurl(0))
	if path.PostCurl(0) > 0.09 {
		t.Fail()
	}
}

func TestDelta(t *testing.T) {
	path, _ := testpath()
	delta1 := delta(path, 1)
	fmt.Printf("delta [1->2] = %g\n", delta1)
	if delta1 != 1-1i {
		t.Fail()
	}
}

func TestD(t *testing.T) {
	path, _ := testpath()
	d := d(path, 2)
	fmt.Printf("d [2->3] = %g\n", d)
	if d < 1.9 {
		t.Fail()
	}
}

func TestPsi(t *testing.T) {
	path, _ := testpath()
	psi := psi(path, 1)
	fmt.Printf("psi [1->2] = %g\n", rad2deg(psi)) // -90.0000001
	if math.Abs(rad2deg(psi)+90.0) > 0.01 {
		t.Fail()
	}
}

func TestPsiCycle(t *testing.T) {
	path, _ := testpath()
	path.cycle = true
	psi := psi(path, 2)
	fmt.Printf("psi [2->3] = %g\n", rad2deg(psi)) // -134.9999997
	if math.Abs(rad2deg(psi)+135.0) > 0.01 {
		t.Fail()
	}
}

func TestPsiCyclePadding(t *testing.T) {
	path, _ := testpath()
	path.cycle = true
	psi1 := psi(path, 1)
	psiN1 := psi(path, path.N()+1)
	if math.Abs(math.Abs(psi1)-math.Abs(psiN1)) > 0.0001 {
		t.Fail()
	}
}

func TestOpen(t *testing.T) {
	path, controls := testpath()
	fmt.Println(PathAsString(path, nil))
	controls = FindHobbyControls(path, controls)
	fmt.Println(PathAsString(path, controls))
}

func TestCycle(t *testing.T) {
	p, _ := testpath()
	path, controls := p.Knot(P(2, 0)).Curve().Cycle()
	fmt.Println(PathAsString(path, nil))
	controls = FindHobbyControls(path, controls)
	//fmt.Println(PathAsString(path, controls))
}

func TestSplices(t *testing.T) {
	path, controls := Nullpath().Knot(P(0, 0)).Curve().Knot(P(2, 3)).TensionCurve(N(1.4), N(1.4)).Knot(P(5, 3)).Curve().DirKnot(P(3, -1), P(-1, 0)).Curve().Cycle()
	controls = FindHobbyControls(path, controls)
	fmt.Println(PathAsString(path, controls))
	//p, controls := testpath()
	//path, _ := p.Curve().Cycle()
	//path.analyzeSplicePaths()
}
