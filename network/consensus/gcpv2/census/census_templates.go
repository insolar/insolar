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
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

var _ localActiveCensus = &PrimingCensusTemplate{}

type PrimingCensusTemplate struct {
	chronicles *localChronicles
	online     copyOnlinePopulationTo
	pd         common.PulseData

	registries     VersionedRegistries
	profileFactory common2.NodeProfileFactory
}

func (c *PrimingCensusTemplate) GetProfileFactory() common2.NodeProfileFactory {
	return c.profileFactory
}

func (c *PrimingCensusTemplate) setVersionedRegistries(vr VersionedRegistries) {
	if vr == nil {
		panic("versioned registries are nil")
	}
	c.registries = vr
}

func (c *PrimingCensusTemplate) getVersionedRegistries() VersionedRegistries {
	return c.registries
}

func NewPrimingCensus(population copyOnlinePopulationTo, pf common2.NodeProfileFactory, registries VersionedRegistries) *PrimingCensusTemplate {
	r := &PrimingCensusTemplate{
		registries:     registries,
		online:         population,
		profileFactory: pf,
		pd:             registries.GetVersionPulseData(),
	}
	return r
}

func (c *PrimingCensusTemplate) SetAsActiveTo(chronicles LocalConsensusChronicles) {
	if c.chronicles != nil {
		panic("illegal state")
	}
	lc := chronicles.(*localChronicles)
	c.chronicles = lc
	lc.makeActive(nil, c)
}

func (*PrimingCensusTemplate) GetCensusState() State {
	return PrimingCensus
}

func (c *PrimingCensusTemplate) GetExpectedPulseNumber() common.PulseNumber {
	switch {
	case c.pd.IsEmpty():
		return common.UnknownPulseNumber
	case c.pd.IsExpectedPulse():
		return c.pd.GetPulseNumber()
	}
	return c.pd.GetNextPulseNumber()
}

func (c *PrimingCensusTemplate) GetPulseNumber() common.PulseNumber {
	return c.pd.GetPulseNumber()
}

func (c *PrimingCensusTemplate) GetPulseData() common.PulseData {
	return c.pd
}

func (*PrimingCensusTemplate) GetGlobulaStateHash() common2.GlobulaStateHash {
	return nil
}

func (*PrimingCensusTemplate) GetCloudStateHash() common2.CloudStateHash {
	return nil
}

func (c *PrimingCensusTemplate) GetOnlinePopulation() OnlinePopulation {
	return c.online
}

func (c *PrimingCensusTemplate) GetOfflinePopulation() OfflinePopulation {
	return c.registries.GetOfflinePopulation()
}

func (c *PrimingCensusTemplate) IsActive() bool {
	return c.chronicles.GetActiveCensus() == c
}

func (c *PrimingCensusTemplate) GetMisbehaviorRegistry() MisbehaviorRegistry {
	return c.registries.GetMisbehaviorRegistry()
}

func (c *PrimingCensusTemplate) GetMandateRegistry() MandateRegistry {
	return c.registries.GetMandateRegistry()
}

func (c *PrimingCensusTemplate) CreateBuilder(pn common.PulseNumber) Builder {
	return newLocalCensusBuilder(c.chronicles, pn, c.online)
}

var _ ActiveCensus = &ActiveCensusTemplate{}

type ActiveCensusTemplate struct {
	PrimingCensusTemplate
	gsh common2.GlobulaStateHash
	csh common2.CloudStateHash
}

func (c *ActiveCensusTemplate) GetExpectedPulseNumber() common.PulseNumber {
	return c.pd.GetNextPulseNumber()
}

func (*ActiveCensusTemplate) GetCensusState() State {
	return SealedCensus
}

func (c *ActiveCensusTemplate) GetPulseNumber() common.PulseNumber {
	return c.pd.PulseNumber
}

func (c *ActiveCensusTemplate) GetPulseData() common.PulseData {
	return c.pd
}

func (c *ActiveCensusTemplate) GetGlobulaStateHash() common2.GlobulaStateHash {
	return c.gsh
}

func (c *ActiveCensusTemplate) GetCloudStateHash() common2.CloudStateHash {
	return c.csh
}

var _ ExpectedCensus = &ExpectedCensusTemplate{}

type ExpectedCensusTemplate struct {
	chronicles *localChronicles
	online     copyOnlinePopulationTo
	prev       ActiveCensus
	gsh        common2.GlobulaStateHash
	csh        common2.CloudStateHash
	pn         common.PulseNumber
}

func (c *ExpectedCensusTemplate) GetExpectedPulseNumber() common.PulseNumber {
	return c.pn
}

func (c *ExpectedCensusTemplate) GetCensusState() State {
	return BuiltCensus
}

func (c *ExpectedCensusTemplate) GetPulseNumber() common.PulseNumber {
	return c.pn
}

func (c *ExpectedCensusTemplate) GetGlobulaStateHash() common2.GlobulaStateHash {
	return c.gsh
}

func (c *ExpectedCensusTemplate) GetCloudStateHash() common2.CloudStateHash {
	return c.csh
}

func (c *ExpectedCensusTemplate) GetOnlinePopulation() OnlinePopulation {
	return c.online
}

func (c *ExpectedCensusTemplate) GetOfflinePopulation() OfflinePopulation {
	// TODO Should be provided via relevant builder
	return c.prev.GetOfflinePopulation()
}

func (c *ExpectedCensusTemplate) GetMisbehaviorRegistry() MisbehaviorRegistry {
	// TODO Should be provided via relevant builder
	return c.prev.GetMisbehaviorRegistry()
}

func (c *ExpectedCensusTemplate) GetMandateRegistry() MandateRegistry {
	// TODO Should be provided via relevant builder
	return c.prev.GetMandateRegistry()
}

func (c *ExpectedCensusTemplate) CreateBuilder(pn common.PulseNumber) Builder {
	return newLocalCensusBuilder(c.chronicles, pn, c.online)
}

func (c *ExpectedCensusTemplate) GetPrevious() ActiveCensus {
	return c.prev
}

func (c *ExpectedCensusTemplate) MakeActive(pd common.PulseData) ActiveCensus {

	a := ActiveCensusTemplate{
		PrimingCensusTemplate: PrimingCensusTemplate{
			chronicles: c.chronicles,
			online:     c.online,
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
