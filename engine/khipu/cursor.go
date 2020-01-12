package khipu

import (
	"errors"
	"fmt"
)

var (
	errorIteratatorEnd = errors.New("Khipu-iterator at end of knot list")
)

// Mark is a type for a position within a khipu.
type Mark interface {
	Position() int
	Knot() Knot
}

type mark struct {
	pos  int
	knot Knot
}

func (m mark) Position() int {
	return m.pos
}

func (m mark) Knot() Knot {
	return m.knot
}

// A Cursor navigates over the knots of a khipu
type Cursor struct {
	khipu *Khipu
	inx   int
}

// NewCursor creates a cursor for a given khipu.
// Usage is unsafe if the referenced khipu changes during lifetime of the cursor.
func NewCursor(kh *Khipu) *Cursor {
	return &Cursor{kh, -1}
}

func (c Cursor) String() string {
	return fmt.Sprintf("[%d]%v", c.inx, c.Knot())
}

// Position returns the current position within the khipu as an integer.
func (c *Cursor) Position() int {
	return c.inx
}

// Next moves the cursor one knot ahead.
// Returns true if the cursor is still at a valid position, false otherwise.
func (c *Cursor) Next() bool {
	c.inx++
	return c.inx < len(c.khipu.knots)
}

// Prev moves the cursor one knot back.
// Returns true if the cursor is still at a valid position, false otherwise.
func (c *Cursor) Prev() bool {
	c.inx--
	return c.inx >= 0
}

// Knot returns the knot at the current position.
func (c *Cursor) Knot() Knot {
	k := c.khipu.knots[c.inx]
	return k
}

// Peek is lookahead 1. Does not advance the cursor.
// Returns true if the lookahead is at a valid position, false otherwise.
func (c *Cursor) Peek() (Knot, bool) {
	if c.IsValidPosition() {
		if c.inx+1 < len(c.khipu.knots) {
			k := c.khipu.knots[c.inx+1]
			return k, true
		}
	}
	return nil, false
}

// Mark returns a mark for the current position/glyph.
func (c *Cursor) Mark() Mark {
	return mark{
		pos:  c.Position(),
		knot: c.Knot(),
	}
}

// Khipu returns the underlying khipu.
// Clients should not modify it, except with the methods of the cursor.
// Different cursors should not modify the same khipu.
// Modifying a khipu may render previously returned marks invalid.
func (c *Cursor) Khipu() *Khipu {
	return c.khipu
}

// IsValidPosition returns true, if the cursor is located at a
// a valid position, false otherwise.
func (c Cursor) IsValidPosition() bool {
	return c.inx >= 0 && c.inx < len(c.khipu.knots)
}

// ReplaceKnot replaces the knot under the cursor, if any.
// It it returns the current knot.
func (c Cursor) ReplaceKnot(knot Knot) Knot {
	if !c.IsValidPosition() {
		return nil
	}
	k := c.Knot()
	c.khipu.knots[c.inx] = knot
	return k
}

// AsGlue returns the current knot as a glue item.
func (c Cursor) AsGlue() Glue {
	return c.Knot().(Glue)
}

// AsPenalty returns the current knot as a penalty.
func (c Cursor) AsPenalty() Penalty {
	return c.Knot().(Penalty)
}

// AsKern returns the current knot as a kern item.
func (c Cursor) AsKern() Kern {
	return c.Knot().(Kern)
}

// AsTextBox returns the current knot as a text box.
func (c Cursor) AsTextBox() *TextBox {
	return c.Knot().(*TextBox)
}

// AsDiscretionary returns the current knot as a discretionary item.
func (c Cursor) AsDiscretionary() *Discretionary {
	return c.Knot().(*Discretionary)
}
