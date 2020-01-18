package bidi

import (
	"testing"

	"golang.org/x/text/unicode/bidi"
)

func TestBidiSimple(t *testing.T) {
	text := "בינלאומי!"
	p := &bidi.Paragraph{}
	p.SetBytes([]byte(text))
	if p.Direction() != bidi.RightToLeft {
		t.Errorf("test did not recognize RtL in Hebrew text")
	}
}
