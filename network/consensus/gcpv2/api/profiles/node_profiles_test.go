package profiles

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"

	"github.com/stretchr/testify/require"
)

func TestNewMembershipProfile(t *testing.T) {
	nsh := proofs.NewNodeStateHashEvidenceMock(t)
	nas := proofs.NewMemberAnnouncementSignatureMock(t)
	index := member.AsIndex(1)
	power := member.Power(2)
	ep := member.Power(3)
	mp := NewMembershipProfile(member.ModeNormal, power, index, nsh, nas, ep)
	require.Equal(t, index, mp.Index)

	require.Equal(t, power, mp.Power)

	require.Equal(t, ep, mp.RequestedPower)

	require.Equal(t, nsh, mp.StateEvidence)

	require.Equal(t, nas, mp.AnnounceSignature)
}

func TestNewMembershipProfileByNode(t *testing.T) {
	np := NewActiveNodeMock(t)
	index := 1
	np.GetIndexMock.Set(func() member.Index { return member.Index(index) })
	np.IsJoinerMock.Set(func() (r bool) {
		return false
	})
	power := member.Power(2)
	np.GetDeclaredPowerMock.Set(func() member.Power { return power })
	np.GetOpModeMock.Set(func() (r member.OpMode) {
		return member.ModeNormal
	})

	nsh := proofs.NewNodeStateHashEvidenceMock(t)
	nas := proofs.NewMemberAnnouncementSignatureMock(t)
	ep := member.Power(3)
	mp := NewMembershipProfileByNode(np, nsh, nas, ep)
	require.Equal(t, member.Index(index), mp.Index)

	require.Equal(t, power, mp.Power)

	require.Equal(t, ep, mp.RequestedPower)

	require.Equal(t, nsh, mp.StateEvidence)

	require.Equal(t, nas, mp.AnnounceSignature)
}

func TestIsEmpty(t *testing.T) {
	mp := MembershipProfile{}
	require.True(t, mp.IsEmpty())

	se := proofs.NewNodeStateHashEvidenceMock(t)
	mp.StateEvidence = se
	require.True(t, mp.IsEmpty())

	mp.StateEvidence = nil
	mp.AnnounceSignature = proofs.NewMemberAnnouncementSignatureMock(t)
	require.True(t, mp.IsEmpty())

	mp.StateEvidence = se
	require.False(t, mp.IsEmpty())
}

func TestEquals(t *testing.T) {
	mp1 := MembershipProfile{}
	mp2 := MembershipProfile{}
	require.False(t, mp1.Equals(mp2))

	mp1.Index = member.AsIndex(1)
	mp1.Power = member.Power(2)
	mp1.RequestedPower = member.Power(3)
	she1 := proofs.NewNodeStateHashEvidenceMock(t)
	mas1 := proofs.NewMemberAnnouncementSignatureMock(t)
	mp1.StateEvidence = she1
	mp1.AnnounceSignature = mas1

	mp2.Index = member.AsIndex(2)
	mp2.Power = mp1.Power
	mp2.RequestedPower = mp1.RequestedPower
	mp2.StateEvidence = mp1.StateEvidence
	mp2.AnnounceSignature = mp1.AnnounceSignature

	require.False(t, mp1.Equals(mp2))

	mp2.Index = mp1.Index
	mp2.Power = member.Power(3)
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
	mp2.RequestedPower = member.Power(4)
	require.False(t, mp1.Equals(mp2))

	mp2.RequestedPower = mp1.RequestedPower
	// TODO
	/*she2 := proofs.NewNodeStateHashEvidenceMock(t)
	mp2.StateEvidence = she2
	nsh := proofs.NewNodeStateHashMock(t)
	she1.GetNodeStateHashMock.Set(func() proofs.NodeStateHash { return nsh })
	she2.GetNodeStateHashMock.Set(func() proofs.NodeStateHash { return nsh })
	nsh.EqualsMock.Set(func(cryptkit.DigestHolder) bool { return false })
	require.False(t, mp1.Equals(mp2))

	nsh.EqualsMock.Set(func(cryptkit.DigestHolder) bool { return true })
	sh := cryptkit.NewSignatureHolderMock(t)
	sh.EqualsMock.Set(func(cryptkit.SignatureHolder) bool { return false })
	she1.GetGlobulaNodeStateSignatureMock.Set(func() cryptkit.SignatureHolder { return sh })
	she2.GetGlobulaNodeStateSignatureMock.Set(func() cryptkit.SignatureHolder { return sh })
	require.False(t, mp1.Equals(mp2))

	sh.EqualsMock.Set(func(cryptkit.SignatureHolder) bool { return true })
	require.True(t, mp1.Equals(mp2))

	mp2.StateEvidence = she1
	mas2 := proofs.NewMemberAnnouncementSignatureMock(t)
	mp2.AnnounceSignature = mas2
	mas1.EqualsMock.Set(func(cryptkit.SignatureHolder) bool { return false })
	require.False(t, mp1.Equals(mp2))

	mas1.EqualsMock.Set(func(cryptkit.SignatureHolder) bool { return true })
	require.True(t, mp1.Equals(mp2))*/
}

