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

package transport

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

func NewNodeAnnouncement(np profiles.ActiveNode, ma profiles.MembershipAnnouncement, nodeCount int,
	pn pulse.Number, joiner *JoinerAnnouncement) *NodeAnnouncementProfile {
	return &NodeAnnouncementProfile{
		static:    np.GetStatic(),
		nodeID:    np.GetNodeID(),
		nodeCount: uint16(nodeCount),
		ma:        ma,
		pn:        pn,
		joiner:    joiner,
	}
}

var _ MembershipAnnouncementReader = &NodeAnnouncementProfile{}

type NodeAnnouncementProfile struct {
	static    profiles.StaticProfile
	ma        profiles.MembershipAnnouncement
	nodeID    insolar.ShortNodeID
	pn        pulse.Number
	nodeCount uint16
	joiner    *JoinerAnnouncement
}

func (c *NodeAnnouncementProfile) GetRequestedPower() member.Power {
	return c.ma.Membership.RequestedPower
}

func (c *NodeAnnouncementProfile) IsLeaving() bool {
	return c.ma.IsLeaving
}

func (c *NodeAnnouncementProfile) GetLeaveReason() uint32 {
	return c.ma.LeaveReason
}

func (c *NodeAnnouncementProfile) GetJoinerID() insolar.ShortNodeID {
	return c.ma.JoinerID
}

func (c *NodeAnnouncementProfile) GetJoinerAnnouncement() JoinerAnnouncementReader {
	if c.joiner == nil {
		return nil
	}

	if !c.ma.JoinerID.IsAbsent() && c.joiner.GetBriefIntroduction().GetStaticNodeID() != c.ma.JoinerID {
		panic("illegal state")
	}
	return c.joiner
}

func (c *NodeAnnouncementProfile) GetNodeRank() member.Rank {
	return c.ma.Membership.AsRankUint16(c.nodeCount)
}

func (c *NodeAnnouncementProfile) GetAnnouncementSignature() proofs.MemberAnnouncementSignature {
	return c.ma.Membership.AnnounceSignature
}

func (c *NodeAnnouncementProfile) GetNodeID() insolar.ShortNodeID {
	return c.nodeID
}

func (c *NodeAnnouncementProfile) GetNodeCount() uint16 {
	return c.nodeCount
}

func (c *NodeAnnouncementProfile) GetNodeStateHashEvidence() proofs.NodeStateHashEvidence {
	return c.ma.Membership.StateEvidence
}

func (c NodeAnnouncementProfile) String() string {
	announcement := ""
	if c.IsLeaving() {
		announcement = fmt.Sprintf(" leave:%d", c.GetLeaveReason())
	} else if !c.GetJoinerID().IsAbsent() {
		joinerIntro := ""
		if c.joiner != nil {
			joinerIntro = "+intro"
		}
		announcement = fmt.Sprintf(" join:%d%s", c.GetJoinerID(), joinerIntro)
	}
	return fmt.Sprintf("{id:%d %03d/%d%s %s}", c.nodeID, c.ma.Membership.Index, c.nodeCount, announcement, c.ma.Membership.StringParts())
}

func (c *NodeAnnouncementProfile) GetMembershipProfile() profiles.MembershipProfile {
	return c.ma.Membership
}

func (c *NodeAnnouncementProfile) GetPulseNumber() pulse.Number {
	return c.pn
}

func (c *NodeAnnouncementProfile) GetStatic() profiles.StaticProfile {
	return c.static
}
