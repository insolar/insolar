//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package insolar

import (
	"context"
	"time"
)

// Cascade contains routing data for cascade sending
type Cascade struct {
	// NodeIds contains the slice of node identifiers that will receive the message
	NodeIds []Reference
	// GeneratedEntropy is used for pseudorandom cascade building
	Entropy Entropy
	// Replication factor is the number of children nodes of the each node of the cascade
	ReplicationFactor uint
}

// RemoteProcedure is remote procedure call function.
type RemoteProcedure func(ctx context.Context, args []byte) ([]byte, error)

// HealthChecker interface provides method to check network health
type HealthChecker interface {
	// IsAlive returns true if todo: fix requirements
	IsAlive() bool
}

type StatusReply struct {
	NetworkState    NetworkState
	Origin          NetworkNode
	ActiveListSize  int
	WorkingListSize int
	// Nodes from active list
	Nodes []NetworkNode
	// Pulse from network pulse storage
	Pulse     Pulse
	Version   string
	Timestamp time.Time
	// node start timestamp for uptime duration
	StartTime time.Time
}

type NetworkStatus interface {
	GetNetworkStatus() StatusReply
}

//go:generate minimock -i github.com/insolar/insolar/insolar.Network -o ../testutils -s _mock.go -g

// Network is interface for network modules facade.
type Network interface {
	// SendParcel sends a message.
	SendMessage(nodeID Reference, method string, msg Parcel) ([]byte, error)
	// SendCascadeMessage sends a message.
	SendCascadeMessage(data Cascade, method string, msg Parcel) error
	// RemoteProcedureRegister is remote procedure register func.
	RemoteProcedureRegister(name string, method RemoteProcedure)
	// Leave notify other nodes that this node want to leave and doesn't want to receive new tasks
	Leave(ctx context.Context, ETA PulseNumber)
	// GetState returns our current thoughs about whole network
	GetState() NetworkState
	GetCert(context.Context, *Reference) (Certificate, error)
}

//go:generate minimock -i github.com/insolar/insolar/insolar.PulseDistributor -o ../testutils -s _mock.go -g

// PulseDistributor is interface for pulse distribution.
type PulseDistributor interface {
	// Distribute distributes a pulse across the network.
	Distribute(context.Context, Pulse)
}

// NetworkState type for bootstrapping process
type NetworkState int

//go:generate stringer -type=NetworkState
const (
	// NoNetworkState state means that nodes doesn`t match majority_rule
	NoNetworkState NetworkState = iota
	JoinerBootstrap
	WaitConsensus
	WaitMajority
	WaitMinRoles
	WaitPulsar
	CompleteNetworkState
)
