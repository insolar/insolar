package pulse

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/longbits"
)

func TestNewLeftGapRange(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pd := pdLeft.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	rg := NewLeftGapRange(pd.PulseNumber, pd.PrevPulseDelta, pd)
	require.True(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdLeft}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	rg = NewLeftGapRange(pdLeft.PulseNumber, pdLeft.PrevPulseDelta, pd)
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	require.Panics(t, func() { NewLeftGapRange(pdAfter.PulseNumber, pdAfter.PrevPulseDelta, pd) })
}

func TestNewSequenceRange(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pd := pdLeft.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	require.Panics(t, func() { NewSequenceRange(nil) })

	rg := NewSequenceRange([]Data{pd})
	require.True(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdLeft}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	rg = NewSequenceRange([]Data{pdLeft, pd})
	require.False(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))
	require.IsType(t, seqPulseRange{}, rg)

	rg = NewSequenceRange([]Data{pdLeft, pd, pdAfter})
	require.False(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter.CreateNextPulse(emptyEntropyFn)}))
	require.IsType(t, seqPulseRange{}, rg)

	require.Panics(t, func() { NewSequenceRange([]Data{pdLeft, pdAfter}) })

	withExpected := []Data{
		{pdLeft.PulseNumber,
			DataExt{
				PrevPulseDelta: pdLeft.PrevPulseDelta,
			}},
		pd}

	require.Panics(t, func() { NewSequenceRange(withExpected) })
	withExpected[0].PulseEpoch = pdLeft.PulseEpoch
	require.Panics(t, func() { NewSequenceRange(withExpected) })
}

func TestNewPulseRange(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pd := pdLeft.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	require.Panics(t, func() { NewPulseRange(nil) })

	rg := NewPulseRange([]Data{pd})
	require.True(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdLeft}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	rg = NewPulseRange([]Data{pdLeft, pd})
	require.False(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))
	require.IsType(t, seqPulseRange{}, rg)

	rg = NewPulseRange([]Data{pdLeft, pd, pdAfter})
	require.False(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter.CreateNextPulse(emptyEntropyFn)}))
	require.IsType(t, seqPulseRange{}, rg)

	rg = NewPulseRange([]Data{pdLeft, pdAfter})
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter.CreateNextPulse(emptyEntropyFn)}))
	require.IsType(t, sparsePulseRange{}, rg)

	withExpected := []Data{
		{pdLeft.PulseNumber,
			DataExt{
				PrevPulseDelta: pdLeft.PrevPulseDelta,
			}},
		pd}

	require.Panics(t, func() { NewPulseRange(withExpected) })

	withExpected[0].PulseEpoch = pdLeft.PulseEpoch
	withExpected[0].PulseNumber++
	require.Panics(t, func() { NewPulseRange(withExpected) })
	withExpected[0].PulseNumber--

	rg = NewPulseRange(withExpected)
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))
	require.IsType(t, gapPulseRange{}, rg)

	rg = NewPulseRange(append(withExpected, pdAfter))
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter.CreateNextPulse(emptyEntropyFn)}))
	require.IsType(t, sparsePulseRange{}, rg)

	withExpected[1] = pdAfter
	rg = NewPulseRange(withExpected)
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())
	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter.CreateNextPulse(emptyEntropyFn)}))
	require.IsType(t, gapPulseRange{}, rg)
}

func Test_checkSequence(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pd := pdLeft.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	//pdBefore.NextPulseDelta = 0

	require.True(t, checkSequence([]Data{pdLeft}))
	require.True(t, checkSequence([]Data{pdLeft}))
	require.True(t, checkSequence([]Data{pdLeft, pd}))
	require.True(t, checkSequence([]Data{pdLeft, pd, pdAfter}))
	require.True(t, checkSequence([]Data{pd, pdAfter}))
	require.True(t, checkSequence([]Data{pdAfter}))
	require.True(t, checkSequence([]Data{pd}))

	require.True(t, checkSequence([]Data{pdBefore, pdLeft, pd, pdAfter}))
	require.False(t, checkSequence([]Data{pdBefore, pdLeft, pdAfter}))
	require.False(t, checkSequence([]Data{pdBefore, pdAfter}))

	require.False(t, checkSequence([]Data{pdLeft, pdAfter}))

	pdLeft.PulseNumber++
	require.Panics(t, func() { checkSequence([]Data{pdLeft, pdAfter}) })
	pdLeft.PulseNumber--

	pdAfter.PulseNumber++
	require.False(t, checkSequence([]Data{pdLeft, pdAfter}))
	require.Panics(t, func() { checkSequence([]Data{pd, pdAfter}) })
	pdAfter.PulseNumber--

	pd.PulseNumber++
	require.Panics(t, func() { checkSequence([]Data{pd, pdAfter}) })
	require.Panics(t, func() { checkSequence([]Data{pdLeft, pd}) })
}

