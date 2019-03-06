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

package network

import (
	"context"
	"time"

	"github.com/insolar/insolar/component"
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
)

type BootstrapResult struct {
	Host *host.Host
	// FirstPulseTime    time.Time
	ReconnectRequired bool
}

// Controller contains network logic.
type Controller interface {
	component.Initer
	// SendParcel send message to nodeID.
	SendMessage(nodeID core.RecordRef, name string, msg core.Parcel) ([]byte, error)
	// RemoteProcedureRegister register remote procedure that will be executed when message is received.
	RemoteProcedureRegister(name string, method core.RemoteProcedure)
	// SendCascadeMessage sends a message from MessageBus to a cascade of nodes.
	SendCascadeMessage(data core.Cascade, method string, msg core.Parcel) error
	// Bootstrap init complex bootstrap process. Blocks until bootstrap is complete.
	Bootstrap(ctx context.Context) (*BootstrapResult, error)

	// TODO: workaround methods, should be deleted once network consensus is alive

	// SetLastIgnoredPulse set pulse number after which we will begin setting new pulses to PulseManager
	SetLastIgnoredPulse(number core.PulseNumber)
	// GetLastIgnoredPulse get last pulse that will be ignored
	GetLastIgnoredPulse() core.PulseNumber
}

// RequestHandler handler function to process incoming requests from network.
type RequestHandler func(context.Context, Request) (Response, error)

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
	SendRequest(ctx context.Context, request Request, receiver core.RecordRef) (Future, error)
	// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
	RegisterRequestHandler(t types.PacketType, handler RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() RequestBuilder
	// BuildResponse create response to an incoming request with Data set to responseData.
	BuildResponse(ctx context.Context, request Request, responseData interface{}) Response
}

type ConsensusPacketHandler func(incomingPacket consensus.ConsensusPacket, sender core.RecordRef)

//go:generate minimock -i github.com/insolar/insolar/network.ConsensusNetwork -o ../testutils/network -s _mock.go
type ConsensusNetwork interface {
	component.Starter
	component.Stopper
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string
	// GetNodeID get current node ID.
	GetNodeID() core.RecordRef

	// SignAndSendPacket send request to a remote node.
	SignAndSendPacket(packet consensus.ConsensusPacket, receiver core.RecordRef, service core.CryptographyService) error
	// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
	RegisterPacketHandler(t consensus.PacketType, handler ConsensusPacketHandler)
}

// RequestID is 64 bit unsigned int request id.
type RequestID uint64

