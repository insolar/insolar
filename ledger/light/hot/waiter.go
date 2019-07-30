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
	Wait(ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/light/hot.JetReleaser -o ../../../testutils -s _mock.go -g

// JetReleaser provides methods for releasing jet waiters.
type JetReleaser interface {
	Unlock(ctx context.Context, jetID insolar.ID) error
	ThrowTimeout(ctx context.Context, pulse insolar.PulseNumber)
}

// ChannelWaiter implements methods for locking and unlocking a certain jet id.
type ChannelWaiter struct {
	lock    sync.Mutex
	waiters map[insolar.PulseNumber]*pulseWaiter
}

type pulseWaiter struct {
	pulse   insolar.PulseNumber
	waiters map[insolar.ID]waiter
	timeout chan struct{}
}

func (pw *pulseWaiter) getOrCreate(jetID insolar.ID) waiter {
	if _, ok := pw.waiters[jetID]; !ok {
		pw.waiters[jetID] = make(waiter)
	}
	return pw.waiters[jetID]
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
		waiters: map[insolar.PulseNumber]*pulseWaiter{},
	}
}

// Wait waits for the raising one of two channels.
// If hotDataChannel or timeoutChannel was raised, the method returns error
// Either nil or ErrHotDataTimeout
func (w *ChannelWaiter) Wait(ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber) error {
	w.lock.Lock()
	pWaiter := w.getOrCreate(pulse)
	timeout := pWaiter.timeout
	waitCh := pWaiter.getOrCreate(jetID)
	w.lock.Unlock()

	select {
	case <-waitCh:
		return nil
	case <-timeout:
		return insolar.ErrHotDataTimeout
	}
}

// Unlock raises hotDataChannel
func (w *ChannelWaiter) Unlock(ctx context.Context, pulse insolar.PulseNumber, jetID insolar.ID) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	waitCh := w.getOrCreate(pulse).getOrCreate(jetID)
	if waitCh.isClosed() {
		return ErrWaiterNotLocked
	}
	close(waitCh)
	return nil
}

// ThrowTimeout raises timeouts on all waiters for pulse.
func (w *ChannelWaiter) ThrowTimeout(ctx context.Context, pn insolar.PulseNumber) {
	w.lock.Lock()
	defer w.lock.Unlock()

	inslogger.FromContext(ctx).Debug("raising timeout for requests")
	w.close(pn)
}

func (w *ChannelWaiter) getOrCreate(pn insolar.PulseNumber) *pulseWaiter {
	pWaiter, ok := w.waiters[pn]
	if ok {
		return pWaiter
	}

	pWaiter = &pulseWaiter{
		pulse:   pn,
		waiters: map[insolar.ID]waiter{},
		timeout: make(chan struct{}),
	}
	w.waiters[pn] = pWaiter
	return pWaiter
}

func (w *ChannelWaiter) close(pn insolar.PulseNumber) {
	pWaiter, ok := w.waiters[pn]
	if !ok {
		return
	}

	close(pWaiter.timeout)
	delete(w.waiters, pn)
}
