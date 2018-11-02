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

package auth

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/version"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

// AuthorizationController is intended
type AuthorizationController struct {
	options             *common.Options
	bootstrapController common.BootstrapController
	signer              *Signer
	transport           hostnetwork.InternalTransport
}

func (ac *AuthorizationController) Authorize() error {
	hosts := ac.bootstrapController.GetBootstrapHosts()
	if len(hosts) == 0 {
		return errors.New("Empty list of bootstrap hosts")
	}

	return nil
}

// RequestGetNonce
type RequestGetNonce struct{}

// ResponseGetNonce
type ResponseGetNonce struct {
	Nonce Nonce
	Error string
}

// RequestAuthorize
type RequestAuthorize struct {
	SignedNonce []byte
	NodeRoles   []core.NodeRole
	Version     string
}

// ResponseAuthorize
type ResponseAuthorize struct {
	ActiveNodes []core.Node
	Error       string
}

// authorizeOnHost send all authorize requests to host and get list of active nodes
func (ac *AuthorizationController) authorizeOnHost(h *host.Host) ([]core.Node, error) {
	log.Infof("Authorizing on host: %s", h)

	nonce, err := ac.sendNonceRequest(h)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting nonce from discovery node %s", h)
	}
	log.Debugf("Got nonce from discovery node: %s", base58.Encode(nonce))
	signedNonce, err := ac.signer.SignNonce(nonce)
	if err != nil {
		return nil, errors.Wrapf(err, "Error signing received nonce from node %s", h)
	}
	nodes, err := ac.sendAuthorizeRequest(signedNonce, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Error authorizing on discovery node %s", h)
	}
	return nodes, nil
}

func (ac *AuthorizationController) sendNonceRequest(h *host.Host) (Nonce, error) {
	log.Debugf("Sending nonce request to host: %s", h)

	request := ac.transport.NewRequestBuilder().Type(types.GetNonce).Data(&RequestGetNonce{}).Build()
	future, err := ac.transport.SendRequestPacket(request, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending nonce request")
	}
	response, err := future.GetResponse(ac.options.AuthorizeTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for nonce request")
	}
	data := response.GetData().(*ResponseGetNonce)
	if data.Nonce == nil {
		return nil, errors.New("Discovery node returned error for nonce request: " + data.Error)
	}
	return data.Nonce, nil
}

func (ac *AuthorizationController) sendAuthorizeRequest(signedNonce []byte, h *host.Host) ([]core.Node, error) {
	log.Debugf("Sending authorize request to host: %s", h)

	request := ac.transport.NewRequestBuilder().Type(types.Authorize).Data(&RequestAuthorize{
		SignedNonce: signedNonce,
		NodeRoles:   []core.NodeRole{core.RoleUnknown},
		Version:     version.Version,
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
	if data.ActiveNodes == nil {
		return nil, errors.New("Discovery node returned error for authorize request: " + data.Error)
	}
	return data.ActiveNodes, nil
}

func (ac *AuthorizationController) processNonceRequest(request network.Request) (network.Response, error) {
	return nil, errors.New("not implemented")
}

func (ac *AuthorizationController) processAuthorizeRequest(request network.Request) (network.Response, error) {
	return nil, errors.New("not implemented")
}

func (ac *AuthorizationController) Start(components core.Components) {
	ac.signer = NewSigner(components.Certificate)

	ac.transport.RegisterPacketHandler(types.GetNonce, ac.processNonceRequest)
	ac.transport.RegisterPacketHandler(types.Authorize, ac.processAuthorizeRequest)
}

func NewAuthorizationController(options *common.Options, bootstrapController common.BootstrapController,
	transport hostnetwork.InternalTransport) *AuthorizationController {
	return &AuthorizationController{options: options, bootstrapController: bootstrapController, transport: transport}
}
