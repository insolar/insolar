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

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
)

type PulseController interface {
	component.Starter
}

type pulseController struct {
	PulseHandler network.PulseHandler `inject:""`

	hostNetwork  network.HostNetwork
	routingTable network.RoutingTable
}

func (pc *pulseController) Start(ctx context.Context) error {
	pc.hostNetwork.RegisterRequestHandler(types.Pulse, pc.processPulse)
	pc.hostNetwork.RegisterRequestHandler(types.GetRandomHosts, pc.processGetRandomHosts)
	return nil
}

func (pc *pulseController) processPulse(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*packet.RequestPulse)
	go pc.PulseHandler.HandlePulse(context.Background(), data.Pulse)
	return pc.hostNetwork.BuildResponse(ctx, request, &packet.ResponsePulse{Success: true, Error: ""}), nil
}

func (pc *pulseController) processGetRandomHosts(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*packet.RequestGetRandomHosts)
	randomHosts := pc.routingTable.GetRandomNodes(data.HostsNumber)
	return pc.hostNetwork.BuildResponse(ctx, request, &packet.ResponseGetRandomHosts{Hosts: randomHosts}), nil
}

func NewPulseController(hostNetwork network.HostNetwork, routingTable network.RoutingTable) *pulseController {
	return &pulseController{hostNetwork: hostNetwork, routingTable: routingTable}
}
