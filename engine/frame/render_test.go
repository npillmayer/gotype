package frame

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/aymerick/douceur/parser"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/domdbg"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/frame/layout"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

var myhtml = `
<html><head>
	<style>
	body { border-color: red; }
	</style>
</head><body>
	<p>The quick brown fox jumps over the lazy dog.</p>
	<p id="world">Hello <b>World</b>!</p>
	<p style="padding-left: 5px;">This is a test.</p>
</body>
`

var mycss = `
p {
	margin-bottom: 10pt;
}
#world {
	padding-top: 20pt;
}
`

// T sets up an engine tracer.
func T() tracing.Trace {
	return tracing.EngineTracer
}

func Test0(t *testing.T) {
	tracing.EngineTracer = gologadapter.New()
	//tracing.EngineTracer.SetTraceLevel(tracing.LevelDebug)
}

func TestVizDom1(t *testing.T) {
	tracing.EngineTracer = gologadapter.New()
	sn := prepareStyledTree()
	log.Println("--- Styling done --------------------------")
	doc := dom.NewRONode(sn, styledtree.Creator().ToStyler)
	gvz, _ := ioutil.TempFile(".", "gvz-style-*.dot")
	defer gvz.Close()
	domdbg.ToGraphViz(doc, gvz, nil)
	tracing.EngineTracer.SetTraceLevel(tracing.LevelDebug)
	layouter := layout.NewLayouter(sn, styledtree.Creator())
	c := layouter.Layout(&dimen.Rect{
		TopL: dimen.Point{X: 0, Y: 0},
		BotR: dimen.Point{X: 500, Y: 500},
	})
	print(c.String())
	gvz2, _ := ioutil.TempFile(".", "gvz-render-*.dot")
	defer gvz2.Close()
	ToGraphViz(c, gvz2)
}

func prepareStyledTree() *tree.Node {
	h, errhtml := html.Parse(strings.NewReader(myhtml))
	styles := douceuradapter.ExtractStyleElements(h)
	log.Printf("Extracted %d <style> elements", len(styles))
	c, errcss := parser.Parse(mycss)
	if errhtml != nil || errcss != nil {
		log.Println("Cannot create test document")
	}
	s := cssom.NewCSSOM(nil)
	for _, sty := range styles {
		s.AddStylesForScope(nil, sty, cssom.Script)
	}
	s.AddStylesForScope(nil, douceuradapter.Wrap(c), cssom.Author)
	doc, err := s.Style(h, styledtree.Creator())
	if err != nil {
		log.Printf("Cannot style test document: %s", err.Error())
	}
	return doc
}
