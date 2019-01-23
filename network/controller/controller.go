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
	Bootstrapper  bootstrap.NetworkBootstrapper `inject:"subcomponent"`
	RPCController RPCController                 `inject:"subcomponent"`

	network network.HostNetwork
}

func (c *Controller) SetLastIgnoredPulse(number core.PulseNumber) {
	c.Bootstrapper.SetLastPulse(number)
}

func (c *Controller) GetLastIgnoredPulse() core.PulseNumber {
	return c.Bootstrapper.GetLastPulse()
}

// SendParcel send message to nodeID.
func (c *Controller) SendMessage(nodeID core.RecordRef, name string, msg core.Parcel) ([]byte, error) {
	return c.RPCController.SendMessage(nodeID, name, msg)
}

// RemoteProcedureRegister register remote procedure that will be executed when message is received.
func (c *Controller) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	c.RPCController.RemoteProcedureRegister(name, method)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes.
func (c *Controller) SendCascadeMessage(data core.Cascade, method string, msg core.Parcel) error {
	return c.RPCController.SendCascadeMessage(data, method, msg)
}

// Bootstrap init bootstrap process: 1. Connect to discovery node; 2. Reconnect to new discovery node if redirected.
func (c *Controller) Bootstrap(ctx context.Context) error {
	return c.Bootstrapper.Bootstrap(ctx)
}

// Inject inject components.
func (c *Controller) Start(ctx context.Context) error {
	c.network.RegisterRequestHandler(types.Ping, func(ctx context.Context, request network.Request) (network.Response, error) {
		return c.network.BuildResponse(ctx, request, nil), nil
	})
	return nil
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
func NewNetworkController(net network.HostNetwork) network.Controller {
	return &Controller{network: net}
}
