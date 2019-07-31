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

	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"

	"github.com/stretchr/testify/require"
)

func TestNewEvictedPopulation(t *testing.T) {
	require.Len(t, newEvictedPopulation(nil, 0).profiles, 0)

	require.Len(t, newEvictedPopulation(make([]*updatableSlot, 0), 0).profiles, 0)

	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	evicts := []*updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	ep := newEvictedPopulation(evicts, 0)
	require.Equal(t, 1, ep.GetCount())
}

func TestEPString(t *testing.T) {
	ep := evictedPopulation{}
	require.Equal(t, "[]", ep.String())

	sp1 := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp1.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	sp2 := profiles.NewStaticProfileMock(t)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	evicts := []*updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp1}},
		{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp2}}}
	ep = newEvictedPopulation(evicts, 0)
	require.NotEmpty(t, ep.String())

	ep = newEvictedPopulation(evicts, 1)
	require.NotEmpty(t, ep.String())

	for i := insolar.ShortNodeID(10); i < 60; i++ {
		ep.profiles[i] = &evictedSlot{}
	}
	require.NotEmpty(t, ep.String())
}

func TestEPIsValid(t *testing.T) {
	ep := evictedPopulation{}
	require.False(t, ep.IsValid())

	ep.detectedErrors = 1
	require.True(t, ep.IsValid())
}

func TestGetDetectedErrors(t *testing.T) {
	ep := evictedPopulation{}
	require.Zero(t, ep.GetDetectedErrors())

	derr := census.RecoverableErrorTypes(1)
	ep.detectedErrors = derr
	require.Equal(t, derr, ep.detectedErrors)
}

func TestEPFindProfile(t *testing.T) {
	sp1 := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp1.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	sp2 := profiles.NewStaticProfileMock(t)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	evicts := []*updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp1}},
		{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp2}}}
	ep := newEvictedPopulation(evicts, 0)
	en := ep.FindProfile(nodeID)
	require.NotNil(t, en)

	en = ep.FindProfile(2)
	require.Nil(t, en)
}

func TestEPGetCount(t *testing.T) {
	sp1 := profiles.NewStaticProfileMock(t)
	sp1.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	sp2 := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return *(&nodeID) })
	evicts := []*updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp1}},
		{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp2}}}
	ep := newEvictedPopulation(evicts, 0)
	require.Equal(t, 1, ep.GetCount())

	nodeID = 1
	ep = newEvictedPopulation(evicts, 0)
	require.Equal(t, 2, ep.GetCount())
}

func TestEPGetProfiles(t *testing.T) {
	sp1 := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp1.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	sp2 := profiles.NewStaticProfileMock(t)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	evicts := []*updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp1}},
		{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp2}}}
	ep := newEvictedPopulation(evicts, 0)
	require.Len(t, ep.GetProfiles(), len(evicts))
}

func TestEPGetNodeID(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	es := evictedSlot{StaticProfile: sp}
	require.Equal(t, nodeID, es.GetNodeID())
}

func TestEPGetStatic(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	es := evictedSlot{StaticProfile: sp}
	require.Equal(t, sp, es.GetStatic())
}

func TestEPGetSignatureVerifier(t *testing.T) {
	sv := cryptkit.NewSignatureVerifierMock(t)
	es := evictedSlot{sf: sv}
	require.Equal(t, sv, es.GetSignatureVerifier())
}

func TestEPGetOpMode(t *testing.T) {
	opMode := member.ModeSuspected
	es := evictedSlot{mode: opMode}
	require.Equal(t, opMode, es.GetOpMode())
}

func TestEPGetLeaveReason(t *testing.T) {
	opMode := member.ModeEvictedGracefully
	leaveReason := uint32(1)
	es := evictedSlot{mode: opMode, leaveReason: leaveReason}
	require.Equal(t, leaveReason, es.GetLeaveReason())

	es.mode = member.ModeSuspected
	require.Zero(t, es.GetLeaveReason())
}
