//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package hot

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// JetWaiter provides method for locking on jet id.
type JetWaiter interface {
	Wait(ctx context.Context, jetID insolar.ID) error
}

// JetReleaser provides methods for releasing jet waiters.
type JetReleaser interface {
	Unlock(ctx context.Context, jetID insolar.ID) error
	ThrowTimeout(ctx context.Context)
}

// ChannelWaiter implements methods for locking and unlocking a certain jet id.
type ChannelWaiter struct {
	lock    sync.Mutex
	waiters map[insolar.ID]waiter
	timeout chan struct{}
}

type waiter chan struct{}

func (w waiter) isClosed() bool {
	select {
	case <-w:
		return true
	default:
	}
	return false
}

// NewChannelWaiter creates new waiter instance.
func NewChannelWaiter() *ChannelWaiter {
	return &ChannelWaiter{
		waiters: map[insolar.ID]waiter{},
		timeout: make(chan struct{}),
	}
}

func (w *ChannelWaiter) waiterForJet(jetID insolar.ID) waiter {
	if _, ok := w.waiters[jetID]; !ok {
		w.waiters[jetID] = make(waiter)
	}
	return w.waiters[jetID]
}

// Wait waits for the raising one of two channels.
// If hotDataChannel or timeoutChannel was raised, the method returns error
// Either nil or ErrHotDataTimeout
func (w *ChannelWaiter) Wait(ctx context.Context, jetID insolar.ID) error {
	w.lock.Lock()
	waiter := w.waiterForJet(jetID)
	timeout := w.timeout
	w.lock.Unlock()

	select {
	case <-waiter:
		return nil
	case <-timeout:
		return insolar.ErrHotDataTimeout
	}
}

// Unlock raises hotDataChannel
func (w *ChannelWaiter) Unlock(ctx context.Context, jetID insolar.ID) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	waiter := w.waiterForJet(jetID)
	if waiter.isClosed() {
		return ErrWaiterNotLocked
	}
	close(waiter)
	return nil
}

// ThrowTimeout raises all timeoutChannel
func (w *ChannelWaiter) ThrowTimeout(ctx context.Context) {
	w.lock.Lock()
	defer w.lock.Unlock()

	inslogger.FromContext(ctx).Debug("raising timeout for requests")
	close(w.timeout)
	w.timeout = make(chan struct{})
	w.waiters = map[insolar.ID]waiter{}
}
