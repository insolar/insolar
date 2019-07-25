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
	vf.GetSignatureVerifierWithPKSMock.Set(func(cryptkit.PublicKeyStore) cryptkit.SignatureVerifier { return sv })
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

func TestFindProfile(t *testing.T) {
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
