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
	NID        core.RecordRef
	NShortID   uint32
	NRoles     []core.NodeRole
	NPublicKey *ecdsa.PublicKey

	NPulseNum core.PulseNumber
	NState    core.NodeState

	NPhysicalAddress string
	NVersion         string
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
		NID:              id,
		NShortID:         generateShortID(id),
		NRoles:           roles,
		NPublicKey:       publicKey,
		NPulseNum:        pulseNum,
		NState:           state,
		NPhysicalAddress: physicalAddress,
		NVersion:         version,
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
	return n.NID
}

func (n *node) ShortID() uint32 {
	return n.NShortID
}

func (n *node) Pulse() core.PulseNumber {
	return n.NPulseNum
}

func (n *node) State() core.NodeState {
	return n.NState
}

func (n *node) Roles() []core.NodeRole {
	return n.NRoles
}

func (n *node) Role() core.NodeRole {
	return n.NRoles[0]
}

func (n *node) PublicKey() *ecdsa.PublicKey {
	// TODO: make a copy of pk
	return n.NPublicKey
}

func (n *node) PhysicalAddress() string {
	return n.NPhysicalAddress
}

func (n *node) Version() string {
	return n.NVersion
}

func (n *node) SetState(state core.NodeState) {
	n.NState = state
}

func (n *node) SetPulse(pulseNum core.PulseNumber) {
	n.NPulseNum = pulseNum
}

func (n *node) SetShortID(id uint32) {
	n.NShortID = id
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
