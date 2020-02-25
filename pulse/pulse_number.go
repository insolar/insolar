// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"encoding/binary"
	"errors"
	"strconv"
	"time"
)

// Number is a type for pulse numbers.
//
// Special values:
// 0 					Unknown
// 1 .. 256				Reserved for package internal usage
// 257 .. 65535			Reserved for platform wide usage
// 65536				Local relative pulse number
// 65537 .. 1<<30 - 1	Regular time based pulse numbers
//
// NB! Range 0..256 IS RESERVED for internal operations
// There MUST BE NO references with PN < 256 ever visible to contracts / users.
type Number uint32

// =========================================================
// NB! To ADD a special pulse - see special_pulse_numbers.go
// =========================================================
const (
	Unknown       Number = 0
	localRelative        = 65536
	LocalRelative Number = localRelative

	// MinTimePulse is the hardcoded first pulse number. Because first 65536 numbers are saved for the system's needs
	MinTimePulse = localRelative + 1
	MaxTimePulse = 1<<30 - 1

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

func (n Number) IsBefore(pn Number) bool {
	return n >= MinTimePulse && n < pn
}

func (n Number) IsAfter(pn Number) bool {
	return n > pn && n <= MaxTimePulse
}

func (n Number) IsBeforeOrEq(pn Number) bool {
	return n >= MinTimePulse && n <= pn
}

func (n Number) IsEqOrAfter(pn Number) bool {
	return n >= pn && n <= MaxTimePulse
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
	switch n, ok := n.TryNext(delta); {
	case ok:
		return n
	case n.IsUnknown():
		panic("not a time pulse")
	default:
		panic("overflow")
	}
}

func (n Number) TryNext(delta uint16) (Number, bool) {
	if !n.IsTimePulse() {
		return Unknown, false
	}
	n += Number(delta)
	if n > MaxTimePulse {
		return MaxTimePulse, false
	}
	return n, true
}

func (n Number) Prev(delta uint16) Number {
	switch n, ok := n.TryPrev(delta); {
	case ok:
		return n
	case n.IsUnknown():
		panic("not a time pulse")
	default:
		panic("underflow")
	}
}

func (n Number) TryPrev(delta uint16) (Number, bool) {
	if !n.IsTimePulse() {
		return Unknown, false
	}
	n -= Number(delta)
	if n < MinTimePulse {
		return MinTimePulse, false
	}
	return n, true
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

func (n Number) EnsureTimePulse() Number {
	if n.IsTimePulse() {
		return n
	}
	panic("illegal value")
}

func (n Number) AsEpoch() Epoch {
	if n.IsTimePulse() {
		return Epoch(n)
	}
	panic("illegal value")
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
