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
	"sync"

	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

func newLocalCensusBuilder(chronicles *localChronicles, pn common.PulseNumber,
	population copyOnlinePopulationTo) *LocalCensusBuilder {

	r := &LocalCensusBuilder{chronicles: chronicles, pulseNumber: pn, population: NewDynamicPopulation(population)}
	r.populationBuilder.census = r
	r.populationView.census = r
	return r
}

var _ Builder = &LocalCensusBuilder{}

type LocalCensusBuilder struct {
	mutex             sync.RWMutex
	chronicles        *localChronicles
	pulseNumber       common.PulseNumber
	population        DynamicPopulation
	state             State
	populationView    LockedPopulation
	populationBuilder DynamicPopulationBuilder
	gsh               common2.GlobulaStateHash
	csh               common2.CloudStateHash

	// content WorkingCensusTemplate
}

func (c *LocalCensusBuilder) GetCensusState() State {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.state
}

func (c *LocalCensusBuilder) GetPulseNumber() common.PulseNumber {
	return c.pulseNumber
}

func (c *LocalCensusBuilder) GetGlobulaStateHash() common2.GlobulaStateHash {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.gsh
}

func (c *LocalCensusBuilder) SetGlobulaStateHash(gsh common2.GlobulaStateHash) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state.IsSealed() {
		panic("illegal state")
	}

	c.gsh = gsh
}

func (c *LocalCensusBuilder) SealCensus() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state.IsSealed() {
		return
	}
	if c.gsh == nil {
		panic("illegal state: GSH is nil")
	}
	c.state = SealedCensus
}

func (c *LocalCensusBuilder) IsSealed() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.state.IsSealed()
}

func (c *LocalCensusBuilder) GetOnlinePopulationView() OnlinePopulation {
	return &c.populationView
}

func (c *LocalCensusBuilder) GetOnlinePopulationBuilder() OnlinePopulationBuilder {
	return &c.populationBuilder
}

func (c *LocalCensusBuilder) build(csh common2.CloudStateHash) copyOnlinePopulationTo {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if csh == nil {
		panic("illegal state: CSH is nil")
	}

	if !c.state.IsSealed() {
		panic("illegal state: not sealed")
	}

	if c.state.IsBuilt() {
		panic("illegal state: was built")
	}
	c.state = BuiltCensus
	c.csh = csh

	r := c.population.CopyAndSortDefault()
	return &r
}

func (c *LocalCensusBuilder) BuildAndMakeExpected(csh common2.CloudStateHash) ExpectedCensus {
	pop := c.build(csh)

	r := &ExpectedCensusTemplate{
		chronicles: c.chronicles,
		prev:       c.chronicles.active,
		csh:        c.csh,
		gsh:        c.gsh,
		pn:         c.pulseNumber,
		online:     pop,
	}

	c.chronicles.makeExpected(r)
	return r
}

var _ OnlinePopulation = &LockedPopulation{}

type LockedPopulation struct {
	census *LocalCensusBuilder
}

func (c *LockedPopulation) FindProfile(nodeID common.ShortNodeID) common2.NodeProfile {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.FindProfile(nodeID)
}

func (c *LockedPopulation) GetProfiles() []common2.NodeProfile {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.GetProfiles()
}

func (c *LockedPopulation) GetCount() int {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.GetCount()
}

func (c *LockedPopulation) GetLocalProfile() common2.LocalNodeProfile {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.GetLocalProfile()
}

var _ OnlinePopulationBuilder = &DynamicPopulationBuilder{}

type DynamicPopulationBuilder struct {
	census *LocalCensusBuilder
}

func (c *DynamicPopulationBuilder) GetUnorderedProfiles() []common2.UpdatableNodeProfile {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.GetUnorderedProfiles()
}

func (c *DynamicPopulationBuilder) GetCount() int {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.GetCount()
}

func (c *DynamicPopulationBuilder) GetLocalProfile() common2.UpdatableNodeProfile {
	return c.FindProfile(c.census.population.GetLocalProfile().GetShortNodeID())
}

func (c *DynamicPopulationBuilder) FindProfile(nodeID common.ShortNodeID) common2.UpdatableNodeProfile {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.FindUpdatableProfile(nodeID)
}

func (c *DynamicPopulationBuilder) AddJoinerProfile(intro common2.NodeIntroProfile) common2.UpdatableNodeProfile {
	c.census.mutex.Lock()
	defer c.census.mutex.Unlock()

	if c.census.state.IsSealed() {
		panic("illegal state")
	}
	return c.census.population.AddJoinerProfile(intro)
}

func (c *DynamicPopulationBuilder) RemoveProfile(nodeID common.ShortNodeID) {
	c.census.mutex.Lock()
	defer c.census.mutex.Unlock()

	if c.census.state.IsSealed() {
		panic("illegal state")
	}
	c.census.population.RemoveProfile(nodeID)
}
