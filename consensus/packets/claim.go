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

package packets

import (
	"crypto"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type ClaimType uint8

const (
	TypeNodeJoinClaim = ClaimType(iota + 1)
	TypeNodeAnnounceClaim
	TypeCapabilityPollingAndActivation
	TypeNodeViolationBlame
	TypeNodeBroadcast
	TypeNodeLeaveClaim
	TypeChangeNetworkClaim
)

const claimHeaderSize = 2

// ChangeNetworkClaim uses to change network state.
type ChangeNetworkClaim struct {
}

func (cnc *ChangeNetworkClaim) Type() ClaimType {
	return TypeChangeNetworkClaim
}

type ReferendumClaim interface {
	Serializer
	Type() ClaimType
}

type ClaimSupplementary interface {
	AddSupplementaryInfo(nodeID core.RecordRef)
}

type SignedClaim interface {
	GetNodeID() core.RecordRef
	GetPublicKey() (crypto.PublicKey, error)
	SerializeRaw() ([]byte, error)
	GetSignature() []byte
}

// NodeBroadcast is a broadcast of info. Must be brief and only one entry per node.
// Type 4.
type NodeBroadcast struct {
	EmergencyLevel uint8
}

func (nb *NodeBroadcast) Type() ClaimType {
	return TypeNodeBroadcast
}

// CapabilityPoolingAndActivation is a type 3.
type CapabilityPoolingAndActivation struct {
	PollingFlags   uint16
	CapabilityType uint16
	CapabilityRef  [ReferenceLength]byte
}

func (cpa *CapabilityPoolingAndActivation) Type() ClaimType {
	return TypeCapabilityPollingAndActivation
}

// NodeViolationBlame is a type 2.
type NodeViolationBlame struct {
	BlameNodeID   uint32
	TypeViolation uint8
}

func (nvb *NodeViolationBlame) Type() ClaimType {
	return TypeNodeViolationBlame
}

const NodeAddressSize = 20

// TODO: create heterogeneous structure for variuos types of adresses (IPv4, IPv6, etc.)
type NodeAddress [NodeAddressSize]byte

func NewNodeAddress(address string) NodeAddress {
	var result NodeAddress
	result.Set(address)
	return result
}

func (address *NodeAddress) Set(s string) {
	copy(address[:], []byte(s)[:NodeAddressSize])
}

func (address NodeAddress) Get() string {
	return string(address[:])
}

// NodeJoinClaim is a type 1, len == 272.
type NodeJoinClaim struct {
	ShortNodeID             core.ShortNodeID
	RelayNodeID             core.ShortNodeID
	ProtocolVersionAndFlags uint32
	JoinsAfter              uint32
	NodeRoleRecID           core.StaticRole
	NodeRef                 core.RecordRef
	NodeAddress             NodeAddress
	NodePK                  [PublicKeyLength]byte
	Signature               [SignatureLength]byte
}

func (njc *NodeJoinClaim) GetNodeID() core.RecordRef {
	return njc.NodeRef
}

func (njc *NodeJoinClaim) GetPublicKey() (crypto.PublicKey, error) {
	keyProc := platformpolicy.NewKeyProcessor()
	return keyProc.ImportPublicKey(njc.NodePK[:])
}

func (njc *NodeJoinClaim) GetSignature() []byte {
	return njc.Signature[:]
}

func (njc *NodeJoinClaim) Type() ClaimType {
	return TypeNodeJoinClaim
}

// NodeJoinClaim is a type 5, len == 272.
type NodeAnnounceClaim struct {
	NodeJoinClaim

	NodeAnnouncerIndex uint16
	NodeJoinerIndex    uint16
	NodeCount          uint16

	// mapper is used to fill three fields above, is not serialized
	BitSetMapper BitSetMapper
}

func (nac *NodeAnnounceClaim) Type() ClaimType {
	return TypeNodeAnnounceClaim
}

// NodeLeaveClaim can be the only be issued by the node itself and must be the only claim record.
// Should be executed with the next pulse. Type 1, len == 0.
type NodeLeaveClaim struct {
	// additional field that is not serialized and is set from transport layer on packet receive
	NodeID core.RecordRef
}

func (nlc *NodeLeaveClaim) AddSupplementaryInfo(nodeID core.RecordRef) {
	nlc.NodeID = nodeID
}

func (nlc *NodeLeaveClaim) Type() ClaimType {
	return TypeNodeLeaveClaim
}

func getClaimSize(claim ReferendumClaim) uint16 {
	return claimSizeMap[claim.Type()]
}

func getClaimWithHeaderSize(claim ReferendumClaim) uint16 {
	return getClaimSize(claim) + claimHeaderSize
}

func NodeToClaim(node core.Node) (*NodeJoinClaim, error) {
	keyProc := platformpolicy.NewKeyProcessor()
	exportedKey, err := keyProc.ExportPublicKey(node.PublicKey())
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeToClaim ] failed to export a public key")
	}
	var keyData [PublicKeyLength]byte
	copy(keyData[:], exportedKey[:PublicKeyLength])

	var s [SignatureLength]byte
	return &NodeJoinClaim{
		ShortNodeID:             node.ShortID(),
		RelayNodeID:             node.ShortID(),
		ProtocolVersionAndFlags: 0,
		JoinsAfter:              0,
		NodeRoleRecID:           node.Role(),
		NodeRef:                 node.ID(),
		NodePK:                  keyData,
		NodeAddress:             NewNodeAddress(node.PhysicalAddress()),
		Signature:               s,
	}, nil
}
