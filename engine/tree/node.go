package tree

/*
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"fmt"
	"sync"
)

// Node is the base type our tree is built of.
type Node struct {
	parent   *Node         // parent node of this node
	children childrenSlice // mutex-protected slice of children nodes
	Payload  interface{}   // nodes may carry a payload of arbitrary type
	Rank     uint32        // rank is used for preserving sequence
}

// NewNode creates a new tree node with a given payload.
func NewNode(payload interface{}) *Node {
	return &Node{Payload: payload}
}

// String is a simple Stringer which the node's Payload packaged in a string.
func (node *Node) String() string {
	return fmt.Sprintf("(Node #ch=%d)", node.ChildCount())
}

// AddChild inserts a new child node into the tree.
// The newly inserted node is connected to this node as its parent.
// It returns the parent node to allow for chaining.
//
// This operation is concurrency-safe.
func (node *Node) AddChild(ch *Node) *Node {
	if ch != nil {
		node.children.addChild(ch, node)
	}
	return node
}

// SetChildAt inserts a new child node into the tree.
// The newly inserted node is connected to this node as its parent.
// The child is set at a given position in relation to other children,
// replacing the child at position i if it exists.
// It returns the parent node to allow for chaining.
//
// This operation is concurrency-safe.
func (node *Node) SetChildAt(i int, ch *Node) *Node {
	if ch != nil {
		node.children.setChild(i, ch, node)
	}
	return node
}

// InsertChildAt inserts a new child node into the tree.
// The newly inserted node is connected to this node as its parent.
// The child is set at a given position in relation to other children,
// shifting children at later positions.
// It returns the parent node to allow for chaining.
//
// This operation is concurrency-safe.
func (node *Node) InsertChildAt(i int, ch *Node) *Node {
	if ch != nil {
		node.children.insertChildAt(i, ch, node)
	}
	return node
}

// Parent returns the parent node or nil (for the root of the tree).
func (node *Node) Parent() *Node {
	return node.parent
}

// Isolate removes a node from its parent.
func (node *Node) Isolate() {
	if node == nil || node.parent == nil {
		return
	}
	node.parent.children.remove(node)
}

// ChildCount returns the number of children-nodes for a node
// (concurrency-safe).
func (node *Node) ChildCount() int {
	return node.children.length()
}

// Child is a concurrency-safe way to get a children-node of a node.
func (node *Node) Child(n int) (*Node, bool) {
	if node.children.length() <= n {
		return nil, false
	}
	return node.children.child(n), true
}

// Children returns a slice with all children of a node.
func (node *Node) Children() []*Node {
	return node.children.asSlice()
}

// IndexOfChild returns the index of a child within the list of children
// of its parent.
func (node *Node) IndexOfChild(ch *Node) int {
	if node.ChildCount() > 0 {
		children := node.Children()
		for i, child := range children {
			if ch == child {
				return i
			}
		}
	}
	return -1
}

// --- Slices of concurrency-safe sets of children ----------------------

type childrenSlice struct {
	sync.RWMutex
	slice []*Node
}

func (chs *childrenSlice) length() int {
	chs.RLock()
	defer chs.RUnlock()
	return len(chs.slice)
}

func (chs *childrenSlice) addChild(child *Node, parent *Node) {
	if child == nil {
		return
	}
	chs.Lock()
	defer chs.Unlock()
	chs.slice = append(chs.slice, child)
	child.parent = parent
}

func (chs *childrenSlice) setChild(i int, child *Node, parent *Node) {
	if child == nil {
		return
	}
	chs.Lock()
	defer chs.Unlock()
	if len(chs.slice) <= i {
		l := len(chs.slice)
		chs.slice = append(chs.slice, make([]*Node, i-l+1)...)
		/*
			c := make([]*Node, i+1)
			for j := 0; j <= i; j++ {
				if j < l {
					c[j] = chs.slice[j]
				} else {
					c[j] = nil
				}
			}
			chs.slice = c
		*/
	}
	chs.slice[i] = child
	child.parent = parent
}

func (chs *childrenSlice) insertChildAt(i int, child *Node, parent *Node) {
	if child == nil {
		return
	}
	chs.Lock()
	defer chs.Unlock()
	if len(chs.slice) <= i {
		l := len(chs.slice)
		chs.slice = append(chs.slice, make([]*Node, i-l+1)...)
	} else {
		chs.slice = append(chs.slice, nil)   // make room for one child
		copy(chs.slice[i+1:], chs.slice[i:]) // shift i+1..n
	}
	chs.slice[i] = child
	child.parent = parent
}

func (chs *childrenSlice) remove(node *Node) {
	chs.Lock()
	defer chs.Unlock()
	for i, ch := range chs.slice {
		if ch == node {
			chs.slice[i] = nil
			node.parent = nil
			break
		}
	}
}

func (chs *childrenSlice) child(n int) *Node {
	if chs.length() == 0 || n < 0 || n >= chs.length() {
		return nil
	}
	chs.RLock()
	defer chs.RUnlock()
	return chs.slice[n]
}

func (chs *childrenSlice) asSlice() []*Node {
	chs.RLock()
	defer chs.RUnlock()
	children := make([]*Node, chs.length())
	for i, ch := range chs.slice {
		children[i] = ch
	}
	return children
}
