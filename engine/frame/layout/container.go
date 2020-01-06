package layout

import (
	"bytes"
	"fmt"

	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/dom/w3cdom"
	"github.com/npillmayer/gotype/engine/frame/box"
	"github.com/npillmayer/gotype/engine/tree"
)

// DisplayMode is a type for CSS property "display".
type DisplayMode uint16

// Flags for box context and display mode (outer and inner).
//go:generate stringer -type=DisplayMode
const (
	NoMode       DisplayMode = iota   // unset or error condition
	DisplayNone  DisplayMode = 0x0001 // CSS outer display = none
	FlowMode     DisplayMode = 0x0002 // CSS inner display = flow
	BlockMode    DisplayMode = 0x0004 // CSS block context (inner or outer)
	InlineMode   DisplayMode = 0x0008 // CSS inline context
	ListItemMode DisplayMode = 0x0010 // CSS list-item display
	FlowRoot     DisplayMode = 0x0020 // CSS flow-root display property
	FlexMode     DisplayMode = 0x0040 // CSS inner display = flex
	GridMode     DisplayMode = 0x0080 // CSS inner display = grid
	TableMode    DisplayMode = 0x0100 // CSS table display property (inner or outer)
	ContentsMode DisplayMode = 0x0200 // CSS contents display mode, experimental !
)

var allDisplayModes = []DisplayMode{
	DisplayNone, FlowMode, BlockMode, InlineMode, ListItemMode, FlowRoot, FlexMode,
	GridMode, TableMode, ContentsMode,
}

// Set sets a given atomic mode within this display mode.
func (disp *DisplayMode) Set(d DisplayMode) {
	*disp = (*disp) | d
}

// Contains checks if a display mode contains a given atomic mode.
// Returns false for d = NoMode.
func (disp DisplayMode) Contains(d DisplayMode) bool {
	return d != NoMode && (disp&d > 0)
}

// Overlaps returns true if a given display mode shares at least one atomic
// mode flag with disp (excluding NoMode).
func (disp DisplayMode) Overlaps(d DisplayMode) bool {
	for _, m := range allDisplayModes {
		if disp.Contains(m) && d.Contains(m) {
			return true
		}
	}
	return false
}

// FullString returns all atomic modes set in a display mode.
func (disp DisplayMode) FullString() string {
	var b bytes.Buffer
	first := true
	for _, m := range allDisplayModes {
		if disp.Contains(m) {
			if !first {
				b.WriteString(" ")
			}
			first = false
			b.WriteString(m.String())
		}
	}
	return b.String()
}

// Symbol returns a Unicode symbol for a mode.
func (disp DisplayMode) Symbol() string {
	if disp == FlowMode {
		return "\u25a7"
	} else if disp.Contains(BlockMode) {
		return "\u25a9"
	} else if disp.Contains(InlineMode) {
		return "\u25ba"
	} else if disp.Contains(FlexMode) {
		return "\u25a4"
	} else if disp.Contains(GridMode) {
		return "\u25f0"
	} else if disp.Contains(ListItemMode) {
		return "\u25a3"
	} else if disp.Contains(TableMode) {
		return "\u25a5"
	}
	return "?"
}

// --- Container -----------------------------------------------------------------------

// Container is an interface type for render tree nodes, i.e., boxes.
type Container interface {
	DOMNode() w3cdom.Node
	TreeNode() *tree.Node
	IsAnonymous() bool
	IsText() bool
	DisplayModes() (DisplayMode, DisplayMode)
	ChildIndices() (uint32, uint32)
}

var _ Container = &PrincipalBox{}
var _ Container = &AnonymousBox{}
var _ Container = &TextBox{}

// --- PrincipalBox --------------------------------------------------------------------

// PrincipalBox is a (CSS-)styled box which may contain other boxes.
// It references a node in the styled tree, i.e., a stylable DOM element node.
type PrincipalBox struct {
	tree.Node                // a container is a node within the layout tree
	Box       *box.StyledBox // styled box for a DOM node
	domNode   *dom.W3CNode   // the DOM node this PrincipalBox refers to
	outerMode DisplayMode    // container lives in this mode (block or inline)
	innerMode DisplayMode    // context of children (block or inline)
	ChildInx  uint32         // this box represents child #ChildInx of the parent principal box
	anonMask  runlength      // mask for anonymous box children
}

// newPrincipalBox creates either a block-level container or an inline-level container
func newPrincipalBox(domnode *dom.W3CNode, outerMode DisplayMode, innerMode DisplayMode) *PrincipalBox {
	pbox := &PrincipalBox{
		domNode:   domnode,
		outerMode: outerMode,
		innerMode: innerMode,
	}
	pbox.Payload = pbox // always points to itself: tree node -> box
	return pbox
}

