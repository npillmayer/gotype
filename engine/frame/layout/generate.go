package layout

// This module should have knowledge about:
// - which mini-hierarchy of boxes to create for each HTML element
// - which context the element should span for its children

import (
	"fmt"

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
	if l.domRoot == nil {
		return nil, errDOMRootIsNull
	}
	domWalker := l.domRoot.Walk()
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
	action := func(node *tree.Node, parentNode *tree.Node, chpos int) (*tree.Node, error) {
		domnode, err := dom.NodeFromTreeNode(node)
		if err != nil {
			return nil, err
		}
		parent, err2 := dom.NodeFromTreeNode(parentNode)
		if err2 != nil {
			return nil, err2
		}
		return makeBoxNode(domnode, parent, chpos, dom2box)
	}
	return action
}

func makeBoxNode(domnode *dom.W3CNode, parent *dom.W3CNode, chpos int, dom2box *domToBoxAssoc) (
	*tree.Node, error) {
	//
	T().Infof("making box for %s", domnode.NodeName())
	box := BoxForNode(domnode)
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
	possiblyCreateMiniHierarchy(box)
	return &box.Node, nil
}

func requiresAnonBox(inner DisplayMode, outer DisplayMode) bool {
	return ((inner == BlockMode && outer == InlineMode) ||
		(inner == InlineMode && outer == BlockMode))
}

func wrapInAnonymousBox(c *Container) *Container {
	mode := BlockMode
	if c.outerMode == BlockMode {
		mode = InlineMode
	}
	anon := newContainer(mode, c.outerMode)
	anon.AddChild(&c.Node)
	return anon
}

// ----------------------------------------------------------------------

// BoxForNode creates an adequately initialized box for a given DOM node.
func BoxForNode(domnode *dom.W3CNode) *Container {
	if domnode == nil || domnode.HTMLNode() == nil {
		return nil
	}
	var innerMode, outerMode DisplayMode
	display := domnode.ComputedStyles().GetPropertyValue("display")
	if display.String() == "initial" {
		outerMode, innerMode = DefaultDisplayModeForHTMLNode(domnode.HTMLNode())
	} else {
		var err error
		outerMode, innerMode, err = ParseDisplay(display.String())
		if err != nil {
			T().Errorf("unrecognized display property: %s", display)
			outerMode, innerMode = BlockMode, BlockMode
		}
	}
	T().Infof("display modes = %s | %s", outerMode.String(), innerMode.String())
	if outerMode == DisplayNone {
		return nil // do not produce box for display = "none"
	}
	c := newContainer(outerMode, innerMode)
	c.DOMNode = domnode
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

// DefaultDisplayModeForHTMLNode returns the default display mode for a HTML node type,
// as described by the CSS specification.
//
// TODO possibly move this to package style (= part of browser defaults)
// If, then return a string.
func DefaultDisplayModeForHTMLNode(h *html.Node) (DisplayMode, DisplayMode) {
	if h == nil {
		return NoMode, NoMode
	}
	switch h.Type {
	case html.DocumentNode:
		return BlockMode, BlockMode
	case html.TextNode:
		return InlineMode, InlineMode
	case html.ElementNode:
		switch h.Data {
		case "table":
			return BlockMode, TableMode
		case "ul", "ol":
			return BlockMode, ListItemMode
		case "li":
			return ListItemMode, BlockMode
		case "html", "body", "div", "section", "article", "nav":
			return BlockMode, BlockMode
		case "p":
			return BlockMode, FlowMode
		case "span", "i", "b", "strong", "em":
			return InlineMode, InlineMode
		case "h1", "h2", "h3", "h4", "h5", "h6":
			return BlockMode, BlockMode
		default:
			return BlockMode, BlockMode
		}
	default:
		T().Errorf("Have styled node for non-element ?!?")
		T().Errorf(" type of node = %d", h.Type)
		T().Errorf(" name of node = %s", h.Data)
		T().Infof("unknown HTML element will stack children vertically")
		return BlockMode, BlockMode
	}
}
