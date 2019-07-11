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

package nodeset

import (
	common2 "github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/common"
)

/*
Contains copy of NodeAppearance fields that can be changed, to avoid possible racing
*/

type VectorEntryData struct {
	NodeID common2.ShortNodeID

	Role           common.NodePrimaryRole
	RequestedPower common.MemberPower
	TrustLevel     common.NodeTrustLevel
	Mode           common.MemberOpMode

	//Node *NodeAppearance
	//common.MembershipAnnouncement
	common.NodeAnnouncedState
}

type VectorEntryScanner interface {
	GetIndexedCount() int
	GetSortedCount() int
	ScanIndexed(apply func(index int, nodeData VectorEntryData))
	ScanSorted(apply func(nodeData VectorEntryData, filter uint32), filterValue uint32)
	ScanSortedWithFilter(apply func(nodeData VectorEntryData, filter uint32),
		filter func(index int, nodeData VectorEntryData) (bool, uint32))
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
	AddEntry(cursor VectorCursor, entry VectorEntryData)
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

func NewDualVectorBuilder(baseBuilder VectorBuilder, masks ...uint16) DualVectorBuilder {
	builders, _ := newFilteredVectorBuilders(baseBuilder, masks)
	return DualVectorBuilder{
		builders: [...]filteredVectorBuilder{builders[0], builders[1]},
		//lazyBuilders: lazyBuilders,
	}
}

func (p *DualVectorBuilder) Get(index int) (int, VectorBuilder, VectorCursor) {
	b := &p.builders[index]
	if b.Builder != nil {
		return -1, b.Builder, b.Cursor
	}
	return int(b.ReuseBuilder), b.Builder, b.Cursor
}

func (p *DualVectorBuilder) ApplyEntry(entry VectorEntryData, isExcluded bool) {
	//
	//filter := fn(ve.node.GetIndex(), ve.nodeData)
	//
	//if se.id == ve.rank.id {
	//	// member
	//	p.scanVectorEntry(filter, ve.role, ve.node, ve.rank.GetPower(), ve.nodeData.NodeAnnouncedState)
	//	continue
	//}
	//
	//// joiner
	//if !ve.HasJoiner() {
	panic("illegal state")
	//}
	//n := ve.joiner
	//p.scanVectorEntry(filter, entry)
}

func (p *DualVectorBuilder) scanVectorEntry(filter uint16, entry VectorEntryData) {
	for i := range p.builders {
		b := &p.builders[i]
		b.LastFilter = b.FilterMask&filter == b.FilterRes
	}
	if p.builders[1].Builder == nil && p.builders[1].LastFilter != p.builders[0].LastFilter {
		p.builders[1].Builder = p.builders[0].Builder.Fork()
		p.builders[1].Cursor = p.builders[0].Cursor
	}
	for i := range p.builders {
		b := &p.builders[i]
		if !b.LastFilter || b.Builder == nil {
			continue
		}
		b.Cursor.BeforeNext(entry.Role)
		b.Builder.AddEntry(b.Cursor, entry)
		b.Cursor.AfterNext(entry.RequestedPower)
	}
}
