package polygon

import (
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"

	a "github.com/npillmayer/gotype/core/arithmetic"
)

func TestBuilder(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	pg := NullPolygon().Knot(a.Pr(0, 0)).Knot(a.Pr(1, 3)).Knot(a.Pr(3, 0)).Cycle()
	L().Infof("pg = %s", AsString(pg))
	if pg.N() != 3 {
		t.Fail()
	}
}

func TestBox(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	box := Box(a.Pr(0, 5), a.Pr(4, 1))
	L().Infof("box = %s", AsString(box))
	if box.N() != 4 {
		t.Fail()
	}
}
