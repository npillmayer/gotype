package main

import (
	"flag"
	"fmt"
)

var typPtr = flag.String("type", "", "Go type to generate a set type for")
var pkgPtr = flag.String("package", ".", "Go package to scan for type to operate on")

// Generator for converting EBNF into grammar builder calls, reading the EBNF
// definition from comments within Go source files.
func main() {
	flag.Parse()
	flag.Usage()
	//os.Exit(1)
	fmt.Printf("=== Code ===========================================\n\n")
}

// func findType(typ string, pkgname string) bool {
// 	cfg := &packages.Config{Mode: packages.LoadSyntax}
// 	pkgs, err := packages.Load(cfg, pkgname)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "load: %v\n", err)
// 		packages.PrintErrors(pkgs)
// 		return nil, err
// 	}
// }
