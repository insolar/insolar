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
	common2 "github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
	"sort"
)

type RealmVectorHelper struct {
	populationVersion uint32

	indexed []VectorEntry
	sorted  []sortedEntry
}

/*
Contains copy of NodeAppearance fields that could be changed to avoid possible racing
*/
type VectorEntry struct {
	rank       entryRank
	role       common.NodePrimaryRole
	joinerRole common.NodePrimaryRole
	//joinerRank entryRank
	trustLevel packets.NodeTrustLevel

	filter uint8

	node   *NodeAppearance
	joiner *NodeAppearance //if joiner == node then this node is leaving

	nodeData common.NodeAnnouncedState
}

func NewRealmVectorHelper() *RealmVectorHelper {
	return &RealmVectorHelper{}
}

func (v *RealmVectorHelper) setNodes(nodeIndex []*NodeAppearance, joinerCount int, populationVersion uint32) {

	indCount := len(nodeIndex)
	if joinerCount < 0 {
		joinerCount = indCount
	}

	v.populationVersion = populationVersion
	v.indexed = make([]VectorEntry, indCount)
	v.sorted = make([]sortedEntry, indCount+indCount+1) //avoid hitting the bound

	sortedCount := indCount
	for i, n := range nodeIndex {
		ve := &v.indexed[i]
		se := &v.sorted[i]
		se.index = uint16(i)

		if ve.setValues(n, uint16(i), &v.sorted[sortedCount]) {
			sortedCount++
		}
	}

	v.sorted = v.sorted[:sortedCount]
	sort.Sort(&vectorSorter{values: v.sorted})
}

func (*RealmVectorHelper) updateNodes(nodeIndex []*NodeAppearance, joinerCount int, populationVersion uint32) *RealmVectorHelper {
	//TODO rescan and update existing entries for possible reuse of hashing data?
	v := NewRealmVectorHelper()
	v.setNodes(nodeIndex, joinerCount, populationVersion)
	return v
}

func (p *VectorEntry) setValues(n *NodeAppearance, index uint16, se *sortedEntry) bool {

	p.node = n

	if n == nil {
		p.role = 0
		p.rank.powerRole = 0
		return false
	}

	p.role = n.profile.GetPrimaryRole()
	leaving, _, joiner, membership, trust := n.GetRequestedState()

	p.nodeData.StateEvidence = membership.StateEvidence
	p.nodeData.AnnounceSignature = membership.AnnounceSignature

	p.trustLevel = trust
	if leaving {
		p.rank.id = 0
		p.rank.powerRole = 0
		p.joiner = n //indication of leaving
	} else {
		p.rank.id = n.profile.GetShortNodeID()
		if p.nodeData.StateEvidence == nil {
			p.rank.powerRole = 0
		} else {
			p.rank.powerRole = powerRoleOf(p.role, membership.RequestedPower)
		}
		p.joiner = joiner
	}

	if joiner == nil {
		return false
	}

	p.joinerRole = joiner.profile.GetPrimaryRole()
	se.id = joiner.profile.GetShortNodeID()
	_, _, _, membership, _ = n.GetRequestedState()
	se.powerRole = powerRoleOf(p.joinerRole, membership.RequestedPower)
	se.index = index
	return true
}

func (p *VectorEntry) HasJoiner() bool {
	return p.joiner != nil && p.joiner != p.node
}

func (p *VectorEntry) IsLeaving() bool {
	return p.joiner == p.node
}

type entryRank struct {
	id        common2.ShortNodeID
	powerRole uint16
}

func powerRoleOf(role common.NodePrimaryRole, power common.MemberPower) uint16 {
	if power == 0 {
		return 0
	}
	return uint16(power) | uint16(role)<<8
}

func (p *entryRank) GetPower() common.MemberPower {
	return common.MemberPower(p.powerRole & 0xFF)
}

type sortedEntry struct {
	entryRank
	index uint16 //points to the same for both member and joiner, but joiner has different id in the entryRank
}

func (v sortedEntry) GetEntryIndex() int {
	return int(v.index &^ 0x8000)
}

//role of zero-power nodes is ignored for sorting
//sorting is REVERSED - it makes the most powerful nodes of a role to be first in the list
//nodeID are also reversed to put leaving nodes (id = 0) at the very end
func (v entryRank) less(o entryRank) bool {
	if v.powerRole < o.powerRole {
		return false
	}
	if v.powerRole > o.powerRole {
		return true
	}
	return v.id > o.id
}

type vectorSorter struct {
	values []sortedEntry
}

func (c *vectorSorter) Len() int {
	return len(c.values)
}

