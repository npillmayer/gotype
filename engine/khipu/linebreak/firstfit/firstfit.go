/*
Package firstfit implements a straightforward line-breaking algorithm
where lines are broken at the first suitable breakpoint.

Wikipedia:

	1. |  SpaceLeft := LineWidth
	2. |  for each Word in Text
	3. |      if (Width(Word) + SpaceWidth) > SpaceLeft
	4. |           insert line break before Word in Text
	5. |           SpaceLeft := LineWidth - Width(Word)
	6. |      else
	7. |           SpaceLeft := SpaceLeft - (Width(Word) + SpaceWidth)

With khipus, we have space before words, as a rule, not after them. Additionally, we
break at penalty knots, not on whitespace. However, the implementation looks
roughly like this:

	 1. |  SpaceUsed := 0; firstInLine := true
	 2. |  for each Knot in Khipu
	 3. |      if Knot.Type == Word
	 4. |          if (SpaceUsed + Width(Word)).MinWidth > LineWidth
	 5. |             insert line break before Word in Text
	 6. |             SpaceUsed := 0
	 7. |             firstInLine := true
	 8. |          else
	 9. |             SpaceUsed := SpaceUsed + (Width(Word) + SpaceWidth)
	10. |             firstInLine := false
	11. |      else if not firstInLine
	12. |          SpaceUsed := SpaceUsed + Width(Knot)

What counts as a word is not so clear with international scripts. We rely on the
khipukamayuq to insert appropriate penalties before line-breaking happens.

BSD License

Copyright (c) 2017–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */
package firstfit

import (
	"errors"
	"fmt"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/engine/khipu"
	"github.com/npillmayer/gotype/engine/khipu/linebreak"
)

// T traces to the core tracer.
func T() tracing.Trace {
	return gtrace.CoreTracer
}

// We use a small object to manage line-breaking information for a client call.
// We do not expose it to clients right now, but this may change as soon as we
// make line-breaking scriptable.
type linebreaker struct {
	buffer    []khipu.Mark
	knot      khipu.Knot
	params    *linebreak.Parameters
	cursor    linebreak.Cursor
	parshape  linebreak.ParShape
	eof       bool
	pos       int
	check     int
	linecount int
}

// newLinebreaker creates an adequately initialized first-fit linebreaker.
// If params are not given, default parameters will be used.
// cursor and parshape are mandatory.
func newLinebreaker(cursor linebreak.Cursor, parshape linebreak.ParShape,
	params *linebreak.Parameters) (*linebreaker, error) {
	if cursor == nil || parshape == nil {
		return nil, errors.New("cannot break a paragraph without cursor or parshape")
	}
	if params == nil {
		params = linebreak.DefaultParameters
	}
	lb := &linebreaker{
		cursor:   cursor,
		parshape: parshape,
		params:   params,
		buffer:   make([]khipu.Mark, 0, 20), // backtracking buffer
		pos:      -1,
		check:    -1,
		//segment:  &segment{},
	}
	return lb, nil
}

// BreakParagraph will find 'first fit' breakpoints for a paragraph of text.
// The method is similar to the one usually used by web browsers. It simply collects
// line material until the current line-length is exhausted, then continues on a
// new line.
func BreakParagraph(cursor linebreak.Cursor, parshape linebreak.ParShape,
	params *linebreak.Parameters) ([]khipu.Mark, error) {
	//
	lb, err := newLinebreaker(cursor, parshape, params)
	if err != nil {
		return nil, err
	}
	return lb.FindBreakpoints()
}

func (lb *linebreaker) FindBreakpoints() ([]khipu.Mark, error) {
	firstInLine := true
	breakpoints := make([]khipu.Mark, 1, 10)
	breakpoints[0] = provisionalMark(-1)
	lineno := 0
	spaceUsed := &segment{}
	knot := lb.next()
	xxx := 0
	last := lb.mark()
	for knot != nil {
		xxx++
		if xxx > 28 {
			//break
		}
		linelen := lb.parshape.LineLength(lineno)
		gtrace.CoreTracer.Debugf("_______________ %v ___________________", knot)
		if knot.Type() == khipu.KTPenalty { // TODO discretionaries
			last = lb.mark()
			penalty := lb.penalty()
			spaceUsed.append(knot)
			segm := spaceUsed.width(lb.params)
			if penalty.Demerits() < linebreak.InfinityDemerits {
				T().Debugf("penalty %v is acceptable", penalty.Demerits())
				T().Debugf("segm=%v", segm)
				if segm.Min > linelen { // overshoot
					if frag := lb.backtrack(); frag != nil {
						// start new line
						T().Infof("backtracked to %v", lb.knot)
						breakpoints = lb.linebreak(breakpoints)
						//spaceUsed.reset(frag)
						firstInLine = true
					} else {
						// overfull hbox
						T().Infof("Overfull box at line %d", lb.linecount+1)
						// start new line
						breakpoints = lb.linebreak(breakpoints)
						spaceUsed.reset(linebreak.WSS{})
						firstInLine = true
					}
				} else if segm.Max >= linelen { // can reach EOL
					if penalty.Demerits() <= lb.params.Tolerance {
						if !lb.checkpoint() {
							panic("CANNOT SET CHECKPOINT")
						}
					} else {
						// save for later review, if no other feasible break is found ? How ?
						if !lb.checkpoint() {
							panic("CANNOT SET CHECKPOINT")
						}
						gtrace.CoreTracer.Infof("setting checkpoint with demerits=%v", penalty.Demerits())
					}
				} else {
					if !lb.checkpoint() {
						panic("CANNOT SET CHECKPOINT")
					}
					gtrace.CoreTracer.Infof("setting emergency with demerits=%v", penalty.Demerits())
				}
			} else {
				// no feasible break
			}
		} else if firstInLine {
			if !knot.IsDiscardable() { // do not add space to front of line
				spaceUsed.append(knot)
				firstInLine = false
			}
		} else { // append everything to line
			spaceUsed.append(knot)
		}
		T().Debugf("segment = %v", spaceUsed)
		knot = lb.next()
	}
	breakpoints = append(breakpoints, last)
	return breakpoints, nil
}

