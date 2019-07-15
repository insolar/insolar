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

package profiles

import (
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemberPowerOf(t *testing.T) {
	require.Equal(t, member.Power(1), member.PowerOf(1))

	require.Equal(t, member.Power(0x1F), member.PowerOf(0x1F))

	require.Equal(t, member.Power(0xFF), member.PowerOf(member.MaxLinearMemberPower))

	require.Equal(t, member.Power(0xFF), member.PowerOf(member.MaxLinearMemberPower+1))

	require.Equal(t, member.Power(0x1F+1), member.PowerOf(0x1F+1))

	require.Equal(t, member.Power(0x2F), member.PowerOf(0x1F<<1))
}

func TestToLinearValue(t *testing.T) {
	require.Equal(t, uint16(0), member.PowerOf(0).ToLinearValue())

	require.Equal(t, uint16(0x1F), member.PowerOf(0x1F).ToLinearValue())

	require.Equal(t, uint16(0x1F+1), member.PowerOf(0x1F+1).ToLinearValue())

	require.Equal(t, uint16(0x3e), member.PowerOf(0x1F<<1).ToLinearValue())
}

func TestPercentAndMin(t *testing.T) {
	require.Equal(t, ^member.Power(0), member.PowerOf(member.MaxLinearMemberPower).PercentAndMin(100, member.PowerOf(0)))

	require.Equal(t, member.Power(2), member.PowerOf(3).PercentAndMin(1, member.PowerOf(2)))

	require.Equal(t, member.Power(2), member.PowerOf(3).PercentAndMin(80, member.PowerOf(1)))
}

func TestNormalize(t *testing.T) {
	zero := member.PowerSet([...]member.Power{0, 0, 0, 0})
	require.Equal(t, zero, zero.Normalize())

	require.Equal(t, zero, member.PowerSet([...]member.Power{1, 0, 0, 0}).Normalize())

	m := member.PowerSet([...]member.Power{1, 1, 1, 1})
	require.Equal(t, m, m.Normalize())
}

// Illegal cases:
// [ x,  y,  z,  0] - when any !=0 value of x, y, z
// [ 0,  x,  0,  y] - when x != 0 and y != 0
// any combination of non-zero x, y such that x > y and y > 0 and position(x) < position(y)
// And cases from the function logic.
func TestIsValid(t *testing.T) {
	require.True(t, member.PowerSet([...]member.Power{0, 0, 0, 0}).IsValid())

	require.False(t, member.PowerSet([...]member.Power{1, 0, 0, 0}).IsValid())

	require.False(t, member.PowerSet([...]member.Power{0, 1, 0, 0}).IsValid())

	require.False(t, member.PowerSet([...]member.Power{0, 0, 1, 0}).IsValid())

	require.False(t, member.PowerSet([...]member.Power{0, 1, 0, 1}).IsValid())

	require.False(t, member.PowerSet([...]member.Power{2, 1, 2, 2}).IsValid())

	require.True(t, member.PowerSet([...]member.Power{1, 0, 0, 1}).IsValid())

	require.True(t, member.PowerSet([...]member.Power{1, 1, 0, 1}).IsValid())

	require.False(t, member.PowerSet([...]member.Power{1, 1, 2, 1}).IsValid())

	require.True(t, member.PowerSet([...]member.Power{1, 0, 2, 2}).IsValid())

	require.True(t, member.PowerSet([...]member.Power{1, 1, 2, 2}).IsValid())
}

func TestNewPowerRequestByLevel(t *testing.T) {
	require.Equal(t, -power.Request(capacity.LevelMinimal), power.NewRequestByLevel(capacity.LevelMinimal))
}

func TestNewPowerRequest(t *testing.T) {
	require.Equal(t, power.Request(1), power.NewRequest(member.Power(1)))
}

func TestAsCapacityLevel(t *testing.T) {
	b, l := power.Request(-1).AsCapacityLevel()
	require.True(t, b)
	require.Equal(t, capacity.Level(1), l)

	b, l = power.Request(1).AsCapacityLevel()
	require.False(t, b)

	r := power.Request(1)
	require.Equal(t, capacity.Level(-r), l)

	b, l = power.Request(0).AsCapacityLevel()
	require.False(t, b)
	require.Equal(t, capacity.Level(0), l)
}

func TestAsMemberPower(t *testing.T) {
	b, l := power.Request(1).AsMemberPower()
	require.True(t, b)
	require.Equal(t, member.Power(1), l)

	b, l = power.Request(-1).AsMemberPower()
	require.False(t, b)

	r := power.Request(-1)
	require.Equal(t, member.Power(r), l)

	b, l = power.Request(0).AsMemberPower()
	require.True(t, b)
	require.Equal(t, member.Power(0), l)
}

//func TestNewMembershipProfile(t *testing.T) {
//	nsh := common2.NewNodeStateHashEvidenceMock(t)
//	nas := common2.NewMemberAnnouncementSignatureMock(t)
//	index := uint16(1)
//	power := Power(2)
//	ep := Power(3)
//	mp := NewMembershipProfile(MemberModeNormal, power, index, nsh, nas, ep)
//	require.Equal(t, index, mp.Index)
//
//	require.Equal(t, power, mp.Power)
//
//	require.Equal(t, ep, mp.RequestedPower)
//
//	require.Equal(t, nsh, mp.StateEvidence)
//
//	require.Equal(t, nas, mp.AnnounceSignature)
//}
//
//func TestNewMembershipProfileByNode(t *testing.T) {
//	np := NewNodeProfileMock(t)
//	index := 1
//	np.GetIndexMock.Set(func() int { return index })
//	power := Power(2)
//	np.GetDeclaredPowerMock.Set(func() Power { return power })
//	np.GetOpModeMock.Set(func() (r MemberOpMode) {
//		return MemberModeNormal
//	})
//
//	nsh := common2.NewNodeStateHashEvidenceMock(t)
//	nas := common2.NewMemberAnnouncementSignatureMock(t)
//	ep := Power(3)
//	mp := NewMembershipProfileByNode(np, nsh, nas, ep)
//	require.Equal(t, uint16(index), mp.Index)
//
//	require.Equal(t, power, mp.Power)
//
//	require.Equal(t, ep, mp.RequestedPower)
//
//	require.Equal(t, nsh, mp.StateEvidence)
//
//	require.Equal(t, nas, mp.AnnounceSignature)
//}
//
//func TestIsEmpty(t *testing.T) {
//	mp := MembershipProfile{}
//	require.True(t, mp.IsEmpty())
//
//	se := common2.NewNodeStateHashEvidenceMock(t)
//	mp.StateEvidence = se
//	require.True(t, mp.IsEmpty())
//
//	mp.StateEvidence = nil
//	mp.AnnounceSignature = common2.NewMemberAnnouncementSignatureMock(t)
//	require.True(t, mp.IsEmpty())
//
//	mp.StateEvidence = se
//	require.False(t, mp.IsEmpty())
//}

//func TestEquals(t *testing.T) {
//	mp1 := MembershipProfile{}
//	mp2 := MembershipProfile{}
//	require.False(t, mp1.Equals(mp2))
//
//	mp1.Index = uint16(1)
//	mp1.Power = Power(2)
//	mp1.RequestedPower = Power(3)
//	she1 := common2.NewNodeStateHashEvidenceMock(t)
//	mas1 := common2.NewMemberAnnouncementSignatureMock(t)
//	mp1.StateEvidence = she1
//	mp1.AnnounceSignature = mas1
//
//	mp2.Index = uint16(2)
//	mp2.Power = mp1.Power
//	mp2.RequestedPower = mp1.RequestedPower
//	mp2.StateEvidence = mp1.StateEvidence
//	mp2.AnnounceSignature = mp1.AnnounceSignature
//
//	require.False(t, mp1.Equals(mp2))
//
//	mp2.Index = mp1.Index
//	mp2.Power = Power(3)
//	require.False(t, mp1.Equals(mp2))
//
//	mp2.Power = mp1.Power
//	mp2.StateEvidence = nil
//	require.False(t, mp1.Equals(mp2))
//
//	mp2.StateEvidence = mp1.StateEvidence
//	mp2.AnnounceSignature = nil
//	require.False(t, mp1.Equals(mp2))
//
//	mp2.AnnounceSignature = mp1.AnnounceSignature
//	mp1.StateEvidence = nil
//	require.False(t, mp1.Equals(mp2))
//
//	mp1.StateEvidence = mp2.StateEvidence
//	mp1.AnnounceSignature = nil
//	require.False(t, mp1.Equals(mp2))
//
//	mp1.AnnounceSignature = mp2.AnnounceSignature
//	mp2.RequestedPower = Power(4)
//	require.False(t, mp1.Equals(mp2))
//
//	mp2.RequestedPower = mp1.RequestedPower
//	she2 := common2.NewNodeStateHashEvidenceMock(t)
//	mp2.StateEvidence = she2
//	nsh := common2.NewNodeStateHashMock(t)
//	common2.Set(func() NodeStateHash { return nsh })
//	common2.Set(func() NodeStateHash { return nsh })
//	common2.Set(func(common.DigestHolder) bool { return false })
//	require.False(t, mp1.Equals(mp2))
//
//	common2.Set(func(common.DigestHolder) bool { return true })
//	sh := common.NewSignatureHolderMock(t)
//	sh.EqualsMock.Set(func(common.SignatureHolder) bool { return false })
//	common2.Set(func() common.SignatureHolder { return sh })
//	common2.Set(func() common.SignatureHolder { return sh })
//	require.False(t, mp1.Equals(mp2))
//
//	sh.EqualsMock.Set(func(common.SignatureHolder) bool { return true })
//	require.True(t, mp1.Equals(mp2))
//
//	mp2.StateEvidence = she1
//	mas2 := common2.NewMemberAnnouncementSignatureMock(t)
//	mp2.AnnounceSignature = mas2
//	common2.Set(func(common.SignatureHolder) bool { return false })
//	require.False(t, mp1.Equals(mp2))
//
//	common2.Set(func(common.SignatureHolder) bool { return true })
//	require.True(t, mp1.Equals(mp2))
//}
//
//func TestStringParts(t *testing.T) {
//	mp := MembershipProfile{}
//	require.True(t, len(mp.StringParts()) > 0)
//	mp.Power = Power(1)
//	require.True(t, len(mp.StringParts()) > 0)
//}
//
//func TestMembershipProfileString(t *testing.T) {
//	mp := MembershipProfile{}
//	require.True(t, len(mp.String()) > 0)
//}
//
//func TestEqualIntroProfiles(t *testing.T) {
//	require.False(t, EqualIntroProfiles(nil, nil))
//	p := common2.NewNodeIntroProfileMock(t)
//	require.False(t, EqualIntroProfiles(p, nil))
//
//	require.False(t, EqualIntroProfiles(nil, p))
//
//	require.True(t, EqualIntroProfiles(p, p))
//
//	snID1 := common.ShortNodeID(1)
//	common2.Set(func() common.ShortNodeID { return *(&snID1) })
//	primaryRole1 := PrimaryRoleNeutral
//	common2.Set(func() NodePrimaryRole { return *(&primaryRole1) })
//	specialRole1 := SpecialRoleDiscovery
//	common2.Set(func() NodeSpecialRole { return *(&specialRole1) })
//	power1 := Power(1)
//	common2.Set(func() Power { return *(&power1) })
//	skh := common.NewSignatureKeyHolderMock(t)
//	signHoldEq := true
//	skh.EqualsMock.Set(func(common.SignatureKeyHolder) bool { return *(&signHoldEq) })
//	common2.Set(func() common.SignatureKeyHolder { return skh })
//
//	o := common2.NewNodeIntroProfileMock(t)
//	snID2 := common.ShortNodeID(2)
//	common2.Set(func() common.ShortNodeID { return *(&snID2) })
//	primaryRole2 := primaryRole1
//	common2.Set(func() NodePrimaryRole { return *(&primaryRole2) })
//	specialRole2 := specialRole1
//	common2.Set(func() NodeSpecialRole { return *(&specialRole2) })
//	power2 := power1
//	common2.Set(func() Power { return *(&power2) })
//	common2.Set(func() common.SignatureKeyHolder { return skh })
//	require.False(t, EqualIntroProfiles(p, o))
//
//	snID2 = snID1
//	primaryRole2 = PrimaryRoleHeavyMaterial
//	require.False(t, EqualIntroProfiles(p, o))
//
//	primaryRole2 = primaryRole1
//	specialRole2 = SpecialRoleNone
//	require.False(t, EqualIntroProfiles(p, o))
//
//	specialRole2 = specialRole1
//	power2 = Power(2)
//	require.False(t, EqualIntroProfiles(p, o))
//
//	power1 = power2
//	signHoldEq = false
//	require.False(t, EqualIntroProfiles(p, o))
//
//	signHoldEq = true
//	ne1 := common.NewNodeEndpointMock(t)
//	ne1.GetEndpointTypeMock.Set(func() common.NodeEndpointType { return common.NameEndpoint })
//	ne1.GetNameAddressMock.Set(func() common.HostAddress { return common.HostAddress("test1") })
//	common2.Set(func() common.NodeEndpoint { return ne1 })
//	ne2 := common.NewNodeEndpointMock(t)
//	ne2.GetEndpointTypeMock.Set(func() common.NodeEndpointType { return common.NameEndpoint })
//	ne2.GetNameAddressMock.Set(func() common.HostAddress { return common.HostAddress("test2") })
//	common2.Set(func() common.NodeEndpoint { return ne2 })
//	require.False(t, EqualIntroProfiles(p, o))
//
//	common2.Set(func() common.NodeEndpoint { return ne1 })
//	require.True(t, EqualIntroProfiles(p, o))
//}
