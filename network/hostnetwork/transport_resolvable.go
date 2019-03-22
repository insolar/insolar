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

package hostnetwork

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

// TransportResolvable is implementation of HostNetwork interface that is capable of address resolving.
type TransportResolvable struct {
	Resolver  network.RoutingTable      `inject:""`
	Transport network.InternalTransport `inject:""`
}

// PublicAddress returns public address that can be published for all nodes.
func (tr *TransportResolvable) PublicAddress() string {
	return tr.Transport.PublicAddress()
}

// GetNodeID get current node ID.
func (tr *TransportResolvable) GetNodeID() insolar.Reference {
	return tr.Transport.GetNodeID()
}

// SendRequest send request to a remote node.
func (tr *TransportResolvable) SendRequest(ctx context.Context, request network.Request, receiver insolar.Reference) (network.Future, error) {
	h, err := tr.Resolver.Resolve(receiver)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving NodeID -> Address")
	}
	return tr.Transport.SendRequestPacket(ctx, request, h)
}

// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
func (tr *TransportResolvable) RegisterRequestHandler(t types.PacketType, handler network.RequestHandler) {
	f := func(ctx context.Context, request network.Request) (network.Response, error) {
		tr.Resolver.AddToKnownHosts(request.GetSenderHost())
		return handler(ctx, request)
	}
	tr.Transport.RegisterPacketHandler(t, f)
}

// NewRequestBuilder create packet Builder for an outgoing request with sender set to current node.
func (tr *TransportResolvable) NewRequestBuilder() network.RequestBuilder {
	return tr.Transport.NewRequestBuilder()
}

// BuildResponse create response to an incoming request with Data set to responseData.
func (tr *TransportResolvable) BuildResponse(ctx context.Context, request network.Request, responseData interface{}) network.Response {
	return tr.Transport.BuildResponse(ctx, request, responseData)
}

func NewHostTransport(transport network.InternalTransport) network.HostNetwork {
	return &TransportResolvable{Transport: transport}
}
