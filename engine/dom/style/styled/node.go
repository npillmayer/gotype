package styled

import (
	"errors"
	"fmt"
	"strings"

	"github.com/npillmayer/gotype/engine/dom/style"
	"golang.org/x/net/html"
)

// StyledNodes are the building blocks of the styled tree.
type Node struct {
	node           *html.Node
	computedStyles *style.PropertyMap
	parent         *Node
	children       []*Node
}

// Interface style.StyledNode.
func (sn Node) Parent() style.StyledNode {
	return sn.parent
}

// LinkToParent links a styled node to another styled node.
// Interface style.StyledNode.
//
// Will panic if parent is not of type *Node.
func (sn *Node) LinkToParent(parent style.StyledNode) {
	// TODO
	// Warning: must be thread-safe
	var ok bool
	sn.parent, ok = parent.(*Node)
	if !ok {
		panic("LinkToParent: cannot link to unknown type of styled node")
	}
	sn.parent.children = append(sn.parent.children, sn) // TODO this is forbidden !
}

// Interface style.StyledNode.
func (sn Node) ComputedStyles() *style.PropertyMap {
	return sn.computedStyles
}

// Interface style.StyledNode.
func (sn *Node) SetComputedStyles(styles *style.PropertyMap) {
	sn.computedStyles = styles
}

var _ style.StyledNode = &Node{}

// Factory is a trivial implementation of style.NodeFactory.
type Factory struct{}

// Interface style.NodeFactory.
func (f Factory) NodeFor(n *html.Node) style.StyledNode {
	sn := &Node{}
	sn.node = n
	return sn
}

func (sn *Node) String() string {
	return fmt.Sprintf("styled(%s)", sn.node.Data)
}

// ----------------------------------------------------------------------

type StyledNodeQuery struct {
	styledTree *Node
	selection  Selection
}

type Selection []*Node

func (sel Selection) First() *Node {
	if len(sel) > 0 {
		return sel[0]
	}
	return nil
}

func QueryFor(sn style.StyledNode) (*StyledNodeQuery, error) {
	node, ok := sn.(*Node)
	if !ok {
		return nil, errors.New("Cannot query unknown type of styled node")
	}
	return &StyledNodeQuery{node, make([]*Node, 0, 10)}, nil
}

func (snq *StyledNodeQuery) FindElement(e string) Selection {
	snq.selection = collect(snq.styledTree, func(n *Node) bool {
		return n.node.Type == html.ElementNode && strings.EqualFold(n.node.Data, e)
	})
	return snq.selection
}

func (snq *StyledNodeQuery) FindStyledNodeFor(htmlNode *html.Node) *Node {
	sn := find(snq.styledTree, func(n *Node) bool {
		return n.node == htmlNode
	})
	if sn != nil {
		snq.selection = []*Node{sn}
		return sn
	}
	snq.selection = nil
	return nil
}

func (snq *StyledNodeQuery) findThisNode(nodeToFind *Node) *Node {
	return find(snq.styledTree, func(n *Node) bool {
		return n == nodeToFind
	})
}

// Helper to find nodes matching a predicate. Currently works recursive.
// Returns a node or nil.
func find(node *Node, matcher func(n *Node) bool) *Node {
	if node == nil {
		return nil
	}
	if matcher(node) {
		return node
	}
	for _, c := range node.children {
		if f := find(c, matcher); f != nil {
			return f
		}
	}
	return nil
}

// Helper to collect nodes matching a predicate. Currently works recursive.
func collect(node *Node, matcher func(n *Node) bool) (sel Selection) {
	if node == nil {
		return nil
	}
	if matcher(node) {
		sel = append(sel, node)
		return
	}
	for _, c := range node.children {
		sel = collect(c, matcher)
	}
	return
}
