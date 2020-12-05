package html

// github.com/andybalholm/cascadia
// github.com/PuerkitoBio/goquery

import (
	"io"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/PuerkitoBio/goquery"
)

// CT traces to the core-tracer.
func CT() tracing.Trace {
	return gtrace.CoreTracer
}

func ReadHTMLBook(r io.Reader) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	//doc, err = html.Parse(r)
	if err != nil {
		CT().Errorf("Unable to parse HTML file: %s", err)
	}
	return doc, err
}
