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
	"fmt"
	"time"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"

	"github.com/insolar/insolar/certificate"

	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"

	"github.com/insolar/insolar/insolar"
)

func newComplete(b *Base) *Complete {
	return &Complete{Base: b}
}

type Complete struct {
	*Base
}

func (g *Complete) Run(ctx context.Context) {
	g.HostNetwork.RegisterRequestHandler(types.SignCert, g.signCertHandler)
	metrics.NetworkComplete.Set(float64(time.Now().Unix()))
}

func (g *Complete) GetState() insolar.NetworkState {
	return insolar.CompleteNetworkState
}

func (g *Complete) OnPulseFromPulsar(ctx context.Context, pu insolar.Pulse, originalPacket network.ReceivedPacket) {
	// forward pulse to Consensus
	g.ConsensusPulseHandler.HandlePulse(ctx, pu, originalPacket)
}

func (g *Complete) NeedLockMessageBus() bool {
	return false
}

// ValidateCert validates node certificate
func (g *Complete) ValidateCert(ctx context.Context, certificate insolar.AuthorizationCertificate) (bool, error) {
	return g.CertificateManager.VerifyAuthorizationCertificate(certificate)
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
		return sign, nil
	}

	request := &packet.SignCertRequest{
		NodeRef: *registeredNodeRef,
	}
	future, err := g.HostNetwork.SendRequest(ctx, types.SignCert, request, *discoveryNode.GetNodeRef())
	if err != nil {
		return nil, err
	}

	p, err := future.WaitResponse(10 * time.Second)
	if err != nil {
		return nil, err
	} else if p.GetResponse().GetError() != nil {
		return nil, fmt.Errorf("[requestCertSign] Remote (%s) said %s", p.GetSender(), p.GetResponse().GetError().Error)
	}

	sign = p.GetResponse().GetSignCert().Sign

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

// signCertHandler is handler that signs certificate for some node with node own key
func (g *Complete) signCertHandler(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetSignCert() == nil {
		inslogger.FromContext(ctx).Warnf("process SignCert: got invalid request protobuf message: %s", request)
	}
	sign, err := g.signCert(ctx, &request.GetRequest().GetSignCert().NodeRef)
	if err != nil {
		return g.HostNetwork.BuildResponse(ctx, request, &packet.ErrorResponse{Error: err.Error()}), nil
	}

	return g.HostNetwork.BuildResponse(ctx, request, &packet.SignCertResponse{Sign: sign}), nil
}

func (g *Complete) ShouldIgnorePulse(context.Context, insolar.Pulse) bool {
	return false
}
