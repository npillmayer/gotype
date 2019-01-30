package dom

import (
	"strings"
	"testing"

	"github.com/aymerick/douceur/parser"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
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
	PrintTree(&sn.Node, t, domFmt)
	tracing.EngineTracer.Debugf("-------------------------------------------")
}

// --- Helpers ----------------------------------------------------------

func domFmt(dn RODomNode) string {
	return dn.String()
}

func PrintTree(n *tree.Node, t *testing.T, fmtnode func(RODomNode) string) {
	indent := 0
	dn := NewRONode(n, styledtree.Creator().ToStyler)
	printNode(dn, indent, t, fmtnode)
}

func printNode(dn RODomNode, w int, t *testing.T, fmtnode func(RODomNode) string) {
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
