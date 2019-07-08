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

package common

import (
	"math"
	"testing"

	"github.com/insolar/insolar/network/consensus/common"

	"github.com/stretchr/testify/require"
)

func TestMemberPowerOf(t *testing.T) {
	require.Equal(t, MemberPowerOf(1), MemberPower(1))

	require.Equal(t, MemberPowerOf(0x1F), MemberPower(0x1F))

	require.Equal(t, MemberPowerOf(MaxLinearMemberPower), MemberPower(0xFF))

	require.Equal(t, MemberPowerOf(MaxLinearMemberPower+1), MemberPower(0xFF))

	require.Equal(t, MemberPowerOf(0x1F+1), MemberPower(0x1F+1))

	require.Equal(t, MemberPowerOf(0x1F<<1), MemberPower(0x2F))
}

func TestToLinearValue(t *testing.T) {
	require.Equal(t, MemberPowerOf(0).ToLinearValue(), uint16(0))

	require.Equal(t, MemberPowerOf(0x1F).ToLinearValue(), uint16(0x1F))

	require.Equal(t, MemberPowerOf(0x1F+1).ToLinearValue(), uint16(0x1F+1))

	require.Equal(t, MemberPowerOf(0x1F<<1).ToLinearValue(), uint16(0x3e))
}

func TestPercentAndMin(t *testing.T) {
	require.Equal(t, MemberPowerOf(MaxLinearMemberPower).PercentAndMin(100, MemberPowerOf(0)), ^MemberPower(0))

	require.Equal(t, MemberPowerOf(3).PercentAndMin(1, MemberPowerOf(2)), MemberPower(2))

	require.Equal(t, MemberPowerOf(3).PercentAndMin(80, MemberPowerOf(1)), MemberPower(2))
}

func TestNormalize(t *testing.T) {
	zero := MemberPowerSet([...]MemberPower{0, 0, 0, 0})
	require.Equal(t, zero.Normalize(), zero)

	require.Equal(t, MemberPowerSet([...]MemberPower{1, 0, 0, 0}).Normalize(), zero)

	m := MemberPowerSet([...]MemberPower{1, 1, 1, 1})
	require.Equal(t, m.Normalize(), m)
}

// Illegal cases:
// [ x,  y,  z,  0] - when any !=0 value of x, y, z
// [ 0,  x,  0,  y] - when x != 0 and y != 0
// any combination of non-zero x, y such that x > y and y > 0 and position(x) < position(y)
// And cases from the function logic.
func TestIsValid(t *testing.T) {
	require.True(t, MemberPowerSet([...]MemberPower{0, 0, 0, 0}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{1, 0, 0, 0}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{0, 1, 0, 0}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{0, 0, 1, 0}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{0, 1, 0, 1}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{2, 1, 2, 2}).IsValid())

	require.True(t, MemberPowerSet([...]MemberPower{1, 0, 0, 1}).IsValid())

	require.True(t, MemberPowerSet([...]MemberPower{1, 1, 0, 1}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{1, 1, 2, 1}).IsValid())

	require.True(t, MemberPowerSet([...]MemberPower{1, 0, 2, 2}).IsValid())

	require.True(t, MemberPowerSet([...]MemberPower{1, 1, 2, 2}).IsValid())
}

func TestNewPowerRequestByLevel(t *testing.T) {
	require.Equal(t, NewPowerRequestByLevel(common.LevelMinimal), -PowerRequest(common.LevelMinimal))
}

func TestNewPowerRequest(t *testing.T) {
	require.Equal(t, NewPowerRequest(MemberPower(1)), PowerRequest(1))
}

func TestAsCapacityLevel(t *testing.T) {
	b, l := PowerRequest(-1).AsCapacityLevel()
	require.True(t, b)
	require.Equal(t, l, common.CapacityLevel(1))

	b, l = PowerRequest(1).AsCapacityLevel()
	require.False(t, b)

	r := PowerRequest(1)
	require.Equal(t, l, common.CapacityLevel(-r))

	b, l = PowerRequest(0).AsCapacityLevel()
	require.False(t, b)
	require.Equal(t, l, common.CapacityLevel(0))
}

func TestAsMemberPower(t *testing.T) {
	b, l := PowerRequest(1).AsMemberPower()
	require.True(t, b)
	require.Equal(t, l, MemberPower(1))

	b, l = PowerRequest(-1).AsMemberPower()
	require.False(t, b)

	r := PowerRequest(-1)
	require.Equal(t, l, MemberPower(r))

	b, l = PowerRequest(0).AsMemberPower()
	require.True(t, b)
	require.Equal(t, l, MemberPower(0))
}

func TestIsSuspect(t *testing.T) {
	require.True(t, MembershipState(-2).IsSuspect())
	require.True(t, Suspected.IsSuspect())
	require.False(t, Undefined.IsSuspect())
	require.False(t, Joining.IsSuspect())
	require.False(t, Working.IsSuspect())
	require.False(t, JustJoined.IsSuspect())
}

func TestIsJustJoined(t *testing.T) {
	require.False(t, Suspected.IsJustJoined())
	require.False(t, Undefined.IsJustJoined())
	require.False(t, Joining.IsJustJoined())
	require.False(t, Working.IsJustJoined())
	require.True(t, JustJoined.IsJustJoined())
	require.True(t, MembershipState(4).IsJustJoined())
}

func TestGetCountInSuspected(t *testing.T) {
	require.Equal(t, Joining.GetCountInSuspected(), 0)
	require.Equal(t, Suspected.GetCountInSuspected(), 1)
	require.Equal(t, MembershipState(-2).GetCountInSuspected(), 2)
}

