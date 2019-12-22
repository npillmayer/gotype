package tree

/*
BSD License

Copyright (c) 2017–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer nor the names of its contributors
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
*/

import (
	"errors"
	"sync"
)

// ErrInvalidFilter is thrown if a pipeline filter step is defunct.
var ErrInvalidFilter = errors.New("Filter stage is invalid")

// ErrEmptyTree is thrown if a Walker is called with an empty tree. Refer to
// the documentation of NewWalker() for details about this scenario.
var ErrEmptyTree = errors.New("Cannot walk empty tree")

// ErrNoMoreFiltersAccepted is thrown if a client already called Promise(), but tried to
// re-use a walker with another filter.
var ErrNoMoreFiltersAccepted = errors.New("In promise mode; will not accept new filters; use a new walker")

// Walker holds information for operating on trees: finding nodes and
// doing work on them. Clients usually create a Walker for a (sub-)tree
// to search for a selection of nodes matching certain criteria, and
// then perform some operation on this selection.
//
// A Walker will eventually return two client-level values:
// A slice of tree nodes and the last error occured.
// Often these fields are accessed through a
// Promise-object, which represents future values for the two fields.
//
// A typical usage of a Walker looks like this ("FindNodesAndDoSomething()" is
// a placeholder for a sequence of function calls, see below):
//
//    w := NewWalker(node)
//    futureResult := w.FindNodesAndDoSomething(...).Promise()
//    nodes, err := futureResult()
//
// Walker support a set of search & filter functions. Clients will chain
// some of these to perform tasks on tree nodes (see examples).
// You may think of the set of operations to form a small
// Domain Specific Language (DSL), similar in concept to JQuery.
//
// ATTENTION: Clients must call Promise() as the final link of the
// DSL expression chain, even if they do not expect the expression to
// return a non-empty set of nodes. Firstly, they need to check for errors,
// and secondly without fetching the (possibly empty) result set by calling
// the promise, the Walker may leak goroutines.
type Walker struct {
	sync.Mutex
	initial          *Node     // initial node of (sub-)tree
	pipe             *pipeline // pipeline of filters to perform work on tree nodes.
	promising        bool      // client has called Promise()
	attributeHandler AttributeHandler
}

// NewWalker creates a Walker for the initial node of a (sub-)tree.
// The first subsequent call to a node filter function will have this
// initial node as input.
//
// If initial is nil, NewWalker will return a nil-Walker, resulting
// in a NOP-pipeline of operations, resulting in an empty set of nodes
// and an error (ErrEmptyTree).
func NewWalker(initial *Node) *Walker {
	if initial == nil {
		return nil
	}
	return &Walker{initial: initial, pipe: newPipeline()}
}

// appendFilterForTask will create a new filter for a task and append
// that filter at the end of the pipeline. If processing has not
// been started yet, it will be started.
func (w *Walker) appendFilterForTask(task workerTask, udata interface{}, buflen int) error {
	if w.promising {
		return ErrNoMoreFiltersAccepted
	}
	newFilter := newFilter(task, udata, buflen)
	if w.pipe.empty() { // quick check, may be false positive when in if-block
		// now we know the new filter might be the first one
		w.startProcessing() // this will check again, and startup if pipe empty
	}
	w.pipe.appendFilter(newFilter) // insert filter in running pipeline
	return nil
}

// startProcessing should be called as soon as the first filter is inserted
// into the pipeline. It will put the initial tree node onto the front input
// channel.
func (w *Walker) startProcessing() {
	doStart := false
	w.pipe.RLock()
	if w.pipe.filters == nil { // no processing up to now => start with initial node
		w.pipe.pushSync(w.initial) // input is buffered, will return immediately
		doStart = true             // yes, we will have to start the pipeline
	}
	w.pipe.RUnlock()
	if doStart { // ok to be outside mutex as other goroutines will check pipe.empty()
		w.pipe.startProcessing() // must be outside of mutex lock
	}
}

