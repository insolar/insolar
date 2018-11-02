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
	"reflect"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeDefaultPacketHeader() *PacketHeader {
	packetHeader := &PacketHeader{}
	packetHeader.HasRouting = true
	packetHeader.SubType = 3
	packetHeader.PacketT = Referendum
	// -------------------
	packetHeader.f00 = true
	packetHeader.f01 = true
	// -------------------
	packetHeader.Pulse = uint32(22)
	packetHeader.OriginNodeID = uint32(42)
	packetHeader.TargetNodeID = uint32(62)

	return packetHeader
}

func serializeData(t *testing.T, serializer Serializer) []byte {
	data, err := serializer.Serialize()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	return data
}

func checkSerializationDeserialization(t *testing.T, orig Serializer) {
	newObj := reflect.New(reflect.TypeOf(orig).Elem()).Interface()

	data := serializeData(t, orig)
	r := bytes.NewReader(data)
	err := newObj.(Serializer).Deserialize(r)
	require.NoError(t, err)
	require.Equal(t, orig, newObj)
}

func checkBadDataSerialization(t *testing.T, orig interface{}, target interface{}, msg string) {
	require.Equal(t, reflect.TypeOf(orig), reflect.TypeOf(target), "Types must be the same")
	data := serializeData(t, orig.(Serializer))
	r := bytes.NewReader(data[:len(data)-1])
	err := target.(Serializer).Deserialize(r)
	require.EqualError(t, err, msg)
}

func TestPacketHeaderReadWrite(t *testing.T) {
	checkSerializationDeserialization(t, makeDefaultPacketHeader())
}

func TestPacketHeaderReadWrite_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeDefaultPacketHeader(), &PacketHeader{},
		"[ PacketHeader.Deserialize ] Can't read TargetNodeID: unexpected EOF")
}

func makeDefaultPulseDataExt() *PulseDataExt {
	pulseDataExt := &PulseDataExt{}
	pulseDataExt.NextPulseDelta = uint16(11)
	pulseDataExt.PrevPulseDelta = uint16(12)
	pulseDataExt.Entropy = core.Entropy{}
	pulseDataExt.Entropy[1] = '3'
	pulseDataExt.EpochPulseNo = uint32(21)
	pulseDataExt.PulseTimestamp = uint32(33)
	pulseDataExt.OriginID = uint16(43)

	return pulseDataExt
}

func TestPulseDataExtReadWrite(t *testing.T) {
	checkSerializationDeserialization(t, makeDefaultPulseDataExt())
}

func TestPulseDataExtReadWrite_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeDefaultPulseDataExt(), &PulseDataExt{},
		"[ PulseDataExt.Deserialize ] Can't read Entropy: unexpected EOF")
}

func TestPulseDataReadWrite(t *testing.T) {
	pulseData := &PulseData{}
	pulseData.PulseNumber = uint32(32)
	pulseData.Data = makeDefaultPulseDataExt()

	checkSerializationDeserialization(t, pulseData)
}

func TestPulseDataReadWrite_BadData(t *testing.T) {
	pulseData := &PulseData{}
	pulseData.PulseNumber = uint32(32)
	pulseData.Data = makeDefaultPulseDataExt()
	checkBadDataSerialization(t, pulseData, &PulseData{},
		"[ PulseData.Deserialize ] Can't read PulseDataExt: [ PulseDataExt.Deserialize ] Can't read Entropy: unexpected EOF")
}

func TestNodePulseProofReadWrite(t *testing.T) {
	nodePulseProof := &NodePulseProof{}
	nodePulseProof.NodeSignature = uint64(63)
	nodePulseProof.NodeStateHash = uint64(64)
	checkSerializationDeserialization(t, nodePulseProof)
}

func TestNodePulseProofReadWrite_BadData(t *testing.T) {
	nodePulseProof := &NodePulseProof{}
	nodePulseProof.NodeSignature = uint64(63)
	nodePulseProof.NodeStateHash = uint64(64)
	checkBadDataSerialization(t, nodePulseProof, &NodePulseProof{},
		"[ NodePulseProof.Deserialize ] Can't read NodeSignature: unexpected EOF")
}

func makeNodeBroadCast() *NodeBroadcast {
	nodeBroadcast := &NodeBroadcast{}
	nodeBroadcast.length = uint16(3)
	nodeBroadcast.EmergencyLevel = uint8(4)

	return nodeBroadcast
}

func TestNodeBroadcast(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeBroadCast())
}

func TestNodeBroadcast_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeNodeBroadCast(), &NodeBroadcast{},
		"[ NodeBroadcast.Deserialize ] Can't read length: unexpected EOF")
}

func makeCapabilityPoolingAndActivation() *CapabilityPoolingAndActivation {
	capabilityPoolingAndActivation := &CapabilityPoolingAndActivation{}
	capabilityPoolingAndActivation.PollingFlags = uint16(10)
	capabilityPoolingAndActivation.length = uint16(7)
	capabilityPoolingAndActivation.CapabilityType = uint16(11)
	capabilityPoolingAndActivation.CapabilityRef = uint64(13)

	return capabilityPoolingAndActivation
}

func TestCapabilityPoolingAndActivation(t *testing.T) {
	checkSerializationDeserialization(t, makeCapabilityPoolingAndActivation())
}

func TestCapabilityPoolingAndActivation_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeCapabilityPoolingAndActivation(), &CapabilityPoolingAndActivation{},
		"[ CapabilityPoolingAndActivation.Deserialize ] Can't read length: unexpected EOF")
}

