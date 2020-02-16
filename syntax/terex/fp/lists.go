package fp

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config/gtrace"

	"github.com/npillmayer/gotype/syntax/terex"
)

/*
The current implementation always pre-fetches the first value.
This could be optimized. It would be a problem with long-running ops in the
atom-creation, in case the value is never fetched by an output call.
For now, we will leave it this way.
*/

type ListSeq struct {
	atom terex.Atom
	seq  ListGenerator
}

func Seq(l *terex.GCons) ListSeq {
	var S ListGenerator
	S = func() ListSeq {
		if l == nil {
			return ListSeq{terex.NilAtom, nil}
		}
		atom := l.Car
		l = l.Cdr
		return ListSeq{atom, S}
	}
	atom := l.Car
	return ListSeq{atom, S}
}

func (seq *ListSeq) Break() {
	seq.seq = nil
}

func (seq *ListSeq) Done() bool {
	return seq.seq == nil
}

func (seq ListSeq) First() (terex.Atom, ListSeq) {
	return seq.atom, seq
}

func (seq *ListSeq) Next() terex.Atom {
	if seq.Done() {
		return terex.NilAtom
	}
	next := seq.seq()
	seq.atom = next.atom
	if seq.atom == terex.NilAtom {
		seq.seq = nil
	} else {
		seq.seq = next.seq
	}
	return seq.atom
}

type ListGenerator func() ListSeq

func NSeq() ListSeq {
	var n int64
	var S ListGenerator
	S = func() ListSeq {
		n++
		atom := terex.Atomize(n)
		return ListSeq{atom, S}
	}
	atom := terex.Atomize(n)
	return ListSeq{atom, S}
}

type ListMapper func(terex.Atom) terex.Atom

func (seq ListSeq) Map(mapper ListMapper) ListSeq {
	var F ListGenerator
	//inner := seq
	atom, inner := seq.atom, seq
	//n, inner := seq.First()
	v := mapper(atom)
	F = func() ListSeq {
		//fmt.Printf("F  called, n=%d\n", n)
		atom = inner.Next()
		v = mapper(atom)
		//fmt.Printf("F' n=%d, v=%d\n", n, v)
		return ListSeq{v, F}
	}
	return ListSeq{v, F}
}

func (seq ListSeq) List() *terex.GCons {
	if seq.Done() {
		return nil
	}
	var start, end *terex.GCons
	//atom, S := seq.First()
	//fmt.Printf("first atom=%s\n", atom)
	S := seq
	// var atom terex.Atom
	for atom := seq.Next(); !S.Done(); atom = S.Next() {
		fmt.Printf("next atom=%s, S=%v\n", atom, S)
		fmt.Printf("  done=%v\n", S.Done())
		if start == nil {
			start = terex.Cons(atom, nil)
			end = start
		} else {
			end.Cdr = terex.Cons(atom, nil)
			end = end.Cdr
		}
		fmt.Printf("result list = %s\n", start.ListString())
	}
	return start
}

// --- Trees -----------------------------------------------------------------

type TreeSeq struct {
	node    *terex.GCons
	channel <-chan *terex.GCons
	seq     TreeGenerator
}

type TreeGenerator func() TreeSeq

type treeTraverser []*terex.GCons

