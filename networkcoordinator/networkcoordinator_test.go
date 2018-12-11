/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package networkcoordinator

import (
	"context"
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewNetworkCoordinator(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	require.Equal(t, &NetworkCoordinator{}, nc)
}

func TestNetworkCoordinator_Start(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}
	nc.MessageBus = mb
	ctx := context.Background()
	err = nc.Start(ctx)
	require.NoError(t, err)
}

func TestNetworkCoordinator_GetCoordinator_Zero(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.NoNetworkState
	}
	nc.NetworkSwitcher = ns
	crd := nc.getCoordinator()
	require.Equal(t, nc.zeroCoordinator, crd)
}

func TestNetworkCoordinator_GetCoordinator_Real(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}
	nc.NetworkSwitcher = ns
	crd := nc.getCoordinator()
	require.Equal(t, nc.realCoordinator, crd)
}

func getNode() *reply.CallMethod {
	node, err := core.MarshalArgs(struct {
		PublicKey string
		Role      core.StaticRole
	}{
		PublicKey: "test_node_public_key",
		Role:      core.StaticRoleVirtual,
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &reply.CallMethod{
		Result: []byte(node),
	}
}

func TestNetworkCoordinator_GetCertSelfDiscoveryNode(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (core.Reply, error) {
		return getNode(), nil
	}
	nc.ContractRequester = cr

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}
	nc.NetworkSwitcher = ns

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}
	nc.MessageBus = mb

	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: "test_reference",
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					PublicKey:   "test_public_key",
					Host:        "test_reference",
					NetworkSign: []byte("test_network_sign"),
				},
			},
		}
	}
	cm.NewUnsignedCertificateFunc = func(key string, role string, nodeRef string) (core.Certificate, error) {
		require.Equal(t, "test_node_public_key", key)
		require.Equal(t, "virtual", role)

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
				certificate.BootstrapNode{
					PublicKey:   "test_discovery_public_key",
					Host:        "test_discovery_host",
					NetworkSign: []byte("test_network_sign"),
					NodeRef:     "test_bootstrap_node_ref",
				},
			},
		}, nil
	}
	nc.CertificateManager = cm

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		sig := core.SignatureFromBytes([]byte("test_sig"))
		return &sig, nil
	}

	nc.CS = cs

	ctx := context.Background()
	nc.Start(ctx)

	ref := core.NewRefFromBase58("test_ref")
	certInt, err := nc.GetCert(ctx, &ref)
	require.NoError(t, err)

	cert := certInt.(*certificate.Certificate)
	require.Equal(t, "test_node_public_key", cert.PublicKey)
	require.Equal(t, "1111111111111111111111111111111111111111111111111111111111111111", cert.Reference)
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
	require.Equal(t, "test_bootstrap_node_ref", cert.BootstrapNodes[0].NodeRef)
}

func TestNetworkCoordinator_GetCertOtherDiscoveryNode(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(p context.Context, p1 *core.RecordRef, p2 string, p3 []interface{}) (core.Reply, error) {
		return getNode(), nil
	}
	nc.ContractRequester = cr

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}
	nc.NetworkSwitcher = ns

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}
	nc.MessageBus = mb

	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: testutils.RandomRef().String(),
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					PublicKey:   "test_discovery_public_key",
					Host:        "test_discovery_reference",
					NetworkSign: []byte("test_network_sign"),
				},
			},
		}
	}
	cm.NewUnsignedCertificateFunc = func(key string, role string, nodeRef string) (core.Certificate, error) {
		require.Equal(t, "test_node_public_key", key)
		require.Equal(t, "virtual", role)

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
				certificate.BootstrapNode{
					PublicKey:   "test_discovery_public_key",
					Host:        "test_discovery_host",
					NetworkSign: []byte("test_network_sign"),
					NodeRef:     "test_bootstrap_node_ref",
				},
			},
		}, nil
	}
	nc.CertificateManager = cm

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		sig := core.SignatureFromBytes([]byte("test_sig"))
		return &sig, nil
	}

	nc.CS = cs

	ctx := context.Background()
	nc.Start(ctx)

	ref := core.NewRefFromBase58("test_ref")
	certInt, err := nc.GetCert(ctx, &ref)
	require.NoError(t, err)

	cert := certInt.(*certificate.Certificate)
	require.Equal(t, "test_node_public_key", cert.PublicKey)
	require.Equal(t, "1111111111111111111111111111111111111111111111111111111111111111", cert.Reference)
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
	require.Equal(t, "test_bootstrap_node_ref", cert.BootstrapNodes[0].NodeRef)
}
