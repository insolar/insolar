package pulse

import (
	"math/rand"
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

func TestNewPulsarDataInvalid(t *testing.T) {
	require.Panics(t, func() {
		NewPulsarData(0, 1, 1, longbits.Bits256{})
	})
	require.Panics(t, func() {
		NewPulsarData(1, 1, 1, longbits.Bits256{})
	})
	require.Panics(t, func() {
		NewPulsarData(MinTimePulse-1, 1, 1, longbits.Bits256{})
	})
	require.Panics(t, func() {
		NewPulsarData(MaxTimePulse+1, 1, 1, longbits.Bits256{})
	})
}

func TestNewPulsarData(t *testing.T) {
	pn := Number(MinTimePulse)
	deltaNext := uint16(2)
	deltaPrev := uint16(3)
	entropy := longbits.Bits256{4}
	pd := NewPulsarData(pn, deltaNext, deltaPrev, entropy)
	require.Equal(t, pn, pd.PulseNumber)

	require.Equal(t, pn.AsEpoch(), pd.DataExt.PulseEpoch)

	require.Equal(t, entropy, pd.DataExt.PulseEntropy)

	require.Equal(t, deltaNext, pd.DataExt.NextPulseDelta)

	require.Equal(t, deltaPrev, pd.PrevPulseDelta)
}

func Test_newPulsarData(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.Panics(t, func() { NewFirstPulsarData(0, entropy) })

	require.Equal(t, pn, pd.PulseNumber)

	require.Equal(t, pn.AsEpoch(), pd.DataExt.PulseEpoch)

	require.Equal(t, entropy, pd.DataExt.PulseEntropy)

	require.Equal(t, delta, pd.DataExt.NextPulseDelta)
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
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse, delta, entropy)
	pd.PulseNumber--
	require.Panics(t, func() { pd.EnsurePulseData() })

	pd.PulseNumber = Number(MinTimePulse + 1)
	pd.PulseEpoch = MinTimePulse - 1
	require.Panics(t, func() { pd.EnsurePulseData() })
	pd.PulseEpoch = MaxTimePulse + 1
	require.Panics(t, func() { pd.EnsurePulseData() })
	pd.PulseEpoch = pd.PulseNumber.AsEpoch() + 1
	require.Panics(t, func() { pd.EnsurePulseData() })

	pd.PulseEpoch = pd.PulseNumber.AsEpoch()
	pd.EnsurePulseData()

	pd.NextPulseDelta = 0
	require.Panics(t, func() { pd.EnsurePulseData() })
	pd.NextPulseDelta = 1
	pd.EnsurePulseData()
}

func TestIsValidPulseData(t *testing.T) {
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse, delta, entropy)
	pd.PulseNumber--
	require.False(t, pd.IsValidPulseData())

	pd.PulseNumber = Number(MinTimePulse)

	pd.PulseEpoch = EphemeralPulseEpoch
	require.True(t, pd.IsValidPulseData())

	pd.PulseEpoch = MinTimePulse - 1
	require.False(t, pd.IsValidPulseData())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch()
	require.True(t, pd.IsValidPulseData())

	pd.PulseEpoch++
	require.False(t, pd.IsValidPulseData())

	pd.PulseNumber++
	require.True(t, pd.IsValidPulseData())

	pd.NextPulseDelta = 0
	require.False(t, pd.IsValidPulseData())

	pd.NextPulseDelta = 1
	require.True(t, pd.IsValidPulseData())
}

func TestIsEmpty(t *testing.T) {
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse, delta, entropy)
	require.False(t, pd.IsEmpty())

	pd.PulseNumber = Unknown
	require.True(t, pd.IsEmpty())
}

