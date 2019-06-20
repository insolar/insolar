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

package phases

import (
	"context"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

func NewPhase0Controller() *Phase0Controller {
	return &Phase0Controller{}
}

func NewPhase1Controller(packetPrepareOptions core.PacketSendOptions, queueNshReady chan<- *core.NodeAppearance) *Phase1Controller {
	return &Phase1Controller{packetPrepareOptions: packetPrepareOptions, qNshReady: queueNshReady}
}

func NewReqPhase1Controller(packetPrepareOptions core.PacketSendOptions, delegate *Phase1Controller) *ReqPhase1Controller {
	return &ReqPhase1Controller{packetPrepareOptions: packetPrepareOptions, delegate: delegate}
}

var _ core.PhaseController = &Phase0Controller{}

type Phase0Controller struct {
	core.PhaseControllerPerMemberTemplate
}

func (*Phase0Controller) GetPacketType() packets.PacketType {
	return packets.PacketPhase0
}

func (c *Phase0Controller) HandleMemberPacket(p packets.MemberPacketReader, n *core.NodeAppearance) error {
	p0 := p.AsPhase0Packet()
	pp := p0.GetEmbeddedPulsePacket()

	err := n.SetReceivedWithDupCheck(c.GetPacketType())
	return handleEmbeddedPulsePacket(p, pp, n, c.R, err)
}

func handleEmbeddedPulsePacket(p packets.MemberPacketReader, pp packets.PulsePacketReader, n *core.NodeAppearance,
	r *core.FullRealm, defErr error) error {

	// TODO validate pulse data
	pp.GetPulseDataEvidence()
	p.GetPacketSignature()

	if r.GetPulseData() == pp.GetPulseData() {
		return defErr
	}
	return r.Blames().NewMismatchedPulsePacket(n.GetProfile(), r.GetOriginalPulse(), pp.GetPulseDataEvidence())
}

var _ core.PhaseController = &Phase1Controller{}

type Phase1Controller struct {
	core.PhaseControllerPerMemberTemplate
	qNshReady            chan<- *core.NodeAppearance
	packetPrepareOptions core.PacketSendOptions
}

func (*Phase1Controller) GetPacketType() packets.PacketType {
	return packets.PacketPhase1
}

func (c *Phase1Controller) HandleMemberPacket(p packets.MemberPacketReader, n *core.NodeAppearance) error {
	p1 := p.AsPhase1Packet()
	err := c.handleNodeData(p1, n)

	if err == nil && p1.HasPulseData() {
		pp := p1.GetEmbeddedPulsePacket()
		err = handleEmbeddedPulsePacket(p, pp, n, c.R, nil)
	}
	return err
}

/* Is also used by ReqPhase1Controller */
func (c *Phase1Controller) handleNodeData(p1 packets.Phase1PacketReader, n *core.NodeAppearance) error {
	send, err := c._handleNodeData(p1, n)
	if send {
		c.qNshReady <- n
	}
	return err
}

func (c *Phase1Controller) _handleNodeData(p1 packets.Phase1PacketReader, n *core.NodeAppearance) (bool, error) {
	dupErr := n.SetReceivedWithDupCheck(c.GetPacketType())

	//if p1.HasSelfIntro() {
	// TODO register protocol misbehavior - IntroClaim was not expected
	//}
	// if c.R.GetNodeCount() != int(p1.GetNodeCount()) {
	// 	//TODO SEND fraud state to others (to Phase2)
	// 	return false, n.RegisterFraud(R.Frauds().NewMismatchedRank(n.GetProfile(), p1.GetNodeStateHashEvidence()))
	// }

	mp := common.NewMembershipProfile(p1.GetNodeIndex(), p1.GetNodePower(), p1.GetNodeStateHash())
	modified, err := n.ApplyNodeMembership(mp, p1.GetNodeStateHashEvidence(), c.R.GetMisbehaviorFactories())
	if modified && dupErr != nil {
		c.R.Log().Warnf("unexpected state: Phase1 was received, but NSH is unset: node=%v", n)
	}

	if err != nil {
		return modified, err
	} else {
		return modified, dupErr
	}
}

func (c *Phase1Controller) StartWorker(ctx context.Context) {
	go c.workerPhase01(ctx)
}

func (c *Phase1Controller) workerPhase01(ctx context.Context) {
	nsh, startIndex := c.workerSendPhase0(ctx)
	if startIndex < 0 {
		return
	}
	c.R.PrepareAndSetLocalNodeStateHashEvidence(nsh)

	c.workerSendPhase1(ctx, startIndex)
}

func (c *Phase1Controller) workerSendPhase0(ctx context.Context) (common.NodeStateHash, int) {

	nshChannel := c.R.UpstreamPreparePulseChange()
	// batchSize := c.R.strategy.GetPhase01SendBatching()

	var nsh common.NodeStateHash

	select {
	case <-ctx.Done():
		return nil, -1
	case <-time.After(c.R.AdjustedAfter(c.R.GetTimings().StartPhase0At)):
		break
	case nsh = <-nshChannel:
		return nsh, 0
	}

	p0 := c.R.GetPacketBuilder().PreparePhase0Packet(c.R.GetLocalProfile(), c.R.GetOriginalPulse(), c.packetPrepareOptions)

	for lastIndex, target := range c.R.GetShuffledOtherNodes() {
		if !target.HasReceivedAnyPhase() {
			p0.SendTo(target.GetProfile(), 0, c.R.GetPacketSender())
			target.SetSentPhase(packets.Phase0)
		}
		select {
		case <-ctx.Done():
			return nil, -1
		case nsh = <-nshChannel:
			return nsh, lastIndex + 1
		default:
		}
	}

	select {
	case <-ctx.Done():
		return nil, -1
	case nsh = <-nshChannel:
		return nsh, 0
	}
}

func (c *Phase1Controller) workerSendPhase1(ctx context.Context, startIndex int) {

	p1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.GetLocalProfile(), c.R.GetOriginalPulse(),
		c.R.GetSelf().GetNodeStateHashEvidence(), c.packetPrepareOptions)

	otherNodes := c.R.GetShuffledOtherNodes()
	for i := range otherNodes {
		index := (startIndex + i) % len(otherNodes)
		target := otherNodes[index]
		var sendOptions core.PacketSendOptions

		if target.HasReceivedAnyPhase() {
			// if something was received from this node, then we don't need to send a copy of pulse data to it
			sendOptions |= core.SendWithoutPulseData
		}
		p1.SendTo(target.GetProfile(), sendOptions, c.R.GetPacketSender())
		target.SetSentByPacketType(c.GetPacketType())
		select {
		case <-ctx.Done():
			return // ctx.Err() ?
		default:
		}
	}
}

