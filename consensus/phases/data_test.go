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

	"crypto/rand"

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

func checkBadDataSerializationDeserialization(t *testing.T, orig Serializer, msg string) {
	newObj := reflect.New(reflect.TypeOf(orig).Elem()).Interface()
	data := serializeData(t, orig)
	r := bytes.NewReader(data[:len(data)-1])
	err := newObj.(Serializer).Deserialize(r)
	require.EqualError(t, err, msg)
}

func TestPacketHeaderReadWrite(t *testing.T) {
	checkSerializationDeserialization(t, makeDefaultPacketHeader())
}

func TestPacketHeaderReadWrite_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeDefaultPacketHeader(),
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
	pulseDataExt.OriginID = [16]byte{}

	return pulseDataExt
}

func TestPulseDataExtReadWrite(t *testing.T) {
	checkSerializationDeserialization(t, makeDefaultPulseDataExt())
}

func TestPulseDataExtReadWrite_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeDefaultPulseDataExt(),
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
	checkBadDataSerializationDeserialization(t, pulseData,
		"[ PulseData.Deserialize ] Can't read PulseDataExt: [ PulseDataExt.Deserialize ] Can't read Entropy: unexpected EOF")
}

func genRandomSlice(n int) []byte {
	var buf = make([]byte, n)
	_, err := rand.Read(buf[:])
	if err != nil {
		panic(buf)
	}

	return buf[:]
}

func randomArray64() [64]byte {
	var buf [64]byte
	copy(buf[:], genRandomSlice(64))
	return buf
}

func randomArray32() [32]byte {
	const n = 32
	var buf [n]byte
	copy(buf[:], genRandomSlice(n))
	return buf
}

func makeNodePulseProof() *NodePulseProof {
	nodePulseProof := &NodePulseProof{}
	nodePulseProof.NodeSignature = randomArray64()
	nodePulseProof.NodeStateHash = randomArray64()

	return nodePulseProof
}

func TestNodePulseProofReadWrite(t *testing.T) {
	checkSerializationDeserialization(t, makeNodePulseProof())
}

func TestNodePulseProofReadWrite_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodePulseProof(),
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
	checkBadDataSerializationDeserialization(t, makeNodeBroadCast(),
		"[ NodeBroadcast.Deserialize ] Can't read length: unexpected EOF")
}

func makeCapabilityPoolingAndActivation() *CapabilityPoolingAndActivation {
	capabilityPoolingAndActivation := &CapabilityPoolingAndActivation{}
	capabilityPoolingAndActivation.PollingFlags = uint16(10)
	capabilityPoolingAndActivation.length = uint16(7)
	capabilityPoolingAndActivation.CapabilityType = uint16(7)
	capabilityPoolingAndActivation.CapabilityRef = randomArray64()

	return capabilityPoolingAndActivation
}

func TestCapabilityPoolingAndActivation(t *testing.T) {
	checkSerializationDeserialization(t, makeCapabilityPoolingAndActivation())
}

func TestCapabilityPoolingAndActivation_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeCapabilityPoolingAndActivation(),
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
	checkBadDataSerializationDeserialization(t, makeNodeViolationBlame(),
		"[ NodeViolationBlame.Deserialize ] Can't read length: unexpected EOF")
}

func makeNodeJoinClaim() *NodeJoinClaim {
	nodeJoinClaim := &NodeJoinClaim{}
	nodeJoinClaim.NodeID = uint32(77)
	nodeJoinClaim.RelayNodeID = uint32(26)
	nodeJoinClaim.ProtocolVersionAndFlags = uint32(99)
	nodeJoinClaim.JoinsAfter = uint32(67)
	nodeJoinClaim.NodeRoleRecID = uint32(32)
	nodeJoinClaim.NodeRef = testutils.RandomRef()
	// nodeJoinClaim.NodePK = // TODO:

	return nodeJoinClaim
}

func TestNodeJoinClaim(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeJoinClaim())
}

func TestNodeJoinClaim_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodeJoinClaim(),
		"[ NodeJoinClaim.Deserialize ] Can't read NodeRef: unexpected EOF")
}

