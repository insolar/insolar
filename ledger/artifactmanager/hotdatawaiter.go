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

package artifactmanager

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// HotDataWaiter provides waiting system for a specific jet
// We tend to think, that it will be used for waiting hot-data in handler
// Also, because of the some jet pitfalls, we need to have an instrument
// to handler edge-cases from pulse manager.
// The main case is when a light material executes a jet for more then 1 pulse
// If it happens, we need to stop waiters from raising and waiting
//go:generate minimock -i github.com/insolar/insolar/ledger/storage.HotDataWaiter -o ./ -s _mock.go
type HotDataWaiter interface {
	Wait(ctx context.Context, jetID insolar.RecordID) error
	Unlock(ctx context.Context, jetID insolar.RecordID) error
	ThrowTimeout(ctx context.Context)
}

// HotDataWaiterConcrete is an implementation of HotDataWaiter
type HotDataWaiterConcrete struct {
	lock    sync.Mutex
	waiters map[insolar.RecordID]waiter
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

// NewHotDataWaiterConcrete is a constructor
func NewHotDataWaiterConcrete() *HotDataWaiterConcrete {
	return &HotDataWaiterConcrete{
		waiters: map[insolar.RecordID]waiter{},
		timeout: make(chan struct{}),
	}
}

func (w *HotDataWaiterConcrete) waiterForJet(jetID insolar.RecordID) waiter {
	if _, ok := w.waiters[jetID]; !ok {
		w.waiters[jetID] = make(waiter)
	}
	return w.waiters[jetID]
}

// Wait waits for the raising one of two channels.
// If hotDataChannel or timeoutChannel was raised, the method returns error
// Either nil or ErrHotDataTimeout
func (w *HotDataWaiterConcrete) Wait(ctx context.Context, jetID insolar.RecordID) error {
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
func (w *HotDataWaiterConcrete) Unlock(ctx context.Context, jetID insolar.RecordID) error {
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
func (w *HotDataWaiterConcrete) ThrowTimeout(ctx context.Context) {
	w.lock.Lock()
	defer w.lock.Unlock()

	inslogger.FromContext(ctx).Debug("raising timeout for requests")
	close(w.timeout)
	w.timeout = make(chan struct{})
	w.waiters = map[insolar.RecordID]waiter{}
}
