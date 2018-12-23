package style_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/aymerick/douceur/parser"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/logrusadapter"
	"github.com/npillmayer/gotype/engine/dom/style"
	"github.com/npillmayer/gotype/engine/dom/style/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/style/styled"
	"golang.org/x/net/html"
)

var T tracing.Trace

const (
	html1 string = `<body><p class="hello">Hello World</p></body>`
	html2 string = `<body><p id="single">Hello</p><p>World</p></body>`
	html3 string = `<body><p>Links:</p><ul><li><a href="foo">Foo</a><li>
<a href="/bar/baz">BarBaz</a></ul></body>`
	css1 string = `p { padding: 10px; } p.hello { color: blue; } #single { margin: 7px; }`
)

func Test0(t *testing.T) {
	tracing.EngineTracer = logrusadapter.New()
	tracing.EngineTracer.SetTraceLevel(tracing.LevelDebug)
	T = tracing.EngineTracer
}

func Test1(t *testing.T) {
	q, tree := setupTest(html1, css1)
	if tree == nil {
		t.Error("failed to setup test")
	}
	paras := findNodesFor("p", q, tree)
	if !assertProperty(paras, "padding-top").equals("10px") {
		t.Error("padding-top of paragraphs should be 10px")
	}
}

func Test2(t *testing.T) {
	q, tree := setupTest(html1, css1)
	if tree == nil {
		t.Error("failed to setup test")
	}
	paras := findNodesFor("p.hello", q, tree)
	if !assertProperty(paras, "color").equals("blue") {
		t.Error("color of paragraph with class=hello should be blue")
	}
}

// --- Helpers ----------------------------------------------------------

func getTestDOM(s string) *html.Node {
	doc, _ := html.Parse(strings.NewReader(s))
	return doc
}

func getTestCSS(s string) style.StyleSheet {
	css, _ := parser.Parse(s)
	return douceuradapter.Wrap(css)
}

func setupTest(htmlStr string, cssStr string) (*goquery.Document, *styled.Node) {
	dom := getTestDOM(htmlStr)
	css := getTestCSS(cssStr)
	cssom := style.NewCSSOM(nil)
	cssom.AddStylesFor(nil, css, style.Author)
	styledTree, err := cssom.Style(dom, styled.Factory{})
	if err != nil {
		T.Errorf("error: %s", err)
		return nil, nil
	}
	tree := styledTree.(*styled.Node)
	doc := goquery.NewDocumentFromNode(dom)
	return doc, tree
}

func findNodesFor(sel string, doc *goquery.Document, tree *styled.Node) []*styled.Node {
	q, _ := styled.QueryFor(tree)
	s := doc.Find(sel)
	var nodes []*styled.Node
	for _, n := range s.Nodes {
		sn := q.FindStyledNodeFor(n)
		if sn == nil {
			return nil
		}
		nodes = append(nodes, sn)
	}
	T.Debugf("found styled nodes: %v", nodes)
	return nodes
}

type props []style.Property

func assertProperty(nodes []*styled.Node, key string) props {
	if nodes == nil {
		return nil
	}
	var pp props
	for _, sn := range nodes {
		p := style.GetCascadedProperty(sn, key)
		T.Debugf("property %s of %s = %s", key, sn, p)
		pp = append(pp, p)
	}
	return pp
}

func (pp props) equals(property style.Property) bool {
	for _, p := range pp {
		if p != property {
			return false
		}
	}
	return true
}
