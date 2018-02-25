package parser

import (
	"testing"
)

func TestTree1(t *testing.T) {
	traceOn()
	tree := NewParseTree()
	node1 := tree.NewNode()
	tree.AddChild(tree.Root(), node1)
	T.Debugf("children = %v", tree.Children(tree.Root()))
	if len(tree.Children(tree.Root())) != 1 {
		t.Fail()
	}
}

func TestTree2(t *testing.T) {
	traceOn()
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
	traceOn()
	tree := NewParseTree()
	node1 := tree.NewNode()
	tree.AddChild(tree.Root(), node1)
	tree.AddChild(tree.Root(), node1)
	T.Debugf("children = %v", tree.Children(tree.Root()))
	if len(tree.Children(tree.Root())) != 2 {
		t.Fail()
	}
}
