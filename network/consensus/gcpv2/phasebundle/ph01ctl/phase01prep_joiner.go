// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package ph01ctl

import (
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/pulsectl"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func NewJoinerPhase01PrepController(s pulsectl.PulseSelectionStrategy) *JoinerPhase01PrepController {
	return &JoinerPhase01PrepController{pulseStrategy: s}
}

var _ core.PrepPhaseController = &JoinerPhase01PrepController{}

type JoinerPhase01PrepController struct {
	core.PrepPhaseControllerTemplate
	core.HostPacketDispatcherTemplate

	realm         *core.PrepRealm
	pulseStrategy pulsectl.PulseSelectionStrategy
}

func (c *JoinerPhase01PrepController) CreatePacketDispatcher(pt phases.PacketType, realm *core.PrepRealm) population.PacketDispatcher {
	c.realm = realm
	return c
}

func (c *JoinerPhase01PrepController) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase0, phases.PacketPhase1}
}

func (c *JoinerPhase01PrepController) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags coreapi.PacketVerifyFlags) error {

	var pp transport.PulsePacketReader
	var nr member.Rank
	mp := packet.GetMemberPacket()

	switch packet.GetPacketType() {
	case phases.PacketPhase0:
		p0 := mp.AsPhase0Packet()
		// can take Phase0 packets only for nodes who has sent Phase1 ...
		nr = p0.GetNodeRank()
	case phases.PacketPhase1:
		p1 := mp.AsPhase1Packet()
		nr = p1.GetAnnouncementReader().GetNodeRank()
		if p1.HasPulseData() {
			pp = p1.GetEmbeddedPulsePacket()
		}
	default:
		panic("not expected")
	}
	if nr.IsJoiner() {
		if pp != nil {
			return fmt.Errorf("pulse data is not allowed from a joiner: from=%v", from)
		}
		return nil // postpone the packet
	}
	if packet.GetPacketType() != phases.PacketPhase1 {
		return nil // postpone the packet
	}

	p1 := mp.AsPhase1Packet()
	if !p1.HasFullIntro() || !p1.HasCloudIntro() {
		return fmt.Errorf("joiner expects full & cloud intro in Phase1: from=%v", from)
	}

	mr := c.realm.GetMandateRegistry()
	ci := p1.GetCloudIntroduction()

	if !mr.GetCloudIdentity().Equals(ci.GetCloudIdentity()) {
		return fmt.Errorf("mismatched cloud identity: from=%v", from)
	}

	// TODO collect a few proposals and choose only by getting some threshold
	// if p1.HasJoinerSecret() {
	//	c.realm.IsValidJoinerSecret()
	// }
	//
	//

	populationCount := int(nr.GetTotalCount())
	lastCloudStateHash := ci.GetLastCloudStateHash()

	if populationCount == 0 {
		return fmt.Errorf("node count cant be zero: from=%v", from)
	}
	if lastCloudStateHash == nil {
		return fmt.Errorf("packet error - last cloud state hash is missing: from=%v", from)
	}

	c.realm.ApplyCloudIntro(lastCloudStateHash, int(nr.GetTotalCount()), from)

	// TODO joiner should wait for CloudIntro also!
	ok, err := c.pulseStrategy.HandlePulsarPacket(ctx, pp, from, false)
	if err != nil || !ok || pp == nil {
		return err
	}

	startedAt := time.Now() // TODO get packet's receive time
	return c.realm.ApplyPulseData(ctx, startedAt, pp, false, from)
}
