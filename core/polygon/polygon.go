/*
Package polygon deals with polygons, ie paths with straight lines.

This is, besides other functionality, a proxy class for a polyclip package from

   https://github.com/akavel/polyclip-go

(See CREDITS file for details)


BSD License

Copyright (c) 2017-20, Norbert Pillmayer (norbert@pillmayer.com)

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package polygon

import (
	"bytes"
	"fmt"

	"github.com/npillmayer/gotype/core/config/gtrace"

	pc "github.com/akavel/polyclip-go"
	a "github.com/npillmayer/gotype/core/arithmetic"
	"github.com/npillmayer/gotype/core/config/tracing"
	dec "github.com/shopspring/decimal"
)

// L is tracing to the syntax tracer.
func L() tracing.Trace {
	return gtrace.SyntaxTracer
}

// === Interface Polygon =====================================================

// Polygon is an interface for immutable polygons.
type Polygon interface {
	IsCycle() bool // is this a cyclic polygon, i.e, is it complete?
	N() int        // number of knots in the polygon
	Pt(int) a.Pair // knot #i modulo N
}

// AsString pretty-prints a polygon.
func AsString(pg Polygon) string {
	var s bytes.Buffer
	start := true
	for i := 0; i < pg.N(); i++ {
		pt := pg.Pt(i)
		if !start {
			s.WriteString(" -- ")
		}
		s.WriteString(fmt.Sprintf("%s", pt))
		start = false
	}
	if pg.IsCycle() {
		s.WriteString(" -- cycle")
	}
	return s.String()
}

// === Polygon Implementation ================================================

// GPPolygon is a concrete implementation of Polygon, based on Go-Polygon's types.
type GPPolygon struct {
	contours pc.Polygon
	isCycle  bool
}

// --- Interface Implementation ----------------------------------------------

// Pt is par of interface Polygon.
func (pg *GPPolygon) Pt(i int) a.Pair {
	i = i % pg.N()
	pt := (*pg.getContour())[i]
	return Pt2Pr(pt)
}

/* Interface Path.
func (pg *Polygon) DeletePoint(i int) {
	c := *pg.getContour()
	c = append(c[:i], c[i+1:]...)
	pg.contours[0] = c
}
*/

// N is part of interface Polygon.
func (pg *GPPolygon) N() int {
	return pg.contours.NumVertices()
}

/* Interface Path.
func (pg *Polygon) Cycle() {
	pg.isCycle = true
}
*/

// IsCycle is part of iterface Polygon.
func (pg *GPPolygon) IsCycle() bool {
	return pg.isCycle
}

// === Builder ===============================================================

// NullPolygon creates an empty polygon.
// This is the starting point for a builder-like functionality. Clients start
// with a null-polygon and subsequently add knots to it.
//
// Example (package qualifiers omitted for clarity and brevity):
//
//    polygon := NullPolygon().Knot(P(1,2)).Knot(P(12,5)).Knot(P(20,3.5)).Cycle()
//
func NullPolygon() *GPPolygon {
	pg := &GPPolygon{}
	return pg
}

// Box creates a box given the top-left and bottom-right corner.
func Box(topleft a.Pair, bottomright a.Pair) *GPPolygon {
	pg := &GPPolygon{}
	pg.Knot(topleft)
	pg.Knot(a.MakePair(bottomright.XPart(), topleft.YPart()))
	pg.Knot(bottomright)
	pg.Knot(a.MakePair(topleft.XPart(), bottomright.YPart()))
	return pg.Cycle()
}

// Internal: create a polygon from a polyclip-polygon
func polygonFromContour(c pc.Polygon) Polygon {
	return &GPPolygon{contours: c, isCycle: true}
}

// Internal: get a polyclip-polygon out of interface Polygon
func getOrMakeContours(pg Polygon) pc.Polygon {
	var contours pc.Polygon
	if p, ok := pg.(*GPPolygon); ok {
		contours = p.contours
	} else {
		p := NullPolygon()
		for i := 0; i < pg.N(); i++ {
			knot := pg.Pt(i)
			p.Knot(knot)
		}
		p.isCycle = true
		contours = p.contours
	}
	return contours
}

