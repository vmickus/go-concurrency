package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Simulate an upstream doing work

type Upstream struct {
}

func (up Upstream) doWork(job int, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	workFor := time.Duration(rand.Intn(job)) * time.Second

	log.Printf("job %d - working for %v\n", job, workFor)
	time.Sleep(workFor)
	log.Printf("job %d - finshed for %v\n", job, workFor)

	ch <- fmt.Sprintf("\nJob %d done", job)
}
