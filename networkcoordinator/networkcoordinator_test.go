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

func TestNetworkCoordinator_GetCert(t *testing.T) {
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
					PublicKey:   "test_discovery_public_key",
					Host:        "test_discovery_host",
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
				PublicKey: "test_public_key",
				Reference: "test_reference",
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					PublicKey:   "test_discovery_public_key",
					Host:        "test_discovery_host",
					NetworkSign: []byte("test_network_sign"),
				},
			},
		}, nil
	}
	nc.CertificateManager = cm

	ctx := context.Background()
	nc.Start(ctx)
	ref := core.NewRefFromBase58("test_ref")
	nc.GetCert(ctx, &ref)

}
