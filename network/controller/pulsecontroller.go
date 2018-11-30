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

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
)

type PulseController struct {
	pulseHandler network.PulseHandler
	hostNetwork  network.HostNetwork
	routingTable network.RoutingTable
}

func (pc *PulseController) Start() {
	pc.hostNetwork.RegisterRequestHandler(types.Pulse, pc.processPulse)
	pc.hostNetwork.RegisterRequestHandler(types.GetRandomHosts, pc.processGetRandomHosts)
}

func (pc *PulseController) processPulse(request network.Request) (network.Response, error) {
	data := request.GetData().(*packet.RequestPulse)
	go pc.pulseHandler.HandlePulse(context.Background(), data.Pulse)
	return pc.hostNetwork.BuildResponse(request, &packet.ResponsePulse{Success: true, Error: ""}), nil
}

func (pc *PulseController) processGetRandomHosts(request network.Request) (network.Response, error) {
	data := request.GetData().(*packet.RequestGetRandomHosts)
	randomHosts := pc.routingTable.GetRandomNodes(data.HostsNumber)
	return pc.hostNetwork.BuildResponse(request, &packet.ResponseGetRandomHosts{Hosts: randomHosts}), nil
}

func NewPulseController(pulseHandler network.PulseHandler, hostNetwork network.HostNetwork,
	routingTable network.RoutingTable) *PulseController {
	return &PulseController{pulseHandler: pulseHandler, hostNetwork: hostNetwork, routingTable: routingTable}
}
