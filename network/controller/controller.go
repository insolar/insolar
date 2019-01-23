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

func (c *Controller) SetLastIgnoredPulse(number core.PulseNumber) {
	c.bootstrapper.SetLastPulse(number)
}

func (c *Controller) GetLastIgnoredPulse() core.PulseNumber {
	return c.bootstrapper.GetLastPulse()
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

	c.network.RegisterRequestHandler(types.Ping, func(ctx context.Context, request network.Request) (network.Response, error) {
		return c.network.BuildResponse(ctx, request, nil), nil
	})
	c.bootstrapper.Start(cryptographyService, networkCoordinator, nodeKeeper)
	c.pulseController.Start()
	c.rpcController.Start()
}

// ConfigureOptions convert daemon configuration to controller options
func ConfigureOptions(config configuration.HostNetwork) *common.Options {
	return &common.Options{
		InfinityBootstrap:   config.InfinityBootstrap,
		TimeoutMult:         time.Duration(config.TimeoutMult) * time.Second,
		MinTimeout:          time.Duration(config.MinTimeout) * time.Second,
		MaxTimeout:          time.Duration(config.MaxTimeout) * time.Second,
		PingTimeout:         1 * time.Second,
		PacketTimeout:       10 * time.Second,
		BootstrapTimeout:    10 * time.Second,
		HandshakeSessionTTL: time.Duration(config.HandshakeSessionTTL) * time.Millisecond,
	}
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
