///
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
///

package api

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
)

type ConsensusChronicles interface {
	GetActiveCensus() ActiveCensus
	GetExpectedCensus() ExpectedCensus
	GetLatestCensus() OperationalCensus
	GetRecentCensus(pn pulse_data.PulseNumber) OperationalCensus
	// FindArchivedCensus(pn common.PulseNumber) ArchivedCensus
}

type PulseCensus interface {
	// GetCensusType() CensusType
	GetCensusState() State
	GetPulseNumber() pulse_data.PulseNumber
	GetExpectedPulseNumber() pulse_data.PulseNumber
	GetGlobulaStateHash() gcp_types.GlobulaStateHash
	GetCloudStateHash() gcp_types.CloudStateHash
}

type ArchivedCensus interface {
	PulseCensus
	GetPulseData() pulse_data.PulseData
}

type OperationalCensus interface {
	PulseCensus
	GetOnlinePopulation() OnlinePopulation
	GetEvictedPopulation() EvictedPopulation
	GetOfflinePopulation() OfflinePopulation
	CreateBuilder(pn pulse_data.PulseNumber, fullCopy bool) Builder
	IsActive() bool

	GetMisbehaviorRegistry() MisbehaviorRegistry
	GetMandateRegistry() MandateRegistry
}

type ActiveCensus interface {
	OperationalCensus
	GetPulseData() pulse_data.PulseData
	GetProfileFactory(ksf cryptography_containers.KeyStoreFactory) gcp_types.NodeProfileFactory
}

type ExpectedCensus interface {
	OperationalCensus
	GetPrevious() ActiveCensus
	MakeActive(pd pulse_data.PulseData) ActiveCensus
}

type Builder interface {
	GetPopulationBuilder() PopulationBuilder

	GetCensusState() State
	GetPulseNumber() pulse_data.PulseNumber

	GetGlobulaStateHash() gcp_types.GlobulaStateHash
	SetGlobulaStateHash(gsh gcp_types.GlobulaStateHash)

	SealCensus()
	IsSealed() bool

	BuildAndMakeExpected(csh gcp_types.CloudStateHash) ExpectedCensus
}

// type CensusType uint8
//
// const (
// 	ConsensusCensusType = iota
// 	GenesisCensusType
// 	ProvidedCensusType
// 	DetachedCensusType
// 	SimulatedCensusType
// )

type State uint8

const (
	DraftCensus State = iota
	SealedCensus
	BuiltCensus
	PrimingCensus
)

func (v State) HasPulseNumber() bool {
	return v > DraftCensus
}

func (v State) IsSealed() bool {
	return v != DraftCensus
}

func (v State) IsBuilt() bool {
	return v >= BuiltCensus
}

type OnlinePopulation interface {
	FindProfile(nodeID insolar.ShortNodeID) gcp_types.NodeProfile
	GetCount() int
	GetProfiles() []gcp_types.NodeProfile
	GetLocalProfile() gcp_types.LocalNodeProfile
}

type EvictedPopulation interface {
	FindProfile(nodeID insolar.ShortNodeID) gcp_types.EvictedNodeProfile
	GetCount() int
	GetProfiles() []gcp_types.EvictedNodeProfile
}

type PopulationBuilder interface {
	GetCount() int
	AddJoinerProfile(intro gcp_types.NodeIntroProfile) gcp_types.UpdatableNodeProfile
	RemoveProfile(nodeID insolar.ShortNodeID)
	GetUnorderedProfiles() []gcp_types.UpdatableNodeProfile
	FindProfile(nodeID insolar.ShortNodeID) gcp_types.UpdatableNodeProfile
	GetLocalProfile() gcp_types.UpdatableNodeProfile
	RemoveOthers()
}
