// Copyright (C) 2011, Ross Light

package pdfapi

import (
	"bytes"
)

// Text is a PDF text object.  The zero value is an empty text object.
type Text struct {
	buf       bytes.Buffer
	usedFonts map[*Font]bool
	cursor    Point
	currFont  *Font
	fontSize  Unit
	leading   Unit
}

// Text adds a string to the text object.
func (text *Text) AddGlyphs(s string) {
	writeCommand(&text.buf, "Tj", s)
	if widths := getFontWidths(text.currFont.fontIdent); widths != nil {
		text.cursor.X += computeStringWidth(s, widths, text.fontSize)
	}
}

// Cursor returns the current cursor location.
// This is where new glyphs will be positioned.
func (text *Text) Cursor() Point {
	return text.cursor
}

// MoveCursorTo moves the text cursor to a point.
func (text *Text) MoveCursorTo(pt Point) {
	text.cursor.X = pt.X
	text.cursor.Y = pt.Y
	writeCommand(&text.buf, "Td", text.cursor.X, text.cursor.Y)
}

// MoveCursor moves the text cursor by a vector.
func (text *Text) MoveCursor(diff Point) {
	text.cursor.X += diff.X
	text.cursor.Y += diff.Y
	writeCommand(&text.buf, "Td", text.cursor.X, text.cursor.Y)
}

// AdvanceCursor moves the text cursor in x-direction.
func (text *Text) AdvanceCursor(x Unit) {
	text.cursor.X += x
	writeCommand(&text.buf, "Td", text.cursor.X, text.cursor.Y)
}

// SetFont changes the current font to a standard font.
func (text *Text) SetFont(font *Font, size Unit) {
	if text.usedFonts == nil {
		text.usedFonts = make(map[*Font]bool)
	}
	text.usedFonts[font] = true
	text.currFont, text.fontSize = font, size
	writeCommand(&text.buf, "Tf", font.fontIdent, size)
}

// SetLeading changes the amount of space between lines.
/*
func (text *Text) SetLeading(leading Unit) {
	writeCommand(&text.buf, "TL", leading)
	text.leading = leading
}
*/

// NextLine advances the current text position to the next line, based on the
// current leading.
/*
func (text *Text) NextLine() {
	writeCommand(&text.buf, "T*")
	text.x = 0
	text.y -= text.currLeading
}
*/

// NextLineOffset moves the current text position to an offset relative to the
// beginning of the line.
/*
func (text *Text) NextLineOffset(tx, ty Unit) {
	writeCommand(&text.buf, "Td", tx, ty)
	text.x = tx
	text.y += ty
}
*/

// --- Fonts ------------------------------------------------------------

// Font is a type representing fonts in PDF documents.
type Font struct {
	Name      string
	fontIdent name
	data      []byte
	format    FontFormat
}

type FontFormat int8

const (
	unsupported FontFormat = iota
	Standard
	Type1
	Type3
	TrueType_Mac
	TrueType_Win
	OpenType
	WOFF
	WOFF2
)

// NewInternalFont creates a font for one of the PDF standard fonts,
// given a font's name (e.g., "Helveltica").
func NewInternalFont(fname string) *Font {
	fn := &Font{}
	fn.Name = fname
	fn.fontIdent = name(fname)
	fn.format = Type1
	return fn
}

// String returns the Font's indentifier string.
func (fn *Font) String() string {
	return string(fn.fontIdent)
}

// Standard 14 fonts
const (
	Courier            = "Courier"
	CourierBold        = "Courier-Bold"
	CourierOblique     = "Courier-Oblique"
	CourierBoldOblique = "Courier-BoldOblique"

	Helvetica            = "Helvetica"
	HelveticaBold        = "Helvetica-Bold"
	HelveticaOblique     = "Helvetica-Oblique"
	HelveticaBoldOblique = "Helvetica-BoldOblique"

	Symbol = "Symbol"

	Times           = "Times-Roman"
	TimesBold       = "Times-Bold"
	TimesItalic     = "Times-Italic"
	TimesBoldItalic = "Times-BoldItalic"

	ZapfDingbats = "ZapfDingbats"
)

func getFontWidths(fontIdent name) []uint16 {
	switch fontIdent {
	case Courier:
		return courierWidths
	case CourierBold:
		return courierBoldWidths
	case CourierOblique:
		return courierObliqueWidths
	case CourierBoldOblique:
		return courierBoldObliqueWidths
	case Helvetica:
		return helveticaWidths
	case HelveticaBold:
		return helveticaBoldWidths
	case HelveticaOblique:
		return helveticaObliqueWidths
	case HelveticaBoldOblique:
		return helveticaBoldObliqueWidths
	case Symbol:
		return symbolWidths
	case Times:
		return timesRomanWidths
	case TimesBold:
		return timesBoldWidths
	case TimesItalic:
		return timesItalicWidths
	case TimesBoldItalic:
		return timesBoldItalicWidths
	case ZapfDingbats:
		return zapfDingbatsWidths
	}
	return nil
}

func computeStringWidth(s string, widths []uint16, fontSize Unit) Unit {
	width := Unit(0)
	for _, r := range s {
		if int(r) < len(widths) {
			width += Unit(widths[r])
		}
	}
	return width * fontSize / 1000
}
