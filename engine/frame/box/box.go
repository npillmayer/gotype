/*
Package box deals with typesetting boxes.

Typesetting may be understood as the process of placing boxes within
larger boxes. The smallest type of box is a glyph, i.e. a printable
letter. The largest type of box is a page—or even a book, where
page-boxes are placed into.

The box model is very versatile. Nevertheless we will generalize the
notion of a box to mean the bounding box of a polygon. Typesetting in
irregular shapes is a feature available in most modern systems, e.g.
when letting text flow around a non-rectangular illustration.

This module deals with rectangular boxes, starting at the glyph level.
Boxes follow the CSS box model. Nevertheless, the notation oftentimes follows
the one introduced by the TeX typesetting system.

BSD License

Copyright (c) 2017–2020, Norbert Pillmayer

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */
package box

import (
	"image/color"

	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/core/font"
)

// Box type, following the CSS box model.
type Box struct {
	dimen.Rect
	Min     dimen.Point
	Max     dimen.Point
	Padding [4]dimen.Dimen // inside of border
	Margins [4]dimen.Dimen // outside of border
}

// For padding, margins
const (
	Top int = iota
	Right
	Bottom
	Left
)

// Box styling: We follow the CSS paradigm for boxes. Boxes are stylable
// objects which have dimensions, spacing, borders and colors.
//
// Some boxes may just implement a subset of the styling parameters. Most
// notably this holds for glyphs: Glyphs may have styled their content only.
// No border or additional spacing is possible with glyphs.

// ColorStyle is a type for styling with color.
type ColorStyle struct {
	Foreground color.Color
	Background color.Color // may be (semi-)transparent
}

// TextStyle is a type for styling text.
type TextStyle struct {
	Typecase *font.TypeCase
}

// BorderStyle is a type for simple borders.
type BorderStyle struct {
	LineColor     color.Color
	LineThickness dimen.Dimen
	LineStyle     int8
	CornerRadius  dimen.Dimen
}

// LineStyle is a type for border line styles.
type LineStyle int8

// We support these line styles only
const (
	LSSolid  LineStyle = 0
	LSDashed LineStyle = 1
	LSDotted LineStyle = 2
)

// Styling rolls all styling options into one type.
type Styling struct {
	TextStyle TextStyle
	Colors    ColorStyle
	Border    BorderStyle
}

// StyledBox is a type for a fully stylable box.
type StyledBox struct {
	Box
	Styling *Styling
}

// Glyph is a box for glyphs. Glyphs currently are content-stylable only (no borders).
//
// Wikipedia: In typography, a glyph [...] is an elemental symbol within
// an agreed set of symbols, intended to represent a readable character
// for the purposes of writing. [ Copyright (c) Wikipedia.com, 2017 ]
type Glyph struct {
	TextStyle *TextStyle
	Colors    *ColorStyle
	CharPos   rune
}

// Normalize sorts the corner coordinates into correct order.
func (box *Box) Normalize() *Box {
	if box.TopL.X > box.BotR.X {
		box.TopL.X, box.BotR.X = box.BotR.X, box.TopL.X
	}
	if box.TopL.Y > box.BotR.Y {
		box.TopL.Y, box.BotR.Y = box.BotR.Y, box.TopL.Y
	}
	return box
}

// Shift a box along a vector. The size of the box is unchanged.
func (box *Box) Shift(vector dimen.Point) *Box {
	box.TopL.Shift(vector)
	box.BotR.Shift(vector)
	return box
}

// Enlarge a box in x- and y-direction. For shrinking, use negative
// argument(s).
func (box *Box) Enlarge(scales dimen.Point) *Box {
	box.BotR.X = box.BotR.X + scales.X
	box.BotR.Y = box.BotR.Y + scales.Y
	return box
}

// Method for boxing content into a horizontal box. Content is given as a
// node list. The nodes will be enclosed into a new box.
// The box may be set to a target size.
// Parameters for styling class and/or identifier may be provided.
/*
func HBoxKhipu(nl *Khipu, target p.Dimen, identifier string, class string) *TypesetBox {
	box := &TypesetBox{}
	box.Cord = nl
	box.Style.StylingIdentifier = identifier
	box.Style.StylingClass = class
	box.Width = target
	_, max, min := nl.Measure(0, -1)
	if min > target {
		fmt.Println("overfull hbox")
	} else if max < target {
		fmt.Println("underfull hbox")
	}
	box.Height, box.Depth = nl.MaxHeightAndDepth(0, -1)
	return box
}
*/
