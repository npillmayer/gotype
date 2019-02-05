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
	"runtime"
	"sync"
)

// Tree operations will be carried out by concurrent worker goroutines.
// As tree operations may be chained, a pipeline of filter stages is
// constructed. Every chained operation is reflected by a filter stage.
// Filters read Nodes from an input channel and put processed Nodes on
// an output channel. This way we create a little pipes&filter design.
//
// Filter stages operate concurrently. Every filter is free to launch
// as many worker goroutines as it sees fit. An overall counter is used
// to track the number of active work-packages (i.e. Nodes) in the
// pipeline. As soon as the number of nodes is zero, all channels (pipes)
// are closed and the workers will terminate.
//
// Every filter performs a specific task, reflected by a workerTask function.
// Filter tasks may use additional data, which may be provided as an
// untyped udata ("user data") argument. Filter task functions are responsible
// for decoding their specific udata.
// Errors occuring in filter tasks will be sent to a pipeline-global error
// channel.

// Minimum and maximum number of concurrent workers for a tree operation
// (filter).
const (
	minWorkerCount int = 3
	maxWorkerCount int = 10
)

// Maxmimum length of internal buffer channel for a filter.
const maxBufferLength int = 128

// Workers will be tasked a series of workerTasks.
type workerTask func(*Node, bool, userdata, func(*Node), func(*Node, interface{})) error

// filter is part of a pipeline, i.e. a stage of the overall pipeline to
// process input (Nodes) and produce results (Nodes).
// filters will perform concurrently.
type filter struct {
	next       *filter          // filters are chained
	results    chan<- *Node     // results of this filter (pipeline stage)
	queue      chan nodeSupport // helper queue if necessary
	task       workerTask       // the task this filter performs
	filterdata interface{}      // additional information needed to perform task
	env        *filterenv       // connection to outside world
}

// nodeSupport is a small helper type to let clients store arbitrary
// user data together with nodes in a buffer queue.
type nodeSupport struct {
	node     *Node       // buffered node
	nodedata interface{} // arbitrary user dataarbitrary user data
}

// filterenv holds information about the outside world to be referenced by
// a filter. This includes input workload, error destination and a counter
// for overall work on an pipeline.
type filterenv struct {
	input        <-chan *Node    // work to do for this filter, connected to predecessor
	errors       chan<- error    // where errors are reported to
	queuecounter *sync.WaitGroup // counter for overall work load
}

// userdata is a container to provide user data for both information global
// to a filter, as well as information companying a single node.
// The user data information will be provided to filter actions.
type userdata struct {
	filterdata interface{}
	nodedata   interface{}
}

// newFilter creates a new pipeline stage, i.e. a filter fed from an input
// channel (workload) and putting processed nodes into an output channel (results).
// Errors are reported to an error channel.
func newFilter(task workerTask, filterdata interface{}, buflen int) *filter {
	f := &filter{}
	if buflen > 0 {
		if buflen > maxBufferLength {
			buflen = maxBufferLength
		}
		f.queue = make(chan nodeSupport, buflen)
	}
	f.task = task
	f.filterdata = filterdata
	return f
}

// This method signature is a bit strange, but for now it does the job.
// Sets an environment for a filter an gets the results-channel in return.
func (f *filter) start(env *filterenv) chan *Node {
	f.env = env
	res := make(chan *Node, 3) // output channel has to be in place before workers start
	f.results = res            // be careful to set write-only for the filter
	n := runtime.NumCPU()
	if n > maxWorkerCount {
		n = maxWorkerCount
	} else if n < minWorkerCount {
		n = minWorkerCount
	}
	for i := 0; i < n; i++ {
		wno := i + 1
		if f.queue == nil {
			go filterWorker(f, wno) // startup worker no. #wno
		} else {
			go filterWorkerWithQueue(f, wno) // startup worker no. #wno
		}
	}
	return res // needed r/w for next filter in pipe
}

// filterWorker is the default worker function. Each filter is free to start
// as many of them as seems adequate.
func filterWorker(f *filter, wno int) {
	//  defer func() {
	//	log.Printf("finished worker #%d\n", wno) // TODO eliminate this
	//}()
	push := func(node *Node) { // worker will use this to hand result to next stage
		f.pushResult(node)
	}
	for inNode := range f.env.input { // get workpackages until drained
		udata := userdata{f.filterdata, nil}
		err := f.task(inNode, false, udata, push, nil) // perform task on workpackage
		if err != nil {
			f.env.errors <- err // signal error to caller
		}
		f.env.queuecounter.Done() // worker has finished a workpackage
	}
}

// filterWorkerWithQueue is a worker function which uses a separate support
// queue, the 'buffer queue'.
func filterWorkerWithQueue(f *filter, wno int) {
	push := func(node *Node) { // worker will use this to hand result to next stage
		f.pushResult(node)
	}
	pushBuf := func(sup *Node, udata interface{}) { // worker will use this to queue work internally
		f.pushBuffer(sup, udata)
	}
	var buffered bool
	var node *Node
	var udata userdata
	for {
		select { // get upstream workpackages and buffered workpackages until drained
		case node = <-f.env.input:
			udata.filterdata = f.filterdata
			buffered = false
		case supdata := <-f.queue:
			node = supdata.node
			udata.filterdata = f.filterdata
			udata.nodedata = supdata.nodedata
			buffered = true
		}
		if node != nil {
			err := f.task(node, buffered, udata, push, pushBuf) // perform filter task
			if err != nil {
				f.env.errors <- err // signal error to caller
			}
			f.env.queuecounter.Done() // worker has finished a workpackage
		} else {
			break // no more work to do
		}
	}
}

