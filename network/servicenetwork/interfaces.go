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

package servicenetwork

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/packet"
)

type NetworkController interface {
	RemoteProcedureRegister(name string, method core.RemoteProcedure)
	SendMessage(nodeID core.RecordRef, method string, msg core.Message) ([]byte, error)
	Bootstrap() error
	AnalyzeNetwork() error
	Authorize() error

	SetNodeKeeper(keeper consensus.NodeKeeper)
}

// RequestHandler handler function to process incoming requests from network
type RequestHandler func(request *packet.Packet) (response *packet.Packet, err error)

// HostNetwork simple interface to send network requests and process network responses
type HostNetwork interface {
	// Listen start listening to network requests, should be started in goroutine
	Listen() error
	// Disconnect stop listening to network requests
	Disconnect() error
	// PublicAddress returns public address that can be published for all nodes
	PublicAddress() string

	// SendRequest send request to a remote node
	SendRequest(*packet.Packet) (transport.Future, error)
	// RegisterRequestHandler register a handler function to process incoming requests of a specific type
	RegisterRequestHandler(t packet.PacketType, handler RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node
	NewRequestBuilder() *packet.Builder
}
