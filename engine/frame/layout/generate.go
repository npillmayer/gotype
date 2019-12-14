package layout

// This module should have knowledge about:
// - which mini-hierarchy of boxes to create for each HTML element
// - which context the element should span for its children

import (
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// Helper struct to pack a node of a styled tree.
type stylednode struct {
	treenode *tree.Node
	toStyler style.StyleInterf
}

// TODO split into 2 runs:
//    1st: generate box nodes
//    2nd: re-order them
// otherwise there is a danger that a child overtakes a parent
//
func (l *Layouter) buildBoxTree() (*Container, error) {
	if l.styleroot == nil {
		return nil, nil
	}
	styles := tree.NewWalker(l.styleroot)
	T().Debugf("Creating box tree")
	styleToBox := newAssoc()
	createBoxForEach := prepareBoxCreator(l.styleCreator.ToStyler, styleToBox)
	//reorder := prepareReorderer() // TODO
	future := styles.TopDown(createBoxForEach).Promise()
	//future := styles.TopDown(createBoxForEach).Filter(reorder).Promise()
	_, err := future()
	if err != nil {
		return nil, err
	}
	var ok bool
	l.boxroot, ok = styleToBox.Get(l.styleroot)
	if !ok {
		T().Errorf("No box created for root style node")
		l.boxroot = nil
	}
	return l.boxroot, nil
}

func prepareBoxCreator(toStyler style.StyleInterf, styleToBox *styleToBoxAssoc) tree.Action {
	action := func(n *tree.Node, parent *tree.Node, i int) (*tree.Node, error) {
		sn := stylednode{
			treenode: n,
			toStyler: toStyler,
		}
		return makeBoxNode(sn, styleToBox)
	}
	return action
}

func makeBoxNode(sn stylednode, styleToBox *styleToBoxAssoc) (*tree.Node, error) {
	box := boxForNode(sn.treenode, sn.toStyler)
	if box == nil { // legit, e.g. for "display:none"
		return nil, nil // will not descend to children of sn
	}
	styleToBox.Put(sn.treenode, box) // associate the styled tree not to this box
	if parent := sn.treenode.Parent(); parent != nil {
		if parentbox, ok := styleToBox.Get(sn.treenode); ok {
			parentbox.Add(box)
		}
	}
	return &box.Node, nil
}

func requiresAnonBox(context uint8, level uint8) bool {
	return ((context == BlockMode && level == InlineMode) ||
		(context == InlineMode && level == BlockMode))
}

func wrapInAnonymousBox(c *Container) *Container {
	mode := BlockMode
	if c.displayLevel == BlockMode {
		mode = InlineMode
	}
	anon := newContainer(mode, c.contextOrientation)
	anon.AddChild(&c.Node)
	return anon
}

// ----------------------------------------------------------------------

func boxForNode(sn *tree.Node, toStyler style.StyleInterf) *Container {
	if sn == nil {
		return nil
	}
	context := GetFormattingContextForStyledNode(sn, toStyler)
	if context == DisplayNone { // should not be displayed (display = none)
		return nil
	}
	var c *Container
	if context == BlockMode {
		c = newBlockBox(sn, toStyler)
	} else {
		c = newInlineBox(sn, toStyler)
	}
	c.displayLevel, _ = getDisplayLevelForStyledNode(sn, toStyler)
	possiblyCreateMiniHierarchy(c)
	return c
}

func possiblyCreateMiniHierarchy(c *Container) {
	htmlnode := c.styleInterf(c.styleNode).HtmlNode()
	//propertyMap := styler.ComputedStyles()
	switch htmlnode.Data {
	case "li":
		markertype, _ := style.GetCascadedProperty(c.styleNode, "list-style-type", c.styleInterf)
		if markertype != "none" {
			markerbox := newInlineBox(nil, nil)
			// TODO: fill box with correct marker symbol
			c.Add(markerbox)
		}
	}
}

// --- Default Display Properties ---------------------------------------

// GetFormattingContextForStyledNode gets the formatting context for a
// container resulting from a
// styled node. The context denotes the orientation in which a box's content
// is layed out. It may be either InlineMode, BlockMode or NoMode.
func GetFormattingContextForStyledNode(sn *tree.Node, toStyler style.StyleInterf) uint8 {
	if sn == nil {
		return NoMode
	}
	styler := toStyler(sn)
	pmap := styler.ComputedStyles()
	if val, _ := style.GetLocalProperty(pmap, "display"); val == "none" {
		return DisplayNone
	}
	htmlnode := styler.HtmlNode()
	if htmlnode.Type != html.ElementNode {
		T().Errorf("Have styled node for non-element ?!?")
		return InlineMode
	}
	switch htmlnode.Data {
	case "body", "div", "ul", "ol", "section":
		return BlockMode
	//case "p", "span", "it", "h1", "h2", "h3", "h4", "h5", "h6",
	//	"h7", "b", "i", "strong":
	//	return InlineMode
	default:
		return InlineMode
	}
	tracing.EngineTracer.Infof("unknown HTML element %s will stack children vertically",
		htmlnode.Data)
	return BlockMode
}
