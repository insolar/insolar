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

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/longbits"
)

func TestNewFirstPulsarData(t *testing.T) {
	delta := uint16(1)
	entropy := longbits.Bits256{2}
	require.Panics(t, func() { NewFirstPulsarData(0, entropy) })

	pd := NewFirstPulsarData(delta, entropy)
	require.Equal(t, entropy, pd.DataExt.PulseEntropy)

	require.Equal(t, delta, pd.DataExt.NextPulseDelta)

	require.Zero(t, pd.DataExt.PrevPulseDelta)
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

	require.Zero(t, pd.Timestamp)

	require.Equal(t, uint16(1), pd.NextPulseDelta)

	require.Zero(t, pd.PrevPulseDelta)
}

func TestString(t *testing.T) {
	delta := uint16(1)
	entropy := longbits.Bits256{2}

	pd := NewFirstPulsarData(delta, entropy)
	require.NotEmpty(t, pd.String())

	pd.PulseNumber = MaxTimePulse + 2
	require.NotEmpty(t, pd.String())

	pd.PulseNumber = MaxTimePulse
	pd.PrevPulseDelta = pd.NextPulseDelta
	require.NotEmpty(t, pd.String())

	pd.NextPulseDelta = 0
	require.NotEmpty(t, pd.String())
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

	require.Zero(t, pd.Timestamp)

	require.Equal(t, uint16(1), pd.NextPulseDelta)

	require.Zero(t, pd.PrevPulseDelta)
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
	require.Panics(t, func() { pd.EnsurePulsarData() })

	pd.PulseNumber = Number(MinTimePulse)
	pd.PulseEpoch = MaxTimePulse + 1
	require.Panics(t, func() { pd.EnsurePulsarData() })

	pd.PulseEpoch = MaxTimePulse
	pd.NextPulseDelta = 0
	require.Panics(t, func() { pd.EnsurePulsarData() })

	pd.NextPulseDelta = 1
	require.NotPanics(t, func() { pd.EnsurePulsarData() })
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

func TestEnsureEphemeralData(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Panics(t, func() { pd.EnsureEphemeralData() })

	pd.PulseEpoch = EphemeralPulseEpoch
	pd.NextPulseDelta = 0
	require.Panics(t, func() { pd.EnsureEphemeralData() })

	pd.NextPulseDelta = 1
	require.NotPanics(t, func() { pd.EnsureEphemeralData() })
}

func TestIsValidEphemeralData(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsValidEphemeralData())

	pd.PulseEpoch = EphemeralPulseEpoch
	pd.NextPulseDelta = 0
	require.False(t, pd.IsValidEphemeralData())

	pd.NextPulseDelta = 1
	require.True(t, pd.IsValidEphemeralData())
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

	pd.PulseEpoch = MaxTimePulse
	pd.PulseNumber = MaxTimePulse
	require.Equal(t, Number(MaxTimePulse), pd.GetStartOfEpoch())
}

func entropyGenTest() longbits.Bits256 {
	return longbits.Bits256{3}
}

func TestCreateNextPulse(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Panics(t, func() { pd.CreateNextPulse(entropyGenTest) })

	pd.PulseNumber = Number(MinTimePulse)
	pd.PulseEpoch = EphemeralPulseEpoch
	d := pd.CreateNextPulse(entropyGenTest)
	require.Equal(t, d.PrevPulseDelta, pd.NextPulseDelta)

	require.Zero(t, d.Timestamp)

	pd.PulseEpoch = MaxTimePulse
	d = pd.CreateNextPulse(entropyGenTest)
	require.Equal(t, d.PrevPulseDelta, pd.NextPulseDelta)
	require.NotZero(t, d.Timestamp)
}

func TestIsValidNext(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd1 := newPulsarData(pn, delta, entropy)
	pd2 := newPulsarData(pn+Number(delta), delta, entropy)
	pd2.PrevPulseDelta = delta
	require.True(t, pd1.IsValidNext(pd2))

	pd2.PrevPulseDelta = 1
	require.False(t, pd1.IsValidNext(pd2))

	pd2.PrevPulseDelta = delta
	pd2.NextPulseDelta = 0
	require.False(t, pd1.IsValidNext(pd2))

	pd2.NextPulseDelta = delta
	pd1.NextPulseDelta = delta + 1
	require.False(t, pd1.IsValidNext(pd2))

	pd1.NextPulseDelta = delta
	pd2.PulseNumber = pn + 3
	require.False(t, pd1.IsValidNext(pd2))

	pd2.PulseNumber = pn + 2
	pd1.PulseNumber = MinTimePulse - 1
	require.Panics(t, func() { pd1.IsValidNext(pd2) })

	delta = 1
	pd1 = newEphemeralData(pn)
	pd2 = newEphemeralData(pn + Number(delta))
	pd2.PrevPulseDelta = delta
	require.True(t, pd1.IsValidNext(pd2))
}

