package gtlocate

import "os"

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
	}
	return path
}
