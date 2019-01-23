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
	logger.Debugf("[breakermiddleware] [getBreaker] jetID - %v", jetID.DebugString())

	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.breakers[jetID]; !ok {
		logger.Debugf("[breakermiddleware] [getBreaker] create new  - %v", jetID.DebugString())
		b.breakers[jetID] = &requestCircuitBreakerProvider{
			hotDataChannel: make(chan struct{}),
			timeoutChannel: make(chan struct{}),
		}
	}

	return b.breakers[jetID]
}

func (b *earlyRequestCircuitBreakerProvider) onTimeoutHappened(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[breakermiddleware] [onTimeoutHappened] start method. breakers should be cleared")

	b.lock.Lock()
	defer b.lock.Unlock()

	for jetID, breaker := range b.breakers {
		logger.Debugf("[breakermiddleware] [onTimeoutHappened] shutdown timeout channel for jetID - %v", jetID.DebugString())
		close(breaker.timeoutChannel)
	}

	b.breakers = map[core.RecordID]*requestCircuitBreakerProvider{}
}

func (m *middleware) checkEarlyRequestBreaker(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		logger.Debugf("[breakermiddleware] [checkEarlyRequestBreaker] for parcel with pulse %v", parcel.Pulse())

		// TODO: 15.01.2019 @egorikas
		// Hack is needed for genesis
		if parcel.Pulse() == core.FirstPulseNumber {
			return handler(ctx, parcel)
		}

		// If the call is a call in redirect-chain
		// skip waiting for the hot records
		if parcel.DelegationToken() != nil {
			logger.Debugf("[breakermiddleware] [checkEarlyRequestBreaker] parcel.DelegationToken() != nil")
			return handler(ctx, parcel)
		}

		jetID := jetFromContext(ctx)
		requestBreaker := m.earlyRequestCircuitBreakerProvider.getBreaker(ctx, jetID)

		logger.Debugf(
			"[breakermiddleware] [checkEarlyRequestBreaker] before pause of request with jet - %v, pulse - %v, type - %v",
			jetID.DebugString(),
			parcel.Pulse(),
			parcel.Message().Type(),
		)
		select {
		case <-requestBreaker.hotDataChannel:
			logger.Debugf("[breakermiddleware] [checkEarlyRequestBreaker] before handler exec - %v", time.Now())
			return handler(ctx, parcel)
		case <-requestBreaker.timeoutChannel:
			logger.Errorf("[breakermiddleware] [checkEarlyRequestBreaker] timeout happened for %v with pulse  %v", jetID.DebugString(), parcel.Pulse())
			return &reply.Error{ErrType: reply.ErrHotDataTimeout}, nil
		}
	}
}

func (m *middleware) closeEarlyRequestBreaker(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		logger.Debugf("[breakermiddleware] [closeEarlyRequestBreaker] pulse %v starts %v", parcel.Pulse(), time.Now())

		hotDataMessage := parcel.Message().(*message.HotData)
		jetID := hotDataMessage.Jet.Record()

		logger.Debugf("[breakermiddleware] [closeEarlyRequestBreaker] hot data for jet happens - %v, pulse - %v", jetID.DebugString(), parcel.Pulse())
		breaker := m.earlyRequestCircuitBreakerProvider.getBreaker(ctx, *jetID)
		defer close(breaker.hotDataChannel)

		logger.Debugf("[breakermiddleware] [closeEarlyRequestBreaker] before handler for jet - %v, pulse - %v", jetID.DebugString(), parcel.Pulse())
		return handler(ctx, parcel)
	}
}

func (m *middleware) closeEarlyRequestBreakerForJet(ctx context.Context, jetID core.RecordID) {
	inslogger.FromContext(ctx).Debugf("[breakermiddleware] [closeEarlyRequestBreakerForJet] jetID - %v", jetID.DebugString())
	breaker := m.earlyRequestCircuitBreakerProvider.getBreaker(ctx, jetID)
	close(breaker.hotDataChannel)
}
