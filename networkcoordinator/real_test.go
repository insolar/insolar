package networkcoordinator

import (
	"context"
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestRealNetworkCoordinator_New(t *testing.T) {
	coord := newRealNetworkCoordinator(nil, nil, nil, nil, nil)
	require.Equal(t, &realNetworkCoordinator{}, coord)
}

func TestRealNetworkCoordinator_GetCert(t *testing.T) {
	nodeRef := testutils.RandomRef()
	certNodeRef := testutils.RandomRef().String()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}

	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef,
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     certNodeRef,
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
					NodeRef:     certNodeRef,
				},
			},
		}, nil
	}

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		sig := core.SignatureFromBytes([]byte("test_sig"))
		return &sig, nil
	}

	coord := newRealNetworkCoordinator(cm, cr, mb, cs, nil)
	ctx := context.Background()
	result, err := coord.GetCert(ctx, &nodeRef)
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
	require.Equal(t, certNodeRef, cert.BootstrapNodes[0].NodeRef)
}

func TestRealNetworkCoordinator_GetCert_getNodeInfoError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return nil, errors.New("test_error")
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, nil, nil)
	ctx := context.Background()
	_, err := coord.GetCert(ctx, &nodeRef)
	require.EqualError(t, err, "[ GetCert ] Couldn't get node info: [ GetCert ] Couldn't call GetNodeInfo: test_error")
}

func TestRealNetworkCoordinator_GetCert_DeserializeError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return &reply.CallMethod{
			Result: []byte(""),
		}, nil
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, nil, nil)
	ctx := context.Background()
	_, err := coord.GetCert(ctx, &nodeRef)
	require.EqualError(t, err, "[ GetCert ] Couldn't get node info: [ GetCert ] Couldn't extract response: [ NodeInfoResponse ] Can't unmarshal response: [ UnMarshalResponse ]: [ Deserialize ]: EOF")
}

func TestRealNetworkCoordinator_GetCert_UnsignedCertificateError(t *testing.T) {
	nodeRef := testutils.RandomRef()
	certNodeRef := testutils.RandomRef().String()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}

	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef,
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     certNodeRef,
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

		return nil, errors.New("test_error")
	}
	coord := newRealNetworkCoordinator(cm, cr, nil, nil, nil)
	ctx := context.Background()
	_, err := coord.GetCert(ctx, &nodeRef)
	require.EqualError(t, err, "[ GetCert ] Couldn't create certificate: test_error")
}

func TestRealNetworkCoordinator_GetCert_SignCertError(t *testing.T) {
	nodeRef := testutils.RandomRef()
	certNodeRef := testutils.RandomRef().String()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}

	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef,
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     certNodeRef,
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

		return nil, nil
	}

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		return nil, errors.New("test_error")
	}

	coord := newRealNetworkCoordinator(cm, cr, nil, cs, nil)
	ctx := context.Background()
	_, err := coord.GetCert(ctx, &nodeRef)
	require.EqualError(t, err, "[ GetCert ] Couldn't request cert sign: [ SignCert ] Couldn't sign: test_error")
}

func TestRealNetworkCoordinator_requestCertSignSelfDiscoveryNode(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}
	mb.SendFunc = func(p context.Context, p1 core.Message, p2 core.Pulse, p3 *core.MessageSendOptions) (core.Reply, error) {
		return &reply.NodeSign{
			Sign: []byte("test_sig"),
		}, nil
	}

	cm := testutils.NewCertificateManagerMock(t)
	certNodeRef := testutils.RandomRef().String()
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef,
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     certNodeRef,
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
					NodeRef:     certNodeRef,
				},
			},
		}, nil
	}

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		sig := core.SignatureFromBytes([]byte("test_sig"))
		return &sig, nil
	}

	coord := newRealNetworkCoordinator(cm, cr, mb, cs, nil)
	ctx := context.Background()
	dNode := certificate.BootstrapNode{
		PublicKey:   "test_discovery_public_key",
		Host:        "test_discovery_host",
		NetworkSign: []byte("test_network_sign"),
		NodeRef:     certNodeRef,
	}
	result, err := coord.requestCertSign(ctx, &dNode, &nodeRef)
	require.NoError(t, err)
	require.Equal(t, []byte("test_sig"), result)
}

