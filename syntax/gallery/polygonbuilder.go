package gallery

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtcore/polygon"
)

/*
----------------------------------------------------------------------

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
3. Neither the name of Norbert Pillmayer or the names of its contributors
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

/*
Helper type and methods for building polygons from pairs and
sub-polygons.
*/
type PolygonBuilder struct {
	q       *dll.List          // linked list of polyg fragments
	polyg   *polygon.GPPolygon // the polygon to build
	iscycle bool
}

/* Stack functionality.
func (ps *PathStack) IsEmpty() bool {
	return ps.stack.Empty()
}
*/

/* Stack functionality.
func (ps *PathStack) Size() int {
	return ps.stack.Size()
}
*/

// Create a new polygon builder, fully initialized.
func NewPolygonBuilder() *PolygonBuilder {
	pb := &PolygonBuilder{}
	pb.q = dll.New()
	return pb
}

func (pb *PolygonBuilder) CollectKnot(pr arithm.Pair) *PolygonBuilder {
	pb.q.Prepend(pr)
	return pb
}

func (pb *PolygonBuilder) CollectSubpolyg(sp polygon.Polygon) *PolygonBuilder {
	pb.q.Prepend(sp)
	return pb
}

func (pb *PolygonBuilder) Cycle() *PolygonBuilder {
	pb.iscycle = true
	return pb
}

func (pb *PolygonBuilder) MakePolygon() polygon.Polygon {
	pb.polyg = polygon.NullPolygon()
	it := pb.q.Iterator()
	for it.Next() {
		fragment := it.Value()
		if pr, ispair := fragment.(arithm.Pair); ispair { // add pair
			pb.polyg.Knot(pr)
		} else { // add subpolyg
			subpolyg, issubp := fragment.(polygon.Polygon)
			if issubp && subpolyg != nil {
				pb.polyg.AppendSubpath(subpolyg)
			} else {
				T.Error("strange polygon fragment detected, ignoring")
			}
		}
	}
	if pb.polyg.N() > 1 && pb.iscycle {
		pb.polyg.Cycle()
	}
	T.Infof("new polygon = %s", polygon.PolygonAsString(pb.polyg))
	return pb.polyg
}
