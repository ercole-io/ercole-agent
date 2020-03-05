package utils

import (
	"sync"

	"github.com/ercole-io/ercole-agent/config"
)

// RunRoutine will run function in a separeted goroutine and wait or not due to config.
func RunRoutine(configuration config.Configuration, function func()) {
	done := make(chan bool, 1)

	go func(done chan<- bool) {
		function()
		done <- true
	}(done)

	if !configuration.ParallelizeRequests {
		<-done
	}
}

// RunRoutineInGroup will run function in a separeted goroutine and wait or not due to config.
// Increment waitGroup counter and notify when done.
func RunRoutineInGroup(configuration config.Configuration, function func(), waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	done := make(chan bool, 1)

	go func(done chan<- bool) {
		function()
		done <- true
		waitGroup.Done()
	}(done)

	if !configuration.ParallelizeRequests {
		<-done
	}
}
