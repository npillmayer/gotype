// Copyright (C) 2011, Ross Light

package pdfapi

import (
	"math"
	"testing"
)

func floatEq(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestTextX(t *testing.T) {
	text := new(Text)
	if text.Cursor().X != 0 {
		t.Errorf("Text does not start at X=0 (got %.5f)", text.Cursor().X)
	}

	text.SetFont(NewInternalFont(Helvetica), 12)

	text.AddGlyphs("Hello!")
	if !floatEq(float64(text.Cursor().X), 30.672, 1e-5) {
		t.Errorf("\"Hello!\" has wrong X (=%.5f) when %.5f is desired", text.Cursor().X, 30.672)
	}

}

func TestTextY(t *testing.T) {
	text := new(Text)
	if text.Cursor().Y != 0 {
		t.Errorf("Text does not start at Y=0 (got %.5f)", text.Cursor().Y)
	}

	text.SetFont(NewInternalFont(Helvetica), 12)
	text.AddGlyphs("Hello!")
	if text.Cursor().Y != 0 {
		t.Errorf("\"Hello!\" changes baseline (got %.5f)", text.Cursor().Y)
	}

}
