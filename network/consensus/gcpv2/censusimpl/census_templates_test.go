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
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

func TestPCTGetProfileFactory(t *testing.T) {
	pct := PrimingCensusTemplate{CensusTemplate: CensusTemplate{chronicles: &localChronicles{}}}
	pf := pct.GetProfileFactory(nil)
	require.Nil(t, pf)
}

func TestSetVersionedRegistries(t *testing.T) {
	pct := PrimingCensusTemplate{}
	require.Panics(t, func() { pct.setVersionedRegistries(nil) })

	vr := census.NewVersionedRegistriesMock(t)
	pct.setVersionedRegistries(vr)
	require.Equal(t, vr, pct.registries)
}

func TestGetVersionedRegistries(t *testing.T) {
	pct := PrimingCensusTemplate{}
	vr := census.NewVersionedRegistriesMock(t)
	pct.setVersionedRegistries(vr)
	require.Equal(t, vr, pct.getVersionedRegistries())
}

func TestNewPrimingCensusForJoiner(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	pks := cryptkit.NewPublicKeyStoreMock(t)
	sp.GetPublicKeyStoreMock.Set(func() cryptkit.PublicKeyStore { return pks })
	registries := census.NewVersionedRegistriesMock(t)
	pn := pulse.Number(1)
	registries.GetVersionPulseDataMock.Set(func() pulse.Data { return pulse.Data{PulseNumber: pn} })
	vf := cryptkit.NewSignatureVerifierFactoryMock(t)
	sv := cryptkit.NewSignatureVerifierMock(t)
	vf.CreateSignatureVerifierWithPKSMock.Set(func(cryptkit.PublicKeyStore) cryptkit.SignatureVerifier { return sv })
	pcj := NewPrimingCensusForJoiner(sp, registries, vf, true)

	// TODO: investigate
	// require.Equal(t, pn, pcj.GetPulseNumber())
	require.EqualValues(t, 0, pcj.GetPulseNumber())
}

func TestNewPrimingCensus(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return *(&nodeID) })
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleNeutral })
	pks := cryptkit.NewPublicKeyStoreMock(t)
	sp.GetPublicKeyStoreMock.Set(func() cryptkit.PublicKeyStore { return pks })
	registries := census.NewVersionedRegistriesMock(t)
	pn := pulse.Number(1)
	registries.GetVersionPulseDataMock.Set(func() pulse.Data { return pulse.Data{PulseNumber: pn} })
	var sps []profiles.StaticProfile
	sps = append(sps, sp)
	vf := cryptkit.NewSignatureVerifierFactoryMock(t)
	sv := cryptkit.NewSignatureVerifierMock(t)
	vf.CreateSignatureVerifierWithPKSMock.Set(func(cryptkit.PublicKeyStore) cryptkit.SignatureVerifier { return sv })
	require.Panics(t, func() { NewPrimingCensus(nil, sp, registries, vf, true) })

	require.Panics(t, func() { NewPrimingCensus(sps, sp, registries, vf, true) })
	nodeID = 1
	pc := NewPrimingCensus(sps, sp, registries, vf, true)

	// TODO: investigate
	// require.Equal(t, pn, pc.GetPulseNumber())
	require.EqualValues(t, 0, pc.GetPulseNumber())

	require.Panics(t, func() { NewPrimingCensus(nil, sp, registries, vf, true) })
}

func TestSetAsActiveTo(t *testing.T) {
	pct := PrimingCensusTemplate{}
	chronicles := &localChronicles{}
	require.Panics(t, func() { pct.SetAsActiveTo(chronicles) })

	pct.chronicles = nil
	pct.registries = census.NewVersionedRegistriesMock(t)
	require.Panics(t, func() { pct.IsActive() })

	pct.SetAsActiveTo(chronicles)
	require.Equal(t, chronicles, pct.chronicles)

	require.True(t, pct.IsActive())

	require.Panics(t, func() { pct.SetAsActiveTo(chronicles) })
}

func TestPCTGetCensusState(t *testing.T) {
	pct := PrimingCensusTemplate{}
	require.Equal(t, census.PrimingCensus, pct.GetCensusState())
}

