/*
Package gfx implements a backend for graphics. The backend connects to various
drawing engines (e.g., raster, SVG, PDF, ...).

We abstract the details away with an interface Canvas, which represents
a rectangular drawing area. The main operation on a Canvas is AddContour,
the drawing or filling of a path.

Clients often won't use Canvas directly, but rather an enclosing struct type
Picture, which holds a Canvas plus some administrative information.


BSD License

Copyright (c) 2017-21, Norbert Pillmayer <norbert@pillmayer.com>

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
package gfx

import (
	"image/color"
	"math/cmplx"

	"github.com/npillmayer/arithm"
	"github.com/npillmayer/arithm/path"
	"github.com/npillmayer/arithm/polygon"
	"github.com/npillmayer/schuko/tracing"
)

// G traces to the graphics tracer.
func G() tracing.Trace {
	return tracing.GraphicsTracer
}

// Canvas is an interface type for a drawing canvas.
type Canvas interface {
	W() float64                                                    // width
	H() float64                                                    // height
	AddContour(DrawableContour, float64, color.Color, color.Color) // filldraw with pen with colors
	SetOption(int)                                                 // set a drawing option
}

// Picture is the type for the backend to work with.
type Picture struct {
	Name         string
	canvas       Canvas
	currentColor color.Color
	currentPen   *Pen
	output       OutputRoutine
}

// Pen is a drawing pen.
type Pen struct {
	diameter float64
	style    int
}

// PENCIRCLE is a pen with a round nib
const PENCIRCLE int = 1

// PENSQUARE is a pen with a square nib
const PENSQUARE int = 2

// NewPicture creates a
// new Picture. Caller has to provide an identifier, a width,
// a height and a Canvas.
func NewPicture(name string, canvas Canvas) *Picture {
	pic := &Picture{
		Name:         name,
		canvas:       canvas,
		currentColor: color.Black,
		currentPen:   NewPencircle(1.0),
	}
	return pic
}

// NewPencircle creates a new round pen.
func NewPencircle(diam float64) *Pen {
	pen := &Pen{
		diameter: diam,
		style:    PENCIRCLE,
	}
	return pen
}

// NewPensquare create a new square pen.
func NewPensquare(diam float64) *Pen {
	pen := &Pen{
		diameter: diam,
		style:    PENSQUARE,
	}
	return pen
}

// Draw a line. Uses the current pen.
func (pic *Picture) Draw(contour DrawableContour) {
	pendiam, _ := pic.currentPen.diameter.Float64()
	pic.canvas.SetOption(pic.currentPen.style)
	pic.canvas.AddContour(contour, pendiam, pic.currentColor, nil)
}

// Fill a closed path. Uses the current pen.
func (pic *Picture) Fill(contour DrawableContour) {
	pendiam, _ := pic.currentPen.diameter.Float64()
	pic.canvas.AddContour(contour, pendiam, nil, pic.currentColor)
}

// FillDraw is
// MetaPost's filldraw command. Uses the current pen and color.
func (pic *Picture) FillDraw(contour DrawableContour) {
	pendiam, _ := pic.currentPen.diameter.Float64()
	pic.canvas.SetOption(pic.currentPen.style)
	pic.canvas.AddContour(contour, pendiam, pic.currentColor, pic.currentColor)
}

// SetPen sets
// the current pen. Will be used for subsequent drawing operations.
func (pic *Picture) SetPen(pen *Pen) {
	pic.currentPen = pen
}

// SetColor sets
// the current color. Will be used for subsequent drawing operations.
func (pic *Picture) SetColor(color color.Color) {
	pic.currentColor = color
}

// === Drawable Contour ======================================================

// DrawableContour is an interface for a stroke.
// The central operation on Canvases is AddContour, the drawing and/or filling
// of a path. Different clients may have a different understanding of what a
// path is and how to store it. Interface DrawableContour tries to be a common
// denomiator for drawing operations on paths.
//
// The 3 parameters for ToNextKnot() are the target point, together with 2
// optional control points for spline curves (quadratic or cubic).
//
type DrawableContour interface {
	IsCycle() bool
	Start() arithm.Pair
	ToNextKnot() (k arithm.Pair, c1 arithm.Pair, c2 arithm.Pair)
}

// NewDrawablePath creates an immutable adapter to contours from cubic splines.
func NewDrawablePath(path path.HobbyPath, controls path.SplineControls) DrawableContour {
	pdrw := &pathdrawer{p: path, c: controls}
	if path.IsCycle() {
		pdrw.n = path.N() + 1
	} else {
		pdrw.n = path.N()
	}
	return pdrw
}

// internal type: immutable adapter to contours from cubic splines
type pathdrawer struct {
	p       path.HobbyPath
	c       path.SplineControls
	current int
	n       int
}

// implement interface DrawableContour
func (pdrw *pathdrawer) IsCycle() bool {
	return pdrw.p.IsCycle()
}

// implement interface DrawableContour
func (pdrw *pathdrawer) Start() arithm.Pair {
	G().Debugf("path start at %s", arithm.Pair(pdrw.p.Z(0)))
	pdrw.current = 0
	return pdrw.p.Z(0)
}

// implement interface DrawableContour
func (pdrw *pathdrawer) ToNextKnot() (arithm.Pair, arithm.Pair, arithm.Pair) {
	pdrw.current++
	if pdrw.current >= pdrw.n {
		G().Debugf("path has no more knots")
		return nil, nil, nil
	}
	c1, c2 := pdrw.c.PostControl(pdrw.current-1), pdrw.c.PreControl(pdrw.current%(pdrw.p.N()))
	if pdrw.current < pdrw.n && !cmplx.IsNaN(c1) {
		G().Debugf("path next  at %s", arithm.Pair(pdrw.p.Z(pdrw.current)))
		G().Debugf("     controls %s and %s", arithm.Pair(c1), arithm.Pair(c2))
	} else {
		G().Debugf("path next  at %s", arithm.Pair(pdrw.p.Z(pdrw.current)))
	}
	return pdrw.p.Z(pdrw.current), c1, c2
}

// NewDrawablePolygon creates an immutable bridge to contours from polygons.
func NewDrawablePolygon(pg polygon.Polygon) DrawableContour {
	panic("not yet implemented: Contour for Polygon")
	return nil
}

// === Output Routine ========================================================

// OutputRoutine is an interface for output routines to ship out
// completed images.
type OutputRoutine interface {
	Shipout(pic *Picture, gfxFormat string) bool
}
