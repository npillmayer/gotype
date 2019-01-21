package pdf

/*
BSD License

Copyright (c) 2017–18, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"fmt"
	"io"
	"sync"

	"github.com/npillmayer/gotype/backend/print/pdf/pdfapi"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/dimen"
)

// https://helpx.adobe.com/de/indesign/using/preparing-pdfs-service-providers.html#
//
// https://blog.flipsnack.com/5-ways-to-reduce-the-size-of-a-pdf-file/

//PDF Backend
//
// https://godoc.org/bitbucket.org/zombiezen/gopdf/pdf#pkg-files
//
// https://godoc.org/github.com/jung-kurt/gofpdf
// https://github.com/jung-kurt/gofpdf

// contPageNo is a contiguous page number.
type contPageNo int16

type PageStatus int8

// Pages queued for printing may be in these states:
const (
	Queued     PageStatus = iota // queued for printing
	Printing                     // in print, i.e, PDF code is generated
	Assembling                   // PDF for page complete, in doc assembly
	Printed                      // already encoded as part of the document
)

// ----------------------------------------------------------------------

// PdfPrinter is used for outputting a render tree in PDF format.
type PdfPrinter struct {
	doc       *pdfapi.Document // PDF document to assemble
	Proofing  bool             // are we in proof mode?
	Colormode bool             // color or b/w ?
	papersize pdfapi.Point     // paper geometry in PDF points
	scale     float64          // scale factor for pages
	pageQ     *pageQueue       // input queue for pages to print
	assemblyQ *pageheap        // queue to collect and assemble completed pages
	done      chan bool        // signal to abort (regular or interrupt)
	errch     chan error       // channel to collect errors
	err       error            // error to return after printing
	mtx       sync.RWMutex     // guards 'running' and page fields
	running   bool             // is this printer still printing?
	maxPageNo contPageNo       // page number of last page to print
	pagecount contPageNo       // number of pages already completed
}

// Printer creates a new PdfPrinter, given a paper format and scale factor.
func Printer(papersize dimen.Point, scale float64) *PdfPrinter {
	if papersize.X <= 0 || papersize.Y <= 0 {
		return nil
	}
	pr := &PdfPrinter{}
	pr.Proofing = true
	pr.Colormode = true
	if scale <= 0 {
		scale = 1
	}
	pr.scale = scale
	pr.doc = pdfapi.NewDocument()
	pr.papersize = dpt2upt(papersize)
	T().Debugf("Printing on paper %v", pr.papersize)
	pr.pageQ = newQueue(pr)
	pr.assemblyQ = newPageAssemblyQueue()
	pr.done = make(chan bool, 1)
	pr.errch = make(chan error)
	pr.maxPageNo = -1
	return pr
}

// pageComplete signals the page assembler that this page is complete,
// i.e, all the rendering for this page has been done.
func (pr *PdfPrinter) pageComplete(page *Page) {
	pr.mtx.Lock()
	if page != nil {
		page.status = Assembling
	}
	pr.mtx.Unlock()
	queued := false
	pr.mtx.RLock()
	if pr.running {
		queued = true
		T().Debugf("Putting page [%d] to assembly queue", page.pageNo)
		pr.assemblyQ.AppendPage(page)
	}
	pr.mtx.RUnlock()
	if !queued {
		T().Debugf("Printer already stopped, page [%d] not queued for assembly",
			page.pageNo)
	}
}

// pageFailed signals the page assembler that this page failed to render.
// We will create an error page instead and queue it for the assembler.
func (pr *PdfPrinter) pageFailed(page *Page, err error) {
	page.pdfcanvas = pr.doc.NewPage(pr.papersize.X, pr.papersize.Y)
	errtext := &pdfapi.Text{}
	cv := makeConv(pr.papersize, page.pageGeom, pr.scale)
	errtext.MoveCursor(cv.toPdfPoint(pdfapi.Point{20, 20}))
	errtext.SetFont(pdfapi.NewInternalFont(pdfapi.Helvetica), 12)
	errtext.AddGlyphs(fmt.Sprintf("Error rendering page [%d]", page.pageNo))
	if err != nil {
		errtext.MoveCursor(cv.toPdfPoint(pdfapi.Point{20, 40}))
		errtext.AddGlyphs(err.Error())
	}
	page.pdfcanvas.DrawText(errtext)
	pr.pageComplete(page)
}

// IsRunning returns true if the printer is accepting input, false otherwise.
// Printers will be running after calls to Start(...) until either the
// expected number of pages has been printed or Abort(...) has been called.
func (pr *PdfPrinter) IsRunning() bool {
	pr.mtx.RLock()
	defer pr.mtx.RUnlock()
	return pr.running
}

// Status returns the print status of a page in the printer's queue.
func (pr *PdfPrinter) StatusOf(page *Page) PageStatus {
	pr.mtx.RLock()
	defer pr.mtx.RUnlock()
	return page.status
}

// assemblePagesToDocument is intended to be executed by a single assembly
// goroutine. It is only safe to be used with many page creation
// workers and a single page assembly worker.
func assemblePagesToDocument(pr *PdfPrinter, w io.Writer) {
	var nextPage *Page
	for pr.IsRunning() { // wait for next page to assemble until stopped by signal
		select {
		case pageno := <-pr.assemblyQ.pagesDone:
			if pageno <= pr.pagecount+1 {
				// try to get next bunch of suitable pages from the queue
				for pr.assemblyQ.LowestPageNo() <= pr.pagecount+1 {
					nextPage = pr.assemblyQ.NextPage()
					err := pr.appendToDocument(nextPage)
					if err != nil { // stop by signaling error to printer
						go func(e error) { // must be async
							pr.Abort(e)
						}(err)
					}
					pr.mtx.Lock()
					if nextPage.pageNo == pr.pagecount+1 {
						pr.pagecount++
						if pr.maxPageNo >= 0 && pr.pagecount >= pr.maxPageNo {
							go func() { // being careful inside mutex
								pr.done <- true
							}()
						}
					}
					pr.mtx.Unlock()
				}
			}
		case <-pr.done:
			pr.pageQ.Stop() // do not accept any more pages
			pr.mtx.Lock()
			pr.running = false
			pr.mtx.Unlock()
		}
		if !pr.IsRunning() {
			pr.drainAssemblyQueueAndStop()
			pr.serialize(w)
			close(pr.errch) // signal to printing clients
		}
	}
	tracing.EngineTracer.Debugf("Assembly worker stopped")
}

func (pr *PdfPrinter) drainAssemblyQueueAndStop() {
	// We are fetching pending page events from the assembly queue.
	// This is for the regular case only. For an interrupt, we currently
	// have no means to detect how many page workers will return.
	//
	// In the regular case, we received a done event. The select may
	// have returned the done event first, with other page events
	// pending. We skip over the events and then read the pages from
	// the priority queue.
	stopped := false
	for !stopped {
		select { // fetch pending page events, if any
		case <-pr.assemblyQ.pagesDone: // do nothing
		default:
			stopped = true
		}
	}
	// Now there may be still other pages in progress, but we currently
	// ignore them. Drain page assembly priority queue.
	for pr.assemblyQ.Len() > 0 {
		// no other goroutine will fetch from the queue, so this gap is safe
		nextPage := pr.assemblyQ.NextPage() // page with lowest page number
		err := pr.appendToDocument(nextPage)
		if pr.err == nil { // if error unset, overwrite
			pr.err = err
		}
	}
	pr.assemblyQ.Close()
}

// SetMaxPage must be called from clients as soon as they know how many
// pages there will be in the print job.
// Page numbers range from 0 to 32767.
func (pr *PdfPrinter) SetMaxPage(n int) {
	pr.mtx.Lock()
	defer pr.mtx.Unlock()
	pr.maxPageNo = contPageNo(n)
	if pr.pagecount >= pr.maxPageNo { // already printed all pages
		pr.done <- true // is buffered, 1 enough ? TODO
	}
}

// Start starts printing to w. It returns a promise/future.
// Clients will call this promise to synchronously wait for printing
// to finish.
func (pr *PdfPrinter) Start(w io.Writer) func() error {
	pr.mtx.Lock()
	pr.running = true
	pr.mtx.Unlock()
	go func(pp *PdfPrinter) {
		assemblePagesToDocument(pp, w)
	}(pr)
	return func() error {
		// fetch all errors, remember just 1
		for err := range pr.errch {
			if pr.err == nil {
				pr.err = err
			}
		}
		return pr.err
	}
}

// Abort stops the printer. The call will block until the printer stopped.
// err should be nil to signal a successful completion of printing.
func (pr *PdfPrinter) Abort(err error) {
	go func() {
		pr.errch <- err
	}()
	pr.done <- true
}

// PrintPage creates a new page with given page number, dimensions and content.
// The call automatically enqueues it into the printer queue.
//
// The page number is intended to be a part of a contiguous series of page
// numbers. The printer expects pages ranging from 1..n, where n is set
// by SetMaxPage(n). If this constraint is violated, the printer may stall
// waiting for non-existent pages. Page numbers range from 0 to 32767.
func (pr *PdfPrinter) PrintPage(pageno int, pageGeom dimen.Rect, content *RenderTree) *Page {
	page := &Page{}
	page.pageNo = contPageNo(pageno)
	page.pageGeom.Min = dpt2upt(pageGeom.TopL)
	page.pageGeom.Max = dpt2upt(pageGeom.BotR)
	page.content = content
	page.pdfcanvas = pr.doc.NewPage(pr.papersize.X, pr.papersize.Y)
	pr.pageQ.enqueue(page)
	return page
}

// PageCount returns the current number of completed pages.
func (pr *PdfPrinter) PageCount() int {
	pr.mtx.RLock()
	defer pr.mtx.RUnlock()
	return int(pr.pagecount)
}

func (pr *PdfPrinter) appendToDocument(page *Page) error {
	tracing.EngineTracer.Infof("OUTPUT PAGE [%d]", page.pageNo)
	page.pdfcanvas.Close()
	pr.doc.Assemble(page.pdfcanvas)
	return nil
}

// serialize outputs the PDF structure to a writer.
func (pr *PdfPrinter) serialize(w io.Writer) {
	T().Debugf("Serializing document")
	pr.doc.Encode(w)
}

// TODO
type RenderTree struct{}

// Page represents a page in the printer queue, as part of a print job.
// Pages will be created by Printer.PrintPage(...).
type Page struct {
	// https://www.prepressure.com/pdf/basics/page-boxes
	pageGeom  pdfapi.Rectangle // position and size on paper
	status    PageStatus       // print status
	pageNo    contPageNo       // page number
	content   *RenderTree      // page contents
	pdfcanvas *pdfapi.Canvas   // canvas to render onto
}

// PageNo returns the contiguous page number of a page.
func (page *Page) PageNo() int {
	return int(page.pageNo)
}

// --- Page Queue -------------------------------------------------------

// A pageQueue is the input queue for the printer. New pages to print
// will be enqueued into this queue.
type pageQueue struct {
	sync.RWMutex
	printer *PdfPrinter // the printer this queue is attached to
	running bool        // guarded by mutex
	working bool        // guarded by mutex
	pages   chan *Page  // pages in print
}

//Create a new printer queue. pr must not be nil.
func newQueue(pr *PdfPrinter) *pageQueue {
	q := &pageQueue{}
	q.printer = pr
	q.pages = make(chan *Page, 10)
	q.running = true
	return q
}

// enqueue will not block, is concurrency safe.
// Will start workers if not already done.
func (pq *pageQueue) enqueue(page *Page) {
	page.status = Queued
	queued := true
	pq.RLock()
	defer pq.RUnlock()
	if pq.running {
		select { // first try to enqueue synchronously
		case pq.pages <- page:
		default:
			queued = false
		}
		if !queued { // didn't work, now async
			go func(p *Page) {
				pq.pages <- p
			}(page)
		}
		if !pq.working {
			pq.startWorkers()
			pq.working = true
		}
	}
}

// We will start this many render workers.
const workerCount = 3

// startWorkers starts the render workers. These workers read pages from
// the printer's intput queue and render them.
func (pq *pageQueue) startWorkers() {
	for i := 0; i < workerCount; i++ {
		go func(pages <-chan *Page, printer *PdfPrinter) {
			for page := range pages {
				err := printer.render(page)
				if err == nil {
					pq.printer.pageComplete(page)
				} else {
					pq.printer.pageFailed(page, err)
				}
			}
			T().Debugf("Render worker stopped")
		}(pq.pages, pq.printer)
	}
}

func (pq *pageQueue) Stop() {
	pq.Lock()
	pq.running = false
	working := pq.working
	pq.Unlock()
	if working {
		close(pq.pages)
	}
}

func (pq *pageQueue) IsRunning() bool {
	pq.RLock()
	defer pq.RUnlock()
	return pq.running
}

// --- Helpers ----------------------------------------------------------

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

// We are tracing to the EngineTracer.
func T() tracing.Trace {
	return tracing.EngineTracer
}
