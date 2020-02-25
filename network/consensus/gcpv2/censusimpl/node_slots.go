// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package censusimpl

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

var _ profiles.LocalNode = &NodeProfileSlot{}

type NodeProfileSlot struct {
	profiles.StaticProfile
	verifier cryptkit.SignatureVerifier
	index    member.Index
	mode     member.OpMode
	power    member.Power
}

func NewNodeProfile(index member.Index, p profiles.StaticProfile, verifier cryptkit.SignatureVerifier, pw member.Power) NodeProfileSlot {

	return NodeProfileSlot{index: index.Ensure(), StaticProfile: p, verifier: verifier, power: pw}
}

func NewJoinerProfile(p profiles.StaticProfile, verifier cryptkit.SignatureVerifier) NodeProfileSlot {

	return NodeProfileSlot{index: member.JoinerIndex, StaticProfile: p, verifier: verifier}
}

func NewNodeProfileExt(index member.Index, p profiles.StaticProfile, verifier cryptkit.SignatureVerifier, pw member.Power,
	mode member.OpMode) NodeProfileSlot {

	return NodeProfileSlot{index: index.Ensure(), StaticProfile: p, verifier: verifier, power: pw, mode: mode}
}

func (c *NodeProfileSlot) GetDeclaredPower() member.Power {
	return c.power
}

func (c *NodeProfileSlot) GetOpMode() member.OpMode {
	return c.mode
}

func (c *NodeProfileSlot) LocalNodeProfile() {
}

func (c *NodeProfileSlot) GetIndex() member.Index {
	return c.index.Ensure()
}

func (c *NodeProfileSlot) IsJoiner() bool {
	return c.index.IsJoiner()
}

func (c *NodeProfileSlot) IsPowered() bool {
	return !c.index.IsJoiner() && !c.mode.IsPowerless() && c.power > 0
}

func (c *NodeProfileSlot) IsVoter() bool {
	return !c.index.IsJoiner() && c.mode.CanVote()
}

func (c *NodeProfileSlot) IsStateful() bool {
	return !c.index.IsJoiner() && c.mode.CanHaveState()
}

func (c *NodeProfileSlot) GetSignatureVerifier() cryptkit.SignatureVerifier {
	return c.verifier
}

func (c *NodeProfileSlot) CanIntroduceJoiner() bool {
	return c.mode.CanIntroduceJoiner(c.index.IsJoiner())
}

func (c *NodeProfileSlot) GetNodeID() insolar.ShortNodeID {
	return c.GetStaticNodeID()
}

func (c *NodeProfileSlot) GetStatic() profiles.StaticProfile {
	return c.StaticProfile
}

func (c *NodeProfileSlot) HasFullProfile() bool {
	return c.StaticProfile.GetExtension() != nil
}

func (c NodeProfileSlot) String() string {
	if c.IsJoiner() {
		return fmt.Sprintf("id:%04d joiner", c.GetNodeID())
	}
	return fmt.Sprintf("id:%04d idx:%d %v", c.GetNodeID(), c.index, c.mode)
}

var _ profiles.Updatable = &updatableSlot{}

type updatableSlot struct {
	NodeProfileSlot
	leaveReason uint32
}

func (c *updatableSlot) AsActiveNode() profiles.ActiveNode {
	return &c.NodeProfileSlot
}

func (c *updatableSlot) SetRank(index member.Index, m member.OpMode, power member.Power) {
	c.index = index.Ensure()
	c.power = power
	c.mode = m
}

func (c *updatableSlot) SetPower(power member.Power) {
	c.power = power
}

func (c *updatableSlot) SetOpMode(m member.OpMode) {
	c.mode = m
}

func (c *updatableSlot) SetOpModeAndLeaveReason(index member.Index, leaveReason uint32) {
	c.index = index.Ensure()
	c.power = 0
	c.mode = member.ModeEvictedGracefully
	c.leaveReason = leaveReason
}

func (c *updatableSlot) GetLeaveReason() uint32 {
	if !c.mode.IsEvictedGracefully() {
		return 0
	}
	return c.leaveReason
}

func (c *updatableSlot) SetIndex(index member.Index) {
	c.index = index.Ensure()
}

func (c *updatableSlot) SetSignatureVerifier(verifier cryptkit.SignatureVerifier) {
	c.verifier = verifier
}

func (c *updatableSlot) IsEmpty() bool {
	return c.StaticProfile == nil
}
