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
