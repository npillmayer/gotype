package khipus

import (
	"bytes"
	"errors"
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

// Every knot has a type
type KnotType int8

// A knot has a width and may be discardable
type Knot interface {
	Type() KnotType      // type identifier of this knot
	W() p.Dimen          // width
	MinW() p.Dimen       // minimum width
	MaxW() p.Dimen       // maximum width
	IsDiscardable() bool // is this knot discardable?
}

// Knot types
const (
	KTKern KnotType = iota
	KTGlue
	KTTextBox
	KTPenalty
	KTDiscretionary
	KTUserDefined // clients should use custom knot types above this
)

// Factory method to create a knot. Parameter is a valid knot type.
func NewKnot(knottype KnotType) Knot {
	switch knottype {
	case KTKern:
		return Kern(0)
	case KTGlue:
		return Glue{0, 0, 0}
	case KTPenalty:
		return Penalty(0)
	case KTDiscretionary:
		d := &Discretionary{}
		//d.pre = "-"
		return d
	case KTTextBox:
		box := &TextBox{}
		return box
	}
	return nil
}

func KnotString(k Knot) string {
	switch k.Type() {
	case KTKern:
		return fmt.Sprintf("[kern %s]", k.W())
	case KTGlue:
		//return fmt.Sprintf("[glue %v]", k)
		g := k.(Glue)
		return g.String()
	case KTPenalty:
		p := k.(Penalty)
		return fmt.Sprintf("[penalty %d]", p)
	case KTTextBox:
		s := k.(*TextBox)
		return s.String()
	default:
		return "yes, it is a knot"
	}
	return "TODO"
}

// --- Kern ------------------------------------------------------------------

// A kern is an unshrinkable space
type Kern p.Dimen // fixed width

// Interface Knot.
func (k Kern) Type() KnotType {
	return KTKern
}

// Interface Knot. Width of the kern.
func (k Kern) W() p.Dimen {
	return p.Dimen(k)
}

// Interface Knot. Kerns do not shrink.
func (k Kern) MinW() p.Dimen {
	return p.Dimen(k)
}

// Interface Knot. Kerns do not stretch.
func (k Kern) MaxW() p.Dimen {
	return p.Dimen(k)
}

// Interface Knot. Kerns are discardable.
func (k Kern) IsDiscardable() bool {
	return true
}

// --- Glue ------------------------------------------------------------------

// A glue is a space which can shrink and expand
type Glue [3]p.Dimen

// Interface Knot.
func (g Glue) Type() KnotType {
	return KTGlue
}

func (g Glue) String() string {
	return fmt.Sprintf("[glue %s <%s >%s]", g.W().String(), g.MinW().String(),
		g.MaxW().String())
}

// Interface Knot. Natural width of the glue.
func (g Glue) W() p.Dimen {
	return g[0]
}

// Interface Knot. Minimum width of the glue.
func (g Glue) MinW() p.Dimen {
	return g[0] + g[1]
}

// Interface Knot. Maximum width of the glue.
func (g Glue) MaxW() p.Dimen {
	return g[0] + g[2]
}

// Interface Knot. Glue is discardable.
func (g Glue) IsDiscardable() bool {
	return true
}

// Create a new drop of glue with stretch and shrink.
func NewGlue(w p.Dimen, shrink p.Dimen, stretch p.Dimen) Glue {
	g := Glue{w, shrink, stretch}
	return g
}

// Create a drop of infinitely stretchable glue.
func NewFill(f int) Glue {
	var stretch p.Dimen
	switch f {
	case 2:
		stretch = p.Fill
	case 3:
		stretch = p.Filll
	default:
		stretch = p.Fil
	}
	return NewGlue(0, 0, stretch)
}

// --- Discretionary ---------------------------------------------------------

// A discretionary is a hyphenation opportunity
type Discretionary struct {
	nobreak TextBox // text if not hyphenated
	pre     TextBox // pre-hyphen text
	post    TextBox // post-hyphen text
}

// Interface Knot.
func (d *Discretionary) Type() KnotType {
	return KTDiscretionary
}

/* Prints the dimension (width) of the kern.
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

// A TextBox is a fixed unit of text
type TextBox struct {
	Width  p.Dimen // width
	Height p.Dimen // height
	Depth  p.Dimen // depth
	text   string  // text, if available
	//knotlist Khipu // content, if available
}

func NewWordBox(s string) *TextBox {
	box := &TextBox{}
	box.text = s
	return box
}

// Interface Knot.
func (b *TextBox) Type() KnotType {
	return KTTextBox
}

/* Interface Knot.
 */
func (b *TextBox) String() string {
	return fmt.Sprintf("\\box{%s}", b.text)
}

/* Interface Knot. Width of the glue.
 */
func (b *TextBox) W() p.Dimen {
	return b.Width
}

/* Interface Knot. Width of the glue.
 */
func (b *TextBox) MinW() p.Dimen {
	return b.Width
}

/* Interface Knot. Width of the glue.
 */
func (b *TextBox) MaxW() p.Dimen {
	return b.Width
}

/* Interface Knot. Glue is discardable.
 */
func (b *TextBox) IsDiscardable() bool {
	return false
}

var _ Knot = &TextBox{}

// --- Penalty ---------------------------------------------------------------

// A penalty contributes to demerits, i.e. the quality index of paragraphs
type Penalty int

// Interface Knot.
func (p Penalty) Type() KnotType {
	return KTPenalty
}

// Interface Knot. Returns 0.
func (p Penalty) W() p.Dimen {
	return 0
}

// Interface Knot. Returns 0.
func (p Penalty) MinW() p.Dimen {
	return 0
}

// Interface Knot. Returns 0.
func (p Penalty) MaxW() p.Dimen {
	return 0
}

// Interface Knot. Penalties are discardable.
func (p Penalty) IsDiscardable() bool {
	return true
}

// === Khipus ================================================================

// We handle text/paragraphcs as khipus, i.e. string of knots
type Khipu struct {
	typ   int    // hlist, vlist or mlist
	knots []Knot // array of knots of different type
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
	kh := &Khipu{}
	kh.knots = make([]Knot, 0, 50)
	return kh
}

/* Number of knots in the list.
 */
func (kh *Khipu) Length() int {
	return len(kh.knots)
}

/* Append a knot at the end of the list.
 */
func (kh *Khipu) AppendKnot(knot Knot) *Khipu {
	kh.knots = append(kh.knots, knot)
	return kh
}

func (kh *Khipu) AppendKhipu(k *Khipu) *Khipu {
	for _, knot := range k.knots {
		kh.knots = append(kh.knots, knot)
	}
	return kh
}

/* Return the widths of a subset of this knot list. The subset runs from
 * index [from ... to-1]. The method returns natural, maximum and minimum
 * width.
 */
func (kh *Khipu) Measure(from, to int) (p.Dimen, p.Dimen, p.Dimen) {
	var w, max, min p.Dimen
	to = iMax(to, len(kh.knots))
	for i := from; i < to; i++ {
		knot := kh.knots[i]
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
func (kh *Khipu) Reach(start int, distance p.Dimen) (int, int) {
	l := len(kh.knots)
	var max, min p.Dimen
	var from, to int = -1, -1
	for i := start; i < l; i++ {
		knot := kh.knots[i]
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
func (kh *Khipu) MaxWidth(from, to int) p.Dimen {
	to = iMax(to, len(kh.knots))
	var w p.Dimen
	for i := from; i < to; i++ {
		knot := kh.knots[i]
		if knot.W() > w {
			w = knot.W()
		}
	}
	return w
}

/* Find the maximum height and depth of the knots in the range [from ... to-1].
 * Only knots of type TextBox are considered.
 */
func (kh *Khipu) MaxHeightAndDepth(from, to int) (p.Dimen, p.Dimen) {
	to = iMax(to, len(kh.knots))
	var h, d p.Dimen
	for i := from; i < to; i++ {
		if knot, ok := kh.knots[i].(*TextBox); ok {
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
func (kh *Khipu) String() string {
	buf := make([]byte, 30)
	w := bytes.NewBuffer(buf)
	switch kh.typ {
	case HList:
		w.WriteString("\\hlist{")
	case VList:
		w.WriteString("\\vlist{")
	case MList:
		w.WriteString("\\mlist{")
	}
	for _, knot := range kh.knots {
		w.WriteString(KnotString(knot))
	}
	w.WriteString("}")
	return w.String()
}

// ----------------------------------------------------------------------

var (
	errorIteratatorEnd error = errors.New("Khipu-iterator at end of knot list")
)

type khipuIterator struct {
	khipu *Khipu
	inx   int
}

func (kh *Khipu) Iterator() *khipuIterator {
	return &khipuIterator{kh, -1}
}

func (khit *khipuIterator) Next() bool {
	khit.inx++
	return khit.inx < len(khit.khipu.knots)
}

func (khit *khipuIterator) Knot() Knot {
	k := khit.khipu.knots[khit.inx]
	return k
}

func (khit *khipuIterator) AsGlue() Glue {
	return khit.Knot().(Glue)
}

func (khit *khipuIterator) AsPenalty() Penalty {
	return khit.Knot().(Penalty)
}

func (khit *khipuIterator) AsKern() Kern {
	return khit.Knot().(Kern)
}

func (khit *khipuIterator) AsTextBox() *TextBox {
	return khit.Knot().(*TextBox)
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
