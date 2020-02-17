// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package core

import (
	"context"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

type PhaseControllerTemplate struct {
}

func (c *PhaseControllerTemplate) BeforeStart(ctx context.Context, realm *FullRealm) {
}

func (*PhaseControllerTemplate) StartWorker(ctx context.Context, realm *FullRealm) {
}

type PrepPhaseControllerTemplate struct {
}

func (c *PrepPhaseControllerTemplate) BeforeStart(ctx context.Context, realm *PrepRealm) {
}

func (*PrepPhaseControllerTemplate) StartWorker(ctx context.Context, realm *PrepRealm) {
}

// var _ PacketDispatcher = &HostPacketDispatcherTemplate{}

type HostPacketDispatcherTemplate struct {
}

func (*HostPacketDispatcherTemplate) TriggerUnknownMember(ctx context.Context, memberID insolar.ShortNodeID,
	packet transport.MemberPacketReader, from endpoints.Inbound) (bool, error) {
	return false, nil
}

func (*HostPacketDispatcherTemplate) HasCustomVerifyForHost(from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) bool {
	return false
}

func (*HostPacketDispatcherTemplate) DispatchMemberPacket(ctx context.Context, packet transport.MemberPacketReader,
	source *population.NodeAppearance) error {
	panic("illegal state")
}

// var _ PacketDispatcher = &MemberPacketDispatcherTemplate{}

type MemberPacketDispatcherTemplate struct {
}

func (*MemberPacketDispatcherTemplate) TriggerUnknownMember(ctx context.Context, memberID insolar.ShortNodeID,
	packet transport.MemberPacketReader, from endpoints.Inbound) (bool, error) {
	return false, nil
}

func (*MemberPacketDispatcherTemplate) HasCustomVerifyForHost(from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) bool {
	return false
}

func (*MemberPacketDispatcherTemplate) DispatchHostPacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, flags coreapi.PacketVerifyFlags) error {
	panic("illegal state")
}
