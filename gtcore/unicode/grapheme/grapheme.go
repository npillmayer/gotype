/*
Package grapheme implements Unicode Annex #29 grapheme breaking.

BSD License

Copyright (c) 2017–18, Norbert Pillmayer

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
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRETC, INDIRETC, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRATC, STRITC LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Content

UAX#29 is the Unicode Annex for breaking text into graphemes, words
and sentences.
It defines code-point classes and sets of rules
for how to place break points and break inhibitors.
This file is about grapheme breaking.

Typical Usage

Clients instantiate a grapheme object and use it as the
breaking engine for a segmenter.

  onGraphemes := grapheme.NewBreaker()
  segmenter := segment.NewSegmenter(onGraphemes)
  segmenter.Init(...)
  for segmenter.Next() ...

Attention

Before using grapheme breakers, clients will have to initialize the
classes and rules.

  SetupGraphemeClasses()

This initializes all the code-point range tables. Initialization is
not done beforehand, as it consumes quite some memory.
As grapheme breaking involves knowledge of emoji classes, a call to
SetupGraphemeClasses() will in turn call

  SetupEmojisClasses()

This UnicodeBreaker successfully passes all 672 tests for grapheme
breaking of UAX#29 (GraphemeBreakTest.txt).
*/
package grapheme

import (
	"sync"
	"unicode"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/gtcore/unicode/emoji"
	"github.com/npillmayer/gotype/gtcore/unicode/segment"
)

// We trace to the core-tracer.
var TC tracing.Trace = tracing.CoreTracer

// Top-level client function:
// Get the line grapheme class for a Unicode code-point
func GraphemeClassForRune(r rune) GraphemeClass {
	if r == rune(0) {
		return eot
	}
	for c := GraphemeClass(0); c <= ZWJClass; c++ {
		urange := rangeFromGraphemeClass[c]
		if urange == nil {
			TC.P("class", c).Errorf("no range for grapheme class")
		} else if unicode.Is(urange, r) {
			return c
		}
	}
	return Any
}

var setupOnce sync.Once

// Top-level preparation function:
// Create code-point classes for grapheme breaking.
// Will in turn set up emoji classes as well.
// (Concurrency-safe).
func SetupGraphemeClasses() {
	setupOnce.Do(setupGraphemeClasses)
	emoji.SetupEmojisClasses()
}

// === Grapheme Breaker ==============================================

// Objects of this type are used by a segment.Segmenter to break text
// up according to UAX#29 / Graphemes.
// It implements the segment.UnicodeBreaker interface.
type GraphemeBreaker struct {
	publisher    segment.RunePublisher
	longestMatch int
	penalties    []int
	rules        map[GraphemeClass][]segment.NfaStateFn
	emojirules   map[int][]segment.NfaStateFn
	blocked      map[GraphemeClass]bool
}

// Create a new UAX#14 line breaker.
//
// Usage:
//
//   onGraphemes := NewBreaker()
//   segmenter := segment.NewSegmenter(onGraphemes)
//   segmenter.Init(...)
//   for segmenter.Next() ...
//
func NewBreaker() *GraphemeBreaker {
	gb := &GraphemeBreaker{}
	gb.publisher = segment.NewRunePublisher()
	//gb.publisher.SetPenaltyAggregator(segment.MaxPenalties)
	gb.rules = map[GraphemeClass][]segment.NfaStateFn{
		//eot:                   {rule_GB2},
		CRClass:                 {rule_NewLine},
		LFClass:                 {rule_NewLine},
		ControlClass:            {rule_Control},
		LClass:                  {rule_GB6},
		VClass:                  {rule_GB7},
		LVClass:                 {rule_GB7},
		LVTClass:                {rule_GB8},
		TClass:                  {rule_GB8},
		ExtendClass:             {rule_GB9},
		ZWJClass:                {rule_GB9},
		SpacingMarkClass:        {rule_GB9a},
		PrependClass:            {rule_GB9b},
		emojiPictographic:       {rule_GB11},
		Regional_IndicatorClass: {rule_GB12},
	}
	gb.blocked = make(map[GraphemeClass]bool)
	return gb
}

// We introduce an offest for Emoji code-point classes
// to be able to tell them apart from grapheme classes.
const emojiPictographic GraphemeClass = ZWJClass + 1

// Return the grapheme code-point class for a rune (= code-point).
// (Interface segment.UnicodeBreaker)
func (gb *GraphemeBreaker) CodePointClassFor(r rune) int {
	c := GraphemeClassForRune(r)
	if c == Any {
		if unicode.Is(emoji.Extended_Pictographic, r) {
			return int(emojiPictographic)
		}
	}
	return int(c)
}

// Start all recognizers where the starting symbol is rune r.
// r is of code-point-class cpClass.
// (Interface segment.UnicodeBreaker)
//
// TODO merge this with ProceedWithRune(), it is unnecessary
func (gb *GraphemeBreaker) StartRulesFor(r rune, cpClass int) {
	c := GraphemeClass(cpClass)
	if !gb.blocked[c] {
		if rules := gb.rules[c]; len(rules) > 0 {
			TC.P("class", c).Debugf("starting %d rule(s)", c)
			for _, rule := range rules {
				rec := segment.NewPooledRecognizer(cpClass, rule)
				rec.UserData = gb
				gb.publisher.SubscribeMe(rec)
			}
		}
	}
}

