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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/controller/pinger"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type Bootstrapper struct {
	options   *common.Options
	transport network.InternalTransport
	pinger    *pinger.Pinger
	cert      core.Certificate
	keeper    network.NodeKeeper

	chosenDiscoveryNode *host.Host
}

type NodeBootstrapRequest struct{}

type NodeBootstrapResponse struct {
	Code         Code
	RedirectHost string
	RejectReason string
}

type GenesisRequest struct {
	Certificate core.Certificate
}

type GenesisResponse struct {
	Discovery core.Node
	Error     string
}

type StartSessionRequest struct{}

type StartSessionResponse struct {
	SessionID SessionID
}

type Code uint8

const (
	Accepted = Code(iota + 1)
	Rejected
	Redirected
)

func init() {
	gob.Register(&NodeBootstrapRequest{})
	gob.Register(&NodeBootstrapResponse{})
	gob.Register(&StartSessionRequest{})
	gob.Register(&StartSessionResponse{})
	gob.Register(&GenesisRequest{})
	gob.Register(&GenesisResponse{})
}

func (bc *Bootstrapper) GetChosenDiscoveryNode() *host.Host {
	return bc.chosenDiscoveryNode
}

func (bc *Bootstrapper) Bootstrap(ctx context.Context) error {
	ch := bc.getBootstrapHostsChannel(ctx, 1)
	host := bc.waitResultFromChannel(ctx, ch)
	if host == nil {
		return errors.New("Failed to bootstrap to any of discovery nodes")
	}
	bc.chosenDiscoveryNode = host
	return nil
}

func (bc *Bootstrapper) StartSession(ctx context.Context) (SessionID, error) {
	request := bc.transport.NewRequestBuilder().Type(types.StartSession).Data(nil).Build()
	future, err := bc.transport.SendRequestPacket(request, bc.chosenDiscoveryNode)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to send StartSession request to discovery node %s", bc.chosenDiscoveryNode)
	}
	response, err := future.GetResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to get StartSession response from discovery node %s", bc.chosenDiscoveryNode)
	}
	data := response.GetData().(*StartSessionResponse)
	return data.SessionID, nil
}

func (bc *Bootstrapper) checkActiveNode(node core.Node) error {
	// TODO: implement check ID, ShortID collision (maybe some other checks)
	return nil
}

func (bc *Bootstrapper) BootstrapDiscovery(ctx context.Context) error {
	discoveryNodes := bc.cert.GetBootstrapNodes()
	var err error
	discoveryNodes, err = RemoveOrigin(discoveryNodes, *bc.cert.GetRef())
	if err != nil {
		return errors.Wrapf(err, "Discovery bootstrap failed")
	}
	discoveryCount := len(discoveryNodes)
	if discoveryCount == 0 {
		return nil
	}

	var hosts []*host.Host
	for {
		ch := bc.getBootstrapHostsChannel(ctx, discoveryCount)
		hosts = bc.waitResultsFromChannel(ctx, ch, discoveryCount)
		if len(hosts) == discoveryCount {
			// we connected to all discovery nodes
			break
		}
	}
	activeNodes := make([]core.Node, 0)
	for _, h := range hosts {
		activeNode, err := bc.sendGenesisRequest(ctx, h)
		if err != nil {
			return errors.Wrapf(err, "Discovery bootstrap to host %s failed", h)
		}
		err = bc.checkActiveNode(activeNode)
		if err != nil {
			return errors.Wrapf(err, "Discovery check of host %s failed", h)
		}
		activeNodes = append(activeNodes, activeNode)
	}
	bc.keeper.AddActiveNodes(activeNodes)
	return nil
}

func (bc *Bootstrapper) sendGenesisRequest(ctx context.Context, h *host.Host) (core.Node, error) {
	request := bc.transport.NewRequestBuilder().Type(types.Genesis).Data(&GenesisRequest{
		Certificate: bc.cert,
	}).Build()
	future, err := bc.transport.SendRequestPacket(request, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send genesis request to address %s", h)
	}
	response, err := future.GetResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to genesis request from address %s", h)
	}
	data := response.GetData().(*GenesisResponse)
	if data.Discovery == nil {
		return nil, errors.New("Error genesis response from discovery node: " + data.Error)
	}
	return data.Discovery, nil
}

