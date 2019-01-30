package styledtree

/*
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
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
)

// StyledNodes are the building blocks of the styled tree.
type StyNode struct {
	tree.Node      // we build on top of general purpose tree
	htmlNode       *html.Node
	computedStyles *style.PropertyMap
}

// NewNodeForHtmlNode creates a new styled node linked to an HTML DOM node.
func NewNodeForHtmlNode(html *html.Node) *tree.Node {
	sn := &StyNode{}
	sn.Payload = sn // Payload will always reference the node itself
	sn.htmlNode = html
	return &sn.Node
}

// Node gets the styled node from a generic tree node.
func Node(n *tree.Node) *StyNode {
	if n == nil {
		return nil
	}
	return n.Payload.(*StyNode)
}

//
func Creator() style.Creator {
	return creator{}
}

type creator struct{}

func (c creator) ToStyler(n *tree.Node) style.Styler {
	return Node(n)
}

func (c creator) StyleForHtmlNode(htmlnode *html.Node) *tree.Node {
	return NewNodeForHtmlNode(htmlnode)
}

// ----------------------------------------------------------------------

// HtmlNode gets the HTML DOM node corresponding to this styled node.
func (sn StyNode) HtmlNode() *html.Node {
	return sn.Payload.(*StyNode).htmlNode
}

/*
// ParentNode returns the parent node or nil (for the root of the tree).
func (sn Node) ParentNode() *Node {
	return sn.parent
}

// ChildCount returns the number of children-nodes for a styled node.
func (sn Node) ChildCount() int {
	return sn.children.length()
}

// Child is a safe way to get a children-node of a styled node.
func (sn Node) ChildNode(n int) (*Node, bool) {
	if sn.children.length() <= n {
		return nil, false
	}
	return sn.children.child(n), true
}

// AddChild inserts a new child node into the tree.
// The newly inserted node is connected to this node as its parent.
//
// This operation is concurrency-safe.
func (sn *Node) AddChild(ch *Node) {
	if ch == nil {
		return
	}
	sn.children.addChild(ch, sn)
}

// ParentNode returns the parent node of this node, as a style.TreeNode.
//
// Interface style.TreeNode.
func (sn StyNode) ParentNode() style.TreeNode {
	return sn.Parent().Payload.(*StyNode)
}

// ChildNode returns the child at position n, as a style.TreeNode
//
// Interface style.TreeNode.
func (sn StyNode) ChildNode(n int) style.TreeNode {
	ch, _ := sn.Child(n)
	return ch.Payload.(*StyNode)
}
*/

// ----------------------------------------------------------------------

// Interface style.TreeNode.
func (sn StyNode) ComputedStyles() *style.PropertyMap {
	return sn.computedStyles
}

// Interface cssom.StyledNode.
func (sn *StyNode) SetComputedStyles(styles *style.PropertyMap) {
	sn.computedStyles = styles
}

// ----------------------------------------------------------------------

/*
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

func (chs *childrenSlice) child(n int) *Node {
	if chs.length() == 0 || n < 0 || n >= chs.length() {
		return nil
	}
	chs.RLock()
	defer chs.RUnlock()
	return chs.slice[n]
}
*/

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
