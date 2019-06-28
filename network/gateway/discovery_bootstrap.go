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

package gateway

import (
	"context"
	"errors"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

func newDiscoveryBootstrap(b *Base) *DiscoveryBootstrap {
	return &DiscoveryBootstrap{b}
}

// DiscoveryBootstrap void network state
type DiscoveryBootstrap struct {
	*Base
}

func (g *DiscoveryBootstrap) Run(ctx context.Context) {

	authorizeRes, err := g.authorize(ctx)
	if err != nil {
		// log warn
		g.Gatewayer.SwitchState(insolar.NoNetworkState)
		return
	}

	// TODO: check authorize result and switch to JoinerBootstrap if other network is complete
	//if err == nil && !insolar.IsEphemeralPulse(&p) {
	//	g.Gatewayer.SwitchState(insolar.JoinerBootstrap)
	//	return
	//}

	g.NodeKeeper.GetConsensusInfo().SetIsJoiner(false)

	_, err = g.PulseAccessor.Latest(ctx)
	pp := pulse.FromProto(authorizeRes.Pulse)
	if err != nil {
		g.PulseAppender.Append(ctx, *pp)
	}

	resp, err := g.BootstrapRequester.Bootstrap(ctx, authorizeRes.Permit, g.joinClaim, pp)
	if err != nil {

	}

	//  ConsensusWaiting, ETA
	g.bootstrapETA = insolar.PulseNumber(resp.ETA)
	g.Gatewayer.SwitchState(insolar.WaitConsensus)
	return

	// Authorize(utc) permit, check version
	// process response: trueAccept, redirect with permit, posibleAccept(regen shortId, updateScedule, update time utc)
	// check majority
	// handle reconect to other network
	// fake pulse

}

func (g *DiscoveryBootstrap) GetState() insolar.NetworkState {
	return insolar.DiscoveryBootstrap
}

func (g *DiscoveryBootstrap) authorize(ctx context.Context) (*packet.AuthorizeResponse, error) {
	cert := g.CertificateManager.GetCertificate()
	discoveryNodes := network.ExcludeOrigin(cert.GetDiscoveryNodes(), g.NodeKeeper.GetOrigin().ID())
	// todo: sort discoveryNodes

	logger := inslogger.FromContext(ctx)
	for _, n := range discoveryNodes {
		if g.NodeKeeper.GetAccessor().GetActiveNode(*n.GetNodeRef()) != nil {
			logger.Info("Skip discovery already in active list: ", n.GetNodeRef().String())
			continue
		}

		h, err := host.NewHostN(n.GetHost(), *n.GetNodeRef())
		if err != nil {
			inslogger.FromContext(ctx).Error(err.Error())
			continue
		}

		res, err := g.BootstrapRequester.Authorize(ctx, h, cert)
		if err != nil {
			logger.Errorf("Error authorizing to discovery node %s: %s", h.String(), err.Error())
			continue
		}

		if res.Permit == nil {
			logger.Error("Error authorizing, got nil permit.")
			continue
		}

		gotPulse := pulse.FromProto(res.Pulse)
		localPulse, err := g.PulseAccessor.Latest(ctx)
		if err != nil {
			localPulse = *insolar.EphemeralPulse
		}

		if gotPulse.PulseNumber < localPulse.PulseNumber {
			logger.Errorf("Skip authorize response with pulse number %d", gotPulse.PulseNumber)
			continue
		}

		//if err == nil && insolar.IsEphemeralPulse(&localPulse) && gotPulse.PulseNumber > localPulse.PulseNumber {
		//	logger.Info("Last stored pulse.")
		//	continue
		//}

		return res, nil
	}

	return nil, errors.New("Failed to authorize to any discovery node.")
}
