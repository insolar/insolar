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

package network

import (
	"context"
	"time"

	"github.com/insolar/insolar/component"
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/node"
)

type BootstrapResult struct {
	Host *host.Host
	// FirstPulseTime    time.Time
	ReconnectRequired bool
	NetworkSize       int
}

// Controller contains network logic.
type Controller interface {
	component.Initer
	// SendParcel send message to nodeID.
	SendMessage(nodeID insolar.Reference, name string, msg insolar.Parcel) ([]byte, error)
	// RemoteProcedureRegister register remote procedure that will be executed when message is received.
	RemoteProcedureRegister(name string, method insolar.RemoteProcedure)
	// SendCascadeMessage sends a message from MessageBus to a cascade of nodes.
	SendCascadeMessage(data insolar.Cascade, method string, msg insolar.Parcel) error
	// Bootstrap init complex bootstrap process. Blocks until bootstrap is complete.
	Bootstrap(ctx context.Context) (*BootstrapResult, error)
	// SetLastIgnoredPulse set pulse number after which we will begin setting new pulses to PulseManager
	SetLastIgnoredPulse(number insolar.PulseNumber)
	// GetLastIgnoredPulse get last pulse that will be ignored
	GetLastIgnoredPulse() insolar.PulseNumber
	AuthenticateToDiscoveryNode(ctx context.Context, discovery insolar.DiscoveryNode) error
}

// RequestHandler handler function to process incoming requests from network.
type RequestHandler func(context.Context, Request) (Response, error)

//go:generate minimock -i github.com/insolar/insolar/network.HostNetwork -o ../testutils/network -s _mock.go

// HostNetwork simple interface to send network requests and process network responses.
type HostNetwork interface {
	component.Starter
	component.Stopper
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string
	// GetNodeID get current node ID.
	GetNodeID() insolar.Reference

	// SendRequest send request to a remote node.
	SendRequest(ctx context.Context, request Request, receiver insolar.Reference) (Future, error)
	// SendRequestToHost send request packet to a remote node.
	SendRequestToHost(ctx context.Context, request Request, receiver *host.Host) (Future, error)
	// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
	RegisterRequestHandler(t types.PacketType, handler RequestHandler)
	// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
	NewRequestBuilder() RequestBuilder
	// BuildResponse create response to an incoming request with Data set to responseData.
	BuildResponse(ctx context.Context, request Request, responseData interface{}) Response
}

type ConsensusPacketHandler func(incomingPacket consensus.ConsensusPacket, sender insolar.Reference)

//go:generate minimock -i github.com/insolar/insolar/network.ConsensusNetwork -o ../testutils/network -s _mock.go
type ConsensusNetwork interface {
	component.Starter
	component.Stopper
	// PublicAddress returns public address that can be published for all nodes.
	PublicAddress() string
	// GetNodeID get current node ID.
	GetNodeID() insolar.Reference

	// SignAndSendPacket send request to a remote node.
	SignAndSendPacket(packet consensus.ConsensusPacket, receiver insolar.Reference, service insolar.CryptographyService) error
	// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
	RegisterPacketHandler(t consensus.PacketType, handler ConsensusPacketHandler)
}

// RequestID is 64 bit unsigned int request id.
type RequestID uint64

// Packet is a packet that is transported via network by HostNetwork.
type Packet interface {
	GetSender() insolar.Reference
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

//go:generate minimock -i github.com/insolar/insolar/network.PulseHandler -o ../testutils/network -s _mock.go

// PulseHandler interface to process new pulse.
type PulseHandler interface {
	HandlePulse(ctx context.Context, pulse insolar.Pulse)
}

//go:generate minimock -i github.com/insolar/insolar/network.NodeKeeper -o ../testutils/network -s _mock.go

// NodeKeeper manages unsync, sync and active lists.
type NodeKeeper interface {
	insolar.NodeNetwork

	// TODO: remove this interface when bootstrap mechanism completed
	insolar.SwitcherWorkAround

	// GetCloudHash returns current cloud hash
	GetCloudHash() []byte
	// SetCloudHash set new cloud hash
	SetCloudHash([]byte)
	// SetInitialSnapshot set initial snapshot for nodekeeper
	SetInitialSnapshot(nodes []insolar.NetworkNode)
	// GetAccessor get accessor to the internal snapshot for the current pulse
	// TODO: add pulse to the function signature to get data of various pulses
	GetAccessor() Accessor
	// GetOriginJoinClaim get origin NodeJoinClaim
	GetOriginJoinClaim() (*consensus.NodeJoinClaim, error)
	// GetOriginAnnounceClaim get origin NodeAnnounceClaim
	GetOriginAnnounceClaim(mapper consensus.BitSetMapper) (*consensus.NodeAnnounceClaim, error)
	// GetClaimQueue get the internal queue of claims
	GetClaimQueue() ClaimQueue
	// GetSnapshotCopy get copy of the current nodekeeper snapshot
	GetSnapshotCopy() *node.Snapshot
	// Sync move unsync -> sync
	Sync(context.Context, []insolar.NetworkNode, []consensus.ReferendumClaim) error
	// MoveSyncToActive merge sync list with active nodes
	MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) error
	// GetConsensusInfo get additional info for the current consensus process
	GetConsensusInfo() ConsensusInfo
}

