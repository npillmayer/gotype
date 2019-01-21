package pdf

import (
	"image/color"
	"os"
	"testing"

	"github.com/npillmayer/gotype/backend/print/pdf/pdfapi"
)

func TestSimple(t *testing.T) {
	pdffile, err := os.Create("apitest.pdf")
	if err != nil {
		t.Error(err)
	}
	t.Log("Writing PDF to: apitest.pdf")
	defer pdffile.Close()
	helvetica := pdfapi.NewInternalFont(pdfapi.Helvetica)
	api := pdfapi.NewDocument()
	canvas := api.NewPage(pdfapi.Unit(300), pdfapi.Unit(300))
	canvas.DrawLine(pdfapi.Point{0, 0}, pdfapi.Point{300, 300})
	p := &pdfapi.Path{}
	r := pdfapi.Rectangle{pdfapi.Point{50, 100}, pdfapi.Point{250, 200}}
	center := pdfapi.Point{150, 150}
	p.Rectangle(r, 10)
	canvas.PushState()
	canvas.SetStrokeColor(color.Black)
	transf := pdfapi.Identity().Shifted(pdfapi.Point{-center.X, -center.Y}).Rotated(
		pdfapi.Deg2Rad(30)).Shifted(center)
	t.Logf("T = %v", transf)
	canvas.Transform(transf)
	canvas.FillStroke(p)
	canvas.PopState()
	txt := &pdfapi.Text{}
	txt.MoveCursor(pdfapi.Point{20, 200})
	txt.SetFont(helvetica, 24)
	txt.AddGlyphs("Hello world!")
	canvas.DrawText(txt)
	canvas.Close()
	api.Assemble(canvas)
	api.Encode(pdffile)
}
