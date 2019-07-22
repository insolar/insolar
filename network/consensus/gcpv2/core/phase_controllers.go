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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetrecorder"
)

type PacketDispatcher interface {
	DispatchHostPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound, flags packetrecorder.PacketVerifyFlags) error
	DispatchMemberPacket(ctx context.Context, packet transport.MemberPacketReader, source *NodeAppearance) error
	DispatchUnknownMemberPacket(ctx context.Context, memberID insolar.ShortNodeID, packet transport.MemberPacketReader,
		from endpoints.Inbound) (bool, error)
	HasCustomVerifyForHost(from endpoints.Inbound, strict bool) bool
}

type MemberPacketSender interface {
	transport.TargetProfile
	SetPacketSent(pt phases.PacketType) bool
}
type MemberPacketReceiver interface {
	GetNodeID() insolar.ShortNodeID
	CanReceivePacket(pt phases.PacketType) bool
	VerifyPacketAuthenticity(packetSignature cryptkit.SignedDigest, from endpoints.Inbound, strictFrom bool) error
	SetPacketReceived(pt phases.PacketType) bool
	DispatchMemberPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound, flags packetrecorder.PacketVerifyFlags,
		pd PacketDispatcher) error
}

type PhasePerNodePacketFunc func(ctx context.Context, packet transport.MemberPacketReader, from *NodeAppearance, realm *FullRealm) error
type PerNodePacketDispatcherFactory interface {
	// PhasePerNodePacketFunc
	CreatePerNodePacketHandler(perNodeContext context.Context, node *NodeAppearance) (context.Context, PhasePerNodePacketFunc)
}

// type PrepPhasePacketHandler func(ctx context.Context, reader transport.PacketParser, from endpoints.Inbound) (postpone bool, err error)
type PrepPhaseController interface {
	GetPacketType() []phases.PacketType
	CreatePacketDispatcher(pt phases.PacketType, realm *PrepRealm) PacketDispatcher

	// HandleHostPacket(ctx context.Context, reader transport.PacketParser, from endpoints.Inbound) (postpone bool, err error)

	BeforeStart(realm *PrepRealm)
	StartWorker(ctx context.Context, realm *PrepRealm)
}

/* realm is provided for this handler to avoid being replicated in individual handlers */

type PhaseController interface {
	GetPacketType() []phases.PacketType
	CreatePacketDispatcher(pt phases.PacketType, ctlIndex int, realm *FullRealm) (PacketDispatcher, PerNodePacketDispatcherFactory)

	// HandleHostPacket(ctx context.Context, reader transport.PacketParser, from endpoints.Inbound) error                                   // GetHandlerType() == PacketHandlerTypeHost
	// HandleMemberPacket(ctx context.Context, reader transport.MemberPacketReader, src *NodeAppearance) error                              // GetHandlerType() == PacketHandlerTypeMember OR PacketHandlerTypeMemberFromUnknown
	// HandleUnknownMemberPacket(ctx context.Context, reader transport.MemberPacketReader, from endpoints.Inbound) (*NodeAppearance, error) // GetHandlerType() == PacketHandlerTypeMemberFromUnknown

	BeforeStart(realm *FullRealm)
	StartWorker(ctx context.Context, realm *FullRealm)
}

type PhaseControllersBundle interface {
	IsEphemeralPulseAllowed() bool
	IsDynamicPopulationRequired() bool
	CreatePrepPhaseControllers() []PrepPhaseController
	CreateFullPhaseControllers(nodeCount int) ([]PhaseController, NodeUpdateCallback)
}

type PhaseControllersBundleFactory interface {
	CreateControllersBundle(population census.OnlinePopulation, config api.LocalNodeConfiguration) PhaseControllersBundle
}

type NodeUpdateCallback interface {
	OnTrustUpdated(populationVersion uint32, n *NodeAppearance, before, after member.TrustLevel)
	OnNodeStateAssigned(populationVersion uint32, n *NodeAppearance)
	OnDynamicNodeAdded(populationVersion uint32, n *NodeAppearance, fullIntro bool)
	OnPurgatoryNodeAdded(populationVersion uint32, n *NodePhantom)
	OnCustomEvent(populationVersion uint32, n *NodeAppearance, event interface{})
	OnDynamicPopulationCompleted(populationVersion uint32, indexedCount int)
}
