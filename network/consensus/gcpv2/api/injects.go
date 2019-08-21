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

package api

import (
	"context"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/pulse"
)

type ConsensusController interface {
	Prepare()
	ProcessPacket(ctx context.Context, payload transport.PacketParser, from endpoints.Inbound) error

	/* Ungraceful stop */
	Abort()
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api.CandidateControlFeeder -o . -s _mock.go -g
type CandidateControlFeeder interface {
	PickNextJoinCandidate() (profiles.CandidateProfile /* joinerSecret */, cryptkit.DigestHolder)
	RemoveJoinCandidate(candidateAdded bool, nodeID insolar.ShortNodeID) bool
}

type TrafficControlFeeder interface {
	/* Application traffic should be stopped or throttled down for the given duration
	LevelMax and LevelNormal should be considered equal, and duration doesnt apply to them
	*/
	SetTrafficLimit(level capacity.Level, duration time.Duration)

	/* Application traffic can be resumed at full */
	ResumeTraffic()
}

type EphemeralMode uint8

const (
	EphemeralNotAllowed EphemeralMode = iota
	EphemeralAllowed                  // can generate ephemeral pulses
)

func (mode EphemeralMode) IsEnabled() bool {
	return mode != EphemeralNotAllowed
}

type PulseControlFeeder interface {
	CanStopOnHastyPulse(pn pulse.Number, expectedEndOfConsensus time.Time) bool
	CanFastForwardPulse(expected, received pulse.Number, lastPulseData pulse.Data) bool
}

type EphemeralControlFeeder interface {
	PulseControlFeeder
	GetEphemeralTimings(LocalNodeConfiguration) RoundTimings
	/* Minimum time after the last ephemeral round before checking for another candidate */
	GetMinDuration() time.Duration
	/* Maximum time to wait for a candidate before starting a next ephemeral round */
	GetMaxDuration() time.Duration

	/* if true, then a new round can be triggered by a joiner candidate */
	IsActive() bool
	CreateEphemeralPulsePacket(census census.Operational) proofs.OriginalPulsarPacket

	OnNonEphemeralPacket(ctx context.Context, parser transport.PacketParser, inbound endpoints.Inbound) error

	/* Applied when an ephemeral node gets a non-ephemeral pulse data from another member */
	CanStopEphemeralByPulse(pd pulse.Data, localNode profiles.ActiveNode) bool
	/* Applied when an ephemeral node finishes consensus */
	CanStopEphemeralByCensus(expected census.Expected) bool

	EphemeralConsensusFinished(isNextEphemeral bool, roundStartedAt time.Time, expected census.Operational)
	/* is called:
		(1) immediately after TryConvertFromEphemeral returned true
	    (2) at start of full realm, when ephemeral mode was cancelled by Phase0/Phase1 packets
	*/
	OnEphemeralCancelled()
}

type ConsensusControlFeeder interface {
	TrafficControlFeeder
	PulseControlFeeder

	GetRequiredPowerLevel() power.Request
	OnAppliedMembershipProfile(mode member.OpMode, pw member.Power, effectiveSince pulse.Number)

	GetRequiredGracefulLeave() (bool, uint32)
	OnAppliedGracefulLeave(exitCode uint32, effectiveSince pulse.Number)

	OnPulseDetected() // this method is not currently invoked
}

type MaintenancePollFunc func(ctx context.Context) bool

type RoundStateCallback interface {
	UpstreamController

	/* Called on receiving seem-to-be-valid Pulsar or Phase0 packets. Can be called multiple time in sequence.
	Application MUST NOT consider it as a new pulse. */
	OnPulseDetected()

	OnFullRoundStarting()

	// A special case for a stateless, as it doesnt request NSG with PreparePulseChange
	CommitPulseChangeByStateless(report UpstreamReport, pd pulse.Data, activeCensus census.Operational)

	/* Called by the longest phase worker on termination */
	OnRoundStopped(ctx context.Context)
}

type RoundControlCode uint8

const (
	KeepRound RoundControlCode = iota
	StartNextRound
	//	NextRoundPrepare
	NextRoundTerminate
)

type RoundController interface {
	PrepareConsensusRound(upstream UpstreamController)
	StopConsensusRound()
	HandlePacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound) (RoundControlCode, error)
}

type RoundControllerFactory interface {
	CreateConsensusRound(chronicle ConsensusChronicles, controlFeeder ConsensusControlFeeder, candidateFeeder CandidateControlFeeder,
		ephemeralFeeder EphemeralControlFeeder, prevPulseRound RoundController) RoundController
	GetLocalConfiguration() LocalNodeConfiguration
}

type LocalNodeConfiguration interface {
	GetConsensusTimings(nextPulseDelta uint16) RoundTimings
	GetSecretKeyStore() cryptkit.SecretKeyStore
	GetParentContext() context.Context
	GetNodeCountHint() int
}
