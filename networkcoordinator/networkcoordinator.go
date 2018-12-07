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

	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/pkg/errors"
)

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator struct {
	CertificateManager  core.CertificateManager  `inject:""`
	NetworkSwitcher     core.NetworkSwitcher     `inject:""`
	ContractRequester   core.ContractRequester   `inject:""`
	GenesisDataProvider core.GenesisDataProvider `inject:""`
	MessageBus          core.MessageBus          `inject:""`
	CS                  core.CryptographyService `inject:""`

	realCoordinator Coordinator
	zeroCoordinator Coordinator
}

// New creates new NetworkCoordinator
func New() (*NetworkCoordinator, error) {
	return &NetworkCoordinator{}, nil
}

// Init implements interface of Component
func (nc *NetworkCoordinator) Init(ctx context.Context) error {
	nc.zeroCoordinator = newZeroNetworkCoordinator()
	nc.realCoordinator = newRealNetworkCoordinator()
	return nil
}

// Start implements interface of Component
func (nc *NetworkCoordinator) Start(ctx context.Context) error {
	nc.MessageBus.MustRegister(core.NetworkCoordinatorNodeSignRequest, nc.signCertHandler)
	return nil
}

func (nc *NetworkCoordinator) getCoordinator() core.NetworkCoordinator {
	if nc.NetworkSwitcher.GetState() == core.CompleteNetworkState {
		return nc.realCoordinator
	}
	return nc.zeroCoordinator
}

// GetCert method returns node certificate by requesting sign from discovery nodes
func (nc *NetworkCoordinator) GetCert(ctx context.Context, registeredNodeRef *core.RecordRef) (core.Certificate, error) {
	pKey, role, err := nc.getNodeInfo(ctx, registeredNodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't get node info")
	}

	currentNodeCert := nc.CertificateManager.GetCertificate()
	registeredNodeCert, err := nc.CertificateManager.NewUnsignedCertificate(pKey, role, registeredNodeRef.String())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't create certificate")
	}

	for i, discoveryNode := range currentNodeCert.GetDiscoveryNodes() {
		sign, err := nc.requestCertSign(ctx, discoveryNode, registeredNodeRef)
		if err != nil {
			return nil, errors.Wrap(err, "[ GetCert ] Couldn't request cert sign")
		}
		registeredNodeCert.(*certificate.Certificate).BootstrapNodes[i].NodeSign = sign
	}
	return registeredNodeCert, nil
}

func (nc *NetworkCoordinator) requestCertSign(ctx context.Context, discoveryNode core.DiscoveryNode, registeredNodeRef *core.RecordRef) ([]byte, error) {
	var sign []byte
	var err error

	currentNodeCert := nc.CertificateManager.GetCertificate()

	if discoveryNode.GetNodeRef() == currentNodeCert.GetNodeRef() {
		sign, err = nc.signCert(ctx, registeredNodeRef)
		if err != nil {
			return nil, err
		}
	} else {
		msg := &message.NodeSignPayload{
			NodeRef: registeredNodeRef,
		}
		opts := &core.MessageSendOptions{
			Receiver: discoveryNode.GetNodeRef(),
		}
		r, err := nc.MessageBus.Send(ctx, msg, opts)
		if err != nil {
			return nil, err
		}
		sign = r.(reply.NodeSignInt).GetSign()
	}

	return sign, nil
}

// ValidateCert validates node certificate
func (nc *NetworkCoordinator) ValidateCert(ctx context.Context, certificate core.AuthorizationCertificate) (bool, error) {
	return nc.CertificateManager.VerifyAuthorizationCertificate(certificate)
}

// signCertHandler is MsgBus handler that signs certificate for some node with node own key
func (nc *NetworkCoordinator) signCertHandler(ctx context.Context, p core.Parcel) (core.Reply, error) {
	nodeRef := p.Message().(message.NodeSignPayloadInt).GetNodeRef()
	sign, err := nc.signCert(ctx, nodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't extract response")
	}
	return &reply.NodeSign{
		Sign: sign,
	}, nil
}

// signCert is MsgBus handler that signs certificate for some node with node own key
// func (nc *NetworkCoordinator) signCert(ctx context.Context, p core.Parcel) (core.Reply, error) {
// 	nodeRef := p.Message().(message.NodeSignPayloadInt).GetNodeRef()
// 	sign, err := nc.signCert(ctx, nodeRef)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "[ SignCert ] Couldn't extract response")
// 	}
// 	return &reply.NodeSign{
// 		Sign: sign,
// 	}, nil
// }

func (nc *NetworkCoordinator) signCert(ctx context.Context, registeredNodeRef *core.RecordRef) ([]byte, error) {
	pKey, role, err := nc.getNodeInfo(ctx, registeredNodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't extract response")
	}

	data := []byte(pKey + registeredNodeRef.String() + role)
	sign, err := nc.CS.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't sign")
	}

	return sign.Bytes(), nil
}

func (nc *NetworkCoordinator) getNodeInfo(ctx context.Context, nodeRef *core.RecordRef) (string, string, error) {
	res, err := nc.ContractRequester.SendRequest(ctx, nodeRef, "GetNodeInfo", []interface{}{})
	if err != nil {
		return "", "", errors.Wrap(err, "[ GetCert ] Couldn't call GetNodeInfo")
	}
	pKey, role, err := extractor.NodeInfoResponse(res.(*reply.CallMethod).Result)
	if err != nil {
		return "", "", errors.Wrap(err, "[ GetCert ] Couldn't extract response")
	}
	return pKey, role, nil
}

// WriteActiveNodes writes active nodes to ledger
func (nc *NetworkCoordinator) WriteActiveNodes(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) error {
	return nc.getCoordinator().WriteActiveNodes(ctx, number, activeNodes)
}

// SetPulse writes pulse data on local storage
func (nc *NetworkCoordinator) SetPulse(ctx context.Context, pulse core.Pulse) error {
	return nc.getCoordinator().SetPulse(ctx, pulse)
}
