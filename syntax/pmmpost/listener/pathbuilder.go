package listener

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/stacks/linkedliststack"
	arithm "github.com/npillmayer/gotype/core/arithmetic"
	"github.com/npillmayer/gotype/core/path"
	"github.com/npillmayer/gotype/syntax/runtime"
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

 * Quick and dirty sketch of a stack for path fragments. Intended to
 * construct paths from points and shorter paths.

*/

// The type we will push onto the stack
type pathNode struct {
	Symbol runtime.Symbol // a pair- or path-variable
	Path   *path.Path     // a MetaFont-like path
	Pair   arithm.Pair    // a 2D-point
}

type pathStack struct {
	stack *linkedliststack.Stack // a stack of paths
}

/* Create a new path stack. It is fully initialized and empty.
 */
func newPathStack() *pathStack {
	pst := &pathStack{
		stack: linkedliststack.New(), // stack of interface{}
	}
	return pst
}

/* Stack functionality. Will return an invalid nullpath if stack is empty.
 */
func (ps *pathStack) Top() *pathNode {
	tos, ok := ps.stack.Peek()
	if !ok {
		tos = &pathNode{
			Symbol: nil,
			Path:   nil,
			Pair:   nil,
		}
	}
	return tos.(*pathNode)
}

/* Stack functionality.
 */
func (ps *pathStack) Pop() (*pathNode, bool) {
	tos, ok := ps.stack.Pop()
	return tos.(*pathNode), ok
}

/* Stack functionality.
 */
func (ps *pathStack) Push(pn *pathNode) *pathStack {
	ps.stack.Push(pn)
	//T.Debugf("TOS is now %v", ps.Top())
	return ps
}

/* Push a path variable. Both arguments may be nil.
 */
func (ps *pathStack) PushPath(sym runtime.Symbol, path *path.Path) {
	pn := &pathNode{
		Symbol: sym,
		Path:   path,
	}
	ps.Push(pn)
}

/* Push a pair variable. Both arguments may be nil.
 */
func (ps *pathStack) PushPair(sym runtime.Symbol, pr arithm.Pair) {
	pn := &pathNode{
		Symbol: sym,
		Pair:   pr,
	}
	ps.Push(pn)
}

/* Stack functionality.
 */
func (ps *pathStack) IsEmpty() bool {
	return ps.stack.Empty()
}

/* Stack functionality.
 */
func (ps *pathStack) Size() int {
	return ps.stack.Size()
}

// === Operations on Paths ===================================================

type pathBuilder struct {
	q       *dll.List  // linked list of path fragments
	path    *path.Path // the path to build
	iscycle bool
}

func newPathBuilder() *pathBuilder {
	pb := &pathBuilder{}
	pb.q = dll.New()
	return pb
}

func (pb *pathBuilder) CollectKnot(pr arithm.Pair) *pathBuilder {
	pb.q.Prepend(pr)
	return pb
}

func (pb *pathBuilder) CollectSubpath(sp *path.Path) *pathBuilder {
	pb.q.Prepend(sp)
	return pb
}

func (pb *pathBuilder) Cycle() *pathBuilder {
	pb.iscycle = true
	return pb
}

func (pb *pathBuilder) MakePath() *path.Path {
	pb.path = path.Nullpath()
	it := pb.q.Iterator()
	for it.Next() {
		fragment := it.Value()
		if pr, ispair := fragment.(arithm.Pair); ispair { // add pair
			pb.path.Knot(pr)
		} else { // add subpath
			subpath, issubp := fragment.(*path.Path)
			if issubp && subpath != nil {
				pb.path.AppendSubpath(subpath)
			} else {
				T().Errorf("strange path fragment detected, ignoring")
			}
		}
	}
	if pb.path.N() > 1 {
		if pb.iscycle {
			pb.path.Cycle()
		} else {
			pb.path.End()
		}
	}
	T().Infof("new path = %s", path.PathAsString(pb.path, nil))
	return pb.path
}
