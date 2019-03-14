/*
 *    Copyright 2019 Insolar Technologies
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

package core

import (
	"crypto"
)

// ShortNodeID is the shortened ID of node that is unique inside the globe
type ShortNodeID uint32

// GlobuleID is the ID of the globe
type GlobuleID uint32

// NodeState is the state of the node
type NodeState uint8

//go:generate stringer -type=NodeState
const (
	NodeUndefined NodeState = iota
	NodePending
	NodeReady
	NodeLeaving
)

//go:generate minimock -i github.com/insolar/insolar/core.Node -o ../testutils/network -s _mock.go
type Node interface {
	// ID is the unique identifier of the node
	ID() RecordRef
	// ShortID get short ID of node
	ShortID() ShortNodeID
	// Role is the candidate Role for the node
	Role() StaticRole
	// PublicKey is the public key of the node
	PublicKey() crypto.PublicKey
	// Address is the network address of the node
	Address() string
	// ConsensusAddress is the network address of the node for consensus packets exchange
	ConsensusAddress() string
	// GetGlobuleID returns node current globule id
	GetGlobuleID() GlobuleID
	// Version of node software
	Version() string
	// LeavingETA is pulse number, after which node leave
	LeavingETA() PulseNumber
	// GetState get state of the node
	GetState() NodeState
}

//go:generate minimock -i github.com/insolar/insolar/core.NodeNetwork -o ../testutils/network -s _mock.go
type NodeNetwork interface {
	// GetOrigin get origin node for the current insolard. Returns nil if the current insolard is not a working node.
	GetOrigin() Node
	// GetWorkingNode get working node by its reference. Returns nil if node is not found.
	GetWorkingNode(ref RecordRef) Node
	// GetWorkingNodes get working nodes.
	GetWorkingNodes() []Node
	// GetWorkingNodesByRole get working nodes by role
	GetWorkingNodesByRole(role DynamicRole) []RecordRef
}

// TODO: remove this interface when bootstrap mechanism completed
// SwitcherWorkAround temp interface for NetworkSwitcher interaction
//go:generate minimock -i github.com/insolar/insolar/core.SwitcherWorkAround -o ../testutils/network -s _mock.go
type SwitcherWorkAround interface {
	// IsBootstrapped method shows that all DiscoveryNodes finds each other
	IsBootstrapped() bool
	// SetIsBootstrapped method set is bootstrap completed
	SetIsBootstrapped(isBootstrap bool)
}
