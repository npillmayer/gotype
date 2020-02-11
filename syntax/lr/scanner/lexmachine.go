package scanner

import (
	"strings"

	"github.com/npillmayer/gotype/core/config/gtrace"

	lex "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

// lexmachine adapter

// LMAdapter is a lexmachine adapter to use lexmachine as a scanner.
type LMAdapter struct {
	Lexer *lex.Lexer
}

// NewLMAdapter creates a new lexmachine adapter. It receives a list of
// literals ('[', ';', …), a list of keywords ("if", "for", …) and a
// map for translating token strings to their values.
//
// NewLMAdapter will return an error if compiling the DFA failed.
func NewLMAdapter(init func(*lex.Lexer), literals []string, keywords []string, tokenIds map[string]int) (*LMAdapter, error) {
	adapter := &LMAdapter{}
	adapter.Lexer = lex.NewLexer()
	init(adapter.Lexer)
	for _, lit := range literals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		gtrace.SyntaxTracer.Debugf("adding literal %s", r)
		adapter.Lexer.Add([]byte(r), MakeToken(lit, tokenIds[lit]))
	}
	for _, name := range keywords {
		adapter.Lexer.Add([]byte(strings.ToLower(name)), MakeToken(name, tokenIds[name]))
	}
	if err := adapter.Lexer.Compile(); err != nil {
		gtrace.SyntaxTracer.Errorf("Error compiling DFA: %v", err)
		return nil, err
	}
	return adapter, nil
}

// Scanner creates a scanner for a given input. The scanner will implement the
// Tokenizer interface.
func (lm *LMAdapter) Scanner(input string) (LMScanner, error) {
	s, err := lm.Lexer.Scanner([]byte(input))
	if err != nil {
		return LMScanner{}, err
	}
	return LMScanner{s}, nil
}

// LMScanner is a scanner type for lexmachine scanners, implementing the
// Tokenizer interface.
type LMScanner struct {
	scanner *lex.Scanner
}

// NextToken is part of the Tokenizer interface.
//
// Warning: The current implementation will ignore the 'expected'-argument.
func (lms *LMScanner) NextToken(expected []int) (int, interface{}, uint64, uint64) {
	tok, err, eof := lms.scanner.Next()
	if err != nil {
		gtrace.SyntaxTracer.Errorf(err.Error())
	}
	if eof {
		return EOF, nil, 0, 0
	}
	token := tok.(*lex.Token)
	tokval := token.Type
	start := uint64(token.StartColumn)
	length := uint64(len(token.Lexeme))
	return tokval, token, start, length
}

// ---------------------------------------------------------------------------

// Skip is a pre-defined action which ignores the scanned match.
func Skip(*lex.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

// MakeToken is a pre-defined action which wraps a scanned match into a token.
func MakeToken(name string, id int) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(id, string(m.Bytes), m), nil
	}
}
