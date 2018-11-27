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

type Bootstrapper struct {
	options   *common.Options
	transport network.InternalTransport
	pinger    *pinger.Pinger

	chosenDiscoveryNode *host.Host
}

type NodeBootstrapRequest struct {
	// pass node certificate, i guess
}

type NodeBootstrapResponse struct {
	Code         Code
	RedirectHost string
	RejectReason string
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
}

func (bc *Bootstrapper) GetChosenDiscoveryNode() *host.Host {
	return bc.chosenDiscoveryNode
}

func (bc *Bootstrapper) Bootstrap(ctx context.Context) error {
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

func (bc *Bootstrapper) getBootstrapHostsChannel(ctx context.Context) <-chan *host.Host {
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

func (bc *Bootstrapper) waitResultFromChannel(ctx context.Context, ch <-chan *host.Host) *host.Host {
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

func (bc *Bootstrapper) bootstrap(address string) (*host.Host, error) {
	// TODO: add infinite bootstrap
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
	// TODO: check certificate and redirect logic
	return bc.transport.BuildResponse(request, &NodeBootstrapResponse{Code: Accepted}), nil
}

func (bc *Bootstrapper) Start() {
	bc.transport.RegisterPacketHandler(types.Bootstrap, bc.processBootstrap)
}

func NewBootstrapController(options *common.Options, transport network.InternalTransport) *Bootstrapper {
	return &Bootstrapper{options: options, transport: transport, pinger: pinger.NewPinger(transport)}
}
