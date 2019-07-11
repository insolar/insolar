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
	"fmt"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"math"
	"sync"

	"github.com/insolar/insolar/network/consensus/gcpv2/errors"

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"

	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/common"
)

func NewNodeAppearanceAsSelf(np api.LocalNodeProfile, callback *nodeContext) *NodeAppearance {
	np.LocalNodeProfile() // to avoid linter's paranoia

	r := &NodeAppearance{
		state: api.NodeStateLocalActive,
		trust: common2.SelfTrust,
	}
	r.init(np, callback, 0)
	return r
}

func (c *NodeAppearance) init(np api.NodeProfile, callback NodeContextHolder, baselineWeight uint32) {
	if np == nil {
		panic("node profile is nil")
	}
	c.profile = np
	c.callback = callback
	c.neighbourWeight = baselineWeight
}

type NodeAppearance struct {
	mutex sync.Mutex

	/* Provided externally at construction. Don't need mutex */
	profile  api.NodeProfile // set by construction
	callback *nodeContext
	handlers []PhasePerNodePacketFunc

	/* Other fields - need mutex */

	//membership common2.MembershipProfile // one-time set
	announceSignature api.MemberAnnouncementSignature // one-time set
	stateEvidence     api.NodeStateHashEvidence       // one-time set
	requestedPower    api.MemberPower                 // one-time set

	requestedJoiner      *NodeAppearance // one-time set
	requestedLeave       bool            // one-time set
	requestedLeaveReason uint32          // one-time set

	firstFraudDetails *errors.FraudError

	neighbourWeight uint32

	state           api.NodeState
	trust           common2.NodeTrustLevel
	neighborReports uint8
}

func (c *NodeAppearance) String() string {
	return fmt.Sprintf("node:{%v}", c.profile)
}

// Unsafe
func LessByNeighbourWeightForNodeAppearance(n1, n2 interface{}) bool {
	return n1.(*NodeAppearance).neighbourWeight < n2.(*NodeAppearance).neighbourWeight
}

// LOCK - self, target must be safe
func (c *NodeAppearance) copySelfTo(target *NodeAppearance) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	/* Ensure that the target is LocalNode */
	target.profile.(api.LocalNodeProfile).LocalNodeProfile()

	target.stateEvidence = c.stateEvidence
	target.announceSignature = c.announceSignature
	target.requestedPower = c.requestedPower
	target.requestedJoiner = c.requestedJoiner
	target.requestedLeave = c.requestedLeave
	target.requestedLeaveReason = c.requestedLeaveReason
	target.firstFraudDetails = c.firstFraudDetails

	target.state = c.state
	target.trust = c.trust
	target.callback.updatePopulationVersion()
}

func (c *NodeAppearance) IsJoiner() bool {
	return c.profile.IsJoiner()
}

func (c *NodeAppearance) GetIndex() int {
	return c.profile.GetIndex()
}

func (c *NodeAppearance) GetShortNodeID() common.ShortNodeID {
	return c.profile.GetShortNodeID()
}

func (c *NodeAppearance) GetTrustLevel() common2.NodeTrustLevel {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.trust
}

func (c *NodeAppearance) GetProfile() api.NodeProfile {
	return c.profile
}

func (c *NodeAppearance) VerifyPacketAuthenticity(packet packets.PacketParser, from endpoints.HostIdentityHolder, strictFrom bool) error {
	return VerifyPacketAuthenticityBy(packet, c.profile, c.profile.GetSignatureVerifier(), from, strictFrom)
}

func (c *NodeAppearance) SetReceivedPhase(phase api.PhaseNumber) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.callback.updatePopulationVersion()
	return c.state.UpdReceivedPhase(phase)
}

func (c *NodeAppearance) SetReceivedByPacketType(pt api.PacketType) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.callback.updatePopulationVersion()
	return c.state.UpdReceivedPacket(pt)
}

/* Explicit use of SetSentPhase is NOT recommended. Please use SetSentByPacketType */
func (c *NodeAppearance) SetSentPhase(phase api.PhaseNumber) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.callback.updatePopulationVersion()
	return c.state.UpdSentPhase(phase)
}

