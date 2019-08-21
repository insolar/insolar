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

package population

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
)

func NewNodeAppearanceAsSelf(np profiles.LocalNode, powerRequest power.Request, hook *Hook) NodeAppearance {
	np.LocalNodeProfile() // to avoid linter's paranoia

	sp := np.GetStatic()
	pw := member.Power(0)
	switch {
	case np.IsJoiner():
		pw = sp.GetStartPower()
	case np.GetOpMode().IsPowerless():
		break
	default:
		powerRequest.Update(&pw, sp.GetExtension().GetPowerLevels())
	}

	return NodeAppearance{
		profile:         np,
		limiter:         phases.NewLocalPacketLimiter(),
		hook:            hook,
		neighbourWeight: 0,
		requestedPower:  pw,
	}
}

// LOCK - self, target must be safe
func (c *NodeAppearance) CopySelfTo(target *NodeAppearance) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	/* Ensure that the target is LocalNode */
	target.profile.(profiles.LocalNode).LocalNodeProfile()

	if c.isStateAvailable() || c.announceSignature != nil || target.isStateAvailable() || target.announceSignature != nil {
		panic("illegal state")
	}
	// target.stateEvidence = c.stateEvidence
	// target.announceSignature = c.announceSignature

	target.requestedPower = c.requestedPower
	target.requestedJoinerID = c.requestedJoinerID
	target.requestedLeave = c.requestedLeave
	target.requestedLeaveReason = c.requestedLeaveReason
	target.firstFraudDetails = c.firstFraudDetails
	target.trust = c.trust

	// target.limiter = c.limiter
	// target.hook.UpdatePopulationVersion()
}

func NewEmptyNodeAppearance(np profiles.ActiveNode) NodeAppearance {

	if np == nil {
		panic("illegal value")
	}
	return NodeAppearance{
		profile: np,
	}
}

func NewNodeAppearance(np profiles.ActiveNode, baselineWeight uint32,
	limiter phases.PacketLimiter, hook *Hook, handlers []DispatchMemberPacketFunc) NodeAppearance {

	if np == nil {
		panic("illegal value")
	}

	return NodeAppearance{
		profile:         np,
		limiter:         limiter,
		hook:            hook,
		neighbourWeight: baselineWeight,
		handlers:        handlers,
		// requestedPower:  np.GetDeclaredPower(),
	}
}

type NodeAppearance struct {
	mutex sync.Mutex

	/* Provided externally at construction. Don't need mutex */
	profile  profiles.ActiveNode // set by construction
	hook     *Hook
	handlers []DispatchMemberPacketFunc

	/* Other fields - need mutex */

	// membership common2.MembershipProfile // one-time set
	announceSignature proofs.MemberAnnouncementSignature // one-time set
	stateEvidence     proofs.NodeStateHashEvidence       // one-time set
	requestedPower    member.Power                       // one-time set

	// statelessDigest cryptkit.DigestHolder

	// joinerSecret         cryptkit.Digest     // TODO implement
	requestedJoinerID    insolar.ShortNodeID // one-time set
	requestedLeave       bool                // one-time set
	requestedLeaveReason uint32              // one-time set

	firstFraudDetails *misbehavior.FraudError

	neighbourWeight uint32

	limiter         phases.PacketLimiter
	trust           member.TrustLevel
	neighborReports uint8
}

func (c *NodeAppearance) EncryptJoinerSecret(joinerSecret cryptkit.DigestHolder) cryptkit.DigestHolder {
	// TODO encryption of joinerSecret
	return joinerSecret
}

func (c *NodeAppearance) GetStatic() profiles.StaticProfile {
	return c.profile.GetStatic()
}

func (c *NodeAppearance) CanIntroduceJoiner() bool {
	return c.profile.CanIntroduceJoiner()
}

func (c *NodeAppearance) GetReportProfile() profiles.BaseNode {
	return c.profile
}

func (c *NodeAppearance) GetRank(nodeCount int) member.Rank {
	return profiles.ProfileAsRank(c.profile, nodeCount)
}

func (c *NodeAppearance) DispatchMemberPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags coreapi.PacketVerifyFlags, pd PacketDispatcher) error {

	return pd.DispatchMemberPacket(ctx, packet.GetMemberPacket(), c)
}

func (c *NodeAppearance) String() string {
	return fmt.Sprintf("node:{%v}", c.profile)
}

