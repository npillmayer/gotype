package pdf

import (
	"bitbucket.org/zombiezen/gopdf/pdf"
	"github.com/npillmayer/gotype/core/dimen"
)

//PDF Backend
//
// https://godoc.org/bitbucket.org/zombiezen/gopdf/pdf#pkg-files
//
// https://godoc.org/github.com/jung-kurt/gofpdf
// https://github.com/jung-kurt/gofpdf

type PdfPrinter struct {
	doc          *pdf.Document
	Proofing     bool
	Colormode    bool
	Pagecount    int
	papersize    pdf.Point
	scale        float64
	pageSettings pdfPageSettings
}

func Printer(papersize dimen.Point, scale float64) *PdfPrinter {
	if papersize.X <= 0 || papersize.Y <= 0 {
		return nil
	}
	pp := &PdfPrinter{}
	pp.Proofing = true
	pp.Colormode = true
	if scale <= 0 {
		scale = 1
	}
	pp.scale = scale
	pp.doc = pdf.New()
	pp.papersize = dpt2upt(papersize)
	return pp
}

type pdfPageSettings struct {
	pagesize pdf.Point
	position pdf.Point
}

// https://www.prepressure.com/pdf/basics/page-boxes
func SetupPages(geom dimen.Rect) {
	//
}

func upt2dpt(p pdf.Point) dimen.Point {
	return dimen.Point{
		X: unit2dimen(p.X),
		Y: unit2dimen(p.Y),
	}
}

func dpt2upt(p dimen.Point) pdf.Point {
	return pdf.Point{
		X: dimen2unit(p.X),
		Y: dimen2unit(p.Y),
	}
}

func dimen2unit(d dimen.Dimen) pdf.Unit {
	return pdf.Unit(d) / pdf.Unit(dimen.BP)
}

func unit2dimen(u pdf.Unit) dimen.Dimen {
	return dimen.Dimen(float64(u) * float64(dimen.BP))
}
