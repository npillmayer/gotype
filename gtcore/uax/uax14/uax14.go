/*
Package uax14 implements Unicode Annex #14 line breaking.

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


Contents

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

Before using line breakers, clients usually will want to initialize the UAX#14
classes and rules.

  SetupUAX14Classes()

This initializes all the code-point range tables. Initialization is
not done beforehand, as it consumes quite some memory, and using UAX#14
is not mandatory for line breaking.
*/
package uax14

import (
	"fmt"
	"math"
	"sync"
	"unicode"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/gtcore/uax"
)

const (
	sot       UAX14Class = 1000 // pseudo class
	eot       UAX14Class = 1001 // pseudo class
	optSpaces UAX14Class = 1002 // pseudo class
)

// We trace to the core-tracer.
var TC tracing.Trace = tracing.CoreTracer

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
	publisher    uax.RunePublisher
	longestMatch int   // longest active match of a rule
	penalties    []int // returned to the segmenter: penalties to insert
	rules        map[UAX14Class][]uax.NfaStateFn
	lastClass    UAX14Class // we have to remember the last code-point class
	blockedRI    bool       // are rules for Regional_Indicator currently blocked?
}

// Create a new UAX#14 line breaker.
//
// Usage:
//
//   linewrap := NewUAX14LineWrap()
//   segmenter := segment.NewSegmenter(linewrap)
//   segmenter.Init(...)
//   for segmenter.Next() ...
//
func NewLineWrap() *UAX14LineWrap {
	uax14 := &UAX14LineWrap{}
	uax14.publisher = uax.NewRunePublisher()
	uax14.rules = map[UAX14Class][]uax.NfaStateFn{
		NLClass: {rule_05_NewLine},
		LFClass: {rule_05_NewLine},
		BKClass: {rule_05_NewLine},
		CRClass: {rule_05_NewLine},
		SPClass: {rule_LB7, rule_LB18},
		ZWClass: {rule_LB7, rule_LB8},
		WJClass: {rule_LB11},
		GLClass: {rule_LB12},
		CLClass: {rule_LB13, rule_LB16},
		CPClass: {rule_LB13, rule_LB16, rule_LB30_2},
		EXClass: {rule_LB13, rule_LB22},
		ISClass: {rule_LB13, rule_LB29},
		SYClass: {rule_LB13, rule_LB21b},
		OPClass: {rule_LB14},
		QUClass: {rule_LB15, rule_LB19},
		B2Class: {rule_LB17},
		BAClass: {rule_LB21},
		HYClass: {rule_LB21},
		NSClass: {rule_LB21},
		BBClass: {rule_LB21x},
		ALClass: {rule_LB22, rule_LB23_1, rule_LB24_2, rule_LB28, rule_LB30_1},
		HLClass: {rule_LB21a, rule_LB22, rule_LB23_1, rule_LB24_2, rule_LB28, rule_LB30_1},
		IDClass: {rule_LB22},
		EBClass: {rule_LB22, rule_LB23a_2, rule_LB30b},
		EMClass: {rule_LB22, rule_LB23a_2},
		INClass: {rule_LB22},
		NUClass: {rule_LB22, rule_LB23_2, rule_LB30_1},
		RIClass: {rule_LB30a},
		PRClass: {rule_LB23a_2},
	}
	if rangeFromUAX14Class == nil {
		TC.Info("UAX#14 classes not yet initialized -> initializing")
	}
	SetupUAX14Classes()
	uax14.lastClass = sot
	return uax14
}

// Return the UAX#14 code-point class for a rune (= code-point).
//
// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) CodePointClassFor(r rune) int {
	c := UAX14ClassForRune(r)
	c = substitueSomeClasses(c, uax14.lastClass)
	return int(c)
}

// Start all recognizers where the starting symbol is rune r.
// r is of code-point-class cpClass.
//
// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) StartRulesFor(r rune, cpClass int) {
	c := UAX14Class(cpClass)
	if c != RIClass || !uax14.blockedRI {
		if rules := uax14.rules[c]; len(rules) > 0 {
			TC.P("class", c).Debugf("starting %d rule(s) for class %s", len(rules), c)
			for _, rule := range rules {
				rec := uax.NewPooledRecognizer(cpClass, rule)
				rec.UserData = uax14
				uax14.publisher.SubscribeMe(rec)
			}
		} else {
			TC.P("class", c).Debugf("starting no rule")
		}
	}
}

// LB9: Do not break a combining character sequence;
// treat it as if it has the line breaking class of the base character in all
// of the following rules. Treat ZWJ as if it were CM.
//
//    X (CM | ZWJ)* ⟼ X.
//
// where X is any line break class except BK, CR, LF, NL, SP, or ZW.
//
// LB10: Treat any remaining combining mark or ZWJ as AL.
func substitueSomeClasses(c UAX14Class, lastClass UAX14Class) UAX14Class {
	orig := c
	switch lastClass {
	case sot, BKClass, CRClass, LFClass, NLClass, SPClass, ZWClass:
		if c == CMClass || c == ZWJClass {
			c = ALClass
		}
	default:
		if c == CMClass || c == ZWJClass {
			c = lastClass
		}
	}
	if orig != c {
		TC.Debugf("subst %+q -> %+q", orig, c)
	}
	return c
}

// A new code-point has been read and this breaker receives a message to
// consume it.
//
// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) ProceedWithRune(r rune, cpClass int) {
	c := UAX14Class(cpClass)
	uax14.longestMatch, uax14.penalties = uax14.publisher.PublishRuneEvent(r, int(c))
	x := uax14.penalties
	for i, p := range uax14.penalties {
		if p > 0 {
			p += 1000
			x[i] = p
		}
	}
	uax14.lastClass = c
}

// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) LongestActiveMatch() int {
	return uax14.longestMatch
}

// Get all active penalties for all active recognizers combined.
// Index 0 belongs to the most recently read rune.
//
// Interface unicode.UnicodeBreaker
func (uax14 *UAX14LineWrap) Penalties() []int {
	return uax14.penalties
}

// Helper: do not start any recognizers for class RI, until
// unblocked again.
func (uax14 *UAX14LineWrap) block() {
	uax14.blockedRI = true
}

// Helper: stop blocking new recognizers for class RI.
func (uax14 *UAX14LineWrap) unblock() {
	uax14.blockedRI = false
}

// Penalties (optional break, suppress break and mandatory break).
var (
	PenaltyForBreak        int = 50
	PenaltyToSuppressBreak int = 5000
	PenaltyForMustBreak    int = -10000
)

// TODO temporaty penalty function
func p(w int) int {
	q := 31 - w
	r := int(math.Pow(1.3, float64(q)))
	TC.P("rule", w).Infof("penalty %d => %d", w, r)
	return r
}

func ps(w int, first int, l int) []int {
	pp := make([]int, l)
	pp[0] = first
	for i := 1; i < l; i++ {
		pp[i] = p(w)
	}
	return pp
}
