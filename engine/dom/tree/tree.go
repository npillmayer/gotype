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
var errInvalidFilter error = errors.New("Filter stage is invalid")

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
// some of these to perform tasks on tree nodes.
type Walker struct {
	sync.Mutex
	initial *Node     // initial node of (sub-)tree
	pipe    *pipeline // pipeline of filters to perform work on tree nodes.
}

// NewWalker creates a Walker for the initial node of a (sub-)tree.
func NewWalker(initial *Node) *Walker {
	return &Walker{initial: initial, pipe: newPipeline()}
}

// LastError returns the last error encountered during previous operations.
//
// If the previous operation generated a promise, a call to LastError()
// will block and collect the result of all spawned goroutines.
//
// If w is nil, LastError() will return nil.
/*
func (w *Walker) LastError() error {
	if w == nil {
		return nil
	}
	w.waitForCompletion()
	return w.lasterror
}
*/

// Selection returns the current selection of tree nodes (which may be nil).
//
// If the previous operation generated a promise, a call to Selection()
// will block and collect the result of all spawned goroutines.
//
// If w is nil, Selection() will return nil.
/*
func (w *Walker) Selection() []*Node {
	if w == nil {
		return nil
	}
	w.waitForCompletion()
	return w.selection
}
*/

// ResetSelection sets the current selection to an empty set.
// Does nothing if w is nil.
//
// If the previous operation generated a promise, a call to ResetSelection()
// will block and wait for the completion of all spawned goroutines.
// It will then collect the last error, clear the selection and return.
/*
func (w *Walker) ResetSelection() {
	if w != nil {
		w.waitForCompletion()
		w.resetSelection()
	}
}
*/

// Initialize and/or clear the current selection.
/*
func (w *Walker) resetSelection() {
	if w.selection == nil {
		w.selection = make([]*Node, 0, 10)
	} else {
		w.selection = w.selection[:0]
	}
}
*/

// waitForCompletion waits for all spawned goroutines to finish.
// It will then set the client-level fields to be fetched by
// Selection() and LastError().
//
// Does nothing if w is nil.
/*
func waitForCompletion(results <-chan *Node, errch <-chan error) ([]*Node, error) {
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
*/

func (w *Walker) appendFilterForTask(task workerTask, udata interface{}) {
	newFilter := newFilter(task, udata)
	if w.pipe.empty() { // quick check, may be false positive when in if-block
		// now we know the new filter might be the first one
		w.startProcessing() // this will check again, and startup if pipe empty
	}
	w.pipe.appendFilter(newFilter) // insert filter in running pipeline
}

func (w *Walker) startProcessing() {
	doStart := false
	w.pipe.RLock()
	if w.pipe.filters == nil { // no processing up to now => start with initial node
		w.pipe.pushSync(w.initial) // input is buffered, so this will return immediately
		doStart = true             // yes, we will have to start the pipeline
	}
	w.pipe.RUnlock()
	if doStart { // ok to be outside mutex as other goroutines will check pipe.empty()
		w.pipe.startProcessing() // must be outside of mutex lock
	}
}

// We need this protected by a mutex because it is outside of our control
// how often this will be called. Every client-level call to Selection(),
// ResetSelection() and LastError() will trigger this.
/*
func (w *Walker) closeErrorChannel() {
	w.Lock()
	defer w.Unlock()
	if !w.errorsClosed {
		w.errorsClosed = true
		close(w.errorch)
	}
}
*/

// start will
// - reset the selection
// - create channels
// - initialize the waitgroup for the workload
// - start the worker goroutines waiting for workload
/*
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
*/

// Walkers may decide to perform certain tasks asynchronously. This
// will result in a promise being created.
/*
type Promise struct {
	selection []*Node      // protected against premature client access
	lasterror error        // protected against premature client access
	results   <-chan *Node // results to collect
	errch     <-chan error // channel for errors from pipeline
}
*/

// Walkers may decide to perform certain tasks asynchronously. This
// will result in a promise being created.
func (w *Walker) Promise() func() ([]*Node, error) {
	errch := w.pipe.errors
	results := w.pipe.results
	signal := make(chan struct{}, 1)
	counter := &w.pipe.queuecount
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

// Predicate is a function type to match against nodes of a tree.
// Is is used as an argument for various Walker functions to
// collect a selection of nodes.
type Predicate func(*Node) (matches bool, err error)

// Whatever is a predicate to match anything. See type Predicate.
var Whatever Predicate = func(*Node) (bool, error) {
	return true, nil
}

// Impossible is a predicate to match nothing. See type Predicate.
var Impossible Predicate = func(*Node) (bool, error) {
	return false, nil
}

// ----------------------------------------------------------------------

// AncesterWith finds an ancestor matching the given predicate.
// The search does not include the start node.
//
// If w is nil, AncestorWith will return nil.
func (w *Walker) AncestorWith(predicate Predicate) *Walker {
	if w == nil {
		return nil
	}
	if w.initial == nil || predicate == nil {
		w.pipe.errors <- errInvalidFilter
	} else {
		w.appendFilterForTask(ancestorWith, predicate) // hook in this filter
	}
	return w
}

// ancestorWith searches iteratively for a node matching a predicate.
// node is at least the parent of the start node.
func ancestorWith(node *Node, udata interface{}, push func(*Node)) error {
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
	return nil // no matching ancestor found
}

// DescendentWith searches for a single descendent matching a predicate.
// The search will be performed sequentially and recursively. Search will
// abort as soon as a matching descendent has been found.
// The search does not include the start node.
//
// If w is nil, DescendentWith will return nil.
/*
func (w *Walker) DescendentWith(predicate Predicate) *Walker {
	if w == nil {
		return nil
	}
	w.ResetSelection()
	if w.initial != nil && predicate != nil {
		chcnt := w.initial.ChildCount()
		for i := 0; i < chcnt; i++ {
			ch, _ := w.initial.Child(i)
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
*/

// descendentWith searches recursively for children matching a predicate.
// It will stop searching as soon as a descendent matches.
/*
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
*/
