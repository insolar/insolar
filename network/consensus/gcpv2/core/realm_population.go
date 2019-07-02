package core

import (
	"context"
	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"sync"
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

	//	purgatory	map[common.ShortNodeID]*NodeAppearance
	rw      sync.RWMutex
	joiners map[common.ShortNodeID]*NodeAppearance

	purgatoryByPK map[string] /* used as string(PK.Bytes())  */ *NodeAppearance
	purgatoryByID map[common.ShortNodeID][]*NodeAppearance
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

func (r *MemberRealmPopulation) createNode(ctx context.Context, np common2.NodeProfile) *NodeAppearance {
	//np.GetNodePublicKey().AsByteString()

	n := &NodeAppearance{}
	n.init(np, nil, r.baselineWeight)
	r.nodeInit(ctx, n)
	return n
}
