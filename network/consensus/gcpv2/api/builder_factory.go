///
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
///

package api

import (
	"context"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

type PreparedPacketSender interface {
	SendTo(ctx context.Context, target gcp_types.NodeProfile, sendOptions PacketSendOptions, sender PacketSender)

	/* Allows to control parallelism. Can return nil to skip a target */
	SendToMany(ctx context.Context, targetCount int, sender PacketSender,
		filter func(ctx context.Context, targetIndex int) (gcp_types.NodeProfile, PacketSendOptions))
}

type PacketBuilder interface {
	GetNeighbourhoodSize() gcp_types.NeighbourhoodSizes

	//PrepareIntro

	PreparePhase0Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket packets.OriginalPulsarPacket,
		options PacketSendOptions) PreparedPacketSender
	PreparePhase1Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket packets.OriginalPulsarPacket,
		options PacketSendOptions) PreparedPacketSender

	/* Prepare receives all introductions at once, but PreparedSendPacket.SendTo MUST:
	1. exclude all intros when target is not joiner
	2. exclude the intro of the target
	*/
	PreparePhase2Packet(sender *packets.NodeAnnouncementProfile,
		neighbourhood []packets.MembershipAnnouncementReader, options PacketSendOptions) PreparedPacketSender

	PreparePhase3Packet(sender *packets.NodeAnnouncementProfile, vectors gcp_types.HashedNodeVector,
		options PacketSendOptions) PreparedPacketSender
}

type TransportFactory interface {
	GetPacketSender() PacketSender
	GetPacketBuilder(signer cryptography_containers.DigestSigner) PacketBuilder
	GetCryptographyFactory() TransportCryptographyFactory
}

type TransportCryptographyFactory interface {
	cryptography_containers.SignatureVerifierFactory
	cryptography_containers.KeyStoreFactory
	GetDigestFactory() cryptography_containers.DigestFactory
	GetNodeSigner(sks cryptography_containers.SecretKeyStore) cryptography_containers.DigestSigner
}
