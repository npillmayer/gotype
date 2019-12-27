package layout

import (
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
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
	styleroot    *tree.Node    // input styled tree
	boxroot      *Container    // layout tree to contruct
	styleCreator style.Creator // to create a style node
	err          error         // remember last error
}

// NewLayouter creates a new layout engine for a given style tree.
// The tree's styled nodes will be accessed using styler(node).
func NewLayouter(styles *tree.Node, creator style.Creator) *Layouter {
	//
	l := &Layouter{
		styleroot:    styles,
		styleCreator: creator,
	}
	return l
}

// Layout produces a render tree from walking the nodes
// of a styled tree.
func (l *Layouter) Layout(viewport *dimen.Rect) *Container {
	// First create the tree without calculating the dimensions
	boxtree, err := l.buildBoxTree()
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
	return boxtree
}
