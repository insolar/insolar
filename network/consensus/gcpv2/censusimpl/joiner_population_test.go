// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package censusimpl

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

func TestNewJoinerPopulation(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	pks := cryptkit.NewPublicKeyStoreMock(t)
	sp.GetPublicKeyStoreMock.Set(func() cryptkit.PublicKeyStore { return pks })
	vf := cryptkit.NewSignatureVerifierFactoryMock(t)
	sv := cryptkit.NewSignatureVerifierMock(t)
	vf.CreateSignatureVerifierWithPKSMock.Set(func(cryptkit.PublicKeyStore) cryptkit.SignatureVerifier { return sv })
	ojp := NewJoinerPopulation(sp, vf)
	require.Zero(t, ojp.localNode.mode)
}

func TestOJPGetSuspendedCount(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Zero(t, ojp.GetSuspendedCount())
}

func TestOJPGetMistrustedCount(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Zero(t, ojp.GetMistrustedCount())
}

func TestOJPGetIdleProfiles(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Nil(t, ojp.GetIdleProfiles())
}

func TestOJPGetIdleCount(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Zero(t, ojp.GetIdleCount())
}

func TestOJPGetIndexedCount(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Zero(t, ojp.GetIndexedCount())
}

func TestOJPGetIndexedCapacity(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Zero(t, ojp.GetIndexedCapacity())
}

func TestOJPIsValid(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.True(t, ojp.IsValid())
}

func TestOJPGetRolePopulation(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Nil(t, ojp.GetRolePopulation(member.PrimaryRoleNeutral))
}

func TestOJPGetWorkingRoles(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Nil(t, ojp.GetWorkingRoles())
}

func TestOJPCopyTo(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	index := member.JoinerIndex
	ojp := OneJoinerPopulation{localNode: updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp, index: index}}}
	population := &DynamicPopulation{}
	ojp.copyTo(population)

	require.Equal(t, index, ojp.localNode.index)

	require.Zero(t, population.local.index)
}

func TestOJPFindProfile(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	ojp := OneJoinerPopulation{localNode: updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}

	require.Nil(t, ojp.FindProfile(1))

	require.NotNil(t, ojp.FindProfile(nodeID))
}

func TestOJPGetProfiles(t *testing.T) {
	ojp := OneJoinerPopulation{}
	require.Len(t, ojp.GetProfiles(), 0)
}

func TestOJPGetLocalProfile(t *testing.T) {
	ojp := OneJoinerPopulation{localNode: updatableSlot{NodeProfileSlot: NodeProfileSlot{}}}
	require.NotNil(t, ojp.GetLocalProfile())
}
