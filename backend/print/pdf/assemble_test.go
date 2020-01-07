package pdf

// TODO remove generated PDF files

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/core/dimen"
)

var a5page = dimen.Rect{
	dimen.Point{3 * dimen.CM, 4 * dimen.CM},
	dimen.Point{3*dimen.CM + dimen.DINA5.X, 4*dimen.CM + dimen.DINA5.Y},
}

func TestAssemble1(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	printer := NewPrinter(dimen.DINA4, 0.5)
	f := tmpPdfFile(t)
	defer f.Close()
	future := printer.Start(f)
	printer.SetMaxPage(2) // set target early
	printer.PrintPage(1, a5page, nil)
	printer.PrintPage(2, a5page, nil)
	err := future()
	if err != nil {
		t.Error(err.Error())
	}
	T().Debugf("Printer returned error = %v", err)
	printer.PrintPage(3, a5page, nil) // should do nothing, not panic
}

/*
func TestAssemble2(t *testing.T) {
	printer := Printer(dimen.DINA4, 0.5)
	f := tmpPdfFile(t)
	defer f.Close()
	future := printer.Start(f)
	page1 := printer.PrintPage(1, a5page, nil)
	page2 := printer.PrintPage(2, a5page, nil)
	printer.pageComplete(page1)
	printer.pageFailed(page2)
	printer.SetMaxPage(2) // set target late
	err := future()
	if err != nil {
		t.Error(err.Error())
	}
	T().Debugf("Printer returned error = %v", err)
	if printer.PageCount() != 2 {
		t.Error("Expected 2 pages to have been printed")
	}
}
*/

func tmpPdfFile(t *testing.T) *os.File {
	tmpfile, err := ioutil.TempFile(".", "test-*.pdf")
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("Writing PDF to %s", tmpfile.Name())
	return tmpfile
}
