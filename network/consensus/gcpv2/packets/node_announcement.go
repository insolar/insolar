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

package packets

import (
	"fmt"

	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

func NewNodeAnnouncement(np common2.NodeProfile, membership common2.MembershipProfile, nodeCount int, pn common.PulseNumber) *NodeAnnouncementProfile {
	return &NodeAnnouncementProfile{
		nodeID:     np.GetShortNodeID(),
		nodeCount:  uint16(nodeCount),
		membership: membership,
		pn:         pn,
	}
}

func NewNodeAnnouncementOf(na MembershipAnnouncementReader, pn common.PulseNumber) *NodeAnnouncementProfile {
	nr := na.GetNodeRank()
	return &NodeAnnouncementProfile{
		nodeID:    na.GetNodeID(),
		nodeCount: nr.GetTotalCount(),
		pn:        pn,
		membership: common2.NewMembershipProfile(
			nr.GetIndex(),
			nr.GetPower(),
			na.GetNodeStateHashEvidence(),
			na.GetAnnouncementSignature(),
			na.GetRequestedPower(),
		),
	}
}

var _ MembershipAnnouncementReader = &NodeAnnouncementProfile{}

type NodeAnnouncementProfile struct {
	nodeID     common.ShortNodeID
	membership common2.MembershipProfile
	pn         common.PulseNumber
	nodeCount  uint16
}

func (c *NodeAnnouncementProfile) GetRequestedPower() common2.MemberPower {
	return c.membership.RequestedPower
}

func (c *NodeAnnouncementProfile) IsLeaving() bool {
	return false
}

func (c *NodeAnnouncementProfile) GetLeaveReason() uint32 {
	return 0
}

func (c *NodeAnnouncementProfile) GetJoinerID() common.ShortNodeID {
	return common.AbsentShortNodeID
}

func (c *NodeAnnouncementProfile) GetJoinerAnnouncement() JoinerAnnouncementReader {
	return nil
}

func (c *NodeAnnouncementProfile) GetNodeRank() common2.MembershipRank {
	return common2.NewMembershipRank(c.membership.Power, c.membership.Index, c.nodeCount, 0)
}

func (c *NodeAnnouncementProfile) GetAnnouncementSignature() common2.MemberAnnouncementSignature {
	return c.membership.AnnounceSignature
}

func (c *NodeAnnouncementProfile) GetNodeID() common.ShortNodeID {
	return c.nodeID
}

func (c *NodeAnnouncementProfile) GetNodeCount() uint16 {
	return c.nodeCount
}

func (c *NodeAnnouncementProfile) GetNodeStateHashEvidence() common2.NodeStateHashEvidence {
	return c.membership.StateEvidence
}

func (c NodeAnnouncementProfile) String() string {
	return fmt.Sprintf("{id:%d %03d/%d %s}", c.nodeID, c.membership.Index, c.nodeCount, c.membership.StringParts())
}

func (c *NodeAnnouncementProfile) GetMembershipProfile() common2.MembershipProfile {
	return c.membership
}

func (c *NodeAnnouncementProfile) GetPulseNumber() common.PulseNumber {
	return c.pn
}
