/*
Package uax29 implements Unicode Annex #29 word breaking.

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
This file is about word breaking.

Typical Usage

Clients instantiate a WordBreaker object and use it as the
breaking engine for a segmenter.

  onWords := uax29.NewWordBreaker()
  segmenter := segment.NewSegmenter(onWords)
  segmenter.Init(...)
  for segmenter.Next() ...

Attention

Before using word breakers, clients will have to initialize the
classes and rules.

  SetupUAX29Classes()

This initializes all the code-point range tables. Initialization is
not done beforehand, as it consumes quite some memory.
*/
package uax29

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
// Get the line word class for a Unicode code-point
func UAX29ClassForRune(r rune) UAX29Class {
	if r == rune(0) {
		return eot
	}
	for c := UAX29Class(0); c <= ZWJClass; c++ {
		urange := rangeFromUAX29Class[c]
		if urange == nil {
			TC.P("class", c).Errorf("no range for UAX29 class")
		} else if unicode.Is(urange, r) {
			return c
		}
	}
	return Other
}

var setupOnce sync.Once

// Top-level preparation function:
// Create code-point classes for word breaking.
// Will in turn set up emoji classes as well.
// (Concurrency-safe).
func SetupUAX29Classes() {
	setupOnce.Do(setupUAX29Classes)
	emoji.SetupEmojisClasses()
}

// === Word Breaker ==============================================

// Objects of this type are used by a segment.Segmenter to break text
// up according to UAX#29 / Words.
// It implements the segment.UnicodeBreaker interface.
type WordBreaker struct {
	publisher    segment.RunePublisher
	longestMatch int
	penalties    []int
	rules        map[UAX29Class][]segment.NfaStateFn
	emojirules   map[int][]segment.NfaStateFn
	blocked      map[UAX29Class]bool
}

// Create a new UAX#14 line breaker.
//
// Usage:
//
//   onWords := NewBreaker()
//   segmenter := segment.NewSegmenter(onWords)
//   segmenter.Init(...)
//   for segmenter.Next() ...
//
func NewBreaker() *WordBreaker {
	gb := &WordBreaker{}
	gb.publisher = segment.NewRunePublisher()
	//gb.publisher.SetPenaltyAggregator(segment.MaxPenalties)
	gb.rules = map[UAX29Class][]segment.NfaStateFn{
		//eot:                   {rule_GB2},
		CRClass:                 {rule_NewLine},
		LFClass:                 {rule_NewLine},
		NewlineClass:            {rule_NewLine},
		ZWJClass:                {rule_WB3c},
		WSegSpaceClass:          {rule_WB3d},
		ALetterClass:            {rule_WB5, rule_WB6_7, rule_WB13a},
		Hebrew_LetterClass:      {rule_WB5, rule_WB6_7, rule_WB7a, rule_WB7bc, rule_WB13a},
		NumericClass:            {rule_WB8, rule_WB10, rule_WB11, rule_WB13a},
		ExtendNumLetClass:       {rule_WB13a},
		KatakanaClass:           {rule_WB13, rule_WB13a},
		Regional_IndicatorClass: {rule_WB15},
	}
	gb.blocked = make(map[UAX29Class]bool)
	return gb
}

const emojiPictographic UAX29Class = ZWJClass + 1