func LessByNeighbourWeightForNodeAppearance(n1, n2 interface{}) bool {
	return n1.(*NodeAppearance).neighbourWeight < n2.(*NodeAppearance).neighbourWeight
}

func (c *NodeAppearance) IsJoiner() bool {
	return c.profile.IsJoiner()
}

func (c *NodeAppearance) IsLocal() bool {
	return c.profile.GetNodeID() == c.hook.GetLocalNodeID()
}

func (c *NodeAppearance) GetIndex() member.Index {
	return c.profile.GetIndex()
}

func (c *NodeAppearance) GetNodeID() insolar.ShortNodeID {
	return c.profile.GetNodeID()
}

func (c *NodeAppearance) GetTrustLevel() member.TrustLevel {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.trust
}

func (c *NodeAppearance) GetProfile() profiles.ActiveNode {
	return c.profile
}

func (c *NodeAppearance) VerifyPacketAuthenticity(ps cryptkit.SignedDigest, from endpoints.Inbound, strictFrom bool) error {
	return coreapi.VerifyPacketAuthenticityBy(ps, c.profile.GetStatic(), c.profile.GetSignatureVerifier(), from, strictFrom)
}

func (c *NodeAppearance) SetPacketReceived(pt phases.PacketType) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	upd := false
	upd, c.limiter = c.limiter.SetPacketReceived(pt)
	if upd {
		c.hook.UpdatePopulationVersion()
		return true
	}
	return false
}

func (c *NodeAppearance) CanReceivePacket(pt phases.PacketType) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.limiter.CanReceivePacket(pt)
}

func (c *NodeAppearance) SetPacketSent(pt phases.PacketType) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	upd := false
	upd, c.limiter = c.limiter.SetPacketSent(pt)
	if upd {
		c.hook.UpdatePopulationVersion()
		return true
	}
	return false
}

func (c *NodeAppearance) GetSignatureVerifier() cryptkit.SignatureVerifier {
	v := c.profile.GetSignatureVerifier()
	if v != nil {
		return v
	}
	vFactory := c.hook.GetCryptographyAssistant()
	return vFactory.CreateSignatureVerifierWithPKS(c.profile.GetStatic().GetPublicKeyStore())
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) ApplyNodeMembership(mp profiles.MemberAnnouncement, applyAfterChecks MembershipApplyFunc) (bool, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	modified, err := c._applyState(mp, applyAfterChecks)
	c.updateNodeTrustLevel(c.trust, member.TrustBySelf)
	return modified, err
}

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) ApplyNeighbourEvidence(witness *NodeAppearance, mp profiles.MemberAnnouncement,
	cappedTrust bool, applyAfterChecks MembershipApplyFunc) (bool, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	var updVersion uint32
	trustBefore := c.trust
	modified, err := c._applyState(mp, applyAfterChecks)

	if err == nil {
		switch {
		case witness.GetNodeID() != c.GetNodeID(): /* a node can't be a witness to itself */
			updVersion = c.incNeighborReports(cappedTrust)
		case mp.JoinerID == c.hook.GetLocalNodeID() && c.hook.GetLocalProfile().IsJoiner():
			// we trust to those who has introduced us
			// It is also REQUIRED as vector calculation requires at least one trusted node to work properly
			if c.trust.Update(member.TrustBySome) {
				modified = true
				updVersion = c.hook.UpdatePopulationVersion()
				break
			}
			fallthrough
		default:
			updVersion = c.hook.GetPopulationVersion()
		}
	}

	if trustBefore != c.trust {
		c.hook.OnTrustUpdated(updVersion, c, trustBefore, c.trust, c.profile.HasFullProfile())
	}

	return modified, err
}

func (c *NodeAppearance) incNeighborReports(cappedTrust bool) uint32 {
	switch {
	case c.neighborReports == 0:
		c.trust.UpdateKeepNegative(member.TrustBySome)
	case cappedTrust:
		// we can't increase trust higher than basic
		return c.hook.GetPopulationVersion()
	case c.neighborReports == uint8(math.MaxUint8):
		// panic("overflow")
		return c.hook.GetPopulationVersion()
	case c.neighborReports > c.GetNeighborTrustThreshold():
		break // to allow the next statement to fire only once
	case c.neighborReports+1 > c.GetNeighborTrustThreshold():
		c.trust.UpdateKeepNegative(member.TrustByNeighbors)
	}

	c.neighborReports++
	return c.hook.UpdatePopulationVersion()
}

