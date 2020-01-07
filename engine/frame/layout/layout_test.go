package layout_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/engine/frame/layout"
	"github.com/npillmayer/gotype/engine/frame/renderdbg"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/engine/dom"
	"golang.org/x/net/html"
)

var graphviz = false

var Xmyhtml = `<div><b>bold</b><table></table><span>spanned</span></div>`

var myhtml = `
<html><head>
<style>
  body { border-color: red; }
</style>
</head><body>
  <p>The quick brown fox jumps over the lazy</p><b>dog.</b>
  <p id="world">Hello <b>World</b>!</p>
  <p style="padding-left: 5px; position: fixed;">This is a test.</p>
</body>
`

func buildDOM(t *testing.T) *dom.W3CNode {
	h, err := html.Parse(strings.NewReader(myhtml))
	if err != nil {
		t.Errorf("Cannot create test document")
	}
	dom := dom.FromHTMLParseTree(h, nil) // nil = no external stylesheet
	if dom == nil {
		t.Errorf("Could not build DOM from HTML")
	}
	return dom
}

func TestBoxGeneration(t *testing.T) {
	//gtrace.EngineTracer = gologadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.EngineTracer.SetTraceLevel(tracing.LevelInfo)
	domroot := buildDOM(t)
	boxes, err := layout.BuildBoxTree(domroot)
	if err != nil {
		t.Errorf(err.Error())
	} else if boxes == nil {
		t.Errorf("Render tree root is null")
	} else {
		t.Logf("root node is %s", boxes.DOMNode().NodeName())
		if boxes.DOMNode().NodeName() != "#document" {
			t.Errorf("name of root element expected to be '#document")
		}
		if graphviz {
			gvz, _ := ioutil.TempFile(".", "layout-*.dot")
			defer gvz.Close()
			renderdbg.ToGraphViz(boxes.(*layout.PrincipalBox), gvz)
		}
	}
}

func TestReorderingSimple(t *testing.T) {
	//viewport := &dimen.Rect{TopL: dimen.Origin, BotR: dimen.DINA4}
	gtrace.EngineTracer = gologadapter.New()
	gtrace.EngineTracer.SetTraceLevel(tracing.LevelError)
	domroot := buildDOM(t)
	boxes, err := layout.BuildBoxTree(domroot)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		t.Logf("root node is %s", boxes.DOMNode().NodeName())
		gtrace.EngineTracer.SetTraceLevel(tracing.LevelDebug)
		gtrace.EngineTracer.Debugf("=== Reordering Boxes =========================")
		renderRoot := layout.TreeNodeAsPrincipalBox(boxes.TreeNode())
		if err = layout.ReorderBoxTree(renderRoot); err != nil {
			t.Errorf(err.Error())
		} else if graphviz {
			gvz, _ := ioutil.TempFile(".", "reorder-*.dot")
			defer gvz.Close()
			renderdbg.ToGraphViz(boxes.(*layout.PrincipalBox), gvz)
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
