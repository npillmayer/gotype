package layout

import (
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/dom"
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

// T traces to the engine tracer.
func T() tracing.Trace {
	return gtrace.EngineTracer
}

// Layouter is a layout engine.
type Layouter struct {
	domRoot *dom.W3CNode // input styled tree
	boxRoot *Container   // layout tree to contruct
	err     error        // remember last error
	//styleCreator style.Creator // to create a style node
}

// NewLayouter creates a new layout engine for a given style tree.
// The tree's styled nodes will be accessed using styler(node).
func NewLayouter(dom *dom.W3CNode) *Layouter {
	//
	l := &Layouter{
		domRoot: dom,
		//styleCreator: creator,
	}
	return l
}

// Layout produces a render tree from walking the nodes
// of a styled tree (DOM).
func (l *Layouter) Layout(viewport *dimen.Rect) *Container {
	// First create the tree without calculating the dimensions
	var err error
	l.boxRoot, err = l.buildBoxTree()
	if err != nil {
		T().Errorf("Error building box tree")
	}
	/* TODO
	 */
	// Next calculate position and dimensions for every box
	/* TODO
	layoutBoxes(layoutTree, viewport)
	*/
	//renderTree := layoutBoxes(layoutTree, viewport)
	return l.boxRoot
}

// BoxRoot returns the root node of the render tree.
func (l *Layouter) BoxRoot() *Container {
	return l.boxRoot
}
