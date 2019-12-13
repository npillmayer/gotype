package dom_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aymerick/douceur/parser"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/domdbg"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

func T() tracing.Trace {
	return tracing.EngineTracer
}

func Test0(t *testing.T) {
	tracing.EngineTracer = gologadapter.New()
	//tracing.EngineTracer.SetTraceLevel(tracing.LevelDebug)
}

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

func prepareStyledTree(t *testing.T) *tree.Node {
	h, errhtml := html.Parse(strings.NewReader(myhtml))
	styles := douceuradapter.ExtractStyleElements(h)
	t.Logf("Extracted %d <style> elements", len(styles))
	c, errcss := parser.Parse(mycss)
	if errhtml != nil || errcss != nil {
		T().Errorf("Cannot create test document")
	}
	s := cssom.NewCSSOM(nil)
	for _, sty := range styles {
		s.AddStylesForScope(nil, sty, cssom.Script)
	}
	s.AddStylesForScope(nil, douceuradapter.Wrap(c), cssom.Author)
	doc, err := s.Style(h, styledtree.Creator())
	if err != nil {
		t.Errorf("Cannot style test document: %s", err.Error())
	}
	return doc
}

func TestDom1(t *testing.T) {
	tracing.EngineTracer.Debugf("===========================================")
	sn := prepareStyledTree(t)
	PrintTree(sn, t, domFmt)
	tracing.EngineTracer.Debugf("-------------------------------------------")
}

func TestDom2(t *testing.T) {
	tracing.EngineTracer.SetTraceLevel(tracing.LevelDebug)
	tracing.EngineTracer.Debugf("===========================================")
	sn := prepareStyledTree(t)
	tracing.EngineTracer.Debugf("--- Styling done --------------------------")
	doc := dom.NewRONode(sn, styledtree.Creator().ToStyler)
	gvz, _ := ioutil.TempFile(".", "graphviz-*.dot")
	defer gvz.Close()
	domdbg.ToGraphViz(doc, gvz, nil)
}

// --- Helpers ----------------------------------------------------------

func domFmt(dn dom.RODomNode) string {
	return dn.String()
}

func PrintTree(n *tree.Node, t *testing.T, fmtnode func(dom.RODomNode) string) {
	indent := 0
	dn := dom.NewRONode(n, styledtree.Creator().ToStyler)
	printNode(dn, indent, t, fmtnode)
}

func printNode(dn dom.RODomNode, w int, t *testing.T, fmtnode func(dom.RODomNode) string) {
	if dn.IsText() {
		t.Logf("%s%s", indent(w), fmtnode(dn))
	} else {
		t.Logf("%s%s = {", indent(w), fmtnode(dn))
		it := dn.ChildrenIterator()
		ch := it()
		for ch != nil {
			printNode(ch, w+1, t, fmtnode)
			ch = it()
		}
		t.Logf("%s}", indent(w))
	}
}

func indent(w int) string {
	s := ""
	for w > 0 {
		s += "   "
		w--
	}
	return s
}