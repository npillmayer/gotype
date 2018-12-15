// Package box deals with typesetting boxes.
//
/*
Typesetting may be understood as the process of placing boxes within
larger boxes. The smallest type of box is a glyph, i.e. a printable
letter. The largest type of box is a page -- or even a book, where
page-boxes are placed into.

The box model is very versatile. Nevertheless we will generalize the
notion of a box to mean the bounding box of a polygon. Typesetting in
irregular shapes is a feature available in most modern systems, e.g.
when letting text flow around a non-rectangular illustration.

This module deals with rectangular boxes, starting at the glyph level.
The notation oftentimes follows the one introduced by the TeX typesetting
system.

BSD License

Copyright (c) 2017–2018, Norbert Pillmayer

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
*/
package box

import (
	"image/color"

	"github.com/npillmayer/gotype/gtcore/dimen"
	"github.com/npillmayer/gotype/gtcore/font"
)

/*
---------------------------------------------------------------------------


----------------------------------------------------------------------


/* Box styling: We follow the CSS paradigm for boxes. Boxes are stylable
 * objects which have dimensions, spacing, borders and colors.
 *
 * Some boxes may just implement a subset of the styling parameters. Most
 * notably this holds for glyphs: Glyphs may have styled their content only.
 * No border or additional spacing is possible with glyphs.
*/

// Type for styled content (color only)
type ContentStyling struct {
	Foreground        color.Color
	Background        color.Color
	StylingIdentifier string // e.g., a CSS ID
	StylingClass      string // e.g., a CSS class
}

// Type for a fully stylable box
type StyledBox struct {
	//Box
	Style   *ContentStyling
	Border  *BorderStyle
	Padding dimen.Dimen // inside of border
	Spacing dimen.Dimen // outside of border
}

// A type for simple borders. Currently no line style variants are supported.
type BorderStyle struct {
	LineColor     color.Color
	LineThickness dimen.Dimen
}

/* Wikipedia: In typography, a glyph [...] is an elemental symbol within
 * an agreed set of symbols, intended to represent a readable character
 * for the purposes of writing. [ Copyright (c) Wikipedia.com, 2017 ]
 */

// A box for glyphs. Glyphs are content-stylable only (no borders).
type Glyph struct {
	//Box
	Style    *ContentStyling
	Typecase *font.TypeCase
	CharPos  rune
}

// Boxed for typesetting, similar to TeX's \hbox and \vbox.
type TypesetBox struct {
	//Box
	Style *ContentStyling
	//Cord  *Khipu
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
