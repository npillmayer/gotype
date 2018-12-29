package layout

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/frame/box"
	"golang.org/x/net/html"
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

// We trace to the engine tracer. Will be set when calling Layout().
var T tracing.Trace

// Layout produces a render tree from walking the nodes
// of a styled tree.
func Layout(styledTree style.TreeNode, viewport *dimen.Rect) *Container {
	T = tracing.EngineTracer
	// First create the tree without calculating the dimensions
	layoutTree := buildLayoutTree(styledTree)
	// Next calculate position and dimensions for every box
	layoutBoxes(layoutTree, viewport)
	//renderTree := layoutBoxes(layoutTree, viewport)
	return nil
}

// BuildLayoutTree constructs a tree of containers and boxes from a tree
// of styled tree nodes. The boxes are not layouted yet, i.e. they have
// neither position nor size.
func buildLayoutTree(styledTree style.TreeNode) *Container {
	if styledTree == nil {
		return nil // nothing to layout
	}
	rootContainer := newVBox(nil) // outermost display context is block
	runner := &constructionRunner{rootContainer}
	runner.construct(styledTree, rootContainer) // start at root node of styled tree
	return rootContainer
}

type constructionRunner struct {
	root *Container
}

func (crun *constructionRunner) construct(sn style.TreeNode, parent *Container) {
	if sn == nil || parent == nil { // should not happen, but be safe
		return
	}
	c := boxForNode(sn, parent)
	if c != nil {
		if requiresAnonBox(c.mode, parent.orientation) {
			T.Debugf("creating anon box for %s in %s", sn.HtmlNode().Data, parent)
			if len(parent.content) > 0 &&
				parent.content[len(parent.content)-1].styleNode == nil {
				// re-use previous anon box
				prevSibling := parent.content[len(parent.content)-1]
				prevSibling.Add(c)
			} else {
				c = wrapInAnonymousBox(c)
				parent.Add(c)
			}
		} else {
			parent.Add(c)
		}
		// now recurse into styled children nodes
		chcnt := sn.ChildCount()
		for i := 0; i < chcnt; i++ {
			crun.construct(sn.Child(i), c)
		}
	}
}

func requiresAnonBox(mode uint8, orientation uint8) bool {
	return (mode == VMODE && orientation == HBOX) || (mode == HMODE && orientation == VBOX)
}

func wrapInAnonymousBox(c *Container) *Container {
	if c.mode == VMODE {
		vbox := newVBox(nil)
		vbox.Add(c)
		return vbox
	}
	hbox := newHBox(nil)
	hbox.Add(c)
	return hbox
}

// LayoutBoxes finds the positions and sizes of boxes of a previously constructed
// layout tree.
func layoutBoxes(renderTree *Container, viewport *dimen.Rect) *Container {
	// TODO
	return nil
}

type layoutRunner struct {
	root *Container
}

// Flags for box context and display mode.
const (
	VBOX uint8 = iota
	VMODE
	HBOX
	HMODE
	NONE
)

func boxForNode(sn style.TreeNode, parent *Container) *Container {
	disp := GetFormattingContextForStyledNode(sn)
	if disp == NONE {
		return nil
	}
	var c *Container
	if disp == VBOX {
		c = newVBox(sn)
	} else {
		c = newHBox(sn)
	}
	c.mode = getDisplayPropertyForStyledNode(sn)
	return c
}

// newVBox creates a container with context VBOX.
func newVBox(sn style.TreeNode) *Container {
	c := NewContainer(VBOX)
	c.styleNode = sn
	return c
}

// newHBox creates a container with context HBOX.
func newHBox(sn style.TreeNode) *Container {
	c := NewContainer(HBOX)
	c.styleNode = sn
	return c
}

// Container is a (CSS-)styled box which may contain other boxes and/or
// containers.
type Container struct {
	box.StyledBox
	//sync.RWMutex       // used to protect non-threadsafe code
	orientation uint8 // context of children V or H
	mode        uint8 // container lives in this mode
	styleNode   style.TreeNode
	content     []*Container
}

// NewContainer creates a container with orientation of either VBOX or HBOX.
func NewContainer(orientation uint8) *Container {
	return &Container{orientation: orientation}
}

func (c *Container) String() string {
	n := "_"
	if c.styleNode != nil {
		n = c.styleNode.HtmlNode().Data
	}
	b := "hbox"
	if c.orientation == VBOX {
		b = "vbox"
	}
	return fmt.Sprintf("\\%s<%s>", b, n)
}

// TODO make concurrency-safe (see styledtree.Node)
func (c *Container) Add(child *Container) {
	tracing.EngineTracer.Debugf("Adding child to %s", c)
	c.content = append(c.content, child)
}

// ----------------------------------------------------------------------

// GetFormattingContextForStyledNode gets the formatting context for a
// container resulting from a
// styled node. The context denotes the orientation in which a box's content
// is layed out. It may be either HBOX, VBOX or NONE.
func GetFormattingContextForStyledNode(sn style.TreeNode) uint8 {
	if sn == nil {
		return NONE
	}
	if val, _ := style.GetLocalProperty(sn, "display"); val == "none" {
		return NONE
	}
	if sn.HtmlNode().Type != html.ElementNode {
		T.Debugf("Have styled node for non-element ?!?")
		return HBOX
	}
	switch sn.HtmlNode().Data {
	case "body":
	case "div":
	case "ul":
	case "ol":
	case "section":
		return VBOX
	case "p":
	case "span":
	case "it":
	case "h1":
	case "h2":
	case "h3":
	case "h4":
	case "h5":
	case "h6":
	case "h7":
	case "b":
	case "i":
	case "strong":
		return HBOX
	}
	tracing.EngineTracer.Infof("unknown HTML element %s will stack children vertically",
		sn.HtmlNode().Data)
	return VBOX
}

func getDisplayPropertyForStyledNode(sn style.TreeNode) uint8 {
	if sn == nil {
		return NONE
	}
	dispProp, isSet := style.GetLocalProperty(sn, "display")
	if isSet {
		if dispProp == "None" {
			return NONE
		} else if dispProp == "block" {
			return VMODE
		} else if dispProp == "inline" {
			return HMODE
		}
	}
	if sn.HtmlNode().Type != html.ElementNode {
		T.Debugf("Have styled node for non-element ?!?")
		return HMODE
	}
	switch sn.HtmlNode().Data {
	case "body": // should not be contained
	case "p":
	case "div":
	case "ul":
	case "ol":
	case "it":
	case "section":
	case "h1":
	case "h2":
	case "h3":
	case "h4":
	case "h5":
	case "h6":
	case "h7":
		return VMODE
	case "span":
	case "b":
	case "i":
	case "strong":
		return HMODE
	}
	tracing.EngineTracer.Infof("unknown HTML element %s will be set to display: block",
		sn.HtmlNode().Data)
	return VMODE
}