func TestNodeLeaveClaim(t *testing.T) {
	nodeLeaveClaim := &NodeLeaveClaim{}
	nodeLeaveClaim.length = uint16(333)
	checkSerializationDeserialization(t, nodeLeaveClaim)
}

func TestNodeLeaveClaim_BadData(t *testing.T) {
	nodeLeaveClaim := &NodeLeaveClaim{}
	nodeLeaveClaim.length = uint16(333)
	checkBadDataSerializationDeserialization(t, nodeLeaveClaim,
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
	checkBadDataSerializationDeserialization(t, makeReferendumVote(),
		"[ ReferendumVote.Deserialize ] Can't read Length: unexpected EOF")
}

func makeNodeListVote() *NodeListVote {
	nodeListVote := &NodeListVote{}
	nodeListVote.NodeListHash = randomArray32()
	nodeListVote.NodeListCount = uint16(77)

	return nodeListVote
}

func TestNodeListVote(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeListVote())
}

func TestNodeListVote_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodeListVote(),
		"[ NodeListVote.Deserialize ] Can't read NodeListHash: unexpected EOF")
}

func makeDeviantBitSet() *DeviantBitSet {
	deviantBitSet := &DeviantBitSet{}
	deviantBitSet.CompressedSet = true
	deviantBitSet.HighBitLengthFlag = true
	deviantBitSet.LowBitLength = uint8(3)
	//-----------------
	deviantBitSet.HighBitLength = uint8(9)

	// TODO: uncomment it when we support reading payload
	//deviantBitSet.Payload = []byte("Hello, World!")

	return deviantBitSet
}

func TestDeviantBitSet(t *testing.T) {
	checkSerializationDeserialization(t, makeDeviantBitSet())
}

func _TestDeviantBitSet_BadData(t *testing.T) {
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

func makePhase1Packet() *Phase1Packet {
	phase1Packet := &Phase1Packet{}
	phase1Packet.packetHeader = *makeDefaultPacketHeader()
	phase1Packet.pulseData = *makeDefaultPulseDataExt()
	phase1Packet.proofNodePulse = NodePulseProof{NodeSignature: randomArray64(), NodeStateHash: randomArray64()}

	phase1Packet.claims = append(phase1Packet.claims, makeNodeJoinClaim())
	phase1Packet.claims = append(phase1Packet.claims, makeNodeViolationBlame())
	phase1Packet.claims = append(phase1Packet.claims, &NodeLeaveClaim{length: 22})

	phase1Packet.signature = 987

	return phase1Packet
}

func TestPhase1Packet_Deserialize(t *testing.T) {
	checkSerializationDeserialization(t, makePhase1Packet())
}

func TestPhase1Packet_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makePhase1Packet(),
		"[ Phase1Packet.Deserialize ] Can't parseReferendumClaim: [ PacketHeader.parseReferendumClaim ] "+
			"Can't deserialize claim.: [ NodeLeaveClaim.Deserialize ] Can't read length: unexpected EOF")

}

func makePhase2Packet() *Phase2Packet {
	phase2Packet := &Phase2Packet{}
	phase2Packet.packetHeader = *makeDefaultPacketHeader()
	phase2Packet.globuleHashSignature = randomArray64()
	phase2Packet.deviantBitSet = *makeDeviantBitSet()
	phase2Packet.signatureHeaderSection1 = randomArray64()
	phase2Packet.signatureHeaderSection2 = randomArray64()

	// TODO: uncomment when support ser\deser of ReferendumVote
	// phase2Packet.votesAndAnswers = append(phase2Packet.votesAndAnswers,*makeReferendumVote())
	// phase2Packet.votesAndAnswers = append(phase2Packet.votesAndAnswers,*makeReferendumVote())

	return phase2Packet
}

func TestPhase2Packet_Deserialize(t *testing.T) {
	checkSerializationDeserialization(t, makePhase2Packet())
}

func TestPhase2Packet_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makePhase2Packet(),
		"[ Phase2Packet.Deserialize ] Can't read signatureHeaderSection2: unexpected EOF")

}
