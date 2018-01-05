package gtcore

import (
	"bytes"
	"fmt"
	"math"
)

/*
---------------------------------------------------------------------------

BSD License
Copyright (c) 2017, Norbert Pillmayer

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. Neither the name of Norbert Pillmayer nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----------------------------------------------------------------------

/* Nodes implement items for typesetting paragraphs. The various node
 * types more or less implement the corresponding types from the TeX
 * typesetting system.

*/

// === Nodes =================================================================

// A node has a width and may be discardable
type Node interface {
	fmt.Stringer
	W() Dimen            // width
	MinW() Dimen         // minimum width
	MaxW() Dimen         // maximum width
	IsDiscardable() bool // is this node discardable?
}

// Node types
const (
	NTKern int = iota
	NTGlue
	NTBox
	NTPenalty
	NTDiscretionary
)

/* Factory method to create a node. Parameter is a valid node type.
 */
func NewNode(nodetype int) Node {
	switch nodetype {
	case NTKern:
		return &Kern{0}
	case NTGlue:
		return &Glue{}
	case NTPenalty:
		return &Penalty{0}
	case NTDiscretionary:
		d := &Discretionary{}
		//d.pre = "-"
		return d
	case NTBox:
		box := &Box{}
		return box
	}
	return nil
}

// --- Kern ------------------------------------------------------------------

// A kern is an unshrinkable space
type Kern struct {
	Width Dimen // fixed width
}

/* Interface Node. Prints the dimension (width) of the kern.
 */
func (k *Kern) String() string {
	return fmt.Sprintf("[kern %s]", k.Width.String())
}

/* Interface Node. Width of the kern.
 */
func (k *Kern) W() Dimen {
	return k.Width
}

/* Interface Node. Kerns do not shrink.
 */
func (k *Kern) MinW() Dimen {
	return k.Width
}

/* Interface Node. Kerns do not stretch.
 */
func (k *Kern) MaxW() Dimen {
	return k.Width
}

/* Interface Node. Kerns are discardable.
 */
func (k *Kern) IsDiscardable() bool {
	return true
}

var _ Node = &Kern{}

// --- Glue ------------------------------------------------------------------

// A glue is a space which can shrink and expand
type Glue struct {
	Width    Dimen // natural width
	MaxWidth Dimen // maximum width
	MinWidth Dimen // minimum width
}

/* Interface Node.
 */
func (g *Glue) String() string {
	return fmt.Sprintf("[glue %s <%s >%s]", g.W().String(), g.MaxW().String(),
		g.MinW().String())
}

/* Interface Node. Natural width of the glue.
 */
func (g *Glue) W() Dimen {
	return g.Width
}

/* Interface Node. Minimum width of the glue.
 */
func (g *Glue) MinW() Dimen {
	return g.MinWidth
}

/* Interface Node. Maximum width of the glue.
 */
func (g *Glue) MaxW() Dimen {
	return g.MaxWidth
}

/* Interface Node. Glue is discardable.
 */
func (g *Glue) IsDiscardable() bool {
	return true
}

var _ Node = &Glue{}

/* Create a new drop of glue with stretch and shrink.
 */
func NewGlue(w Dimen, stretch Dimen, shrink Dimen) *Glue {
	glue := NewNode(NTGlue).(*Glue)
	glue.Width = w
	glue.MaxWidth = w + stretch
	glue.MinWidth = w - shrink
	return glue
}

/* Create a drop of infinitely stretchable glue.
 */
func NewFill(f int) *Glue {
	var stretch Dimen
	switch f {
	case 2:
		stretch = Fill
	case 3:
		stretch = Fill
	default:
		stretch = Fil
	}
	return NewGlue(0, stretch, 0)
}

// --- Discretionary ---------------------------------------------------------

// A discretionary is a hyphenation opportunity
type Discretionary struct {
	nobreak Box // text if not hyphenated
	pre     Box // pre-hyphen text
	post    Box // post-hyphen text
}

/* Interface Node. Prints the dimension (width) of the kern.
 */
func (d *Discretionary) String() string {
	return fmt.Sprintf("\\discretionary{%s}{%s}{%s}", d.nobreak.text,
		d.pre.text, d.post.text)
}

/* Interface Node. Returns the width of the un-hyphenated text.
 */
func (d *Discretionary) W() Dimen {
	return d.nobreak.W()
}

/* Interface Node. Returns the width of the pre-hyphen text.
 */
func (d *Discretionary) MinW() Dimen {
	return d.pre.W()
}

/* Interface Node. Returns the width of the post-hyphen text.
 */
func (d *Discretionary) MaxW() Dimen {
	return d.post.W()
}

/* Interface Node. Discretionaries are not discardable.
 */
func (d *Discretionary) IsDiscardable() bool {
	return false
}

var _ Node = &Discretionary{}

// --- Boxes -----------------------------------------------------------------

