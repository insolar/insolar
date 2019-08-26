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

	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/rules"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/certificate"

	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"

	"github.com/insolar/insolar/insolar"
)

func newComplete(b *Base) *Complete {
	return &Complete{
		Base: b,
	}
}

type Complete struct {
	*Base
}

func (g *Complete) Run(ctx context.Context, pulse insolar.Pulse) {
	if g.bootstrapTimer != nil {
		g.bootstrapTimer.Stop()
	}

	g.HostNetwork.RegisterRequestHandler(types.SignCert, g.signCertHandler)
	metrics.NetworkComplete.Set(float64(time.Now().Unix()))
}

func (g *Complete) GetState() insolar.NetworkState {
	return insolar.CompleteNetworkState
}

func (g *Complete) BeforeRun(ctx context.Context, pulse insolar.Pulse) {
	err := g.PulseManager.Set(ctx, pulse)
	if err != nil {
		inslogger.FromContext(ctx).Panicf("failed to set start pulse: %d, %s", pulse.PulseNumber, err.Error())
	}
}

// GetCert method generates cert by requesting signs from discovery nodes
func (g *Complete) GetCert(ctx context.Context, registeredNodeRef *insolar.Reference) (insolar.Certificate, error) {
	pKey, role, err := g.getNodeInfo(ctx, registeredNodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCert ] Couldn't get node info")
	}

	currentNodeCert := g.CertificateManager.GetCertificate()
	registeredNodeCert, err := certificate.NewUnsignedCertificate(currentNodeCert, pKey, role, registeredNodeRef.String())
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
	currentNodeCert := g.CertificateManager.GetCertificate()

	if *discoveryNode.GetNodeRef() == *currentNodeCert.GetNodeRef() {
		sign, err := g.signCert(ctx, registeredNodeRef)
		if err != nil {
			return nil, err
		}
		return sign.Bytes(), nil
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

	return p.GetResponse().GetSignCert().Sign, nil
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

func (g *Complete) signCert(ctx context.Context, registeredNodeRef *insolar.Reference) (*insolar.Signature, error) {
	pKey, role, err := g.getNodeInfo(ctx, registeredNodeRef)
	if err != nil {
		return nil, errors.Wrap(err, "[ SignCert ] Couldn't extract response")
	}
	return certificate.SignCert(g.CryptographyService, pKey, role, registeredNodeRef.String())
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

	return g.HostNetwork.BuildResponse(ctx, request, &packet.SignCertResponse{Sign: sign.Bytes()}), nil
}

func (g *Complete) EphemeralMode(nodes []insolar.NetworkNode) bool {
	return false
}

func (g *Complete) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	workingNodes := node.Select(nodes, node.ListWorking)

	if ok, _ := rules.CheckMajorityRule(g.CertificateManager.GetCertificate(), workingNodes); !ok {
		g.Gatewayer.FailState(ctx, "MajorityRule failed")
	}

	if !rules.CheckMinRole(g.CertificateManager.GetCertificate(), workingNodes) {
		g.Gatewayer.FailState(ctx, "MinRoles failed")
	}

	g.Base.UpdateState(ctx, pulseNumber, nodes, cloudStateHash)
}

func (g *Complete) OnPulseFromConsensus(ctx context.Context, pulse insolar.Pulse) {
	g.Base.OnPulseFromConsensus(ctx, pulse)

	done := make(chan struct{})
	defer close(done)
	pulseProcessingWatchdog(ctx, pulse, done)

	logger := inslogger.FromContext(ctx)

	logger.Infof("Got new pulse number: %d", pulse.PulseNumber)
	ctx, span := instracer.StartSpan(ctx, "ServiceNetwork.Handlepulse")
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(pulse.PulseNumber)),
	)
	defer span.End()

	err := g.PulseManager.Set(ctx, pulse)
	if err != nil {
		logger.Fatalf("Failed to set new pulse: %s", err.Error())
	}
	logger.Infof("Set new current pulse number: %d", pulse.PulseNumber)
	stats.Record(ctx, statPulse.M(int64(pulse.PulseNumber)))
}

func pulseProcessingWatchdog(ctx context.Context, pulse insolar.Pulse, done chan struct{}) {
	logger := inslogger.FromContext(ctx)

	go func() {
		select {
		case <-time.After(time.Second * time.Duration(pulse.NextPulseNumber-pulse.PulseNumber)):
			logger.Errorf("Node stopped due to long pulse processing, pulse:%v", pulse.PulseNumber)
		case <-done:
		}
	}()
}
