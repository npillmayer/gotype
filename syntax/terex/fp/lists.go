package fp

/*
BSD License

Copyright (c) 2019–20, Norbert Pillmayer

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

import (
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/syntax/terex"
)

/*
Note:
=====
The current implementation always pre-fetches the first value.
This could be optimized. It would be a problem with long-running ops in the
atom-creation, in case the value is never fetched by an output call.
For now, we will leave it this way.
*/

// ListSeq is a sequence on TeREx lists.
// It moves over the atoms of concrete or virtual lists.
type ListSeq struct {
	atom terex.Atom
	seq  ListGenerator
}

// Seq wraps a TeREx list into a sequence.
func Seq(l *terex.GCons) ListSeq {
	var S ListGenerator
	S = func() ListSeq {
		if l == nil {
			return ListSeq{terex.NilAtom, nil}
		}
		atom := l.Car
		l = l.Cdr
		return ListSeq{atom, S}
	}
	atom := l.Car
	return ListSeq{atom, S}
}

// Break signals a sequene to stop iterating.
func (seq *ListSeq) Break() {
	seq.seq = nil
}

// Done returns true if a sequence stopped iterating.
func (seq *ListSeq) Done() bool {
	return seq.seq == nil
}

// First returns the first atom of a list, together with a sequence successor.
func (seq ListSeq) First() (terex.Atom, ListSeq) {
	return seq.atom, seq
}

// Next returns the next atom of a list-sequence.
func (seq *ListSeq) Next() terex.Atom {
	if seq.Done() {
		return terex.NilAtom
	}
	next := seq.seq()
	seq.atom = next.atom
	if seq.atom == terex.NilAtom {
		seq.seq = nil
	} else {
		seq.seq = next.seq
	}
	return seq.atom
}

// ListGenerator is a function type to generate a list.
type ListGenerator func() ListSeq

// NSeq is an infinite sequence over whole number 0...
func NSeq() ListSeq {
	var n int64
	var S ListGenerator
	S = func() ListSeq {
		n++
		atom := terex.Atomize(n)
		return ListSeq{atom, S}
	}
	atom := terex.Atomize(n)
	return ListSeq{atom, S}
}

// A ListMapper represents an operation on an atom, resulting in a modified atom.
type ListMapper func(terex.Atom) terex.Atom

// Map creates new values from elements/atoms in a list.
func (seq ListSeq) Map(mapper ListMapper) ListSeq {
	var F ListGenerator
	//inner := seq
	atom, inner := seq.atom, seq
	//n, inner := seq.First()
	v := mapper(atom)
	F = func() ListSeq {
		//fmt.Printf("F  called, n=%d\n", n)
		atom = inner.Next()
		v = mapper(atom)
		//fmt.Printf("F' n=%d, v=%d\n", n, v)
		return ListSeq{v, F}
	}
	return ListSeq{v, F}
}

// List returns all the atoms of a sequence as an instantiated list.
func (seq ListSeq) List() *terex.GCons {
	if seq.Done() {
		return nil
	}
	var start, end *terex.GCons
	S := seq
	for atom := seq.Next(); !S.Done(); atom = S.Next() {
		//fmt.Printf("next atom=%s, S=%v\n", atom, S)
		if start == nil {
			start = terex.Cons(atom, nil)
			end = start
		} else {
			end.Cdr = terex.Cons(atom, nil)
			end = end.Cdr
		}
		//fmt.Printf("result list = %s\n", start.ListString())
	}
	return start
}

// --- Trees -----------------------------------------------------------------

// TreeSeq is a type which represents a tree walk as a sequence.
type TreeSeq struct {
	node    TreeNode
	channel <-chan TreeNode
	seq     TreeGenerator
}

// A TreeNode represents a homogenous tree node. Its parent node is available with a
// call to Parent().
type TreeNode struct {
	Node   *terex.GCons
	parent *terex.GCons
	UData  interface{}
}

// internal shortcut for creating a node
func node(node *terex.GCons, parent *terex.GCons) TreeNode {
	return TreeNode{Node: node, parent: parent}
}