// A Box is a fixed unit of text
type Box struct {
	Width  Dimen  // width
	Height Dimen  // height
	Depth  Dimen  // depth
	text   string // text, if available
	//nodelist Nodelist // content, if available
}

/* Interface Node.
 */
func (b *Box) String() string {
	return fmt.Sprintf("\\box{%s}", b.text)
}

/* Interface Node. Width of the glue.
 */
func (b *Box) W() Dimen {
	return b.Width
}

/* Interface Node. Width of the glue.
 */
func (b *Box) MinW() Dimen {
	return b.Width
}

/* Interface Node. Width of the glue.
 */
func (b *Box) MaxW() Dimen {
	return b.Width
}

/* Interface Node. Glue is discardable.
 */
func (b *Box) IsDiscardable() bool {
	return false
}

var _ Node = &Box{}

// --- Penalty ---------------------------------------------------------------

// A penalty contributes to demerits, i.e. the quality index of paragraphs
type Penalty struct {
	P int
}

/* Interface Node.
 */
func (p *Penalty) String() string {
	return fmt.Sprintf("[penalty %d]", p.P)
}

/* Interface Node. Returns 0.
 */
func (p *Penalty) W() Dimen {
	return 0
}

/* Interface Node. Returns 0.
 */
func (p *Penalty) MinW() Dimen {
	return 0
}

/* Interface Node. Returns 0.
 */
func (p *Penalty) MaxW() Dimen {
	return 0
}

/* Interface Node. Penalties are discardable.
 */
func (p *Penalty) IsDiscardable() bool {
	return true
}

var _ Node = &Penalty{}

// === Node lists ============================================================

// We handle text/paragraphcs as a list of nodes
type Nodelist struct {
	Listtype int    // hlist, vlist or mlist
	Nodes    []Node // array of nodes of different type
}

// List types
const (
	HList int = iota
	VList
	MList
)

/* Create a new node list.
 */
func NewNodelist() *Nodelist {
	nl := &Nodelist{}
	nl.Nodes = make([]Node, 0, 50)
	return nl
}

/* Number of nodes in the list.
 */
func (nl *Nodelist) Length() int {
	return len(nl.Nodes)
}

/* Append a node at the end of the list.
 */
func (nl *Nodelist) AppendNode(node Node) *Nodelist {
	nl.Nodes = append(nl.Nodes, node)
	return nl
}

/* Return the widths of a subset of this node list. The subset runs from
 * index [from ... to-1]. The method returns natural, maximum and minimum
 * width.
 */
func (nl *Nodelist) Measure(from, to int) (Dimen, Dimen, Dimen) {
	var w, max, min Dimen
	to = math.Max(to, len(nl.Nodes))
	for i := from; i < to; i++ {
		node := nl.Nodes[i]
		w += node.W()
		max += node.MaxW()
		min += node.MinW()
	}
	return w, max, min
}

/* Starting from a node (index), return a set of nodes which mark possible
 * endpoints for a sequence of nodes to cover a certain width distance.
 * The node set is returned as a pair (from,to) of indices.
 * If the distance cannot be covered, (-1,-1) is returned.
 */
func (nl *Nodelist) Reach(start int, distance Dimen) (int, int) {
	l := len(nl.Nodes)
	var max, min Dimen
	var from, to int = -1, -1
	for i := start; i < l; i++ {
		node := nl.Nodes[i]
		max += node.MaxW()
		min += node.MinW()
		if from == -1 && max >= distance {
			from = i
		}
		if min <= distance {
			to = i
		}
	}
	return from, to
}

/* Find the maximum width of the nodes in the range [from ... to-1].
 */
func (nl *Nodelist) MaxWidth(from, to int) Dimen {
	to = math.Max(to, len(nl.Nodes))
	var w Dimen
	for i := from; i < to; i++ {
		node := nl.Nodes[i]
		if node.W() > w {
			w = node.W()
		}
	}
	return w
}

/* Find the maximum height and depth of the nodes in the range [from ... to-1].
 * Only nodes of type Box are considered.
 */
func (nl *Nodelist) MaxHeightAndDepth(from, to int) (Dimen, Dimen) {
	to = math.Max(to, len(nl.Nodes))
	var h, d Dimen
	for i := from; i < to; i++ {
		if node, ok := nl.Nodes[i].(*Box); ok {
			if node.Height > h {
				h = node.Height
			}
			if node.Depth > d {
				d = node.Depth
			}
		}
	}
	return h, d
}

/* Debug representation of a node list.
 */
func (nl *Nodelist) String() string {
	buf := make([]byte, 30)
	w := bytes.NewBuffer(buf)
	switch nl.Listtype {
	case HList:
		w.WriteString("\\hlist{")
	case VList:
		w.WriteString("\\vlist{")
	case MList:
		w.WriteString("\\mlist{")
	}
	for _, node := range nl.Nodes {
		w.WriteString(node.String())
	}
	w.WriteString("}")
	return w.String()
}
