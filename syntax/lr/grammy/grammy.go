package main

import (
	"fmt"
	"os"
	"strings"

	"go/ast"

	"golang.org/x/exp/ebnf"
	"golang.org/x/tools/go/packages"
)

// EBNFGrammar is a type for a grammar defined in EBNF.
type EBNFGrammar struct {
	pkgname string
	Name    string
	Start   string
	ebnf    ebnf.Grammar
}

// LoadEBNFFromComment searches for an EBNF grammar in a comments section of
// Go source files, given a package name. If found, will return an instance
// of EBNFGrammar with the grammar parsed.
// If no grammar has been found or parsing an existing grammar resulted in an
// error, an error is returned.
func LoadEBNFFromComment(grammarname string, pkgname string) (*EBNFGrammar, error) {
	cfg := &packages.Config{Mode: packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, pkgname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		packages.PrintErrors(pkgs)
		return nil, err
	}
	for _, pkg := range pkgs {
		for _, ast := range pkg.Syntax {
			for _, commentgroup := range ast.Comments {
				ebnfsnippet, gname, ok := findEBNF(commentgroup)
				if ok && grammarname == gname {
					reader := strings.NewReader(ebnfsnippet)
					g, err := ebnf.Parse(grammarname, reader)
					if err != nil {
						return nil, fmt.Errorf("Could not parse EBNF: %v", err)
					}
					startsym := strings.TrimSpace(ebnfsnippet[:strings.Index(ebnfsnippet, "=")])
					// if err = ebnf.Verify(g, startsym); err != nil {
					// 	return nil, fmt.Errorf("Error verifying grammar: %v", err)
					// }
					return &EBNFGrammar{
						pkgname: pkg.ID,
						Name:    gname,
						Start:   startsym,
						ebnf:    g,
					}, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("No EBNF comments found in package")
}

// findEBNF scans through a comment for an EBNF section.
func findEBNF(cg *ast.CommentGroup) (string, string, bool) {
	directive := cg.List[0].Text
	if strings.HasPrefix(directive, "//grammy ") {
		cg.List[0].Text = "//" // purge the directive from EBNF code
		fmt.Printf("EBNF found:\n%s\n", cg.Text())
		return cg.Text(), parseDirective(directive), true
	}
	return "", "", false
}

// parseDirective searches for "-grammar <Name>" the directive string
// and returns <name>, or "G" if none has been found.
func parseDirective(directive string) string {
	s := strings.Split(directive, " ")
	var grammarname string
	for i, seg := range s {
		if seg == "-grammar" && len(s) > i+1 {
			grammarname = s[i+1]
		}
	}
	if len(grammarname) == 0 {
		return "G"
	}
	return grammarname
}

// ======================================================================
// This is a test comment for the accompanying test cases.
// Please to not edit!

//grammy -grammar G -parser -package github.com/npillmayer/gotype/syntax/lr/ebnfcom
// Production  = name "=" [ Expression ] "." .
// Expression  = Alternative { "|" Alternative } .
// Alternative = Term { Term } .
// Term        = name | token [ "…" token ] | Group | Option | Repetition .
// Group       = "(" Expression ")" .
// Option      = "[" Expression "]" .
// Repetition  = "{" Expression "}" .
