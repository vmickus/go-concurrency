package main

import (
	"fmt"
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
		case result := <-chResult:
			if result == "" {
				t.Failed()
			}
			counter++
		case <-chErr:
			counter++
		}
	}

	fmt.Println(counter)
	if counter != (WORKERS - 1) {
		t.Failed()
	}
	t.Failed()

}
