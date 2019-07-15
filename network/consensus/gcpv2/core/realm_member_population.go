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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/consensuskit"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

func NewMemberRealmPopulation(strategy RoundStrategy, population census.OnlinePopulation, phase2ExtLimit uint8,
	fn NodeInitFunc) *MemberRealmPopulation {

	nodeCount := population.GetCount()

	r := &MemberRealmPopulation{
		population: population,
		dynPop: dynPop{DynamicRealmPopulation{
			nodeInit:       fn,
			baselineWeight: strategy.RandUint32(),
			indexedCount:   nodeCount,
			nodeIndex:      make([]*NodeAppearance, nodeCount),
			nodeShuffle:    make([]*NodeAppearance, nodeCount-1),
		}},
		bftMajorityCount: consensuskit.BftMajority(nodeCount),
	}
	r.initPopulation(phase2ExtLimit)
	ShuffleNodeProjections(strategy, r.nodeShuffle)

	return r
}

var _ RealmPopulation = &MemberRealmPopulation{}

type dynPop struct{ DynamicRealmPopulation }

type MemberRealmPopulation struct {
	dynPop
	population census.OnlinePopulation

	bftMajorityCount int
}

func (r *MemberRealmPopulation) IsComplete() bool {
	return true
}

func (r *MemberRealmPopulation) initPopulation(phase2ExtLimit uint8) {
	activeProfiles := r.population.GetProfiles()
	thisNodeID := r.population.GetLocalProfile().GetShortNodeID()

	nodes := make([]NodeAppearance, r.indexedCount)

	var j = 0
	for i, p := range activeProfiles {
		n := &nodes[i]
		r.nodeIndex[i] = n

		if p.GetOpMode().IsEvicted() {
			panic("illegal state")
		}

		n.init(p, nil, r.baselineWeight, phase2ExtLimit)
		r.nodeInit(context.Background(), n)

		if p.GetShortNodeID() == thisNodeID {
			if r.self != nil {
				panic("schizophrenia")
			}
			r.self = n
		} else {
			if j == len(activeProfiles) {
				panic("didnt find myself among active nodes")
			}
			r.nodeShuffle[j] = n
			j++
		}
	}
}

func (r *MemberRealmPopulation) GetNodeCount() int {
	return r.indexedCount
}

func (r *MemberRealmPopulation) GetOthersCount() int {
	return r.indexedCount - 1
}

func (r *MemberRealmPopulation) GetBftMajorityCount() int {
	return r.bftMajorityCount
}

func (r *MemberRealmPopulation) GetActiveNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	np := r.population.FindProfile(id)
	if np != nil && !np.IsJoiner() {
		return r.GetNodeAppearanceByIndex(np.GetIndex())
	}
	return nil
}

func (r *MemberRealmPopulation) GetNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	na := r.GetActiveNodeAppearance(id)
	if na != nil {
		return na
	}
	return r.GetJoinerNodeAppearance(id)
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

func (r *MemberRealmPopulation) AddToDynamics(n *NodeAppearance) (*NodeAppearance, []*NodeAppearance) {
	if !n.profile.IsJoiner() {
		panic("illegal value")
	}
	return r.dynPop.AddToDynamics(n)
}

func (r *MemberRealmPopulation) CreateVectorHelper() *RealmVectorHelper {

	v := r.DynamicRealmPopulation.CreateVectorHelper()
	v.realmPopulation = r
	return v
}

var _ profiles.ActiveNode = &joiningNodeProfile{}

type joiningNodeProfile struct {
	profiles.NodeIntroProfile
}

func (p *joiningNodeProfile) IsJoiner() bool {
	return true
}

func (p *joiningNodeProfile) GetOpMode() member.OpMode {
	return member.ModeNormal
}

func (p *joiningNodeProfile) GetIndex() int {
	return 0
}

func (p *joiningNodeProfile) GetDeclaredPower() member.Power {
	return p.GetStartPower()
}

func (*joiningNodeProfile) GetSignatureVerifier() cryptkit.SignatureVerifier {
	return nil
}
