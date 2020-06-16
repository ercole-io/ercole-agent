//+build rhel5

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
	"errors"
	"sync"
)

// ContextWithCancel return stdlib context
func WithCancel() (Context, func()) {
	c := new(cancelCtx)

	return c, func() { c.cancel(Canceled) }
}

type cancelCtx struct {
	Context

	mu   sync.Mutex    // protects following fields
	done chan struct{} // created lazily, closed by first cancel call
	err  error         // set to non-nil by the first cancel call
}

var Canceled = errors.New("context canceled")

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

func (c *cancelCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()
	return d
}

func (c *cancelCtx) cancel(err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}

	c.mu.Lock()

	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}

	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)
	}

	c.mu.Unlock()
}