func (c *NodeAppearance) UpdateNodeTrustLevel(trust member.TrustLevel) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.updateNodeTrustLevel(c.trust, trust)
}

func (c *NodeAppearance) updateNodeTrustLevel(trustBefore, trust member.TrustLevel) bool {

	updVersion := c.hook.GetPopulationVersion()
	modified := c.trust.Update(trust)
	if modified {
		updVersion = c.hook.UpdatePopulationVersion()
	}
	if trustBefore != c.trust {
		c.hook.OnTrustUpdated(updVersion, c, trustBefore, c.trust, c.profile.HasFullProfile())
	}
	return modified
}

func (c *NodeAppearance) Frauds() misbehavior.FraudFactory {
	return c.hook.GetFraudFactory()
}

func (c *NodeAppearance) Blames() misbehavior.BlameFactory {
	return c.hook.GetBlameFactory()
}

func (c *NodeAppearance) IsStateful() bool {
	return !c.hook.GetEphemeralMode().IsEnabled() && c.profile.IsStateful()
}

type MembershipApplyFunc func(ma profiles.MemberAnnouncement) error

/* Evidence MUST be verified before this call */
func (c *NodeAppearance) _applyState(ma profiles.MemberAnnouncement,
	applyAfterChecks MembershipApplyFunc) (bool, error) {

	if ma.Membership.IsEmpty() {
		panic(fmt.Sprintf("membership evidence is nil: %v", c.GetNodeID()))
	}

	if ma.MemberID != c.GetNodeID() {
		panic(fmt.Sprintf("member announcement is for a wrong node: %v %v", c.GetNodeID(), ma.MemberID))
	}

	if c.isStateAvailable() {
		lmp := c.getMembership()
		// var lma profiles.MembershipAnnouncement
		if ma.Membership.Equals(lmp) && ma.IsLeaving == c.requestedLeave {
			switch {
			case c.requestedLeave:
				if ma.LeaveReason == c.requestedLeaveReason {
					return false, nil
				}
			default:
				if c.requestedJoinerID == ma.JoinerID {
					return false, nil
				}
			}
		}
		lma := c.getMembershipAnnouncement()
		return c.registerFraud(c.Frauds().NewInconsistentMembershipAnnouncement(c.GetProfile(), lma, ma.MembershipAnnouncement))
	}

	updVersion := c.hook.UpdatePopulationVersion()

	switch {
	case ma.IsLeaving:
		switch {
		case c.IsJoiner():
			c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "joiner can't request leave"))
		case !ma.JoinerID.IsAbsent():
			c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "leaver can't introduce a joiner"))
		case !ma.Joiner.IsEmpty():
			c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "joiner announcement was not expected"))
		default:
			c.requestedLeave = true
			c.requestedLeaveReason = ma.LeaveReason
		}
	case ma.JoinerID.IsAbsent() != ma.Joiner.IsEmpty():
		if ma.JoinerID.IsAbsent() && !c.IsJoiner() {
			c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "joiner announcement was provided but a joiner was not declared"))
			// } else {
			//	c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "joiner was declared but an announcement was not provided"))
		}
	case c.CanIntroduceJoiner() || ma.JoinerID.IsAbsent():
	case c.IsJoiner():
		c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "joiner can't add a joiner"))
	default:
		c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "restricted/suspended nodes can't add a joiner"))
	}

	switch {
	case c.IsJoiner():
		sp := c.profile.GetStatic().GetStartPower()
		if ma.Membership.RequestedPower != sp {
			c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "start power is different"))
		}
		ma.Membership.RequestedPower = sp
	case ma.Membership.RequestedPower == 0:
		break
	case c.profile.GetStatic().GetExtension() == nil:
		if ma.Membership.RequestedPower != c.profile.GetStatic().GetStartPower() {
			c.RegisterBlame(c.Blames().NewProtocolViolation(c.profile, "unable to verify power"))
			// return false, nil // let the node to be "unset" // TODO handle properly
		}
	case !c.profile.GetStatic().GetExtension().GetPowerLevels().IsAllowed(ma.Membership.RequestedPower):
		return false, c.RegisterFraud(c.Frauds().NewInvalidPowerLevel(c.profile))
	}

	if applyAfterChecks != nil {
		err := applyAfterChecks(ma)
		if err != nil {
			return false, err
		}
	}

	c.stateEvidence = ma.Membership.StateEvidence
	c.announceSignature = ma.Membership.AnnounceSignature
	c.neighbourWeight ^= longbits.FoldUint64(c.stateEvidence.GetDigestHolder().FoldToUint64())

	c.requestedPower = ma.Membership.RequestedPower
	c.requestedJoinerID = ma.JoinerID

	c.hook.OnNodeStateAssigned(updVersion, c)

	return true, nil
}

