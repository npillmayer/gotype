package pdf

import (
	"os"
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/dimen"
)

func Test0(t *testing.T) {
	tracing.EngineTracer = gologadapter.New()
	tracing.EngineTracer.SetTraceLevel(tracing.LevelDebug)
}

var stdpage = dimen.Rect{
	dimen.Point{100, 100},
	dimen.Point{10000, 10000},
}

func TestAssemble1(t *testing.T) {
	papersize := dimen.Point{20000, 20000}
	printer := Printer(papersize, 1.0)
	future := printer.Start(os.Stdout)
	printer.SetMaxPage(2) // set target early
	page1 := printer.PrintPage(1, stdpage, nil)
	printer.pageComplete(page1)
	page2 := printer.PrintPage(2, stdpage, nil)
	printer.pageComplete(page2)
	err := future()
	if err != nil {
		t.Error(err.Error())
	}
	T().Debugf("Printer returned error = %v", err)
}

func TestAssemble2(t *testing.T) {
	papersize := dimen.Point{20000, 20000}
	printer := Printer(papersize, 1.0)
	future := printer.Start(os.Stdout)
	page1 := printer.PrintPage(1, stdpage, nil)
	page2 := printer.PrintPage(2, stdpage, nil)
	printer.pageComplete(page1)
	printer.pageComplete(page2)
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
