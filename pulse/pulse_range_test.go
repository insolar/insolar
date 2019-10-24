//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package pulse

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

//func TestNewLeftGapRange(t *testing.T) {
//	type args struct {
//		left          Number
//		leftPrevDelta uint16
//		right         Data
//	}
//	tests := []struct {
//		name string
//		args args
//		want Range
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewLeftGapRange(tt.args.left, tt.args.leftPrevDelta, tt.args.right); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewLeftGapRange() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNewMultiPulseRange(t *testing.T) {
//	type args struct {
//		data []Data
//	}
//	tests := []struct {
//		name string
//		args args
//		want Range
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewMultiPulseRange(tt.args.data); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewMultiPulseRange() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_checkSequence(t *testing.T) {
//	type args struct {
//		data []Data
//	}
//	tests := []struct {
//		name string
//		args args
//		want bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := checkSequence(tt.args.data); got != tt.want {
//				t.Errorf("checkSequence() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func Test_onePulseRange(t *testing.T) {

}

func Test_gapPulseRange(t *testing.T) {
}

func Test_seqPulseRange(t *testing.T) {

}

func Test_sparsePulseRange(t *testing.T) {
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
