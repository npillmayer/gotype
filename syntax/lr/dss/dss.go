/*
Package dss implements variants of a DAG-structured stack (DSS).
It is used for GLR-parsing and for substring parsing.
A DSS is suitable for parsing with ambiguous grammars, so the parser can
execute a breadth-first (quasi-)parallel shift and/or reduce operation
in ambiguous parse states.

Each parser (logical or real) sees it's own linear stack, but all stacks
together form a DAG. Stacks share common fragments, making a DSS more
space-efficient than a separeted forest of stacks.
All stacks are anchored at the root of the DSS.

	root := NewRoot("G")           // represents the DSS
	stack1 := NewStack(root)       // a linear stack within the DSS
	stack2 := NewStack(root)

Please note that a DSS (this one, at least) is not well suited for general
purpose stack operations. The non-deterministic concept of GLR-parsing
will always lurk in the detail of this implementation. There are other
stack implementations around which are much better suited for general
purpose stack operations.

API

The main API of this DSS consists of

	root          = NewRoot("...")
	stack         = NewStack(root)
	state, symbol = stack.Peek()          // peek at top state and symbol of stack
	newstack      = stack.Fork()          // duplicate stack
	stacks        = stack.Reduce(handle)  // reduce with RHS of a production
	stack.Push(state, symbol)             // transition to new parse state, i.e. a shift

Other methods of the API are rarely used in parsing and exist more or less
to complete a conventional stack API. Note that a method for determining the
size of a stack is missing.

	stacks = stack.Pop()

GLR Parsing

A GLR parser forks a stack whenever an ambiguous state on top of the
parse stack is processed.
In a GLR-parser shift/reduce- and reduce/reduce-conflicts are
not uncommon, as potentially ambiguous grammars are used.

Assume that a shift/reduce-conflict is signalled by the TOS.
Then the parser will duplicate the stack (stack.Fork()) and carry out both
operations: for one stack a symbol will be shifted, the other will be
used for the reduce-operations.

Further Information

For further information see for example

	https://people.eecs.berkeley.edu/~necula/Papers/elkhound_cc04.pdf

Status

This is experimental software, currently not intended for production use.

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
*/
package dss

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	ssl "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
)

// Configurable trace
var T tracing.Trace = tracing.SyntaxTracer

// --- Stack Nodes -----------------------------------------------------------

// Type for a node in a DSS stack
type DSSNode struct {
	State   int        // identifier of a parse state
	Sym     lr.Symbol  // symbol on the stack (between states)
	preds   []*DSSNode // predecessors of this node (inverse join)
	succs   []*DSSNode // successors of this node (inverse fork)
	pathcnt int        // count of paths through this node
}

// create an empty stack node
func newDSSNode(state int, sym lr.Symbol) *DSSNode {
	n := &DSSNode{}
	n.State = state
	n.Sym = sym
	return n
}

// Simple Stringer for a stack node.
func (n *DSSNode) String() string {
	return fmt.Sprintf("<-(%v)-[%d]", n.Sym, n.State)
}

// Make an identical copy of a node
func (n *DSSNode) clone() *DSSNode {
	nn := newDSSNode(n.State, n.Sym)
	for _, p := range n.preds {
		nn.prepend(p)
	}
	for _, s := range n.succs {
		nn.append(s)
	}
	return nn
}

// Append a node, i.e., insert it into the list of successors
func (n *DSSNode) append(succ *DSSNode) {
	if n.succs == nil {
		n.succs = make([]*DSSNode, 0, 3)
	}
	n.succs = append(n.succs, succ)
}

// Prepend a node, i.e., insert it into the list of predecessors
func (n *DSSNode) prepend(pred *DSSNode) {
	if n.preds == nil {
		n.preds = make([]*DSSNode, 0, 3)
	}
	n.preds = append(n.preds, pred)
}

// Unlink a node
func (n *DSSNode) isolate() {
	T.Debugf("isolating %v", n)
	for _, p := range n.preds {
		p.succs, _ = remove(p.succs, n)
	}
	for _, s := range n.succs {
		s.preds, _ = remove(s.preds, n)
	}
}

func (n *DSSNode) is(state int, sym lr.Symbol) bool {
	return n.State == state && n.Sym == sym
}

func (n *DSSNode) isInverseJoin() bool {
	return n.preds != nil && len(n.preds) > 1
}

func (n *DSSNode) isInverseFork() bool {
	return n.succs != nil && len(n.succs) > 1
}

func (n *DSSNode) successorCount() int {
	if n.succs != nil {
		return len(n.succs)
	}
	return 0
}

