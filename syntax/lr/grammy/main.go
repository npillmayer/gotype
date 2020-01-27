package main

import (
	"flag"
	"fmt"
	"os"
)

var gnamePtr = flag.String("grammar", "G", "Name of a grammar, specified in EBNF")
var pkgPtr = flag.String("package", ".", "Go package to scan for an EBNF grammar")
var target = flag.Bool("parser", false, "Generate code for a parser for the grammar")

// Generator for converting EBNF into grammar builder calls, reading the EBNF
// definition from comments within Go source files.
func main() {
	flag.Parse()
	g, err := LoadEBNFFromComment(*gnamePtr, *pkgPtr)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Printf("Found grammar %s in comments\n", g.Name)
	code, err := GenerateBuilder(g)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
	fmt.Printf("=== Code ===========================================\n%s\n", code)
}