func TestIsValidExpectedPulseData(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsValidExpectedPulseData())

	pd = pd.CreateNextExpected()
	require.True(t, pd.IsValidExpectedPulseData())

	pd.NextPulseDelta = 1
	require.False(t, pd.IsValidExpectedPulseData())
	pd.NextPulseDelta = 0

	pd.PrevPulseDelta = 0
	require.True(t, pd.IsValidExpectedPulseData())

	pd.PulseEpoch = MinTimePulse - 1
	require.False(t, pd.IsValidExpectedPulseData())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch() + 1
	require.False(t, pd.IsValidExpectedPulseData())

	pd.PulseEpoch = ArticulationPulseEpoch
	require.True(t, pd.IsValidExpectedPulseData())

	pd.PulseEpoch = EphemeralPulseEpoch
	require.True(t, pd.IsValidExpectedPulseData())
}

func TestIsValidExpectedPulsarData(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.False(t, pd.IsValidExpectedPulsarData())

	pd = pd.CreateNextExpected()
	require.True(t, pd.IsValidExpectedPulsarData())

	pd.NextPulseDelta = 1
	require.False(t, pd.IsValidExpectedPulsarData())
	pd.NextPulseDelta = 0

	pd.PrevPulseDelta = 0
	require.True(t, pd.IsValidExpectedPulsarData())

	pd.PulseEpoch = MinTimePulse - 1
	require.False(t, pd.IsValidExpectedPulsarData())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch() + 1
	require.False(t, pd.IsValidExpectedPulsarData())

	pd.PulseEpoch = ArticulationPulseEpoch
	require.False(t, pd.IsValidExpectedPulsarData())

	pd.PulseEpoch = EphemeralPulseEpoch
	require.False(t, pd.IsValidExpectedPulsarData())
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

	pd.PulseEpoch = MinTimePulse - 1
	require.Panics(t, func() { pd.EnsurePulsarData() })

	pd.PulseEpoch = pd.PulseNumber.AsEpoch() + 1
	require.Panics(t, func() { pd.EnsurePulsarData() })

	pd.PulseEpoch = pd.PulseNumber.AsEpoch()
	pd.NextPulseDelta = 0
	require.Panics(t, func() { pd.EnsurePulsarData() })

	pd.NextPulseDelta = 1
	pd.EnsurePulsarData()
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

	pd.PulseEpoch = MinTimePulse - 1
	require.False(t, pd.IsValidPulsarData())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch() - 1
	require.False(t, pd.IsValidPulsarData())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch()
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
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse, delta, entropy)
	pd.PulseNumber--
	require.False(t, pd.IsFromPulsar())

	pd.PulseNumber = MinTimePulse
	pd.PulseEpoch = MaxTimePulse + 1
	require.False(t, pd.IsFromPulsar())

	pd.PulseEpoch = MaxTimePulse
	require.True(t, pd.IsFromPulsar())
}

func TestIsFromEphemeral(t *testing.T) {
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse, delta, entropy)
	pd.PulseNumber--
	pd.PulseEpoch = EphemeralPulseEpoch
	require.False(t, pd.IsFromEphemeral())

	pd.PulseNumber = MinTimePulse
	pd.PulseEpoch = MaxTimePulse
	require.False(t, pd.IsFromEphemeral())

	pd.PulseEpoch = EphemeralPulseEpoch
	require.True(t, pd.IsFromEphemeral())
}

func TestGetStartOfEpoch(t *testing.T) {
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse+1, delta, entropy)
	require.Equal(t, pd.PulseNumber, pd.GetStartOfEpoch())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch() + 1
	require.Equal(t, pd.PulseNumber, pd.GetStartOfEpoch())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch() - 1
	require.Equal(t, pd.PulseNumber-1, pd.GetStartOfEpoch())

	pd.PulseEpoch = MaxTimePulse
	pd.PulseNumber = MaxTimePulse
	require.Equal(t, pd.PulseNumber, pd.GetStartOfEpoch())

	pd.PulseNumber++
	require.Equal(t, Unknown, pd.GetStartOfEpoch())

	pd.PulseNumber = MaxTimePulse
	pd.PulseEpoch = 0
	require.Equal(t, pd.PulseNumber, pd.GetStartOfEpoch())

	pd.PulseEpoch = MinTimePulse - 1
	require.Equal(t, pd.PulseNumber, pd.GetStartOfEpoch())

	pd.PulseEpoch = MinTimePulse
	pd.PulseNumber = MinTimePulse - 1
	require.Equal(t, Unknown, pd.GetStartOfEpoch())
}