// penaltyAt iterates over all penalties, starting at the current cursor check, and
// collects penalties, searching for the most significant one.
// Will return
//
//        -10000, if present
//        max(p1, p2, ..., pn) otherwise
//
// Returns the most significant penalty. Advances the cursor over all adjacent penalties.
// After this, the cursor check may not reflect the position of the significant penalty.
func (lb *linebreaker) penalty() khipu.Penalty {
	if lb.eof || lb.knot.Type() != khipu.KTPenalty {
		return khipu.Penalty(linebreak.InfinityDemerits)
	}
	penalty := lb.knot.(khipu.Penalty)
	ignore := false // final penalty found, ignore all other penalties
	knot, ok := lb.peek()
	for ok {
		if knot.Type() == khipu.KTPenalty {
			lb.next() // move over knot
			if ignore {
				break
			} else {
				p := knot.(khipu.Penalty)
				if p.Demerits() <= linebreak.InfinityMerits { // -10000 must break (like in TeX)
					penalty = p
					ignore = true
				} else if p.Demerits() > penalty.Demerits() {
					penalty = p
				}
			}
			knot, ok = lb.peek()
		} else {
			ok = false
		}
	}
	p := khipu.Penalty(linebreak.CapDemerits(penalty.Demerits()))
	return p
}

type segment struct {
	length       linebreak.WSS //
	breakDiscard linebreak.WSS // sum of discardable space while lookinf for next breakpoint
}

// append the width information of a knot at the end of a segment.
// if the knot is a discardable item, s.breakDiscard is updated as well.
func (s *segment) append(knot khipu.Knot) {
	s.length = s.length.Add(linebreak.WSS{}.SetFromKnot(knot))
	if knot.IsDiscardable() {
		s.breakDiscard = s.breakDiscard.Add(linebreak.WSS{}.SetFromKnot(knot))
	} else {
		s.breakDiscard = linebreak.WSS{}
	}
}

// width returns the widths of the current partial line, subtracting
// space at the end of the segment (as this will be dropped).
// width is respecting params.LeftSkip & RightSkip.
func (s *segment) width(params *linebreak.Parameters) linebreak.WSS {
	return s.length
	// segw := s.length.Subtract(s.breakDiscard)
	// w := linebreak.WSS{}.SetFromKnot(params.LeftSkip)
	// segw = segw.Add(w)
	// w = linebreak.WSS{}.SetFromKnot(params.RightSkip)
	// segw = segw.Add(w)
	// return segw
}

func (s *segment) reset(wss linebreak.WSS) {
	s.length = wss
	s.breakDiscard = linebreak.WSS{}
}

func (s *segment) String() string {
	return fmt.Sprintf("|---%.2f---->", s.length.W.Points())
}

func (lb *linebreaker) linebreak(breakpoints []khipu.Mark) []khipu.Mark {
	//panic("LINEBREAK")
	lb.linecount++
	gtrace.CoreTracer.Debugf("new line #%d", lb.linecount)
	// reset pos, check, buffer
	breakpoints = append(breakpoints, lb.mark())
	lb.buffer = lb.buffer[:0]
	lb.check = -1
	lb.pos = -1
	return breakpoints
}

func (lb *linebreaker) next() khipu.Knot {
	if lb.pos < 0 {
		if lb.cursor.Next() {
			lb.knot = lb.cursor.Knot()
			return lb.knot
		}
		return nil // end of khipu
	}
	lb.pos++
	if lb.pos >= len(lb.buffer) {
		if lb.cursor.Next() {
			lb.buffer = append(lb.buffer, lb.cursor.Mark())
		} else {
			lb.eof = true
			return nil
		}
	}
	lb.knot = lb.buffer[lb.pos].Knot()
	return lb.knot
}

func (lb *linebreaker) mark() khipu.Mark {
	if lb.pos < 0 {
		return lb.cursor.Mark()
	}
	return lb.buffer[lb.pos]
}

func (lb *linebreaker) peek() (khipu.Knot, bool) {
	if lb.pos < 0 {
		return lb.cursor.Peek()
	}
	if lb.pos+1 >= len(lb.buffer) {
		return lb.cursor.Peek()
	}
	return lb.buffer[lb.pos+1].Knot(), true
}

func (lb *linebreaker) checkpoint() bool {
	if lb.pos > -1 {
		lb.buffer = lb.buffer[lb.pos:]
		lb.check, lb.pos = 0, 0
		lb.knot = lb.buffer[lb.pos].Knot()
		return true
	} else if lb.cursor.Knot() != nil {
		lb.buffer = lb.buffer[:0]
		lb.buffer = append(lb.buffer, lb.cursor.Mark())
		lb.check, lb.pos = 0, 0
		lb.knot = lb.buffer[lb.pos].Knot()
		return true
	}
	return false
}

func (lb *linebreaker) backtrack() khipu.Knot {
	if lb.check < 0 {
		return nil
	}
	lb.pos = lb.check
	lb.knot = lb.buffer[lb.pos].Knot()
	return lb.knot
}

type provisionalMark int // provisional mark from an integer position

func (m provisionalMark) Position() int    { return int(m) }
func (m provisionalMark) Knot() khipu.Knot { return khipu.Penalty(-10000) }
