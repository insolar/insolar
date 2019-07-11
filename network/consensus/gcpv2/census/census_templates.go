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

package census

import (
	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api_2"
)

var _ localActiveCensus = &PrimingCensusTemplate{}

type copyToOnlinePopulation interface {
	copyToPopulation
	api_2.OnlinePopulation
}

type PrimingCensusTemplate struct {
	chronicles *localChronicles
	online     copyToOnlinePopulation
	evicted    api_2.EvictedPopulation
	pd         pulse_data.PulseData

	registries     api_2.VersionedRegistries
	profileFactory api.NodeProfileFactory
}

func (c *PrimingCensusTemplate) GetProfileFactory(ksf cryptography_containers.KeyStoreFactory) api.NodeProfileFactory {
	return c.profileFactory
}

func (c *PrimingCensusTemplate) setVersionedRegistries(vr api_2.VersionedRegistries) {
	if vr == nil {
		panic("versioned registries are nil")
	}
	c.registries = vr
}

func (c *PrimingCensusTemplate) getVersionedRegistries() api_2.VersionedRegistries {
	return c.registries
}

func NewPrimingCensus(population copyToOnlinePopulation, pf api.NodeProfileFactory, registries api_2.VersionedRegistries) *PrimingCensusTemplate {
	//TODO HACK - ugly sorting impl to establish initial node ordering
	dp := NewDynamicPopulation(population)
	sortedPopulation := ManyNodePopulation{}
	sortedPopulation.makeCopyOfMapAndSort(dp.slotByID, dp.local, lessForNodeProfile)

	r := &PrimingCensusTemplate{
		registries:     registries,
		online:         &sortedPopulation,
		evicted:        &evictedPopulation{},
		profileFactory: pf,
		pd:             registries.GetVersionPulseData(),
	}
	return r
}

func nodeProfileOrdering(np api.NodeProfile) (api.NodePrimaryRole, api.MemberPower, common.ShortNodeID) {
	p := np.GetDeclaredPower()
	r := np.GetPrimaryRole()
	if p == 0 || !np.GetOpMode().IsPowerful() {
		return api.PrimaryRoleInactive, 0, np.GetShortNodeID()
	}
	return r, p, np.GetShortNodeID()
}

func lessForNodeProfile(c api.NodeProfile, o api.NodeProfile) bool {
	cR, cP, cI := nodeProfileOrdering(c)
	oR, oP, oI := nodeProfileOrdering(o)

	/* Reversed order */
	if cR < oR {
		return false
	} else if cR > oR {
		return true
	}

	if cP < oP {
		return true
	} else if cP > oP {
		return false
	}

	return cI < oI
}

func (c *PrimingCensusTemplate) SetAsActiveTo(chronicles LocalConsensusChronicles) {
	if c.chronicles != nil {
		panic("illegal state")
	}
	lc := chronicles.(*localChronicles)
	c.chronicles = lc
	lc.makeActive(nil, c)
}

func (*PrimingCensusTemplate) GetCensusState() api_2.State {
	return api_2.PrimingCensus
}

func (c *PrimingCensusTemplate) GetExpectedPulseNumber() pulse_data.PulseNumber {
	switch {
	case c.pd.IsEmpty():
		return pulse_data.UnknownPulseNumber
	case c.pd.IsExpectedPulse():
		return c.pd.GetPulseNumber()
	}
	return c.pd.GetNextPulseNumber()
}

func (c *PrimingCensusTemplate) GetPulseNumber() pulse_data.PulseNumber {
	return c.pd.GetPulseNumber()
}

func (c *PrimingCensusTemplate) GetPulseData() pulse_data.PulseData {
	return c.pd
}

func (*PrimingCensusTemplate) GetGlobulaStateHash() api.GlobulaStateHash {
	return nil
}

func (*PrimingCensusTemplate) GetCloudStateHash() api.CloudStateHash {
	return nil
}

func (c *PrimingCensusTemplate) GetOnlinePopulation() api_2.OnlinePopulation {
	return c.online
}

func (c *PrimingCensusTemplate) GetEvictedPopulation() api_2.EvictedPopulation {
	return c.evicted
}

func (c *PrimingCensusTemplate) GetOfflinePopulation() api_2.OfflinePopulation {
	return c.registries.GetOfflinePopulation()
}

func (c *PrimingCensusTemplate) IsActive() bool {
	return c.chronicles.GetActiveCensus() == c
}

