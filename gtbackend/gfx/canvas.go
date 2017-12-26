/*
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

 * Backend for graphics. The backend connects to various drawing
 * engines (e.g., raster, SVG, PDF, ...).
 * We abstract the details away with an interface Canvas, which represents
 * a rectangular drawing area. The main operation on a Canvas is AddContour,
 * the drawing or filling of a path.
 *
 * We do not use Canvas directly, but rather an enclosing struct type Picture,
 * which holds a Canvas plus some administrative information.

*/

package gfx

import (
	"image"
	"image/color"

	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
	dec "github.com/shopspring/decimal"
)

var G tracing.Trace = tracing.GraphicsTracer

// Interface for output routines to ship out completed images
type OutputRoutine interface {
	Shipout(string, image.Image)
}

// We need to create drawing canvases
var GlobalCanvasFactory CanvasFactory

// We need to create drawing canvases
type CanvasFactory interface {
	New(string, float64, float64) Canvas
}

// The interface type for drawing canvases
type Canvas interface {
	W() float64                                                // width
	H() float64                                                // height
	AddContour(arithm.Path, float64, color.Color, color.Color) // filldraw with pen with colors
	AsImage() image.Image                                      // get a stdlib image
	SetOption(int)                                             // set a drawing option
}

// This is the type for the backend to work with
type Picture struct {
	Name         string
	canvas       Canvas
	currentColor color.Color
	currentPen   *Pen
}

// A drawing pen
type Pen struct {
	diameter dec.Decimal
	style    int
}

// A pen with a round nib
const PENCIRCLE int = 1

// A pen with a square nib
const PENSQUARE int = 2

/* Create a new Picture. Caller has to provide an itentifier, a width and
 * height.
 */
func NewPicture(name string, w float64, h float64) *Picture {
	if GlobalCanvasFactory == nil {
		G.Error("no canvas factory set")
		panic("cannot create picture: no canvas factory")
		return nil
	} else {
		canvas := GlobalCanvasFactory.New(name, w, h)
		pic := &Picture{
			Name:         name,
			canvas:       canvas,
			currentColor: color.Black,
			currentPen:   NewPencircle(arithm.ConstOne),
		}
		return pic
	}
}

/* Create a new round pen.
 */
func NewPencircle(diam dec.Decimal) *Pen {
	pen := &Pen{
		diameter: diam,
		style:    PENCIRCLE,
	}
	return pen
}

/* Create a new square pen.
 */
func NewPensquare(diam dec.Decimal) *Pen {
	pen := &Pen{
		diameter: diam,
		style:    PENSQUARE,
	}
	return pen
}

/* Draw a line. Uses the current pen.
 */
func (pic *Picture) Draw(path arithm.Path) {
	pendiam, _ := pic.currentPen.diameter.Float64()
	pic.canvas.SetOption(pic.currentPen.style)
	pic.canvas.AddContour(path, pendiam, pic.currentColor, nil)
}

/* Fill a closed path. Uses the current pen.
 */
func (pic *Picture) Fill(path arithm.Path) {
	pendiam, _ := pic.currentPen.diameter.Float64()
	pic.canvas.AddContour(path, pendiam, nil, pic.currentColor)
}

/* MetaPost's filldraw command. Uses the current pen and color.
 */
func (pic *Picture) FillDraw(path arithm.Path) {
	pendiam, _ := pic.currentPen.diameter.Float64()
	pic.canvas.SetOption(pic.currentPen.style)
	pic.canvas.AddContour(path, pendiam, pic.currentColor, pic.currentColor)
}

/* Set the current pen. Will be used for future drawing operations.
 */
func (pic *Picture) SetPen(pen *Pen) {
	pic.currentPen = pen
}

/* Set the current color. Will be used for future drawing operations.
 */
func (pic *Picture) SetColor(color color.Color) {
	pic.currentColor = color
}

/* Return the current canvas as an image (standard library).
 */
func (pic *Picture) AsImage() image.Image {
	return pic.canvas.AsImage()
}
