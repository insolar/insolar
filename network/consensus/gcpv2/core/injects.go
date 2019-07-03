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

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/nodeset"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

type PhaseControllersBundle interface {
	GetPrepPhaseControllers() []PrepPhaseController
	GetFullPhaseControllers(nodeCount int) ([]PhaseController, NodeUpdateCallback)
}

type NodeUpdateCallback interface {
	OnTrustUpdated(n *NodeAppearance, before, after packets.NodeTrustLevel)
	OnNodeStateAssigned(n *NodeAppearance)
	OnCustomEvent(n *NodeAppearance, event interface{})
}

type ConsensusController interface {
	ProcessPacket(ctx context.Context, payload packets.PacketParser, from common.HostIdentityHolder) error

	/* Ungraceful stop */
	Abort()
	/* Graceful exit, actual moment of leave will be indicated via Upstream */
	//RequestLeave()

	/* This node power in the active population, and pulse number of such. Without active population returns (0,0) */
	GetActivePowerLimit() (common2.MemberPower, common.PulseNumber)
}

type CandidateControlFeeder interface {
	PickNextJoinCandidate() common2.CandidateProfile
	RemoveJoinCandidate(candidateAdded bool, nodeID common.ShortNodeID) bool
}

type ConsensusControlFeeder interface {
	//To Be used

	GetRequiredPowerLevel() common2.PowerRequest
	OnAppliedPowerLevel(pw common2.MemberPower, effectiveSince common.PulseNumber)

	GetRequiredGracefulLeave() (bool, uint32)
	OnAppliedGracefulLeave(exitCode uint32, effectiveSince common.PulseNumber)

	//OnAppliedPopulation()
}

type RoundController interface {
	HandlePacket(ctx context.Context, packet packets.PacketParser, from common.HostIdentityHolder) error
	StopConsensusRound()
	StartConsensusRound(upstream UpstreamPulseController)
}

type RoundControllerFactory interface {
	CreateConsensusRound(chronicle census.ConsensusChronicles, controlFeeder ConsensusControlFeeder,
		candidateFeeder CandidateControlFeeder, prevPulseRound RoundController) RoundController
	GetLocalConfiguration() LocalNodeConfiguration
}

type RoundStrategyFactory interface {
	CreateRoundStrategy(chronicle census.ConsensusChronicles, config LocalNodeConfiguration) RoundStrategy
}

type RoundStrategy interface {
	PhaseControllersBundle

	RandUint32() uint32
	ShuffleNodeSequence(n int, swap func(i, j int))
	IsEphemeralPulseAllowed() bool
	ConfigureRoundContext(ctx context.Context, expectedPulse common.PulseNumber, self common2.LocalNodeProfile) context.Context
	AdjustConsensusTimings(timings *common2.RoundTimings)
}

type LocalNodeConfiguration interface {
	GetConsensusTimings(nextPulseDelta uint16, isJoiner bool) common2.RoundTimings
	GetSecretKeyStore() common.SecretKeyStore
	GetParentContext() context.Context
}

type PacketSender interface {
	SendPacketToTransport(ctx context.Context, t common2.NodeProfile, sendOptions PacketSendOptions, payload interface{})
}

type PacketSendOptions uint32

const (
	SendWithoutPulseData PacketSendOptions = 1 << iota
	RequestForPhase1
)

type PreparedPacketSender interface {
	SendTo(ctx context.Context, target common2.NodeProfile, sendOptions PacketSendOptions, sender PacketSender)
}

//type PreparedIntro interface {}

type PacketBuilder interface {
	GetNeighbourhoodSize() common2.NeighbourhoodSizes

	//PrepareIntro

	PreparePhase0Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket common2.OriginalPulsarPacket,
		options PacketSendOptions) PreparedPacketSender
	PreparePhase1Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket common2.OriginalPulsarPacket,
		options PacketSendOptions) PreparedPacketSender

	/* Prepare receives all introductions at once, but PreparedSendPacket.SendTo MUST:
	1. exclude all intros when target is not joiner
	2. exclude the intro of the target
	*/
	PreparePhase2Packet(sender *packets.NodeAnnouncementProfile,
		neighbourhood []packets.MembershipAnnouncementReader, options PacketSendOptions) PreparedPacketSender

	PreparePhase3Packet(sender *packets.NodeAnnouncementProfile,
		bitset nodeset.NodeBitset, gshTrusted common2.GlobulaStateHash, gshDoubted common2.GlobulaStateHash,
		options PacketSendOptions) PreparedPacketSender
}

type TransportFactory interface {
	GetPacketSender() PacketSender
	GetPacketBuilder(signer common.DigestSigner) PacketBuilder
	GetCryptographyFactory() TransportCryptographyFactory
}

type TransportCryptographyFactory interface {
	common.SignatureVerifierFactory
	common.KeyStoreFactory
	GetDigestFactory() common.DigestFactory
	GetNodeSigner(sks common.SecretKeyStore) common.DigestSigner
}
