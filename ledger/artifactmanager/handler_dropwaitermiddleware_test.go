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
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func TestMiddleware_waitForDrop(t *testing.T){
	// there is a case, when there is only the first pulse (65537)
	// when it happens, we have no data
	// also the same case happens, when we don't know anything about the jet
	t.Run("jetDropTimeout is nil.", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		jetID := core.NewRecordID(core.FirstPulseNumber, []byte{1})

		middleware := newMiddleware(nil, nil, nil)
		expectedParcel := message.Parcel{PulseNumber:8888}
		handler := func(context context.Context, parcel core.Parcel) (rep core.Reply, e error) {
				require.Equal(t, &expectedParcel, parcel)
				return &reply.OK{}, nil
		}

		internal := middleware.waitForDrop(handler)
		rep, err := internal(contextWithJet(ctx, *jetID), &expectedParcel)

		require.Equal(t, &reply.OK{}, rep)
		require.Nil(t, err)
	})


	t.Run("timeout works well", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		jetID := core.NewRecordID(core.FirstPulseNumber, []byte{1})

		middleware := newMiddleware(nil, nil, nil)
		middleware.jetDropTimeoutProvider.waiters[*jetID] = &jetDropTimeout{
			jetDropLocker: make(chan struct{}),
			timeoutLocker: make(chan struct{}),
			lastJdPulse:1,
		}
		expectedParcel := message.Parcel{PulseNumber:8888}
		handler := func(context context.Context, parcel core.Parcel) (rep core.Reply, e error) {
			require.Equal(t, &expectedParcel, parcel)
			return &reply.Object{IsPrototype:true}, nil
		}

		internal := middleware.waitForDrop(handler)
		rep, err := internal(contextWithJet(ctx, *jetID), &expectedParcel)

		require.NoError(t, err)
		require.Equal(t, &reply.Object{IsPrototype:true}, rep)
		require.Equal(t, false, middleware.jetDropTimeoutProvider.waiters[*jetID].isTimeoutRun)
	})

	t.Run("timeout works well, but drop unlock happened", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		jetID := core.NewRecordID(core.FirstPulseNumber, []byte{1})

		middleware := newMiddleware(nil, nil, nil)
		middleware.jetDropTimeoutProvider.waiters[*jetID] = &jetDropTimeout{
			jetDropLocker: make(chan struct{}),
			timeoutLocker: make(chan struct{}),
			lastJdPulse:1,
		}
		expectedParcel := message.Parcel{PulseNumber:8888}
		handler := func(context context.Context, parcel core.Parcel) (rep core.Reply, e error) {
			require.Equal(t, &expectedParcel, parcel)
			return &reply.Object{IsPrototype:true}, nil
		}

		internal := middleware.waitForDrop(handler)
		go func() {
			time.Sleep(300 * time.Millisecond)
			close(middleware.jetDropTimeoutProvider.waiters[*jetID].jetDropLocker)
		}()
		rep, err := internal(contextWithJet(ctx, *jetID), &expectedParcel)

		require.NoError(t, err)
		require.Equal(t, &reply.Object{IsPrototype:true}, rep)
		require.Equal(t, false, middleware.jetDropTimeoutProvider.waiters[*jetID].isTimeoutRun)
	})

	t.Run("unlockDropWaiters", func(t *testing.T) {
		t.Run("init if nil", func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			jetID := core.NewRecordID(core.FirstPulseNumber, []byte{1})

			middleware := newMiddleware(nil, nil, nil)
			expectedParcel := message.Parcel{PulseNumber:8888}
			handler := func(context context.Context, parcel core.Parcel) (rep core.Reply, e error) {
				require.Equal(t, &expectedParcel, parcel)
				return &reply.Object{IsPrototype:true}, nil
			}

			internal := middleware.unlockDropWaiters(handler)
			rep, err := internal(contextWithJet(ctx, *jetID), &expectedParcel)

			require.Nil(t, err)
			require.Equal(t, &reply.Object{IsPrototype:true}, rep)

			waiter := middleware.jetDropTimeoutProvider.waiters[*jetID]
			require.NotNil(t, waiter)
			require.NotNil(t, waiter.jetDropLocker)
			require.NotNil(t, waiter.timeoutLocker)
			require.Equal(t, 8888, int(waiter.lastJdPulse))
		})

		t.Run("works well, when jetDropTimeout for jet isn't inited", func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			jetID := core.NewRecordID(core.FirstPulseNumber, []byte{1})

			middleware := newMiddleware(nil, nil, nil)
			middleware.jetDropTimeoutProvider.waiters[*jetID] = &jetDropTimeout{
				jetDropLocker: make(chan struct{}),
				timeoutLocker: make(chan struct{}),
				lastJdPulse: core.PulseNumber(7777),
			}
			expectedParcel := message.Parcel{PulseNumber:8888}
			handler := func(context context.Context, parcel core.Parcel) (rep core.Reply, e error) {
				require.Equal(t, &expectedParcel, parcel)
				return &reply.Object{IsPrototype:true}, nil
			}

			internal := middleware.unlockDropWaiters(handler)
			rep, err := internal(contextWithJet(ctx, *jetID), &expectedParcel)

			require.Nil(t, err)
			require.Equal(t, &reply.Object{IsPrototype:true}, rep)

			waiter := middleware.jetDropTimeoutProvider.waiters[*jetID]
			require.NotNil(t, waiter)
			require.NotNil(t, waiter.jetDropLocker)
			require.NotNil(t, waiter.timeoutLocker)
			require.Equal(t, 8888, int(waiter.lastJdPulse))
		})
	})


}