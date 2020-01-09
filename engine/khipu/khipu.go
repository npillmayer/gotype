package khipu

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/npillmayer/gotype/core/dimen"
)

/*
BSD License
Copyright (c) 2017-20, Norbert Pillmayer

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
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
*/

// Knots implement items for typesetting paragraphs. The various knot
// types more or less implement the corresponding node types from the TeX
// typesetting system.

// === Knots =================================================================

// KnotType is a type for the different flavours of knots.
type KnotType int8

// A Knot has a width and may be discardable
type Knot interface {
	Type() KnotType      // type identifier of this knot
	W() dimen.Dimen      // width
	MinW() dimen.Dimen   // minimum width
	MaxW() dimen.Dimen   // maximum width
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

// NewKnot is a factory method to create a knot. Parameter is a valid knot type.
func NewKnot(knottype KnotType) Knot {
	switch knottype {
	case KTKern:
		return Kern(0)
	case KTGlue:
		return Glue{0, 0, 0}
	case KTPenalty:
		return Penalty(0)
	case KTDiscretionary:
		return Discretionary('-') // TODO should be hyphenchar of current font or -1
	case KTTextBox:
		box := &TextBox{}
		return box
	}
	return nil
}

// KnotString is a debugging helper and returns a textual representation of a knot.
func KnotString(k Knot) string {
	switch k.Type() {
	case KTKern:
		return k.(Kern).String()
	case KTGlue:
		//return fmt.Sprintf("[glue %v]", k)
		//g := k.(Glue)
		return k.(Glue).String()
	case KTPenalty:
		return k.(Penalty).String()
	case KTTextBox:
		s, ok := k.(*TextBox)
		if ok {
			return s.String()
		}
		return k.(TextBox).String()
	case KTDiscretionary:
		return "\u2af6"
	default:
		return "yes, it is a knot"
	}
}

// --- Kern ------------------------------------------------------------------

// A Kern is an unshrinkable space
type Kern dimen.Dimen // fixed width

// Type is part of interface Knot.
func (k Kern) Type() KnotType {
	return KTKern
}

func (k Kern) String() string {
	return fmt.Sprintf("[\u29e6%s]", k.W())
}

// W is part of interface Knot. Width of the kern.
func (k Kern) W() dimen.Dimen {
	return dimen.Dimen(k)
}

// MinW is part of interface Knot. Kerns do not shrink.
func (k Kern) MinW() dimen.Dimen {
	return dimen.Dimen(k)
}

// MaxW is part of interface Knot. Kerns do not stretch.
func (k Kern) MaxW() dimen.Dimen {
	return dimen.Dimen(k)
}

// IsDiscardable is part of interface Knot. Kerns are discardable.
func (k Kern) IsDiscardable() bool {
	return true
}

// --- Glue ------------------------------------------------------------------

// A Glue is a space which can shrink and expand
type Glue [3]dimen.Dimen

// Type is part of interface Knot.
func (g Glue) Type() KnotType {
	return KTGlue
}

func (g Glue) String() string {
	minus := g.W() - g.MinW()
	plus := g.MaxW() - g.W()
	return fmt.Sprintf("(\u29df %s-%s+%s)", g.W().String(), minus, plus)
}

// W is part of interface Knot. Natural width of the glue.
func (g Glue) W() dimen.Dimen {
	return g[0]
}

// MinW is part of interface Knot. Minimum width of the glue.
func (g Glue) MinW() dimen.Dimen {
	return g[0] + g[1]
}

// MaxW is part of interface Knot. Maximum width of the glue.
func (g Glue) MaxW() dimen.Dimen {
	return g[0] + g[2]
}

// IsDiscardable is part of interface Knot. Glue is discardable.
func (g Glue) IsDiscardable() bool {
	return true
}

// NewGlue creates a new drop of glue with stretch and shrink.
func NewGlue(w dimen.Dimen, shrink dimen.Dimen, stretch dimen.Dimen) Glue {
	g := Glue{w, shrink, stretch}
	return g
}

// NewFill creates a drop of infinitely stretchable glue.
func NewFill(f int) Glue {
	var stretch dimen.Dimen
	switch f {
	case 2:
		stretch = dimen.Fill
	case 3:
		stretch = dimen.Filll
	default:
		stretch = dimen.Fil
	}
	return NewGlue(0, 0, stretch)
}

// --- Discretionary ---------------------------------------------------------

// A Discretionary is a hyphenation opportunity
type Discretionary rune

// Type is part of interface Knot.
func (d Discretionary) Type() KnotType {
	return KTDiscretionary
}

// W is part of interface Knot. Returns the width of the un-hyphenated text.
func (d Discretionary) W() dimen.Dimen {
	return 0
}

// MinW is part of interface Knot. Returns the width of the pre-hyphen text.
func (d Discretionary) MinW() dimen.Dimen {
	return 0
}

// MaxW is part of interface Knot. Returns the width of the post-hyphen text.
func (d Discretionary) MaxW() dimen.Dimen {
	return 5 * dimen.PT // TODO
}

// IsDiscardable is part of interface Knot. Discretionaries are not discardable.
func (d Discretionary) IsDiscardable() bool {
	return false
}

// --- Boxes -----------------------------------------------------------------

// A TextBox is a fixed unit of text
type TextBox struct {
	Width  dimen.Dimen // width
	Height dimen.Dimen // height
	Depth  dimen.Dimen // depth
	text   string      // text, if available
	//knotlist Khipu // content, if available
}

// NewTextBox creates a text box.
func NewTextBox(s string) *TextBox {
	box := &TextBox{}
	box.text = s
	return box
}

// Type is part of interface Knot.
func (b TextBox) Type() KnotType {
	return KTTextBox
}

func (b TextBox) String() string {
	//return fmt.Sprintf("\u25a1\u00ab%s\u00bb", b.text)
	return fmt.Sprintf("\u00ab%s\u00bb", b.text)
}

// W is part of interface Knot. Width of the glue.
func (b TextBox) W() dimen.Dimen {
	return b.Width
}

// MinW is part of interface Knot. Width of the glue.
func (b TextBox) MinW() dimen.Dimen {
	return b.Width
}

// MaxW is part of interface Knot. Width of the glue.
func (b TextBox) MaxW() dimen.Dimen {
	return b.Width
}

// IsDiscardable is part of interface Knot. Text is not discardable.
func (b TextBox) IsDiscardable() bool {
	return false
}

var _ Knot = &TextBox{}

// --- Penalty ---------------------------------------------------------------

// A Penalty contributes to demerits, i.e. the quality index of paragraphs
type Penalty int

// Type is part of interface Knot.
func (p Penalty) Type() KnotType {
	return KTPenalty
}

func (p Penalty) String() string {
	return fmt.Sprintf("(\u2af2 %d)", p)
}

// W is part of interface Knot. Returns 0.
func (p Penalty) W() dimen.Dimen {
	return 0
}

// MinW is part of interface Knot. Returns 0.
func (p Penalty) MinW() dimen.Dimen {
	return 0
}

// MaxW is part of interface Knot. Returns 0.
func (p Penalty) MaxW() dimen.Dimen {
	return 0
}

// IsDiscardable is part of interface Knot. Penalties are discardable.
func (p Penalty) IsDiscardable() bool {
	return true
}

// === Khipus ================================================================

// Khipu is a string of knots.
// We handle text/paragraphs as khipus.
type Khipu struct {
	typ   int    // hlist, vlist or mlist
	knots []Knot // array of knots of different type
}

// List types
const (
	HList int = iota // horizontal list
	VList            // vertical list
	MList            // math list
)

// NewKhipu creates a new knot list.
func NewKhipu() *Khipu {
	kh := &Khipu{}
	kh.knots = make([]Knot, 0, 50)
	return kh
}

// Length gives the number of knots in the list.
func (kh *Khipu) Length() int {
	return len(kh.knots)
}

// AppendKnot appends a knot at the end of the list.
func (kh *Khipu) AppendKnot(knot Knot) *Khipu {
	kh.knots = append(kh.knots, knot)
	return kh
}

// AppendKhipu concatenates two khipus.
func (kh *Khipu) AppendKhipu(k *Khipu) *Khipu {
	kh.knots = append(kh.knots, k.knots...)
	// for _, knot := range k.knots {
	// 	kh.knots = append(kh.knots, knot)
	// }
	return kh
}

// Measure returns the widths of a subset of this knot list. The subset runs from
// index [from ... to-1]. The method returns natural, maximum and minimum
// width.
func (kh *Khipu) Measure(from, to int) (dimen.Dimen, dimen.Dimen, dimen.Dimen) {
	var w, max, min dimen.Dimen
	to = iMax(to, len(kh.knots))
	for i := from; i < to; i++ {
		knot := kh.knots[i]
		w += knot.W()
		max += knot.MaxW()
		min += knot.MinW()
	}
	return w, max, min
}

// Reach iterates over a khipu to find a point beyond a given distance.
// Starting from a knot (index), return a set of knots which mark possible
// endpoints for a sequence of knots to cover a certain width distance.
// The knot set is returned as a pair (from,to) of indices.
// If the distance cannot be covered, (-1,-1) is returned.
func (kh *Khipu) Reach(start int, distance dimen.Dimen) (int, int) {
	l := len(kh.knots)
	var max, min dimen.Dimen
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

// MaxWidth finds the maximum width of the knots in the range [from ... to-1].
func (kh *Khipu) MaxWidth(from, to int) dimen.Dimen {
	to = iMax(to, len(kh.knots))
	var w dimen.Dimen
	for i := from; i < to; i++ {
		knot := kh.knots[i]
		if knot.W() > w {
			w = knot.W()
		}
	}
	return w
}

// MaxHeightAndDepth finds the maximum height and depth of the knots in the range
// [from ... to-1].
// Only knots of type TextBox are considered.
func (kh *Khipu) MaxHeightAndDepth(from, to int) (dimen.Dimen, dimen.Dimen) {
	to = iMax(to, len(kh.knots))
	var h, d dimen.Dimen
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

// Debug representation of a knot list.
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
	errorIteratatorEnd = errors.New("Khipu-iterator at end of knot list")
)

type KhipuIterator struct {
	khipu *Khipu
	inx   int
}

// Iterator returns an interator over a khipu.
func (kh *Khipu) Iterator() *KhipuIterator {
	return &KhipuIterator{kh, -1}
}

func (khit *KhipuIterator) String() string {
	return fmt.Sprintf("[%d]%v", khit.inx, khit.Knot())
}

func (khit *KhipuIterator) Next() bool {
	khit.inx++
	return khit.inx < len(khit.khipu.knots)
}

func (khit *KhipuIterator) Index() int {
	return khit.inx
}

func (khit *KhipuIterator) Knot() Knot {
	k := khit.khipu.knots[khit.inx]
	return k
}

func (khit *KhipuIterator) AsGlue() Glue {
	return khit.Knot().(Glue)
}

func (khit *KhipuIterator) AsPenalty() Penalty {
	return khit.Knot().(Penalty)
}

func (khit *KhipuIterator) AsKern() Kern {
	return khit.Knot().(Kern)
}

func (khit *KhipuIterator) AsTextBox() *TextBox {
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
