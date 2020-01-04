package layout

import (
	"fmt"
	"strings"

	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/frame/box"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// DisplayMode is a type for CSS property "display".
type DisplayMode uint16

// Flags for box context and display mode (outer and inner).
//go:generate stringer -type=DisplayMode
const (
	NoMode       DisplayMode = iota   // unset or error condition
	DisplayNone  DisplayMode = 0x0001 // CSS outer display = none
	FlowMode     DisplayMode = 0x0002 // CSS inner display = flow
	BlockMode    DisplayMode = 0x0004 // CSS block context (inner or outer)
	InlineMode   DisplayMode = 0x0008 // CSS inline context
	ListItemMode DisplayMode = 0x0010 // CSS list-item display
	FlowRoot     DisplayMode = 0x0020 // CSS flow-root display property
	FlexMode     DisplayMode = 0x0040 // CSS inner display = flex
	GridMode     DisplayMode = 0x0080 // CSS inner display = grid
	TableMode    DisplayMode = 0x0100 // CSS table display property (inner or outer)
	ContentsMode DisplayMode = 0x0200 // CSS contents display, experimental !
)

// Container is a (CSS-)styled box which may contain other boxes or
// text.
type Container struct {
	tree.Node                // a container is a node within the layout tree
	Box       *box.StyledBox // styled box for a DOM node
	DOMNode   *dom.W3CNode   // the DOM node this Container refers to
	innerMode DisplayMode    // context of children (block or inline)
	outerMode DisplayMode    // container lives in this mode (block or inline)
	ChildInx  uint32         // this box represents child(ren) #ChildInx of the principal box
}

// newContainer creates either a block-level container or an inline-level container
func newContainer(innerMode DisplayMode, outerMode DisplayMode) *Container {
	c := &Container{
		innerMode: innerMode,
		outerMode: outerMode,
	}
	c.Payload = c // always points to itself
	return c
}

// IsPrincipal returns true if this is a principal box.
//
// Some HTML elements create a mini-hierachy of boxes for rendering. The outermost box
// is called the principal box. It will always refer to the styled node.
// An example would be an "li"-element: it will create two sub-boxes, one for the
// list item marker and one for the item's text/content. Another example are anonymous
// boxes, which will be generated for reconciling context/level-discrepancies.
func (c *Container) IsPrincipal() bool {
	return (c.DOMNode != nil)
}

// IsBlock returns true if a container has context orientation 'block'.
func (c *Container) IsBlock() bool {
	return c.innerMode == BlockMode
}

// newVBox creates a block-level container with block context.
func newBlockBox(domnode *dom.W3CNode) *Container {
	c := newContainer(BlockMode, BlockMode)
	c.DOMNode = domnode
	return c
}

// newHBox creates an inline-level container with inline context.
func newInlineBox(domnode *dom.W3CNode) *Container {
	c := newContainer(InlineMode, InlineMode)
	c.DOMNode = domnode
	return c
}

// ContainerFromNode retrieves the payload of a tree node as a Container.
// Will be called from clients as
//
//    container := layout.ContainerFromNode(n)
//
func ContainerFromNode(n *tree.Node) *Container {
	if n == nil {
		return nil
	}
	return n.Payload.(*Container)
}

func (c *Container) String() string {
	if c == nil {
		return "<empty>"
	}
	n := "_"
	b := "inline-box"
	if c.innerMode == BlockMode {
		b = "block-box"
	}
	return fmt.Sprintf("\\%s<%s>", b, n)
}

// Add appends a child container as the last sibling of existing
// child containers for c. Does nothing if child is nil.
// Will wrap the child in an anonymous container, if necessary.
//
// Returns the container itself (for chaining).
func (c *Container) Add(child *Container) *Container {
	T().Debugf("Adding child to %s", c)
	if child == nil {
		return c
	}
	if requiresAnonBox(c.innerMode, child.outerMode) {
		anon := newContainer(child.outerMode, c.innerMode)
		anon.AddChild(&child.Node)
		c.AddChild(&anon.Node)
		return c
	}
	c.AddChild(&child.Node)
	return c
}

func getDisplayLevelForStyledNode(domnode *dom.W3CNode) (DisplayMode, DisplayMode) {
	if domnode == nil {
		return NoMode, NoMode
	}
	dispProp := domnode.ComputedStyles().GetPropertyValue("display")
	if dispProp != style.NullStyle {
		T().Debugf("node %s has set display = ?", dbgNodeString(domnode), dispProp)
		inner, outer, _ := ParseDisplay(dispProp.String())
		return inner, outer
	}
	if domnode.NodeType() == html.TextNode {
		return InlineMode, InlineMode
	}
	dispProp = style.DisplayPropertyForHTMLNode(domnode.HTMLNode())
	if strings.HasPrefix(dispProp.String(), "none") {
		return DisplayNone, NoMode
	} else if strings.HasPrefix(dispProp.String(), "block") {
		return BlockMode, BlockMode
	} else if strings.HasPrefix(dispProp.String(), "inline") {
		return InlineMode, InlineMode
	}
	return NoMode, NoMode
}
