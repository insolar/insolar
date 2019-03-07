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
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type MutableNode interface {
	core.Node

	SetShortID(shortID core.ShortNodeID)
	SetState(state core.NodeState)
	ChangeState()
	SetLeavingETA(number core.PulseNumber)
	SetVersion(version string)
}

type node struct {
	NodeID        core.RecordRef
	NodeShortID   uint32
	NodeRole      core.StaticRole
	NodePublicKey crypto.PublicKey

	NodePulseNum core.PulseNumber

	NodeAddress string
	CAddress    string

	versionMutex sync.RWMutex
	NodeVersion  string

	leavingMutex   sync.RWMutex
	NodeLeaving    bool
	NodeLeavingETA core.PulseNumber

	state uint32
}

func (n *node) SetVersion(version string) {
	n.versionMutex.Lock()
	defer n.versionMutex.Unlock()

	n.NodeVersion = version
}

func (n *node) SetState(state core.NodeState) {
	atomic.StoreUint32(&n.state, uint32(state))
}

func (n *node) GetState() core.NodeState {
	return core.NodeState(atomic.LoadUint32(&n.state))
}

func (n *node) ChangeState() {
	// we don't expect concurrent changes, so do not CAS

	currentState := atomic.LoadUint32(&n.state)
	if currentState == uint32(core.NodeReady) {
		return
	}
	atomic.StoreUint32(&n.state, currentState+1)
}

func newMutableNode(
	id core.RecordRef,
	role core.StaticRole,
	publicKey crypto.PublicKey,
	address, version string) MutableNode {

	consensusAddress, err := incrementPort(address)
	if err != nil {
		panic(err)
	}
	return &node{
		NodeID:        id,
		NodeShortID:   utils.GenerateUintShortID(id),
		NodeRole:      role,
		NodePublicKey: publicKey,
		NodeAddress:   address,
		CAddress:      consensusAddress,
		NodeVersion:   version,
		state:         uint32(core.NodeReady),
	}
}

func NewNode(
	id core.RecordRef,
	role core.StaticRole,
	publicKey crypto.PublicKey,
	address, version string) core.Node {
	return newMutableNode(id, role, publicKey, address, version)
}

func (n *node) ID() core.RecordRef {
	return n.NodeID
}

func (n *node) ShortID() core.ShortNodeID {
	return core.ShortNodeID(atomic.LoadUint32(&n.NodeShortID))
}

func (n *node) Role() core.StaticRole {
	return n.NodeRole
}

func (n *node) PublicKey() crypto.PublicKey {
	return n.NodePublicKey
}

func (n *node) Address() string {
	return n.NodeAddress
}

func (n *node) ConsensusAddress() string {
	return n.CAddress
}

func (n *node) GetGlobuleID() core.GlobuleID {
	return 0
}

func (n *node) Version() string {
	n.versionMutex.RLock()
	defer n.versionMutex.RUnlock()

	return n.NodeVersion
}

func (n *node) IsWorking() bool {
	return atomic.LoadUint32(&n.state) == uint32(core.NodeReady)
}

func (n *node) SetShortID(id core.ShortNodeID) {
	atomic.StoreUint32(&n.NodeShortID, uint32(id))
}

func (n *node) Leaving() bool {
	n.leavingMutex.RLock()
	defer n.leavingMutex.RUnlock()
	return n.NodeLeaving
}
func (n *node) LeavingETA() core.PulseNumber {
	n.leavingMutex.RLock()
	defer n.leavingMutex.RUnlock()
	return n.NodeLeavingETA
}

func (n *node) SetLeavingETA(number core.PulseNumber) {
	n.leavingMutex.Lock()
	defer n.leavingMutex.Unlock()

	n.NodeLeaving = true
	n.NodeLeavingETA = number
}

func init() {
	gob.Register(&node{})
}

func ClaimToNode(version string, claim *packets.NodeJoinClaim) (core.Node, error) {
	keyProc := platformpolicy.NewKeyProcessor()
	key, err := keyProc.ImportPublicKeyBinary(claim.NodePK[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ ClaimToNode ] failed to import a public key")
	}
	node := newMutableNode(
		claim.NodeRef,
		claim.NodeRoleRecID,
		key,
		claim.NodeAddress.Get(),
		version)
	node.SetShortID(claim.ShortNodeID)
	return node, nil
}

// incrementPort increments port number if it not equals 0
func incrementPort(address string) (string, error) {
	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		return address, errors.New("failed to get port from address " + address)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return address, err
	}

	if port != 0 {
		port++
	}
	return fmt.Sprintf("%s:%d", parts[0], port), nil
}