func (n *DSSNode) predecessorCount() int {
	if n.preds != nil {
		return len(n.preds)
	}
	return 0
}

func (n *DSSNode) findUplink(sym lr.Symbol) *DSSNode {
	s := contains(n.succs, sym)
	return s
}

func (n *DSSNode) findDownlink(state int) *DSSNode {
	s := hasState(n.preds, state)
	return s
}

// === DSS Data Structure ====================================================

/*
Root for a DSS-stack. All stacks of a DSS stack structure share a common
root. Clients have to create a root before stacks can be created.
*/
type DSSRoot struct {
	Name      string     // identifier for the DSS
	roots     *DSSNode   // points to a chain of sibling root nodes
	toss      []*DSSNode // list of all TOSs
	bottom    *DSSNode   // stopper
	stacks    []*Stack   // housekeeping
	reservoir *ssl.List  // list of nodes to be re-used
}

// Create a named root for a DSS-stack.
func NewRoot(name string) *DSSRoot {
	root := &DSSRoot{Name: name}
	root.bottom = newDSSNode(-999, &pseudo{"bottom"})
	root.bottom.pathcnt = 1
	root.stacks = make([]*Stack, 0, 10)
	root.reservoir = ssl.New()
	return root
}

// As a slight optimization, we do not throw away popped nodes, but rather
// append them to a free-list for re-use.
// TODO: create an initial pool of nodes.
func (root *DSSRoot) recycleNode(node *DSSNode) {
	T.Debugf("recycling node %v", node)
	node.State = 0
	node.Sym = nil
	node.preds = node.preds[0:0]
	node.succs = node.succs[0:0]
	root.reservoir.Append(node) // list of nodes free for grab
}

// Create a new stack node or fetch a recycled one.
func (root *DSSRoot) newNode(state int, sym lr.Symbol) *DSSNode {
	var node *DSSNode
	n, ok := root.reservoir.Get(0)
	if ok {
		root.reservoir.Remove(0)
		node = n.(*DSSNode)
		node.State = state
		node.Sym = sym
	} else {
		node = newDSSNode(state, sym)
	}
	return node
}

// Wrap nodes of a stack as heads of new stacks.
func (root *DSSRoot) makeStackHeadsFrom(nodes []*DSSNode) (ret []*Stack) {
	for _, n := range nodes {
		s := NewStack(root)
		s.tos = n
		ret = append(ret, s)
	}
	return
}

// Find a TOS of any stack in the DSS which carries state and sym.
// This is needed for suffix merging. Returns nil if no TOS meets the criteria.
func (root *DSSRoot) findTOSofAnyStack(state int, sym lr.Symbol) *DSSNode {
	for _, stack := range root.stacks {
		if stack.tos.is(state, sym) {
			return stack.tos
		}
	}
	return nil
}

/*
A data structure for a linear stack within a DSS (DAG structured stack = DSS).
All stacks together form a DSS, i.e. they
may share portions of stacks. Each client carries its own stack, without
noticing the other stacks. The data structure hides the fact that stack-
fragments may be shared.
*/
type Stack struct {
	root *DSSRoot // common root for a set of sub-stacks
	tos  *DSSNode // top of stack
	//height int      // height of the stack
}

/*
Create a new linear stack within a DSS.
*/
func NewStack(root *DSSRoot) *Stack {
	s := &Stack{root: root, tos: root.bottom}
	s.root.stacks = append(s.root.stacks, s)
	return s
}

/* Height of a stack (number of nodes).
func (stack *Stack) Size() int {
	return stack.height
}

// re-calculate the height of a stack
func (stack *Stack) calculateHeight() {
	stack.height = stack.tos.height(stack.root.bottom) // recurse until hits bottom
}

func (node *DSSNode) height(stopper *DSSNode) int {
	if node == stopper {
		return 0
	}
	max := 0
	for _, p := range node.preds {
		h := p.height(stopper)
		if h > max {
			max = h
		}
	}
	return max + 1
}
*/

/*
Duplicate the head of a stack, resulting in a new stack.
The new stack is implicitely registered with the (common) root.
*/
func (stack *Stack) Fork() *Stack {
	s := NewStack(stack.root)
	//s.height = stack.height
	s.tos = stack.tos
	s.tos.pathcnt++
	return s
}

/*
Return TOS.Symbol and TOS.State without popping it.
*/
func (stack *Stack) Peek() (int, lr.Symbol) {
	return stack.tos.State, stack.tos.Sym
}

