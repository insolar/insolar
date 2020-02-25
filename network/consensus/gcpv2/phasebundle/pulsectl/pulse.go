// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulsectl

import (
	"context"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

// TODO HACK - network doesnt have information about pulsars to validate packets, the next line must be removed when fixed

func NewPulsePrepController(s PulseSelectionStrategy, ignoreHostVerificationForPulses bool) *PulsePrepController {
	return &PulsePrepController{pulseStrategy: s, ignoreHostVerificationForPulses: ignoreHostVerificationForPulses}
}

func NewPulseController(ignoreHostVerificationForPulses bool) *PulseController {
	return &PulseController{ignoreHostVerificationForPulses: ignoreHostVerificationForPulses}
}

func (p *PulsePrepController) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags coreapi.PacketVerifyFlags) error {

	pp := packet.GetPulsePacket()
	ok, err := p.pulseStrategy.HandlePulsarPacket(ctx, pp, from, true)
	if err != nil || !ok {
		return err
	}
	startedAt := time.Now() // TODO get packet's receive time
	return p.R.ApplyPulseData(ctx, startedAt, pp, true, from)
}

func (p *PulseController) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags coreapi.PacketVerifyFlags) error {

	pp := packet.GetPulsePacket()
	// FullRealm already has a pulse data, so should only check it
	pd := pp.GetPulseData()
	if p.R.GetPulseData() == pd {
		return nil
	}
	return p.R.MonitorOtherPulses(pp, from)
}

func (p *PulsePrepController) HasCustomVerifyForHost(from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) bool {
	return p.ignoreHostVerificationForPulses
}

func (p *PulseController) HasCustomVerifyForHost(from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) bool {
	return p.ignoreHostVerificationForPulses
}

var _ core.PrepPhaseController = &PulsePrepController{}

type PulsePrepController struct {
	core.PrepPhaseControllerTemplate
	core.HostPacketDispatcherTemplate
	R                               *core.PrepRealm
	pulseStrategy                   PulseSelectionStrategy
	ignoreHostVerificationForPulses bool
}

func (p *PulsePrepController) CreatePacketDispatcher(pt phases.PacketType, realm *core.PrepRealm) population.PacketDispatcher {
	p.R = realm
	return p
}

func (*PulsePrepController) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPulsarPulse}
}

var _ core.PhaseController = &PulseController{}

type PulseController struct {
	core.PhaseControllerTemplate
	core.HostPacketDispatcherTemplate
	R                               *core.FullRealm
	ignoreHostVerificationForPulses bool
}

func (p *PulseController) CreatePacketDispatcher(pt phases.PacketType, ctlIndex int, realm *core.FullRealm) (population.PacketDispatcher, core.PerNodePacketDispatcherFactory) {
	p.R = realm
	return p, nil
}

func (*PulseController) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPulsarPulse}
}
