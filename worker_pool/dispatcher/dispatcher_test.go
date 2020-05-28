package dispatcher

import (
	"fmt"
	"sync/atomic"
	"testing"

	"worker-pool/worker_pool/job"
)

const (
	WorkersCount = 10
	InputCount   = 100
)

type testJob struct {
	doCallCount *uint32
}

func (t *testJob) Do() {
	atomic.AddUint32(t.doCallCount, 1)
}

func TestDispatcher(t *testing.T) {
	doCallCount := uint32(0)
	expectedCountOfDoneJobs := uint32(WorkersCount * InputCount)

	inputChan := generateTestInput(&doCallCount)

	dispatcher := NewDispatcher(inputChan, WorkersCount)

	doneChan := dispatcher.Start()
	<-doneChan

	if doCallCount != expectedCountOfDoneJobs {
		t.Error(fmt.Printf("Count of done jobs less than expected: done %v expected %v", doCallCount, expectedCountOfDoneJobs))
	}
}

func generateTestInput(doCallCount *uint32) job.Queue {
	dataChan := make(chan job.Job)

	go func() {
		for i := 0; i < InputCount; i++ {
			dataChan <- &testJob{doCallCount}
		}

		close(dataChan)
	}()

	return dataChan
}