/*
Push a state and a symbol on the stack.
Interpretation is as follows: a transition has been found consuming symbol sym,
leading to state.

The method looks for a TOS in the DSS representing the combination of
(state, sym) and -- if present -- uses (reverse joins) the TOSs. Otherwise it
creates a new DAG node.
*/
func (stack *Stack) Push(state int, sym lr.Symbol) *Stack {
	// create or find a node
	// - create: new and let node.prev be tos
	// - find: return existing one
	// update references / pathcnt
	if succ := stack.tos.findUplink(sym); succ != nil {
		T.Debugf("state already present: %v", succ)
		stack.tos = succ // pushed state is already there, upchain
	} else {
		succ := stack.root.findTOSofAnyStack(state, sym)
		if succ == nil { // not found in DSS
			succ = stack.root.newNode(state, sym)
			T.Debugf("creating state: %v", succ)
			succ.pathcnt = stack.tos.pathcnt
		} else {
			T.Debugf("found state on other stack: %v", succ)
			succ.pathcnt++
		}
		succ.prepend(stack.tos)
		stack.tos.append(succ)
		stack.tos = succ
	}
	return stack
}

// For reduce operations on a stack it is necessary to identify the
// nodes which correspond to the symbols of the RHS of the production
// to reduce. The nodes (some or all) may overlap with other stacks of
// the DSS.
// For a RHS there may be more than one path downwards the stack marked
// with the symbols of RHS. It is not deterministic which one this method
// will find and return.
func (stack *Stack) findHandleBranch(handle []lr.Symbol, skip int) []*DSSNode {
	path, ok := collectHandleBranch(stack.tos, handle, len(handle), &skip)
	if ok {
		T.Debugf("found a handle %v", path)
	}
	return path
}

// Recursive function to collect nodes corresponding to a list of symbols.
// The bottom-most node, i.e. the one terminating the recursion, will allocate
// the result array. This array will then be handed upwards the call stack,
// being filled with node links on the way.
func collectHandleBranch(n *DSSNode, handleRest []lr.Symbol, handleLen int, skip *int) ([]*DSSNode, bool) {
	l := len(handleRest)
	if l > 0 {
		if n.Sym == handleRest[l-1] {
			T.Debugf("handle symbol match at %d = %v", l-1, n.Sym)
			for _, pred := range n.preds {
				branch, found := collectHandleBranch(pred, handleRest[:l-1], handleLen, skip)
				if found {
					if *skip == 0 {
						T.Debugf("partial branch: %v", branch)
						if branch != nil { // a-ha, deepest node has created collector
							branch[l-1] = n // collect n in branch
						}
						return branch, true // return with partially filled branch
					} else {
						*skip = max(0, *skip-1)
					}
				}
			}
		}
	} else {
		branchCreated := make([]*DSSNode, handleLen) // create the handle-path collector
		return branchCreated, true                   // search/recursion terminates
	}
	return nil, false // abort search, handle-path not found
}

// This is probably never used for a real parser
func (stack *Stack) splitOff(path []*DSSNode) *Stack {
	l := len(path) // pre-condition: there are at least l nodes on the stack-path
	var node, mynode, upperNode *DSSNode = nil, nil, nil
	mystack := NewStack(stack.root)
	//mystack.height = stack.height
	for i := l - 1; i >= 0; i-- { // walk path from back to font, i.e. top-down
		node = path[i]
		mynode = stack.root.newNode(node.State, node.Sym)
		//mynode = stack.root.newNode(node.State+100, node.Sym)
		T.Debugf("split off node %v", mynode)
		if upperNode != nil { // we are not the top-node of the stack
			mynode.append(upperNode)
			upperNode.prepend(mynode)
			upperNode.pathcnt++
		} else { // is the newly created top node
			mystack.tos = mynode // make it TOS of new stack
		}
		upperNode = mynode
	}
	if mynode != nil && node.preds != nil {
		for _, p := range node.preds {
			mynode.prepend(p)
			mynode.pathcnt++
		}
	}
	return mystack
}

/*
// find join(s) and split(s) within handleLen
// if merged: de-merge handleLen symbols
// remove handle symbols
// push lhs symbol
// return new top symbol (should be lhs)

The handle nodes must exist in the DSS (not checked again).
*/
func (stack *Stack) reduce(handleNodes []*DSSNode) (ret []*Stack) {
	dodelete := true
	if stack.tos.successorCount() > 0 { // there are nodes above
		dodelete = false
	}
	haveDeleted := false
	l := len(handleNodes)
	var node *DSSNode
	for i := l - 1; i >= 0; i-- { // iterate handle symbols back to front
		node = handleNodes[i]
		T.Debugf("reducing node %v (now cnt=%d)", node, node.pathcnt)
		T.Debugf("         node %v has %d succs", node, len(node.succs))
		if node.isInverseJoin() {
			T.Debugf("is join: %v", node)
			dodelete = true
			node.pathcnt--
		} else if haveDeleted && node.successorCount() > 0 {
			T.Debugf("is in line: %v", node)
			dodelete = false
		} else {
			dodelete = true
			node.pathcnt--
		}
		if i == 0 {
			ret = stack.root.makeStackHeadsFrom(node.preds)
		}
		if dodelete && node.pathcnt == 0 {
			node.isolate()
			stack.root.recycleNode(node)
			haveDeleted = true
		} else {
			haveDeleted = false
		}
	}
	return
}

