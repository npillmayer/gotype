package gtlocate

/*
TODO for fonts:

- use "fontconfig" CLIs (via shell)

- use Google fonts API to locate remote fonts

- use some webfont library interface
*/

import (
	"os"
	"path/filepath"
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
	}
	return path
}
