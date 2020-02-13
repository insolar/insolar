// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package transport

import (
	"context"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
)

type Factory interface {
	GetPacketSender() PacketSender
	GetPacketBuilder(signer cryptkit.DigestSigner) PacketBuilder
	GetCryptographyFactory() CryptographyAssistant
}

type TargetProfile interface {
	GetNodeID() insolar.ShortNodeID
	GetStatic() profiles.StaticProfile
	IsJoiner() bool
	// GetOpMode() member.OpMode
	EncryptJoinerSecret(joinerSecret cryptkit.DigestHolder) cryptkit.DigestHolder
}

type ProfileFilter func(ctx context.Context, targetIndex int) (TargetProfile, PacketSendOptions)

type PreparedPacketSender interface {
	SendTo(ctx context.Context, target TargetProfile, sendOptions PacketSendOptions, sender PacketSender)

	/* Allows to control parallelism. Can return nil to skip a target */
	SendToMany(ctx context.Context, targetCount int, sender PacketSender, filter ProfileFilter)
}

type PacketBuilder interface {
	GetNeighbourhoodSize() NeighbourhoodSizes

	// PrepareIntro

	PreparePhase0Packet(sender *NodeAnnouncementProfile, pulsarPacket proofs.OriginalPulsarPacket,
		options PacketPrepareOptions) PreparedPacketSender
	PreparePhase1Packet(sender *NodeAnnouncementProfile, pulsarPacket proofs.OriginalPulsarPacket,
		welcome *proofs.NodeWelcomePackage, options PacketPrepareOptions) PreparedPacketSender

	/* Prepare receives all introductions at once, but PreparedSendPacket.SendTo MUST:
	1. exclude all intros when target is not joiner
	2. exclude the intro of the target
	*/
	PreparePhase2Packet(sender *NodeAnnouncementProfile, welcome *proofs.NodeWelcomePackage,
		neighbourhood []MembershipAnnouncementReader, options PacketPrepareOptions) PreparedPacketSender

	PreparePhase3Packet(sender *NodeAnnouncementProfile, vectors statevector.Vector,
		options PacketPrepareOptions) PreparedPacketSender
}

type PacketSender interface {
	SendPacketToTransport(ctx context.Context, t TargetProfile, sendOptions PacketSendOptions, payload interface{})
}

type PacketPrepareOptions uint32
type PacketSendOptions uint32

func (o PacketSendOptions) HasAny(mask PacketSendOptions) bool {
	return (o & mask) != 0
}

func (o PacketSendOptions) HasAll(mask PacketSendOptions) bool {
	return (o & mask) == mask
}

func (o PacketPrepareOptions) HasAny(mask PacketPrepareOptions) bool {
	return (o & mask) != 0
}

func (o PacketPrepareOptions) HasAll(mask PacketPrepareOptions) bool {
	return (o & mask) == mask
}

func (o PacketPrepareOptions) AsSendOptions() PacketSendOptions {
	return PacketSendOptions(o) & SharedPrepareSendOptionsMask
}

const (
	SendWithoutPulseData PacketSendOptions = 1 << iota
	TargetNeedsIntro
)

const SharedPrepareSendOptionsMask = SendWithoutPulseData | TargetNeedsIntro
const PrepareWithoutPulseData = PacketPrepareOptions(SendWithoutPulseData)
const PrepareWithIntro = PacketPrepareOptions(TargetNeedsIntro)

const (
	AlternativePhasePacket PacketPrepareOptions = 1 << (16 + iota)
	OnlyBriefIntroAboutJoiner
)