func TestPCTGetExpectedPulseNumber(t *testing.T) {
	pct := PrimingCensusTemplate{}
	require.Equal(t, pulse.Unknown, pct.GetExpectedPulseNumber())

	pct.pd.PulseNumber = 1
	require.Panics(t, func() { pct.GetExpectedPulseNumber() })

	pct.pd.PulseNumber = pulse.MinTimePulse
	require.Equal(t, pulse.Number(pulse.MinTimePulse), pct.GetExpectedPulseNumber())

	pct.pd.NextPulseDelta = 1
	require.Equal(t, pulse.MinTimePulse+pulse.Number(pct.pd.NextPulseDelta), pct.GetExpectedPulseNumber())
}

func TestPCTMakeExpected(t *testing.T) {
	pd := pulse.NewFirstEphemeralData()
	pct := PrimingCensusTemplate{}
	pct.pd.PulseEpoch = pulse.EphemeralPulseEpoch
	require.Panics(t, func() { pct.BuildCopy(pd, nil, nil).MakeExpected() })

	csh := proofs.NewCloudStateHashMock(t)
	require.Panics(t, func() { pct.BuildCopy(pd, csh, nil).MakeExpected() })

	gsh := proofs.NewGlobulaStateHashMock(t)
	pct.pd.PulseNumber = pulse.MinTimePulse
	pct.pd.NextPulseDelta = 1
	require.Panics(t, func() { pct.BuildCopy(pd, csh, gsh).MakeExpected() })

	pct.chronicles = &localChronicles{}
	next := pd.CreateNextEphemeralPulse()
	r := pct.BuildCopy(next, csh, gsh).MakeExpected()
	require.Equal(t, next.PulseNumber, r.GetPulseNumber())
}

func TestPCTGetPulseNumber(t *testing.T) {
	pct := PrimingCensusTemplate{}
	pn := pulse.Number(1)
	pct.pd.PulseNumber = pn
	require.Equal(t, pn, pct.GetPulseNumber())
}

func TestPCTGetPulseData(t *testing.T) {
	pct := PrimingCensusTemplate{}
	pn := pulse.Number(1)
	pct.pd.PulseNumber = pn
	pd := pct.GetPulseData()
	require.Equal(t, pn, pd.GetPulseNumber())
}

func TestPCTGetGlobulaStateHash(t *testing.T) {
	pct := PrimingCensusTemplate{}
	require.Nil(t, pct.GetGlobulaStateHash())
}

func TestPCTGetCloudStateHash(t *testing.T) {
	pct := PrimingCensusTemplate{}
	require.Nil(t, pct.GetCloudStateHash())
}

func TestPCTString(t *testing.T) {
	pct := PrimingCensusTemplate{}
	require.NotEmpty(t, pct.String())
}

func TestPCTGetOnlinePopulation(t *testing.T) {
	pct := PrimingCensusTemplate{}
	require.Nil(t, pct.GetOnlinePopulation())
}

func TestPCTGetEvictedPopulation(t *testing.T) {
	ep := census.NewEvictedPopulationMock(t)
	pct := PrimingCensusTemplate{CensusTemplate: CensusTemplate{evicted: ep}}
	require.Equal(t, ep, pct.GetEvictedPopulation())
}

func TestPCTGetOfflinePopulation(t *testing.T) {
	registries := census.NewVersionedRegistriesMock(t)
	op := census.NewOfflinePopulationMock(t)
	registries.GetOfflinePopulationMock.Set(func() census.OfflinePopulation { return op })
	pct := PrimingCensusTemplate{CensusTemplate: CensusTemplate{registries: registries}}
	require.Equal(t, op, pct.GetOfflinePopulation())
}

func TestPCTIsActive(t *testing.T) {
	chronicles := &localChronicles{}
	pct := PrimingCensusTemplate{CensusTemplate: CensusTemplate{chronicles: chronicles}}
	require.False(t, pct.IsActive())

	pct.chronicles = nil
	registries := census.NewVersionedRegistriesMock(t)
	pct.registries = registries
	pct.SetAsActiveTo(chronicles)
	require.True(t, pct.IsActive())
}