var emptyEntropyFn = func() longbits.Bits256 {
	return longbits.Bits256{}
}

func Test_onePulseRange(t *testing.T) {
	pdLeft := NewPulsarData(MinTimePulse<<1, 10, 10, longbits.Bits256{})
	pd := pdLeft.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	rg := onePulseRange{pd}
	require.Equal(t, pd, rg.RightBoundData())
	require.Equal(t, pd.PulseNumber, rg.LeftBoundNumber())
	require.Equal(t, pd.PrevPulseDelta, rg.LeftPrevDelta())
	require.True(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())

	require.Equal(t, []Data{pd}, testRangeEnumNonArticulatedData(t, rg))
	require.Equal(t, []Number{pd.PulseNumber}, testRangeEnumSegments(t, rg))
	require.Equal(t, []Data{pd}, testRangeEnumData(t, rg))

	require.True(t, rg.IsValidPrev(onePulseRange{pdLeft}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	pdLeft.PulseNumber--
	require.False(t, rg.IsValidPrev(onePulseRange{pdLeft}))

	pdAfter.PulseNumber++
	require.False(t, rg.IsValidNext(onePulseRange{pdAfter}))
}

func Test_gapPulseRange(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pd := pdLeft.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	rg := gapPulseRange{start: pdLeft.PulseNumber, prevDelta: pdLeft.PrevPulseDelta, end: pd}

	require.Equal(t, pd, rg.RightBoundData())
	require.Equal(t, pdLeft.PulseNumber, rg.LeftBoundNumber())
	require.Equal(t, pdLeft.PrevPulseDelta, rg.LeftPrevDelta())
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())

	require.Equal(t, []Data{pd}, testRangeEnumNonArticulatedData(t, rg))

	require.Equal(t, []Number{pdLeft.PulseNumber, pd.PulseNumber},
		testRangeEnumSegments(t, rg))

	require.Equal(t, []Data{{pdLeft.PulseNumber,
		DataExt{
			PulseEpoch:     ArticulationPulseEpoch,
			PrevPulseDelta: pdLeft.PrevPulseDelta,
			NextPulseDelta: pdLeft.NextPulseDelta,
		},
	}, pd}, testRangeEnumData(t, rg))

	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	pdBefore.PulseNumber--
	require.False(t, rg.IsValidPrev(onePulseRange{pdBefore}))

	pdAfter.PulseNumber++
	require.False(t, rg.IsValidNext(onePulseRange{pdAfter}))
}

func Test_gapPulseRange_wideGap(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pdInterim := NewPulsarData(pdLeft.PulseNumber+math.MaxUint16, 10, 10, longbits.Bits256{})
	pd := pdInterim.CreateNextPulse(emptyEntropyFn)

	rg := gapPulseRange{start: pdLeft.PulseNumber, prevDelta: pdLeft.PrevPulseDelta, end: pd}

	require.Equal(t, pd, rg.RightBoundData())
	require.Equal(t, pdLeft.PulseNumber, rg.LeftBoundNumber())
	require.Equal(t, pdLeft.PrevPulseDelta, rg.LeftPrevDelta())
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())

	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))

	require.Equal(t, []Number{pdLeft.PulseNumber, pdInterim.PulseNumber, pd.PulseNumber},
		testRangeEnumSegments(t, rg))

	require.Equal(t, []Data{
		{pdLeft.PulseNumber,
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: pdLeft.PrevPulseDelta,
				NextPulseDelta: math.MaxUint16,
			}},
		{pdInterim.PulseNumber,
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: math.MaxUint16,
				NextPulseDelta: pd.PrevPulseDelta,
			}},
		pd}, testRangeEnumData(t, rg))

	pdLeft.PulseNumber--
	require.False(t, rg.IsValidPrev(onePulseRange{pdLeft}))
}

