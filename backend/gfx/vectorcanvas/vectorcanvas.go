/*
Package vectorcanvas implements a vanilla canvas type for vector drawings.

The canvas stores the strokes in a generic fashion.
For output, clients convert the strokes into other formats (SVG, PDF, ...)
by calling the appropriate ToXXX(...) function.


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
package vectorcanvas

import (
	"fmt"
	"image/color"

	"github.com/npillmayer/gotype/backend/gfx"
	"github.com/npillmayer/gotype/backend/print/pdf/pdfapi"
	"github.com/npillmayer/gotype/core/arithmetic"
)

// VCanvas is a generic vector canvas.
type VCanvas struct {
	w, h     float64
	strokes  []*stroke
	gfxState *graphicsState
}

// New creates a new canvas with a given size.
// Size is interpreted as 'big points', i.e. 1/72 of an inch.
func New(w, h float64) *VCanvas {
	canv := &VCanvas{}
	canv.w = w
	canv.h = h
	return canv
}

// W returns the width of the canvas in 'big points'.
func (canv *VCanvas) W() float64 {
	return canv.w
}

// W returns the height of the canvas in 'big points'.
func (canv *VCanvas) H() float64 {
	return canv.h
}

// AddContour adds a drawable path to the canvas.
// Clients provide the drawing parameters, where colors may be nil.
//
func (canv *VCanvas) AddContour(c gfx.DrawableContour, linew float64,
	linecol color.Color, fillcol color.Color) {
	//
	g := canv.pushGraphicsStateIfNecessary(linew, linecol, fillcol)
	canv.strokes = append(canv.strokes, &stroke{
		contour: c,
		gfxCtx:  g,
	})
	fmt.Printf("new contour with gfx = %v\n", g)
}

// SetOption currently does nothing.
func (canv *VCanvas) SetOption(int) {
	// TODO
}

func (canv *VCanvas) pushGraphicsStateIfNecessary(linew float64, linecol color.Color,
	fillcol color.Color) *graphicsState {
	//
	g := canv.gfxState
	if g == nil ||
		linew != g.penwidth || linecol != g.color || fillcol != g.fillcolor {
		//
		g = &graphicsState{next: canv.gfxState}
		g.color = linecol
		g.fillcolor = fillcol
		g.penwidth = linew
	}
	return g
}

// ----------------------------------------------------------------------

type graphicsState struct {
	color     color.Color
	fillcolor color.Color
	penwidth  float64
	next      *graphicsState
}

type stroke struct {
	contour gfx.DrawableContour
	gfxCtx  *graphicsState
}

// ----------------------------------------------------------------------

// ToPDF converts a generic vector drawing into PDF format.
// Callers provide a PDF canvas to draw into.
func (canv *VCanvas) ToPDF(pdfcanv *pdfapi.Canvas) error {
	for _, st := range canv.strokes {
		contour := st.contour
		pr := contour.Start()
		if pr == nil {
			continue // skip empty contour
		}
		pt := pr2PdfPt(pr)
		P := &pdfapi.Path{}
		P.MoveTo(pt)
		for pr != nil {
			var c1, c2 arithmetic.Pair
			pr, c1, c2 = contour.ToNextKnot()
			pdfpathC(P, pr, c1, c2)
		}
		ctx := st.gfxCtx
		if ctx != nil && contour.IsCycle() {
			if ctx.color != nil {
				pdfcanv.SetStrokeColor(ctx.color)
				pdfcanv.SetLineWidth(pdfapi.Unit(ctx.penwidth))
				if ctx.fillcolor != nil {
					pdfcanv.SetFillColor(ctx.fillcolor)
					pdfcanv.FillStroke(P)
				} else {
					pdfcanv.Stroke(P)
				}
			} else {
				pdfcanv.SetFillColor(ctx.fillcolor)
				pdfcanv.Fill(P)
			}
		} else {
			pdfcanv.Stroke(P)
		}
	}
	return nil
}

func pdfpathC(P *pdfapi.Path, pr arithmetic.Pair, c1 arithmetic.Pair, c2 arithmetic.Pair) {
	if pr != nil {
		pt := pr2PdfPt(pr)
		if c1 == nil && c2 == nil {
			P.LineTo(pt)
		} else {
			c1pt := pr2PdfPt(c1)
			c2pt := pr2PdfPt(c2)
			P.CurveTo(c1pt, c2pt, pt)
		}
	}
}

func pr2PdfPt(pr arithmetic.Pair) pdfapi.Point {
	x, y := arithmetic.Pr2Pt(pr)
	return pdfapi.Point{pdfapi.Unit(x), pdfapi.Unit(y)}
}
