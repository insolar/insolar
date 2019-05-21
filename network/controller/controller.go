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

package controller

import (
	"context"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/bootstrap"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

// Controller contains network logic.
type Controller struct {
	Bootstrapper  bootstrap.NetworkBootstrapper `inject:""`
	RPCController RPCController                 `inject:""`
	Network       network.HostNetwork           `inject:""`
}

func (c *Controller) SetLastIgnoredPulse(number insolar.PulseNumber) {
	c.Bootstrapper.SetLastPulse(number)
}

func (c *Controller) GetLastIgnoredPulse() insolar.PulseNumber {
	return c.Bootstrapper.GetLastPulse()
}

// SendParcel send message to nodeID.
func (c *Controller) SendMessage(nodeID insolar.Reference, name string, msg insolar.Parcel) ([]byte, error) {
	return c.RPCController.SendMessage(nodeID, name, msg)
}

// SendBytes send message to nodeID.
func (c *Controller) SendBytes(ctx context.Context, nodeID insolar.Reference, name string, msgBytes []byte) ([]byte, error) {
	return c.RPCController.SendBytes(ctx, nodeID, name, msgBytes)
}

// RemoteProcedureRegister register remote procedure that will be executed when message is received.
func (c *Controller) RemoteProcedureRegister(name string, method insolar.RemoteProcedure) {
	c.RPCController.RemoteProcedureRegister(name, method)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes.
func (c *Controller) SendCascadeMessage(data insolar.Cascade, method string, msg insolar.Parcel) error {
	return c.RPCController.SendCascadeMessage(data, method, msg)
}

// Bootstrap init bootstrap process: 1. Connect to discovery node; 2. Reconnect to new discovery node if redirected.
func (c *Controller) Bootstrap(ctx context.Context) (*network.BootstrapResult, error) {
	return c.Bootstrapper.Bootstrap(ctx)
}

// Inject inject components.
func (c *Controller) Init(ctx context.Context) error {
	c.Network.RegisterRequestHandler(types.Ping, func(ctx context.Context, request network.Packet) (network.Packet, error) {
		return c.Network.BuildResponse(ctx, request, nil), nil
	})
	return nil
}

func (c *Controller) AuthenticateToDiscoveryNode(ctx context.Context, discovery insolar.DiscoveryNode) error {
	return c.Bootstrapper.AuthenticateToDiscoveryNode(ctx, nil)
}

// ConfigureOptions convert daemon configuration to controller options
func ConfigureOptions(conf configuration.Configuration) *common.Options {
	config := conf.Host
	return &common.Options{
		InfinityBootstrap:      config.InfinityBootstrap,
		TimeoutMult:            time.Duration(config.TimeoutMult) * time.Second,
		MinTimeout:             time.Duration(config.MinTimeout) * time.Second,
		MaxTimeout:             time.Duration(config.MaxTimeout) * time.Second,
		PingTimeout:            1 * time.Second,
		PacketTimeout:          10 * time.Second,
		BootstrapTimeout:       10 * time.Second,
		HandshakeSessionTTL:    time.Duration(config.HandshakeSessionTTL) * time.Millisecond,
		FakePulseDuration:      time.Duration(conf.Pulsar.PulseTime) * time.Millisecond,
		CyclicBootstrapEnabled: false,
	}
}

// NewNetworkController create new network controller.
func NewNetworkController() network.Controller {
	return &Controller{}
}
