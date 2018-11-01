/*
Package textshaping implements the task of selecting glyphs for
Unicode characters.

To understand what a text shaper does, please have a look at
http://www.manpagez.com/html/harfbuzz/harfbuzz-/what-is-harfbuzz.php

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
	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/gtcore/font"
)

// We trace to the commands-tracer. // TODO: create a font tracer?
var T tracing.Trace = tracing.CommandTracer

// Script IDs are copied from Harfbuzz. We don't have any sensible
// preference for IDs, thus there's no harm in being in sync with
// our main shaping engine.  Just do not rely on it outside of
// Harfbuzz's scope.
type ScriptID uint32

//go:generate stringer -type=ScriptID
const (
	Invalid               ScriptID = 0
	Arabic                ScriptID = 1098015074
	ImperialAramaic       ScriptID = 1098018153
	Armenian              ScriptID = 1098018158
	Avestan               ScriptID = 1098281844
	Balinese              ScriptID = 1113681001
	Bamum                 ScriptID = 1113681269
	Batak                 ScriptID = 1113683051
	Bengali               ScriptID = 1113943655
	Bopomofo              ScriptID = 1114599535
	Brahmi                ScriptID = 1114792296
	Braille               ScriptID = 1114792297
	Buginese              ScriptID = 1114990441
	Buhid                 ScriptID = 1114990692
	Chakma                ScriptID = 1130457965
	CanadianSyllabics     ScriptID = 1130458739
	Carian                ScriptID = 1130459753
	Cham                  ScriptID = 1130914157
	Cherokee              ScriptID = 1130915186
	Coptic                ScriptID = 1131376756
	Cypriot               ScriptID = 1131442804
	Cyrillic              ScriptID = 1132032620
	Devanagari            ScriptID = 1147500129
	Deseret               ScriptID = 1148416628
	EgyptianHieroglyphs   ScriptID = 1164409200
	Ethiopic              ScriptID = 1165256809
	Georgian              ScriptID = 1197830002
	Glagolitic            ScriptID = 1198285159
	Gothic                ScriptID = 1198486632
	Greek                 ScriptID = 1198679403
	Gujarati              ScriptID = 1198877298
	Gurmukhi              ScriptID = 1198879349
	Hangul                ScriptID = 1214344807
	Han                   ScriptID = 1214344809
	Hanunoo               ScriptID = 1214344815
	Hebrew                ScriptID = 1214603890
	Hiragana              ScriptID = 1214870113
	OldItalic             ScriptID = 1232363884
	Javanese              ScriptID = 1247901281
	KayahLi               ScriptID = 1264675945
	Katakana              ScriptID = 1264676449
	Kharoshthi            ScriptID = 1265131890
	Khmer                 ScriptID = 1265134962
	Kannada               ScriptID = 1265525857
	Kaithi                ScriptID = 1265920105
	TaiTham               ScriptID = 1281453665
	Lao                   ScriptID = 1281453935
	Latin                 ScriptID = 1281455214
	Lepcha                ScriptID = 1281716323
	Limbu                 ScriptID = 1281977698
	LinearB               ScriptID = 1281977954
	Lisu                  ScriptID = 1281979253
	Lycian                ScriptID = 1283023721
	Lydian                ScriptID = 1283023977
	Mandaic               ScriptID = 1298230884
	MeroiticCursive       ScriptID = 1298494051
	MeroiticHieroglyphs   ScriptID = 1298494063
	Malayalam             ScriptID = 1298954605
	Mongolian             ScriptID = 1299148391
	MeeteiMayek           ScriptID = 1299473769
	Myanmar               ScriptID = 1299803506
	Nko                   ScriptID = 1315663727
	Ogham                 ScriptID = 1332175213
	OlChiki               ScriptID = 1332503403
	OldTurkic             ScriptID = 1332898664
	Oriya                 ScriptID = 1332902241
	Osmanya               ScriptID = 1332964705
	PhagsPa               ScriptID = 1349017959
	InscriptionalPahlavi  ScriptID = 1349020777
	Phoenician            ScriptID = 1349021304
	Miao                  ScriptID = 1349284452
	InscriptionalParthian ScriptID = 1349678185
	Rejang                ScriptID = 1382706791
	Runic                 ScriptID = 1383427698
	Samaritan             ScriptID = 1398893938
	OldSouthArabian       ScriptID = 1398895202
	Saurashtra            ScriptID = 1398895986
	Shavian               ScriptID = 1399349623
	Sharada               ScriptID = 1399353956
	Sinhala               ScriptID = 1399418472
	SoraSompeng           ScriptID = 1399812705
	Sundanese             ScriptID = 1400204900
	SylotiNagri           ScriptID = 1400466543
	Syriac                ScriptID = 1400468067
	Tagbanwa              ScriptID = 1415669602
	Takri                 ScriptID = 1415670642
	TaiLe                 ScriptID = 1415670885
	NewTaiLue             ScriptID = 1415670901
	Tamil                 ScriptID = 1415671148
	TaiViet               ScriptID = 1415673460
	Telugu                ScriptID = 1415933045
	Tifinagh              ScriptID = 1415999079
	Tagalog               ScriptID = 1416064103
	Thaana                ScriptID = 1416126817
	Thai                  ScriptID = 1416126825
	Tibetan               ScriptID = 1416192628
	Ugaritic              ScriptID = 1432838514
	Vai                   ScriptID = 1449224553
	OldPersian            ScriptID = 1483761007
	Cuneiform             ScriptID = 1483961720
	Yi                    ScriptID = 1500080489
	Inherited             ScriptID = 1516858984
	Common                ScriptID = 1517910393
	Unknown               ScriptID = 1517976186
)

// Direction to typeset text in.
type TextDirection int32

// Direction to typeset text in. We use a generator to produce a stringer
// for this enum.
//
//go:generate stringer -type=TextDirection
const (
	LeftToRight TextDirection = 0
	RightToLeft               = 1
	TopToBottom               = 2
	BottomToTop               = 3
)

// Interface for single glyphs in a glyph sequence.
// The 'cluster' field corresponds to the code-point position within the
// original string.
type GlyphInfo interface {
	Glyph() rune
	Cluster() int
	XAdvance() float64  // should be internal typesetter dimen
	YAdvance() float64  // should be internal typesetter dimen
	XPosition() float64 // should be internal typesetter dimen
	YPosition() float64 // should be internal typesetter dimen
}

// Interface for a sequence of glyphs as returned by a text shaper.
type GlyphSequence interface {
	GlyphCount() int
	GetGlyphInfoAt(pos int) GlyphInfo
}

// A text shaper creates a sequence of glyphs from a sequence of
// Unicode code-points. Glyphs are taken from a font (a set of
// of sorts from a type case), with specific dimens.
//
// Sometimes the shaping of text depends on the language in use.
// Shapers may be able to react to this kind of information.
type TextShaper interface {
	Shape(text string, typecase *font.TypeCase) GlyphSequence
	SetScript(scr ScriptID)
	SetDirection(dir TextDirection)
	SetLanguage() // TODO: what is a language in Go?
}
