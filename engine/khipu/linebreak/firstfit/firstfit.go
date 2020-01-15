/*
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
	"github.com/npillmayer/gotype/engine/khipu"
	"github.com/npillmayer/gotype/engine/khipu/linebreak"
)

/*
Wikipedia:

	1. |  SpaceLeft := LineWidth
	2. |  for each Word in Text
	3. |      if (Width(Word) + SpaceWidth) > SpaceLeft
	4. |         insert line break before Word in Text
	5. |         SpaceLeft := LineWidth - Width(Word)
	6. |      else
	7. |         SpaceLeft := SpaceLeft - (Width(Word) + SpaceWidth)
*/

// Line 3 seems to be an error: We need to just fit in the word, not word+space.
// All depends wether the space is counted before the word or afterwards. With khipus
// we have space before words, as a rule.
// Therefore we implement:
/*
	1. |  SpaceUsed := 0; firstInLine := true
	2. |  for each Knot in Khipu
	2. |      if Knot.Type == Word
	3. |         if (SpaceUsed + Width(Word)).MinW > LineWidth
	4. |            insert line break before Word in Text
	5. |            SpaceUsed := 0
	3. |            firstInLine := true
	6. |         else
	7. |            SpaceUsed := SpaceUsed + (Width(Word) + SpaceWidth)
	8. |            firstInLine := false
	8. |      else if not firstInLine
	8. |         SpaceUsed := SpaceUsed + Width(Knot)
*/

// What counts as a word is not so clear with international scripts. We rely on the
// khipukamayuq to insert appropriate penalties before line-breaking happens.

// BreakParagraph will find 'first fit' breakpoints for a paragraph of text.
// The method is similar to the one usually used by web browsers. It simply collects
// line material until the current line-length is exhausted, then continues on a
// new line.
func BreakParagraph(cursor linebreak.Cursor, parshape linebreak.Parshape,
	params *linebreak.Parameters) ([]khipu.Mark, error) {
	//
	spaceUsed := linebreak.WSS
	firstInLine := true
	breakpoints := make([]khipu.Mark, 0, 10)
	lineno := 0
	for cursor.Next() {
		if cursor.Mark().Knot().Type() == khipu.KTPenalty { // TODO discretionaries
			penalty, mark := penaltyAt(cursor)
			linelen := parshape.LineWidth(lineno)
			if spaceUsed.MinW() >= linelen {
			}
			if penalty < linebreak.InfinityDemerits {
				if penalty <= params.Tolerance {
					// immediate break
					firstInLine = true
				} else {
					// save for later review, if no other feasible break is found ? How ?
				}
			} else {
				// no feasible break
			}
		} else if !firstInLine {
			// TODO
			// append to line
		}
	}
	return nil, breakpoints
}

// penaltyAt iterates over all penalties, starting at the current cursor mark, and
// collects penalties, searching for the most significant one.
// Will return
//
//        -10000, if present
//        max(p1, p2, ..., pn) otherwise
//
// Returns the most significant penalty. Advances the cursor over all adjacent penalties.
// After this, the cursor mark may not reflect the position of the significant penalty.
func penaltyAt(cursor linebreak.Cursor) (khipu.Penalty, khipu.Mark) {
	if cursor.Knot().Type() != khipu.KTPenalty {
		return khipu.Penalty(linebreak.InfinityDemerits), cursor.Mark()
	}
	penalty := cursor.Knot().(khipu.Penalty)
	ignore := false // final penalty found, ignore all other penalties
	knot, ok := cursor.Peek()
	for ok {
		if knot.Type() == khipu.KTPenalty {
			cursor.Next() // advance to next penalty
			if ignore {
				break // just skip over adjacent penalties
			}
			p := knot.(khipu.Penalty)
			if p.Demerits() <= linebreak.InfinityMerits { // -10000 must break (like in TeX)
				penalty = p
				ignore = true
			} else if p.Demerits() > penalty.Demerits() {
				penalty = p
			}
			knot, ok = cursor.Peek() // now check next knot
		} else {
			ok = false
		}
	}
	p := khipu.Penalty(linebreak.CapDemerits(penalty.Demerits()))
	return p, cursor.Mark()
}
