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

func (p *EmuRoundStrategyFactory) CreateRoundStrategy(pop census.OnlinePopulation,
	config api.LocalNodeConfiguration) (core.RoundStrategy, core.PhaseControllersBundle) {

	if p.bundleFactory == nil {
		p.bundleFactory = phasebundle.NewStandardBundleFactoryDefault()
	}

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
