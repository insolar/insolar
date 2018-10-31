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

import "github.com/damnever/bitarray"

// ----------------------------------PHASE 1--------------------------------

// PulseDataExt is a pulse data extension.
type PulseDataExt struct {
	NextPulseDelta uint16
	PrevPulseDelta uint16
	EpochPulseNo   uint32
	PulseTimestamp uint32
	OriginID       *bitarray.BitArray
	Entropy        *bitarray.BitArray
}

// PulseData is a pulse data.
type PulseData struct {
	PulseNumer uint32
	Data       *PulseDataExt
}

// NewPulseData creates and returns a pulse data.
func NewPulseData() *PulseData {
	return &PulseData{
		Data: NewPulseDataExt(),
	}
}

// NewPulseDataExt creates and returns a pulse data extension.
func NewPulseDataExt() *PulseDataExt {
	return &PulseDataExt{
		OriginID: bitarray.New(128),
		Entropy:  bitarray.New(256),
	}
}

// --------------REFERENDUM--------------

type ReferendumClaim struct {
	Type   *bitarray.BitArray
	Length *bitarray.BitArray
}

// NodeBroadcast is a broadcast of info. Must be brief and only one entry per node.
// Type 4.
type NodeBroadcast struct {
	EmergencyLevel *bitarray.BitArray
}

// CapabilityPoolingAndActivation is type 3.
type CapabilityPoolingAndActivation struct {
	PollingFlags   uint16
	CapabilityType uint16
	CapabilityRef  *bitarray.BitArray
}

// NodeViolationBlame is type 2.
type NodeViolationBlame struct {
	BlameNodeID   uint32
	TypeViolation *bitarray.BitArray
}

// NodeJoinClaim is type 1, len == 272.
type NodeJoinClaim struct {
	NodeID                  uint32
	RelayNodeID             uint32
	ProtocolVersionAndFlags uint32
	JoinsAfter              uint32
	NodeRoleRecID           *bitarray.BitArray
	NodeRef                 *bitarray.BitArray
	NodePK                  *bitarray.BitArray
}

// NodeLeaveClaim can be the only be issued by the node itself and must be the only claim record.
// Should be executed with the next pulse. Type 1, len == 0.
type NodeLeaveClaim struct {
}

func NewNodeJoinClaim() *NodeJoinClaim {
	return &NodeJoinClaim{
		NodeRoleRecID: bitarray.New(256),
		NodeRef:       bitarray.New(512),
		NodePK:        bitarray.New(512),
	}
}

func NewNodViolationBlame() *NodeViolationBlame {
	return &NodeViolationBlame{
		TypeViolation: bitarray.New(8),
	}
}

func NewCapabilityPoolingAndActivation() *CapabilityPoolingAndActivation {
	return &CapabilityPoolingAndActivation{
		CapabilityRef: bitarray.New(512),
	}
}

func NewNodeBroadcast() *NodeBroadcast {
	return &NodeBroadcast{
		EmergencyLevel: bitarray.New(8),
	}
}

func NewReferendumClaim() *ReferendumClaim {
	return &ReferendumClaim{
		Type:   bitarray.New(6),
		Length: bitarray.New(10),
	}
}

// ----------------------------------PHASE 2--------------------------------

type ReferendumVote struct {
	Type   *bitarray.BitArray
	Length *bitarray.BitArray
}

type NodeListVote struct {
	NodeListCount uint16
	NodeListHash  *bitarray.BitArray
}

// NewNodeListVote creates and returns a node list vote.
func NewNodeListVote() *NodeListVote {
	return &NodeListVote{
		NodeListHash: bitarray.New(256),
	}
}

// NewReferendumVote creates and returns a referendum vote.
func NewReferendumVote() *ReferendumVote {
	return &ReferendumVote{
		Type:   bitarray.New(6),
		Length: bitarray.New(10),
	}
}