// pipeline is a chain of filters to perform tasks on Nodes.
// Filters, i.e., pipeline stages are connected by channels.
type pipeline struct {
	sync.RWMutex                // to sychronize access to various fields
	queuecount   sync.WaitGroup // overall count of work packages
	errors       chan error     // collector channel for error messages
	filters      *filter        // chain of filters
	input        chan *Node     // initial workload
	results      chan *Node     // where final output of this pipeline goes to
	running      bool           // is this pipeline processing?
}

// newPipeline creates an empty pipeline.
func newPipeline() *pipeline {
	pipe := &pipeline{}
	pipe.errors = make(chan error, 3)
	pipe.input = make(chan *Node, 10)
	pipe.results = pipe.input // short-curcuit, will be filled with filters
	return pipe
}

// Is this pipeline empty, i.e., has no filter stages yet?
func (pipe *pipeline) empty() bool {
	pipe.RLock()
	defer pipe.RUnlock()
	return pipe.filters == nil
}

// pushResult puts a node on the results channel of a filter stage (non-blocking).
// It is used by filter workers to communicate a result to the next stage
// of a pipeline.
func (f *filter) pushResult(node *Node) {
	f.env.queuecounter.Add(1)
	written := true
	select { // try to send it synchronously without blocking
	case f.results <- node:
	default:
		written = false
	}
	if !written { // nope, we'll have to go async
		go func(node *Node) {
			f.results <- node
		}(node)
	}
}

// pushBuffer puts a node on the buffer queue of a filter
// (non-blocking).
func (f *filter) pushBuffer(node *Node, udata interface{}) {
	nodesup := nodeSupport{node, udata}
	f.env.queuecounter.Add(1) // overall workload increases
	written := true
	select { // try to send it synchronously without blocking
	case f.queue <- nodesup:
	default:
		written = false
	}
	if !written { // nope, we'll have to go async
		go func(sup nodeSupport) {
			f.queue <- sup
		}(nodesup)
	}
}

// appendFilter appends a filter to a pipeline, i.e. as the last stage of
// the pipeline. Connects input- and output-channels appropriately and
// sets an environment for the filter.
func (pipe *pipeline) appendFilter(f *filter) {
	pipe.Lock()
	defer pipe.Unlock()
	if pipe.filters == nil {
		pipe.filters = f
	} else { // append at end of filter chain
		ff := pipe.filters
		for f.next != nil {
			ff = ff.next
		}
		ff.next = f
	}
	env := &filterenv{} // now set the environment for the filter
	env.errors = pipe.errors
	env.queuecounter = &pipe.queuecount
	env.input = pipe.results    // current output is input to new filter stage
	pipe.results = f.start(env) // remember new final output
}

// startProcessing starts a pipeline. It will start a watchdog goroutine
// to wait for the overall number of work packages to become zero.
// The watchdog will close all channels as soon as no more work
// packages (i.e., Nodes) are in the pipeline.
// Pre-requisite: at least one node/task in the front input channel.
func (pipe *pipeline) startProcessing() {
	pipe.Lock()
	defer pipe.Unlock()
	if !pipe.running {
		pipe.running = true
		go func() { // cleanup function
			pipe.queuecount.Wait() // wait for empty queues
			close(pipe.errors)
			close(pipe.input)
			f := pipe.filters
			i := 1
			for f != nil {
				close(f.results)
				if f.queue != nil {
					close(f.queue)
				}
				f = f.next
				i++
			}
			pipe.running = false
		}()
	}
}

// pushSync synchronously puts a node on the input channel of a pipeline.
func (pipe *pipeline) pushSync(node *Node) {
	pipe.queuecount.Add(1)
	pipe.input <- node // input q is buffered
}

// pushAsync asynchronously puts a node on the input channel of a pipeline.
func (pipe *pipeline) pushAsync(node *Node) {
	pipe.queuecount.Add(1)
	go func(node *Node) {
		pipe.input <- node
	}(node)
}

// waitForCompletion blocks until all work packages of a pipeline are done.
// It will receive the results of the final filter stage of the pipeline
// and collect them into a slice of Nodes. The slice will be a set, i.e.
// not contain duplicate Nodes.
func waitForCompletion(results <-chan *Node, errch <-chan error, counter *sync.WaitGroup) ([]*Node, error) {
	// Collect all results from the pipeline
	var selection []*Node
	m := make(map[*Node]struct{}) // intermediate map to suppress duplicates
	for node := range results {   // drain results channel
		m[node] = struct{}{}
		counter.Done() // we removed a value => count down
	}
	for node, _ := range m {
		selection = append(selection, node) // collect unique return values
	}
	// Get last error from error channel
	var lasterror error
	for err := range errch {
		if err != nil {
			lasterror = err // throw away all errors but the last one
		}
	}
	return selection, lasterror
}
