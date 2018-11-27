package grapheme

import (
	"bytes"
	"fmt"
	"testing"
	u "unicode"

	"github.com/npillmayer/gotype/gtcore/unicode/segment"
)

func TestGraphemeClasses(t *testing.T) {
	c1 := LClass
	if c1.String() != "LClass" {
		t.Errorf("String(LClass) should be 'LClass', is %s", c1)
	}
	SetupGraphemeClasses()
	if !u.Is(Control, '\t') {
		t.Error("<TAB> should be identified as control character")
	}
	hangsyl := '\uac1c'
	if c := GraphemeClassForRune(hangsyl); c != LVClass {
		t.Errorf("Hang syllable GAE should be of class LV, is %s", c)
	}
}

func TestGraphemes1(t *testing.T) {
	onGraphemes := NewBreaker()
	input := bytes.NewReader([]byte("Hello\tWorld"))
	seg := segment.NewSegmenter(onGraphemes)
	seg.Init(input)
	seg.Next()
	fmt.Printf("Next() = %s\n", seg.Text())
}

func TestGraphemes2(t *testing.T) {
	onGraphemes := NewBreaker()
	input := bytes.NewReader([]byte("Hello\tWorld"))
	seg := segment.NewSegmenter(onGraphemes)
	seg.Init(input)
	for seg.Next() {
		fmt.Printf("Next() = %s\n", seg.Text())
	}
	if seg.Err() != nil {
		t.Errorf("segmenter.Next() failed with error: %s", seg.Err())
	}
}
