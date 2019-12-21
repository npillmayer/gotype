package dom

import (
	"strings"

	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// ----------------------------------------------------------------------
//
// Material:
//
// DOM:
// https://www.data2type.de/xml-xslt-xslfo/xml/xml-in-a-nutshell/document-object-model/interface-node/
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

// RODomNode is an interface type which represents nodes of a DOM tree
// in read-only mode, i.e. the node may not be modified by the client.
type RODomNode interface {
	String() string                            // default string representation
	HtmlNode() *html.Node                      // return the underlying HTML node
	IsText() bool                              // is this an HTML CDATA node?
	Parent() RODomNode                         // give the parent node
	ChildrenIterator() DomNodeChildrenIterator // iterate over children
	ComputedStyles() *style.PropertyMap        // all CSS style properties
	Style(string) style.Property               // return a CSS style property
	//TreeNode() *tree.Node
}

// DomNodeChildrenIterator iterates over the children nodes of a DOM node.
// Clients should call it until it returns nil.
type DomNodeChildrenIterator func() RODomNode

type domnode struct {
	stylednode *tree.Node
	toStyler   style.StyleInterf
}

// NewRONode wraps a tree node into an RODomNode. Callers must supply
// a style.Styler interface converters, which is responsible for
// giving a Styler for a generic tree node.
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
	} else if dn.HtmlNode().Type == html.DocumentNode {
		return ":root"
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

func (dn domnode) ComputedStyles() *style.PropertyMap {
	return dn.styler().ComputedStyles()
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
		return nil // ignore other types of nodes // TODO bug
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

func (tn textNode) ComputedStyles() *style.PropertyMap {
	return nil
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
