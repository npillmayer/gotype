// Copyright (C) 2011, Ross Light

package pdfapi

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
)

// writeCommand writes a PDF graphics command.
func writeCommand(w io.Writer, op string, args ...interface{}) error {
	for _, arg := range args {
		// TODO: Use the same buffer for all arguments
		if m, err := marshal(nil, arg); err == nil {
			if _, err := w.Write(append(m, ' ')); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if _, err := io.WriteString(w, op); err != nil {
		return err
	}
	if _, err := w.Write([]byte{'\n'}); err != nil {
		return err
	}
	return nil
}

// Canvas is a two-dimensional drawing region on a single page.  You can obtain
// a canvas once you have created a document.
type Canvas struct {
	doc          *Document
	page         *pageDict
	ref          Reference
	contents     *stream
	imageCounter uint
}

// Document returns the document the canvas is attached to.
func (canvas *Canvas) Document() *Document {
	return canvas.doc
}

// Close flushes the page's stream to the document.  This must be called once
// drawing has completed or else the document will be inconsistent.
func (canvas *Canvas) Close() error {
	return canvas.contents.Close()
}

// Size returns the page's media box (the size of the physical medium).
func (canvas *Canvas) Size() (width, height Unit) {
	mbox := canvas.page.MediaBox
	return mbox.Dx(), mbox.Dy()
}

// SetSize changes the page's media box (the size of the physical medium).
/*
func (canvas *Canvas) SetSize(width, height Unit) {
	canvas.page.MediaBox = Rectangle{Point{0, 0}, Point{width, height}}
}
*/

// CropBox returns the page's crop box.
func (canvas *Canvas) CropBox() Rectangle {
	return canvas.page.CropBox
}

// SetCropBox changes the page's crop box.
/*
func (canvas *Canvas) SetCropBox(crop Rectangle) {
	canvas.page.CropBox = crop
}
*/

// FillStroke fills then strokes the given path.  This operation has the same
// effect as performing a fill then a stroke, but does not repeat the path in
// the file.
func (canvas *Canvas) FillStroke(p *Path) {
	io.Copy(canvas.contents, &p.buf)
	writeCommand(canvas.contents, "B")
}

// Fill paints the area enclosed by the given path using the current fill color.
func (canvas *Canvas) Fill(p *Path) {
	io.Copy(canvas.contents, &p.buf)
	writeCommand(canvas.contents, "f")
}

// Stroke paints a line along the given path using the current stroke color.
func (canvas *Canvas) Stroke(p *Path) {
	io.Copy(canvas.contents, &p.buf)
	writeCommand(canvas.contents, "S")
}

// SetLineWidth changes the stroke width to the given value.
func (canvas *Canvas) SetLineWidth(w Unit) {
	writeCommand(canvas.contents, "w", w)
}

// SetLineDash changes the line dash pattern in the current graphics state.
// Examples:
//
//   c.SetLineDash(0, []Unit{})     // solid line
//   c.SetLineDash(0, []Unit{3})    // 3 units on, 3 units off...
//   c.SetLineDash(0, []Unit{2, 1}) // 2 units on, 1 unit off...
//   c.SetLineDash(1, []Unit{2})    // 1 unit on, 2 units off, 2 units on...
func (canvas *Canvas) SetLineDash(phase Unit, dash []Unit) {
	writeCommand(canvas.contents, "d", dash, phase)
}

// SetFillColor changes the current fill color to the given RGB triple (in device
// RGB space).
func (canvas *Canvas) SetFillColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	rf := float32(r) / float32(0xffff)
	gf := float32(g) / float32(0xffff)
	bf := float32(b) / float32(0xffff)
	writeCommand(canvas.contents, "rg", rf, gf, bf)
}

// SetStrokeColor changes the current stroke color to the given RGB triple (in
// device RGB space).
func (canvas *Canvas) SetStrokeColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	rf := float32(r) / float32(0xffff)
	gf := float32(g) / float32(0xffff)
	bf := float32(b) / float32(0xffff)
	writeCommand(canvas.contents, "RG", rf, gf, bf)
}

