package terexlang

/*
BSD License

Copyright (c) 2019–20, Norbert Pillmayer

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

import (
	"fmt"
	"sync"

	"github.com/npillmayer/gotype/syntax/lr/scanner"
	"github.com/timtadh/lexmachine"
)

// The tokens representing literal one-char lexemes
var literals = []string{"'", "(", ")", "[", "]"}
var ops = []string{"+", "-", "*", "/", "=", "=", "!", "$", "%", "&", "?",
	"<", ">", "≤", "≥", "≠", ".", ",", "^"}

// The keyword tokens
var keywords = []string{"nil", "t"}

// All of the tokens (including literals and keywords)
var tokens = []string{"COMMENT", "ID", "NUM", "STRING", "VAR"}

// tokenIds will be set in initTokens()
var tokenIds map[string]int // A map from the token names to their int ids

var initOnce sync.Once // monitors one-time initialization
func initTokens() {
	initOnce.Do(func() {
		// var toks []string
		// toks = append(toks, tokens...)
		// toks = append(toks, ops...)
		// toks = append(toks, keywords...)
		// toks = append(toks, literals...)
		tokenIds = make(map[string]int)
		tokenIds["COMMENT"] = scanner.Comment
		tokenIds["ID"] = scanner.Ident
		tokenIds["NUM"] = scanner.Float
		tokenIds["STRING"] = scanner.String
		tokenIds["VAR"] = -9
		tokenIds["nil"] = 1
		tokenIds["t"] = 2
		for _, lit := range literals {
			r := lit[0]
			tokenIds[lit] = int(r)
		}
		for _, op := range ops {
			tokenIds[op] = scanner.Ident
		}
	})
}

// Token returns a token name and its value.
func Token(t string) (string, int) {
	id, ok := tokenIds[t]
	if !ok {
		panic(fmt.Errorf("unknown token: %s", t))
	}
	return t, id
}

// Lexer creates a new lexmachine lexer.
func Lexer() (*scanner.LMAdapter, error) {
	initTokens()
	init := func(lexer *lexmachine.Lexer) {
		lexer.Add([]byte(`;[^\n]*\n?`), scanner.Skip) // skip comments
		lexer.Add([]byte(`\"[^"]*\"`), makeToken("STRING"))
		lexer.Add([]byte(`([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_|-)*[!\?]?`), makeToken("ID"))
		lexer.Add([]byte(`#([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_|-)*[!\?]?`), makeToken("VAR"))
		lexer.Add([]byte(`[\+\-]?[0-9]+(\.[0-9]+)?`), makeToken("NUM"))
		lexer.Add([]byte(`( |\,|\t|\n|\r)+`), scanner.Skip)
		//lexer.Add([]byte(`.`), makeToken("ID"))
	}
	adapter, err := scanner.NewLMAdapter(init, append(literals, ops...), keywords, tokenIds)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

func makeToken(s string) lexmachine.Action {
	id, ok := tokenIds[s]
	if !ok {
		panic(fmt.Errorf("unknown token: %s", s))
	}
	return scanner.MakeToken(s, id)
}
