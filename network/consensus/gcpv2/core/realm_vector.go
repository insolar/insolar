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
	//prepareFilter PrepareFilterFunc

	indexed       []VectorEntry
	poweredSorted []sortedEntry
	//nopowerSorted []sortedEntry
}

type VectorScanner interface {
}

/*
Contains copy of NodeAppearance fields that could be changed to avoid possible racing
*/

type VectorEntryData struct {
	NodeID common2.ShortNodeID

	Role       common.NodePrimaryRole
	Power      common.MemberPower
	TrustLevel packets.NodeTrustLevel
	Mode       common.MemberMode

	Node *NodeAppearance
	common.NodeAnnouncedState
}

type VectorEntry struct {
	VectorEntryData
	filterBy *VectorEntry
	//joiner 	 *NodeAppearance
}

type PrepareFilterFunc func(index int, nodeData VectorEntryData) uint16

func NewRealmVectorHelper() *RealmVectorHelper {
	return &RealmVectorHelper{}
}

func (p *RealmVectorHelper) setNodes(nodeIndex []*NodeAppearance, joinerCount int, populationVersion uint32) {

	indCount := len(nodeIndex)
	if joinerCount < 0 {
		joinerCount = indCount
	}

	p.populationVersion = populationVersion
	p.indexed = make([]VectorEntry, indCount)
	p.poweredSorted = make([]sortedEntry, indCount+joinerCount) //avoid hitting the bound at the setValue call

	sortedCount := 0
	for i, n := range nodeIndex {
		if n == nil {
			continue
		}

		joiner := p.addNode(i, n, sortedCount, nil)
		sortedCount++

		if joiner != nil {
			p.addNode(i+indCount, n, sortedCount, &p.indexed[i])
			sortedCount++
		}
	}

	p.poweredSorted = p.poweredSorted[:sortedCount]
	sort.Sort(&vectorSorter{values: p.poweredSorted})
}

func (p *RealmVectorHelper) addNode(index int, n *NodeAppearance, sortedCount int, filterBy *VectorEntry) *NodeAppearance {
	ve := &p.indexed[index]
	joiner := ve.setValues(n, uint16(index))
	ve.filterBy = filterBy

	p.poweredSorted[sortedCount].setValues(ve, index)
	return joiner
}

func (p *VectorEntry) setValues(n *NodeAppearance, index uint16) *NodeAppearance {

	p.Node = n
	p.NodeID = n.profile.GetShortNodeID()
	p.Role = n.profile.GetPrimaryRole()
	leaving, _, joiner, membership, trust := n.GetRequestedState()

	p.TrustLevel = trust
	p.StateEvidence = membership.StateEvidence
	p.AnnounceSignature = membership.AnnounceSignature

	if leaving || p.StateEvidence == nil {
		return nil
	}
	p.Power = membership.RequestedPower
	return joiner
}

//func (p *VectorEntry) HasJoiner() bool {
//	return p.joiner != nil && p.joiner != p.node
//}
//
//func (p *VectorEntry) IsLeaving() bool {
//	return p.joiner == p.node
//}

func (p *RealmVectorHelper) setOrUpdateNodes(nodeIndex []*NodeAppearance, joinerCount int, populationVersion uint32) *RealmVectorHelper {
	if p.HasSameVersion(populationVersion) {
		return p
	}

	//TODO rescan and update existing entries for possible reuse of hashing data?
	v := NewRealmVectorHelper()
	v.setNodes(nodeIndex, joinerCount, populationVersion)
	return v
}

func (p *RealmVectorHelper) HasSameVersion(version uint32) bool {
	return len(p.indexed) > 0 && p.populationVersion == version
}

func (p *RealmVectorHelper) GetIndexedCount() int {
	return len(p.indexed)
}

func (p *RealmVectorHelper) GetSortedCount() int {
	return len(p.poweredSorted)
}

//func (p *RealmVectorHelper) ScanVector(scanner VectorScanner, fn PrepareFilterFunc) {
//	for _, se := range p.poweredSorted {
//
//		filter := fn(ve.node.GetIndex(), ve.nodeData)
//
//		if se.id == ve.rank.id {
//			// member
//			scanner.ScanEntry(filter, ve.role, ve.node, ve.rank.GetPower(), ve.nodeData.NodeAnnouncedState)
//			continue
//		}
//
//		// joiner
//		if !ve.HasJoiner() {
//			panic("illegal state")
//		}
//		n := ve.joiner
//		scanner.ScanEntry(filter, ve.joinerRole, n, se.GetPower(), n.GetNodeMembershipProfileOrEmpty().NodeAnnouncedState)
//	}
//}