func Test_gapPulseRange_wideGap2(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pdInterim := NewPulsarData(pdLeft.PulseNumber+math.MaxUint16+1, 10, 10, longbits.Bits256{})
	pd := pdInterim.CreateNextPulse(emptyEntropyFn)

	rg := gapPulseRange{start: pdLeft.PulseNumber, prevDelta: pdLeft.PrevPulseDelta, end: pd}

	require.Equal(t, pd, rg.RightBoundData())
	require.Equal(t, pdLeft.PulseNumber, rg.LeftBoundNumber())
	require.Equal(t, pdLeft.PrevPulseDelta, rg.LeftPrevDelta())
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())

	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))

	require.Equal(t, []Number{pdLeft.PulseNumber, pdInterim.PulseNumber - minSegmentPulseDelta, pdInterim.PulseNumber, pd.PulseNumber},
		testRangeEnumSegments(t, rg))

	require.Equal(t, []Data{
		{pdLeft.PulseNumber,
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: pdLeft.PrevPulseDelta,
				NextPulseDelta: uint16(pdInterim.PulseNumber - pdLeft.PulseNumber - minSegmentPulseDelta),
			}},
		{pdInterim.PulseNumber - minSegmentPulseDelta,
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: uint16(pdInterim.PulseNumber - pdLeft.PulseNumber - minSegmentPulseDelta),
				NextPulseDelta: minSegmentPulseDelta,
			}},
		{pdInterim.PulseNumber,
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: minSegmentPulseDelta,
				NextPulseDelta: pd.PrevPulseDelta,
			}},
		pd}, testRangeEnumData(t, rg))

	pdLeft.PulseNumber--
	require.False(t, rg.IsValidPrev(onePulseRange{pdLeft}))
}

func Test_seqPulseRange(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.createNextPulsarPulse(11, emptyEntropyFn)
	pdInterim := pdLeft.createNextPulsarPulse(12, emptyEntropyFn)
	pd := pdInterim.createNextPulsarPulse(13, emptyEntropyFn)
	pdAfter := pd.createNextPulsarPulse(14, emptyEntropyFn)

	rg := seqPulseRange{data: []Data{pdLeft, pdInterim, pd}}

	require.Equal(t, pd, rg.RightBoundData())
	require.Equal(t, pdLeft.PulseNumber, rg.LeftBoundNumber())
	require.Equal(t, pdLeft.PrevPulseDelta, rg.LeftPrevDelta())
	require.False(t, rg.IsSingular())
	require.False(t, rg.IsArticulated())

	require.Equal(t, []Data{pdLeft, pdInterim, pd}, testRangeEnumNonArticulatedData(t, rg))
	require.Equal(t, []Number{pdLeft.PulseNumber, pdInterim.PulseNumber, pd.PulseNumber},
		testRangeEnumSegments(t, rg))

	require.Equal(t, []Data{pdLeft, pdInterim, pd}, testRangeEnumData(t, rg))

	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	pdBefore.PulseNumber--
	require.False(t, rg.IsValidPrev(onePulseRange{pdBefore}))

	pdAfter.PulseNumber++
	require.False(t, rg.IsValidNext(onePulseRange{pdAfter}))
}

