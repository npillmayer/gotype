package dom

import (
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

type DOMNode interface {
	hasAttributes() bool
	hasChildNodes() bool
}

type DOMNodeList interface {
	Length() int
	Item(int) DOMNode
}

// --------------------------------------------------------------------------------

type w3cNode struct {
	styledtree.StyNode
}

// NodeAsTreeNode returns the underlying tree.Node from a DOM node.
func NodeAsTreeNode(domnode DOMNode) *tree.Node {
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

// W3CDOM is a data structure for a W3C Document Object Model.
type W3CDOM struct {
	styledTree *styledtree.StyNode
}

func FromHTMLParseTree(h *html.Node) *W3CDOM {
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
	dom := &W3CDOM{}
	dom.styledTree = styledtree.Node(stytree)
	return dom
}
