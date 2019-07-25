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

package censusimpl

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/stretchr/testify/require"
)

func TestNPSGetNodeID(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	nps := NodeProfileSlot{StaticProfile: sp}
	require.Equal(t, nodeID, nps.GetNodeID())
}

func TestNPSGetStatic(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	nps := NodeProfileSlot{StaticProfile: sp}
	require.Equal(t, sp, nps.GetStatic())
}

func TestNewNodeProfile(t *testing.T) {
	index := member.Index(1)
	sp := profiles.NewStaticProfileMock(t)
	sv := cryptkit.NewSignatureVerifierMock(t)
	power := member.Power(1)
	nps := NewNodeProfile(index, sp, sv, power)
	require.Equal(t, index, nps.index)

	require.Equal(t, sp, nps.StaticProfile)

	require.Equal(t, sv, nps.verifier)

	require.Equal(t, power, nps.power)

	require.Panics(t, func() { NewNodeProfile(member.MaxNodeIndex+1, sp, sv, power) })
}

func TestNewJoinerProfile(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	sv := cryptkit.NewSignatureVerifierMock(t)
	nps := NewJoinerProfile(sp, sv)
	require.Equal(t, member.JoinerIndex, nps.index)

	require.Equal(t, sp, nps.StaticProfile)

	require.Equal(t, sv, nps.verifier)

	require.Zero(t, nps.power)
}

func TestNewNodeProfileExt(t *testing.T) {
	index := member.Index(1)
	sp := profiles.NewStaticProfileMock(t)
	sv := cryptkit.NewSignatureVerifierMock(t)
	power := member.Power(2)
	mode := member.ModeSuspected
	nps := NewNodeProfileExt(index, sp, sv, power, mode)
	require.Equal(t, index, nps.index)

	require.Equal(t, sp, nps.StaticProfile)

	require.Equal(t, sv, nps.verifier)

	require.Equal(t, power, nps.power)

	require.Equal(t, mode, nps.mode)

	require.Panics(t, func() { NewNodeProfileExt(member.MaxNodeIndex+1, sp, sv, power, mode) })
}

func TestGetDeclaredPower(t *testing.T) {
	power := member.Power(1)
	nps := NodeProfileSlot{power: power}
	require.Equal(t, power, nps.GetDeclaredPower())
}

func TestNPSGetOpMode(t *testing.T) {
	mode := member.ModeSuspected
	nps := NodeProfileSlot{mode: mode}
	require.Equal(t, mode, nps.GetOpMode())
}

func TestLocalNodeProfile(t *testing.T) {
	nps := NodeProfileSlot{}
	require.NotPanics(t, func() { nps.LocalNodeProfile() })
}

func TestGetIndex(t *testing.T) {
	index := member.Index(1)
	nps := NodeProfileSlot{index: index}
	require.Equal(t, index, nps.GetIndex())

	nps.index = member.MaxNodeIndex + 1
	require.Panics(t, func() { nps.GetIndex() })
}

func TestIsJoiner(t *testing.T) {
	index := member.Index(1)
	nps := NodeProfileSlot{index: index}
	require.False(t, nps.IsJoiner())

	nps.index = member.JoinerIndex
	require.True(t, nps.IsJoiner())
}

func TestNPGetSignatureVerifier(t *testing.T) {
	sv := cryptkit.NewSignatureVerifierMock(t)
	nps := NodeProfileSlot{verifier: sv}
	require.Equal(t, sv, nps.GetSignatureVerifier())
}

func TestString(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	nps := NodeProfileSlot{StaticProfile: sp, index: 1}
	require.NotEmpty(t, nps.String())

	nps.index = member.JoinerIndex
	require.NotEmpty(t, nps.String())
}

func TestAsActiveNode(t *testing.T) {
	nps := NodeProfileSlot{index: 1}
	us := updatableSlot{NodeProfileSlot: nps}
	act := us.AsActiveNode()
	require.Equal(t, &nps, act)

	require.Implements(t, (*profiles.ActiveNode)(nil), act)
}

func TestSetRank(t *testing.T) {
	index := member.Index(1)
	mode := member.ModeSuspected
	power := member.Power(2)
	us := updatableSlot{}
	us.SetRank(index, mode, power)

	require.Equal(t, index, us.index)

	require.Equal(t, power, us.power)

	require.Equal(t, mode, us.mode)

	require.Panics(t, func() { us.SetRank(member.MaxNodeIndex+1, mode, power) })
}

func TestSetPower(t *testing.T) {
	us := updatableSlot{}
	power := member.Power(1)
	us.SetPower(power)
	require.Equal(t, power, us.power)
}

func TestSetOpMode(t *testing.T) {
	us := updatableSlot{}
	mode := member.ModeSuspected
	us.SetOpMode(mode)
	require.Equal(t, mode, us.GetOpMode())
}

func TestSetOpModeAndLeaveReason(t *testing.T) {
	index := member.Index(1)
	leaveReason := uint32(2)
	us := updatableSlot{}
	us.SetOpModeAndLeaveReason(index, leaveReason)
	require.Equal(t, index, us.index)

	require.Zero(t, us.power)

	require.Equal(t, member.ModeEvictedGracefully, us.mode)

	require.Equal(t, leaveReason, us.leaveReason)

	require.Panics(t, func() { us.SetOpModeAndLeaveReason(member.MaxNodeIndex+1, leaveReason) })
}

func TestUSGetLeaveReason(t *testing.T) {
	leaveReason := uint32(1)
	us := updatableSlot{leaveReason: leaveReason}
	require.Zero(t, us.GetLeaveReason())

	us.mode = member.ModeEvictedGracefully
	require.Equal(t, leaveReason, us.GetLeaveReason())
}

func TestSetIndex(t *testing.T) {
	us := updatableSlot{}
	index := member.Index(1)
	us.SetIndex(index)
	require.Equal(t, index, us.index)

	require.Panics(t, func() { us.SetIndex(member.MaxNodeIndex + 1) })
}

func TestSetSignatureVerifier(t *testing.T) {
	us := updatableSlot{}
	sv := cryptkit.NewSignatureVerifierMock(t)
	us.SetSignatureVerifier(sv)
	require.Equal(t, sv, us.verifier)
}
