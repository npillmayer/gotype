package dom

import (
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// Node represents W3C-type Node
type Node interface {
	HasAttributes() bool
	HasChildNodes() bool
	ChildNodes() NodeList
	FirstChild() Node
}

// NodeList represents W3C-type NodeList
type NodeList interface {
	Length() int
	Item(int) Node
}

// --- Node -----------------------------------------------------------------------

type w3cNode struct {
	styledtree.StyNode
}

var _ Node = &w3cNode{}

// NodeAsTreeNode returns the underlying tree.Node from a DOM node.
func NodeAsTreeNode(w3cNode *Node) (*tree.Node, bool) {
	w, ok := domnode.(*w3cNode)
	if !ok {
		T().Errorf("DOM node has not been created from w3cdom.go")
		return nil, false
	}
	return &w.Node, true
}

func (w *w3cNode) HasAttributes() bool {
	return false
}

func (w *w3cNode) HasChildNodes() bool {
	tn, ok := NodeAsTreeNode(w)
	if ok {
		return tn.ChildCount > 0
	}
	return false
}

func (w *w3cNode) ChildNodes() NodeList {
	tn, ok := NodeAsTreeNode(w)
	if ok {
		children := tn.Children()
		childnodes := make([]styledtree.StyNode, children.length())
		for i := children.length() - 1; i >= 0; i-- {
			childnodes[i] = &w3cNode{styledtree.Node(children[i])}
		}
		return &w3cNodeList{childnodes}
	}
	return nil
}

func (w *w3cNode) FirstChild() Node {
	tn, ok := NodeAsTreeNode(w)
	if ok && tn.ChildCount>0 {
		ch, ok := tn.Child(0) {
			return &w3cNode{styledtree.Node(ch)}
		}
	}
	return nil
}

// --- NodeList -------------------------------------------------------------------

type w3cNodeList struct {
	nodes []w3cNode
}

var _ NodeList = &w3cNodeList{}

func (wl *w3cNodeList) Length() int {
	return wl.nodes.length()
}

func (wl *w3cNodeList) Item(int i) Node {
	return wl.nodes[i]
}

// --------------------------------------------------------------------------------

// FromHTMLParseTree returns a DOM from parsed HTML.
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
	}
	dom := &w3cNode{*styledtree.Node(stytree)}
	return dom
}
