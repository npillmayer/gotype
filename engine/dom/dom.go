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

type W3CNode struct {
	stylednode *styledtree.StyNode
}

var _ w3cdom.Node = &W3CNode{}

// NodeFromStyledNode creates a new DOM node from a styled node.
func NodeFromStyledNode(sn *styledtree.StyNode) w3cdom.Node {
	return &W3CNode{sn}
}

// NodeFromTreeNode creates a new DOM node from a tree node, which should
// be the inner node of a styledtree.Node.
func NodeFromTreeNode(tn *tree.Node) (w3cdom.Node, error) {
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

// NodeType returns the type of the underlying HTML node, something like
// html.ElementNode, html.TextNode, etc.
func (w *W3CNode) NodeType() html.NodeType {
	if w == nil {
		return html.ErrorNode
	}
	return w.stylednode.HtmlNode().Type
}

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

func (w *W3CNode) ComputedStyles() *style.PropertyMap {
	if w == nil {
		return nil
	}
	return w.stylednode.ComputedStyles()
}

// --- Attributes -----------------------------------------------------------------

type W3CAttr struct {
	attr *html.Attribute
}

var _ w3cdom.Attr = &W3CAttr{}

func (a *W3CAttr) Namespace() string {
	return a.attr.Namespace
}

func (a *W3CAttr) Key() string {
	return a.attr.Key
}

func (a *W3CAttr) Value() string {
	return a.attr.Val
}

var _ w3cdom.Node = &W3CAttr{} // Attributes are W3C DOM nodes as well

// AttrNode is an additional node type, complementing those defined in
// standard-package html.
const AttrNode = html.NodeType(77)

func (a *W3CAttr) NodeName() string {
	if a == nil {
		return ""
	}
	return a.attr.Key
}

func (a *W3CAttr) NodeValue() string {
	if a == nil {
		return ""
	}
	return a.attr.Val
}

func (a *W3CAttr) NodeType() html.NodeType            { return AttrNode }
func (a *W3CAttr) HasAttributes() bool                { return false }
func (a *W3CAttr) HasChildNodes() bool                { return false }
func (a *W3CAttr) ParentNode() w3cdom.Node            { return nil }
func (a *W3CAttr) ChildNodes() w3cdom.NodeList        { return nil }
func (a *W3CAttr) Children() w3cdom.NodeList          { return nil }
func (a *W3CAttr) FirstChild() w3cdom.Node            { return nil }
func (a *W3CAttr) Attributes() w3cdom.NamedNodeMap    { return nil }
func (a *W3CAttr) ComputedStyles() *style.PropertyMap { return nil }

// --- NamedNodeMap ---------------------------------------------------------------

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

func (wm *W3CMap) Length() int {
	if wm == nil {
		return 0
	}
	return len(wm.forNode.HtmlNode().Attr)
}

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

type W3CNodeList struct {
	nodes []*W3CNode
}

var _ w3cdom.NodeList = &W3CNodeList{}

func (wl *W3CNodeList) Length() int {
	if wl == nil {
		return 0
	}
	return len(wl.nodes)
}

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
func FromHTMLParseTree(h *html.Node) w3cdom.Node {
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