type entryRank struct {
	id        common2.ShortNodeID
	powerRole uint16
}

type sortedEntry struct {
	entryRank
	index uint16 //points to the same for both member and joiner, but joiner has different id in the entryRank
}

func (v sortedEntry) GetEntryIndex() int {
	return int(v.index)
}

func (p *entryRank) GetPower() common.MemberPower {
	return common.MemberPower(p.powerRole & 0xFF)
}

func (v *sortedEntry) setValues(ve *VectorEntry, index int) {
	v.id = ve.NodeID
	v.index = uint16(index)
	//if v.IsSuspended() || ve.Power == 0 {
	//	v.powerRole = 0
	//} else {
	v.powerRole = uint16(ve.Power) | uint16(ve.Role)<<8
	//}
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
	FilterMask   uint16
	FilterRes    uint16
	Index        uint8
	ReuseBuilder int8
	LastFilter   bool
}

type DualVectorBuilder struct {
	builders [2]filteredVectorBuilder
	//lazyBuilders []*filteredVectorBuilder
	//bitVectors []byte
}

func newFilteredVectorBuilders(baseBuilder VectorBuilder, masks []uint16) ([]filteredVectorBuilder, []*filteredVectorBuilder) {
	if len(masks) != 2 && len(masks) != 4 {
		panic("illegal value")
	}

	builders := make([]filteredVectorBuilder, len(masks)/2)
	builders[0].Builder = baseBuilder
	for i := range builders {
		builders[i].FilterMask = masks[i*2]
		builders[i].FilterRes = masks[i*2+1]
		builders[i].Index = uint8(i)
	}
	if len(builders) == 1 {
		return builders, nil
	}

	lazyBuilders := make([]*filteredVectorBuilder, len(builders)-1)
	for i := range lazyBuilders {
		lazyBuilders[i] = &builders[i+1]
	}
	return builders, lazyBuilders
}

//func NewDualVectorBuilder(baseBuilder VectorBuilder, masks ...uint16) DualVectorBuilder {
//	builders, _ := newFilteredVectorBuilders(baseBuilder, masks)
//	return DualVectorBuilder{
//		builders: [...]filteredVectorBuilder{builders[0], builders[1]},
//		//lazyBuilders: lazyBuilders,
//	}
//}
//
//func (p *DualVectorBuilder) Get(index int) (int, VectorBuilder, VectorCursor) {
//	b := &p.builders[index]
//	if b.Builder != nil {
//		return -1, b.Builder, b.Cursor
//	}
//	return int(b.ReuseBuilder), b.Builder, b.Cursor
//}
//
//func (p *DualVectorBuilder) ScanVector(helper *RealmVectorHelper, fn PrepareFilterFunc) {
//	for _, se := range helper.poweredSorted {
//		ve := &helper.indexed[se.index]
//
//		filter := fn(ve.node.GetIndex(), ve.nodeData)
//
//		if se.id == ve.rank.id {
//			// member
//			p.scanVectorEntry(filter, ve.role, ve.node, ve.rank.GetPower(), ve.nodeData.NodeAnnouncedState)
//			continue
//		}
//
//		// joiner
//		if !ve.HasJoiner() {
//			panic("illegal state")
//		}
//		n := ve.joiner
//		p.scanVectorEntry(filter, ve.joinerRole, n, se.GetPower(), n.GetNodeMembershipProfileOrEmpty().NodeAnnouncedState)
//	}
//}
//
//func (p *DualVectorBuilder) scanVectorEntry(filter uint16, role common.NodePrimaryRole, n *NodeAppearance, power common.MemberPower, nodeData common.NodeAnnouncedState) {
//	for i := range p.builders {
//		b := &p.builders[i]
//		b.LastFilter = b.FilterMask&filter == b.FilterRes
//	}
//	if p.builders[1].Builder == nil && p.builders[1].LastFilter != p.builders[0].LastFilter {
//		p.builders[1].Builder = p.builders[0].Builder.Fork()
//		p.builders[1].Cursor = p.builders[0].Cursor
//	}
//	for i := range p.builders {
//		b := &p.builders[i]
//		if !b.LastFilter || b.Builder == nil {
//			continue
//		}
//		b.Cursor.BeforeNext(role)
//		b.Builder.AddEntry(n, power, b.Cursor, &nodeData)
//		b.Cursor.AfterNext(power)
//	}
//}
