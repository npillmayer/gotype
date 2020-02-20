package bidi

import (
	"bufio"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"golang.org/x/text/unicode/bidi"
)

// Scanner implements the scanner.Tokenizer interface.
// It will read runs of text as a unit, as long as all runes therein have the same Bidi_Class.
type Scanner struct {
	runeScanner *bufio.Scanner // we're using an embedded rune reader
	currClz     bidi.Class     // the current Bidi_Class (of the last rune read)
	lookahead   []byte         // lookahead rune
	buffer      []byte         // character buffer for token lexeme
	level       int            // embedding level
	strong      bidi.Class     // Bidi_Class of last strong character encountered
	pos         uint64         // position in input string
	length      uint64         // length of current lexeme
	done        bool           // at EOF?
	mode        uint           // scanner modes, set with options
}

// NewScanner creates a scanner for bidi formatting. It will read runs of text
// as a unit, as long as all runes therein have the same Bidi_Class.
func NewScanner(input io.Reader, opts ...ScannerOption) *Scanner {
	sc := &Scanner{}
	sc.runeScanner = bufio.NewScanner(input)
	sc.runeScanner.Split(bufio.ScanRunes)
	sc.currClz = bidi.LRI
	sc.buffer = make([]byte, 0, 4096)
	sc.lookahead = make([]byte, 0, 32)
	for _, opt := range opts {
		opt(sc)
	}
	return sc
}

// NextToken reads the next run of input text with identical bidi_class,
// returning a token for it.
//
// The token's value will be set to the bidi_class, the token itself will be
// set to the corresponding input string.
func (sc *Scanner) NextToken(expected []int) (int, interface{}, uint64, uint64) {
	if len(sc.lookahead) > 0 {
		sc.prepareNewRun()
	}
	for sc.runeScanner.Scan() {
		rune := sc.runeScanner.Bytes()
		clz, sz := sc.bidic(rune)
		//T().Debugf("'%s' has class %s", string(rune), ClassString(clz))
		if clz != sc.currClz {
			sc.lookahead = sc.lookahead[:0]
			sc.lookahead = append(sc.lookahead, rune...)
			r := sc.currClz  // tmp for returning current class
			sc.currClz = clz // change current class to class of LA
			T().Debugf("Token '%s' as :%s", string(sc.buffer), ClassString(r))
			return int(r), sc.buffer, sc.pos, uint64(len(sc.buffer))
		}
		sc.buffer = append(sc.buffer, rune...)
		sc.length += uint64(sz)
	}
	if len(sc.lookahead) > 0 {
		// sc.prepareNewRun()
		sc.lookahead = sc.lookahead[:0]
		sc.pos += sc.length            // set new input position
		clz, sz := sc.bidic(sc.buffer) // calculate current bidi class
		sc.currClz = clz
		sc.length += uint64(sz) // include len(LA) in run's length
		T().Debugf("Token '%s' as :%s", string(sc.buffer), ClassString(sc.currClz))
		return int(sc.currClz), sc.buffer, sc.pos, uint64(len(sc.buffer))
	}
	if !sc.done {
		sc.done = true
		T().Debugf("Token :%s", ClassString(bidi.PDI))
		return int(bidi.PDI), "", sc.pos, 0
	}
	return scanner.EOF, "", sc.pos, 0
}

func (sc *Scanner) prepareNewRun() {
	var sz int
	sc.buffer = sc.buffer[:0]                      // reset buffer
	sc.buffer = append(sc.buffer, sc.lookahead...) // move LA to buffer
	sc.pos += sc.length                            // set new input position
	sc.currClz, sz = sc.bidic(sc.buffer)           // calculate current bidi class
	sc.length += uint64(sz)                        // include len(LA) in run's length
}

