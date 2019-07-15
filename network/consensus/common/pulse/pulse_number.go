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

import "time"

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

func (n Number) AsApproximateTime() time.Time {
	return timeOfMinTimePulse.Add(time.Second * time.Duration(n))
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

func (n Number) IsUnknown() bool {
	return n == Unknown
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

func IsValidAsPulseNumber(n int) bool {
	return n >= MinTimePulse && n <= MaxTimePulse
}

func OfInt(n int) Number {
	return Number(n) & MaxTimePulse
}

func OfUint32(n uint32) Number {
	return Number(n) & MaxTimePulse
}

func FlagsOf(n int) uint {
	return uint(n) >> 30
}
