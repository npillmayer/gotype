package khipus

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/blevesearch/segment"
	params "github.com/npillmayer/gotype/gtcore/parameters"
	"golang.org/x/text/unicode/norm"
)

/*
Interfaces and methods to create "bead chains" of typesetting items.
Bead chains are lists of items, such as glyphs, kerns, glue, etc.
Sometimes we will call these kinds of items "beads".

Bead chains are the input for linebreakers.

Finding possible line-wrap points (Unicode spec):
http://unicode.org/reports/tr14/

JavaScript implementation: https://github.com/foliojs/linebreak (basiert
auf deprecated pair table)
Besser: https://github.com/niklasvh/css-line-break  (JavaScript)
Python: https://uniseg-python.readthedocs.io/en/latest/ bzw.
https://bitbucket.org/emptypage/uniseg-python/src/e4077d17d026c36999b89c10081a85b219e1fb7b/uniseg/?at=default
online: http://unicode.org/cldr/utility/breaks.jsp

Golang package unicode provides all sorts of code-point ranges:
https://golang.org/pkg/unicode/

Moreover: https://godoc.org/golang.org/x/text/unicode

Moreover: https://godoc.org/golang.org/x/text/width

Moreover: "golang.org/x/text/internal/ucd":
Package ucd provides a parser for Unicode Character Database files,
the format of which is defined in https://www.unicode.org/reports/tr44/.
See https://www.unicode.org/Public/UCD/latest/ucd/ for example files.
It currently does not support substitutions of missing fields.

Es wird offenbar an einem "segment" package für x/text gearbeitet, so dass
man nicht so viel Aufwand in eine Eigenlösung investieren sollte.

// -----------------------------------------------------------------------

(0) Normalize -> Word Boundaries
    Graphemes ?

A: Finding feasible break positions

(1) Mandatory breaks + prohibiting no-break points

(2) Natural text wrap (words in many scripts, syllables/character in east asia, etc.)
    -> Unicode UA#14 Line Breaking (http://unicode.org/reports/tr14/)
    Algorithm: http://unicode.org/reports/tr14/#Algorithm

(3) Bidi
    https://godoc.org/golang.org/x/text/unicode/bidi

(4) Hyphenation
    Lliang patterns + language-specific code

(5) Translate feasible breakpoints to penalties, glue and dicretionaries

B: Deciding break positions

(1) Shape text -> Glyphs
    + alternative glyphs (end-of-line condensed in Arabic, etc.)

(2) Translate everything to node list

(3) Apply line-breaking algorithm (simple, K&P, etc.)

C: Justify text

(1) Set glue


// ----------------------------------------------------------------------

Khipus: https://www.sapiens.org/technology/khipu-incas-knotty-history/

the khipukamayuqs (Quechua for “knot-makers/animators”) encoded administrative
data such as census figures and tax allocation in the twisted strings of these
ancient spreadsheets.
 ...
In fact, the use of knotted cords was an important adaptation to living in the
Andes, one of the most challenging geographies on Earth. Chaskis (Inca
messengers) navigated the steep slopes of the Andes on foot, carrying one of
the world’s most durable and portable envelopes: a khipu draped over
each shoulder.

*/

type TypesettingPipeline struct {
	todo int
}

type Khipukamayuq interface {
	KnotEncode(text string, pipeline *TypesettingPipeline, regs *params.TypesettingRegisters) *Khipu
}

func UAX14LineWrap(text string, regs *params.TypesettingRegisters) *Khipu {
	sread := strings.NewReader(text)     // wrap a reader around the CDATA string
	nfcread := norm.NFC.Reader(sread)    // wrap a normalization-reader around it
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

// https://github.com/blevesearch/segment
/*
A Go library for performing Unicode Text Segmentation as described in
Unicode Standard Annex #29
*/
// Alternativen: https://github.com/go-ego/gse
//
func iterateOverWords(s *strings.Reader) {
	segmenter := segment.NewWordSegmenter(s)
	for segmenter.Segment() {
		tokenBytes := segmenter.Bytes()
		tokenType := segmenter.Type()
		r := string(tokenBytes[:])
		fmt.Printf("%s %v\n", r, tokenType)
	}
}
