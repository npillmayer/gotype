package douceuradapter

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var myhtml = `
<html><head>
<style>
  body { border-color: red; }
</style>
</head><body>
  <p>The quick brown fox jumps over the lazy dog.</p>
  <p id="world">Hello <b>World</b>!</p>
  <p style="padding-left: 5px;">This is a test.</p>
</body>
`

func TestExtract1(t *testing.T) {
	h, errhtml := html.Parse(strings.NewReader(myhtml))
	if errhtml != nil {
		t.Error(errhtml)
	}
	css := ExtractStyleElements(h)
	if len(css) != 1 {
		t.Error("Should extract 1 stylesheet")
	}
}