func Test_sparsePulseRange(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pdInterim := NewPulsarData(pdLeft.PulseNumber+math.MaxUint16, 10, 10, longbits.Bits256{})
	pd := pdInterim.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	rg := sparsePulseRange{data: []Data{pdLeft, pdInterim, pd}}

	require.Equal(t, pd, rg.RightBoundData())
	require.Equal(t, pdLeft.PulseNumber, rg.LeftBoundNumber())
	require.Equal(t, pdLeft.PrevPulseDelta, rg.LeftPrevDelta())
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())

	require.Equal(t, []Data{pdLeft, pdInterim, pd}, testRangeEnumNonArticulatedData(t, rg))

	require.Equal(t, []Number{pdLeft.PulseNumber, pdLeft.NextPulseNumber(),
		pdInterim.PrevPulseNumber(), pdInterim.PulseNumber, pd.PulseNumber},
		testRangeEnumSegments(t, rg))

	require.Equal(t, []Data{pdLeft,
		{pdLeft.NextPulseNumber(),
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: pdLeft.NextPulseDelta,
				NextPulseDelta: math.MaxUint16 - pdLeft.NextPulseDelta - pdInterim.PrevPulseDelta,
			}},
		{pdInterim.PrevPulseNumber(),
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: math.MaxUint16 - pdLeft.NextPulseDelta - pdInterim.PrevPulseDelta,
				NextPulseDelta: pdInterim.PrevPulseDelta,
			}},
		pdInterim,
		pd}, testRangeEnumData(t, rg))

	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	pdBefore.PulseNumber--
	require.False(t, rg.IsValidPrev(onePulseRange{pdBefore}))

	pdAfter.PulseNumber++
	require.False(t, rg.IsValidNext(onePulseRange{pdAfter}))
}

func Test_sparsePulseRange_leftGap(t *testing.T) {
	pdBefore := NewPulsarData(MinTimePulse<<1, 10, 1, longbits.Bits256{})
	pdLeft := pdBefore.CreateNextPulse(emptyEntropyFn)
	pdInterim := NewPulsarData(pdLeft.PulseNumber+math.MaxUint16, 10, 10, longbits.Bits256{})
	pd := pdInterim.CreateNextPulse(emptyEntropyFn)
	pdAfter := pd.CreateNextPulse(emptyEntropyFn)

	rg := sparsePulseRange{data: []Data{{pdLeft.PulseNumber,
		DataExt{
			PrevPulseDelta: pdLeft.NextPulseDelta,
		},
	}, pdInterim, pd}}

	require.Equal(t, pd, rg.RightBoundData())
	require.Equal(t, pdLeft.PulseNumber, rg.LeftBoundNumber())
	require.Equal(t, pdLeft.PrevPulseDelta, rg.LeftPrevDelta())
	require.False(t, rg.IsSingular())
	require.True(t, rg.IsArticulated())

	require.Equal(t, []Data{pdInterim, pd}, testRangeEnumNonArticulatedData(t, rg))

	require.Equal(t, []Number{pdLeft.PulseNumber, pdInterim.PrevPulseNumber(), pdInterim.PulseNumber, pd.PulseNumber},
		testRangeEnumSegments(t, rg))

	require.Equal(t, []Data{
		{pdLeft.PulseNumber,
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: pdLeft.PrevPulseDelta,
				NextPulseDelta: math.MaxUint16 - pdInterim.PrevPulseDelta,
			}},
		{pdInterim.PrevPulseNumber(),
			DataExt{
				PulseEpoch:     ArticulationPulseEpoch,
				PrevPulseDelta: math.MaxUint16 - pdInterim.PrevPulseDelta,
				NextPulseDelta: pdInterim.PrevPulseDelta,
			}},
		pdInterim,
		pd}, testRangeEnumData(t, rg))

	require.True(t, rg.IsValidPrev(onePulseRange{pdBefore}))
	require.True(t, rg.IsValidNext(onePulseRange{pdAfter}))

	pdBefore.PulseNumber--
	require.False(t, rg.IsValidPrev(onePulseRange{pdBefore}))

	pdAfter.PulseNumber++
	require.False(t, rg.IsValidNext(onePulseRange{pdAfter}))
}

