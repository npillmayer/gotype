package layout

import (
	"testing"
)

func TestIntervals(t *testing.T) {
	// teardown := gotestingadapter.RedirectTracing(t)
	// defer teardown()
	rl := runlength{}
	rl = append(rl, intv{1, 2})
	rl = append(rl, intv{4, 2})
	rl = append(rl, intv{8, 4})
	cond := rl.Condense()
	if len(cond) != 3 || cond[2] != 6 {
		t.Logf("runlength = %s", rl.String())
		t.Errorf("condensed = %v, should be [1 3 6]", cond)
	}
}

func TestIntervalTranslation(t *testing.T) {
	// teardown := gotestingadapter.RedirectTracing(t)
	// defer teardown()
	rl := runlength{}
	rl = append(rl, intv{2, 2})
	rl = append(rl, intv{5, 2})
	rl = append(rl, intv{9, 4})
	t.Logf("runlength = %s", rl.String())
	indexes := []uint32{0, 2, 4, 5, 7, 10, 15}
	expected := []uint32{0, 0, 3, 0, 5, 1, 10}
	for n, inx := range indexes {
		anon, ino, j := rl.Translate(inx)
		t.Logf("translate(%d) = %v, intv=%d, inx=%d", inx, anon, ino, j)
		if expected[n] != j {
			t.Errorf("local index for Translate(%d) is %d, should be %d", inx, j, expected[n])
		}
	}
	t.Fail()
}

func TestDisplayMode(t *testing.T) {
	// teardown := gotestingadapter.RedirectTracing(t)
	// defer teardown()
	disp := InlineMode | FlowMode
	expected := "FlowMode InlineMode"
	if disp.FullString() != expected {
		t.Errorf("display = '%s', should be '%s'", disp.FullString(), expected)
	}
	if !disp.Contains(FlowMode) {
		t.Errorf("display expected to contain 'flow', is missing")
	}
	other := BlockMode | ContentsMode | FlowMode
	if !disp.Overlaps(other) {
		t.Errorf("displays should both contain 'flow', don't")
	}
}
