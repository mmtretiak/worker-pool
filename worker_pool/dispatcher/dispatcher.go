package dispatcher

import (
	"sync"

	"worker-pool/worker_pool/job"
	worker2 "worker-pool/worker_pool/worker"
)

type doneChan chan struct{}

type Dispatcher interface {
	Start() doneChan
}

func NewDispatcher(inputChan job.Queue, workersCount int) Dispatcher {
	return &dispatcher{
		workersCount: workersCount,
		workersWG:    &sync.WaitGroup{},
		workersChan:  []job.Queue{},

		inputChan: inputChan,
		doneChan:  make(chan struct{}),
	}
}

type dispatcher struct {
	workersCount int
	workersChan  []job.Queue
	workersWG    *sync.WaitGroup

	inputChan job.Queue
	doneChan  doneChan
}

func (d *dispatcher) Start() doneChan {
	d.workersWG.Add(d.workersCount)

	for i := 0; i < d.workersCount; i++ {
		worker := worker2.NewWorker(d.workersWG)

		workerChan := worker.GetInputChan()
		d.workersChan = append(d.workersChan, workerChan)

		worker.Start()
	}

	go d.dispatch()

	return d.doneChan
}

func (d *dispatcher) dispatch() {
	for {
		sleepTime, ok := <-d.inputChan
		// chan closed
		if !ok {
			for _, workerChan := range d.workersChan {
				close(workerChan)
			}
			break
		}

		for _, workerChan := range d.workersChan {
			workerChan <- sleepTime
		}
	}

	d.workersWG.Wait()
	d.doneChan <- struct{}{}
}
