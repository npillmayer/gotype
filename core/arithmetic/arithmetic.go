/*
Package arithmetic implements basic arithmetic objects.


BSD License

Copyright (c) 2017–20, Norbert Pillmayer

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */
package arithmetic

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/npillmayer/gotype/core/config/gtrace"

	"github.com/npillmayer/gotype/core/config/tracing"
	dec "github.com/shopspring/decimal"
)

// T traces to the equations-tracer.
func T() tracing.Trace {
	return gtrace.EngineTracer
}

// === Numeric Data Type =====================================================

// ConstZero represents a frequently used constant 0.
var ConstZero = dec.Zero

// ConstOne represents a frequently used constant 1.0
var ConstOne = dec.New(1, 0)

// MinusOne represents a frequently used constant -1.0
var MinusOne = dec.New(-1, 0)

// Deg2Rad is a constant for converting from DEG to RAD or vice versa
var Deg2Rad, _ = dec.NewFromString("0.01745329251")

// numerics below epsilon are considered Zero
var epsilon = dec.New(1, -6)

// Zero is a predicate: is n zero ?
func Zero(n dec.Decimal) bool {
	return n.Abs().LessThanOrEqual(epsilon)
}

// One is a predicate: is n = 1.0 ?
func One(n dec.Decimal) bool {
	return n.Sub(ConstOne).Abs().LessThanOrEqual(epsilon)
}

// Zap makes n = 0 if n "means" to be zero
func Zap(n dec.Decimal) dec.Decimal {
	if Zero(n) {
		n = ConstZero
	}
	return n
}

// Round to epsilon.
func Round(n dec.Decimal) dec.Decimal {
	return n.Round(7)
}

// === Pair Data Type ========================================================

// Pair is an interface for pairs / 2D-points
type Pair interface {
	fmt.Stringer
	XPart() dec.Decimal
	YPart() dec.Decimal
}

// SimplePair is a concrete implementation of interface Pair
type SimplePair struct {
	X dec.Decimal
	Y dec.Decimal
}

// Origin represents the frequently used constant (0,0).
var Origin = MakePair(ConstZero, ConstZero)

// --- Constructing Pairs ----------------------------------------------------

// MakePair is a constructor for simple pairs.
func MakePair(x, y dec.Decimal) Pair {
	return SimplePair{
		X: x,
		Y: y,
	}
}

// Pretty Stringer for simple pairs.
func (p SimplePair) String() string {
	return fmt.Sprintf("(%s,%s)", p.X.Round(3).String(), p.Y.Round(3).String())
}

// AsCmplx returns a SimplePair as a complex number.
func (p SimplePair) AsCmplx() complex128 {
	x, _ := p.X.Float64()
	y, _ := p.Y.Float64()
	return complex(x, y)
}

// C2Pr returns a Pair from a complex number.
func C2Pr(c complex128) Pair {
	if cmplx.IsNaN(c) || cmplx.IsInf(c) {
		gtrace.EquationsTracer.Errorf("created pair for complex.NaN")
		return MakePair(ConstZero, ConstZero)
	}
	return Pr(real(c), imag(c))
}

// Pr is a quick notation for contructing a pair from floats.
func Pr(x, y float64) Pair {
	return MakePair(dec.NewFromFloat(x), dec.NewFromFloat(y))
}

// Pr2Pt is a quick notation for getting float values of a pair.
func Pr2Pt(pr Pair) (float64, float64) {
	px, _ := pr.XPart().Float64()
	py, _ := pr.YPart().Float64()
	return px, py
}

// XPart is for interface Pair.
func (p SimplePair) XPart() dec.Decimal {
	return p.X
}

// YPart is for interface Pair.
func (p SimplePair) YPart() dec.Decimal {
	return p.Y
}

// Zap rounds x-part and y-part to epsilon.
func (p SimplePair) Zap() Pair {
	p = SimplePair{
		X: Zap(p.X),
		Y: Zap(p.Y),
	}
	return p
}

// Zero is a predicate: is this pair origin?
func (p SimplePair) Zero() bool {
	return p.Equal(Origin)
}

// Equal compares 2 pairs.
func (p SimplePair) Equal(p2 Pair) bool {
	p = p.Zap().(SimplePair)
	return p.X.Equal(p2.XPart()) && p.Y.Equal(p2.YPart())
}

// Add 2 pairs.
func (p SimplePair) Add(p2 Pair) Pair {
	p = SimplePair{
		X: Zap(p.X.Add(p2.XPart())),
		Y: Zap(p.Y.Add(p2.YPart())),
	}
	return p.Zap()
}

