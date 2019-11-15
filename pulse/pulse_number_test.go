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

package pulse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOfNow(t *testing.T) {
	require.True(t, OfNow() > 0)
}

func TestOfTime(t *testing.T) {
	require.Equal(t, Unknown, OfTime(time.Time{}))

	require.True(t, OfTime(time.Now()) > MinTimePulse)
}

func TestOfUnixTime(t *testing.T) {
	require.Equal(t, Unknown, OfUnixTime(0))

	require.Equal(t, Number(MinTimePulse), OfUnixTime(UnixTimeOfMinTimePulse))
}

func TestAsApproximateTime_FromOfUnixTime(t *testing.T) {
	ts := int64(1595808000) // 27.07.2020
	number := OfUnixTime(ts)
	newTs, err := number.AsApproximateTime()
	require.NoError(t, err)
	require.Equal(t, ts, newTs.Unix())
}

func TestAsApproximateTime(t *testing.T) {
	t.Run("pulse less than minimal", func(t *testing.T) {
		n := Number(0)
		_, err := n.AsApproximateTime()
		require.Error(t, err)
	})

	t.Run("pulse greater than maximal", func(t *testing.T) {
		n := Number(0xFFFFFFFF)
		_, err := n.AsApproximateTime()
		require.Error(t, err)
	})

	t.Run("ordinary", func(t *testing.T) {
		n := Number(MinTimePulse)
		approx, err := n.AsApproximateTime()
		require.NoError(t, err)
		require.Equal(t, timeOfMinTimePulse, approx)
	})
}

func TestIsTimePulse(t *testing.T) {
	n := Number(MinTimePulse)
	require.True(t, n.IsTimePulse())

	n = Number(MinTimePulse - 1)
	require.False(t, n.IsTimePulse())

	n = Number(MaxTimePulse + 1)
	require.False(t, n.IsTimePulse())
}

func TestAsUint32(t *testing.T) {
	n := Number(MinTimePulse)
	require.Equal(t, uint32(MinTimePulse), n.AsUint32())
}

func TestIsSpecialOrTimePulse(t *testing.T) {
	n := Unknown
	require.False(t, n.IsSpecialOrTimePulse())

	n = Number(MaxTimePulse + 1)
	require.False(t, n.IsSpecialOrTimePulse())

	n = MinTimePulse
	require.True(t, n.IsSpecialOrTimePulse())
}

func TestIsUnknown(t *testing.T) {
	n := Unknown
	require.True(t, n.IsUnknown())

	n = Number(MaxTimePulse)
	require.False(t, n.IsUnknown())
}

func TestIsUnknownOrEqualTo(t *testing.T) {
	n1 := Unknown
	n2 := Number(MaxTimePulse)
	require.True(t, n1.IsUnknownOrEqualTo(n2))

	n1 = Number(MinTimePulse)
	require.False(t, n1.IsUnknownOrEqualTo(n2))

	n2 = Number(MinTimePulse)
	require.True(t, n1.IsUnknownOrEqualTo(n2))
}

func TestNext(t *testing.T) {
	delta := uint16(2)
	require.Panics(t, func() { Number(MinTimePulse - 1).Next(delta) })

	require.Panics(t, func() { Number(MaxTimePulse - 1).Next(delta) })

	require.Equal(t, Number(MaxTimePulse), (MaxTimePulse - Number(delta)).Next(delta))
}

func TestPrev(t *testing.T) {
	delta := uint16(2)
	require.Panics(t, func() { Number(MinTimePulse - 1).Prev(delta) })

	require.Panics(t, func() { Number(MinTimePulse).Prev(delta) })

	require.Equal(t, Number(MinTimePulse), (MinTimePulse + Number(delta)).Prev(delta))
}

func TestIsValidAsPulseNumber(t *testing.T) {
	require.False(t, IsValidAsPulseNumber(MinTimePulse-1))

	require.False(t, IsValidAsPulseNumber(MaxTimePulse+1))

	require.True(t, IsValidAsPulseNumber(MinTimePulse+1))
}

func TestOfInt(t *testing.T) {
	require.Zero(t, OfInt(MaxTimePulse+1))

	require.Equal(t, Number(MaxTimePulse-1), OfInt(MaxTimePulse-1))
}

func TestOfUint32(t *testing.T) {
	require.Zero(t, OfUint32(MaxTimePulse+1))

	require.Equal(t, Number(MaxTimePulse-1), OfUint32(MaxTimePulse-1))
}

func TestFlagsOf(t *testing.T) {
	require.Zero(t, FlagsOf(MaxTimePulse))

	require.Equal(t, uint(1), FlagsOf(MaxTimePulse+2))
}
