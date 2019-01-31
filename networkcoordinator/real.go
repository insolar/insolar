/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
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
}

func newRealNetworkCoordinator(
	manager core.CertificateManager,
	requester core.ContractRequester,
	msgBus core.MessageBus,
	cs core.CryptographyService,
) *realNetworkCoordinator {
	return &realNetworkCoordinator{
		CertificateManager: manager,
		ContractRequester:  requester,
		MessageBus:         msgBus,
		CS:                 cs,
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
		r, err := rnc.MessageBus.Send(ctx, msg, opts)
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

// SetPulse uses PulseManager component for saving pulse info
func (rnc *realNetworkCoordinator) SetPulse(ctx context.Context, pulse core.Pulse) error {
	return errors.New("not implemented")
}
