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

3. Neither the name of this software nor the names of its contributors
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
	"sort"
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
//
// node: input tree node
// isbuffered: is the input node from this stage's buffer queue?
// udata: user provided additional data
// emit:  // function to emit result node to next stage
// buffer: function to queue node in local buffer
//
// Does not return anything except a possible error condition.
type workerTask func(node *Node, isbuffered bool, udata userdata,
	emit func(*Node, uint32), buffer func(*Node, interface{}, uint32)) error

// filter is part of a pipeline, i.e. a stage of the overall pipeline to
// process input (Nodes) and produce results (Nodes).
// filters will perform concurrently.
type filter struct {
	next       *filter            // filters are chained
	results    chan<- nodePackage // results of this filter (pipeline stage)
	queue      chan nodePackage   // helper queue if necessary
	task       workerTask         // the task this filter performs
	filterdata interface{}        // user-provided information needed to perform task
	env        *filterenv         // connection to outside world
}

// nodePackage is the type which is transported in a pipeline.
// Each pipeline stage emits an instance of this type to the next stage.
//
// 'nodelocal' lets clients store arbitrary user data together with the node.
// It will be set to 'nil' as soon as the nodepackage is transferred to the next stage,
// i.e., this type is local to a pipeline-stage/filter.
type nodePackage struct {
	node      *Node       // tree node
	nodelocal interface{} // arbitrary user data
	serial    uint32      // serial number of node for ordering
}

// filterenv holds information about the outside world to be referenced by
// a filter. This includes input workload, error destination and a counter
// for overall work on an pipeline.
type filterenv struct {
	input        <-chan nodePackage // work to do for this filter, connected to predecessor
	errors       chan<- error       // where errors are reported to
	queuecounter *sync.WaitGroup    // counter for overall work load
}

// userdata is a container managed by the pipeline mechanism. It will contain
// two types of data availble for filters:
// information global to a filter (filterdata),
// and information acompanying a single node (nodelocal & serial).
// The pipeline mechanism will construct this from the filter environment and from
// node-local user-managed data, and it will deconstruct it for calls to a 'task()'.
type userdata struct {
	filterdata interface{}
	nodelocal  interface{}
	serial     uint32
}

// newFilter creates a new pipeline stage, i.e. a filter fed from an input
// channel (workload). the filter is expected to put processed nodes into an
// output channel (results).
//
// Errors are reported to an error channel.
func newFilter(task workerTask, filterdata interface{}, buflen int) *filter {
	f := &filter{}
	if buflen > 0 {
		if buflen > maxBufferLength {
			buflen = maxBufferLength
		}
		f.queue = make(chan nodePackage, buflen)
	}
	f.task = task
	f.filterdata = filterdata
	return f
}

