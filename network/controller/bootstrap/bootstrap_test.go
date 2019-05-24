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

package bootstrap

import (
	"context"
	"crypto"
	"strconv"
	"testing"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/controller/common"
	gateway2 "github.com/insolar/insolar/network/gateway"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	networkTest "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestCert = "../../../certificate/testdata/cert.json"
const TestKeys = "../../../certificate/testdata/keys.json"
const activeNodesCount = 5

type requestMock struct {
	perm Permission
}

func (rm *requestMock) GetSender() insolar.Reference {
	return testutils.RandomRef()
}

func (rm *requestMock) GetSenderHost() *host.Host {
	return nil
}

func (rm *requestMock) GetType() types.PacketType {
	return types.Bootstrap
}

func (rm *requestMock) GetData() interface{} {
	return &NodeBootstrapRequest{
		JoinClaim:     packets.NodeJoinClaim{},
		LastNodePulse: 123,
		Permission:    rm.perm,
	}
}

func (rm *requestMock) GetRequestID() network.RequestID {
	return 1
}

func getBootstrapResults(t *testing.T, ips []string) []*network.BootstrapResult {
	results := make([]*network.BootstrapResult, activeNodesCount)
	for i := 0; i < activeNodesCount; i++ {
		host, err := host.NewHost(ips[i])
		assert.NoError(t, err)

		results[i] = &network.BootstrapResult{
			Host:              host,
			ReconnectRequired: false,
			NetworkSize:       activeNodesCount,
		}
	}
	results[activeNodesCount-1].NetworkSize = activeNodesCount + 1
	return results
}

func getOptions(infinity bool) *common.Options {
	return &common.Options{
		TimeoutMult:            2 * time.Millisecond,
		InfinityBootstrap:      infinity,
		MinTimeout:             100 * time.Millisecond,
		MaxTimeout:             200 * time.Millisecond,
		PingTimeout:            1 * time.Second,
		PacketTimeout:          10 * time.Second,
		BootstrapTimeout:       10 * time.Second,
		CyclicBootstrapEnabled: false,
	}
}

func TestCyclicBootstrap(t *testing.T) {
	ctx := context.Background()

	cs, _ := cryptography.NewStorageBoundCryptographyService(TestKeys)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()
	cert, err := certificate.ReadCertificate(pk, kp, TestCert)
	require.NoError(t, err)
	require.NotEmpty(t, cert.PublicKey)

	activeNodes := make([]insolar.NetworkNode, activeNodesCount)
	ips := make([]string, activeNodesCount)
	for i := 0; i < activeNodesCount; i++ {
		ip := "127.0.0.1:" + strconv.Itoa(i) + strconv.Itoa(i)
		ips[i] = ip
		activeNodes[i] = node.NewNode(insolar.Reference{}, insolar.StaticRoleUnknown, nil, ip, "")
	}

	node := node.NewNode(insolar.Reference{}, insolar.StaticRoleUnknown, nil, "127.0.0.1:8432", "")
	nodekeeper := nodenetwork.NewNodeKeeper(node)
	nodekeeper.SetInitialSnapshot(activeNodes)

	origin := bootstrapper{
		options:                 getOptions(false),
		bootstrapLock:           make(chan struct{}),
		genesisRequestsReceived: make(map[insolar.Reference]*GenesisRequest),
		Certificate:             cert,
		NodeKeeper:              nodekeeper,
	}

	index := origin.getLagerNetorkIndex(ctx, getBootstrapResults(t, ips))
	reconnectRequired := false
	if index >= 0 {
		reconnectRequired = true
	}
	assert.True(t, reconnectRequired)
}

func TestBootstrapRedirect(t *testing.T) {
	ctx := context.Background()

	activeNodes := make([]insolar.NetworkNode, activeNodesCount)
	refs := make([]insolar.Reference, activeNodesCount)
	ips := make([]string, activeNodesCount)
	for i := 0; i < activeNodesCount; i++ {
		ip := "127.0.0.1:" + strconv.Itoa(i) + strconv.Itoa(i)
		ips[i] = ip
		refs[i] = testutils.RandomRef()
		activeNodes[i] = node.NewNode(refs[i], insolar.StaticRoleUnknown, nil, ip, "")
		activeNodes[i].(node.MutableNode).SetState(insolar.NodeReady)
	}

	node := node.NewNode(insolar.Reference{}, insolar.StaticRoleUnknown, nil, "127.0.0.1:8432", "")

	origin := bootstrapper{
		options:                 getOptions(false),
		bootstrapLock:           make(chan struct{}),
		genesisRequestsReceived: make(map[insolar.Reference]*GenesisRequest),
	}

	pulseAccessor := pulse.NewAccessorMock(t)
	pulseAccessor.LatestFunc = func(ctx context.Context) (insolar.Pulse, error) {
		return *insolar.GenesisPulse, nil
	}

	cryptoService := testutils.NewCryptographyServiceMock(t)
	cryptoService.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
		return &insolar.Signature{}, nil
	}
	cryptoService.VerifyFunc = func(crypto.PublicKey, insolar.Signature, []byte) bool {
		return true
	}

	gateway := networkTest.NewGatewayerMock(t)
	gateway.GatewayFunc = func() (r network.Gateway) {
		return gateway2.NewNoNetwork(networkTest.NewGatewayerMock(t), testutils.NewGlobalInsolarLockMock(t),
			networkTest.NewNodeKeeperMock(t), testutils.NewContractRequesterMock(t),
			testutils.NewCryptographyServiceMock(t), testutils.NewMessageBusMock(t),
			testutils.NewCertificateManagerMock(t))
	}

	hostNetwork := networkTest.NewHostNetworkMock(t)
	hostNetwork.BuildResponseFunc = func(p context.Context, p1 network.Request, p2 interface{}) (r network.Response) {
		sender := p1.GetSenderHost()
		host, err := host.NewHost(node.Address())
		assert.NoError(t, err)
		r = packet.NewBuilder(host).Type(p1.GetType()).Receiver(sender).RequestID(p1.GetRequestID()).
			Response(p2).TraceID(inslogger.TraceID(ctx)).Build()
		return r
	}

	accessor := networkTest.NewAccessorMock(t)
	accessor.GetActiveNodesFunc = func() (r []insolar.NetworkNode) {
		return activeNodes
	}
	accessor.GetActiveNodeByShortIDFunc = func(p insolar.ShortNodeID) (r insolar.NetworkNode) {
		return activeNodes[0]
	}
	accessor.GetActiveNodeFunc = func(r insolar.Reference) insolar.NetworkNode {
		for _, n := range activeNodes {
			if n.ID().Equal(r) {
				return n
			}
		}
		return nil
	}

	keeper := networkTest.NewNodeKeeperMock(t)
	keeper.GetAccessorFunc = func() network.Accessor {
		return accessor
	}
	keeper.GetOriginFunc = func() insolar.NetworkNode {
		return node
	}

	discoveryNodes := make([]insolar.DiscoveryNode, activeNodesCount)
	for i := range activeNodes {
		n := testutils.NewDiscoveryNodeMock(t)
		n.GetNodeRefFunc = func() *insolar.Reference {
			return &refs[i]
		}
		n.GetHostFunc = func() string {
			return activeNodes[i].Address()
		}
		discoveryNodes[i] = n
	}

	cert := testutils.NewCertificateMock(t)
	cert.GetDiscoveryNodesFunc = func() []insolar.DiscoveryNode {
		return discoveryNodes
	}

	mngr := component.NewManager(nil)
	mngr.Register(&origin)
	mngr.Inject(cert, hostNetwork, keeper, gateway, pulseAccessor, cryptoService)

	request := requestMock{Permission{}}
	response, err := origin.processBootstrap(ctx, &request)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.GetData().(*NodeBootstrapResponse).Permission.ReconnectTo)
	assert.NotEqual(t, origin.NodeKeeper.GetOrigin().Address(), response.GetData().(*NodeBootstrapResponse).Permission.ReconnectTo)
}

