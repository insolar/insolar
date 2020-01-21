// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package member

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/common/capacity"

	"github.com/stretchr/testify/require"
)

func TestPowerOf(t *testing.T) {
	require.Equal(t, Power(1), PowerOf(1))

	require.Equal(t, Power(0x1F), PowerOf(0x1F))

	require.Equal(t, Power(0xFF), PowerOf(MaxLinearMemberPower))

	require.Equal(t, Power(0xFF), PowerOf(MaxLinearMemberPower+1))

	require.Equal(t, Power(0x1F+1), PowerOf(0x1F+1))

	require.Equal(t, Power(0x2F), PowerOf(0x1F<<1))
}

func TestToLinearValue(t *testing.T) {
	require.Equal(t, uint16(0), PowerOf(0).ToLinearValue())

	require.Equal(t, uint16(0x1F), PowerOf(0x1F).ToLinearValue())

	require.Equal(t, uint16(0x1F+1), PowerOf(0x1F+1).ToLinearValue())

	require.Equal(t, uint16(0x3e), PowerOf(0x1F<<1).ToLinearValue())
}

func TestPercentAndMin(t *testing.T) {
	require.Equal(t, ^Power(0), PowerOf(MaxLinearMemberPower).PercentAndMin(100, PowerOf(0)))

	require.Equal(t, Power(2), PowerOf(3).PercentAndMin(1, PowerOf(2)))

	require.Equal(t, Power(2), PowerOf(3).PercentAndMin(80, PowerOf(1)))
}

func TestDelta(t *testing.T) {
	require.Zero(t, Power(1).Delta(Power(1)))

	require.Equal(t, uint16(1), Power(2).Delta(Power(1)))
}

func TestAbsDelta(t *testing.T) {
	require.Zero(t, Power(1).AbsDelta(Power(1)))

	require.Equal(t, uint16(1), Power(2).AbsDelta(Power(1)))

	require.Equal(t, uint16(1), Power(1).AbsDelta(Power(2)))
}

func TestPowerSetOf(t *testing.T) {
	require.Equal(t, PowerSet([...]Power{0, 0, 0, 0}), PowerSetOf(0))

	require.Equal(t, PowerSet([...]Power{1, 2, 3, 4}), PowerSetOf(67305985))
}

func TestPowerSetAsUint32(t *testing.T) {
	require.Zero(t, PowerSet([...]Power{0, 0, 0, 0}).AsUint32())

	require.Equal(t, uint32(67305985), PowerSet([...]Power{1, 2, 3, 4}).AsUint32())
}

func TestNormalize(t *testing.T) {
	zero := PowerSet([...]Power{0, 0, 0, 0})
	require.Equal(t, zero, zero.Normalize())

	require.Equal(t, zero, PowerSet([...]Power{1, 0, 0, 0}).Normalize())

	m := PowerSet([...]Power{1, 1, 1, 1})
	require.Equal(t, m, m.Normalize())
}

// Illegal cases:
// [ x,  y,  z,  0] - when any !=0 value of x, y, z
// [ 0,  x,  0,  y] - when x != 0 and y != 0
// any combination of non-zero x, y such that x > y and y > 0 and position(x) < position(y)
// And cases from the function logic.
func TestIsValid(t *testing.T) {
	require.True(t, PowerSet([...]Power{0, 0, 0, 0}).IsValid())

	require.False(t, PowerSet([...]Power{1, 0, 0, 0}).IsValid())

	require.False(t, PowerSet([...]Power{0, 1, 0, 0}).IsValid())

	require.False(t, PowerSet([...]Power{0, 0, 1, 0}).IsValid())

	require.False(t, PowerSet([...]Power{0, 1, 0, 1}).IsValid())

	require.False(t, PowerSet([...]Power{2, 1, 2, 2}).IsValid())

	require.True(t, PowerSet([...]Power{1, 0, 0, 1}).IsValid())

	require.True(t, PowerSet([...]Power{1, 1, 0, 1}).IsValid())

	require.False(t, PowerSet([...]Power{1, 1, 2, 1}).IsValid())

	require.True(t, PowerSet([...]Power{1, 0, 2, 2}).IsValid())

	require.True(t, PowerSet([...]Power{1, 1, 2, 2}).IsValid())
}

