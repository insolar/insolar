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
func (tr *TransportResolvable) Start(ctx context.Context) error {
	return tr.internalTransport.Start(ctx)
}

// Stop listening to network requests.
func (tr *TransportResolvable) Stop(ctx context.Context) error {
	return tr.internalTransport.Stop(ctx)
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
func (tr *TransportResolvable) SendRequest(ctx context.Context, request network.Request, receiver core.RecordRef) (network.Future, error) {
	h, err := tr.resolver.Resolve(receiver)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving NodeID -> Address")
	}
	return tr.internalTransport.SendRequestPacket(ctx, request, h)
}

// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
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
func (tr *TransportResolvable) BuildResponse(ctx context.Context, request network.Request, responseData interface{}) network.Response {
	return tr.internalTransport.BuildResponse(ctx, request, responseData)
}

func NewHostTransport(transport network.InternalTransport, resolver network.RoutingTable) network.HostNetwork {
	return &TransportResolvable{internalTransport: transport, resolver: resolver}
}