/*
Perform a reduce operation, given a handle, i.e. a right hand side of a
grammar rule. Strictly speaking, it performs not a complete reduce operation,
but just the first part: popping the RHS symbols off the stack.
Clients will have to push the LHS symbol separately.

With a DSS, reduce may result in a multitude of new stack configurations.
Whenever there is a reduce/reduce conflict or a shift/reduce-conflict, a GLR
parser will perform both reduce-operations. To this end each possible operation
(i.e, parser) will (conceptually) use its own stack.
Thus multiple return values of a reduce operation correspond to (temporary
or real) ambiguity of a grammar.

Example:  X ::= A + A  (grammar rule)

	Stack 1:  ... a A + A
	Stack 2:  ... A + A + A

	as a DSS: -[a]-------
	                    [A] [+]-[A]        now reduce this: A + A  to  X
	          -[A]-[+]---

This will result in 2 new stack heads:

	as a DSS: -[a]-
	          -[A]-[+]-

After pushing X onto the stack, the 2 stacks may be merged (on 'X'), thus
resulting in a single stack head again.

	as a DSS: -[a]-----
	                   [X]
	          -[A]-[+]-

*/
func (stack *Stack) Reduce(handle []lr.Symbol) (ret []*Stack) {
	haveReduced := true
	foundCnt := 0
	for haveReduced {
		haveReduced = false
		handleNodes := stack.findHandleBranch(handle, foundCnt)
		if handleNodes != nil {
			haveReduced = true
			foundCnt++
			s := stack.reduce(handleNodes)
			ret = append(ret, s...)
		}
	}
	return
}

/*
Pop the TOS of a stack. This is straigtforward for (non-empty) linear stacks
without shared nodes. For stacks with common suffixes, i.e. with inverse joins,
it is more tricky. The popping may result in multiple new stacks, one for
each predecessor.

During parsing Pop() may not be very useful. It is included here for
conformance with the general contract of a stack. With parsing, popping of
states happens during reductions, and the API offers more convenient
functions for this.
*/
func (stack *Stack) Pop() (ret []*Stack) {
	if stack.tos != stack.root.bottom { // ensure stack not empty
		// If tos is part of another chain: return node and go down one node
		// If shared by another stack: return node and go down one node
		// If not shared: remove node
		// Increase pathcnt at new TOS
		var oldTOS *DSSNode = stack.tos
		//var r []*Stack
		for i, n := range stack.tos.preds { // at least 1 predecessor
			if i == 0 { // 1st one: keep this stack
				stack.tos = n
				//stack.calculateHeight()
			} else { // further predecessors: create stack for each one
				s := NewStack(stack.root)
				s.tos = n
				//s.calculateHeight()
				//T.Debugf("creating new stack for %v (of height=%d)", n, s.height)
				ret = append(ret, s)
			}
		}
		if oldTOS.succs == nil || len(oldTOS.succs) == 0 {
			oldTOS.isolate()
			stack.root.recycleNode(oldTOS)
		}
		return ret
	}
	return nil
}

func (stack *Stack) pop(toNode *DSSNode, deleteNode bool, collectStacks bool) ([]*Stack, error) {
	var r []*Stack                      // return value
	var err error                       // return value
	if stack.tos != stack.root.bottom { // ensure stack not empty
		var oldTOS *DSSNode = stack.tos
		var found bool
		stack.tos.preds, found = remove(stack.tos.preds, toNode)
		if !found {
			err = errors.New("unable to pop TOS: 2OS not appropriate")
		} else {
			stack.tos.pathcnt--
			toNode.pathcnt++
			stack.tos = toNode
			//stack.calculateHeight()
			if deleteNode && oldTOS.pathcnt == 0 &&
				(oldTOS.succs == nil || len(oldTOS.succs) == 0) {
				oldTOS.isolate()
				stack.root.recycleNode(oldTOS)
			}
		}
	} else {
		T.Error("unable to pop TOS: stack empty")
		err = errors.New("unable to pop TOS: stack empty")
	}
	return r, err
}

