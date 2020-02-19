// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
		p.digester = p.digestFactory.CreateAnnouncementDigester()
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
	// if s.IsEmpty() {
	//	panic("illegal value")
	// }

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
		p.digester = p.digestFactory.CreateGlobulaStateDigester()
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