// Packet is a packet that is transported via network by HostNetwork.
type Packet interface {
	GetSender() core.RecordRef
	GetSenderHost() *host.Host
	GetType() types.PacketType
	GetData() interface{}
	GetRequestID() RequestID
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

// NodeKeeper manages unsync, sync and active lists.
//go:generate minimock -i github.com/insolar/insolar/network.NodeKeeper -o ../testutils/network -s _mock.go
type NodeKeeper interface {
	core.NodeNetwork

	// TODO: remove this interface when bootstrap mechanism completed
	core.SwitcherWorkAround

	// GetCloudHash returns current cloud hash
	GetCloudHash() []byte
	// SetCloudHash set new cloud hash
	SetCloudHash([]byte)
	// GetActiveNode returns active node.
	GetActiveNode(ref core.RecordRef) core.Node
	// AddActiveNodes add active nodes.
	AddActiveNodes([]core.Node)
	// GetActiveNodes returns active nodes.
	GetActiveNodes() []core.Node
	// GetActiveNodeByShortID get active node by short ID. Returns nil if node is not found.
	GetActiveNodeByShortID(shortID core.ShortNodeID) core.Node
	// SetState set state of the NodeKeeper
	SetState(core.NodeNetworkState)
	// GetOriginJoinClaim get origin NodeJoinClaim
	GetOriginJoinClaim() (*consensus.NodeJoinClaim, error)
	// GetOriginAnnounceClaim get origin NodeAnnounceClaim
	GetOriginAnnounceClaim(mapper consensus.BitSetMapper) (*consensus.NodeAnnounceClaim, error)
	// NodesJoinedDuringPreviousPulse returns true if the last Sync call contained approved Join claims
	NodesJoinedDuringPreviousPulse() bool
	// AddPendingClaim add pending claim to the internal queue of claims
	AddPendingClaim(consensus.ReferendumClaim) bool
	// GetClaimQueue get the internal queue of claims
	GetClaimQueue() ClaimQueue
	// GetUnsyncList get unsync list for current pulse. Has copy of active node list from nodekeeper as internal state.
	// Should be called when nodekeeper state is ReadyNodeNetworkState.
	GetUnsyncList() UnsyncList
	// GetSparseUnsyncList get sparse unsync list for current pulse with predefined length of active node list.
	// Does not contain active list, should collect active list during its lifetime via AddClaims.
	// Should be called when nodekeeper state is WaitingNodeNetworkState.
	GetSparseUnsyncList(length int) UnsyncList
	// Sync move unsync -> sync
	Sync(list UnsyncList)
	// MoveSyncToActive merge sync list with active nodes
	MoveSyncToActive(ctx context.Context) error
	// AddTemporaryMapping add temporary mapping till the next pulse for consensus
	AddTemporaryMapping(nodeID core.RecordRef, shortID core.ShortNodeID, address string) error
	// ResolveConsensus get temporary mapping by short ID
	ResolveConsensus(shortID core.ShortNodeID) *host.Host
	// ResolveConsensusRef get temporary mapping by node ID
	ResolveConsensusRef(nodeID core.RecordRef) *host.Host
}

// UnsyncList is interface to manage unsync list
//go:generate minimock -i github.com/insolar/insolar/network.UnsyncList -o ../testutils/network -s _mock.go
type UnsyncList interface {
	consensus.BitSetMapper
	// ApproveSync
	ApproveSync([]core.RecordRef)
	// AddClaims
	AddClaims(map[core.RecordRef][]consensus.ReferendumClaim) error
	// AddNode
	AddNode(node core.Node, bitsetIndex uint16)
	// GetClaims
	GetClaims(nodeID core.RecordRef) []consensus.ReferendumClaim
	// InsertClaims
	InsertClaims(core.RecordRef, []consensus.ReferendumClaim)
	// AddProof
	AddProof(nodeID core.RecordRef, proof *consensus.NodePulseProof)
	// GetProof
	GetProof(nodeID core.RecordRef) *consensus.NodePulseProof
	// GetGlobuleHashSignature
	GetGlobuleHashSignature(ref core.RecordRef) (consensus.GlobuleHashSignature, bool)
	// SetGlobuleHashSignature
	SetGlobuleHashSignature(core.RecordRef, consensus.GlobuleHashSignature)
	// GetActiveNode get active node by reference ID for current consensus
	GetActiveNode(ref core.RecordRef) core.Node
	// GetActiveNodes get active nodes for current consensus
	GetActiveNodes() []core.Node
	// GetMergedCopy returns copy of unsyncList with claims applied
	GetMergedCopy() (*MergedListCopy, error)
	//
	RemoveNode(nodeID core.RecordRef)
}

type MergedListCopy struct {
	ActiveList                 map[core.RecordRef]core.Node
	NodesJoinedDuringPrevPulse bool
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
	// ResolveConsensus ShortID -> NodeID, Address for node inside current globe for current consensus.
	ResolveConsensus(core.ShortNodeID) (*host.Host, error)
	// ResolveConsensusRef NodeID -> ShortID, Address for node inside current globe for current consensus.
	ResolveConsensusRef(core.RecordRef) (*host.Host, error)
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
	SendRequestPacket(ctx context.Context, request Request, receiver *host.Host) (Future, error)
	// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
	RegisterPacketHandler(t types.PacketType, handler RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() RequestBuilder
	// BuildResponse create response to an incoming request with Data set to responseData.
	BuildResponse(ctx context.Context, request Request, responseData interface{}) Response
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
