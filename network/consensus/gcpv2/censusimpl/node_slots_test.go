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

func TestIsPowered(t *testing.T) {
	index := member.JoinerIndex
	nps := NodeProfileSlot{index: index}
	require.False(t, nps.IsPowered())

	nps.index = 1
	nps.power = 1
	require.True(t, nps.IsPowered())

	nps.mode = member.ModeFlagSuspendedOps
	require.False(t, nps.IsPowered())
}

func TestIsVoter(t *testing.T) {
	index := member.JoinerIndex
	nps := NodeProfileSlot{index: index}
	require.False(t, nps.IsVoter())

	nps.index = 1
	nps.mode = member.ModeNormal
	require.True(t, nps.IsVoter())

	nps.mode = member.ModeFlagSuspendedOps
	require.False(t, nps.IsVoter())
}

func TestIsStateful(t *testing.T) {
	index := member.JoinerIndex
	nps := NodeProfileSlot{index: index}
	require.False(t, nps.IsStateful())

	nps.index = 1
	nps.mode = member.ModeNormal
	require.True(t, nps.IsStateful())

	// TODO: uncomment later
	// nps.mode = member.ModeFlagSuspendedOps
	// require.False(t, nps.IsStateful())
}

func TestCanIntroduceJoiner(t *testing.T) {
	index := member.JoinerIndex
	nps := NodeProfileSlot{index: index}
	require.False(t, nps.CanIntroduceJoiner())

	nps.index = 1
	nps.mode = member.ModeNormal
	require.True(t, nps.CanIntroduceJoiner())

	nps.mode = member.ModeFlagSuspendedOps
	require.False(t, nps.CanIntroduceJoiner())
}

func TestNPGetSignatureVerifier(t *testing.T) {
	sv := cryptkit.NewSignatureVerifierMock(t)
	nps := NodeProfileSlot{verifier: sv}
	require.Equal(t, sv, nps.GetSignatureVerifier())
}

func TestHasFullProfile(t *testing.T) {
	nps := NodeProfileSlot{}
	sp := profiles.NewStaticProfileMock(t)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return nil })
	nps.StaticProfile = sp
	require.False(t, nps.HasFullProfile())

	spe := profiles.NewStaticProfileExtensionMock(t)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return spe })
	require.True(t, nps.HasFullProfile())
}

func TestNPSString(t *testing.T) {
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
