package main

import (
	"fmt"
	"sync"
)

type Worker interface {
	Start()
	GetInputChan() JobQueue
}

func NewWorker(wg *sync.WaitGroup) Worker {
	inputChan := make(chan Job)

	return &worker{
		inputChan: inputChan,
		wg:        wg,
	}
}

type worker struct {
	inputChan JobQueue
	wg        *sync.WaitGroup
}

func (w *worker) Start() {
	go func() {
		for {
			job, ok := <-w.inputChan
			// chan closed
			if !ok {
				fmt.Println("Exit")
				w.wg.Done()
				return
			}

			job.Do()
		}
	}()
}

func (w *worker) GetInputChan() JobQueue {
	return w.inputChan
}
