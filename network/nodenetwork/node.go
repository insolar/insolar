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
	"crypto"
	"encoding/gob"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/utils"
)

type MutableNode interface {
	core.Node

	SetShortID(shortID core.ShortNodeID)
}

type node struct {
	NodeID        core.RecordRef
	NodeShortID   core.ShortNodeID
	NodeRole      core.StaticRole
	NodePublicKey crypto.PublicKey

	NodePulseNum core.PulseNumber

	NodePhysicalAddress string
	NodeVersion         string
}

func newMutableNode(
	id core.RecordRef,
	role core.StaticRole,
	publicKey crypto.PublicKey,
	physicalAddress,
	version string) MutableNode {
	return &node{
		NodeID:              id,
		NodeShortID:         utils.GenerateShortID(id),
		NodeRole:            role,
		NodePublicKey:       publicKey,
		NodePhysicalAddress: physicalAddress,
		NodeVersion:         version,
	}
}

func NewNode(
	id core.RecordRef,
	role core.StaticRole,
	publicKey crypto.PublicKey,
	physicalAddress,
	version string) core.Node {
	return newMutableNode(id, role, publicKey, physicalAddress, version)
}

func (n *node) ID() core.RecordRef {
	return n.NodeID
}

func (n *node) ShortID() core.ShortNodeID {
	return n.NodeShortID
}

func (n *node) Role() core.StaticRole {
	return n.NodeRole
}

func (n *node) PublicKey() crypto.PublicKey {
	return n.NodePublicKey
}

func (n *node) PhysicalAddress() string {
	return n.NodePhysicalAddress
}

func (n *node) GetGlobuleID() core.GlobuleID {
	return 0
}

func (n *node) Version() string {
	return n.NodeVersion
}

func (n *node) SetShortID(id core.ShortNodeID) {
	n.NodeShortID = id
}

func init() {
	gob.Register(&node{})
}
