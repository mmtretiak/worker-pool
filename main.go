package main

import (
	"math/rand"
	"sync"
	"time"
)

const (
	WorkersCount = 2
	InputCount   = 100
	MinNumber    = 1
	MaxNumber    = 10000
)

func main() {
	inputChan := generateInput()

	var workersChannels []chan int
	var workersWG sync.WaitGroup

	workersWG.Add(WorkersCount)

	for i := 0; i < WorkersCount; i++ {
		workerInput := make(chan int)
		workersChannels = append(workersChannels, workerInput)

		go func() {
			for {
				sleepTime, ok := <-workerInput
				// chan closed
				if !ok {
					workersWG.Done()
					return
				}

				time.Sleep(time.Millisecond * time.Duration(sleepTime))
			}
		}()
	}

	for {
		sleepTime, ok := <-inputChan
		// chan closed
		if !ok {
			for _, workerChan := range workersChannels {
				close(workerChan)
			}
			break
		}

		for _, workerChan := range workersChannels {
			workerChan <- sleepTime
		}
	}

	workersWG.Wait()
}

func generateInput() chan int {
	dataChan := make(chan int)

	go func() {
		for i := 0; i < InputCount; i++ {
			dataChan <- generateNumber()
		}

		close(dataChan)
	}()

	return dataChan
}

func generateNumber() int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(MaxNumber-MinNumber) + MinNumber
}
