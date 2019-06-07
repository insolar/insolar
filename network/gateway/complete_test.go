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
	"errors"
	"testing"

	"github.com/insolar/insolar/testutils/network"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func mockCryptographyService(t *testing.T, ok bool) insolar.CryptographyService {
	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*insolar.Signature, error) {
		if ok {
			sig := insolar.SignatureFromBytes([]byte("test_sig"))
			return &sig, nil
		}
		return nil, errors.New("test_error")
	}
	return cs
}

func mockCertificateManager(t *testing.T, certNodeRef *insolar.Reference, discoveryNodeRef *insolar.Reference, unsignCertOk bool) *testutils.CertificateManagerMock {
	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() insolar.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef.String(),
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				{
					NodeRef:     discoveryNodeRef.String(),
					PublicKey:   "test_discovery_public_key",
					Host:        "test_discovery_host",
					NetworkSign: []byte("test_network_sign"),
				},
			},
		}
	}
	cm.NewUnsignedCertificateFunc = func(key string, role string, nodeRef string) (insolar.Certificate, error) {
		require.Equal(t, "test_node_public_key", key)
		require.Equal(t, "virtual", role)

		if unsignCertOk {
			return &certificate.Certificate{
				AuthorizationCertificate: certificate.AuthorizationCertificate{
					PublicKey: key,
					Reference: nodeRef,
					Role:      role,
				},
				RootDomainReference: "test_root_domain_ref",
				MajorityRule:        0,
				PulsarPublicKeys:    []string{},
				BootstrapNodes: []certificate.BootstrapNode{
					{
						PublicKey:   "test_discovery_public_key",
						Host:        "test_discovery_host",
						NetworkSign: []byte("test_network_sign"),
						NodeRef:     discoveryNodeRef.String(),
					},
				},
			}, nil
		}
		return nil, errors.New("test_error")
	}
	return cm
}

func mockMessageBus(t *testing.T, ok bool, ref *insolar.Reference, discovery *insolar.Reference) *testutils.MessageBusMock {
	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p insolar.MessageType, handler insolar.MessageHandler) {
		require.Equal(t, p, insolar.TypeNodeSignRequest)
	}
	mb.SendFunc = func(p context.Context, msg insolar.Message, options *insolar.MessageSendOptions) (insolar.Reply, error) {
		require.Equal(t, ref, msg.(*message.NodeSignPayload).NodeRef)
		require.Equal(t, discovery, options.Receiver)
		if ok {
			return &reply.NodeSign{
				Sign: []byte("test_sig"),
			}, nil
		}
		return nil, errors.New("test_error")
	}
	return mb
}

func mockReply(t *testing.T) []byte {
	node, err := insolar.MarshalArgs(struct {
		PublicKey string
		Role      insolar.StaticRole
	}{
		PublicKey: "test_node_public_key",
		Role:      insolar.StaticRoleVirtual,
	}, nil)
	require.NoError(t, err)
	return []byte(node)
}

func mockContractRequester(t *testing.T, nodeRef insolar.Reference, ok bool, r []byte) insolar.ContractRequester {
	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *insolar.Reference, method string, args []interface{}) (insolar.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		if ok {
			return &reply.CallMethod{
				Result: r,
			}, nil
		}
		return nil, errors.New("test_error")
	}
	return cr
}

func TestComplete_GetCert(t *testing.T) {
	nodeRef := testutils.RandomRef()
	certNodeRef := testutils.RandomRef()

	gatewayer := network.NewGatewayerMock(t)
	GIL := testutils.NewGlobalInsolarLockMock(t)
	GIL.AcquireMock.Return()
	nodekeeper := network.NewNodeKeeperMock(t)

	cr := mockContractRequester(t, nodeRef, true, mockReply(t))
	mb := mockMessageBus(t, true, &nodeRef, &certNodeRef)
	cm := mockCertificateManager(t, &certNodeRef, &certNodeRef, true)
	cs := mockCryptographyService(t, true)

	ge := NewNoNetwork(gatewayer, GIL, nodekeeper, cr, cs, mb, cm)
	ge = ge.NewGateway(insolar.CompleteNetworkState)
	ctx := context.Background()
	result, err := ge.Auther().GetCert(ctx, &nodeRef)
	require.NoError(t, err)

	cert := result.(*certificate.Certificate)
	require.Equal(t, "test_node_public_key", cert.PublicKey)
	require.Equal(t, nodeRef.String(), cert.Reference)
	require.Equal(t, "virtual", cert.Role)
	require.Equal(t, 0, cert.MajorityRule)
	require.Equal(t, uint(0), cert.MinRoles.Virtual)
	require.Equal(t, uint(0), cert.MinRoles.HeavyMaterial)
	require.Equal(t, uint(0), cert.MinRoles.LightMaterial)
	require.Equal(t, []string{}, cert.PulsarPublicKeys)
	require.Equal(t, "test_root_domain_ref", cert.RootDomainReference)
	require.Equal(t, 1, len(cert.BootstrapNodes))
	require.Equal(t, "test_discovery_public_key", cert.BootstrapNodes[0].PublicKey)
	require.Equal(t, []byte("test_network_sign"), cert.BootstrapNodes[0].NetworkSign)
	require.Equal(t, "test_discovery_host", cert.BootstrapNodes[0].Host)
	require.Equal(t, []byte("test_sig"), cert.BootstrapNodes[0].NodeSign)
	require.Equal(t, certNodeRef.String(), cert.BootstrapNodes[0].NodeRef)
}

func TestComplete_handler(t *testing.T) {
	nodeRef := testutils.RandomRef()
	certNodeRef := testutils.RandomRef()

	gatewayer := network.NewGatewayerMock(t)
	GIL := testutils.NewGlobalInsolarLockMock(t)
	GIL.AcquireMock.Return()
	nodekeeper := network.NewNodeKeeperMock(t)

	cr := mockContractRequester(t, nodeRef, true, mockReply(t))
	mb := mockMessageBus(t, true, &nodeRef, &certNodeRef)
	cm := mockCertificateManager(t, &certNodeRef, &certNodeRef, true)
	cs := mockCryptographyService(t, true)

	ge := NewNoNetwork(gatewayer, GIL, nodekeeper, cr, cs, mb, cm)
	ge = ge.NewGateway(insolar.CompleteNetworkState)
	ctx := context.Background()

	result, err := ge.(*Complete).signCertHandler(ctx, &message.Parcel{Msg: &message.NodeSignPayload{NodeRef: &nodeRef}})
	require.NoError(t, err)
	require.Equal(t, []byte("test_sig"), result.(*reply.NodeSign).Sign)
}
