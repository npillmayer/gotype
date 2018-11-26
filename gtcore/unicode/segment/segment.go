/*
Package segment is about Unicode and text segmenting.

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

Typical Usage

Clients instantiate a UnicodeBreaker object and use it as the
breaking engine for a segmenter. Multiple breaking engines may be
supplied.

  breaker1 := ...
  breaker2 := ...
  segmenter := unicode.NewSegmenter(breaker1, breaker2)
  segmenter.Init(...)
  match, length, err := segmenter.Next()

An example for an UnicodeBreaker is "uax14.LineWrap", a breaker
implementing the UAX#14 line breaking algorithm.

This package provides a variety of types to support authors of
UnicodeBreakers: RunePublisher/RuneSubscriber and Recognizer.
*/
package segment

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gammazero/deque"
	pool "github.com/jolestar/go-commons-pool"
)

// Object of type UnicodeBreaker represent logic components to split up
// Unicode sequences into smaller parts. They are used by Segmenters
// to supply breaking logic.
type UnicodeBreaker interface {
	CodePointClassFor(rune) int
	StartRulesFor(rune, int)
	ProceedWithRune(rune, int)
	LongestMatch() int
	Penalties() []int
}

type atom struct {
	r       rune
	length  int
	penalty int
}

var eotAtom atom = atom{rune(0), 0, 0}

func (a *atom) String() string {
	return fmt.Sprintf("[%+q p=%d]", a.r, a.penalty)
}

// A Segmenter receives a sequence of code-points from an io.RuneReader and
// segments it into smaller parts.
type Segmenter struct {
	deque *deque.Deque
	//publisher             RunePublisher
	breakers              []UnicodeBreaker
	reader                io.RuneReader
	nullPenalty           int
	longestActiveMatch    int
	distanceToNextPenalty int
	atEOF                 bool
}

// Create a new Segmenter by providing breaking logic (UnicodeBreaker).
func NewSegmenter(breakers ...UnicodeBreaker) *Segmenter {
	s := &Segmenter{}
	s.breakers = breakers
	return s
}

// Initialize a Segmenter with an io.RuneReader to read from.
func (s *Segmenter) Init(reader io.RuneReader) {
	if reader == nil {
		reader = strings.NewReader("")
	}
	s.reader = reader
	if s.deque == nil {
		s.deque = &deque.Deque{} // Q of atoms
		//s.publisher = NewRunePublisher() // for publishing rune events to breakers
		//for _, breaker := range s.breakers {
		//breaker.InitFor(s.publisher)
		//}
	} else {
		s.deque.Clear()
		s.longestActiveMatch = 0
		s.atEOF = false
	}
	s.distanceToNextPenalty = 10000
}

func (s *Segmenter) SetNullPenalty(null int) {
	s.nullPenalty = null
}

// Get the next segment, together with the accumulated penalty for this break.
func (s *Segmenter) Next() ([]byte, int, error) {
	if s.reader == nil {
		return nil, 0, errors.New("segmenter not initialized: no input; must call Init()")
	}
	var match []byte
	l := 0
	err := s.readEnoughInput()
	if err != nil && err != io.EOF {
		return nil, 0, err
	}
	if s.distanceToNextPenalty < 1000 {
		fmt.Printf("distance to next penalty (%d) = %d\n", s.deque.At(s.distanceToNextPenalty-1), s.distanceToNextPenalty)
	}
	match, l = s.getTrailSegment()
	return match, l, err
}

func (s *Segmenter) frontAtom() *atom {
	return s.deque.Front().(*atom)
}

func (s *Segmenter) trailAtom() *atom {
	return s.deque.Back().(*atom)
}

func (s *Segmenter) readRune() (err error) {
	fmt.Println("-------- reading next rune -----------")
	if s.atEOF {
		err = io.EOF
	} else {
		var r rune
		r, _, err = s.reader.ReadRune()
		fmt.Printf("rune = %+q\n", r)
		if err == nil {
			a := &atom{} // TODO get from pool
			a.r = r
			a.length = 1
			s.deque.PushFront(a)
		} else if err == io.EOF {
			s.deque.PushFront(&eotAtom)
			s.atEOF = true
			err = nil
		} else {
			fmt.Printf("ReadRune() returned error = %s\n", err)
			s.atEOF = true
		}
	}
	return err
}

