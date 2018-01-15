/*
Package png implements a canvas type as a bridge to the Go Graphics
drawing package.

----------------------------------------------------------------------

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

----------------------------------------------------------------------
*/
package png

import (
	"image"
	"image/color"
	ospng "image/png"
	"io"

	"github.com/fogleman/gg"
	"github.com/npillmayer/gotype/gtbackend/gfx"
	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
)

// We are tracing to the graphics trace
var T tracing.Trace = tracing.GraphicsTracer

// A canvas type as a bridge to the Go Graphics drawing package.
type GGCanvas struct {
	context *gg.Context
}

// check assignability
var _ gfx.Canvas = &GGCanvas{}

func init() {
	gfx.RegisterCanvasCreator("PNG", NewCanvas)
}

func NewCanvas(width, height float64) gfx.Canvas {
	ggc := &GGCanvas{}
	ggc.context = newGC(width, height)
	return ggc
}

// Interface Canvas. Return graphic's width.
func (ggc *GGCanvas) W() float64 {
	return float64(ggc.context.Width())
}

// Interface Canvas. Return graphic's height.
func (ggc *GGCanvas) H() float64 {
	return float64(ggc.context.Height())
}

/*
Interface Canvas. Add a stroke to the graphics context.
If linecol is non-void, a line will be drawn. If fillcol is non-void and
the path is closed, the path's inside area will be filled.
*/
func (ggc *GGCanvas) AddContour(contour gfx.DrawableContour, linethickness float64,
	linecol color.Color, fillcol color.Color) {
	//
	T.P("fmt", "PNG").Debug("add contour")
	pt := contour.Start()
	px, _ := pt.XPart().Float64()
	py, _ := pt.YPart().Float64()
	ggc.context.MoveTo(px, ggc.yflip(py))
	for pt != nil {
		var c1, c2 arithm.Pair
		pt, c1, c2 = contour.ToNextKnot()
		if pt != nil {
			px, _ = pt.XPart().Float64()
			py, _ = pt.YPart().Float64()
			if c1 == nil && c2 == nil {
				ggc.context.LineTo(px, ggc.yflip(py))
			} else {
				c1x, _ := c1.XPart().Float64()
				c1y, _ := c1.YPart().Float64()
				c2x, _ := c2.XPart().Float64()
				c2y, _ := c2.YPart().Float64()
				ggc.context.CubicTo(c1x, ggc.yflip(c1y), c2x, ggc.yflip(c2y), px, ggc.yflip(py))
			}
		}
	}
	if contour.IsCycle() {
		// close path ?
		if fillcol != nil {
			ggc.context.SetColor(fillcol)
			ggc.context.Fill()
		}
	}
	if linecol != nil {
		ggc.context.SetLineWidth(linethickness)
		ggc.context.SetColor(linecol)
		ggc.context.Stroke()
	}
}

/*
Interface Canvas. Add a stroke to the graphics context.
If linecol is non-void, a line will be drawn. If fillcol is non-void and
the path is closed, the path's inside area will be filled.
func (ggc *GGCanvas) _addContour(path gfx.DrawableContour, linethickness float64,
	linecol color.Color, fillcol color.Color) {
	//
	T.P("fmt", "PNG").Debug("add contour")
	l := path.Length()
	pt := arithm.Pr2Pt(path.GetPoint(0))
	ggc.context.MoveTo(pt.X, ggc.yflip(pt.Y))
	for i := 1; i < l; i++ {
		pt = arithm.Pr2Pt(path.GetPoint(i))
		ggc.context.LineTo(pt.X, ggc.yflip(pt.Y))
	}
	if path.IsCycle() {
		pt = arithm.Pr2Pt(path.GetPoint(0))
		ggc.context.LineTo(pt.X, ggc.yflip(pt.Y))
		if fillcol != nil {
			ggc.context.SetColor(fillcol)
			ggc.context.Fill()
		}
	}
	if linecol != nil {
		ggc.context.SetLineWidth(linethickness)
		ggc.context.SetColor(linecol)
		ggc.context.Stroke()
	}
}
*/

// GG draw uses a top-down y-coordinate. Flip it.
func (ggc *GGCanvas) yflip(y float64) float64 {
	return ggc.H() - y
}

// Interface Canvas. Return the drawing as a stdlib image.
func (ggc *GGCanvas) AsImage() image.Image {
	T.P("fmt", "PNG").Debug("converting image")
	return ggc.context.Image()
}

// Shipout the picture to an output device.
func (ggc *GGCanvas) Shipout(w io.Writer) bool {
	T.Debugf("shipping out PNG image")
	img := ggc.AsImage()
	if err := ospng.Encode(w, img); err != nil {
		T.Errorf("file error: %v", err.Error())
		return false
	}
	return true
}

/*
Interface Canvas. Set drawing options. Currently only gfx.PENSQUARE and
gfx.PENCIRCLE are understood. They determine linecap and and linejoin
parameters.
*/
func (ggc *GGCanvas) SetOption(opts int) {
	if opts == gfx.PENSQUARE {
		ggc.context.SetLineCapSquare()
		ggc.context.SetLineJoinBevel()
	} else {
		ggc.context.SetLineCapRound()
		ggc.context.SetLineJoinRound()
	}
}

// create a new GG graphics context.
func newGC(w float64, h float64) *gg.Context {
	dc := gg.NewContext(int(w), int(h))
	dc.SetColor(color.Black)
	dc.SetLineWidth(1.0)
	dc.MoveTo(0, 0) // for debugging purposes only, TODO remove this
	dc.LineTo(w, 0) // should be w-1 and h-1, but seems to be a bug
	dc.LineTo(w, h)
	dc.LineTo(0, h)
	dc.LineTo(0, 0)
	dc.Stroke()
	return dc
}
