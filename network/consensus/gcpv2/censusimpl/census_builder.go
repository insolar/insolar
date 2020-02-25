// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
