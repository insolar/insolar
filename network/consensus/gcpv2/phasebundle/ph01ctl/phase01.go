// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package ph01ctl

import (
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/announce"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

/*
Does Phase0/Phase1/Phase1Rq processing
*/
func NewPhase01Controller(packetPrepareOptions transport.PacketPrepareOptions, qIntro <-chan population.MemberPacketSender) *Phase01Controller {
	return &Phase01Controller{packetPrepareOptions: packetPrepareOptions, qIntro: qIntro}
}

func (p *packetPhase0Dispatcher) DispatchMemberPacket(ctx context.Context, packet transport.MemberPacketReader, n *population.NodeAppearance) error {

	p0 := packet.AsPhase0Packet()
	nr := p0.GetNodeRank()

	if n.GetRank(p.ctl.R.GetNodeCount()) != nr {
		return n.Frauds().NewMismatchedNeighbourRank(n.GetReportProfile())
	}

	pp := p0.GetEmbeddedPulsePacket()
	return p.ctl.handlePulseData(ctx, pp, n)
}

func (c *packetPhase1Dispatcher) DispatchMemberPacket(ctx context.Context, packet transport.MemberPacketReader, n *population.NodeAppearance) error {

	p1 := packet.AsPhase1Packet()
	_, _, err := announce.ApplyMemberAnnouncement(ctx, p1, nil, true, n, c.ctl.R)
	if err != nil {
		return err
	}

	if p1.HasPulseData() {
		pp := p1.GetEmbeddedPulsePacket()
		err = c.ctl.handlePulseData(ctx, pp, n)
	}
	return err
}

func (c *packetPhase1Dispatcher) TriggerUnknownMember(ctx context.Context, memberID insolar.ShortNodeID,
	packet transport.MemberPacketReader, from endpoints.Inbound) (bool, error) {

	p1 := packet.AsPhase1Packet()

	// if p1.HasPulseData() {
	//	return false, fmt.Errorf("pulse data is not expected")
	// }

	// TODO check endpoint and PK

	// na := p1.GetAnnouncementReader()
	// nr := na.GetNodeRank()
	// if !c.ctl.isJoiner && !nr.IsJoiner() {
	//	return false, fmt.Errorf("member to member intro")
	// }

	return announce.ApplyUnknownAnnouncement(ctx, memberID, p1, nil, true, c.ctl.R)
}

func (c *packetPhase1RqDispatcher) DispatchMemberPacket(ctx context.Context, packet transport.MemberPacketReader, n *population.NodeAppearance) error {

	p1 := packet.AsPhase1Packet()

	_, _, err := announce.ApplyMemberAnnouncement(ctx, p1, nil, false, n, c.ctl.R)
	if err != nil {
		return err
	}

	if p1.HasPulseData() {
		return fmt.Errorf("pulse data is not expected") // TODO blame
	}

	if !c.ctl.R.GetSelf().IsNSHRequired() {
		c.ctl.sendReqReply(ctx, n)
	} else {
		inslogger.FromContext(ctx).Warn("got Phase1Req, but NSH is still unavailable")
	}
	return nil
}

func (c *Phase01Controller) handlePulseData(ctx context.Context, pp transport.PulsePacketReader, n *population.NodeAppearance) error {

	// TODO when PulseDataExt is bigger than ~1.0KB - ignore it, as we will not be able to resend it within Ph0/Ph1 packets
	// TODO validate pulse data
	pp.GetPulseDataEvidence()
	_ = ctx.Err()

	if c.R.GetPulseData() == pp.GetPulseData() {
		return nil
	}
	return n.Blames().NewMismatchedPulsePacket(n.GetProfile(), c.R.GetOriginalPulse(), pp.GetPulseDataEvidence())
}

var _ core.PhaseController = &Phase01Controller{}

type Phase01Controller struct {
	core.PhaseControllerTemplate
	packetPrepareOptions transport.PacketPrepareOptions
	qIntro               <-chan population.MemberPacketSender
	R                    *core.FullRealm
}

func (c *Phase01Controller) CreatePacketDispatcher(pt phases.PacketType, ctlIndex int, realm *core.FullRealm) (population.PacketDispatcher, core.PerNodePacketDispatcherFactory) {
	switch pt {
	case phases.PacketPhase0:
		return &packetPhase0Dispatcher{packetPhase01Dispatcher{ctl: c}}, nil
	case phases.PacketPhase1:
		return &packetPhase1Dispatcher{packetPhase01Dispatcher{ctl: c}}, nil
	case phases.PacketReqPhase1:
		return &packetPhase1RqDispatcher{packetPhase01Dispatcher{ctl: c}}, nil
	default:
		panic("illegal value")
	}
}

func (*Phase01Controller) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase0, phases.PacketPhase1, phases.PacketReqPhase1}
}