func (c *NodeAppearance) SetLocalNodeState(ma profiles.MemberAnnouncement) bool {

	if !c.IsLocal() {
		panic(fmt.Sprintf("illegal state - not local: %v", c.GetNodeID()))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isStateAvailable() {
		panic(fmt.Sprintf("illegal state - already set: %v", c.GetNodeID()))
	}

	trustBefore := c.trust

	updated, err := c._applyState(ma, nil)
	if err != nil {
		panic(err)
	}

	c.trust.Update(member.LocalSelfTrust)
	if trustBefore != c.trust {
		c.hook.OnTrustUpdated(c.hook.UpdatePopulationVersion(), c, trustBefore, c.trust, c.profile.HasFullProfile())
	}

	return updated
}

func (c *NodeAppearance) GetNodeMembershipProfile() profiles.MembershipProfile {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isStateAvailable() {
		panic(fmt.Sprintf("illegal state: for=%v", c.GetNodeID()))
	}
	return c.getMembership()
}

func (c *NodeAppearance) GetNodeTrustAndMembershipOrEmpty() (profiles.MembershipProfile, member.TrustLevel) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.getMembership(), c.trust
}

func (c *NodeAppearance) GetNodeMembershipProfileOrEmpty() profiles.MembershipProfile {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getMembership()
}

func (c *NodeAppearance) isStateAvailable() bool {
	return c.stateEvidence != nil
}

func (c *NodeAppearance) GetStatelessAnnouncementEvidence() cryptkit.SignedDigestHolder {

	if c.IsStateful() {
		panic("illegal state")
	}
	// TODO caching
	return c.calcStatelessAnnouncementDigest()
}

func (c *NodeAppearance) calcStatelessAnnouncementDigest() cryptkit.SignedDigestHolder {
	sp := c.profile.GetStatic()
	introDigest := sp.GetBriefIntroSignedDigest()

	// d := c.hook.GetCryptographyAssistant().GetDigestFactory().CreateAnnouncementDigester()
	// d.AddNext(c.hook.GetPulseData().GetPulseDataDigest())
	// d.AddNext(introDigest.GetDigestHolder())
	// return d.FinishSequence().AsDigestHolder()
	return introDigest
}

func (c *NodeAppearance) onAddedToPopulation(fixedInit bool) {

	flags := FlagCreated
	if fixedInit {
		flags |= FlagFixedInit
	}

	full := c.profile.HasFullProfile()
	if full {
		flags |= FlagUpdatedProfile
	}
	pv := c.hook.GetPopulationVersion()
	c.hook.OnDynamicNodeUpdate(pv, c, flags)

	trust := c.trust // this is safe as this method is called before any concurrent access
	if trust != member.UnknownTrust {
		c.hook.OnTrustUpdated(pv, c, member.UnknownTrust, trust, full)
	}
}

func (c *NodeAppearance) IsNSHRequired() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return !c.isStateAvailable()
}

func (c *NodeAppearance) HasAnyPacketReceived() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.limiter.HasAnyPacketReceived()
}

func (c *NodeAppearance) GetNeighbourWeight() uint32 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.neighbourWeight
}

func (c *NodeAppearance) registerFraud(fraud misbehavior.FraudError) (bool, error) {
	if fraud.IsUnknown() {
		panic("empty fraud")
	}

	prevTrust := c.trust
	if c.trust.Update(member.FraudByThisNode) {
		updVersion := c.hook.UpdatePopulationVersion()
		c.firstFraudDetails = &fraud
		c.hook.OnTrustUpdated(updVersion, c, prevTrust, c.trust, c.profile.HasFullProfile())
		return true, fraud
	}
	return false, fraud
}