// Parent returns the parent of a tree node
func (n TreeNode) Parent() *terex.GCons {
	return n.parent
}

// ReplaceWith replaces a node with a new node, altering the parent node (if present).
func (n TreeNode) ReplaceWith(new *terex.GCons) TreeNode {
	if n.Node == nil {
		if n.parent == nil {
			return node(new, nil)
		}
		panic("inconsistent parent-node combination; node is nil")
	}
	l, r := children(n.parent) // will return (nil,nil) for parent=nil
	if l == n.Node {
		n.parent.Car.Data = new // therefore this is safe
	} else if r == n.Node {
		n.parent.Cdr = new // and this, too
	} else {
		panic("inconsistent parent-node combination; node is no valid child")
	}
	return node(new, n.parent)
}

func (n TreeNode) String() string {
	if n.Node == nil {
		return "<nil>"
	}
	return n.Node.String()
}

// TreeGenerator is a generator function type to iterate over trees.
type TreeGenerator func() TreeSeq

type treeTraverser []*terex.GCons

func (t treeTraverser) tos() *terex.GCons {
	if len(t) > 0 {
		return t[len(t)-1]
	}
	return nil
}

// Traverse creates a sequence from a TeREx tree structure. The sequence traverses the
// tree in depth-first post-order. Internally it uses a goroutine to produce the sequence
// of nodes, receiving them in a channel.
//
// Warning: Currently a goroutine will leak if not all of the nodes of the list are fetched
// by the client.
func Traverse(l *terex.GCons) TreeSeq {
	channel := TreeIteratorCh(l)
	if channel == nil {
		return TreeSeq{}
	}
	var T TreeGenerator
	T = func() TreeSeq {
		var ok bool
		tseq := TreeSeq{node(nil, nil), channel, T}
		if tseq.node, ok = <-channel; !ok {
			tseq.seq = nil
		}
		return tseq
	}
	var ok bool
	var node TreeNode
	tseq := TreeSeq{node, channel, T}
	if tseq.node, ok = <-channel; !ok {
		tseq.seq = nil
	}
	return tseq
}

/*
TreeIteratorCh creates a goroutine and a channel to produce a sequence of nodes from
a depth-first tree walk.

https://www.geeksforgeeks.org/iterative-postorder-traversal-using-stack/

	1.1 Create an empty stack
	2.1 Do following while root is not NULL
		a) Push root's right child and then root to stack.
		b) Set root as root's left child.
	2.2 Pop an item from stack and set it as root.
		a) If the popped item has a right child and the right child
		is at top of stack, then remove the right child from stack,
		push the root back and set root as root's right child.
		b) Else print root's data and set root as NULL.
	2.3 Repeat steps 2.1 and 2.2 while stack is not empty.

For TeREx' pre-order, a node's content is Car, left child is Cdar, right child is Cddr.
The tree from the above website's example

    1---2---4
    |   \---5
    \---3---6
	    \---7

is

	(1 (2 (4) 5) 3 (6) 7)

in TeREx pre-order format. A depth-first traversal will yield

	(4 5 2 6 7 3 1)

*/
func TreeIteratorCh(l *terex.GCons) <-chan TreeNode {
	// 1.1 Create an empty stack
	t := treeTraverser(make([]*terex.GCons, 0, 32))
	if l == nil {
		return nil
	}
	channel := make(chan TreeNode)
	go func(l *terex.GCons) {
		defer close(channel)
		root := l // set root
		for {
			// 2.1 Do following while root is not NULL
			for root != nil {
				left, right := children(root)
				// a) Push root's right child and then root to stack.
				if right != nil {
					t = append(t, right) // push right child node
				}
				t = append(t, root) // push root
				// b) Set root as root's left child.
				root = left
			} // now root == nil
			// 2.2 Pop an item from stack and set it as root.
			root, t = t[len(t)-1], t[:len(t)-1]
			_, right := children(root)
			if len(t) > 0 && right != nil && right == t[len(t)-1] {
				// a) If the popped item has a right child and the right child
				// is at top of stack, then remove the right child from stack,
				// push the root back and set root as root's right child.
				t = t[:len(t)-1]    // pop right child
				t = append(t, root) // push root
				root = right        // root <- right child
			} else {
				// b) Else print root's data and set root as NULL.
				gtrace.SyntaxTracer.Debugf("Node=%s, parent=%s", root, t.tos())
				channel <- node(root, t.tos())
				root = nil
			}
			// 2.3 Repeat steps 2.1 and 2.2 while stack is not empty.
			if len(t) == 0 {
				break
			}
		}
	}(l)
	return channel
}

