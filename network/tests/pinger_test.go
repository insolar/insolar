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
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/pulsenetwork"
	"github.com/insolar/insolar/network/transport"
	mock "github.com/insolar/insolar/testutils/network"
)

const (
	PULSENUMBER = 155
	ID1         = "4K2V1kpVycZ6qSFsNdz2FtpNxnJs17eBNzf9rdCMcKoe"
	DOMAIN      = ".4F7BsTMVPKFshM1MwLf6y23cid6fL3xMpazVoF9krzUw"
)

type testnode struct {
	name        string
	ctx         context.Context
	cm          *component.Manager
	Factory     transport.Factory `inject:""`
	hostNetwork network.HostNetwork
	transport   transport.StreamTransport
	host        *host.Host
}

func (n *testnode) start() error {
	err := n.cm.Start(n.ctx)
	if err != nil {
		return err
	}

	publicAddress := n.hostNetwork.PublicAddress()

	h, err := host.NewHost(publicAddress)
	if err != nil {
		return errors.Wrap(err, "[ NewDistributor ] failed to create pulsar host")
	}
	h.NodeID = insolar.Reference{}

	n.host = h
	return nil
}

func (n *testnode) stop() error {
	return n.cm.Stop(n.ctx)
}

func newNode(t *testing.T, name string) *testnode {
	var err error
	n := &testnode{name: name}
	n.ctx = context.Background()
	n.cm = component.NewManager(nil)
	n.hostNetwork, err = hostnetwork.NewHostNetwork(ID1 + DOMAIN)
	require.NoError(t, err)

	routingTable := mock.NewRoutingTableMock(t)
	routingTable.AddToKnownHostsMock.Set(func(p *host.Host) {
		log.Infof("AddToKnownHostsMock: %s", p.String())
	})

	n.cm.Inject(
		n,
		routingTable,
		n.hostNetwork,
		transport.NewFactory(configuration.NewHostNetwork().Transport),
	)
	err = n.cm.Init(n.ctx)
	require.NoError(t, err)

	n.hostNetwork.RegisterRequestHandler(types.Ping, func(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
		inslogger.FromContext(ctx).Infof("[ HandlePing ] response to ping from %s", request.GetSenderHost().String())
		return n.hostNetwork.BuildResponse(ctx, request, &packet.Ping{}), nil
	})
	return n
}

func (d *testnode) pingHost(ctx context.Context, host *host.Host) error {
	logger := inslogger.FromContext(ctx)

	pingCall, err := d.hostNetwork.SendRequestToHost(ctx, types.Ping, &packet.Ping{}, host)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "[ pingHost ] failed to send ping request")
	}

	logger.Debugf("before ping request")
	result, err := pingCall.WaitResponse(time.Second * 2)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "[ pingHost ] failed to get ping result")
	}

	host.NodeID = result.GetSender()

	return nil
}

func TestPingerStress(t *testing.T) {
	t.Skip("skip until INS-2978 will fixed")
	// defer leaktest.Check(t)()

	nodes := make([]*testnode, 0)
	nodes = append(nodes, newNode(t, "A"))
	nodes = append(nodes, newNode(t, "B"))
	nodes = append(nodes, newNode(t, "C"))
	nodes = append(nodes, newNode(t, "D"))
	nodes = append(nodes, newNode(t, "E"))

	for _, n := range nodes {
		require.NoError(t, n.start())
	}
	// defer func() {
	// 	for _, n := range nodes {
	// 		require.NoError(t, n.stop())
	// 	}
	// }()

	pinger := newNode(t, "Pinger")
	require.NoError(t, pinger.start())
	// defer require.NoError(t, pinger.stop())

	// <-time.After(time.Second)

	for i := 0; i < 2; i++ {
		for _, n := range nodes {
			err := pinger.pingHost(context.Background(), n.host)
			assert.NoError(t, err)
		}
		require.NoError(t, pinger.hostNetwork.Stop(pinger.ctx))
		require.NoError(t, pinger.hostNetwork.Start(pinger.ctx))
	}
}

func TestDistributorStress(t *testing.T) {
	cfg := configuration.NewPulsar().PulseDistributor

	nodes := make([]*testnode, 0)
	nodes = append(nodes, newNode(t, "A"))
	nodes = append(nodes, newNode(t, "B"))
	nodes = append(nodes, newNode(t, "C"))
	nodes = append(nodes, newNode(t, "D"))
	nodes = append(nodes, newNode(t, "E"))

	for _, n := range nodes {
		require.NoError(t, n.start())
		cfg.BootstrapHosts = append(cfg.BootstrapHosts, n.host.Address.String())
	}

	pulsarCtx := context.Background()
	connFactory := transport.NewFactory(configuration.NewHostNetwork().Transport)
	pulsar, err := pulsenetwork.NewDistributor(cfg)
	require.NoError(t, err)

	cm := component.NewManager(nil)
	cm.Inject(pulsar, connFactory)
	require.NoError(t, cm.Init(pulsarCtx))
	require.NoError(t, cm.Start(pulsarCtx))

	for i := 0; i < 50; i++ {
		pulsar.Distribute(pulsarCtx, *insolar.GenesisPulse)
	}
}
