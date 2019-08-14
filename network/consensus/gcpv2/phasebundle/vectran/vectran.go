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

package vectran

import (
	"github.com/insolar/insolar/network/consensus/common/longbits"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
)

func NewMatchingVectorAnalyzer(entryCount int) VectorAnalyzer {
	return VectorAnalyzer{
		byMasks:        make(map[longbits.ByteString]*MaskInfo, 1+entryCount>>5),
		indexedEntries: make([]VectorEntry, entryCount),
	}
}

type VectorAnalyzer struct {
	byMasks        map[longbits.ByteString]*MaskInfo
	indexedEntries []VectorEntry
	totalMissCount uint16
}

type MaskInfo struct {
	mask              longbits.ByteString
	maskBits          uint16
	memberCount       uint16
	missingEntryCount uint16
	nextMissMask      *MaskInfo

	isReady         bool
	calculatedEntry CalculatedEntry
	pendingEntries  []InspectedEntry
}

type VectorEntry struct {
	nextMiss *MaskInfo
	data     nodeset.VectorEntryData
}

type CalculatedEntry interface {
}

type InspectedEntry interface {
}

func (p *VectorAnalyzer) AddInspectedEntry(mask longbits.ByteString, maskBits uint16, entry InspectedEntry) {
	if maskBits == 0 {
		panic("illegal value")
	}
	existing := p.byMasks[mask]
	if existing == nil {
		existing = p.fillByMask(mask, maskBits)
		p.byMasks[mask] = existing
	}

	existing.AddInspectedEntry(entry, p)
}

func (p *VectorAnalyzer) AddVectorEntry(entryData nodeset.VectorEntryData) {
	index := entryData.Profile.GetIndex()
	entry := &p.indexedEntries[index]
	if !entry.data.IsEmpty() {
		panic("illegal state")
	}
	entry.data = entryData

	head := entry.nextMiss
	entry.nextMiss = nil

	for head != nil {
		next := head.nextMissMask
		head.applyEntryUpdate(index.AsInt(), entryData, p)
		head = next
	}
}

func (p *VectorAnalyzer) fillByMask(mask longbits.ByteString, maskBits uint16) *MaskInfo {

	maskInfo := &MaskInfo{mask: mask, maskBits: maskBits}

	var dataEntries []*nodeset.VectorEntryData
	if maskBits <= uint16(len(p.indexedEntries))-p.totalMissCount {
		dataEntries = make([]*nodeset.VectorEntryData, 0, maskBits)
	}

	for nextBit := 0; nextBit >= 0; nextBit = mask.SearchBit(nextBit+1, true) {
		entry := &p.indexedEntries[nextBit]

		if entry.data.IsEmpty() {
			if p.totalMissCount == 0 {
				panic("illegal state")
			}

			maskInfo.missingEntryCount++

			if maskInfo.nextMissMask == nil {
				maskInfo.nextMissMask = entry.nextMiss
				entry.nextMiss = maskInfo
			}
			dataEntries = nil
		} else if cap(dataEntries) > 0 {
			dataEntries = append(dataEntries, &entry.data)
		}
		maskBits--
		if maskBits == 0 {
			break
		}
	}

	if maskInfo.missingEntryCount == 0 {
		if len(dataEntries) != int(maskInfo.maskBits) {
			panic("illegal state")
		}
		p.CalculateVector(dataEntries, maskInfo)
	}

	return maskInfo
}

func (p *VectorAnalyzer) CalculateVector(dataEntries []*nodeset.VectorEntryData, m *MaskInfo) {
	// do calc
	//m.setReady(calculatedEntry, p)
}

func (p *VectorAnalyzer) InspectEntry(entry InspectedEntry, calcEntry CalculatedEntry) {
	// do calc
}

func (p *MaskInfo) setReady(calcEntry CalculatedEntry, va *VectorAnalyzer) {
	if p.isReady {
		panic("illegal state")
	}
	p.calculatedEntry = calcEntry
	p.isReady = true

	pendingEntries := p.pendingEntries
	p.pendingEntries = nil
	for _, pending := range pendingEntries {
		va.InspectEntry(pending, p.calculatedEntry)
	}
}

func (p *MaskInfo) applyEntryUpdate(index int, data nodeset.VectorEntryData, va *VectorAnalyzer) {

	nextBit := index

outer:
	for {
		switch p.missingEntryCount {
		case 0:
			panic("illegal state")
		case 1:
			p.missingEntryCount = 0
			break outer
		default:
			p.missingEntryCount--
			nextBit := p.mask.SearchBit(nextBit+1, true)
			if nextBit < 0 {
				panic("illegal state")
			}

			entry := &va.indexedEntries[nextBit]

			if entry.data.IsEmpty() {
				p.nextMissMask = entry.nextMiss
				entry.nextMiss = p
				return
			}
		}
	}

	maskBits := p.maskBits
	dataEntries := make([]*nodeset.VectorEntryData, 0, maskBits)

	for nextBit := 0; nextBit >= 0; nextBit = p.mask.SearchBit(nextBit+1, true) {
		entry := &va.indexedEntries[nextBit]

		if entry.data.IsEmpty() {
			panic("illegal state")
		}
		dataEntries = append(dataEntries, &entry.data)

		maskBits--
		if maskBits == 0 {
			break
		}
	}

	va.CalculateVector(dataEntries, p)
}

func (p *MaskInfo) AddInspectedEntry(entry InspectedEntry, va *VectorAnalyzer) {

	p.memberCount++
	if p.isReady {
		va.InspectEntry(entry, p.calculatedEntry)
	} else {
		p.pendingEntries = append(p.pendingEntries, entry)
	}
}