func TestRealNetworkCoordinator_requestCertSignOtherDiscoveryNode(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}
	mb.SendFunc = func(p context.Context, p1 core.Message, p2 core.Pulse, p3 *core.MessageSendOptions) (core.Reply, error) {
		return &reply.NodeSign{
			Sign: []byte("test_sig"),
		}, nil
	}

	cm := testutils.NewCertificateManagerMock(t)
	certNodeRef := testutils.RandomRef().String()
	discoveryNodeRef := testutils.RandomRef().String()
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef,
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     discoveryNodeRef,
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
					NodeRef:     discoveryNodeRef,
				},
			},
		}, nil
	}

	pm := testutils.NewPulseManagerMock(t)
	pm.CurrentFunc = func(ctx context.Context) (*core.Pulse, error) {
		return &core.Pulse{}, nil
	}

	coord := newRealNetworkCoordinator(cm, cr, mb, nil, pm)
	ctx := context.Background()
	dNode := certificate.BootstrapNode{
		PublicKey:   "test_discovery_public_key",
		Host:        "test_discovery_host",
		NetworkSign: []byte("test_network_sign"),
		NodeRef:     discoveryNodeRef,
	}
	result, err := coord.requestCertSign(ctx, &dNode, &nodeRef)
	require.NoError(t, err)
	require.Equal(t, []byte("test_sig"), result)
}

func TestRealNetworkCoordinator_requestCertSignSelfDiscoveryNode_signCertError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return nil, errors.New("test_error")
	}

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}

	cm := testutils.NewCertificateManagerMock(t)
	certNodeRef := testutils.RandomRef().String()
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef,
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     certNodeRef,
					PublicKey:   "test_discovery_public_key",
					Host:        "test_discovery_host",
					NetworkSign: []byte("test_network_sign"),
				},
			},
		}
	}
	coord := newRealNetworkCoordinator(cm, cr, nil, nil, nil)
	ctx := context.Background()
	dNode := certificate.BootstrapNode{
		PublicKey:   "test_discovery_public_key",
		Host:        "test_discovery_host",
		NetworkSign: []byte("test_network_sign"),
		NodeRef:     certNodeRef,
	}
	_, err := coord.requestCertSign(ctx, &dNode, &nodeRef)
	require.EqualError(t, err, "[ SignCert ] Couldn't extract response: [ GetCert ] Couldn't call GetNodeInfo: test_error")
}

func TestRealNetworkCoordinator_requestCertSignOtherDiscoveryNode_CurrentPulseError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}

	cm := testutils.NewCertificateManagerMock(t)
	certNodeRef := testutils.RandomRef().String()
	discoveryNodeRef := testutils.RandomRef().String()
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef,
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     discoveryNodeRef,
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
					NodeRef:     discoveryNodeRef,
				},
			},
		}, nil
	}

	pm := testutils.NewPulseManagerMock(t)
	pm.CurrentFunc = func(ctx context.Context) (*core.Pulse, error) {
		return nil, errors.New("test_error")
	}

	coord := newRealNetworkCoordinator(cm, cr, mb, nil, pm)
	ctx := context.Background()
	dNode := certificate.BootstrapNode{
		PublicKey:   "test_discovery_public_key",
		Host:        "test_discovery_host",
		NetworkSign: []byte("test_network_sign"),
		NodeRef:     discoveryNodeRef,
	}
	_, err := coord.requestCertSign(ctx, &dNode, &nodeRef)
	require.EqualError(t, err, "test_error")
}

func TestRealNetworkCoordinator_requestCertSignOtherDiscoveryNode_SendError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}
	mb.SendFunc = func(p context.Context, p1 core.Message, p2 core.Pulse, p3 *core.MessageSendOptions) (core.Reply, error) {
		return nil, errors.New("test_error")
	}

	cm := testutils.NewCertificateManagerMock(t)
	certNodeRef := testutils.RandomRef().String()
	discoveryNodeRef := testutils.RandomRef().String()
	cm.GetCertificateFunc = func() core.Certificate {
		return &certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: "test_public_key",
				Reference: certNodeRef,
				Role:      "virtual",
			},
			MajorityRule: 0,
			BootstrapNodes: []certificate.BootstrapNode{
				certificate.BootstrapNode{
					NodeRef:     discoveryNodeRef,
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
					NodeRef:     discoveryNodeRef,
				},
			},
		}, nil
	}

	pm := testutils.NewPulseManagerMock(t)
	pm.CurrentFunc = func(ctx context.Context) (*core.Pulse, error) {
		return &core.Pulse{}, nil
	}

	coord := newRealNetworkCoordinator(cm, cr, mb, nil, pm)
	ctx := context.Background()
	dNode := certificate.BootstrapNode{
		PublicKey:   "test_discovery_public_key",
		Host:        "test_discovery_host",
		NetworkSign: []byte("test_network_sign"),
		NodeRef:     discoveryNodeRef,
	}
	_, err := coord.requestCertSign(ctx, &dNode, &nodeRef)
	require.EqualError(t, err, "test_error")
}

