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

package packets

import (
	"bytes"
	"crypto"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

//go:generate stringer -type=ClaimType

type ClaimType uint8

const (
	TypeNodeJoinClaim      = ClaimType(1)
	TypeNodeAnnounceClaim  = ClaimType(2)
	TypeNodeLeaveClaim     = ClaimType(3)
	TypeChangeNetworkClaim = ClaimType(4)
)

const claimHeaderSize = 2

// ChangeNetworkClaim uses to change network state.
type ChangeNetworkClaim struct {
	Address string
}

func (cnc *ChangeNetworkClaim) Type() ClaimType {
	return TypeChangeNetworkClaim
}

func (cnc *ChangeNetworkClaim) Clone() ReferendumClaim {
	result := *cnc
	return &result
}

type ReferendumClaim interface {
	Serializer
	Type() ClaimType
	Clone() ReferendumClaim
}

type ClaimSupplementary interface {
	AddSupplementaryInfo(nodeID insolar.Reference)
}

type SignedClaim interface {
	GetNodeID() insolar.Reference
	GetPublicKey() (crypto.PublicKey, error)
	SerializeRaw() ([]byte, error)
	GetSignature() []byte
}

// NodeJoinClaim is a type 1, len == 272.
type NodeJoinClaim struct {
	ShortNodeID             insolar.ShortNodeID
	RelayNodeID             insolar.ShortNodeID
	ProtocolVersionAndFlags uint32
	JoinsAfter              uint32
	NodeRoleRecID           insolar.StaticRole
	NodeRef                 insolar.Reference
	NodeAddress             NodeAddress
	NodePK                  [PublicKeyLength]byte
	Signature               [SignatureLength]byte
}

func (njc *NodeJoinClaim) Clone() ReferendumClaim {
	result := *njc
	return &result
}

func (njc *NodeJoinClaim) GetNodeID() insolar.Reference {
	return njc.NodeRef
}

func (njc *NodeJoinClaim) GetPublicKey() (crypto.PublicKey, error) {
	keyProc := platformpolicy.NewKeyProcessor()
	return keyProc.ImportPublicKeyBinary(njc.NodePK[:])
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
	CloudHash          [HashLength]byte

	// mapper is used to fill three fields above, is not serialized
	BitSetMapper BitSetMapper
}

func (nac *NodeAnnounceClaim) Clone() ReferendumClaim {
	result := *nac
	return &result
}

func (nac *NodeAnnounceClaim) Type() ClaimType {
	return TypeNodeAnnounceClaim
}

func (nac *NodeAnnounceClaim) SetCloudHash(cloudHash []byte) {
	copy(nac.CloudHash[:], cloudHash[:HashLength])
}

// NodeLeaveClaim can be the only be issued by the node itself and must be the only claim record.
// Should be executed with the next pulse. Type 1, len == 0.
type NodeLeaveClaim struct {
	// additional field that is not serialized and is set from transport layer on packet receive
	NodeID insolar.Reference
	ETA    insolar.PulseNumber
}

func (nlc *NodeLeaveClaim) Clone() ReferendumClaim {
	result := *nlc
	return &result
}

func (nlc *NodeLeaveClaim) AddSupplementaryInfo(nodeID insolar.Reference) {
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

// If you need a signed join claim use NodeKeeper.GetOriginJoinClaim()
func NodeToClaim(node insolar.NetworkNode) (*NodeJoinClaim, error) {
	keyProc := platformpolicy.NewKeyProcessor()
	exportedKey, err := keyProc.ExportPublicKeyBinary(node.PublicKey())
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeToClaim ] failed to export a public key")
	}

	address, err := NewNodeAddress(node.Address())
	if err != nil {
		return nil, errors.Wrap(err, "[ NodeToClaim ] failed to convert node address")
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
		NodeAddress:             address,
		Signature:               s,
	}, nil
}

func (njc *NodeJoinClaim) Marshal() ([]byte, error) {
	return njc.Serialize()
}

func (njc *NodeJoinClaim) MarshalTo(data []byte) (int, error) {
	tmp, err := njc.Serialize()
	if err != nil {
		return 0, errors.New("Error serializing NodeJoinClaim")
	}
	copy(data, tmp)
	return len(tmp), nil
}

func (njc *NodeJoinClaim) Unmarshal(data []byte) error {
	if len(data) != njc.Size() {
		return errors.New("Not enough bytes to unpack NodeJoinClaim")
	}
	return njc.Deserialize(bytes.NewReader(data))
}

func (njc *NodeJoinClaim) Size() int {
	return int(claimSizeMap[TypeNodeJoinClaim])
}

func (njc *NodeJoinClaim) Compare(other NodeJoinClaim) int {
	return njc.NodeRef.Compare(other.NodeRef)
}

func (njc *NodeJoinClaim) Equal(other NodeJoinClaim) bool {
	return njc.Compare(other) == 0
}
