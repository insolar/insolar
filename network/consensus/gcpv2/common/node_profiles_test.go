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

func TestNewMembershipProfile(t *testing.T) {
	nsh := NewNodeStateHashEvidenceMock(t)
	nas := NewMemberAnnouncementSignatureMock(t)
	index := uint16(1)
	power := MemberPower(2)
	ep := MemberPower(3)
	mp := NewMembershipProfile(MemberModeNormal, power, index, nsh, nas, ep)
	require.Equal(t, mp.Index, index)

	require.Equal(t, mp.Power, power)

	require.Equal(t, mp.RequestedPower, ep)

	require.Equal(t, mp.StateEvidence, nsh)

	require.Equal(t, mp.AnnounceSignature, nas)
}

func TestNewMembershipProfileByNode(t *testing.T) {
	np := NewNodeProfileMock(t)
	index := 1
	np.GetIndexMock.Set(func() int { return index })
	power := MemberPower(2)
	np.GetDeclaredPowerMock.Set(func() MemberPower { return power })
	np.GetOpModeMock.Set(func() (r MemberOpMode) {
		return MemberModeNormal
	})

	nsh := NewNodeStateHashEvidenceMock(t)
	nas := NewMemberAnnouncementSignatureMock(t)
	ep := MemberPower(3)
	mp := NewMembershipProfileByNode(np, nsh, nas, ep)
	require.Equal(t, mp.Index, uint16(index))

	require.Equal(t, mp.Power, power)

	require.Equal(t, mp.RequestedPower, ep)

	require.Equal(t, mp.StateEvidence, nsh)

	require.Equal(t, mp.AnnounceSignature, nas)
}

func TestIsEmpty(t *testing.T) {
	mp := MembershipProfile{}
	require.True(t, mp.IsEmpty())

	se := NewNodeStateHashEvidenceMock(t)
	mp.StateEvidence = se
	require.True(t, mp.IsEmpty())

	mp.StateEvidence = nil
	mp.AnnounceSignature = NewMemberAnnouncementSignatureMock(t)
	require.True(t, mp.IsEmpty())

	mp.StateEvidence = se
	require.False(t, mp.IsEmpty())
}

func TestEquals(t *testing.T) {
	mp1 := MembershipProfile{}
	mp2 := MembershipProfile{}
	require.False(t, mp1.Equals(mp2))

	mp1.Index = uint16(1)
	mp1.Power = MemberPower(2)
	mp1.RequestedPower = MemberPower(3)
	she1 := NewNodeStateHashEvidenceMock(t)
	mas1 := NewMemberAnnouncementSignatureMock(t)
	mp1.StateEvidence = she1
	mp1.AnnounceSignature = mas1

	mp2.Index = uint16(2)
	mp2.Power = mp1.Power
	mp2.RequestedPower = mp1.RequestedPower
	mp2.StateEvidence = mp1.StateEvidence
	mp2.AnnounceSignature = mp1.AnnounceSignature

	require.False(t, mp1.Equals(mp2))

	mp2.Index = mp1.Index
	mp2.Power = MemberPower(3)
	require.False(t, mp1.Equals(mp2))

	mp2.Power = mp1.Power
	mp2.StateEvidence = nil
	require.False(t, mp1.Equals(mp2))

	mp2.StateEvidence = mp1.StateEvidence
	mp2.AnnounceSignature = nil
	require.False(t, mp1.Equals(mp2))

	mp2.AnnounceSignature = mp1.AnnounceSignature
	mp1.StateEvidence = nil
	require.False(t, mp1.Equals(mp2))

	mp1.StateEvidence = mp2.StateEvidence
	mp1.AnnounceSignature = nil
	require.False(t, mp1.Equals(mp2))

	mp1.AnnounceSignature = mp2.AnnounceSignature
	mp2.RequestedPower = MemberPower(4)
	require.False(t, mp1.Equals(mp2))

	mp2.RequestedPower = mp1.RequestedPower
	she2 := NewNodeStateHashEvidenceMock(t)
	mp2.StateEvidence = she2
	nsh := NewNodeStateHashMock(t)
	she1.GetNodeStateHashMock.Set(func() NodeStateHash { return nsh })
	she2.GetNodeStateHashMock.Set(func() NodeStateHash { return nsh })
	nsh.EqualsMock.Set(func(common.DigestHolder) bool { return false })
	require.False(t, mp1.Equals(mp2))

	nsh.EqualsMock.Set(func(common.DigestHolder) bool { return true })
	sh := common.NewSignatureHolderMock(t)
	sh.EqualsMock.Set(func(common.SignatureHolder) bool { return false })
	she1.GetGlobulaNodeStateSignatureMock.Set(func() common.SignatureHolder { return sh })
	she2.GetGlobulaNodeStateSignatureMock.Set(func() common.SignatureHolder { return sh })
	require.False(t, mp1.Equals(mp2))

	sh.EqualsMock.Set(func(common.SignatureHolder) bool { return true })
	require.True(t, mp1.Equals(mp2))

	mp2.StateEvidence = she1
	mas2 := NewMemberAnnouncementSignatureMock(t)
	mp2.AnnounceSignature = mas2
	mas1.EqualsMock.Set(func(common.SignatureHolder) bool { return false })
	require.False(t, mp1.Equals(mp2))

	mas1.EqualsMock.Set(func(common.SignatureHolder) bool { return true })
	require.True(t, mp1.Equals(mp2))
}

