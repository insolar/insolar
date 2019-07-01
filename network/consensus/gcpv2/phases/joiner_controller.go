package phases

import (
	"context"
	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

/*
	Provides handling of packets _from_ joiners.
*/
type PhaseControllerWithJoinersTemplate struct {
	core.PhaseControllerPerMemberFromUnknownTemplate
}

func (*PhaseControllerWithJoinersTemplate) GetHandlerType() core.PhaseControllerHandlerType {
	return core.PacketHandlerTypeEnableUnknown | core.PacketHandlerTypeMember | core.PacketHandlerTypePerNode
}

func (c *PhaseControllerWithJoinersTemplate) HandleUnknownMemberPacket(ctx context.Context, reader packets.MemberPacketReader,
	from common.HostIdentityHolder) (*core.NodeAppearance, error) {

	jc := newJoinerController()
	return jc.handleUnknownJoinerPacket(ctx, reader, from, c.R)
}

func (c *PhaseControllerWithJoinersTemplate) createPerNodePacketHandler(ctlIndex int, node *core.NodeAppearance,
	realm *core.FullRealm, ctx context.Context, fn JoinerControllerPacketFunc) (core.PhasePerNodePacketFunc, context.Context) {

	if !node.IsJoiner() {
		return nil, ctx
	}

	if jc, ok := ctx.Value(contextKeyValue).(*JoinerController); ok {
		jc.EnsureEnvironment(node, realm)
		return jc.getJoinerPacketHandler(ctlIndex, fn), ctx
	}

	return nil, ctx
}

type contextKeyType struct{}

var contextKeyValue = contextKeyType{}

type postponedPacket struct {
	packet packets.PacketParser
}

type JoinerControllerPacketFunc func(ctx context.Context, reader packets.MemberPacketReader, from *JoinerController) error

type JoinerController struct {
	node *core.NodeAppearance
	//realm *core.FullRealm

	handlerIndices   []int //to cleanup when joiner is confirmed
	postponedPackets []postponedPacket
}

func newJoinerController() *JoinerController {
	return &JoinerController{}
}

func (p *JoinerController) EnsureEnvironment(n *core.NodeAppearance, r *core.FullRealm) {
	if p.node != n /* || p.realm != r */ {
		panic("illegal value")
	}
}

func (p *JoinerController) getJoinerPacketHandler(ctlIndex int, fn JoinerControllerPacketFunc) core.PhasePerNodePacketFunc {
	p.handlerIndices = append(p.handlerIndices, ctlIndex)
	return p.createPacketHandler(fn)
}

func (p *JoinerController) createPacketHandler(fn JoinerControllerPacketFunc) core.PhasePerNodePacketFunc {
	return func(ctx context.Context, reader packets.MemberPacketReader, from *core.NodeAppearance, realm *core.FullRealm) error {
		p.EnsureEnvironment(from, realm)

		//if err == nil {
		//	p.addPostponedPacket(reader)
		//}
		//return na, err
		//
		return fn(ctx, reader, p)
	}
}

func (p *JoinerController) handleUnknownJoinerPacket(ctx context.Context, reader packets.MemberPacketReader,
	from common.HostIdentityHolder, r *core.FullRealm) (*core.NodeAppearance, error) {

	switch reader.GetPacketType() {

	case packets.PacketPhase2:
		p2 := reader.AsPhase2Packet()
		return p.createPurgatoryNode(ctx, p2.GetBriefIntroduction(), from, r)

	case packets.PacketPhase1:
		p1 := reader.AsPhase1Packet()
		intro := p1.GetFullIntroduction()
		na, err := p.createPurgatoryNode(ctx, intro, from, r)
		if err != nil {
			return nil, err
		}
		err = p.convertPurgatoryToDynamicNode(intro, p1.GetCloudIntroduction(), r)
		return na, err
	}
	return nil, nil
}

func (p *JoinerController) createPurgatoryNode(ctx context.Context, intro packets.BriefIntroductionReader,
	from common.HostIdentityHolder, r *core.FullRealm) (*core.NodeAppearance, error) {

	ctx = context.WithValue(ctx, contextKeyValue, p)

	panic("not implemented")

	//intro.
	//
	//if err != nil {
	//	return nil, err
	//}
	//p.node = na //createPurgatoryNode guarantees that the node will not be duplicated, so there is no racing here without a lock
	//return na, nil

}

func (p *JoinerController) convertPurgatoryToDynamicNode(intro packets.FullIntroductionReader,
	cloudIntro packets.CloudIntroductionReader, realm *core.FullRealm) error {

	panic("not implemented")
}