func (bc *Bootstrapper) getBootstrapHostsChannel(ctx context.Context, capacity int) <-chan *host.Host {
	// we need only one host to bootstrap
	bootstrapHosts := make(chan *host.Host, capacity)
	for _, bootstrapNode := range bc.cert.GetBootstrapNodes() {
		go func(ctx context.Context, address string, ch chan<- *host.Host) {
			inslogger.FromContext(ctx).Infof("Starting bootstrap to address %s", address)
			bootstrapHost, err := bc.bootstrap(address)
			if err != nil {
				inslogger.FromContext(ctx).Errorf("Error bootstrapping to address %s: %s", address, err.Error())
				return
			}
			bootstrapHosts <- bootstrapHost
		}(ctx, bootstrapNode.GetHost(), bootstrapHosts)
	}
	return bootstrapHosts
}

func (bc *Bootstrapper) waitResultFromChannel(ctx context.Context, ch <-chan *host.Host) *host.Host {
	for {
		select {
		case bootstrapHost := <-ch:
			return bootstrapHost
		case <-time.After(bc.options.BootstrapTimeout):
			inslogger.FromContext(ctx).Warn("Bootstrap timeout")
			return nil
		}
	}
}

func (bc *Bootstrapper) waitResultsFromChannel(ctx context.Context, ch <-chan *host.Host, count int) []*host.Host {
	result := make([]*host.Host, 0)
	for {
		select {
		case bootstrapHost := <-ch:
			result = append(result, bootstrapHost)
			if len(result) == count {
				return result
			}
		case <-time.After(bc.options.BootstrapTimeout):
			inslogger.FromContext(ctx).Warnf("Bootstrap timeout, successful bootstraps: %d/%d", len(result), count)
			return result
		}
	}
}

func (bc *Bootstrapper) bootstrap(address string) (*host.Host, error) {
	// TODO: add infinite bootstrap option
	bootstrapHost, err := bc.pinger.Ping(address, bc.options.PingTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to ping address %s", address)
	}
	request := bc.transport.NewRequestBuilder().Type(types.Bootstrap).Data(&NodeBootstrapRequest{}).Build()
	future, err := bc.transport.SendRequestPacket(request, bootstrapHost)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send bootstrap request to address %s", address)
	}
	response, err := future.GetResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to bootstrap request from address %s", address)
	}
	data := response.GetData().(*NodeBootstrapResponse)
	if data.Code == Rejected {
		return nil, errors.New("Rejected: " + data.RejectReason)
	}
	if data.Code == Redirected {
		return bc.bootstrap(data.RedirectHost)
	}
	return response.GetSenderHost(), nil
}

func (bc *Bootstrapper) processBootstrap(request network.Request) (network.Response, error) {
	// TODO: redirect logic
	return bc.transport.BuildResponse(request, &NodeBootstrapResponse{Code: Accepted}), nil
}

func (bc *Bootstrapper) checkGenesisCert(cert core.Certificate) error {
	// TODO: check certificate
	return nil
}

func (bc *Bootstrapper) processGenesis(request network.Request) (network.Response, error) {
	data := request.GetData().(*GenesisRequest)
	err := bc.checkGenesisCert(data.Certificate)
	if err != nil {
		return bc.transport.BuildResponse(request, &GenesisResponse{Error: err.Error()}), nil
	}
	return bc.transport.BuildResponse(request, &GenesisResponse{Discovery: bc.keeper.GetOrigin()}), nil
}

func (bc *Bootstrapper) Start() {
	bc.transport.RegisterPacketHandler(types.Bootstrap, bc.processBootstrap)
	bc.transport.RegisterPacketHandler(types.Genesis, bc.processGenesis)
}

func NewBootstrapController(options *common.Options, certificate core.Certificate, transport network.InternalTransport) *Bootstrapper {
	return &Bootstrapper{options: options, cert: certificate, transport: transport, pinger: pinger.NewPinger(transport)}
}
