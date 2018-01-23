package runtime

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/stacks/linkedliststack"
	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtcore/path"
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
type PathNode struct {
	Symbol Symbol      // a pair- or path-variable
	Path   *path.Path  // a MetaFont-like path
	Pair   arithm.Pair // a 2D-point
}

type PathStack struct {
	stack *linkedliststack.Stack // a stack of paths
}

/* Create a new path stack. It is fully initialized and empty.
 */
func NewPathStack() *PathStack {
	pst := &PathStack{
		stack: linkedliststack.New(), // stack of interface{}
	}
	return pst
}

/* Stack functionality. Will return an invalid nullpath if stack is empty.
 */
func (ps *PathStack) Top() *PathNode {
	tos, ok := ps.stack.Peek()
	if !ok {
		tos = &PathNode{
			Symbol: nil,
			Path:   nil,
			Pair:   nil,
		}
	}
	return tos.(*PathNode)
}

/* Stack functionality.
 */
func (ps *PathStack) Pop() (*PathNode, bool) {
	tos, ok := ps.stack.Pop()
	return tos.(*PathNode), ok
}

/* Stack functionality.
 */
func (ps *PathStack) Push(pn *PathNode) *PathStack {
	ps.stack.Push(pn)
	//T.Debugf("TOS is now %v", ps.Top())
	return ps
}

/* Push a path variable. Both arguments may be nil.
 */
func (ps *PathStack) PushPath(sym Symbol, path *path.Path) {
	pn := &PathNode{
		Symbol: sym,
		Path:   path,
	}
	ps.Push(pn)
}

/* Push a pair variable. Both arguments may be nil.
 */
func (ps *PathStack) PushPair(sym Symbol, pr arithm.Pair) {
	pn := &PathNode{
		Symbol: sym,
		Pair:   pr,
	}
	ps.Push(pn)
}

/* Stack functionality.
 */
func (ps *PathStack) IsEmpty() bool {
	return ps.stack.Empty()
}

/* Stack functionality.
 */
func (ps *PathStack) Size() int {
	return ps.stack.Size()
}

// === Operations on Paths ===================================================

type PathBuilder struct {
	q       *dll.List  // linked list of path fragments
	path    *path.Path // the path to build
	iscycle bool
}

func NewBuilder() *PathBuilder {
	pb := &PathBuilder{}
	pb.q = dll.New()
	return pb
}

func (pb *PathBuilder) CollectKnot(pr arithm.Pair) *PathBuilder {
	pb.q.Prepend(pr)
	return pb
}

func (pb *PathBuilder) CollectSubpath(sp *path.Path) *PathBuilder {
	pb.q.Prepend(sp)
	return pb
}

func (pb *PathBuilder) Cycle() *PathBuilder {
	pb.iscycle = true
	return pb
}

func (pb *PathBuilder) MakePath() *path.Path {
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
				T.Error("strange path fragment detected, ignoring")
			}
		}
	}
	if pb.path.N() > 1 && pb.iscycle {
		pb.path.Cycle()
	}
	T.Infof("new path = %s", path.PathAsString(pb.path, nil))
	return pb.path
}
