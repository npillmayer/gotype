/*
Package w3cdom defines an interface type for W3C Document Object Models.

Status

Early draft—API may change frequently. Please stay patient.


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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */
package w3cdom

import (
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"golang.org/x/net/html"
)

// Node represents W3C-type Node
type Node interface {
	NodeType() html.NodeType  // type of the underlying HTML node (ElementNode, TextNode, etc.)
	NodeName() string         // node name output depends on the node's type
	NodeValue() string        // node value output depends on the node's type
	HasAttributes() bool      // check for existence of attributes
	ParentNode() Node         // get the parent node, if any
	HasChildNodes() bool      // check for existende of sub-nodes
	ChildNodes() NodeList     // get a list of all children-nodes
	Children() NodeList       // get a list of element child-nodes
	FirstChild() Node         // get the first children-node
	NextSibling() Node        // get the Node's next sibling or nil if last
	Attributes() NamedNodeMap // get all attributes of a node
	TextContent() string      // get text from node and all descendents
	ComputedStyles() *style.PropertyMap
}

// NodeList represents W3C-type NodeList
type NodeList interface {
	Length() int
	Item(int) Node
	String() string
}

// Attr represents W3C-type Attr
type Attr interface {
	Namespace() string
	Key() string
	Value() string
}

// NamedNodeMap represents w3C-type NamedNodeMap
type NamedNodeMap interface {
	Length() int
	Item(int) Attr
	GetNamedItem(string) Attr
}
