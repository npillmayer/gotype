package parser

import "fmt"

/*
----------------------------------------------------------------------

BSD License

Copyright (c) 2017–2020, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software or the names of its contributors
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
*/

/*
Wir brauchen einen gerichteten Multigraph (DAG mit mehreren geordneten Kanten
zwischen denselben Knoten).

Kanten werden durch Versionierung verändert, sind also selbst versionierte
Objekte. => 2 Maps pro Graph: To und From.

Kanten liegen in einem Array. Geht das? Werden die Kanten von einem Knoten
weg bzw zu einem Knoten hin immer zusammen angelegt? Beides zusammen geht
eigentlich gar nicht. Umsortieren ist zu aufwändig. From würde gehen.
To nicht.

simple/graph legt eine Map von Maps an. Die 2. Maps sind jedoch so klein, dass
ich das als suboptimal ansehe. Eine Treemap? müsste sortiert sein, da die
Reihenfolge der Kanten relevant ist. Wenn man davon ausgeht, dass 3 ein
guter Durchnitt der Länge ist, sollte es besser sein, ein Array zu allokieren.

Die Eltern sind nicht geordnet. Spart das etwas ein? Die From-Einträge
könnten gehasht sein. Gibt's ein multi-hash (mit Verkettung der Überläufer)
in Go? https://godoc.org/github.com/hit9/htree  (sieht gut aus, kann aber
wohl keine multi-insert ohne den Key int32 zu manipulieren. Wie? 32 bits sind
wohl zu wenig, um Anzahl und ID eines Knotens unterzubringen).

Wie implementiert man Subgraph-Sharing und ambiguous Parent Merging?
Höhe des Knotens muss gehalten werden und Span der Blätter. Höhe ist
relevant für Subgraphs: (State,h) wird identifiziert. Span ist relevant
für Merging: (Symbol,span) wird merged. => ParseTree dient auch als Rope:
Halten des Text-Anteils. https://en.wikipedia.org/wiki/Rope_(data_structure)
(symbol,span) kann man log suchen, also erst mal kein zusätzl Hash.
(state,h) kann irgendwo im Forest sein, also hashing.
*/

const (
	StateUpated int8 = iota
	Changed
	Deleted
)

type Span interface {
	Start() int64
	End() int64
}

type Delta interface {
	Version() int64
	Inserted() Span
	Removed() Span
}

// ---------------------------------------------------------------------------

// Every node in the DAG has a unique ID
type NodeID int64

// A node in a parse tree (which is in fact a forest/DAG)
type ParseNode struct {
	ids               []NodeID // unique ID or list of IDs (for merged nodes)
	leftPos, rightPos int64    // span of lexeme
	//Payload           interface{}
	//PreSpanLen        int64 // for rope-like bookkeeping
	//Rank              int    // height of subtree up to node
}

func (pnode *ParseNode) String() string {
	if pnode.ids[0] == 0 {
		return "root"
	}
	return fmt.Sprintf("N%v", pnode.ids)
	//return strconv.FormatInt(int64(pnode.ids[0]), 10)
}

type daglinks []*ParseNode

// This is in fact a hierarchical multi-DAG with ordered children
type ParseTree struct {
	root      *ParseNode
	edgesfrom map[NodeID]daglinks
	edgesto   map[NodeID]daglinks
	seqcnt    NodeID
}

// Create an empty parse tree
func NewParseTree() *ParseTree {
	ptree := &ParseTree{}
	r := ptree.NewNode()
	ptree.root = r
	ptree.edgesfrom = make(map[NodeID]daglinks)
	ptree.edgesto = make(map[NodeID]daglinks)
	return ptree
}

// Get the root node of a parse tree/forest, a sentinel node.
func (ptree *ParseTree) Root() *ParseNode {
	return ptree.root
}

// Create a node for a tree
func (ptree *ParseTree) NewNode() *ParseNode {
	pnode := &ParseNode{}
	pnode.ids = make([]NodeID, 1)
	pnode.ids[0] = ptree.seqcnt
	ptree.seqcnt++
	return pnode
}

// Number of parents for a node.
func (ptree *ParseTree) ParentCount(nid NodeID) int {
	return len(ptree.edgesto[nid])
}

// Get all parents of a node.
func (ptree *ParseTree) Parents(nid NodeID) []*ParseNode {
	plinks := ptree.edgesto[nid]
	return plinks.asNodeList()
}

func (ptree *ParseTree) ChildCount(nid NodeID) int {
	return len(ptree.edgesfrom[nid])
}

// Get all children of a node.
func (ptree *ParseTree) Children(nid NodeID) []*ParseNode {
	plinks := ptree.edgesfrom[nid]
	return plinks.asNodeList()
}

// Add a child to a node. The node will be added as a parent to the child.
func (ptree *ParseTree) AddChild(pnid NodeID, chnid NodeID) {
	// addLink(ptree.edgesfrom, pnode.ids[0], pnid, chnid)
	// addLink(ptree.edgesto, ch.ids[0], chnid, pnid) // TODO
}

func addLink(m map[NodeID]daglinks, outID NodeID, outnode *ParseNode, target *ParseNode) {
	plinks := m[outID]
	if plinks == nil {
		plinks = make(daglinks, 0, 3)
		plinks = append(plinks, target)
		m[outID] = plinks
	} else {
		plinks = append(plinks, target)
		delete(m, outID)
		m[outID] = plinks
	}
}

// Pack nodes 1 and 2 into 1.
func (ptree *ParseTree) packNodes(pn1 *ParseNode, pn2 *ParseNode) {
	pn1.ids = append(pn1.ids, pn2.ids...)
	for _, id := range pn2.ids {
		T().Debugf("node %v has id = %v", pn2, id)
		for _, pnode := range ptree.edgesto[id] { // all nodes pointing to an ID of pn2
			T().Debugf("node %v points to %v", pnode, pn2)
			for _, otherid := range pnode.ids {
				for i := 0; i < len(ptree.edgesfrom[otherid]); i++ {
					target := ptree.edgesfrom[otherid][i]
					if target == pn2 { // that's a back link to pn2
						ptree.edgesfrom[otherid][i] = pn1 // now link to the packed node
					}
				}
			}
		}
	}
}

// make a clone of a slice of links
func (dlinks daglinks) asNodeList() []*ParseNode {
	if dlinks == nil {
		return make(daglinks, 0)
	}
	nodes := make([]*ParseNode, len(dlinks))
	for i, l := range dlinks {
		nodes[i] = l
	}
	return nodes
}

// === A Versioned Variant of a DAG ==========================================

/*
type VersionedNode struct {
	BaseNode
	State int8
}

type BaseLeaf struct {
	graph *BaseDAG
	id    simple.Node
	State int8
	log   []Delta
}

func (vl *BaseLeaf) Children() []DAGNode {
	return nil
}

func (vl *BaseLeaf) addChild(ch DAGNode) {
	panic("internal error: cannot add child to versioned leaf")
}
*/
