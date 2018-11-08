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

package nodenetwork

import (
	"crypto/ecdsa"
	"encoding/gob"
	"hash/crc32"

	"github.com/insolar/insolar/core"
)

type mutableNode interface {
	core.Node

	SetState(core.NodeState)
	SetPulse(core.PulseNumber)
}

type node struct {
	NodeID        core.RecordRef
	NodeShortID   uint32
	NodeRoles     []core.NodeRole
	NodePublicKey *ecdsa.PublicKey

	NodePulseNum core.PulseNumber
	NodeState    core.NodeState

	NodePhysicalAddress string
	NodeVersion         string
}

func newMutableNode(
	id core.RecordRef,
	roles []core.NodeRole,
	publicKey *ecdsa.PublicKey,
	pulseNum core.PulseNumber,
	state core.NodeState,
	physicalAddress,
	version string) mutableNode {
	return &node{
		NodeID:              id,
		NodeShortID:         generateShortID(id),
		NodeRoles:           roles,
		NodePublicKey:       publicKey,
		NodePulseNum:        pulseNum,
		NodeState:           state,
		NodePhysicalAddress: physicalAddress,
		NodeVersion:         version,
	}
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
	return n.NodeID
}

func (n *node) ShortID() uint32 {
	return n.NodeShortID
}

func (n *node) Pulse() core.PulseNumber {
	return n.NodePulseNum
}

func (n *node) State() core.NodeState {
	return n.NodeState
}

func (n *node) Roles() []core.NodeRole {
	return n.NodeRoles
}

func (n *node) Role() core.NodeRole {
	return n.NodeRoles[0]
}

func (n *node) PublicKey() *ecdsa.PublicKey {
	// TODO: make a copy of pk
	return n.NodePublicKey
}

func (n *node) PhysicalAddress() string {
	return n.NodePhysicalAddress
}

func (n *node) Version() string {
	return n.NodeVersion
}

func (n *node) SetState(state core.NodeState) {
	n.NodeState = state
}

func (n *node) SetPulse(pulseNum core.PulseNumber) {
	n.NodePulseNum = pulseNum
}

func (n *node) SetShortID(id uint32) {
	n.NodeShortID = id
}

type mutableNodes []mutableNode

func (mn mutableNodes) Export() []core.Node {
	nodes := make([]core.Node, len(mn))
	for i := range mn {
		nodes[i] = mn[i]
	}
	return nodes
}

func generateShortID(ref core.RecordRef) uint32 {
	return crc32.ChecksumIEEE(ref[:])
}

func init() {
	gob.Register(&node{})
}