func entropyGenTest() longbits.Bits256 {
	return longbits.Bits256{3}
}

func TestCreateNextPulse(t *testing.T) {
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse, delta, entropy)
	pd.PulseNumber--
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

	pd.PulseEpoch = ArticulationPulseEpoch
	require.Panics(t, func() { pd.CreateNextPulse(entropyGenTest) })

	pd.PulseEpoch = 0
	require.Panics(t, func() { pd.CreateNextPulse(entropyGenTest) })
}

func TestIsValidNext(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd1 := newPulsarData(pn, delta, entropy)
	pd2 := newPulsarData(pn+Number(delta), delta, entropy)
	pd2.PrevPulseDelta = delta
	require.True(t, pd1.IsValidNext(pd2))

	pd2.PulseEpoch = 1
	require.False(t, pd1.IsValidNext(pd2))
	pd2.PulseEpoch = pd2.PulseNumber.AsEpoch()

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
	require.False(t, pd1.IsValidNext(pd2))

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

	pd2.PulseNumber = 1
	require.False(t, pd1.IsValidPrev(pd2))

	pd2.PulseNumber = pn - Number(delta)
	pd2.PulseEpoch = 1
	require.False(t, pd1.IsValidPrev(pd2))
	pd2.PulseEpoch = pd2.PulseNumber.AsEpoch()

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
	require.False(t, pd1.IsValidPrev(pd2))

	delta = 1
	pd1 = newEphemeralData(pn)
	pd2 = newEphemeralData(pn - Number(delta))
	pd1.PrevPulseDelta = delta
	require.True(t, pd1.IsValidPrev(pd2))
}

func TestGetNextPulseNumber(t *testing.T) {
	pd := newPulsarData(MinTimePulse, 2, longbits.Bits256{3})

	pn, ok := pd.GetNextPulseNumber()
	require.True(t, ok)
	require.Equal(t, Number(MinTimePulse+2), pn)

	epd := pd.CreateNextExpected()
	pn, ok = epd.GetNextPulseNumber()
	require.False(t, ok)
	require.Equal(t, Number(MinTimePulse+2), pn)

	pd = newPulsarData(MaxTimePulse, 2, longbits.Bits256{3})
	pn, ok = pd.GetNextPulseNumber()
	require.False(t, ok)
	require.Equal(t, Number(MaxTimePulse), pn)

	pd = newPulsarData(MaxTimePulse-1, 2, longbits.Bits256{3})
	pn, ok = pd.GetNextPulseNumber()
	require.False(t, ok)
	require.Equal(t, Number(MaxTimePulse), pn)

	pd = newPulsarData(MaxTimePulse-2, 2, longbits.Bits256{3})
	pn, ok = pd.GetNextPulseNumber()
	require.True(t, ok)
	require.Equal(t, Number(MaxTimePulse), pn)
}

func TestGetPrevPulseNumber(t *testing.T) {
	pd := NewFirstEphemeralData()

	pn, ok := pd.GetPrevPulseNumber()
	require.False(t, ok)
	require.Equal(t, Number(MinTimePulse), pn)

	pd = newPulsarData(MinTimePulse+2, 2, longbits.Bits256{3})
	pd.PrevPulseDelta = 2
	pn, ok = pd.GetPrevPulseNumber()
	require.True(t, ok)
	require.Equal(t, Number(MinTimePulse), pn)

	pd = newPulsarData(MinTimePulse+1, 2, longbits.Bits256{3})
	pd.PrevPulseDelta = 2
	pn, ok = pd.GetPrevPulseNumber()
	require.False(t, ok)
	require.Equal(t, Number(MinTimePulse), pn)

	pd = newPulsarData(MinTimePulse, 2, longbits.Bits256{3})
	pd.PrevPulseDelta = 2
	pn, ok = pd.GetPrevPulseNumber()
	require.False(t, ok)
	require.Equal(t, Number(MinTimePulse), pn)
}

