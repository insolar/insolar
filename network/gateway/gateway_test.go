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
	"testing"

	"github.com/insolar/insolar/certificate"

	"github.com/insolar/insolar/insolar/reply"

	"github.com/insolar/insolar/network"
	testnet "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/insolar"
)

func emtygateway(t *testing.T) network.Gateway {
	return NewNoNetwork(testnet.NewGatewayerMock(t), testutils.NewGlobalInsolarLockMock(t),
		testnet.NewNodeKeeperMock(t), testutils.NewContractRequesterMock(t),
		testutils.NewCryptographyServiceMock(t), testutils.NewMessageBusMock(t),
		testutils.NewCertificateManagerMock(t))
}

func TestSWitch(t *testing.T) {
	ctx := context.Background()

	nodekeeper := testnet.NewNodeKeeperMock(t)
	gatewayer := testnet.NewGatewayerMock(t)
	GIL := testutils.NewGlobalInsolarLockMock(t)
	GIL.AcquireMock.Return()
	MB := testutils.NewMessageBusMock(t)

	MB.MustRegisterFunc = func(p insolar.MessageType, p1 insolar.MessageHandler) {}

	ge := NewNoNetwork(gatewayer, GIL,
		nodekeeper, testutils.NewContractRequesterMock(t),
		testutils.NewCryptographyServiceMock(t), MB,
		testutils.NewCertificateManagerMock(t))

	require.NotNil(t, ge)
	require.Equal(t, "NoNetworkState", ge.GetState().String())

	ge.Run(ctx)

	nodekeeper.IsBootstrappedFunc = func() (r bool) { return true }
	gatewayer.GatewayFunc = func() (r network.Gateway) { return ge }
	gatewayer.SetGatewayFunc = func(p network.Gateway) { ge = p }
	gilreleased := false
	GIL.ReleaseFunc = func(p context.Context) { gilreleased = true }

	ge.OnPulse(ctx, insolar.Pulse{})

	require.Equal(t, "CompleteNetworkState", ge.GetState().String())
	require.True(t, gilreleased)
	cref := testutils.RandomRef()

	for _, state := range []insolar.NetworkState{insolar.NoNetworkState,
		insolar.AuthorizationNetworkState, insolar.JetlessNetworkState, insolar.VoidNetworkState} {
		ge = ge.NewGateway(state)
		require.Equal(t, state, ge.GetState())
		ge.Run(ctx)
		au := ge.Auther()

		_, err := au.GetCert(ctx, &cref)
		require.Error(t, err)

		_, err = au.ValidateCert(ctx, &certificate.Certificate{})
		require.Error(t, err)

		ge.OnPulse(ctx, insolar.Pulse{})

	}

}

func TestDumbComplete_GetCert(t *testing.T) {
	ctx := context.Background()

	nodekeeper := testnet.NewNodeKeeperMock(t)
	gatewayer := testnet.NewGatewayerMock(t)
	GIL := testutils.NewGlobalInsolarLockMock(t)
	GIL.AcquireMock.Return()
	MB := testutils.NewMessageBusMock(t)

	MB.MustRegisterFunc = func(p insolar.MessageType, p1 insolar.MessageHandler) {}

	CR := testutils.NewContractRequesterMock(t)
	CM := testutils.NewCertificateManagerMock(t)
	ge := NewNoNetwork(gatewayer, GIL,
		nodekeeper, CR,
		testutils.NewCryptographyServiceMock(t), MB,
		CM)

	require.NotNil(t, ge)
	require.Equal(t, "NoNetworkState", ge.GetState().String())

	ge.Run(ctx)

	nodekeeper.IsBootstrappedFunc = func() (r bool) { return true }
	gatewayer.GatewayFunc = func() (r network.Gateway) { return ge }
	gatewayer.SetGatewayFunc = func(p network.Gateway) { ge = p }
	gilreleased := false
	GIL.ReleaseFunc = func(p context.Context) { gilreleased = true }

	ge.OnPulse(ctx, insolar.Pulse{})

	require.Equal(t, "CompleteNetworkState", ge.GetState().String())
	require.True(t, gilreleased)

	cref := testutils.RandomRef()

	CR.SendRequestFunc = func(ctx context.Context, ref *insolar.Reference, method string, argsIn []interface{},
	) (r insolar.Reply, r1 error) {
		require.Equal(t, &cref, ref)
		require.Equal(t, "GetNodeInfo", method)
		repl, _ := insolar.Serialize(struct {
			PublicKey string
			Role      insolar.StaticRole
		}{"LALALA", insolar.StaticRoleVirtual})
		return &reply.CallMethod{
			Result: repl,
		}, nil
	}

	CM.GetCertificateFunc = func() (r insolar.Certificate) { return &certificate.Certificate{} }
	CM.NewUnsignedCertificateFunc = func(p string, p1 string, p2 string) (r insolar.Certificate, r1 error) {
		return &certificate.Certificate{}, nil
	}
	cert, err := ge.Auther().GetCert(ctx, &cref)

	require.NoError(t, err)
	require.NotNil(t, cert)
	require.Equal(t, cert, &certificate.Certificate{})
}
