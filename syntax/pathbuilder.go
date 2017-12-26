package syntax

import (
	//  dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/stacks/linkedliststack"
	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
)

// The type we will push onto the stack
type PathNode struct {
	Symbol Symbol
	Path   arithm.Path
	Pair   arithm.Pair
}

type PathStack struct {
	stack *linkedliststack.Stack // a stack of paths
	//Q     *dll.List
}

/* Create a new path stack. It is fully initialized and empty.
 */
func NewPathStack() *PathStack {
	pst := &PathStack{
		stack: linkedliststack.New(), // stack of interface{}
	}
	return pst
}

/* Stack functionality. Will return an invalid nullpath if stack is empty.
 */
func (ps *PathStack) Top() *PathNode {
	tos, ok := ps.stack.Peek()
	if !ok {
		tos = &PathNode{
			Symbol: nil,
			Path:   nil,
			Pair:   nil,
		}
	}
	return tos.(*PathNode)
}

/* Stack functionality.
 */
func (ps *PathStack) Pop() (*PathNode, bool) {
	tos, ok := ps.stack.Pop()
	return tos.(*PathNode), ok
}

/* Stack functionality.
 */
func (ps *PathStack) Push(pn *PathNode) *PathStack {
	ps.stack.Push(pn)
	//T.Debugf("TOS is now %v", ps.Top())
	return ps
}

/* Push a path variable. Both arguments may be nil.
 */
func (ps *PathStack) PushPath(sym Symbol, path arithm.Path) {
	pn := &PathNode{
		Symbol: sym,
		Path:   path,
	}
	ps.Push(pn)
}

/* Push a pair variable. Both arguments may be nil.
 */
func (ps *PathStack) PushPair(sym Symbol, pr arithm.Pair) {
	pn := &PathNode{
		Symbol: sym,
		Pair:   pr,
	}
	ps.Push(pn)
}

/* Stack functionality.
 */
func (ps *PathStack) IsEmpty() bool {
	return ps.stack.Empty()
}

/* Stack functionality.
 */
func (ps *PathStack) Size() int {
	return ps.stack.Size()
}

// === Operations on Paths ===================================================