func Test_enumSegmentData(t *testing.T) {
	for prevDelta := uint16(1); prevDelta <= 20; prevDelta++ {
		for base := Number(MinTimePulse); base <= Number(MinTimePulse*3); base += MinTimePulse {
			for widthStep := 0; widthStep <= 3; widthStep++ {
				end := base + Number(widthStep*(math.MaxUint16)+int(prevDelta))
				t.Run(fmt.Sprint(prevDelta, base, end), func(t *testing.T) {
					endData := Data{end, DataExt{
						NextPulseDelta: 10,
						PrevPulseDelta: prevDelta,
						PulseEpoch:     end.AsEpoch(),
					}}

					{
						numbers := make([]Number, widthStep+1)
						for i := 0; i <= widthStep; i++ {
							numbers[i] = base + math.MaxUint16*Number(i)
						}
						require.Equal(t, d(0, numbers, endData), testEnumData(t, base, 0, endData))
						require.Equal(t, d(prevDelta, numbers, endData), testEnumData(t, base, prevDelta, endData))
					}

					if widthStep == 0 {
						require.Equal(t, d(prevDelta, nil, endData), testEnumData(t, endData.PulseNumber, prevDelta, endData))
						return
					}
					{ //step back doesn't cause change to a number of segments
						endData.PulseNumber--
						endData.PulseEpoch--

						numbers := make([]Number, widthStep+1)
						for i := 0; i <= widthStep; i++ {
							numbers[i] = base + math.MaxUint16*Number(i)
						}
						numbers[widthStep]--
						require.Equal(t, d(0, numbers, endData), testEnumData(t, base, 0, endData))
						require.Equal(t, d(prevDelta, numbers, endData), testEnumData(t, base, prevDelta, endData))
					}

					{ // step forward will result in a segment less than minSegmentPulseDelta
						// so one more segment is added with reduction of a previous segment
						endData.PulseNumber += 2
						endData.PulseEpoch += 2

						numbers := make([]Number, widthStep+2)
						for i := 0; i <= widthStep; i++ {
							numbers[i] = base + math.MaxUint16*Number(i)
						}
						numbers[widthStep+1] = numbers[widthStep] + 1
						numbers[widthStep] = numbers[widthStep+1] - minSegmentPulseDelta

						require.Equal(t, d(0, numbers, endData), testEnumData(t, base, 0, endData))
						require.Equal(t, d(prevDelta, numbers, endData), testEnumData(t, base, prevDelta, endData))
					}
				})
			}
		}
	}
}

func Test_enumSegments(t *testing.T) {
	for base := Number(MinTimePulse); base <= Number(MinTimePulse+20); base++ {
		for prevDelta := uint16(0); prevDelta <= uint16(base-MinTimePulse); prevDelta++ {
			for nextDelta := uint16(1); nextDelta <= uint16(base-MinTimePulse); nextDelta += 2 {
				require.Equal(t, n(base, base+1),
					testEnumSegments(t, base, prevDelta, base+1, nextDelta))

				require.Equal(t, n(base, base+math.MaxUint16),
					testEnumSegments(t, base, prevDelta, base+math.MaxUint16, nextDelta))

				require.Equal(t, n(base, base+math.MaxUint16-minSegmentPulseDelta+1, base+math.MaxUint16+1),
					testEnumSegments(t, base, prevDelta, base+math.MaxUint16+1, nextDelta))

				require.Equal(t, n(base, base+math.MaxUint16-1, base+math.MaxUint16+minSegmentPulseDelta-1),
					testEnumSegments(t, base, prevDelta, base+math.MaxUint16+minSegmentPulseDelta-1, nextDelta))

				require.Equal(t, n(base, base+math.MaxUint16, base+math.MaxUint16+minSegmentPulseDelta),
					testEnumSegments(t, base, prevDelta, base+math.MaxUint16+minSegmentPulseDelta, nextDelta))

				require.Equal(t, n(base, base+math.MaxUint16, base+math.MaxUint16*2),
					testEnumSegments(t, base, prevDelta, base+math.MaxUint16*2, nextDelta))

				require.Equal(t, n(base, base+math.MaxUint16, base+math.MaxUint16*2-minSegmentPulseDelta+1, base+math.MaxUint16*2+1),
					testEnumSegments(t, base, prevDelta, base+math.MaxUint16*2+1, nextDelta))
			}
		}
	}
}