func (s *Segmenter) readEnoughInput() (err error) {
	//for i := 0; s.deque.Len()-s.longestActiveMatch <= 0; {
	for s.distanceToNextPenalty+s.longestActiveMatch >= s.deque.Len() {
		err = s.readRune()
		if err != nil {
			break
		}
		if s.deque.Len() > 0 {
			s.longestActiveMatch = 0
			a := s.frontAtom()
			for _, breaker := range s.breakers {
				cpClass := breaker.CodePointClassFor(a.r)
				breaker.StartRulesFor(a.r, cpClass)
				breaker.ProceedWithRune(a.r, cpClass)
				if breaker.LongestMatch() > s.longestActiveMatch {
					s.longestActiveMatch = breaker.LongestMatch()
				}
				s.insertPenalties(breaker.Penalties())
			}
			s.printQ()
		} else {
			fmt.Println("Q empty")
		}
	}
	return err
}

func (s *Segmenter) insertPenalties(penalties []int) {
	l := s.deque.Len()
	for i, p := range penalties {
		atom := s.deque.At(i).(*atom)
		atom.penalty += p
		fmt.Printf("l = %d, i = %d\n", l, i)
		if atom.penalty != s.nullPenalty && l-(i+1) < s.distanceToNextPenalty {
			s.distanceToNextPenalty = l - (i + 1)
			fmt.Printf("=> distance to penalty = %d\n", s.distanceToNextPenalty)
		}
	}
}

func (s *Segmenter) getTrailSegment() ([]byte, int) {
	fmt.Println("Collecting trail match:")
	buf := new(bytes.Buffer)
	l := s.deque.Len() - 1 // index of last
	var last *atom
	seglen := 0
	for i := 0; i < s.distanceToNextPenalty; i++ {
		fmt.Printf(".at[%d] = %s\n", l-i, s.deque.At(l-i))
		last = s.deque.PopBack().(*atom)
		written, _ := buf.WriteRune(last.r)
		seglen += written
	}
	s.distanceToNextPenalty = 10000
	for i := 0; i < s.deque.Len(); i++ {
		atom := s.deque.At(i).(*atom)
		if atom.penalty != s.nullPenalty {
			s.distanceToNextPenalty = s.deque.Len() - i // TODO may be off by 1
		}
	}
	fmt.Printf("new distance to penalty = %d\n", s.distanceToNextPenalty)
	return buf.Bytes(), seglen
}

func (s *Segmenter) printQ() {
	fmt.Printf("Q #%d: ", s.deque.Len())
	for i := 0; i < s.deque.Len(); i++ {
		fmt.Printf(" - %s", s.deque.At(i))
	}
	fmt.Println(" .")
}

// ----------------------------------------------------------------------

type NfaStateFn func(*Recognizer, rune, int) NfaStateFn

type Recognizer struct {
	Expect       int
	DistanceToGo int
	MatchLen     int
	Penalties    []int
	nextStep     NfaStateFn
}

// Create a new Recognizer. This is rarely used, as clients rather should
// call NewPooledRecognizer().
func NewRecognizer(codePointClass int, distance int, p []int, next NfaStateFn) *Recognizer {
	rec := &Recognizer{}
	rec.Expect = codePointClass
	rec.DistanceToGo = distance
	rec.Penalties = p
	rec.nextStep = next
	return rec
}

type recognizerPool struct {
	opool *pool.ObjectPool
	ctx   context.Context
}

var globalRecognizerPool *recognizerPool

func init() {
	globalRecognizerPool = &recognizerPool{}
	factory := pool.NewPooledObjectFactorySimple(
		func(context.Context) (interface{}, error) {
			rec := &Recognizer{}
			rec.Penalties = make([]int, 5)
			return rec, nil
		})
	globalRecognizerPool.ctx = context.Background()
	config := pool.NewDefaultPoolConfig()
	config.LIFO = false
	config.MaxTotal = -1 // infinity
	config.BlockWhenExhausted = false
	globalRecognizerPool.opool = pool.NewObjectPool(globalRecognizerPool.ctx, factory, config)
}

// Returns a new Recognizer, pre-filled with an expected code-point class
// and a state function. The Recognizer is pooled for efficiency.
func NewPooledRecognizer(cpClass int, stateFn NfaStateFn) *Recognizer {
	o, _ := globalRecognizerPool.opool.BorrowObject(globalRecognizerPool.ctx)
	rec := o.(*Recognizer)
	rec.Expect = cpClass
	rec.nextStep = stateFn
	return rec
}

