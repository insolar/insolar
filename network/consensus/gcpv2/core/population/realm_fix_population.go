// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package population

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
)

func NewFixedRealmPopulation(population census.OnlinePopulation, phase2ExtLimit uint8,
	shuffleFn args.ShuffleFunc, baselineWeight uint32, hookCfg SharedNodeContext, fn DispatchFactoryFunc) *FixedRealmPopulation {

	if !population.IsValid() || population.GetIndexedCount() != population.GetIndexedCapacity() {
		panic("illegal value - fixed realm population can't be initialized on incomplete/invalid population")
	}
	nodeCount := population.GetIndexedCount()
	otherCount := nodeCount
	if otherCount > 0 && !population.GetLocalProfile().IsJoiner() {
		otherCount-- // remove self when it is not a joiner
	}

	r := &FixedRealmPopulation{
		population: population,
		dynPop: dynPop{DynamicRealmPopulation{
			dispatchFactory: fn,
			shuffleFunc:     shuffleFn,
			baselineWeight:  baselineWeight,
			phase2ExtLimit:  phase2ExtLimit,
			indexedCount:    nodeCount,
			nodeIndex:       make([]*NodeAppearance, nodeCount),
			nodeShuffle:     make([]*NodeAppearance, otherCount),
			shuffledCount:   otherCount,
			dynamicNodes:    make(map[insolar.ShortNodeID]*NodeAppearance),
			indexedLenSet:   true, // locks down SealIndexed
			hook:            NewHook(nil, nil, hookCfg),
		}},
	}
	r.initPopulation()
	ShuffleNodeAppearances(shuffleFn, r.nodeShuffle)

	return r
}

var _ RealmPopulation = &FixedRealmPopulation{}

type dynPop struct{ DynamicRealmPopulation }

type FixedRealmPopulation struct {
	dynPop
	population census.OnlinePopulation
}

func (r *FixedRealmPopulation) GetSealedCapacity() (int, bool) {
	return len(r.nodeIndex), true
}

func (r *FixedRealmPopulation) initPopulation() {

	ctx := context.Background()

	activeProfiles := r.population.GetProfiles()
	thisNodeID := r.population.GetLocalProfile().GetNodeID()

	nodes := make([]NodeAppearance, r.indexedCount)

	var j = 0
	for i, p := range activeProfiles {
		if p.GetOpMode().IsEvicted() {
			panic("illegal state")
		}
		if p.IsJoiner() {
			panic("illegal state")
		}

		n := &nodes[i]
		*n = NewNodeAppearance(p, r.baselineWeight,
			r.CreatePacketLimiter(false), &r.hook,
			nil)
		r.nodeIndex[i] = n
		n.handlers = r.dispatchFactory(ctx, n)

		if p.GetNodeID() == thisNodeID {
			if r.self != nil {
				panic("schizophrenia")
			}
			r.self = n
		} else {
			r.nodeShuffle[j] = n
			j++
		}
	}
	if r.self == nil {
		panic("illegal state")
	}
	r.initHook()
}

func (r *FixedRealmPopulation) NotifyAllOnAdded() {
	r.rw.RLock()
	defer r.rw.RUnlock()

	for _, n := range r.nodeIndex {
		if n == nil {
			continue
		}
		n.onAddedToPopulation(true)
	}

	for _, n := range r.dynamicNodes {
		n.onAddedToPopulation(false)
	}
}

func (r *FixedRealmPopulation) GetIndexedCount() int {
	return r.indexedCount
}

func (r *FixedRealmPopulation) GetActiveNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	np := r.population.FindProfile(id)
	if np != nil && !np.IsJoiner() {
		return r.GetNodeAppearanceByIndex(np.GetIndex().AsInt())
	}
	return nil
}

func (r *FixedRealmPopulation) GetNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	na := r.GetActiveNodeAppearance(id)
	if na != nil {
		return na
	}
	return r.GetJoinerNodeAppearance(id)
}

func (r *FixedRealmPopulation) GetNodeAppearanceByIndex(idx int) *NodeAppearance {
	return r.nodeIndex[idx]
}

func (r *FixedRealmPopulation) GetShuffledOtherNodes() []*NodeAppearance {
	return r.nodeShuffle
}

func (r *FixedRealmPopulation) GetIndexedNodes() []*NodeAppearance {
	return r.nodeIndex
}

func (r *FixedRealmPopulation) GetIndexedNodesAndHasNil() ([]*NodeAppearance, bool) {
	return r.nodeIndex, true
}

func (r *FixedRealmPopulation) SealIndexed(indexedCountLimit int) bool {
	return r.indexedCount == indexedCountLimit
}

func (r *FixedRealmPopulation) AddToDynamics(ctx context.Context, n *NodeAppearance) (*NodeAppearance, error) {
	// if !n.profile.IsJoiner() {
	//	panic("illegal value")
	// }
	return r.dynPop.AddToDynamics(ctx, n)
}

func (r *FixedRealmPopulation) CreateVectorHelper() *RealmVectorHelper {

	v := r.DynamicRealmPopulation.CreateVectorHelper()
	v.realmPopulation = r
	return v
}

func (r *FixedRealmPopulation) appendDynamicNodes(nodes []*NodeAppearance) []*NodeAppearance {

	r.rw.RLock()
	defer r.rw.RUnlock()

	index := len(nodes)
	nodes = append(nodes, make([]*NodeAppearance, len(r.dynamicNodes))...)
	for _, v := range r.dynamicNodes {
		nodes[index] = v
		index++
	}
	return nodes
}

func (r *FixedRealmPopulation) GetAnyNodes(includeIndexed bool, shuffle bool) []*NodeAppearance {

	var nodes []*NodeAppearance
	joinerCount := r.GetJoinersCount()

	if !shuffle {
		if includeIndexed {
			nodes = append(make([]*NodeAppearance, 0, r.indexedCount+joinerCount), r.nodeIndex...)
		}
		nodes = r.appendDynamicNodes(nodes)
		return nodes
	}

	if includeIndexed {
		nodes = append(make([]*NodeAppearance, 0, r.indexedCount+joinerCount), r.nodeShuffle...)
		before := len(nodes)

		if !r.self.IsJoiner() {
			nodes = append(nodes, r.self)
		}
		nodes = r.appendDynamicNodes(nodes)
		if len(nodes) > before+reshuffleTolerance {
			ShuffleNodeAppearances(r.shuffleFunc, nodes)
		}
	}
	return nodes
}
