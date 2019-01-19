package pdf

import (
	"github.com/npillmayer/gotype/backend/print/pdf/pdfapi"
	"github.com/npillmayer/gotype/core/dimen"
)

// https://helpx.adobe.com/de/indesign/using/preparing-pdfs-service-providers.html#

//PDF Backend
//
// https://godoc.org/bitbucket.org/zombiezen/gopdf/pdf#pkg-files
//
// https://godoc.org/github.com/jung-kurt/gofpdf
// https://github.com/jung-kurt/gofpdf

// PageNo is a page number
type PageNo int16

type PageStatus int8

const (
	Queued PageStatus = iota
	Printing
	Assembled
	Printed
)

type PageQueue struct {
	//sync.RWMutex
	pages chan *Page // pages in print
	//buzzer       chan PageNo  // signal channel
}

func newQueue() *PageQueue {
	q := &PageQueue{}
	q.pages = make(chan *Page, 10)
	return q
}

func (pq *PageQueue) enqueue(page *Page) {
	page.status = Queued
	queued := true
	select {
	case pq.pages <- page:
	default:
		queued = false
	}
	if !queued {
		go func(p *Page) {
		}(page)
	}
}

type PdfPrinter struct {
	doc       *pdfapi.Document
	Proofing  bool   // are we in proof mode?
	Colormode bool   // color or b/w ?
	Pagecount PageNo // # of pages already completed
	papersize pdfapi.Point
	q         *PageQueue
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
	pp.doc = pdfapi.NewDocument()
	pp.papersize = dpt2upt(papersize)
	pp.q = newQueue()
	return pp
}

func (pr *PdfPrinter) NewPage(pageGeom dimen.Rect) *Page {
	page := &Page{}
	page.pageGeom.Min = dpt2upt(pageGeom.TopL)
	page.pageGeom.Max = dpt2upt(pageGeom.BotR)
	pr.q.enqueue(page)
	return page
}

// TODO
type BoxTree struct{}

type Page struct {
	// https://www.prepressure.com/pdf/basics/page-boxes
	pageGeom pdfapi.Rectangle // position and size on paper
	status   PageStatus       // print status
	pageNo   PageNo           // page number
	content  *BoxTree         // page contents
}

// ----------------------------------------------------------------------

func upt2dpt(p pdfapi.Point) dimen.Point {
	return dimen.Point{
		X: unit2dimen(p.X),
		Y: unit2dimen(p.Y),
	}
}

func dpt2upt(p dimen.Point) pdfapi.Point {
	return pdfapi.Point{
		X: dimen2unit(p.X),
		Y: dimen2unit(p.Y),
	}
}

func dimen2unit(d dimen.Dimen) pdfapi.Unit {
	return pdfapi.Unit(d) / pdfapi.Unit(dimen.BP)
}

func unit2dimen(u pdfapi.Unit) dimen.Dimen {
	return dimen.Dimen(float64(u) * float64(dimen.BP))
}
