/*
Package tree implements an all-purpose tree type.

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
package tree

import (
	"log"
	"sync"
)

// Node is the base type our tree is built of.
type Node struct {
	parent   *Node
	children childrenSlice
	Payload  interface{}
}

// NewNode creates a new tree node with a payload.
func NewNode(payload interface{}) *Node {
	return &Node{Payload: payload}
}

// AddChild inserts a new child node into the tree.
// The newly inserted node is connected to this node as its parent.
//
// This operation is concurrency-safe.
func (sn *Node) AddChild(ch *Node) {
	if ch == nil {
		return
	}
	sn.children.addChild(ch, sn)
}

// ParentNode returns the parent node or nil (for the root of the tree).
func (sn Node) Parent() *Node {
	return sn.parent
}

// ChildCount returns the number of children-nodes for a styled node
// (concurrency-safe).
func (sn Node) ChildCount() int {
	return sn.children.length()
}

// Child is a concurrency-safe way to get a children-node of a styled node.
func (sn Node) Child(n int) (*Node, bool) {
	if sn.children.length() <= n {
		return nil, false
	}
	return sn.children.child(n), true
}

// ----------------------------------------------------------------------

// Walker holds information for operating on trees: finding nodes and
// doing work on them. Clients usually create a Walker for a (sub-)tree
// to search for a selection of nodes matching certain criteria, and
// then perform some operation on this selection.
//
// A Walker has two client-level fields, accessed by getter methods:
// Selection and LastError.
type Walker struct {
	sync.Mutex
	root         *Node      // root node of (sub-)tree
	selection    []*Node    // protected against premature client access
	lasterror    error      // protected against premature client access
	promises     *promise   // pipeline of chained promises
	errorch      chan error // channel for errors from workers
	errorsClosed bool       // is the error channel closed, i.e., all workers done?
}

// NewWalker creates a Walker for the root node of a (sub-)tree.
func NewWalker(root *Node) *Walker {
	return &Walker{root: root}
}

// LastError returns the last error encountered during previous operations.
//
// If the previous operation generated a promise, a call to LastError()
// will block and collect the result of all spawned goroutines.
//
// If w is nil, LastError() will return nil.
func (w *Walker) LastError() error {
	if w == nil {
		return nil
	}
	w.waitForCompletion()
	return w.lasterror
}

// Selection returns the current selection of tree nodes (which may be nil).
//
// If the previous operation generated a promise, a call to Selection()
// will block and collect the result of all spawned goroutines.
//
// If w is nil, Selection() will return nil.
func (w *Walker) Selection() []*Node {
	if w == nil {
		return nil
	}
	w.waitForCompletion()
	return w.selection
}

// ResetSelection sets the current selection to an empty set.
// Does nothing if w is nil.
//
// If the previous operation generated a promise, a call to ResetSelection()
// will block and wait for the completion of all spawned goroutines.
// It will then collect the last error, clear the selection and return.
func (w *Walker) ResetSelection() {
	if w != nil {
		w.waitForCompletion()
		w.resetSelection()
	}
}

// Initialize and/or clear the current selection.
func (w *Walker) resetSelection() {
	if w.selection == nil {
		w.selection = make([]*Node, 0, 10)
	} else {
		w.selection = w.selection[:0]
	}
}

// waitForCompletion waits for all spawned goroutines to finish.
// It will then set the client-level fields to be fetched by
// Selection() and LastError().
//
// Does nothing if w is nil.
func (w *Walker) waitForCompletion() {
	if w == nil {
		return
	}
	go func() {
		w.workers.queuecount.Wait()          // wait for workload queue to become empty
		log.Printf("all workers are done\n") // TODO
		w.closeErrorChannel()
	}()
	w.lasterror = nil
	for err := range w.errorch {
		if err != nil {
			w.lasterror = err // throw away all errors but the last one
		}
	}
	// Now collect all results from the pipeline of promises
	finalPromise := w.promises[len(w.promises)-1]
	for node := range finalPromise.results {
		w.selection = append(w.selection, node)
	}
	w.cleanupPromises()
}

// We need this protected by a mutex because it is outside of our control
// how often this will be called. Every client-level call to Selection(),
// ResetSelection() and LastError() will trigger this.
func (w *Walker) closeErrorChannel() {
	w.Lock()
	defer w.Unlock()
	if !w.errorsClosed {
		w.errorsClosed = true
		close(w.errorch)
	}
}

// start will
// - reset the selection
// - create channels
// - initialize the waitgroup for the workload
// - start the worker goroutines waiting for workload
func (w *Walker) start() {
	w.Lock()
	defer w.Unlock()
	if w.errorsClosed {
		w.resetSelection()
		w.errorch = make(chan error, 10)
		w.workers.queuecount = sync.WaitGroup{}
		w.workers.workload = make(chan workPackage)
		w.workers.errorch = w.errorch // write-only copy for workers
		w.errorsClosed = false
		for i := 0; i < 2; i++ { // TODO where to define # of workers ?
			wno := i + 1 // going to start worker #wno
			go func(workload <-chan workPackage, errch chan<- error) {
				defer func() {
					log.Printf("finished worker #%d\n", wno)
				}()
				for wp := range workload { // get workpackages until drained
					node, err := wp.todo(wp.data) // perform task on workpackage
					if err != nil {
						errch <- err // signal error to Walker
					}
					if node != nil {
						wp.promise.results <- node
					}
					w.workers.queuecount.Done() // worker has finished a workpackage
				}
			}(w.workers.workload, w.errorch)
		}
	}
}

type pipe struct {
	workload   chan *Node     // data pipe
	queuecount sync.WaitGroup // count of work packages
}

func (p *pipe) push(node *Node) {
	p.queuecount.Add(1)
	go func(wp workPackage) {
		p.workload <- wp
	}(wp)
}

// Walkers may decide to perform certain tasks asynchronously. This
// will result in a promise being created and handed to the succeeding
// task in a pipeline.
type promise struct {
	workers workergroup // a pool of worker goroutines
	next    *promise    // the next promise in a pipeline
}

func makePromise(workercnt int, workload <-chan *Node, errch chan<- error) *promise {
	p := &promise{}
	p.workers.queuecount = sync.WaitGroup{} // should be unnecessary
	p.workers.workload = workload           // promise gets input from here...
	p.workers.results = make(chan *Node)    // ... and puts results here
	p.workers.errorch = errch               // this is were error messages go to
	return p
}

func (w *Walker) cleanupPromises() {
	for _, p := range w.promises {
		close(p.results)
	}
}

func (p *promise) order(workers *workergroup, task workerTask, data interface{}) {
	wp := workPackage{
		todo:    task,
		data:    data,
		promise: p,
	}
	workers.queuecount.Add(1) // must be before put
	go func(wp workPackage) {
		workers.workload <- wp
	}(wp)
}

type nodePredicateTask struct {
	predicate Predicate
}

func matchNode(wpData interface{}) (*Node, error) {
	data := wpData.(*nodePredicateTask)
	var matches bool
	var err error
	if matches, err = data.predicate(data.node); matches {
		return data.node, err
	}
	return nil, err
}

// Predicate is a function type to match against nodes of a tree.
// Is is used as an argument for various Walker functions to
// collect a selection of nodes.
type Predicate func(*Node) (matches bool, err error)

// Whatever is a predicate to match anything. See type Predicate.
func Whatever(*Node) (bool, error) {
	return true, nil
}

// Impossible is a predicate to match nothing. See type Predicate.
func Impossible(*Node) (bool, error) {
	return false, nil
}

// AncesterWith finds an ancestor matching the given predicate.
// The search does not include the start node.
//
// If w is nil, AncestorWith will return nil.
func (w *Walker) AncestorWith(predicate Predicate) *Walker {
	if w == nil {
		return nil
	}
	if w.root != nil && predicate != nil {
		promise := makePromise()
		w.promises = append(w.promises, promise) // TODO make conc-safe
		go func(pre Predicate, output chan<- *Node, errch chan<- error) {
			node, err := ancestorWith(w.root, pre)
			if err != nil {
				errch <- err // may block
			}
			if node != nil {
				output <- node // may block
			}
		}(predicate, promise.results, w.workers.errorch)
	}
	return w
}

// ancestorWith searches recursively for a node matching a predicate.
// node is at least the parent of the start node.
func ancestorWith(node *Node, predicate Predicate) (anc *Node, err error) {
	matched := false
	matched, err = predicate(node)
	if err == nil {
		if matched {
			anc = node
		} else if parent := node.Parent(); parent != nil {
			anc, err = ancestorWith(node, predicate)
		}
	}
	return
}

// DescendentWith searches for a single descendent matching a predicate.
// The search will be performed sequentially and recursively. Search will
// abort as soon as a matching descendent has been found.
// The search does not include the start node.
//
// If w is nil, DescendentWith will return nil.
func (w *Walker) DescendentWith(predicate Predicate) *Walker {
	if w == nil {
		return nil
	}
	w.ResetSelection()
	if w.root != nil && predicate != nil {
		chcnt := w.root.ChildCount()
		for i := 0; i < chcnt; i++ {
			ch, _ := w.root.Child(i)
			node, err := descendentWith(ch, predicate)
			if err != nil {
				w.lasterror = err
				break
			}
			if node != nil {
				w.selection = append(w.selection, node)
				break
			}
		}
	}
	return w
}

// descendentWith searches recursively for children matching a predicate.
// It will stop searching as soon as a descendent matches.
func descendentWith(node *Node, predicate Predicate) (desc *Node, err error) {
	matched := false
	matched, err = predicate(node)
	if err == nil {
		if matched {
			desc = node
		} else {
			chcnt := node.ChildCount()
			for i := 0; i < chcnt; i++ {
				ch, _ := node.Child(i)
				desc, err = descendentWith(ch, predicate)
				if err != nil || node != nil {
					break
				}
			}
		}
	}
	return
}

// DescendentsWith searches a sub-tree for nodes matching a predicate and
// collects these nodes.
//
// The function may decide to spawn goroutines for traversing the tree.
// In this case, the returned Walker will contain a (transparentyl wrapped)
// promise.
//
// If w is nil, DescendentWith will return nil.
func (w *Walker) DescendentsWith(predicate Predicate) *Walker {
	if w == nil {
		return nil
	}
	w.ResetSelection()
	if w.root != nil && predicate != nil {
		chcnt := w.root.ChildCount()
		if chcnt > 0 {
			for i := 0; i < chcnt; i++ {
				ch, _ := w.root.Child(i)
				node, err := descendentWith(ch, predicate)
				if err != nil {
					w.lasterror = err
					break
				}
				if node != nil {
					w.selection = append(w.selection, node)
					break
				}
			}
		}
	}
	return w
}

// --- Slices of concurrency-safe sets of children ----------------------

type childrenSlice struct {
	sync.RWMutex
	slice []*Node
}

func (chs *childrenSlice) length() int {
	chs.RLock()
	defer chs.RUnlock()
	return len(chs.slice)
}

func (chs *childrenSlice) addChild(child *Node, parent *Node) {
	if child == nil {
		return
	}
	chs.Lock()
	defer chs.Unlock()
	chs.slice = append(chs.slice, child)
	child.parent = parent
}

func (chs *childrenSlice) child(n int) *Node {
	if chs.length() == 0 || n < 0 || n >= chs.length() {
		return nil
	}
	chs.RLock()
	defer chs.RUnlock()
	return chs.slice[n]
}