func (c *Phase01Controller) sendReqReply(ctx context.Context, target *population.NodeAppearance) { // nolint:interfacer

	p1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.CreateLocalAnnouncement(),
		c.R.GetOriginalPulse(), c.R.GetWelcomePackage(), transport.PrepareWithoutPulseData|c.packetPrepareOptions)

	p1.SendTo(ctx, target, 0, c.R.GetPacketSender())
	target.SetPacketSent(phases.PacketPhase1)
}

func (c *Phase01Controller) StartWorker(ctx context.Context, realm *core.FullRealm) {
	c.R = realm
	go c.workerPhase01(ctx)
}

func (c *Phase01Controller) workerPhase01(ctx context.Context) {

	nodes := c.R.GetPopulation().GetShuffledOtherNodes()

	var nsh proofs.NodeStateHash
	startIndex := 0

	if ok, nshChannel := c.R.PreparePulseChange(); ok {
		nsh, startIndex = c.workerSendPhase0(ctx, nodes, nshChannel)
		if startIndex < 0 {
			// stopped via context
			inslogger.FromContext(ctx).Error(">>>>>>workerPhase01: was stopped via context")
			return
		}
		if nsh == nil {
			panic(">>>>>>workerPhase01: empty NSH")
			// return
		}
		c.R.CommitPulseChange()
	}

	if nsh == nil {
		inslogger.FromContext(ctx).Debugf(">>>>>>workerPhase01: NSH is empty: stateful=%v", c.R.IsLocalStateful())
	}
	inslogger.FromContext(ctx).Debugf(">>>>>>workerPhase01: before NSH update: nsh=%v, self=%+v", nsh, c.R.GetSelf())
	updated := c.R.ApplyLocalState(nsh)
	inslogger.FromContext(ctx).Debugf(">>>>>>workerPhase01: after NSH update: updated=%v, nsh=%v, self=%+v", updated, nsh, c.R.GetSelf())

	go c.workerSendPhase1ToFixed(ctx, startIndex, nodes)
	c.workerSendPhase1ToDynamics(ctx)
}

func (c *Phase01Controller) workerSendPhase0(ctx context.Context, nodes []*population.NodeAppearance,
	nshChannel <-chan api.UpstreamState) (proofs.NodeStateHash, int) {

	/*
		TODO when PulseDataExt is bigger than ~0.7KB then it is too big for Ph1 and can only be distributed
		with Ph0 packets, so Ph0 phase should continue to run
		Also size of Ph1 claims should be considered too.
	*/

	select {
	case <-ctx.Done():
		return nil, -1
	case <-time.After(c.R.AdjustedAfter(c.R.GetTimings().StartPhase0At)):
		break
	case nsh := <-nshChannel:
		return nsh.NodeState, 0
	}

	p0 := c.R.GetPacketBuilder().PreparePhase0Packet(c.R.CreateLocalPhase0Announcement(), c.R.GetOriginalPulse(),
		c.packetPrepareOptions)

	sendOptions := c.packetPrepareOptions.AsSendOptions() &^ transport.SendWithoutPulseData

	for lastIndex, target := range nodes {
		if target.HasAnyPacketReceived() {
			continue
		}

		// DONT use SendToMany here, as this is "gossip" style and parallelism is provided by multiplicity of nodes
		// Instead we have a chance to save some traffic.
		p0.SendTo(ctx, target, sendOptions, c.R.GetPacketSender())
		target.SetPacketSent(phases.PacketPhase0)

		select {
		case <-ctx.Done():
			return nil, -1
		case nsh := <-nshChannel:
			return nsh.NodeState, lastIndex + 1
		default:
		}
	}

	select {
	case <-ctx.Done():
		return nil, -1
	case nsh := <-nshChannel:
		return nsh.NodeState, 0
	}
}

