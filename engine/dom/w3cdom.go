package dom

import (
	"fmt"

	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// Node represents W3C-type Node
type Node interface {
	NodeType() html.NodeType  // type of the underlying HTML node (ElementNode, TextNode, etc.)
	NodeName() string         // node name output depends on the node's type
	NodeValue() string        // node value output depends on the node's type
	HasAttributes() bool      // check for existence of attributes
	HasChildNodes() bool      // check for existende of sub-nodes
	ChildNodes() NodeList     // get a list of all children-nodes
	Children() NodeList       // get a list of element child-nodes
	FirstChild() Node         // get the first children-node
	Attributes() NamedNodeMap // get all attributes of a node
	ComputedStyles() *style.PropertyMap
}

// NodeList represents W3C-type NodeList
type NodeList interface {
	Length() int
	Item(int) Node
}

// Attr represents W3C-type Attr
type Attr interface {
	Namespace() string
	Key() string
	Value() string
}

// NamedNodeMap represents w3C-type NamedNodeMap
type NamedNodeMap interface {
	Length() int
	Item(int) Attr
	GetNamedItem(string) Attr
}

// --- Node -----------------------------------------------------------------------

type w3cNode struct {
	stylednode *styledtree.StyNode
}

var _ Node = &w3cNode{}

// NodeFromStyledNode creates a new DOM node from a styled node.
func NodeFromStyledNode(sn *styledtree.StyNode) Node {
	return &w3cNode{sn}
}

// NodeFromTreeNode creates a new DOM node from a tree node, which should
// be the inner node of a styledtree.Node.
func NodeFromTreeNode(tn *tree.Node) (Node, error) {
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

func domify(tn *tree.Node) *w3cNode {
	sn := styledtree.Node(tn)
	if sn != nil {
		return &w3cNode{sn}
	}
	return nil
}

// NodeAsTreeNode returns the underlying tree.Node from a DOM node.
func NodeAsTreeNode(domnode Node) (*tree.Node, bool) {
	if domnode == nil {
		return nil, false
	}
	w, ok := domnode.(*w3cNode)
	if !ok {
		T().Errorf("DOM node has not been created from w3cdom.go")
		return nil, false
	}
	return &w.stylednode.Node, true
}

// NodeType returns the type of the underlying HTML node, something like
// html.ElementNode, html.TextNode, etc.
func (w *w3cNode) NodeType() html.NodeType {
	if w == nil {
		return html.ErrorNode
	}
	return w.stylednode.HtmlNode().Type
}

func (w *w3cNode) NodeName() string {
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

func (w *w3cNode) NodeValue() string {
	if w == nil {
		return ""
	}
	h := w.stylednode.HtmlNode()
	if h.Type == html.TextNode {
		return h.Data
	}
	return ""
}

func (w *w3cNode) HasAttributes() bool {
	if w == nil {
		return false
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		return len(styledtree.Node(tn).HtmlNode().Attr) > 0
	}
	return false
}

func (w *w3cNode) HasChildNodes() bool {
	if w == nil {
		return false
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		return tn.ChildCount() > 0
	}
	return false
}

func (w *w3cNode) ChildNodes() NodeList {
	if w == nil {
		return nil
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		children := tn.Children()
		childnodes := make([]*w3cNode, len(children))
		for i, ch := range children {
			childnodes[i] = &w3cNode{styledtree.Node(ch)}
		}
		return &w3cNodeList{childnodes}
	}
	return nil
}

func (w *w3cNode) Children() NodeList {
	if w == nil {
		return nil
	}
	tn, ok := NodeAsTreeNode(w)
	if ok {
		children := tn.Children()
		childnodes := make([]*w3cNode, len(children))
		j := 0
		for _, ch := range children {
			sn := styledtree.Node(ch)
			if sn.HtmlNode().Type == html.ElementNode {
				childnodes[j] = &w3cNode{sn}
				j++
			}
		}
		return &w3cNodeList{childnodes}
	}
	return nil
}

func (w *w3cNode) FirstChild() Node {
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

func (w *w3cNode) NextSibling() Node {
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

func (w *w3cNode) Attributes() NamedNodeMap {
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

func (w *w3cNode) ComputedStyles() *style.PropertyMap {
	if w == nil {
		return nil
	}
	return w.stylednode.ComputedStyles()
}

// --- Attributes -----------------------------------------------------------------

type w3cAttr struct {
	attr *html.Attribute
}

var _ Attr = &w3cAttr{}

func (a *w3cAttr) Namespace() string {
	return a.attr.Namespace
}

func (a *w3cAttr) Key() string {
	return a.attr.Key
}

func (a *w3cAttr) Value() string {
	return a.attr.Val
}

var _ Node = &w3cAttr{} // Attributes are W3C DOM nodes as well

// AttrNode is an additional node type, complementing those defined in
// standard-package html.
const AttrNode = html.NodeType(77)

func (a *w3cAttr) NodeName() string {
	if a == nil {
		return ""
	}
	return a.attr.Key
}

func (a *w3cAttr) NodeValue() string {
	if a == nil {
		return ""
	}
	return a.attr.Val
}

func (a *w3cAttr) NodeType() html.NodeType            { return AttrNode }
func (a *w3cAttr) HasAttributes() bool                { return false }
func (a *w3cAttr) HasChildNodes() bool                { return false }
func (a *w3cAttr) ChildNodes() NodeList               { return nil }
func (a *w3cAttr) Children() NodeList                 { return nil }
func (a *w3cAttr) FirstChild() Node                   { return nil }
func (a *w3cAttr) Attributes() NamedNodeMap           { return nil }
func (a *w3cAttr) ComputedStyles() *style.PropertyMap { return nil }

// --- NamedNodeMap ---------------------------------------------------------------

type w3cMap struct {
	forNode *styledtree.StyNode
}

var _ NamedNodeMap = &w3cMap{}

var emptyNodeMap = &w3cMap{}

func nodeMapFor(sn *styledtree.StyNode) NamedNodeMap {
	if sn != nil {
		return &w3cMap{sn}
	}
	return nil
}

func (wm *w3cMap) Length() int {
	if wm == nil {
		return 0
	}
	return len(wm.forNode.HtmlNode().Attr)
}

func (wm *w3cMap) Item(i int) Attr {
	if wm == nil {
		return nil
	}
	attrs := wm.forNode.HtmlNode().Attr
	if len(attrs) <= i || i < 0 {
		return nil
	}
	return &w3cAttr{&attrs[i]}
}

func (wm *w3cMap) GetNamedItem(key string) Attr {
	if wm == nil {
		return nil
	}
	attrs := wm.forNode.HtmlNode().Attr
	for _, a := range attrs {
		if a.Key == key {
			return &w3cAttr{&a}
		}
	}
	return nil
}

// --- NodeList -------------------------------------------------------------------

type w3cNodeList struct {
	nodes []*w3cNode
}

var _ NodeList = &w3cNodeList{}

func (wl *w3cNodeList) Length() int {
	if wl == nil {
		return 0
	}
	return len(wl.nodes)
}

func (wl *w3cNodeList) Item(i int) Node {
	if wl == nil {
		return nil
	}
	if i >= len(wl.nodes) || i < 0 {
		return nil
	}
	return wl.nodes[i]
}

// --------------------------------------------------------------------------------

// FromHTMLParseTree returns a W3C DOM from parsed HTML.
func FromHTMLParseTree(h *html.Node) Node {
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
