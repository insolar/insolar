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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
)

// InternalTransport simple interface to send network requests and process network responses.
type InternalTransport interface {
	// Start listening to network requests, should be started in goroutine.
	Start()
	// Stop listening to network requests.
	Stop()
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string
	// GetNodeID get current node ID.
	GetNodeID() core.RecordRef

	// SendRequestPacket send request packet to a remote node.
	SendRequestPacket(request network.Request, receiver *host.Host) (network.Future, error)
	// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
	RegisterPacketHandler(t types.PacketType, handler network.RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() network.RequestBuilder
	// BuildResponse create response to an incoming request with Data set to responseData.
	BuildResponse(request network.Request, responseData interface{}) network.Response
}
