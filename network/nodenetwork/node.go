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
	"fmt"
	"strconv"
	"strings"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
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

	NodeAddress string
	CAddress    string
	NodeVersion string
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
		NodeShortID:   utils.GenerateShortID(id),
		NodeRole:      role,
		NodePublicKey: publicKey,
		NodeAddress:   address,
		CAddress:      consensusAddress,
		NodeVersion:   version,
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
	return n.NodeShortID
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
	return n.NodeVersion
}

func (n *node) SetShortID(id core.ShortNodeID) {
	n.NodeShortID = id
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
		return address, errors.New("failed to get port from address")
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
