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
	"encoding/binary"
	"errors"
	"strconv"
	"time"
)

type Number uint32

const (
	Unknown       Number = 0
	LocalRelative        = 65536
	// MinTimePulse is the hardcoded first pulse number. Because first 65536 numbers are saved for the system's needs
	MinTimePulse = LocalRelative + 1
	MaxTimePulse = 1<<30 - 1
	// Jet is a special pulse number value that signifies jet ID.
	Jet Number = 1
	// BuiltinContract declares special pulse number that creates namespace for builtin contracts
	BuiltinContract Number = 200
	// PulseNumberSize declares the number of bytes in the pulse number
	NumberSize int = 4
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

// Bytes serializes pulse number.
func (n Number) Bytes() []byte {
	var buf [NumberSize]byte
	binary.BigEndian.PutUint32(buf[:], uint32(n))
	return buf[:]
}

func (n Number) String() string {
	return strconv.FormatUint(uint64(n), 10)
}

func (n Number) MarshalTo(data []byte) (int, error) {
	if len(data) < NumberSize {
		return 0, errors.New("not enough bytes to marshal pulse.Number")
	}
	binary.BigEndian.PutUint32(data, uint32(n))
	return NumberSize, nil
}

func (n *Number) Unmarshal(data []byte) error {
	if len(data) < NumberSize {
		return errors.New("not enough bytes to unmarshal pulse.Number")
	}
	*n = Number(binary.BigEndian.Uint32(data))
	return nil
}

func (n Number) Equal(other Number) bool {
	return n == other
}

func (n Number) Size() int {
	return NumberSize
}

func (n Number) IsJet() bool {
	return n == Jet
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
