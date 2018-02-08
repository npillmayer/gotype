/*
Package dss implements variants of a DAG-structured state (DSS).
It is used for GLR-parsing and for substring parsing.
A DSS is suitable with non-deterministic grammars, so the parser can
execute a breadth-first (quasi-)parallel shift and/or reduce operations.

Each parser (logical or real) sees it's own linear stack, but all stacks
together form a DAG. Stacks share common fragments, making a DSS more
space-efficient than a separeted forest of stacks.
All stacks are anchored at the root of the stack.

	root := NewRoot("G")           // represents the DSS
	stack1 := NewStack(root)       // a linear stack within the DSS
	stack2 := NewStack(root)

For further information see for example

	https://people.eecs.berkeley.edu/~necula/Papers/elkhound_cc04.pdf

*/
package dss

import (
	"fmt"
	"io"
	"strconv"

	ssl "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
)

// Configurable trace
var T tracing.Trace = tracing.SyntaxTracer

// Type for a node in a DSS stack
type DSSNode struct {
	State   int        // identifier of a parse state
	Sym     lr.Symbol  // symbol on the stack (between states)
	Preds   []*DSSNode // predecessors of this node (join)
	Succs   []*DSSNode // successors of this node (fork)
	linkcnt int        // count of stacks referencing this as TOS
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
	for _, p := range n.Preds {
		nn.prepend(p)
	}
	for _, s := range n.Succs {
		nn.append(s)
	}
	return nn
}

// Append a node, i.e., insert it into the list of successors
func (n *DSSNode) append(succ *DSSNode) {
	if n.Succs == nil {
		n.Succs = make([]*DSSNode, 0, 3)
	}
	n.Succs = append(n.Succs, succ)
}

// Prepend a node, i.e., insert it into the list of predecessors
func (n *DSSNode) prepend(pred *DSSNode) {
	if n.Preds == nil {
		n.Preds = make([]*DSSNode, 0, 3)
	}
	n.Preds = append(n.Preds, pred)
}

// Unlink a node
func (n *DSSNode) isolate() {
	for _, p := range n.Preds {
		p.Succs, _ = remove(p.Succs, n)
	}
	for _, s := range n.Succs {
		s.Preds, _ = remove(s.Preds, n)
	}
}

func (n *DSSNode) is(state int, sym lr.Symbol) bool {
	return n.State == state && n.Sym == sym
}

func (n *DSSNode) findUplink(sym lr.Symbol) *DSSNode {
	s := contains(n.Succs, sym)
	return s
}

func (n *DSSNode) findDownlink(state int) *DSSNode {
	s := hasState(n.Preds, state)
	return s
}

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
	root.stacks = make([]*Stack, 0, 10)
	root.reservoir = ssl.New()
	return root
}

// As a slight optimization, we do not throw away popped nodes, but rather
// append them to a free-list for re-use.
func (root *DSSRoot) recycleNode(node *DSSNode) {
	node.State = 0
	node.Sym = nil
	node.Preds = node.Preds[0:0]
	node.Succs = node.Succs[0:0]
	root.reservoir.Append(node) // list of nodes free for grab
}

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

/*
Find a TOS of any stack in the DSS which carries state and sym.
This is needed for suffix merging. Returns nil if no TOS meets the criteria.
*/
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
	root   *DSSRoot // common root for a set of sub-stacks
	tos    *DSSNode // top of stack
	height int      // height of the stack
}

/*
Create a new linear stack within a DSS.
*/
func NewStack(root *DSSRoot) *Stack {
	s := &Stack{root: root, tos: root.bottom, height: 0}
	s.root.stacks = append(s.root.stacks, s)
	return s
}

// Height of a stack (number of nodes).
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
	for _, p := range node.Preds {
		h := p.height(stopper)
		if h > max {
			max = h
		}
	}
	return max + 1
}

/*
Duplicate the head of a stack, resulting in a new stack.
The new stack is implicitely registered with the (now common) root.
*/
func (stack *Stack) Fork() *Stack {
	s := NewStack(stack.root)
	s.height = stack.height
	s.tos = stack.tos
	s.tos.linkcnt++
	return s
}

func (stack *Stack) moveUp(succ *DSSNode) int {
	// successor of TOS must exist as succ
	stack.tos.linkcnt--
	stack.tos = succ
	stack.tos.linkcnt++
	stack.height++
	return stack.tos.State
}

// For reduce operations on a stack it is necessary to identify the
// nodes, which correspond to the symbols of the RHS of the production
// to reduce. The nodes (some or all) may overlap with other stacks of
// the DSS.
func (stack *Stack) findHandleBranch(handle []lr.Symbol) []*DSSNode {
	path, ok := stack.tos.collectHandleBranch(handle, len(handle))
	if ok {
		T.Infof("FOUND %v", path)
	}
	return path
}

