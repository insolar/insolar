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

func NewNodeAppearanceAsSelf(np common2.LocalNodeProfile) *NodeAppearance {
	if np == nil {
		panic("node profile is nil")
	}
	np.LocalNodeProfile() //to avoid linter's paranoia

	return &NodeAppearance{
		profile: np,
		state:   packets.NodeStateLocalActive,
		trust:   packets.SelfTrust,
	}
}

func (c *NodeAppearance) init(np common2.NodeProfile) {
	if np == nil {
		panic("node profile is nil")
	}
	c.profile = np
}

type NodeAppearance struct {
	mutex sync.Mutex

	/* Provided externally at construction. Don't need mutex */
	profile common2.NodeProfile // set by construction
	// errorFactory errors.MisbehaviorFactories
	handlers               []PhasePerNodePacketHandler
	neighborTrustThreshold uint8

	/* Other fields - need mutex */
	Power             common2.MemberPower
	nshEvidence       common2.NodeStateHashEvidence // one-time set
	firstFraudDetails *errors.FraudError

	neighbourWeight uint32

	state           packets.NodeState
	trust           packets.NodeTrustLevel
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

	target.nshEvidence = c.nshEvidence
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

func (c *NodeAppearance) GetShortNodeId() common.ShortNodeID {
	return c.profile.GetShortNodeId()
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
func (c *NodeAppearance) ApplyNodeMembership(mp common2.MembershipProfile, evidence common2.NodeStateHashEvidence,
	errorFactory errors.MisbehaviorFactories) (bool, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c._applyNodeMembership(mp, evidence, true, errorFactory)
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) ApplyNeighbourEvidence(witness *NodeAppearance, mp common2.MembershipProfile, evidence common2.NodeStateHashEvidence,
	errorFactory errors.MisbehaviorFactories) (modifiedNsh bool, trustBefore, trustAfter packets.NodeTrustLevel) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	trustBefore = c.trust
	modifiedNsh, _ = c._applyNodeMembership(mp, evidence, false, errorFactory)

	if witness.GetShortNodeId() != c.GetShortNodeId() { // a node can't be a witness to itself
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
	}

	return modifiedNsh, trustBefore, c.trust
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) _applyNodeMembership(mp common2.MembershipProfile, evidence common2.NodeStateHashEvidence,
	direct bool, errorFactory errors.MisbehaviorFactories) (bool, error) {

	// TODO rank check
	// if c.GetIndex() != int(mp.Index) || c.GetPower() != mp.Power {
	// 	return false, c.registerFraud(errorFactory.GetFraudFactory().NewMismatchedRank(c.GetProfile(), evidence))
	// }

	if c.nshEvidence == nil {
		if evidence == nil {
			panic(fmt.Sprintf("evidence is nil: for=%v", c.GetShortNodeId()))
		}
		c.neighbourWeight ^= common.FoldUint64(mp.Nsh.FoldToUint64())
		c.nshEvidence = evidence
		return true, nil
	} else if c.nshEvidence.GetNodeStateHash().Equals(mp.Nsh) && mp.Nsh.Equals(evidence.GetNodeStateHash()) &&
		c.nshEvidence.GetGlobulaNodeStateSignature().Equals(evidence.GetGlobulaNodeStateSignature()) {

		return false, nil
	}

	return false, c.registerFraud(errorFactory.GetFraudFactory().NewMultipleNsh(c.GetProfile(), c.nshEvidence, evidence))
}

func (c *NodeAppearance) GetNodeStateHashEvidence() common2.NodeStateHashEvidence {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.nshEvidence == nil {
		panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeId()))
	}
	return c.nshEvidence
}

func (c *NodeAppearance) GetNodeMembershipAndEvidence() (common2.MembershipProfile, common2.NodeStateHashEvidence) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.nshEvidence == nil {
		panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeId()))
	}
	return c.getMembership(), c.nshEvidence
}

func (c *NodeAppearance) SetLocalNodeStateHashEvidence(evidence common2.NodeStateHashEvidence) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.nshEvidence != nil {
		panic(fmt.Sprintf("illegal state: for=%v", c.GetShortNodeId()))
	}
	c.neighbourWeight ^= common.FoldUint64(evidence.GetNodeStateHash().FoldToUint64())
	c.nshEvidence = evidence
}

func (c *NodeAppearance) GetNodeMembershipAndTrust() (common2.MembershipProfile, packets.NodeTrustLevel) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getMembership(), c.trust
}

func (c *NodeAppearance) GetNodeMembership() common2.MembershipProfile {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getMembership()
}

func (c *NodeAppearance) getMembership() common2.MembershipProfile {
	var nsh common2.NodeStateHash = nil
	if c.nshEvidence != nil {
		nsh = c.nshEvidence.GetNodeStateHash()
	}
	return common2.NewMembershipProfile(uint16(c.GetIndex()), c.GetPower(), nsh)
}

func (c *NodeAppearance) IsNshRequired() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.nshEvidence == nil
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

	if c.trust.Update(packets.FraudByThisNode) {
		c.firstFraudDetails = &fraud
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
