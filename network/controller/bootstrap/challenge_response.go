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

package bootstrap

import (
	"context"
	"encoding/gob"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

type ChallengeResponseController interface {
	component.Starter

	Execute(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) (*ChallengePayload, error)
}

type challengeResponseController struct {
	SessionManager SessionManager           `inject:""`
	Cryptography   core.CryptographyService `inject:""`
	NodeKeeper     network.NodeKeeper       `inject:""`

	options   *common.Options
	transport network.InternalTransport
}

type Nonce []byte
type SignedNonce []byte

// Node                           Discovery Node
//  1| ------ ChallengeRequest -----> |
//  2| <-- SignedChallengeResponse -- |
//  3| --- SignedChallengeRequest --> |
//  4| <----- ChallengeResponse ----- |
// ------------------------------------

type ChallengeResponseHeader struct {
	Success bool
	Error   string
}

type ChallengeRequest struct {
	SessionID SessionID

	Nonce Nonce
}

type SignedChallengeResponse struct {
	Header  ChallengeResponseHeader
	Payload *SignedChallengePayload
}

type SignedChallengePayload struct {
	SignedNonce       SignedNonce
	XorDiscoveryNonce Nonce
	DiscoveryNonce    Nonce
}

type SignedChallengeRequest struct {
	SessionID SessionID

	SignedDiscoveryNonce SignedNonce
	XorNonce             Nonce
}

type ChallengeResponse struct {
	Header  ChallengeResponseHeader
	Payload *ChallengePayload
}

type ChallengePayload struct {
	// CurrentPulse  core.Pulse
	// State         core.NetworkState
	AssignShortID core.ShortNodeID
}

func init() {
	gob.Register(&ChallengeRequest{})
	gob.Register(&SignedChallengeResponse{})
	gob.Register(&SignedChallengeRequest{})
	gob.Register(&ChallengeResponse{})
}

func (cr *challengeResponseController) processChallenge1(ctx context.Context, request network.Request) (network.Response, error) {
	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.processChallenge1")
	defer span.End()
	data := request.GetData().(*ChallengeRequest)
	// CheckSession is performed in SetDiscoveryNonce too, but we want to return early if the request is invalid
	err := cr.SessionManager.CheckSession(data.SessionID, Authorized)
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, err.Error()), nil
	}
	xorNonce, err := GenerateNonce()
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, "error generating discovery xor nonce: "+err.Error()), nil
	}
	sign, err := cr.Cryptography.Sign(Xor(data.Nonce, xorNonce))
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, "error signing nonce: "+err.Error()), nil
	}
	discoveryNonce, err := GenerateNonce()
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, "error generating discovery nonce: "+err.Error()), nil
	}
	err = cr.SessionManager.SetDiscoveryNonce(data.SessionID, discoveryNonce)
	if err != nil {
		return cr.buildChallenge1ErrorResponse(ctx, request, err.Error()), nil
	}
	response := cr.transport.BuildResponse(ctx, request, &SignedChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: true,
		},
		Payload: &SignedChallengePayload{
			SignedNonce:       sign.Bytes(),
			XorDiscoveryNonce: xorNonce,
			DiscoveryNonce:    discoveryNonce,
		},
	})
	return response, nil
}

func (cr *challengeResponseController) buildChallenge1ErrorResponse(ctx context.Context, request network.Request, err string) network.Response {
	log.Warn(err)
	return cr.transport.BuildResponse(ctx, request, &ChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: false,
			Error:   err,
		},
	})
}

func (cr *challengeResponseController) processChallenge2(ctx context.Context, request network.Request) (network.Response, error) {
	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.processChallenge2")
	defer span.End()
	data := request.GetData().(*SignedChallengeRequest)
	cert, discoveryNonce, err := cr.SessionManager.GetChallengeData(data.SessionID)
	if err != nil {
		return cr.buildChallenge2ErrorResponse(ctx, request, err.Error()), nil
	}
	sign := core.SignatureFromBytes(data.SignedDiscoveryNonce)
	success := cr.Cryptography.Verify(cert.GetPublicKey(), sign, Xor(data.XorNonce, discoveryNonce))
	if !success {
		return cr.buildChallenge2ErrorResponse(ctx, request, "node %s signature check failed"), nil
	}
	err = cr.SessionManager.ChallengePassed(data.SessionID)
	if err != nil {
		return cr.buildChallenge2ErrorResponse(ctx, request, err.Error()), nil
	}
	response := cr.transport.BuildResponse(ctx, request, &ChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: true,
		},
		Payload: &ChallengePayload{
			AssignShortID: GenerateShortID(cr.NodeKeeper, *cert.GetNodeRef()),
		},
	})
	return response, nil
}