func TestPCTGetMisbehaviorRegistry(t *testing.T) {
	registries := census.NewVersionedRegistriesMock(t)
	mr := census.NewMisbehaviorRegistryMock(t)
	registries.GetMisbehaviorRegistryMock.Set(func() census.MisbehaviorRegistry { return mr })
	pct := PrimingCensusTemplate{CensusTemplate: CensusTemplate{registries: registries}}
	require.Equal(t, mr, pct.GetMisbehaviorRegistry())
}

func TestPCTGetMandateRegistry(t *testing.T) {
	registries := census.NewVersionedRegistriesMock(t)
	mr := census.NewMandateRegistryMock(t)
	registries.GetMandateRegistryMock.Set(func() census.MandateRegistry { return mr })
	pct := PrimingCensusTemplate{CensusTemplate: CensusTemplate{registries: registries}}
	require.Equal(t, mr, pct.GetMandateRegistry())
}

func TestPCTCreateBuilder(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(65537)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}},
		slots: []updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}}

	pct := PrimingCensusTemplate{CensusTemplate: CensusTemplate{chronicles: chronicles, online: population}}
	builder := pct.CreateBuilder(context.Background(), pn)
	require.Equal(t, pn, builder.GetPulseNumber())
}

func TestCTGetNearestPulseData(t *testing.T) {
	ct := CensusTemplate{}
	pulseNumber := pulse.Number(1)
	pd := pulse.Data{PulseNumber: pulseNumber}
	ct.pd = pd
	b, data := ct.GetNearestPulseData()
	require.True(t, b)

	require.Equal(t, pd, data)
}

func TestCTGetProfileFactory(t *testing.T) {
	ksf := cryptkit.NewKeyStoreFactoryMock(t)
	ct := CensusTemplate{}
	ct.chronicles = &localChronicles{}
	f := profiles.NewFactoryMock(t)
	ct.chronicles.profileFactory = f
	require.Equal(t, f, ct.GetProfileFactory(ksf))
}

func TestCTSetVersionedRegistries(t *testing.T) {
	ct := CensusTemplate{}
	require.Panics(t, func() { ct.setVersionedRegistries(nil) })

	vr := census.NewVersionedRegistriesMock(t)
	ct.setVersionedRegistries(vr)
	require.Equal(t, vr, ct.registries)
}

func TestCTGetVersionedRegistries(t *testing.T) {
	ct := CensusTemplate{}
	vr := census.NewVersionedRegistriesMock(t)
	ct.setVersionedRegistries(vr)
	require.Equal(t, vr, ct.getVersionedRegistries())
}

func TestGetOnlinePopulation(t *testing.T) {
	ct := CensusTemplate{}
	population := &ManyNodePopulation{}
	ct.online = population
	require.Equal(t, population, ct.GetOnlinePopulation())
}

func TestCTGetEvictedPopulation(t *testing.T) {
	ep := census.NewEvictedPopulationMock(t)
	ct := CensusTemplate{evicted: ep}
	require.Equal(t, ep, ct.GetEvictedPopulation())
}

func TestCTGetOfflinePopulation(t *testing.T) {
	offPop := census.NewOfflinePopulationMock(t)
	registries := census.NewVersionedRegistriesMock(t)
	registries.GetOfflinePopulationMock.Set(func() census.OfflinePopulation { return offPop })
	ct := CensusTemplate{registries: registries}
	require.Equal(t, offPop, ct.GetOfflinePopulation())
}

func TestCTGetMisbehaviorRegistry(t *testing.T) {
	registries := census.NewVersionedRegistriesMock(t)
	mr := census.NewMisbehaviorRegistryMock(t)
	registries.GetMisbehaviorRegistryMock.Set(func() census.MisbehaviorRegistry { return mr })
	ct := CensusTemplate{registries: registries}
	require.Equal(t, mr, ct.GetMisbehaviorRegistry())
}

func TestCTGetMandateRegistry(t *testing.T) {
	registries := census.NewVersionedRegistriesMock(t)
	mr := census.NewMandateRegistryMock(t)
	registries.GetMandateRegistryMock.Set(func() census.MandateRegistry { return mr })
	ct := CensusTemplate{registries: registries}
	require.Equal(t, mr, ct.GetMandateRegistry())
}

