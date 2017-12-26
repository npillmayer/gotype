package arithmetic

import (
	"fmt"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	numeric "github.com/shopspring/decimal"
)

/*
BSD License
Copyright (c) 2017, Norbert Pillmayer

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. Neither the name of Tom Everett nor the names of its contributors
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

 * Basic arithmetic objects and interface.

*/

// === Tracing ==========================================================

var T tracing.Trace = tracing.EquationsTracer

// --- Numeric data type ------------------------------------------------

var ConstZero = numeric.Zero      // often used constant 0
var ConstOne = numeric.New(1, 0)  // often used constant 1.0
var MinusOne = numeric.New(-1, 0) // often used constant -1.0

var Deg2Rad, _ = numeric.NewFromString("0.01745329251")

var epsilon numeric.Decimal = numeric.New(1, -6) // numerics below epsilon are considered Zero

// is n Zero ?
func Zero(n numeric.Decimal) bool {
	return n.Abs().LessThanOrEqual(epsilon)
}

// is n = 1.0 ?
func One(n numeric.Decimal) bool {
	return n.Sub(ConstOne).Abs().LessThanOrEqual(epsilon)
}

// Make n = 0 if n "means" to be Zero
func Zap(n numeric.Decimal) numeric.Decimal {
	if Zero(n) {
		n = ConstZero
	}
	return n
}

// Round to epsilon
func Round(n numeric.Decimal) numeric.Decimal {
	return n.Round(7)
}

// --- Pair data type ---------------------------------------------------

type Pair interface {
	fmt.Stringer
	XPart() numeric.Decimal
	YPart() numeric.Decimal
}

type SimplePair struct {
	X numeric.Decimal
	Y numeric.Decimal
}

func MakePair(x, y numeric.Decimal) Pair {
	return SimplePair{
		X: x,
		Y: y,
	}
}

var Origin Pair = MakePair(ConstZero, ConstZero) // often used constant

func (p SimplePair) String() string {
	return fmt.Sprintf("(%s,%s)", p.X.Round(3).String(), p.Y.Round(3).String())
}

func (p SimplePair) XPart() numeric.Decimal {
	return p.X
}

func (p SimplePair) YPart() numeric.Decimal {
	return p.Y
}

func (p SimplePair) Zap() Pair {
	p = SimplePair{
		X: Zap(p.X),
		Y: Zap(p.Y),
	}
	return p
}

func (p SimplePair) Zero() bool {
	return p.Equal(Origin)
}

func (p SimplePair) Equal(p2 Pair) bool {
	p = p.Zap().(SimplePair)
	return p.X.Equal(p2.XPart()) && p.Y.Equal(p2.YPart())
}

func (p SimplePair) Add(p2 Pair) Pair {
	p = SimplePair{
		X: Zap(p.X.Add(p2.XPart())),
		Y: Zap(p.Y.Add(p2.YPart())),
	}
	return p.Zap()
}

func (p SimplePair) Subtract(p2 Pair) Pair {
	p = SimplePair{
		X: Zap(p.X.Sub(p2.XPart())),
		Y: Zap(p.Y.Sub(p2.YPart())),
	}
	return p.Zap()
}

func (p SimplePair) Multiply(n numeric.Decimal) Pair {
	p = SimplePair{
		X: Zap(p.X.Mul(n)),
		Y: Zap(p.Y.Mul(n)),
	}
	return p
}

func (p SimplePair) Divide(n numeric.Decimal) Pair {
	p = SimplePair{
		X: Zap(p.X.Div(n)),
		Y: Zap(p.Y.Div(n)),
	}
	return p
}
