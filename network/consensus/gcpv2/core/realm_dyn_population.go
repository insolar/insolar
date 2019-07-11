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
	"fmt"
	"sync"

	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

func NewDynamicRealmPopulation(baselineWeight uint32, local common2.NodeProfile, nodeCountHint int,
	fn NodeInitFunc) *DynamicRealmPopulation {

	r := &DynamicRealmPopulation{
		nodeInit:       fn,
		baselineWeight: baselineWeight,
	}
	r.initPopulation(local, nodeCountHint)

	return r
}

var _ RealmPopulation = &DynamicRealmPopulation{}

type DynamicRealmPopulation struct {
	nodeInit       NodeInitFunc
	baselineWeight uint32
	self           *NodeAppearance

	rw sync.RWMutex

	joinerCount  int
	indexedCount int

	nodeIndex     []*NodeAppearance
	nodeShuffle   []*NodeAppearance // excluding self
	dynamicNodes  map[common.ShortNodeID]*NodeAppearance
	purgatoryByPK map[string]*NodeAppearance
	purgatoryByID map[common.ShortNodeID]*[]*NodeAppearance
}

func (r *DynamicRealmPopulation) initPopulation(local common2.NodeProfile, nodeCountHint int) {
	r.self = r.CreateNodeAppearance(context.Background(), local)
}

func (r *DynamicRealmPopulation) GetSelf() *NodeAppearance {
	return r.self
}

func (r *DynamicRealmPopulation) GetNodeCount() int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return len(r.nodeIndex)
}

func (r *DynamicRealmPopulation) GetJoinersCount() int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.joinerCount
}

func (r *DynamicRealmPopulation) GetOthersCount() int {
	return r.GetNodeCount() - 1
}

func (r *DynamicRealmPopulation) GetBftMajorityCount() int {
	return common.BftMajority(r.GetNodeCount())
}

func (r *DynamicRealmPopulation) GetNodeAppearance(id common.ShortNodeID) *NodeAppearance {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.dynamicNodes[id]
}

func (r *DynamicRealmPopulation) GetActiveNodeAppearance(id common.ShortNodeID) *NodeAppearance {
	na := r.GetNodeAppearance(id)
	if !na.GetProfile().IsJoiner() {
		return na
	}
	return nil
}

func (r *DynamicRealmPopulation) GetJoinerNodeAppearance(id common.ShortNodeID) *NodeAppearance {
	na := r.GetNodeAppearance(id)
	if !na.GetProfile().IsJoiner() {
		return nil
	}
	return na
}

func (r *DynamicRealmPopulation) GetNodeAppearanceByIndex(idx int) *NodeAppearance {
	if idx < 0 {
		panic("illegal value")
	}

	r.rw.RLock()
	defer r.rw.RUnlock()

	if idx >= len(r.nodeIndex) {
		return nil
	}
	return r.nodeIndex[idx]
}

func (r *DynamicRealmPopulation) GetShuffledOtherNodes() []*NodeAppearance {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.nodeShuffle
}

func (r *DynamicRealmPopulation) IsComplete() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return len(r.nodeIndex) == r.indexedCount
}

func (r *DynamicRealmPopulation) GetIndexedNodes() []*NodeAppearance {
	cp, _ := r.GetIndexedNodesWithCheck()
	//if !ok {
	//	panic("node set is incomplete")
	//}
	return cp
}

func (r *DynamicRealmPopulation) GetIndexedNodesWithCheck() ([]*NodeAppearance, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	cp := make([]*NodeAppearance, len(r.nodeIndex))
	copy(cp, r.nodeIndex)

	return cp, len(r.nodeIndex) == r.indexedCount
}

func (r *DynamicRealmPopulation) CreateNodeAppearance(ctx context.Context, np common2.NodeProfile) *NodeAppearance {

	n := &NodeAppearance{}
	n.init(np, nil, r.baselineWeight)
	r.nodeInit(ctx, n)

	return n
}

func (r *DynamicRealmPopulation) AddToPurgatory(n *NodeAppearance) (*NodeAppearance, PurgatoryNodeState) {
	if n.profile.HasIntroduction() {
		panic("illegal value")
	}

	id := n.profile.GetShortNodeID()
	na := r.GetActiveNodeAppearance(id)
	if na != nil {
		return na, PurgatoryExistingMember
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	nn := r.dynamicNodes[id]
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
		return n, PurgatoryNodeState(len(*nodes) - 1)
	}
}

func (r *DynamicRealmPopulation) AddToDynamics(n *NodeAppearance) (*NodeAppearance, []*NodeAppearance) {
	np := n.profile

	if !np.HasIntroduction() {
		panic("illegal value")
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	id := np.GetShortNodeID()

	delete(r.purgatoryByPK, np.GetNodePublicKey().AsByteString())
	nodes := r.purgatoryByID[id]
	if nodes != nil {
		delete(r.purgatoryByID, id)
	} else {
		nodes = &[]*NodeAppearance{}
	}

	na := r.GetActiveNodeAppearance(id)
	if na != nil {
		return na, *nodes
	}

	na = r.dynamicNodes[id]
	if na != nil {
		return na, *nodes
	}

	if np.IsJoiner() {
		r.joinerCount++
	} else {
		ni := np.GetIndex()
		switch {
		case ni == len(r.nodeIndex):
			r.nodeIndex = append(r.nodeIndex, n)
		case ni > len(r.nodeIndex):
			r.nodeIndex = append(r.nodeIndex, make([]*NodeAppearance, 1+ni-len(r.nodeIndex))...)
			r.nodeIndex[ni] = n
		default:
			if r.nodeIndex[ni] != nil {
				panic(fmt.Sprintf("duplicate node id(%v)", ni))
			}
			r.nodeIndex[ni] = n
		}
		r.indexedCount++
		r.nodeShuffle = append(r.nodeShuffle, n)
	}
	return n, *nodes
}

//
func (r *DynamicRealmPopulation) SetOrUpdateVectorHelper(v *RealmVectorHelper) *RealmVectorHelper {
	if v.HasSameVersion(r.self.callback.GetPopulationVersion()) {
		return v
	}

	r.rw.RLock()
	defer r.rw.RUnlock()

	return v.SetOrUpdateNodes(r.nodeIndex, r.joinerCount, r.self.callback.GetPopulationVersion())
}
