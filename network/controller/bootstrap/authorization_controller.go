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

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

// AuthorizationController is intended
type AuthorizationController struct {
	options             *common.Options
	bootstrapController common.BootstrapController
	transport           network.InternalTransport
	keeper              network.NodeKeeper
}

// RequestAuthorize
type RequestAuthorize struct {
	SessionID SessionID
	JoinClaim *packets.NodeJoinClaim
}

type AuthorizeCode uint8

const (
	AuthConfirmed AuthorizeCode = iota + 1
	AuthRejected
)

// ResponseAuthorize
type ResponseAuthorize struct {
	AuthorizeCode AuthorizeCode
	Error         string
}

func init() {
	gob.Register(&RequestAuthorize{})
	gob.Register(&ResponseAuthorize{})
}

// authorizeOnHost send all authorize requests to host and get list of active nodes
func (ac *AuthorizationController) AuthorizeOnHost(ctx context.Context, sessionID SessionID, h *host.Host) (*ResponseAuthorize, error) {
	inslogger.FromContext(ctx).Infof("Authorizing on host: %s", h)

	request := ac.transport.NewRequestBuilder().Type(types.Authorize).Data(&RequestAuthorize{
		SessionID: sessionID,
	}).Build()
	future, err := ac.transport.SendRequestPacket(request, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending authorize request")
	}
	response, err := future.GetResponse(ac.options.AuthorizeTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for authorize request")
	}
	data := response.GetData().(*ResponseAuthorize)
	return data, nil
}

func (ac *AuthorizationController) checkClaim(sessionID SessionID, claim *packets.NodeJoinClaim) error {
	// TODO: check ID, signature and sessionID
	return nil
}

func (ac *AuthorizationController) processAuthorizeRequest(request network.Request) (network.Response, error) {
	data := request.GetData().(*RequestAuthorize)
	err := ac.checkClaim(data.SessionID, data.JoinClaim)
	if err != nil {
		responseAuthorize := &ResponseAuthorize{AuthorizeCode: AuthRejected, Error: err.Error()}
		return ac.transport.BuildResponse(request, responseAuthorize), nil
	}
	ac.keeper.AddPendingClaim(data.JoinClaim)
	return ac.transport.BuildResponse(request, &ResponseAuthorize{AuthorizeCode: AuthConfirmed}), nil
}

func (ac *AuthorizationController) Start(cryptographyService core.CryptographyService,
	networkCoordinator core.NetworkCoordinator, nodeKeeper network.NodeKeeper) {

	ac.keeper = nodeKeeper
	ac.transport.RegisterPacketHandler(types.Authorize, ac.processAuthorizeRequest)
}

func NewAuthorizationController(options *common.Options, bootstrapController common.BootstrapController,
	transport network.InternalTransport) *AuthorizationController {
	return &AuthorizationController{options: options, bootstrapController: bootstrapController, transport: transport}
}
