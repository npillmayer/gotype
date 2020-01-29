package scanner

import (
	"fmt"
	"io"
	"text/scanner"
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
}

// DefaultTokenizer is a default implementation, backed by scanner.Scanner.
// Create one with GoTokenizer.
type DefaultTokenizer struct {
	scanner.Scanner
	lastToken rune
}

// GoTokenizer creates a scanner/tokenizer accepting tokens similar to the Go language.
func GoTokenizer(sourceID string, input io.Reader) *DefaultTokenizer {
	t := &DefaultTokenizer{}
	t.Init(input)
	t.Filename = sourceID
	return t
}

// NextToken is part of the Tokenizer interface.
func (t *DefaultTokenizer) NextToken(exp []int) (int, interface{}, uint64, uint64) {
	t.lastToken = t.Scan()
	if t.lastToken == scanner.EOF {
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
