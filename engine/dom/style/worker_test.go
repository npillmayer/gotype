package style

import "time"
import "fmt"
import "errors"
import "testing"

func mytask(wp workPackage) error {
	duration := time.Duration(wp)
	time.Sleep(duration)
	return errors.New(fmt.Sprintf("slept %d", duration))
}

func Test3Workers(t *testing.T) {
	workload := make(chan workPackage)
	errorch := make(chan error)
	workers := launch(3).workers(mytask, workload, errorch)
	watch(workers)
	go func() {
		workload <- workPackage(1000)
		workload <- workPackage(2000)
		workload <- workPackage(3000)
		close(workload)
	}()
	e := collect(errorch)
	if len(e) != 3 {
		t.Errorf("expected to receive 3 messages from 3 workers")
	}
}
