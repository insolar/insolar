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
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type earlyRequestCircuitBreakerProvider struct {
	lock     sync.Mutex
	breakers map[core.RecordID]*requestCircuitBreakerProvider
}

type requestCircuitBreakerProvider struct {
	hotDataChannel chan struct{}
	timeoutChannel chan struct{}
}

func (b *earlyRequestCircuitBreakerProvider) getBreaker(jetID core.RecordID) *requestCircuitBreakerProvider {
	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.breakers[jetID]; !ok {
		b.breakers[jetID] = &requestCircuitBreakerProvider{
			hotDataChannel: make(chan struct{}),
			timeoutChannel: make(chan struct{}),
		}
	}

	return b.breakers[jetID]
}

func (b *earlyRequestCircuitBreakerProvider) onTimeoutHappened() {
	b.lock.Lock()
	defer b.lock.Unlock()

	for _, breaker := range b.breakers {
		close(breaker.timeoutChannel)
	}

	b.breakers = map[core.RecordID]*requestCircuitBreakerProvider{}
}

func (m *middleware) checkBreaker(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		logger.Debugf("[waitForDrop] pulse %v starts %v", parcel.Pulse(), time.Now())

		// TODO: 15.01.2019 @egorikas
		// Hack is needed for genesis
		if parcel.Pulse() == core.FirstPulseNumber {
			return handler(ctx, parcel)
		}

		// If the call is a call in redirect-chain
		// skip waiting for the hot records
		if parcel.DelegationToken() != nil {
			logger.Debugf("[waitForDrop] parcel.DelegationToken() != nil")
			return handler(ctx, parcel)
		}

		jetID := jetFromContext(ctx)
		requestBreaker := m.earlyRequestCircuitBreakerProvider.getBreaker(jetID)

		select {
		case <-requestBreaker.hotDataChannel:
		case <-requestBreaker.timeoutChannel:
			return &reply.Error{ErrType: reply.ErrHotDataTimeout}, nil
		}

		logger.Debugf("[waitForDrop] before handler exec - %v", time.Now())
		return handler(ctx, parcel)
	}
}

func (m *middleware) closeBreaker(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		logger.Debugf("[unlockDropWaiters] pulse %v starts %v", parcel.Pulse(), time.Now())

		hotDataMessage := parcel.Message().(*message.HotData)
		jetID := hotDataMessage.DropJet

		breaker := m.earlyRequestCircuitBreakerProvider.getBreaker(jetID)

		logger.Debugf("[unlockDropWaiters] before handler %v", time.Now())
		resp, err := handler(ctx, parcel)
		logger.Debugf("[unlockDropWaiters] after handler %v", time.Now())

		close(breaker.hotDataChannel)

		return resp, err
	}
}
