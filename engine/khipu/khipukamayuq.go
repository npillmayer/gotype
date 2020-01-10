package khipu

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"bufio"
	"io"
	"strings"

	"github.com/npillmayer/gotype/core/dimen"
	params "github.com/npillmayer/gotype/core/parameters"
	"github.com/npillmayer/gotype/core/uax/segment"
	"github.com/npillmayer/gotype/core/uax/uax14"
	"github.com/npillmayer/gotype/core/uax/uax29"
	"github.com/npillmayer/gotype/gtlocate"
	"golang.org/x/text/unicode/norm"
)

// A TypesettingPipeline consists of steps to produce a khipu from text.
type TypesettingPipeline struct {
	input       io.RuneReader
	linewrap    *uax14.LineWrap
	wordbreaker *uax29.WordBreaker
	segmenter   *segment.Segmenter
	words       *segment.Segmenter
}

// KnotEncode transforms an input text into a khipu.
func KnotEncode(text io.Reader, pipeline *TypesettingPipeline, regs *params.TypesettingRegisters) *Khipu {
	if regs == nil {
		regs = params.NewTypesettingRegisters()
	}
	pipeline = PrepareTypesettingPipeline(text, pipeline)
	khipu := NewKhipu()
	seg := pipeline.segmenter
	for seg.Next() {
		fragment := seg.Text()
		CT().Infof("next segment = '%s'\twith penalties %v", fragment, seg.Penalties())
		k := createPartialKhipuFromSegment(seg, pipeline, regs)
		if regs.N(params.P_MINHYPHENLENGTH) < dimen.Infty {
			HyphenateTextBoxes(k, pipeline, regs)
		}
		khipu.AppendKhipu(k)
	}
	CT().Infof("resulting khipu = %s\n", khipu)
	return khipu
}

// Call this for creating a sub-khipu from a segment. The fist parameter
// is a segmenter which already has detected a segment, i.e. seg.Next()
// has been called successfully.
//
// Calls to createPartialKhipuFromSegment will panic if one of its
// arguments is invalid.
//
// Returns a khipu consisting of text-boxes, glues and penalties.
func createPartialKhipuFromSegment(seg *segment.Segmenter, pipeline *TypesettingPipeline, regs *params.TypesettingRegisters) *Khipu {
	khipu := NewKhipu()
	if seg.Penalties()[0] < 1000 { // broken by primary breaker // TODO: this API is aweful...
		// fragment is terminated by possible line wrap opportunity
		if seg.Penalties()[1] < 1000 { // broken by secondary breaker, too
			if seg.Penalties()[1] == segment.PenaltyAfterWhitespace {
				g := NewGlue(5*dimen.BP, 1*dimen.BP, 2*dimen.BP)
				p := Penalty(seg.Penalties()[1])
				khipu.AppendKnot(g).AppendKnot(p)
			} else {
				b := NewTextBox(seg.Text())
				p := Penalty(dimen.Infty)
				khipu.AppendKnot(b).AppendKnot(p)
			}
		}
	} else { // segments is broken by secondary breaker
		// fragment is start or end of a span of whitespace
		if seg.Penalties()[1] == segment.PenaltyBeforeWhitespace {
			// close a text box which is not a possible line wrap position
			b := NewTextBox(seg.Text())
			p := Penalty(dimen.Infty)
			khipu.AppendKnot(b).AppendKnot(p)
		} else {
			// close a span of whitespace
			g := NewGlue(5*dimen.PT, 1*dimen.PT, 2*dimen.PT)
			p := Penalty(seg.Penalties()[1])
			khipu.AppendKnot(g).AppendKnot(p)
		}
	}
	return khipu
}

// HyphenateTextBoxes hypenates all the words in a khipu.
// Words are contained inside TextBox knots.
//
// Hyphenation is governed by the typesetting registers provided.
// If regs is nil, no hyphenation is done.
func HyphenateTextBoxes(khipu *Khipu, pipeline *TypesettingPipeline, regs *params.TypesettingRegisters) {
	if regs == nil || khipu == nil {
		return
	}
	k := make([]Knot, 0, khipu.Length())
	iterator := NewCursor(khipu)
	for iterator.Next() {
		if iterator.Knot().Type() != KTTextBox { // can only hyphenate text knots
			k = append(k, iterator.Knot())
			continue
		}
		CT().Infof("knot = %v | %v", iterator.Knot(), iterator.Knot())
		text := iterator.AsTextBox().text
		pipeline.words.Init(strings.NewReader(text))
		for pipeline.words.Next() {
			word := pipeline.words.Text()
			CT().Infof("   word = '%s'", word)
			var syllables []string
			isHyphenated := false
			if len(word) >= regs.N(params.P_MINHYPHENLENGTH) {
				if syllables, isHyphenated = HyphenateWord(word, regs); isHyphenated {
					hyphen := NewKnot(KTDiscretionary)
					for _, sy := range syllables[:len(syllables)-1] {
						k = append(k, NewTextBox(sy))
						k = append(k, hyphen)
					}
					k = append(k, NewTextBox(syllables[len(syllables)-1]))
				}
			}
			if !isHyphenated {
				if word == text {
					k = append(k, iterator.Knot())
				} else {
					k = append(k, NewTextBox(word))
				}
			}
		}
	}
	khipu.knots = k
}

