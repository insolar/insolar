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
	"math"
	"sync"

	"github.com/insolar/insolar/network/consensus/gcpv2/errors"

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"

	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/common"
)

func NewNodeAppearanceAsSelf(np common2.LocalNodeProfile, callback *nodeCallback) *NodeAppearance {
	if np == nil {
		panic("node profile is nil")
	}
	np.LocalNodeProfile() // to avoid linter's paranoia

	return &NodeAppearance{
		profile:  np,
		state:    packets.NodeStateLocalActive,
		trust:    packets.SelfTrust,
		callback: callback,
	}
}

func (c *NodeAppearance) init(np common2.NodeProfile, callback *nodeCallback) {
	if np == nil {
		panic("node profile is nil")
	}
	c.profile = np
	c.callback = callback
}

type NodeAppearance struct {
	mutex sync.Mutex

	/* Provided externally at construction. Don't need mutex */
	profile                common2.NodeProfile // set by construction
	callback               *nodeCallback
	handlers               []PhasePerNodePacketHandler
	neighborTrustThreshold uint8

	/* Other fields - need mutex */

	//membership common2.MembershipProfile // one-time set
	claimSignature common2.NodeClaimSignature    // one-time set
	stateEvidence  common2.NodeStateHashEvidence // one-time set

	firstFraudDetails *errors.FraudError

	neighbourWeight uint32

	state           packets.NodeState
	trust           packets.NodeTrustLevel
	neighborReports uint8
	claimHash       common2.NodeClaimSignature
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
	target.profile.(common2.LocalNodeProfile).LocalNodeProfile()

	target.stateEvidence = c.stateEvidence
	target.claimSignature = c.claimSignature

	target.state = c.state
	target.trust = c.trust
}

// func (c *NodeAppearance) Frauds() errors.FraudFactory {
// 	return c.errorFactory.GetFraudFactory()
// }
//
// func (c *NodeAppearance) Blames() errors.BlameFactory {
// 	return c.errorFactory.GetBlameFactory()
// }

func (c *NodeAppearance) IsJoiner() bool {
	return c.profile.IsJoiner()
}

func (c *NodeAppearance) GetIndex() int {
	return c.profile.GetIndex()
}

func (c *NodeAppearance) GetShortNodeID() common.ShortNodeID {
	return c.profile.GetShortNodeID()
}

func (c *NodeAppearance) GetTrustLevel() packets.NodeTrustLevel {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.trust
}

func (c *NodeAppearance) GetProfile() common2.NodeProfile {
	return c.profile
}

func (c *NodeAppearance) VerifyPacketAuthenticity(packet packets.PacketParser, from common.HostIdentityHolder, preVerified bool) error {
	if preVerified {
		return nil
	}
	return VerifyPacketAuthenticityBy(packet, c.profile, c.profile.GetSignatureVerifier(), from)
}

func (c *NodeAppearance) SetReceivedPhase(phase packets.PhaseNumber) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.state.UpdReceivedPhase(phase)
}

func (c *NodeAppearance) SetReceivedByPacketType(pt packets.PacketType) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.state.UpdReceivedPacket(pt)
}

/* Explicit use of SetSentPhase is NOT recommended. Please use SetSentByPacketType */
func (c *NodeAppearance) SetSentPhase(phase packets.PhaseNumber) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.state.UpdSentPhase(phase)
}

func (c *NodeAppearance) SetSentByPacketType(pt packets.PacketType) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.state.UpdSentPacket(pt)
}

func (c *NodeAppearance) SetReceivedWithDupCheck(pt packets.PacketType) error {
	if c.SetReceivedByPacketType(pt) {
		return nil
	}
	return errors.ErrRepeatedPhasePacket
}

func (c *NodeAppearance) GetSignatureVerifier(vFactory common.SignatureVerifierFactory) common.SignatureVerifier {
	v := c.profile.GetSignatureVerifier()
	if v != nil {
		return v
	}
	return c.CreateSignatureVerifier(vFactory)
}

func (c *NodeAppearance) CreateSignatureVerifier(vFactory common.SignatureVerifierFactory) common.SignatureVerifier {
	return vFactory.GetSignatureVerifierWithPKS(c.profile.GetNodePublicKeyStore())
}

