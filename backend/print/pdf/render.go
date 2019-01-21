package pdf

import (
	"image/color"

	"github.com/npillmayer/gotype/backend/print/pdf/pdfapi"
	"github.com/npillmayer/gotype/core/dimen"
)

// render will be executed by concurrent render workers.
// TODO Walk the render tree of the pages content
// TODO and create pdfapi calls on it.
func (pr *PdfPrinter) render(page *Page) error {
	T().Debugf("Rendering page [%d]", page.pageNo)
	cv := makeConv(pr.papersize, page.pageGeom, pr.scale) // set up conversion
	renderPrinterMarks(cv, page.pdfcanvas, pr.Proofing)
	return nil
}

// Render printer marks:
// - crop marks
// - color scales    TODO
// - registers       TODO
// - additional info TODO
//
func renderPrinterMarks(cv *conv, canvas *pdfapi.Canvas, proofing bool) {
	canvas.PushState()
	canvas.SetStrokeColor(color.Black)
	canvas.SetLineWidth(0.2) // thin pen for crop marks
	if proofing {
		P := &pdfapi.Path{}
		pagebox := pdfapi.Rectangle{
			cv.toPdfPoint(pdfapi.Point{0, 0}),
			cv.toPdfPoint(cv.pageDim),
		}
		P.Rectangle(pagebox, 0)
		canvas.Stroke(P)
	}
	canvas.DrawLine(cv.upt(-20, 0), cv.upt(-5, 0))
	canvas.DrawLine(cv.upt(0, -20), cv.upt(0, -5))
	canvas.DrawLine(cv.upt(cv.pageDim.X+5, 0), cv.upt(cv.pageDim.X+20, 0))
	canvas.DrawLine(cv.upt(cv.pageDim.X, -20), cv.upt(cv.pageDim.X, -5))
	canvas.DrawLine(cv.upt(-20, cv.pageDim.Y), cv.upt(-5, cv.pageDim.Y))
	canvas.DrawLine(cv.upt(0, cv.pageDim.Y+20), cv.upt(0, cv.pageDim.Y+5))
	canvas.DrawLine(cv.upt(cv.pageDim.X+5, cv.pageDim.Y), cv.upt(cv.pageDim.X+20, cv.pageDim.Y))
	canvas.DrawLine(cv.upt(cv.pageDim.X, cv.pageDim.Y+20), cv.upt(cv.pageDim.X, cv.pageDim.Y+5))
	canvas.PopState()
}

// --- co-ordinates conversion ------------------------------------------

// PDF coordinate systems have their origin in the bottom left corner,
// while TySE pages start in the top left.
// We set up a conversion helper, based on the papersize and page geometry.
// Additionally we'll have to convert between TySE dimensions, which are
// in scaled points, to default PDF units.
//
type conv struct {
	paper      pdfapi.Point
	pageOrigin pdfapi.Point
	pageDim    pdfapi.Point
	scale      pdfapi.Unit
}

// Create a conversion helper.
func makeConv(papersize pdfapi.Point, pagegeom pdfapi.Rectangle, scale float64) *conv {
	c := &conv{}
	c.scale = pdfapi.Unit(scale)
	c.paper = papersize
	c.pageOrigin = pagegeom.Min
	c.pageDim = pdfapi.Point{pagegeom.Dx() * c.scale, pagegeom.Dy() * c.scale}
	T().Debugf("PDF page origin is = (%v,%v)", c.pageDim.X, c.pageDim.Y)
	return c
}

// Scaled points ➝ PDF units
func (cv *conv) ScaledUnit(d dimen.Dimen) pdfapi.Unit {
	u := dimen2unit(d) * cv.scale
	return u
}

// Scaled points ➝ PDF units
func (cv *conv) Point(d dimen.Point) pdfapi.Point {
	x := cv.ScaledUnit(d.X)
	y := cv.ScaledUnit(d.Y)
	return cv.toPdfPoint(pdfapi.Point{x, y})
}

// Scaled points ➝ PDF units
func (cv *conv) Rect(r dimen.Rect) pdfapi.Rectangle {
	rect := pdfapi.Rectangle{
		cv.Point(r.TopL),
		cv.Point(r.BotR),
	}
	return rect
}

// Big point to PDF unit
func (cv *conv) toPdfPoint(p pdfapi.Point) pdfapi.Point {
	x := cv.pageOrigin.X + p.X
	y := cv.paper.Y - cv.pageOrigin.Y - p.Y
	return pdfapi.Point{x, y}
}

// BP point to PDF unit with scale
func (cv *conv) toScaledPdfPoint(p pdfapi.Point) pdfapi.Point {
	x := cv.pageOrigin.X + (p.X * cv.scale)
	y := cv.paper.Y - cv.pageOrigin.Y - (p.Y * cv.scale)
	return pdfapi.Point{x, y}
}

// BP units pair to PDF unit point (without scale)
func (cv *conv) upt(x, y pdfapi.Unit) pdfapi.Point {
	return cv.toPdfPoint(pdfapi.Point{x, y})
}
