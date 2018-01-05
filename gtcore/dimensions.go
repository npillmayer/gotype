package gtcore

import "fmt"

/*
---------------------------------------------------------------------------

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

 * Dimensions and units. Values are in scaled big points (different from TeX).

*/

// A dimension type.
type Dimen int32

var SP Dimen = 1     // scaled point = BP / 65536
var BP Dimen = 65536 // big point = 1/72 inch
var PT Dimen         // printers point
var MM Dimen         // millimeters
var CM Dimen         // centimeters
var PC Dimen         // pica
var CC Dimen         // cicero
var IN Dimen         // inch

// Infinite dimensions
var Fil Dimen = BP * 10000
var Fill Dimen = 2 * BP * 10000
var Filll Dimen = 3 * BP * 10000

// An infinite numeric
var Infty int32 = 100000000

/* Stringer implementation.
 */
func (d Dimen) String() string {
	return fmt.Sprintf("%dsp", int32(d))
}
