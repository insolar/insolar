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
	"sort"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network/consensus/gcpv2/nodeset"
)

type RealmVectorHelper struct {
	populationVersion uint32

	indexed []VectorEntry
	joiners []VectorEntry

	poweredSorted []sortedEntry
}

type VectorEntry struct {
	nodeset.VectorEntryData
	filterBy uint16
	// joiner 	 *NodeAppearance
}

func NewRealmVectorHelper() *RealmVectorHelper {
	return &RealmVectorHelper{}
}

func (p *RealmVectorHelper) SetOrUpdateNodes(nodeIndex []*NodeAppearance, joinerCount int, populationVersion uint32) *RealmVectorHelper {
	if p.HasSameVersion(populationVersion) {
		return p
	}

	// TODO rescan and update existing entries for possible reuse of hashing data?
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

func (p *RealmVectorHelper) ScanIndexed(apply func(index int, nodeData nodeset.VectorEntryData)) {
	for i := range p.indexed {
		apply(i, p.indexed[i].VectorEntryData)
	}
}

func (p *RealmVectorHelper) ScanSorted(apply func(nodeData nodeset.VectorEntryData, filter uint32), filterValue uint32) {
	for _, se := range p.poweredSorted {
		_, ve := se.chooseEntry(p.indexed, p.joiners)
		apply(ve.VectorEntryData, filterValue)
	}
}

func (p *RealmVectorHelper) ScanSortedWithFilter(apply func(nodeData nodeset.VectorEntryData, filter uint32),
	filter func(index int, nodeData nodeset.VectorEntryData) (bool, uint32)) {

	type postponedEntry struct {
		ve     *VectorEntry
		filter uint32
	}
	var skipped []postponedEntry

	for _, se := range p.poweredSorted {
		_, ve := se.chooseEntry(p.indexed, p.joiners)
		postpone, filterValue := filter(int(ve.filterBy), p.indexed[ve.filterBy].VectorEntryData)

		if postpone {
			if skipped == nil {
				skipped = make([]postponedEntry, 1, 1+len(p.poweredSorted)>>1)
				skipped[0] = postponedEntry{ve, filterValue}
			} else {
				skipped = append(skipped, postponedEntry{ve, filterValue})
			}
			continue
		}
		apply(ve.VectorEntryData, filterValue)
	}

	for _, pe := range skipped {
		apply(pe.ve.VectorEntryData, pe.filter)
	}
}

func (p *RealmVectorHelper) setNodes(nodeIndex []*NodeAppearance, joinerCountHint int, populationVersion uint32) {

	indCount := len(nodeIndex)
	if joinerCountHint < 0 {
		joinerCountHint = indCount
	}

	p.populationVersion = populationVersion
	p.indexed = make([]VectorEntry, indCount)
	p.joiners = make([]VectorEntry, joinerCountHint)
	p.poweredSorted = make([]sortedEntry, indCount+joinerCountHint)

	sortedCount := 0
	joinerCount := 0

	for i, n := range nodeIndex {
		if n == nil {
			continue
		}

		ve := &p.indexed[i]
		joiner := ve.setValues(n)
		ve.filterBy = uint16(i)

		p.poweredSorted[sortedCount].setValues(ve, int16(i))
		sortedCount++

		if joiner == nil {
			continue
		}

		if joinerCount >= len(p.joiners) {
			// got more joiners than expected - it is possible
			p.joiners = append(p.joiners, VectorEntry{})
			p.poweredSorted = append(p.poweredSorted, sortedEntry{})
		}
		je := &p.joiners[joinerCount]
		joinerCount++
		joiner = je.setValues(n)
		if joiner != nil {
			panic("illegal state")
		}

		je.filterBy = uint16(i)

		p.poweredSorted[sortedCount].setValues(je, -int16(joinerCount))
		sortedCount++
	}

	p.poweredSorted = p.poweredSorted[:sortedCount]
	p.joiners = p.joiners[:joinerCount]
	sort.Sort(&vectorPowerSorter{values: p.poweredSorted})
}

func (p *RealmVectorHelper) GetEntry(index int) nodeset.VectorEntryData {
	return p.indexed[index].VectorEntryData
}

func (p *VectorEntry) setValues(n *NodeAppearance) *NodeAppearance {

	np := n.GetProfile()
	p.NodeID = np.GetShortNodeID()
	p.Role = np.GetPrimaryRole()
	leaving, _, joiner, membership, trust := n.GetRequestedState()

	p.TrustLevel = trust
	p.StateEvidence = membership.StateEvidence
	p.AnnounceSignature = membership.AnnounceSignature

	if leaving || p.StateEvidence == nil {
		return nil
	}
	p.RequestedPower = membership.RequestedPower
	return joiner
}

type sortedEntry struct {
	id        insolar.ShortNodeID
	powerRole uint16
	index     int16 // points to the same for both member and joiner, but joiner has different id in the entryRank
}

func (v *sortedEntry) isJoiner() bool {
	return v.index < 0
}

func (v *sortedEntry) chooseEntry(indexed, joiners []VectorEntry) (bool, *VectorEntry) {
	if v.isJoiner() {
		return true, &joiners[-(v.index + 1)]
	}
	return false, &indexed[v.index]
}

func (v *sortedEntry) setValues(ve *VectorEntry, index int16) {
	v.id = ve.NodeID
	v.index = index
	// role of zero-power nodes is ignored for sorting
	if ve.RequestedPower == 0 {
		v.powerRole = 0
	} else {
		v.powerRole = uint16(ve.RequestedPower) | uint16(ve.Role)<<8
	}
}

func (v sortedEntry) lessByPowerRole(o sortedEntry) bool {
	if v.powerRole > o.powerRole {
		return false
	}
	if v.powerRole < o.powerRole {
		return true
	}
	return v.id < o.id
}

type vectorPowerSorter struct {
	values []sortedEntry
}

func (c *vectorPowerSorter) Len() int {
	return len(c.values)
}

// sorting is REVERSED - it makes the most powerful nodes of a role to be first in the list
func (c *vectorPowerSorter) Less(i, j int) bool {
	return c.values[j].lessByPowerRole(c.values[i])
}

func (c *vectorPowerSorter) Swap(i, j int) {
	c.values[i], c.values[j] = c.values[j], c.values[i]
}
