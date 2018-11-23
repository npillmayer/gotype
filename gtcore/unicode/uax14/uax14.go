/*
Package for implemeting Unicode Annex #14 line breaking.

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

UAX#14 is the Unicode Annex for Line Breaking (Line Wrap).
It defines a bunch of code-point classes and a set of rules
for how to place break points / break inhibitors.

Typical Usage

Clients instantiate a line wrapper object and use it as the
breaking engine for a segmenter.

  linewrap := NewUAX14LineWrap()
  segmenter := unicode.NewSegmenter(linewrap)
  segmenter.Init(...)
  match, length, err := segmenter.Next()

*/
package uax14

import (
	"fmt"
	"log"
	"sync"
	"time"
	"unicode"

	u "github.com/npillmayer/gotype/gtcore/unicode"
)

const (
	sot       UAX14Class = 1000 // pseudo class
	eot       UAX14Class = 1001 // pseudo class
	optSpaces UAX14Class = 1002 // pseudo class
)

// Top-level client function:
// Get the line breaking/wrap class for a Unicode code-point
func UAX14ClassForRune(r rune) UAX14Class {
	if r == rune(0) {
		return eot
	}
	for lbc := UAX14Class(0); lbc <= ZWJClass; lbc++ {
		urange := rangeFromUAX14Class[lbc]
		if urange == nil {
			fmt.Printf("-- no range for class %s\n", lbc)
		} else if unicode.Is(urange, r) {
			return lbc
		}
	}
	return XXClass
}

var setupOnce sync.Once

// Top-level preparation function:
// Create 'constants' for UAX#14 line breaking/wrap.
func SetupUAX14Classes() {
	defer timeTrack(time.Now(), "setup of UAX#14 code-point ranges")
	setupOnce.Do(setupUAX14Classes)
}

// === UAX#14 Rules =====================================================

// Objects of this type are used by a unicode.Segmenter to break lines
// up according to UAX#14. It implements the unicode.UnicodeBreaker interface.
type UAX14LineWrap struct {
	publisher    u.RunePublisher
	longestMatch int
	penalties    []int
	rules        map[UAX14Class][]*u.Recognizer // TODO map of StateFn s
}

// Create a new (un-initialized) UAX#14 line breaker.
//
// Usage:
//
//   linewrap := NewUAX14LineWrap()
//   segmenter := unicode.NewSegmenter(linewrap)
//   segmenter.Init(...)
//   match, length, err := segmenter.Next()
//
func NewUAX14LineWrap() *UAX14LineWrap {
	uax14 := &UAX14LineWrap{}
	uax14.rules = map[UAX14Class][]*u.Recognizer{
		NLClass:   {u.NewRecognizer(int(NLClass), 2, []int{-1000}, UAX14NewLine)}, // TODO: create new one for every StartRulesFor(...), even better from pool
		QUClass:   {u.NewRecognizer(int(QUClass), 2, []int{1000}, UAX14Quote)},    // just keep StateFn here
		optSpaces: {u.NewRecognizer(int(SPClass), 1, nil, UAX14OptSpaces)},
	}
	return uax14
}

// Initialize a line breaker for a rune-publisher (normally from a
// unicode.Segmenter).
func (uax14 *UAX14LineWrap) InitFor(rpub u.RunePublisher) {
	uax14.publisher = rpub
}

func (uax14 *UAX14LineWrap) CodePointClassFor(r rune) int {
	return int(UAX14ClassForRune(r))
}

func (uax14 *UAX14LineWrap) StartRulesFor(r rune, cpClass int) {
	uax14c := UAX14Class(cpClass)
	rules := uax14.rules[uax14c]
	if len(rules) > 0 {
		fmt.Printf("starting rules for class = %s\n", uax14c)
	}
	for _, rule := range rules {
		uax14.publisher.SubscribeMe(rule)
	}
}

func (uax14 *UAX14LineWrap) ProceedWithRune(r rune, cpClass int) {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("proceeding with rune = %+q / %s\n", r, uax14c)
	uax14.longestMatch, uax14.penalties = uax14.publisher.PublishRuneEvent(r, int(uax14c))
	fmt.Printf("longest match = %d\n", uax14.LongestMatch())
}

func (uax14 *UAX14LineWrap) LongestMatch() int {
	return uax14.longestMatch
}

func (uax14 *UAX14LineWrap) Penalties() []int {
	return uax14.penalties
}

// --- Recognizer Rules -------------------------------------------------

func doAbort(rec *u.Recognizer) u.NfaStateFn {
	rec.DistanceToGo = 0
	rec.MatchLen = 0
	return abort
}

func abort(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	rec.DistanceToGo = 0
	rec.MatchLen = 0
	fmt.Println("-> abort")
	return abort
}

func accept(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("-> accept %s with lookahead = %s\n", UAX14Class(rec.Expect), uax14c)
	rec.DistanceToGo = 0
	return abort
}

func UAX14NewLine(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("fire rule NewLine for lookahead = %s\n", uax14c)
	if uax14c == NLClass || uax14c == LFClass {
		rec.DistanceToGo = 1
		rec.MatchLen++
		return accept
	}
	return doAbort(rec)
}

func UAX14Quote(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("fire rule Quote for lookahead = %s\n", uax14c)
	if uax14c == QUClass {
		rec.DistanceToGo = 2
		rec.MatchLen++
		rec.Expect = int(OPClass)
		return UAX14OptSpaces
	}
	return doAbort(rec)
}

func UAX14ThisClass(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("fire generic rule for %s with lookahead = %s\n", UAX14Class(rec.Expect), uax14c)
	if uax14c == UAX14Class(rec.Expect) {
		rec.DistanceToGo--
		rec.MatchLen++
		return accept
	}
	return doAbort(rec)
}

func UAX14OptSpaces(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	if uax14c == SPClass {
		fmt.Println("ignoring optional space")
		rec.MatchLen++
		return UAX14OptSpaces // repeat
	} else if uax14c == UAX14Class(rec.Expect) {
		rec.MatchLen++
		return accept // accept rec.Expect in next rec
	}
	return doAbort(rec)
}

// --- Util -------------------------------------------------------------

// Little helper for testing
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("timing: %s took %s\n", name, elapsed)
}
