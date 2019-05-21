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

package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/platformpolicy"
)

const (
	registrationRetries = 20
)

type AuthorizationController interface {
	component.Initer

	Authorize(ctx context.Context, discoveryNode *DiscoveryNode, cert insolar.AuthorizationCertificate) (*packet.AuthorizationData, error)
	Register(ctx context.Context, discoveryNode *DiscoveryNode, sessionID SessionID) error
}

type authorizationController struct {
	NodeKeeper     network.NodeKeeper  `inject:""`
	Gatewayer      network.Gatewayer   `inject:""`
	SessionManager SessionManager      `inject:""`
	Network        network.HostNetwork `inject:""`

	options *common.Options
}

// Authorize node on the discovery node (step 2 of the bootstrap process)
func (ac *authorizationController) Authorize(ctx context.Context, discoveryNode *DiscoveryNode, cert insolar.AuthorizationCertificate) (*packet.AuthorizationData, error) {
	inslogger.FromContext(ctx).Infof("Authorizing on host: %s", discoveryNode.Host)
	inslogger.FromContext(ctx).Infof("cert: %s", cert)

	ctx, span := instracer.StartSpan(ctx, "AuthorizationController.Authorize")
	span.AddAttributes(
		trace.StringAttribute("node", discoveryNode.Node.GetNodeRef().String()),
	)
	defer span.End()
	serializedCert, err := certificate.Serialize(cert)
	if err != nil {
		return nil, errors.Wrap(err, "Error serializing certificate")
	}

	auth := &packet.AuthorizeRequest{Certificate: serializedCert}
	future, err := ac.Network.SendRequestToHost(ctx, types.Authorize, auth, discoveryNode.Host)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending authorize request")
	}
	response, err := future.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for authorize request")
	}
	if response.GetResponse() == nil || response.GetResponse().GetAuthorize() == nil {
		return nil, errors.Errorf("Authorize failed: got incorrect response: %s", response)
	}
	data := response.GetResponse().GetAuthorize()
	if data.Code == packet.Denied {
		return nil, errors.New("Authorize rejected: " + data.Error)
	}
	return data.Data, nil
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
	request := &packet.RegisterRequest{
		Version:   ac.NodeKeeper.GetOrigin().Version(),
		SessionID: uint64(sessionID),
		JoinClaim: originClaim,
	}
	future, err := ac.Network.SendRequestToHost(ctx, types.Register, request, discoveryNode.Host)
	if err != nil {
		return errors.Wrapf(err, "Error sending register request")
	}
	response, err := future.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return errors.Wrapf(err, "Error getting response for register request")
	}
	if response.GetResponse() == nil || response.GetResponse().GetRegister() == nil {
		return errors.Errorf("Register failed: got incorrect response: %s", response)
	}
	data := response.GetResponse().GetRegister()
	if data.Code == packet.Denied {
		return errors.New("Register rejected: " + data.Error)
	}
	if data.Code == packet.Retry {
		if attempt >= registrationRetries {
			return errors.Errorf("Exceeded maximum number of registration retries (%d)", registrationRetries)
		}
		log.Warnf("Failed to register on discovery node %s. Reason: node %s is already in network active list. "+
			"Retrying registration in %v", discoveryNode.Host, ac.NodeKeeper.GetOrigin().ID(), data.RetryIn)
		time.Sleep(time.Duration(data.RetryIn))
		return ac.register(ctx, discoveryNode, sessionID, attempt+1)
	}
	return nil
}

func (ac *authorizationController) buildRegistrationResponse(sessionID SessionID, claim *packets.NodeJoinClaim) *packet.RegisterResponse {
	session, err := ac.getSession(sessionID, claim)
	if err != nil {
		return &packet.RegisterResponse{Code: packet.Denied, Error: err.Error()}
	}
	if node := ac.NodeKeeper.GetAccessor().GetActiveNode(claim.NodeRef); node != nil {
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
		return &packet.RegisterResponse{Code: packet.Retry, RetryIn: int64(retryIn)}
	}
	return &packet.RegisterResponse{Code: packet.Confirmed}
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

func (ac *authorizationController) processRegisterRequest(ctx context.Context, request network.Packet) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetRegister() == nil {
		return nil, errors.Errorf("process register: got invalid protobuf request message: %s", request)
	}
	data := request.GetRequest().GetRegister()
	if data.Version != ac.NodeKeeper.GetOrigin().Version() {
		response := &packet.RegisterResponse{Code: packet.Denied,
			Error: fmt.Sprintf("Joiner version %s does not match discovery version %s",
				data.Version, ac.NodeKeeper.GetOrigin().Version())}
		return ac.Network.BuildResponse(ctx, request, response), nil
	}
	// TODO: remove SessionID convertions
	response := ac.buildRegistrationResponse(SessionID(data.SessionID), data.JoinClaim)
	if response.Code != packet.Confirmed {
		return ac.Network.BuildResponse(ctx, request, response), nil
	}

	// TODO: fix Short ID assignment logic
	if CheckShortIDCollision(ac.NodeKeeper, data.JoinClaim.ShortNodeID) {
		response = &packet.RegisterResponse{Code: packet.Denied,
			Error: "Short ID of the joiner node conflicts with active node short ID"}
		return ac.Network.BuildResponse(ctx, request, response), nil
	}

	inslogger.FromContext(ctx).Infof("Added join claim from node %s", request.GetSender())
	ac.NodeKeeper.GetClaimQueue().Push(data.JoinClaim)
	return ac.Network.BuildResponse(ctx, request, response), nil
}

func (ac *authorizationController) processAuthorizeRequest(ctx context.Context, request network.Packet) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetAuthorize() == nil {
		return nil, errors.Errorf("process authorize: got invalid protobuf request message: %s", request)
	}
	data := request.GetRequest().GetAuthorize()
	cert, err := certificate.Deserialize(data.Certificate, platformpolicy.NewKeyProcessor())
	if err != nil {
		return ac.Network.BuildResponse(ctx, request, &packet.AuthorizeResponse{Code: packet.Denied, Error: err.Error()}), nil
	}
	valid, err := ac.Gatewayer.Gateway().Auther().ValidateCert(ctx, cert)
	if !valid {
		if err == nil {
			err = errors.New("Certificate validation failed")
		}
		return ac.Network.BuildResponse(ctx, request, &packet.AuthorizeResponse{Code: packet.Denied, Error: err.Error()}), nil
	}
	session := ac.SessionManager.NewSession(request.GetSender(), cert, ac.options.HandshakeSessionTTL)
	return ac.Network.BuildResponse(ctx, request, &packet.AuthorizeResponse{
		Code: packet.Confirmed,
		Data: &packet.AuthorizationData{
			SessionID:     uint64(session),
			AssignShortID: uint32(GenerateShortID(ac.NodeKeeper, *cert.GetNodeRef())),
		},
	}), nil
}

func (ac *authorizationController) Init(ctx context.Context) error {
	ac.Network.RegisterRequestHandler(types.Register, ac.processRegisterRequest)
	ac.Network.RegisterRequestHandler(types.Authorize, ac.processAuthorizeRequest)
	return nil
}

func NewAuthorizationController(options *common.Options) AuthorizationController {
	return &authorizationController{options: options}
}