func (c *NodeAppearance) SetSentByPacketType(pt api.PacketType) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.callback.updatePopulationVersion()
	return c.state.UpdSentPacket(pt)
}

func (c *NodeAppearance) SetReceivedWithDupCheck(pt api.PacketType) error {
	if c.SetReceivedByPacketType(pt) {
		return nil
	}
	return errors.ErrRepeatedPhasePacket
}

func (c *NodeAppearance) GetSignatureVerifier(vFactory cryptography_containers.SignatureVerifierFactory) cryptography_containers.SignatureVerifier {
	v := c.profile.GetSignatureVerifier()
	if v != nil {
		return v
	}
	return c.CreateSignatureVerifier(vFactory)
}

func (c *NodeAppearance) CreateSignatureVerifier(vFactory cryptography_containers.SignatureVerifierFactory) cryptography_containers.SignatureVerifier {
	return vFactory.GetSignatureVerifierWithPKS(c.profile.GetNodePublicKeyStore())
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) ApplyNodeMembership(mp api.MembershipAnnouncement) (bool, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c._applyNodeMembership(mp)
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) ApplyNeighbourEvidence(witness *NodeAppearance, mp api.MembershipAnnouncement) (bool, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	trustBefore := c.trust
	modified, err := c._applyNodeMembership(mp)

	if err == nil && witness.GetShortNodeID() != c.GetShortNodeID() { // a node can't be a witness to itself
		switch {
		case c.neighborReports == 0:
			c.trust.UpdateKeepNegative(common2.TrustBySome)
		case c.neighborReports == uint8(math.MaxUint8):
			panic("overflow")
		case c.neighborReports > c.GetNeighborTrustThreshold():
			break // to allow the next statement to fire only once
		case c.neighborReports+1 > c.GetNeighborTrustThreshold():
			c.trust.UpdateKeepNegative(common2.TrustByNeighbors)
		}

		c.neighborReports++
		c.callback.updatePopulationVersion()
	}
	if trustBefore != c.trust {
		c.callback.onTrustUpdated(c, trustBefore, c.trust)
	}

	return modified, err
}

func (c *NodeAppearance) Frauds() errors.FraudFactory {
	return c.callback.GetFraudFactory()
}

func (c *NodeAppearance) Blames() errors.BlameFactory {
	return c.callback.GetBlameFactory()
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) _applyNodeMembership(ma api.MembershipAnnouncement) (bool, error) {

	if ma.Membership.IsEmpty() {
		panic(fmt.Sprintf("membership evidence is nil: for=%v", c.GetShortNodeID()))
	}

	if c.stateEvidence != nil {
		lmp := c.getMembership()
		var lma api.MembershipAnnouncement
		if ma.Membership.Equals(lmp) && ma.IsLeaving == c.requestedLeave {
			switch {
			case c.requestedLeave:
				if ma.LeaveReason == c.requestedLeaveReason {
					return false, nil
				}
				lma = api.NewMembershipAnnouncementWithLeave(lmp, c.requestedLeaveReason)
			case c.requestedJoiner == nil:
				if ma.Joiner == nil {
					return false, nil
				}
				lma = api.NewMembershipAnnouncement(lmp)
			default:
				if api.EqualIntroProfiles(c.requestedJoiner.GetProfile(), ma.Joiner) {
					return false, nil
				}
				lma = api.NewMembershipAnnouncementWithJoiner(lmp, c.requestedJoiner.profile)
			}
		}
		return c.registerFraud(c.Frauds().NewInconsistentMembershipAnnouncement(c.GetProfile(), lma, ma))
	}

	c.callback.updatePopulationVersion()

	if ma.IsLeaving {
		c.requestedLeave = true
		c.requestedLeaveReason = ma.LeaveReason
	} else if ma.Joiner != nil {
		panic("not implemented") //TODO implement
	}

	c.neighbourWeight ^= long_bits.FoldUint64(ma.Membership.StateEvidence.GetNodeStateHash().FoldToUint64())
	c.stateEvidence = ma.Membership.StateEvidence
	c.announceSignature = ma.Membership.AnnounceSignature
	c.requestedPower = ma.Membership.RequestedPower

	c.callback.onNodeStateAssigned(c)

	return true, nil
}

