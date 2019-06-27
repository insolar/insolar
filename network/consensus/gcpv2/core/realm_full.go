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

package core

import (
	"fmt"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"

	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/gcpv2/census"

	"github.com/insolar/insolar/network/consensus/common"
)

type FullRealm struct {
	coreRealm
	/* Derived from the ones provided externally - set at init() or start(). Don't need mutex */
	handlers []packetRoute

	nodes    []NodeAppearance
	nodeRefs []*NodeAppearance

	timings  common2.RoundTimings
	nbhSizes common2.NeighbourhoodSizes

	othersCount      int
	joinersCount     int
	bftMajorityCount int

	census     census.ActiveCensus
	population census.OnlinePopulation

	/* Other fields - need mutex */
	isFinished bool
}

/* LOCK - runs under RoundController lock */
func (r *FullRealm) start() {
	r.roundContext, _ = inslogger.WithFields(r.roundContext, map[string]interface{}{
		"node_id":   r.GetSelfNodeID(),
		"pulse":     r.GetPulseNumber(),
		"is_joiner": r.IsJoiner(),
	})

	r.census = r.chronicle.GetActiveCensus()
	r.population = r.census.GetOnlinePopulation()

	r.initBasics()
	allCtls, perNodeCtls := r.initHandlers()
	r.initProjections(perNodeCtls)
	r.startWorkers(allCtls)
}

func (r *FullRealm) initBasics() {

	nodeCount := r.population.GetCount()

	r.timings = r.config.GetConsensusTimings(r.pulseData.NextPulseDelta, r.IsJoiner())
	r.strategy.AdjustConsensusTimings(&r.timings)

	r.nbhSizes = r.packetBuilder.GetNeighbourhoodSize(nodeCount)

	r.bftMajorityCount = common.BftMajority(nodeCount)

	r.nodes = make([]NodeAppearance, nodeCount)
	nodeCount-- // remove self
	r.nodeRefs = make([]*NodeAppearance, nodeCount)

	r.othersCount = nodeCount
}

func (r *FullRealm) initHandlers() (allControllers []PhaseController, perNodeControllers []PhaseController) {
	r.handlers = make([]packetRoute, packets.MaxPacketType)

	controllers := r.strategy.GetFullPhaseControllers(r.othersCount + 1)
	if len(controllers) == 0 {
		panic("no phase controllers")
	}
	individualHandlers := make([]PhaseController, 0, len(controllers))

	for _, ctl := range controllers {
		pt := ctl.GetPacketType()
		if !r.handlers[pt].IsEmpty() {
			panic("multiple handlers for packet type")
		}
		r.handlers[pt].realm = r
		switch ctl.IsPerNode() {
		case HandlerTypeHostPacket:
			r.handlers[pt].handlerHost = ctl.HandleHostPacket
			continue
		case HandlerTypeMemberPacket:
			r.handlers[pt].handlerMember = ctl.HandleMemberPacket
		case HandlerTypePerNodePacket:
			r.handlers[pt].setRedirectHandler(len(individualHandlers))
			individualHandlers = append(individualHandlers, ctl)
		}
		if !pt.IsMemberPacket() {
			panic("only member packet types can be handled as member/per-node")
		}
	}

	return controllers, individualHandlers
}

func (r *FullRealm) initProjections(individualHandlers []PhaseController) {

	thisNodeID := r.population.GetLocalProfile().GetShortNodeID()
	profiles := r.population.GetProfiles()
	baselineWeight := r.strategy.RandUint32()

	neighborTrustThreshold := uint8(r.nbhSizes.NeighbourhoodTrustThreshold)

	r.joinersCount = 0
	var j = 0
	prevSelf := r.self
	r.self = nil // resets self set on prep
	for i, p := range profiles {
		if p.IsJoiner() {
			r.joinersCount++
		}
		n := &r.nodes[i]
		n.init(p, &r.coreRealm.nodeCallback)
		n.neighborTrustThreshold = neighborTrustThreshold
		n.neighbourWeight = baselineWeight
		if p.GetShortNodeID() == thisNodeID {
			if r.self != nil {
				panic("schizophrenia")
			}
			r.self = n
		} else {
			if j == len(profiles) {
				panic("didnt find myself among active nodes")
			}
			r.nodeRefs[j] = n
			j++
		}
		if len(individualHandlers) > 0 {
			n.handlers = make([]PhasePerNodePacketHandler, len(individualHandlers))
			for _, ctl := range individualHandlers {
				ph := ctl.CreatePerNodePacketHandler(n)
				if ph == nil {
					panic("nil packet handler")
				}
				n.handlers[ctl.GetPacketType()] = ph
			}
		}
	}
	r.ShuffleNodeProjections(r.nodeRefs)

	// Transition data from prev self
	if prevSelf.IsJoiner() != r.self.IsJoiner() || prevSelf.GetShortNodeID() != r.self.GetShortNodeID() {
		panic("inconsistent transition of self between realms")
	}
	prevSelf.copySelfTo(r.self)
}

