package dom

import (
	"errors"

	"github.com/antchfx/xpath"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
)

/* -----------------------------------------------------------------
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
----------------------------------------------------------------- */

// Provide tree walking for the styled document with XPath syntax.

var errInvalidXPathExpr error = errors.New("Invalid XPath expression")
var errNavigator error = errors.New("Invalid XPath tree navigator")

// XPath is used to select nodes of a styled tree by XPath expressions.
type XPath struct {
	navigator xpath.NodeNavigator
	extractor NodeExtractorFunc
}

// NodeExtracotorFunc is a supporting function to get the current node
// from an xpath.NodeNavigator.
type NodeExtractorFunc func(xpath.NodeNavigator) (style.TreeNode, error)

// NewXPath creates a new XPath walker for a styled tree, with the root
// wrapped into an xpath.NodeNavigator.
// extract is used to extract the current tree node (cursor position) from
// the NodeNavigator. The extract function is required to be able to handle a
// NodeNavigator of the type underlying nav.
func NewXPath(nav xpath.NodeNavigator, extract NodeExtractorFunc) (*XPath, error) {
	xp := &XPath{}
	if nav == nil {
		return nil, errNavigator
	}
	xp.navigator = nav
	return xp, nil
}

// Find searches the styled nodes that match the specified XPath expr.
// Find returns a slice of tree nodes, together with the last non-nil error
// that occured.
func (xp *XPath) Find(xpathexpr string) ([]style.TreeNode, error) {
	var elems []style.TreeNode
	compiled, err := xpath.Compile(xpathexpr)
	if err != nil {
		return nil, errInvalidXPathExpr
	}
	t := compiled.Select(xp.navigator)
	var lasterr error
	for t.MoveNext() {
		node, err := xp.getCurrentNode(t)
		if err != nil {
			lasterr = err
		} else {
			elems = append(elems, node)
		}
	}
	return elems, lasterr
}

// FindOne searches the styled nodes that match the specified XPath expr
// and returns the first match.
func (xp *XPath) FindOne(xpathexpr string) (style.TreeNode, error) {
	compiled, comperr := xpath.Compile(xpathexpr)
	if comperr != nil {
		return nil, errInvalidXPathExpr
	}
	t := compiled.Select(xp.navigator)
	var elem style.TreeNode
	var err error
	if t.MoveNext() {
		elem, err = xp.getCurrentNode(t)
	}
	return elem, err
}

// Each searches the styled tree and calls a callback on each node.
// If brkOnErr (break on error) is set, Each will stop traversing the nodes
// and return immediately, as soon as callback returns an error.
// If brkOnErr is set to false,
// Each will return the last non-nil error returned by the callback function,
// after calling callback for each selected node.
func (xp *XPath) Each(xpathexpr string, callback func(int, style.TreeNode) error, brkOnErr bool) error {
	compiled, err := xpath.Compile(xpathexpr)
	if err != nil {
		return errInvalidXPathExpr
	}
	t := compiled.Select(xp.navigator)
	i := 0
	var lasterr error
	for t.MoveNext() {
		node, err := xp.getCurrentNode(t)
		if err != nil {
			lasterr = err
		} else {
			if cberr := callback(i, node); cberr != nil {
				lasterr = cberr
			}
		}
		if lasterr != nil && brkOnErr {
			break
		}
		i++
	}
	return lasterr
}

// Returns the next styled node pointed to by the navigator.
func (xp *XPath) getCurrentNode(it *xpath.NodeIterator) (style.TreeNode, error) {
	nav := it.Current()
	node, err := xp.extractor(nav)
	return node, err
}

/*
func getCurrentNode(it *xpath.NodeIterator) style.TreeNode {
	n := it.Current()
	if n.NodeType() == xpath.AttributeNode {
		childNode := &html.Node{
			Type: html.TextNode,
			Data: n.Value(),
		}
		return &html.Node{
			Type:       html.ElementNode,
			Data:       n.LocalName(),
			FirstChild: childNode,
			LastChild:  childNode,
		}

	}
	return n.curr
}
*/

// SelectAttr returns the attribute value with the specified name.
/*
func SelectAttr(n *styledtree.Node, name string) (val style.Property) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode && n.Parent == nil && name == n.Data {
		return innerText(n)
	}
	for _, attr := range n.Attr {
		if attr.Key == name {
			val = attr.Val
			break
		}
	}
	return
}
*/
