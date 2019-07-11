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

package phases

import (
	"context"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"

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
	from endpoints.HostIdentityHolder) (*core.NodeAppearance, error) {

	jc := newJoinerController()
	return jc.handleUnknownJoinerPacket(ctx, reader, from, c.R)
}

func (c *PhaseControllerWithJoinersTemplate) createPerNodePacketHandler(ctx context.Context, ctlIndex int, node *core.NodeAppearance,
	realm *core.FullRealm, fn JoinerControllerPacketFunc) (core.PhasePerNodePacketFunc, context.Context) {

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

//type postponedPacket struct {
//	packet packets.PacketParser
//}

type JoinerControllerPacketFunc func(ctx context.Context, reader packets.MemberPacketReader, from *JoinerController) error

type JoinerController struct {
	node *core.NodeAppearance
	//realm *core.FullRealm

	handlerIndices []int //to cleanup when joiner is confirmed
	//postponedPackets []postponedPacket
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
	from endpoints.HostIdentityHolder, r *core.FullRealm) (*core.NodeAppearance, error) {

	//all packets will also be processed by main handlers

	switch reader.GetPacketType() {
	case gcp_types.PacketPhase1:
		p1 := reader.AsPhase1Packet()
		//r.GetProfileFactory().CreateBriefIntroProfile(intro, intro.GetJoinerSignature())
		//nip := r.profileFactory.CreateBriefIntroProfile(intro, intro.GetJoinerSignature())
		//if fIntro, ok := intro.(packets.FullIntroductionReader); ok && !fIntro.GetIssuerID().IsAbsent() {
		//	nip = r.profileFactory.CreateFullIntroProfile(nip, fIntro)
		//}
		//na := r.population.CreateNodeAppearance(r.roundContext, nip)
		//
		//cIntro := p1.GetCloudIntroduction()
		////checkJoinerSecret(cIntro.GetCloudIdentity(), cIntro.GetJoinerSecret())
		//cIntro.GetLastCloudStateHash()

		return p.applyBriefInfo(ctx, p1.GetFullIntroduction(), from, r)

	case gcp_types.PacketPhase2:
		p2 := reader.AsPhase2Packet()
		intro := p2.GetBriefIntroduction()
		r.GetProfileFactory().CreateBriefIntroProfile(intro)
		return p.applyBriefInfo(ctx, p2.GetBriefIntroduction(), from, r)
	}
	return nil, nil
}

func (p *JoinerController) applyBriefInfo(ctx context.Context, intro packets.BriefIntroductionReader,
	from endpoints.HostIdentityHolder, r *core.FullRealm) (*core.NodeAppearance, error) {

	ctx = context.WithValue(ctx, contextKeyValue, p)
	return r.CreatePurgatoryNode(ctx, intro, from)
}
