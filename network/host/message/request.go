/*
 *    Copyright 2018 INS Ecosystem
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

package message

// CommandType - type for commands.
type CommandType int

const (
	// Unknown - unknown command.
	Unknown = CommandType(iota + 1)
	// StartRelay - command start relay.
	StartRelay
	// StopRelay - command stop relay.
	StopRelay
	// BeginAuth - begin authentication.
	BeginAuth
	// RevokeAuth - revoke authentication.
	RevokeAuth
)

// RequestDataFindNode is data for FindNode request.
type RequestDataFindNode struct {
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

// RequestRelay is data for relay request (commands: start/stop relay).
type RequestRelay struct {
	Command CommandType
}

// RequestAuth is data for authentication.
type RequestAuth struct {
	Command CommandType
}

// RequestCheckOrigin is data to check originality.
type RequestCheckOrigin struct {
}

// RequestObtainIP is data to obtain a self IP.
type RequestObtainIP struct {
}

// RequestRelayOwnership is data to notify that current node can be a relay.
type RequestRelayOwnership struct {
	Ready bool
}

// RequestKnownOuterNodes is data to notify home subnet about known outer nodes.
type RequestKnownOuterNodes struct {
	ID         string // origin ID
	OuterNodes int    // number of known outer nodes
}
