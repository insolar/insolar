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

package phases

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type PacketType uint8
type ClaimType uint8
type ReferendumType uint8

const (
	Type1 = PacketType(iota + 1)
)

const (
	TypeNodeClaim = ClaimType(iota + 1)
	TypeNodeViolationBlame
	TypeCapabilityPollingAndActivation
	TypeNodeBroadcast
)

// ----------------------------------PHASE 1--------------------------------

var defaultByteOrder = binary.BigEndian

type PacketHeader struct {
	Routing      uint8
	Pulse        uint32
	OriginNodeID uint32
	TargetNodeID uint32
}

func (ph *PacketHeader) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &ph.Routing)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read Routing")
	}

	err = binary.Read(data, defaultByteOrder, &ph.Pulse)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read Pulse")
	}

	err = binary.Read(data, defaultByteOrder, &ph.OriginNodeID)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read OriginNodeID")
	}

	err = binary.Read(data, defaultByteOrder, &ph.TargetNodeID)
	if err != nil {
		return errors.Wrap(err, "[ PacketHeader.Deserialize ] Can't read TargetNodeID")
	}

	return nil
}

func (ph *PacketHeader) Serialize() ([]byte, error) {
	result := new(bytes.Buffer)
	err := binary.Write(result, defaultByteOrder, ph.Routing)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't read Routing")
	}

	err = binary.Write(result, defaultByteOrder, ph.Pulse)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't read Pulse")
	}

	err = binary.Write(result, defaultByteOrder, ph.OriginNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't read OriginNodeID")
	}

	err = binary.Write(result, defaultByteOrder, ph.TargetNodeID)
	if err != nil {
		return nil, errors.Wrap(err, "[ PacketHeader.Serialize ] Can't read TargetNodeID")
	}

	return result.Bytes(), nil
}

// PulseDataExt is a pulse data extension.
type PulseDataExt struct {
	NextPulseDelta uint16
	PrevPulseDelta uint16
	OriginID       uint16
	EpochPulseNo   uint32
	PulseTimestamp uint32
	Entropy        core.Entropy
}

func (pde *PulseDataExt) Deserialize(data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &pde.NextPulseDelta)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read NextPulseDelta")
	}

	err = binary.Read(data, defaultByteOrder, &pde.PrevPulseDelta)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read PrevPulseDelta")
	}

	err = binary.Read(data, defaultByteOrder, &pde.OriginID)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read OriginID")
	}

	err = binary.Read(data, defaultByteOrder, &pde.EpochPulseNo)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read EpochPulseNo")
	}

	err = binary.Read(data, defaultByteOrder, &pde.PulseTimestamp)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read PulseTimestamp")
	}

	err = binary.Read(data, defaultByteOrder, &pde.Entropy)
	if err != nil {
		return errors.Wrap(err, "[ PulseDataExt.Deserialize ] Can't read Entropy")
	}

	return nil
}

func (pde *PulseDataExt) Serialize() ([]byte, error) {
	result := new(bytes.Buffer)
	err := binary.Write(result, defaultByteOrder, pde.NextPulseDelta)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't read NextPulseDelta")
	}

	err = binary.Write(result, defaultByteOrder, pde.OriginID)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't read OriginID")
	}

	err = binary.Write(result, defaultByteOrder, pde.EpochPulseNo)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't read EpochPulseNo")
	}

	err = binary.Write(result, defaultByteOrder, pde.PulseTimestamp)
	if err != nil {
		return nil, errors.Wrap(err, "[ PulseDataExt.Serialize ] Can't read PulseTimestamp")
	}

	return result.Bytes(), nil
}

// PulseData is a pulse data.
type PulseData struct {
	PulseNumer uint32
	Data       *PulseDataExt
}

type NodePulseProof struct {
	NodeStateHash uint64
	NodeSignature uint64
}

// --------------REFERENDUM--------------

type ReferendumClaim interface {
	Type() ClaimType
	Length() uint16
}

// NodeBroadcast is a broadcast of info. Must be brief and only one entry per node.
// Type 4.
type NodeBroadcast struct {
	EmergencyLevel uint8
	claimType      ClaimType
	length         uint16
}

func (nb *NodeBroadcast) Type() ClaimType {
	return nb.claimType
}

func (nb *NodeBroadcast) Length() uint16 {
	return nb.length
}

// CapabilityPoolingAndActivation is a type 3.
type CapabilityPoolingAndActivation struct {
	PollingFlags   uint16
	CapabilityType uint16
	CapabilityRef  uint64
	claimType      ClaimType
	length         uint16
}

func (cpa *CapabilityPoolingAndActivation) Type() ClaimType {
	return cpa.claimType
}

func (cpa *CapabilityPoolingAndActivation) Length() uint16 {
	return cpa.length
}

// NodeViolationBlame is a type 2.
type NodeViolationBlame struct {
	BlameNodeID   uint32
	TypeViolation uint8
	claimType     ClaimType
	length        uint16
}

func (nvb *NodeViolationBlame) Type() ClaimType {
	return nvb.claimType
}

func (nvb *NodeViolationBlame) Length() uint16 {
	return nvb.length
}

// NodeJoinClaim is a type 1, len == 272.
type NodeJoinClaim struct {
	NodeID                  uint32
	RelayNodeID             uint32
	ProtocolVersionAndFlags uint32
	JoinsAfter              uint32
	NodeRoleRecID           uint32
	NodeRef                 core.RecordRef
	NodePK                  ecdsa.PrivateKey
	claimType               ClaimType
	length                  uint16
}

func (njc *NodeJoinClaim) Type() ClaimType {
	return njc.claimType
}

func (njc *NodeJoinClaim) Length() uint16 {
	return njc.length
}

// NodeLeaveClaim can be the only be issued by the node itself and must be the only claim record.
// Should be executed with the next pulse. Type 1, len == 0.
type NodeLeaveClaim struct {
	claimType ClaimType
	length    uint16
}

func (nlc *NodeLeaveClaim) Type() ClaimType {
	return nlc.claimType
}

func (nlc *NodeLeaveClaim) Length() uint16 {
	return nlc.length
}

func NewNodeLeaveClaim() *NodeLeaveClaim {
	return &NodeLeaveClaim{
		claimType: TypeNodeClaim,
	}
}

func NewNodeJoinClaim() *NodeJoinClaim {
	return &NodeJoinClaim{
		claimType: TypeNodeClaim,
		length:    272,
	}
}

func NewNodViolationBlame() *NodeViolationBlame {
	return &NodeViolationBlame{
		claimType: TypeNodeViolationBlame,
	}
}

func NewCapabilityPoolingAndActivation() *CapabilityPoolingAndActivation {
	return &CapabilityPoolingAndActivation{
		claimType: TypeCapabilityPollingAndActivation,
	}
}

func NewNodeBroadcast() *NodeBroadcast {
	return &NodeBroadcast{
		claimType: TypeNodeBroadcast,
	}
}

// ----------------------------------PHASE 2--------------------------------

type ReferendumVote struct {
	Type   ReferendumType
	Length uint16
}

type NodeListVote struct {
	NodeListCount uint16
	NodeListHash  uint32
}