// bidic returns the Bidi_Class for a rune. It will apply certain UAX#9 rules
// immediately to relief the parser.
func (sc *Scanner) bidic(rune []byte) (bidi.Class, int) {
	r, sz := utf8.DecodeRune(rune)
	if sz > 0 {
		if sc.hasMode(optionTesting) && unicode.IsUpper(r) {
			return bidi.R, sz // during testing, UPPERCASE is R2L
		}
		props, sz := bidi.Lookup(rune)
		clz := props.Class()
		sc.setStrong(clz)
		switch clz { // do some pre-processing
		case bidi.NSM: // rule W1, handle accents
			switch sc.currClz {
			case bidi.LRI:
				return bidi.L, sz
			case bidi.RLI:
				return bidi.R, sz
			case bidi.PDI:
				return bidi.ON, sz
			}
			return sc.currClz, sz
		case bidi.EN: // rule W2 and pretext to W7
			switch sc.strong {
			case bidi.AL:
				return bidi.AN, sz
			case bidi.L:
				return LEN, sz
			}
		case bidi.ON:
			if props.IsBracket() {
				if props.IsOpeningBracket() {
					return LPAREN, sz // TODO different kinds of brackets
				}
				return RPAREN, sz
			}
		case bidi.LRI, bidi.RLI:
			sc.level++
		case bidi.PDI:
			sc.level--
		}
		return props.Class(), sz
	}
	return bidi.L, 0
}

func (sc *Scanner) setStrong(c bidi.Class) bidi.Class {
	switch c {
	case bidi.R, bidi.RLI:
		sc.strong = bidi.R
		return bidi.R
	case bidi.L, bidi.LRI:
		sc.strong = bidi.L
		return bidi.L
	case bidi.AL:
		sc.strong = bidi.AL
		return bidi.AL
	}
	return ILLEGAL
}

// --- Bidi_Class ------------------------------------------------------------

// We use some additional Bidi_Classes, which reflects additional knowledge about
// a character. Our scanner will process some Bidi rules before the parser is
// going to see the tokens.
const (
	LEN     bidi.Class = iota + 100 // left biased european number (EN)
	LPAREN                          // opening bracket
	RPAREN                          // closing bracket
	LSOS                            // start of sequence with direction L
	RSOS                            // start of sequence with direction R
	ILLEGAL bidi.Class = 999        // in-band value denoting illegal class
)

const claszname = "LRENESETANCSBSWSONBNNSMALControlNumLRORLOLRERLEPDFLRIRLIFSIPDI----------"
const claszadd = "LENLPARENRPARENLSOSRSOS"

var claszindex = [...]uint8{0, 1, 2, 4, 6, 8, 10, 12, 13, 14, 16, 18, 20, 23, 25, 32, 35, 38, 41, 44, 47, 50, 53, 56, 59, 62}
var claszaddinx = [...]uint8{0, 3, 9, 15, 19}

// ClassString returns a bidi class as a string.
func ClassString(i bidi.Class) string {
	if i == ILLEGAL {
		return "bidi_class(none)"
	}
	if i >= bidi.Class(len(claszindex)-1) {
		if i >= LEN && i < LEN+bidi.Class(len(claszaddinx)) {
			j := i - LEN
			return claszadd[claszaddinx[j]:claszaddinx[j+1]]
		}
		return "bidi_class(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return claszname[claszindex[i]:claszindex[i+1]]
}

// --- Scanner options -------------------------------------------------------

// ScannerOption configures a bidi scanner
type ScannerOption func(p *Scanner)

const (
	optionRecognizeLegacy uint = 1 << 1 // recognize LRM, RLM, ALM, LRE, RLE, LRO, RLO, PDF
	optionOuterR2L        uint = 1 << 2 // set outer direction as RtoL
	optionTesting         uint = 1 << 3 // test mode: recognize uppercase as class R
)

// RecognizeLegacy sets an option to recognize legacy formatting, i.e.
// LRM, RLM, ALM, LRE, RLE, LRO, RLO, PDF.
func RecognizeLegacy(b bool) ScannerOption {
	return func(sc *Scanner) {
		if !sc.hasMode(optionRecognizeLegacy) && b ||
			sc.hasMode(optionRecognizeLegacy) && !b {
			sc.mode |= optionRecognizeLegacy
		}
	}
}

// Testing will set up the scanner to recognize UPPERCASE letters as having R2L class.
// This is a common pattern in bidi algorithm development.
func Testing(b bool) ScannerOption {
	return func(sc *Scanner) {
		if !sc.hasMode(optionTesting) && b ||
			sc.hasMode(optionTesting) && !b {
			sc.mode |= optionTesting
		}
	}
}

func (sc *Scanner) hasMode(m uint) bool {
	return sc.mode&m > 0
}
