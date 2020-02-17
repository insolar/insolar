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

func NewPhase01PrepController(s pulsectl.PulseSelectionStrategy) *Phase01PrepController {
	return &Phase01PrepController{pulseStrategy: s}
}

var _ core.PrepPhaseController = &Phase01PrepController{}

type Phase01PrepController struct {
	core.PrepPhaseControllerTemplate
	core.HostPacketDispatcherTemplate

	realm         *core.PrepRealm
	pulseStrategy pulsectl.PulseSelectionStrategy
}

func (c *Phase01PrepController) CreatePacketDispatcher(pt phases.PacketType, realm *core.PrepRealm) population.PacketDispatcher {
	c.realm = realm
	return c
}

func (c *Phase01PrepController) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase0, phases.PacketPhase1}
}

func (c *Phase01PrepController) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags coreapi.PacketVerifyFlags) error {

	var pp transport.PulsePacketReader
	var nr member.Rank

	switch packet.GetPacketType() {
	case phases.PacketPhase0:
		p0 := packet.GetMemberPacket().AsPhase0Packet()
		nr = p0.GetNodeRank()
		pp = p0.GetEmbeddedPulsePacket()
	case phases.PacketPhase1:
		p1 := packet.GetMemberPacket().AsPhase1Packet()
		// if p1.HasFullIntro() || p1.HasCloudIntro() || p1.HasJoinerSecret() {
		//	return fmt.Errorf("introduction data were not expected: from=%v", from)
		// }
		nr = p1.GetAnnouncementReader().GetNodeRank()
		if p1.HasPulseData() {
			pp = p1.GetEmbeddedPulsePacket()
		}
	default:
		panic("illegal value")
	}
	if nr.IsJoiner() && pp != nil {
		return fmt.Errorf("pulse data in Phase0/Phas1 is not allowed from a joiner: from=%v", from)
	}
	// if { TODO check ranks? }

	ok, err := c.pulseStrategy.HandlePulsarPacket(ctx, pp, from, false)
	if err != nil || !ok {
		return err
	}

	startedAt := time.Now()
	return c.realm.ApplyPulseData(ctx, startedAt, pp, false, from)
}