// TreeNodeAsPrincipalBox retrieves the payload of a tree node as a PrincipalBox.
// Will be called from clients as
//
//    box := layout.PrincipalBoxFromNode(n)
//
func TreeNodeAsPrincipalBox(n *tree.Node) *PrincipalBox {
	if n == nil {
		return nil
	}
	pbox, ok := n.Payload.(*PrincipalBox)
	if ok {
		return pbox
	}
	return nil
}

// TreeNode returns the underlying tree node for a box.
func (pbox *PrincipalBox) TreeNode() *tree.Node {
	return &pbox.Node
}

// DOMNode returns the underlying DOM node for a render tree element.
func (pbox *PrincipalBox) DOMNode() w3cdom.Node {
	return pbox.domNode
}

// IsPrincipal returns true if this is a principal box.
//
// Some HTML elements create a mini-hierachy of boxes for rendering. The outermost box
// is called the principal box. It will always refer to the styled node.
// An example would be an "li"-element: it will create two sub-boxes, one for the
// list item marker and one for the item's text/content. Another example are anonymous
// boxes, which will be generated for reconciling context/level-discrepancies.
func (pbox *PrincipalBox) IsPrincipal() bool {
	return (pbox.domNode != nil)
}

// IsAnonymous will always return false for a container.
func (pbox *PrincipalBox) IsAnonymous() bool {
	return false
}

// IsText will always return false for a principal box.
func (pbox *PrincipalBox) IsText() bool {
	return false
}

// DisplayModes returns outer and inner display mode of this box.
func (pbox *PrincipalBox) DisplayModes() (DisplayMode, DisplayMode) {
	return pbox.outerMode, pbox.innerMode
}

func (pbox *PrincipalBox) String() string {
	if pbox == nil {
		return "<empty box>"
	}
	name := pbox.DOMNode().NodeName()
	innerSym := pbox.innerMode.Symbol()
	outerSym := pbox.outerMode.Symbol()
	return fmt.Sprintf("%s %s %s", outerSym, innerSym, name)
}

// ChildIndices returns the positional index of this box reference to
// the parent principal box. To comply with the PrincipalBox interface, it returns
// the index twice (from, to).
func (pbox *PrincipalBox) ChildIndices() (uint32, uint32) {
	return pbox.ChildInx, pbox.ChildInx
}

func (pbox *PrincipalBox) prepareAnonymousBoxes() {
	if pbox.domNode.HasChildNodes() {
		if pbox.innerMode.Contains(InlineMode) {
			// In inline mode all block-children have to be wrapped in an anon box.
			blockChPos := pbox.checkForChildrenWithDisplayMode(BlockMode)
			if !blockChPos.Empty() { // yes, found
				// At least one block child present => need anon box for block children
				pbox.anonMask = blockChPos
				anonpos := blockChPos.Condense()
				for i, intv := range blockChPos {
					anon := newAnonymousBox(InlineMode, BlockMode)
					anon.ChildInxFrom = intv.from
					anon.ChildInxTo = intv.from + intv.len - 1
					pbox.SetChildAt(int(anonpos[i]), anon.TreeNode())
				}
			}
		}
		if pbox.innerMode.Contains(BlockMode) {
			// In flow mode all children must have the same outer display mode,
			// either block or inline.
			// TODO This holds for flow and grid, too ?! others?
			inlineChPos := pbox.checkForChildrenWithDisplayMode(InlineMode)
			if !(pbox.checkForChildrenWithDisplayMode(BlockMode).Empty() ||
				inlineChPos.Empty()) { // found both
				// Both inline and block children => need anon boxes for inline children
				T().Debugf("Creating inline anon boxes at %s", inlineChPos)
				pbox.anonMask = inlineChPos
				anonpos := inlineChPos.Condense()
				for i, intv := range inlineChPos {
					anon := newAnonymousBox(BlockMode, InlineMode)
					anon.ChildInxFrom = intv.from
					anon.ChildInxTo = intv.from + intv.len - 1
					pbox.SetChildAt(int(anonpos[i]), anon.TreeNode())
				}
			}
		}
	}
}

func (pbox *PrincipalBox) checkForChildrenWithDisplayMode(dispMode DisplayMode) runlength {
	domchildren := pbox.domNode.ChildNodes()
	var rl runlength
	var openintv intv
	for i := 0; i < domchildren.Length(); i++ {
		domchild := domchildren.Item(i).(*dom.W3CNode)
		outerMode, _ := DisplayModesForDOMNode(domchild)
		if outerMode.Overlaps(dispMode) {
			if openintv != nullintv {
				openintv.len++
			} else {
				openintv = intv{uint32(i), uint32(1)}
			}
		} else {
			if openintv.len > 0 {
				rl = append(rl, openintv)
			}
			openintv = nullintv
		}
	}
	if openintv.len > 0 {
		rl = append(rl, openintv)
	}
	return rl
}

