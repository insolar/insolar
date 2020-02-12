// Copyright 2020 Insolar Network Ltd.
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

package artifacts

import (
	"context"
	"errors"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/stretchr/testify/require"
)

func TestNewPulseAccessorLRU_NewPulseAccessorLRU(t *testing.T) {
	t.Parallel()

	t.Run("init successfully", func(t *testing.T) {
		lru := NewPulseAccessorLRU(nil, nil, 1)

		require.NotNil(t, lru)
	})

	t.Run("fails with 0 size", func(t *testing.T) {
		require.Panics(t, func() {
			NewPulseAccessorLRU(nil, nil, 0)
		})
	})
}

func TestNewPulseAccessorLRU_Latest(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		pulses := pulse.NewAccessorMock(t)
		pulses.LatestMock.Return(insolar.Pulse{PulseNumber: 1234}, nil)
		lru := NewPulseAccessorLRU(pulses, nil, 1)

		p, err := lru.Latest(context.Background())

		require.NoError(t, err)
		require.Equal(t, uint32(1234), p.PulseNumber.AsUint32())
	})

	t.Run("no error", func(t *testing.T) {
		pulses := pulse.NewAccessorMock(t)
		pulses.LatestMock.Return(insolar.Pulse{}, errors.New("problems"))
		lru := NewPulseAccessorLRU(pulses, nil, 1)

		_, err := lru.Latest(context.Background())

		require.Error(t, err)
		require.Equal(t, "problems", err.Error())
	})
}

func TestNewPulseAccessorLRU_ForPulseNumber(t *testing.T) {
	t.Parallel()

	t.Run("pulse exists in cache", func(t *testing.T) {
		pn := insolar.PulseNumber(123)
		p := insolar.Pulse{PulseNumber: pn}
		lru := NewPulseAccessorLRU(nil, nil, 1)
		lru.cache.Add(pn, p)

		saved, err := lru.ForPulseNumber(context.TODO(), pn)

		require.NoError(t, err)
		require.Equal(t, p, saved)
	})

	t.Run("pulse doesn't exist in cache and exists in store", func(t *testing.T) {
		pn := insolar.PulseNumber(123)
		p := insolar.Pulse{PulseNumber: pn}
		pulses := pulse.NewAccessorMock(t)
		pulses.ForPulseNumberMock.Return(p, nil)
		lru := NewPulseAccessorLRU(pulses, nil, 1)

		saved, err := lru.ForPulseNumber(context.TODO(), pn)

		require.NoError(t, err)
		require.Equal(t, p, saved)
	})

	t.Run("pulse doesn't exist in cache and doesn't exist in store", func(t *testing.T) {
		pn := insolar.PulseNumber(123)
		p := insolar.Pulse{PulseNumber: pn}
		pulses := pulse.NewAccessorMock(t)
		pulses.ForPulseNumberMock.Return(insolar.Pulse{}, errors.New("custom error"))
		client := NewClientMock(t)
		client.GetPulseMock.When(context.TODO(), pn).Then(p, nil)
		lru := NewPulseAccessorLRU(pulses, client, 1)

		saved, err := lru.ForPulseNumber(context.TODO(), pn)

		require.NoError(t, err)
		require.Equal(t, p, saved)

		val, ok := lru.cache.Get(pn)
		require.True(t, ok)
		require.Equal(t, p, val.(insolar.Pulse))
	})
}
