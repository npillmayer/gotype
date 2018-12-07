package textshaping

import (
	"fmt"
	"testing"

	"github.com/npillmayer/gotype/gtcore/font"
	"github.com/npillmayer/gotype/gtlocate"
)

func TestOpenHarfbuzzFontCreation(t *testing.T) {
	var hbfont uintptr
	fontpath := gtlocate.FileResource("GentiumPlus-R.ttf", "font")
	if f, err := font.LoadOpenTypeFont(fontpath); err == nil {
		if tc, err2 := f.PrepareCase(12.0); err2 == nil {
			hbfont = makeHBFont(tc)
		}
	}
	if hbfont == 0 {
		t.Fail()
	}
}

func TestHarfbuzzShape(t *testing.T) {
	var hb *Harfbuzz
	fontpath := gtlocate.FileResource("GentiumPlus-R.ttf", "font")
	if f, err := font.LoadOpenTypeFont(fontpath); err == nil {
		if tc, err2 := f.PrepareCase(12.0); err2 == nil {
			hb = NewHarfbuzz()
			seq := hb.Shape("Wäffle", tc)
			if seq == nil {
				t.Fail()
			}
			fmt.Printf("Input is \"Wäffle\",\nHarfbuzz returns %d glyphs: %s\n",
				seq.GlyphCount(), hb.GlyphSequenceString(tc, seq))
			if seq.GlyphCount() != 4 {
				t.Fail()
			}
		}
	}
	if hb == nil {
		t.Fail()
	}
}

func TestHarfbuzzShapeResult(t *testing.T) {
	var seq GlyphSequence
	fontpath := gtlocate.FileResource("GentiumPlus-R.ttf", "font")
	if f, err := font.LoadOpenTypeFont(fontpath); err == nil {
		if tc, err2 := f.PrepareCase(12.0); err2 == nil {
			if hb := NewHarfbuzz(); hb != nil {
				if seq = hb.Shape("Fifig", tc); seq != nil {
					fmt.Printf("Input is \"Fifig\",\nHarfbuzz returns %d glyphs: %s\n",
						seq.GlyphCount(), hb.GlyphSequenceString(tc, seq))
					cnt := seq.GlyphCount()
					for i := 0; i < cnt; i++ {
						gi := seq.GetGlyphInfoAt(i)
						fmt.Printf("glyph info #%d/%d: x-advance %.2f\n", i, gi.Cluster(), gi.XAdvance())
					}
				}
			}
		}
	}
	if seq == nil {
		t.Fail()
	}
}
