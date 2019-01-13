/*
Package svg implements operations on/for SVG graphics.


BSD License

Copyright (c) 2017, Norbert Pillmayer <norbert@pillmayer.com>

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
package svg

import (
	"fmt"
	"image/color"
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/npillmayer/gotype/backend/gfx"
	"github.com/npillmayer/gotype/core/arithmetic"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/dimen"
)

func T() tracing.Trace {
	return tracing.GraphicsTracer
}

type SvgPicture struct {
	svgCanvas *svg.SVG
	w         dimen.Dimen
	h         dimen.Dimen
}

func NewPicture(wr io.Writer, w int, h int) *SvgPicture {
	pic := &SvgPicture{}
	pic.svgCanvas = svg.New(wr)
	pic.svgCanvas.Start(w, h)
	pic.svgCanvas.Grid(0, 0, w, h, int(10*dimen.BP))
	return pic
}

func (pic *SvgPicture) W() float64 {
	return float64(pic.w)
}

func (pic *SvgPicture) H() float64 {
	return float64(pic.h)
}

func (pic *SvgPicture) Shipout() {
	pic.svgCanvas.End()
}

// Filldraw with pen with colors:
// Add a stroke to the graphics context.
// If linecol is non-void, a line will be drawn. If fillcol is non-void and
// the path is closed, the path's inside area will be filled.
//
// If tracing mode is DebugLevel red dots are drawn for the contour's knots.
//
func (pic *SvgPicture) AddContour(contour gfx.DrawableContour, linethickness float64,
	linecol color.Color, fillcol color.Color) {
	//
	T().P("fmt", "SVG").Debugf("add contour")
	pt := contour.Start()
	if pt == nil {
		return
	}
	px, py := arithmetic.Pr2Pt(pt)
	svgpath := fmt.Sprintf("M%d,%d ", int(px), int(py)) // move to
	for pt != nil {
		var c1, c2 arithmetic.Pair
		pt, c1, c2 = contour.ToNextKnot()
		svgpath += svgpathC(pt, c1, c2)
	}
	T().Debugf("SVG drawing path '%s'", svgpath)
	style := fmt.Sprintf("stroke:#%s;stroke-width:%d", hexcolor(linecol), int(linethickness))
	if contour.IsCycle() { // close path ?
		style += fmt.Sprintf(";fill:#%s", hexcolor(fillcol))
	}
	pic.svgCanvas.Path(svgpath, style)
}

func svgpathC(pt arithmetic.Pair, c1 arithmetic.Pair, c2 arithmetic.Pair) string {
	if pt != nil {
		x, y := arithmetic.Pr2Pt(pt)
		if c1 == nil && c2 == nil {
			return fmt.Sprintf("L%d,%d ", int(x), int(y))
		} else {
			c1x, c1y := arithmetic.Pr2Pt(c1)
			c2x, c2y := arithmetic.Pr2Pt(c2)
			return fmt.Sprintf("C%d,%d %d,%d %d,%d ",
				int(c1x), int(c1y), int(c2x), int(c2y), int(x), int(y))
		}
	}
	return ""
}

func svgpathL(pt arithmetic.Pair) string {
	x, y := arithmetic.Pr2Pt(pt)
	return fmt.Sprintf("L%d,%d ", int(x), int(y))
}

func (pic *SvgPicture) SetOption(int) {
	//
}

func hexcolor(c color.Color) string {
	if c == nil {
		return "ffffff"
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%02x%02x%02x", r, g, b)
}
