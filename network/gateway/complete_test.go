// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/testutils"
	mock "github.com/insolar/insolar/testutils/network"
)

func mockCryptographyService(t *testing.T, ok bool) insolar.CryptographyService {
	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignMock.Set(func(data []byte) (*insolar.Signature, error) {
		if ok {
			sig := insolar.SignatureFromBytes([]byte("test_sig"))
			return &sig, nil
		}
		return nil, errors.New("test_error")
	})
	return cs
}

func mockCertificateManager(t *testing.T, certNodeRef *insolar.Reference, discoveryNodeRef *insolar.Reference, unsignCertOk bool) *testutils.CertificateManagerMock {
	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateMock.Set(func() insolar.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef.String(),
				Role:      "virtual",
			},
			RootDomainReference: "test_root_domain_ref",
			MajorityRule:        0,
			PulsarPublicKeys:    []string{},
			BootstrapNodes: []certificate.BootstrapNode{
				{
					NodeRef:     discoveryNodeRef.String(),
					PublicKey:   "test_discovery_public_key",
					Host:        "test_discovery_host",
					NetworkSign: []byte("test_network_sign"),
				},
			},
		}
	})
	return cm
}

func mockReply(t *testing.T) []byte {
	res := struct {
		PublicKey string
		Role      insolar.StaticRole
	}{
		PublicKey: "test_node_public_key",
		Role:      insolar.StaticRoleVirtual,
	}
	node, err := foundation.MarshalMethodResult(res, nil)
	require.NoError(t, err)
	return node
}

func mockContractRequester(t *testing.T, nodeRef insolar.Reference, ok bool, r []byte) insolar.ContractRequester {
	cr := testutils.NewContractRequesterMock(t)
	cr.CallMock.Set(func(ctx context.Context, ref *insolar.Reference, method string, argsIn []interface{}, p insolar.PulseNumber) (r1 insolar.Reply, r2 *insolar.Reference, err error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(argsIn))
		if ok {
			return &reply.CallMethod{
				Result: r,
			}, nil, nil
		}
		return nil, nil, errors.New("test_error")
	})
	return cr
}

func mockPulseManager(t *testing.T) insolar.PulseManager {
	pm := testutils.NewPulseManagerMock(t)
	return pm
}

func TestComplete_GetCert(t *testing.T) {
	nodeRef := gen.Reference()
	certNodeRef := gen.Reference()

	gatewayer := mock.NewGatewayerMock(t)
	nodekeeper := mock.NewNodeKeeperMock(t)
	hn := mock.NewHostNetworkMock(t)

	cr := mockContractRequester(t, nodeRef, true, mockReply(t))
	cm := mockCertificateManager(t, &certNodeRef, &certNodeRef, true)
	cs := mockCryptographyService(t, true)
	pm := mockPulseManager(t)
	pa := mock.NewPulseAccessorMock(t)

	var ge network.Gateway
	ge = newNoNetwork(&Base{
		Gatewayer:           gatewayer,
		NodeKeeper:          nodekeeper,
		HostNetwork:         hn,
		ContractRequester:   cr,
		CertificateManager:  cm,
		CryptographyService: cs,
		PulseManager:        pm,
		PulseAccessor:       pa,
	})
	ge = ge.NewGateway(context.Background(), insolar.CompleteNetworkState)
	ctx := context.Background()

	pa.GetLatestPulseMock.Expect(ctx).Return(*insolar.GenesisPulse, nil)

	result, err := ge.Auther().GetCert(ctx, &nodeRef)
	require.NoError(t, err)

	cert := result.(*certificate.Certificate)
	assert.Equal(t, "test_node_public_key", cert.PublicKey)
	assert.Equal(t, nodeRef.String(), cert.Reference)
	assert.Equal(t, "virtual", cert.Role)
	assert.Equal(t, 0, cert.MajorityRule)
	assert.Equal(t, uint(0), cert.MinRoles.Virtual)
	assert.Equal(t, uint(0), cert.MinRoles.HeavyMaterial)
	assert.Equal(t, uint(0), cert.MinRoles.LightMaterial)
	assert.Equal(t, []string{}, cert.PulsarPublicKeys)
	assert.Equal(t, "test_root_domain_ref", cert.RootDomainReference)
	assert.Equal(t, 1, len(cert.BootstrapNodes))
	assert.Equal(t, "test_discovery_public_key", cert.BootstrapNodes[0].PublicKey)
	assert.Equal(t, []byte("test_network_sign"), cert.BootstrapNodes[0].NetworkSign)
	assert.Equal(t, "test_discovery_host", cert.BootstrapNodes[0].Host)
	assert.Equal(t, []byte("test_sig"), cert.BootstrapNodes[0].NodeSign)
	assert.Equal(t, certNodeRef.String(), cert.BootstrapNodes[0].NodeRef)
}

func TestComplete_handler(t *testing.T) {
	nodeRef := gen.Reference()
	certNodeRef := gen.Reference()

	gatewayer := mock.NewGatewayerMock(t)
	nodekeeper := mock.NewNodeKeeperMock(t)

	cr := mockContractRequester(t, nodeRef, true, mockReply(t))
	cm := mockCertificateManager(t, &certNodeRef, &certNodeRef, true)
	cs := mockCryptographyService(t, true)
	pm := mockPulseManager(t)
	pa := mock.NewPulseAccessorMock(t)

	hn := mock.NewHostNetworkMock(t)

	var ge network.Gateway
	ge = newNoNetwork(&Base{
		Gatewayer:           gatewayer,
		NodeKeeper:          nodekeeper,
		HostNetwork:         hn,
		ContractRequester:   cr,
		CertificateManager:  cm,
		CryptographyService: cs,
		PulseManager:        pm,
		PulseAccessor:       pa,
	})

	ge = ge.NewGateway(context.Background(), insolar.CompleteNetworkState)
	ctx := context.Background()
	pa.GetLatestPulseMock.Expect(ctx).Return(*insolar.GenesisPulse, nil)

	p := packet.NewReceivedPacket(packet.NewPacket(nil, nil, types.SignCert, 1), nil)
	p.SetRequest(&packet.SignCertRequest{NodeRef: nodeRef})

	hn.BuildResponseMock.Set(func(ctx context.Context, request network.Packet, responseData interface{}) (p1 network.Packet) {
		r := packet.NewPacket(nil, nil, types.SignCert, 1)
		r.SetResponse(&packet.SignCertResponse{Sign: []byte("test_sig")})
		return r
	})
	result, err := ge.(*Complete).signCertHandler(ctx, p)

	require.NoError(t, err)
	require.Equal(t, []byte("test_sig"), result.GetResponse().GetSignCert().Sign)
}
