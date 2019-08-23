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
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

func TestNewLocalCensusBuilder(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })

	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}},
		slots: []updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	require.Equal(t, pn, lcb.populationBuilder.census.pulseNumber)
}

func TestLCBGetCensusState(t *testing.T) {
	st := census.SealedCensus
	lcb := LocalCensusBuilder{state: st}
	require.Equal(t, st, lcb.GetCensusState())
}

func TestLCBGetPulseNumber(t *testing.T) {
	pn := pulse.Number(1)
	lcb := LocalCensusBuilder{pulseNumber: pn}
	require.Equal(t, pn, lcb.GetPulseNumber())
}

func TestLCBGetGlobulaStateHash(t *testing.T) {
	gsh := proofs.NewGlobulaStateHashMock(t)
	lcb := LocalCensusBuilder{gsh: gsh}
	require.Equal(t, gsh, lcb.GetGlobulaStateHash())
}

func TestSetGlobulaStateHash(t *testing.T) {
	gsh := proofs.NewGlobulaStateHashMock(t)
	lcb := LocalCensusBuilder{state: census.PrimingCensus}
	require.Panics(t, func() { lcb.SetGlobulaStateHash(gsh) })

	lcb.state = census.DraftCensus
	lcb.SetGlobulaStateHash(gsh)
	require.Equal(t, gsh, lcb.gsh)
}

func TestSealCensus(t *testing.T) {
	st := census.PrimingCensus
	lcb := LocalCensusBuilder{state: census.PrimingCensus}
	lcb.SealCensus()
	require.Equal(t, st, lcb.state)

	lcb.state = census.DraftCensus
	require.Panics(t, func() { lcb.SealCensus() })

	lcb.gsh = proofs.NewGlobulaStateHashMock(t)
	lcb.SealCensus()
	require.Equal(t, census.SealedCensus, lcb.state)
}

func TestIsSealed(t *testing.T) {
	st := census.SealedCensus
	lcb := LocalCensusBuilder{state: st}
	require.True(t, lcb.IsSealed())

	lcb.state = census.DraftCensus
	require.False(t, lcb.IsSealed())
}

func TestGetPopulationBuilder(t *testing.T) {
	lcb := LocalCensusBuilder{}
	require.NotPanics(t, func() { lcb.GetPopulationBuilder() })
}

func TestBuild(t *testing.T) {
	lcb := LocalCensusBuilder{ctx: context.Background(), state: census.CompleteCensus}
	csh := proofs.NewCloudStateHashMock(t)
	require.Panics(t, func() { lcb.buildPopulation(true, csh) })

	lcb.state = census.DraftCensus
	require.Panics(t, func() { lcb.buildPopulation(false, csh) })

	lcb.state = census.SealedCensus
	require.Panics(t, func() { lcb.buildPopulation(false, nil) })

	lcb.buildPopulation(false, csh)
	require.Equal(t, census.CompleteCensus, lcb.state)

	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleNeutral })
	spe := profiles.NewStaticProfileExtensionMock(t)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return spe })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	lc := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	lc.state = census.SealedCensus
	lc.buildPopulation(true, csh)
	require.Equal(t, census.CompleteCensus, lc.state)

	require.Equal(t, csh, lc.csh)
}

func TestBuildAndMakeExpected(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleNeutral })
	spe := profiles.NewStaticProfileExtensionMock(t)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return spe })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	lcb.state = census.SealedCensus
	csh := proofs.NewCloudStateHashMock(t)
	ce := lcb.Build(csh).MakeExpected()
	require.Equal(t, csh, lcb.csh)

	require.Equal(t, pn, ce.GetPulseNumber())
}

func TestBuildAndMakeBrokenExpected(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleNeutral })
	spe := profiles.NewStaticProfileExtensionMock(t)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return spe })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	lcb.state = census.SealedCensus
	csh := proofs.NewCloudStateHashMock(t)
	ce := lcb.BuildAsBroken(csh).MakeExpected()
	require.Equal(t, csh, lcb.csh)

	require.Equal(t, pn, ce.GetPulseNumber())
}

func TestLCBMakeExpected(t *testing.T) {
	t.Skip("merge")

	// chronicles := &localChronicles{}
	// pn := pulse.Number(1)
	// sp := profiles.NewStaticProfileMock(t)
	// sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	// sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleNeutral })
	// spe := profiles.NewStaticProfileExtensionMock(t)
	// sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return spe })
	// population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	// lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	// lcb.state = census.SealedCensus
	// csh := proofs.NewCloudStateHashMock(t)
	// pop, evicts := lcb.buildPopulation(true, csh)
	//
	// ce := lcb.makeExpected(pop, evicts)
	// require.Equal(t, csh, lcb.csh)
	//
	// require.Equal(t, pn, ce.GetPulseNumber())
}

func TestDPBRemoveOthers(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)

	dpb := DynamicPopulationBuilder{census: lcb}
	lcb.population.slotByID[1] = nil
	require.Len(t, lcb.population.slotByID, 2)

	dpb.RemoveOthers()
	require.Len(t, lcb.population.slotByID, 1)
}

func TestDPBGetUnorderedProfiles(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)

	dpb := DynamicPopulationBuilder{census: lcb}
	lcb.population.slotByID[1] = nil
	length := 2
	require.Len(t, lcb.population.slotByID, length)

	up := dpb.GetUnorderedProfiles()
	require.Len(t, up, length)
}

func TestDPBGetCount(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)

	dpb := DynamicPopulationBuilder{census: lcb}
	lcb.population.slotByID[1] = nil
	require.Equal(t, 2, dpb.GetCount())
}

func TestDPBGetLocalProfile(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	dpb := DynamicPopulationBuilder{census: lcb}
	require.Zero(t, dpb.GetLocalProfile().GetLeaveReason())
}

func TestDPBFindProfile(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	dpb := DynamicPopulationBuilder{census: lcb}
	require.Zero(t, dpb.FindProfile(0).GetLeaveReason())

	require.Panics(t, func() { dpb.FindProfile(1).GetLeaveReason() })
}

func TestDPBAddProfile(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp1 := profiles.NewStaticProfileMock(t)
	sp1.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp1}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	dpb := DynamicPopulationBuilder{census: lcb}

	sp2 := profiles.NewStaticProfileMock(t)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	sp2.GetStartPowerMock.Set(func() member.Power { return 0 })
	lcb.state = census.SealedCensus
	require.Panics(t, func() { dpb.AddProfile(sp2) })

	lcb.state = census.DraftCensus
	require.Len(t, lcb.population.slotByID, 1)

	dpb.AddProfile(sp2)
	require.Len(t, lcb.population.slotByID, 2)
}

func TestDPBRemoveProfile(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp1 := profiles.NewStaticProfileMock(t)
	sp1.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp1}}}
	lcb := newLocalCensusBuilder(context.Background(), chronicles, pn, population)
	dpb := DynamicPopulationBuilder{census: lcb}

	sp2 := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	sp2.GetStartPowerMock.Set(func() member.Power { return 0 })
	lcb.state = census.DraftCensus
	dpb.AddProfile(sp2)

	lcb.state = census.SealedCensus
	require.Panics(t, func() { dpb.RemoveProfile(nodeID) })

	lcb.state = census.DraftCensus
	require.Len(t, lcb.population.slotByID, 2)

	dpb.RemoveProfile(nodeID)
	require.Len(t, lcb.population.slotByID, 1)
}
