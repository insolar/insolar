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

package tests

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle"
	"github.com/insolar/insolar/pulse"
)

func NewConsensusHost(hostAddr endpoints.Name) *EmuHostConsensusAdapter {
	return &EmuHostConsensusAdapter{hostAddr: hostAddr}
}

type EmuHostConsensusAdapter struct {
	controller api.ConsensusController

	hostAddr endpoints.Name
	inbound  <-chan Packet
	outbound chan<- Packet
}

func (h *EmuHostConsensusAdapter) ConnectTo(chronicles api.ConsensusChronicles, network *EmuNetwork,
	strategyFactory core.RoundStrategyFactory, candidateFeeder api.CandidateControlFeeder,
	controlFeeder api.ConsensusControlFeeder, ephemeralFeeder api.EphemeralControlFeeder, config api.LocalNodeConfiguration) {

	ctx := network.ctx
	// &EmuConsensusStrategy{ctx: ctx}
	upstream := NewEmuUpstreamPulseController(ctx, defaultNshGenerationDelay)

	h.controller = gcpv2.NewConsensusMemberController(
		chronicles, upstream,
		core.NewPhasedRoundControllerFactory(config, NewEmuTransport(h), strategyFactory),
		candidateFeeder,
		controlFeeder,
		ephemeralFeeder,
	)

	h.inbound, h.outbound = network.AddHost(ctx, h.hostAddr)
	go h.run(ctx)
}

func (h *EmuHostConsensusAdapter) run(ctx context.Context) {
	defer func() {
		// r := recover()
		// inslogger.FromContext(ctx).Errorf("host has died: %v, %v", h.hostAddr, r)
		// TODO print stacktrace
		close(h.outbound)
	}()

	for {
		var err error
		payload, from, err := h.receive(ctx)
		if err == nil {
			if payload == nil && from == nil {
				h.controller.Abort()
				return
			}

			var packet transport.PacketParser

			packet, err = h.parsePayload(payload)
			if err == nil {
				if packet != nil {
					hostFrom := endpoints.InboundConnection{Addr: *from}

					sourceID := packet.GetSourceID()
					targetID := packet.GetTargetID()

					if sourceID != 0 && sourceID == targetID { // TODO for debugging
						panic("must not")
					}

					err = h.controller.ProcessPacket(ctx, packet, &hostFrom)
				}
			}
		}

		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
	}
}

func (h *EmuHostConsensusAdapter) SendPacketToTransport(ctx context.Context, t transport.TargetProfile, sendOptions transport.PacketSendOptions, payload interface{}) {
	h.send(t.GetStatic().GetDefaultEndpoint(), payload)
}

func (h *EmuHostConsensusAdapter) receive(ctx context.Context) (payload interface{}, from *endpoints.Name, err error) {
	packet, ok := <-h.inbound
	if !ok {
		inslogger.FromContext(ctx).Debugf("host is dead: %s", h.hostAddr)
		return nil, nil, nil
	}
	inslogger.FromContext(ctx).Debugf("receivedBy: %s - %+v", h.hostAddr, packet)
	if packet.Payload == nil {
		return nil, &packet.Host, errors.New("missing payload")
	}
	err, ok = packet.Payload.(error)
	if ok {
		return nil, &packet.Host, err
	}
	return packet.Payload, &packet.Host, nil
}

func (h *EmuHostConsensusAdapter) send(target endpoints.Outbound, payload interface{}) {
	defer func() {
		_ = recover()
	}()
	parser := payload.(transport.PacketParser)
	pkt := Packet{Host: target.GetNameAddress(), Payload: WrapPacketParser(parser)}
	// fmt.Println(">SEND> ", pkt)
	h.outbound <- pkt
}

func (h *EmuHostConsensusAdapter) parsePayload(payload interface{}) (transport.PacketParser, error) {
	return UnwrapPacketParser(payload), nil
}

func (h *EmuHostConsensusAdapter) TransportPacketSender() {
}

type EmuRoundStrategyFactory struct {
	roundStrategy EmuRoundStrategy
	bundleFactory core.PhaseControllersBundleFactory
}