func (r *FullRealm) startWorkers(controllers []PhaseController) {
	for _, ctl := range controllers {
		ctl.BeforeStart(r)
	}
	for _, ctl := range controllers {
		ctl.StartWorker(r.roundContext)
	}
}

func (r *FullRealm) GetJoinersCount() int {
	return r.joinersCount
}

func (r *FullRealm) GetNodeCount() int {
	return len(r.nodes)
}

func (r *FullRealm) GetOthersCount() int {
	return r.othersCount
}

func (r *FullRealm) GetBftMajorityCount() int {
	return r.bftMajorityCount
}

func (r *FullRealm) FindActiveNode(id common.ShortNodeID) common2.NodeProfile {
	return r.population.FindProfile(id)
}

func (r *FullRealm) GetActiveNode(id common.ShortNodeID) (common2.NodeProfile, error) {
	np := r.population.FindProfile(id)
	if np == nil {
		return nil, fmt.Errorf("unknown ShortNodeID: %v", id)
	}
	return np, nil
}

func (r *FullRealm) GetNodeAppearance(id common.ShortNodeID) (*NodeAppearance, error) {
	np, err := r.GetActiveNode(id)
	if err != nil {
		return nil, err
	}
	return r.GetNodeAppearanceByIndex(np.GetIndex()), nil
}

func (r *FullRealm) GetNodeAppearanceByIndex(idx int) *NodeAppearance {
	return &r.nodes[idx]
}

func (r *FullRealm) GetPulseNumber() common.PulseNumber {
	return r.pulseData.PulseNumber
}

func (r *FullRealm) GetNextPulseNumber() common.PulseNumber {
	return r.pulseData.GetNextPulseNumber()
}

func (r *FullRealm) GetOriginalPulse() common2.OriginalPulsarPacket {
	// NB! locks for this field are only needed for PrepRealm
	return r.coreRealm.originalPulse
}

func (r *FullRealm) GetPulseData() common.PulseData {
	return r.pulseData
}

func (r *FullRealm) GetLastCloudStateHash() common2.CloudStateHash {
	return r.census.GetCloudStateHash()
}

func (r *coreRealm) UpstreamPreparePulseChange() <-chan common2.NodeStateHash {
	if !r.pulseData.PulseNumber.IsTimePulse() {
		panic("pulse number was not set")
	}

	sp := r.GetSelf().GetProfile()
	report := MembershipUpstreamReport{
		PulseNumber:     r.pulseData.PulseNumber,
		MemberPower:     sp.GetPower(),
		MembershipState: sp.GetState(),
	}
	return r.upstream.PreparePulseChange(report)
}

func (r *FullRealm) GetTimings() common2.RoundTimings {
	return r.timings
}

func (r *FullRealm) GetNeighbourhoodSizes() common2.NeighbourhoodSizes {
	return r.nbhSizes
}

/* Shuffled only once, when the round is created */
func (r *FullRealm) GetShuffledOtherNodes() []*NodeAppearance {
	return r.nodeRefs
}

func (r *FullRealm) GetLocalProfile() common2.LocalNodeProfile {
	return r.population.GetLocalProfile()
}

func (r *FullRealm) PrepareAndSetLocalNodeStateHashEvidence(nsh common2.NodeStateHash, nch common2.NodeClaimSignature) {
	// TODO use r.GetLastCloudStateHash() + digest(PulseData) + r.digest.GetGshDigester() to build digest for signing
	v := nsh.SignWith(r.signer)
	r.self.SetLocalNodeStateHashEvidence(common2.NewNodeStateHashEvidence(v), nch)
}

func (r *FullRealm) GetIndexedNodes() []NodeAppearance {
	return r.nodes
}

func (r *FullRealm) CreateNextPopulationBuilder() census.Builder {
	return r.chronicle.GetActiveCensus().CreateBuilder(r.GetNextPulseNumber())
}

func (r *FullRealm) FinishRound(builder census.Builder, csh common2.CloudStateHash) {
	r.Lock()
	defer r.Unlock()

	if r.isFinished {
		panic("illegal state")
	}
	r.isFinished = true

	r.prepareNewMembers(builder.GetOnlinePopulationBuilder())
	builder.BuildAndMakeExpected(csh)
}