func (pg *GPPolygon) getContour() *pc.Contour {
	if len(pg.contours) == 0 {
		pg.contours.Add(pc.Contour{})
	}
	return &pg.contours[0]
}

// Knot appends a knot to a polygon.
func (pg *GPPolygon) Knot(pr a.Pair) *GPPolygon {
	L().Debugf("add knot to polygon: %v", pr)
	pg.getContour().Add(Pr2Pt(pr))
	return pg
}

// AppendSubpath appends a polygon to a polygon.
func (pg *GPPolygon) AppendSubpath(p Polygon) *GPPolygon {
	L().Debugf("add subpath: %s", AsString(p))
	for i := 0; i < p.N(); i++ {
		pg.Knot(p.Pt(i))
	}
	return pg
}

// Cycle closes a polygon.
func (pg *GPPolygon) Cycle() *GPPolygon {
	pg.isCycle = true
	return pg
}

// Copy gets a copy of a polygon.
func (pg *GPPolygon) Copy() *GPPolygon {
	cs := pg.contours.Clone()
	p := &GPPolygon{contours: cs}
	return p
}

// Reverse a polygon. Destructive operation.
func (pg *GPPolygon) Reverse() {
	c := *pg.getContour()
	for i := 0; i < len(c)/2; i++ { // swap points in place
		j := len(c) - i - 1
		c[i], c[j] = c[j], c[i]
	}
}

// Subpath clips out a subpath of a polygon. Destructive operation.
func (pg *GPPolygon) Subpath(from, to int) {
	var contour = *pg.getContour()
	pg.contours[0] = contour[from : to+1]
}

// check assignability
//var _ Path = &Polygon{}
var _ Polygon = &GPPolygon{}

// Stringer
func (pg *GPPolygon) String() string {
	return AsString(pg)
}

// === Operations on Polygons ================================================

// Transform applies an affine transform to (all knots of) a polygon.
// Returns a new polygon (input parameters are unchanged).
func Transform(pg Polygon, t *a.AffineTransform) Polygon {
	ptransformed := NullPolygon()
	for i := 0; i < pg.N(); i++ {
		knot := pg.Pt(i)
		knot = t.Transform(knot)
		ptransformed.Knot(knot)
	}
	ptransformed.isCycle = pg.IsCycle()
	return ptransformed
}

// Union constructs the union of 2 polygons. Returns a new polygon.
func Union(pg1 Polygon, pg2 Polygon) Polygon {
	contour1 := getOrMakeContours(pg1)
	contour2 := getOrMakeContours(pg2)
	pg := contour1.Construct(pc.UNION, contour2)
	return polygonFromContour(pg)
}

// Intersection constructs the intersection of 2 polygons. Returns a new polygon.
func Intersection(pg1 Polygon, pg2 Polygon) Polygon {
	contour1 := getOrMakeContours(pg1)
	contour2 := getOrMakeContours(pg2)
	pg := contour1.Construct(pc.INTERSECTION, contour2)
	return polygonFromContour(pg)
}

// Difference constructs the difference of 2 polygons. Returns a new polygon.
func Difference(pg1 Polygon, pg2 Polygon) Polygon {
	contour1 := getOrMakeContours(pg1)
	contour2 := getOrMakeContours(pg2)
	pg := contour1.Construct(pc.DIFFERENCE, contour2)
	return polygonFromContour(pg)
}

// === Utilities =============================================================

// Pr2Pt is a quick convertion to go-polygon's point type.
func Pr2Pt(pr a.Pair) pc.Point {
	px, _ := pr.XPart().Float64()
	py, _ := pr.YPart().Float64()
	return pc.Point{
		X: px,
		Y: py,
	}
}

// Pt2Pr is a quick convertion from go-polygon's point type.
func Pt2Pr(pt pc.Point) a.Pair {
	px := dec.NewFromFloat(pt.X)
	py := dec.NewFromFloat(pt.Y)
	return a.MakePair(px, py)
}
