package linebreak

import (
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/khipu"
)

// FixedWidthCursor is a linebreak-cursor for assigning a fixed width to
// letters and spaces.
// It is intended to wrap a khipu.Cursor or another linebreak.Cursor.
type FixedWidthCursor struct {
	cursor     Cursor
	glyphWidth dimen.Dimen
}

var _ Cursor = &FixedWidthCursor{}

// NewFixedWidthCursor creates a FixedWidthCursor, given a width dimension for
// every glyph it will read.
func NewFixedWidthCursor(cursor Cursor, glyphWidth dimen.Dimen) FixedWidthCursor {
	return FixedWidthCursor{
		cursor:     cursor,
		glyphWidth: glyphWidth,
	}
}

// Next is part of interface Cursor.
func (fwc FixedWidthCursor) Next() bool {
	ok := fwc.cursor.Next()
	if ok {
		knot := fwc.cursor.Knot()
		var isChanged bool
		knot, isChanged = setTextDimens(knot, fwc.glyphWidth)
		if isChanged {
			pos := fwc.cursor.Mark().Position()
			fwc.cursor.Khipu().ReplaceKnot(pos, knot)
		}
	}
	return ok
}

// Knot is part of interface Cursor.
func (fwc FixedWidthCursor) Knot() khipu.Knot {
	return fwc.cursor.Knot()
}

// Peek is part of interface Cursor.
func (fwc FixedWidthCursor) Peek() (khipu.Knot, bool) {
	peek, ok := fwc.cursor.Peek()
	if ok {
		peek, _ = setTextDimens(peek, fwc.glyphWidth)
	}
	return peek, ok
}

// Mark is part of interface Cursor.
func (fwc FixedWidthCursor) Mark() khipu.Mark {
	return fwc.cursor.Mark()
}

// Khipu is part of interface Cursor.
func (fwc FixedWidthCursor) Khipu() *khipu.Khipu {
	return fwc.cursor.Khipu()
}

func setTextDimens(knot khipu.Knot, glyphWidth dimen.Dimen) (khipu.Knot, bool) {
	isChanged := false
	switch knot.Type() {
	case khipu.KTDiscretionary:
		d := knot.(khipu.Discretionary)
		isChanged = (d.Width != glyphWidth)
		d.Width = glyphWidth
	case khipu.KTTextBox:
		b := knot.(*khipu.TextBox)
		newW := dimen.Dimen(len(b.Text())) * glyphWidth
		isChanged = (b.Width != newW || b.Height != glyphWidth)
		b.Width = newW
		b.Height = glyphWidth
	case khipu.KTGlue:
		g := knot.(khipu.Glue)
		g[0] = max(1, glyphWidth)
		g[1] = 0
		g[2] = 0
		return g, true
	}
	return knot, isChanged
}

func max(d1, d2 dimen.Dimen) dimen.Dimen {
	if d1 > d2 {
		return d1
	}
	return d2
}