var _ core.PhaseController = &ReqPhase1Controller{}

type ReqPhase1Controller struct {
	core.PhaseControllerPerMemberTemplate
	delegate             *Phase1Controller
	packetPrepareOptions core.PacketSendOptions
}

func (c *ReqPhase1Controller) GetPacketType() packets.PacketType {
	return packets.PacketReqPhase1
}

func (c *ReqPhase1Controller) HandleMemberPacket(p packets.MemberPacketReader, n *core.NodeAppearance) error {
	p1 := p.AsPhase1Packet()
	err := c.delegate.handleNodeData(p1, n)
	if err != nil {
		return err
	}
	if !c.R.GetSelf().IsNshRequired() {
		c.sendReqPhase1Reply(n)
	} else {
		c.R.Log().Warn("got Phase1 request, but NSH is still unavailable")
	}
	return nil
}

func (c *ReqPhase1Controller) sendReqPhase1Reply(target *core.NodeAppearance) {

	p1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.GetLocalProfile(), c.R.GetOriginalPulse(),
		c.R.GetSelf().GetNodeStateHashEvidence(), core.SendWithoutPulseData|c.packetPrepareOptions)

	p1.SendTo(target.GetProfile(), 0, c.R.GetPacketSender())
	target.SetSentByPacketType(c.GetPacketType())
}
