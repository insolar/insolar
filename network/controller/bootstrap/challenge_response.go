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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

type ChallengeResponseController struct {
	options        *common.Options
	transport      network.InternalTransport
	cryptoSrv      core.CryptographyService
	sessionManager *SessionManager
	keeper         network.NodeKeeper
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

func (cr *ChallengeResponseController) processChallenge1(request network.Request) (network.Response, error) {
	data := request.GetData().(*ChallengeRequest)
	// CheckSession is performed in SetDiscoveryNonce too, but we want to return early if the request is invalid
	err := cr.sessionManager.CheckSession(data.SessionID, Authorized)
	if err != nil {
		return cr.buildChallenge1ErrorResponse(request, err.Error()), nil
	}
	xorNonce, err := GenerateNonce()
	if err != nil {
		return cr.buildChallenge1ErrorResponse(request, "error generating discovery xor nonce: "+err.Error()), nil
	}
	sign, err := cr.cryptoSrv.Sign(Xor(data.Nonce, xorNonce))
	if err != nil {
		return cr.buildChallenge1ErrorResponse(request, "error signing nonce: "+err.Error()), nil
	}
	discoveryNonce, err := GenerateNonce()
	if err != nil {
		return cr.buildChallenge1ErrorResponse(request, "error generating discovery nonce: "+err.Error()), nil
	}
	err = cr.sessionManager.SetDiscoveryNonce(data.SessionID, discoveryNonce)
	if err != nil {
		return cr.buildChallenge1ErrorResponse(request, err.Error()), nil
	}
	response := cr.transport.BuildResponse(request, &SignedChallengeResponse{
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

func (cr *ChallengeResponseController) buildChallenge1ErrorResponse(request network.Request, err string) network.Response {
	log.Warn(err)
	return cr.transport.BuildResponse(request, &ChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: false,
			Error:   err,
		},
	})
}

func (cr *ChallengeResponseController) processChallenge2(request network.Request) (network.Response, error) {
	data := request.GetData().(*SignedChallengeRequest)
	cert, discoveryNonce, err := cr.sessionManager.GetChallengeData(data.SessionID)
	if err != nil {
		return cr.buildChallenge2ErrorResponse(request, err.Error()), nil
	}
	sign := core.SignatureFromBytes(data.SignedDiscoveryNonce)
	success := cr.cryptoSrv.Verify(cert.GetPublicKey(), sign, Xor(data.XorNonce, discoveryNonce))
	if !success {
		return cr.buildChallenge2ErrorResponse(request, "node %s signature check failed"), nil
	}
	err = cr.sessionManager.ChallengePassed(data.SessionID)
	if err != nil {
		return cr.buildChallenge2ErrorResponse(request, err.Error()), nil
	}
	response := cr.transport.BuildResponse(request, &ChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: true,
		},
		Payload: &ChallengePayload{
			AssignShortID: GenerateShortID(cr.keeper, *cert.GetNodeRef()),
		},
	})
	return response, nil
}

func (cr *ChallengeResponseController) buildChallenge2ErrorResponse(request network.Request, err string) network.Response {
	log.Warn(err)
	return cr.transport.BuildResponse(request, &SignedChallengeResponse{
		Header: ChallengeResponseHeader{
			Success: false,
			Error:   err,
		},
	})
}

func (cr *ChallengeResponseController) Start(cryptoSrv core.CryptographyService, keeper network.NodeKeeper) {
	cr.keeper = keeper
	cr.cryptoSrv = cryptoSrv
	cr.transport.RegisterPacketHandler(types.Challenge1, cr.processChallenge1)
	cr.transport.RegisterPacketHandler(types.Challenge2, cr.processChallenge2)
}

func (cr *ChallengeResponseController) sendRequest1(discoveryHost *host.Host,
	sessionID SessionID, nonce Nonce) (*SignedChallengePayload, error) {

	request := cr.transport.NewRequestBuilder().Type(types.Challenge1).Data(&ChallengeRequest{
		SessionID: sessionID, Nonce: nonce}).Build()
	future, err := cr.transport.SendRequestPacket(request, discoveryHost)
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

func (cr *ChallengeResponseController) sendRequest2(discoveryHost *host.Host,
	sessionID SessionID, signedDiscoveryNonce SignedNonce, xorNonce Nonce) (*ChallengePayload, error) {

	request := cr.transport.NewRequestBuilder().Type(types.Challenge2).Data(&SignedChallengeRequest{
		SessionID: sessionID, XorNonce: xorNonce, SignedDiscoveryNonce: signedDiscoveryNonce}).Build()
	future, err := cr.transport.SendRequestPacket(request, discoveryHost)
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
func (cr *ChallengeResponseController) Execute(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) (*ChallengePayload, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, errors.Wrap(err, "error generating nonce")
	}
	inslogger.FromContext(ctx).Debugf("Generated nonce: %s", base58.Encode(nonce))

	data, err := cr.sendRequest1(discoveryNode.Host, sessionID, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "error executing challenge response (step 1)")
	}

	inslogger.FromContext(ctx).Debugf("Discovery SignedNonce: %s", base58.Encode(data.SignedNonce))
	inslogger.FromContext(ctx).Debugf("Discovery DiscoveryNonce: %s", base58.Encode(data.DiscoveryNonce))
	inslogger.FromContext(ctx).Debugf("Discovery XorDiscoveryNonce: %s", base58.Encode(data.XorDiscoveryNonce))

	sign := core.SignatureFromBytes(data.SignedNonce)
	success := cr.cryptoSrv.Verify(discoveryNode.Node.GetPublicKey(), sign, Xor(nonce, data.XorDiscoveryNonce))
	if !success {
		return nil, errors.New("Error checking signed nonce from discovery node")
	}

	xorNonce, err := GenerateNonce()
	if err != nil {
		return nil, errors.Wrap(err, "error generating xor nonce")
	}
	signedDiscoveryNonce, err := cr.cryptoSrv.Sign(Xor(xorNonce, data.DiscoveryNonce))
	if err != nil {
		return nil, errors.Wrap(err, "error signing discovery nonce")
	}
	payload, err := cr.sendRequest2(discoveryNode.Host, sessionID, signedDiscoveryNonce.Bytes(), xorNonce)
	if err != nil {
		return nil, errors.Wrap(err, "error executing challenge response (step 2)")
	}
	return payload, nil
}

func NewChallengeResponseController(options *common.Options, transport network.InternalTransport,
	sessionManager *SessionManager) *ChallengeResponseController {

	return &ChallengeResponseController{
		options:        options,
		transport:      transport,
		sessionManager: sessionManager,
	}
}
