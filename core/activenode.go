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

package core

// JetRoleMask is bitmask for the set of candidate JetRoles
type JetRoleMask uint64

// NodeState is the state of the node
type NodeState uint8

// TODO: document all node states
const (
	// Joined
	NodeJoined = NodeState(iota + 1)
	// Prepared
	NodePrepared
	// Active
	NodeActive
	// Leaved
	NodeLeaved
	// Suspended
	NodeSuspended
)

func (mask JetRoleMask) IsSet(role JetRole) bool {
	return uint64(mask)&(1<<(uint64(role)-1)) != 0
}

func (mask *JetRoleMask) Set(role JetRole) {
	*mask |= 1 << (uint64(role) - 1)
}

func (mask *JetRoleMask) Unset(role JetRole) {
	var n JetRoleMask = ^(1 << (uint64(role) - 1))
	*mask &= n
}

type ActiveNode struct {
	// NodeID is the unique identifier of the node
	NodeID RecordRef
	// PulseNum is the pulse number after which the new state is assigned to the node
	PulseNum PulseNumber
	// State is the node state
	State NodeState
	// JetRoles is the set of candidate JetRoles for the node
	JetRoles JetRoleMask
	// PublicKey is the public key of the node
	PublicKey []byte
}