func children(node *terex.GCons) (*terex.GCons, *terex.GCons) {
	if node == nil {
		return nil, nil
	}
	if node.Car.Type() == terex.ConsType {
		// anonymous node
		panic("anonymous nodes not yet implemented")
	}
	if node.Cdr == nil {
		return nil, nil
	}
	left := node.Cdr.Tee()
	right := node.Cddr()
	return left, right
}

func (t treeTraverser) printStack() {
	for i, n := range t {
		gtrace.SyntaxTracer.Debugf("   [%d] %s", i, terex.Elem(n).String())
	}
}

// Break stops a traversing sequence.
func (seq *TreeSeq) Break() {
	seq.seq = nil
}

// Done returns true if a traversing sequence is stopped.
func (seq *TreeSeq) Done() bool {
	return seq.seq == nil
}

// First returns the first node of a tree traversal.
func (seq TreeSeq) First() (TreeNode, TreeSeq) {
	return seq.node, seq
}

// Next returns the next node of a tree traversal.
func (seq *TreeSeq) Next() TreeNode {
	if seq.Done() {
		return node(nil, nil)
	}
	next := seq.seq()
	node := next.node
	seq.seq = next.seq
	return node
}

// List returns all the nodes of a tree walk as a instantiated list.
func (seq TreeSeq) List() *terex.GCons {
	if seq.Done() {
		return nil
	}
	var start, end *terex.GCons
	for node, T := seq.First(); !T.Done(); node = T.Next() {
		if start == nil {
			start = terex.Cons(node.Node.Car, nil)
			end = start
		} else {
			end.Cdr = terex.Cons(node.Node.Car, nil)
			end = end.Cdr
		}
	}
	return start
}

// A NodeFilter filteres nodes from a sequence of tree traversal nodes.
type NodeFilter func(node TreeNode) bool

// IsLeaf is a filter for tree nodes which only accepts leaf nodes.
func IsLeaf() NodeFilter {
	return func(node TreeNode) bool {
		l, r := children(node.Node)
		return l == nil && r == nil
	}
}

// Where applies a filter to a sequence of integers.
func (seq TreeSeq) Where(filt NodeFilter) TreeSeq {
	var T TreeGenerator
	node, inner := seq.node, seq
	T = func() TreeSeq {
		node = inner.Next()
		for !inner.Done() && !filt(node) {
			node = inner.Next()
		}
		if inner.Done() {
			return TreeSeq{node, nil, nil}
		}
		return TreeSeq{node, nil, T}
	}
	return TreeSeq{node, nil, T}
}

// NodeMapper is a function returning an integer from an input integer.
type NodeMapper func(node TreeNode) TreeNode

// Print prints a node to the syntax tracer and returns the input node.
func Print() NodeMapper {
	return func(node TreeNode) TreeNode {
		gtrace.SyntaxTracer.Debugf("tree node = %s", node)
		return node
	}
}

// Map applies a mapper to all elements of an integer sequence.
func (seq TreeSeq) Map(mapper NodeMapper) TreeSeq {
	var T TreeGenerator
	node, inner := seq.node, seq
	v := mapper(node)
	T = func() TreeSeq {
		node = inner.Next()
		if !inner.Done() {
			v = mapper(node)
			return TreeSeq{v, nil, T}
		}
		return TreeSeq{node, nil, nil}
	}
	return TreeSeq{v, nil, T}
}

func (seq TreeSeq) Range() (int, TreeNode) {
	node := seq.node.Node
	if pattern.Match(seq.node.Node, env) {
		return NodeEnv{
			node:   node,
			parent: seq.node.Parent(),
			env:    env,
		}
	}
	return NodeEnv{}
}
