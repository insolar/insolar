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

type PacketContext interface {
	PacketHeaderAccessor
	context.Context
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
