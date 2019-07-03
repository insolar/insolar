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

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"

	"github.com/insolar/insolar/network/consensus/common"
)

type PrepPhaseController interface {
	GetPacketType() packets.PacketType
	HandleHostPacket(ctx context.Context, reader packets.PacketParser, from common.HostIdentityHolder) (postpone bool, err error)
	BeforeStart(realm *PrepRealm)
	StartWorker(ctx context.Context)
}

type PrepPhasePacketHandler func(ctx context.Context, reader packets.PacketParser, from common.HostIdentityHolder) (postpone bool, err error)

type PhaseControllerHandlerType uint8

const (
	HandlerTypeHostPacket PhaseControllerHandlerType = iota
	HandlerTypeMemberPacket
	HandlerTypePerNodePacket // This mode allows to attach a custom object(state) to each NodeAppearance
)

type PhaseController interface {
	GetPacketType() packets.PacketType
	IsPerNode() PhaseControllerHandlerType
	HandleHostPacket(ctx context.Context, reader packets.PacketParser, from common.HostIdentityHolder) error // IsPerNode() == HandlerTypeHostPacket
	HandleMemberPacket(ctx context.Context, reader packets.MemberPacketReader, src *NodeAppearance) error    // IsPerNode() == HandlerTypeMemberPacket
	CreatePerNodePacketHandler(node *NodeAppearance) PhasePerNodePacketHandler                               // IsPerNode() == HandlerTypePerNodePacket
	BeforeStart(realm *FullRealm)
	StartWorker(ctx context.Context)
}

type PhaseHostPacketHandler func(ctx context.Context, reader packets.PacketParser, from common.HostIdentityHolder) error
type PhaseNodePacketHandler func(ctx context.Context, reader packets.MemberPacketReader, from *NodeAppearance) error

/* realm is provided for this handler to avoid being replicated in individual handlers */
type PhasePerNodePacketHandler func(ctx context.Context, reader packets.MemberPacketReader, from *NodeAppearance, realm *FullRealm) error

type PhaseControllerPerMemberTemplate struct {
	R *FullRealm
}

func (c *PhaseControllerPerMemberTemplate) BeforeStart(realm *FullRealm) {
	c.R = realm
}

func (*PhaseControllerPerMemberTemplate) IsPerNode() PhaseControllerHandlerType {
	return HandlerTypeMemberPacket
}

func (*PhaseControllerPerMemberTemplate) HandleHostPacket(ctx context.Context, reader packets.PacketParser, from common.HostIdentityHolder) error {
	panic("illegal call")
}

func (*PhaseControllerPerMemberTemplate) CreatePerNodePacketHandler(node *NodeAppearance) PhasePerNodePacketHandler {
	panic("illegal call")
}

func (*PhaseControllerPerMemberTemplate) StartWorker(ctx context.Context) {
}

// var _ PhaseController = &PhaseControllerPerNodeTemplate{}
type PhaseControllerPerNodeTemplate struct {
	R *FullRealm
}

func (c *PhaseControllerPerNodeTemplate) BeforeStart(realm *FullRealm) {
	c.R = realm
}

func (*PhaseControllerPerNodeTemplate) IsPerNode() PhaseControllerHandlerType {
	return HandlerTypePerNodePacket
}

func (*PhaseControllerPerNodeTemplate) HandleHostPacket(ctx context.Context, reader packets.PacketParser, from common.HostIdentityHolder) error {
	panic("illegal call")
}

func (*PhaseControllerPerNodeTemplate) HandleMemberPacket(ctx context.Context, reader packets.MemberPacketReader, src *NodeAppearance) error {
	panic("illegal call")
}

func (*PhaseControllerPerNodeTemplate) StartWorker(ctx context.Context) {
}

type PhaseControllerPerHostTemplate struct {
	R *FullRealm
}

func (c *PhaseControllerPerHostTemplate) BeforeStart(realm *FullRealm) {
	c.R = realm
}

func (*PhaseControllerPerHostTemplate) HandleMemberPacket(ctx context.Context, reader packets.MemberPacketReader, src *NodeAppearance) error {
	panic("illegal call")
}

func (*PhaseControllerPerHostTemplate) CreatePerNodePacketHandler(node *NodeAppearance) PhasePerNodePacketHandler {
	panic("illegal call")
}

func (*PhaseControllerPerHostTemplate) IsPerNode() PhaseControllerHandlerType {
	return HandlerTypeHostPacket
}

func (*PhaseControllerPerHostTemplate) StartWorker(ctx context.Context) {
}