func TestStringParts(t *testing.T) {
	mp := MembershipProfile{}
	require.True(t, len(mp.StringParts()) > 0)
	mp.Power = MemberPower(1)
	require.True(t, len(mp.StringParts()) > 0)
}

func TestMembershipProfileString(t *testing.T) {
	mp := MembershipProfile{}
	require.True(t, len(mp.String()) > 0)
}

func TestEqualIntroProfiles(t *testing.T) {
	require.False(t, EqualIntroProfiles(nil, nil))
	p := NewNodeIntroProfileMock(t)
	require.False(t, EqualIntroProfiles(p, nil))

	require.False(t, EqualIntroProfiles(nil, p))

	require.True(t, EqualIntroProfiles(p, p))

	snID1 := common.ShortNodeID(1)
	p.GetShortNodeIDMock.Set(func() common.ShortNodeID { return *(&snID1) })
	primaryRole1 := PrimaryRoleNeutral
	p.GetPrimaryRoleMock.Set(func() NodePrimaryRole { return *(&primaryRole1) })
	specialRole1 := SpecialRoleDiscovery
	p.GetSpecialRolesMock.Set(func() NodeSpecialRole { return *(&specialRole1) })
	power1 := MemberPower(1)
	p.GetStartPowerMock.Set(func() MemberPower { return *(&power1) })
	skh := common.NewSignatureKeyHolderMock(t)
	signHoldEq := true
	skh.EqualsMock.Set(func(common.SignatureKeyHolder) bool { return *(&signHoldEq) })
	p.GetNodePublicKeyMock.Set(func() common.SignatureKeyHolder { return skh })

	o := NewNodeIntroProfileMock(t)
	snID2 := common.ShortNodeID(2)
	o.GetShortNodeIDMock.Set(func() common.ShortNodeID { return *(&snID2) })
	primaryRole2 := primaryRole1
	o.GetPrimaryRoleMock.Set(func() NodePrimaryRole { return *(&primaryRole2) })
	specialRole2 := specialRole1
	o.GetSpecialRolesMock.Set(func() NodeSpecialRole { return *(&specialRole2) })
	power2 := power1
	o.GetStartPowerMock.Set(func() MemberPower { return *(&power2) })
	o.GetNodePublicKeyMock.Set(func() common.SignatureKeyHolder { return skh })
	require.False(t, EqualIntroProfiles(p, o))

	snID2 = snID1
	primaryRole2 = PrimaryRoleHeavyMaterial
	require.False(t, EqualIntroProfiles(p, o))

	primaryRole2 = primaryRole1
	specialRole2 = SpecialRoleNone
	require.False(t, EqualIntroProfiles(p, o))

	specialRole2 = specialRole1
	power2 = MemberPower(2)
	require.False(t, EqualIntroProfiles(p, o))

	power1 = power2
	signHoldEq = false
	require.False(t, EqualIntroProfiles(p, o))

	signHoldEq = true
	ne1 := common.NewNodeEndpointMock(t)
	ne1.GetEndpointTypeMock.Set(func() common.NodeEndpointType { return common.NameEndpoint })
	ne1.GetNameAddressMock.Set(func() common.HostAddress { return common.HostAddress("test1") })
	p.GetDefaultEndpointMock.Set(func() common.NodeEndpoint { return ne1 })
	ne2 := common.NewNodeEndpointMock(t)
	ne2.GetEndpointTypeMock.Set(func() common.NodeEndpointType { return common.NameEndpoint })
	ne2.GetNameAddressMock.Set(func() common.HostAddress { return common.HostAddress("test2") })
	o.GetDefaultEndpointMock.Set(func() common.NodeEndpoint { return ne2 })
	require.False(t, EqualIntroProfiles(p, o))

	o.GetDefaultEndpointMock.Set(func() common.NodeEndpoint { return ne1 })
	require.True(t, EqualIntroProfiles(p, o))
}
