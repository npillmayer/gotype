package scanner

import (
	"fmt"
	"strings"
	"testing"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/timtadh/lexmachine"
	lex "github.com/timtadh/lexmachine"
)

var inputStrings = []string{
	"1",
	"1+1",
	"Hello #World",
	`x="mystring" // commented `,
	"1,2,3",
}

var tokenCounts = []int{1, 3, 3, 3, 5}

func TestScan1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	for i, input := range inputStrings {
		t.Logf("------+-----------------+--------")
		reader := strings.NewReader(input)
		name := fmt.Sprintf("input #%d", i)
		scanner := GoTokenizer(name, reader)
		tokval, token, pos, _ := scanner.NextToken(AnyToken)
		count := 0
		for tokval != EOF {
			t.Logf(" %4d | %15s | @%5d", tokval, Lexeme(token), pos)
			tokval, token, pos, _ = scanner.NextToken(AnyToken)
			count++
		}
		if count != tokenCounts[i] {
			t.Errorf("Expected token count for #%d to be %d, is %d", i, tokenCounts[i], count)
		}
	}
	t.Logf("------+-----------------+--------")
}

var lispTokenCounts = []int{1, 3, 2, 3, 3}

func TestLM(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	initTokens()
	init := func(lexer *lexmachine.Lexer) {
		lexer.Add([]byte(`//[^\n]*\n?`), Skip)
		lexer.Add([]byte(`\"[^"]*\"`), MakeToken("STRING", tokenIds["STRING"]))
		lexer.Add([]byte(`#?([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_|-)*[!\?]?`), MakeToken("ID", tokenIds["ID"]))
		lexer.Add([]byte(`[1-9][0-9]*`), MakeToken("NUM", tokenIds["NUM"]))
		lexer.Add([]byte(`( |\,|\t|\n|\r)+`), Skip)
	}
	LM, err := NewLMAdapter(init, literals, keywords, tokenIds)
	if err != nil {
		t.Error(err)
	}
	for i, input := range inputStrings {
		t.Logf("------+-----------------+--------")
		scanner, err := LM.Scanner(input)
		if err != nil {
			t.Error(err)
		}
		tokval, token, pos, _ := scanner.NextToken(AnyToken)
		count := 0
		for tokval != EOF {
			t.Logf(" %4d | %15s | @%5d", tokval, token.(*lex.Token).Lexeme, pos)
			tokval, token, pos, _ = scanner.NextToken(AnyToken)
			count++
		}
		if count != lispTokenCounts[i] {
			t.Errorf("Expected token count for #%d to be %d, is %d", i, lispTokenCounts[i], count)
		}
	}
	t.Logf("------+-----------------+--------")
}

var literals []string       // The tokens representing literal strings
var keywords []string       // The keyword tokens
var tokens []string         // All of the tokens (including literals and keywords)
var tokenIds map[string]int // A map from the token names to their int ids

func initTokens() {
	literals = []string{
		"'",
		"(",
		")",
		"[",
		"]",
		"=",
		"+",
		"-",
		"*",
		"/",
	}
	keywords = []string{
		"nil",
		"t",
	}
	tokens = []string{
		"COMMENT",
		"ID",
		"NUM",
		"STRING",
	}
	tokens = append(tokens, keywords...)
	tokens = append(tokens, literals...)
	tokenIds = make(map[string]int)
	tokenIds["COMMENT"] = scanner.Comment
	tokenIds["ID"] = scanner.Ident
	tokenIds["NUM"] = scanner.Int
	tokenIds["STRING"] = scanner.String
	for i, tok := range tokens[4:] {
		tokenIds[tok] = i + 10
	}
}
