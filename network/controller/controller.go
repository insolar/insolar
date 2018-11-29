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

package controller

import (
	"context"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/bootstrap"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/packet/types"
)

// Controller contains network logic.
type Controller struct {
	options *common.Options
	network network.HostNetwork

	bootstrapper    *bootstrap.NetworkBootstrapper
	pulseController *PulseController
	rpcController   *RPCController
}

// SendParcel send message to nodeID.
func (c *Controller) SendMessage(nodeID core.RecordRef, name string, msg core.Parcel) ([]byte, error) {
	return c.rpcController.SendMessage(nodeID, name, msg)
}

// RemoteProcedureRegister register remote procedure that will be executed when message is received.
func (c *Controller) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	c.rpcController.RemoteProcedureRegister(name, method)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes.
func (c *Controller) SendCascadeMessage(data core.Cascade, method string, msg core.Parcel) error {
	return c.rpcController.SendCascadeMessage(data, method, msg)
}

// Bootstrap init bootstrap process: 1. Connect to discovery node; 2. Reconnect to new discovery node if redirected.
func (c *Controller) Bootstrap(ctx context.Context) error {
	return c.bootstrapper.Bootstrap(ctx)
}

// Inject inject components.
func (c *Controller) Inject(cryptographyService core.CryptographyService,
	networkCoordinator core.NetworkCoordinator, nodeKeeper network.NodeKeeper) {

	c.network.RegisterRequestHandler(types.Ping, func(request network.Request) (network.Response, error) {
		return c.network.BuildResponse(request, nil), nil
	})
	c.bootstrapper.Start(cryptographyService, networkCoordinator, nodeKeeper)
	c.pulseController.Start()
	c.rpcController.Start()
}

// ConfigureOptions convert daemon configuration to controller options
func ConfigureOptions(config configuration.HostNetwork) *common.Options {
	options := &common.Options{}
	if options.PingTimeout == 0 {
		options.PingTimeout = time.Second * 1
	}
	if options.PacketTimeout == 0 {
		options.PacketTimeout = time.Second * 10
	}
	if options.BootstrapTimeout == 0 {
		options.BootstrapTimeout = time.Second * 10
	}
	return options
}

// NewNetworkController create new network controller.
func NewNetworkController(
	pulseHandler network.PulseHandler,
	options *common.Options,
	certificate core.Certificate,
	transport network.InternalTransport,
	routingTable network.RoutingTable,
	network network.HostNetwork,
	scheme core.PlatformCryptographyScheme) network.Controller {

	c := Controller{}
	c.network = network
	c.options = options
	c.bootstrapper = bootstrap.NewNetworkBootstrapper(c.options, certificate, transport)
	c.pulseController = NewPulseController(pulseHandler, network, routingTable)
	c.rpcController = NewRPCController(c.options, network, scheme)

	return &c
}
