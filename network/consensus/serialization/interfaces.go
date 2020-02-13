// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package serialization

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
)

var (
	defaultByteOrder = binary.BigEndian
)

type PacketHeaderAccessor interface {
	GetProtocolType() ProtocolType
	GetPacketType() phases.PacketType
	GetSourceID() insolar.ShortNodeID
	HasFlag(flag Flag) bool
	GetFlagRangeInt(from, to uint8) uint8
	IsRelayRestricted() bool
	IsBodyEncrypted() bool
}

type PacketHeaderModifier interface {
	SetFlag(flag Flag)
	ClearFlag(flag Flag)
}

type PacketBody interface {
	ContextSerializerTo
	ContextDeserializerFrom

	String(ctx PacketContext) string
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
	GetNeighbourNodeID() insolar.ShortNodeID
	SetNeighbourNodeID(nodeID insolar.ShortNodeID)
	GetAnnouncedJoinerNodeID() insolar.ShortNodeID
	SetAnnouncedJoinerNodeID(nodeID insolar.ShortNodeID)
}

type SerializeContext interface {
	PacketHeaderModifier
	PacketContext
}

type DeserializeContext interface {
	PacketContext
}

type SerializerTo interface {
	SerializeTo(ctx context.Context, writer io.Writer, digester cryptkit.DataDigester, signer cryptkit.DigestSigner) (int64, error)
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
