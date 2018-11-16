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

	"github.com/insolar/insolar/consensus/packets"
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
	// Bootstrap init bootstrap process: 1. Connect to discovery node; 2. Reconnect to new discovery node if redirected.
	Bootstrap(ctx context.Context) error
	// Authorize start authorization process on discovery node.
	Authorize(ctx context.Context) error
	// ResendPulseToKnownHosts resend pulse when we receive pulse from pulsar daemon.
	// DEPRECATED
	ResendPulseToKnownHosts(pulse core.Pulse)

	// GetNodeID get self node id (should be removed in far future).
	GetNodeID() core.RecordRef

	// Inject inject components.
	Inject(components core.Components)
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
	// Waiting is state of NodeKeeper while it is not part of consensus yet (waits for his join claim to pass)
	Waiting
	// Ready is state of NodeKeeper when it is ready for consensus
	Ready
)

// NodeKeeper manages unsync, sync and active lists.
type NodeKeeper interface {
	core.NodeNetwork
	// AddActiveNodes add active nodes.
	AddActiveNodes([]core.Node)
	// GetActiveNodeByShortID get active node by short ID. Returns nil if node is not found.
	GetActiveNodeByShortID(shortID core.ShortNodeID) core.Node

	// SetState set state of the NodeKeeper
	SetState(NodeKeeperState)
	// GetState get state of the NodeKeeper
	GetState() NodeKeeperState
	// SetOriginClaim set origin NodeJoinClaim. It is needed to join to discovery node or (sometimes) in consensus
	SetOriginClaim(*packets.NodeJoinClaim)
	// GetOriginClaim get origin NodeJoinClaim
	GetOriginClaim() *packets.NodeJoinClaim
	// IsNodeJoinedOnPreviousPulse returns true if the last Sync call contained approved Join claims
	IsNodeJoinedOnPreviousPulse() bool
	// AddClaim add claim to the internal queue of claims
	AddClaim(packets.ReferendumClaim) bool
	// GetClaimQueue get the internal queue of claims
	GetClaimQueue() ClaimQueue
	// Sync move approved claims to sync list
	Sync(approvedClaims []packets.ReferendumClaim, nodesLost []core.RecordRef)
	// MoveSyncToActive merge sync list with active nodes
	MoveSyncToActive(number core.PulseNumber)

	// FIXME: DEPRECATED METHODS
	// SetPulse sets internal PulseNumber to number. Returns true if set was successful, false if number is less
	// or equal to internal PulseNumber. If set is successful, returns collected unsync list and starts collecting new unsync list.
	SetPulse(number core.PulseNumber) (bool, UnsyncList)
	// SyncOld initiates transferring syncCandidates -> sync, sync -> active.
	// If number is less than internal PulseNumber then ignore SyncOld.
	SyncOld(syncCandidates []core.Node, number core.PulseNumber)
	// AddUnsync add unsync node to the unsync list. Returns channel that receives active node on successful sync.
	// Channel will return nil node if added node has not passed the consensus.
	// Returns error if current node is not active and cannot participate in consensus.
	AddUnsync(nodeID core.RecordRef, roles []core.NodeRole, address string,
		version string /*, publicKey *ecdsa.PublicKey*/) (chan core.Node, error)
	// GetUnsyncHolder get unsync list executed in consensus for specific pulse.
	// 1. If pulse is less than internal NodeKeeper pulse, returns error.
	// 2. If pulse is equal to internal NodeKeeper pulse, returns unsync list holder for currently executed consensus.
	// 3. If pulse is more than internal NodeKeeper pulse, blocks till next SetPulse or duration timeout and then acts like in par. 2.
	GetUnsyncHolder(pulse core.PulseNumber, duration time.Duration) (UnsyncList, error)
}

type UnsyncList interface {
	// GetUnsync returns list of local unsync nodes. This list is created.
	GetUnsync() []core.Node
	// GetPulse returns actual pulse for current consensus process.
	GetPulse() core.PulseNumber
	// SetHash sets hash of unsync lists for each node of consensus.
	SetHash([]*NodeUnsyncHash)
	// GetHash get hash of unsync lists for each node of  If hash is not calculated yet, then this call blocks
	// until the hash is calculated with SetHash() call.
	GetHash(blockTimeout time.Duration) ([]*NodeUnsyncHash, error)
	// AddUnsyncList add unsync list for remote ref.
	AddUnsyncList(ref core.RecordRef, unsync []core.Node)
	// AddUnsyncHash add unsync hash for remote ref.
	AddUnsyncHash(ref core.RecordRef, hash []*NodeUnsyncHash)
	// GetUnsyncList get unsync list for remote ref.
	GetUnsyncList(ref core.RecordRef) ([]core.Node, bool)
	// GetUnsyncHash get unsync hash for remote ref.
	GetUnsyncHash(ref core.RecordRef) ([]*NodeUnsyncHash, bool)
}

// NodeUnsyncHash data needed for consensus.
type NodeUnsyncHash struct {
	NodeID core.RecordRef
	Hash   []byte
	// TODO: add signature
}

// PartitionPolicy contains all rules how to initiate globule resharding.
type PartitionPolicy interface {
	ShardsCount() int
}

// RoutingTable contains all routing information of the network.
type RoutingTable interface {
	// Start inject dependencies from components
	Start(components core.Components)
	// Resolve NodeID -> ShortID, Address. Can initiate network requests.
	Resolve(core.RecordRef) (*host.Host, error)
	// ResolveS ShortID -> NodeID, Address for node inside current globe.
	ResolveS(core.ShortNodeID) (*host.Host, error)
	// AddToKnownHosts add host to routing table.
	AddToKnownHosts(*host.Host)
	// Rebalance recreate shards of routing table with known hosts according to new partition policy.
	Rebalance(PartitionPolicy)
	// GetLocalNodes get all nodes from the local globe.
	GetLocalNodes() []core.RecordRef
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
type ClaimQueue interface {
	// Pop takes claim from the queue.
	Pop() packets.ReferendumClaim
	// Front returns claim from the queue without removing it from the queue.
	Front() packets.ReferendumClaim
	// Length returns the length of the queue
	Length() int
}
