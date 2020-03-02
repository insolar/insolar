// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package population

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
)

type PacketDispatcher interface {
	HasCustomVerifyForHost(from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) bool

	DispatchHostPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound, flags coreapi.PacketVerifyFlags) error

	/* This method can validate and create a member, but MUST NOT apply any changes to members etc */
	TriggerUnknownMember(ctx context.Context, memberID insolar.ShortNodeID, packet transport.MemberPacketReader, from endpoints.Inbound) (bool, error)
	DispatchMemberPacket(ctx context.Context, packet transport.MemberPacketReader, source *NodeAppearance) error
}

type DispatchMemberPacketFunc func(ctx context.Context, packet transport.MemberPacketReader, from *NodeAppearance) error
