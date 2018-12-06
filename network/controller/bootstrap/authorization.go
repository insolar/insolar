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

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

// AuthorizationController is intended
type AuthorizationController struct {
	options        *common.Options
	transport      network.InternalTransport
	keeper         network.NodeKeeper
	coordinator    core.NetworkCoordinator
	sessionManager *SessionManager
}

type OperationCode uint8

const (
	OpConfirmed OperationCode = iota + 1
	OpRejected
)

// AuthorizationRequest
type AuthorizationRequest struct {
	Certificate []byte
}

// AuthorizationResponse
type AuthorizationResponse struct {
	Code      OperationCode
	Error     string
	SessionID SessionID
}

// RegistrationRequest
type RegistrationRequest struct {
	SessionID SessionID
	JoinClaim *packets.NodeJoinClaim
}

// RegistrationResponse
type RegistrationResponse struct {
	Code  OperationCode
	Error string
}

func init() {
	gob.Register(&AuthorizationRequest{})
	gob.Register(&AuthorizationResponse{})
	gob.Register(&RegistrationRequest{})
	gob.Register(&RegistrationResponse{})
}

// Authorize node on the discovery node (step 2 of the bootstrap process)
func (ac *AuthorizationController) Authorize(ctx context.Context, discoveryNode *DiscoveryNode, cert core.AuthorizationCertificate) (SessionID, error) {
	inslogger.FromContext(ctx).Infof("Authorizing on host: %s", discoveryNode)

	serializedCert, err := certificate.Serialize(cert)
	if err != nil {
		return 0, errors.Wrap(err, "Error serializing certificate")
	}

	request := ac.transport.NewRequestBuilder().Type(types.Authorize).Data(&AuthorizationRequest{
		Certificate: serializedCert,
	}).Build()
	future, err := ac.transport.SendRequestPacket(request, discoveryNode.Host)
	if err != nil {
		return 0, errors.Wrapf(err, "Error sending authorize request")
	}
	response, err := future.GetResponse(ac.options.PacketTimeout)
	if err != nil {
		return 0, errors.Wrapf(err, "Error getting response for authorize request")
	}
	data := response.GetData().(*AuthorizationResponse)
	if data.Code == OpRejected {
		return 0, errors.New("Authorize rejected: " + data.Error)
	}
	return data.SessionID, nil
}

// Register node on the discovery node (step 4 of the bootstrap process)
func (ac *AuthorizationController) Register(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) error {
	inslogger.FromContext(ctx).Infof("Registering on host: %s", discoveryNode)

	originClaim, err := ac.keeper.GetOriginClaim()
	if err != nil {
		return errors.Wrap(err, "[ Register ] failed to get origin claim")
	}
	request := ac.transport.NewRequestBuilder().Type(types.Register).Data(&RegistrationRequest{
		SessionID: sessionID,
		JoinClaim: originClaim,
	}).Build()
	future, err := ac.transport.SendRequestPacket(request, discoveryNode.Host)
	if err != nil {
		return errors.Wrapf(err, "Error sending register request")
	}
	response, err := future.GetResponse(ac.options.PacketTimeout)
	if err != nil {
		return errors.Wrapf(err, "Error getting response for register request")
	}
	data := response.GetData().(*RegistrationResponse)
	if data.Code == OpRejected {
		return errors.New("Register rejected: " + data.Error)
	}
	return nil
}

func (ac *AuthorizationController) checkClaim(sessionID SessionID, claim *packets.NodeJoinClaim) error {
	session, err := ac.sessionManager.ReleaseSession(sessionID)
	if err != nil {
		return errors.Wrapf(err, "Error getting section %d for authorization", sessionID)
	}
	if !claim.NodeRef.Equal(session.NodeID) {
		return errors.New("Claim node ID is not equal to session node ID")
	}
	// TODO: check claim signature
	return nil
}

func (ac *AuthorizationController) processRegisterRequest(request network.Request) (network.Response, error) {
	data := request.GetData().(*RegistrationRequest)
	err := ac.checkClaim(data.SessionID, data.JoinClaim)
	if err != nil {
		responseAuthorize := &RegistrationResponse{Code: OpRejected, Error: err.Error()}
		return ac.transport.BuildResponse(request, responseAuthorize), nil
	}
	ac.keeper.AddPendingClaim(data.JoinClaim)
	return ac.transport.BuildResponse(request, &RegistrationResponse{Code: OpConfirmed}), nil
}

func (ac *AuthorizationController) processAuthorizeRequest(request network.Request) (network.Response, error) {
	data := request.GetData().(*AuthorizationRequest)
	cert, err := certificate.Deserialize(data.Certificate)
	if err != nil {
		return ac.transport.BuildResponse(request, &AuthorizationResponse{Code: OpRejected, Error: err.Error()}), nil
	}
	valid, err := ac.coordinator.ValidateCert(context.Background(), cert)
	if !valid {
		if err == nil {
			err = errors.New("Certificate validation failed")
		}
		return ac.transport.BuildResponse(request, &AuthorizationResponse{Code: OpRejected, Error: err.Error()}), nil
	}
	session := ac.sessionManager.NewSession(request.GetSender(), cert)
	return ac.transport.BuildResponse(request, &AuthorizationResponse{Code: OpConfirmed, SessionID: session}), nil
}

func (ac *AuthorizationController) Start(networkCoordinator core.NetworkCoordinator, nodeKeeper network.NodeKeeper) {
	ac.keeper = nodeKeeper
	ac.coordinator = networkCoordinator
	ac.transport.RegisterPacketHandler(types.Register, ac.processRegisterRequest)
	ac.transport.RegisterPacketHandler(types.Authorize, ac.processAuthorizeRequest)
}

func NewAuthorizationController(options *common.Options, transport network.InternalTransport, sessionManager *SessionManager) *AuthorizationController {
	return &AuthorizationController{
		options:        options,
		transport:      transport,
		sessionManager: sessionManager,
	}
}
