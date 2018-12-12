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
	"github.com/insolar/insolar/core"
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

// NodeJoinClaim is a type 1, len == 272.
type NodeJoinClaim struct {
	ShortNodeID             core.ShortNodeID
	RelayNodeID             core.ShortNodeID
	ProtocolVersionAndFlags uint32
	JoinsAfter              uint32
	NodeRoleRecID           uint32
	NodeRef                 core.RecordRef
	NodePK                  [PublicKeyLength]byte
	Signature               [SignatureLength]byte
}

func (njc *NodeJoinClaim) Type() ClaimType {
	return TypeNodeJoinClaim
}

// NodeJoinClaim is a type 5, len == 272.
type NodeAnnounceClaim struct {
	NodeJoinClaim

	NodeIndex uint16
	NodeCount uint16
}

func (nac *NodeAnnounceClaim) Type() ClaimType {
	return TypeNodeAnnounceClaim
}

// NodeLeaveClaim can be the only be issued by the node itself and must be the only claim record.
// Should be executed with the next pulse. Type 1, len == 0.
type NodeLeaveClaim struct {
}

func (nlc *NodeLeaveClaim) Type() ClaimType {
	return TypeNodeLeaveClaim
}
