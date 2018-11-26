/*
Package uax14 implements Unicode Annex #14 line breaking.

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

Clients instantiate a UAX#14 line breaker object and use it as the
breaking engine for a segmenter.

  breaker := uax14.NewLineWrap()
  segmenter := unicode.NewSegmenter(breaker)
  segmenter.Init(...)
  match, length, err := segmenter.Next()

Attention

Before using line breakers, clients will have to initialize the UAX#14
classes and rules.

  SetupUAX14Classes()

This initializes all the code-point range tables. Initialization is
not done beforehand, as it consumes quite some memory, and using UAX#14
is not mandatory for line breaking.
*/
package uax14

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/npillmayer/gotype/gtcore/unicode/segment"
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
// Create code-point classes for UAX#14 line breaking/wrap.
// (Concurrency-safe).
func SetupUAX14Classes() {
	setupOnce.Do(setupUAX14Classes)
}

// === UAX#14 Line Breaker ==============================================

// Objects of this type are used by a unicode.Segmenter to break lines
// up according to UAX#14. It implements the unicode.UnicodeBreaker interface.
type UAX14LineWrap struct {
	publisher    segment.RunePublisher
	longestMatch int
	penalties    []int
	rules        map[UAX14Class][]segment.NfaStateFn
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
func NewLineWrap() *UAX14LineWrap {
	uax14 := &UAX14LineWrap{}
	uax14.publisher = segment.NewRunePublisher()
	uax14.rules = map[UAX14Class][]segment.NfaStateFn{
		NLClass: {Rule05_NewLine},
		CRClass: {Rule05_NewLine},
	}
	return uax14
}

// Initialize a line breaker for a rune-publisher (normally from a
// unicode.Segmenter).
//
// Interface unicode.UnicodeBreaker
/*
func (uax14 *UAX14LineWrap) InitFor(rpub segment.RunePublisher) {
	uax14.publisher = rpub // TODO this is wrong
}
*/

// Return the UAX#14 code-point class for a rune (= code-point).
//
// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) CodePointClassFor(r rune) int {
	return int(UAX14ClassForRune(r))
}

// Start all recognizers where the starting symbol is rune r.
// r is of code-point-class cpClass.
//
// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) StartRulesFor(r rune, cpClass int) {
	uax14c := UAX14Class(cpClass)
	rules := uax14.rules[uax14c]
	if len(rules) > 0 {
		fmt.Printf("starting rules for class = %s\n", uax14c)
	}
	for _, rule := range rules {
		rec := segment.NewPooledRecognizer(cpClass, rule)
		uax14.publisher.SubscribeMe(rec)
	}
}

// A new code-point has been read and this breaker receives a message to
// consume it.
//
// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) ProceedWithRune(r rune, cpClass int) {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("proceeding with rune = %+q / %s\n", r, uax14c)
	uax14.longestMatch, uax14.penalties = uax14.publisher.PublishRuneEvent(r, int(uax14c))
	fmt.Printf("longest match = %d\n", uax14.LongestMatch())
}

// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) LongestMatch() int {
	return uax14.longestMatch
}

// Get all active penalties for all active recognizers combined.
// Index 0 belongs to the most recently read rune.
//
// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) Penalties() []int {
	return uax14.penalties
}

// --- Standard Recognizer Rules ----------------------------------------

func doAbort(rec *segment.Recognizer) segment.NfaStateFn {
	rec.DistanceToGo = 0
	rec.MatchLen = 0
	return abort
}

func abort(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	rec.DistanceToGo = 0
	rec.MatchLen = 0
	fmt.Println("-> abort")
	return nil
}

func doAccept(rec *segment.Recognizer, penalties ...int) segment.NfaStateFn {
	fmt.Printf("-> ACCEPT %s\n", UAX14Class(rec.Expect))
	fmt.Printf("   penalties: %v", penalties)
	return nil
}

func accept(rec *segment.Recognizer, r rune, cpClass int) segment.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("-> accept %s with lookahead = %s\n", UAX14Class(rec.Expect), uax14c)
	rec.DistanceToGo = 0
	return abort
}
