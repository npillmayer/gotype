/*
Package vectorcanvas implements a canvas type for vector drawings.

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
	"image/color"

	"github.com/npillmayer/gotype/backend/gfx"
)

type VCanvas struct {
	w, h     float64
	strokes  []*stroke
	gfxState *graphicsState
}

func New(w, h float64) *VCanvas {
	canv := &VCanvas{}
	canv.w = w
	canv.h = h
	return canv
}

func (canv *VCanvas) W() float64 {
	return canv.w
}

func (canv *VCanvas) H() float64 {
	return canv.h
}

func (canv *VCanvas) AddContour(c gfx.DrawableContour, linew float64,
	linecol color.Color, fillcol color.Color) {
	//
	g := canv.pushGraphicsStateIfNecessary(linew, linecol, fillcol)
	canv.strokes = append(canv.strokes, &stroke{
		contour: c,
		gfxCtx:  g,
	})
}
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
		g := &graphicsState{next: canv.gfxState}
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
