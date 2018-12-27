/*
Package xpathadapter implements an xpath.NodeNavigator.

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
package xpathadapter

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/antchfx/xpath"
	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"golang.org/x/net/html"
)

type NodeNavigator struct {
	root, current *styledtree.Node
	chinx         int // index into children slice
	attr          int // attributes index
}

// NewNavigator creates a new xpath.NodeNavigator for a styled tree.
func NewNavigator(styledtree *styledtree.Node) *NodeNavigator {
	return &NodeNavigator{
		current: styledtree,
		root:    styledtree,
		attr:    -1,
	}
}

// CurrentNode implements dom.NodeExtractorFunc
func CurrentNode(nav xpath.NodeNavigator) (dom.TreeNode, error) {
	mynav, ok := nav.(*NodeNavigator)
	if !ok {
		return nil, errors.New("Navigator is not of type xpathadapter.NodeNavigator")
	}
	return mynav.current, nil
}

func (nav *NodeNavigator) NodeType() xpath.NodeType {
	switch nav.current.HtmlNode().Type {
	case html.CommentNode:
		return xpath.CommentNode
	case html.TextNode:
		return xpath.TextNode
	case html.DocumentNode:
		return xpath.RootNode
	case html.ElementNode:
		if nav.attr != -1 {
			return xpath.AttributeNode
		}
		return xpath.ElementNode
	case html.DoctypeNode:
		// ignored <!DOCTYPE HTML> declare and as Root-Node type.
		return xpath.RootNode
	}
	panic(fmt.Sprintf("unknown node type: %v", nav.current.HtmlNode().Type))
}

func (nav *NodeNavigator) LocalName() string {
	if nav.attr != -1 {
		return nav.current.HtmlNode().Attr[nav.attr].Key
	}
	return nav.current.HtmlNode().Data
}

func (*NodeNavigator) Prefix() string {
	return ""
}

func (nav *NodeNavigator) Value() string {
	switch nav.current.HtmlNode().Type {
	case html.CommentNode:
		return "<comment nodes not supported>"
	case html.ElementNode:
		if nav.attr != -1 {
			return nav.current.HtmlNode().Attr[nav.attr].Val
		}
		return innerText(nav.current.HtmlNode())
	case html.TextNode:
		return nav.current.HtmlNode().Data
	}
	return ""
}

func (nav *NodeNavigator) Copy() xpath.NodeNavigator {
	n := *nav
	return &n
}

func (nav *NodeNavigator) MoveToRoot() {
	nav.current = nav.root
}

func (nav *NodeNavigator) MoveToParent() bool {
	if nav.attr != -1 {
		nav.attr = -1 // move from attributes to element
		return true
	}
	if nav.current == nav.root {
		return false
	}
	nav.current = nav.current.Parent
	nav.chinx = 0
	return true
}

func (nav *NodeNavigator) MoveToNextAttribute() bool {
	if nav.attr >= len(nav.current.HtmlNode().Attr)-1 {
		return false
	}
	nav.attr++
	return true
}

func (nav *NodeNavigator) MoveToChild() bool {
	if nav.attr != -1 {
		return false
	}
	if nav.current.ChildCount() == 0 {
		return false
	}
	nav.chinx = 0
	ok := false
	nav.current, ok = nav.current.Child(0)
	return ok
}

func (nav *NodeNavigator) MoveToFirst() bool {
	if nav.attr != -1 || nav.chinx == 0 {
		return false
	}
	nav.chinx = 0
	ok := true
	nav.current, ok = nav.current.Parent.Child(0)
	return ok
}

func (nav *NodeNavigator) String() string {
	return nav.Value()
}

func (nav *NodeNavigator) MoveToNext() bool {
	if nav.attr != -1 {
		return false
	}
	if nav.chinx < nav.current.Parent.ChildCount()-1 {
		nav.chinx++
		nav.current = nav.current.Parent.Children[nav.chinx]
		return true
	}
	return false
}

func (nav *NodeNavigator) MoveToPrevious() bool {
	if nav.attr != -1 {
		return false
	}
	if nav.chinx > 0 {
		nav.chinx--
		nav.current = nav.current.Parent.Children[nav.chinx]
		return true
	}
	return false
}

func (nav *NodeNavigator) MoveTo(other xpath.NodeNavigator) bool {
	n, ok := other.(*NodeNavigator)
	if !ok || n.root != nav.root {
		return false
	}

	nav.current = n.current
	nav.attr = n.attr
	return true
}

var _ xpath.NodeNavigator = &NodeNavigator{}
var _ dom.TreeNode = &styledtree.Node{}

// InnerText returns the text between the start and end tags of the object.
func innerText(n *html.Node) string {
	var output func(*bytes.Buffer, *html.Node)
	output = func(buf *bytes.Buffer, n *html.Node) {
		switch n.Type {
		case html.TextNode:
			buf.WriteString(n.Data)
			return
		case html.CommentNode:
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			output(buf, child)
		}
	}

	var buf bytes.Buffer
	output(&buf, n)
	return buf.String()
}
