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

package hostnetwork

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

// TransportResolvable is implementation of HostNetwork interface that is capable of address resolving.
type TransportResolvable struct {
	internalTransport network.InternalTransport
	resolver          network.RoutingTable
}

// Start listening to network requests.
func (tr *TransportResolvable) Start(ctx context.Context) {
	tr.internalTransport.Start(ctx)
}

// Stop listening to network requests.
func (tr *TransportResolvable) Stop() {
	tr.internalTransport.Stop()
}

// PublicAddress returns public address that can be published for all nodes.
func (tr *TransportResolvable) PublicAddress() string {
	return tr.internalTransport.PublicAddress()
}

// GetNodeID get current node ID.
func (tr *TransportResolvable) GetNodeID() core.RecordRef {
	return tr.internalTransport.GetNodeID()
}

// SendRequest send request to a remote node.
func (tr *TransportResolvable) SendRequest(request network.Request, receiver core.RecordRef) (network.Future, error) {
	h, err := tr.resolver.Resolve(receiver)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving NodeID -> Address")
	}
	return tr.internalTransport.SendRequestPacket(request, h)
}

// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
func (tr *TransportResolvable) RegisterRequestHandler(t types.PacketType, handler network.RequestHandler) {
	f := func(ctx context.Context, request network.Request) (network.Response, error) {
		tr.resolver.AddToKnownHosts(request.GetSenderHost())
		return handler(ctx, request)
	}
	tr.internalTransport.RegisterPacketHandler(t, f)
}

// NewRequestBuilder create packet Builder for an outgoing request with sender set to current node.
func (tr *TransportResolvable) NewRequestBuilder() network.RequestBuilder {
	return tr.internalTransport.NewRequestBuilder()
}

// BuildResponse create response to an incoming request with Data set to responseData.
func (tr *TransportResolvable) BuildResponse(request network.Request, responseData interface{}) network.Response {
	return tr.internalTransport.BuildResponse(request, responseData)
}

func NewHostTransport(transport network.InternalTransport, resolver network.RoutingTable) network.HostNetwork {
	return &TransportResolvable{internalTransport: transport, resolver: resolver}
}
