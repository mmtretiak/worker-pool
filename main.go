package main

import (
	"math/rand"
	"time"
)

const (
	WorkersCount = 5
	InputCount   = 10
	MinNumber    = 1
	MaxNumber    = 10000
)

func main() {
	inputChan := generateInput()

	dispatcher := NewDispatcher(inputChan, WorkersCount)

	doneChan := dispatcher.Start()
	<-doneChan
}

func generateInput() JobQueue {
	dataChan := make(chan Job)

	go func() {
		for i := 0; i < InputCount; i++ {
			dataChan <- NewJob(generateNumber())
		}

		close(dataChan)
	}()

	return dataChan
}

func generateNumber() int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(MaxNumber-MinNumber) + MinNumber
}
