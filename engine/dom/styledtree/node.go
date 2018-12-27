/*
Package styledtree implements a node type for styled trees.

This is an implementation of style.StyledNode.
It is suited to be called from style.CSSOM to create a styled tree
(i.e., a builder type will construct a tree based on types in this package).

This builder constructs a styled tree using type styledtree.Node.

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
package styledtree

// https://github.com/antchfx/xpath       XPath for parser
// https://github.com/antchfx/htmlquery   HTML DOM XPath
// https://github.com/ChrisTrenkamp/goxpath/tree/master/tree
// https://github.com/santhosh-tekuri/xpathparser  XPath parser
//
// https://github.com/beevik/etree        XPath for XML ("easy tree"), does this :-( :
// type Token interface {
//    Parent() *Element
//    // contains filtered or unexported methods
//}
//
// https://github.com/mmbros/treepath     (kind-of-)XPath for tree interface, BROKEN !
//
// https://godoc.org/github.com/jiangmitiao/ebook-go

import (
	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/dom/style"
	"golang.org/x/net/html"
)

// StyledNodes are the building blocks of the styled tree.
type Node struct {
	htmlNode       *html.Node
	computedStyles *style.PropertyMap
	Parent         *Node
	Children       []*Node
}

// NewNodeForHtmlNode creates a new styled node linked to an HTML DOM node.
func NewNodeForHtmlNode(n *html.Node) *Node {
	return &Node{htmlNode: n}
}

// ----------------------------------------------------------------------

// HtmlNode gets the HTML DOM node corresponding to this styled node.
func (sn Node) HtmlNode() *html.Node {
	return sn.htmlNode
}

// ChildCount returns the number of children-nodes for a styled node.
func (sn Node) ChildCount() int {
	return len(sn.Children)
}

// Child is a safe way to get a children-node of a styled node.
func (sn Node) Child(i int) (*Node, bool) {
	if len(sn.Children) <= i {
		return nil, false
	}
	return sn.Children[i], true
}

// ----------------------------------------------------------------------

// Interface style.StyledNode.
func (sn Node) ComputedStyles() *style.PropertyMap {
	return sn.computedStyles
}

// Interface style.StyledNode.
func (sn *Node) SetComputedStyles(styles *style.PropertyMap) {
	sn.computedStyles = styles
}

var _ style.StyledNode = &Node{}
var _ dom.TreeNode = &Node{}

// ----------------------------------------------------------------------

/*
type StyledNodeQuery struct {
	styledTree *Node
	selection  Selection
}

type Selection []*Node

func (sel Selection) First() *Node {
	if len(sel) > 0 {
		return sel[0]
	}
	return nil
}

func QueryFor(sn style.StyledNode) (*StyledNodeQuery, error) {
	node, ok := sn.(*Node)
	if !ok {
		return nil, errors.New("Cannot query unknown type of styled node")
	}
	return &StyledNodeQuery{node, make([]*Node, 0, 10)}, nil
}

func (snq *StyledNodeQuery) FindElement(e string) Selection {
	snq.selection = collect(snq.styledTree, func(n *Node) bool {
		return n.node.Type == html.ElementNode && strings.EqualFold(n.node.Data, e)
	})
	return snq.selection
}

func (snq *StyledNodeQuery) FindStyledNodeFor(htmlNode *html.Node) *Node {
	sn := find(snq.styledTree, func(n *Node) bool {
		return n.node == htmlNode
	})
	if sn != nil {
		snq.selection = []*Node{sn}
		return sn
	}
	snq.selection = nil
	return nil
}

func (snq *StyledNodeQuery) findThisNode(nodeToFind *Node) *Node {
	return find(snq.styledTree, func(n *Node) bool {
		return n == nodeToFind
	})
}

// Helper to find nodes matching a predicate. Currently works recursive.
// Returns a node or nil.
func find(node *Node, matcher func(n *Node) bool) *Node {
	if node == nil {
		return nil
	}
	if matcher(node) {
		return node
	}
	for _, c := range node.children {
		if f := find(c, matcher); f != nil {
			return f
		}
	}
	return nil
}

// Helper to collect nodes matching a predicate. Currently works recursive.
func collect(node *Node, matcher func(n *Node) bool) (sel Selection) {
	if node == nil {
		return nil
	}
	if matcher(node) {
		sel = append(sel, node)
		return
	}
	for _, c := range node.children {
		sel = collect(c, matcher)
	}
	return
}
*/
