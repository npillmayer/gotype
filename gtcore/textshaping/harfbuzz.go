/*
CGo wrapper for the Harfbuzz text shaping library.

----------------------------------------------------------------------

BSD License

Copyright (c) 2017-2018, Norbert Pillmayer

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
package textshaping

import (
	"fmt"

	"github.com/npillmayer/gotype/gtcore/font"
)

// Harfbuzz is the de-facto standard for text shaping.
// For further information see
// https://www.freedesktop.org/wiki/Software/HarfBuzz .
//
// A remark to the use of pointers to Harfbuzz-objects: Harfbuzz
// does its own memory management and we must avoid interfering with it.
// The Go garbage collector will therefore be unaware of the memory
// managed by Harfbuzz (in the worst of cases, a fancy Go garbage collector
// may re-locate memory). To hide Harfbuzz-memory from Go, we will use
// 'uintptr' variables instead of 'unsafe.Pointer's.
//
// The downside of this is the need to free() memory whenever we
// hand a Harfbuzz-shaper to GC.
type Harfbuzz struct {
	buffer    uintptr       // central data structure for Harfbuzz
	direction TextDirection // L-to-R, R-to-L, T-to-B
	script    ScriptID      // i.e., Latin, Arabic, Korean, ...
}

// Create a new Harfbuzz text shaper, fully initialized.
// Defaults are for Latin script, left-to-right.
func NewHarfbuzz() *Harfbuzz {
	hb := &Harfbuzz{}
	hb.buffer = allocHBBuffer()
	hb.direction = LeftToRight
	setHBBufferDirection(hb.buffer, hb.direction)
	hb.script = Latin
	setHBBufferScript(hb.buffer, hb.script)
	return hb
}

// Cache for font structures prepared for Harfbuzz.
// Harfbuzz uses its own font structure, different from ours.
// Unfortunately this duplicates the binary data of the font.
var harfbuzzFontCache map[*font.TypeCase]uintptr

// TODO: make cache thread-safe
// TODO: return error
func (hb *Harfbuzz) findFont(typecase *font.TypeCase) uintptr {
	var hbfont uintptr
	if harfbuzzFontCache == nil {
		harfbuzzFontCache = make(map[*font.TypeCase]uintptr)
	}
	if hbfont = harfbuzzFontCache[typecase]; hbfont == 0 {
		if hbfont = makeHBFont(typecase); hbfont != 0 {
			harfbuzzFontCache[typecase] = hbfont
		}
	}
	return hbfont
}

// Implement TextShaper interface.
func (hb *Harfbuzz) SetScript(scr ScriptID) {
	setHBBufferScript(hb.buffer, scr)
}

// Implement TextShaper interface.
func (hb *Harfbuzz) SetDirection(dir TextDirection) {
	setHBBufferDirection(hb.buffer, dir)
}

// Implement TextShaper interface.
// Harfbuzz doesn't evaluate a language parameter; method is a NOP.
func (hb *Harfbuzz) SetLanguage() {
}

// Implement TextShaper interface.
//
// This is where all the heavy lifting is done. We input a font and a
// string of Unicode code-points, and receive a list of glyphs.
func (hb *Harfbuzz) Shape(text string, typecase *font.TypeCase) GlyphSequence {
	var hbfont uintptr
	hbfont = hb.findFont(typecase)
	if hbfont == 0 {
		panic(fmt.Sprintf("*** cannot find/create Harfbuzz font for [%s]\n",
			typecase.ScalableFontParent().Fontname))
	}
	harfbuzzShape(hb.buffer, text, hbfont)
	seq := getHBGlyphInfo(hb.buffer)
	return seq
}

func (hb *Harfbuzz) GlyphSequenceString(typecase *font.TypeCase, seq GlyphSequence) string {
	var hbfont uintptr
	hbfont = hb.findFont(typecase)
	if hbfont == 0 {
		panic(fmt.Sprintf("*** cannot find/create Harfbuzz font for [%s]\n",
			typecase.ScalableFontParent().Fontname))
	}
	s := hbGlyphString(hbfont, seq.(*hbGlyphSequence))
	return s
}

var _ = TextShaper(&Harfbuzz{})
