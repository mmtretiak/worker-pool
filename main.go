package main

import (
	"math/rand"
	"time"

	"worker-pool/worker_pool/dispatcher"
	"worker-pool/worker_pool/job"
)

const (
	WorkersCount = 5
	InputCount   = 10
	MinNumber    = 1
	MaxNumber    = 10000
)

func main() {
	inputChan := generateInput()

	dispatcher := dispatcher.NewDispatcher(inputChan, WorkersCount)

	doneChan := dispatcher.Start()
	<-doneChan
}

func generateInput() job.Queue {
	dataChan := make(chan job.Job)

	go func() {
		for i := 0; i < InputCount; i++ {
			dataChan <- job.NewJob(generateNumber())
		}

		close(dataChan)
	}()

	return dataChan
}

func generateNumber() int {
	// Use the Seed function to initialize the default Source for different behavior for each run.
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(MaxNumber-MinNumber) + MinNumber
}
