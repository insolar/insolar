//
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
//

package member

import (
	"math/bits"

	"github.com/insolar/insolar/network/consensus/common/capacity"
)

type Power uint8

const MaxLinearMemberPower = (0x1F+32)<<(0xFF>>5) - 32

func PowerOf(linearValue uint16) Power { // TODO tests are needed
	if linearValue <= 0x1F {
		return Power(linearValue)
	}
	if linearValue >= MaxLinearMemberPower {
		return 0xFF
	}

	linearValue += 32
	pwr := uint8(bits.Len16(linearValue))
	if pwr > 6 {
		pwr -= 6
		linearValue >>= pwr
	} else {
		pwr = 0
	}
	return Power((pwr << 5) | uint8(linearValue-32))
}

func (v Power) ToLinearValue() uint16 {
	if v <= 0x1F {
		return uint16(v)
	}
	return uint16(v&0x1F+32)<<(v>>5) - 32
}

func (v Power) PercentAndMin(percent int, min Power) Power {
	vv := (int(v.ToLinearValue()) * percent) / 100
	if vv >= MaxLinearMemberPower {
		return ^Power(0)
	}
	if vv <= int(min.ToLinearValue()) {
		return min
	}
	return PowerOf(uint16(vv))
}

/*
	PowerSet enables power control by both discreet values or ranges.
	Zero level is always allowed by default
		PowerLevels[0] - min power value, must be <= PowerLevels[3], node is not allowed to set power lower than this value, except for zero power
		PowerLevels[3] - max power value, node is not allowed to set power higher than this value

	To define only distinct value, all values must be >0, e.g. (p1 = PowerLevels[1], p2 = PowerLevels[2]):
		[10, 20, 30, 40] - a node can only choose of: 0, 10, 20, 30, 40
		[10, 10, 30, 40] - a node can only choose of: 0, 10, 30, 40
		[10, 20, 20, 40] - a node can only choose of: 0, 10, 20, 40
		[10, 20, 20, 20] - a node can only choose of: 0, 10, 20
		[10, 10, 10, 10] - a node can only choose of: 0, 10

	Presence of 0 values treats nearest non-zero value as range boundaries, e.g.
		[ 0, 20, 30, 40] - a node can choose of: [0..20], 30, 40
		[10,  0, 30, 40] - a node can choose of: 0, [10..30], 40
		[10, 20,  0, 40] - a node can choose of: 0, 10, [20..40]
		[10,  0,  0, 40] - a node can choose of: 0, [10..40] ??? should be a special case?
		[ 0,  0,  0, 40] - a node can choose of: [0..40] ??? should be a special case?

	Special case:
		[ 0,  0,  0,  0] - a node can only use: 0

	Illegal cases:
		[ x,  y,  z,  0] - when any !=0 value of x, y, z
		[ 0,  x,  0,  y] - when x != 0 and y != 0
	    any combination of non-zero x, y such that x > y and y > 0 and position(x) < position(y)
*/

type PowerSet [4]Power

func (v PowerSet) Normalize() PowerSet {
	if v.IsValid() {
		return v
	}
	return [...]Power{0, 0, 0, 0}
}

func (v PowerSet) IsValid() bool {
	if v[3] == 0 {
		return v[0] == 0 && v[1] == 0 && v[2] == 0
	}

	if v[2] == 0 {
		if v[0] == 0 {
			return v[1] == 0
		}
		if v[1] == 0 {
			return v[0] <= v[3]
		}
		return v[0] <= v[1] && v[1] <= v[3]
	}

	if v[2] > v[3] {
		return false
	}
	if v[1] == 0 {
		return v[0] <= v[2]
	}

	return v[0] <= v[1] && v[1] <= v[2]
}

/*
Always true for p=0. Requires normalized ops.
*/func (v PowerSet) IsAllowed(p Power) bool {
	if p == 0 || v[0] == p || v[1] == p || v[2] == p || v[3] == p {
		return true
	}
	if v[0] > p || v[3] < p {
		return false
	}

	if v[2] == 0 { // [min, ?, 0, max]
		if v[0] == 0 || v[1] == 0 {
			return true
		} // [0, ?0, 0, max] or [min, 0, 0, max]

		// [min, p1, 0, max]
		return v[1] <= p
	}

	if v[1] == 0 { // [?, 0, p2, max]
		if v[0] == 0 { // [0, 0, p2, max]
			return p <= v[2] || p == v[3]
		}
		// [min, 0, p2, max]
		return v[3] == p || v[2] >= p
	}
	// [min, p1, p2, max]
	return false
}

/*
Only for normalized
*/func (v PowerSet) IsEmpty() bool {
	return v[0] == 0 && v[3] == 0
}

/*
Only for normalized
*/func (v PowerSet) Max() Power {
	return v[3]
}

/*
Only for normalized
*/func (v PowerSet) Min() Power {
	return v[0]
}

/*
Only for normalized
*/

func (v PowerSet) ForLevel(lvl capacity.Level) Power {
	return v.ForLevelWithPercents(lvl, 20, 60, 80)
}

/*
Only for normalized
*/

func (v PowerSet) ForLevelWithPercents(lvl capacity.Level, pMinimal, pReduced, pNormal int) Power {

	if lvl == capacity.LevelZero || v.IsEmpty() {
		return 0
	}

	switch lvl {
	case capacity.LevelMinimal:
		if v[0] != 0 {
			return v[0]
		}
		vv := v.Max().PercentAndMin(pMinimal, 1)

		if v[1] != 0 {
			if vv >= v[1] {
				return v[1]
			}
			return vv
		}
		if v[2] != 0 && vv >= v[2] {
			return v[2]
		}
		return vv
	case capacity.LevelReduced:
		if v[1] != 0 {
			return v[1]
		}
		vv := v.Max().PercentAndMin(pReduced, 1)

		if v[2] != 0 && vv >= v[2] {
			return v[2]
		}
		if v[0] != 0 && vv <= v[0] {
			return v[0]
		}
		return vv
	case capacity.LevelNormal:
		if v[2] != 0 {
			return v[2]
		}
		vv := v.Max().PercentAndMin(pNormal, 1)

		if v[1] != 0 {
			if vv >= v[1] {
				return vv
			}
			return v[1]
		}
		if v[0] != 0 && vv <= v[0] {
			return v[0]
		}
		return vv
	case capacity.LevelMax:
		return v[3]
	default:
		panic("missing")
	}
}

func (v PowerSet) FindNearestValid(pw Power) Power {

	panic("unsupported")
	if pw == 0 {
		return 0
	}
	if pw >= v[3] {
		return v[3]
	}
}
