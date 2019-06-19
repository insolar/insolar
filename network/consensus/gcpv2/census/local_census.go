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

type ConsensusChronicles interface {
	GetActiveCensus() ActiveCensus
	GetExpectedCensus() ExpectedCensus
	GetLatestCensus() OperationalCensus
	// FindArchivedCensus(pn common.PulseNumber) ArchivedCensus
}

type PulseCensus interface {
	// GetCensusType() CensusType
	GetCensusState() State
	GetPulseNumber() common.PulseNumber
	GetExpectedPulseNumber() common.PulseNumber
	GetGlobulaStateHash() common2.GlobulaStateHash
	GetCloudStateHash() common2.CloudStateHash
}

type ArchivedCensus interface {
	PulseCensus
	GetPulseData() common.PulseData
}

type OperationalCensus interface {
	PulseCensus
	GetOnlinePopulation() OnlinePopulation
	GetOfflinePopulation() OfflinePopulation
	CreateBuilder(pn common.PulseNumber) Builder
	IsActive() bool

	GetMisbehaviorRegistry() MisbehaviorRegistry
	GetMandateRegistry() MandateRegistry
}

type ActiveCensus interface {
	OperationalCensus
	GetPulseData() common.PulseData
}

type ExpectedCensus interface {
	OperationalCensus
	GetPrevious() ActiveCensus
	MakeActive(pd common.PulseData) ActiveCensus
}

type Builder interface {
	GetOnlinePopulationBuilder() OnlinePopulationBuilder
	GetOnlinePopulationView() OnlinePopulation

	GetCensusState() State
	GetPulseNumber() common.PulseNumber

	GetGlobulaStateHash() common2.GlobulaStateHash
	SetGlobulaStateHash(gsh common2.GlobulaStateHash)

	SealCensus()
	IsSealed() bool

	BuildAndMakeExpected(csh common2.CloudStateHash) ExpectedCensus
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
	FindProfile(nodeId common.ShortNodeId) common2.NodeProfile
	GetCount() int
	GetProfiles() []common2.NodeProfile
	GetLocalProfile() common2.LocalNodeProfile
}

type OnlinePopulationBuilder interface {
	GetCount() int
	AddJoinerProfile(intro common2.NodeIntroProfile) common2.UpdatableNodeProfile
	RemoveProfile(nodeId common.ShortNodeId)
	GetUnorderedProfiles() []common2.UpdatableNodeProfile
	FindProfile(nodeId common.ShortNodeId) common2.UpdatableNodeProfile
	GetLocalProfile() common2.LocalNodeProfile
}
