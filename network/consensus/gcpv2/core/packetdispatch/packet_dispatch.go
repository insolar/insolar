package packetdispatch

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"
)

type MemberPacketReceiver interface {
	GetNodeID() insolar.ShortNodeID
	CanReceivePacket(pt phases.PacketType) bool
	VerifyPacketAuthenticity(packetSignature cryptkit.SignedDigest, from endpoints.Inbound, strictFrom bool) error
	SetPacketReceived(pt phases.PacketType) bool
	DispatchMemberPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound, flags coreapi.PacketVerifyFlags,
		pd population.PacketDispatcher) error
}
