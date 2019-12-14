package layout

// This module should have knowledge about:
// - which mini-hierarchy of boxes to create for each HTML element
// - which context the element should span for its children

import (
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
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
