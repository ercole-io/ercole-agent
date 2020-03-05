package utils

import "github.com/ercole-io/ercole-agent/config"

// RunInRoutine will run function in a separeted goroutine and wait or not due to config.
func RunInRoutine(configuration config.Configuration, function func()) {
	done := make(chan bool, 1)

	go func(done chan<- bool) {
		function()
		done <- true
	}(done)

	if !configuration.ParallelizeRequests {
		<-done
	}
}
