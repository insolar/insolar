package serialization

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

var (
	defaultByteOrder = binary.BigEndian
)

type PacketHeaderAccessor interface {
	GetProtocolType() ProtocolType
	GetPacketType() packets.PacketType
	GetSourceID() common.ShortNodeID
	HasFlag(flag Flag) bool
	IsRelayRestricted() bool
	IsBodyEncrypted() bool
}

type PacketHeaderModifier interface {
	SetFlag(flag Flag)
}

type PacketBody interface {
	ContextSerializerTo
	ContextDeserializerFrom
}

type FieldContext uint

const (
	NoContext = FieldContext(iota)
	ContextMembershipAnnouncement
	ContextNeighbourAnnouncement
)

type PacketContext interface {
	PacketHeaderAccessor
	context.Context

	InContext(ctx FieldContext) bool
	SetInContext(ctx FieldContext)
	GetNeighbourNodeID() common.ShortNodeID
	SetNeighbourNodeID(nodeID common.ShortNodeID)
	GetAnnouncedJoinerNodeID() common.ShortNodeID
	SetAnnouncedJoinerNodeID(nodeID common.ShortNodeID)
}

type SerializeContext interface {
	PacketHeaderModifier
	PacketContext
}

type DeserializeContext interface {
	PacketContext
}

type SerializerTo interface {
	SerializeTo(ctx context.Context, writer io.Writer, signer common.DataSigner) (int64, error)
}

type ContextSerializerTo interface {
	SerializeTo(ctx SerializeContext, writer io.Writer) error
}

type DeserializerFrom interface {
	DeserializeFrom(ctx context.Context, reader io.Reader) (int64, error)
}

type ContextDeserializerFrom interface {
	DeserializeFrom(ctx DeserializeContext, reader io.Reader) error
}
