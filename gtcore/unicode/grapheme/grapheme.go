/*
Package grapheme implements Unicode Annex #29 grapheme breaking.

TODO
// https://github.com/nitely/nim-graphemes/tree/master/gen
// Berücksichtigen von emoji-data.txt

---------------------------------------------------------------------------

BSD License

Copyright (c) 2017-18, Norbert Pillmayer

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
  grapheme, length, err := segmenter.Next()

Attention

Before using grapheme breakers, clients will have to initialize the
classes and rules.

  SetupGraphemeClasses()

This initializes all the code-point range tables. Initialization is
not done beforehand, as it consumes quite some memory.
*/
package grapheme

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/npillmayer/gotype/gtcore/unicode/segment"
)

// Top-level client function:
// Get the line grapheme class for a Unicode code-point
func GraphemeClassForRune(r rune) GraphemeClass {
	if r == rune(0) {
		return eot
	}
	for c := GraphemeClass(0); c <= ZWJClass; c++ {
		urange := rangeFromGraphemeClass[c]
		if urange == nil {
			fmt.Printf("-- no range for grapheme class %s\n", c)
		} else if unicode.Is(urange, r) {
			return c
		}
	}
	return Any
}

var setupOnce sync.Once

// Top-level preparation function:
// Create code-point classes for grapheme breaking.
// (Concurrency-safe).
func SetupGraphemeClasses() {
	setupOnce.Do(setupGraphemeClasses)
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
}

// Create a new (un-initialized) UAX#14 line breaker.
//
// Usage:
//
//   onGraphemes := NewBreaker()
//   segmenter := segment.NewSegmenter(onGraphemes)
//   segmenter.Init(...)
//   match, length, err := segmenter.Next() // TODO: 2nd is penalty
//
func NewBreaker() *GraphemeBreaker {
	gb := &GraphemeBreaker{}
	gb.publisher = segment.NewRunePublisher()
	gb.publisher.SetPenaltyAggregator(segment.MaxPenalties)
	gb.rules = map[GraphemeClass][]segment.NfaStateFn{
		//eot:              {rule_GB2},
		CRClass:          {rule_NewLine},
		LFClass:          {rule_NewLine},
		ControlClass:     {rule_Control},
		LClass:           {rule_GB6},
		VClass:           {rule_GB7},
		LVClass:          {rule_GB7},
		LVTClass:         {rule_GB8},
		TClass:           {rule_GB8},
		ExtendClass:      {rule_GB9},
		ZWJClass:         {rule_GB9},
		SpacingMarkClass: {rule_GB9a},
		PrependClass:     {rule_GB9b},
	}
	return gb
}

// Return the grapheme code-point class for a rune (= code-point).
//
// Interface segment.UnicodeBreaker
func (gb *GraphemeBreaker) CodePointClassFor(r rune) int {
	return int(GraphemeClassForRune(r))
}

// Start all recognizers where the starting symbol is rune r.
// r is of code-point-class cpClass.
//
// Interface segment.UnicodeBreaker
//
// TODO merge this with ProceedWithRune(), it is unnecessary
func (gb *GraphemeBreaker) StartRulesFor(r rune, cpClass int) {
	c := GraphemeClass(cpClass)
	rules := gb.rules[c]
	if len(rules) > 0 {
		fmt.Printf("starting rules for class = %s\n", c)
	}
	for _, rule := range rules {
		rec := segment.NewPooledRecognizer(cpClass, rule)
		gb.publisher.SubscribeMe(rec)
	}
}

// A new code-point has been read and this breaker receives a message to
// consume it.
//
// Interface segment.UnicodeBreaker
func (gb *GraphemeBreaker) ProceedWithRune(r rune, cpClass int) {
	gbc := GraphemeClass(cpClass)
	fmt.Printf("proceeding with rune = %+q / %s\n", r, gbc)
	gb.longestMatch, gb.penalties = gb.publisher.PublishRuneEvent(r, int(gbc))
	fmt.Printf("longest match = %d\n", gb.LongestMatch())
}

// Interface segment.UnicodeBreaker
func (gb *GraphemeBreaker) LongestMatch() int {
	return gb.longestMatch
}

// Get all active penalties for all active recognizers combined.
// Index 0 belongs to the most recently read rune, i.e., represents
// the penalty for breaking after it.
//
// Interface segment.UnicodeBreaker
func (gb *GraphemeBreaker) Penalties() []int {
	return gb.penalties
}

// --- Rules ------------------------------------------------------------

const (
	GlueBREAK int = -500
	GlueJOIN  int = 10000
)

// unnecessary ?!
/*
func rule_GB2(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	return segment.DoAccept(rec, 0, GlueBREAK)
}
*/

func rule_NewLine(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	fmt.Printf("fire rule NewLine for lookahead = %s\n", c)
	if c == LFClass {
		return segment.DoAccept(rec, GlueBREAK, GlueBREAK)
	} else if c == CRClass {
		rec.MatchLen++
		return rule_CRLF
	}
	return segment.DoAbort(rec)
}

func rule_CRLF(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	fmt.Printf("fire rule 05_CRLF for lookahead = %s\n", c)
	if c == LFClass {
		return segment.DoAccept(rec, GlueBREAK, GlueBREAK) // accept CR+LF
	}
	return segment.DoAccept(rec, 0, GlueBREAK, GlueBREAK) // accept CR
}

func rule_Control(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	fmt.Printf("fire rule Control for lookahead = %s\n", c)
	return segment.DoAccept(rec, GlueBREAK, GlueBREAK)
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
	fmt.Printf("start rule GB8 LVT|T x T for lookahead = %s\n", c)
	rec.MatchLen++
	return rule_GB8_T
}

func rule_GB8_T(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	fmt.Printf("accept rule GB8 T for lookahead = %s?\n", c)
	if c == TClass {
		return segment.DoAccept(rec, 0, GlueJOIN)
	}
	return segment.DoAbort(rec)
}

func rule_GB9(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	fmt.Printf("fire rule ZWJ|Extend for lookahead = %s\n", c)
	return segment.DoAccept(rec, 0, GlueJOIN)
}

func rule_GB9a(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	fmt.Printf("fire rule SpacingMark for lookahead = %s\n", c)
	return segment.DoAccept(rec, 0, GlueJOIN)
}

func rule_GB9b(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	c := GraphemeClass(cpClass)
	fmt.Printf("fire rule Preprend for lookahead = %s\n", c)
	return segment.DoAccept(rec, GlueJOIN)
}
