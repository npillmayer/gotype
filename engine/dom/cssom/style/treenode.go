package style

import "golang.org/x/net/html"

// TreeNode is the node type for our styled tree. We use an interface, as
// there will probably different implementations for batch and for interactive
// use. This is the least denominator of functions for reading style properties
// from a tree.
//
// Navigating the tree is implemented by other means (see XPath).
type TreeNode interface {
	HtmlNode() *html.Node         // the HTML DOM node related to this styled node
	ComputedStyles() *PropertyMap // get the computed styles of this styled node
	ParentNode() TreeNode         // safe getter (concurrency safe)
	ChildNode(int) TreeNode       // safe getter (concurrency safe)
	ChildCount() int              // number of children
}
