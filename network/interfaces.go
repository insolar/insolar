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

package network

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/packet"
)

// Controller contains network logic
type Controller interface {
	// SendMessage send message to nodeID
	SendMessage(nodeID core.RecordRef, name string, msg core.Message) ([]byte, error)
	// RemoteProcedureRegister register remote procedure that will be executed when message is received
	RemoteProcedureRegister(name string, method core.RemoteProcedure)
	// SendCascadeMessage sends a message from MessageBus to a cascade of nodes
	SendCascadeMessage(data core.Cascade, method string, msg core.Message) error
	// Bootstrap init bootstrap process: 1. Connect to discovery node; 2. Reconnect to new discovery node if redirected.
	Bootstrap() error
	// AnalyzeNetwork legacy method for old DHT network (should be removed in
	AnalyzeNetwork() error
	// Authorize start authorization process on discovery node.
	Authorize() error
	// ResendPulseToKnownHosts resend pulse when we receive pulse from pulsar daemon
	ResendPulseToKnownHosts(pulse core.Pulse)

	// GetNodeID get self node id (should be removed in far future)
	GetNodeID() core.RecordRef

	// Inject inject components
	Inject(components core.Components)

	// GetConsensus get consensus processor. Should be deleted and implemented on the same level as network.Controller
	GetConsensus() consensus.Processor
}

// RequestHandler handler function to process incoming requests from network.
type RequestHandler func(request *packet.Packet) (response *packet.Packet, err error)

// HostNetwork simple interface to send network requests and process network responses.
type HostNetwork interface {
	// Listen start listening to network requests, should be started in goroutine.
	Listen() error
	// Disconnect stop listening to network requests.
	Disconnect() error
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string

	// SendRequest send request to a remote node.
	SendRequest(*packet.Packet) (transport.Future, error)
	// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
	RegisterRequestHandler(t packet.PacketType, handler RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() *packet.Builder
}

// OnPulse callback function to process new pulse from pulsar
type OnPulse func(pulse core.Pulse)
