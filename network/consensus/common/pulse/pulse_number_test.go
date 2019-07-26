//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
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

func TestAsApproximateTime(t *testing.T) {
	n := Number(0)
	approx := n.AsApproximateTime()
	require.Equal(t, timeOfMinTimePulse, approx)
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
	n := Number(Unknown)
	require.False(t, n.IsSpecialOrTimePulse())

	n = Number(MaxTimePulse + 1)
	require.False(t, n.IsSpecialOrTimePulse())

	n = MinTimePulse
	require.True(t, n.IsSpecialOrTimePulse())
}

func TestIsUnknown(t *testing.T) {
	n := Number(Unknown)
	require.True(t, n.IsUnknown())

	n = Number(MaxTimePulse)
	require.False(t, n.IsUnknown())
}

func TestIsUnknownOrEqualTo(t *testing.T) {
	n1 := Number(Unknown)
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
