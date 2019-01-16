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

func (b *earlyRequestCircuitBreakerProvider) getBreaker(ctx context.Context, jetID core.RecordID) *requestCircuitBreakerProvider {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[getBreaker] jetID - %v", jetID.JetIDString())

	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.breakers[jetID]; !ok {
		logger.Debugf("[getBreaker] create new  - %v", jetID.JetIDString())
		b.breakers[jetID] = &requestCircuitBreakerProvider{
			hotDataChannel: make(chan struct{}),
			timeoutChannel: make(chan struct{}),
		}
	}

	return b.breakers[jetID]
}

func (b *earlyRequestCircuitBreakerProvider) onTimeoutHappened(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[onTimeoutHappened] start method. breakers should be cleared")

	b.lock.Lock()
	defer b.lock.Unlock()

	for jetID, breaker := range b.breakers {
		logger.Debugf("[onTimeoutHappened] shutdown timeout channel for jetID - %v", jetID.JetIDString())
		close(breaker.timeoutChannel)
	}

	b.breakers = map[core.RecordID]*requestCircuitBreakerProvider{}
}

func (m *middleware) checkEarlyRequestBreaker(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		logger.Debugf("[checkEarlyRequestBreaker] pulse %v starts %v", parcel.Pulse(), time.Now())

		// TODO: 15.01.2019 @egorikas
		// Hack is needed for genesis
		if parcel.Pulse() == core.FirstPulseNumber {
			return handler(ctx, parcel)
		}

		// If the call is a call in redirect-chain
		// skip waiting for the hot records
		if parcel.DelegationToken() != nil {
			logger.Debugf("[checkEarlyRequestBreaker] parcel.DelegationToken() != nil")
			return handler(ctx, parcel)
		}

		jetID := jetFromContext(ctx)
		requestBreaker := m.earlyRequestCircuitBreakerProvider.getBreaker(ctx, jetID)

		logger.Debugf("[checkEarlyRequestBreaker] before select, jet - %v", jetID.JetIDString())
		select {
		case <-requestBreaker.hotDataChannel:
			logger.Debugf("[checkEarlyRequestBreaker] before handler exec - %v", time.Now())
			return handler(ctx, parcel)
		case <-requestBreaker.timeoutChannel:
			logger.Errorf("[checkEarlyRequestBreaker] timeout happened for %v with pulse  %v", jetID.JetIDString(), parcel.Pulse())
			return &reply.Error{ErrType: reply.ErrHotDataTimeout}, nil
		}
	}
}

func (m *middleware) closeEarlyRequestBreaker(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		logger.Debugf("[closeEarlyRequestBreaker] pulse %v starts %v", parcel.Pulse(), time.Now())

		hotDataMessage := parcel.Message().(*message.HotData)
		jetID := hotDataMessage.Jet.Record()

		logger.Debugf("[closeEarlyRequestBreaker] wait jet - %v", jetID.JetIDString())
		breaker := m.earlyRequestCircuitBreakerProvider.getBreaker(ctx, *jetID)
		defer close(breaker.hotDataChannel)

		logger.Debugf("[closeEarlyRequestBreaker] before handler %v", time.Now())
		return handler(ctx, parcel)
	}
}

func (m *middleware) closeEarlyRequestBreakerForJet(ctx context.Context, jetID core.RecordID) {
	inslogger.FromContext(ctx).Debugf("[closeEarlyRequestBreakerForJet] jetID - %v", jetID.JetIDString())
	breaker := m.earlyRequestCircuitBreakerProvider.getBreaker(ctx, jetID)
	close(breaker.hotDataChannel)
}