func TestIsAllowed(t *testing.T) {
	ps := PowerSet([...]Power{10, 20, 30, 40})
	require.True(t, ps.IsAllowed(0))

	require.True(t, ps.IsAllowed(10))

	require.True(t, ps.IsAllowed(20))

	require.True(t, ps.IsAllowed(30))

	require.True(t, ps.IsAllowed(40))

	require.False(t, ps.IsAllowed(1))

	require.False(t, ps.IsAllowed(41))

	require.False(t, ps.IsAllowed(31))

	require.True(t, PowerSet([...]Power{0, 20, 0, 40}).IsAllowed(1))

	require.True(t, PowerSet([...]Power{10, 0, 0, 40}).IsAllowed(11))

	require.True(t, PowerSet([...]Power{10, 20, 0, 40}).IsAllowed(21))

	require.False(t, PowerSet([...]Power{10, 20, 0, 40}).IsAllowed(19))

	require.True(t, PowerSet([...]Power{0, 0, 30, 40}).IsAllowed(29))

	require.True(t, PowerSet([...]Power{0, 0, 30, 40}).IsAllowed(40))

	require.False(t, PowerSet([...]Power{0, 0, 30, 40}).IsAllowed(31))

	require.True(t, PowerSet([...]Power{10, 0, 30, 40}).IsAllowed(29))

	require.True(t, PowerSet([...]Power{10, 0, 30, 40}).IsAllowed(40))

	require.True(t, PowerSet([...]Power{0, 20, 30, 40}).IsAllowed(1))

	require.False(t, PowerSet([...]Power{0, 20, 30, 40}).IsAllowed(21))
}

func TestIsEmpty(t *testing.T) {
	require.True(t, PowerSet([...]Power{0, 0, 0, 0}).IsEmpty())

	require.False(t, PowerSet([...]Power{1, 0, 0, 1}).IsEmpty())
}

func TestMax(t *testing.T) {
	require.Equal(t, Power(1), PowerSet([...]Power{0, 0, 0, 1}).Max())
}

func TestMin(t *testing.T) {
	require.Equal(t, Power(1), PowerSet([...]Power{1, 0, 0, 0}).Min())
}

func TestForLevel(t *testing.T) {
	ps := PowerSet([...]Power{10, 20, 30, 40})
	require.Zero(t, ps.ForLevel(capacity.LevelZero))

	require.Equal(t, ps[0], ps.ForLevel(capacity.LevelMinimal))

	require.Equal(t, ps[1], ps.ForLevel(capacity.LevelReduced))

	require.Equal(t, ps[2], ps.ForLevel(capacity.LevelNormal))

	require.Equal(t, ps[3], ps.ForLevel(capacity.LevelMax))

	require.Panics(t, func() { ps.ForLevel(capacity.LevelMax + 1) })
}

