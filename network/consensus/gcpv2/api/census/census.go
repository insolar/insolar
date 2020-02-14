// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
