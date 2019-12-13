package layout

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
	reorder := prepareReorderer() // TODO
	future := styles.TopDown(createBoxForEach).Promise()
	//future := styles.TopDown(createBoxForEach).Filter(reorder).Promise()
	_, err := future()
	if err == nil {
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
	createBoxForEach := func(n *tree.Node, parent *tree.Node, i int) (*tree.Node, error) {
		//createBoxForEach := tree.Action {
		sn := stylednode{
			treenode: n,
			toStyler: toStyler,
		}
		return makeBoxNode(sn, styleToBox)
	}
	return createBoxForEach
}

func makeBoxNode(sn stylednode, styleToBox *styleToBoxAssoc) (*tree.Node, error) {
	box := boxForNode(sn.treenode, sn.toStyler)
	if box == nil { // legit, e.g. for "display:none"
		return nil, nil // will not descend to children of sn
	}
	styleToBox.Put(sn.treenode, box)
	if parent := sn.treenode.Parent(); parent != nil {
		if parentbox, ok := styleToBox.Get(sn.treenode); ok {
			parentbox.Add(box)
		}
	}
	return &box.Node, nil
}

// TODO
func reorder(*tree.Node) (*tree.Node, error) {
	return nil, nil
}

// TODO
func prepareReorderer() func(*tree.Node) (*tree.Node, error) {
	return reorder
}

/*
type nodeContext struct {
	viewport     *Container
	boxes        map[*tree.Node]*Container
	styleCreator style.Creator
}

func makeBox(sn *tree.Node, ctx *nodeContext) (*tree.Node, error) {
	//
	var c *Container
	styler := ctx.styleCreator.ToStyler(sn)
	p, err := style.GetProperty(sn, "display", ctx.styleCreator.ToStyler)
	if err == nil {
		if p == "inline" {
			c = newHBox(sn, ctx.styleCreator.ToStyler)
		} else if p == "none" {
			T().Errorf("NONE display")
			return nil, nil
		} else {
			c = newVBox(sn, ctx.styleCreator.ToStyler)
		}
		ctx.assoc[sn] = c
		// --- region flow ---------------
		p, err = style.GetProperty(sn, "flow-into", ctx.styleCreator.ToStyler) // CSS region
		if err != nil && p != style.NullStyle {
			T().Debugf("Attaching box node to region %s", p)
			nch := ctx.viewport.ChildCount()
			var flowfrom style.Property
			found := false
			for i := 0; i < nch; i++ {
				ch, ok := ctx.viewport.Child(i)
				if ok {
					rsty := ch.Payload.(*Container).styleNode // may be region style
					pmap := ctx.styleCreator.ToStyler(rsty).ComputedStyles()
					flowfrom, ok = style.GetLocalProperty(pmap, "flow-from")
					if ok && flowfrom == p {
						ch.Payload.(*Container).Add(c)
					}
				}
			}
			if !found { // then create a new region
				T().Debugf("Creating new region %s", flowfrom)
				h := ctx.styleCreator.ToStyler(ctx.viewport.styleNode).HtmlNode()
				rsty := ctx.styleCreator.StyleForHtmlNode(h)
				ctx.viewport.styleNode.AddChild(rsty)
				pmap := ctx.styleCreator.ToStyler(rsty).ComputedStyles()
				pmap.Add("flow-from", flowfrom)
				region := newVBox(rsty, ctx.styleCreator.ToStyler)
				ctx.viewport.Add(region)
				region.Add(c)
			}
		}
		// TODO skip this if box has been created for region
		// --- position --------------------
		if err == nil {
			p, err = style.GetProperty(sn, "position", ctx.styleCreator.ToStyler)
			if err == nil {
				switch p {
				case "fixed": // attach to viewport
					ctx.viewport.AddChild(&c.Node)
				case "absolute": // attach to ancestor with position ≠ static
					// TODO: I want to say 'findAncestorWith(PropIsNot("static"))'
					anc := sn.Parent()
					styler := ctx.styleCreator.ToStyler(anc)
					pp, err := style.GetLocalProperty(styler.ComputedStyles(), "position")
					for !err && (pp == "static" || pp == style.NullStyle) {
						anc = anc.Parent() // stopper is style of viewport, which has position=fixed
						styler = ctx.styleCreator.ToStyler(anc)
						pp, err = style.GetLocalProperty(styler.ComputedStyles(), "position")
					}
					ctx.assoc[anc].Add(c)
				case "relative": // attach to parent
					ctx.assoc[sn.Parent()].Add(c)
				case "static": // attach to parent
					ctx.assoc[sn.Parent()].Add(c)
				default:
					// TODO what else? table? grid + flex?
					ctx.assoc[sn.Parent()].Add(c)
				}
			}
		}
	}
	if err != nil {
		T().Errorf("Cannot create container for element %s: %v",
			styler.HtmlNode().Data, err.Error())
		return nil, err
	}
	return &c.Node, err
}
*/

// ----------------------------------------------------------------------

func boxForNode(sn *tree.Node, toStyler style.StyleInterf) *Container {
	disp := GetFormattingContextForStyledNode(sn, toStyler)
	if disp == NONE {
		return nil
	}
	var c *Container
	if disp == VBOX {
		c = newVBox(sn, toStyler)
	} else {
		c = newHBox(sn, toStyler)
	}
	c.mode = getDisplayPropertyForStyledNode(sn, toStyler)
	return c
}

// --- Default Display Properties ---------------------------------------

// GetFormattingContextForStyledNode gets the formatting context for a
// container resulting from a
// styled node. The context denotes the orientation in which a box's content
// is layed out. It may be either HBOX, VBOX or NONE.
func GetFormattingContextForStyledNode(sn *tree.Node, toStyler style.StyleInterf) uint8 {
	if sn == nil {
		return NONE
	}
	styler := toStyler(sn)
	pmap := styler.ComputedStyles()
	if val, _ := style.GetLocalProperty(pmap, "display"); val == "none" {
		return NONE
	}
	htmlnode := styler.HtmlNode()
	if htmlnode.Type != html.ElementNode {
		T().Debugf("Have styled node for non-element ?!?")
		return HBOX
	}
	switch htmlnode.Data {
	case "body", "div", "ul", "ol", "section":
		return VBOX
	case "p", "span", "it", "h1", "h2", "h3", "h4", "h5", "h6",
		"h7", "b", "i", "strong":
		return HBOX
	}
	tracing.EngineTracer.Infof("unknown HTML element %s will stack children vertically",
		htmlnode.Data)
	return VBOX
}

func getDisplayPropertyForStyledNode(sn *tree.Node, toStyler style.StyleInterf) uint8 {
	if sn == nil {
		return NONE
	}
	styler := toStyler(sn)
	pmap := styler.ComputedStyles()
	dispProp, isSet := style.GetLocalProperty(pmap, "display")
	if !isSet {
		if styler.HtmlNode().Type != html.ElementNode {
			T().Debugf("Have styled node for non-element ?!?")
			return HMODE
		}
	}
	dispProp = style.DisplayPropertyForHtmlNode(styler.HtmlNode())
	if dispProp == "none" {
		return NONE
	} else if dispProp == "block" {
		return VMODE
	} else if dispProp == "inline" {
		return HMODE
	}
	return NONE
}
