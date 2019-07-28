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

package internal

import (
	"context"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/constestus/cloud"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/transport"
)

type Nodes []Node

func (ns Nodes) networkNodes() []insolar.NetworkNode {
	networkNodes := make([]insolar.NetworkNode, len(ns))
	for i, n := range ns {
		networkNodes[i] = n.networkNode
	}

	return networkNodes
}

func (ns Nodes) CreateActiveNodes(ctx context.Context, mode consensus.Mode, config cloud.Network) (ActiveNodes, error) {
	knownNodes := ns.networkNodes()
	activeNodes := make([]ActiveNode, len(ns))
	for i, n := range ns {
		activeNode, err := n.CreateActiveNode(ctx, knownNodes, mode, config)
		if err != nil {
			return nil, err
		}

		activeNodes[i] = *activeNode
	}

	return activeNodes, nil
}

type Components struct {
	controller   consensus.Controller
	nodeKeeper   network.NodeKeeper
	transport    transport.DatagramTransport
	pulseHandler network.PulseHandler
}

type Node struct {
	identity    Identity
	networkNode insolar.NetworkNode
	profile     *adapters.StaticProfile
	certificate insolar.Certificate
}

func (n Node) createComponents(
	ctx context.Context,
	knownNodes []insolar.NetworkNode,
	mode consensus.Mode,
	config cloud.Network,
) (*Components, error) {
	nodeKeeper := nodenetwork.NewNodeKeeper(n.networkNode)
	nodeKeeper.SetInitialSnapshot(knownNodes)

	datagramHandler := adapters.NewDatagramHandler()
	pulseHandler := adapters.NewPulseHandler(n.identity.id)

	conf := configuration.NewHostNetwork().Transport
	conf.Address = n.identity.addr
	transportFactory := transport.NewFactory(conf)
	datagramTransport, err := transportFactory.CreateDatagramTransport(datagramHandler)
	if err != nil {
		return nil, err
	}

	strategy := NewDelayNetStrategy(config.Delay)
	delayTransport := strategy.GetLink(datagramTransport)

	controller := consensus.New(ctx, consensus.Dep{
		PrimingCloudStateHash: [64]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		EphemeralPulseAllowed: func() bool { return config.Pulse.Ephemeral },
		KeyProcessor:          keyProcessor,
		Scheme:                scheme,
		CertificateManager:    certificate.NewCertificateManager(n.certificate),
		KeyStore:              keystore.NewInplaceKeyStore(n.identity.privateKey),
		NodeKeeper:            nodeKeeper,
		StateGetter:           &nshGen{nshDelay: config.Consensus.NodeStateHashGenerationDelay},
		PulseChanger: &pulseChanger{
			nodeKeeper: nodeKeeper,
			ctx:        ctx,
		},
		StateUpdater: &stateUpdater{
			nodeKeeper: nodeKeeper,
			ctx:        ctx,
		},
		DatagramTransport: delayTransport,
	}).ControllerFor(mode, datagramHandler, pulseHandler)

	components := Components{
		controller:   controller,
		nodeKeeper:   nodeKeeper,
		transport:    delayTransport,
		pulseHandler: pulseHandler,
	}

	return &components, nil
}

func (n Node) CreateActiveNode(
	ctx context.Context,
	knownNodes []insolar.NetworkNode,
	mode consensus.Mode,
	config cloud.Network,
) (*ActiveNode, error) {

	nodeCtx, _ := inslogger.WithFields(ctx, map[string]interface{}{
		"node_id":   n.identity.id,
		"node_addr": n.identity.addr,
	})

	components, err := n.createComponents(nodeCtx, knownNodes, mode, config)
	if err != nil {
		return nil, err
	}

	return &ActiveNode{
		node:       n,
		components: *components,
		ctx:        nodeCtx,
	}, nil
}