func (c *NodeAppearance) Locked(fn func() error) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return fn()
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) ApplyNodeMembership(mp common2.MembershipProfile) (bool, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c._applyNodeMembership(mp)
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) ApplyNeighbourEvidence(witness *NodeAppearance, mp common2.MembershipProfile) (modifiedNsh bool, trustBefore, trustAfter packets.NodeTrustLevel) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error
	trustBefore = c.trust
	modifiedNsh, err = c._applyNodeMembership(mp)

	if err == nil && witness.GetShortNodeID() != c.GetShortNodeID() { // a node can't be a witness to itself
		trustBefore = c.trust
		switch {
		case c.neighborReports == 0:
			c.trust.UpdateKeepNegative(packets.TrustBySome)
		case c.neighborReports == uint8(math.MaxUint8):
			panic("overflow")
		case c.neighborReports > c.neighborTrustThreshold:
			break // to allow the next statement to fire only once
		case c.neighborReports+1 > c.neighborTrustThreshold:
			c.trust.UpdateKeepNegative(packets.TrustByNeighbors)
		}
		c.neighborReports++

		if trustBefore != c.trust {
			c.callback.onTrustUpdated(c, trustBefore, c.trust)
		}
	}

	return modifiedNsh, trustBefore, c.trust
}

func (c *NodeAppearance) Frauds() errors.FraudFactory {
	return c.callback.GetFraudFactory()
}

func (c *NodeAppearance) Blames() errors.BlameFactory {
	return c.callback.GetBlameFactory()
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) _applyNodeMembership(mp common2.MembershipProfile) (bool, error) {

	if c.stateEvidence == nil {
		if mp.IsEmpty() {
			panic(fmt.Sprintf("membership evidence is nil: for=%v", c.GetShortNodeID()))
		}
		if c.GetIndex() != int(mp.Index) || c.GetPower() != mp.Power {
			return false, c.registerFraud(c.Frauds().NewMismatchedMembershipRank(c.GetProfile(), mp))
		}

		c.neighbourWeight ^= common.FoldUint64(mp.StateEvidence.GetNodeStateHash().FoldToUint64())
		c.stateEvidence = mp.StateEvidence
		c.claimSignature = mp.ClaimSignature

		c.callback.onNodeStateAssigned(c)

		return true, nil
	}

	lmp := c.getMembership()
	if mp.Equals(lmp) {
		return false, nil
	}

	return false, c.registerFraud(c.Frauds().NewMultipleMembershipProfiles(c.GetProfile(), lmp, mp))
}

//func (c *NodeAppearance) GetNodeStateHashEvidence() common2.NodeStateHashEvidence {
//	c.mutex.Lock()
//	defer c.mutex.Unlock()
//
//	if c.stateEvidence == nil {
//		panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeID()))
//	}
//	return c.membership
//}

func (c *NodeAppearance) GetNodeMembershipProfile() common2.MembershipProfile {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.stateEvidence == nil {
		panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeID()))
	}
	return c.getMembership()
}

func (c *NodeAppearance) GetNodeMembershipProfileOrEmpty() common2.MembershipProfile {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getMembership()
}

func (c *NodeAppearance) SetLocalNodeStateHashEvidence(evidence common2.NodeStateHashEvidence, claims common2.NodeClaimSignature) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.stateEvidence != nil {
		panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeID()))
	}
	if claims == nil {
		panic("illegal param")
	}
	c.neighbourWeight ^= common.FoldUint64(evidence.GetNodeStateHash().FoldToUint64())
	c.stateEvidence = evidence
	c.claimSignature = claims
}

func (c *NodeAppearance) GetNodeMembershipAndTrust() (common2.MembershipProfile, packets.NodeTrustLevel) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getMembership(), c.trust
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

func (c *NodeAppearance) GetPower() common2.MemberPower {
	return c.profile.GetPower()
}

func (c *NodeAppearance) registerFraud(fraud errors.FraudError) error {
	if fraud.IsUnknown() {
		panic("empty fraud")
	}

	prevTrust := c.trust
	if c.trust.Update(packets.FraudByThisNode) {
		c.firstFraudDetails = &fraud
		c.callback.onTrustUpdated(c, prevTrust, c.trust)
	}
	return fraud
}

func (c *NodeAppearance) RegisterFraud(fraud errors.FraudError) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	/* Here the pointer comparison is intentional to ensure exact NodeProfile, as it may change across rounds etc */
	if fraud.ViolatorNode() != c.GetProfile() {
		panic("misplaced fraud")
	}

	return c.registerFraud(fraud)
}

/*
deprecated
*/
func (c *NodeAppearance) RegisterFraudWithTrust(fraud errors.FraudError) (before, after packets.NodeTrustLevel, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	/* Here the pointer comparison is intentional to ensure exact NodeProfile, as it may change across rounds etc */
	if fraud.ViolatorNode() != c.GetProfile() {
		panic("misplaced fraud")
	}

	before = c.trust
	err = c.registerFraud(fraud)
	after = c.trust

	return
}

func (c *NodeAppearance) getMembership() common2.MembershipProfile {
	return common2.NewMembershipProfileByNode(c.profile, c.stateEvidence, c.claimSignature)
}
