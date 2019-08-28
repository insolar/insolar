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
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

func newLocalCensusBuilder(ctx context.Context, chronicles *localChronicles, pn pulse.Number,
	population copyToPopulation) *LocalCensusBuilder {

	r := &LocalCensusBuilder{chronicles: chronicles, pulseNumber: pn, ctx: ctx}
	r.population = NewDynamicPopulationCopySelf(population)
	r.populationBuilder.census = r
	return r
}

var _ census.Builder = &LocalCensusBuilder{}

type LocalCensusBuilder struct {
	ctx               context.Context
	mutex             sync.RWMutex
	chronicles        *localChronicles
	pulseNumber       pulse.Number
	population        DynamicPopulation
	state             census.State
	populationBuilder DynamicPopulationBuilder
	gsh               proofs.GlobulaStateHash
	csh               proofs.CloudStateHash
}

func (c *LocalCensusBuilder) GetCensusState() census.State {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.state
}

func (c *LocalCensusBuilder) GetPulseNumber() pulse.Number {
	return c.pulseNumber
}

func (c *LocalCensusBuilder) GetGlobulaStateHash() proofs.GlobulaStateHash {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.gsh
}

func (c *LocalCensusBuilder) SetGlobulaStateHash(gsh proofs.GlobulaStateHash) {
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
	c.state = census.SealedCensus
}

func (c *LocalCensusBuilder) IsSealed() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.state.IsSealed()
}

func (c *LocalCensusBuilder) GetPopulationBuilder() census.PopulationBuilder {
	return &c.populationBuilder
}

func (c *LocalCensusBuilder) buildPopulation(markBroken bool, csh proofs.CloudStateHash) (copyToOnlinePopulation, census.EvictedPopulation) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state.IsBuilt() {
		panic("illegal state: was built")
	}

	if !markBroken {
		if !c.state.IsSealed() {
			panic("illegal state: not sealed")
		}

		if csh == nil {
			panic("illegal state: CSH is nil")
		}
	}
	c.csh = csh
	c.state = census.CompleteCensus
	log := inslogger.FromContext(c.ctx)
	pop, evicts := c.population.CopyAndSeparate(markBroken, func(e census.RecoverableErrorTypes, msg string, args ...interface{}) {
		log.Debugf(msg, args...)
	})
	return pop, evicts
}

func (c *LocalCensusBuilder) Build(csh proofs.CloudStateHash) census.Built {
	return c.buildCensus(csh, false)
}

func (c *LocalCensusBuilder) BuildAsBroken(csh proofs.CloudStateHash) census.Built {
	return c.buildCensus(csh, true)
}

func (c *LocalCensusBuilder) buildCensus(csh proofs.CloudStateHash, markBroken bool) census.Built {

	pop, evicts := c.buildPopulation(markBroken, csh)
	return &BuiltCensusTemplate{ExpectedCensusTemplate{
		c.chronicles, pop, evicts, c.chronicles.active, c.csh, c.gsh,
		c.pulseNumber,
	}}
}

var _ census.PopulationBuilder = &DynamicPopulationBuilder{}

type DynamicPopulationBuilder struct {
	census *LocalCensusBuilder
}

func (c *DynamicPopulationBuilder) RemoveOthers() {
	c.census.mutex.Lock()
	defer c.census.mutex.Unlock()

	c.census.population.RemoveOthers()
}

func (c *DynamicPopulationBuilder) GetUnorderedProfiles() []profiles.Updatable {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.GetUnorderedProfiles()
}

func (c *DynamicPopulationBuilder) GetCount() int {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.GetCount()
}

func (c *DynamicPopulationBuilder) GetLocalProfile() profiles.Updatable {
	return c.FindProfile(c.census.population.GetLocalProfile().GetNodeID())
}

func (c *DynamicPopulationBuilder) FindProfile(nodeID insolar.ShortNodeID) profiles.Updatable {
	c.census.mutex.RLock()
	defer c.census.mutex.RUnlock()

	return c.census.population.FindUpdatableProfile(nodeID)
}

func (c *DynamicPopulationBuilder) AddProfile(intro profiles.StaticProfile) profiles.Updatable {
	c.census.mutex.Lock()
	defer c.census.mutex.Unlock()

	if c.census.state.IsSealed() {
		panic("illegal state")
	}
	return c.census.population.AddProfile(intro)
}

func (c *DynamicPopulationBuilder) RemoveProfile(nodeID insolar.ShortNodeID) {
	c.census.mutex.Lock()
	defer c.census.mutex.Unlock()

	if c.census.state.IsSealed() {
		panic("illegal state")
	}
	c.census.population.RemoveProfile(nodeID)
}
