///
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
///

package pulse

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/common/longbits"

	"github.com/stretchr/testify/require"
)

func TestNewFirstPulsarData(t *testing.T) {
	delta := uint16(1)
	entropy := longbits.Bits256{2}
	require.Panics(t, func() { NewFirstPulsarData(0, entropy) })

	pd := NewFirstPulsarData(delta, entropy)
	require.Equal(t, entropy, pd.DataExt.PulseEntropy)

	require.Equal(t, delta, pd.DataExt.NextPulseDelta)

	require.Equal(t, uint16(0), pd.DataExt.PrevPulseDelta)
}

func TestNewPulsarData(t *testing.T) {
	pn := Number(1)
	deltaNext := uint16(2)
	deltaPrev := uint16(3)
	entropy := longbits.Bits256{4}
	pd := NewPulsarData(pn, deltaNext, deltaPrev, entropy)
	require.Equal(t, pn, pd.PulseNumber)

	require.Equal(t, uint32(pn), pd.DataExt.PulseEpoch)

	require.Equal(t, entropy, pd.DataExt.PulseEntropy)

	require.Equal(t, deltaNext, pd.DataExt.NextPulseDelta)

	require.Equal(t, deltaPrev, pd.PrevPulseDelta)
}

func TestNewFirstEphemeralData(t *testing.T) {
	pd := NewFirstEphemeralData()
	require.Equal(t, Number(MinTimePulse), pd.PulseNumber)

	require.Equal(t, EphemeralPulseEpoch, pd.PulseEpoch)

	require.Equal(t, uint32(0), pd.Timestamp)

	require.Equal(t, uint16(1), pd.NextPulseDelta)

	require.Equal(t, uint16(0), pd.PrevPulseDelta)
}

func TestString(t *testing.T) {
	delta := uint16(1)
	entropy := longbits.Bits256{2}

	pd := NewFirstPulsarData(delta, entropy)
	require.True(t, pd.String() != "")
}

func TestNnewPulsarData(t *testing.T) {
	pn := Number(1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Panics(t, func() { NewFirstPulsarData(0, entropy) })

	require.Equal(t, pn, pd.PulseNumber)

	require.Equal(t, uint32(pn), pd.DataExt.PulseEpoch)

	require.Equal(t, entropy, pd.DataExt.PulseEntropy)

	require.Equal(t, delta, pd.DataExt.NextPulseDelta)
}

func TestNewEphemeralData(t *testing.T) {
	pn := Number(1)
	pd := newEphemeralData(pn)
	require.Equal(t, pn, pd.PulseNumber)

	require.Equal(t, EphemeralPulseEpoch, pd.PulseEpoch)

	require.Equal(t, uint32(0), pd.Timestamp)

	require.Equal(t, uint16(1), pd.NextPulseDelta)

	require.Equal(t, uint16(0), pd.PrevPulseDelta)
}

func TestFixedPulseEntropy(t *testing.T) {
	bits := longbits.Bits256{1}
	pn := Number(1)
	require.NotPanics(t, func() { fixedPulseEntropy(&bits, pn) })
}

func TestEnsurePulseData(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Panics(t, func() { pd.EnsurePulseData() })

	pd.PulseNumber = Number(MinTimePulse)
	pd.PulseEpoch = MaxTimePulse + 1
	require.Panics(t, func() { pd.EnsurePulseData() })

	pd.PulseEpoch = MaxTimePulse
	pd.NextPulseDelta = 0
	require.Panics(t, func() { pd.EnsurePulseData() })

	pd.NextPulseDelta = 1
	require.NotPanics(t, func() { pd.EnsurePulseData() })
}

func TestIsValidPulseData(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsValidPulseData())

	pd.PulseNumber = Number(MinTimePulse)
	pd.PulseEpoch = MaxTimePulse + 1
	require.False(t, pd.IsValidPulseData())

	pd.PulseEpoch = MaxTimePulse
	pd.NextPulseDelta = 0
	require.False(t, pd.IsValidPulseData())

	pd.NextPulseDelta = 1
	require.True(t, pd.IsValidPulseData())
}

func TestIsEmpty(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsEmpty())

	pd.PulseNumber = Unknown
	require.True(t, pd.IsEmpty())
}

func TestIsValidExpectedPulseData(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsValidExpectedPulseData())

	pd.PulseNumber = Number(MinTimePulse)
	pd.PulseEpoch = MaxTimePulse + 1
	require.False(t, pd.IsValidExpectedPulseData())

	pd.PulseEpoch = MaxTimePulse
	pd.PrevPulseDelta = 1
	require.False(t, pd.IsValidExpectedPulseData())

	pd.PrevPulseDelta = 0
	require.True(t, pd.IsValidExpectedPulseData())
}

func TestEnsurePulsarData(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.PulseEpoch = MaxTimePulse + 1
	require.Panics(t, func() { pd.EnsurePulsarData() })

	pd.PulseEpoch = MaxTimePulse
	pd.PulseNumber = Number(MinTimePulse - 1)
	require.Panics(t, func() { pd.EnsurePulseData() })

	pd.PulseNumber = Number(MinTimePulse)
	pd.PulseEpoch = MaxTimePulse + 1
	require.Panics(t, func() { pd.EnsurePulseData() })

	pd.PulseEpoch = MaxTimePulse
	pd.NextPulseDelta = 0
	require.Panics(t, func() { pd.EnsurePulseData() })

	pd.NextPulseDelta = 1
	require.NotPanics(t, func() { pd.EnsurePulseData() })
}

func TestIsValidPulsarData(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.PulseEpoch = MaxTimePulse + 1
	require.False(t, pd.IsValidPulsarData())

	pd.PulseEpoch = MaxTimePulse
	pd.PulseNumber = Number(MinTimePulse - 1)
	require.False(t, pd.IsValidPulsarData())

	pd.PulseNumber = Number(MinTimePulse)
	pd.PulseEpoch = MaxTimePulse + 1
	require.False(t, pd.IsValidPulsarData())

	pd.PulseEpoch = MaxTimePulse
	pd.NextPulseDelta = 0
	require.False(t, pd.IsValidPulsarData())

	pd.NextPulseDelta = 1
	require.True(t, pd.IsValidPulsarData())
}

func TestIsValidEphemeralData(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsValidEphemeralData())

	pd.PulseEpoch = EphemeralPulseEpoch
	pd.NextPulseDelta = 0
	require.False(t, pd.IsValidPulseData())

	pd.NextPulseDelta = 1
	require.True(t, pd.IsValidPulseData())
}

func TestIsFromPulsar(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsFromPulsar())

	pd.PulseNumber = MinTimePulse
	pd.PulseEpoch = MaxTimePulse + 1
	require.False(t, pd.IsFromPulsar())

	pd.PulseEpoch = MaxTimePulse
	require.True(t, pd.IsFromPulsar())
}

func TestIsFromEphemeral(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.PulseEpoch = EphemeralPulseEpoch
	require.False(t, pd.IsFromEphemeral())

	pd.PulseNumber = MinTimePulse
	pd.PulseEpoch = MaxTimePulse
	require.False(t, pd.IsFromEphemeral())

	pd.PulseEpoch = EphemeralPulseEpoch
	require.True(t, pd.IsFromEphemeral())
}

func TestGetStartOfEpoch(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.PulseEpoch = MaxTimePulse + 1
	require.Equal(t, Number(1<<16), pd.GetStartOfEpoch())

	pd.PulseNumber = MaxTimePulse + 1
	require.Equal(t, Number(MaxTimePulse+1), pd.GetStartOfEpoch())
}
