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

package purgatory

type RealmPurgatory struct {
	//nodeInit       core.NodeInitFunc
	//baselineWeight uint32
	//phase2ExtLimit uint8
	//self           *core.NodeAppearance
	//
	//rw sync.RWMutex
	//
	//joinerCount  int
	//indexedCount int
	//
	//nodeIndex    []*core.NodeAppearance
	//nodeShuffle  []*core.NodeAppearance // excluding self
	//dynamicNodes map[insolar.ShortNodeID]*core.NodeAppearance
	////	purgatoryByEP map[string]*NodeAppearance
	////purgatoryByPK map[string]*core.NodeAppearance
	////purgatoryByID map[insolar.ShortNodeID]*[]*core.NodeAppearance
	////purgatoryOuts map[insolar.ShortNodeID]*core.NodeAppearance
}

//func NewRealmPurgatory(baselineWeight uint32, local profiles.ActiveNode, nodeCountHint int, phase2ExtLimit uint8,
//	fn core.NodeInitFunc) *RealmPurgatory {
//
//	r := &RealmPurgatory{
//		nodeInit:       fn,
//		baselineWeight: baselineWeight,
//		phase2ExtLimit: phase2ExtLimit,
//	}
//	r.initPopulation(local, nodeCountHint)
//
//	return r
//}
//
//
//func (r *RealmPurgatory) initPopulation(local profiles.ActiveNode, nodeCountHint int) {
//	//r.self = r.CreateNodeAppearance(context.Background(), local)
//	r.dynamicNodes = make(map[insolar.ShortNodeID]*core.NodeAppearance, nodeCountHint)
//}
//
//func (r *RealmPurgatory) GetSelf() *core.NodeAppearance {
//	return r.self
//}
//
//func (r *RealmPurgatory) GetNodeCount() int {
//	r.rw.RLock()
//	defer r.rw.RUnlock()
//	return len(r.nodeIndex)
//}
//
//func (r *RealmPurgatory) GetJoinersCount() int {
//	r.rw.RLock()
//	defer r.rw.RUnlock()
//	return r.joinerCount
//}
//
//func (r *RealmPurgatory) GetOthersCount() int {
//	return r.GetNodeCount() - 1
//}
//
//func (r *RealmPurgatory) GetBftMajorityCount() int {
//	return consensuskit.BftMajority(r.GetNodeCount())
//}
//
//func (r *RealmPurgatory) GetNodeAppearance(id insolar.ShortNodeID) *core.NodeAppearance {
//	r.rw.RLock()
//	defer r.rw.RUnlock()
//
//	return r.dynamicNodes[id]
//}
//
//func (r *RealmPurgatory) GetActiveNodeAppearance(id insolar.ShortNodeID) *core.NodeAppearance {
//	na := r.GetNodeAppearance(id)
//	if !na.GetProfile().IsJoiner() {
//		return na
//	}
//	return nil
//}
//
//func (r *RealmPurgatory) GetJoinerNodeAppearance(id insolar.ShortNodeID) *core.NodeAppearance {
//	na := r.GetNodeAppearance(id)
//	if !na.GetProfile().IsJoiner() {
//		return nil
//	}
//	return na
//}
//
//func (r *RealmPurgatory) GetNodeAppearanceByIndex(idx int) *core.NodeAppearance {
//	if idx < 0 {
//		panic("illegal value")
//	}
//
//	r.rw.RLock()
//	defer r.rw.RUnlock()
//
//	if idx >= len(r.nodeIndex) {
//		return nil
//	}
//	return r.nodeIndex[idx]
//}
//
//func (r *RealmPurgatory) GetShuffledOtherNodes() []*core.NodeAppearance {
//	r.rw.RLock()
//	defer r.rw.RUnlock()
//
//	return r.nodeShuffle
//}
//
//func (r *RealmPurgatory) IsComplete() bool {
//	r.rw.RLock()
//	defer r.rw.RUnlock()
//
//	return len(r.nodeIndex) == r.indexedCount
//}
//
//func (r *RealmPurgatory) GetIndexedNodes() []*core.NodeAppearance {
//	cp, _ := r.GetIndexedNodesWithCheck()
//	// if !ok {
//	//	panic("node set is incomplete")
//	// }
//	return cp
//}
//
//func (r *RealmPurgatory) GetIndexedNodesWithCheck() ([]*core.NodeAppearance, bool) {
//	r.rw.RLock()
//	defer r.rw.RUnlock()
//
//	cp := make([]*core.NodeAppearance, len(r.nodeIndex))
//	copy(cp, r.nodeIndex)
//
//	return cp, len(r.nodeIndex) == r.indexedCount
//}
//
////func (r *RealmPurgatory) CreateNodeAppearance(ctx context.Context, np profiles.ActiveNode) *core.NodeAppearance {
////
////	n := &core.NodeAppearance{}
////	core.init(np, nil, r.baselineWeight, r.phase2ExtLimit)
////	r.nodeInit(ctx, n)
////
////	return n
////}
////
////func (r *RealmPurgatory) AddToPurgatory(n *core.NodeAppearance) (*core.NodeAppearance, core.PurgatoryNodeState) {
////	nip := n.profile.GetStatic()
////	if nip.GetIntroduction() != nil {
////		panic("illegal value")
////	}
////
////	id := nip.GetStaticNodeID()
////	na := r.GetActiveNodeAppearance(id)
////	if na != nil {
////		return na, core.PurgatoryExistingMember
////	}
////
////	r.rw.Lock()
////	defer r.rw.Unlock()
////
////	nn := r.dynamicNodes[id]
////	if nn != nil {
////		return nn, core.PurgatoryExistingMember
////	}
////
////	if r.purgatoryByPK == nil {
////		r.purgatoryByPK = make(map[string]*core.NodeAppearance)
////		r.purgatoryByID = make(map[insolar.ShortNodeID]*[]*core.NodeAppearance)
////
////		r.purgatoryByPK[nip.GetNodePublicKey().AsByteString()] = n
////		r.purgatoryByID[nip.GetStaticNodeID()] = &[]*core.NodeAppearance{n}
////		return n, 0
////	}
////
////	pk := nip.GetNodePublicKey().AsByteString()
////	nn = r.purgatoryByPK[pk]
////	if nn != nil {
////		return nn, core.PurgatoryDuplicatePK
////	}
////
////	nodes := r.purgatoryByID[id]
////
////	if nodes == nil {
////		nodes = &[]*core.NodeAppearance{n}
////		r.purgatoryByID[id] = nodes
////		return n, 0
////	}
////	*nodes = append(*nodes, n)
////	return n, core.PurgatoryNodeState(len(*nodes) - 1)
////}
////
////func (r *RealmPurgatory) AddToDynamics(n *core.NodeAppearance) (*core.NodeAppearance, []*core.NodeAppearance) {
////	nip := n.profile.GetStatic()
////
////	if nip.GetIntroduction() == nil {
////		panic("illegal value")
////	}
////
////	r.rw.Lock()
////	defer r.rw.Unlock()
////
////	id := nip.GetStaticNodeID()
////
////	delete(r.purgatoryByPK, nip.GetNodePublicKey().AsByteString())
////	nodes := r.purgatoryByID[id]
////	if nodes != nil {
////		delete(r.purgatoryByID, id)
////	} else {
////		nodes = &[]*core.NodeAppearance{}
////	}
////
////	na := r.GetActiveNodeAppearance(id)
////	if na != nil {
////		return na, *nodes
////	}
////
////	na = r.dynamicNodes[id]
////	if na != nil {
////		return na, *nodes
////	}
////
////	if n.profile.IsJoiner() {
////		r.joinerCount++
////	} else {
////		ni := n.profile.GetIndex()
////		switch {
////		case ni.AsInt() == len(r.nodeIndex):
////			r.nodeIndex = append(r.nodeIndex, n)
////		case ni.AsInt() > len(r.nodeIndex):
////			r.nodeIndex = append(r.nodeIndex, make([]*core.NodeAppearance, 1+ni.AsInt()-len(r.nodeIndex))...)
////			r.nodeIndex[ni] = n
////		default:
////			if r.nodeIndex[ni] != nil {
////				panic(fmt.Sprintf("duplicate node id(%v)", ni))
////			}
////			r.nodeIndex[ni] = n
////		}
////		r.indexedCount++
////		r.nodeShuffle = append(r.nodeShuffle, n)
////	}
////	return n, *nodes
////}
