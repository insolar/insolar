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

package adapters

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/pulse"
)

const defaultEphemeralPulseDuration = 2 * time.Second

type EphemeralController interface {
	EphemeralMode(nodes []insolar.NetworkNode) bool
}

var _ api.ConsensusControlFeeder = &ConsensusControlFeeder{}

type ConsensusControlFeeder struct {
	mu            *sync.RWMutex
	capacityLevel capacity.Level
	leave         bool
	leaveReason   uint32
}

func NewConsensusControlFeeder() *ConsensusControlFeeder {
	return &ConsensusControlFeeder{
		mu:            &sync.RWMutex{},
		capacityLevel: capacity.LevelNormal,
	}
}

func (cf *ConsensusControlFeeder) SetRequiredGracefulLeave(reason uint32) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	cf.leave = true
	cf.leaveReason = reason
}

func (cf *ConsensusControlFeeder) SetRequiredPowerLevel(level capacity.Level) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	cf.capacityLevel = level
}

func (cf *ConsensusControlFeeder) GetRequiredGracefulLeave() (bool, uint32) {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	return cf.leave, cf.leaveReason
}

func (cf *ConsensusControlFeeder) GetRequiredPowerLevel() power.Request {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	return power.NewRequestByLevel(capacity.LevelNormal)
}

func (cf *ConsensusControlFeeder) CanFastForwardPulse(expected, received pulse.Number, lastPulseData pulse.Data) bool {
	return true
}

func (cf *ConsensusControlFeeder) CanStopOnHastyPulse(pn pulse.Number, expectedEndOfConsensus time.Time) bool {
	return true
}

func (cf *ConsensusControlFeeder) OnPulseDetected() {
}

func (cf *ConsensusControlFeeder) OnAppliedMembershipProfile(mode member.OpMode, pw member.Power, effectiveSince pulse.Number) {
}

func (cf *ConsensusControlFeeder) OnAppliedGracefulLeave(exitCode uint32, effectiveSince pulse.Number) {
}

func (cf *ConsensusControlFeeder) SetTrafficLimit(level capacity.Level, duration time.Duration) {
}

func (cf *ConsensusControlFeeder) ResumeTraffic() {
}

func InterceptConsensusControl(originalFeeder *ConsensusControlFeeder) *ControlFeederInterceptor {
	r := ControlFeederInterceptor{}
	r.internal.ConsensusControlFeeder = originalFeeder
	r.internal.mu = &sync.Mutex{}
	return &r
}

type ControlFeederInterceptor struct {
	internal InternalControlFeederAdapter
}

func (i *ControlFeederInterceptor) Feeder() *InternalControlFeederAdapter {
	return &i.internal
}

func (i *ControlFeederInterceptor) PrepareLeave() <-chan struct{} {
	i.internal.mu.Lock()
	defer i.internal.mu.Unlock()

	if i.internal.zeroReadyChannel != nil {
		panic("illegal state")
	}
	i.internal.zeroReadyChannel = make(chan struct{})
	if i.internal.hasZero || i.internal.zeroPending {
		i.internal.hasZero = true
		close(i.internal.zeroReadyChannel)
	}
	return i.internal.zeroReadyChannel
}

func (i *ControlFeederInterceptor) Leave(leaveReason uint32) <-chan struct{} {
	i.internal.mu.Lock()
	defer i.internal.mu.Unlock()

	if i.internal.leavingChannel != nil {
		panic("illegal state")
	}
	i.internal.leaveReason = leaveReason
	i.internal.leavingChannel = make(chan struct{})
	if i.internal.hasLeft {
		i.internal.setHasZero()
		close(i.internal.leavingChannel)
	}
	return i.internal.leavingChannel
}

var _ api.ConsensusControlFeeder = &InternalControlFeederAdapter{}

type InternalControlFeederAdapter struct {
	*ConsensusControlFeeder

	mu *sync.Mutex

	hasLeft bool
	hasZero bool

	zeroPending bool

	leaveReason      uint32
	zeroReadyChannel chan struct{}
	leavingChannel   chan struct{}
}

func (cf *InternalControlFeederAdapter) GetRequiredPowerLevel() power.Request {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	if cf.zeroReadyChannel != nil || cf.leavingChannel != nil {
		return power.NewRequestByLevel(capacity.LevelZero)
	}
	return cf.ConsensusControlFeeder.GetRequiredPowerLevel()
}

