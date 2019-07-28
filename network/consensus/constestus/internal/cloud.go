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
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/constestus/cloud"
	"github.com/insolar/insolar/network/consensus/constestus/internal/interfaces"
	"github.com/pkg/errors"
)

type cloudConfig struct {
	*cloud.Config
}

func (c cloudConfig) VerifyAndReplace() error {
	if c.Network.Delay.Min > c.Network.Delay.Max {
		return errors.New("invalid Network.Delay - Min must <= Max")
	}

	if c.Network.Delay.Variance < 0 {
		return errors.New("invalid Network.Delay.Variance -  must be in [0, Inf)")
	}

	if c.Network.Delay.SpikeProbability < 0 || c.Network.Delay.SpikeProbability > 1 {
		return errors.New("invalid Network.Delay.SpikeProbability - must be in [0, 1]")
	}

	if c.Network.Delay.Spike == 0 {
		c.Network.Delay.Spike = c.Network.Delay.Max
	}

	if c.Network.Pulse.Duration%time.Second != 0 {
		return errors.New("invalid Network.Pulse.Duration - must be a multiple of a second")
	}

	if c.DiscoveryNodes.GetTotal() == 0 {
		return errors.New("invalid DiscoveryNodes - discovery activeNodes must exist")
	}

	if c.Identity.BaseAddr == "" {
		c.Identity.BaseAddr = "127.0.0.1"
	}

	if c.Identity.BasePort == 0 {
		c.Identity.BasePort = 10000
	}

	if c.Identity.BaseID == 0 {
		c.Identity.BaseID = 1000
	}

	return nil
}

type InitializedCloud struct {
	config  cloud.Config
	factory IdentityFactory

	nodes Nodes
}

func NewCloud(config cloud.Config) (*InitializedCloud, error) {
	cloudConfig := cloudConfig{&config}
	err := cloudConfig.VerifyAndReplace()
	if err != nil {
		return nil, err
	}

	factory := newIdentityFactory(config.Identity)

	discoveryIdentities, err := NodeCounts(config.DiscoveryNodes).createIdentities(factory)
	if err != nil {
		return nil, err
	}

	discoveryNodes, err := discoveryIdentities.CreateNodes(discoveryIdentities)
	if err != nil {
		return nil, err
	}

	nodeIdentities, err := NodeCounts(config.Nodes).createIdentities(factory)
	if err != nil {
		return nil, err
	}

	nodes, err := nodeIdentities.CreateNodes(discoveryIdentities)
	if err != nil {
		return nil, err
	}

	allNodes := make([]Node, 0, len(discoveryNodes)+len(nodes))
	allNodes = append(allNodes, discoveryNodes...)
	allNodes = append(allNodes, nodes...)

	return &InitializedCloud{
		config:  config,
		factory: factory,

		nodes: allNodes,
	}, nil
}

func (ic InitializedCloud) Start(ctx context.Context) (*Cloud, error) {
	activeNodes, err := ic.nodes.CreateActiveNodes(ctx, consensus.ReadyNetwork, ic.config.Network)
	if err != nil {
		return nil, err
	}

	if err := activeNodes.Connect(); err != nil {
		return nil, err
	}

	return &Cloud{
		ctx:         ctx,
		activeNodes: activeNodes,
		pulsar:      NewPulsar(ic.config.Network.Pulse),
	}, nil
}

type Cloud struct {
	ctx context.Context

	pulsar         Pulsar
	activeNodes    ActiveNodes
	activeNodesMap map[insolar.ShortNodeID]ActiveNode
}

func (c Cloud) Intercept(nodes ...interfaces.Node) interfaces.TypedInterceptor {
	panic("implement me")
}

func (c Cloud) Pulse() error {
	return c.pulsar.Pulse(c.ctx, c.activeNodes, 4+len(c.activeNodes)/10)
}

func (c Cloud) GetNode(id insolar.ShortNodeID) ActiveNode {
	return c.activeNodes[id]
}
