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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/log" // TODO remove before merge

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
)

// Base is abstract class for gateways

type Base struct {
	Self                network.Gateway
	Network             network.Gatewayer
	Nodekeeper          network.NodeKeeper
	ContractRequester   insolar.ContractRequester
	CryptographyService insolar.CryptographyService
	CertificateManager  insolar.CertificateManager
	GIL                 insolar.GlobalInsolarLock
	MessageBus          insolar.MessageBus
}

// NewGateway creates new gateway on top of existing
func (g *Base) NewGateway(state insolar.NetworkState) network.Gateway {
	log.Infof("NewGateway %s", state.String())
	switch state {
	case insolar.NoNetworkState:
		g.Self = &NoNetwork{g}
	case insolar.VoidNetworkState:
		g.Self = NewVoid(g)
	case insolar.JetlessNetworkState:
		g.Self = NewJetless(g)
	case insolar.AuthorizationNetworkState:
		g.Self = NewAuthorisation(g)
	case insolar.CompleteNetworkState:
		g.Self = NewComplete(g)
	default:
		panic("Try to switch network to unknown state. Memory of process is inconsistent.")
	}
	return g.Self
}

func (g *Base) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	if g.Nodekeeper == nil {
		return nil
	}
	if g.Nodekeeper.IsBootstrapped() {
		g.Network.SetGateway(g.Network.Gateway().NewGateway(insolar.CompleteNetworkState))
		g.Network.Gateway().Run(ctx)
	}
	return nil
}

// Auther casts us to Auther or obtain it in another way
func (g *Base) Auther() network.Auther {
	if ret, ok := g.Self.(network.Auther); ok {
		return ret
	}
	panic("Our network gateway suddenly is not an Auther")
}

// GetCert method returns node certificate by requesting sign from discovery nodes
func (g *Base) GetCert(ctx context.Context, ref *insolar.Reference) (insolar.Certificate, error) {
	return nil, errors.New("GetCert() in non active mode")
}

// ValidateCert validates node certificate
func (g *Base) ValidateCert(ctx context.Context, certificate insolar.AuthorizationCertificate) (bool, error) {
	return false, errors.New("ValidateCert() in non active mode")
}

func (g *Base) FilterJoinerNodes(certificate insolar.Certificate, nodes []insolar.NetworkNode) []insolar.NetworkNode {
	dNodes := make(map[insolar.Reference]struct{}, len(certificate.GetDiscoveryNodes()))
	for _, dn := range certificate.GetDiscoveryNodes() {
		dNodes[*dn.GetNodeRef()] = struct{}{}
	}
	ret := []insolar.NetworkNode{}
	for _, n := range nodes {
		if _, ok := dNodes[n.ID()]; ok {
			ret = append(ret, n)
		}
	}
	return ret
}