// Helper: do not start any recognizers for this grapheme class, until
// unblocked again.
func (gb *GraphemeBreaker) block(c GraphemeClass) {
	gb.blocked[c] = true
}

// Helper: stop blocking new recognizers for this grapheme class.
func (gb *GraphemeBreaker) unblock(c GraphemeClass) {
	gb.blocked[c] = false
}

// A new code-point has been read and this breaker receives a message to
// consume it.
// (Interface segment.UnicodeBreaker)
func (gb *GraphemeBreaker) ProceedWithRune(r rune, cpClass int) {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("proceeding with rune %+q", r)
	gb.longestMatch, gb.penalties = gb.publisher.PublishRuneEvent(r, int(c))
	TC.P("class", c).Infof("...done with |match|=%d and %v", gb.longestMatch, gb.penalties)
	/*
		if c == Any { // rule GB999
			if len(gb.penalties) > 1 {
				gb.penalties[1] = segment.AddPenalties(gb.penalties[1], penaltyForAny)
			} else if len(gb.penalties) > 0 {
				gb.penalties = append(gb.penalties, penaltyForAny)
			} else {
				gb.penalties = penaltyForAnyAsSlice
			}
		}
	*/
}

// Collect form all active recognizers information about current match length
// and return the longest one for all still active recognizers.
// (Interface segment.UnicodeBreaker)
func (gb *GraphemeBreaker) LongestActiveMatch() int {
	return gb.longestMatch
}

// Get all active penalties for all active recognizers combined.
// Index 0 belongs to the most recently read rune, i.e., represents
// the penalty for breaking after it.
// (Interface segment.UnicodeBreaker)
func (gb *GraphemeBreaker) Penalties() []int {
	return gb.penalties
}

// --- Rules ------------------------------------------------------------

const (
	GlueBREAK int = -500
	GlueJOIN  int = 5000
	GlueBANG  int = -10000
)

// This is the break penalty for rule Any ÷ Any
const penaltyForAny = GlueBREAK

var penaltyForAnyAsSlice = []int{0, penaltyForAny}

// unnecessary ?!
/*
func rule_GB2(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	return segment.DoAccept(rec, 0, GlueBREAK)
}
*/

func rule_NewLine(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("fire rule NewLine")
	if c == LFClass {
		return segment.DoAccept(rec, GlueBANG, GlueBANG)
	} else if c == CRClass {
		rec.MatchLen++
		return rule_CRLF
	}
	return segment.DoAbort(rec)
}

func rule_CRLF(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("fire rule 05_CRLF")
	if c == LFClass {
		return segment.DoAccept(rec, GlueBANG, 3*GlueJOIN) // accept CR+LF
	}
	return segment.DoAccept(rec, 0, GlueBANG, GlueBANG) // accept CR
}

func rule_Control(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("fire rule Control")
	return segment.DoAccept(rec, GlueBANG, GlueBANG)
}

func rule_GB6(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	//c := GraphemeClass(cpClass)
	rec.MatchLen++
	return rule_GB6_L_V_LV_LVT
}

func rule_GB6_L_V_LV_LVT(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	if c == LClass || c == VClass || c == LVClass || c == LVTClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

func rule_GB7(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	//c := GraphemeClass(cpClass)
	rec.MatchLen++
	return rule_GB7_V_T
}

func rule_GB7_V_T(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	if c == VClass || c == TClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

func rule_GB8(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("start rule GB8 LVT|T x T")
	rec.MatchLen++
	return rule_GB8_T
}

func rule_GB8_T(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("accept rule GB8 T")
	if c == TClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

func rule_GB9(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("fire rule ZWJ|Extend")
	return segment.DoAccept(rec, 0, GlueJOIN)
}

func rule_GB9a(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("fire rule SpacingMark")
	return segment.DoAccept(rec, 0, GlueJOIN)
}

func rule_GB9b(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	TC.P("class", c).Infof("fire rule Preprend")
	return segment.DoAccept(rec, GlueJOIN)
}

func rule_GB11(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	TC.P("class", cpClass).Infof("fire rule Emoji Pictographic")
	return rule_GB11Cont
}

func rule_GB11Cont(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	if cpClass == int(ZWJClass) {
		rec.MatchLen++
		return rule_GB11Finish
	} else if cpClass == int(ExtendClass) {
		rec.MatchLen++
		return rule_GB11Cont
	}
	return segment.DoAbort(rec)
}

func rule_GB11Finish(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	if cpClass == int(emojiPictographic) {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

func rule_GB12(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	TC.P("class", cpClass).Infof("fire rule RI")
	gb := rec.UserData.(*GraphemeBreaker)
	gb.block(Regional_IndicatorClass)
	return rule_GB12Cont
}

func rule_GB12Cont(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	gb := rec.UserData.(*GraphemeBreaker)
	gb.unblock(Regional_IndicatorClass)
	if c == Regional_IndicatorClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}
