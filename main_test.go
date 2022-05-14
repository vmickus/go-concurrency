package main

import (
	"testing"
)

func TestRunWorker(t *testing.T) {

	up := Upstream{}
	chResult := make(chan string)
	chErr := make(chan error)

	go runWorkers(up, chResult, chErr)

	counter := 0
	for i := 0; i < WORKERS; i++ {
		select {
		case <-chResult:
			counter++
		case <-chErr:
			counter++
		}
	}

	// This test make sense? It seems to me that when it reaches here it will always have counter == WORKERS
	if counter != WORKERS {
		t.Errorf("test failed")
	}
}
