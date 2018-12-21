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
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewNetworkCoordinator(t *testing.T) {
	certificateManager := testutils.NewCertificateManagerMock(t)
	networkSwitcher := testutils.NewNetworkSwitcherMock(t)
	contractRequester := testutils.NewContractRequesterMock(t)
	messageBus := testutils.NewMessageBusMock(t)
	cs := testutils.NewCryptographyServiceMock(t)
	ps := testutils.NewPulseStorageMock(t)

	nc, err := New()
	require.NoError(t, err)
	require.Equal(t, &NetworkCoordinator{}, nc)

	cm := &component.Manager{}
	cm.Inject(certificateManager, networkSwitcher, contractRequester, messageBus, cs, ps, nc)
	require.Equal(t, certificateManager, nc.CertificateManager)
	require.Equal(t, networkSwitcher, nc.NetworkSwitcher)
	require.Equal(t, contractRequester, nc.ContractRequester)
	require.Equal(t, messageBus, nc.MessageBus)
	require.Equal(t, cs, nc.CS)
	require.Equal(t, ps, nc.PS)
}

func TestNetworkCoordinator_Start(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	nc.MessageBus = mockMessageBus(t, true, nil, nil)
	ctx := context.Background()
	err = nc.Start(ctx)
	require.NoError(t, err)
	require.NotNil(t, nc.realCoordinator)
	require.NotNil(t, nc.zeroCoordinator)
}

func TestNetworkCoordinator_GetCoordinator_Zero(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.NoNetworkState
	}
	nc.NetworkSwitcher = ns
	nc.MessageBus = mockMessageBus(t, true, nil, nil)
	ctx := context.Background()
	nc.Start(ctx)
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
	nc.MessageBus = mockMessageBus(t, true, nil, nil)
	ctx := context.Background()
	nc.Start(ctx)
	crd := nc.getCoordinator()
	require.Equal(t, nc.realCoordinator, crd)
}

func mockReply(t *testing.T) []byte {
	node, err := core.MarshalArgs(struct {
		PublicKey string
		Role      core.StaticRole
	}{
		PublicKey: "test_node_public_key",
		Role:      core.StaticRoleVirtual,
	}, nil)
	require.NoError(t, err)
	return []byte(node)
}

func mockMessageBus(t *testing.T, ok bool, ref *core.RecordRef, discovery *core.RecordRef) *testutils.MessageBusMock {
	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.TypeNodeSignRequest)
	}
	mb.SendFunc = func(p context.Context, msg core.Message, options *core.MessageSendOptions) (core.Reply, error) {
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

func mockCertificateManager(t *testing.T, certNodeRef *core.RecordRef, discoveryNodeRef *core.RecordRef, unsignCertOk bool) *testutils.CertificateManagerMock {
	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef.String(),
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     discoveryNodeRef.String(),
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
					certificate.BootstrapNode{
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

func mockCryptographyService(t *testing.T, ok bool) core.CryptographyService {
	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		if ok {
			sig := core.SignatureFromBytes([]byte("test_sig"))
			return &sig, nil
		}
		return nil, errors.New("test_error")
	}
	return cs
}
