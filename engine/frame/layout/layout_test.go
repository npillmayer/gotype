package layout_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/dimen"

	"github.com/npillmayer/gotype/engine/frame/layout"
	"github.com/npillmayer/gotype/engine/frame/renderdbg"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/engine/dom"
	"golang.org/x/net/html"
)

var graphviz = false

// func T() tracing.Trace {
// 	return gtrace.EngineTracer
// }

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

func buildDOM(t *testing.T) *dom.W3CNode {
	h, err := html.Parse(strings.NewReader(myhtml))
	if err != nil {
		t.Errorf("Cannot create test document")
	}
	return dom.FromHTMLParseTree(h)
}

func TestLayout1(t *testing.T) {
	gtrace.EngineTracer = gologadapter.New()
	// teardown := gotestingadapter.RedirectTracing(t)
	// // defer teardown()
	gtrace.EngineTracer.SetTraceLevel(tracing.LevelInfo)
	root := buildDOM(t)
	gtrace.EngineTracer.Infof("===================================================")
	layout := layout.NewLayouter(root)
	viewport := &dimen.Rect{TopL: dimen.Origin, BotR: dimen.DINA4}
	layout.Layout(viewport)
	if layout.BoxRoot() == nil {
		t.Errorf("Layout failed, render tree root is null")
	} else {
		if root.NodeName() != "#document" {
			t.Errorf("name of root element expected to be '#document")
		}
		t.Logf("root node is %s", root.NodeName())
		if graphviz {
			gvz, _ := ioutil.TempFile(".", "layout-*.dot")
			defer gvz.Close()
			renderdbg.ToGraphViz(layout, gvz)
		}
	}
}

// ----------------------------------------------------------------------

/*
func Test0(t *testing.T) {
	tracing.EngineTracer = gologadapter.New()
}

func TestLayout1(t *testing.T) {
	root := styledtree.NewNodeForHtmlNode(makeHtmlNode("body"))
	p1 := styledtree.NewNodeForHtmlNode(makeHtmlNode("p"))
	p2 := styledtree.NewNodeForHtmlNode(makeHtmlNode("p"))
	p1.HtmlNode().AppendChild(makeCData("Hello World!"))
	p2.HtmlNode().AppendChild(makeCData("More text."))
	root.AddChild(p1)
	root.AddChild(p2)
	tree := buildLayoutTree(root)
	PrintTree(tree, t)
	if l := len(tree.content[0].content); l != 2 {
		t.Errorf("<body> should have 2 content children, has %d", l)
	}
	group := style.NewPropertyGroup("Display")
	group.Set("display", "none")
	pmap := style.NewPropertyMap()
	pmap.AddAllFromGroup(group, true)
	p2.SetComputedStyles(pmap)
	t.Log("-----------------------")
	d, _ := style.GetLocalProperty(p2, "display")
	t.Logf("Now p2.display = %s\n", d)
	t.Log("-----------------------")
	tree = buildLayoutTree(root)
	if l := len(tree.content[0].content); l != 1 {
		t.Errorf("<body> should have 1 content child, has %d", l)
	}
	PrintTree(tree, t)
	group.Set("display", "inline")
	pmap.AddAllFromGroup(group, true)
	t.Log("-----------------------")
	d, _ = style.GetLocalProperty(p2, "display")
	t.Logf("Now p2.display = %s\n", d)
	t.Log("-----------------------")
	tree = buildLayoutTree(root)
	PrintTree(tree, t)
}
*/

// ----------------------------------------------------------------------

/*
func makeHTMLNode(e string) *html.Node {
	n := &html.Node{}
	n.Type = html.ElementNode
	n.Data = e
	return n
}

func makeCData(s string) *html.Node {
	n := &html.Node{}
	n.Type = html.TextNode
	n.Data = s
	return n
}
*/

/*
func PrintTree(c *Container, t *testing.T) {
	indent := 0
	printContainer(c, indent, t)
}

func printContainer(c *Container, w int, t *testing.T) {
	t.Logf("%s%s{\n", indent(w), c)
	for _, ch := range c.content {
		printContainer(ch, w+1, t)
	}
	t.Logf("%s}", indent(w))
}

func indent(w int) string {
	s := ""
	for w > 0 {
		s += "   "
		w--
	}
	return s
}
*/
