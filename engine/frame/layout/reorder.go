package layout

import (
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
)

// ReorderBoxTree reorders box nodes of a render tree to account for
// "position" CSS properties.
// In a future version, CSS regions will be supported as well.
func ReorderBoxTree(boxRoot *PrincipalBox) error {
	if boxRoot == nil {
		return nil
	}
	walker := tree.NewWalker(boxRoot.TreeNode())
	_, err := walker.DescendentsWith(absolutePositioning).AncestorWith(anchor).Promise()()
	return err
}

// Tree filter predicate: box has position "fixed" or "absolute".
func absolutePositioning(node *tree.Node, unused *tree.Node) (match *tree.Node, err error) {
	pbox := TreeNodeAsPrincipalBox(node)
	if pbox != nil {
		position := pbox.DOMNode().ComputedStyles().GetPropertyValue("position")
		if position == "fixed" || position == "absolute" {
			match = pbox.TreeNode()
		}
	}
	return
}

// Tree filter predicate with side effect: attaches node to anchor, if suited.
func anchor(test *tree.Node, node *tree.Node) (match *tree.Node, err error) {
	absoultePosChild := TreeNodeAsPrincipalBox(node)
	possibleAnchor := TreeNodeAsPrincipalBox(test)
	if absoultePosChild != nil && possibleAnchor != nil {
		var anchor *PrincipalBox
		position := absoultePosChild.DOMNode().ComputedStyles().GetPropertyValue("position")
		switch position {
		case "fixed": // searching for viewport
			if possibleAnchor.DOMNode().NodeName() == "#document" {
				anchor = possibleAnchor
			}
		case "absolute":
			anchorPos := possibleAnchor.DOMNode().ComputedStyles().GetPropertyValue("position")
			if anchorPos != style.NullStyle && anchorPos != "static" {
				anchor = possibleAnchor
			}
		}
		if anchor != nil {
			anchor.AppendChild(TreeNodeAsPrincipalBox(absoultePosChild.Isolate()))
		}
	}
	return
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
