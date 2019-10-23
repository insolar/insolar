///
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
///

package pulse

import "math"

//var _ DataReader = Range{}

type Range struct {
	start          Number
	startPrevDelta uint16

	end   Data
	epoch *Data
}

func (v Range) IsZero() bool {
	return v.start.IsUnknown() && v.start == v.end.PulseNumber
}

func (v Range) IsArticulated() bool {
	return v.start < v.end.PulseNumber
}

//func (v Range) ListNumbers(fn func (n Number, prevDelta uint16)) {
//	if !v.IsArticulated() {
//		fn(v.end.PulseNumber)
//	}
//}

func (v Range) BuildPulseChain() []Data {
	if !v.IsArticulated() {
		return []Data{v.end}
	}

	switch {
	case v.epoch == nil || v.start > v.epoch.PulseNumber: // epoch is out of range
		return _appendSegments(nil, v.start, v.startPrevDelta, v.end)

	case v.start == v.epoch.PulseNumber: // epoch is the start of the range
		// TODO check [epoch, end] case
		chain := append(make([]Data, 0, 2), *v.epoch)
		return _appendSegments(chain, v.epoch.NextPulseNumber(), v.epoch.NextPulseDelta, v.end)

	case v.epoch.PulseNumber >= v.end.PulseNumber:
		panic("illegal state")

	default: // epoch is inside this range
		chain := make([]Data, 0, 3)
		chain = _appendSegments(nil, v.start, v.startPrevDelta, *v.epoch)
		return _appendSegments(chain, v.epoch.NextPulseNumber(), v.epoch.NextPulseDelta, v.end)
	}
}

const minSegmentPulseDelta = 10

func _appendSegments(chain []Data, next Number, prevDelta uint16, end Data) []Data {
	if next != end.PulseNumber {
		_enumSegments(next, prevDelta, end.PrevPulseNumber(), end.PrevPulseDelta, func(n Number, prevDelta, nextDelta uint16) {
			chain = append(chain, Data{
				next, DataExt{
					PulseEpoch:     ArticulationPulseEpoch,
					NextPulseDelta: nextDelta,
					PrevPulseDelta: prevDelta,
				}})
		})
	}
	chain = append(chain, end)
	return chain
}

func _enumSegments(next Number, prevDelta uint16, beforeEnd Number, endNextDelta uint16,
	fn func(n Number, prevDelta, nextDelta uint16),
) {
	for {
		switch {
		case next < beforeEnd:
			delta := beforeEnd - next
			switch {
			case delta <= math.MaxUint16:
			case delta < math.MaxUint16+minSegmentPulseDelta:
				delta -= minSegmentPulseDelta
			default:
				delta = math.MaxUint16
			}
			fn(next, prevDelta, uint16(delta))
			prevDelta = uint16(delta)
			next = next.Next(prevDelta)
			continue
		case next == beforeEnd:
			fn(next, prevDelta, endNextDelta)
			return
		default:
			panic("illegal state")
		}
	}
}
