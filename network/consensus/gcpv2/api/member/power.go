// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package member

import (
	"math"
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

	linearValue += 0x20
	pwr := uint8(bits.Len16(linearValue))
	pwr -= 6 // linearValue is always >= 0x40, so pwr > 6
	linearValue >>= pwr
	return Power((pwr << 5) | uint8(linearValue-0x20))
}

func (v Power) ToLinearValue() uint16 {
	if v <= 0x1F {
		return uint16(v)
	}
	return uint16(v&0x1F+0x20)<<(v>>5) - 0x20
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

func (v Power) Delta(o Power) uint16 {
	if v == o {
		return 0
	}
	return v.ToLinearValue() - o.ToLinearValue()
}

func (v Power) AbsDelta(o Power) uint16 {
	if v == o {
		return 0
	}
	if v > o {
		return v.ToLinearValue() - o.ToLinearValue()
	}
	return o.ToLinearValue() - v.ToLinearValue()
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

func PowerSetOf(v uint32) PowerSet {
	return PowerSet{Power(v), Power(v >> 8), Power(v >> 16), Power(v >> 24)}
}

func (v PowerSet) AsUint32() uint32 {
	return uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24
}

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

	// [min, p1, p2, max] - was tested at entry
	// [0, p1, p2, max]

	return v[0] == 0 && v[1] > p
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

/*
Chooses the nearest allowed value. Prefers higher values. Returns zero only for a zero value or for a zero range.
Only for normalized set.
*/

func (v PowerSet) FindNearestValid(p Power) Power {

	left := int8(0)
	right := int8(3)
	switch {
	case p == 0 || v[0] == p || v[1] == p || v[2] == p || v[3] == p:
		return p
	case p >= v[3]:
		return v[3]
	case p <= v[0]:
		return v[0]
	case v[2] == 0: // [min, ?, 0, max]
		if v[0] == 0 || v[1] == 0 { // [0, _, 0, max] or [min, 0, 0, max]
			return p
		}
		// [v0, min, 0, max]
		if v[1] < p {
			return p
		}
		// get nearest of v[0], v[1]
		right = 1
	case v[1] == 0: // [?, 0, p2, max]
		if v[2] < p {
			// get nearest of v[2], v[3]
			left = 2
			break
		}
		if v[0] < p { // [min, 0, p2, _]
			return p
		}
		return v[0]
	case v[0] == 0: // [0, p1, p2, p3]
		if v[1] > p { // [0, p1, p2, p3]
			return p
		}
		// get nearest of v[1], v[2], v[3]
		left = 1
	default: // [p0, p1, p2, p3]
	}

	var nearest Power
	for delta := uint16(math.MaxUint16); right >= left; right-- {
		next := v[right]
		nextDelta := p.AbsDelta(next)
		if delta > nextDelta {
			delta = nextDelta
			nearest = next
		} else {
			return nearest
		}
	}
	if nearest == 0 {
		panic("impossible")
	}
	return nearest
}