// Return the word code-point class for a rune (= code-point).
// (Interface segment.UnicodeBreaker)
func (gb *WordBreaker) CodePointClassFor(r rune) int {
	c := UAX29ClassForRune(r)
	if c == Other {
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
func (gb *WordBreaker) StartRulesFor(r rune, cpClass int) {
	c := UAX29Class(cpClass)
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

// Helper: do not start any recognizers for this word class, until
// unblocked again.
func (gb *WordBreaker) block(c UAX29Class) {
	gb.blocked[c] = true
}

// Helper: stop blocking new recognizers for this word class.
func (gb *WordBreaker) unblock(c UAX29Class) {
	gb.blocked[c] = false
}

// A new code-point has been read and this breaker receives a message to
// consume it.
// (Interface segment.UnicodeBreaker)
func (gb *WordBreaker) ProceedWithRune(r rune, cpClass int) {
	c := UAX29Class(cpClass)
	TC.P("class", c).Infof("proceeding with rune %+q", r)
	gb.longestMatch, gb.penalties = gb.publisher.PublishRuneEvent(r, int(c))
	TC.P("class", c).Infof("...done with |match|=%d and %v", gb.longestMatch, gb.penalties)
}

// Collect form all active recognizers information about current match length
// and return the longest one for all still active recognizers.
// (Interface segment.UnicodeBreaker)
func (gb *WordBreaker) LongestActiveMatch() int {
	return gb.longestMatch
}

// Get all active penalties for all active recognizers combined.
// Index 0 belongs to the most recently read rune, i.e., represents
// the penalty for breaking after it.
// (Interface segment.UnicodeBreaker)
func (gb *WordBreaker) Penalties() []int {
	return gb.penalties
}

// --- Rules ------------------------------------------------------------

const (
	GlueBREAK int = -500
	GlueJOIN  int = 5000
	GlueBANG  int = -10000
)

func rule_NewLine(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if c == LFClass || c == NewlineClass {
		return segment.DoAccept(rec, GlueBANG, GlueBANG)
	} else if c == CRClass {
		rec.MatchLen++
		return rule_CRLF
	}
	return segment.DoAbort(rec)
}

func rule_CRLF(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if c == LFClass {
		return segment.DoAccept(rec, GlueBANG, 3*GlueJOIN) // accept CR+LF
	}
	return segment.DoAccept(rec, 0, GlueBANG, GlueBANG) // accept CR
}

func rule_WB3c(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return rule_Pictography
}

func rule_Pictography(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if c == emojiPictographic {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

func rule_WB3d(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB3dCont
}

func WB3dCont(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if c == WSegSpaceClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

func checkIgnoredCharacters(rec *segment.Recognizer, c UAX29Class) bool {
	if c == ExtendClass || c == FormatClass || c == ZWJClass {
		rec.MatchLen++
		return true
	}
	return false
}

// start AHLetter x AHLetter
func rule_WB5(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB5_10Cont
}

// ... x AHLetter
func WB5_10Cont(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB5_10Cont
	}
	if c == ALetterClass || c == Hebrew_LetterClass {
		rec.MatchLen++
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

// AHLetter x (MidLetter | MidNumLet | Single_Quote) x AHLetter
func rule_WB6_7(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB6_7Cont1
}

func WB6_7Cont1(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB6_7Cont1
	}
	if c == MidLetterClass || c == MidNumLetClass || c == Single_QuoteClass {
		rec.MatchLen++
		return WB6_7Cont2
	}
	return segment.DoAbort(rec)
}

func WB6_7Cont2(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB6_7Cont2
	}
	if c == ALetterClass || c == Hebrew_LetterClass {
		return segment.DoAccept(rec, 0, GlueJOIN, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

// Hebrew_Letter x Single_Quote
func rule_WB7a(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB7aCont
}

func WB7aCont(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB7aCont
	}
	if c == Hebrew_LetterClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

// Hebrew_Letter x Double_Quote x Hebrew_Letter
func rule_WB7bc(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB7bcCont1
}

func WB7bcCont1(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB7bcCont1
	}
	if c == Double_QuoteClass {
		rec.MatchLen++
		return WB7bcCont2
	}
	return segment.DoAbort(rec)
}

func WB7bcCont2(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB7bcCont2
	}
	if c == Hebrew_LetterClass {
		return segment.DoAccept(rec, 0, GlueJOIN, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

// start Numeric x Numeric
func rule_WB8(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB89Cont
}

// ... x Numeric
func WB89Cont(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB89Cont
	}
	if c == NumericClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

// start (ALetter | Hebrew_Letter) x Numeric
func rule_WB9(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB89Cont
}

// start Numeric x AHLetter
func rule_WB10(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB5_10Cont
}

// start Numeric x (MidNum | MidNumLet | Single_Quote) x Numeric
func rule_WB11(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB11Cont1
}

// ... x (MidNum | MidNumLet | Single_Quote) x Numeric
func WB11Cont1(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB11Cont1
	}
	if c == MidNumClass || c == MidNumLetClass || c == Single_QuoteClass {
		rec.MatchLen++
		return WB11Cont2
	}
	return segment.DoAbort(rec)
}

// ... x ... x Numeric
func WB11Cont2(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB11Cont2
	}
	if c == NumericClass {
		return segment.DoAccept(rec, 0, GlueJOIN, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

// start Katakana x Katakana
func rule_WB13(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB13Cont
}

// ... x Katakana
func WB13Cont(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB89Cont
	}
	if c == KatakanaClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

// start (AHLetter | Numeric | Katakana | ExtendNumLet) xExtendNumLet x (AHLetter | Numeric | Katakana)
func rule_WB13a(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.MatchLen++
	return WB13aCont1
}

// ... x (MidNum | MidNumLet | Single_Quote) x Numeric
func WB13aCont1(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB13aCont1
	}
	if c == ExtendNumLetClass {
		rec.MatchLen++
		return WB13aCont2
	}
	return segment.DoAbort(rec)
}

// ... x ... x Numeric
func WB13aCont2(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	if checkIgnoredCharacters(rec, c) {
		return WB13aCont2
	}
	if c == ALetterClass || c == Hebrew_LetterClass || c == NumericClass || c == KatakanaClass {
		return segment.DoAccept(rec, 0, GlueJOIN, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

// start RI x RI (blocking)
func rule_WB15(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	gb := rec.UserData.(*WordBreaker)
	gb.block(Regional_IndicatorClass)
	return WB15Cont
}

// ... x RI
func WB15Cont(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := UAX29Class(cpClass)
	gb := rec.UserData.(*WordBreaker)
	gb.unblock(Regional_IndicatorClass)
	if checkIgnoredCharacters(rec, c) {
		return WB15Cont
	}
	if c == Regional_IndicatorClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}