func TestIsValidPrev(t *testing.T) {
	pn := Number(MaxTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd1 := newPulsarData(pn, delta, entropy)
	pd2 := newPulsarData(pn-Number(delta), delta, entropy)
	pd1.PrevPulseDelta = delta
	require.True(t, pd1.IsValidPrev(pd2))

	pd2.NextPulseDelta = 1
	require.False(t, pd1.IsValidPrev(pd2))

	pd2.NextPulseDelta = delta
	pd1.PrevPulseDelta = 0
	require.False(t, pd1.IsValidPrev(pd2))

	pd1.PrevPulseDelta = delta
	pd2.NextPulseDelta = delta - 1
	require.False(t, pd1.IsValidPrev(pd2))

	pd2.NextPulseDelta = delta
	pd2.PulseNumber = pn - 3
	require.False(t, pd1.IsValidPrev(pd2))

	pd2.PulseNumber = MaxTimePulse + 1
	require.Panics(t, func() { pd1.IsValidPrev(pd2) })

	delta = 1
	pd1 = newEphemeralData(pn)
	pd2 = newEphemeralData(pn - Number(delta))
	pd1.PrevPulseDelta = delta
	require.True(t, pd1.IsValidPrev(pd2))
}

func TestGetNextPulseNumber(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Panics(t, func() { pd.GetNextPulseNumber() })

	pd.PulseNumber = MaxTimePulse + 1
	require.Panics(t, func() { pd.GetNextPulseNumber() })

	pd.PulseNumber = MaxTimePulse - 1
	require.Panics(t, func() { pd.GetNextPulseNumber() })

	pd.NextPulseDelta = 0
	require.Panics(t, func() { pd.GetNextPulseNumber() })

	pd.NextPulseDelta = delta
	pd.PulseNumber = MinTimePulse
	require.Equal(t, MinTimePulse+Number(delta), pd.GetNextPulseNumber())
}

func TestGetPrevPulseNumber(t *testing.T) {
	pn := Number(MinTimePulse - 1)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Panics(t, func() { pd.GetPrevPulseNumber() })

	pd.PulseNumber = MaxTimePulse + 1
	require.Panics(t, func() { pd.GetPrevPulseNumber() })

	pd.PulseNumber = MinTimePulse + 1
	require.Panics(t, func() { pd.GetPrevPulseNumber() })

	pd.PrevPulseDelta = 0
	require.Panics(t, func() { pd.GetPrevPulseNumber() })

	pd.PrevPulseDelta = delta
	pd.PulseNumber = MaxTimePulse
	require.Equal(t, MaxTimePulse-Number(delta), pd.GetPrevPulseNumber())
}

func TestCreateNextExpected(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	cne := pd.CreateNextExpected()
	require.Equal(t, MinTimePulse+Number(delta), cne.PulseNumber)

	require.Equal(t, delta, cne.PrevPulseDelta)

	pd.PulseEpoch = EphemeralPulseEpoch
	cne = pd.CreateNextExpected()
	require.Equal(t, MinTimePulse+Number(delta), cne.PulseNumber)

	require.Equal(t, delta, cne.PrevPulseDelta)

	require.Equal(t, EphemeralPulseEpoch, cne.PulseEpoch)
}

func TestCreateNextEphemeralPulse(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Panics(t, func() { pd.CreateNextEphemeralPulse() })

	pd.PulseEpoch = EphemeralPulseEpoch
	cne := pd.CreateNextEphemeralPulse()
	require.Equal(t, EphemeralPulseEpoch, cne.PulseEpoch)

	require.Equal(t, pn+Number(delta), cne.PulseNumber)
}

func TestCcreateNextEphemeralPulse(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.PulseEpoch = EphemeralPulseEpoch
	cne := pd.CreateNextEphemeralPulse()
	require.Equal(t, EphemeralPulseEpoch, cne.PulseEpoch)

	require.Equal(t, pn+Number(delta), cne.PulseNumber)
}

func TestCreateNextPulsarPulse(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	cnp := pd.CreateNextPulsarPulse(delta, entropyGenTest)
	require.NotEqual(t, EphemeralPulseEpoch, cnp.PulseEpoch)

	require.Equal(t, pn+Number(delta), cnp.PulseNumber)

	pd.PulseEpoch = EphemeralPulseEpoch
	require.Panics(t, func() { pd.CreateNextPulsarPulse(delta, entropyGenTest) })
}

func TestCcreateNextPulsarPulse(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	cnp := pd.createNextPulsarPulse(delta, entropyGenTest)
	require.NotEqual(t, EphemeralPulseEpoch, cnp.PulseEpoch)

	require.Equal(t, pn+Number(delta), cnp.PulseNumber)
}

func TestGetPulseNumber(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Equal(t, pn, pd.GetPulseNumber())
}

func TestGetNextPulseDelta(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.NextPulseDelta = 3
	require.Equal(t, uint16(3), pd.GetNextPulseDelta())
}

func TestGetPrevPulseDelta(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.PrevPulseDelta = 5
	require.Equal(t, uint16(5), pd.GetPrevPulseDelta())
}

func TestGetTimestamp(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.Timestamp = 5
	require.Equal(t, uint64(5), pd.GetTimestamp())
}

func TestIsExpectedPulse(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsExpectedPulse())

	pd.NextPulseDelta = 0
	pd.PulseNumber = Number(MinTimePulse - 1)
	require.False(t, pd.IsExpectedPulse())

	pd.PulseNumber = Number(MinTimePulse)
	require.True(t, pd.IsExpectedPulse())
}

func TestIsFirstPulse(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	pd.PrevPulseDelta = 1
	require.False(t, pd.IsFirstPulse())

	pd.PrevPulseDelta = 0
	pd.PulseNumber = Number(MinTimePulse - 1)
	require.False(t, pd.IsFirstPulse())

	pd.PulseNumber = Number(MinTimePulse)
	require.True(t, pd.IsFirstPulse())
}
