/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
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
