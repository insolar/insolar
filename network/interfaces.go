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
	"context"
	"time"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
)

// Controller contains network logic.
type Controller interface {
	// SendParcel send message to nodeID.
	SendMessage(nodeID core.RecordRef, name string, msg core.Parcel) ([]byte, error)
	// RemoteProcedureRegister register remote procedure that will be executed when message is received.
	RemoteProcedureRegister(name string, method core.RemoteProcedure)
	// SendCascadeMessage sends a message from MessageBus to a cascade of nodes.
	SendCascadeMessage(data core.Cascade, method string, msg core.Parcel) error
	// Bootstrap init complex bootstrap process. Blocks until bootstrap is complete.
	Bootstrap(ctx context.Context) error

	// Inject inject components.
	Inject(cryptographyService core.CryptographyService,
		networkCoordinator core.NetworkCoordinator, nodeKeeper NodeKeeper)
}

// RequestHandler handler function to process incoming requests from network.
type RequestHandler func(Request) (Response, error)

// HostNetwork simple interface to send network requests and process network responses.
//go:generate minimock -i github.com/insolar/insolar/network.HostNetwork -o ../testutils/network -s _mock.go
type HostNetwork interface {
	// Start listening to network requests.
	Start(ctx context.Context)
	// Stop listening to network requests.
	Stop()
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string
	// GetNodeID get current node ID.
	GetNodeID() core.RecordRef

	// SendRequest send request to a remote node.
	SendRequest(request Request, receiver core.RecordRef) (Future, error)
	// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
	RegisterRequestHandler(t types.PacketType, handler RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() RequestBuilder
	// BuildResponse create response to an incoming request with Data set to responseData.
	BuildResponse(request Request, responseData interface{}) Response
}

type ConsensusRequestHandler func(Request)

//go:generate minimock -i github.com/insolar/insolar/network.ConsensusNetwork -o ../testutils/network -s _mock.go
type ConsensusNetwork interface {
	// Start listening to network requests.
	Start(ctx context.Context)
	// Stop listening to network requests.
	Stop()
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string
	// GetNodeID get current node ID.
	GetNodeID() core.RecordRef

	// SendRequest send request to a remote node.
	SendRequest(request Request, receiver core.RecordRef) error
	// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
	RegisterRequestHandler(t types.PacketType, handler ConsensusRequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() RequestBuilder
}

// Packet is a packet that is transported via network by HostNetwork.
type Packet interface {
	GetSender() core.RecordRef
	GetSenderHost() *host.Host
	GetType() types.PacketType
	GetData() interface{}
}

// Request is a packet that is sent from the current node.
type Request Packet

// Response is a packet that is received in response to a previously sent Request.
type Response Packet

// Future allows to handle responses to a previously sent request.
type Future interface {
	GetRequest() Request
	Response() <-chan Response
	GetResponse(duration time.Duration) (Response, error)
}

// RequestBuilder allows to build a Request.
type RequestBuilder interface {
	Type(packetType types.PacketType) RequestBuilder
	Data(data interface{}) RequestBuilder
	Build() Request
}

// PulseHandler interface to process new pulse.
//go:generate minimock -i github.com/insolar/insolar/network.PulseHandler -o ../testutils/network -s _mock.go
type PulseHandler interface {
	HandlePulse(ctx context.Context, pulse core.Pulse)
}

type NodeKeeperState uint8

const (
	// Undefined is state of NodeKeeper while it is not valid
	Undefined NodeKeeperState = iota + 1
	// Waiting is state of NodeKeeper while it is not part of consensus yet (waits for its join claim to pass)
	Waiting
	// Ready is state of NodeKeeper when it is ready for consensus
	Ready
)

// NodeKeeper manages unsync, sync and active lists.
//go:generate minimock -i github.com/insolar/insolar/network.NodeKeeper -o ../testutils/network -s _mock.go
type NodeKeeper interface {
	core.NodeNetwork

	// TODO: remove this interface when bootstrap mechanism completed
	core.SwitcherWorkAround

	// SetCloudHash set new cloud hash
	SetCloudHash([]byte)
	// AddActiveNodes add active nodes.
	AddActiveNodes([]core.Node)
	// GetActiveNodeByShortID get active node by short ID. Returns nil if node is not found.
	GetActiveNodeByShortID(shortID core.ShortNodeID) core.Node
	// SetState set state of the NodeKeeper
	SetState(NodeKeeperState)
	// GetState get state of the NodeKeeper
	GetState() NodeKeeperState
	// GetOriginClaim get origin NodeJoinClaim
	GetOriginClaim() (*consensus.NodeJoinClaim, error)
	// NodesJoinedDuringPreviousPulse returns true if the last Sync call contained approved Join claims
	NodesJoinedDuringPreviousPulse() bool
	// AddPendingClaim add pending claim to the internal queue of claims
	AddPendingClaim(consensus.ReferendumClaim) bool
	// GetClaimQueue get the internal queue of claims
	GetClaimQueue() ClaimQueue
	// GetUnsyncList get unsync list for current pulse. Has copy of active node list from nodekeeper as internal state.
	// Should be called when nodekeeper state is Ready.
	GetUnsyncList() UnsyncList
	// GetSparseUnsyncList get sparse unsync list for current pulse with predefined length of active node list.
	// Does not contain active list, should collect active list during its lifetime via AddClaims.
	// Should be called when nodekeeper state is Waiting.
	GetSparseUnsyncList(length int) UnsyncList
	// Sync move unsync -> sync
	Sync(list UnsyncList)
	// MoveSyncToActive merge sync list with active nodes
	MoveSyncToActive()
}

// UnsyncList is interface to manage unsync list
//go:generate minimock -i github.com/insolar/insolar/network.UnsyncList -o ../testutils/network -s _mock.go
type UnsyncList interface {
	consensus.BitSetMapper
	// RemoveClaims
	RemoveClaims(core.RecordRef)
	// AddClaims
	AddClaims(map[core.RecordRef][]consensus.ReferendumClaim, map[core.RecordRef]string)
	// CalculateHash calculate node list hash based on active node list and claims
	CalculateHash() ([]byte, error)
	// GetActiveNode get active node by reference ID for current consensus
	GetActiveNode(ref core.RecordRef) core.Node
	// GetActiveNodes get active nodes for current consensus
	GetActiveNodes() []core.Node
}

// PartitionPolicy contains all rules how to initiate globule resharding.
type PartitionPolicy interface {
	ShardsCount() int
}

// RoutingTable contains all routing information of the network.
type RoutingTable interface {
	// Inject inject dependencies from components
	Inject(nodeKeeper NodeKeeper)
	// Resolve NodeID -> ShortID, Address. Can initiate network requests.
	Resolve(core.RecordRef) (*host.Host, error)
	// ResolveS ShortID -> NodeID, Address for node inside current globe.
	ResolveS(core.ShortNodeID) (*host.Host, error)
	// AddToKnownHosts add host to routing table.
	AddToKnownHosts(*host.Host)
	// Rebalance recreate shards of routing table with known hosts according to new partition policy.
	Rebalance(PartitionPolicy)
	// GetRandomNodes get a specified number of random nodes. Returns less if there are not enough nodes in network.
	GetRandomNodes(count int) []host.Host
}

// InternalTransport simple interface to send network requests and process network responses.
type InternalTransport interface {
	// Start listening to network requests, should be started in goroutine.
	Start(ctx context.Context)
	// Stop listening to network requests.
	Stop()
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string
	// GetNodeID get current node ID.
	GetNodeID() core.RecordRef

	// SendRequestPacket send request packet to a remote node.
	SendRequestPacket(request Request, receiver *host.Host) (Future, error)
	// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
	RegisterPacketHandler(t types.PacketType, handler RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() RequestBuilder
	// BuildResponse create response to an incoming request with Data set to responseData.
	BuildResponse(request Request, responseData interface{}) Response
}

// ClaimQueue is the queue that contains consensus claims.
//go:generate minimock -i github.com/insolar/insolar/network.ClaimQueue -o ../testutils/network -s _mock.go
type ClaimQueue interface {
	// Pop takes claim from the queue.
	Pop() consensus.ReferendumClaim
	// Front returns claim from the queue without removing it from the queue.
	Front() consensus.ReferendumClaim
	// Length returns the length of the queue
	Length() int
}
