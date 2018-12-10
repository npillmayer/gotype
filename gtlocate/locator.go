package gtlocate

/*
TODO for fonts:

- use "fontconfig" CLIs (via shell):
https://www.freedesktop.org/wiki/Software/fontconfig/

- use Google fonts API to locate remote fonts

- use some webfont library interface
*/

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/npillmayer/gotype/gtcore/hyphenation"
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
		path = "/Users/npi/prg/go/gotype/etc/hyph-en-us.tex"
	}
	return path
}

var dicts map[string]*hyphenation.Dictionnary

func Dictionnary(loc string) *hyphenation.Dictionnary {
	if dicts == nil {
		dicts = make(map[string]*hyphenation.Dictionnary)
	}
	if dicts[loc] == nil {
		pname := "hyph-en-us.tex"
		d := hyphenation.LoadPatterns(FileResource("pattern", pname))
		dicts[loc] = d
	}
	if dicts[loc] == nil {
		panic(fmt.Sprintf("No dictionnary found for %s", loc))
	}
	return dicts[loc]
}
