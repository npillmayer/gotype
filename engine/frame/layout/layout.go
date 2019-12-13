package layout

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/frame/box"
	"github.com/npillmayer/gotype/engine/tree"
)

// TODO:
// - create a "child-iterator" in package dom: iterate through CData & Element
//   children. Element children are accessible from their styled node.
// - correctly implement anonymous boxes (see link below).

// Invaluable:
// https://developer.mozilla.org/en-US/docs/Web/CSS/Visual_formatting_model
//
// Regions:
// http://cna.mamk.fi/Public/FJAK/MOAC_MTA_HTML5_App_Dev/c06.pdf

// T traces to the engine tracer.
func T() tracing.Trace {
	return tracing.EngineTracer
}

// Layouter is a layout engine.
type Layouter struct {
	styleroot    *tree.Node    // input styled tree
	boxroot      *Container    // layout tree to contruct
	styleCreator style.Creator // to create a style node
	err          error         // remember last error
}

// NewLayouter creates a new layout engine for a given style tree.
// The tree's styled nodes will be accessed using styler(node).
func NewLayouter(styles *tree.Node, creator style.Creator) *Layouter {
	//
	l := &Layouter{
		styleroot:    styles,
		styleCreator: creator,
	}
	return l
}

// Layout produces a render tree from walking the nodes
// of a styled tree.
func (l *Layouter) Layout(viewport *dimen.Rect) *Container {
	// First create the tree without calculating the dimensions
	/* TODO
	layoutTree := buildLayoutTree(styledTree)
	*/
	// Next calculate position and dimensions for every box
	/* TODO
	layoutBoxes(layoutTree, viewport)
	*/
	//renderTree := layoutBoxes(layoutTree, viewport)
	return nil
}

func requiresAnonBox(mode uint8, orientation uint8) bool {
	return (mode == VMODE && orientation == HBOX) || (mode == HMODE && orientation == VBOX)
}

func wrapInAnonymousBox(c *Container) *Container {
	if c.mode == VMODE {
		vbox := newVBox(nil, c.styleInterf)
		vbox.Add(c)
		return vbox
	}
	hbox := newHBox(nil, c.styleInterf)
	hbox.Add(c)
	return hbox
}

// Flags for box context and display mode.
const (
	VBOX uint8 = iota
	VMODE
	HBOX
	HMODE
	NONE
)

// newVBox creates a container with context VBOX.
//
// Sometimes I miss subclassing dearly.
func newVBox(sn *tree.Node, styleInterf style.StyleInterf) *Container {
	c := newContainer(VBOX)
	c.styleNode = sn
	c.styleInterf = styleInterf
	return c
}

// newHBox creates a container with context HBOX.
//
// Sometimes I miss subclassing dearly.
func newHBox(sn *tree.Node, styleInterf style.StyleInterf) *Container {
	c := newContainer(HBOX)
	c.styleNode = sn
	c.styleInterf = styleInterf
	return c
}

// Container is a (CSS-)styled box which may contain other boxes and/or
// containers.
type Container struct {
	tree.Node
	Box         box.StyledBox
	orientation uint8 // context of children V or H
	mode        uint8 // container lives in this mode
	styleNode   *tree.Node
	styleInterf style.StyleInterf
}

// newContainer creates a container with orientation of either VBOX or HBOX.
func newContainer(orientation uint8) *Container {
	c := &Container{orientation: orientation}
	c.Payload = c // always points to itself
	return c
}

// Node gets the payload of a tree node as a Container.
// Called from clients as
//
//    container := layout.Node(n)
//
func Node(n *tree.Node) *Container {
	if n == nil {
		return nil
	}
	return n.Payload.(*Container)
}

func (c *Container) String() string {
	n := "_"
	if c.styleNode != nil {
		n = c.styleInterf(c.styleNode).HtmlNode().Data
	}
	b := "hbox"
	if c.orientation == VBOX {
		b = "vbox"
	}
	return fmt.Sprintf("\\%s<%s>", b, n)
}

// Add appends a child container as the last sibling of existing
// child containers for c. Does nothing if child is nil.
func (c *Container) Add(child *Container) {
	T().Debugf("Adding child to %s", c)
	if child == nil {
		return
	}
	c.AddChild(&child.Node)
}
