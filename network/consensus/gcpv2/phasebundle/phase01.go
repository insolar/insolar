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

package phasebundle

import (
	"context"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func NewPhase0Controller() *Phase0Controller {
	return &Phase0Controller{}
}

func NewPhase1Controller(packetPrepareOptions transport.PacketSendOptions) *Phase1Controller {
	return &Phase1Controller{packetPrepareOptions: packetPrepareOptions}
}

func NewReqPhase1Controller(packetPrepareOptions transport.PacketSendOptions, delegate *Phase1Controller) *ReqPhase1Controller {
	return &ReqPhase1Controller{packetPrepareOptions: packetPrepareOptions, delegate: delegate}
}

var _ core.PhaseController = &Phase0Controller{}

type Phase0Controller struct {
	core.PhaseControllerTemplate
}

func (*Phase0Controller) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase0}
}

func (c *Phase0Controller) CreatePacketDispatcher(pt phases.PacketType, ctlIndex int,
	realm *core.FullRealm) (core.PacketDispatcher, core.PerNodePacketDispatcherFactory) {

	return &packetPhase0Dispatcher{}, nil
}

type packetPhase0Dispatcher struct {
}

func (p *packetPhase0Dispatcher) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags core.PacketVerifyFlags) error {

	p0 := packet.GetMemberPacket().AsPhase0Packet()
	pp := p0.GetEmbeddedPulsePacket()

	//TODO check NodeRank - especially for suspected node

	//TODO when PulseDataExt is bigger than ~1.5KB - ignore it, as we will not be able to resend it within Ph0/Ph1 packets

	err := n.SetPacketReceivedWithDupError(c.GetPacketType())
	return handleEmbeddedPulsePacket(ctx, p, pp, n, c.R, err)
}

func handleEmbeddedPulsePacket(ctx context.Context, p transport.MemberPacketReader, pp transport.PulsePacketReader, n *core.NodeAppearance,
	r *core.FullRealm, defErr error) error {

	// TODO validate pulse data
	pp.GetPulseDataEvidence()
	p.GetPacketSignature()
	_ = ctx.Err()

	if r.GetPulseData() == pp.GetPulseData() {
		return defErr
	}
	return n.Blames().NewMismatchedPulsePacket(n.GetProfile(), r.GetOriginalPulse(), pp.GetPulseDataEvidence())
}

var _ core.PhaseController = &Phase1Controller{}

type Phase1Controller struct {
	packetPrepareOptions transport.PacketSendOptions
}

func (*Phase1Controller) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase1, phases.PacketReqPhase1}
}

func (c *Phase1Controller) HandleMemberPacket(ctx context.Context, p transport.MemberPacketReader, n *core.NodeAppearance) error {
	p1 := p.AsPhase1Packet()
	err := c.handleNodeData(p1, n)

	if err == nil && p1.HasPulseData() {
		pp := p1.GetEmbeddedPulsePacket()
		err = handleEmbeddedPulsePacket(ctx, p, pp, n, c.R, nil)
	}
	return err
}

