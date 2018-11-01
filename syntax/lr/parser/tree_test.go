package parser

import (
	"testing"
)

func TestTree1(t *testing.T) {
	tree := NewParseTree()
	node1 := tree.NewNode()
	tree.AddChild(tree.Root(), node1)
	if len(tree.Children(tree.Root())) != 1 {
		T.Debugf("children = %v", tree.Children(tree.Root()))
		t.Fail()
	}
	if tree.ParentCount(node1) != 1 {
		T.Debugf("parents = %v", tree.Parents(node1))
		t.Fail()
	}
}

func TestTree2(t *testing.T) {
	tree := NewParseTree()
	node1 := tree.NewNode()
	node2 := tree.NewNode()
	tree.AddChild(tree.Root(), node1)
	tree.AddChild(tree.Root(), node2)
	T.Debugf("children = %v", tree.Children(tree.Root()))
	T.Debugf("parents  = %v", tree.Parents(node1))
	if len(tree.Children(tree.Root())) != 2 {
		t.Fail()
	}
}

func TestTree3(t *testing.T) {
	tree := NewParseTree()
	node1 := tree.NewNode()
	tree.AddChild(tree.Root(), node1)
	tree.AddChild(tree.Root(), node1)
	T.Debugf("children = %v", tree.Children(tree.Root()))
	if len(tree.Children(tree.Root())) != 2 {
		t.Fail()
	}
}

func TestTreePack1(t *testing.T) {
	traceOn()
	tree := NewParseTree()
	node1 := tree.NewNode()
	node2 := tree.NewNode()
	node3 := tree.NewNode()
	node4 := tree.NewNode()
	tree.AddChild(tree.Root(), node1)
	tree.AddChild(node1, node3)
	tree.AddChild(node2, node4)
	tree.packNodes(node1, node2)
	T.Debugf("backlinks to %v = %v", node2, tree.edgesto[node2.ids[0]])
	T.Debugf("     IDs of node1 = %v", node1)
	T.Debugf("children of node1 = %v", tree.Children(node1))
	T.Debugf("children of node2 = %v", tree.Children(node2))
	if len(tree.Children(node1)) != 2 {
		t.Fail()
	}
}
