package layout

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/frame/box"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// Flags for box context and display level.
const (
	NoMode uint8 = iota
	BlockMode
	InlineMode
)

// Container is a (CSS-)styled box which may contain other boxes and/or
// containers.
type Container struct {
	tree.Node
	Box                box.StyledBox
	contextOrientation uint8      // context of children (block or inline)
	levelMode          uint8      // container lives in this mode (block or inline)
	styleNode          *tree.Node // the DOM node this Container refers to
	styleInterf        style.StyleInterf
}

// newContainer creates either a block-level container or an inline-level container
func newContainer(orientation uint8, level uint8) *Container {
	c := &Container{
		contextOrientation: orientation,
		levelMode:          level,
	}
	c.Payload = c // always points to itself
	return c
}

// isPrincipal is a predicate wether this is a principal box.
func (c *Container) isPrincipal() bool {
	return (c.styleNode != nil)
}

// newVBox creates a block-level container with block context.
func newBlockBox(sn *tree.Node, styleInterf style.StyleInterf) *Container {
	c := newContainer(BlockMode, BlockMode)
	c.styleNode = sn
	c.styleInterf = styleInterf
	return c
}

// newHBox creates an inline-level container with inline context.
func newInlineBox(sn *tree.Node, styleInterf style.StyleInterf) *Container {
	c := newContainer(InlineMode, InlineMode)
	c.styleNode = sn
	c.styleInterf = styleInterf
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
	if c.contextOrientation == BlockMode {
		b = "vbox"
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
	if requiresAnonBox(c.contextOrientation, child.levelMode) {
		anon := newContainer(child.levelMode, c.contextOrientation)
		anon.AddChild(&child.Node)
		c.AddChild(&anon.Node)
		return c
	}
	c.AddChild(&child.Node)
	return c
}

func requiresAnonBox(context uint8, level uint8) bool {
	return ((context == BlockMode && level == InlineMode) ||
		(context == InlineMode && level == BlockMode))
}

func wrapInAnonymousBox(c *Container) *Container {
	mode := BlockMode
	if c.levelMode == BlockMode {
		mode = InlineMode
	}
	anon := newContainer(mode, c.contextOrientation)
	anon.AddChild(&c.Node)
	return anon
}

// ----------------------------------------------------------------------

func boxForNode(sn *tree.Node, toStyler style.StyleInterf) *Container {
	context := GetFormattingContextForStyledNode(sn, toStyler)
	if context == NoMode {
		return nil
	}
	var c *Container
	if context == BlockMode {
		c = newBlockBox(sn, toStyler)
	} else {
		c = newInlineBox(sn, toStyler)
	}
	c.levelMode = getLevelModeForStyledNode(sn, toStyler)
	return c
}

// --- Default Display Properties ---------------------------------------

// GetFormattingContextForStyledNode gets the formatting context for a
// container resulting from a
// styled node. The context denotes the orientation in which a box's content
// is layed out. It may be either HBOX, VBOX or NONE.
func GetFormattingContextForStyledNode(sn *tree.Node, toStyler style.StyleInterf) uint8 {
	if sn == nil {
		return NoMode
	}
	styler := toStyler(sn)
	pmap := styler.ComputedStyles()
	if val, _ := style.GetLocalProperty(pmap, "display"); val == "none" {
		return NoMode
	}
	htmlnode := styler.HtmlNode()
	if htmlnode.Type != html.ElementNode {
		T().Debugf("Have styled node for non-element ?!?")
		return InlineMode
	}
	switch htmlnode.Data {
	case "body", "div", "ul", "ol", "section":
		return BlockMode
	case "p", "span", "it", "h1", "h2", "h3", "h4", "h5", "h6",
		"h7", "b", "i", "strong":
		return InlineMode
	}
	tracing.EngineTracer.Infof("unknown HTML element %s will stack children vertically",
		htmlnode.Data)
	return BlockMode
}

func getLevelModeForStyledNode(sn *tree.Node, toStyler style.StyleInterf) uint8 {
	if sn == nil {
		return NoMode
	}
	styler := toStyler(sn)
	pmap := styler.ComputedStyles()
	dispProp, isSet := style.GetLocalProperty(pmap, "display")
	if !isSet {
		if styler.HtmlNode().Type != html.ElementNode {
			T().Debugf("Have styled node for non-element ?!?")
			return InlineMode
		}
	}
	dispProp = style.DisplayPropertyForHtmlNode(styler.HtmlNode())
	if dispProp == "none" {
		return NoMode
	} else if dispProp == "block" {
		return BlockMode
	} else if dispProp == "inline" {
		return InlineMode
	}
	return NoMode
}
