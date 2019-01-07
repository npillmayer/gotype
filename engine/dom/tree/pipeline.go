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
	"fmt"
	"log"
	"sync"
)

// Workers will be tasked a series of workerTasks.
type workerTask func(*Node, interface{}, func(*Node)) error

// filter is part of a pipeline, i.e. a stage of the overall pipeline to
// process input (Nodes) and produce results (Nodes).
// filters will perform concurrently.
type filter struct {
	next    *filter      // filters are chained
	results chan<- *Node // results of this filter (pipeline stage)
	task    workerTask   // the task this filter performs
	udata   interface{}  // additional information needed to perform task
	env     *filterenv   // connection to outside world
}

// filterenv holds information about the outside world to be referenced by
// a filter. This includes input workload, error destination and a counter
// for overall work on an pipeline.
type filterenv struct {
	input        <-chan *Node    // work to do for this filter
	errors       chan<- error    // where errors are reported to
	queuecounter *sync.WaitGroup // counter for overall work load
}

// newFilter creates a new pipeline stage, i.e. a filter fed from an input
// channel (workload) and putting processed nodes into an output channel (results).
// Errors are reported to an error channel.
func newFilter(task workerTask, udata interface{}) *filter {
	f := &filter{}
	f.task = task
	f.udata = udata
	return f
}

// This method signature is a bit strange, but for now it does the job.
// Sets an environment for a filter an gets the results-channel in return.
func (f *filter) start(env *filterenv) chan *Node {
	f.env = env
	res := make(chan *Node)  // output channel has to be in place before workers start
	f.results = res          // be careful to set write-only for the filter
	for i := 0; i < 1; i++ { // TODO how many workers?
		wno := i + 1
		go filterWorker(f, wno) // startup worker no. #wno
	}
	return res // needed r/w for next filter in pipe
}

// filterWorker is the default worker function. Each filter is free to start
// as many of them as seems adequate.
func filterWorker(f *filter, wno int) {
	defer func() {
		log.Printf("finished worker #%d\n", wno)
	}()
	counter := f.env.queuecounter
	for inNode := range f.env.input { // get workpackages until drained
		push := func(node *Node) {
			pushResult(node, f.results, counter)
		}
		err := f.task(inNode, f.udata, push) // perform task on workpackage
		if err != nil {
			f.env.errors <- err // signal error to caller
		}
		f.env.queuecounter.Done() // worker has finished a workpackage
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

// pre-requisite: at least one node/task in the front input channel.
func (pipe *pipeline) startProcessing() {
	pipe.Lock()
	defer pipe.Unlock()
	if !pipe.running {
		pipe.running = true
		go func() { // cleanup function
			pipe.queuecount.Wait() // wait for empty queues
			fmt.Printf("all work is done\n")
			close(pipe.errors)
			close(pipe.input)
			f := pipe.filters
			i := 1
			for f != nil {
				close(f.results)
				f = f.next
				i++
			}
			pipe.running = false
		}()
	}
}

func (pipe *pipeline) pushSync(node *Node) {
	pipe.queuecount.Add(1)
	pipe.input <- node // input q is buffered
}

func (pipe *pipeline) pushAsync(node *Node) {
	pipe.queuecount.Add(1)
	go func(node *Node) {
		pipe.input <- node
	}(node)
}

func pushResult(node *Node, output chan<- *Node, counter *sync.WaitGroup) {
	counter.Add(1)
	go func(node *Node) {
		output <- node
	}(node)
}

func waitForCompletion(results <-chan *Node, errch <-chan error, counter *sync.WaitGroup) ([]*Node, error) {
	// Collect all results from the pipeline
	var selection []*Node
	for node := range results {
		selection = append(selection, node)
		counter.Done() // we removed a value => count down
	}
	// Get last error
	var lasterror error
	for err := range errch {
		if err != nil {
			lasterror = err // throw away all errors but the last one
		}
	}
	return selection, lasterror
}