func TestAsJustJoinedRemainingCount(t *testing.T) {
	require.Equal(t, Working.AsJustJoinedRemainingCount(), 0)
	require.Equal(t, JustJoined.AsJustJoinedRemainingCount(), 1)
	require.Equal(t, MembershipState(4).AsJustJoinedRemainingCount(), 2)
}

func TestInSuspectedExceeded(t *testing.T) {
	require.Panics(t, func() { Working.InSuspectedExceeded(-1) })
	require.Panics(t, func() { Working.InSuspectedExceeded(int(Suspected) - math.MinInt8 + 1) })
	require.False(t, Working.InSuspectedExceeded(0))
	require.True(t, Suspected.InSuspectedExceeded(0))
	require.False(t, Suspected.InSuspectedExceeded(1))
}

func TestSetJustJoined(t *testing.T) {
	require.Panics(t, func() { Working.SetJustJoined(0) })
	require.Panics(t, func() { Working.SetJustJoined(math.MaxInt8 - int(JustJoined) + 1) })
	require.Equal(t, Working.SetJustJoined(1), JustJoined+MembershipState(1)-1)
}

func TestIncrementSuspected(t *testing.T) {
	require.Panics(t, func() { Undefined.IncrementSuspected() })
	require.Panics(t, func() { MembershipState(math.MinInt8).IncrementSuspected() })
	require.Equal(t, Suspected.IncrementSuspected(), MembershipState(-2))
	require.Equal(t, Working.IncrementSuspected(), Suspected)
}

func TestDecrementJustJoined(t *testing.T) {
	require.Panics(t, func() { Undefined.DecrementJustJoined() })
	require.Equal(t, JustJoined.DecrementJustJoined(), MembershipState(JustJoined-1))
	require.Equal(t, Working.DecrementJustJoined(), Working)
}

func TestUpdateOnNextPulse(t *testing.T) {
	require.Panics(t, func() { Undefined.UpdateOnNextPulse(0) })
	require.Equal(t, Joining.UpdateOnNextPulse(0), Working)
	require.Panics(t, func() { Joining.UpdateOnNextPulse(math.MaxInt8 - int(JustJoined) + 1) })
	require.Equal(t, Joining.UpdateOnNextPulse(1), JustJoined)
	require.Equal(t, Suspected.UpdateOnNextPulse(2), Suspected-1)
	require.Equal(t, JustJoined.UpdateOnNextPulse(2), JustJoined-1)
	require.Equal(t, Working.UpdateOnNextPulse(2), Working)
}

func TestIsUndefined(t *testing.T) {
	require.False(t, MembershipState(Suspected-2).IsUndefined())
	require.False(t, Suspected.IsUndefined())
	require.True(t, Undefined.IsUndefined())
	require.False(t, Joining.IsUndefined())
	require.False(t, Working.IsUndefined())
	require.False(t, JustJoined.IsUndefined())
	require.False(t, MembershipState(JustJoined+2).IsUndefined())
}

func TestIsActive(t *testing.T) {
	require.True(t, MembershipState(Suspected-2).IsActive())
	require.True(t, Suspected.IsActive())
	require.False(t, Undefined.IsActive())
	require.False(t, Joining.IsActive())
	require.True(t, Working.IsActive())
	require.True(t, JustJoined.IsActive())
	require.True(t, MembershipState(JustJoined+2).IsActive())
}

func TestIsWorking(t *testing.T) {
	require.False(t, MembershipState(Suspected-2).IsWorking())
	require.False(t, Suspected.IsWorking())
	require.False(t, Undefined.IsWorking())
	require.False(t, Joining.IsWorking())
	require.True(t, Working.IsWorking())
	require.True(t, JustJoined.IsWorking())
	require.True(t, MembershipState(JustJoined+2).IsWorking())
}

func TestIsJoining(t *testing.T) {
	require.False(t, MembershipState(Suspected-2).IsJoining())
	require.False(t, Suspected.IsJoining())
	require.False(t, Undefined.IsJoining())
	require.True(t, Joining.IsJoining())
	require.False(t, Working.IsJoining())
	require.False(t, JustJoined.IsJoining())
	require.False(t, MembershipState(JustJoined+2).IsJoining())
}

func TestNodeProfileOrdering(t *testing.T) {
	np := NewNodeProfileMock(t)
	power1 := MemberPower(0)
	np.GetDeclaredPowerMock.Set(func() MemberPower { return power1 })
	role := PrimaryRoleNeutral
	np.GetPrimaryRoleMock.Set(func() NodePrimaryRole { return role })
	shortNodeID := common.ShortNodeID(2)
	np.GetShortNodeIDMock.Set(func() common.ShortNodeID { return shortNodeID })
	r, p, id := nodeProfileOrdering(np)
	require.Equal(t, r, PrimaryRoleInactive)

	require.Equal(t, p, MemberPower(0))

	require.Equal(t, id, shortNodeID)

	power2 := MemberPower(1)
	np.GetDeclaredPowerMock.Set(func() MemberPower { return power2 })
	np.GetStateMock.Set(func() MembershipState { return Suspected })
	r, p, id = nodeProfileOrdering(np)
	require.Equal(t, r, PrimaryRoleInactive)

	require.Equal(t, p, MemberPower(0))

	require.Equal(t, id, shortNodeID)

	np.GetStateMock.Set(func() MembershipState { return Working })
	r, p, id = nodeProfileOrdering(np)
	require.Equal(t, r, role)

	require.Equal(t, p, power2)

	require.Equal(t, id, shortNodeID)
}
