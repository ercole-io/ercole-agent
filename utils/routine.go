package utils

import (
	"sync"

	"github.com/ercole-io/ercole-agent/config"
)

// RunRoutine will run function in a separeted goroutine and wait or not due to config.
func RunRoutine(configuration config.Configuration, function func()) {
	if configuration.ParallelizeRequests {
		go function()
	} else {
		function()
	}
}

// RunRoutineInGroup will run function in a separeted goroutine and wait or not due to config.
// Increment waitGroup counter and notify when done.
func RunRoutineInGroup(configuration config.Configuration, function func(), waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)

	if configuration.ParallelizeRequests {
		go func() {
			function()
			waitGroup.Done()
		}()
	} else {
		function()
		waitGroup.Done()
	}
}
