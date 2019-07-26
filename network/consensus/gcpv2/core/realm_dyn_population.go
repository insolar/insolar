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

	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

func NewDynamicRealmPopulation(strategy RoundStrategy, population census.OnlinePopulation, nodeCountHint int, phase2ExtLimit uint8,
	shuffleFn args.ShuffleFunc, fn NodeInitFunc) *DynamicRealmPopulation {

	nodeCount := population.GetIndexedCapacity()
	if nodeCount > nodeCountHint {
		nodeCountHint = nodeCount + nodeCount>>2 // 125%
	}

	r := &DynamicRealmPopulation{
		nodeInit:       fn,
		shuffleFunc:    shuffleFn,
		baselineWeight: strategy.GetBaselineWeightForNeighbours(),
		phase2ExtLimit: phase2ExtLimit,
	}
	r.initPopulation(population, nodeCountHint)

	return r
}

var _ RealmPopulation = &DynamicRealmPopulation{}

const reshuffleTolerance int = 3

type DynamicRealmPopulation struct {
	nodeInit       NodeInitFunc
	shuffleFunc    args.ShuffleFunc
	baselineWeight uint32
	phase2ExtLimit uint8
	self           *NodeAppearance

	rw sync.RWMutex

	joinerCount   int
	indexedCount  int
	indexedLenSet bool
	shuffledCount int

	nodeIndex    []*NodeAppearance
	nodeShuffle  []*NodeAppearance // excluding self
	dynamicNodes map[insolar.ShortNodeID]*NodeAppearance
}

func (r *DynamicRealmPopulation) SealIndexed(indexedCountLimit int) bool {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.indexedLenSet {
		return len(r.nodeIndex) == indexedCountLimit
	}
	if len(r.nodeIndex) > indexedCountLimit {
		return false
	}

	if len(r.nodeIndex) != indexedCountLimit {
		cp := make([]*NodeAppearance, indexedCountLimit)
		copy(cp, r.nodeIndex)
		r.nodeIndex = cp
	}

	if indexedCountLimit > cap(r.nodeShuffle) {
		r.nodeShuffle = append(make([]*NodeAppearance, 0, indexedCountLimit), r.nodeShuffle...)
	}

	r.indexedLenSet = true
	if r.indexedCount == indexedCountLimit {
		r.onDynamicPopulationCompleted()
	}
	return true
}

func (r *DynamicRealmPopulation) initPopulation(population census.OnlinePopulation, nodeCountHint int) {

	r.dynamicNodes = make(map[insolar.ShortNodeID]*NodeAppearance, nodeCountHint)
	r.nodeIndex = make([]*NodeAppearance, 0, nodeCountHint)
	r.nodeShuffle = make([]*NodeAppearance, 0, nodeCountHint)

	local := population.GetLocalProfile()
	r.self = r.CreateNodeAppearance(context.Background(), local)
	r.self, _ = r.AddToDynamics(r.self)

	for _, np := range population.GetProfiles() {
		na := r.CreateNodeAppearance(context.Background(), np)
		_, _ = r.AddToDynamics(na) // repeated addition will leave the initial node
	}
	self := r.dynamicNodes[local.GetNodeID()]
	if r.self == nil || r.self != self {
		panic("illegal state")
	}
}

func (r *DynamicRealmPopulation) GetSelf() *NodeAppearance {
	return r.self
}

func (r *DynamicRealmPopulation) GetIndexedCount() int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return len(r.nodeIndex)
}

func (r *DynamicRealmPopulation) GetJoinersCount() int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.joinerCount
}

func (r *DynamicRealmPopulation) GetNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.dynamicNodes[id]
}

func (r *DynamicRealmPopulation) GetActiveNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	na := r.GetNodeAppearance(id)
	if na == nil || na.GetProfile().IsJoiner() {
		return nil
	}
	return na
}

func (r *DynamicRealmPopulation) GetJoinerNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	na := r.GetNodeAppearance(id)
	if na == nil || !na.GetProfile().IsJoiner() {
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

func (r *DynamicRealmPopulation) readShuffledOtherNodes() (bool, []*NodeAppearance) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.shuffleFunc == nil || r.shuffledCount+reshuffleTolerance >= len(r.nodeShuffle) {
		return true, r.nodeShuffle
	}
	return false, nil
}