// PushState saves a copy of the current graphics state.  The state can later be
// restored using Pop.
func (canvas *Canvas) PushState() {
	writeCommand(canvas.contents, "q")
}

// PopState restores the most recently saved graphics state by popping it from the
// stack.
func (canvas *Canvas) PopState() {
	writeCommand(canvas.contents, "Q")
}

// Transform concatenates a 3x3 matrix with the current transformation matrix.
// See type Transform.
func (canvas *Canvas) Transform(t Transform) {
	writeCommand(canvas.contents, "cm", t[0], t[1], t[2], t[3], t[4], t[5])
}

// DrawText paints a text object onto the canvas.
func (canvas *Canvas) DrawText(text *Text) {
	for font := range text.usedFonts {
		fontName := font.fontIdent
		if _, ok := canvas.page.Resources.Font[fontName]; !ok {
			canvas.page.Resources.Font[fontName] = canvas.doc.standardFont(fontName)
		}
	}
	writeCommand(canvas.contents, "BT")
	io.Copy(canvas.contents, &text.buf)
	writeCommand(canvas.contents, "ET")
}

// DrawImage paints a raster image at the given location and scaled to the
// given dimensions.  If you want to render the same image multiple times in
// the same document, use DrawImageReference.
func (canvas *Canvas) DrawImage(img image.Image, rect Rectangle) {
	canvas.DrawImageReference(canvas.doc.AddImage(img), rect)
}

// DrawImageReference paints the raster image referenced in the document at the
// given location and scaled to the given dimensions.
func (canvas *Canvas) DrawImageReference(ref Reference, rect Rectangle) {
	name := canvas.nextImageName()
	canvas.page.Resources.XObject[name] = ref

	canvas.PushState()
	t := Identity().Shifted(rect.Min).Scaled(float32(rect.Dx()), float32(rect.Dy()))
	canvas.Transform(t)
	//canvas.Transform(float32(rect.Dx()), 0, 0, float32(rect.Dy()), float32(rect.Min.X), float32(rect.Min.Y))
	writeCommand(canvas.contents, "Do", name)
	canvas.PopState()
}

// DrawLine paints a straight line from pt1 to pt2 using the current stroke
// color and line width.
func (canvas *Canvas) DrawLine(pt1, pt2 Point) {
	var path Path
	path.Move(pt1)
	path.Line(pt2)
	canvas.Stroke(&path)
}

const anonymousImageFormat = "__image%d__"

func (canvas *Canvas) nextImageName() name {
	var n name
	for {
		n = name(fmt.Sprintf(anonymousImageFormat, canvas.imageCounter))
		canvas.imageCounter++
		if _, ok := canvas.page.Resources.XObject[n]; !ok {
			break
		}
	}
	return n
}

func (canvas *Canvas) MoveTo(pt Point) {
	writeCommand(canvas.contents, "m", pt.X, pt.Y)
}

// Path is a shape that can be painted on a canvas.  The zero value is an empty
// path.
type Path struct {
	buf bytes.Buffer
	at  Point
}

// MoveTo begins a new subpath by moving the current point to the given location.
func (path *Path) MoveTo(pt Point) {
	path.at = pt
	writeCommand(&path.buf, "m", pt.X, pt.Y)
}

// Move begins a new subpath by moving the current point by a vector.
func (path *Path) Move(pt Point) {
	path.at.X += pt.X
	path.at.Y += pt.Y
	writeCommand(&path.buf, "m", pt.X, pt.Y)
}

// Line appends a line segment from the current point to the given location.
func (path *Path) LineTo(pt Point) {
	path.at = pt
	writeCommand(&path.buf, "l", pt.X, pt.Y)
}

// Line appends a line segment from the current point to the given location.
func (path *Path) Line(pt Point) {
	path.at.X += pt.X
	path.at.Y += pt.Y
	writeCommand(&path.buf, "l", pt.X, pt.Y)
}

// Curve appends a cubic Bezier curve to the path.
func (path *Path) CurveTo(pt1, pt2, pt3 Point) {
	writeCommand(&path.buf, "c", pt1.X, pt1.Y, pt2.X, pt2.Y, pt3.X, pt3.Y)
}

