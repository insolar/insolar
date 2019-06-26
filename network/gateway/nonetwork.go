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

// TODO: spans, metrics

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
)

func newNoNetwork(b *Base) *NoNetwork {
	return &NoNetwork{Base: b}
}

// NoNetwork initial state
type NoNetwork struct {
	*Base
}

func (g *NoNetwork) Run(ctx context.Context) {

	cert := g.CertificateManager.GetCertificate()
	discoveryNodes := network.ExcludeOrigin(cert.GetDiscoveryNodes(), g.NodeKeeper.GetOrigin().ID())

	// TODO: clear NodeKeeper state

	if len(discoveryNodes) == 0 {
		g.zeroBootstrap(ctx)
		// create complete network
		g.Gatewayer.SwitchState(insolar.CompleteNetworkState)
		return
	}

	// run bootstrap
	if !network.OriginIsDiscovery(cert) {
		g.Gatewayer.SwitchState(insolar.JoinerBootstrap)
		return
	}

	pulse, err := g.PulseAccessor.Latest(ctx)
	if err != nil {
		g.Gatewayer.SwitchState(insolar.JoinerBootstrap)
		return
	}

	if pulse.PulseNumber > insolar.FirstPulseNumber {
		g.Gatewayer.SwitchState(insolar.JoinerBootstrap)
		return
	}

	g.Gatewayer.SwitchState(insolar.DiscoveryBootstrap)
}

func (g *NoNetwork) GetState() insolar.NetworkState {
	return insolar.NoNetworkState
}

// func (g *NoNetwork) connectToNewNetwork(ctx context.Context, address string) {
// 	g.NodeKeeper.GetClaimQueue().Push(&packets.ChangeNetworkClaim{Address: address})
// 	logger := inslogger.FromContext(ctx)
//
// 	// node, err := findNodeByAddress(address, g.CertificateManager.GetCertificate().GetDiscoveryNodes())
// 	// if err != nil {
// 	// 	logger.Warnf("Failed to find a discovery node: ", err)
// 	// }
//
// 	err := g.Bootstrapper.AuthenticateToDiscoveryNode(ctx, nil /*node*/)
// 	if err != nil {
// 		logger.Errorf("Failed to authenticate a node: " + err.Error())
// 	}
// }

// todo: remove this workaround
func (g *NoNetwork) zeroBootstrap(ctx context.Context) {
	inslogger.FromContext(ctx).Info("[ Bootstrap ] Zero bootstrap")
	g.NodeKeeper.SetInitialSnapshot([]insolar.NetworkNode{g.NodeKeeper.GetOrigin()})
}

func findNodeByAddress(address string, nodes []insolar.DiscoveryNode) (insolar.DiscoveryNode, error) {
	for _, node := range nodes {
		if node.GetHost() == address {
			return node, nil
		}
	}
	return nil, errors.New("Failed to find a discovery node with address: " + address)
}
