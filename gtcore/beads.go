package gtcore

import (
	"bytes"
	"fmt"

	p "github.com/npillmayer/gotype/gtcore/parameters"
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

/* Beads implement items for typesetting paragraphs. The various bead
 * types more or less implement the corresponding bead types from the TeX
 * typesetting system.

*/

// === Beads =================================================================

// A bead has a width and may be discardable
type Bead interface {
	fmt.Stringer
	W() p.Dimen          // width
	MinW() p.Dimen       // minimum width
	MaxW() p.Dimen       // maximum width
	IsDiscardable() bool // is this bead discardable?
}

// Bead types
const (
	NTKern int = iota
	NTGlue
	NTBox
	NTPenalty
	NTDiscretionary
)

/* Factory method to create a bead. Parameter is a valid bead type.
 */
func NewBead(beadtype int) Bead {
	switch beadtype {
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
	Width p.Dimen // fixed width
}

/* Interface Bead. Prints the dimension (width) of the kern.
 */
func (k *Kern) String() string {
	return fmt.Sprintf("[kern %s]", k.Width.String())
}

/* Interface Bead. Width of the kern.
 */
func (k *Kern) W() p.Dimen {
	return k.Width
}

/* Interface Bead. Kerns do not shrink.
 */
func (k *Kern) MinW() p.Dimen {
	return k.Width
}

/* Interface Bead. Kerns do not stretch.
 */
func (k *Kern) MaxW() p.Dimen {
	return k.Width
}

/* Interface Bead. Kerns are discardable.
 */
func (k *Kern) IsDiscardable() bool {
	return true
}

var _ Bead = &Kern{}

// --- Glue ------------------------------------------------------------------

// A glue is a space which can shrink and expand
type Glue struct {
	Width    p.Dimen // natural width
	MaxWidth p.Dimen // maximum width
	MinWidth p.Dimen // minimum width
}

/* Interface Bead.
 */
func (g *Glue) String() string {
	return fmt.Sprintf("[glue %s <%s >%s]", g.W().String(), g.MaxW().String(),
		g.MinW().String())
}

/* Interface Bead. Natural width of the glue.
 */
func (g *Glue) W() p.Dimen {
	return g.Width
}

/* Interface Bead. Minimum width of the glue.
 */
func (g *Glue) MinW() p.Dimen {
	return g.MinWidth
}

/* Interface Bead. Maximum width of the glue.
 */
func (g *Glue) MaxW() p.Dimen {
	return g.MaxWidth
}

/* Interface Bead. Glue is discardable.
 */
func (g *Glue) IsDiscardable() bool {
	return true
}

var _ Bead = &Glue{}

/* Create a new drop of glue with stretch and shrink.
 */
func NewGlue(w p.Dimen, stretch p.Dimen, shrink p.Dimen) *Glue {
	glue := NewBead(NTGlue).(*Glue)
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

/* Interface Bead. Prints the dimension (width) of the kern.
 */
func (d *Discretionary) String() string {
	return fmt.Sprintf("\\discretionary{%s}{%s}{%s}", d.nobreak.text,
		d.pre.text, d.post.text)
}

/* Interface Bead. Returns the width of the un-hyphenated text.
 */
func (d *Discretionary) W() p.Dimen {
	return d.nobreak.W()
}

/* Interface Bead. Returns the width of the pre-hyphen text.
 */
func (d *Discretionary) MinW() p.Dimen {
	return d.pre.W()
}

/* Interface Bead. Returns the width of the post-hyphen text.
 */
func (d *Discretionary) MaxW() p.Dimen {
	return d.post.W()
}

/* Interface Bead. Discretionaries are not discardable.
 */
func (d *Discretionary) IsDiscardable() bool {
	return false
}

var _ Bead = &Discretionary{}

// --- Boxes -----------------------------------------------------------------

// A Box is a fixed unit of text
type Box struct {
	Width  p.Dimen // width
	Height p.Dimen // height
	Depth  p.Dimen // depth
	text   string  // text, if available
	//beadlist BeadChain // content, if available
}

/* Interface Bead.
 */
func (b *Box) String() string {
	return fmt.Sprintf("\\box{%s}", b.text)
}

/* Interface Bead. Width of the glue.
 */
func (b *Box) W() p.Dimen {
	return b.Width
}

/* Interface Bead. Width of the glue.
 */
func (b *Box) MinW() p.Dimen {
	return b.Width
}

/* Interface Bead. Width of the glue.
 */
func (b *Box) MaxW() p.Dimen {
	return b.Width
}

/* Interface Bead. Glue is discardable.
 */
func (b *Box) IsDiscardable() bool {
	return false
}

var _ Bead = &Box{}

// --- Penalty ---------------------------------------------------------------

// A penalty contributes to demerits, i.e. the quality index of paragraphs
type Penalty struct {
	P int
}

/* Interface Bead.
 */
func (p *Penalty) String() string {
	return fmt.Sprintf("[penalty %d]", p.P)
}

/* Interface Bead. Returns 0.
 */
func (p *Penalty) W() p.Dimen {
	return 0
}

/* Interface Bead. Returns 0.
 */
func (p *Penalty) MinW() p.Dimen {
	return 0
}

/* Interface Bead. Returns 0.
 */
func (p *Penalty) MaxW() p.Dimen {
	return 0
}

/* Interface Bead. Penalties are discardable.
 */
func (p *Penalty) IsDiscardable() bool {
	return true
}

var _ Bead = &Penalty{}

// === Bead lists ============================================================

// We handle text/paragraphcs as a list of beads
type BeadChain struct {
	Listtype int    // hlist, vlist or mlist
	Beads    []Bead // array of beads of different type
}

// List types
const (
	HList int = iota
	VList
	MList
)

/* Create a new bead list.
 */
func NewBeadChain() *BeadChain {
	nl := &BeadChain{}
	nl.Beads = make([]Bead, 0, 50)
	return nl
}

/* Number of beads in the list.
 */
func (nl *BeadChain) Length() int {
	return len(nl.Beads)
}

/* Append a bead at the end of the list.
 */
func (nl *BeadChain) AppendBead(bead Bead) *BeadChain {
	nl.Beads = append(nl.Beads, bead)
	return nl
}

/* Return the widths of a subset of this bead list. The subset runs from
 * index [from ... to-1]. The method returns natural, maximum and minimum
 * width.
 */
func (nl *BeadChain) Measure(from, to int) (p.Dimen, p.Dimen, p.Dimen) {
	var w, max, min p.Dimen
	to = iMax(to, len(nl.Beads))
	for i := from; i < to; i++ {
		bead := nl.Beads[i]
		w += bead.W()
		max += bead.MaxW()
		min += bead.MinW()
	}
	return w, max, min
}

/* Starting from a bead (index), return a set of beads which mark possible
 * endpoints for a sequence of beads to cover a certain width distance.
 * The bead set is returned as a pair (from,to) of indices.
 * If the distance cannot be covered, (-1,-1) is returned.
 */
func (nl *BeadChain) Reach(start int, distance p.Dimen) (int, int) {
	l := len(nl.Beads)
	var max, min p.Dimen
	var from, to int = -1, -1
	for i := start; i < l; i++ {
		bead := nl.Beads[i]
		max += bead.MaxW()
		min += bead.MinW()
		if from == -1 && max >= distance {
			from = i
		}
		if min <= distance {
			to = i
		}
	}
	return from, to
}

/* Find the maximum width of the beads in the range [from ... to-1].
 */
func (nl *BeadChain) MaxWidth(from, to int) p.Dimen {
	to = iMax(to, len(nl.Beads))
	var w p.Dimen
	for i := from; i < to; i++ {
		bead := nl.Beads[i]
		if bead.W() > w {
			w = bead.W()
		}
	}
	return w
}

/* Find the maximum height and depth of the beads in the range [from ... to-1].
 * Only beads of type Box are considered.
 */
func (nl *BeadChain) MaxHeightAndDepth(from, to int) (p.Dimen, p.Dimen) {
	to = iMax(to, len(nl.Beads))
	var h, d p.Dimen
	for i := from; i < to; i++ {
		if bead, ok := nl.Beads[i].(*Box); ok {
			if bead.Height > h {
				h = bead.Height
			}
			if bead.Depth > d {
				d = bead.Depth
			}
		}
	}
	return h, d
}

/* Debug representation of a bead list.
 */
func (nl *BeadChain) String() string {
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
	for _, bead := range nl.Beads {
		w.WriteString(bead.String())
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
