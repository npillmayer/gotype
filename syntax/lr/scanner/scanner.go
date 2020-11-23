package scanner

import (
	"fmt"
	"io"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
)

// AnyToken is a helper flag: expect any token from the scanner.
var AnyToken []int = nil

// EOF is identical to text/scanner.EOF.
// Token types are replicated here for practical reasons.
const (
	EOF       = scanner.EOF
	Ident     = scanner.Ident
	Int       = scanner.Int
	Float     = scanner.Float
	Char      = scanner.Char
	String    = scanner.String
	RawString = scanner.RawString
	Comment   = scanner.Comment
)

// Tokenizer is a scanner interface.
type Tokenizer interface {
	NextToken(expected []int) (tokval int, token interface{}, start, len uint64)
	SetErrorHandler(func(error))
}

// DefaultTokenizer is a default implementation, backed by scanner.Scanner.
// Create one with GoTokenizer.
type DefaultTokenizer struct {
	scanner.Scanner
	lastToken    rune        // last token this scanner has produced
	Error        func(error) // error handler
	unifyStrings bool        // convert single chars to strings
}

// Defailt error reporting function for scanners
func logError(e error) {
	gtrace.SyntaxTracer.Errorf(e.Error())
}

// GoTokenizer creates a scanner/tokenizer accepting tokens similar to the Go language.
func GoTokenizer(sourceID string, input io.Reader, opts ...Option) *DefaultTokenizer {
	t := &DefaultTokenizer{}
	t.Error = logError
	t.Init(input)
	t.Filename = sourceID
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// SetErrorHandler sets an error handler for the scanner.
func (t *DefaultTokenizer) SetErrorHandler(h func(error)) {
	if h == nil {
		t.Error = logError
		return
	}
	t.Error = h
}

// NextToken is part of the Tokenizer interface.
func (t *DefaultTokenizer) NextToken(exp []int) (int, interface{}, uint64, uint64) {
	t.lastToken = t.Scan()
	if t.lastToken == scanner.EOF {
		gtrace.SyntaxTracer.Debugf("DefaultTokenizer reached end of input")
	}
	if t.unifyStrings &&
		(t.lastToken == scanner.RawString || t.lastToken == scanner.Char) {
		t.lastToken = scanner.String
	}
	return int(t.lastToken), t.TokenText(), uint64(t.Position.Offset),
		uint64(t.Pos().Offset - t.Position.Offset)
}

// Lexeme is a helper function to receive a string from a token.
func Lexeme(token interface{}) string {
	switch t := token.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprintf("%v", t)
	}
}

// Option configures a default tokenier.
type Option func(p *DefaultTokenizer)

const (
	optionSkipComments uint = 1 << 1 // do not pass comments
	optionUnifyStrings uint = 1 << 2 // treat raw strings and single chars as strings
)

// SkipComments set or clears mode-flag SkipComments.
func SkipComments(b bool) Option {
	return func(t *DefaultTokenizer) {
		if !t.hasmode(optionSkipComments) && b ||
			t.hasmode(optionSkipComments) && !b {
			t.Mode |= scanner.SkipComments
		}
	}
}

// UnifyStrings sets or clears option UnifyStrings:
// treat raw strings and single chars as strings.
func UnifyStrings(b bool) Option {
	return func(t *DefaultTokenizer) {
		t.unifyStrings = b
	}
}

func (t *DefaultTokenizer) hasmode(m uint) bool {
	if m == optionUnifyStrings {
		return t.unifyStrings
	}
	return t.Mode&m > 0
}
