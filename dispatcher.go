package main

import "sync"

type doneChan chan struct{}

type Dispatcher interface {
	Start() doneChan
}

func NewDispatcher(inputChan JobQueue, workersCount int) Dispatcher {
	return &dispatcher{
		workersCount: workersCount,
		workersWG:    &sync.WaitGroup{},
		workersChan:  []JobQueue{},

		inputChan: inputChan,
		doneChan:  make(chan struct{}),
	}
}

type dispatcher struct {
	workersCount int
	workersChan  []JobQueue
	workersWG    *sync.WaitGroup

	inputChan JobQueue
	doneChan  doneChan
}

func (d *dispatcher) Start() doneChan {
	d.workersWG.Add(d.workersCount)

	for i := 0; i < d.workersCount; i++ {
		worker := NewWorker(d.workersWG)

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