func (r *DynamicRealmPopulation) GetShuffledOtherNodes() []*NodeAppearance {

	ok, nodes := r.readShuffledOtherNodes()
	if ok {
		return nodes
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	if r.shuffledCount+3 >= len(r.nodeShuffle) {
		return r.nodeShuffle
	}

	cp := append(make([]*NodeAppearance, 0, len(r.nodeShuffle)), r.nodeShuffle...)
	ShuffleNodeAppearances(r.shuffleFunc, cp)
	r.shuffledCount = len(r.nodeShuffle)
	r.nodeShuffle = cp
	return cp
}

func ShuffleNodeAppearances(shuffleFunc args.ShuffleFunc, nodeRefs []*NodeAppearance) {
	shuffleFunc(len(nodeRefs),
		func(i, j int) { nodeRefs[i], nodeRefs[j] = nodeRefs[j], nodeRefs[i] })
}

func (r *DynamicRealmPopulation) GetIndexedNodes() []*NodeAppearance {
	cp, _ := r.GetIndexedNodesAndHasNil()
	return cp
}

func (r *DynamicRealmPopulation) GetIndexedNodesAndHasNil() ([]*NodeAppearance, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	cp := make([]*NodeAppearance, len(r.nodeIndex))
	copy(cp, r.nodeIndex)

	return cp, len(r.nodeIndex) > r.indexedCount
}

func (r *DynamicRealmPopulation) GetSealedCapacity() (int, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return len(r.nodeIndex), r.indexedLenSet
}

func (r *DynamicRealmPopulation) GetCountAndCompleteness(includeJoiners bool) (int, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	count := r.indexedCount
	if includeJoiners {
		count += r.joinerCount
	}

	return count, r.indexedLenSet && len(r.nodeIndex) == r.indexedCount
}

func (r *DynamicRealmPopulation) CreatePacketLimiter() phases.PacketLimiter {
	return phases.NewPacketLimiter(r.phase2ExtLimit)
}

func (r *DynamicRealmPopulation) CreateNodeAppearance(ctx context.Context, np profiles.ActiveNode) *NodeAppearance {

	n := &NodeAppearance{}
	n.init(np, nil, r.baselineWeight, r.CreatePacketLimiter())
	n.requestedPower = np.GetDeclaredPower()
	r.nodeInit(ctx, n)

	return n
}

func (r *DynamicRealmPopulation) AddToDynamics(na *NodeAppearance) (*NodeAppearance, error) {
	nna := r.addToDynamics(na)

	if na == nna {
		na.onNodeAdded(context.TODO()) // TODO context?
	} else if !profiles.EqualBriefProfiles(nna.profile.GetStatic(), na.profile.GetStatic()) {
		return nil, fmt.Errorf("multiple joiners on same id(%v): %v", na.GetNodeID(), []*NodeAppearance{na, nna})
	}

	return nna, nil
}

func (r *DynamicRealmPopulation) addToDynamics(n *NodeAppearance) *NodeAppearance {
	nip := n.profile.GetStatic()

	// if nip.GetExtension() == nil {
	//	panic("illegal value")
	// }

	r.rw.Lock()
	defer r.rw.Unlock()

	id := nip.GetStaticNodeID()
	na := r.dynamicNodes[id]
	if na != nil {
		return na
	}

	if n.profile.IsJoiner() {
		r.joinerCount++
	} else {
		ni := n.profile.GetIndex()
		switch {
		case ni.AsInt() == len(r.nodeIndex):
			r.nodeIndex = append(r.nodeIndex, n)
		case ni.AsInt() > len(r.nodeIndex):
			r.nodeIndex = append(r.nodeIndex, make([]*NodeAppearance, 1+ni.AsInt()-len(r.nodeIndex))...)
			r.nodeIndex[ni] = n
		default:
			if r.nodeIndex[ni] != nil {
				panic(fmt.Sprintf("duplicate node index(%v)", ni))
			}
			r.nodeIndex[ni] = n
		}
		r.indexedCount++
		r.nodeShuffle = append(r.nodeShuffle, n)

		if r.indexedLenSet && r.indexedCount == len(r.nodeIndex) {
			r.onDynamicPopulationCompleted()
		}
	}
	r.dynamicNodes[id] = n

	flags := FlagCreated
	if nip.GetExtension() != nil {
		flags |= FlagProfileUpdated
	}
	n.callback.onDynamicNodeUpdate(n.callback.updatePopulationVersion(), n, flags)

	return n
}

func (r *DynamicRealmPopulation) appendAnyNodes(includeIndexed bool, nodes []*NodeAppearance) []*NodeAppearance {

	r.rw.RLock()
	defer r.rw.RUnlock()

	delta := len(r.dynamicNodes)
	if includeIndexed {
		delta -= r.indexedCount
	}
	if delta < 0 {
		panic("illegal state")
	}
	if delta == 0 {
		return nodes
	}

	index := len(nodes)
	nodes = append(nodes, make([]*NodeAppearance, delta)...)
	for _, v := range r.dynamicNodes {
		if !includeIndexed && !v.IsJoiner() {
			continue
		}
		nodes[index] = v
		index++
	}
	return nodes
}

func (r *DynamicRealmPopulation) GetAnyNodes(includeIndexed bool, shuffle bool) []*NodeAppearance {

	nodes := r.appendAnyNodes(includeIndexed, nil)
	if shuffle && r.shuffleFunc != nil {
		ShuffleNodeAppearances(r.shuffleFunc, nodes)
	}
	return nodes
}

//
func (r *DynamicRealmPopulation) CreateVectorHelper() *RealmVectorHelper {
	r.rw.RLock()
	defer r.rw.RUnlock()

	v := &RealmVectorHelper{realmPopulation: r}
	v.setArrayNodes(r.nodeIndex, r.dynamicNodes, r.self.callback.GetPopulationVersion())
	v.realmPopulation = r
	return v
}

func (r *DynamicRealmPopulation) onDynamicPopulationCompleted() {
	go r.self.callback.onDynamicPopulationCompleted(r.self.callback.GetPopulationVersion(), r.indexedCount)
}
