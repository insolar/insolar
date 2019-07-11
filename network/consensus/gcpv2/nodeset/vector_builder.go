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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
)

/*
Contains copy of NodeAppearance fields that can be changed, to avoid possible racing
*/

type VectorEntryData struct {
	NodeID insolar.ShortNodeID

	Role           gcp_types.NodePrimaryRole
	RequestedPower gcp_types.MemberPower
	TrustLevel     gcp_types.NodeTrustLevel
	Mode           gcp_types.MemberOpMode

	// Node *NodeAppearance
	// common.MembershipAnnouncement
	gcp_types.NodeAnnouncedState
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

	LastRole gcp_types.NodePrimaryRole
}

func (p *VectorCursor) BeforeNext(role gcp_types.NodePrimaryRole) {
	if p.LastRole == role {
		return
	}
	p.RoleIndex = 0
	p.PowIndex = 0
	p.LastRole = role
}

func (p *VectorCursor) AfterNext(power gcp_types.MemberPower) {
	p.RoleIndex++
	p.PowIndex += power.ToLinearValue()
	p.NodeIndex++
}
