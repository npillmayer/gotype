package terexlang

import (
	"sync"

	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/timtadh/lexmachine"
)

// The tokens representing literal strings
var literals = []string{"'", "(", ")", "[", "]", "=", "+", "-", "*", "/"}

// The keyword tokens
var keywords = []string{"nil", "t"}

// All of the tokens (including literals and keywords)
var tokens = []string{"COMMENT", "ID", "NUM", "STRING"}

// tokenIds will be set in initTokens()
var tokenIds map[string]int // A map from the token names to their int ids

var initOnce sync.Once // monitors one-time initialization
func initTokens() {
	initOnce.Do(func() {
		var toks []string
		toks = append(toks, tokens...)
		toks = append(toks, keywords...)
		toks = append(toks, literals...)
		tokenIds = make(map[string]int)
		tokenIds["COMMENT"] = scanner.Comment
		tokenIds["ID"] = scanner.Ident
		tokenIds["NUM"] = scanner.Int
		tokenIds["STRING"] = scanner.String
		for i, tok := range toks[4:] {
			tokenIds[tok] = i + 10
		}
	})
}

// Token returns a token name and its value.
func Token(t string) (string, int) {
	return t, tokenIds[t]
}

// Lexer creates a new lexmachine lexer.
func Lexer() (*scanner.LMAdapter, error) {
	initTokens()
	init := func(lexer *lexmachine.Lexer) {
		lexer.Add([]byte(`//[^\n]*\n?`), scanner.Skip)
		lexer.Add([]byte(`\"[^"]*\"`), makeToken("STRING"))
		lexer.Add([]byte(`#?([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_|-)*[!\?]?`), makeToken("ID"))
		lexer.Add([]byte(`[1-9][0-9]*`), makeToken("NUM"))
		lexer.Add([]byte(`( |\,|\t|\n|\r)+`), scanner.Skip)
	}
	adapter, err := scanner.NewLMAdapter(init, literals, keywords, tokenIds)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

func makeToken(s string) lexmachine.Action {
	return scanner.MakeToken(s, tokenIds[s])
}