/* Is also used by ReqPhase1Controller */
func (c *Phase1Controller) handleNodeData(p1 transport.Phase1PacketReader, n *core.NodeAppearance) error {
	dupErr := n.SetPacketReceivedWithDupError(c.GetPacketType())

	// if p1.HasSelfIntro() {
	// TODO register protocol misbehavior - IntroClaim was not expected

	na := p1.GetAnnouncementReader()
	nr := na.GetNodeRank()
	mp := profiles.NewMembershipProfile(nr.GetMode(), nr.GetPower(), nr.GetIndex(),
		na.GetNodeStateHashEvidence(), na.GetAnnouncementSignature(), na.GetRequestedPower())

	if c.R.GetNodeCount() != int(nr.GetTotalCount()) {
		_, err := n.RegisterFraud(n.Frauds().NewMismatchedMembershipRank(n.GetProfile(), mp))
		return err
	}

	var ma profiles.MembershipAnnouncement
	switch {
	case na.IsLeaving():
		ma = profiles.NewMembershipAnnouncementWithLeave(mp, na.GetLeaveReason())
	case na.GetJoinerID().IsAbsent():
		ma = profiles.NewMembershipAnnouncement(mp)
	default:
		panic("not implemented") //TODO implement
		//jar := na.GetJoinerAnnouncement()
		//ma = common.NewMembershipAnnouncementWithJoiner(mp)
	}

	_, err := n.ApplyNodeMembership(ma)

	if err != nil {
		return err
	}
	return dupErr
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

func (c *Phase1Controller) workerSendPhase0(ctx context.Context) (proofs.NodeStateHash, int) {

	nshChannel := c.R.UpstreamPreparePulseChange()
	/*
		TODO when PulseDataExt is bigger than ~1KB then it is too big for Ph1 and can only be distributed with Ph0 packets, so Ph0 phase should continue to run
		Also size of Ph1 claims should be considered too.
	*/
	var nsh proofs.NodeStateHash

	select {
	case <-ctx.Done():
		return nil, -1
	case <-time.After(c.R.AdjustedAfter(c.R.GetTimings().StartPhase0At)):
		break
	case nsh = <-nshChannel:
		return nsh, 0
	}

	p0 := c.R.GetPacketBuilder().PreparePhase0Packet(c.R.CreateLocalAnnouncement(), c.R.GetOriginalPulse(),
		c.packetPrepareOptions)

	for lastIndex, target := range c.R.GetPopulation().GetShuffledOtherNodes() {
		if target.HasAnyPacketReceived() {
			continue
		}

		//DONT use SendToMany here, as this is "gossip" style and parallelism is provided by multiplicity of nodes
		//Instead we have a chance to save some traffic.
		p0.SendTo(ctx, target.GetProfile(), 0, c.R.GetPacketSender())
		target.SetPacketSent(phases.PacketPhase0)

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

	p1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.CreateLocalAnnouncement(),
		c.R.GetOriginalPulse(), nil, c.packetPrepareOptions)

	otherNodes := c.R.GetPopulation().GetShuffledOtherNodes()

	//first, send to nodes not covered by Phase 0
	p1.SendToMany(ctx, len(otherNodes)-startIndex, c.R.GetPacketSender(),
		func(ctx context.Context, targetIdx int) (profiles.ActiveNode, transport.PacketSendOptions) {
			target := otherNodes[targetIdx+startIndex]
			target.SetPacketSent(c.GetPacketType())

			if target.HasAnyPacketReceived() {
				return target.GetProfile(), transport.SendWithoutPulseData
			}
			return target.GetProfile(), 0
		})

	p1.SendToMany(ctx, startIndex, c.R.GetPacketSender(),
		func(ctx context.Context, targetIdx int) (profiles.ActiveNode, transport.PacketSendOptions) {
			target := otherNodes[targetIdx]
			target.SetPacketSent(c.GetPacketType())

			if target.HasAnyPacketReceived() {
				return target.GetProfile(), transport.SendWithoutPulseData
			}
			return target.GetProfile(), 0
		})
}

var _ core.PhaseController = &ReqPhase1Controller{}

type ReqPhase1Controller struct {
	core.PhaseControllerPerMemberTemplate
	delegate             *Phase1Controller
	packetPrepareOptions transport.PacketSendOptions
}

func (c *ReqPhase1Controller) GetPacketType() phases.PacketType {
	return phases.PacketReqPhase1
}

func (c *ReqPhase1Controller) HandleMemberPacket(ctx context.Context, p transport.MemberPacketReader, n *core.NodeAppearance) error {
	p1 := p.AsPhase1Packet()
	err := c.delegate.handleNodeData(p1, n)
	if err != nil {
		return err
	}
	if !c.R.GetSelf().IsNshRequired() {
		c.sendReqPhase1Reply(ctx, n)
	} else {
		inslogger.FromContext(ctx).Warn("got Phase1 request, but NSH is still unavailable")
	}
	return nil
}

func (c *ReqPhase1Controller) sendReqPhase1Reply(ctx context.Context, target *core.NodeAppearance) {

	p1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.CreateLocalAnnouncement(),
		c.R.GetOriginalPulse(), nil, transport.SendWithoutPulseData|c.packetPrepareOptions)

	p1.SendTo(ctx, target.GetProfile(), 0, c.R.GetPacketSender())
	target.SetPacketSent(c.GetPacketType())
}