func TestForLevelWithPercents(t *testing.T) {
	psBase := PowerSet([...]Power{10, 20, 30, 40})
	require.Zero(t, psBase.ForLevelWithPercents(capacity.LevelZero, 0, 0, 0))

	require.Zero(t, PowerSet([...]Power{0, 0, 0, 0}).ForLevelWithPercents(capacity.LevelMinimal,
		0, 0, 0))

	require.Equal(t, psBase[0], psBase.ForLevelWithPercents(capacity.LevelMinimal, 0, 0, 0))

	ps := psBase
	level := capacity.LevelMinimal
	ps[0] = 0
	require.Equal(t, ps[1], ps.ForLevelWithPercents(level, 100, 0, 0))

	require.Equal(t, Power(1), ps.ForLevelWithPercents(level, 1, 0, 0))

	ps[1] = 0
	require.Equal(t, ps[2], ps.ForLevelWithPercents(level, 100, 0, 0))

	require.Equal(t, Power(1), ps.ForLevelWithPercents(level, 1, 0, 0))

	ps[2] = 0
	require.Equal(t, Power(1), ps.ForLevelWithPercents(level, 1, 0, 0))

	ps = psBase
	level = capacity.LevelReduced
	require.Equal(t, ps[1], ps.ForLevelWithPercents(level, 0, 0, 0))

	ps[1] = 0
	require.Equal(t, ps[2], ps.ForLevelWithPercents(level, 0, 100, 0))

	require.Equal(t, ps[0], ps.ForLevelWithPercents(level, 0, 1, 0))

	ps[2] = 0
	require.Equal(t, ps[3], ps.ForLevelWithPercents(level, 0, 100, 0))

	ps[0] = 0
	require.Equal(t, Power(1), ps.ForLevelWithPercents(level, 0, 1, 0))

	ps = psBase
	level = capacity.LevelNormal
	require.Equal(t, ps[2], ps.ForLevelWithPercents(level, 0, 0, 0))

	ps[2] = 0
	require.Equal(t, ps[3], ps.ForLevelWithPercents(level, 0, 0, 100))

	require.Equal(t, ps[1], ps.ForLevelWithPercents(level, 0, 0, 1))

	ps[1] = 0
	require.Equal(t, ps[0], ps.ForLevelWithPercents(level, 0, 0, 1))

	require.Equal(t, ps[3], ps.ForLevelWithPercents(level, 0, 0, 100))

	ps[0] = 0
	require.Equal(t, Power(1), ps.ForLevelWithPercents(level, 0, 0, 1))

	level = capacity.LevelMax
	require.Equal(t, ps[3], ps.ForLevelWithPercents(level, 0, 0, 0))

	level = capacity.LevelMax + 1
	require.Panics(t, func() { ps.ForLevelWithPercents(level, 0, 0, 0) })
}

func TestFindNearestValid(t *testing.T) {
	psBase := PowerSet([...]Power{10, 20, 30, 40})
	ps := psBase
	require.Zero(t, ps.FindNearestValid(0))

	require.Equal(t, ps[0], ps.FindNearestValid(ps[0]))

	require.Equal(t, ps[1], ps.FindNearestValid(ps[1]))

	require.Equal(t, ps[2], ps.FindNearestValid(ps[2]))

	require.Equal(t, ps[3], ps.FindNearestValid(ps[3]))

	require.Equal(t, ps[3], ps.FindNearestValid(ps[3]+1))

	require.Equal(t, ps[0], ps.FindNearestValid(ps[0]-1))

	ps[2] = 0
	ps[0] = 0
	require.Equal(t, ps[1]-1, ps.FindNearestValid(ps[1]-1))

	ps[0] = psBase[0]
	ps[1] = 0
	require.Equal(t, ps[0]+1, ps.FindNearestValid(ps[0]+1))

	ps[1] = psBase[1]

	require.Equal(t, ps[1]+1, ps.FindNearestValid(ps[1]+1))

	require.Equal(t, ps[1], ps.FindNearestValid(ps[1]-1))

	require.Equal(t, ps[0], ps.FindNearestValid(ps[0]+1))

	ps = psBase
	ps[1] = 0
	require.Equal(t, ps[2], ps.FindNearestValid(ps[2]+1))

	require.Equal(t, ps[3], ps.FindNearestValid(ps[3]-1))

	require.Equal(t, ps[2]-1, ps.FindNearestValid(ps[2]-1))

	require.Equal(t, ps[0], ps.FindNearestValid(ps[0]-1))

	ps = psBase
	ps[0] = 0
	require.Equal(t, ps[1]-1, ps.FindNearestValid(ps[1]-1))

	require.Equal(t, ps[1], ps.FindNearestValid(ps[1]+1))

	require.Equal(t, ps[2], ps.FindNearestValid(ps[2]-1))

	ps = psBase
	require.Equal(t, ps[0], ps.FindNearestValid(ps[0]+1))

	require.Equal(t, ps[1], ps.FindNearestValid(ps[1]-1))

	require.Equal(t, ps[1], ps.FindNearestValid(ps[1]+1))

	require.Equal(t, ps[2], ps.FindNearestValid(ps[2]-1))

	require.Equal(t, ps[2], ps.FindNearestValid(ps[2]+1))

	require.Equal(t, ps[3], ps.FindNearestValid(ps[3]-1))
}
