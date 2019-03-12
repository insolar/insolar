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
	"fmt"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

const (
	registrationRetries = 10
)

type AuthorizationController interface {
	component.Initer

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
	OpRetry
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
	Version   string
	JoinClaim *packets.NodeJoinClaim
}

// RegistrationResponse
type RegistrationResponse struct {
	Code    OperationCode
	RetryIn time.Duration
	Error   string
}

func init() {
	gob.Register(&AuthorizationRequest{})
	gob.Register(&AuthorizationResponse{})
	gob.Register(&RegistrationRequest{})
	gob.Register(&RegistrationResponse{})
}

// Authorize node on the discovery node (step 2 of the bootstrap process)
func (ac *authorizationController) Authorize(ctx context.Context, discoveryNode *DiscoveryNode, cert core.AuthorizationCertificate) (SessionID, error) {
	inslogger.FromContext(ctx).Infof("Authorizing on host: %s", discoveryNode.Host)
	inslogger.FromContext(ctx).Infof("cert: %s", cert)

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
	return ac.register(ctx, discoveryNode, sessionID, 0)
}

func (ac *authorizationController) register(ctx context.Context, discoveryNode *DiscoveryNode,
	sessionID SessionID, attempt int) error {

	if attempt == 0 {
		inslogger.FromContext(ctx).Infof("Registering on host: %s", discoveryNode.Host)
	} else {
		inslogger.FromContext(ctx).Infof("Registering on host: %s; attempt: %d", discoveryNode.Host, attempt+1)
	}

	ctx, span := instracer.StartSpan(ctx, "AuthorizationController.Register")
	span.AddAttributes(
		trace.StringAttribute("node", discoveryNode.Node.GetNodeRef().String()),
	)
	defer span.End()
	originClaim, err := ac.NodeKeeper.GetOriginJoinClaim()
	if err != nil {
		return errors.Wrap(err, "Failed to get origin claim")
	}
	request := ac.transport.NewRequestBuilder().Type(types.Register).Data(&RegistrationRequest{
		Version:   ac.NodeKeeper.GetOrigin().Version(),
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
	if data.Code == OpRetry {
		if attempt >= registrationRetries {
			return errors.Errorf("Exceeded maximum number of registration retries (%d)", registrationRetries)
		}
		log.Warnf("Failed to register on discovery node %s. Reason: node %s is already in network active list. "+
			"Retrying registration in %v", discoveryNode.Host, ac.NodeKeeper.GetOrigin().ID(), data.RetryIn)
		time.Sleep(data.RetryIn)
		return ac.register(ctx, discoveryNode, sessionID, attempt+1)
	}
	return nil
}

func (ac *authorizationController) buildRegistrationResponse(sessionID SessionID, claim *packets.NodeJoinClaim) *RegistrationResponse {
	session, err := ac.getSession(sessionID, claim)
	if err != nil {
		return &RegistrationResponse{Code: OpRejected, Error: err.Error()}
	}
	if node := ac.NodeKeeper.GetActiveNode(claim.NodeRef); node != nil {
		retryIn := session.TTL / 2

		keyProc := platformpolicy.NewKeyProcessor()
		// little hack: ignoring error, because it never fails in current implementation
		nodeKey, _ := keyProc.ExportPublicKeyBinary(node.PublicKey())

		log.Warnf("Joiner node (ID: %s, PK: %s) conflicts with node (ID: %s, PK: %s) in active list, sending request to reconnect in %v",
			claim.NodeRef, base58.Encode(claim.NodePK[:]), node.ID(), base58.Encode(nodeKey), retryIn)

		statsErr := stats.RecordWithTags(context.Background(), []tag.Mutator{
			tag.Upsert(tagNodeRef, claim.NodeRef.String()),
		}, statBootstrapReconnectRequired.M(1))
		if statsErr != nil {
			log.Warn("Failed to record reconnection retries metric: " + statsErr.Error())
		}

		ac.SessionManager.ProlongateSession(sessionID, session)
		return &RegistrationResponse{Code: OpRetry, RetryIn: retryIn}
	}
	return &RegistrationResponse{Code: OpConfirmed}
}

func (ac *authorizationController) getSession(sessionID SessionID, claim *packets.NodeJoinClaim) (*Session, error) {
	session, err := ac.SessionManager.ReleaseSession(sessionID)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting session %d for authorization", sessionID)
	}
	if !claim.NodeRef.Equal(session.NodeID) {
		return nil, errors.New("Claim node ID is not equal to session node ID")
	}
	// TODO: check claim signature
	return session, nil
}

func (ac *authorizationController) processRegisterRequest(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*RegistrationRequest)
	if data.Version != ac.NodeKeeper.GetOrigin().Version() {
		response := &RegistrationResponse{Code: OpRejected,
			Error: fmt.Sprintf("Joiner version %s does not match discovery version %s",
				data.Version, ac.NodeKeeper.GetOrigin().Version())}
		return ac.transport.BuildResponse(ctx, request, response), nil
	}
	response := ac.buildRegistrationResponse(data.SessionID, data.JoinClaim)
	if response.Code != OpConfirmed {
		return ac.transport.BuildResponse(ctx, request, response), nil
	}

	// TODO: fix Short ID assignment logic
	if CheckShortIDCollision(ac.NodeKeeper, data.JoinClaim.ShortNodeID) {
		response = &RegistrationResponse{Code: OpRejected,
			Error: "Short ID of the joiner node conflicts with active node short ID"}
		return ac.transport.BuildResponse(ctx, request, response), nil
	}

	inslogger.FromContext(ctx).Infof("Added join claim from node %s", request.GetSender())
	ac.NodeKeeper.AddPendingClaim(data.JoinClaim)
	return ac.transport.BuildResponse(ctx, request, response), nil
}

func (ac *authorizationController) processAuthorizeRequest(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*AuthorizationRequest)
	cert, err := certificate.Deserialize(data.Certificate, platformpolicy.NewKeyProcessor())
	if err != nil {
		return ac.transport.BuildResponse(ctx, request, &AuthorizationResponse{Code: OpRejected, Error: err.Error()}), nil
	}
	valid, err := ac.NetworkCoordinator.ValidateCert(ctx, cert)
	if !valid {
		if err == nil {
			err = errors.New("Certificate validation failed")
		}
		return ac.transport.BuildResponse(ctx, request, &AuthorizationResponse{Code: OpRejected, Error: err.Error()}), nil
	}
	session := ac.SessionManager.NewSession(request.GetSender(), cert, ac.options.HandshakeSessionTTL)
	return ac.transport.BuildResponse(ctx, request, &AuthorizationResponse{Code: OpConfirmed, SessionID: session}), nil
}

func (ac *authorizationController) Init(ctx context.Context) error {
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
