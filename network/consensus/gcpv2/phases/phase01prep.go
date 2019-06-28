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

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func NewPhase0PrepController(s PulseSelectionStrategy) *Phase0PrepController {
	return &Phase0PrepController{pulseStrategy: s}
}

func NewPhase1PrepController(s PulseSelectionStrategy) *Phase1PrepController {
	return &Phase1PrepController{pulseStrategy: s}
}

var _ core.PrepPhaseController = &Phase0PrepController{}

type Phase0PrepController struct {
	pulseStrategy PulseSelectionStrategy
	realm         *core.PrepRealm
}

func (r *Phase0PrepController) BeforeStart(realm *core.PrepRealm) {
	r.realm = realm
}

func (*Phase0PrepController) GetPacketType() packets.PacketType {
	return packets.PacketPhase0
}

func (r *Phase0PrepController) HandleHostPacket(reader packets.PacketParser, from common.HostIdentityHolder) (postpone bool, err error) {
	p := reader.GetMemberPacket().AsPhase0Packet()
	err = r.pulseStrategy.HandlePrepPulsarPacket(p.GetEmbeddedPulsePacket(), from, r.realm, false)
	return err == nil, err
}

func (*Phase0PrepController) StartWorker(ctx context.Context) {
}

var _ core.PrepPhaseController = &Phase1PrepController{}

type Phase1PrepController struct {
	pulseStrategy PulseSelectionStrategy
	realm         *core.PrepRealm
}

func (r *Phase1PrepController) BeforeStart(realm *core.PrepRealm) {
	r.realm = realm
}

func (*Phase1PrepController) GetPacketType() packets.PacketType {
	return packets.PacketPhase1
}

func (r *Phase1PrepController) HandleHostPacket(reader packets.PacketParser, from common.HostIdentityHolder) (postpone bool, err error) {
	p := reader.GetMemberPacket().AsPhase1Packet()
	err = r.pulseStrategy.HandlePrepPulsarPacket(p.GetEmbeddedPulsePacket(), from, r.realm, false)
	return err == nil, err
}

func (*Phase1PrepController) StartWorker(ctx context.Context) {
}
