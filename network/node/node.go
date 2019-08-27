//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package node

import (
	"crypto"
	"hash/crc32"
	"sync"
	"sync/atomic"

	"github.com/insolar/insolar/insolar"
)

type MutableNode interface {
	insolar.NetworkNode

	SetShortID(shortID insolar.ShortNodeID)
	SetState(state insolar.NodeState)
	GetSignature() ([]byte, insolar.Signature)
	SetSignature(digest []byte, signature insolar.Signature)
	ChangeState()
	SetLeavingETA(number insolar.PulseNumber)
	SetVersion(version string)
	SetPower(power insolar.Power)
	SetAddress(address string)
}

// GenerateUintShortID generate short ID for node without checking collisions
func GenerateUintShortID(ref insolar.Reference) uint32 {
	return crc32.ChecksumIEEE(ref.Bytes())
}

type node struct {
	NodeID        insolar.Reference
	NodeShortID   uint32
	NodeRole      insolar.StaticRole
	NodePublicKey crypto.PublicKey
	NodePower     uint32

	NodeAddress string

	mutex          sync.RWMutex
	digest         []byte
	signature      insolar.Signature
	NodeVersion    string
	NodeLeavingETA uint32
	state          uint32
}

func (n *node) SetVersion(version string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.NodeVersion = version
}

func (n *node) SetState(state insolar.NodeState) {
	atomic.StoreUint32(&n.state, uint32(state))
}

func (n *node) GetState() insolar.NodeState {
	return insolar.NodeState(atomic.LoadUint32(&n.state))
}

func (n *node) ChangeState() {
	// we don't expect concurrent changes, so do not CAS
	currentState := atomic.LoadUint32(&n.state)
	if currentState >= uint32(insolar.NodeReady) {
		return
	}
	atomic.StoreUint32(&n.state, currentState+1)
}

func newMutableNode(
	id insolar.Reference,
	role insolar.StaticRole,
	publicKey crypto.PublicKey,
	state insolar.NodeState,
	address, version string) MutableNode {

	return &node{
		NodeID:        id,
		NodeShortID:   GenerateUintShortID(id),
		NodeRole:      role,
		NodePublicKey: publicKey,
		NodeAddress:   address,
		NodeVersion:   version,
		state:         uint32(state),
	}
}

func NewNode(
	id insolar.Reference,
	role insolar.StaticRole,
	publicKey crypto.PublicKey,
	address, version string) insolar.NetworkNode {
	return newMutableNode(id, role, publicKey, insolar.NodeReady, address, version)
}

func (n *node) ID() insolar.Reference {
	return n.NodeID
}

func (n *node) ShortID() insolar.ShortNodeID {
	return insolar.ShortNodeID(atomic.LoadUint32(&n.NodeShortID))
}

func (n *node) Role() insolar.StaticRole {
	return n.NodeRole
}

func (n *node) PublicKey() crypto.PublicKey {
	return n.NodePublicKey
}

func (n *node) Address() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.NodeAddress
}

func (n *node) GetGlobuleID() insolar.GlobuleID {
	return 0
}

func (n *node) GetPower() insolar.Power {
	return insolar.Power(atomic.LoadUint32(&n.NodePower))
}

func (n *node) SetPower(power insolar.Power) {
	atomic.StoreUint32(&n.NodePower, uint32(power))
}

func (n *node) Version() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.NodeVersion
}

func (n *node) GetSignature() ([]byte, insolar.Signature) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.digest, n.signature
}

func (n *node) SetSignature(digest []byte, signature insolar.Signature) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.signature = signature
	n.digest = digest
}

func (n *node) SetShortID(id insolar.ShortNodeID) {
	atomic.StoreUint32(&n.NodeShortID, uint32(id))
}

func (n *node) LeavingETA() insolar.PulseNumber {
	return insolar.PulseNumber(atomic.LoadUint32(&n.NodeLeavingETA))
}

func (n *node) SetLeavingETA(number insolar.PulseNumber) {
	n.SetState(insolar.NodeLeaving)
	atomic.StoreUint32(&n.NodeLeavingETA, uint32(number))
}

func (n *node) SetAddress(address string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.NodeAddress = address
}