func (c *vectorSorter) Less(i, j int) bool {
	return c.values[i].less(c.values[j].entryRank)
}

func (c *vectorSorter) Swap(i, j int) {
	c.values[i], c.values[j] = c.values[j], c.values[i]
}

type VectorCursor struct {
	NodeIndex uint16
	RoleIndex uint16
	PowIndex  uint16

	LastRole common.NodePrimaryRole
}

func (p *VectorCursor) BeforeNext(role common.NodePrimaryRole) {
	if p.LastRole == role {
		return
	}
	p.RoleIndex = 0
	p.PowIndex = 0
	p.LastRole = role
}

func (p *VectorCursor) AfterNext(power common.MemberPower) {
	p.RoleIndex++
	p.PowIndex += power.ToLinearValue()
	p.NodeIndex++
}

type VectorBuilder interface {
	AddEntry(n *NodeAppearance, power common.MemberPower, cursor VectorCursor, nodeData *common.NodeAnnouncedState)
	Fork() VectorBuilder
}

type filteredVectorBuilder struct {
	Cursor       VectorCursor
	Builder      VectorBuilder
	Index        uint8
	ReuseBuilder int8
	FilterMask   uint8
	FilterRes    uint8
	LastFilter   bool
}

type DualVectorBuilder struct {
	helper   *RealmVectorHelper
	builders [2]filteredVectorBuilder
	//lazyBuilders []*filteredVectorBuilder
	//bitVectors []byte
}

func newFilteredVectorBuilders(baseBuilder VectorBuilder, masks []uint8) ([]filteredVectorBuilder, []*filteredVectorBuilder) {
	if len(masks) == 0 || len(masks)%2 != 0 {
		panic("illegal value")
	}

	builders := make([]filteredVectorBuilder, len(masks)/2)
	builders[0].Builder = baseBuilder
	for i := range builders {
		builders[i].FilterMask = masks[i*2]
		builders[i].FilterRes = masks[i*2+1]
		builders[i].Index = uint8(i)
	}
	if len(masks) == 1 {
		return builders, nil
	}

	lazyBuilders := make([]*filteredVectorBuilder, len(masks)-1)
	for i := range lazyBuilders {
		lazyBuilders[i] = &builders[i+1]
	}
	return builders, lazyBuilders
}

func NewDualVectorBuilder(helper *RealmVectorHelper, baseBuilder VectorBuilder, masks ...uint8) DualVectorBuilder {
	builders, _ := newFilteredVectorBuilders(baseBuilder, masks)
	return DualVectorBuilder{
		helper:   helper,
		builders: [...]filteredVectorBuilder{builders[0], builders[1]},
		//lazyBuilders: lazyBuilders,
	}
}

//func (p *DualVectorBuilder) smth() {
//	for _, se := range p.helper.sorted {
//		ve := &p.helper.indexed[se.index]
//
//		var (
//			n        *NodeAppearance
//			power    common.MemberPower
//			role     common.NodePrimaryRole
//			nodeData *common.NodeAnnouncedState
//		)
//		if se.id == ve.rank.id {
//			// member
//			n = ve.node
//			power = ve.rank.GetPower()
//			role = ve.role
//			nodeData = &ve.nodeData
//		} else {
//			// joiner
//			if !ve.HasJoiner() {
//				panic("illegal state")
//			}
//			n = ve.joiner
//			power = se.GetPower()
//			nas := n.GetNodeMembershipProfileOrEmpty().NodeAnnouncedState
//			role = ve.joinerRole
//			nodeData = &nas
//		}
//
//		for i := range p.builders {
//			b := &p.builders[i]
//			b.LastFilter = b.FilterMask&ve.filter == b.FilterRes
//
//		}
//
//		if len(p.lazyBuilders) > 0 {
//			j := 0
//			for i, b := range p.lazyBuilders {
//				base := &p.builders[b.ReuseBuilder]
//				if b.LastFilter != base.LastFilter {
//					b.Builder = base.Builder.Fork()
//					for k := 0; k < i; k++ {
//						bb := p.lazyBuilders[k]
//						if bb.ReuseBuilder == b.ReuseBuilder && bb.LastFilter == b.LastFilter {
//							bb.ReuseBuilder =
//						}
//						if
//					}
//				}
//
//				//if b.LastFilter !=
//			}
//
//
//		} else {
//			for i := range p.builders {
//				b := &p.builders[i]
//				if b.FilterMask & ve.filter == 0 { continue }
//				b.Cursor.BeforeNext(role)
//				b.Builder.AddEntry(n, power, b.Cursor, nodeData)
//				b.Cursor.AfterNext(power)
//			}
//		}
//	}
//}
