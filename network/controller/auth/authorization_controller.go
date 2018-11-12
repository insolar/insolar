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
	"encoding/gob"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/nodenetwork"
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
	transport           network.InternalTransport
	keeper              network.NodeKeeper
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
	Address     string
	Version     string
}

// ResponseAuthorize
type ResponseAuthorize struct {
	ActiveNodes []core.Node
	Error       string
}

func init() {
	gob.Register(&RequestGetNonce{})
	gob.Register(&ResponseGetNonce{})
	gob.Register(&RequestAuthorize{})
	gob.Register(&ResponseAuthorize{})
}

func (ac *AuthorizationController) Authorize() error {
	hosts := ac.bootstrapController.GetBootstrapHosts()
	if len(hosts) == 0 {
		log.Info("Empty list of bootstrap hosts")
		return nil
	}

	ch := make(chan []core.Node, len(hosts))
	for _, h := range hosts {
		go func(ch chan<- []core.Node, h *host.Host) {
			activeNodes, err := ac.authorizeOnHost(h)
			if err != nil {
				log.Error(err)
				return
			}
			ch <- activeNodes
		}(ch, h)
	}

	activeLists := ac.collectActiveLists(ch, len(hosts))
	if len(activeLists) == 0 {
		return errors.New("Failed to authorize on any of discovery nodes")
	}
	activeNodes, success := MajorityRuleCheck(activeLists, ac.options.MajorityRule)
	if !success {
		return errors.New("Majority rule check failed")
	}
	ac.keeper.AddActiveNodes(activeNodes)
	return nil
}

func (ac *AuthorizationController) collectActiveLists(ch <-chan []core.Node, count int) [][]core.Node {
	receivedResults := make([][]core.Node, 0)
	for {
		select {
		case activeNodeList := <-ch:
			receivedResults = append(receivedResults, activeNodeList)
			if len(receivedResults) == count {
				return receivedResults
			}
		case <-time.After(ac.options.AuthorizeTimeout):
			log.Warnf("Authorize timeout, successful auths: %d/%d", len(receivedResults), count)
			return receivedResults
		}
	}
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
	request := ac.transport.NewRequestBuilder().Type(types.Authorize).Data(&RequestAuthorize{
		SignedNonce: signedNonce,
		NodeRoles:   []core.NodeRole{core.RoleUnknown},
		Address:     ac.transport.PublicAddress(),
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
	if ac.signer == nil {
		return ac.getNonceErrorResponse(request, "Signer is not initialized"), nil
	}
	nonce, err := ac.signer.AddPendingNode(request.GetSender())
	if err != nil {
		return ac.getNonceErrorResponse(request, err.Error()), nil
	}
	return ac.transport.BuildResponse(request, &ResponseGetNonce{Nonce: nonce}), nil
}

func (ac *AuthorizationController) getNonceErrorResponse(request network.Request, err string) network.Response {
	log.Warn("Error processing nonce request: " + err)
	return ac.transport.BuildResponse(request, &ResponseGetNonce{Error: err})
}

func (ac *AuthorizationController) processAuthorizeRequest(request network.Request) (network.Response, error) {
	if ac.signer == nil {
		return ac.getAuthErrorResponse(request, "Signer is not initialized"), nil
	}
	data := request.GetData().(*RequestAuthorize)
	err := ac.signer.AuthorizeNode(request.GetSender(), data.SignedNonce)
	if err != nil {
		return ac.getAuthErrorResponse(request, "Signed nonce check failed: "+err.Error()), nil
	}
	if ac.keeper == nil {
		return ac.getAuthErrorResponse(request, "NodeKeeper is not initialized"), nil
	}

	// TODO: fix this with new consensus
	// waitToken, err := ac.keeper.AddUnsync(request.GetSender(), data.NodeRoles, data.Address, data.Version)
	// if err != nil {
	// 	return ac.getAuthErrorResponse(request, "Error adding to unsync list: "+err.Error()), nil
	// }
	// select {
	// case d := <-waitToken:
	// 	if d == nil {
	// 		return ac.getAuthErrorResponse(request, "Error adding to unsync list: channel closed"), nil
	// 	}
	// case <-time.After(ac.options.AuthorizeTimeout):
	// 	return ac.getAuthErrorResponse(request, "Error adding to unsync list: timeout"), nil
	// }

	node := nodenetwork.NewNode(
		request.GetSender(),
		data.NodeRoles,
		nil,
		core.PulseNumber(0),
		data.Address,
		data.Version)
	if CheckShortIDCollision(ac.keeper, node.ShortID()) {
		CorrectShortIDCollision(ac.keeper, node)
	}
	ac.keeper.AddActiveNodes([]core.Node{node})

	return ac.transport.BuildResponse(request, &ResponseAuthorize{ActiveNodes: ac.keeper.GetActiveNodes()}), nil
}

func (ac *AuthorizationController) getAuthErrorResponse(request network.Request, err string) network.Response {
	log.Warn("Error processing auth request: " + err)
	return ac.transport.BuildResponse(request, &ResponseAuthorize{Error: err})
}

func (ac *AuthorizationController) Start(components core.Components) {
	ac.signer = NewSigner(components.Certificate, components.NetworkCoordinator)
	ac.keeper = components.NodeNetwork.(network.NodeKeeper)

	ac.transport.RegisterPacketHandler(types.GetNonce, ac.processNonceRequest)
	ac.transport.RegisterPacketHandler(types.Authorize, ac.processAuthorizeRequest)
}

func NewAuthorizationController(options *common.Options, bootstrapController common.BootstrapController,
	transport network.InternalTransport) *AuthorizationController {
	return &AuthorizationController{options: options, bootstrapController: bootstrapController, transport: transport}
}
