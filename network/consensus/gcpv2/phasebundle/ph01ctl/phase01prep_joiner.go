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

package ph01ctl

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetrecorder"
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

func (c *JoinerPhase01PrepController) CreatePacketDispatcher(pt phases.PacketType, realm *core.PrepRealm) core.PacketDispatcher {
	c.realm = realm
	return c
}

func (c *JoinerPhase01PrepController) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase0, phases.PacketPhase1}
}

func (c *JoinerPhase01PrepController) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags packetrecorder.PacketVerifyFlags) error {

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
	//if p1.HasJoinerSecret() {
	//	c.realm.IsValidJoinerSecret()
	//}
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

	return c.realm.ApplyPulseData(pp, false, from)
}
