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

// +build networktest

package tests

import (
	"context"
	"sync/atomic"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/bootstrap"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/version"
)

type fakeBootstrap struct {
	NodeKeeper network.NodeKeeper `inject:""`

	suite *fixture
}

func (b *fakeBootstrap) Init(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("Using fakeBootstrap")
	return nil
}

func (b *fakeBootstrap) BootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error) {
	if atomic.LoadUint32(&b.suite.discoveriesAreBootstrapped) == 1 {
		return nil, bootstrap.ErrReconnectRequired
	}
	nodes := make([]insolar.NetworkNode, 0)
	for _, bNode := range b.suite.bootstrapNodes {
		if bNode.id.Equal(b.NodeKeeper.GetOrigin().ID()) {
			continue
		}
		n := convertNode(bNode)
		n.(node.MutableNode).SetState(insolar.NodeUndefined)
		nodes = append(nodes, n)
	}
	b.NodeKeeper.GetOrigin().(node.MutableNode).SetState(insolar.NodeUndefined)
	nodes = append(nodes, b.NodeKeeper.GetOrigin())
	b.NodeKeeper.SetInitialSnapshot(nodes)
	// this result is currently not used in discovery bootstrap scenario, so do not fill the Host
	return &network.BootstrapResult{
		Host:              nil,
		ReconnectRequired: false,
		NetworkSize:       len(b.suite.bootstrapNodes),
	}, nil
}

func (b *fakeBootstrap) SetLastPulse(number insolar.PulseNumber) {}

func (b *fakeBootstrap) GetLastPulse() insolar.PulseNumber {
	return 0
}

func convertNode(n *networkNode) insolar.NetworkNode {
	pk, err := n.cryptographyService.GetPublicKey()
	if err != nil {
		panic(err)
	}
	return node.NewNode(n.id, n.role, pk, n.host, version.Version)
}

func newFakeBootstrap(suite *fixture) *fakeBootstrap {
	return &fakeBootstrap{
		suite: suite,
	}
}
