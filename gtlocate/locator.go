/*
This is currentyl just a stand-in for a real implementation.

It grows whenever I add some functionality needed for tests. Everything here
is quick and dirty right now.

TODO for fonts:

- use "fontconfig" CLIs (via shell):
https://www.freedesktop.org/wiki/Software/fontconfig/

- use Google fonts API to locate remote fonts

- use some webfont library interface
*/
package gtlocate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/npillmayer/gotype/core/hyphenation"
)

func gtrootdir() string {
	gtroot := os.Getenv("GTROOT")
	if gtroot == "" {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gtroot = os.Getenv("HOME") + "/gotype"
		} else {
			gtroot = gopath + "/src/github.com/npillmayer/gotype/GTROOT"
		}
	}
	return gtroot
}

// Return path for a resource file
func FileResource(item string, typ string) string {
	gtroot := gtrootdir()
	var path string
	switch typ {
	case "lua":
		path = gtroot + "/lib/lua/" + item + ".lua"
	case "font":
		path = filepath.Join(os.Getenv("HOME"), "Library", "Fonts", item)
	case "pattern":
		path = filepath.Join(gtroot, "pattern")
		//path = "/Users/npi/prg/go/gotype/etc/" + item
	}
	return path
}

var dicts map[string]*hyphenation.Dictionary

func Dictionary(loc string) *hyphenation.Dictionary {
	if dicts == nil {
		dicts = make(map[string]*hyphenation.Dictionary)
	}
	if dicts[loc] == nil {
		pname := "hyph-en-us.tex"
		d := hyphenation.LoadPatterns(FileResource(pname, "pattern"))
		dicts[loc] = d
	}
	if dicts[loc] == nil {
		panic(fmt.Sprintf("No dictionnary found for %s", loc))
	}
	return dicts[loc]
}
