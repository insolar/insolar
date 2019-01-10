/*
 *    Copyright 2018 Insolar
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
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type jetDropTimeoutProvider struct {
	waitersLock sync.RWMutex
	waiters     map[core.RecordID]*jetDropTimeout

	waitersInitLocksLock sync.Mutex
	waitersInitLocks     map[core.RecordID]*sync.RWMutex
}

func (p *jetDropTimeoutProvider) getLock(jetID core.RecordID) *sync.RWMutex {
	p.waitersInitLocksLock.Lock()
	defer p.waitersInitLocksLock.Unlock()

	if _, ok := p.waitersInitLocks[jetID]; !ok {
		p.waitersInitLocks[jetID] = &sync.RWMutex{}
	}

	return p.waitersInitLocks[jetID]
}

func (p *jetDropTimeoutProvider) getWaiter(jetID core.RecordID) *jetDropTimeout {
	p.waitersLock.RLock()
	defer p.waitersLock.RUnlock()

	return p.waiters[jetID]
}

type jetDropTimeout struct {
	lastJdPulseLock sync.RWMutex
	lastJdPulse     core.PulseNumber

	jetDropLocker chan struct{}
	timeoutLocker chan struct{}

	isTimeoutRunLock sync.Mutex
	isTimeoutRun     bool
}

func (jdw *jetDropTimeout) getLastPulse() core.PulseNumber {
	jdw.lastJdPulseLock.RLock()
	defer jdw.lastJdPulseLock.RUnlock()

	return jdw.lastJdPulse
}

func (jdw *jetDropTimeout) setLastPulse(pn core.PulseNumber) {
	jdw.lastJdPulseLock.Lock()
	defer jdw.lastJdPulseLock.Unlock()

	jdw.lastJdPulse = pn
}

func (m *middleware) waitForDrop(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		inslogger.FromContext(ctx).Debugf("[waitForDrop] pulse %v starts %v", parcel.Pulse(), time.Now())
		// If the call is a call in redirect-chain
		// skip waiting for the hot records
		if parcel.DelegationToken() != nil {
			inslogger.FromContext(ctx).Debugf("[waitForDrop] parcel.DelegationToken() != nil")
			return handler(ctx, parcel)
		}

		jetID := jetFromContext(ctx)
		lock := m.jetDropTimeoutProvider.getLock(jetID)
		waiter := m.jetDropTimeoutProvider.getWaiter(jetID)

		lock.RLock()
		if waiter == nil {
			inslogger.FromContext(ctx).Debugf("[waitForDrop] waiter is nil for %v", jetID)
			lock.RUnlock()
			return handler(ctx, parcel)
		}
		lock.RUnlock()

		if waiter.getLastPulse() < parcel.Pulse() {
			inslogger.FromContext(ctx).Debugf("[waitForDrop] waiter.getLastPulse() != parcel.Pulse(), %v - %v,", waiter.getLastPulse(), parcel.Pulse())
			waiter.runDropWaitingTimeout()

			select {
			case <-waiter.jetDropLocker:
			case <-waiter.timeoutLocker:
			}

			inslogger.FromContext(ctx).Debugf("[waitForDrop] after select - %v", time.Now())

			waiter.isTimeoutRunLock.Lock()
			waiter.isTimeoutRun = false
			waiter.isTimeoutRunLock.Unlock()
		}

		inslogger.FromContext(ctx).Debugf("[waitForDrop] before handler exec - %v", time.Now())
		fmt.Println("waiter.getLastJdPulse() - ", waiter.getLastPulse())
		fmt.Println("parcel.Pulse() - ", parcel.Pulse())
		fmt.Println("jetID - ", jetID)
		fmt.Println("waitForDrop, handle now")
		return handler(ctx, parcel)
	}
}

func (jdw *jetDropTimeout) runDropWaitingTimeout() {
	jdw.isTimeoutRunLock.Lock()
	defer jdw.isTimeoutRunLock.Unlock()

	if jdw.isTimeoutRun {
		return
	}

	jdw.isTimeoutRun = true
	jdw.timeoutLocker = make(chan struct{})

	go func() {
		time.Sleep(2 * time.Second)

		close(jdw.timeoutLocker)

		jdw.isTimeoutRunLock.Lock()
		jdw.isTimeoutRun = false
		jdw.isTimeoutRunLock.Unlock()
	}()
}

func (m *middleware) unlockDropWaiters(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		inslogger.FromContext(ctx).Debugf("[unlockDropWaiters] pulse %v starts %v", parcel.Pulse(), time.Now())
		jetID := jetFromContext(ctx)
		lock := m.jetDropTimeoutProvider.getLock(jetID)
		waiter := m.jetDropTimeoutProvider.getWaiter(jetID)
		inslogger.FromContext(ctx).Debugf("[unlockDropWaiters] jetID %v", jetID)

		lock.Lock()
		defer lock.Unlock()

		if waiter == nil {
			inslogger.FromContext(ctx).Debugf("[unlockDropWaiters] waiter == nil, %v", jetID)
			waiter = &jetDropTimeout{
				jetDropLocker: make(chan struct{}),
				timeoutLocker: make(chan struct{}),
			}
			m.jetDropTimeoutProvider.waiters[jetID] = waiter
		}

		inslogger.FromContext(ctx).Debugf("[unlockDropWaiters] before handler %v", time.Now())
		resp, err := handler(ctx, parcel)
		inslogger.FromContext(ctx).Debugf("[unlockDropWaiters] after handler %v", time.Now())

		waiter.setLastPulse(parcel.Pulse())
		close(waiter.jetDropLocker)

		inslogger.FromContext(ctx).Debugf("[unlockDropWaiters] channel unlocked %v", time.Now())

		waiter.jetDropLocker = make(chan struct{})

		return resp, err
	}
}
