package vectorcanvas

import (
	"image/color"
	"os"
	"testing"

	"github.com/npillmayer/gotype/backend/gfx"
	"github.com/npillmayer/gotype/backend/print/pdf/pdfapi"
	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/testadapter"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/path"
)

func TestInit0(t *testing.T) {
	config.Initialize(testadapter.New())
	tracing.GraphicsTracer.SetTraceLevel(tracing.LevelDebug)
}

func TestPdf(t *testing.T) {
	pdfdoc := pdfapi.NewDocument()
	f, err := os.Create("test_path.pdf")
	if err != nil {
		t.Error(err.Error())
	}
	defer f.Close()
	vecc := New(100, 100)
	p, controls := path.Nullpath().Knot(path.P(10, 50)).Curve().Knot(path.P(50, 90)).Curve().
		Knot(path.P(90, 50)).Curve().Knot(path.P(50, 10)).Curve().Cycle()
	controls = path.FindHobbyControls(p, controls)
	t.Logf("path1 = %v", p)
	red := color.RGBA{200, 200, 200, 250}
	vecc.AddContour(gfx.NewDrawablePath(p, controls), 2, color.Black, red)
	page := pdfdoc.NewPage(100, 100)
	vecc.ToPDF(page)
	page.Close()
	pdfdoc.Encode(f)
}