func makeNodeViolationBlame() *NodeViolationBlame {
	nodeViolationBlame := &NodeViolationBlame{}
	nodeViolationBlame.length = uint16(2)
	nodeViolationBlame.claimType = TypeNodeViolationBlame
	nodeViolationBlame.TypeViolation = uint8(4)

	return nodeViolationBlame
}

func TestNodeViolationBlame(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeViolationBlame())
}

func TestNodeViolationBlame_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeNodeViolationBlame(), &NodeViolationBlame{},
		"[ NodeViolationBlame.Deserialize ] Can't read length: unexpected EOF")
}

func makeNodeJoinClaim() *NodeJoinClaim {
	nodeJoinClaim := &NodeJoinClaim{}
	nodeJoinClaim.length = uint16(2)
	nodeJoinClaim.NodeRoleRecID = uint32(32)
	nodeJoinClaim.JoinsAfter = uint32(67)
	nodeJoinClaim.ProtocolVersionAndFlags = uint32(99)
	nodeJoinClaim.NodeID = uint32(77)
	nodeJoinClaim.NodeRef = testutils.RandomRef()
	nodeJoinClaim.RelayNodeID = uint32(26)
	// nodeJoinClaim.NodePK = // TODO:

	return nodeJoinClaim
}

func TestNodeJoinClaim(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeJoinClaim())
}

func TestNodeJoinClaim_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeNodeJoinClaim(), &NodeJoinClaim{},
		"[ NodeJoinClaim.Deserialize ] Can't read length: unexpected EOF")
}

func TestNodeLeaveClaim(t *testing.T) {
	nodeLeaveClaim := &NodeLeaveClaim{}
	nodeLeaveClaim.length = uint16(333)
	checkSerializationDeserialization(t, nodeLeaveClaim)
}

func TestNodeLeaveClaim_BadData(t *testing.T) {
	nodeLeaveClaim := &NodeLeaveClaim{}
	nodeLeaveClaim.length = uint16(333)
	checkBadDataSerialization(t, nodeLeaveClaim, &NodeLeaveClaim{},
		"[ NodeLeaveClaim.Deserialize ] Can't read length: unexpected EOF")
}

// ----------------------------------PHASE 2--------------------------------

func makeReferendumVote() *ReferendumVote {
	referendumVote := &ReferendumVote{}
	referendumVote.Length = uint16(44)
	referendumVote.Type = ReferendumType(23)

	return referendumVote
}

func TestReferendumVote(t *testing.T) {
	checkSerializationDeserialization(t, makeReferendumVote())
}

func TestReferendumVote_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeReferendumVote(), &ReferendumVote{},
		"[ ReferendumVote.Deserialize ] Can't read Length: unexpected EOF")
}

func makeNodeListVote() *NodeListVote {
	nodeListVote := &NodeListVote{}
	nodeListVote.NodeListHash = uint32(13)
	nodeListVote.NodeListCount = uint16(77)

	return nodeListVote
}

func TestNodeListVote(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeListVote())
}

func TestNodeListVote_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeNodeListVote(), &NodeListVote{},
		"[ NodeListVote.Deserialize ] Can't read NodeListHash: unexpected EOF")
}

func makeDeviantBitSet() *DeviantBitSet {
	deviantBitSet := &DeviantBitSet{}
	deviantBitSet.CompressedSet = true
	deviantBitSet.HighBitLengthFlag = true
	deviantBitSet.LowBitLength = uint8(3)
	//-----------------
	deviantBitSet.HighBitLength = uint8(9)
	deviantBitSet.Payload = []byte("Hello, World!")

	return deviantBitSet
}

func TestDeviantBitSet(t *testing.T) {
	checkSerializationDeserialization(t, makeDeviantBitSet())
}

func TestDeviantBitSet_BadData(t *testing.T) {
	deviantBitSet := makeDeviantBitSet()
	newDeviantBitSet := &DeviantBitSet{}

	data := serializeData(t, deviantBitSet)
	r := bytes.NewReader(data[:len(data)-2])
	err := newDeviantBitSet.Deserialize(r)
	assert.NoError(t, err)

	require.NotEqual(t, deviantBitSet.Payload, newDeviantBitSet.Payload)
}

func TestParseAndCompactRouteInfo(t *testing.T) {
	var routInfoTests = []PacketHeader{
		PacketHeader{
			PacketT:    NetworkConsistency,
			SubType:    1,
			HasRouting: true,
		},
		PacketHeader{
			PacketT:    NetworkConsistency,
			SubType:    1,
			HasRouting: false,
		},
		PacketHeader{
			PacketT:    Referendum,
			SubType:    0,
			HasRouting: false,
		},
	}

	for _, ph := range routInfoTests {
		raw := ph.compactRouteInfo()
		newPh := PacketHeader{}
		newPh.parseRouteInfo(raw)
		require.Equal(t, ph, newPh)
	}
}

func TestParseAndCompactPulseAndCustomFlags(t *testing.T) {
	var pulseAndCustomFlagsTests = []PacketHeader{
		PacketHeader{
			f00:   true,
			f01:   true,
			Pulse: 0,
		},
		PacketHeader{
			f00:   false,
			f01:   true,
			Pulse: 1,
		},
		PacketHeader{
			f00:   true,
			f01:   false,
			Pulse: 2,
		},
		PacketHeader{
			f00:   false,
			f01:   false,
			Pulse: 2,
		},
		PacketHeader{
			f00:   false,
			f01:   false,
			Pulse: 0,
		},
	}

	for _, ph := range pulseAndCustomFlagsTests {
		raw := ph.compactPulseAndCustomFlags()
		newPh := PacketHeader{}
		newPh.parsePulseAndCustomFlags(raw)
		require.Equal(t, ph, newPh)
	}

}
