package dom

import (
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// Node represents W3C type Node
type Node interface {
	hasAttributes() bool
	hasChildNodes() bool
}

// NodeList represents W3C type NodeList
type NodeList interface {
	Length() int
	Item(int) Node
}

// --------------------------------------------------------------------------------

type w3cNode struct {
	styledtree.StyNode
}

// NodeAsTreeNode returns the underlying tree.Node from a DOM node.
func NodeAsTreeNode(domnode Node) *tree.Node {
	w, ok := domnode.(*w3cNode)
	if !ok {
		T().Errorf("DOM node has not been created from w3cdom.go")
		return nil
	}
	return &w.Node
}

func (w *w3cNode) hasAttributes() bool {
	return false
}

func (w *w3cNode) hasChildNodes() bool {
	return false
}

// --------------------------------------------------------------------------------

// FromHTMLParseTree returns a
func FromHTMLParseTree(h *html.Node) Node {
	//h, errhtml := html.Parse(strings.NewReader(myhtml))
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