func TestBootstrapRedirectToSelf(t *testing.T) {
	ctx := context.Background()

	activeNodes := make([]insolar.NetworkNode, activeNodesCount)
	refs := make([]insolar.Reference, activeNodesCount)
	ips := make([]string, activeNodesCount)
	for i := 0; i < activeNodesCount; i++ {
		ip := "127.0.0.1:" + strconv.Itoa(i) + strconv.Itoa(i)
		ips[i] = ip
		refs[i] = testutils.RandomRef()
		activeNodes[i] = node.NewNode(refs[i], insolar.StaticRoleUnknown, nil, ip, "")
		activeNodes[i].(node.MutableNode).SetState(insolar.NodePending)
	}

	node := node.NewNode(insolar.Reference{}, insolar.StaticRoleUnknown, nil, "127.0.0.1:8432", "")

	origin := bootstrapper{
		options:                 getOptions(false),
		bootstrapLock:           make(chan struct{}),
		genesisRequestsReceived: make(map[insolar.Reference]*GenesisRequest),
	}

	pulseAccessor := pulse.NewAccessorMock(t)
	pulseAccessor.LatestFunc = func(ctx context.Context) (insolar.Pulse, error) {
		return *insolar.GenesisPulse, nil
	}

	cryptoService := testutils.NewCryptographyServiceMock(t)
	cryptoService.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
		return &insolar.Signature{}, nil
	}
	cryptoService.VerifyFunc = func(crypto.PublicKey, insolar.Signature, []byte) bool {
		return true
	}

	gateway := networkTest.NewGatewayerMock(t)
	gateway.GatewayFunc = func() (r network.Gateway) {
		return gateway2.NewNoNetwork(networkTest.NewGatewayerMock(t), testutils.NewGlobalInsolarLockMock(t),
			networkTest.NewNodeKeeperMock(t), testutils.NewContractRequesterMock(t),
			testutils.NewCryptographyServiceMock(t), testutils.NewMessageBusMock(t),
			testutils.NewCertificateManagerMock(t))
	}

	hostNetwork := networkTest.NewHostNetworkMock(t)
	hostNetwork.BuildResponseFunc = func(p context.Context, p1 network.Request, p2 interface{}) (r network.Response) {
		sender := p1.GetSenderHost()
		host, err := host.NewHost(node.Address())
		assert.NoError(t, err)
		r = packet.NewBuilder(host).Type(p1.GetType()).Receiver(sender).RequestID(p1.GetRequestID()).
			Response(p2).TraceID(inslogger.TraceID(ctx)).Build()
		return r
	}

	accessor := networkTest.NewAccessorMock(t)
	accessor.GetActiveNodesFunc = func() (r []insolar.NetworkNode) {
		return activeNodes
	}
	accessor.GetActiveNodeByShortIDFunc = func(p insolar.ShortNodeID) (r insolar.NetworkNode) {
		return activeNodes[0]
	}
	accessor.GetActiveNodeFunc = func(r insolar.Reference) insolar.NetworkNode {
		for _, n := range activeNodes {
			if n.ID().Equal(r) {
				return n
			}
		}
		return nil
	}

	keeper := networkTest.NewNodeKeeperMock(t)
	keeper.GetAccessorFunc = func() network.Accessor {
		return accessor
	}
	keeper.GetOriginFunc = func() insolar.NetworkNode {
		return node
	}

	discoveryNodes := make([]insolar.DiscoveryNode, activeNodesCount)
	for i := range activeNodes {
		n := testutils.NewDiscoveryNodeMock(t)
		n.GetNodeRefFunc = func() *insolar.Reference {
			return &refs[i]
		}
		n.GetHostFunc = func() string {
			return activeNodes[i].Address()
		}
		discoveryNodes[i] = n
	}

	cert := testutils.NewCertificateMock(t)
	cert.GetDiscoveryNodesFunc = func() []insolar.DiscoveryNode {
		return discoveryNodes
	}

	mngr := component.NewManager(nil)
	mngr.Register(&origin)
	mngr.Inject(cert, hostNetwork, keeper, gateway, pulseAccessor, cryptoService)

	request := requestMock{Permission{}}
	response, err := origin.processBootstrap(ctx, &request)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.GetData().(*NodeBootstrapResponse).Permission.ReconnectTo)
	assert.Equal(t, origin.NodeKeeper.GetOrigin().Address(), response.GetData().(*NodeBootstrapResponse).Permission.ReconnectTo)
}