func TestCTCreateBuilder(t *testing.T) {
	t.Skip("merge")
	// 	chronicles := &localChronicles{}
	// 	pn := pulse.Number(1)
	// 	sp := profiles.NewStaticProfileMock(t)
	// 	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	// 	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}},
	// 		slots: []updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}}
	//
	// 	ct := CensusTemplate{chronicles: chronicles, online: population}
	// 	builder := ct.CreateBuilder(context.Background(), pn)
	// 	require.Equal(t, pn, builder.GetPulseNumber())
}

func TestCTString(t *testing.T) {
	ct := CensusTemplate{}
	require.NotEmpty(t, ct.String())
}

func TestACTIsActive(t *testing.T) {
	act := ActiveCensusTemplate{}
	chronicles := &localChronicles{}
	act.chronicles = chronicles
	require.False(t, act.IsActive())

	act.chronicles.active = &act
	require.True(t, act.IsActive())
}

func TestACTGetExpectedPulseNumber(t *testing.T) {
	act := ActiveCensusTemplate{}
	act.pd.PulseNumber = pulse.MinTimePulse
	require.Panics(t, func() { act.GetExpectedPulseNumber() })

	act.pd.NextPulseDelta = 1
	require.Equal(t, pulse.MinTimePulse+pulse.Number(act.pd.NextPulseDelta), act.GetExpectedPulseNumber())
}

func TestACTGetCensusState(t *testing.T) {
	act := ActiveCensusTemplate{}
	require.Equal(t, census.SealedCensus, act.GetCensusState())
}

func TestACTGetPulseNumber(t *testing.T) {
	act := ActiveCensusTemplate{}
	pn := pulse.Number(1)
	act.pd.PulseNumber = pn
	require.Equal(t, pn, act.GetPulseNumber())
}

func TestACTGetPulseData(t *testing.T) {
	act := ActiveCensusTemplate{}
	pn := pulse.Number(1)
	act.pd.PulseNumber = pn
	pd := act.GetPulseData()
	require.Equal(t, pn, pd.GetPulseNumber())
}

func TestACTGetGlobulaStateHash(t *testing.T) {
	act := ActiveCensusTemplate{}
	require.Nil(t, act.GetGlobulaStateHash())
}

func TestACTGetCloudStateHash(t *testing.T) {
	act := ActiveCensusTemplate{}
	require.Nil(t, act.GetCloudStateHash())
}

func TestACTString(t *testing.T) {
	act := ActiveCensusTemplate{}
	chronicles := &localChronicles{}
	act.chronicles = chronicles
	require.NotEmpty(t, act.String())

	act.activeRef = &act
	act.chronicles.active = &act
	require.NotEmpty(t, act.String())
}

func TestECTGetNearestPulseData(t *testing.T) {
	ect := ExpectedCensusTemplate{}
	act := census.NewActiveMock(t)
	pulseNumber := pulse.Number(1)
	pd := pulse.Data{PulseNumber: pulseNumber}
	act.GetPulseDataMock.Set(func() pulse.Data { return pd })
	ect.prev = act
	b, data := ect.GetNearestPulseData()
	require.False(t, b)

	require.Equal(t, pulseNumber, data.PulseNumber)
}

func TestECTGetProfileFactory(t *testing.T) {
	ect := ExpectedCensusTemplate{chronicles: &localChronicles{}}
	require.Nil(t, ect.GetProfileFactory(cryptkit.NewKeyStoreFactoryMock(t)))
}

func TestECTGetEvictedPopulation(t *testing.T) {
	ep := census.NewEvictedPopulationMock(t)
	ect := ExpectedCensusTemplate{evicted: ep}
	require.Equal(t, ep, ect.GetEvictedPopulation())
}

func TestECTGetExpectedPulseNumber(t *testing.T) {
	pn := pulse.Number(1)
	ect := ExpectedCensusTemplate{pn: pn}
	require.Equal(t, pn, ect.GetExpectedPulseNumber())
}

func TestECTGetCensusState(t *testing.T) {
	ect := ExpectedCensusTemplate{}
	require.Equal(t, census.CompleteCensus, ect.GetCensusState())
}

