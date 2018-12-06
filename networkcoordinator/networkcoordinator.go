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
	Bus                 core.MessageBus          `inject:""`
	CS                  core.CryptographyService `inject:""`

	realCoordinator core.NetworkCoordinator
	zeroCoordinator core.NetworkCoordinator
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
	return nc.Bus.Register(core.NetworkCoordinatorNodeSignRequest, nc.SignNode)
}

func (nc *NetworkCoordinator) getCoordinator() core.NetworkCoordinator {
	if nc.NetworkSwitcher.GetState() == core.CompleteNetworkState {
		return nc.realCoordinator
	}
	return nc.zeroCoordinator
}

// GetCert method returns node certificate
func (nc *NetworkCoordinator) GetCert(ctx context.Context, nodeRef core.RecordRef) (core.Certificate, error) {
	res, err := nc.ContractRequester.SendRequest(ctx, &nodeRef, "GetNodeInfo", []interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't call GetNodeInfo")
	}
	pKey, role, err := extractor.NodeInfoResponse(res.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't extract response")
	}

	currentNodeCert := nc.CertificateManager.GetCertificate()
	cert, err := nc.CertificateManager.NewUnsignedCertificate(pKey, role, nodeRef.String())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't create certificate")
	}

	for i, node := range currentNodeCert.GetDiscoveryNodes() {
		if node.GetNodeRef() == currentNodeCert.GetNodeRef() {
			sign, err := nc.signNode(ctx, node.GetNodeRef())
			if err != nil {
				return nil, err
			}
			currentNodeCert.(*certificate.Certificate).BootstrapNodes[i].NodeSign = sign
		} else {
			msg := message.NodeSignPayload{
				NodeRef: &nodeRef,
			}
			opts := core.MessageSendOptions{
				Receiver: node.GetNodeRef(),
			}
			r, err := nc.Bus.Send(ctx, &msg, &opts)
			if err != nil {
				return nil, err
			}
			sign := r.(reply.NodeSignInt).GetSign()
			currentNodeCert.(*certificate.Certificate).BootstrapNodes[i].NodeSign = sign
		}
	}
	return cert, nil
}

// ValidateCert validates node certificate
func (nc *NetworkCoordinator) ValidateCert(ctx context.Context, certificate core.AuthorizationCertificate) (bool, error) {
	return nc.CertificateManager.VerifyAuthorizationCertificate(certificate)
}

// WriteActiveNodes writes active nodes to ledger
func (nc *NetworkCoordinator) WriteActiveNodes(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) error {
	return nc.getCoordinator().WriteActiveNodes(ctx, number, activeNodes)
}

// SetPulse writes pulse data on local storage
func (nc *NetworkCoordinator) SetPulse(ctx context.Context, pulse core.Pulse) error {
	return nc.getCoordinator().SetPulse(ctx, pulse)
}

// SignNode signs info about some node
func (nc *NetworkCoordinator) SignNode(ctx context.Context, p core.Parcel) (core.Reply, error) {
	nodeRef := p.Message().(message.NodeSignPayloadInt).GetNodeRef()
	sign, err := nc.signNode(ctx, nodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignNode ] Couldn't extract response")
	}
	return &reply.NodeSign{
		Sign: sign,
	}, nil

}

func (nc *NetworkCoordinator) signNode(ctx context.Context, nodeRef *core.RecordRef) ([]byte, error) {
	res, err := nc.ContractRequester.SendRequest(ctx, nodeRef, "GetNodeInfo", []interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "[ SignNode ] Couldn't call GetNodeInfo")
	}
	pKey, role, err := extractor.NodeInfoResponse(res.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignNode ] Couldn't extract response")
	}

	data := []byte(pKey + nodeRef.String() + role)
	sign, err := nc.CS.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignNode ] Couldn't sign")
	}

	return sign.Bytes(), nil
}