// ErrNullChild flags an error condition when a non-nil child has been expected.
var ErrNullChild = fmt.Errorf("Child box max not be null")

// ErrAnonBoxNotFound flags an error condition where an anonymous box should be
// present but could not be found.
var ErrAnonBoxNotFound = fmt.Errorf("No anonymous box found for index")

// AddChild appends a child box to its parent principal box.
// The child is a principal box itself, i.e. references a styleable DOM node.
// The child must have its child index set.
func (pbox *PrincipalBox) AddChild(child *PrincipalBox) error {
	return pbox.addChildContainer(child)
}

// AddTextChild appends a child box to its parent principal box.
// The child is a text box, i.e., references a HTML text node.
// The child must have its child index set.
func (pbox *PrincipalBox) AddTextChild(child *TextBox) error {
	err := pbox.addChildContainer(child)
	// if err == nil {
	// 	if pbox.innerMode.Contains(InlineMode) {
	// 		child.outerMode.Set(InlineMode)
	// 	} else if pbox.innerMode.Contains(BlockMode) {
	// 		child.outerMode.Set(BlockMode)
	// 	}
	// }
	return err
}

func (pbox *PrincipalBox) addChildContainer(child Container) error {
	if child == nil {
		return ErrNullChild
	}
	inx, _ := child.ChildIndices()
	anon, ino, j := pbox.anonMask.Translate(inx)
	T().Debugf("Anon mask of %s is %s, transl child #%d to %v->(%d,%d)",
		pbox.String(), pbox.anonMask, inx, anon, ino, j)
	var node *tree.Node
	var ok bool
	if anon {
		// we will add the child to an anonymous box
		node, ok = pbox.TreeNode().Child(int(ino))
		if !ok { // oops, we expected an anonymous box there
			return ErrAnonBoxNotFound
		}
	} else {
		node = pbox.TreeNode() // we will add the child to the principal box
	}
	node.SetChildAt(int(j), child.TreeNode())
	return nil
}

// AppendChild appends a child box to a principal box.
// The child is a principal box itself, i.e. references a styleable DOM node.
// It is appended as the last child of pbox.
func (pbox *PrincipalBox) AppendChild(child *PrincipalBox) {
	if !pbox.innerMode.Overlaps(child.outerMode) {
		// create an anon box
		anon := newAnonymousBox(pbox.innerMode, child.outerMode)
		anon.TreeNode().AddChild(child.TreeNode())
		pbox.TreeNode().AddChild(anon.TreeNode())
		return
	}
	pbox.TreeNode().AddChild(child.TreeNode())
}

// --- Anonymous Boxes -----------------------------------------------------------------

// AnonymousBox is a type for CSS anonymous boxes.
//
// From the spec: "If a container box (inline or block) has a block-level box inside it,
// then we force it to have only block-level boxes inside it."
//
// These block-level boxes are anonymous boxes. There are anonymous inline-level boxes,
// too. Both are not directly stylable by the user, but rather inherit the styles of
// their principal boxes.
type AnonymousBox struct {
	tree.Node                // an anonymous box is a node within the layout tree
	Box          *box.Box    // an anoymous box cannot be styled
	outerMode    DisplayMode // container lives in this mode (block or inline)
	innerMode    DisplayMode // context of children (block or inline)
	ChildInxFrom uint32      // this box represents children starting at #ChildInxFrom of the principal box
	ChildInxTo   uint32      // this box represents children to #ChildInxTo
}

// DOMNode returns the underlying DOM node for a render tree element.
func (anon *AnonymousBox) DOMNode() w3cdom.Node {
	return nil
}

// TreeNode returns the underlying tree node for a box.
func (anon *AnonymousBox) TreeNode() *tree.Node {
	return &anon.Node
}

// IsAnonymous will always return true for an anonymous box.
func (anon *AnonymousBox) IsAnonymous() bool {
	return true
}

// IsText will always return false for an anonymous box.
func (anon *AnonymousBox) IsText() bool {
	return false
}

// DisplayModes returns outer and inner display mode of this box.
func (anon *AnonymousBox) DisplayModes() (DisplayMode, DisplayMode) {
	return anon.outerMode, anon.innerMode
}

