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

package pulsectl

import (
	"context"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetrecorder"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

// TODO HACK - network doesnt have information about pulsars to validate packets, the next line must be removed when fixed
const ignoreHostVerificationForPulses = true

func NewPulsePrepController(s PulseSelectionStrategy) *PulsePrepController {
	return &PulsePrepController{pulseStrategy: s}
}

func NewPulseController() *PulseController {
	return &PulseController{}
}

func (p *pulsePacketPrepDispatcher) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags packetrecorder.PacketVerifyFlags) error {

	pp := packet.GetPulsePacket()
	ok, err := p.pulseStrategy.HandlePulsarPacket(ctx, pp, from, true)
	if err != nil || !ok {
		return err
	}
	return p.R.ApplyPulseData(pp, true)
}

func (p *pulsePacketDispatcher) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags packetrecorder.PacketVerifyFlags) error {

	pp := packet.GetPulsePacket()
	// FullRealm already has a pulse data, so should only check it
	pd := pp.GetPulseData()
	if p.R.GetPulseData() == pd {
		return nil
	}
	return p.R.GetBlameFactory().NewMismatchedPulsarPacket(from, p.R.GetOriginalPulse(), pp.GetPulseDataEvidence())
}

func (*pulsePacketPrepDispatcher) HasCustomVerifyForHost(from endpoints.Inbound, strict bool) bool {
	//noinspection GoBoolExpressions
	return ignoreHostVerificationForPulses
}

func (*pulsePacketDispatcher) HasCustomVerifyForHost(from endpoints.Inbound, strict bool) bool {
	//noinspection GoBoolExpressions
	return ignoreHostVerificationForPulses
}

var _ core.PrepPhaseController = &PulsePrepController{}

type PulsePrepController struct {
	core.PrepPhaseControllerTemplate
	pulseStrategy PulseSelectionStrategy
}

func (r *PulsePrepController) CreatePacketDispatcher(pt phases.PacketType, realm *core.PrepRealm) core.PacketDispatcher {
	return &pulsePacketPrepDispatcher{pulseStrategy: r.pulseStrategy, R: realm}
}

func (*PulsePrepController) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPulse}
}

type pulsePacketPrepDispatcher struct {
	core.HostPacketDispatcherTemplate
	pulseStrategy PulseSelectionStrategy
	R             *core.PrepRealm
}

var _ core.PhaseController = &PulseController{}

type PulseController struct {
	core.PhaseControllerTemplate
}

func (c *PulseController) CreatePacketDispatcher(pt phases.PacketType, ctlIndex int, realm *core.FullRealm) (core.PacketDispatcher, core.PerNodePacketDispatcherFactory) {
	return &pulsePacketDispatcher{R: realm}, nil
}

func (*PulseController) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPulse}
}

type pulsePacketDispatcher struct {
	core.HostPacketDispatcherTemplate
	R *core.FullRealm
}
