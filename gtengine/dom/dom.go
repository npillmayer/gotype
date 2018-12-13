package dom

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/benbjohnson/css"
	"github.com/npillmayer/gotype/gtcore/dimen"
	"github.com/npillmayer/gotype/gtengine/box"
)

// https://github.com/benbjohnson/css
// This package provides a CSS parser and scanner in pure Go.
// It is an implementation as specified in the W3C's CSS Syntax Module Level 3.

// github.com/andrewstuart/goq
// Package goq was built to allow users to declaratively unmarshal HTML into go
// structs using struct tags composed of css selectors.

// https://de.wikipedia.org/wiki/Mikroformat

// https://github.com/shurcooL/htmlg
// Package for generating and rendering HTML nodes with context-aware escaping.

// https://github.com/stilvoid/please
// Please is a command line utility that makes it easy to integrate web APIs
// into your shell scripts.
// It's called Please because the web works much better if you ask nicely.

// https://gowalker.org/sethwklein.net/go/webutil
// webutil provides Go functions that operate at a level better matching the
// level that I'm working at when I'm using JSON API's and scraping the web.

// https://github.com/tdewolff/minify
// Minify is a minifier package written in Go. It provides HTML5, CSS3, JS,
// JSON, SVG and XML minifiers.

// https://www.smashingmagazine.com/2015/01/designing-for-print-with-css/

type DOM struct {
	doc    *goquery.Document
	design *css.StyleSheet
	//page []media-pages
}

type Page struct {
	size         dimen.Point
	marginsLeft  dimen.Rect
	marginsRight dimen.Rect
	marginBoxes  []*box.StyledBox
}