func (c *NodeAppearance) GetNodeMembershipProfile() api.MembershipProfile {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.stateEvidence == nil {
		panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeID()))
	}
	return c.getMembership()
}

func (c *NodeAppearance) GetNodeTrustAndMembershipOrEmpty() (api.MembershipProfile, common2.NodeTrustLevel) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//if c.stateEvidence == nil {
	//	panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeID()))
	//}
	return c.getMembership(), c.trust
}

func (c *NodeAppearance) GetNodeMembershipProfileOrEmpty() api.MembershipProfile {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getMembership()
}

func (c *NodeAppearance) SetLocalNodeStateHashEvidence(evidence api.NodeStateHashEvidence, announce api.MemberAnnouncementSignature) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.stateEvidence != nil {
		panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeID()))
	}
	if announce == nil {
		panic("illegal param")
	}

	c.neighbourWeight ^= long_bits.FoldUint64(evidence.GetNodeStateHash().FoldToUint64())
	c.stateEvidence = evidence
	c.announceSignature = announce
	c.callback.updatePopulationVersion()
}

func (c *NodeAppearance) IsNshRequired() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.stateEvidence == nil
}

func (c *NodeAppearance) HasReceivedAnyPhase() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.state.HasReceived()
}

func (c *NodeAppearance) GetNeighbourWeight() uint32 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.neighbourWeight
}

func (c *NodeAppearance) registerFraud(fraud errors.FraudError) (bool, error) {
	if fraud.IsUnknown() {
		panic("empty fraud")
	}

	prevTrust := c.trust
	if c.trust.Update(common2.FraudByThisNode) {
		c.firstFraudDetails = &fraud
		c.callback.updatePopulationVersion()
		c.callback.onTrustUpdated(c, prevTrust, c.trust)
		return true, fraud
	}
	return false, fraud
}

func (c *NodeAppearance) RegisterFraud(fraud errors.FraudError) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	/* Here the pointer comparison is intentional to ensure exact NodeProfile, as it may change across rounds etc */
	if fraud.ViolatorNode() != c.GetProfile() {
		panic("misplaced fraud")
	}

	return c.registerFraud(fraud)
}

// MUST BE NO LOCK
func (c *NodeAppearance) getMembership() api.MembershipProfile {
	return api.NewMembershipProfileByNode(c.profile, c.stateEvidence, c.announceSignature, c.requestedPower)
}

func (c *NodeAppearance) GetNeighborTrustThreshold() uint8 {
	return c.callback.GetNeighbourhoodTrustThreshold()
}

func (c *NodeAppearance) NotifyOnCustom(event interface{}) {
	c.callback.onCustomEvent(c, event)
}

func (c *NodeAppearance) getPacketHandler(i int) PhasePerNodePacketFunc {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.handlers) == 0 {
		return nil
	}
	return c.handlers[i]
}

func (c *NodeAppearance) ResetAllPacketHandlers() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers = nil
}

func (c *NodeAppearance) ResetPacketHandlers(indices ...int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.handlers) == 0 {
		return
	}

	for i := range indices {
		c.handlers[i] = nil
	}
	for _, h := range c.handlers {
		if h != nil {
			return
		}
	}
	c.handlers = nil
}

func (c *NodeAppearance) GetRequestedState() (bool, uint32, *NodeAppearance, api.MembershipProfile, common2.NodeTrustLevel) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.requestedLeave, c.requestedLeaveReason, c.requestedJoiner, c.getMembership(), c.trust
}

func (c *NodeAppearance) GetRequestedAnnouncement() api.MembershipAnnouncement {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	mb := c.getMembership()
	switch {
	case c.requestedLeave:
		return api.NewMembershipAnnouncementWithLeave(mb, c.requestedLeaveReason)
	case c.requestedJoiner != nil:
		return api.NewMembershipAnnouncementWithJoiner(mb, c.requestedJoiner.profile)
	default:
		return api.NewMembershipAnnouncement(mb)
	}
}