// This method signature is a bit strange, but for now it does the job.
// Sets an environment for a filter an gets the results-channel in return.
func (f *filter) start(env *filterenv) chan nodePackage {
	f.env = env
	res := make(chan nodePackage, 3) // output channel has to be in place before workers start
	f.results = res                  // be careful to set write-only for the filter
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
//
// Each worker is identified through a worker number 'wno'.
func filterWorker(f *filter, wno int) {
	//  defer func() {
	//	log.Printf("finished worker #%d\n", wno) // for debugging
	//}()
	push := func(node *Node, serial uint32) { // worker will use this to hand result to next stage
		f.pushResult(node, serial)
	}
	for inNode := range f.env.input { // get workpackages until drained
		node := inNode.node
		serial := inNode.serial
		udata := userdata{f.filterdata, nil, serial}
		err := f.task(node, false, udata, push, nil) // perform task on workpackage
		if err != nil {
			f.env.errors <- err // signal error to caller
		}
		f.env.queuecounter.Done() // worker has finished a workpackage
	}
}

// filterWorkerWithQueue is a worker function which uses a separate support
// queue, the 'buffer queue'. This buffer queue may be used to re-schedule nodes
// until they are completely processed.
func filterWorkerWithQueue(f *filter, wno int) {
	push := func(node *Node, serial uint32) { // worker will use this to hand result to next stage
		f.pushResult(node, serial)
	}
	pushBuf := func(sup *Node, udata interface{}, serial uint32) { // worker will use this to queue work internally
		f.pushBuffer(sup, udata, serial)
	}
	var buffered bool
	var node *Node
	var udata userdata
	for {
		select { // get upstream workpackages and buffered workpackages until drained
		case inNode := <-f.env.input:
			node = inNode.node
			udata.serial = inNode.serial
			udata.nodelocal = nil
			udata.filterdata = f.filterdata
			buffered = false
		case supdata := <-f.queue:
			node = supdata.node
			udata.filterdata = f.filterdata
			udata.nodelocal = supdata.nodelocal
			udata.serial = supdata.serial
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
	sync.RWMutex                  // to sychronize access to various fields
	queuecount   sync.WaitGroup   // overall count of work packages
	errors       chan error       // collector channel for error messages
	filters      *filter          // chain of filters
	input        chan nodePackage // initial workload
	results      chan nodePackage // where final output of this pipeline goes to
	running      bool             // is this pipeline processing?
}

// newPipeline creates an empty pipeline.
func newPipeline() *pipeline {
	pipe := &pipeline{}
	pipe.errors = make(chan error, 3)
	pipe.input = make(chan nodePackage, 10)
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
func (f *filter) pushResult(node *Node, serial uint32) {
	f.env.queuecounter.Add(1)
	written := true
	select { // try to send it synchronously without blocking
	case f.results <- nodePackage{node, nil, serial}:
	default:
		written = false
	}
	if !written { // nope, we'll have to go async
		go func(node *Node, serial uint32) {
			f.results <- nodePackage{node, nil, serial}
		}(node, serial)
	}
}

// pushBuffer puts a node on the buffer queue of a filter
// (non-blocking).
func (f *filter) pushBuffer(node *Node, udata interface{}, serial uint32) {
	nodesup := nodePackage{node, udata, serial}
	f.env.queuecounter.Add(1) // overall workload increases
	written := true
	select { // try to send it synchronously without blocking
	case f.queue <- nodesup:
	default:
		written = false
	}
	if !written { // nope, we'll have to go async
		go func(sup nodePackage) {
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
func (pipe *pipeline) pushSync(node *Node, serial uint32) {
	pipe.queuecount.Add(1)
	pipe.input <- nodePackage{node, nil, serial} // input q is buffered
}

// pushAsync asynchronously puts a node on the input channel of a pipeline.
func (pipe *pipeline) pushAsync(node *Node, serial uint32) {
	pipe.queuecount.Add(1)
	go func(node *Node) {
		pipe.input <- nodePackage{node, nil, serial} // input q is buffered
	}(node)
}

// waitForCompletion blocks until all work packages of a pipeline are done.
// It will receive the results of the final filter stage of the pipeline
// and collect them into a slice of Nodes. The slice will be a set, i.e.
// not contain duplicate Nodes.
func waitForCompletion(results <-chan nodePackage, errch <-chan error, counter *sync.WaitGroup) ([]*Node, error) {
	// Collect all results from the pipeline
	var selection []*Node          // slice of nodes -> return value
	var serials []uint32           // slice of serial numbers for ordering
	m := make(map[*Node]uint32)    // intermediate map to suppress duplicates
	for nodepkg := range results { // drain results channel
		m[nodepkg.node] = nodepkg.serial // remember last serial for node (may be random)
		counter.Done()                   // we removed a value => count down
	}
	for node, serial := range m { // extract unique results into slices
		selection = append(selection, node) // collect unique return values
		serials = append(serials, serial)
		// resultSlices is a helper struct for sorting
		// it implements the Sort interface
		if len(selection) > 0 && selection[0].Rank > 0 { // if rank is unset: no sorting possible
			sort.Sort(resultSlices{selection, serials})
		}
		// after this, serials are discarded
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
