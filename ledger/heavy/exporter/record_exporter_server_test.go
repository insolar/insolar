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

package exporter

import (
	"errors"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"
)

func TestRecordIterator_HasNext(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	t.Run("returns false, if LastKnownPosition returns error", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(0, errors.New("some error"))

		iter := newRecordIterator(pn, 0, 0, positionAccessor, nil, nil, nil)

		hasNext := iter.HasNext(ctx)

		require.False(t, hasNext)
	})

	t.Run("returns false, if read all the count", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(156, nil)

		iter := newRecordIterator(pn, 0, 10, positionAccessor, nil, nil, nil)
		iter.read = 11

		hasNext := iter.HasNext(ctx)

		require.False(t, hasNext)
	})

	t.Run("returns true, if read not all the count", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(156, nil)

		iter := newRecordIterator(pn, 0, 10, positionAccessor, nil, nil, nil)
		iter.read = 9

		hasNext := iter.HasNext(ctx)

		require.True(t, hasNext)
	})

	t.Run("cross-pulse situations", func(t *testing.T) {
		t.Run("no data in the current.no further pulses. returns false", func(t *testing.T) {
			pn := gen.PulseNumber()
			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(pn).Return(1, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{}, store.ErrNotFound)

			iter := newRecordIterator(pn, 0, 0, positionAccessor, nil, nil, pulseCalculator)
			iter.currentPosition = 2

			hasNext := iter.HasNext(ctx)

			require.False(t, hasNext)
		})

		t.Run("no data in the current.no more synced pulses. returns false", func(t *testing.T) {
			pn := gen.PulseNumber()
			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(pn).Return(1, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{PulseNumber: 100}, nil)

			jetKeeper := executor.NewJetKeeperMock(t)
			jetKeeper.TopSyncPulseMock.Return(99)

			iter := newRecordIterator(pn, 0, 0, positionAccessor, nil, jetKeeper, pulseCalculator)
			iter.currentPosition = 2

			hasNext := iter.HasNext(ctx)

			require.False(t, hasNext)
		})

		t.Run("no data in the current. has more synce pulses. returns true", func(t *testing.T) {
			pn := gen.PulseNumber()
			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(pn).Return(1, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{PulseNumber: 100}, nil)

			jetKeeper := executor.NewJetKeeperMock(t)
			jetKeeper.TopSyncPulseMock.Return(101)

			iter := newRecordIterator(pn, 0, 0, positionAccessor, nil, jetKeeper, pulseCalculator)
			iter.currentPosition = 2

			hasNext := iter.HasNext(ctx)

			require.True(t, hasNext)
		})

	})
}
