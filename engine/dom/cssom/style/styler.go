package style

import (
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// Styler is an interface all concrete types of styled tree nodes
// will have to implement to be usable for layout, rendering, etc.
type Styler interface {
	HtmlNode() *html.Node
	ComputedStyles() *PropertyMap
}

// StyleInterf is a mapper from a concrete tree node to an interface
// implementation for Styler. You can think of this function type as
// an adapter from a certain tree implementation to a styled tree.
type StyleInterf func(*tree.Node) Styler

// Creator is a function to create a style node for a given
// HTML node.
// You can think of this interface as
// an adapter from a certain tree implementation to a styled tree.
type Creator interface {
	StyleForHtmlNode(*html.Node) *tree.Node
	ToStyler(*tree.Node) Styler
	SetComputedStyles(*tree.Node, *PropertyMap)
}
