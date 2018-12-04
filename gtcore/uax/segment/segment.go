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

Segmenter provides an interface similar to bufio.Scanner for reading data
such as a file of Unicode text.
Similar to Scanner's Scan() function, successive calls to a segmenter's
Next() method will step through the 'segments' of a file.
Clients are able to get runes of the segment by calling Bytes() or Text().
Unlike Scanner, segmenters are calculating a 'penalty' for breaking
at this segment. Penalties are numeric values and reflect costs, where
negative values are to be interpreted as negative costs, i.e. merits.

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
*/
package segment

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/gtcore/uax"
	"github.com/npillmayer/gotype/gtcore/uax/uax29"
)

// We trace to the core-tracer.
var CT tracing.Trace = tracing.CoreTracer

// A Segmenter receives a sequence of code-points from an io.RuneReader and
// segments it into smaller parts, called segments.
//
// The specification of a segment is defined by a breaker function of type
// UnicodeBreaker; the default split function breaks the input into words
// (see package uax29).
type Segmenter struct {
	deque                      *deque
	breakers                   []uax.UnicodeBreaker
	reader                     io.RuneReader
	buffer                     *bytes.Buffer
	activeSegment              []byte
	lastPenalty                int // penalty at last break opportunity
	maxSegmentLen              int // maximum length allowed for segments
	nullPenalty                func(int) bool
	aggregate                  uax.PenaltyAggregator
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
// Clients may provide more than one UnicodeBreaker. Specifying no
// UnicodeBreaker results in using a word breaker.
//
// Before using new segmenters, clients will have to call Init(...) on them.
func NewSegmenter(breakers ...uax.UnicodeBreaker) *Segmenter {
	s := &Segmenter{}
	if len(breakers) == 0 {
		breakers = []uax.UnicodeBreaker{uax29.NewWordBreaker()}
	}
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
		s.aggregate = uax.AddPenalties
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
func (s *Segmenter) SetPenaltyAggregator(pa uax.PenaltyAggregator) {
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
