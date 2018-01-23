package polygon

import (
	"testing"

	a "github.com/npillmayer/gotype/gtcore/arithmetic"
)

func TestBuilder(t *testing.T) {
	pg := NullPolygon().Knot(a.Pr(0, 0)).Knot(a.Pr(1, 3)).Knot(a.Pr(3, 0)).Cycle()
	L.Infof("pg = %s", PolygonAsString(pg))
	if pg.N() != 3 {
		t.Fail()
	}
}

func TestBox(t *testing.T) {
	box := Box(a.Pr(0, 5), a.Pr(4, 1))
	L.Infof("box = %s", PolygonAsString(box))
	if box.N() != 4 {
		t.Fail()
	}
}
