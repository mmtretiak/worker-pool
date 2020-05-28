package worker

import (
	"sync"

	"worker-pool/worker_pool/job"
)

type Worker interface {
	Start()
	GetInputChan() job.Queue
}

func NewWorker(wg *sync.WaitGroup) Worker {
	inputChan := make(chan job.Job)

	return &worker{
		inputChan: inputChan,
		wg:        wg,
	}
}

type worker struct {
	inputChan job.Queue
	wg        *sync.WaitGroup
}

func (w *worker) Start() {
	go func() {
		for {
			job, ok := <-w.inputChan
			// chan closed
			if !ok {
				w.wg.Done()
				return
			}

			job.Do()
		}
	}()
}

func (w *worker) GetInputChan() job.Queue {
	return w.inputChan
}
