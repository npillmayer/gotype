package khipus

import (
	"bytes"
	"fmt"

	p "github.com/npillmayer/gotype/gtcore/parameters"
)

/*
---------------------------------------------------------------------------

BSD License
Copyright (c) 2017-18, Norbert Pillmayer

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

/* Knots implement items for typesetting paragraphs. The various knot
 * types more or less implement the corresponding node types from the TeX
 * typesetting system.

*/

// === Knots =================================================================

// A knot has a width and may be discardable
type Knot interface {
	fmt.Stringer
	W() p.Dimen          // width
	MinW() p.Dimen       // minimum width
	MaxW() p.Dimen       // maximum width
	IsDiscardable() bool // is this knot discardable?
}

// Knot types
const (
	KTKern int = iota
	KTGlue
	KTBox
	KTPenalty
	KTDiscretionary
)

/* Factory method to create a knot. Parameter is a valid knot type.
 */
func NewKnot(knottype int) Knot {
	switch knottype {
	case KTKern:
		return &Kern{0}
	case KTGlue:
		return &Glue{}
	case KTPenalty:
		return &Penalty{0}
	case KTDiscretionary:
		d := &Discretionary{}
		//d.pre = "-"
		return d
	case KTBox:
		box := &Box{}
		return box
	}
	return nil
}

// --- Kern ------------------------------------------------------------------

// A kern is an unshrinkable space
type Kern struct {
	Width p.Dimen // fixed width
}

/* Interface Knot. Prints the dimension (width) of the kern.
 */
func (k *Kern) String() string {
	return fmt.Sprintf("[kern %s]", k.Width.String())
}

/* Interface Knot. Width of the kern.
 */
func (k *Kern) W() p.Dimen {
	return k.Width
}

/* Interface Knot. Kerns do not shrink.
 */
func (k *Kern) MinW() p.Dimen {
	return k.Width
}

/* Interface Knot. Kerns do not stretch.
 */
func (k *Kern) MaxW() p.Dimen {
	return k.Width
}

/* Interface Knot. Kerns are discardable.
 */
func (k *Kern) IsDiscardable() bool {
	return true
}

var _ Knot = &Kern{}

// --- Glue ------------------------------------------------------------------

// A glue is a space which can shrink and expand
type Glue struct {
	Width    p.Dimen // natural width
	MaxWidth p.Dimen // maximum width
	MinWidth p.Dimen // minimum width
}

/* Interface Knot.
 */
func (g *Glue) String() string {
	return fmt.Sprintf("[glue %s <%s >%s]", g.W().String(), g.MaxW().String(),
		g.MinW().String())
}

/* Interface Knot. Natural width of the glue.
 */
func (g *Glue) W() p.Dimen {
	return g.Width
}

/* Interface Knot. Minimum width of the glue.
 */
func (g *Glue) MinW() p.Dimen {
	return g.MinWidth
}

/* Interface Knot. Maximum width of the glue.
 */
func (g *Glue) MaxW() p.Dimen {
	return g.MaxWidth
}

/* Interface Knot. Glue is discardable.
 */
func (g *Glue) IsDiscardable() bool {
	return true
}

var _ Knot = &Glue{}

/* Create a new drop of glue with stretch and shrink.
 */
func NewGlue(w p.Dimen, stretch p.Dimen, shrink p.Dimen) *Glue {
	glue := NewKnot(KTGlue).(*Glue)
	glue.Width = w
	glue.MaxWidth = w + stretch
	glue.MinWidth = w - shrink
	return glue
}

/* Create a drop of infinitely stretchable glue.
 */
