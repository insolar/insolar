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
	"crypto"
)

const (
	ShortNodeIDSize = 4
)

// ShortNodeID is the shortened ID of node that is unique inside the globe
type ShortNodeID uint32 // ZERO is RESERVED

const AbsentShortNodeID ShortNodeID = 0

func (v ShortNodeID) IsAbsent() bool { return v == AbsentShortNodeID }

func (v ShortNodeID) Equal(other ShortNodeID) bool { return v == other }

// GlobuleID is the ID of the globe
type GlobuleID uint32

// NodeState is the state of the node
type NodeState uint8

// Power is node power
type Power uint8

//go:generate stringer -type=NodeState
const (
	// NodeUndefined node started but is not connected to network yet
	NodeUndefined NodeState = iota
	// NodeJoining node is in first pulse of discovery bootstrap or is joining to a bootstrapped network
	NodeJoining
	// NodeReady node is connected to network
	NodeReady
	// NodeLeaving node is about to leave network
	NodeLeaving
)

//go:generate minimock -i github.com/insolar/insolar/insolar.NetworkNode -o ../testutils/network -s _mock.go -g

type NetworkNode interface {
	// ID is the unique identifier of the node
	ID() Reference
	// ShortID get short ID of node
	ShortID() ShortNodeID
	// Role is the candidate Role for the node
	Role() StaticRole
	// PublicKey is the public key of the node
	PublicKey() crypto.PublicKey
	// Address is the network address of the node
	Address() string
	// GetGlobuleID returns node current globule id
	GetGlobuleID() GlobuleID
	// Version of node software
	Version() string
	// LeavingETA is pulse number, after which node leave
	LeavingETA() PulseNumber
	// GetState get state of the node
	GetState() NodeState
	// GetPower get power of node
	GetPower() Power
}