func (c *Phase01Controller) workerSendPhase1ToFixed(ctx context.Context, startIndex int, otherNodes []*population.NodeAppearance) {

	p1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.CreateLocalAnnouncement(),
		c.R.GetOriginalPulse(), c.R.GetWelcomePackage(), c.packetPrepareOptions)

	sendOptions := c.packetPrepareOptions.AsSendOptions()

	from := c.R.GetSelfNodeID()

	// first, send to nodes not covered by Phase 0
	p1.SendToMany(ctx, len(otherNodes)-startIndex, c.R.GetPacketSender(),
		func(ctx context.Context, targetIdx int) (transport.TargetProfile, transport.PacketSendOptions) {
			return prepareTarget(ctx, otherNodes[targetIdx+startIndex], from, sendOptions)
		})

	// then to the rest
	p1.SendToMany(ctx, startIndex, c.R.GetPacketSender(),
		func(ctx context.Context, targetIdx int) (transport.TargetProfile, transport.PacketSendOptions) {
			return prepareTarget(ctx, otherNodes[targetIdx], from, sendOptions)
		})
}

func prepareTarget(ctx context.Context, target *population.NodeAppearance, from insolar.ShortNodeID,
	sendOptions transport.PacketSendOptions) (transport.TargetProfile, transport.PacketSendOptions) {

	if !target.SetPacketSent(phases.PacketPhase1) {
		return nil, 0
	}
	if target.HasAnyPacketReceived() {
		sendOptions |= transport.SendWithoutPulseData
	}
	inslogger.FromContext(ctx).Debugf("sent ph1: from=%d to=%d mode=fix", from, target.GetNodeID())
	return target, sendOptions
}

func (c *Phase01Controller) workerSendPhase1ToDynamics(ctx context.Context) {

	p1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.CreateLocalAnnouncement(),
		c.R.GetOriginalPulse(), c.R.GetWelcomePackage(),
		c.packetPrepareOptions)

	// TODO check if Phase1 packet size is ok to send both intro and pulse data - then, as a backup - send Phase0

	sendOptions := c.packetPrepareOptions.AsSendOptions() | transport.TargetNeedsIntro

	selfID := c.R.GetSelfNodeID()
	for {
		select {
		case <-ctx.Done():
			return
		case introTo := <-c.qIntro:
			nodeID := introTo.GetNodeID()
			if nodeID == selfID {
				continue
			}
			na := c.R.GetPopulation().GetNodeAppearance(nodeID)
			if na != nil {
				introTo = na
			}
			if !introTo.SetPacketSent(phases.PacketPhase1) {
				continue
			}
			inslogger.FromContext(ctx).Debugf("sent ph1: from=%d to=%d mode=dyn", c.R.GetSelfNodeID(), introTo.GetNodeID())
			p1.SendTo(ctx, introTo, sendOptions, c.R.GetPacketSender())
		}
	}
}

type packetPhase01Dispatcher struct {
	core.MemberPacketDispatcherTemplate
	ctl *Phase01Controller
}

type packetPhase0Dispatcher struct {
	packetPhase01Dispatcher
}

type packetPhase1Dispatcher struct {
	packetPhase01Dispatcher
}

type packetPhase1RqDispatcher struct {
	packetPhase01Dispatcher
}
