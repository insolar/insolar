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
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/utils"
)

// TODO Slightly ugly, decide how to inject anything without exporting Base
// TODO Remove message bus here and switch communication to network.rpc
// NewNoNetwork this initial constructor have special signature to be called outside
func NewNoNetwork(n network.Gatewayer, gil insolar.GlobalInsolarLock,
	nk network.NodeKeeper, cr insolar.ContractRequester,
	cs insolar.CryptographyService, mb insolar.MessageBus,
	cm insolar.CertificateManager) network.Gateway {
	return (&Base{
		Gatewayer: n, GIL: gil,
		NodeKeeper: nk, ContractRequester: cr,
		CryptographyService: cs, MessageBus: mb,
		CertificateManager: cm,
	}).NewGateway(insolar.NoNetworkState)
}

// NoNetwork initial state
type NoNetwork struct {
	*Base

	isDiscovery bool
	skip        uint32
}

func (g *NoNetwork) Run(ctx context.Context) {

	cert := g.CertificateManager.GetCertificate()
	if len(cert.GetDiscoveryNodes()) == 0 {
		g.zeroBootstrap(ctx)
		// create complete network
		return
	}

	// run bootstrap
	g.isDiscovery = utils.OriginIsDiscovery(cert)

	log.Info("TODO: remove! Bootstrapping network...")
	_, err := g.Bootstrapper.Bootstrap(ctx)
	if err != nil {
		err = errors.Wrap(err, "Failed to bootstrap network")
		panic(err.Error())
	}

}

func (g *NoNetwork) GetState() insolar.NetworkState {
	return insolar.NoNetworkState
}

func (g *NoNetwork) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	return g.Base.OnPulse(ctx, pu)
}

func (g *NoNetwork) ShoudIgnorePulse(ctx context.Context, newPulse insolar.Pulse) bool {
	if true { //!g.Base.NodeKeeper.IsBootstrapped() { always true here
		g.Bootstrapper.SetLastPulse(newPulse.NextPulseNumber)
		return true
	}

	return g.isDiscovery && !g.NodeKeeper.GetConsensusInfo().IsJoiner() &&
		newPulse.PulseNumber <= g.Bootstrapper.GetLastPulse()+insolar.PulseNumber(g.skip)
}

func (g *NoNetwork) connectToNewNetwork(ctx context.Context, address string) {
	g.NodeKeeper.GetClaimQueue().Push(&packets.ChangeNetworkClaim{Address: address})
	logger := inslogger.FromContext(ctx)

	// node, err := findNodeByAddress(address, g.CertificateManager.GetCertificate().GetDiscoveryNodes())
	// if err != nil {
	// 	logger.Warnf("Failed to find a discovery node: ", err)
	// }

	err := g.Bootstrapper.AuthenticateToDiscoveryNode(ctx, nil /*node*/)
	if err != nil {
		logger.Errorf("Failed to authenticate a node: " + err.Error())
	}
}

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
