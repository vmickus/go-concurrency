package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

const WORKERS = 5

type Worker interface {
	doWork(int, chan<- string, *sync.WaitGroup)
}

func main() {
	up := Upstream{}
	workChan := make(chan string)

	printGoroutines()

	// Initialize goroutines
	go doWork(up, workChan)

	// Initialize timer
	timer := time.NewTimer(1 * time.Second)
	log.Println("Started timer")
	defer timer.Stop()

	// Collect response from goroutines or timeout
	loop := true
	results := make([]string, 0)
	for loop {
		select {
		case result, ok := <-workChan:
			if !ok {
				loop = false
				break
			}
			printGoroutines()
			results = append(results, result)
		case <-timer.C:
			log.Println("timeout")
			loop = false
		}
	}

	// Display work done
	fmt.Printf("Result collected from upstream: %v\n", results)
	printGoroutines()
}

func doWork(upstream Worker, ch chan<- string) {
	var wg sync.WaitGroup

	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)
		go func(i int) {
			upstream.doWork(i, ch, &wg)
		}(i)
	}

	// Only close channel after all goroutines finish
	wg.Wait()
	close(ch)
}

func printGoroutines() {
	log.Printf("Number of goroutines %d\n", runtime.NumGoroutine())
}
