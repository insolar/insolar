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
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

func NewAnnouncementSequenceCalc(digestFactory transport.ConsensusDigestFactory) AnnouncementSequenceCalc {
	return AnnouncementSequenceCalc{digestFactory: digestFactory}
}

type AnnouncementSequenceCalc struct {
	digestFactory transport.ConsensusDigestFactory
	digester      cryptkit.SequenceDigester
}

func (p *AnnouncementSequenceCalc) AddNext(nodeData VectorEntryData, zeroPower bool) {
	if p.digester == nil {
		p.digester = p.digestFactory.GetAnnouncementDigester()
	}
	p.digester.AddNext(nodeData.AnnounceSignature)
}

func (p *AnnouncementSequenceCalc) ForkSequence() VectorEntryDigester {
	if p.digester == nil {
		return &AnnouncementSequenceCalc{digestFactory: p.digestFactory}
	}
	r := AnnouncementSequenceCalc{}
	r.ForkSequenceOf(*p)
	return &r
}

func (p *AnnouncementSequenceCalc) IsEmpty() bool {
	return p.digester == nil && p.digestFactory == nil
}

func (p *AnnouncementSequenceCalc) ForkSequenceOf(s AnnouncementSequenceCalc) {
	if !p.IsEmpty() {
		panic("illegal state")
	}
	//if s.IsEmpty() {
	//	panic("illegal value")
	//}

	if s.digester != nil {
		p.digester = s.digester.ForkSequence()
	} else {
		p.digestFactory = s.digestFactory
	}
}

func (p *AnnouncementSequenceCalc) FinishSequence() cryptkit.DigestHolder {
	if p.digester == nil {
		return nil
	}
	return p.digester.FinishSequence().AsDigestHolder()
}

func NewStateAndRankSequenceCalc(digestFactory transport.ConsensusDigestFactory, nodeID insolar.ShortNodeID, roleCountHint int) StateAndRankSequenceCalc {
	return StateAndRankSequenceCalc{digestFactory: digestFactory, nodeID: nodeID,
		entries: make([]memberEntry, 0, roleCountHint), cursor: member.RankCursor{Role: ^member.PrimaryRole(0)}}
}

type StateAndRankSequenceCalc struct {
	digestFactory transport.ConsensusDigestFactory
	digester      transport.StateDigester

	nodeID       insolar.ShortNodeID
	nodeFullRank member.FullRank

	cursor              member.RankCursor
	roleFirstTotalIndex member.Index
	entries             []memberEntry
}

type memberEntry struct {
	state   proofs.NodeStateHashEvidence
	capture bool

	SpecialRoles member.SpecialRole
	Power        member.Power
	OpMode       member.OpMode

	RolePowerIndex uint32
}

func (p *StateAndRankSequenceCalc) AddNext(nodeData VectorEntryData, zeroPower bool) {
	np := nodeData.Profile.GetStatic()
	orderingRole := np.GetPrimaryRole()
	if orderingRole == member.PrimaryRoleInactive {
		panic("illegal state")
	}
	if zeroPower {
		orderingRole = member.PrimaryRoleInactive
	}
	if p.cursor.Role != orderingRole {
		if p.cursor.Role < orderingRole {
			panic("illegal state")
		}
		p.flushRoleMembers()
		p.roleFirstTotalIndex = p.cursor.TotalIndex
		p.cursor = member.RankCursor{Role: orderingRole, TotalIndex: p.cursor.TotalIndex}
	}

	nodeID := np.GetStaticNodeID()
	me := memberEntry{
		state:        nodeData.StateEvidence,
		capture:      p.nodeID == nodeID,
		SpecialRoles: np.GetSpecialRoles(),
	}

	if orderingRole == member.PrimaryRoleInactive {
		me.OpMode = nodeData.RequestedMode
		p.hashMemberEntry(me, 0)
	} else {
		me.OpMode = nodeData.RequestedMode
		me.Power = nodeData.RequestedPower
		me.RolePowerIndex = p.cursor.RolePowerIndex

		p.entries = append(p.entries, me)

		p.cursor.RolePowerIndex += uint32(nodeData.RequestedPower.ToLinearValue())
		p.cursor.RoleIndex++
	}
	p.cursor.TotalIndex++
}

func (p *StateAndRankSequenceCalc) hashMemberEntry(v memberEntry, roleIndex member.Index) {
	if p.digester == nil {
		p.digester = p.digestFactory.GetGlobulaStateDigester()
	}

	fr := member.FullRank{
		InterimRank: member.InterimRank{
			RankCursor: member.RankCursor{
				Role:           p.cursor.Role,
				RoleIndex:      roleIndex,
				RolePowerIndex: v.RolePowerIndex,
				TotalIndex:     p.roleFirstTotalIndex + roleIndex,
			},
			SpecialRoles: v.SpecialRoles,
			Power:        v.Power,
			OpMode:       v.OpMode,
		},
		RoleCount: p.cursor.RoleIndex.AsUint16(),
		RolePower: p.cursor.RolePowerIndex,
	}
	if v.capture {
		p.nodeFullRank = fr
	}

	if v.state == nil {
		p.digester.AddNext(nil, fr)
	} else {
		p.digester.AddNext(v.state.GetDigestHolder(), fr)
	}
}

func (p *StateAndRankSequenceCalc) flushRoleMembers() {
	if len(p.entries) == 0 {
		return
	}

	for roleIndex, v := range p.entries {
		p.hashMemberEntry(v, member.AsIndex(roleIndex))
	}
	p.entries = p.entries[:0]
}

func (p *StateAndRankSequenceCalc) FinishSequence() (cryptkit.DigestHolder, member.FullRank, member.Index) {
	p.flushRoleMembers()
	if p.digester == nil {
		return nil, p.nodeFullRank, 0
	}
	return p.digester.FinishSequence().AsDigestHolder(), p.nodeFullRank, p.cursor.TotalIndex
}

func (p *StateAndRankSequenceCalc) IsEmpty() bool {
	return p.digestFactory == nil && p.digester == nil
}

func (p *StateAndRankSequenceCalc) ForkSequenceOf(s StateAndRankSequenceCalc) {
	if !p.IsEmpty() {
		panic("illegal state")
	}
	if s.IsEmpty() {
		panic("illegal value")
	}

	*p = s
	if p.digester != nil {
		p.digester = p.digester.ForkSequence()
	}
	p.entries = append(make([]memberEntry, 0, cap(s.entries)), s.entries...)
}

func (p *StateAndRankSequenceCalc) ForkSequence() VectorEntryDigester {
	s := StateAndRankSequenceCalc{}
	s.ForkSequenceOf(*p)
	return &s
}
