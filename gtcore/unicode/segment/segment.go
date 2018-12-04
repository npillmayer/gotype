/*
Package segment is about Unicode text segmenting.

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
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


Typical Usage

Clients instantiate a UnicodeBreaker object and use it as the
breaking engine for a segmenter. Multiple breaking engines may be
supplied.

  breaker1 := ...
  breaker2 := ...
  segmenter := unicode.NewSegmenter(breaker1, breaker2)
  segmenter.Init(...)
  for segmenter.Next() ...

An example for an UnicodeBreaker is "uax29.WordBreak", a breaker
implementing the UAX#29 word breaking algorithm.

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
	"strconv"
	"strings"

	pool "github.com/jolestar/go-commons-pool"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
)

// We trace to the core-tracer.
var CT tracing.Trace = tracing.CoreTracer

// Object of type UnicodeBreaker represent logic components to split up
// Unicode sequences into smaller parts. They are used by Segmenters
// to supply breaking logic.
type UnicodeBreaker interface {
	CodePointClassFor(rune) int
	StartRulesFor(rune, int)
	ProceedWithRune(rune, int)
	LongestActiveMatch() int
	Penalties() []int
}

type atom struct {
	r       rune
	penalty int
}

var eotAtom atom = atom{rune(0), 0}

func (a *atom) String() string {
	return fmt.Sprintf("[%+q p=%d]", a.r, a.penalty)
}

// A Segmenter receives a sequence of code-points from an io.RuneReader and
// segments it into smaller parts, called segments.
//
// Segmenter provides an interface similar to bufio.Scanner for reading data
// such as a file of Unicode text.
// Successive calls to the Next() method will
// step through the 'segments' of a file, calculating a 'penalty' for breaking
// at this segment. Penalties are numeric values and reflect costs, where
// negative values are to be interpreted as negative costs, i.e. merits.
//
// The specification of a segment is defined by a breaker function of type
// UnicodeBreaker; the default split function breaks the input into words
// (see package uax29).
// Clients may provide more than one UnicodeBreaker.
//
// Segmenting stops unrecoverably at EOF, the first I/O error, or a token too
// large to fit in the buffer.
type Segmenter struct {
	deque                      *deque
	breakers                   []UnicodeBreaker
	reader                     io.RuneReader
	buffer                     *bytes.Buffer
	activeSegment              []byte
	lastPenalty                int // penalty at last break opportunity
	maxSegmentLen              int // maximum length allowed for segments
	nullPenalty                func(int) bool
	aggregate                  PenaltyAggregator
	longestActiveMatch         int
	positionOfBreakOpportunity int
	atEOF                      bool
	err                        error
	inUse                      bool // Next() has been called; buffer is in use.
	voidCount                  int
}

// MaxSegmentSize is the maximum size used to buffer a segment
// unless the user provides an explicit buffer with Segmenter.Buffer().
const MaxSegmentSize = 64 * 1024
const startBufSize = 4096 // Size of initial allocation for buffer.

var (
	ErrTooLong        = errors.New("segment.Segmenter: segment too long for buffer")
	ErrNotInitialized = errors.New("Segmenter not initialized: no input; must call Init(...) first")
)

// Create a new Segmenter by providing breaking logic (UnicodeBreaker).
//
// Before using new segmenters, clients will have to call Init(...) on them.
func NewSegmenter(breakers ...UnicodeBreaker) *Segmenter {
	s := &Segmenter{}
	s.breakers = breakers
	return s
}

// Initialize a Segmenter with an io.RuneReader to read from.
// Re-initializes a segmenter already in use.
func (s *Segmenter) Init(reader io.RuneReader) {
	if reader == nil {
		reader = strings.NewReader("")
	}
	s.reader = reader
	if s.deque == nil {
		s.deque = &deque{} // Q of atoms
		s.buffer = bytes.NewBuffer(make([]byte, 0, startBufSize))
		s.maxSegmentLen = MaxSegmentSize
	} else {
		s.deque.Clear()
		s.longestActiveMatch = 0
		s.atEOF = false
		s.buffer.Reset()
		s.inUse = false
	}
	s.positionOfBreakOpportunity = -1
	if s.aggregate == nil {
		s.aggregate = AddPenalties
	}
	if s.nullPenalty == nil {
		s.nullPenalty = tooBad
	}
}

// Buffer sets the initial buffer to use when scanning and the maximum size of
// buffer that may be allocated during segmenting.
// The maximum segment size is the larger of max and cap(buf).
// If max <= cap(buf), Next() will use this buffer only and do no allocation.
//
// By default, Segmenter uses an internal buffer and sets the maximum token size
// to MaxSegmentntSize.
//
// Buffer panics if it is called after scanning has started. Clients will have
// to call Init(...) again to permit re-setting the buffer.
func (s *Segmenter) Buffer(buf []byte, max int) {
	if s.inUse {
		panic("segment.Buffer: buffer already in use; cannot be re-set")
	}
	s.buffer = bytes.NewBuffer(buf)
	s.maxSegmentLen = max
}

// Err returns the first non-EOF error that was encountered by the
// Segmenter.
func (s *Segmenter) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// Set the null-value for penalties. Penalties equal to this value will
// be treated as if no penalty occured (possibly resulting in the
// suppression of a break opportunity).
//
// TODO: API changed
func (s *Segmenter) SetNullPenalty(isNull func(int) bool) {
	if isNull == nil {
		s.nullPenalty = tooBad
	} else {
		s.nullPenalty = isNull
	}
}

func tooBad(p int) bool {
	return p > 1000
}

// Get the next segment, together with the accumulated penalty for this break.
//
// Next() advances the Segmenter to the next segment, which will then be available
// through the Bytes() or Text() method. It returns false when the segmenting
// stops, either by reaching the end of the input or an error.
// After Next() returns false, the Err() method will return any error
// that occurred during scanning, except for io.EOF.
// For the latter case Err() will return nil.
func (s *Segmenter) Next() bool {
	if s.reader == nil {
		s.setErr(ErrNotInitialized)
	}
	s.inUse = true
	if !s.atEOF {
		//fmt.Println("((((")
		err := s.readEnoughInput()
		//fmt.Println("))))")
		if err != nil && err != io.EOF {
			s.setErr(err)
			s.activeSegment = nil
			return false
		}
	}
	if s.positionOfBreakOpportunity >= 0 { // found a break opportunity
		l := s.getFrontSegment(s.buffer)
		s.activeSegment = s.buffer.Bytes()
		CT.P("length", strconv.Itoa(l)).Debugf("Next() = %v", s.activeSegment)
		return true
	} else {
		s.activeSegment = nil
		return false
	}
}

// Bytes() returns the most recent token generated by a call to Next().
// The underlying array may point to data that will be overwritten by a
// subsequent call to Next(). No allocation is performed.
func (s *Segmenter) Bytes() []byte {
	return s.activeSegment
}

// Text() returns the most recent segment generated by a call to Next()
// as a newly allocated string holding its bytes.
func (s *Segmenter) Text() string {
	return string(s.activeSegment)
}

func (s *Segmenter) Penalty() int {
	return s.lastPenalty
}

// setErr() records the first error encountered.
func (s *Segmenter) setErr(err error) {
	if s.err == nil || s.err == io.EOF {
		s.err = err
	}
}

// TODO
func (s *Segmenter) SetPenaltyAggregator(pa PenaltyAggregator) {
	if pa != nil {
		s.aggregate = pa
	}
}

func (s *Segmenter) readRune() (err error) {
	if s.atEOF {
		err = io.EOF
	} else {
		var r rune
		r, _, err = s.reader.ReadRune()
		CT.P("rune", fmt.Sprintf("%+q", r)).Debug("--------------------------------------")
		if err == nil {
			s.deque.PushBack(r, 0)
		} else if err == io.EOF {
			s.deque.PushBack(eotAtom.r, eotAtom.penalty)
			s.atEOF = true
			err = nil
		} else {
			CT.P("rune", r).Errorf("ReadRune() error: %s", err)
			s.atEOF = true
		}
	}
	return err
}

func (s *Segmenter) readEnoughInput() (err error) {
	for s.positionOfBreakOpportunity < 0 {
		l := s.deque.Len()
		err = s.readRune()
		if err != nil {
			break
		}
		if s.deque.Len() > l { // is this necessary ?
			from := max(0, l-1-s.longestActiveMatch) // old longest match limit
			l = s.deque.Len()
			s.longestActiveMatch = 0
			r, _ := s.deque.Back()
			for _, breaker := range s.breakers {
				cpClass := breaker.CodePointClassFor(r)
				breaker.StartRulesFor(r, cpClass)
				breaker.ProceedWithRune(r, cpClass)
				if breaker.LongestActiveMatch() > s.longestActiveMatch {
					s.longestActiveMatch = breaker.LongestActiveMatch()
				}
				s.insertPenalties(breaker.Penalties())
			}
			s.positionOfBreakOpportunity = s.findBreakOpportunity(from, l-1-s.longestActiveMatch)
			CT.Debugf("distance = %d, active match = %d", s.positionOfBreakOpportunity, s.longestActiveMatch)
			s.printQ()
		} else {
			CT.Error("code-point deque did not grow")
		}
	}
	return err
}

func (s *Segmenter) findBreakOpportunity(from int, to int) int {
	pos := -1
	CT.Debugf("searching for break opportunity from %d to %d: ", from, to-1)
	for i := 0; i < to; i++ {
		_, p := s.deque.At(i)
		if !s.nullPenalty(p) {
			pos = i
			break
		}
	}
	CT.Debugf("break opportunity at %d", pos)
	return pos
}

func (s *Segmenter) insertPenalties(penalties []int) {
	l := s.deque.Len()
	if len(penalties) > l {
		penalties = penalties[0:l]
	}
	for i, p := range penalties {
		r, total := s.deque.At(l - 1 - i)
		total = s.aggregate(total, p)
		s.deque.SetAt(l-1-i, r, total)
	}
}

func (s *Segmenter) getFrontSegment(buf *bytes.Buffer) int {
	seglen := 0
	s.lastPenalty = 0
	buf.Reset()
	l := min(s.deque.Len()-1, s.positionOfBreakOpportunity)
	CT.Debugf("cutting front segment of length 0..%d", l)
	if l > buf.Len() {
		if l > s.maxSegmentLen {
			s.setErr(ErrTooLong)
			return 0
		}
		newSize := max(buf.Len()+startBufSize, l+1)
		if newSize > s.maxSegmentLen {
			newSize = s.maxSegmentLen
		}
		buf.Grow(newSize)
	}
	cnt := 0
	for i := 0; i <= l; i++ {
		r, p := s.deque.PopFront()
		written, _ := buf.WriteRune(r)
		seglen += written
		cnt++
		s.lastPenalty = p
	}
	CT.Debugf("front segment is of length %d/%d", seglen, cnt)
	s.positionOfBreakOpportunity = s.findBreakOpportunity(0, s.deque.Len()-1-s.longestActiveMatch)
	s.printQ()
	return seglen
}

// Debugging helper. Print the content of the current queue to the debug log.
func (s *Segmenter) printQ() {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Q #%d: ", s.deque.Len()))
	for i := 0; i < s.deque.Len(); i++ {
		var a atom
		a.r, a.penalty = s.deque.At(i)
		sb.WriteString(fmt.Sprintf(" <- %s", a.String()))
	}
	sb.WriteString(" .")
	CT.Debugf(sb.String())
}

// ----------------------------------------------------------------------

// Functions of type NfaStateFn try to match a rune (Unicode code-point).
// The caller may provide a third argument, which should be a rune class.
// Rune (code-point) classes are described in various Unicode standards
// and annexes. One such standard (UAX#29) desribes classes to help
// splitting up text into graphemes or words. An example class may be
// a class of western language alphabetic characters AL, of which runes
// 'A' and 'é' would be part of.
//
// The first argument is a Recognizer (see the definition of
// type Recognizer in this package), which carries this state function.
//
// NfaStateFn – after they matched a rune – must return another NfaStateFn,
// which will then in turn be called to process the next rune. The process
// of matching a string will stop as soon as a NfaStateFn returns nil.
type NfaStateFn func(*Recognizer, rune, int) NfaStateFn

// A Recognizer represents an automata to recognize sequences of runes
// (i.e. Unicode code-points). Its main functionality is performed by
// an embedded NfaStateFn. The first NfaStateFn to use is provided with
// the constructor.
//
// Recognizer's state functions must be careful to increment MatchLen
// with each matched rune. Failing to do so may result in incorrect splits
// of text.
//
// Semantics of Expect and UserData are up to the client and not used by
// the default mechanism.
//
// It is not mandatory to use Recognizers for segmenting text. The type is
// provided for easier implementation of types implementing UnicodeBreaker.
// Recognizers implement interface RuneSubscriber and UnicodeBreakers will
// use a UnicodePublisher to interact with them.
type Recognizer struct {
	Expect    int         // next code-point to expect; semantics are up to the client
	MatchLen  int         // length of active match
	UserData  interface{} // clients may need to store additional information
	penalties []int       // penalties to return, used internally in DoAccept()
	nextStep  NfaStateFn  // next step of a DFA
}

// Create a new Recognizer. This is rarely used, as clients rather should
// call NewPooledRecognizer().
func NewRecognizer(codePointClass int, next NfaStateFn) *Recognizer {
	rec := &Recognizer{}
	rec.Expect = codePointClass
	rec.nextStep = next
	return rec
}

// Recognizers are short-lived objects. To avoid multiple allocation of
// small objects we will pool them.
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
			return rec, nil
		})
	globalRecognizerPool.ctx = context.Background()
	config := pool.NewDefaultPoolConfig()
	//config.LIFO = false
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
	rec.penalties = nil
	rec.Expect = 0
	rec.MatchLen = 0
	rec.nextStep = nil
	_ = globalRecognizerPool.opool.ReturnObject(globalRecognizerPool.ctx, rec)
}

// Simple stringer for debugging purposes.
func (rec *Recognizer) String() string {
	if rec == nil {
		return "[nil rule]"
	}
	return fmt.Sprintf("[%d -> done=%v]", rec.Expect, rec.Done())
}

// Signal a Recognizer that it has been unsubscribed from a RunePublisher;
// usually after the Recognizer's NfaStateFn has returned nil.
//
// Interface RuneSubscriber
func (rec *Recognizer) Unsubscribed() {
	rec.releaseIntoPool()
}

// A Recognizer signals that it is done matching runes.
// If MatchLength() > 0 is has been accepting a sequence of runes,
// otherwise it has aborted to further try a match.
//
// Interface RuneSubscriber
func (rec *Recognizer) Done() bool {
	return rec.nextStep == nil
}

// Interface RuneSubscriber
func (rec *Recognizer) MatchLength() int {
	return rec.MatchLen
}

// Interface RuneSubscriberyy
func (rec *Recognizer) RuneEvent(r rune, codePointClass int) []int {
	//fmt.Printf("received rune event: %+q / %d\n", r, codePointClass)
	var penalties []int
	if rec.nextStep != nil {
		//CT.Infof("  calling func = %v", rec.nextStep)
		rec.nextStep = rec.nextStep(rec, r, codePointClass)
	} else {
		//CT.Info("  not calling func = nil")
	}
	if rec.Done() && rec.MatchLen > 0 { // accepted a match
		penalties = rec.penalties
	}
	//CT.Infof("    subscriber:      penalites = %v, done = %v, match len = %d", penalties, rec.Done(), rec.MatchLength())
	return penalties
}

// --- Standard Recognizer Rules ----------------------------------------

// Get a state function which signals abort.
func DoAbort(rec *Recognizer) NfaStateFn {
	rec.MatchLen = 0
	return nil
}

// Get a state function which signals accept, together with break
// penalties for matches runes (in reverse sequence).
func DoAccept(rec *Recognizer, penalties ...int) NfaStateFn {
	rec.MatchLen++
	rec.penalties = penalties
	CT.Debugf("ACCEPT with %v", rec.penalties)
	return nil
}

// --- Rune Publishing and Subscription ---------------------------------

// RuneSubscribers are receivers of rune events, i.e. messages to
// process a new code-point (rune). If they can match the rune, they
// will expect further runes, otherwise they abort. To they are finished,
// either by accepting or rejecting input, they set Done() to true.
// A successful acceptance of input is signalled by Done()==true and
// MatchLength()>0.
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
	SubscribeMe(RuneSubscriber) RunePublisher // subscribe an additional rune subscriber
	PublishRuneEvent(r rune, codePointClass int) (longestDistance int, penalties []int)
	SetPenaltyAggregator(pa PenaltyAggregator) // function to aggregate break penalties
}

// Create a new default RunePublisher.
func NewRunePublisher() *DefaultRunePublisher {
	rpub := &DefaultRunePublisher{}
	rpub.aggregate = AddPenalties
	return rpub
}

// Trigger a rune event notification to all subscribers. Rune events
// include the rune (code-point) and an optional code-point class for
// the rune.
//
// Return values are: the longest active match and a slice of penalties.
// These values usually are collected from the RuneSubscribers.
// Penalties will be overwritten by the next call to PublishRuneEvent().
// Clients will have to make a copy if they want to preserve penalty
// values.
//
// Interface RunePublisher
func (rpub *DefaultRunePublisher) PublishRuneEvent(r rune, codePointClass int) (int, []int) {
	longest := 0
	if rpub.penaltiesTotal == nil {
		rpub.penaltiesTotal = make([]int, 1024)
	}
	//CT.Infof("pre-publish(): total penalites = %v", rpub.penaltiesTotal)
	rpub.penaltiesTotal = rpub.penaltiesTotal[:0]
	//CT.Infof("pre-publish(): total penalites = %v", rpub.penaltiesTotal)
	// pre-condition: no subscriber is Done()
	for i := rpub.Len() - 1; i >= 0; i-- {
		subscr := rpub.at(i)
		penalties := subscr.RuneEvent(r, codePointClass)
		//CT.Infof("    publish():       penalites = %v", penalties)
		for j, p := range penalties { // aggregate all penalties
			if j >= len(rpub.penaltiesTotal) {
				rpub.penaltiesTotal = append(rpub.penaltiesTotal, p)
			} else {
				rpub.penaltiesTotal[j] = rpub.aggregate(rpub.penaltiesTotal[j], p)
			}
		}
		//CT.Infof("    publish(): total penalites = %v", rpub.penaltiesTotal)
		if !subscr.Done() { // compare against longest active match
			if d := subscr.MatchLength(); d > longest {
				longest = d
			}
		}
		rpub.Fix(i) // re-order heap if subscr.Done()
	}
	//CT.Infof("pre-publish(): total penalites = %v", rpub.penaltiesTotal)
	// now unsubscribe all done subscribers
	for subscr := rpub.PopDone(); subscr != nil; subscr = rpub.PopDone() {
		subscr.Unsubscribed()
	}
	return longest, rpub.penaltiesTotal
}

// Function type for methods of penalty-aggregation. Aggregates all the
// break penalties each a break-point to a single penalty value at that point.
//
// TODO: for segmenters: add argument for *UnicodeBreaker.
// For unicodeBreakers: leave it this way.
type PenaltyAggregator func(int, int) int

// Interface RunePublisher
func (rpub *DefaultRunePublisher) SetPenaltyAggregator(pa PenaltyAggregator) {
	if pa == nil {
		rpub.aggregate = AddPenalties
	} else {
		rpub.aggregate = pa
	}
}

// The default function to aggregate break-penalties.
// Simply adds up all penalties at each break position, respectively.
func AddPenalties(total int, p int) int {
	return total + p
}

// Alternative function to aggregate break-penalties.
// Returns maximum of all penalties at each break position.
func MaxPenalties(total int, p int) int {
	return max(total, p)
}

// Interface RunePublisher
func (rpub *DefaultRunePublisher) SubscribeMe(rsub RuneSubscriber) RunePublisher {
	if rpub.aggregate == nil { // this is necessary as we allow uninitialzed DefaultRunePublishers
		rpub.aggregate = AddPenalties
	}
	rpub.Push(rsub)
	return rpub
}

// ----------------------------------------------------------------------

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
