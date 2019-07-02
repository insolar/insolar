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

	GetNodeAppearance(id common.ShortNodeID) *NodeAppearance
	GetJoinerNodeAppearance(id common.ShortNodeID) *NodeAppearance
	GetNodeAppearanceByIndex(idx int) *NodeAppearance

	GetShuffledOtherNodes() []*NodeAppearance
	GetIndexedNodes() []*NodeAppearance

	GetSelf() *NodeAppearance
	//CreateDynamicNode(constructionContext context.Context) *NodeAppearance

	//AddToPurgatory(np common2.NodeIntroProfile) (*NodeAppearance, int)
	//AddToJoiners(n *NodeAppearance) (*NodeAppearance, int)
}

func NewMemberRealmPopulation(strategy RoundStrategy, population census.OnlinePopulation,
	fn NodeInitFunc) *MemberRealmPopulation {

	nodeCount := population.GetCount()

	r := &MemberRealmPopulation{
		population:       population,
		nodeInit:         fn,
		baselineWeight:   strategy.RandUint32(),
		nodeCount:        nodeCount,
		bftMajorityCount: common.BftMajority(nodeCount),
		nodeIndex:        make([]*NodeAppearance, nodeCount),
		nodeShuffle:      make([]*NodeAppearance, nodeCount-1),
	}
	r.initPopulation()
	ShuffleNodeProjections(strategy, r.nodeShuffle)

	return r
}

type NodeInitFunc func(ctx context.Context, n *NodeAppearance)

func (r *MemberRealmPopulation) initPopulation() {
	profiles := r.population.GetProfiles()
	thisNodeID := r.population.GetLocalProfile().GetShortNodeID()

	nodes := make([]NodeAppearance, r.nodeCount)

	var j = 0
	for i, p := range profiles {
		n := &nodes[i]
		r.nodeIndex[i] = n

		n.init(p, nil, r.baselineWeight)
		r.nodeInit(context.Background(), n)

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
	}
}

var _ RealmPopulation = &MemberRealmPopulation{}

type MemberRealmPopulation struct {
	population     census.OnlinePopulation
	nodeInit       NodeInitFunc
	baselineWeight uint32

	nodeIndex   []*NodeAppearance
	nodeShuffle []*NodeAppearance // excluding self
	self        *NodeAppearance

	nodeCount        int
	bftMajorityCount int

	rw sync.RWMutex

	joiners       map[common.ShortNodeID]*NodeAppearance
	purgatoryByPK map[string]*NodeAppearance
	purgatoryByID map[common.ShortNodeID]*[]*NodeAppearance
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

func (r *MemberRealmPopulation) GetNodeAppearance(id common.ShortNodeID) *NodeAppearance {
	na := r.getActiveNodeAppearance(id)
	if na != nil {
		return na
	}
	return r.GetJoinerNodeAppearance(id)
}

func (r *MemberRealmPopulation) GetJoinerNodeAppearance(id common.ShortNodeID) *NodeAppearance {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.joiners[id]
}

func (r *MemberRealmPopulation) getActiveNodeAppearance(id common.ShortNodeID) *NodeAppearance {
	np := r.population.FindProfile(id)
	if np != nil {
		return r.GetNodeAppearanceByIndex(np.GetIndex())
	}
	return nil
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

func (r *MemberRealmPopulation) createNode(ctx context.Context, inp common2.NodeIntroProfile) *NodeAppearance {
	//np.GetNodePublicKey().AsByteString()
	np := &common2.JoinerNodeProfile{NodeIntroProfile: inp}

	n := &NodeAppearance{}
	n.init(np, nil, r.baselineWeight)
	r.nodeInit(ctx, n)
	return n
}

const PurgatoryDuplicatePK = -1
const PurgatoryExistingMember = -2

//TODO remember who has sent
func (r *MemberRealmPopulation) AddToPurgatory(n *NodeAppearance) (*NodeAppearance, int) {
	if !n.profile.GetState().IsUndefined() {
		panic("illegal value")
	}

	id := n.profile.GetShortNodeID()
	na := r.getActiveNodeAppearance(id)
	if na != nil {
		return na, PurgatoryExistingMember
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	nn := r.joiners[id]
	if nn != nil {
		return nn, PurgatoryExistingMember
	}

	if r.purgatoryByPK == nil {
		r.purgatoryByPK = make(map[string]*NodeAppearance)
		r.purgatoryByID = make(map[common.ShortNodeID]*[]*NodeAppearance)

		r.purgatoryByPK[n.profile.GetNodePublicKey().AsByteString()] = n
		r.purgatoryByID[n.profile.GetShortNodeID()] = &[]*NodeAppearance{n}
		return n, 0
	}

	pk := n.profile.GetNodePublicKey().AsByteString()
	nn = r.purgatoryByPK[pk]
	if nn != nil {
		return nn, PurgatoryDuplicatePK
	}

	nodes := r.purgatoryByID[id]

	if nodes == nil {
		nodes = &[]*NodeAppearance{n}
		r.purgatoryByID[id] = nodes
		return n, 0
	} else {
		*nodes = append(*nodes, n)
		return n, len(*nodes) - 1
	}
}

func (r *MemberRealmPopulation) AddToJoiners(n *NodeAppearance) (*NodeAppearance, int) {
	if !n.profile.GetState().IsJoining() {
		panic("illegal value")
	}

	id := n.profile.GetShortNodeID()
	na := r.getActiveNodeAppearance(id)
	if na != nil {
		return na, PurgatoryExistingMember
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	nn := r.joiners[id]
	if nn != nil {
		return nn, PurgatoryExistingMember
	}

	delete(r.purgatoryByPK, n.profile.GetNodePublicKey().AsByteString())
	delete(r.purgatoryByID, n.profile.GetShortNodeID())

	return n, 0
}
