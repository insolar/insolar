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
	"bytes"
	"crypto/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeDefaultPacketHeader(packetType PacketType) *PacketHeader {
	packetHeader := &PacketHeader{}
	packetHeader.HasRouting = true
	packetHeader.PacketT = packetType
	// -------------------
	packetHeader.f00 = true
	packetHeader.f01 = true
	// -------------------
	packetHeader.Pulse = uint32(22)
	packetHeader.OriginNodeID = core.ShortNodeID(42)
	packetHeader.TargetNodeID = core.ShortNodeID(62)

	return packetHeader
}

func serializeData(t *testing.T, serializer Serializer) []byte {
	data, err := serializer.Serialize()
	require.NoError(t, err)
	//TODO: require.NotEmpty(t, data) - need to fix test, coz some claims serializes to empty []byte

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
	checkSerializationDeserialization(t, makeDefaultPacketHeader(Phase1))
}

func TestPacketHeaderReadWrite_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeDefaultPacketHeader(Phase1),
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

func randomArray71() [SignatureLength]byte {
	var buf [SignatureLength]byte
	copy(buf[:], genRandomSlice(SignatureLength))
	return buf
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
	nodePulseProof.NodeSignature = randomArray71()
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

func TestPhase1Packet_SetPulseProof(t *testing.T) {
	p := NewPhase1Packet()
	proofStateHash := genRandomSlice(HashLength)
	proofSignature := genRandomSlice(SignatureLength)

	err := p.SetPulseProof(proofStateHash, proofSignature)
	assert.NoError(t, err)

	assert.Equal(t, proofStateHash, p.proofNodePulse.NodeStateHash[:])
	assert.Equal(t, proofSignature, p.proofNodePulse.NodeSignature[:])

	invalidStateHash := genRandomSlice(32)
	invalidSignature := genRandomSlice(128)
	err = p.SetPulseProof(invalidStateHash, invalidSignature)
	assert.Error(t, err)

	assert.NotEqual(t, invalidStateHash, p.proofNodePulse.NodeStateHash[:])
	assert.NotEqual(t, invalidSignature, p.proofNodePulse.NodeSignature[:])

}

// ----------------------------------PHASE 2--------------------------------

func TestParseAndCompactRouteInfo(t *testing.T) {
	var routInfoTests = []PacketHeader{
		PacketHeader{
			PacketT:    Phase1,
			HasRouting: true,
		},
		PacketHeader{
			PacketT:    Phase1,
			HasRouting: false,
		},
		PacketHeader{
			PacketT:    Phase2,
			HasRouting: true,
		},
		PacketHeader{
			PacketT:    Phase2,
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
	phase1Packet := NewPhase1Packet()
	phase1Packet.packetHeader = *makeDefaultPacketHeader(Phase1)
	phase1Packet.pulseData = *makeDefaultPulseDataExt()
	phase1Packet.proofNodePulse = NodePulseProof{NodeSignature: randomArray71(), NodeStateHash: randomArray64()}

	phase1Packet.AddClaim(makeNodeJoinClaim())
	phase1Packet.AddClaim(makeNodeViolationBlame())
	phase1Packet.AddClaim(&NodeLeaveClaim{})

	phase1Packet.Signature = randomArray71()

	return phase1Packet
}

func TestPhase1Packet_Deserialize(t *testing.T) {
	checkSerializationDeserialization(t, makePhase1Packet())
}

func makePhase2Packet() *Phase2Packet {
	phase2Packet := &Phase2Packet{}
	phase2Packet.packetHeader = *makeDefaultPacketHeader(Phase2)
	phase2Packet.globuleHashSignature = randomArray64()
	phase2Packet.SignatureHeaderSection1 = randomArray71()
	phase2Packet.SignatureHeaderSection2 = randomArray71()

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
		"[ Phase2Packet.Deserialize ] Can't deserialize body: [ Phase2Packet.Deserialize ] "+
			"Can't read SignatureHeaderSection2: unexpected EOF")

}

func checkExtractPacket(t *testing.T, packet Serializer) {
	data, err := packet.Serialize()
	require.NoError(t, err)

	buf := bytes.NewReader(data)
	consensusPacket, err := ExtractPacket(buf)
	require.NoError(t, err)

	newRawPacket, err := consensusPacket.Serialize()
	require.NoError(t, err)

	assert.Equal(t, data, newRawPacket)
}

func TestExtractPacket_Phase1(t *testing.T) {
	checkExtractPacket(t, makePhase1Packet())
}

func TestExtractPacket_Phase2(t *testing.T) {
	checkExtractPacket(t, makePhase2Packet())
}

func TestExtractPacket_BadHeader(t *testing.T) {
	reader := strings.NewReader("1")
	_, err := ExtractPacket(reader)
	require.EqualError(t, err, "[ ExtractPacket ] Can't read packet header")
}

func checkWrongPacket(t *testing.T, packet Serializer) {
	data, err := packet.Serialize()
	require.NoError(t, err)

	buf := bytes.NewReader(data[:(len(data)-1)/3])
	_, err = ExtractPacket(buf)
	require.Contains(t, err.Error(), "Can't DeserializeWithoutHeader")
}

func TestExtractPacket_Phase2_BadExtract(t *testing.T) {
	packet := makePhase2Packet()
	packet.packetHeader.PacketT = Phase1
	checkWrongPacket(t, packet)
}

func TestExtractPacket_Phase1_BadExtract(t *testing.T) {
	packet := makePhase1Packet()
	packet.packetHeader.PacketT = Phase2
	checkWrongPacket(t, packet)
}

func TestPhase1Packet_AddClaim(t *testing.T) {
	packet := makePhase1Packet()

	success := packet.AddClaim(makeNodeJoinClaim())
	assert.True(t, success)

	for success {
		success = packet.AddClaim(&NodeLeaveClaim{})
	}
	assert.False(t, success)
}

func TestPhase3Packet_Serialize(t *testing.T) {
	checkSerializationDeserialization(t, getPhase3Packet(t))
}

func getPhase3Packet(t *testing.T) *Phase3Packet {
	packet := &Phase3Packet{}
	packet.packetHeader = *makeDefaultPacketHeader(Phase3)
	packet.globuleHashSignature = randomArray71()
	packet.SignatureHeaderSection1 = randomArray71()
	var err error
	packet.deviantBitSet, err = NewBitSet(100)
	assert.NoError(t, err)

	refs := initRefs()
	cells := initBitCells(refs)
	bitset, err := NewBitSet(len(cells))
	assert.NoError(t, err)

	packet.deviantBitSet = bitset

	return packet
}
