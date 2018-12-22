package styled

import (
	"github.com/npillmayer/gotype/engine/dom/style"
	"golang.org/x/net/html"
)

// --- Styled Node Tree -------------------------------------------------

/*
Creation of node:
- link to html node

Link (double?) when in Child
- add to children
- set parent for child

Set computed styles

Get computed styles (for a group?)

parent.findAncestorWithPropertyGroup(groupname)
need not be a member function

Tree als extra type ?
= viewport ? (normaler styled node)

type StyledNodeTree struct {
	root          *StyledNode
	defaultStyles map[string]*propertyGroup
}
*/

// StyledNodes are the building blocks of the styled tree.
type Node struct {
	node           *html.Node
	computedStyles style.PropertyMap
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
func (sn Node) ComputedStyles() style.PropertyMap {
	return sn.computedStyles
}

// Interface style.StyledNode.
func (sn *Node) SetComputedStyles(styles style.PropertyMap) {
	sn.computedStyles = styles
}

var _ style.StyledNode = &Node{}

// Factory is a trivial implementation of style.NodeFactory.
type Factory struct{}

// Interface style.NodeFactory.
func (f *Factory) NodeFor(n *html.Node) style.StyledNode {
	sn := &Node{}
	sn.node = n
	return sn
}