func NewFill(f int) *Glue {
	var stretch p.Dimen
	switch f {
	case 2:
		stretch = p.Fill
	case 3:
		stretch = p.Filll
	default:
		stretch = p.Fil
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

/* Interface Knot. Prints the dimension (width) of the kern.
 */
func (d *Discretionary) String() string {
	return fmt.Sprintf("\\discretionary{%s}{%s}{%s}", d.nobreak.text,
		d.pre.text, d.post.text)
}

/* Interface Knot. Returns the width of the un-hyphenated text.
 */
func (d *Discretionary) W() p.Dimen {
	return d.nobreak.W()
}

/* Interface Knot. Returns the width of the pre-hyphen text.
 */
func (d *Discretionary) MinW() p.Dimen {
	return d.pre.W()
}

/* Interface Knot. Returns the width of the post-hyphen text.
 */
func (d *Discretionary) MaxW() p.Dimen {
	return d.post.W()
}

/* Interface Knot. Discretionaries are not discardable.
 */
func (d *Discretionary) IsDiscardable() bool {
	return false
}

var _ Knot = &Discretionary{}

// --- Boxes -----------------------------------------------------------------

// A Box is a fixed unit of text
type Box struct {
	Width  p.Dimen // width
	Height p.Dimen // height
	Depth  p.Dimen // depth
	text   string  // text, if available
	//knotlist Khipu // content, if available
}

func NewWordBox(s string) *Box {
	box := &Box{}
	box.text = s
	return box
}

/* Interface Knot.
 */
func (b *Box) String() string {
	return fmt.Sprintf("\\box{%s}", b.text)
}

/* Interface Knot. Width of the glue.
 */
func (b *Box) W() p.Dimen {
	return b.Width
}

/* Interface Knot. Width of the glue.
 */
func (b *Box) MinW() p.Dimen {
	return b.Width
}

/* Interface Knot. Width of the glue.
 */
func (b *Box) MaxW() p.Dimen {
	return b.Width
}

/* Interface Knot. Glue is discardable.
 */
func (b *Box) IsDiscardable() bool {
	return false
}

var _ Knot = &Box{}

// --- Penalty ---------------------------------------------------------------

// A penalty contributes to demerits, i.e. the quality index of paragraphs
type Penalty struct {
	P int
}

/* Interface Knot.
 */
func (p *Penalty) String() string {
	return fmt.Sprintf("[penalty %d]", p.P)
}

/* Interface Knot. Returns 0.
 */
func (p *Penalty) W() p.Dimen {
	return 0
}

/* Interface Knot. Returns 0.
 */
func (p *Penalty) MinW() p.Dimen {
	return 0
}

/* Interface Knot. Returns 0.
 */
func (p *Penalty) MaxW() p.Dimen {
	return 0
}

/* Interface Knot. Penalties are discardable.
 */
func (p *Penalty) IsDiscardable() bool {
	return true
}

var _ Knot = &Penalty{}

// === Knot lists ============================================================

// We handle text/paragraphcs as khips, i.e. string of knots
type Khipu struct {
	Typ   int    // hlist, vlist or mlist
	Knots []Knot // array of knots of different type
}

// List types
const (
	HList int = iota
	VList
	MList
)

/* Create a new knot list.
 */
func NewKhipu() *Khipu {
	nl := &Khipu{}
	nl.Knots = make([]Knot, 0, 50)
	return nl
}

/* Number of knots in the list.
 */
func (nl *Khipu) Length() int {
	return len(nl.Knots)
}

/* Append a knot at the end of the list.
 */
func (nl *Khipu) AppendKnot(knot Knot) *Khipu {
	nl.Knots = append(nl.Knots, knot)
	return nl
}

/* Return the widths of a subset of this knot list. The subset runs from
 * index [from ... to-1]. The method returns natural, maximum and minimum
 * width.
 */
func (nl *Khipu) Measure(from, to int) (p.Dimen, p.Dimen, p.Dimen) {
	var w, max, min p.Dimen
	to = iMax(to, len(nl.Knots))
	for i := from; i < to; i++ {
		knot := nl.Knots[i]
		w += knot.W()
		max += knot.MaxW()
		min += knot.MinW()
	}
	return w, max, min
}

/* Starting from a knot (index), return a set of knots which mark possible
 * endpoints for a sequence of knots to cover a certain width distance.
 * The knot set is returned as a pair (from,to) of indices.
 * If the distance cannot be covered, (-1,-1) is returned.
 */
func (nl *Khipu) Reach(start int, distance p.Dimen) (int, int) {
	l := len(nl.Knots)
	var max, min p.Dimen
	var from, to int = -1, -1
	for i := start; i < l; i++ {
		knot := nl.Knots[i]
		max += knot.MaxW()
		min += knot.MinW()
		if from == -1 && max >= distance {
			from = i
		}
		if min <= distance {
			to = i
		}
	}
	return from, to
}

/* Find the maximum width of the knots in the range [from ... to-1].
 */
func (nl *Khipu) MaxWidth(from, to int) p.Dimen {
	to = iMax(to, len(nl.Knots))
	var w p.Dimen
	for i := from; i < to; i++ {
		knot := nl.Knots[i]
		if knot.W() > w {
			w = knot.W()
		}
	}
	return w
}

/* Find the maximum height and depth of the knots in the range [from ... to-1].
 * Only knots of type Box are considered.
 */
func (nl *Khipu) MaxHeightAndDepth(from, to int) (p.Dimen, p.Dimen) {
	to = iMax(to, len(nl.Knots))
	var h, d p.Dimen
	for i := from; i < to; i++ {
		if knot, ok := nl.Knots[i].(*Box); ok {
			if knot.Height > h {
				h = knot.Height
			}
			if knot.Depth > d {
				d = knot.Depth
			}
		}
	}
	return h, d
}

/* Debug representation of a knot list.
 */
func (nl *Khipu) String() string {
	buf := make([]byte, 30)
	w := bytes.NewBuffer(buf)
	switch nl.Typ {
	case HList:
		w.WriteString("\\hlist{")
	case VList:
		w.WriteString("\\vlist{")
	case MList:
		w.WriteString("\\mlist{")
	}
	for _, knot := range nl.Knots {
		w.WriteString(knot.String())
	}
	w.WriteString("}")
	return w.String()
}

// ----------------------------------------------------------------------

func iMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func iMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}
