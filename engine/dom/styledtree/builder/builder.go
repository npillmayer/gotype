/*
Package builder implements a builder for styled trees.

Type Builder is an implementation of cssom.TreeBuilder.
It is suited to be called from CSSOM.Style() to create a styled tree.

This builder constructs a styled tree using the default node type, i.e.
styledtree.Node.

BSD License

Copyright (c) 2017–18, Norbert Pillmayer

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
package builder

import (
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"golang.org/x/net/html"
)

// Builder is an implementation of style.TreeBuilder.
type Builder struct{}

// MakeNodeFor creates a new styled node corresponding to an HTML DOM node.
//
// Interface cssom.StyledTreeBuilder.
func (b Builder) MakeNodeFor(n *html.Node) cssom.StyledNode {
	sn := styledtree.NewNodeForHtmlNode(n)
	return sn
}

// LinkNodeToParent attaches a styled node to the tree. ATTENTION:
// Tree construction may be concurrent => this method must be thread-safe!
//
// Will panic if nodes are not of type *styledtree.Node (this builder
// builds trees of this node type, after all).
//
// Interface cssom.StyledTreeBuilder.
func (b Builder) LinkNodeToParent(sn cssom.StyledNode, parent cssom.StyledNode) {
	p, ok := parent.(*styledtree.Node)
	if !ok {
		panic("LinkNodeToParent: cannot link to unknown type of styled node")
	}
	this := sn.(*styledtree.Node)
	p.AddChild(this) // concurrency-safe operation
}

// WalkUpwards walks to parent of node.
//
// Interface cssom.StyledTreeBuilder.
func (b Builder) WalkUpwards(sn cssom.StyledNode) cssom.StyledNode {
	this := sn.(*styledtree.Node)
	return this.Parent().(*styledtree.Node)
}
