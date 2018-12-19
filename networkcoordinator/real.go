/*
 *    Copyright 2018 INS Ecosystem
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

type realNetworkCoordinator struct {
	CertificateManager core.CertificateManager
	ContractRequester  core.ContractRequester
	MessageBus         core.MessageBus
	CS                 core.CryptographyService
	PS                 core.PulseStorage
}

func newRealNetworkCoordinator(
	manager core.CertificateManager,
	requester core.ContractRequester,
	msgBus core.MessageBus,
	cs core.CryptographyService,
	ps core.PulseStorage,
) *realNetworkCoordinator {
	return &realNetworkCoordinator{
		CertificateManager: manager,
		ContractRequester:  requester,
		MessageBus:         msgBus,
		CS:                 cs,
		PS:                 ps,
	}
}

// GetCert method generates cert by requesting signs from discovery nodes
func (rnc *realNetworkCoordinator) GetCert(ctx context.Context, registeredNodeRef *core.RecordRef) (core.Certificate, error) {
	pKey, role, err := rnc.getNodeInfo(ctx, registeredNodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't get node info")
	}

	currentNodeCert := rnc.CertificateManager.GetCertificate()
	registeredNodeCert, err := rnc.CertificateManager.NewUnsignedCertificate(pKey, role, registeredNodeRef.String())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't create certificate")
	}

	for i, discoveryNode := range currentNodeCert.GetDiscoveryNodes() {
		sign, err := rnc.requestCertSign(ctx, discoveryNode, registeredNodeRef)
		if err != nil {
			return nil, errors.Wrap(err, "[ GetCert ] Couldn't request cert sign")
		}
		registeredNodeCert.(*certificate.Certificate).BootstrapNodes[i].NodeSign = sign
	}
	return registeredNodeCert, nil
}

// requestCertSign method requests sign from single discovery node
func (rnc *realNetworkCoordinator) requestCertSign(ctx context.Context, discoveryNode core.DiscoveryNode, registeredNodeRef *core.RecordRef) ([]byte, error) {
	var sign []byte
	var err error

	currentNodeCert := rnc.CertificateManager.GetCertificate()

	if *discoveryNode.GetNodeRef() == *currentNodeCert.GetNodeRef() {
		sign, err = rnc.signCert(ctx, registeredNodeRef)
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
		currentPulse, err := rnc.PS.Current(ctx)
		if err != nil {
			return nil, err
		}
		r, err := rnc.MessageBus.Send(ctx, msg, *currentPulse, opts)
		if err != nil {
			return nil, err
		}
		sign = r.(reply.NodeSignInt).GetSign()
	}

	return sign, nil
}

// signCertHandler is MsgBus handler that signs certificate for some node with node own key
func (rnc *realNetworkCoordinator) signCertHandler(ctx context.Context, p core.Parcel) (core.Reply, error) {
	nodeRef := p.Message().(message.NodeSignPayloadInt).GetNodeRef()
	sign, err := rnc.signCert(ctx, nodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't extract response")
	}
	return &reply.NodeSign{
		Sign: sign,
	}, nil
}

// signCert returns certificate sign fore node
func (rnc *realNetworkCoordinator) signCert(ctx context.Context, registeredNodeRef *core.RecordRef) ([]byte, error) {
	pKey, role, err := rnc.getNodeInfo(ctx, registeredNodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't extract response")
	}

	data := []byte(pKey + registeredNodeRef.String() + role)
	sign, err := rnc.CS.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't sign")
	}

	return sign.Bytes(), nil
}

// getNodeInfo request info from ledger
func (rnc *realNetworkCoordinator) getNodeInfo(ctx context.Context, nodeRef *core.RecordRef) (string, string, error) {
	res, err := rnc.ContractRequester.SendRequest(ctx, nodeRef, "GetNodeInfo", []interface{}{})
	if err != nil {
		return "", "", errors.Wrap(err, "[ GetCert ] Couldn't call GetNodeInfo")
	}
	pKey, role, err := extractor.NodeInfoResponse(res.(*reply.CallMethod).Result)
	if err != nil {
		return "", "", errors.Wrap(err, "[ GetCert ] Couldn't extract response")
	}
	return pKey, role, nil
}

// TODO: remove and use SetPulse instead
func (rnc *realNetworkCoordinator) WriteActiveNodes(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) error {
	return errors.New("not implemented")
}

// SetPulse uses PulseManager component for saving pulse info
func (rnc *realNetworkCoordinator) SetPulse(ctx context.Context, pulse core.Pulse) error {
	return errors.New("not implemented")
}
