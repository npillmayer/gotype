package font

import (
	"fmt"
	"testing"

	"github.com/npillmayer/gotype/gtlocate"
)

func TestOpenOpenTypeLoading(t *testing.T) {
	fontpath := gtlocate.FileResource("GentiumPlus-R.ttf", "font")
	f, err := LoadOpenTypeFont(fontpath)
	if err == nil {
		fmt.Printf("loaded font [%s] from \"%s\"\n", f.Fontname, fontpath)
	} else {
		t.Fail()
	}
}

func TestOpenOpenTypeCaseCreation(t *testing.T) {
	fontpath := gtlocate.FileResource("GentiumPlus-R.ttf", "font")
	f, err := LoadOpenTypeFont(fontpath)
	if err != nil {
		t.Fail()
	}
	tc, err2 := f.PrepareCase(12.0)
	if err2 != nil {
		fmt.Printf("cannot create OT face for [%s]\n", f.Fontname)
		t.Fail()
	}
	metrics := tc.font.Metrics()
	fmt.Printf("interline spacing for [%s]@%.1fpt is %s\n", f.Fontname, tc.size, metrics.Height)
}