// Subtract 2 pairs.
func (p SimplePair) Subtract(p2 Pair) Pair {
	p = SimplePair{
		X: Zap(p.X.Sub(p2.XPart())),
		Y: Zap(p.Y.Sub(p2.YPart())),
	}
	return p.Zap()
}

// Multiply 2 pairs.
func (p SimplePair) Multiply(n dec.Decimal) Pair {
	p = SimplePair{
		X: Zap(p.X.Mul(n)),
		Y: Zap(p.Y.Mul(n)),
	}
	return p
}

// Divide 2 pairs.
func (p SimplePair) Divide(n dec.Decimal) Pair {
	p = SimplePair{
		X: Zap(p.X.Div(n)),
		Y: Zap(p.Y.Div(n)),
	}
	return p
}

// === Affine Transformations ================================================

// AffineTransform is a matrix type, used for transforming vectors
type AffineTransform struct {
	matrix []dec.Decimal // a 3x3 matrix, flattened by rows
}

// Internal constructor. Clients can use this as a starting point for
// transform combinations.
func newAffineTransform() *AffineTransform {
	m := &AffineTransform{}
	m.matrix = make([]dec.Decimal, 9)
	return m
}

// Identity transform. Will transform a point onto itself.
func Identity() *AffineTransform {
	m := newAffineTransform()
	m.set(0, 0, ConstOne)
	m.set(1, 1, ConstOne)
	m.set(2, 2, ConstOne)
	return m
}

// Translation transform. Translate a point by (dx,dy).
func Translation(pr Pair) *AffineTransform {
	m := Identity()
	m.set(0, 2, pr.XPart())
	m.set(1, 2, pr.YPart())
	return m
}

// Rotation transform. Rotate a point counter-clockwise around the origin.
// Argument is in radians.
func Rotation(theta dec.Decimal) *AffineTransform {
	m := newAffineTransform()
	f, _ := theta.Float64()
	sin := math.Sin(f)
	cos := math.Cos(f)
	m.set(0, 0, dec.NewFromFloat(cos))
	m.set(0, 1, dec.NewFromFloat(-sin))
	m.set(1, 0, dec.NewFromFloat(sin))
	m.set(1, 1, dec.NewFromFloat(cos))
	m.set(2, 2, ConstOne)
	return m
}

// Debug Stringer for an affine transform.
func (m *AffineTransform) String() string {
	s := fmt.Sprintf("[%s,%s,%s|%s,%s,%s|%s,%s,%s]", m.matrix[0], m.matrix[1],
		m.matrix[2], m.matrix[3], m.matrix[4], m.matrix[5], m.matrix[6],
		m.matrix[7], m.matrix[8])
	return s
}

func (m *AffineTransform) get(row, col int) dec.Decimal {
	return m.matrix[row*3+col]
}

func (m *AffineTransform) set(row, col int, value dec.Decimal) {
	m.matrix[row*3+col] = value
}

func (m *AffineTransform) row(row int) []dec.Decimal {
	return m.matrix[row*3 : (row+1)*3]
}

func (m *AffineTransform) col(col int) []dec.Decimal {
	c := make([]dec.Decimal, 3)
	c[0] = m.matrix[col]
	c[1] = m.matrix[3+col]
	c[2] = m.matrix[6+col]
	return c
}

func dotProd(vec1, vec2 []dec.Decimal) dec.Decimal {
	p1 := vec1[0].Mul(vec2[0])
	p2 := vec1[1].Mul(vec2[1])
	p3 := vec1[2].Mul(vec2[2])
	return p1.Add(p2.Add(p3))
}

// Combine 2 affine transformation to a new one. Returns a new transformation
// without changing the argument(s).
func (m *AffineTransform) Combine(n *AffineTransform) *AffineTransform {
	o := newAffineTransform()
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			o.set(row, col, dotProd(n.row(row), m.col(col)))
		}
	}
	return o
}

func (m *AffineTransform) multiplyVector(v []dec.Decimal) []dec.Decimal {
	c := make([]dec.Decimal, 3)
	c[0] = dotProd(m.row(0), v)
	c[1] = dotProd(m.row(1), v)
	c[2] = dotProd(m.row(2), v)
	return c
}

// Transform a 2D-point. The argument is unchanged and a new pair is returned.
func (m *AffineTransform) Transform(pr Pair) Pair {
	c := make([]dec.Decimal, 3)
	c[0] = pr.XPart()
	c[1] = pr.YPart()
	c[2] = ConstOne
	c = m.multiplyVector(c)
	return MakePair(c[0], c[1])
}