func TestNextPulseNumber(t *testing.T) {
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse, delta, entropy)
	pd.PulseNumber--
	require.Panics(t, func() { pd.NextPulseNumber() })

	pd.PulseNumber = MaxTimePulse + 1
	require.Panics(t, func() { pd.NextPulseNumber() })

	pd.PulseNumber = MaxTimePulse - 1
	require.Panics(t, func() { pd.NextPulseNumber() })

	pd.NextPulseDelta = 0
	require.Panics(t, func() { pd.NextPulseNumber() })

	pd.NextPulseDelta = delta
	pd.PulseNumber = MinTimePulse
	require.Equal(t, MinTimePulse+Number(delta), pd.NextPulseNumber())
}

func TestPrevPulseNumber(t *testing.T) {
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(MinTimePulse, delta, entropy)
	pd.PulseNumber--
	require.Panics(t, func() { pd.PrevPulseNumber() })

	pd.PulseNumber = MaxTimePulse + 1
	require.Panics(t, func() { pd.PrevPulseNumber() })

	pd.PulseNumber = MinTimePulse + 1
	require.Panics(t, func() { pd.PrevPulseNumber() })

	pd.PrevPulseDelta = 0
	require.Panics(t, func() { pd.PrevPulseNumber() })

	pd.PrevPulseDelta = delta
	pd.PulseNumber = MaxTimePulse
	require.Equal(t, MaxTimePulse-Number(delta), pd.PrevPulseNumber())
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

func TestCreateNextExpectedAndValidated(t *testing.T) {
	pn := Number(MinTimePulse)
	delta := uint16(2)
	entropy := longbits.Bits256{3}
	pd := newPulsarData(pn, delta, entropy)
	require.True(t, pd.CreateNextExpected().IsValidExpectedPulseData())
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

func Test_createNextEphemeralPulse(t *testing.T) {
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

func Test_createNextPulsarPulse(t *testing.T) {
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
	require.Equal(t, int64(5), pd.GetTimestamp())
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

func TestSort(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pd := pdLeft.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	sorted := []Data{pdBefore, pdLeft, pd, pdAfter}
	shuffled := append([]Data(nil), sorted...)

	require.Equal(t, sorted, shuffled)
	for {
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		if sorted[0] != shuffled[0] {
			break
		}
	}
	require.NotEqual(t, sorted, shuffled)
	SortData(shuffled)
	require.Equal(t, sorted, shuffled)
}

func TestAsPulse(t *testing.T) {
	pd := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	require.Equal(t, pd, pd.AsRange().RightBoundData())
}

func TestData_IsEmptyCompatibleWith(t *testing.T) {
	pd := Data{}

	require.False(t, pd.IsEmptyCompatibleWith(InvalidPulseEpoch))
	require.False(t, pd.IsEmptyCompatibleWith(EphemeralPulseEpoch))
	require.False(t, pd.IsEmptyCompatibleWith(ArticulationPulseEpoch))
	require.False(t, pd.IsEmptyCompatibleWith(MinTimePulse))

	pd.PulseEpoch = EphemeralPulseEpoch
	require.False(t, pd.IsEmptyCompatibleWith(InvalidPulseEpoch))
	require.True(t, pd.IsEmptyCompatibleWith(EphemeralPulseEpoch))
	require.False(t, pd.IsEmptyCompatibleWith(ArticulationPulseEpoch))
	require.False(t, pd.IsEmptyCompatibleWith(MinTimePulse))
}

func TestData_HasValidTimeEpoch(t *testing.T) {
	pd := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	require.True(t, pd.HasValidTimeEpoch())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch() + 1
	require.False(t, pd.HasValidTimeEpoch())

	pd.PulseEpoch = pd.PulseNumber.AsEpoch() - 1
	require.True(t, pd.HasValidTimeEpoch())

	pd.PulseEpoch = EphemeralPulseEpoch
	require.False(t, pd.HasValidTimeEpoch())
}