// Clears the Recognizer and puts it back into the pool.
func (rec *Recognizer) releaseIntoPool() {
	if rec.Penalties != nil {
		rec.Penalties = rec.Penalties[:0]
	}
	rec.Expect = 0
	rec.DistanceToGo = 0
	rec.MatchLen = 0
	rec.nextStep = nil
	_ = globalRecognizerPool.opool.ReturnObject(globalRecognizerPool.ctx, rec)
}

func (rec *Recognizer) String() string {
	if rec == nil {
		return "[nil rule]"
	}
	return fmt.Sprintf("[%d -> %d steps]", rec.Expect, rec.DistanceToGo)
}

// Interface RuneSubscriber
func (rec *Recognizer) Unsubscribed() {
	rec.releaseIntoPool()
}

// Interface RuneSubscriber
func (rec *Recognizer) Done() bool {
	return rec.nextStep == nil
}

// Interface RuneSubscriber
func (rec *Recognizer) MatchLength() int {
	return rec.MatchLen
}

// Interface RuneSubscriber
func (rec *Recognizer) RuneEvent(r rune, codePointClass int) []int {
	fmt.Printf("received rune event: %+q / %d\n", r, codePointClass)
	var penalties []int
	if rec.nextStep == nil {
		rec.DistanceToGo = 0
	} else {
		rec.nextStep = rec.nextStep(rec, r, codePointClass)
		if rec.MatchLen > 0 && rec.DistanceToGo == 0 { // accepted a match
			penalties = rec.Penalties
		}
	}
	return penalties
}

// --- Rune Publishing and Subscription ---------------------------------

// RuneSubscribers are receivers of rune events, i.e. messages to
// process a new code-point (rune). If they can match the rune, they
// will expect further runes, otherwise they abort. To abort, they
// set Done() to true.
type RuneSubscriber interface {
	RuneEvent(r rune, codePointClass int) []int // receive a new code-point
	MatchLength() int                           // length (in # of code-point) of the match up to now
	Done() bool                                 // is this subscriber done?
	Unsubscribed()                              // this subscriber has been unsubscribed
}

// RunePublishers notify subscribers with rune events: a new rune has been read
// and the subscriber – usually a recognizer rule – has to react to it.
//
// UnicodeBreakers are not required to use the RunePublisher/RuneSubscriber
// pattern, but it is convenient to stick to it. UnicodeBreakers often
// rely on sets of rules, which are tested interleavingly. To releave
// UnicodeBreakers from managing rune-distribution to all the rules, it
// may be advantageous hold a RunePublisher within a UnicodeBreaker and
// let all rules implement the RuneSubscriber interface.
type RunePublisher interface {
	SubscribeMe(RuneSubscriber) RunePublisher
	PublishRuneEvent(r rune, codePointClass int) (longestDistance int, penalties []int)
}

// Create a new default RunePublisher.
func NewRunePublisher() *DefaultRunePublisher { return &DefaultRunePublisher{} }

// Trigger a rune event notification to all subscribers. Rune events
// include the rune (code-point) and an optional code-point class for
// the rune.
//
// Return values are: the longest active match and a slice of penalties.
// These values usually are collected from the RuneSubscribers.
//
// Interface RunePublisher
func (rpub DefaultRunePublisher) PublishRuneEvent(r rune, codePointClass int) (int, []int) {
	longest := 0
	penaltiesTotal := make([]int, 0, 2)
	// pre-condition: no subscriber is Done()
	for i := rpub.Len() - 1; i >= 0; i++ {
		subscr := rpub.at(i)
		penalties := subscr.RuneEvent(r, codePointClass)
		penaltiesTotal = AddPenalties(penaltiesTotal, penalties)
		if d := subscr.MatchLength(); d > longest {
			longest = d
		}
		rpub.Fix(i) // re-order heap if subscr.Done()
	}
	// no unsubscribe all done subscribers
	for subscr := rpub.PopDone(); subscr != nil; {
		subscr.Unsubscribed()
	}
	return longest, penaltiesTotal
}

// The default function to aggregate break-penalties.
// Simply adds up all penalties.
func AddPenalties(total []int, penalties []int) []int {
	for i, p := range penalties {
		if i >= len(total) {
			total = append(total, p)
		} else {
			total[i] += p
		}
	}
	return total
}

// Interface RunePublisher
func (rpub *DefaultRunePublisher) SubscribeMe(rsub RuneSubscriber) RunePublisher {
	rpub.Push(rsub)
	return rpub
}
