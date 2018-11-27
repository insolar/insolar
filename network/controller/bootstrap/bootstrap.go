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

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/controller/pinger"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type BootstrapController struct {
	options   *common.Options
	transport network.InternalTransport
	pinger    *pinger.Pinger

	chosenDiscoveryNode *host.Host
}

type BootstrapRequest struct {
	// pass node certificate, i guess
}

type BootstrapResponse struct {
	Code         BootstrapCode
	RedirectHost string
	RejectReason string
}

type BootstrapCode uint8

const (
	BootstrapAccepted = BootstrapCode(iota + 1)
	BootstrapRejected
	BootstrapRedirected
)

func init() {
	gob.Register(&BootstrapRequest{})
	gob.Register(&BootstrapResponse{})
}

func (bc *BootstrapController) GetChosenDiscoveryNode() *host.Host {
	return bc.chosenDiscoveryNode
}

func (bc *BootstrapController) Bootstrap(ctx context.Context) error {
	if len(bc.options.BootstrapHosts) == 0 {
		return nil
	}
	ch := bc.getBootstrapHostsChannel(ctx)
	host := bc.waitResultFromChannel(ctx, ch)

	if host == nil {
		return errors.New("Failed to bootstrap to any of discovery nodes")
	}
	bc.chosenDiscoveryNode = host
	return nil
}

func (bc *BootstrapController) getBootstrapHostsChannel(ctx context.Context) <-chan *host.Host {
	// we need only one host to bootstrap
	bootstrapHosts := make(chan *host.Host, 1)
	for _, bootstrapAddress := range bc.options.BootstrapHosts {
		go func(ctx context.Context, address string, ch chan<- *host.Host) {
			inslogger.FromContext(ctx).Infof("Starting bootstrap to address %s", address)
			bootstrapHost, err := bc.bootstrap(address)
			if err != nil {
				inslogger.FromContext(ctx).Errorf("Error bootstrapping to address %s: %s", address, err.Error())
				return
			}
			bootstrapHosts <- bootstrapHost
		}(ctx, bootstrapAddress, bootstrapHosts)
	}
	return bootstrapHosts
}

func (bc *BootstrapController) waitResultFromChannel(ctx context.Context, ch <-chan *host.Host) *host.Host {
	for {
		select {
		case bootstrapHost := <-ch:
			return bootstrapHost
		case <-time.After(bc.options.BootstrapTimeout):
			inslogger.FromContext(ctx).Warnf("Bootstrap timeout")
			return nil
		}
	}
}

func (bc *BootstrapController) bootstrap(address string) (*host.Host, error) {
	// TODO: add infinite bootstrap
	bootstrapHost, err := bc.pinger.Ping(address, bc.options.PingTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to ping address %s", address)
	}
	request := bc.transport.NewRequestBuilder().Type(types.Bootstrap).Data(&BootstrapRequest{}).Build()
	future, err := bc.transport.SendRequestPacket(request, bootstrapHost)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to ping address %s", address)
	}
	response, err := future.GetResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get response to bootstrap request")
	}
	data := response.GetData().(*BootstrapResponse)
	if data.Code == BootstrapRejected {
		return nil, errors.New("Rejected: " + data.RejectReason)
	}
	if data.Code == BootstrapRedirected {
		return bc.bootstrap(data.RedirectHost)
	}
	return response.GetSenderHost(), nil
}

func (bc *BootstrapController) processBootstrap(request network.Request) (network.Response, error) {
	// TODO: check certificate and redirect logic
	return bc.transport.BuildResponse(request, &BootstrapResponse{Code: BootstrapAccepted}), nil
}

func (bc *BootstrapController) Start() {
	bc.transport.RegisterPacketHandler(types.Bootstrap, bc.processBootstrap)
}

func NewBootstrapController(options *common.Options, transport network.InternalTransport) *BootstrapController {
	return &BootstrapController{options: options, transport: transport, pinger: pinger.NewPinger(transport)}
}