// Promise is a future synchronisation point.
// Walkers may decide to perform certain tasks asynchronously.
// Clients will not receive the resulting node list immediately, but
// rather get handed a Promise.
// Clients will then—any time after they received the Promise—call the
// Promise (which is of function type) to receive a slice of nodes and
// a possible error value. Calling the Promise will block until all
// concurrent operations on the tree nodes have finished, i.e. it
// is a synchronization point.
func (w *Walker) Promise() func() ([]*Node, error) {
	if w == nil {
		// empty Walker => return nil set and an error
		return func() ([]*Node, error) {
			return nil, ErrEmptyTree
		}
	}
	// drain the result channel and the error channel
	w.promising = true // will block calls to establish new filters
	errch := w.pipe.errors
	results := w.pipe.results
	counter := &w.pipe.queuecount
	signal := make(chan struct{}, 1)
	var selection []*Node
	var lasterror error
	go func() {
		defer close(signal)
		selection, lasterror = waitForCompletion(results, errch, counter)
	}()
	return func() ([]*Node, error) {
		<-signal
		return selection, lasterror
	}
}

// ----------------------------------------------------------------------

// Predicate is a function type to match against nodes of a tree.
// Is is used as an argument for various Walker functions to
// collect a selection of nodes.
type Predicate func(*Node) (matches bool, err error)

// Whatever is a predicate to match anything (see type Predicate).
// It is useful to match the first node in a given direction.
var Whatever Predicate = func(*Node) (bool, error) {
	return true, nil
}

// TraverseAll is a predicate to match nothing (see type Predicate).
// It is useful to traverse a whole tree.
/*
var TraverseAll Predicate = func(*Node) (bool, error) {
	return false, nil
}
*/

// AttributeHandler supports the querying of attributes for a node.
// As we do not know the internal structure of a node's payload, we need
// help from the client.
//
// ■ GetAttribute() should return the attribute value for a given key.
//
// ■ AttributesEqual() should return true iff two values are considered equal.
//
// ■ SetAttribute() should set a new attribute value
type AttributeHandler interface {
	GetAttribute(payload interface{}, key interface{}) interface{}
	AttributesEqual(value1 interface{}, value2 interface{}) bool
	SetAttribute(payload interface{}, key interface{}, value interface{}) bool
}

// SetAttributeHandler sets an attribute getter and setter to support
// nodes' attributes. See type AttributeHandler.
func (w *Walker) SetAttributeHandler(handler AttributeHandler) {
	w.attributeHandler = handler
}

// ----------------------------------------------------------------------

// Parent returns the parent node.
//
// If w is nil, Parent will return nil.
func (w *Walker) Parent() *Walker {
	if w == nil {
		return nil
	}
	if err := w.appendFilterForTask(parent, nil, 0); err != nil {
		T().Errorf(err.Error())
		panic(err)
	}
	return w
}

// parent is a very simple filter task to retrieve the parent of a tree node.
// if the node is the tree root node, parent() will not produce a result.
func parent(node *Node, isBuffered bool, udata userdata, push func(*Node),
	pushBuf func(*Node, interface{})) error {
	//
	p := node.Parent()
	if p != nil {
		push(p) // forward parent node to next pipeline stage
	}
	return nil
}

// AncestorWith finds an ancestor matching the given predicate.
// The search does not include the start node.
//
// If w is nil, AncestorWith will return nil.
func (w *Walker) AncestorWith(predicate Predicate) *Walker {
	if w == nil {
		return nil
	}
	if predicate == nil {
		w.pipe.errors <- ErrInvalidFilter
	} else {
		err := w.appendFilterForTask(ancestorWith, predicate, 0) // hook in this filter
		if err != nil {
			T().Errorf(err.Error())
			panic(err)
		}
	}
	return w
}

// ancestorWith searches iteratively for a node matching a predicate.
// node is at least the parent of the start node.
func ancestorWith(node *Node, isBuffered bool, udata userdata, push func(*Node),
	pushBuf func(*Node, interface{})) error {
	//
	predicate := udata.filterdata.(Predicate)
	anc := node.Parent()
	for anc != nil {
		matches, err := predicate(anc)
		if err != nil {
			return err
		}
		if matches {
			push(anc) // put ancestor on output channel for next pipeline stage
			return nil
		}
	}
	return nil // no matching ancestor found, not an error
}

