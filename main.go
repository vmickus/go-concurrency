package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

const (
	WORKERS = 10
	TIMEOUT = 5
)

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
	timer := time.NewTimer(TIMEOUT * time.Second)
	log.Println("Started timer")
	defer timer.Stop()

	loop := true
	results := make([]string, 0)
	// Collect response from goroutines or timeout
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
			// Log error, but let it continue to collect other results.
			// In this case we are fine with some upstream calls that errored
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
	printResult(results)

	// TODO Is it possible to close the goroutines after timeout?
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

func printResult(results []string) {
	fmt.Printf("Processed %d of %d:\n", len(results), WORKERS)
	for _, result := range results {
		fmt.Printf("\t%s\n", result)
	}
}