func (cr *challengeResponseController) buildChallenge2ErrorResponse(ctx context.Context, request network.Request, err string) network.Response {
	log.Warn(err)
	return cr.transport.BuildResponse(ctx, request, &SignedChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: false,
			Error:   err,
		},
	})
}

func (cr *challengeResponseController) Start(ctx context.Context) error {
	cr.transport.RegisterPacketHandler(types.Challenge1, cr.processChallenge1)
	cr.transport.RegisterPacketHandler(types.Challenge2, cr.processChallenge2)
	return nil
}

func (cr *challengeResponseController) sendRequest1(ctx context.Context, discoveryHost *host.Host,
	sessionID SessionID, nonce Nonce) (*SignedChallengePayload, error) {

	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.sendRequest1")
	defer span.End()
	request := cr.transport.NewRequestBuilder().Type(types.Challenge1).Data(&ChallengeRequest{
		SessionID: sessionID, Nonce: nonce}).Build()
	future, err := cr.transport.SendRequestPacket(ctx, request, discoveryHost)
	if err != nil {
		return nil, errors.Wrap(err, "Error sending challenge request")
	}
	response, err := future.GetResponse(cr.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting response for challenge request")
	}
	data := response.GetData().(*SignedChallengeResponse)
	if !data.Header.Success {
		return nil, errors.Wrap(err, "Discovery node returned error for challenge request: "+data.Header.Error)
	}
	return data.Payload, nil
}

func (cr *challengeResponseController) sendRequest2(ctx context.Context, discoveryHost *host.Host,
	sessionID SessionID, signedDiscoveryNonce SignedNonce, xorNonce Nonce) (*ChallengePayload, error) {

	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.sendRequest2")
	defer span.End()
	request := cr.transport.NewRequestBuilder().Type(types.Challenge2).Data(&SignedChallengeRequest{
		SessionID: sessionID, XorNonce: xorNonce, SignedDiscoveryNonce: signedDiscoveryNonce}).Build()
	future, err := cr.transport.SendRequestPacket(ctx, request, discoveryHost)
	if err != nil {
		return nil, errors.Wrap(err, "Error sending challenge request")
	}
	response, err := future.GetResponse(cr.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting response for challenge request")
	}
	data := response.GetData().(*ChallengeResponse)
	if !data.Header.Success {
		return nil, errors.Wrap(err, "Discovery node returned error for challenge request: "+data.Header.Error)
	}
	return data.Payload, nil
}

// Execute double challenge response between the node and the discovery node (step 3 of the bootstrap process)
func (cr *challengeResponseController) Execute(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) (*ChallengePayload, error) {
	ctx, span := instracer.StartSpan(ctx, "ChallengeResponseController.Execute")
	defer span.End()
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, errors.Wrap(err, "error generating nonce")
	}
	inslogger.FromContext(ctx).Debugf("Generated nonce: %s", base58.Encode(nonce))

	data, err := cr.sendRequest1(ctx, discoveryNode.Host, sessionID, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "error executing challenge response (step 1)")
	}

	inslogger.FromContext(ctx).Debugf("Discovery SignedNonce: %s", base58.Encode(data.SignedNonce))
	inslogger.FromContext(ctx).Debugf("Discovery DiscoveryNonce: %s", base58.Encode(data.DiscoveryNonce))
	inslogger.FromContext(ctx).Debugf("Discovery XorDiscoveryNonce: %s", base58.Encode(data.XorDiscoveryNonce))

	sign := core.SignatureFromBytes(data.SignedNonce)
	success := cr.Cryptography.Verify(discoveryNode.Node.GetPublicKey(), sign, Xor(nonce, data.XorDiscoveryNonce))
	if !success {
		return nil, errors.New("Error checking signed nonce from discovery node")
	}

	xorNonce, err := GenerateNonce()
	if err != nil {
		return nil, errors.Wrap(err, "error generating xor nonce")
	}
	signedDiscoveryNonce, err := cr.Cryptography.Sign(Xor(xorNonce, data.DiscoveryNonce))
	if err != nil {
		return nil, errors.Wrap(err, "error signing discovery nonce")
	}
	payload, err := cr.sendRequest2(ctx, discoveryNode.Host, sessionID, signedDiscoveryNonce.Bytes(), xorNonce)
	if err != nil {
		return nil, errors.Wrap(err, "error executing challenge response (step 2)")
	}
	return payload, nil
}

func NewChallengeResponseController(options *common.Options, transport network.InternalTransport) ChallengeResponseController {
	return &challengeResponseController{
		options:   options,
		transport: transport,
	}
}