// DescendentsWith finds descendents matching a predicate.
// The search does not include the start node.
//
// If w is nil, DescendentsWith will return nil.
func (w *Walker) DescendentsWith(predicate Predicate) *Walker {
	if w == nil {
		return nil
	}
	if predicate == nil {
		w.pipe.errors <- ErrInvalidFilter
	} else {
		err := w.appendFilterForTask(descendentsWith, predicate, 5) // need a helper queue
		if err != nil {
			T().Errorf(err.Error())
			panic(err)
		}
	}
	return w
}

func descendentsWith(node *Node, isBuffered bool, udata userdata, push func(*Node),
	pushBuf func(*Node, interface{})) error {
	//
	if isBuffered {
		predicate := udata.filterdata.(Predicate)
		matches, err := predicate(node)
		T().Debugf("Predicate for node %s returned: %v, err=%v", node, matches, err)
		if err != nil {
			return err // do not descend further
		}
		if matches {
			push(node) // found one, put on output channel for next pipeline stage
		}
		revisitChildrenOf(node, pushBuf)
	} else {
		revisitChildrenOf(node, pushBuf)
	}
	return nil
}

func revisitChildrenOf(node *Node, pushBuf func(*Node, interface{})) {
	chcnt := node.ChildCount()
	for i := 0; i < chcnt; i++ {
		ch, _ := node.Child(i)
		pp := parentAndPosition{node, i}
		pushBuf(ch, pp)
	}
}

// AllDescendents traverses all descendents.
// The traversal does not include the start node.
// This is just a wrapper around `w.DescendentsWith(Whatever)`.
//
// If w is nil, AllDescendents will return nil.
func (w *Walker) AllDescendents() *Walker {
	return w.DescendentsWith(Whatever)
}

// AttributeIs checks a node's attributes and filters all nodes with
// their attributes not matching the requested value.
//
// If w is nil, AncestorWith will return nil.
func (w *Walker) AttributeIs(key interface{}, value interface{}) *Walker {
	if w == nil {
		return nil
	}
	if key == nil {
		w.pipe.errors <- ErrInvalidFilter
	} else {
		attr := attrInfo{w.attributeHandler, key, value}
		err := w.appendFilterForTask(attributeIs, attr, 0) // hook in this filter
		if err != nil {
			T().Errorf(err.Error())
			panic(err)
		}
	}
	return w
}

type attrInfo struct {
	handler AttributeHandler
	key     interface{}
	value   interface{}
}

// attributeIs checks an attribute of a tree node. It relies on an
// attribute handler to perform this task. The attribute handler
// has to be provided by the caller.
// nil is a valid attribute value to compare.
//
// If no attribute handler is provided, no tree node will match.
func attributeIs(node *Node, isBuffered bool, udata userdata, push func(*Node),
	pushBuf func(*Node, interface{})) error {
	//
	attr := udata.filterdata.(attrInfo)
	if attr.handler != nil {
		val := attr.handler.GetAttribute(node.Payload, attr.key)
		if attr.handler.AttributesEqual(val, attr.value) {
			push(node) // forward node to next pipeline stage
		}
	}
	return nil
}

// SetAttribute sets a node's attributes and filters all nodes
// where setting the attribute failed.
//
// If w is nil, SetAttribute will return nil.
func (w *Walker) SetAttribute(key interface{}, value interface{}) *Walker {
	if w == nil {
		return nil
	}
	if key == nil {
		w.pipe.errors <- ErrInvalidFilter
	} else {
		attr := attrInfo{w.attributeHandler, key, value}
		err := w.appendFilterForTask(attributeIs, attr, 0) // hook in this filter
		if err != nil {
			T().Errorf(err.Error())
			panic(err)
		}
	}
	return w
}

// setAttribute uses an attribute handler to set a node's attribute.
// The attribute handler has to be provided by the caller.
//
// If no attribute handler is provided, no tree node will match.
func setAttribute(node *Node, isBuffered bool, udata userdata, push func(*Node),
	pushBuf func(*Node, interface{})) error {
	//
	attr := udata.filterdata.(attrInfo)
	if attr.handler != nil {
		ok := attr.handler.SetAttribute(node.Payload, attr.key, attr.value)
		if ok {
			push(node)
		}
	}
	return nil
}

