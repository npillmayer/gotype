// Package dimen implements dimensions and units.
//
/*
BSD License

Copyright (c) 2017–18, Norbert Pillmayer (norbert@pillmayer.com)

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
*/
package dimen

import "fmt"

// A dimension type.
// Values are in scaled big points (different from TeX).
type Dimen int32

const (
	SP Dimen = 1     // scaled point = BP / 65536
	BP Dimen = 65536 // big point = 1/72 inch
	PT Dimen = 7     // printers point   // TODO
	MM Dimen = 7     // millimeters
	CM Dimen = 7     // centimeters
	PC Dimen = 7     // pica
	CC Dimen = 7     // cicero
	IN Dimen = 7     // inch
)

// Infinite dimensions
const Fil Dimen = BP * 10000
const Fill Dimen = 2 * BP * 10000
const Filll Dimen = 3 * BP * 10000

// An infinite numeric
const Infty int = 100000000 // TODO

// Stringer implementation.
func (d Dimen) String() string {
	return fmt.Sprintf("%dsp", int32(d))
}

// A point on a page
//
// TODO see methods in https://golang.org/pkg/image/#Point
type Point struct {
	X, Y Dimen
}

// A rectangle on a page
type Rect struct {
	TopL, BotR Point
}