func testEnumData(t *testing.T, next Number, prevDeltaOuter uint16, end Data) []Data {
	var nums []Data
	lastNextDelta := prevDeltaOuter
	lastNextNum := next

	_enumSegmentData(next, prevDeltaOuter, end, func(d Data) bool {
		require.Equal(t, lastNextNum, d.PulseNumber)
		require.Equal(t, lastNextDelta, d.PrevPulseDelta)
		require.Greater(t, d.NextPulseDelta, uint16(0))
		nums = append(nums, d)

		lastNextDelta = d.NextPulseDelta
		lastNextNum = d.NextPulseNumber()

		return false
	})

	require.Equal(t, lastNextNum, end.NextPulseNumber())
	require.Equal(t, lastNextDelta, end.NextPulseDelta)

	return nums
}

func testEnumSegments(t *testing.T, next Number, prevDeltaOuter uint16, end Number, endNextDelta uint16) []Number {
	var nums []Number
	lastNextDelta := prevDeltaOuter
	lastNextNum := next

	_enumSegments(next, prevDeltaOuter, end, endNextDelta, func(n Number, prevDelta, nextDelta uint16) bool {
		require.Equal(t, lastNextNum, n)
		require.Equal(t, lastNextDelta, prevDelta)
		require.Greater(t, nextDelta, uint16(0))
		nums = append(nums, n)

		lastNextDelta = nextDelta
		lastNextNum = n.Next(nextDelta)

		return false
	})

	require.Equal(t, lastNextNum, end.Next(endNextDelta))
	require.Equal(t, lastNextDelta, endNextDelta)

	return nums
}

func testRangeEnumData(t *testing.T, rg Range) []Data {
	var nums []Data
	lastNextDelta := rg.LeftPrevDelta()
	lastNextNum := rg.LeftBoundNumber()

	rg.EnumData(func(d Data) bool {
		require.Equal(t, lastNextNum, d.PulseNumber)
		require.Equal(t, lastNextDelta, d.PrevPulseDelta)
		require.Greater(t, d.NextPulseDelta, uint16(0))
		nums = append(nums, d)

		lastNextDelta = d.NextPulseDelta
		lastNextNum = d.NextPulseNumber()

		return false
	})

	end := rg.RightBoundData()
	require.Equal(t, lastNextNum, end.NextPulseNumber())
	require.Equal(t, lastNextDelta, end.NextPulseDelta)

	return nums
}

func testRangeEnumNonArticulatedData(_ *testing.T, rg Range) []Data {
	var nums []Data
	rg.EnumNonArticulatedData(func(d Data) bool {
		nums = append(nums, d)
		return false
	})
	return nums
}

func testRangeEnumSegments(t *testing.T, rg Range) []Number {
	var nums []Number
	lastNextDelta := rg.LeftPrevDelta()
	lastNextNum := rg.LeftBoundNumber()

	rg.EnumNumbers(func(n Number, prevDelta, nextDelta uint16) bool {
		require.Equal(t, lastNextNum, n)
		require.Equal(t, lastNextDelta, prevDelta)
		require.Greater(t, nextDelta, uint16(0))
		nums = append(nums, n)

		lastNextDelta = nextDelta
		lastNextNum = n.Next(nextDelta)

		return false
	})

	end := rg.RightBoundData()
	require.Equal(t, lastNextNum, end.NextPulseNumber())
	require.Equal(t, lastNextDelta, end.NextPulseDelta)

	return nums
}

//type numCollector struct {
//	nums []Number
//}
//
//func (p *numCollector) collect(n Number) {
//	p.nums = append(p.nums, n)
//}

func n(numbers ...Number) []Number {
	return numbers
}

func d(prevDelta uint16, numbers []Number, end Data) []Data {
	data := make([]Data, 0, len(numbers)+1)
	for _, n := range numbers {
		cp := end
		cp.PulseNumber = n
		cp.PulseEpoch = ArticulationPulseEpoch
		data = append(data, cp)
	}
	data = append(data, end)
	data[0].PrevPulseDelta = prevDelta

	for i := 0; i < len(data)-1; i++ {
		delta := data[i+1].PulseNumber - data[i].PulseNumber
		data[i].NextPulseDelta = uint16(delta)
		data[i+1].PrevPulseDelta = uint16(delta)
	}

	return data
}
