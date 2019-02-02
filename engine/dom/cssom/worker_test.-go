package cssom

import "time"
import "fmt"
import "errors"
import "testing"

func makeWorkPackage(n int) workPackage {
	wp := workPackage{udata: n}
	return wp
}

// type workerTask
func mytask(wp workPackage) error {
	duration := time.Duration(wp.udata.(int))
	time.Sleep(duration)
	return errors.New(fmt.Sprintf("slept:%d", duration))
}

func Test3Workers(t *testing.T) {
	workload := make(chan workPackage)
	errors := make(chan error)
	workers := launch(3).workers(mytask, workload, errors)
	watch(workers, makeWorkPackage(1000)) // initial work package
	for i := 0; i < 5; i++ {              // create another 5 work packages
		order(workers, makeWorkPackage(800+i*100))
	}
	e := collect(errors) // wait for workers to complete
	t.Logf("errors = %v", e)
	if len(e) != 6 {
		t.Errorf("expected to receive 6 messages from 3 workers")
	}
}