func TestECTGetPulseNumber(t *testing.T) {
	ect := ExpectedCensusTemplate{}
	pn := pulse.Number(1)
	ect.pn = pn
	require.Equal(t, pn, ect.GetPulseNumber())
}

func TestECTGetGlobulaStateHash(t *testing.T) {
	ect := ExpectedCensusTemplate{}
	require.Nil(t, ect.GetGlobulaStateHash())
}

func TestECTGetCloudStateHash(t *testing.T) {
	ect := ExpectedCensusTemplate{}
	require.Nil(t, ect.GetCloudStateHash())
}

func TestECTGetOnlinePopulation(t *testing.T) {
	ect := ExpectedCensusTemplate{}
	require.Nil(t, ect.GetOnlinePopulation())
}

func TestECTGetOfflinePopulation(t *testing.T) {
	act := census.NewActiveMock(t)
	offPop := census.NewOfflinePopulationMock(t)
	act.GetOfflinePopulationMock.Set(func() census.OfflinePopulation { return offPop })
	ect := ExpectedCensusTemplate{prev: act}
	require.Equal(t, offPop, ect.GetOfflinePopulation())
}

func TestECTGetMisbehaviorRegistry(t *testing.T) {
	act := census.NewActiveMock(t)
	mr := census.NewMisbehaviorRegistryMock(t)
	act.GetMisbehaviorRegistryMock.Set(func() census.MisbehaviorRegistry { return mr })
	ect := ExpectedCensusTemplate{prev: act}
	require.Equal(t, mr, ect.GetMisbehaviorRegistry())
}

func TestECTGetMandateRegistry(t *testing.T) {
	act := census.NewActiveMock(t)
	mr := census.NewMandateRegistryMock(t)
	act.GetMandateRegistryMock.Set(func() census.MandateRegistry { return mr })
	ect := ExpectedCensusTemplate{prev: act}
	require.Equal(t, mr, ect.GetMandateRegistry())
}

func TestECTCreateBuilder(t *testing.T) {
	chronicles := &localChronicles{}
	pn := pulse.Number(1)
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	population := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}},
		slots: []updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}}

	ect := ExpectedCensusTemplate{chronicles: chronicles, online: population}
	builder := ect.CreateBuilder(context.Background(), pn)
	require.Equal(t, pn, builder.GetPulseNumber())
}

func TestGetPrevious(t *testing.T) {
	act := census.NewActiveMock(t)
	ect := ExpectedCensusTemplate{prev: act}
	require.Equal(t, act, ect.GetPrevious())
}

func TestECTMakeActive(t *testing.T) {
	chronicles := &localChronicles{}
	ect := ExpectedCensusTemplate{chronicles: chronicles}
	pd := pulse.Data{PulseNumber: pulse.Number(1)}
	require.Panics(t, func() { ect.MakeActive(pd) })

	ect.chronicles.expected = &ect
	require.Panics(t, func() { ect.MakeActive(pd) })

	ect.chronicles = &localChronicles{}
	ect.chronicles.expected = &ect
	ect.chronicles.active = &ActiveCensusTemplate{}
	require.Panics(t, func() { ect.MakeActive(pd) })

	pn := pulse.Number(pulse.MinTimePulse)
	pd.PulseNumber = pn
	require.Panics(t, func() { ect.MakeActive(pd) })

	pd.PulseEpoch = pulse.MinTimePulse
	require.Panics(t, func() { ect.MakeActive(pd) })

	pd.NextPulseDelta = 1
	registries := census.NewVersionedRegistriesMock(t)
	registries.CommitNextPulseMock.Set(func(pulse.Data, census.OnlinePopulation) census.VersionedRegistries { return registries })
	registries.GetNearestValidPulseDataMock.Set(func() (d1 pulse.Data) { return pd })
	act := &ActiveCensusTemplate{}
	act.setVersionedRegistries(registries)
	ect.chronicles.active = act
	a := ect.MakeActive(pd)
	require.Equal(t, pn, a.GetPulseNumber())
}

func TestECTIsActive(t *testing.T) {
	ect := ExpectedCensusTemplate{}
	require.False(t, ect.IsActive())
}

func TestECTString(t *testing.T) {
	ect := ExpectedCensusTemplate{}
	require.NotEmpty(t, ect.String())
}
