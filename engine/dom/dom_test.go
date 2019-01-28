package dom

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aymerick/douceur/parser"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/dom/styledtree/builder"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

func Test0(t *testing.T) {
	tracing.EngineTracer = gologadapter.New()
	tracing.EngineTracer.SetTraceLevel(tracing.LevelDebug)
}

var myhtml = `
<html><head></head><body>
  <p>The quick brown fox jumps over the lazy dog.</p>
  <p id="world">Hello <b>World</b>!</p>
  <p>This is a test.</p>
</body>
`

var mycss = `
p {
    margin-bottom: 10pt;
}
#world {
    margin-top: 20pt;
}
`

func prepareStyledTree(t *testing.T) *styledtree.StyNode {
	h, errhtml := html.Parse(strings.NewReader(myhtml))
	c, errcss := parser.Parse(mycss)
	if errhtml != nil || errcss != nil {
		T().Errorf("Cannot create test document")
	}
	s := cssom.NewCSSOM(nil)
	s.AddStylesForScope(nil, douceuradapter.Wrap(c), cssom.Author)
	doc, err := s.Style(h, builder.New())
	if err != nil {
		t.Errorf("Cannot style test document: %s", err.Error())
	}
	return doc.(*styledtree.StyNode)
}

func TestDom1(t *testing.T) {
	tracing.EngineTracer.Debugf("-------------------------------------------")
	sn := prepareStyledTree(t)
	PrintTree(&sn.Node, t, stylestring)
	tracing.EngineTracer.Debugf("-------------------------------------------")
}

// --- Helpers ----------------------------------------------------------

func PrintTree(n *tree.Node, t *testing.T, fmtnode func(*tree.Node) string) {
	indent := 0
	printNode(n, indent, t, fmtnode)
}

func stylestring(n *tree.Node) string {
	props := " ["
	pmap := styledtree.Creator().ToStyler(n).ComputedStyles()
	if pmap != nil {
		if p, ok := style.GetLocalProperty(pmap, "margin-bottom"); ok {
			props += fmt.Sprintf("%s: %s;", "margin-bottom", p)
		}
		if p, ok := style.GetLocalProperty(pmap, "margin-top"); ok {
			props += fmt.Sprintf("%s: %s;", "margin-top", p)
		}
	}
	props += "]"
	return n.Payload.(*styledtree.StyNode).HtmlNode().Data + props
}

func printNode(n *tree.Node, w int, t *testing.T, fmtnode func(*tree.Node) string) {
	t.Logf("%s%s = {", indent(w), fmtnode(n))
	for i := 0; i < n.ChildCount(); i++ {
		ch, _ := n.Child(i)
		printNode(ch, w+1, t, fmtnode)
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