func (cf *InternalControlFeederAdapter) OnAppliedMembershipProfile(mode member.OpMode, pw member.Power, effectiveSince pulse.Number) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	cf.zeroPending = pw == 0
	if pw == 0 && cf.zeroReadyChannel != nil {
		cf.setHasZero()
	}

	if mode.IsEvicted() {
		cf.setHasLeft()
	}

	cf.ConsensusControlFeeder.OnAppliedMembershipProfile(mode, pw, effectiveSince)
}

func (cf *InternalControlFeederAdapter) GetRequiredGracefulLeave() (bool, uint32) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	if cf.leavingChannel != nil {
		return true, cf.leaveReason
	}
	return cf.ConsensusControlFeeder.GetRequiredGracefulLeave()
}

func (cf *InternalControlFeederAdapter) OnPulseDetected() {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	if cf.zeroPending {
		cf.setHasZero()
	}
	cf.ConsensusControlFeeder.OnPulseDetected()
}

func (cf *InternalControlFeederAdapter) setHasZero() {
	if !cf.hasZero && cf.zeroReadyChannel != nil {
		close(cf.zeroReadyChannel)
	}
	cf.hasZero = true
}

func (cf *InternalControlFeederAdapter) setHasLeft() {
	cf.setHasZero()

	if !cf.hasLeft && cf.leavingChannel != nil {
		close(cf.leavingChannel)
	}
	cf.hasLeft = true
}

func NewEphemeralControlFeeder(ephemeralController EphemeralController) *EphemeralControlFeeder {
	return &EphemeralControlFeeder{
		ephemeralController: ephemeralController,
		pulseDuration:       defaultEphemeralPulseDuration,
	}
}

type EphemeralControlFeeder struct {
	pulseChanger        PulseChanger
	ephemeralController EphemeralController
	pulseDuration       time.Duration
}

func (f *EphemeralControlFeeder) CanFastForwardPulse(expected, received pulse.Number, lastPulseData pulse.Data) bool {
	return true
}

func (f *EphemeralControlFeeder) CanStopEphemeralByPulse(pd pulse.Data, localNode profiles.ActiveNode) bool {
	return true
}

func (f *EphemeralControlFeeder) OnEphemeralCancelled() {
	// TODO is called on cancellation by both Ph1 packets and TryConvertFromEphemeral
}

func (f *EphemeralControlFeeder) GetMinDuration() time.Duration {
	return f.pulseDuration
}

func (f *EphemeralControlFeeder) GetMaxDuration() time.Duration {
	return math.MaxInt64
}

func (f *EphemeralControlFeeder) OnNonEphemeralPacket(ctx context.Context, parser transport.PacketParser, inbound endpoints.Inbound) error {
	_, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"sender_address": inbound.GetNameAddress().String(),
		"sender_id":      parser.GetSourceID(),
		"packet_type":    parser.GetPacketType().String(),
		"packet_pulse":   parser.GetPulseNumber(),
	})

	logger.Info("non-ephemeral packet")
	return nil
}

func (f *EphemeralControlFeeder) CanStopEphemeralByCensus(expected census.Expected) bool {
	if expected == nil {
		return false
	}

	population := expected.GetOnlinePopulation()
	if !population.IsValid() {
		return false
	}

	networkNodes := NewNetworkNodeList(population.GetProfiles())
	if f.ephemeralController.EphemeralMode(networkNodes) {
		return false
	}

	return true
}

func (f *EphemeralControlFeeder) EphemeralConsensusFinished(isNextEphemeral bool, roundStartedAt time.Time, expected census.Operational) {
}

func (f *EphemeralControlFeeder) GetEphemeralTimings(config api.LocalNodeConfiguration) api.RoundTimings {
	delta := f.pulseDuration / time.Second
	if delta < 2 {
		delta = 2
	}

	return config.GetConsensusTimings(uint16(delta))
}

func (f *EphemeralControlFeeder) IsActive() bool {
	return true
}

func (f *EphemeralControlFeeder) CreateEphemeralPulsePacket(census census.Operational) proofs.OriginalPulsarPacket {
	_, pd := census.GetNearestPulseData()
	if pd.IsEmpty() {
		pd = pulse.NewFirstEphemeralData()
	}
	pd = pd.CreateNextEphemeralPulse()

	return NewPulsePacketParser(pd)
}

func (f *EphemeralControlFeeder) CanStopOnHastyPulse(pn pulse.Number, expectedEndOfConsensus time.Time) bool {
	return false
}