// PrepareTypesettingPipeline checks if a typesetting pipeline is correctly
// initialized and creates a new one if is is invalid.
//
// We use a uax14.LineWrapper as the primary breaker and
// use a segment.SimpleWordBreaker to extract spans of whitespace.
// For the inner loop we use a uax29.WordBreaker.
// This is a default configuration adequate for western languages.
func PrepareTypesettingPipeline(text io.Reader, pipeline *TypesettingPipeline) *TypesettingPipeline {
	// wrap a normalization-reader around the input
	if pipeline == nil {
		pipeline = &TypesettingPipeline{}
	}
	pipeline.input = bufio.NewReader(norm.NFC.Reader(text))
	if pipeline.segmenter == nil {
		pipeline.linewrap = uax14.NewLineWrap()
		pipeline.segmenter = segment.NewSegmenter(pipeline.linewrap, segment.NewSimpleWordBreaker())
		pipeline.segmenter.Init(pipeline.input)
		pipeline.wordbreaker = uax29.NewWordBreaker()
		pipeline.words = segment.NewSegmenter(pipeline.wordbreaker)
	}
	return pipeline
}

// HyphenateWord hyphenates a single word.
func HyphenateWord(word string, regs *params.TypesettingRegisters) ([]string, bool) {
	dict := gtlocate.Dictionary(regs.S(params.P_LANGUAGE))
	ok := false
	if dict == nil {
		panic("TODO not yet implemented: find dictionnary for language")
	}
	CT().Infof("   will try to hyphenate word")
	splitWord := dict.Hyphenate(word)
	if len(splitWord) > 1 {
		ok = true
	}
	CT().Infof("   %v", splitWord)
	return splitWord, ok
}

/*
func UAX14LineWrap(text string, regs *params.TypesettingRegisters) *Khipu {
	sread := strings.NewReader(text) // wrap a reader around the CDATA string
	nfcread := norm.NFC.Reader(sread)
	scanner := bufio.NewScanner(nfcread) // wrap a buffered scanner around it
	scanner.Split(segment.SplitWords)    // split on words according to UAX#29
	for scanner.Scan() {
		tokenBytes := scanner.Bytes()
		fmt.Printf("%s \n", tokenBytes)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("ERROR") // TODO
	}

	// TODO: line wrap
	// https://github.com/gorilla/i18n/blob/master/linebreak/linebreak.go

	return nil
}


// Simple Lösung. Unvollständig.
// Erkennt Emoji-Zeichenketten nicht; nur Zeichen + diakr. Zeichen
func grLen(s string) int { // length in terms of graphemes
	if len(s) == 0 {
		return 0
	}
	gr := 1
	_, s1 := utf8.DecodeRuneInString(s)
	for _, r := range s[s1:] {
		if !unicode.Is(unicode.Mn, r) {
			gr++
		}
	}
	return gr
}

func graphemes(s string) int {
	n := 0
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		fmt.Printf("%c %v\n", r, size)
		s = s[size:]
		n++
	}
	return n
}

func iterateOverGraphemes(s string) {
	var it norm.Iter
	it.InitString(norm.NFC, s)
	for !it.Done() {
		b := it.Next()
		r, size := utf8.DecodeRuneInString(string(b[:]))
		fmt.Printf("%c %v\n", r, size)
	}
}
*/

// https://github.com/blevesearch/segment
/*
A Go library for performing Unicode Text Segmentation as described in
Unicode Standard Annex #29
*/
// Alternativen: https://github.com/go-ego/gse
//
/*
func iterateOverWords(s *strings.Reader) {
	segmenter := segment.NewWordSegmenter(s)
	for segmenter.Segment() {
		tokenBytes := segmenter.Bytes()
		tokenType := segmenter.Type()
		r := string(tokenBytes[:])
		fmt.Printf("%s %v\n", r, tokenType)
	}
}
*/
