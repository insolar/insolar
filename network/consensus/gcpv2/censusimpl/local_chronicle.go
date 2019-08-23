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
		case pda.PulseEpoch == pulse.EphemeralPulseEpoch: // supports empty with ephemeral
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

			if !pda.IsEmpty() && pd.PulseNumber < pda.GetNextPulseNumber() {
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
