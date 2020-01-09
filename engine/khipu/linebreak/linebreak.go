/*
Package linebreak collects types for line-breaking.

BSD License

Copyright (c) 2017–20, Norbert Pillmayer

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package linebreak

import (
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/khipu"
)

// WSS (width stretch & shrink) is a type to hold an elastic width (of text).
type WSS struct {
	W   dimen.Dimen
	Min dimen.Dimen
	Max dimen.Dimen
}

// Spread returns the width's of an elastic WSS.
func (wss WSS) Spread() (w dimen.Dimen, min dimen.Dimen, max dimen.Dimen) {
	return wss.W, wss.Min, wss.Max
}

// SetFromKnot sets the width's of an elastic WSS from a knot.
func (wss WSS) SetFromKnot(knot khipu.Knot) WSS {
	wss.W = knot.W()
	wss.Min = knot.MinW()
	wss.Max = knot.MaxW()
	return wss
}

// Add adds dimensions from a given WSS to wss, creating a new WSS.
func (wss WSS) Add(other WSS) WSS {
	return WSS{
		W:   wss.W + other.W,
		Min: wss.Min + other.Min,
		Max: wss.Max + other.Max,
	}
}

// InfinityDemerits is the worst demerit value possible.
const InfinityDemerits int32 = 10000

// InfinityMerits is the best (most desirable) demerit value possible.
const InfinityMerits int32 = -10000

// CapDemerits caps a demerit value at infinity.
func CapDemerits(d int32) int32 {
	if d > InfinityDemerits {
		d = InfinityDemerits
	}
	return d
}

// --- Interfaces -------------------------------------------------------

// Parshape is a type to return the line length for a given line number.
type Parshape interface {
	LineLength(int) dimen.Dimen
}

// Cursor is a type to iterate over a khipu.
type Cursor interface {
	Next() bool
	Knot() khipu.Knot
	Mark() khipu.Mark
}

type rectParshape int

func (r rectParshape) LineLength(int) dimen.Dimen {
	return (dimen.PT * dimen.Dimen(r))
}

// RectangularParshape returns a Parshape for paragraphs of constant line length.
func RectangularParshape(linelen int) Parshape {
	return rectParshape(linelen)
}

// ----------------------------------------------------------------------

func PrintParagraphBreaks(k *khipu.Khipu, breaks []khipu.Mark) string {
	//cursor := khipu.NewCursor(k)
	return ""
}
