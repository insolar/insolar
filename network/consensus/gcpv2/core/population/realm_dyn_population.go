// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package population

import (
	"context"
	"fmt"
	"sync"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"

	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

func NewDynamicRealmPopulation(population census.OnlinePopulation, nodeCountHint int, phase2ExtLimit uint8,
	shuffleFn args.ShuffleFunc, baselineWeight uint32, hookCfg SharedNodeContext, fn DispatchFactoryFunc) *DynamicRealmPopulation {

	nodeCount := population.GetIndexedCapacity()
	if nodeCount > nodeCountHint {
		nodeCountHint = nodeCount + nodeCount>>2 // 125%
	}

	r := &DynamicRealmPopulation{
		dispatchFactory: fn,
		shuffleFunc:     shuffleFn,
		baselineWeight:  baselineWeight,
		phase2ExtLimit:  phase2ExtLimit,
		hook:            NewHook(nil, nil, hookCfg),
	}
	r.initPopulation(population, nodeCountHint)

	return r
}

var _ RealmPopulation = &DynamicRealmPopulation{}

const reshuffleTolerance int = 3

type DispatchFactoryFunc func(ctx context.Context, n *NodeAppearance) []DispatchMemberPacketFunc

type DynamicRealmPopulation struct {
	dispatchFactory DispatchFactoryFunc
	shuffleFunc     args.ShuffleFunc
	baselineWeight  uint32
	phase2ExtLimit  uint8
	self            *NodeAppearance
	eventStats      AtomicEventStats

	rw sync.RWMutex

	hook         Hook
	externalSink EventDispatcher

	joinerCount   int
	indexedCount  int
	indexedLenSet bool
	shuffledCount int

	nodeIndex    []*NodeAppearance
	nodeShuffle  []*NodeAppearance // excluding self
	dynamicNodes map[insolar.ShortNodeID]*NodeAppearance
	reservations int
	// voters int
}

func (p *DynamicRealmPopulation) SealIndexed(indexedCountLimit int) bool {
	p.rw.Lock()
	defer p.rw.Unlock()

	if p.indexedLenSet {
		return len(p.nodeIndex) == indexedCountLimit
	}
	if len(p.nodeIndex) > indexedCountLimit {
		return false
	}

	if len(p.nodeIndex) != indexedCountLimit {
		cp := make([]*NodeAppearance, indexedCountLimit)
		copy(cp, p.nodeIndex)
		p.nodeIndex = cp
	}

	if indexedCountLimit > cap(p.nodeShuffle) {
		p.nodeShuffle = append(make([]*NodeAppearance, 0, indexedCountLimit), p.nodeShuffle...)
	}

	p.indexedLenSet = true
	// if r.indexedCount == indexedCountLimit {
	//	r.onDynamicPopulationCompleted()
	// }
	return true
}

func (p *DynamicRealmPopulation) initPopulation(population census.OnlinePopulation, nodeCountHint int) {

	p.dynamicNodes = make(map[insolar.ShortNodeID]*NodeAppearance, nodeCountHint)
	p.nodeIndex = make([]*NodeAppearance, 0, nodeCountHint)
	p.nodeShuffle = make([]*NodeAppearance, 0, nodeCountHint)

	local := population.GetLocalProfile()

	ctx := context.Background()

	selfNode := NewNodeAppearanceAsSelf(local, power.EmptyRequest, nil)
	p.self = &selfNode
	_, p.self = p._addToDynamics(ctx, p.self, true, false)

	for _, np := range population.GetProfiles() {
		node := NewEmptyNodeAppearance(np)
		_, _ = p._addToDynamics(ctx, &node, false, false) // repeated addition will leave the initial node
	}

	self := p.dynamicNodes[local.GetNodeID()]
	if p.self == nil || p.self != self {
		panic("illegal state")
	}
	p.initHook()
}

func (p *DynamicRealmPopulation) initHook() {
	p.hook.local = p.self.profile
	p.hook.internalPopulationEventDispatcher = &EventWrapper{p.dispatchEvent}
}

func (p *DynamicRealmPopulation) dispatchEvent(fn EventClosureFunc) {
	fn(&p.eventStats)

	sink := p.externalSink
	if sink != nil {
		fn(sink)
	}
}

func (p *DynamicRealmPopulation) GetTrustCounts() (fraudCount, bySelfCount, bySomeCount, byNeighborsCount uint16) {
	return p.eventStats.GetTrustCounts()
}

func (p *DynamicRealmPopulation) GetDynamicCounts() (briefCount, fullCount uint16) {
	return p.eventStats.GetDynamicCounts()
}

func (p *DynamicRealmPopulation) GetPurgatoryCounts() (addedCount, ascentCount uint16) {
	return p.eventStats.GetPurgatoryCounts()
}

func (p *DynamicRealmPopulation) GetHook() *Hook {
	return &p.hook
}

func (p *DynamicRealmPopulation) NotifyAllOnAdded() {
	p.rw.RLock()
	defer p.rw.RUnlock()

	for _, n := range p.dynamicNodes {
		n.onAddedToPopulation(false)
	}
}

func (p *DynamicRealmPopulation) GetSelf() *NodeAppearance {
	return p.self
}

func (p *DynamicRealmPopulation) GetIndexedCount() int {
	p.rw.RLock()
	defer p.rw.RUnlock()
	return len(p.nodeIndex)
}

func (p *DynamicRealmPopulation) GetJoinersCount() int {
	p.rw.RLock()
	defer p.rw.RUnlock()
	return p.joinerCount
}

func (p *DynamicRealmPopulation) GetNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	p.rw.RLock()
	defer p.rw.RUnlock()

	return p.dynamicNodes[id]
}

