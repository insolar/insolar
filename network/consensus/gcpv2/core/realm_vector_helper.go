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
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
)

type RealmVectorHelper struct {
	realmPopulation RealmPopulation

	members []VectorEntry
	joiners []VectorEntry

	mutex      sync.Mutex
	projection RealmVectorProjection
}

/*
Unsafe for concurrent use
*/
type RealmVectorProjection struct {
	origin              *RealmVectorHelper
	populationVersion   uint32
	indexedRefs         []*VectorEntry // can't be appended, must be copied on setting new entries
	joinersRefs         []*VectorEntry // can be appended in-place, can't set/update entries
	poweredSorted       []sortedEntry  // must be copied on append/sort, can't set/update entries
	sharedIndexedRefs   bool
	sharedPoweredSorted bool
}

type VectorEntry struct {
	nodeset.VectorEntryData
	filterBy uint16
	// joiner 	 *NodeAppearance
}

func (p *RealmVectorProjection) CreateProjection() RealmVectorProjection {

	p.sharedIndexedRefs = true
	p.sharedPoweredSorted = true
	return *p
}

func (p *RealmVectorProjection) HasSameVersion(version uint32) bool {
	return len(p.indexedRefs) > 0 && p.populationVersion == version
}

func (p *RealmVectorProjection) ForceEntryUpdate(index int) bool {
	//if p == &p.origin.projection {
	//	member, _ := p.origin.forceEntryUpdate(index)
	//	return member != nil
	//}
	//member, joiner := p.origin.forceEntryUpdate(index)
	//
	//if member == nil {
	return false
	//}
	//p.updateEntry(index, member, joiner)
	//return true
}

//func (p *RealmVectorHelper) forceEntryUpdate(index int) (*VectorEntry, *VectorEntry) {
//	p.mutex.Lock()
//	defer p.mutex.Lock()
//
//	return nil, nil
//}

func (p *RealmVectorProjection) updateEntry(index int, member, joiner *VectorEntry) { // nolint: unused
	if p.sharedIndexedRefs {
		cp := make([]*VectorEntry, len(p.indexedRefs))
		copy(cp, p.indexedRefs)
		p.indexedRefs = cp
		p.sharedIndexedRefs = false
	}
	p.indexedRefs[index] = member

	if joiner != nil {
		p.joinersRefs = append(p.joinersRefs, joiner)
	}

	pos := len(p.poweredSorted)
	if p.sharedPoweredSorted {
		newLen := pos + 1
		if joiner != nil {
			newLen++
		}
		cp := make([]sortedEntry, newLen)
		copy(cp, p.poweredSorted)
		p.poweredSorted = cp
		p.sharedPoweredSorted = false
	}
	p.poweredSorted[pos].setMember(member, index)
	if joiner != nil {
		pos++
		p.poweredSorted[pos].setJoiner(member, pos)
	}
	sort.Sort(&vectorPowerSorter{p.poweredSorted})
}

func (p *RealmVectorProjection) GetIndexedCount() int {
	return len(p.indexedRefs)
}

func (p *RealmVectorProjection) GetSortedCount() int {
	return len(p.poweredSorted)
}

func (p *RealmVectorProjection) ScanIndexed(apply func(index int, nodeData nodeset.VectorEntryData)) {
	for i := range p.indexedRefs {
		apply(i, p.indexedRefs[i].VectorEntryData)
	}
}

func (p *RealmVectorProjection) GetEntry(index int) nodeset.VectorEntryData {
	return p.indexedRefs[index].VectorEntryData
}

func (p *RealmVectorProjection) ScanSorted(apply nodeset.EntryFilteredScannerFunc, filterValue uint32) {
	for _, se := range p.poweredSorted {
		_, ve := se.chooseEntry(p.indexedRefs, p.joinersRefs)
		apply(ve.VectorEntryData, false, filterValue)
	}
}

type postponedEntry struct {
	ve     *VectorEntry
	filter uint32
}

func (p *RealmVectorProjection) ScanSortedWithFilter(apply nodeset.EntryFilteredScannerFunc, filter nodeset.EntryFilterFunc) {

	var skipped []postponedEntry
	unorderedSkipped := false

	prevID := insolar.AbsentShortNodeID

	for _, se := range p.poweredSorted {
		joiner, valueEntry := se.chooseEntry(p.indexedRefs, p.joinersRefs)
		filterEntry := p.indexedRefs[valueEntry.filterBy]
		postpone, filterValue := filter(int(valueEntry.filterBy), filterEntry.VectorEntryData)

		nodeID := valueEntry.Profile.GetNodeID()
		if joiner {
			if postpone {
				// joiner MUST NOT appear when an introduction node is out
				return
			}
			// joiner may appear multiple times in a powered section via multiple introducing nodes
			// due to sorting all repetitions will come in sequence
			if prevID == nodeID {
				continue
			}
			postpone, _ = filter(-1, valueEntry.VectorEntryData)
		} else if prevID == nodeID {
			// regular nodes MUST NOT be multiplied
			panic("illegal state")
		}
		prevID = nodeID

		if postpone {
			if skipped == nil {
				skipped = make([]postponedEntry, 1, 1+len(p.poweredSorted)>>1)
				skipped[0] = postponedEntry{valueEntry, filterValue}
			} else {
				if skipped[len(skipped)-1].ve.Profile.GetNodeID() >= valueEntry.Profile.GetNodeID() {
					unorderedSkipped = true
				}
				skipped = append(skipped, postponedEntry{valueEntry, filterValue})
			}
			continue
		}
		apply(valueEntry.VectorEntryData, false, filterValue)
	}

	if unorderedSkipped {
		sort.Sort(&vectorIDSorter{skipped})
	}

	for _, pe := range skipped {
		apply(pe.ve.VectorEntryData, true, pe.filter)
	}
}