func (anon *AnonymousBox) String() string {
	if anon == nil {
		return "<empty anon box>"
	}
	innerSym := anon.innerMode.Symbol()
	outerSym := anon.outerMode.Symbol()
	return fmt.Sprintf("%s %s", outerSym, innerSym)
}

// ChildIndices returns the positional indices of all child-boxes in reference to
// the principal box.
func (anon *AnonymousBox) ChildIndices() (uint32, uint32) {
	return anon.ChildInxFrom, anon.ChildInxTo
}

func newAnonymousBox(outer DisplayMode, inner DisplayMode) *AnonymousBox {
	anon := &AnonymousBox{
		outerMode: outer,
		innerMode: inner,
	}
	anon.Payload = anon // always points to itself: tree node -> box
	return anon
}

// --- Anonymous Boxes -----------------------------------------------------------------

// TextBox is a type for CSS inline text boxes.
// It references a text node in the DOM.
// They are not directly stylable by the user, but rather inherit the styles of
// their principal boxes. Text boxes have an inner display type of inline.
type TextBox struct {
	tree.Node              // a text box is a node within the layout tree
	Box       *box.Box     // text box cannot be explicitely styled
	domNode   *dom.W3CNode // the DOM text-node this box refers to
	//outerMode DisplayMode  // container lives in this mode (block or inline)
	ChildInx uint32 // this box represents a text node at #ChildInx of the principal box
}

func newTextBox(domnode *dom.W3CNode) *TextBox {
	tbox := &TextBox{
		domNode: domnode,
		//outerMode: FlowMode,
	}
	tbox.Payload = tbox // always points to itself: tree node -> box
	return tbox
}

// DOMNode returns the underlying DOM node for a render tree element.
func (tbox *TextBox) DOMNode() w3cdom.Node {
	return tbox.domNode
}

// TreeNode returns the underlying tree node for a box.
func (tbox *TextBox) TreeNode() *tree.Node {
	return &tbox.Node
}

// IsAnonymous will always return true for a text box.
func (tbox *TextBox) IsAnonymous() bool {
	return true
}

// IsText will always return true for a text box.
func (tbox *TextBox) IsText() bool {
	return true
}

// DisplayModes always returns inline.
func (tbox *TextBox) DisplayModes() (DisplayMode, DisplayMode) {
	return InlineMode, InlineMode
}

// ChildIndices returns the positional index of the text node in reference to
// the principal box. To comply with the PrincipalBox interface, it returns
// the index twice (from, to).
func (tbox *TextBox) ChildIndices() (uint32, uint32) {
	return tbox.ChildInx, tbox.ChildInx
}

// ----------------------------------------------------------------------------------

type runlength []intv // a list of intervals
type intv struct {    // run-length interval
	from, len uint32
}

var nullintv = intv{} // null-type for intervals

func (rl runlength) Empty() bool {
	return len(rl) == 0
}

// Condense returns a list of positions, where every interval of rl is counted
// as a single position. This gives positional indices for anonymous boxes
// associated with the intervals, usable as indices in the parents child-vector.
func (rl runlength) Condense() (positions []uint32) {
	if rl == nil {
		return positions
	}
	pos := uint32(0)
	next := uint32(0)
	for _, intv := range rl {
		if intv.from > pos {
			for j := pos; j < intv.from; j++ {
				next++
			}
		}
		positions = append(positions, next)
		next++
		pos = intv.from + intv.len
	}
	return positions
}

// Translate takes an input index (of a child node) and returns the real
// position. The boolean return value is true, if the input index lies within
// one of the intervals of rl, otherwise false.
func (rl runlength) Translate(inx uint32) (bool, uint32, uint32) {
	if rl == nil {
		return false, 0, inx // nothing to translate
	}
	last := uint32(0) // max input index processed + 1
	pos := uint32(0)  // next possible output index
	for _, intv := range rl {
		if inx < intv.from { // inx is left of this interval
			pos = pos + inx - last
			return false, uint32(0), pos
		}
		if inx <= intv.from+intv.len-1 { // inx is in this interval
			return true, pos + intv.from - last, inx - intv.from
		}
		// account for positions including the current interval
		pos = pos + intv.from - last + 1
		last = intv.from + intv.len
	}
	// inx is to the right of the last interval
	return false, uint32(0), pos + inx - last
}

func (rl runlength) String() string {
	var b bytes.Buffer
	b.WriteString("(")
	for _, iv := range rl {
		if iv.len == 0 {
			b.WriteString(" []")
		} else {
			b.WriteString(fmt.Sprintf(" [%d..%d]", iv.from, iv.from+iv.len-1))
		}
	}
	b.WriteString(" )")
	return b.String()
}
