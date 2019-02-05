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

package bootstrap

import (
	"context"
	"encoding/gob"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type AuthorizationController interface {
	component.Starter

	Authorize(ctx context.Context, discoveryNode *DiscoveryNode, cert core.AuthorizationCertificate) (SessionID, error)
	Register(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) error
}

type authorizationController struct {
	NodeKeeper         network.NodeKeeper      `inject:""`
	NetworkCoordinator core.NetworkCoordinator `inject:""`
	SessionManager     SessionManager          `inject:""`

	options   *common.Options
	transport network.InternalTransport
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
func (ac *authorizationController) Authorize(ctx context.Context, discoveryNode *DiscoveryNode, cert core.AuthorizationCertificate) (SessionID, error) {
	inslogger.FromContext(ctx).Infof("Authorizing on host: %s", discoveryNode)

	ctx, span := instracer.StartSpan(ctx, "AuthorizationController.Authorize")
	span.AddAttributes(
		trace.StringAttribute("node", discoveryNode.Node.GetNodeRef().String()),
	)
	defer span.End()
	serializedCert, err := certificate.Serialize(cert)
	if err != nil {
		return 0, errors.Wrap(err, "Error serializing certificate")
	}

	request := ac.transport.NewRequestBuilder().Type(types.Authorize).Data(&AuthorizationRequest{
		Certificate: serializedCert,
	}).Build()
	future, err := ac.transport.SendRequestPacket(ctx, request, discoveryNode.Host)
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
func (ac *authorizationController) Register(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) error {
	inslogger.FromContext(ctx).Infof("Registering on host: %s", discoveryNode)

	ctx, span := instracer.StartSpan(ctx, "AuthorizationController.Register")
	span.AddAttributes(
		trace.StringAttribute("node", discoveryNode.Node.GetNodeRef().String()),
	)
	defer span.End()
	originClaim, err := ac.NodeKeeper.GetOriginClaim()
	if err != nil {
		return errors.Wrap(err, "[ Register ] failed to get origin claim")
	}
	request := ac.transport.NewRequestBuilder().Type(types.Register).Data(&RegistrationRequest{
		SessionID: sessionID,
		JoinClaim: originClaim,
	}).Build()
	future, err := ac.transport.SendRequestPacket(ctx, request, discoveryNode.Host)
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

func (ac *authorizationController) checkClaim(sessionID SessionID, claim *packets.NodeJoinClaim) error {
	session, err := ac.SessionManager.ReleaseSession(sessionID)
	if err != nil {
		return errors.Wrapf(err, "Error getting session %d for authorization", sessionID)
	}
	if !claim.NodeRef.Equal(session.NodeID) {
		return errors.New("Claim node ID is not equal to session node ID")
	}
	// TODO: check claim signature
	return nil
}

func (ac *authorizationController) processRegisterRequest(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*RegistrationRequest)
	err := ac.checkClaim(data.SessionID, data.JoinClaim)
	if err != nil {
		responseAuthorize := &RegistrationResponse{Code: OpRejected, Error: err.Error()}
		return ac.transport.BuildResponse(ctx, request, responseAuthorize), nil
	}
	ac.NodeKeeper.AddPendingClaim(data.JoinClaim)
	return ac.transport.BuildResponse(ctx, request, &RegistrationResponse{Code: OpConfirmed}), nil
}

func (ac *authorizationController) processAuthorizeRequest(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*AuthorizationRequest)
	cert, err := certificate.Deserialize(data.Certificate, platformpolicy.NewKeyProcessor())
	if err != nil {
		return ac.transport.BuildResponse(ctx, request, &AuthorizationResponse{Code: OpRejected, Error: err.Error()}), nil
	}
	valid, err := ac.NetworkCoordinator.ValidateCert(context.Background(), cert)
	if !valid {
		if err == nil {
			err = errors.New("Certificate validation failed")
		}
		return ac.transport.BuildResponse(ctx, request, &AuthorizationResponse{Code: OpRejected, Error: err.Error()}), nil
	}
	session := ac.SessionManager.NewSession(request.GetSender(), cert, ac.options.HandshakeSessionTTL)
	return ac.transport.BuildResponse(ctx, request, &AuthorizationResponse{Code: OpConfirmed, SessionID: session}), nil
}

func (ac *authorizationController) Start(ctx context.Context) error {
	ac.transport.RegisterPacketHandler(types.Register, ac.processRegisterRequest)
	ac.transport.RegisterPacketHandler(types.Authorize, ac.processAuthorizeRequest)
	return nil
}

func NewAuthorizationController(options *common.Options, transport network.InternalTransport) AuthorizationController {
	return &authorizationController{
		options:   options,
		transport: transport,
	}
}