func (p *RealmVectorHelper) CreateProjection() RealmVectorProjection {

	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.projection.CreateProjection()
}

func (p *RealmVectorHelper) CreateUnsafeProjection() RealmVectorProjection {

	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.projection
}

func (p *RealmVectorHelper) setArrayNodes(nodeIndex []*NodeAppearance,
	dynamicNodes map[insolar.ShortNodeID]*NodeAppearance, populationVersion uint32) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.projection.origin != nil {
		panic("illegal state")
	}
	p.projection.origin = p

	indCount := len(nodeIndex)

	p.projection.populationVersion = populationVersion
	p.members = make([]VectorEntry, indCount)
	p.joiners = make([]VectorEntry, len(dynamicNodes))

	p.projection.indexedRefs = make([]*VectorEntry, indCount)
	p.projection.poweredSorted = make([]sortedEntry, indCount+len(dynamicNodes))

	sortedCount := 0
	joinerCount := 0

	for i, n := range nodeIndex {
		if n == nil {
			continue
		}

		ve := &p.members[i]
		p.projection.indexedRefs[i] = ve
		joinerID := ve.setValues(n)
		ve.filterBy = uint16(i)

		p.projection.poweredSorted[sortedCount].setMember(ve, i)
		sortedCount++

		if joinerID.IsAbsent() {
			continue
		}
		joiner := dynamicNodes[joinerID]
		if joiner == nil {
			//panic("joiner is missing")
			continue
		}

		if joinerCount >= len(p.joiners) {
			// got more joiners than expected - it is possible
			p.joiners = append(p.joiners, VectorEntry{})
			p.projection.poweredSorted = append(p.projection.poweredSorted, sortedEntry{})
		}
		je := &p.joiners[joinerCount]

		joinerID = je.setValues(joiner)
		if !joinerID.IsAbsent() {
			panic("illegal state")
		}

		je.filterBy = uint16(i)

		p.projection.poweredSorted[sortedCount].setJoiner(je, joinerCount)
		joinerCount++
		sortedCount++
	}

	p.projection.poweredSorted = p.projection.poweredSorted[:sortedCount]
	p.joiners = p.joiners[:joinerCount]
	sort.Sort(&vectorPowerSorter{p.projection.poweredSorted})

	p.projection.joinersRefs = make([]*VectorEntry, joinerCount)
	for i := range p.joiners {
		p.projection.joinersRefs[i] = &p.joiners[i]
	}
}

func (p *VectorEntry) setValues(n *NodeAppearance) insolar.ShortNodeID {

	np := n.GetProfile()
	p.Profile = np
	rs := n.GetRequestedState()

	p.TrustLevel = rs.TrustLevel
	p.StateEvidence = rs.StateEvidence
	p.AnnounceSignature = rs.AnnounceSignature
	p.RequestedMode = rs.RequestedMode
	p.RequestedPower = rs.RequestedPower
	return rs.JoinerID
}

type sortedEntry struct {
	id        insolar.ShortNodeID
	powerRole uint16
	index     int16 // points to the same for both member and joiner, but joiner has different id in the entryRank
}

func (v *sortedEntry) isJoiner() bool {
	return v.index < 0
}

func (v *sortedEntry) chooseEntry(members, joiners []*VectorEntry) (bool, *VectorEntry) {
	if v.isJoiner() {
		return true, joiners[-(v.index + 1)]
	}
	return false, members[v.index]
}

func (v *sortedEntry) setJoiner(ve *VectorEntry, index int) {
	v.setMember(ve, -(index + 1))
}

func (v *sortedEntry) setMember(ve *VectorEntry, index int) {
	v.id = ve.Profile.GetNodeID()
	v.index = int16(index)
	// role of zero-power nodes is ignored for sorting
	if ve.RequestedPower == 0 {
		v.powerRole = 0
	} else {
		v.powerRole = uint16(ve.RequestedPower) | uint16(ve.Profile.GetStatic().GetPrimaryRole())<<8
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

type vectorIDSorter struct {
	values []postponedEntry
}

func (c *vectorIDSorter) Len() int {
	return len(c.values)
}

func (c *vectorIDSorter) Less(i, j int) bool {
	return c.values[i].ve.Profile.GetNodeID() < c.values[j].ve.Profile.GetNodeID()
}

func (c *vectorIDSorter) Swap(i, j int) {
	c.values[i], c.values[j] = c.values[j], c.values[i]
}
