/*
Package dss implements variants of a DAG-structured stack (DSS).
It is used for GLR-parsing and for substring parsing.
A DSS is suitable for parsing with ambiguous grammars, so the parser can
execute a breadth-first (quasi-)parallel shift and/or reduce operation
in inadequate (ambiguous) parse states.

Each parser (logical or real) sees it's own linear stack, but all stacks
together form a DAG. Stacks share common fragments, making a DSS more
space-efficient than a separeted forest of stacks.
All stacks are anchored at the root of the DSS.

	root := NewRoot("G", -1)       // represents the DSS
	stack1 := NewStack(root)       // a linear stack within the DSS
	stack2 := NewStack(root)

Please note that a DSS (this one, at least) is not well suited for general
purpose stack operations. The non-deterministic concept of GLR-parsing
will always lurk in the detail of this implementation. There are other
stack implementations around which are much better suited, especially if
performance matters.

API

The main API of this DSS consists of

	root          = NewRoot("...", impossible)
	stack         = NewStack(root)
	state, symbol = stack.Peek()                 // peek at top state and symbol of stack
	newstack      = stack.Fork()                 // duplicate stack
	stacks        = stack.Reduce(handle)         // reduce with RHS of a production
	stack         = stack.Push(state, symbol)    // transition to new parse state, i.e. a shift
	                stack.Die()                  // end of life for this stack

Additionally, there are some low-level methods, which may help debugging or
implementing your own add-on functionality.

	nodes = stack.FindHandlePath(handle, 0)    // find a string of symbols down the stack
	DSS2Dot(root, path, writer)                  // output to Graphviz DOT format
	WalkDAG(root, worker, arg)                   // execute a function on each DSS node

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

// We often will have to deal with fragments of a stack
type NodePath []*DSSNode

// === DSS Data Structure ====================================================

/*
Root for a DSS-stack. All stacks of a DSS stack structure share a common
root. Clients have to create a root before stacks can be created.
*/
type DSSRoot struct {
	Name      string    // identifier for the DSS
	bottom    *DSSNode  // stopper
	stacks    []*Stack  // housekeeping
	reservoir *ssl.List // list of nodes to be re-used
}

/*
Create a named root for a DSS-stack, given a name.
The second parameter is an implementation detail: clients have to supply
a state the root will use as a stopper. It should be an "impossible" state
for normal execution, to avoid confusion.
*/
func NewRoot(name string, invalidState int) *DSSRoot {
	root := &DSSRoot{Name: name}
	root.bottom = newDSSNode(invalidState, &pseudo{"bottom"})
	root.bottom.pathcnt = 1
	root.stacks = make([]*Stack, 0, 10)
	root.reservoir = ssl.New()
	return root
}

/*
Get the stacks heads currently active in the DSS.
No order is guaranteed.
*/
func (root *DSSRoot) ActiveStacks() []*Stack {
	dup := make([]*Stack, len(root.stacks))
	copy(dup, root.stacks)
	return dup
}

// Remove a stack from the list of stacks
func (root *DSSRoot) removeStack(stack *Stack) {
	for i, s := range root.stacks {
		if s == stack {
			root.stacks[i] = root.stacks[len(root.stacks)-1]
			root.stacks[len(root.stacks)-1] = nil
			root.stacks = root.stacks[:len(root.stacks)-1]
		}
	}
}

// As a slight optimization, we do not throw away popped nodes, but rather
// append them to a free-list for re-use.
// TODO: create an initial pool of nodes.
// TODO: migrate this to stdlib's sync/pool.
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
Check if a stack may safely delete nodes for a reduce operation. There are
exactly two cases when deletion is not permitted: (1) One or more other stacks
are sitting on the same node as this one (2) There are nodes present further
up the stack (presumably operated on by other stacks).
*/
func (root *DSSRoot) IsAlone(stack *Stack) bool {
	if stack.tos.successorCount() > 0 {
		return false
	}
	for _, other := range root.stacks {
		if other.tos == stack.tos {
			return false
		}
	}
	return true
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
}

/*
Create a new linear stack within a DSS.
*/
func NewStack(root *DSSRoot) *Stack {
	s := &Stack{root: root, tos: root.bottom}
	s.root.stacks = append(s.root.stacks, s)
	return s
}

// Calculate the height of a stack
/*
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
Check if a stack may safely delete nodes for a reduce operation. There are
exactly two cases when deletion is not permitted: (1) One or more other stacks
are sitting on the same node as this one (2) There are nodes present further
up the stack (presumably operated on by other stacks).
*/
func (stack *Stack) IsAlone() bool {
	return stack.root.IsAlone(stack)
}

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
leading to state state.

