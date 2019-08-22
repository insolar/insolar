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
	"context"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

type Pulse interface {
	GetCensusState() State
	GetPulseNumber() pulse.Number
	GetExpectedPulseNumber() pulse.Number
	GetGlobulaStateHash() proofs.GlobulaStateHash
	GetCloudStateHash() proofs.CloudStateHash
	// returns true, when PulseData belongs to this census, PulseData can be empty for PrimingCensus
	GetNearestPulseData() (bool, pulse.Data)
}

type Archived interface {
	Pulse
	GetPulseData() pulse.Data
}

type Operational interface {
	Pulse
	GetOnlinePopulation() OnlinePopulation
	GetEvictedPopulation() EvictedPopulation
	GetOfflinePopulation() OfflinePopulation
	CreateBuilder(ctx context.Context, pn pulse.Number) Builder
	IsActive() bool

	GetMisbehaviorRegistry() MisbehaviorRegistry
	GetMandateRegistry() MandateRegistry
	GetProfileFactory(ksf cryptkit.KeyStoreFactory) profiles.Factory
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/census.Active -o . -s _mock.go -g

type Active interface {
	Operational
	GetPulseData() pulse.Data
}

type Prime interface {
	Active
	BuildCopy(pd pulse.Data, csh proofs.CloudStateHash, gsh proofs.GlobulaStateHash) Built
	// MakeExpected(pn pulse.Number, csh proofs.CloudStateHash, gsh proofs.GlobulaStateHash) Expected
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/census.Expected -o . -s _mock.go -g

type Expected interface {
	Operational
	GetPrevious() Active
	MakeActive(pd pulse.Data) Active
	Rebuild(pn pulse.Number) Built
}

type Built interface {
	GetOnlinePopulation() OnlinePopulation
	GetEvictedPopulation() EvictedPopulation
	GetGlobulaStateHash() proofs.GlobulaStateHash
	GetCloudStateHash() proofs.CloudStateHash
	GetNearestPulseData() (bool, pulse.Data)

	Update(csh proofs.CloudStateHash, gsh proofs.GlobulaStateHash) Built

	MakeExpected() Expected
}

type Builder interface {
	GetPopulationBuilder() PopulationBuilder

	// GetCensusState() State
	GetPulseNumber() pulse.Number
	// IsEphemeralAllowed() bool

	GetGlobulaStateHash() proofs.GlobulaStateHash
	SetGlobulaStateHash(gsh proofs.GlobulaStateHash)

	SealCensus()
	IsSealed() bool

	Build(csh proofs.CloudStateHash) Built
	BuildAsBroken(csh proofs.CloudStateHash) Built
	// BuildAndMakeExpected(csh proofs.CloudStateHash) Expected
	// BuildAndMakeBrokenExpected(csh proofs.CloudStateHash) Expected
}

type State uint8

const (
	DraftCensus State = iota
	SealedCensus
	CompleteCensus
	PrimingCensus
)

func (v State) HasPulseNumber() bool {
	return v > DraftCensus
}

func (v State) IsSealed() bool {
	return v != DraftCensus
}

func (v State) IsBuilt() bool {
	return v >= CompleteCensus
}
