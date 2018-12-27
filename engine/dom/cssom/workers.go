package cssom

import (
	"fmt"
	"sync"
)

// We will employ a small helper for managing worker goroutines.
type workerLauncher int

// Prepare for launching n worker goroutines.
// 1 <= n < 100, otherwise will panic.
//
// Call the launcher like this:
//
//    workers := launch(3).workers(task, workload, errors)
//    watch(workers, doSomething)  // async, to wait for workers, called only once
//    order(workers, doSomething)  // create more workload, called many times
//    errlist := collect(errors)   // read errors util closed => no more workers active
//
// This pattern follows the one in "The Go Programming Language", chapter 8.5.
//
func launch(n int) workerLauncher {
	if n <= 0 || n >= 100 {
		panic("internal: illegal number of concurrent workers. 0 < n < 100")
	}
	return workerLauncher(n)
}

// Launches a previously defined number of workers.
// Clients have to supply 2 channels:
// - workload for read and write (workers may produce work packages themselves)
// - errorch for write (to be read by the caller)
// Returns a workergroup (see below).
func (launcher workerLauncher) workers(task workerTask, workload chan workPackage,
	errorch chan<- error) *workergroup {
	workers := &workergroup{}
	workers.queuecount = sync.WaitGroup{}
	workers.workload = workload
	workers.errorch = errorch
	for i := 0; i < int(launcher); i++ {
		wno := i + 1
		go func(workload <-chan workPackage) {
			defer func() {
				fmt.Printf("finished worker #%d\n", wno)
			}()
			for wp := range workload { // get workpackages until drained
				err := task(wp) // perform task on workpackage
				if err != nil {
					errorch <- err // signal error to caller
				}
				workers.queuecount.Done() // worker has finished a workpackage
			}
		}(workload)
	}
	return workers
}

// workergroup is a helper to asynchronously wait for a group of workers
// to complete the workload.
type workergroup struct {
	queuecount sync.WaitGroup   // count of work packages
	workload   chan workPackage // channel for work packages
	errorch    chan<- error     // response channel for error messages
}

// waitfor spawns a goroutine to wait for completion of a worker group.
// It will close input- and output-channel supplies for the call to launch.
// Closing errorch should signal to the caller that no more worker is running.
func watch(workers *workergroup, wp workPackage) {
	go func() {
		order(workers, wp)        // must be before q.Wait()
		workers.queuecount.Wait() // wait for empty queue
		fmt.Printf("all workers are done\n")
		close(workers.errorch)
	}()
}

// Order a new workpackage, i.e. put it on the workload queue.
func order(workers *workergroup, wp workPackage) {
	workers.queuecount.Add(1) // must be before put
	go func(wp workPackage) {
		workers.workload <- wp
	}(wp)
}

// Collect error message from workers, waiting synchronously until
// error channel is drained.
func collect(errorch <-chan error) []error {
	var e []error
	for err := range errorch {
		if err != nil {
			e = append(e, err)
		}
	}
	return e
}

// Workers will be tasked a series of workerTasks.
type workerTask func(wp workPackage) error

// TODO define a styled node work package.
type workPackage int
