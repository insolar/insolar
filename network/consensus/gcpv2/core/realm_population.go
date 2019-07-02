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
	"context"
	"sync"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

type RealmPopulation interface {
	GetNodeCount() int
	GetOthersCount() int
	GetJoinersCount() int
	GetBftMajorityCount() int

	//GetOnlyActiveNode(id common.ShortNodeID) (common2.NodeProfile, error)
	//GetOnlyDynamicNodeAppearance(id common.ShortNodeID) *NodeAppearance
	GetNodeAppearance(id common.ShortNodeID) *NodeAppearance
	//GetOrAddNodeAppearance(id common.ShortNodeID) *NodeAppearance
	GetNodeAppearanceByIndex(idx int) *NodeAppearance

	GetShuffledOtherNodes() []*NodeAppearance
	GetIndexedNodes() []*NodeAppearance

	GetSelf() *NodeAppearance
	//CreateDynamicNode(constructionContext context.Context) *NodeAppearance
}

func NewMemberRealmPopulation(strategy RoundStrategy, population census.OnlinePopulation,
	individualHandlers []PhaseController, nodeContext NodeContextHolder, realm *FullRealm) *MemberRealmPopulation {

	nodeCount := population.GetCount()

	r := &MemberRealmPopulation{
		population:       population,
		nodeCount:        nodeCount,
		bftMajorityCount: common.BftMajority(nodeCount),
		nodeIndex:        make([]*NodeAppearance, nodeCount),
		nodeShuffle:      make([]*NodeAppearance, nodeCount-1),
	}
	r.initPopulation(strategy, individualHandlers, nodeContext, realm)
	return r
}

func (r *MemberRealmPopulation) initPopulation(strategy RoundStrategy, individualHandlers []PhaseController, nodeContext NodeContextHolder, realm *FullRealm) {
	profiles := r.population.GetProfiles()
	thisNodeID := r.population.GetLocalProfile().GetShortNodeID()
	baselineWeight := strategy.RandUint32()

	nodes := make([]NodeAppearance, r.nodeCount)

	var j = 0
	for i, p := range profiles {
		n := &nodes[i]
		n.init(p, nodeContext)
		n.neighbourWeight = baselineWeight
		r.nodeIndex[i] = n

		if p.GetShortNodeID() == thisNodeID {
			if r.self != nil {
				panic("schizophrenia")
			}
			r.self = n
		} else {
			if j == len(profiles) {
				panic("didnt find myself among active nodes")
			}
			r.nodeShuffle[j] = n
			j++
		}

		var sharedContext = context.Background()
		for k, ctl := range individualHandlers {
			var ph PhasePerNodePacketFunc
			ph, sharedContext = ctl.CreatePerNodePacketHandler(k, n, realm, sharedContext)
			if ph == nil {
				continue
			}
			if n.handlers == nil {
				n.handlers = make([]PhasePerNodePacketFunc, len(individualHandlers))
			}
			n.handlers[k] = ph
		}
	}
	ShuffleNodeProjections(strategy, r.nodeShuffle)
}

var _ RealmPopulation = &MemberRealmPopulation{}

type MemberRealmPopulation struct {
	population census.OnlinePopulation

	nodeIndex   []*NodeAppearance
	nodeShuffle []*NodeAppearance // excluding self
	self        *NodeAppearance

	nodeCount        int
	bftMajorityCount int

	//	purgatory	map[common.ShortNodeID]*NodeAppearance
	rw      sync.RWMutex
	joiners map[common.ShortNodeID]*NodeAppearance
}

func (r *MemberRealmPopulation) GetSelf() *NodeAppearance {
	return r.self
}

func (r *MemberRealmPopulation) GetNodeCount() int {
	return r.nodeCount
}

func (r *MemberRealmPopulation) GetJoinersCount() int {
	return 0
}

func (r *MemberRealmPopulation) GetOthersCount() int {
	return r.nodeCount - 1
}

func (r *MemberRealmPopulation) GetBftMajorityCount() int {
	return r.bftMajorityCount
}

func (r *MemberRealmPopulation) GetOnlyActiveNode(id common.ShortNodeID) common2.NodeProfile {
	return r.population.FindProfile(id)
}

func (r *MemberRealmPopulation) GetOnlyDynamicNodeAppearance(id common.ShortNodeID) *NodeAppearance {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.joiners[id]
}

func (r *MemberRealmPopulation) GetNodeAppearance(id common.ShortNodeID) *NodeAppearance {
	np := r.GetOnlyActiveNode(id)
	if np != nil {
		return r.GetNodeAppearanceByIndex(np.GetIndex())
	}
	return r.GetOnlyDynamicNodeAppearance(id)
}

func (r *MemberRealmPopulation) GetNodeAppearanceByIndex(idx int) *NodeAppearance {
	return r.nodeIndex[idx]
}

func (r *MemberRealmPopulation) GetShuffledOtherNodes() []*NodeAppearance {
	return r.nodeShuffle
}

func (r *MemberRealmPopulation) GetIndexedNodes() []*NodeAppearance {
	return r.nodeIndex
}
