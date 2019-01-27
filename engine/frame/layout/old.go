package layout

/*
// BuildLayoutTree constructs a tree of containers and boxes from a tree
// of styled tree nodes. The boxes are not layouted yet, i.e. they have
// neither position nor size.
func buildLayoutTree(styledTree *tree.Node) *Container {
	if styledTree == nil {
		return nil // nothing to layout
	}
	rootContainer := newVBox(nil) // outermost display context is block
	runner := &constructionRunner{rootContainer}
	runner.construct(styledTree, rootContainer) // start at root node of styled tree
	return rootContainer
}
*/

// ----------------------------------------------------------------------

/*
type constructionRunner struct {
	root *Container
}

func (crun *constructionRunner) construct(sn *tree.Node, parent *Container) {
	if sn == nil || parent == nil { // should not happen, but be safe
		return
	}
	c := boxForNode(sn, parent)
	if c != nil {
		if requiresAnonBox(c.mode, parent.orientation) {
			T().Debugf("creating anon box for %s in %s", sn.HtmlNode().Data, parent)
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
*/

/*
// LayoutBoxes finds the positions and sizes of boxes of a previously constructed
// layout tree.
func layoutBoxes(renderTree *Container, viewport *dimen.Rect) *Container {
	// TODO
	return nil
}

type layoutRunner struct {
	root *Container
}
*/
