/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package artifactmanager

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
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
	Wait(ctx context.Context, jetID core.RecordID) error
	Unlock(ctx context.Context, jetID core.RecordID)
	ThrowTimeout(ctx context.Context)
}

// HotDataWaiterConcrete is an implementation of HotDataWaiter
type HotDataWaiterConcrete struct {
	waitersMapLock sync.Mutex
	waiters        map[core.RecordID]*waiter
}

// NewHotDataWaiterConcrete is a constructor
func NewHotDataWaiterConcrete() *HotDataWaiterConcrete {
	return &HotDataWaiterConcrete{waiters: map[core.RecordID]*waiter{}}
}

type waiter struct {
	hotDataChannel chan struct{}
	timeoutChannel chan struct{}
}

func (hdw *HotDataWaiterConcrete) getWaiter(ctx context.Context, jetID core.RecordID) *waiter {
	hdw.waitersMapLock.Lock()
	defer hdw.waitersMapLock.Unlock()

	if _, ok := hdw.waiters[jetID]; !ok {
		hdw.waiters[jetID] = &waiter{
			hotDataChannel: make(chan struct{}),
			timeoutChannel: make(chan struct{}),
		}
	}

	return hdw.waiters[jetID]
}

// Wait waits for the raising one of two channels.
// If hotDataChannel or timeoutChannel was raised, the method returns error
// Either nil or ErrHotDataTimeout
func (hdw *HotDataWaiterConcrete) Wait(ctx context.Context, jetID core.RecordID) error {
	waiter := hdw.getWaiter(ctx, jetID)

	select {
	case <-waiter.hotDataChannel:
		return nil
	case <-waiter.timeoutChannel:
		return core.ErrHotDataTimeout
	}
}

// Unlock raises hotDataChannel
func (hdw *HotDataWaiterConcrete) Unlock(ctx context.Context, jetID core.RecordID) {
	waiter := hdw.getWaiter(ctx, jetID)

	hdw.waitersMapLock.Lock()
	defer hdw.waitersMapLock.Unlock()

	close(waiter.hotDataChannel)
}

// ThrowTimeout raises all timeoutChannel
func (hdw *HotDataWaiterConcrete) ThrowTimeout(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	hdw.waitersMapLock.Lock()
	defer hdw.waitersMapLock.Unlock()

	for jetID, waiter := range hdw.waiters {
		logger.WithField("jetid", jetID.DebugString()).Debug("raising timeout for requests")
		close(waiter.timeoutChannel)
	}

	hdw.waiters = map[core.RecordID]*waiter{}
}
