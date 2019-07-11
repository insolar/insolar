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

package core

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

/*
	PrepRealm is a functionally limited and temporary realm that is used when this node doesn't know pulse or last consensus.
	It can ONLY pre-processed packets, but is disallowed to send them.

	Pre-processed packets as postponed by default and processing will be repeated when FullRealm is activated.
*/
type PrepRealm struct {
	/* Provided externally. Don't need mutex */
	*coreRealm                              // points the core part realms, it is shared between of all Realms of a Round
	completeFn        func(successful bool) //
	postponedPacketFn postponedPacketFunc

	/* Derived from the provided externally - set at init() or start(). Don't need mutex */
	handlers    []PrepPhasePacketHandler
	queueToFull chan postponedPacket

	/* Other fields - need mutex */
	// 	censusBuilder census.Builder
}

/* LOCK - runs under RoundController lock */
func (p *PrepRealm) start(ctx context.Context, controllers []PrepPhaseController, prepToFullQueueSize int) {

	if p.postponedPacketFn != nil {
		p.queueToFull = make(chan postponedPacket, prepToFullQueueSize)
	}

	p.handlers = make([]PrepPhasePacketHandler, gcp_types.MaxPacketType)
	for _, ctl := range controllers {
		pt := ctl.GetPacketType()
		if p.handlers[pt] != nil {
			panic("multiple handlers for packet type")
		}
		p.handlers[pt] = ctl.HandleHostPacket
	}

	for _, ctl := range controllers {
		ctl.BeforeStart(p)
	}
	for _, ctl := range controllers {
		ctl.StartWorker(ctx)
	}
}

/* LOCK - runs under RoundController lock */
func (p *PrepRealm) stop() {
	/*
		NB! do not close p.queueToFull here immediately, as some messages can still be in processing and will be lost
	*/
	if p.postponedPacketFn != nil {
		/* Do not give out a PrepRealm reference to avoid retention in memory */
		go flushQueueTo(p.coreRealm.roundContext, p.queueToFull, p.postponedPacketFn)
	}
}

type postponedPacketFunc func(packet packets.PacketParser, from endpoints.HostIdentityHolder)

type postponedPacket struct {
	packet packets.PacketParser
	from   endpoints.HostIdentityHolder
}

func flushQueueTo(ctx context.Context, in chan postponedPacket, out postponedPacketFunc) {
	for {
		select {
		case p, ok := <-in:
			if !ok {
				return
			}
			out(p.packet, p.from)
		case <-ctx.Done():
			return
		}
	}
}

func (p *PrepRealm) GetOriginalPulse() packets.OriginalPulsarPacket {
	p.RLock()
	defer p.RUnlock()

	// locks are only needed for PrepRealm
	return p.coreRealm.originalPulse
}

func (p *PrepRealm) ApplyPulseData(pp packets.PulsePacketReader, fromPulsar bool) error {
	pd := pp.GetPulseData()

	p.Lock()
	defer p.Unlock()

	valid := false
	switch {
	case p.originalPulse != nil:
		if pd == p.pulseData {
			return nil // got it already
		}
	case fromPulsar || !p.strategy.IsEphemeralPulseAllowed():
		// Pulsars are NEVER ALLOWED to send ephemeral pulses
		valid = pd.IsValidPulsarData()
	default:
		valid = pd.IsValidPulseData()
	}
	if !valid {
		// if fromPulsar
		// TODO blame pulsar and/or node
		return fmt.Errorf("invalid pulse data")
	}
	epn := p.GetExpectedPulseNumber()
	if !epn.IsUnknown() && epn != pd.PulseNumber {
		return fmt.Errorf("unexpected pulse number: expected=%v, received=%v", epn, pd.PulseNumber)
	}

	p.originalPulse = pp.GetPulseDataEvidence()
	p.pulseData = pd

	p.completeFn(true)

	return nil
}

func (p *PrepRealm) GetExpectedPulseNumber() pulse_data.PulseNumber {
	return p.initialCensus.GetExpectedPulseNumber()
}

func (p *PrepRealm) handleHostPacket(ctx context.Context, packet packets.PacketParser, from endpoints.HostIdentityHolder) error {
	pt := packet.GetPacketType()
	h := p.handlers[pt]

	var explicitPostpone = false
	var err error

	if h != nil {
		explicitPostpone, err = h(ctx, packet, from)
		if !explicitPostpone || err != nil {
			return err
		}
	}
	// if packet is not handled, then we may need to leave it for FullRealm
	if p.postponePacket(packet, from) {
		return nil
	}
	if explicitPostpone {
		return fmt.Errorf("unable to postpone packet explicitly: type=%v", pt)
	}
	return errPacketIsNotAllowed
}

func (p *PrepRealm) postponePacket(packet packets.PacketParser, from endpoints.HostIdentityHolder) bool {
	if p.queueToFull == nil {
		return false
	}
	p.queueToFull <- postponedPacket{packet: packet, from: from}
	return true
}