// TODO: refactor code and make it not necessary
// ConsensusInfo additional info for the current consensus process
type ConsensusInfo interface {
	// NodesJoinedDuringPreviousPulse returns true if the last Sync call contained approved Join claims
	NodesJoinedDuringPreviousPulse() bool
	// AddTemporaryMapping add temporary mapping till the next pulse for consensus
	AddTemporaryMapping(nodeID insolar.Reference, shortID insolar.ShortNodeID, address string) error
	// ResolveConsensus get temporary mapping by short ID
	ResolveConsensus(shortID insolar.ShortNodeID) *host.Host
	// ResolveConsensusRef get temporary mapping by node ID
	ResolveConsensusRef(nodeID insolar.Reference) *host.Host
	// SetIsJoiner instruct current node whether it should perform consensus as joiner or not
	SetIsJoiner(isJoiner bool)
	// IsJoiner true if current node should perform consensus as joiner
	IsJoiner() bool
}

// PartitionPolicy contains all rules how to initiate globule resharding.
type PartitionPolicy interface {
	ShardsCount() int
}

// RoutingTable contains all routing information of the network.
type RoutingTable interface {
	// Resolve NodeID -> ShortID, Address. Can initiate network requests.
	Resolve(insolar.Reference) (*host.Host, error)
	// ResolveConsensus ShortID -> NodeID, Address for node inside current globe for current consensus.
	ResolveConsensus(insolar.ShortNodeID) (*host.Host, error)
	// ResolveConsensusRef NodeID -> ShortID, Address for node inside current globe for current consensus.
	ResolveConsensusRef(insolar.Reference) (*host.Host, error)
	// AddToKnownHosts add host to routing table.
	AddToKnownHosts(*host.Host)
	// Rebalance recreate shards of routing table with known hosts according to new partition policy.
	Rebalance(PartitionPolicy)
}

// ClaimQueue is the queue that contains consensus claims.
type ClaimQueue interface {
	// Pop takes claim from the queue.
	Pop() consensus.ReferendumClaim
	// Front returns claim from the queue without removing it from the queue.
	Front() consensus.ReferendumClaim
	// Length returns the length of the queue
	Length() int
	// Push adds claim to the queue.
	Push(claim consensus.ReferendumClaim)
	// Clear removes all claims from queue
	Clear()
}

// Accessor is interface that provides read access to nodekeeper internal snapshot
type Accessor interface {
	// GetWorkingNode get working node by its reference. Returns nil if node is not found or is not working.
	GetWorkingNode(ref insolar.Reference) insolar.NetworkNode
	// GetWorkingNodes returns sorted list of all working nodes.
	GetWorkingNodes() []insolar.NetworkNode
	// GetWorkingNodesByRole get working nodes by role.
	GetWorkingNodesByRole(role insolar.DynamicRole) []insolar.Reference

	// GetActiveNode returns active node.
	GetActiveNode(ref insolar.Reference) insolar.NetworkNode
	// GetActiveNodes returns unsorted list of all active nodes.
	GetActiveNodes() []insolar.NetworkNode
	// GetActiveNodeByShortID get active node by short ID. Returns nil if node is not found.
	GetActiveNodeByShortID(shortID insolar.ShortNodeID) insolar.NetworkNode
}

// Mutator is interface that provides read and write access to a snapshot
type Mutator interface {
	Accessor
	// AddWorkingNode adds active node to index and underlying snapshot so it is accessible via GetActiveNode(s).
	AddWorkingNode(n insolar.NetworkNode)
}

// Gatewayer is a network which can change it's Gateway
//go:generate minimock -i github.com/insolar/insolar/network.Gatewayer -o ../testutils/network -s _mock.go
type Gatewayer interface {
	Gateway() Gateway
	SetGateway(Gateway)
}

// Gateway responds for whole network state
type Gateway interface {
	Run(context.Context)
	GetState() insolar.NetworkState
	OnPulse(context.Context, insolar.Pulse) error
	NewGateway(insolar.NetworkState) Gateway
}
