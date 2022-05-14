package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Simulate an upstream doing work

type Upstream struct {
}

// doWork simulate job done in a upstream service.
// The result will be sent to the channel.
func (up Upstream) doWork(job int, chResult chan<- string, chErr chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	// work duration
	workFor := time.Duration(rand.Intn(TIMEOUT)) * time.Second

	// log.Printf("job %d - working for %v\n", job, workFor)
	time.Sleep(workFor)
	// log.Printf("job %d - finshed for %v\n", job, workFor)

	if job%2 == 0 {
		chErr <- fmt.Errorf("job %d failed", job)
	} else {
		chResult <- fmt.Sprintf("Job %d done", job)
	}
}
