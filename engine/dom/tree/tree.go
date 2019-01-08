package tree

/*
BSD License

Copyright (c) 2017–18, Norbert Pillmayer

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

// Error to be emitted if a pipeline filter step is defunct.
var ErrInvalidFilter error = errors.New("Filter stage is invalid")

// Error to be emitted if a Walker is called with an empty tree
var ErrEmptyTree error = errors.New("Cannot walk empty tree")

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
	initial         *Node     // initial node of (sub-)tree
	pipe            *pipeline // pipeline of filters to perform work on tree nodes.
	attributeGetter AttributeGetter
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
func (w *Walker) appendFilterForTask(task workerTask, udata interface{}, buflen int) {
	newFilter := newFilter(task, udata, buflen)
	if w.pipe.empty() { // quick check, may be false positive when in if-block
		// now we know the new filter might be the first one
		w.startProcessing() // this will check again, and startup if pipe empty
	}
	w.pipe.appendFilter(newFilter) // insert filter in running pipeline
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

// Walkers may decide to perform certain tasks asynchronously.
// Clients will not receive the resulting node list immediately, but
// rather get handed a Promise.
// Clients will then—any time after they received the Promise—call the
// Promise (which is a function type) to receive a slice of nodes and
// a possible error value. Calling the Promise will block until all
// concurrent operations on the tree nodes have finished, i.e. it
// is a synchronization point.
func (w *Walker) Promise() func() ([]*Node, error) {
	if w == nil {
		// empty Walker => return nil set and an error
		return func() ([]*Node, error) {
			return nil, ErrEmptyTree
		}
	} else {
		// drain the result channel and the error channel
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
var TraverseAll Predicate = func(*Node) (bool, error) {
	return false, nil
}

// AttributeGetter supports the querying of attributes for a node.
// As we do not know the internal structure of a node's payload, we need
// help from the client.
//
// GetAttribute() should return the attribute value for a given key.
// AttributesEqual() should return true iff two values are considered equal.
type AttributeGetter interface {
	GetAttribute(payload interface{}, key interface{}) interface{}
	AttributesEqual(value1 interface{}, value2 interface{}) bool
}

// SetAttributeGetter sets an attribute getter to support querying of
// nodes' attributes.
func (w *Walker) SetAttributeGetter(getter AttributeGetter) {
	w.attributeGetter = getter
}

// ----------------------------------------------------------------------

// Parent returns the parent node.
//
// If w is nil, Parent will return nil.
func (w *Walker) Parent() *Walker {
	if w == nil {
		return nil
	}
	w.appendFilterForTask(parent, nil, 0)
	return w
}

// parent is a very simple filter task to retrieve the parent of a tree node.
// if the node is the tree root node, parent() will not produce a result.
func parent(node *Node, isBuffered bool, udata interface{}, push func(*Node),
	pushBuf func(*Node)) error {
	//
	p := node.Parent()
	if p != nil {
		push(p) // forward parent node to next pipeline stage
	}
	return nil
}

// AncesterWith finds an ancestor matching the given predicate.
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
		w.appendFilterForTask(ancestorWith, predicate, 0) // hook in this filter
	}
	return w
}

// ancestorWith searches iteratively for a node matching a predicate.
// node is at least the parent of the start node.
func ancestorWith(node *Node, isBuffered bool, udata interface{}, push func(*Node),
	pushBuf func(*Node)) error {
	//
	predicate := udata.(Predicate)
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
		w.appendFilterForTask(descendentsWith, predicate, 5) // need a helper queue
	}
	return w
}

func descendentsWith(node *Node, isBuffered bool, udata interface{}, push func(*Node),
	pushBuf func(*Node)) error {
	//
	if isBuffered {
		predicate := udata.(Predicate)
		matches, err := predicate(node)
		if err != nil {
			return err // do not descend further
		}
		if matches {
			push(node) // found one, put on output channel for next pipeline stage
		} else {
			revisitChildrenOf(node, pushBuf)
		}
	} else {
		revisitChildrenOf(node, pushBuf)
	}
	return nil
}

func revisitChildrenOf(node *Node, pushBuf func(*Node)) {
	chcnt := node.ChildCount()
	for i := 0; i < chcnt; i++ {
		ch, _ := node.Child(i)
		pushBuf(ch)
	}
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
		attr := attrMatcher{w.attributeGetter, key, value}
		w.appendFilterForTask(attributeIs, attr, 0) // hook in this filter
	}
	return w
}

type attrMatcher struct {
	getter AttributeGetter
	key    interface{}
	value  interface{}
}

// attributeIs checks an attribute of a tree node. It relies on an
// attribute getter function to perform this task. The attribute getter
// has to be provided by the caller.
// nil is a valid attribute value to compare.
//
// If no attribute getter is provided, no tree node will match.
func attributeIs(node *Node, isBuffered bool, udata interface{}, push func(*Node),
	pushBuf func(*Node)) error {
	//
	attr := udata.(attrMatcher)
	if attr.getter != nil {
		val := attr.getter.GetAttribute(node.Payload, attr.key)
		if attr.getter.AttributesEqual(val, attr.value) {
			push(node) // forward node to next pipeline stage
		}
	}
	return nil
}
