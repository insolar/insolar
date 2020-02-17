// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
