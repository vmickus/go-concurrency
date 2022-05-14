package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

const WORKERS = 20

type Worker interface {
	doWork(int, chan<- string, chan<- error, *sync.WaitGroup)
}

func main() {
	up := Upstream{}
	resultChan := make(chan string)
	errChan := make(chan error)

	printGoroutines()

	// Initialize goroutines
	go runWorkers(up, resultChan, errChan)

	// Initialize timer
	timer := time.NewTimer(10 * time.Second)
	log.Println("Started timer")
	defer timer.Stop()

	// Collect response from goroutines or timeout
	loop := true
	results := make([]string, 0)
	for loop {
		select {
		case result, ok := <-resultChan:
			if !ok {
				loop = false
				break
			}
			printGoroutines()
			results = append(results, result)
		case err, ok := <-errChan:
			if !ok {
				loop = false
				break
			}
			printGoroutines()
			log.Printf("error %v", err)
		case <-timer.C:
			// Instead of return, we just exit the loop so we can process the results that we managed to obtain
			log.Println("timeout")
			loop = false
		}
	}

	// Display work done
	fmt.Printf("Processed %d jobs from %d\nResult: %v\n", len(results), WORKERS, results)

	//TODO Is it possible to close the goroutines after timeout?
	printGoroutines()
}

func runWorkers(upstream Worker, resultChan chan<- string, errChan chan<- error) {
	var wg sync.WaitGroup

	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)
		go func(i int) {
			upstream.doWork(i, resultChan, errChan, &wg)
		}(i)
	}

	wg.Wait()
	// Only close channel after all goroutines finish
	close(resultChan)
	close(errChan)
}

func printGoroutines() {
	log.Printf("Number of goroutines %d\n", runtime.NumGoroutine())
}