func (c *PrimingCensusTemplate) GetMisbehaviorRegistry() api_2.MisbehaviorRegistry {
	return c.registries.GetMisbehaviorRegistry()
}

func (c *PrimingCensusTemplate) GetMandateRegistry() api_2.MandateRegistry {
	return c.registries.GetMandateRegistry()
}

func (c *PrimingCensusTemplate) CreateBuilder(pn pulse_data.PulseNumber, fullCopy bool) api_2.Builder {
	return newLocalCensusBuilder(c.chronicles, pn, c.online, fullCopy)
}

var _ api_2.ActiveCensus = &ActiveCensusTemplate{}

type ActiveCensusTemplate struct {
	PrimingCensusTemplate
	gsh api.GlobulaStateHash
	csh api.CloudStateHash
}

func (c *ActiveCensusTemplate) GetExpectedPulseNumber() pulse_data.PulseNumber {
	return c.pd.GetNextPulseNumber()
}

func (*ActiveCensusTemplate) GetCensusState() api_2.State {
	return api_2.SealedCensus
}

func (c *ActiveCensusTemplate) GetPulseNumber() pulse_data.PulseNumber {
	return c.pd.PulseNumber
}

func (c *ActiveCensusTemplate) GetPulseData() pulse_data.PulseData {
	return c.pd
}

func (c *ActiveCensusTemplate) GetGlobulaStateHash() api.GlobulaStateHash {
	return c.gsh
}

func (c *ActiveCensusTemplate) GetCloudStateHash() api.CloudStateHash {
	return c.csh
}

var _ api_2.ExpectedCensus = &ExpectedCensusTemplate{}

type ExpectedCensusTemplate struct {
	chronicles *localChronicles
	online     copyToOnlinePopulation
	evicted    api_2.EvictedPopulation
	prev       api_2.ActiveCensus
	gsh        api.GlobulaStateHash
	csh        api.CloudStateHash
	pn         pulse_data.PulseNumber
}

func (c *ExpectedCensusTemplate) GetEvictedPopulation() api_2.EvictedPopulation {
	return c.evicted
}

func (c *ExpectedCensusTemplate) GetExpectedPulseNumber() pulse_data.PulseNumber {
	return c.pn
}

func (c *ExpectedCensusTemplate) GetCensusState() api_2.State {
	return api_2.BuiltCensus
}

func (c *ExpectedCensusTemplate) GetPulseNumber() pulse_data.PulseNumber {
	return c.pn
}

func (c *ExpectedCensusTemplate) GetGlobulaStateHash() api.GlobulaStateHash {
	return c.gsh
}

func (c *ExpectedCensusTemplate) GetCloudStateHash() api.CloudStateHash {
	return c.csh
}

func (c *ExpectedCensusTemplate) GetOnlinePopulation() api_2.OnlinePopulation {
	return c.online
}

func (c *ExpectedCensusTemplate) GetOfflinePopulation() api_2.OfflinePopulation {
	// TODO Should be provided via relevant builder
	return c.prev.GetOfflinePopulation()
}

func (c *ExpectedCensusTemplate) GetMisbehaviorRegistry() api_2.MisbehaviorRegistry {
	// TODO Should be provided via relevant builder
	return c.prev.GetMisbehaviorRegistry()
}

func (c *ExpectedCensusTemplate) GetMandateRegistry() api_2.MandateRegistry {
	// TODO Should be provided via relevant builder
	return c.prev.GetMandateRegistry()
}

func (c *ExpectedCensusTemplate) CreateBuilder(pn pulse_data.PulseNumber, fullCopy bool) api_2.Builder {
	return newLocalCensusBuilder(c.chronicles, pn, c.online, fullCopy)
}

func (c *ExpectedCensusTemplate) GetPrevious() api_2.ActiveCensus {
	return c.prev
}

func (c *ExpectedCensusTemplate) MakeActive(pd pulse_data.PulseData) api_2.ActiveCensus {

	a := ActiveCensusTemplate{
		PrimingCensusTemplate: PrimingCensusTemplate{
			chronicles: c.chronicles,
			online:     c.online,
			evicted:    c.evicted,
			pd:         pd,
		},
		gsh: c.gsh,
		csh: c.csh,
	} // make copy for thread-safe access

	c.chronicles.makeActive(c, &a)
	return &a
}

func (c *ExpectedCensusTemplate) IsActive() bool {
	return false
}
