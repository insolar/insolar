package core

import (
	"context"
	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

type PhaseControllerPerMemberFromUnknownTemplate struct {
	R *FullRealm
}

func (c *PhaseControllerPerMemberFromUnknownTemplate) BeforeStart(realm *FullRealm) {
	c.R = realm
}

func (*PhaseControllerPerMemberFromUnknownTemplate) GetHandlerType() PhaseControllerHandlerType {
	return PacketHandlerTypeEnableUnknown | PacketHandlerTypeMember
}

func (*PhaseControllerPerMemberFromUnknownTemplate) HandleHostPacket(ctx context.Context, reader packets.PacketParser,
	from common.HostIdentityHolder) error {

	return errPacketIsNotAllowed
}

func (*PhaseControllerPerMemberFromUnknownTemplate) CreatePerNodePacketHandler(ctlIndex int, node *NodeAppearance,
	realm *FullRealm, sharedNodeContext context.Context) (PhasePerNodePacketFunc, context.Context) {

	return nil, sharedNodeContext
}

func (*PhaseControllerPerMemberFromUnknownTemplate) StartWorker(ctx context.Context) {
}

type PhaseControllerPerMemberTemplate struct {
	PhaseControllerPerMemberFromUnknownTemplate
}

func (*PhaseControllerPerMemberTemplate) GetHandlerType() PhaseControllerHandlerType {
	return PacketHandlerTypeMember
}

func (*PhaseControllerPerMemberTemplate) HandleUnknownMemberPacket(ctx context.Context,
	reader packets.MemberPacketReader, from common.HostIdentityHolder) (*NodeAppearance, error) {

	return nil, errPacketIsNotAllowed
}

// var _ PhaseController = &PhaseControllerPerNodeTemplate{}
type PhaseControllerPerNodeTemplate struct {
	R *FullRealm
}

func (c *PhaseControllerPerNodeTemplate) BeforeStart(realm *FullRealm) {
	c.R = realm
}

func (*PhaseControllerPerNodeTemplate) GetHandlerType() PhaseControllerHandlerType {
	return PacketHandlerTypePerNode
}

func (*PhaseControllerPerNodeTemplate) HandleHostPacket(ctx context.Context, reader packets.PacketParser, from common.HostIdentityHolder) error {
	return errPacketIsNotAllowed
}

func (*PhaseControllerPerNodeTemplate) HandleMemberPacket(ctx context.Context, reader packets.MemberPacketReader, src *NodeAppearance) error {
	return errPacketIsNotAllowed
}

func (*PhaseControllerPerNodeTemplate) HandleUnknownMemberPacket(ctx context.Context, reader packets.MemberPacketReader, from common.HostIdentityHolder) (*NodeAppearance, error) {
	return nil, errPacketIsNotAllowed
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
	return errPacketIsNotAllowed
}

func (*PhaseControllerPerHostTemplate) CreatePerNodePacketHandler(ctlIndex int, node *NodeAppearance,
	realm *FullRealm, sharedNodeContext context.Context) (PhasePerNodePacketFunc, context.Context) {

	return nil, sharedNodeContext
}

func (*PhaseControllerPerHostTemplate) HandleUnknownMemberPacket(ctx context.Context, reader packets.MemberPacketReader, from common.HostIdentityHolder) (*NodeAppearance, error) {
	return nil, errPacketIsNotAllowed
}

func (*PhaseControllerPerHostTemplate) GetHandlerType() PhaseControllerHandlerType {
	return PacketHandlerTypeHost
}

func (*PhaseControllerPerHostTemplate) StartWorker(ctx context.Context) {
}
