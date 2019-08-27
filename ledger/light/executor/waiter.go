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

package executor

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// JetWaiter provides method for locking on jet id.
type JetWaiter interface {
	Wait(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) error
}

// HotDataStatusChecker provides methods for checking receiving status of hot data.
type HotDataStatusChecker interface {
	IsReceived(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) bool
}

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.JetReleaser -o ./ -s _mock.go -g
// JetReleaser provides methods for releasing jet waiters.
type JetReleaser interface {
	Unlock(ctx context.Context, pulse insolar.PulseNumber, jetID insolar.JetID) error
	CloseAllUntil(ctx context.Context, pulse insolar.PulseNumber)
}

// ChannelWaiter implements methods for locking and unlocking a certain jet id.
type ChannelWaiter struct {
	lock        sync.Mutex
	closedUntil insolar.PulseNumber
	waiters     map[insolar.PulseNumber]*pulseWaiter
}

type pulseWaiter struct {
	pulse   insolar.PulseNumber
	waiters map[insolar.JetID]waiter
	timeout chan struct{}
}

func (pw *pulseWaiter) getOrCreate(jetID insolar.JetID) waiter {
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
func (w *ChannelWaiter) Wait(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) error {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"pulse": pulse.String(),
		"jet":   jetID.DebugString(),
	})
	logger.Debug("started waiting for hot objects")

	w.lock.Lock()
	if pulse <= w.closedUntil {
		w.lock.Unlock()
		return nil
	}
	pWaiter := w.getOrCreate(pulse)
	timeout := pWaiter.timeout
	waitCh := pWaiter.getOrCreate(jetID)
	w.lock.Unlock()

	select {
	case <-waitCh:
		return nil
	case <-timeout:
		logger.Errorf("timeout while waiting for hot objects")
		return insolar.ErrHotDataTimeout
	}
}

func (w *ChannelWaiter) IsReceived(ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber) bool {
	w.lock.Lock()
	defer w.lock.Unlock()

	pWaiter, ok := w.waiters[pn]
	if !ok {
		return false
	}
	jWaiter, ok := pWaiter.waiters[jetID]
	if !ok {
		return false
	}
	return jWaiter.isClosed()
}

// Unlock raises hotDataChannel
func (w *ChannelWaiter) Unlock(ctx context.Context, pulse insolar.PulseNumber, jetID insolar.JetID) error {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"pulse": pulse.String(),
		"jet":   jetID.DebugString(),
	})

	w.lock.Lock()
	defer w.lock.Unlock()

	waitCh := w.getOrCreate(pulse).getOrCreate(jetID)
	if waitCh.isClosed() {
		return ErrWaiterNotLocked
	}
	close(waitCh)
	logger.Debug("unlocked hot objects")
	return nil
}

// CloseAllUntil raises timeouts on all waiters until pulse.
func (w *ChannelWaiter) CloseAllUntil(ctx context.Context, pulse insolar.PulseNumber) {
	w.lock.Lock()
	defer w.lock.Unlock()

	for pn, pWaiter := range w.waiters {
		if pn > pulse {
			continue
		}

		close(pWaiter.timeout)
		delete(w.waiters, pn)
	}

	w.closedUntil = pulse
}

func (w *ChannelWaiter) getOrCreate(pn insolar.PulseNumber) *pulseWaiter {
	pWaiter, ok := w.waiters[pn]
	if ok {
		return pWaiter
	}

	pWaiter = &pulseWaiter{
		pulse:   pn,
		waiters: map[insolar.JetID]waiter{},
		timeout: make(chan struct{}),
	}
	w.waiters[pn] = pWaiter
	return pWaiter
}