func (p *DynamicRealmPopulation) GetActiveNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	na := p.GetNodeAppearance(id)
	if na == nil || na.GetProfile().IsJoiner() {
		return nil
	}
	return na
}

func (p *DynamicRealmPopulation) GetJoinerNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	na := p.GetNodeAppearance(id)
	if na == nil || !na.GetProfile().IsJoiner() {
		return nil
	}
	return na
}

func (p *DynamicRealmPopulation) GetNodeAppearanceByIndex(idx int) *NodeAppearance {
	if idx < 0 {
		panic("illegal value")
	}

	p.rw.RLock()
	defer p.rw.RUnlock()

	if idx >= len(p.nodeIndex) {
		return nil
	}
	return p.nodeIndex[idx]
}

func (p *DynamicRealmPopulation) readShuffledOtherNodes() (bool, []*NodeAppearance) {
	p.rw.RLock()
	defer p.rw.RUnlock()

	if p.shuffleFunc == nil || p.shuffledCount+reshuffleTolerance >= len(p.nodeShuffle) {
		return true, p.nodeShuffle
	}
	return false, nil
}

func (p *DynamicRealmPopulation) GetShuffledOtherNodes() []*NodeAppearance {

	ok, nodes := p.readShuffledOtherNodes()
	if ok {
		return nodes
	}

	p.rw.Lock()
	defer p.rw.Unlock()

	if p.shuffledCount+3 >= len(p.nodeShuffle) {
		return p.nodeShuffle
	}

	cp := append(make([]*NodeAppearance, 0, len(p.nodeShuffle)), p.nodeShuffle...)
	ShuffleNodeAppearances(p.shuffleFunc, cp)
	p.shuffledCount = len(p.nodeShuffle)
	p.nodeShuffle = cp
	return cp
}

func ShuffleNodeAppearances(shuffleFunc args.ShuffleFunc, nodeRefs []*NodeAppearance) {
	shuffleFunc(len(nodeRefs),
		func(i, j int) { nodeRefs[i], nodeRefs[j] = nodeRefs[j], nodeRefs[i] })
}

func (p *DynamicRealmPopulation) GetIndexedNodes() []*NodeAppearance {
	cp, _ := p.GetIndexedNodesAndHasNil()
	return cp
}

func (p *DynamicRealmPopulation) GetIndexedNodesAndHasNil() ([]*NodeAppearance, bool) {
	p.rw.RLock()
	defer p.rw.RUnlock()

	cp := make([]*NodeAppearance, len(p.nodeIndex))
	copy(cp, p.nodeIndex)

	return cp, len(p.nodeIndex) > p.indexedCount
}

func (p *DynamicRealmPopulation) GetSealedCapacity() (int, bool) {
	p.rw.RLock()
	defer p.rw.RUnlock()

	return len(p.nodeIndex), p.indexedLenSet
}

func (p *DynamicRealmPopulation) GetIndexedCountAndCompleteness() (int, bool) {
	p.rw.RLock()
	defer p.rw.RUnlock()

	return p.indexedCount, p.indexedLenSet && len(p.nodeIndex) == p.indexedCount
}

func (p *DynamicRealmPopulation) CreatePacketLimiter(isJoiner bool) phases.PacketLimiter {
	pl := phases.NewPacketLimiter(p.phase2ExtLimit)
	if isJoiner {
		return pl.ForJoiner()
	}
	return pl
}

