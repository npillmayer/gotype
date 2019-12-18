package layout

// This module should have knowledge about:
// - which mini-hierarchy of boxes to create for each HTML element
// - which context the element should span for its children

import (
	"fmt"

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
	styler := l.styleCreator.ToStyler
	sn := styler(l.styleroot)
	T().Infof("ROOT node of style tree is %s", nodeTypeString(sn.HtmlNode().Type))
	style2BoxDict := newAssoc()
	createBoxForEach := prepareBoxCreator(l.styleCreator.ToStyler, style2BoxDict)
	future := styles.TopDown(createBoxForEach).Promise() // start asynchronous traversal
	_, err := future()                                   // wait for top-down traversal to finish
	if err != nil {
		return nil, err
	}
	var ok bool
	l.boxroot, ok = style2BoxDict.Get(l.styleroot)
	if !ok {
		T().Errorf("No box created for root style node")
		l.boxroot = nil
	}
	if l.boxroot != nil {
		T().Infof("ROOT BOX done!!")
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
	T().Infof("making box for %s", nodeTypeString(sn.toStyler(sn.treenode).HtmlNode().Type))
	box := boxForNode(sn.treenode, sn.toStyler)
	if box == nil { // legit, e.g. for "display:none"
		return nil, nil // will not descend to children of sn
	}
	T().Infof("assoc of %s", nodeTypeString(sn.toStyler(sn.treenode).HtmlNode().Type))
	styleToBox.Put(sn.treenode, box) // associate the styled tree node to this box
	if parent := sn.treenode.Parent(); parent != nil {
		fmt.Printf("parent is %s\n", parent)
		p, found := styleToBox.Get(parent)
		if found {
			fmt.Println("------------------>")
			fmt.Printf("adding new node to parent %s\n", p)
			if parentbox, ok := styleToBox.Get(parent); ok {
				parentbox.Add(box)
			}
			fmt.Printf("parent now has %d children\n", p.ChildCount())
			ch, _ := p.Child(0)
			c := Node(ch)
			fmt.Printf("1st child is %s\n", c)
			fmt.Println("------------------<")
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
		c = newBlockBox(sn)
	} else {
		c = newInlineBox(sn)
	}
	c.displayLevel, _ = getDisplayLevelForStyledNode(sn, toStyler)
	possiblyCreateMiniHierarchy(c, toStyler)
	return c
}

func possiblyCreateMiniHierarchy(c *Container, toStyler style.StyleInterf) {
	htmlnode := toStyler(c.StyleNode).HtmlNode()
	//propertyMap := styler.ComputedStyles()
	switch htmlnode.Data {
	case "li":
		markertype, _ := style.GetCascadedProperty(c.StyleNode, "list-style-type", toStyler)
		if markertype != "none" {
			markerbox := newInlineBox(nil)
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
	switch htmlnode.Type {
	case html.ElementNode:
		switch htmlnode.Data {
		case "html", "body", "div", "ul", "ol", "section":
			return BlockMode
		//case "p", "span", "it", "h1", "h2", "h3", "h4", "h5", "h6",
		//	"h7", "b", "i", "strong":
		//	return InlineMode
		default:
			return InlineMode
		}
	case html.DocumentNode:
		return BlockMode
	case html.TextNode:
		return InlineMode
	default:
		T().Errorf("Have styled node for non-element ?!?")
		T().Errorf("type of node = %s", nodeTypeString(htmlnode.Type))
		T().Errorf("data of node = %s", htmlnode.Data)
		tracing.EngineTracer.Infof("unknown HTML element %s will stack children vertically",
			htmlnode.Data)
		return BlockMode
	}
}
