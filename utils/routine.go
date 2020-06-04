// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
