//package ebnfparse
package main

import "testing"

func TestLoadGrammar(t *testing.T) {
	g, err := LoadEBNFFromComment("G", ".", "-")
	if err != nil {
		t.Fatalf("Could not load EBNF grammar G from comments in local package")
	}
	t.Logf("Parsed EBNF grammar %s", g.Name)
}
