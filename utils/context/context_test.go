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

package context

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	var counter uint64

	statsCtx, cancelStatsCtx := WithCancel()
	var wg sync.WaitGroup

	goroutines := 100
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			<-statsCtx.Done()
			atomic.AddUint64(&counter, 1)
			wg.Done()
		}()
	}

	assert.Equal(t, uint64(0), counter)
	cancelStatsCtx()
	wg.Wait()
	assert.Equal(t, uint64(goroutines), counter)
}
