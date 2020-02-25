// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package censusimpl

import (
	"fmt"
	"sync"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/pulse"
)

func NewLocalChronicles(profileFactory profiles.Factory) LocalConsensusChronicles {
	return &localChronicles{profileFactory: profileFactory}
}

type LocalConsensusChronicles interface {
	api.ConsensusChronicles
	makeActive(ce census.Expected, ca localActiveCensus)
}

var _ api.ConsensusChronicles = &localChronicles{}

type localActiveCensus interface {
	census.Active
	getVersionedRegistries() census.VersionedRegistries
	setVersionedRegistries(vr census.VersionedRegistries)
	onMadeActive()
}

type localChronicles struct {
	rw             sync.RWMutex
	active         localActiveCensus
	expected       census.Expected
	profileFactory profiles.Factory
}

func (c *localChronicles) GetLatestCensus() (census.Operational, bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if c.expected != nil {
		return c.expected, true
	}
	return c.active, false
}

func (c *localChronicles) GetRecentCensus(pn pulse.Number) census.Operational {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if c.expected != nil && pn == c.expected.GetPulseNumber() {
		return c.expected
	}

	if pn == c.active.GetPulseNumber() {
		return c.active
	}
	panic(fmt.Sprintf("recent census is missing for pulse (%v)", pn))
}

func (c *localChronicles) GetActiveCensus() census.Active {
	c.rw.RLock()
	defer c.rw.RUnlock()

	return c.active
}

func (c *localChronicles) GetExpectedCensus() census.Expected {
	c.rw.RLock()
	defer c.rw.RUnlock()

	return c.expected
}

func (c *localChronicles) makeActive(ce census.Expected, ca localActiveCensus) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if c.expected != ce {
		panic("illegal state")
	}

	if c.active == nil {
		// priming
		if ce != nil {
			panic("illegal state")
		}
		if ca.getVersionedRegistries() == nil {
			panic("versioned registries are missing")
		}
	} else {
		pd := ca.GetPulseData()
		if pd.IsEmpty() {
			panic("illegal value")
		}

		lastRealPulse := c.active.getVersionedRegistries().GetNearestValidPulseData()

		pda := c.active.GetPulseData()

		checkExpectedPulse := true
		switch {
		case pda.PulseEpoch.IsEphemeral(): // supports empty with ephemeral
			if pd.IsFromEphemeral() {
				if !pda.IsEmpty() && !pda.IsValidNext(pd) {
					panic("illegal value - ephemeral pulses must be consecutive")
				}
				break
			}
			// we can't check it vs last ephemeral, so lets take the last real one
			pda = lastRealPulse
			checkExpectedPulse = false
			fallthrough
		case pda.IsFromPulsar() || pda.IsEmpty():
			// must be regular pulse
			if !pd.IsValidPulsarData() {
				panic("illegal value")
			}

			if !pda.IsEmpty() && pd.PulseNumber < pda.NextPulseNumber() {
				panic("illegal value - pulse retroactive")
			}
		}

		if checkExpectedPulse && !ce.GetPulseNumber().IsUnknownOrEqualTo(pd.PulseNumber) {
			panic("illegal value")
		}

		registries := c.active.getVersionedRegistries()
		if pd.IsFromPulsar() {
			registries = registries.CommitNextPulse(pd, ca.GetOnlinePopulation())
		}
		ca.setVersionedRegistries(registries)
	}

	c.active = ca
	c.expected = nil
	ca.onMadeActive()
}

func (c *localChronicles) makeExpected(ce census.Expected) census.Expected {
	c.rw.Lock()
	defer c.rw.Unlock()

	if c.active != ce.GetPrevious() {
		panic("illegal state")
	}

	if c.expected != nil && c.expected.GetOnlinePopulation() != ce.GetOnlinePopulation() {
		panic("illegal state")
	}

	c.expected = ce
	return ce
}

func (c *localChronicles) GetProfileFactory(factory cryptkit.KeyStoreFactory) profiles.Factory {
	return c.profileFactory
}
