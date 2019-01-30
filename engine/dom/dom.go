package dom

import (
	"strings"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// ----------------------------------------------------------------------
//
// Material:
//
// github.com/andrewstuart/goq
// Package goq was built to allow users to declaratively unmarshal HTML into go
// structs using struct tags composed of css selectors.
//
// https://de.wikipedia.org/wiki/Mikroformat
//
// https://github.com/shurcooL/htmlg
// Package for generating and rendering HTML nodes with context-aware escaping.
//
// https://github.com/stilvoid/please
// Please is a command line utility that makes it easy to integrate web APIs
// into your shell scripts.
// It's called Please because the web works much better if you ask nicely.
//
// https://gowalker.org/sethwklein.net/go/webutil
// webutil provides Go functions that operate at a level better matching the
// level that I'm working at when I'm using JSON API's and scraping the web.
//
// https://github.com/tdewolff/minify
// Minify is a minifier package written in Go. It provides HTML5, CSS3, JS,
// JSON, SVG and XML minifiers.
//
// https://www.smashingmagazine.com/2015/01/designing-for-print-with-css/
//
// ----------------------------------------------------------------------

// We are tracing to the EngineTracer
func T() tracing.Trace {
	return tracing.EngineTracer
}

type RODomNode interface {
	String() string
	IsText() bool
	Style(string) style.Property
	HtmlNode() *html.Node
	Parent() RODomNode
	ChildrenIterator() DomNodeChildrenIterator
	//TreeNode() *tree.Node
}

type DomNodeChildrenIterator func() RODomNode

type domnode struct {
	stylednode *tree.Node
	toStyler   style.StyleInterf
}

func NewRONode(sn *tree.Node, toStyler style.StyleInterf) RODomNode {
	dn := domnode{}
	dn.stylednode = sn
	dn.toStyler = toStyler
	return dn
}

func (dn domnode) String() string {
	if dn.HtmlNode().Type == html.TextNode {
		return shortText(dn.HtmlNode())
	} else if dn.HtmlNode().Type == html.ElementNode {
		return dn.HtmlNode().Data
	}
	return "<>"
}

func (dn domnode) IsText() bool {
	if dn.HtmlNode().Type == html.TextNode {
		return true
	}
	return false
}

func (dn domnode) Style(key string) style.Property {
	return style.NullStyle
}

func (dn domnode) HtmlNode() *html.Node {
	h := dn.styler().HtmlNode()
	return h
}

func (dn domnode) Parent() RODomNode {
	if parent := dn.stylednode.Parent(); parent != nil {
		return NewRONode(parent, dn.toStyler)
	}
	return nil
}

func (dn domnode) styler() style.Styler {
	return dn.toStyler(dn.stylednode)
}

func (dn domnode) ChildrenIterator() DomNodeChildrenIterator {
	iteratable := dn
	hnode := dn.HtmlNode().FirstChild
	chinx := 0
	return func() RODomNode {
		if hnode == nil {
			return nil
		}
		if hnode.Type == html.ElementNode {
			inx := chinx
			chinx++
			ch, ok := iteratable.stylednode.Child(inx)
			hnode = hnode.NextSibling
			if !ok {
				return nil
			}
			return NewRONode(ch, iteratable.toStyler)
		} else if hnode.Type == html.TextNode {
			old := hnode
			hnode = hnode.NextSibling
			return newTextNode(iteratable, old) // return artificial DOM node
		}
		hnode = hnode.NextSibling
		return nil // ignore other types of nodes
	}
}

// func (dn domnode) TreeNode() *tree.Node {
// 	return dn.stylednode
// }

// --- Text Nodes -------------------------------------------------------

type textNode struct {
	htmlTextNode *html.Node
	parent       RODomNode
}

func newTextNode(parent RODomNode, text *html.Node) RODomNode {
	tn := textNode{}
	tn.htmlTextNode = text
	tn.parent = parent
	return tn
}

func (tn textNode) String() string {
	return shortText(tn.htmlTextNode)
}

func (tn textNode) IsText() bool {
	return true
}

func (tn textNode) Style(string) style.Property {
	return style.NullStyle
}

func (tn textNode) HtmlNode() *html.Node {
	return tn.htmlTextNode
}

func (tn textNode) Parent() RODomNode {
	return tn.parent
}

func (tn textNode) ChildrenIterator() DomNodeChildrenIterator {
	return zeroChildren
}

func zeroChildren() RODomNode {
	return nil
}

// func (tn textNode) TreeNode() *tree.Node {
// 	return nil
// }

// --- Helpers ----------------------------------------------------------

func shortText(h *html.Node) string {
	s := "CDATA\""
	if len(h.Data) > 10 {
		s += h.Data[:10] + "...\""
	} else {
		s += h.Data + "\""
	}
	s = strings.Replace(s, "\n", `\n`, -1)
	return s
}