// --- Debugging -------------------------------------------------------------

/*
Output a DSS in Graphviz DOT format (for debugging purposes).
*/
func DSS2Dot(root *DSSRoot, path []*DSSNode, w io.Writer) {
	istos := map[*DSSNode]bool{}
	for _, stack := range root.stacks {
		istos[stack.tos] = true
	}
	ids := map[*DSSNode]int{}
	idcounter := 1
	io.WriteString(w, "digraph {\n")
	walkDAG(root, func(node *DSSNode, arg interface{}) {
		ids[node] = idcounter
		idcounter++
		styles := nodeDotStyles(node, pathContains(path, node))
		if istos[node] {
			styles += ",shape=box"
		}
		io.WriteString(w, fmt.Sprintf("\"%d[%d]\" [label=%d %s];\n",
			node.State, ids[node], node.State, styles))
	}, nil)
	walkDAG(root, func(node *DSSNode, arg interface{}) {
		ww := w.(io.Writer)
		for _, p := range node.preds {
			io.WriteString(ww, fmt.Sprintf("\"%d[%d]\" -> \"%d[%d]\" [label=\"%v\"];\n",
				node.State, ids[node], p.State, ids[p], node.Sym))
		}
	}, w)
	io.WriteString(w, "}\n")
}

func nodeDotStyles(node *DSSNode, highlight bool) string {
	s := ",style=filled"
	if highlight {
		s = s + fmt.Sprintf(",fillcolor=\"%s\"", hexhlcolors[node.pathcnt])
	} else {
		s = s + fmt.Sprintf(",fillcolor=\"%s\"", hexcolors[node.pathcnt])
	}
	return s
}

var hexhlcolors = [...]string{"#FFEEDD", "#FFDDCC", "#FFCCAA", "#FFBB88", "#FFAA66",
	"#FF9944", "#FF8822", "#FF7700", "#ff6600"}

var hexcolors = [...]string{"white", "#CCDDFF", "#AACCFF", "#88BBFF", "#66AAFF",
	"#4499FF", "#2288FF", "#0077FF", "#0066FF"}

// Debugging
func PrintDSS(root *DSSRoot) {
	walkDAG(root, func(node *DSSNode, arg interface{}) {
		predList := " "
		for _, p := range node.preds {
			predList = predList + strconv.Itoa(p.State) + " "
		}
		fmt.Printf("edge  [%s]%v\n", predList, node)
	}, nil)
}

func walkDAG(root *DSSRoot, worker func(*DSSNode, interface{}), arg interface{}) {
	visited := map[*DSSNode]bool{}
	var walk func(*DSSNode, func(*DSSNode, interface{}), interface{})
	walk = func(node *DSSNode, worker func(*DSSNode, interface{}), arg interface{}) {
		if visited[node] {
			return
		}
		visited[node] = true
		for _, p := range node.preds {
			walk(p, worker, arg)
		}
		worker(node, arg)
	}
	for _, stack := range root.stacks {
		walk(stack.tos, worker, arg)
	}
}

// --- Helpers ---------------------------------------------------------------

func pathContains(s []*DSSNode, node *DSSNode) bool {
	if s != nil {
		for _, n := range s {
			if n == node {
				return true
			}
		}
	}
	return false
}

func contains(s []*DSSNode, sym lr.Symbol) *DSSNode {
	if s != nil {
		for _, n := range s {
			if n.Sym == sym {
				return n
			}
		}
	}
	return nil
}

// Helper : remove an item from a node slice
// without preserving order (i.e, replace by last in slice)
func remove(nodes []*DSSNode, node *DSSNode) ([]*DSSNode, bool) {
	for i, n := range nodes {
		if n == node {
			nodes[i] = nodes[len(nodes)-1]
			nodes[len(nodes)-1] = nil
			nodes = nodes[:len(nodes)-1]
			return nodes, true
		}
	}
	return nodes, false
}

func hasState(s []*DSSNode, state int) *DSSNode {
	for _, n := range s {
		if n.State == state {
			return n
		}
	}
	return nil
}

type pseudo struct {
	name string
}

func pseudosym(name string) lr.Symbol {
	return &pseudo{name: name}
}

func (sy *pseudo) String() string {
	return sy.name
}

func (sy *pseudo) IsTerminal() bool {
	return true
}

func (sy *pseudo) Token() int {
	return 0
}

func (sy *pseudo) GetID() int {
	return 0
}

var _ lr.Symbol = &pseudo{}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
