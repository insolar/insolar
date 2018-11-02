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

package nodekeeper

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/core"
)

type mutableNode interface {
	core.Node

	SetState(core.NodeState)
	SetPulse(core.PulseNumber)
}

type node struct {
	id        core.RecordRef
	roles     []core.NodeRole
	publicKey *ecdsa.PublicKey

	pulseNum core.PulseNumber
	state    core.NodeState

	physicalAddress string
	version         string
}

func NewNode(
	id core.RecordRef,
	roles []core.NodeRole,
	publicKey *ecdsa.PublicKey,
	pulseNum core.PulseNumber,
	state core.NodeState,
	physicalAddress,
	version string) core.Node {
	return newMutableNode(id, roles, publicKey, pulseNum, state, physicalAddress, version)
}

func (n *node) ID() core.RecordRef {
	return n.id
}

func (n *node) Pulse() core.PulseNumber {
	return n.pulseNum
}

func (n *node) State() core.NodeState {
	return n.state
}

func (n *node) Roles() []core.NodeRole {
	return n.roles
}

func (n *node) PublicKey() *ecdsa.PublicKey {
	// TODO: make a copy of pk
	return n.publicKey
}

func (n *node) PhysicalAddress() string {
	return n.physicalAddress
}

func (n *node) Version() string {
	return n.version
}