func TestStringParts(t *testing.T) {
	mp := MembershipProfile{}
	require.True(t, len(mp.StringParts()) > 0)
	mp.Power = member.Power(1)
	require.True(t, len(mp.StringParts()) > 0)
}

func TestMembershipProfileString(t *testing.T) {
	mp := MembershipProfile{}
	require.True(t, len(mp.String()) > 0)
}

func TestEqualIntroProfiles(t *testing.T) {
	require.False(t, EqualBriefProfiles(nil, nil))
	// TODO
	/*p := NewNodeIntroProfileMock(t)
	require.False(t, EqualBriefProfiles(p, nil))

	require.False(t, EqualBriefProfiles(nil, p))

	require.True(t, EqualBriefProfiles(p, p))

	snID1 := insolar.ShortNodeID(1)
	p.GetShortNodeIDMock.Set(func() insolar.ShortNodeID { return *(&snID1) })
	primaryRole1 := member.PrimaryRoleNeutral
	p.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return *(&primaryRole1) })
	specialRole1 := member.SpecialRoleDiscovery
	p.GetSpecialRolesMock.Set(func() member.SpecialRole { return *(&specialRole1) })
	power1 := member.Power(1)
	p.GetStartPowerMock.Set(func() member.Power { return *(&power1) })
	skh := cryptkit.NewSignatureKeyHolderMock(t)
	signHoldEq := true
	skh.EqualsMock.Set(func(cryptkit.SignatureKeyHolder) bool { return *(&signHoldEq) })
	p.GetNodePublicKeyMock.Set(func() cryptkit.SignatureKeyHolder { return skh })

	o := NewNodeIntroProfileMock(t)
	snID2 := insolar.ShortNodeID(2)
	o.GetShortNodeIDMock.Set(func() insolar.ShortNodeID { return *(&snID2) })
	primaryRole2 := primaryRole1
	o.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return *(&primaryRole2) })
	specialRole2 := specialRole1
	o.GetSpecialRolesMock.Set(func() member.SpecialRole { return *(&specialRole2) })
	power2 := power1
	o.GetStartPowerMock.Set(func() member.Power { return *(&power2) })
	o.GetNodePublicKeyMock.Set(func() cryptkit.SignatureKeyHolder { return skh })
	require.False(t, EqualBriefProfiles(p, o))

	snID2 = snID1
	primaryRole2 = member.PrimaryRoleHeavyMaterial
	require.False(t, EqualBriefProfiles(p, o))

	primaryRole2 = primaryRole1
	specialRole2 = member.SpecialRoleNone
	require.False(t, EqualBriefProfiles(p, o))

	specialRole2 = specialRole1
	power2 = member.Power(2)
	require.False(t, EqualBriefProfiles(p, o))

	power1 = power2
	signHoldEq = false
	require.False(t, EqualBriefProfiles(p, o))

	signHoldEq = true
	ne1 := endpoints.NewOutboundMock(t)
	ne1.GetEndpointTypeMock.Set(func() endpoints.NodeEndpointType { return endpoints.NameEndpoint })
	ne1.GetNameAddressMock.Set(func() endpoints.Name { return endpoints.Name("test1") })
	p.GetDefaultEndpointMock.Set(func() endpoints.Outbound { return ne1 })
	ne2 := endpoints.NewOutboundMock(t)
	ne2.GetEndpointTypeMock.Set(func() endpoints.NodeEndpointType { return endpoints.NameEndpoint })
	ne2.GetNameAddressMock.Set(func() endpoints.Name { return endpoints.Name("test2") })
	o.GetDefaultEndpointMock.Set(func() endpoints.Outbound { return ne2 })
	require.False(t, EqualBriefProfiles(p, o))

	o.GetDefaultEndpointMock.Set(func() endpoints.Outbound { return ne1 })
	sh := cryptkit.NewSignatureHolderMock(t)
	equal := false
	sh.EqualsMock.Set(func(cryptkit.SignatureHolder) bool { return *(&equal) })
	p.GetAnnouncementSignatureMock.Set(func() cryptkit.SignatureHolder { return sh })
	o.GetAnnouncementSignatureMock.Set(func() cryptkit.SignatureHolder { return sh })
	require.False(t, EqualBriefProfiles(p, o))

	equal = true
	require.True(t, EqualBriefProfiles(p, o))*/
}
