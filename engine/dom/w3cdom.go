package dom

import (
	"fmt"

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
	w, ok := domnode.(*w3cNode)
	if !ok {
		T().Errorf("DOM node has not been created from w3cdom.go")
		return nil, false
	}
	return &w.stylednode.Node, true
}

func (w *w3cNode) HasAttributes() bool {
	return false
}

func (w *w3cNode) HasChildNodes() bool {
	tn, ok := NodeAsTreeNode(w)
	if ok {
		return tn.ChildCount() > 0
	}
	return false
}

func (w *w3cNode) ChildNodes() NodeList {
	tn, ok := NodeAsTreeNode(w)
	if ok {
		children := tn.Children()
		childnodes := make([]*w3cNode, len(children))
		for i := len(children) - 1; i >= 0; i-- {
			childnodes[i] = &w3cNode{styledtree.Node(children[i])}
		}
		return &w3cNodeList{childnodes}
	}
	return nil
}

func (w *w3cNode) FirstChild() Node {
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

// --- NodeList -------------------------------------------------------------------

type w3cNodeList struct {
	nodes []*w3cNode
}

var _ NodeList = &w3cNodeList{}

func (wl *w3cNodeList) Length() int {
	return len(wl.nodes)
}

func (wl *w3cNodeList) Item(i int) Node {
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
	return domify(stytree)
}