func (p *EmuRoundStrategyFactory) CreateRoundStrategy(chronicle api.ConsensusChronicles,
	config api.LocalNodeConfiguration) (core.RoundStrategy, core.PhaseControllersBundle) {

	if p.bundleFactory == nil {
		p.bundleFactory = phasebundle.NewStandardBundleFactoryDefault()
	}

	lastCensus, _ := chronicle.GetLatestCensus()
	pop := lastCensus.GetOnlinePopulation()
	bundle := p.bundleFactory.CreateControllersBundle(pop, config)
	return &p.roundStrategy, bundle
}

type EmuRoundStrategy struct {
}

func (*EmuRoundStrategy) IsEphemeralPulseAllowed() bool {
	return false
}

func (*EmuRoundStrategy) ConfigureRoundContext(ctx context.Context, expectedPulse pulse.Number, self profiles.LocalNode) context.Context {
	return ctx
}

func (*EmuRoundStrategy) GetBaselineWeightForNeighbours() uint32 {
	return rand.Uint32()
}

func (*EmuRoundStrategy) ShuffleNodeSequence(n int, swap func(i, j int)) {
	rand.Shuffle(n, swap)
}

func (*EmuRoundStrategy) AdjustConsensusTimings(timings *api.RoundTimings) {
}

var _ api.ConsensusControlFeeder = &EmuControlFeeder{}

type EmuControlFeeder struct {
	leaveReason uint32
}

func (p *EmuControlFeeder) CanFastForwardPulse(expected, received pulse.Number, lastPulseData pulse.Data) bool {
	panic("implement me")
}

func (p *EmuControlFeeder) OnPulseDetected() {
	panic("implement me")
}

func (p *EmuControlFeeder) CanStopOnHastyPulse(pn pulse.Number, expectedEndOfConsensus time.Time) bool {
	return false
}

func (p *EmuControlFeeder) OnAppliedMembershipProfile(mode member.OpMode, pw member.Power, effectiveSince pulse.Number) {
}

func (*EmuControlFeeder) SetTrafficLimit(level capacity.Level, duration time.Duration) {
}

func (*EmuControlFeeder) ResumeTraffic() {
}

func (*EmuControlFeeder) GetRequiredPowerLevel() power.Request {
	return power.NewRequestByLevel(capacity.LevelNormal)
}

func (p *EmuControlFeeder) GetRequiredGracefulLeave() (bool, uint32) {
	return p.leaveReason != 0, p.leaveReason
}

func (*EmuControlFeeder) OnAppliedGracefulLeave(exitCode uint32, effectiveSince pulse.Number) {
}

type EmuEphemeralFeeder struct{}

func (e EmuEphemeralFeeder) CanFastForwardPulse(expected, received pulse.Number, lastPulseData pulse.Data) bool {
	panic("implement me")
}

func (e EmuEphemeralFeeder) CanStopEphemeralByPulse(pd pulse.Data, localNode profiles.ActiveNode) bool {
	panic("implement me")
}

func (e EmuEphemeralFeeder) CanStopEphemeralByCensus(expected census.Expected) bool {
	panic("implement me")
}

func (e EmuEphemeralFeeder) GetMaxDuration() time.Duration {
	panic("implement me")
}

func (e EmuEphemeralFeeder) OnEphemeralCancelled() {
	panic("implement me")
}

func (e EmuEphemeralFeeder) GetMinDuration() time.Duration {
	return 2 * time.Second
}

func (e EmuEphemeralFeeder) OnNonEphemeralPacket(ctx context.Context, parser transport.PacketParser, inbound endpoints.Inbound) error {
	return nil
}

func (e EmuEphemeralFeeder) EphemeralConsensusFinished(isNextEphemeral bool, roundStartedAt time.Time, expected census.Operational) {
}

func (e EmuEphemeralFeeder) GetEphemeralTimings(c api.LocalNodeConfiguration) api.RoundTimings {
	return c.GetConsensusTimings(2)
}

func (e EmuEphemeralFeeder) IsActive() bool {
	return false
}

func (e EmuEphemeralFeeder) CreateEphemeralPulsePacket(census census.Operational) proofs.OriginalPulsarPacket {
	return nil
}

func (e EmuEphemeralFeeder) CanStopOnHastyPulse(pn pulse.Number, expectedEndOfConsensus time.Time) bool {
	return false
}