func TestRealNetworkCoordinator_signCertHandler(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		sig := core.SignatureFromBytes([]byte("test_sig"))
		return &sig, nil
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, cs, nil)
	ctx := context.Background()
	result, err := coord.signCertHandler(ctx, &message.Parcel{Msg: &message.NodeSignPayload{NodeRef: &nodeRef}})
	require.NoError(t, err)
	require.Equal(t, []byte("test_sig"), result.(*reply.NodeSign).Sign)
}

func TestRealNetworkCoordinator_signCertHandler_NodeInfoError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return nil, errors.New("test_error")
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, nil, nil)
	ctx := context.Background()
	_, err := coord.signCertHandler(ctx, &message.Parcel{Msg: &message.NodeSignPayload{NodeRef: &nodeRef}})
	require.EqualError(t, err, "[ SignCert ] Couldn't extract response: [ SignCert ] Couldn't extract response: [ GetCert ] Couldn't call GetNodeInfo: test_error")
}

func TestRealNetworkCoordinator_signCertHandler_SignError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		return nil, errors.New("test_error")
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, cs, nil)
	ctx := context.Background()
	_, err := coord.signCertHandler(ctx, &message.Parcel{Msg: &message.NodeSignPayload{NodeRef: &nodeRef}})
	require.EqualError(t, err, "[ SignCert ] Couldn't extract response: [ SignCert ] Couldn't sign: test_error")
}

func TestRealNetworkCoordinator_signCert(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		sig := core.SignatureFromBytes([]byte("test_sig"))
		return &sig, nil
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, cs, nil)
	ctx := context.Background()
	result, err := coord.signCert(ctx, &nodeRef)
	require.NoError(t, err)
	require.Equal(t, []byte("test_sig"), result)
}

func TestRealNetworkCoordinator_signCert_NodeInfoError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return nil, errors.New("test_error")
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, nil, nil)
	ctx := context.Background()
	_, err := coord.signCert(ctx, &nodeRef)
	require.EqualError(t, err, "[ SignCert ] Couldn't extract response: [ GetCert ] Couldn't call GetNodeInfo: test_error")
}

func TestRealNetworkCoordinator_signCert_SignError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignFunc = func(data []byte) (*core.Signature, error) {
		return nil, errors.New("test_error")
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, cs, nil)
	ctx := context.Background()
	_, err := coord.signCert(ctx, &nodeRef)
	require.EqualError(t, err, "[ SignCert ] Couldn't sign: test_error")
}

func TestRealNetworkCoordinator_getNodeInfo(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return test_getNode(), nil
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, nil, nil)
	ctx := context.Background()
	key, role, err := coord.getNodeInfo(ctx, &nodeRef)
	require.NoError(t, err)
	require.Equal(t, "test_node_public_key", key)
	require.Equal(t, "virtual", role)
}

func TestRealNetworkCoordinator_getNodeInfo_SendRequestError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return nil, errors.New("test_error")
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, nil, nil)
	ctx := context.Background()
	_, _, err := coord.getNodeInfo(ctx, &nodeRef)
	require.EqualError(t, err, "[ GetCert ] Couldn't call GetNodeInfo: test_error")
}

func TestRealNetworkCoordinator_getNodeInfo_ExtractError(t *testing.T) {
	nodeRef := testutils.RandomRef()

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(ctx context.Context, ref *core.RecordRef, method string, args []interface{}) (core.Reply, error) {
		require.Equal(t, nodeRef, *ref)
		require.Equal(t, "GetNodeInfo", method)
		require.Equal(t, 0, len(args))
		return &reply.CallMethod{
			Result: []byte(""),
		}, nil
	}

	coord := newRealNetworkCoordinator(nil, cr, nil, nil, nil)
	ctx := context.Background()
	_, _, err := coord.getNodeInfo(ctx, &nodeRef)
	require.EqualError(t, err, "[ GetCert ] Couldn't extract response: [ NodeInfoResponse ] Can't unmarshal response: [ UnMarshalResponse ]: [ Deserialize ]: EOF")
}
