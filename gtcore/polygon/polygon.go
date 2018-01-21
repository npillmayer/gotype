/*
Package polygon deals with polygons, i.e. paths with straigt lines.

This is, besides other functionality, a proxy class for a polyclip package from
https://github.com/akavel/polyclip-go.

---------------------------------------------------------------------------

BSD License

Copyright (c) 2017, Norbert Pillmayer (norbert@pillmayer.com)

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer nor the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----------------------------------------------------------------------

*/
package polygon

import (
	"bytes"
	"fmt"

	pc "github.com/akavel/polyclip-go"
	a "github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
	dec "github.com/shopspring/decimal"
)

// We are tracing to the syntax tracer
var L tracing.Trace = tracing.SyntaxTracer

// === Interface Polygon =====================================================

// An interface for immutable polygons.
type Polygon interface {
	IsCycle() bool // is this a cyclic polygon, i.e, is is complete?
	N() int        // number of knots in the polygon
	Pt(int) a.Pair // knot #i modulo N
}

// Pretty print a polygon
func PolygonAsString(pg Polygon) string {
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

// An interface for paths
/*
type Path interface {
	fmt.Stringer
	AddPoint(pr Pair)
	AddSubpath(p Path)
	GetPoint(i int) Pair
	DeletePoint(i int)
	Subpath(from, to int)
	Length() int
	IsCycle() bool
	Cycle()
	Copy() Path
	Reverse()
}
*/

// === Polygon Implementation ================================================

// A concrete implementation of Polygon, based on Go-Polygon's types
type GPPolygon struct {
	contours pc.Polygon
	isCycle  bool
}

// --- Interface Implementation ----------------------------------------------

// Interface Polygon.
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

// Interface Polygon.
func (pg *GPPolygon) N() int {
	return pg.contours.NumVertices()
}

/* Interface Path.
func (pg *Polygon) Cycle() {
	pg.isCycle = true
}
*/

// Interface Polygon.
func (pg *GPPolygon) IsCycle() bool {
	return pg.isCycle
}

// === Builder ===============================================================

// Create an empty polygon.
func NullPolygon() *GPPolygon {
	pg := &GPPolygon{}
	return pg
}

func Box(lt a.Pair, rb a.Pair) *GPPolygon {
	pg := &GPPolygon{}
	pg.Knot(lt)
	pg.Knot(a.MakePair(rb.XPart(), lt.YPart()))
	pg.Knot(rb)
	pg.Knot(a.MakePair(lt.XPart(), rb.YPart()))
	return pg.Cycle()
}

/*
func NewPath(prs []Pair) Path {
	pg := NullPath()
	for _, pr := range prs {
		pg.AddPoint(pr)
	}
	return pg
}
*/

func (pg *GPPolygon) getContour() *pc.Contour {
	if len(pg.contours) == 0 {
		pg.contours.Add(pc.Contour{})
	}
	return &pg.contours[0]
}

// Append a knot to a polygon.
func (pg *GPPolygon) Knot(pr a.Pair) *GPPolygon {
	L.Debugf("add knot to polygon: %v", pr)
	pg.getContour().Add(Pr2Pt(pr))
	return pg
}

// Append a polygon to a polygon.
func (pg *GPPolygon) AppendSubpath(p Polygon) *GPPolygon {
	L.Debugf("add subpath: %s", PolygonAsString(p))
	for i := 0; i < p.N(); i++ {
		pg.Knot(p.Pt(i))
	}
	return pg
}

// Close a polygon.
func (pg *GPPolygon) Cycle() *GPPolygon {
	pg.isCycle = true
	return pg
}

// Get a copy a polygon.
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

// Clip out a subpath of a polygon. Destructive operation.
func (pg *GPPolygon) Subpath(from, to int) {
	var contour pc.Contour = *pg.getContour()
	pg.contours[0] = contour[from : to+1]
}

// check assignability
//var _ Path = &Polygon{}
var _ Polygon = &GPPolygon{}

// === Utilities =============================================================

// A quick convertion to go-polygon's point type
func Pr2Pt(pr a.Pair) pc.Point {
	px, _ := pr.XPart().Float64()
	py, _ := pr.YPart().Float64()
	return pc.Point{
		X: px,
		Y: py,
	}
}

// A quick convertion from go-polygon's point type
func Pt2Pr(pt pc.Point) a.Pair {
	px := dec.NewFromFloat(pt.X)
	py := dec.NewFromFloat(pt.Y)
	return a.MakePair(px, py)
}