// Filter calls a client-provided function on each node of the selection.
// The user function should return the input node if it is accepted and
// nil otherwise.
//
// If w is nil, Filter will return nil.
func (w *Walker) Filter(f func(*Node) (*Node, error)) *Walker {
	if w == nil {
		return nil
	}
	if f == nil {
		w.pipe.errors <- ErrInvalidFilter
	} else {
		err := w.appendFilterForTask(clientFilter, f, 0) // hook in this filter
		if err != nil {
			T().Errorf(err.Error())
			panic(err)
		}
	}
	return w
}

func clientFilter(node *Node, isBuffered bool, udata userdata, push func(*Node),
	pushBuf func(*Node, interface{})) error {
	//
	userfunc := udata.filterdata.(func(*Node) (*Node, error))
	n, err := userfunc(node)
	if n != nil {
		push(n) // forward filtered node to next pipeline stage
	}
	return err
}

// Action is a function type to operate on tree nodes.
// Resulting nodes will be pushed to the next pipeline stage, if
// no error occured.
type Action func(n *Node, parent *Node, position int) (*Node, error)

// TopDown traverses a tree starting at (and including) the root node.
// The traversal guarantees that parents are always processed before
// their children.
//
// If the action function returns an error for a node,
// descending the branch below this node is aborted.
//
// If w is nil, TopDown will return nil.
func (w *Walker) TopDown(action Action) *Walker {
	if w == nil {
		return nil
	}
	if action == nil {
		w.pipe.errors <- ErrInvalidFilter
	} else {
		err := w.appendFilterForTask(topDown, action, 5) // need a helper queue
		if err != nil {
			T().Errorf(err.Error())
			panic(err) // TODO for debugging purposes until more mature
		}
	}
	return w
}

// ad-hoc container
type parentAndPosition struct {
	parent   *Node
	position int
}

func topDown(node *Node, isBuffered bool, udata userdata, push func(*Node),
	pushBuf func(*Node, interface{})) error {
	//
	if isBuffered {
		action := udata.filterdata.(Action)
		var parent *Node
		var position int
		if udata.nodedata != nil {
			parent = udata.nodedata.(parentAndPosition).parent
			position = udata.nodedata.(parentAndPosition).position
		}
		result, err := action(node, parent, position)
		T().Debugf("Action for node %s returned: %v, err=%v", node, result, err)
		if err != nil {
			return err // do not descend further
		}
		if result != nil {
			push(result) // result -> next pipeline stage
		}
		revisitChildrenOf(node, pushBuf) // hand over node as parent
	} else {
		pushBuf(node, nil) // simply move incoming nodes over to buffer queue
	}
	return nil
}

// TODO
func (w *Walker) BottomUp(action Action) *Walker {
	if w == nil {
		return nil
	}
	if action == nil {
		w.pipe.errors <- ErrInvalidFilter
	} else {
		err := w.appendFilterForTask(bottomUp, action, 5) // need a helper queue
		if err != nil {
			T().Errorf(err.Error())
			panic(err) // TODO for debugging purposes until more mature
		}
	}
	return w
}

// TODO find a way of operating on parent nodes just once, as soon as the las
// child has finished action(...)
// ad-hoc container
// type parentAndPosition struct {
// 	parent   *Node
// 	position int
// }

func bottomUp(node *Node, isBuffered bool, udata userdata, push func(*Node),
	pushBuf func(*Node, interface{})) error {
	//
	if isBuffered {
		action := udata.filterdata.(Action)
		var parent *Node
		var position int
		if udata.nodedata != nil {
			parent = udata.nodedata.(parentAndPosition).parent
			position = udata.nodedata.(parentAndPosition).position
		}
		result, err := action(node, parent, position)
		T().Debugf("Action for node %s returned: %v, err=%v", node, result, err)
		if err != nil {
			return err // do not descend further
		}
		if result != nil {
			push(result) // result -> next pipeline stage
		}
		revisitChildrenOf(node, pushBuf) // hand over node as parent
	} else {
		// TODO: put this into *buffered* branch
		action := udata.filterdata.(Action)
		position := 0
		parent := node.Parent()
		if parent != nil {
			position = parent.IndexOfChild(node)
		}
		result, err := action(node, parent, position)
		if err == nil && result != nil {
			push(result) // result -> next pipeline stage
		}
		if parent != nil {

		}
		pushBuf(parent, udata)
	}
	return nil
}
