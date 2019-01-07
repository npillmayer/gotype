package tree

import (
	"testing"
)

func TestAncestor(t *testing.T) {
	node1, node2 := NewNode(1), NewNode(2)
	node1.AddChild(node2) // simple tree: (1)-->(2)
	w := NewWalker(node2)
	anc := w.AncestorWith(Whatever).Promise()
	t.Log("getting a promise() for an ancestor")
	nodes, err := anc()
	if err != nil {
		t.Error(err)
	}
	for _, n := range nodes {
		t.Logf("ancestor = %v\n", n)
	}
}
