package html

// github.com/andybalholm/cascadia
// github.com/PuerkitoBio/goquery

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
)

// We trace to the core-tracer.
var CT tracing.Trace = tracing.CoreTracer

func ReadHTMLbook(r io.Reader) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	//doc, err := html.Parse(r)
	if err != nil {
		CT.Errorf("Unable to parse HTML file: %s", err)
	}
	return doc, err
}