func (p *DynamicRealmPopulation) AddToDynamics(ctx context.Context, na *NodeAppearance) (*NodeAppearance, error) {

	added, nna := p._addToDynamics(ctx, na, false, true)

	if !added && !profiles.EqualBriefProfiles(nna.GetStatic(), na.GetStatic()) {
		return nil, fmt.Errorf("multiple joiners on same id(%v): %v", na.GetNodeID(), []*NodeAppearance{na, nna})
	}

	return nna, nil
}

func (p *DynamicRealmPopulation) AddReservation(id insolar.ShortNodeID) (bool, *NodeAppearance) {

	p.rw.Lock()
	defer p.rw.Unlock()

	na, ok := p.dynamicNodes[id]
	if ok || na != nil {
		return false, na
	}

	p.dynamicNodes[id] = nil
	p.reservations++
	return true, nil
}

func (p *DynamicRealmPopulation) FindReservation(id insolar.ShortNodeID) (bool, *NodeAppearance) {

	p.rw.RLock()
	defer p.rw.RUnlock()

	na, ok := p.dynamicNodes[id]
	return ok && na == nil, na
}

func (p *DynamicRealmPopulation) _addToDynamics(ctx context.Context, n *NodeAppearance, isLocal, callOnAdded bool) (bool, *NodeAppearance) {

	if !isLocal && n.IsJoiner() && n.GetAnnouncementAsJoiner() == nil {
		panic("illegal state")
	}

	nip := n.GetStatic()
	handlers := p.dispatchFactory(ctx, n)

	p.rw.Lock()
	defer p.rw.Unlock()

	id := nip.GetStaticNodeID()
	na, ok := p.dynamicNodes[id]
	if na != nil {
		return false, na
	}
	if ok {
		p.reservations--
	}

	n.handlers = handlers
	n.hook = &p.hook
	n.neighbourWeight = p.baselineWeight

	n.limiter = p.CreatePacketLimiter(n.IsJoiner()).MergeSent(n.limiter)

	if n.IsJoiner() {
		p.joinerCount++
	} else {
		ni := n.GetIndex()
		switch {
		case ni.AsInt() == len(p.nodeIndex):
			p.nodeIndex = append(p.nodeIndex, n)
		case ni.AsInt() > len(p.nodeIndex):
			p.nodeIndex = append(p.nodeIndex, make([]*NodeAppearance, 1+ni.AsInt()-len(p.nodeIndex))...)
			p.nodeIndex[ni] = n
		default:
			if p.nodeIndex[ni] != nil {
				panic(fmt.Sprintf("duplicate node index(%v)", ni))
			}
			p.nodeIndex[ni] = n
		}
		p.indexedCount++
		p.nodeShuffle = append(p.nodeShuffle, n)
	}

	if callOnAdded {
		p._addNodeToMap(id, n)
	} else {
		p.dynamicNodes[id] = n
	}

	// if r.indexedLenSet && r.indexedCount == len(r.nodeIndex) {
	//	r.onDynamicPopulationCompleted()
	// }

	return true, n
}

func (p *DynamicRealmPopulation) _addNodeToMap(id insolar.ShortNodeID, n *NodeAppearance) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	p.dynamicNodes[id] = n
	n.onAddedToPopulation(false)
}

func (p *DynamicRealmPopulation) appendAnyNodes(includeIndexed bool, nodes []*NodeAppearance) []*NodeAppearance {

	p.rw.RLock()
	defer p.rw.RUnlock()

	delta := len(p.dynamicNodes)
	if includeIndexed {
		delta -= p.indexedCount
	}
	if delta < 0 {
		panic("illegal state")
	}
	if delta == 0 {
		return nodes
	}

	index := len(nodes)
	nodes = append(nodes, make([]*NodeAppearance, delta)...)
	for _, v := range p.dynamicNodes {
		if !includeIndexed && !v.IsJoiner() {
			continue
		}
		nodes[index] = v
		index++
	}
	return nodes
}

func (p *DynamicRealmPopulation) GetAnyNodes(includeIndexed bool, shuffle bool) []*NodeAppearance {

	nodes := p.appendAnyNodes(includeIndexed, nil)
	if shuffle && p.shuffleFunc != nil {
		ShuffleNodeAppearances(p.shuffleFunc, nodes)
	}
	return nodes
}

func (p *DynamicRealmPopulation) CreateVectorHelper() *RealmVectorHelper {
	p.rw.RLock()
	defer p.rw.RUnlock()

	v := &RealmVectorHelper{realmPopulation: p}
	v.setNodes(p.nodeIndex, p.dynamicNodes, p.hook.GetPopulationVersion())
	v.realmPopulation = p
	return v
}

/* must be set before parallel use */
func (p *DynamicRealmPopulation) InitCallback(callback EventDispatcher) {
	if p.externalSink != nil {
		panic("illegal state")
	}
	p.externalSink = callback
}