The method looks for a TOS in the DSS representing the combination of
(state, sym) and -- if present -- uses (reverse joins) the TOSs. Otherwise it
creates a new DAG node.
*/
func (stack *Stack) Push(state int, sym lr.Symbol) *Stack {
	if sym == nil {
		return stack
	}
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

/*
For reduce operations on a stack it is necessary to identify the
nodes which correspond to the symbols of the RHS of the production
to reduce. The nodes (some or all) may overlap with other stacks of
the DSS.

For a RHS there may be more than one path downwards the stack marked
with the symbols of RHS. It is not deterministic which one this method
will find and return.
Clients may call this method multiple times. If a clients wants to find
all simultaneously existing paths, it should increment the parameter skip
every time. This determines the number of branches have been found
previously and should be skipped now.

Will return nil if no (more) path is found.
*/
func (stack *Stack) FindHandlePath(handle []lr.Symbol, skip int) NodePath {
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
func collectHandleBranch(n *DSSNode, handleRest []lr.Symbol, handleLen int, skip *int) (NodePath, bool) {
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
func (stack *Stack) splitOff(path NodePath) *Stack {
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
Pop nodes corresponding to a handle from a stack.

The handle path nodes must exist in the DSS (not checked again).

The decision wether to delete the nodes on the way down is not trivial.
If the destructive-flag is unset, we do not delete anything. If it is set
to true, we take this as a general permission to not have to regard other
stacks. But we must be careful not to burn our own bridges. We may need a
node for other (amybiguity-)paths we'll reduce within the same operation.
The logic is as follows:

(1) If we have already deleted the single predecessor of this node and
    this is a regular node (no join, no fork), then we may delete this one
    as well.

(2) If this is a reverse join, we are decrementing the linkcnt and have
    permission to delete it, if it is the last run through this node.

(3) If this is still a reverse fork (although we possbily have deleted one
    successor), we do not delete.
*/
func (stack *Stack) reduce(path NodePath, destructive bool) (ret []*Stack) {
	maydelete := true
	haveDeleted := false
	for i := len(path) - 1; i >= 0; i-- { // iterate over handle symbols back to front
		node := path[i]
		T.Debugf("reducing node %v (now cnt=%d)", node, node.pathcnt)
		T.Debugf("         node %v has %d succs", node, len(node.succs))
		if node.isInverseJoin() {
			T.Debugf("is join: %v", node)
			maydelete = true
			node.pathcnt--
		} else if haveDeleted && node.successorCount() > 0 {
			T.Debugf("is or has been fork: %v", node)
			maydelete = false
		} else {
			maydelete = true
			node.pathcnt--
		}
		if i == 0 { // when popped every node: every predecessor is a stack head now
			ret = stack.root.makeStackHeadsFrom(node.preds)
		}
		if destructive && maydelete && node.pathcnt == 0 {
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

	as a DSS: -[a]------
	                    [A] [+]-[A]     // now reduce this: A + A  to  X
	          -[A]-[+]--

This will result in 2 new stack heads:

	as a DSS: -[a]                // return stack #1 of Reduce()
	          -[A]-[+]            // return stack #2 of Reduce()

After pushing X onto the stack, the 2 stacks may be merged (on 'X'), thus
resulting in a single stack head again.

	as a DSS: -[a]-----
	                   [X]
	          -[A]-[+]-

The stack triggering the reduce will be the first stack within the returning
array of stacks, i.e. the Reduce() return the stack itself plus possibly
other stacks created during the reduce operation.
*/
func (stack *Stack) Reduce(handle []lr.Symbol) (ret []*Stack) {
	if len(handle) == 0 {
		return
	}
	var paths []NodePath // first collect all possible handle paths
	skip := 0            // how many paths already found?
	path := stack.FindHandlePath(handle, skip)
	for path != nil {
		paths = append(paths, path)
		skip++
		path = stack.FindHandlePath(handle, skip)
	}
	destructive := stack.IsAlone() // give general permission to delete reduced nodes?
	for _, path = range paths {    // now reduce along every path
		stacks := stack.reduce(path, destructive)
		ret = append(ret, stacks...) // collect returned stack heads
	}
	if len(ret) > 0 { // if avail, replace 1st stack with this
		stack.tos = ret[0].tos // make us a lookalike of the 1st returned one
		ret[0].Die()           // the 1st returned one has to die
		ret[0] = stack         // and we replace it
	}
	return
}

func (stack *Stack) reduceHandle(handle []lr.Symbol, destructive bool) (ret []*Stack) {
	haveReduced := true
	skipCnt := 0
	for haveReduced { // as long as a reduction has been done
		haveReduced = false
		handleNodes := stack.FindHandlePath(handle, skipCnt)
		if handleNodes != nil {
			haveReduced = true
			if !destructive {
				skipCnt++
			}
			s := stack.reduce(handleNodes, destructive)
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

For parsers Pop() may not be very useful. It is included here for
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

If parameter path is given, it will be highlighted in the output.
*/
func DSS2Dot(root *DSSRoot, path []*DSSNode, w io.Writer) {
	istos := map[*DSSNode]bool{}
	for _, stack := range root.stacks {
		istos[stack.tos] = true
	}
	ids := map[*DSSNode]int{}
	idcounter := 1
	io.WriteString(w, "digraph {\n")
	WalkDAG(root, func(node *DSSNode, arg interface{}) {
		ids[node] = idcounter
		idcounter++
		styles := nodeDotStyles(node, pathContains(path, node))
		if istos[node] {
			styles += ",shape=box"
		}
		io.WriteString(w, fmt.Sprintf("\"%d[%d]\" [label=%d %s];\n",
			node.State, ids[node], node.State, styles))
	}, nil)
	WalkDAG(root, func(node *DSSNode, arg interface{}) {
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
	WalkDAG(root, func(node *DSSNode, arg interface{}) {
		predList := " "
		for _, p := range node.preds {
			predList = predList + strconv.Itoa(p.State) + " "
		}
		fmt.Printf("edge  [%s]%v\n", predList, node)
	}, nil)
}

/*
Walk all nodes of a DSS and execute a worker function on each.
parameter arg is presented as a second argument to each worker execution.
*/
func WalkDAG(root *DSSRoot, worker func(*DSSNode, interface{}), arg interface{}) {
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

/*
End of lfe for this stack.
It will be detached from the DSS root and will let go of its TOS.
*/
func (stack *Stack) Die() {
	stack.root.removeStack(stack)
	stack.tos = nil
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
