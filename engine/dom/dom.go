package dom

import (
	"bytes"
	"fmt"

	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/dom/w3cdom"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// --- Node -----------------------------------------------------------------------

// A W3CNode is common type from which various types of DOM API objects inherit.
// This allows these types to be treated similarly.
type W3CNode struct {
	stylednode *styledtree.StyNode
}

var _ w3cdom.Node = &W3CNode{}

// NodeFromStyledNode creates a new DOM node from a styled node.
func NodeFromStyledNode(sn *styledtree.StyNode) *W3CNode {
	return &W3CNode{sn}
}

// NodeFromTreeNode creates a new DOM node from a tree node, which should
// be the inner node of a styledtree.Node.
func NodeFromTreeNode(tn *tree.Node) (*W3CNode, error) {
	if tn == nil {
		return nil, ErrNotAStyledNode
	}
	w := domify(tn)
	if w == nil {
		return nil, ErrNotAStyledNode
	}
	return w, nil
}

// ErrNotAStyledNode is returned if a tree node does not belong to a styled tree node.
var ErrNotAStyledNode = fmt.Errorf("Tree node is not a styled node")

func domify(tn *tree.Node) *W3CNode {
	sn := styledtree.Node(tn)
	if sn != nil {
		return &W3CNode{sn}
	}
	return nil
}

// NodeAsTreeNode returns the underlying tree.Node from a DOM node.
func NodeAsTreeNode(domnode w3cdom.Node) (*tree.Node, bool) {
	if domnode == nil {
		return nil, false
	}
	w, ok := domnode.(*W3CNode)
	if !ok {
		T().Errorf("DOM node has not been created from w3cdom.go")
		return nil, false
	}
	return &w.stylednode.Node, true
}

// HTMLNode returns the HTML parse node this DOM node is derived from.
func (w *W3CNode) HTMLNode() *html.Node {
	return w.stylednode.HtmlNode()
}

// NodeType returns the type of the underlying HTML node, something like
// html.ElementNode, html.TextNode, etc.
func (w *W3CNode) NodeType() html.NodeType {
	if w == nil {
		return html.ErrorNode
	}
	return w.stylednode.HtmlNode().Type
}

// NodeName read-only property returns the name of the current Node as a string.
//
//      Node         NodeName value
//      ------------+----------------------------
//      Attr         The value of Attr.name
//      Document     "#document"
//      Element      The value of Element.TagName
//      Text         "#text"
//
func (w *W3CNode) NodeName() string {
	if w == nil {
		return ""
	}
	h := w.stylednode.HtmlNode()
	switch h.Type {
	case html.DocumentNode:
		return "#document"
	case html.ElementNode:
		return h.Data
	case html.TextNode:
		return "#text"
	}
	return "<node>"
}

// NodeValue returns textual content for text/CData-Nodes, and an empty string for any other
// Node type.
func (w *W3CNode) NodeValue() string {
	if w == nil {
		return ""
	}
	h := w.stylednode.HtmlNode()
	if h.Type == html.TextNode {
		return h.Data
	}
	return ""
}

// HasAttributes returns a boolean indicating whether the current element has any
// attributes or not.
func (w *W3CNode) HasAttributes() bool {
	if w == nil {
		return false
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		return len(styledtree.Node(tn).HtmlNode().Attr) > 0
	}
	return false
}

// ParentNode read-only property returns the parent of the specified node in the DOM tree.
func (w *W3CNode) ParentNode() w3cdom.Node {
	if w == nil {
		return nil
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		p := tn.Parent()
		if p != nil {
			return domify(p)
		}
	}
	return nil
}

// HasChildNodes method returns a boolean value indicating whether the given Node
// has child nodes or not.
func (w *W3CNode) HasChildNodes() bool {
	if w == nil {
		return false
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		return tn.ChildCount() > 0
	}
	return false
}

// ChildNodes read-only property returns a live NodeList of child nodes of
// the given element.
func (w *W3CNode) ChildNodes() w3cdom.NodeList {
	if w == nil {
		return nil
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		children := tn.Children()
		childnodes := make([]*W3CNode, len(children))
		for i, ch := range children {
			childnodes[i] = &W3CNode{styledtree.Node(ch)}
		}
		return &W3CNodeList{childnodes}
	}
	return nil
}

// Children is a read-only property that returns a node list which contains all of
// the child *elements* of the node upon which it was called
func (w *W3CNode) Children() w3cdom.NodeList {
	if w == nil {
		return nil
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		children := tn.Children()
		childnodes := make([]*W3CNode, len(children))
		j := 0
		for _, ch := range children {
			sn := styledtree.Node(ch)
			if sn.HtmlNode().Type == html.ElementNode {
				childnodes[j] = &W3CNode{sn}
				j++
			}
		}
		return &W3CNodeList{childnodes}
	}
	return nil
}

// FirstChild read-only property returns the node's first child in the tree,
// or nil if the node has no children.
func (w *W3CNode) FirstChild() w3cdom.Node {
	if w == nil {
		return nil
	}
	tn, ok := NodeAsTreeNode(w)
	if ok && tn.ChildCount() > 0 {
		ch, ok := tn.Child(0)
		if ok {
			return domify(ch)
		}
	}
	return nil
}

// NextSibling read-only property returns the node immediately following the
// specified one in their parent's childNodes,
// or returns nil if the specified node is the last child in the parent element.
func (w *W3CNode) NextSibling() w3cdom.Node {
	if w == nil {
		return nil
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		if parent := tn.Parent(); parent != nil {
			if i := parent.IndexOfChild(tn); i >= 0 {
				sibling, ok := parent.Child(i + 1)
				if ok {
					return domify(sibling)
				}
			}
		}
	}
	return nil
}

// Attributes property returns a collection of all attribute nodes registered
// to the specified node. It is a NamedNodeMap, not an array.
func (w *W3CNode) Attributes() w3cdom.NamedNodeMap {
	if w == nil {
		return emptyNodeMap
	}
	h := w.stylednode.HtmlNode()
	switch h.Type {
	case html.DocumentNode:
	case html.ElementNode:
		return nodeMapFor(w.stylednode)
	}
	return emptyNodeMap
}

// TextContent property of the Node interface represents the text content of
// the node and its descendants.
//
// This implementation will include error strings in the text output, if errors occur.
// They will be flagged as "(ERROR: ... )".
func (w *W3CNode) TextContent() (string, error) {
	future := w.Walk().DescendentsWith(NodeIsText).Promise()
	textnodes, err := future()
	if err != nil {
		T().Errorf(err.Error())
		return "(ERROR: " + err.Error() + " )", err
	}
	var b bytes.Buffer
	var domnode *W3CNode
	for _, t := range textnodes {
		domnode, err = NodeFromTreeNode(t)
		if err != nil {
			b.WriteString("(ERROR: " + err.Error() + " )")
		} else {
			b.WriteString(domnode.NodeValue())
		}
	}
	return b.String(), err
}

// ComputedStyles returns a map of style properties for a given (stylable) Node.
func (w *W3CNode) ComputedStyles() *ComputedStyles {
	if w == nil {
		return nil
	}
	return &ComputedStyles{w, w.stylednode.ComputedStyles()}
}

// ComputedStyles is a proxy type for a node's styles.
type ComputedStyles struct {
	domnode  *W3CNode
	propsMap *style.PropertyMap
}

// Styles returns the underlying style.PropertyMap.
func (cstyles *ComputedStyles) Styles() *style.PropertyMap {
	return cstyles.propsMap
}

// HTMLNode returns the underlying html.Node.
func (cstyles *ComputedStyles) HTMLNode() *style.PropertyMap {
	retrn cstyles.domnode.HTMLNode()
}

var _ style.Styler = &ComputedStyles{} // implementing style.Styler may be useful

// Helper implementing style.Interf
func styler(n *tree.Node) style.Styler {
	return styledtree.Node(n)
}

// GetPropertyValue returns the property value for a given key.
func (cstyles *ComputedStyles) GetPropertyValue(key string) style.Property {
	if cstyles == nil {
		return style.NullStyle
	}
	node := &cstyles.domnode.stylednode.Node
	return cstyles.propsMap.GetPropertyValue(key, node, styler)
}

// --- Attributes -----------------------------------------------------------------

// A W3CAttr represents a single attribute of an element Node.
type W3CAttr struct {
	attr *html.Attribute
}

var _ w3cdom.Attr = &W3CAttr{}

// Namespace returns the namespace prefix of an attribute.
func (a *W3CAttr) Namespace() string {
	return a.attr.Namespace
}

// Key is the name of an attribute.
func (a *W3CAttr) Key() string {
	return a.attr.Key
}

// Value is the string value of an attribute.
func (a *W3CAttr) Value() string {
	return a.attr.Val
}

var _ w3cdom.Node = &W3CAttr{} // Attributes are W3C DOM nodes as well

// AttrNode is an additional node type, complementing those defined in
// standard-package html.
const AttrNode = html.NodeType(77)

// NodeName for an attribute is the attribute key
func (a *W3CAttr) NodeName() string {
	if a == nil {
		return ""
	}
	return a.attr.Key
}

// NodeValue for an attribute is the attribute value
func (a *W3CAttr) NodeValue() string {
	if a == nil {
		return ""
	}
	return a.attr.Val
}

// NodeType returns type AttrNode
func (a *W3CAttr) NodeType() html.NodeType { return AttrNode }

// HasAttributes returns false
func (a *W3CAttr) HasAttributes() bool { return false }

// HasChildNodes returns false
func (a *W3CAttr) HasChildNodes() bool { return false }

// ParentNode returns nil
func (a *W3CAttr) ParentNode() w3cdom.Node { return nil }

// ChildNodes returns nil
func (a *W3CAttr) ChildNodes() w3cdom.NodeList { return nil }

// Children returns nil
func (a *W3CAttr) Children() w3cdom.NodeList { return nil }

// FirstChild returns nil
func (a *W3CAttr) FirstChild() w3cdom.Node { return nil }

// NextSibling returns nil
func (a *W3CAttr) NextSibling() w3cdom.Node { return nil }

// Attributes returns nil
func (a *W3CAttr) Attributes() w3cdom.NamedNodeMap { return nil }

// TextContent returns an empty string
func (a *W3CAttr) TextContent() (string, error) { return "", nil }

// --- NamedNodeMap ---------------------------------------------------------------

// A W3CMap represents a key-value map
type W3CMap struct {
	forNode *styledtree.StyNode
}

var _ w3cdom.NamedNodeMap = &W3CMap{}

var emptyNodeMap = &W3CMap{}

func nodeMapFor(sn *styledtree.StyNode) w3cdom.NamedNodeMap {
	if sn != nil {
		return &W3CMap{sn}
	}
	return nil
}

// Length returns the number of entries in a key-value map
func (wm *W3CMap) Length() int {
	if wm == nil {
		return 0
	}
	return len(wm.forNode.HtmlNode().Attr)
}

// Item returns the i.th item in a key-value map
func (wm *W3CMap) Item(i int) w3cdom.Attr {
	if wm == nil {
		return nil
	}
	attrs := wm.forNode.HtmlNode().Attr
	if len(attrs) <= i || i < 0 {
		return nil
	}
	return &W3CAttr{&attrs[i]}
}

// GetNamedItem returns the attribute with key key.
func (wm *W3CMap) GetNamedItem(key string) w3cdom.Attr {
	if wm == nil {
		return nil
	}
	attrs := wm.forNode.HtmlNode().Attr
	for _, a := range attrs {
		if a.Key == key {
			return &W3CAttr{&a}
		}
	}
	return nil
}

// --- NodeList -------------------------------------------------------------------

// A W3CNodeList is a type for a list of nodes
type W3CNodeList struct {
	nodes []*W3CNode
}

var _ w3cdom.NodeList = &W3CNodeList{}

// Length returns the number of Nodes in a list
func (wl *W3CNodeList) Length() int {
	if wl == nil {
		return 0
	}
	return len(wl.nodes)
}

// Item returns the i.th Node
func (wl *W3CNodeList) Item(i int) w3cdom.Node {
	if wl == nil {
		return nil
	}
	if i >= len(wl.nodes) || i < 0 {
		return nil
	}
	return wl.nodes[i]
}

func (wl *W3CNodeList) String() string {
	var s bytes.Buffer
	s.WriteString("[ ")
	if wl != nil {
		for _, n := range wl.nodes {
			s.WriteString(n.NodeName())
			s.WriteString(" ")
		}
	}
	s.WriteString("]")
	return s.String()
}

// --------------------------------------------------------------------------------

// FromHTMLParseTree returns a W3C DOM from parsed HTML.
func FromHTMLParseTree(h *html.Node) *W3CNode {
	styles := douceuradapter.ExtractStyleElements(h)
	T().Debugf("Extracted %d <style> elements", len(styles))
	//c, errcss := parser.Parse(mycss)
	//if errhtml != nil || errcss != nil {
	//	T().Errorf("Cannot create test document")
	//}
	s := cssom.NewCSSOM(nil)
	for _, sty := range styles {
		s.AddStylesForScope(nil, sty, cssom.Script)
	}
	//s.AddStylesForScope(nil, douceuradapter.Wrap(c), cssom.Author)
	stytree, err := s.Style(h, styledtree.Creator())
	if err != nil {
		T().Errorf("Cannot style test document: %s", err.Error())
		return nil
	}
	return domify(stytree)
}

// Walk creates a tree walker set up to traverse the DOM.
func (w *W3CNode) Walk() *tree.Walker {
	if w == nil {
		return nil
	}
	return tree.NewWalker(&w.stylednode.Node)
}
