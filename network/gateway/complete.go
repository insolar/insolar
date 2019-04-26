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
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/certificate"

	"github.com/insolar/insolar/insolar/message"

	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"

	"github.com/insolar/insolar/insolar"
)

func NewComplete(b *Base) *Complete {
	return &Complete{Base: b}
}

type Complete struct {
	*Base
}

func (g *Complete) Run(ctx context.Context) {
	g.GIL.Release(ctx)
	g.MessageBus.MustRegister(insolar.TypeNodeSignRequest, g.signCertHandler)
	metrics.NetworkComplete.Set(float64(time.Now().Unix()))
}

func (g *Complete) GetState() insolar.NetworkState {
	return insolar.CompleteNetworkState
}

func (g *Complete) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	inslogger.FromContext(ctx).Debugf("Gateway.Complete: pulse happens %d", pu.PulseNumber)
	return nil
}

// GetCert method generates cert by requesting signs from discovery nodes
func (g *Complete) GetCert(ctx context.Context, registeredNodeRef *insolar.Reference) (insolar.Certificate, error) {
	pKey, role, err := g.getNodeInfo(ctx, registeredNodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't get node info")
	}

	currentNodeCert := g.CertificateManager.GetCertificate()
	registeredNodeCert, err := g.CertificateManager.NewUnsignedCertificate(pKey, role, registeredNodeRef.String())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't create certificate")
	}

	for i, discoveryNode := range currentNodeCert.GetDiscoveryNodes() {
		sign, err := g.requestCertSign(ctx, discoveryNode, registeredNodeRef)
		if err != nil {
			return nil, errors.Wrap(err, "[ GetCert ] Couldn't request cert sign")
		}
		registeredNodeCert.(*certificate.Certificate).BootstrapNodes[i].NodeSign = sign
	}
	return registeredNodeCert, nil
}

// requestCertSign method requests sign from single discovery node
func (g *Complete) requestCertSign(ctx context.Context, discoveryNode insolar.DiscoveryNode, registeredNodeRef *insolar.Reference) ([]byte, error) {
	var sign []byte
	var err error

	currentNodeCert := g.CertificateManager.GetCertificate()

	if *discoveryNode.GetNodeRef() == *currentNodeCert.GetNodeRef() {
		sign, err = g.signCert(ctx, registeredNodeRef)
		if err != nil {
			return nil, err
		}
	} else {
		msg := &message.NodeSignPayload{
			NodeRef: registeredNodeRef,
		}
		opts := &insolar.MessageSendOptions{
			Receiver: discoveryNode.GetNodeRef(),
		}
		r, err := g.MessageBus.Send(ctx, msg, opts)
		if err != nil {
			return nil, err
		}
		sign = r.(reply.NodeSignInt).GetSign()
	}

	return sign, nil
}

func (g *Complete) getNodeInfo(ctx context.Context, nodeRef *insolar.Reference) (string, string, error) {
	res, err := g.ContractRequester.SendRequest(ctx, nodeRef, "GetNodeInfo", []interface{}{})
	if err != nil {
		return "", "", errors.Wrap(err, "[ GetCert ] Couldn't call GetNodeInfo")
	}
	pKey, role, err := extractor.NodeInfoResponse(res.(*reply.CallMethod).Result)
	if err != nil {
		return "", "", errors.Wrap(err, "[ GetCert ] Couldn't extract response")
	}
	return pKey, role, nil
}

// signCert returns certificate sign fore node
func (g *Complete) signCert(ctx context.Context, registeredNodeRef *insolar.Reference) ([]byte, error) {
	pKey, role, err := g.getNodeInfo(ctx, registeredNodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't extract response")
	}

	data := []byte(pKey + registeredNodeRef.String() + role)
	sign, err := g.CryptographyService.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't sign")
	}

	return sign.Bytes(), nil
}

// signCertHandler is MsgBus handler that signs certificate for some node with node own key
func (g *Complete) signCertHandler(ctx context.Context, p insolar.Parcel) (insolar.Reply, error) {
	nodeRef := p.Message().(message.NodeSignPayloadInt).GetNodeRef()
	sign, err := g.signCert(ctx, nodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't extract response")
	}
	return &reply.NodeSign{
		Sign: sign,
	}, nil
}

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