// Rectangle appends a complete rectangle to the path. If corner > 0, the
// corners of the rectangle will be rounded. All four corners will have the
// same radius.
func (path *Path) Rectangle(rect Rectangle, corner Unit) {
	if corner <= 0 {
		writeCommand(&path.buf, "re", rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy())
	} else {
		off := 0.45 * corner
		corner = min(corner, rect.Dy()/2)
		corner = min(corner, rect.Dx()/2)
		path.Move(Point{rect.Max.X - corner, rect.Min.Y})
		path.CurveTo(
			Point{rect.Max.X - corner + off, rect.Min.Y},
			Point{rect.Max.X, rect.Min.Y + corner - off},
			Point{rect.Max.X, rect.Min.Y + corner})
		path.Line(Point{rect.Max.X, rect.Max.Y - corner})
		path.CurveTo(
			Point{rect.Max.X, rect.Max.Y - corner + off},
			Point{rect.Max.X - corner + off, rect.Max.Y},
			Point{rect.Max.X - corner, rect.Max.Y})
		path.Line(Point{rect.Min.X + corner, rect.Max.Y})
		path.CurveTo(
			Point{rect.Min.X + corner - off, rect.Max.Y},
			Point{rect.Min.X, rect.Max.Y - corner + off},
			Point{rect.Min.X, rect.Max.Y - corner})
		path.Line(Point{rect.Min.X, rect.Min.Y + corner})
		path.CurveTo(
			Point{rect.Min.X, rect.Min.Y + corner - off},
			Point{rect.Min.X + corner - off, rect.Min.Y},
			Point{rect.Min.X + corner, rect.Min.Y})
		path.Close()
	}
}

// Close appends a line segment from the current point to the starting point of
// the subpath.
func (path *Path) Close() {
	writeCommand(&path.buf, "h")
}

// ----------------------------------------------------------------------

// Transform represents an affine transformation (on a PDF graphics state).
// The six indices map to values in a 3x3-matrix as shown below:
//
//  ⎛ a b 0 ⎞
//  ⎜ c d 0 ⎥
//  ⎝ e f 1 ⎠
//
// For more information, see Section 8.3.4 of ISO 32000-1.
type Transform [6]float32

// Identity is the neutral transformation. This is the starting point for
// chaining transform operations.
// To create a transformation which shifts a point by vector (a, b) and
// then rotates around origin by 30 degrees:
//
//    vector := Point{a, b}
//    T := Identity().Shifted(vector).Rotated(Deg2Rad(30))
//    mycanvas.Transform(T)
//
func Identity() Transform {
	return Transform{1, 0, 0, 1, 0, 0}
}

// Shifted moves the canvas's coordinates system by the given offset.
func (t Transform) Shifted(by Point) Transform {
	t[4] += float32(by.X)
	t[5] += float32(by.Y)
	return t
}

// Rotate rotates the canvas's coordinate system by a given angle (in radians),
// counterclockwise.
func (t Transform) Rotated(theta float32) Transform {
	tt := Identity()
	s, c := float32(math.Sin(float64(theta))), float32(math.Cos(float64(theta)))
	tt[0] = t[0]*c - t[1]*s
	tt[1] = t[0]*s + t[1]*c
	tt[2] = t[2]*c - t[3]*s
	tt[3] = t[2]*s + t[3]*c
	tt[4] = t[4]*c - t[5]*s
	tt[5] = t[4]*s + t[5]*c
	return tt
}

// Scale multiplies the canvas's coordinate system by the given scalars.
func (t Transform) Scaled(x, y float32) Transform {
	t[0] *= x
	t[3] *= y
	return t
}

// ----------------------------------------------------------------------

func min(u1 Unit, u2 Unit) Unit {
	if u1 < u2 {
		return u1
	}
	return u2
}

// Deg2Rad returns degree angle coordinates to radians.
func Deg2Rad(theta float32) float32 {
	return theta * 0.01745329252
}
