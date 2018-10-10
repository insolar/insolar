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

package packet

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/hostnetwork/id"
)

// CommandType - type for commands.
type CommandType int

const (
	// Unknown - unknown command.
	Unknown = CommandType(iota + 1)
	// StartRelay - command start relay.
	StartRelay
	// StopRelay - command stop relay.
	StopRelay
	// BeginAuthentication - begin authentication.
	BeginAuthentication
	// RevokeAuthentication - revoke authentication.
	RevokeAuthentication
)

// RequestCheckNodePriv is data for check node privileges.
type RequestCheckNodePriv struct {
	RoleKey string
}

// RequestDataFindHost is data for FindHost request.
type RequestDataFindHost struct {
	Target []byte
}

// RequestDataFindValue is data for FindValue request.
type RequestDataFindValue struct {
	Target []byte
}

// RequestDataStore is data for Store request.
type RequestDataStore struct {
	Data       []byte
	Publishing bool // Whether or not we are the original publisher.
}

// RequestDataRPC is data for RPC request.
type RequestDataRPC struct {
	Method string
	Args   [][]byte
}

// RequestCascadeSend is data for cascade sending feature
type RequestCascadeSend struct {
	RPC  RequestDataRPC
	Data core.Cascade
}

// RequestPulse is data received from a pulsar
type RequestPulse struct {
	Pulse core.Pulse
}

// RequestGetRandomHosts is data for the call that returns random hosts of the DHT network
type RequestGetRandomHosts struct {
	HostsNumber int
}

// RequestRelay is data for relay request (commands: start/stop relay).
type RequestRelay struct {
	Command CommandType
}

// RequestAuthentication is data for authentication.
type RequestAuthentication struct {
	Command CommandType
}

// RequestCheckOrigin is data to check originality.
type RequestCheckOrigin struct {
}

// RequestObtainIP is data to obtain a self IP.
type RequestObtainIP struct {
}

// RequestRelayOwnership is data to notify that current host can be a relay.
type RequestRelayOwnership struct {
	Ready bool
}

// RequestKnownOuterHosts is data to notify home subnet about known outer hosts.
type RequestKnownOuterHosts struct {
	ID         string // origin ID
	OuterHosts int    // number of known outer hosts
}

// RequestCheckPublicKey is data to check a public key.
type RequestCheckPublicKey struct {
	NodeID core.RecordRef
	HostID id.ID
}

// RequestCheckSignedNonce is data to check a signed nonce.
type RequestCheckSignedNonce struct {
	Signed []byte
}

// RequestActiveNodes is request to get active nodes.
type RequestActiveNodes struct {
}

// RequestExchangeUnsyncLists is request to exchange unsync lists during consensus
type RequestExchangeUnsyncLists struct {
	Pulse      core.PulseNumber
	UnsyncList []*core.ActiveNode
}

// RequestExchangeUnsyncHash is request to exchange hash of merged unsync lists during consensus
type RequestExchangeUnsyncHash struct {
	Pulse      core.PulseNumber
	UnsyncHash []*consensus.NodeUnsyncHash
}

// RequestDisconnect is request to disconnect from active list.
type RequestDisconnect struct {
}
