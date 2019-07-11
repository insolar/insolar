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

package api

//import (
//	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
//	"github.com/insolar/insolar/network/consensus/common/endpoints"
//	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
//	"testing"
//
//	"github.com/insolar/insolar/network/consensus/common"
//
//	"github.com/stretchr/testify/require"
//)
//
//func TestMemberPowerOf(t *testing.T) {
//	require.Equal(t, MemberPowerOf(1), MemberPower(1))
//
//	require.Equal(t, MemberPowerOf(0x1F), MemberPower(0x1F))
//
//	require.Equal(t, MemberPowerOf(MaxLinearMemberPower), MemberPower(0xFF))
//
//	require.Equal(t, MemberPowerOf(MaxLinearMemberPower+1), MemberPower(0xFF))
//
//	require.Equal(t, MemberPowerOf(0x1F+1), MemberPower(0x1F+1))
//
//	require.Equal(t, MemberPowerOf(0x1F<<1), MemberPower(0x2F))
//}
//
//func TestToLinearValue(t *testing.T) {
//	require.Equal(t, ToLinearValue(), uint16(0))
//
//	require.Equal(t, ToLinearValue(), uint16(0x1F))
//
//	require.Equal(t, ToLinearValue(), uint16(0x1F+1))
//
//	require.Equal(t, ToLinearValue(), uint16(0x3e))
//}
//
//func TestPercentAndMin(t *testing.T) {
//	require.Equal(t, PercentAndMin(100, MemberPowerOf(0)), ^MemberPower(0))
//
//	require.Equal(t, PercentAndMin(1, MemberPowerOf(2)), MemberPower(2))
//
//	require.Equal(t, PercentAndMin(80, MemberPowerOf(1)), MemberPower(2))
//}
//
//func TestNormalize(t *testing.T) {
//	zero := MemberPowerSet([...]MemberPower{0, 0, 0, 0})
//	require.Equal(t, Normalize(), zero)
//
//	require.Equal(t, Normalize(), zero)
//
//	m := MemberPowerSet([...]MemberPower{1, 1, 1, 1})
//	require.Equal(t, Normalize(), m)
//}
//
//// Illegal cases:
//// [ x,  y,  z,  0] - when any !=0 value of x, y, z
//// [ 0,  x,  0,  y] - when x != 0 and y != 0
//// any combination of non-zero x, y such that x > y and y > 0 and position(x) < position(y)
//// And cases from the function logic.
//func TestIsValid(t *testing.T) {
//	require.True(t, IsValid())
//
//	require.False(t, IsValid())
//
//	require.False(t, IsValid())
//
//	require.False(t, IsValid())
//
//	require.False(t, IsValid())
//
//	require.False(t, IsValid())
//
//	require.True(t, IsValid())
//
//	require.True(t, IsValid())
//
//	require.False(t, IsValid())
//
//	require.True(t, IsValid())
//
//	require.True(t, IsValid())
//}
//
//func TestNewPowerRequestByLevel(t *testing.T) {
//	require.Equal(t, NewPowerRequestByLevel(common.LevelMinimal), -PowerRequest(common.LevelMinimal))
//}
//
//func TestNewPowerRequest(t *testing.T) {
//	require.Equal(t, NewPowerRequest(MemberPower(1)), PowerRequest(1))
//}
//
//func TestAsCapacityLevel(t *testing.T) {
//	b, l := AsCapacityLevel()
//	require.True(t, b)
//	require.Equal(t, l, common.CapacityLevel(1))
//
//	b, l = AsCapacityLevel()
//	require.False(t, b)
//
//	r := PowerRequest(1)
//	require.Equal(t, l, common.CapacityLevel(-r))
//
//	b, l = AsCapacityLevel()
//	require.False(t, b)
//	require.Equal(t, l, common.CapacityLevel(0))
//}
//
//func TestAsMemberPower(t *testing.T) {
//	b, l := AsMemberPower()
//	require.True(t, b)
//	require.Equal(t, l, MemberPower(1))
//
//	b, l = AsMemberPower()
//	require.False(t, b)
//
//	r := PowerRequest(-1)
//	require.Equal(t, l, MemberPower(r))
//
//	b, l = AsMemberPower()
//	require.True(t, b)
//	require.Equal(t, l, MemberPower(0))
//}
//
//func TestNewMembershipProfile(t *testing.T) {
//	nsh := common2.NewNodeStateHashEvidenceMock(t)
//	nas := common2.NewMemberAnnouncementSignatureMock(t)
//	index := uint16(1)
//	power := MemberPower(2)
//	ep := MemberPower(3)
//	mp := NewMembershipProfile(MemberModeNormal, power, index, nsh, nas, ep)
//	require.Equal(t, Index, index)
//
//	require.Equal(t, Power, power)
//
//	require.Equal(t, RequestedPower, ep)
//
//	require.Equal(t, common2.StateEvidence, nsh)
//
//	require.Equal(t, common2.AnnounceSignature, nas)
//}
//
//func TestNewMembershipProfileByNode(t *testing.T) {
//	np := common2.NewNodeProfileMock(t)
//	index := 1
//	common2.Set(func() int { return index })
//	power := MemberPower(2)
//	common2.Set(func() MemberPower { return power })
//	common2.Set(func() (r MemberOpMode) {
//		return MemberModeNormal
//	})
//
//	nsh := common2.NewNodeStateHashEvidenceMock(t)
//	nas := common2.NewMemberAnnouncementSignatureMock(t)
//	ep := MemberPower(3)
//	mp := NewMembershipProfileByNode(np, nsh, nas, ep)
//	require.Equal(t, Index, uint16(index))
//
//	require.Equal(t, Power, power)
//
//	require.Equal(t, RequestedPower, ep)
//
//	require.Equal(t, common2.StateEvidence, nsh)
//
//	require.Equal(t, common2.AnnounceSignature, nas)
//}
//
//func TestIsEmpty(t *testing.T) {
//	mp := MembershipProfile{}
//	require.True(t, IsEmpty())
//
//	se := common2.NewNodeStateHashEvidenceMock(t)
//	common2.StateEvidence = se
//	require.True(t, IsEmpty())
//
//	common2.StateEvidence = nil
//	common2.AnnounceSignature = common2.NewMemberAnnouncementSignatureMock(t)
//	require.True(t, IsEmpty())
//
//	common2.StateEvidence = se
//	require.False(t, IsEmpty())
//}
//
//func TestEquals(t *testing.T) {
//	mp1 := MembershipProfile{}
//	mp2 := MembershipProfile{}
//	require.False(t, Equals(mp2))
//
//	Index = uint16(1)
//	Power = MemberPower(2)
//	RequestedPower = MemberPower(3)
//	she1 := common2.NewNodeStateHashEvidenceMock(t)
//	mas1 := common2.NewMemberAnnouncementSignatureMock(t)
//	common2.StateEvidence = she1
//	common2.AnnounceSignature = mas1
//
//	Index = uint16(2)
//	Power = Power
//	RequestedPower = RequestedPower
//	common2.StateEvidence = common2.StateEvidence
//	common2.AnnounceSignature = common2.AnnounceSignature
//
//	require.False(t, Equals(mp2))
//
//	Index = Index
//	Power = MemberPower(3)
//	require.False(t, Equals(mp2))
//
//	Power = Power
//	common2.StateEvidence = nil
//	require.False(t, Equals(mp2))
//
//	common2.StateEvidence = common2.StateEvidence
//	common2.AnnounceSignature = nil
//	require.False(t, Equals(mp2))
//
//	common2.AnnounceSignature = common2.AnnounceSignature
//	common2.StateEvidence = nil
//	require.False(t, Equals(mp2))
//
//	common2.StateEvidence = common2.StateEvidence
//	common2.AnnounceSignature = nil
//	require.False(t, Equals(mp2))
//
//	common2.AnnounceSignature = common2.AnnounceSignature
//	RequestedPower = MemberPower(4)
//	require.False(t, Equals(mp2))
//
//	RequestedPower = RequestedPower
//	she2 := common2.NewNodeStateHashEvidenceMock(t)
//	common2.StateEvidence = she2
//	nsh := common2.NewNodeStateHashMock(t)
//	common2.Set(func() NodeStateHash { return nsh })
//	common2.Set(func() NodeStateHash { return nsh })
//	common2.Set(func(cryptography_containers.DigestHolder) bool { return false })
//	require.False(t, Equals(mp2))
//
//	common2.Set(func(cryptography_containers.DigestHolder) bool { return true })
//	sh := cryptography_containers.NewSignatureHolderMock(t)
//	sh.EqualsMock.Set(func(cryptography_containers.SignatureHolder) bool { return false })
//	common2.Set(func() cryptography_containers.SignatureHolder { return sh })
//	common2.Set(func() cryptography_containers.SignatureHolder { return sh })
//	require.False(t, Equals(mp2))
//
//	sh.EqualsMock.Set(func(cryptography_containers.SignatureHolder) bool { return true })
//	require.True(t, Equals(mp2))
//
//	common2.StateEvidence = she1
//	mas2 := common2.NewMemberAnnouncementSignatureMock(t)
//	common2.AnnounceSignature = mas2
//	common2.Set(func(cryptography_containers.SignatureHolder) bool { return false })
//	require.False(t, Equals(mp2))
//
//	common2.Set(func(cryptography_containers.SignatureHolder) bool { return true })
//	require.True(t, Equals(mp2))
//}
//
//func TestStringParts(t *testing.T) {
//	mp := MembershipProfile{}
//	require.True(t, len(StringParts()) > 0)
//	Power = MemberPower(1)
//	require.True(t, len(StringParts()) > 0)
//}
//
//func TestMembershipProfileString(t *testing.T) {
//	mp := MembershipProfile{}
//	require.True(t, len(String()) > 0)
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
//	power1 := MemberPower(1)
//	common2.Set(func() MemberPower { return *(&power1) })
//	skh := cryptography_containers.NewSignatureKeyHolderMock(t)
//	signHoldEq := true
//	skh.EqualsMock.Set(func(cryptography_containers.SignatureKeyHolder) bool { return *(&signHoldEq) })
//	common2.Set(func() cryptography_containers.SignatureKeyHolder { return skh })
//
//	o := common2.NewNodeIntroProfileMock(t)
//	snID2 := common.ShortNodeID(2)
//	common2.Set(func() common.ShortNodeID { return *(&snID2) })
//	primaryRole2 := primaryRole1
//	common2.Set(func() NodePrimaryRole { return *(&primaryRole2) })
//	specialRole2 := specialRole1
//	common2.Set(func() NodeSpecialRole { return *(&specialRole2) })
//	power2 := power1
//	common2.Set(func() MemberPower { return *(&power2) })
//	common2.Set(func() cryptography_containers.SignatureKeyHolder { return skh })
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
//	power2 = MemberPower(2)
//	require.False(t, EqualIntroProfiles(p, o))
//
//	power1 = power2
//	signHoldEq = false
//	require.False(t, EqualIntroProfiles(p, o))
//
//	signHoldEq = true
//	ne1 := endpoints.NewNodeEndpointMock(t)
//	ne1.GetEndpointTypeMock.Set(func() endpoints.NodeEndpointType { return endpoints.NameEndpoint })
//	ne1.GetNameAddressMock.Set(func() endpoints.HostAddress { return endpoints.HostAddress("test1") })
//	common2.Set(func() endpoints.NodeEndpoint { return ne1 })
//	ne2 := endpoints.NewNodeEndpointMock(t)
//	ne2.GetEndpointTypeMock.Set(func() endpoints.NodeEndpointType { return endpoints.NameEndpoint })
//	ne2.GetNameAddressMock.Set(func() endpoints.HostAddress { return endpoints.HostAddress("test2") })
//	common2.Set(func() endpoints.NodeEndpoint { return ne2 })
//	require.False(t, EqualIntroProfiles(p, o))
//
//	common2.Set(func() endpoints.NodeEndpoint { return ne1 })
//	require.True(t, EqualIntroProfiles(p, o))
//}