/*
Tree creates a sequence from a TeREx tree structure. The sequence traverses the
tree in depth-first post-order.

https://www.geeksforgeeks.org/iterative-postorder-traversal-using-stack/

	1.1 Create an empty stack
	2.1 Do following while root is not NULL
		a) Push root's right child and then root to stack.
		b) Set root as root's left child.
	2.2 Pop an item from stack and set it as root.
		a) If the popped item has a right child and the right child
		is at top of stack, then remove the right child from stack,
		push the root back and set root as root's right child.
		b) Else print root's data and set root as NULL.
	2.3 Repeat steps 2.1 and 2.2 while stack is not empty.

For TeREx' pre-order, a node's content is Car, left child is Cdar, right child is Cddr.
The tree from the above website's example

    1---2---4
    |   \---5
    \---3---6
	    \---7

is

	(1 (2 (4) 5) 3 (6) 7)

in TeREx pre-order format. A depth-first traversal will yield

	(4 5 2 6 7 3 1)

*/
func TreeIteratorCh(l *terex.GCons) <-chan *terex.GCons {
	// 1.1 Create an empty stack
	t := make([]*terex.GCons, 0, 32)
	if l == nil {
		return nil
	}
	channel := make(chan *terex.GCons)
	go func(l *terex.GCons) {
		defer close(channel)
		node := l // set root
		for {
			// 2.1 Do following while root is not NULL
			for node != nil {
				left, right := children(node)
				// a) Push root's right child and then root to stack.
				if right != nil {
					t = append(t, right) // push right child node
				}
				t = append(t, node) // push node
				// b) Set root as root's left child.
				node = left
			} // now node == nil
			// 2.2 Pop an item from stack and set it as root.
			node, t = t[len(t)-1], t[:len(t)-1]
			_, right := children(node)
			if len(t) > 0 && right != nil && right == t[len(t)-1] {
				// a) If the popped item has a right child and the right child
				// is at top of stack, then remove the right child from stack,
				// push the root back and set root as root's right child.
				t = t[:len(t)-1]    // pop right child
				t = append(t, node) // push root
				node = right        // root <- right child
			} else {
				// b) Else print root's data and set root as NULL.
				gtrace.SyntaxTracer.Debugf("Node=%s", node)
				channel <- node
				node = nil
			}
			// 2.3 Repeat steps 2.1 and 2.2 while stack is not empty.
			if len(t) == 0 {
				break
			}
		}
	}(l)
	return channel
}

func Traverse(l *terex.GCons) TreeSeq {
	channel := TreeIteratorCh(l)
	if channel == nil {
		return TreeSeq{}
	}
	var T TreeGenerator
	T = func() TreeSeq {
		var ok bool
		tseq := TreeSeq{nil, channel, T}
		if tseq.node, ok = <-channel; !ok {
			tseq.seq = nil
		}
		return tseq
	}
	var ok bool
	var node *terex.GCons
	tseq := TreeSeq{node, channel, T}
	if tseq.node, ok = <-channel; !ok {
		tseq.seq = nil
	}
	return tseq
}

func children(node *terex.GCons) (*terex.GCons, *terex.GCons) {
	if node == nil {
		return nil, nil
	}
	if node.Car.Type() == terex.ConsType {
		// anonymous node
		panic("anonymous nodes not yet implemented")
	}
	if node.Cdr == nil {
		return nil, nil
	}
	left := node.Cdr.Tee()
	right := node.Cddr()
	return left, right
}

func (t treeTraverser) printStack() {
	for i, n := range t {
		gtrace.SyntaxTracer.Debugf("   [%d] %s", i, terex.Elem(n).String())
	}
}

func (seq *TreeSeq) Break() {
	seq.seq = nil
}

func (seq *TreeSeq) Done() bool {
	return seq.seq == nil
}

func (seq TreeSeq) First() (*terex.GCons, TreeSeq) {
	return seq.node, seq
}

func (seq *TreeSeq) Next() *terex.GCons {
	if seq.Done() {
		return nil
	}
	next := seq.seq()
	node := next.node
	seq.seq = next.seq
	return node
}

func (seq TreeSeq) List() *terex.GCons {
	if seq.Done() {
		return nil
	}
	var start, end *terex.GCons
	for node, T := seq.First(); !T.Done(); node = T.Next() {
		if start == nil {
			start = terex.Cons(node.Car, nil)
			end = start
		} else {
			end.Cdr = terex.Cons(node.Car, nil)
			end = end.Cdr
		}
	}
	return start
}
