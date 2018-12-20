package style

import (
	"log"
	"strings"

	//	"testing"

	"github.com/aymerick/douceur/css"
	"github.com/aymerick/douceur/parser"

	//"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/tracing"
	//"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"golang.org/x/net/html"
)

func getTestDOM() *html.Node {
	s := `<body><p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul></body>`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func getTestCSS() *css.Stylesheet {
	testCSS := `table, ul, td {
  padding: 0;
}
body, h1, ul, p   {
  margin-top: 10px;
}`
	stylesheet, err := parser.Parse(testCSS)
	if err != nil {
		tracing.EngineTracer.Errorf("could not parse test CSS")
		panic("error parsing test CSS")
	}
	return stylesheet
}

/*
func Test1(t *testing.T) {
	config.InitTracing(gologadapter.GetAdapter())
	styles := getTestCSS()
	tracing.With(tracing.CoreTracer).Dump("styles", styles)
}

func TestCSSOM1(t *testing.T) {
	dom := getTestDOM()
	cssom := NewCSSOM(nil)
	cssom.AddStylesFor(dom, getTestCSS(), Author)
}
*/

/*
func Test2(t *testing.T) {
	dom := getTestDOM()
	cssom := getTestCSS()
	rulestree := &RulesTree{cssom, nil}
	ConstructStyledNodeTree(dom, rulestree)
}
*/
