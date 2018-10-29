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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/packet/types"
)

// Controller contains network logic.
type Controller interface {
	// SendMessage send message to nodeID.
	SendMessage(nodeID core.RecordRef, name string, msg core.SignedMessage) ([]byte, error)
	// RemoteProcedureRegister register remote procedure that will be executed when message is received.
	RemoteProcedureRegister(name string, method core.RemoteProcedure)
	// SendCascadeMessage sends a message from MessageBus to a cascade of nodes.
	SendCascadeMessage(data core.Cascade, method string, msg core.SignedMessage) error
	// Bootstrap init bootstrap process: 1. Connect to discovery node; 2. Reconnect to new discovery node if redirected.
	Bootstrap() error
	// AnalyzeNetwork legacy method for old DHT network (should be removed in new network).
	AnalyzeNetwork() error
	// Authorize start authorization process on discovery node.
	Authorize() error
	// ResendPulseToKnownHosts resend pulse when we receive pulse from pulsar daemon.
	ResendPulseToKnownHosts(pulse core.Pulse)

	// GetNodeID get self node id (should be removed in far future).
	GetNodeID() core.RecordRef

	// Inject inject components.
	Inject(components core.Components)
}

// RequestHandler handler function to process incoming requests from network.
type RequestHandler func(Request) (Response, error)

// HostNetwork simple interface to send network requests and process network responses.
type HostNetwork interface {
	// Listen start listening to network requests, should be started in goroutine.
	Listen() error
	// Disconnect stop listening to network requests.
	Disconnect() error
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string

	// SendRequest send request to a remote node.
	SendRequest(Request) (Future, error)
	// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
	RegisterRequestHandler(t types.PacketType, handler RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() RequestBuilder
}

// Packet is a packet that is transported via network by HostNetwork.
type Packet interface {
	GetSender() core.RecordRef
	GetReceiver() core.RecordRef
	GetType() types.PacketType
	GetData() interface{}
}

// Request is a packet that is sent from the current node.
type Request Packet

// Response is a packet that is received in response to a previously sent Request.
type Response Packet

// Future allows to handle responses to a previously sent request.
type Future interface {
	Result() <-chan Response
	GetResult(duration time.Duration) Response
}

// RequestBuilder allows to build a Request.
type RequestBuilder interface {
	Receiver(ref core.RecordRef) RequestBuilder
	Type(packetType types.PacketType) RequestBuilder
	Data(data interface{}) RequestBuilder
	Build() Request
}

// OnPulse callback function to process new pulse from pulsar.
type OnPulse func(pulse core.Pulse)

// NodeKeeper manages unsync, sync and active lists.
type NodeKeeper interface {
	core.NodeNetwork
	// AddActiveNodes add active nodes.
	AddActiveNodes([]*core.Node)
	// SetPulse sets internal PulseNumber to number. Returns true if set was successful, false if number is less
	// or equal to internal PulseNumber. If set is successful, returns collected unsync list and starts collecting new unsync list.
	SetPulse(number core.PulseNumber) (bool, UnsyncList)
	// Sync initiates transferring syncCandidates -> sync, sync -> active.
	// If number is less than internal PulseNumber then ignore Sync.
	Sync(syncCandidates []*core.Node, number core.PulseNumber)
	// AddUnsync add unsync node to the unsync list. Returns channel that receives active node on successful sync.
	// Channel will return nil node if added node has not passed the consensus.
	// Returns error if current node is not active and cannot participate in consensus.
	AddUnsync(nodeID core.RecordRef, roles []core.NodeRole, address string,
		version string /*, publicKey *ecdsa.PublicKey*/) (chan *core.Node, error)
	// GetUnsyncHolder get unsync list executed in consensus for specific pulse.
	// 1. If pulse is less than internal NodeKeeper pulse, returns error.
	// 2. If pulse is equal to internal NodeKeeper pulse, returns unsync list holder for currently executed consensus.
	// 3. If pulse is more than internal NodeKeeper pulse, blocks till next SetPulse or duration timeout and then acts like in par. 2.
	GetUnsyncHolder(pulse core.PulseNumber, duration time.Duration) (UnsyncList, error)
}

type UnsyncList interface {
	// GetUnsync returns list of local unsync nodes. This list is created.
	GetUnsync() []*core.Node
	// GetPulse returns actual pulse for current consensus process.
	GetPulse() core.PulseNumber
	// SetHash sets hash of unsync lists for each node of consensus.
	SetHash([]*NodeUnsyncHash)
	// GetHash get hash of unsync lists for each node of  If hash is not calculated yet, then this call blocks
	// until the hash is calculated with SetHash() call.
	GetHash(blockTimeout time.Duration) ([]*NodeUnsyncHash, error)
	// AddUnsyncList add unsync list for remote ref.
	AddUnsyncList(ref core.RecordRef, unsync []*core.Node)
	// AddUnsyncHash add unsync hash for remote ref.
	AddUnsyncHash(ref core.RecordRef, hash []*NodeUnsyncHash)
	// GetUnsyncList get unsync list for remote ref.
	GetUnsyncList(ref core.RecordRef) ([]*core.Node, bool)
	// GetUnsyncHash get unsync hash for remote ref.
	GetUnsyncHash(ref core.RecordRef) ([]*NodeUnsyncHash, bool)
}

// NodeUnsyncHash data needed for consensus.
type NodeUnsyncHash struct {
	NodeID core.RecordRef
	Hash   []byte
	// TODO: add signature
}
