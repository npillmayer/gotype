package dagstack

import (
	"fmt"
	"io"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/lr"
)

// Configurable trace
var T tracing.Trace = tracing.SyntaxTracer

// Node in a DAG stack
type DSNode struct {
	State   int       // identifier of a parse state
	Sym     lr.Symbol // symbol on the stack (between states)
	Preds   []*DSNode // predecessors of this node (join)
	Succs   []*DSNode // successors of this node (fork)
	linkcnt int       // count of stacks referencing this as TOS
}

// create an empty stack node
func newDSNode(state int, sym lr.Symbol) *DSNode {
	n := &DSNode{}
	n.State = state
	n.Sym = sym
	return n
}

// Simple Stringer for a stack node.
func (n *DSNode) String() string {
	return fmt.Sprintf("<-(%v)-[%d]", n.Sym, n.State)
}

func (n *DSNode) append(succ *DSNode) {
	if n.Succs == nil {
		n.Succs = make([]*DSNode, 0, 3)
	}
	n.Succs = append(n.Succs, succ)
}

func (n *DSNode) prepend(pred *DSNode) {
	if n.Preds == nil {
		n.Preds = make([]*DSNode, 0, 3)
	}
	n.Preds = append(n.Preds, pred)
}

func (n *DSNode) followedBy(sym lr.Symbol) *DSNode {
	s := contains(n.Succs, sym)
	return s
}

/*
Root for a DAG-stack. All stacks of a DAG stack structure share a common
root.
*/
type DSRoot struct {
	roots  *DSNode     // points to a chain of sibling root nodes
	toss   []*DSNode   // list of all TOSs
	bottom *DSNode     // stopper
	name   string      // identifier of the root
	stacks []*DAGStack // housekeeping
}

// Create a named root for a DAG-stack.
func DAGStackRoot(name string) *DSRoot {
	root := &DSRoot{name: name}
	root.bottom = newDSNode(-999, &pseudo{"bottom"})
	root.stacks = make([]*DAGStack, 0, 10)
	return root
}

/*
A data structure for a stack. All stacks together form a DAG, i.e. they
may share portions of stacks. Each client carries its own stack, without
noticing the other stacks. The data structure hides the fact that stack-
fragments may be shared.
*/
type DAGStack struct {
	root   *DSRoot
	tos    *DSNode
	height int
}

func NewDAGStack(root *DSRoot) *DAGStack {
	s := &DAGStack{root: root, tos: root.bottom, height: 0}
	s.root.stacks = append(s.root.stacks, s)
	return s
}

func (stack *DAGStack) Size() int {
	return stack.height
}

func (stack *DAGStack) Fork() *DAGStack {
	// create a new stack clone
	// update all linkcounts down to root, if necessary
	return nil
}

func (stack *DAGStack) moveUp(succ *DSNode) int {
	// successor of TOS must exist as succ
	stack.tos.linkcnt--
	stack.tos = succ
	stack.tos.linkcnt++
	stack.height++
	return stack.tos.State
}

func (stack *DAGStack) Push(state int, sym lr.Symbol) {
	// create or find a node
	// - create: new and let node.prev be tos
	// - find: return existing one
	// update references / linkcnt
	if succ := stack.tos.followedBy(sym); succ != nil {
		stack.moveUp(succ)
	} else {
		succ := newDSNode(state, sym)
		succ.prepend(stack.tos)
		stack.tos.append(succ)
		stack.moveUp(succ)
	}
}

func (stack *DAGStack) TOS() (lr.Symbol, int) {
	return stack.tos.Sym, stack.tos.State
}

func (stack *DAGStack) Pop() lr.Symbol {
	// if tos is part of another chain: return node and go up one node
	// if shared by another stack: return node and go up one node
	// if not shared: remove node
	return nil
}

func (stack *DAGStack) Reduce(handleLen int, lhs lr.Symbol) lr.Symbol {
	// find join(s) and split(s) within handleLen
	// if merged: de-merge handleLen symbols
	// remove handle symbols
	// push lhs symbol
	// return new top symbol (should be lhs)
	return nil
}

// --- Debugging -------------------------------------------------------------

func Stack2Dot(root *DSRoot, w io.Writer) {
	io.WriteString(w, "digraph {\n")
	for _, stack := range root.stacks {
		tos := stack.tos
		for tos != stack.root.bottom {
			T.Debugf("TOS = %v", tos)
			for _, p := range tos.Preds {
				io.WriteString(w, fmt.Sprintf("%d -> %d [label=\"%v\"];", tos.State, p.State, tos.Sym))
			}
			tos = tos.Preds[0]
		}
	}
	io.WriteString(w, "}\n")
}

// --- Helpers ---------------------------------------------------------------

func contains(s []*DSNode, sym lr.Symbol) *DSNode {
	for _, n := range s {
		if n.Sym == sym {
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
