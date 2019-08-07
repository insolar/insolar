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
	"bytes"
	"context"
	"crypto/rand"
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

func newJoinerBootstrap(b *Base) *JoinerBootstrap {
	return &JoinerBootstrap{b}
}

// JoinerBootstrap void network state
type JoinerBootstrap struct {
	*Base
}

func (g *JoinerBootstrap) Run(ctx context.Context, pulse insolar.Pulse) {
	logger := inslogger.FromContext(ctx)
	permit, err := g.authorize(ctx)
	if err != nil {
		logger.Error(err.Error())
		g.Gatewayer.SwitchState(ctx, insolar.NoNetworkState, pulse)
		return
	}

	resp, err := g.BootstrapRequester.Bootstrap(ctx, permit, *g.originCandidate, &pulse)
	if err != nil {
		logger.Error(err.Error())
		g.Gatewayer.SwitchState(ctx, insolar.NoNetworkState, pulse)
		return
	}

	g.bootstrapETA = time.Second * time.Duration(resp.ETASeconds)
	g.Gatewayer.SwitchState(ctx, insolar.WaitConsensus, pulse)
}

func (g *JoinerBootstrap) GetState() insolar.NetworkState {
	return insolar.JoinerBootstrap
}

func (g *JoinerBootstrap) authorize(ctx context.Context) (*packet.Permit, error) {
	cert := g.CertificateManager.GetCertificate()
	discoveryNodes := network.ExcludeOrigin(cert.GetDiscoveryNodes(), g.NodeKeeper.GetOrigin().ID())

	entropy := make([]byte, insolar.RecordRefSize)
	if _, err := rand.Read(entropy); err != nil {
		panic("Failed to get bootstrap entropy")
	}

	sort.Slice(discoveryNodes, func(i, j int) bool {
		return bytes.Compare(
			xor(*discoveryNodes[i].GetNodeRef(), entropy),
			xor(*discoveryNodes[j].GetNodeRef(), entropy)) < 0
	})

	bestResult := &packet.AuthorizeResponse{}

	for _, n := range discoveryNodes {
		h, _ := host.NewHostN(n.GetHost(), *n.GetNodeRef())

		res, err := g.BootstrapRequester.Authorize(ctx, h, cert)
		if err != nil {
			inslogger.FromContext(ctx).Warnf("Error authorizing to host %s: %s", h.String(), err.Error())
			continue
		}

		if int(res.DiscoveryCount) < cert.GetMajorityRule() {
			inslogger.FromContext(ctx).Infof(
				"Check MajorityRule failed on authorize, expect %d, got %d",
				cert.GetMajorityRule(),
				res.DiscoveryCount,
			)

			if res.DiscoveryCount > bestResult.DiscoveryCount {
				bestResult = res
			}

			continue
		}

		return res.Permit, nil
	}

	if network.OriginIsDiscovery(cert) && bestResult.Permit != nil {
		return bestResult.Permit, nil
	}

	return nil, errors.New("failed to authorize to any discovery node")
}

func xor(ref insolar.Reference, entropy []byte) []byte {
	for i, d := range ref {
		ref[i] = entropy[i] ^ d
	}
	return ref[:]
}
