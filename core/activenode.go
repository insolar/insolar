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

// BMJetRole is bitmask for the set of candidate JetRoles
type BMJetRole uint64

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

const (
	// BMRoleVirtualExecutor is responsible for current pulse CPU operations.
	BMRoleVirtualExecutor = 1 << uint(RoleVirtualExecutor-1)
	// BMRoleVirtualValidator is responsible for previous pulse CPU operations.
	BMRoleVirtualValidator = 1 << uint(RoleVirtualValidator-1)
	// BMRoleLightExecutor is responsible for current pulse Disk operations.
	BMRoleLightExecutor = 1 << uint(RoleLightExecutor-1)
	// BMRoleLightValidator is responsible for previous pulse Disk operations.
	BMRoleLightValidator = 1 << uint(RoleLightValidator-1)
	// BMRoleHeavyExecutor is responsible for permanent Disk operations.
	BMRoleHeavyExecutor = 1 << uint(RoleHeavyExecutor-1)
)

type ActiveNode struct {
	// NodeID is the unique identifier of the node
	NodeID RecordRef
	// PulseNum is the pulse number after which the new state is assigned to the node
	PulseNum PulseNumber
	// State is the node state
	State NodeState
	// JetRoles is the set of candidate JetRoles for the node
	JetRoles BMJetRole
	// PublicKey is the public key of the node
	PublicKey []byte
}
