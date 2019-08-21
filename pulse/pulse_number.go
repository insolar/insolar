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
	"errors"
	"time"
)

type Number uint32

const (
	Unknown       Number = 0
	LocalRelative        = 65536
	MinTimePulse         = LocalRelative + 1
	MaxTimePulse         = 1<<30 - 1
)
const UnixTimeOfMinTimePulse = 1546300800                                           // 2019-01-01 00:00:00 +0000 UTC
const UnixTimeOfMaxTimePulse = UnixTimeOfMinTimePulse - MinTimePulse + MaxTimePulse // 2053-01-08 19:24:46 +0000 UTC

var timeOfMinTimePulse = time.Unix(UnixTimeOfMinTimePulse, 0)

func OfNow() Number {
	return OfTime(time.Now())
}

func OfTime(t time.Time) Number {
	return OfUnixTime(t.Unix())
}

func OfUnixTime(u int64) Number {
	if u < UnixTimeOfMinTimePulse || u > UnixTimeOfMaxTimePulse {
		return Unknown
	}
	return MinTimePulse + Number(u-UnixTimeOfMinTimePulse)
}

func (n Number) AsApproximateTime() (time.Time, error) {
	if !n.IsTimePulse() {
		return timeOfMinTimePulse, errors.New("illegal state")
	}

	return timeOfMinTimePulse.Add(time.Second * time.Duration(n-MinTimePulse)), nil
}

func (n Number) IsTimePulse() bool {
	return IsValidAsPulseNumber(int(n))
}

func (n Number) AsUint32() uint32 {
	return uint32(n)
}

func (n Number) IsSpecialOrTimePulse() bool {
	return n > Unknown && n <= MaxTimePulse
}

func (n Number) IsSpecial() bool {
	return n > Unknown && n < MinTimePulse
}

func (n Number) IsUnknown() bool {
	return n == Unknown
}

func (n Number) IsUnknownOrTimePulse() bool {
	return n == Unknown || n >= MinTimePulse && n <= MaxTimePulse
}

func (n Number) IsUnknownOrEqualTo(o Number) bool {
	return n.IsUnknown() || n == o
}

func (n Number) Next(delta uint16) Number {
	if !n.IsTimePulse() {
		panic("not a time pulse")
	}
	n += Number(delta)
	if n > MaxTimePulse {
		panic("overflow")
	}
	return n
}

func (n Number) Prev(delta uint16) Number {
	if !n.IsTimePulse() {
		panic("not a time pulse")
	}
	n -= Number(delta)
	if n < MinTimePulse {
		panic("underflow")
	}
	return n
}

func (n Number) WithFlags(flags uint8) uint32 {
	if n > MaxTimePulse {
		panic("illegal value")
	}
	if flags > 3 {
		panic("illegal value")
	}
	return n.AsUint32() | uint32(flags)<<30
}

func IsValidAsPulseNumber(n int) bool {
	return n >= MinTimePulse && n <= MaxTimePulse
}

func OfInt(n int) Number {
	return Number(n) & MaxTimePulse
}

func OfUint32(n uint32) Number {
	return Number(n) & MaxTimePulse
}

func FlagsOf(n uint32) uint {
	return uint(n) >> 30
}