func (c *NodeAppearance) RegisterFraud(fraud misbehavior.FraudError) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	/* Here the pointer comparison is intentional to ensure exact ActiveNode, as it may change across rounds etc */
	if fraud.ViolatorNode() != c.GetProfile() {
		panic("misplaced fraud")
	}

	_, err := c.registerFraud(fraud)
	return err
}

func (c *NodeAppearance) RegisterBlame(blame misbehavior.BlameError) {
	// TODO RegisterBlame
	// inslogger.FromContext(ctx).Error(blame)
}

func (c *NodeAppearance) getMembership() profiles.MembershipProfile {
	return profiles.NewMembershipProfileByNode(c.profile, c.stateEvidence, c.announceSignature, c.requestedPower)
}

func (c *NodeAppearance) GetNeighborTrustThreshold() uint8 {
	return c.hook.GetNeighbourhoodTrustThreshold()
}

func (c *NodeAppearance) NotifyOnCustom(event interface{}) {
	c.hook.OnCustomEvent(c.hook.GetPopulationVersion(), c, event)
}

func (c *NodeAppearance) GetPacketHandler(i int) DispatchMemberPacketFunc {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.handlers) == 0 {
		return nil
	}
	return c.handlers[i]
}

type NodeRequestedState struct {
	profiles.MembershipProfile
	LeaveReason   uint32
	TrustLevel    member.TrustLevel
	IsLeaving     bool
	RequestedMode member.OpMode
	JoinerID      insolar.ShortNodeID
}

func (c *NodeAppearance) GetRequestedState() NodeRequestedState {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.getMembership()
	if m.Mode.IsEvicted() {
		panic("illegal state")
	}
	if !c.requestedJoinerID.IsAbsent() && !m.CanIntroduceJoiner() {
		panic("illegal state")
	}

	reqMode := member.ModeNormal
	switch {
	case c.requestedLeave:
		reqMode = member.ModeEvictedGracefully
	case c.IsJoiner():
		reqMode = member.ModeRestrictedAnnouncement
	}

	return NodeRequestedState{
		m, c.requestedLeaveReason, c.trust,
		c.requestedLeave, reqMode, c.requestedJoinerID,
	}
}

func (c *NodeAppearance) getMembershipAnnouncement() profiles.MembershipAnnouncement {
	mb := c.getMembership()
	switch {
	case c.requestedLeave:
		return profiles.NewMembershipAnnouncementWithLeave(mb, c.requestedLeaveReason)
	default:
		return profiles.NewMembershipAnnouncementWithJoinerID(mb, c.requestedJoinerID, nil) // TODO joiner secret
	}
}

func (c *NodeAppearance) GetRequestedAnnouncement() profiles.MembershipAnnouncement {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.getMembershipAnnouncement()
}

/* deprecated */ // replace with DispatchAnnouncement
func (c *NodeAppearance) UpgradeDynamicNodeProfile(ctx context.Context, full transport.FullIntroductionReader) bool {
	return c.upgradeDynamicNodeProfile(ctx, full, full)
}

func (c *NodeAppearance) upgradeDynamicNodeProfile(ctx context.Context, brief profiles.BriefCandidateProfile, ext profiles.CandidateProfileExtension) bool {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	match, created := profiles.UpgradeStaticProfile(c.profile.GetStatic(), brief, ext)
	if match && created != nil {
		// here we should check/apply all related attributes
		// TODO handle possible mismatch
		// c.requestedPower = c.profile.GetStatic().GetExtension().GetPowerLevels().FindNearestValid(c.requestedPower)

		inslogger.FromContext(ctx).Debugf("Node profile was upgraded: s=%d, t=%d",
			c.hook.GetLocalNodeID(), c.GetNodeID())

		v := c.hook.UpdatePopulationVersion()
		c.hook.OnDynamicNodeUpdate(v, c, FlagUpdatedProfile)
		c.hook.OnTrustUpdated(v, c, c.trust, c.trust, true)
	}
	return match
}

func (c *NodeAppearance) DispatchAnnouncement(ctx context.Context, rank member.Rank, profile profiles.StaticProfile,
	announcement profiles.MemberAnnouncement) error {

	// TODO additional checks

	if args.IsNil(profile) {
		return nil
	}
	if !c.upgradeDynamicNodeProfile(ctx, profile, profile.GetExtension()) {
		return fmt.Errorf("mismatch")
	}

	return nil
}

func (c *NodeAppearance) GetRequestedLeave() (bool, uint32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.requestedLeave, c.requestedLeaveReason
}