// Recursive function to collect nodes corresponding to a list of symbols.
// The bottom-most node, i.e. the one terminating the recursion, will allocate
// the result array. This array will then be handed upwards the call stack,
// being filled with node links on the way.
func (n *DSSNode) collectHandleBranch(handleRest []lr.Symbol, handleLen int) ([]*DSSNode, bool) {
	l := len(handleRest)
	if l > 0 {
		if n.Sym == handleRest[l-1] {
			T.Debugf("handle symbol match at %d = %v", l-1, n.Sym)
			for _, pred := range n.Preds {
				branch, found := pred.collectHandleBranch(handleRest[:l-1], handleLen)
				if found {
					T.Debugf("partial branch: %v", branch)
					if branch != nil { // a-ha, deepest node has created collector
						branch[l-1] = n // collect n in branch
					}
					return branch, true // return with partially filled branch
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
	mystack.height = stack.height
	for i := l - 1; i >= 0; i-- { // walk path from back to font, i.e. top-down
		node = path[i]
		//mynode = stack.root.newNode(node.State, node.Sym)
		mynode = stack.root.newNode(node.State+100, node.Sym)
		T.Debugf("split off node %v", mynode)
		if upperNode != nil { // we are not the top-node of the stack
			mynode.append(upperNode)
			upperNode.prepend(mynode)
		} else {
			mynode.linkcnt++ // this is entry point for mystack
			mystack.tos = mynode
		}
		upperNode = mynode
	}
	if mynode != nil && node.Preds != nil {
		for _, p := range node.Preds {
			mynode.prepend(p)
		}
	}
	return mystack
}

/*
Push a state and a symbol on the stack.
Interpretation is as follows: a transition has been found consuming symbol sym,
leading to state.

The method looks for a TOS in the DSS representing the combination of
(state, sym) and -- if present -- uses (joins) the TOSs. Otherwise it
creates a new DAG node.
*/
func (stack *Stack) Push(state int, sym lr.Symbol) *Stack {
	// create or find a node
	// - create: new and let node.prev be tos
	// - find: return existing one
	// update references / linkcnt
	if succ := stack.tos.findUplink(sym); succ != nil {
		T.Debugf("state already present: %v", succ)
		stack.moveUp(succ) // pushed state is already there, upchain
	} else {
		succ := stack.root.findTOSofAnyStack(state, sym)
		if succ == nil { // not found in DSS
			succ = stack.root.newNode(state, sym)
			T.Debugf("creating state: %v", succ)
		} else {
			T.Debugf("found state on other stack: %v", succ)
		}
		succ.prepend(stack.tos)
		stack.tos.append(succ)
		stack.moveUp(succ)
	}
	return stack
}

/*
Return TOS.Symbol and TOS.State without popping it.
*/
func (stack *Stack) Peek() (lr.Symbol, int) {
	return stack.tos.Sym, stack.tos.State
}

// if tos is part of another chain: return node and go down one node
// if shared by another stack: return node and go down one node
// if not shared: remove node
// increase linkcnt at new TOS
func (stack *Stack) Pop() []*Stack {
	if stack.tos != stack.root.bottom { // ensure stack not empty
		var oldTOS *DSSNode = stack.tos
		var r []*Stack
		// remove TOS and go down 1 node
		stack.tos.linkcnt--
		for i, n := range stack.tos.Preds { // at least 1 predecessor
			if i == 0 { // 1st one: keep this stack
				stack.tos = n
				stack.tos.linkcnt++
				stack.height--
			} else { // further predecessors: create stack for each one
				s := NewStack(stack.root)
				s.tos = n
				s.tos.linkcnt++
				s.calculateHeight()
				T.Debugf("creating new stack for %v (of height=%d)", n, s.height)
				r = append(r, s)
			}
		}
		if oldTOS.linkcnt == 0 &&
			(oldTOS.Succs == nil || len(oldTOS.Succs) == 0) {
			oldTOS.isolate()
			stack.root.recycleNode(oldTOS)
		}
		return r
	}
	return nil
}

/*
// find join(s) and split(s) within handleLen
// if merged: de-merge handleLen symbols
// remove handle symbols
// push lhs symbol
// return new top symbol (should be lhs)

The handle nodes must exist in the DSS (not checked again).
*/
func (stack *Stack) Reduce(handleNodes []*DSSNode) []*Stack {
	return nil
}

// --- Debugging -------------------------------------------------------------

/*
Output a DSS in Graphviz DOT format (for debugging purposes).
*/
func DSS2Dot(root *DSSRoot, w io.Writer) {
	io.WriteString(w, "digraph {\n")
	walkDAG(root, func(node *DSSNode, arg interface{}) {
		ww := w.(io.Writer)
		if node.linkcnt > 0 {
			io.WriteString(ww, fmt.Sprintf("\"%d\" [style=filled];\n", node.State))
		} else {
			io.WriteString(ww, fmt.Sprintf("\"%d\";\n", node.State))
		}
	}, w)
	walkDAG(root, func(node *DSSNode, arg interface{}) {
		ww := w.(io.Writer)
		for _, p := range node.Preds {
			io.WriteString(ww, fmt.Sprintf("\"%d\" -> \"%d\" [label=\"%v\"];\n",
				node.State, p.State, node.Sym))
		}
	}, w)
	io.WriteString(w, "}\n")
}

// Debugging
func PrintDSS(root *DSSRoot) {
	walkDAG(root, func(node *DSSNode, arg interface{}) {
		predList := " "
		for _, p := range node.Preds {
			predList = predList + strconv.Itoa(p.State) + " "
		}
		fmt.Printf("edge  [%s]%v\n", predList, node)
	}, nil)
}

func walkDAG(root *DSSRoot, worker func(*DSSNode, interface{}), arg interface{}) {
	seen := map[*DSSNode]bool{}
	var walk func(*DSSNode, func(*DSSNode, interface{}), interface{})
	walk = func(node *DSSNode, worker func(*DSSNode, interface{}), arg interface{}) {
		if seen[node] {
			return
		}
		seen[node] = true
		for _, p := range node.Preds {
			walk(p, worker, arg)
		}
		worker(node, arg)
	}
	for _, stack := range root.stacks {
		walk(stack.tos, worker, arg)
	}
}

// --- Helpers ---------------------------------------------------------------

func contains(s []*DSSNode, sym lr.Symbol) *DSSNode {
	for _, n := range s {
		if n.Sym == sym {
			return n
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
