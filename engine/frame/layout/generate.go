package layout

// This module should have knowledge about:
// - which mini-hierarchy of boxes to create for each HTML element
// - which context the element should span for its children

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config/gtrace"

	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// Helper struct to pack a node of a styled tree.
// type stylednode struct {
// 	treenode *tree.Node
// 	toStyler style.Interf
// }

var errDOMRootIsNull = fmt.Errorf("DOM root is null")
var errDOMNodeNotSuitable = fmt.Errorf("DOM node is not suited for layout")

// TODO split into 2 runs:
//    1st: generate box nodes
//    2nd: re-order them
// otherwise there is a danger that a child overtakes a parent
//
func (l *Layouter) buildBoxTree() (*Container, error) {
	gtrace.EngineTracer.Infof("SAY HI !")
	fmt.Printf("SAY HI to Layout()\n")
	if l.domRoot == nil {
		return nil, errDOMRootIsNull
	}
	domWalker := l.domRoot.Walk()
	gtrace.EngineTracer.Errorf("SAY HI !")
	T().Debugf("Creating box tree")
	T().Infof("ROOT node of style tree is %s", dbgNodeString(l.domRoot))
	dom2box := newAssoc()
	createBoxForEach := prepareBoxCreator(dom2box)
	future := domWalker.TopDown(createBoxForEach).Promise() // start asynchronous traversal
	renderNodes, err := future()                            // wait for top-down traversal to finish
	if err != nil {
		return nil, err
	}
	T().Infof("Walker returned %d render nodes", len(renderNodes))
	for _, rnode := range renderNodes {
		n := ContainerFromNode(rnode)
		T().Infof("  node for %s", n.DOMNode.NodeName())
	}
	T().Infof("dom2box contains %d entries", dom2box.Length())
	T().Errorf("domRoot/2 = %s", dbgNodeString(l.domRoot))
	var ok bool
	l.boxRoot, ok = dom2box.Get(l.domRoot)
	T().Errorf("box for domRoot = %v", l.boxRoot)
	if !ok {
		T().Errorf("No box created for root style node")
	}
	if l.boxRoot != nil {
		T().Infof("ROOT BOX done!!")
	}
	return l.boxRoot, nil
}

func prepareBoxCreator(dict *domToBoxAssoc) tree.Action {
	dom2box := dict
	action := func(node *tree.Node, parent *tree.Node, pos int) (*tree.Node, error) {
		domnode, err := dom.NodeFromTreeNode(node)
		if err != nil {
			return nil, err
		}
		return makeBoxNode(domnode, dom2box)
	}
	return action
}

func makeBoxNode(domnode *dom.W3CNode, dom2box *domToBoxAssoc) (*tree.Node, error) {
	T().Infof("making box for %s", domnode.NodeName())
	box := boxForNode(domnode)
	if box == nil { // legit, e.g. for "display:none"
		return nil, nil // will not descend to children of sn
	}
	T().Infof("assoc of %d/%s", domnode.NodeType(), domnode.NodeName())
	dom2box.Put(domnode, box) // associate the styled tree node to this box
	if !domnode.IsDocument() {
		if parentNode := domnode.ParentNode(); parentNode != nil {
			parent := parentNode.(*dom.W3CNode)
			fmt.Printf("parent is %s\n", parent.NodeName())
			pbox, found := dom2box.Get(parent)
			if found {
				fmt.Println("------------------>")
				fmt.Printf("adding new box node to parent %s\n", pbox)
				pbox.Add(box)
				// fmt.Printf("parent now has %d children\n", parent.Children().Length())
				// ch := parent.FirstChild()
				// // fmt.Printf("1st child is %s\n", ch)
				fmt.Println("------------------<")
			}
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

func boxForNode(domnode *dom.W3CNode) *Container {
	if domnode == nil {
		return nil
	}
	context := GetFormattingContextForStyledNode(domnode)
	T().Infof("display context = %s", ModeString(context))
	if context == DisplayNone { // should not be displayed (display = none)
		return nil // do not produce box
	}
	var c *Container
	if context == BlockMode {
		c = newBlockBox(domnode)
	} else {
		c = newInlineBox(domnode)
	}
	c.displayLevel, _ = getDisplayLevelForStyledNode(domnode)
	T().Infof("display level = %s", ModeString(c.displayLevel))
	possiblyCreateMiniHierarchy(c)
	return c
}

func possiblyCreateMiniHierarchy(c *Container) {
	htmlnode := c.DOMNode.HTMLNode()
	//propertyMap := styler.ComputedStyles()
	switch htmlnode.Data {
	case "li":
		//markertype, _ := style.GetCascadedProperty(c.DOMNode, "list-style-type", toStyler)
		markertype := c.DOMNode.ComputedStyles().GetPropertyValue("list-style-type")
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
func GetFormattingContextForStyledNode(domnode *dom.W3CNode) uint8 {
	if domnode == nil {
		return NoMode
	}
	dispProp := domnode.ComputedStyles().GetPropertyValue("display")
	mode := ModeFromString(dispProp.String())
	if mode != NoMode {
		return mode
	}
	switch domnode.HTMLNode().Type {
	case html.ElementNode:
		switch domnode.NodeName() {
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
		T().Errorf(" type of node = %d", domnode.NodeType())
		T().Errorf(" name of node = %s", domnode.NodeName())
		T().Infof("unknown HTML element will stack children vertically")
		return BlockMode
	}
}
