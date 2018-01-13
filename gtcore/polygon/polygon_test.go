package polygon

import (
	"testing"

	a "github.com/npillmayer/gotype/gtcore/arithmetic"
)

func TestBuilder(t *testing.T) {
	pg := NullPolygon().AddKnot(a.P(0, 0)).AddKnot(a.P(1, 3)).AddKnot(a.P(3, 0)).Cycle()
	L.Infof("pg = %s", PolygonAsString(pg))
	if pg.N() != 3 {
		t.Fail()
	}
}

func TestBox(t *testing.T) {
	box := Box(a.P(0, 5), a.P(4, 1))
	L.Infof("box = %s", PolygonAsString(box))
	if box.N() != 4 {
		t.Fail()
	}
}
