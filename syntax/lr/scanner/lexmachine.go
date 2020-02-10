package scanner

import (
	"strings"
	"sync"

	"github.com/npillmayer/gotype/core/config/gtrace"

	lex "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

// lexmachine adapter

var Literals []string       // The tokens representing literal strings
var Keywords []string       // The keyword tokens
var Tokens []string         // All of the tokens (including literals and keywords)
var TokenIds map[string]int // A map from the token names to their int ids
var lexer *lex.Lexer

var initOnce sync.Once // monitors one-time initialization

func initTokens() {
	initOnce.Do(func() {
		Literals = []string{
			"(",
			")",
			"[",
			"]",
			"=",
			",",
			"+",
			"-",
			"*",
			"/",
		}
		Keywords = []string{
			"NODE",
			"EDGE",
			"GRAPH",
			"DIGRAPH",
			"SUBGRAPH",
			"STRICT",
		}
		Tokens = []string{
			"COMMENT",
			"ID",
		}
		Tokens = append(Tokens, Keywords...)
		Tokens = append(Tokens, Literals...)
		TokenIds = make(map[string]int)
		for i, tok := range Tokens {
			TokenIds[tok] = i
		}
	})
}

type LMAdapter struct {
	lexer *lex.Lexer
}

func NewLMAdapter() *LMAdapter {
	adapter := &LMAdapter{}
	adapter.lexer = lex.NewLexer()
	for _, lit := range Literals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		adapter.lexer.Add([]byte(r), token(lit))
	}
	for _, name := range Keywords {
		adapter.lexer.Add([]byte(strings.ToLower(name)), token(name))
	}
	adapter.lexer.Add([]byte(`//[^\n]*\n?`), token("COMMENT"))
	adapter.lexer.Add([]byte(`#?([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_|-)*[!\?]?`), token("ID"))
	adapter.lexer.Add([]byte("( |;|\t|\n|\r)+"), skip)
	if err := adapter.lexer.Compile(); err != nil {
		gtrace.SyntaxTracer.Errorf("Error compiling DFA: %v", err)
		panic(err)
	}
	return adapter
}

func (lm *LMAdapter) NextToken(expected []int) (tokval int, token interface{}, start, len uint64) {
	//
}

func skip(*lex.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

func token(name string) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(TokenIds[name], string(m.Bytes), m), nil
	}
}
