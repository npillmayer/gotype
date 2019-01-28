package tree

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
)

func Test0(t *testing.T) {
	tracing.EngineTracer = gologadapter.New()
}

func TestEmptyWalker(t *testing.T) {
	n := checkRuntime(t, -1)
	w := NewWalker(nil)
	future := w.Parent().Promise()
	nodes, err := future()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("Walker for empty tree should return an error")
	}
	if len(nodes) != 0 {
		t.Errorf("result set of empty pipeline should be empty")
	}
	checkRuntime(t, n)
}

func TestParent(t *testing.T) {
	n := checkRuntime(t, -1)
	node1, node2 := NewNode(1), NewNode(2)
	node1.AddChild(node2) // simple tree: (1)-->(2)
	w := NewWalker(node2)
	future := w.Parent().Promise()
	nodes, err := future()
	if err != nil {
		t.Error(err)
	}
	if len(nodes) != 1 || !checkNodes(nodes, 1) {
		t.Errorf("did not find parent, nodes = %v", nodes)
	}
	checkRuntime(t, n)
}

func TestParentOfRoot(t *testing.T) {
	n := checkRuntime(t, -1)
	node1 := NewNode(1)
	w := NewWalker(node1)
	future := w.Parent().Promise()
	nodes, err := future()
	if err != nil {
		t.Error(err)
	}
	if len(nodes) != 0 {
		t.Errorf("root should have no parent, nodes = %v", nodes)
	}
	checkRuntime(t, n)
}

func TestAncestor(t *testing.T) {
	n := checkRuntime(t, -1)
	node1, node2 := NewNode(1), NewNode(2)
	node1.AddChild(node2) // simple tree: (1)-->(2)
	w := NewWalker(node2)
	future := w.AncestorWith(Whatever).Promise()
	nodes, err := future()
	if err != nil {
		t.Error(err)
	}
	for _, node := range nodes {
		t.Logf("ancestor = %v\n", node)
	}
	if len(nodes) != 1 || !checkNodes(nodes, 1) {
		t.Errorf("did not find single ancestor (parent), nodes = %v", nodes)
	}
	checkRuntime(t, n)
}

func TestDescendents(t *testing.T) {
	n := checkRuntime(t, -1)
	node1, node2, node3, node4 := NewNode(1), NewNode(2), NewNode(3), NewNode(4)
	// build tree:
	// (1)
	//  +---(2)
	//  |    +---(3)
	//  |
	//  +---(4)
	node1.AddChild(node2)
	node2.AddChild(node3)
	node1.AddChild(node4)
	gr3 := func(node *Node) (bool, error) {
		val := node.Payload.(int)
		return (val >= 3), nil // match nodes (3) and (4)
	}
	w := NewWalker(node1)
	future := w.DescendentsWith(gr3).Promise()
	nodes, err := future()
	if err != nil {
		t.Error(err)
	}
	for _, node := range nodes {
		t.Logf("descendent = %v\n", node)
	}
	if len(nodes) != 2 || !checkNodes(nodes, 3, 4) {
		t.Errorf("did not find descendents (3) and (4), nodes = %v", nodes)
	}
	checkRuntime(t, n)
}

func ExampleWalker_Promise() {
	// Build a tree:
	//
	//                 (root:1)
	//          (n2:2)----+----(n4:10)
	//  (n3:10)----+
	//
	// Then query for nodes with value > 5
	//
	root, n2, n3, n4 := NewNode(1), NewNode(2), NewNode(10), NewNode(10)
	root.AddChild(n2).AddChild(n4)
	n2.AddChild(n3)
	// Define our ad-hoc predicate
	greater5 := func(node *Node) (bool, error) {
		val := node.Payload.(int)
		return (val > 5), nil // match nodes with value > 5
	}
	// Now navigate the tree (concurrently)
	future := NewWalker(root).DescendentsWith(greater5).Promise()
	// Any time later call the promise ...
	nodes, err := future() // will block until walking is finished
	if err != nil {
		fmt.Print(err)
	}
	for _, node := range nodes {
		fmt.Printf("matching descendent found: (Node %d)\n", node.Payload.(int))
	}
	// Output:
	// matching descendent found: (Node 10)
	// matching descendent found: (Node 10)
}

func TestAttribute1(t *testing.T) {
	n := checkRuntime(t, -1)
	node1 := NewNode(1)
	w := NewWalker(node1)
	w.SetAttributeHandler(attr{})
	future := w.AttributeIs("num", 1).Promise()
	nodes, err := future()
	if err != nil {
		t.Error(err)
	}
	if len(nodes) != 1 {
		t.Errorf("attribute for node 1 should have been detected, nodes = %v", nodes)
	}
	checkRuntime(t, n)
}

type attr struct{}

func (a attr) GetAttribute(payload interface{}, key interface{}) interface{} {
	val := payload.(int)
	return val
}

func (a attr) AttributesEqual(val1 interface{}, val2 interface{}) bool {
	v1 := val1.(int)
	v2 := val2.(int)
	return v1 == v2
}

func (a attr) SetAttribute(payload interface{}, key interface{}, value interface{}) bool {
	return false
}

// ----------------------------------------------------------------------

// Helper to check if result nodes are the expected ones.
func checkNodes(nodes []*Node, vals ...int) bool {
	var found bool
	for _, node := range nodes {
		found = false
		v := node.Payload.(int)
		for _, val := range vals {
			if v == val {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Helper to check for leaked goroutines.
func checkRuntime(t *testing.T, N int) int {
	if N < 1 {
		return runtime.NumGoroutine()
	}
	time.Sleep(10 * time.Millisecond)
	n := runtime.NumGoroutine()
	if n > N {
		t.Logf("still %d goroutines alive", n)
	}
	return n
}
