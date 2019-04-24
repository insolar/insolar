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
	"crypto/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/insolar/insolar/insolar"
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
	packetHeader.OriginNodeID = insolar.ShortNodeID(42)
	packetHeader.TargetNodeID = insolar.ShortNodeID(62)

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
	require.Error(t, err)
	assert.Contains(t, err.Error(), msg)
}

func TestPacketHeaderReadWrite(t *testing.T) {
	checkSerializationDeserialization(t, makeDefaultPacketHeader(Phase1))
}

func TestPacketHeaderReadWrite_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeDefaultPacketHeader(Phase1), "unexpected EOF")
}

func makeDefaultPulseDataExt() PulseDataExt {
	pulseDataExt := PulseDataExt{}
	pulseDataExt.NextPulseDelta = uint16(11)
	pulseDataExt.PrevPulseDelta = uint16(12)
	pulseDataExt.Entropy = insolar.Entropy{}
	pulseDataExt.Entropy[1] = '3'
	pulseDataExt.EpochPulseNo = uint32(21)
	pulseDataExt.PulseTimestamp = int64(33)
	pulseDataExt.OriginID = [16]byte{}

	return pulseDataExt
}

func TestPulseDataExtReadWrite(t *testing.T) {
	pulseDataExt := makeDefaultPulseDataExt()
	checkSerializationDeserialization(t, &pulseDataExt)
}

func TestPulseDataExtReadWrite_BadData(t *testing.T) {
	pulseDataExt := makeDefaultPulseDataExt()
	checkBadDataSerializationDeserialization(t, &pulseDataExt, "unexpected EOF")
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
	checkBadDataSerializationDeserialization(t, pulseData, "unexpected EOF")
}

func genRandomSlice(n int) []byte {
	var buf = make([]byte, n)
	_, err := rand.Read(buf[:])
	if err != nil {
		panic(buf)
	}

	return buf[:]
}

func randomArray66() [SignatureLength]byte {
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
	nodePulseProof.NodeSignature = randomArray66()
	nodePulseProof.NodeStateHash = randomArray64()

	return nodePulseProof
}

func TestNodePulseProofReadWrite(t *testing.T) {
	checkSerializationDeserialization(t, makeNodePulseProof())
}

func TestNodePulseProofReadWrite_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodePulseProof(), "unexpected EOF")
}

func TestPhase1Packet_SetPulseProof(t *testing.T) {
	p := NewPhase1Packet(insolar.Pulse{})
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
	phase1Packet := NewPhase1Packet(insolar.Pulse{})
	phase1Packet.packetHeader = *makeDefaultPacketHeader(Phase1)
	phase1Packet.pulseData = makeDefaultPulseDataExt()
	phase1Packet.proofNodePulse = NodePulseProof{NodeSignature: randomArray66(), NodeStateHash: randomArray64()}

	phase1Packet.AddClaim(makeNodeJoinClaim(true))
	phase1Packet.AddClaim(makeNodeViolationBlame())
	phase1Packet.AddClaim(&NodeLeaveClaim{})

	phase1Packet.Signature = randomArray66()

	return phase1Packet
}

func TestPhase1Packet_Deserialize(t *testing.T) {
	checkSerializationDeserialization(t, makePhase1Packet())
}

func makePhase2Packet() *Phase2Packet {
	phase2Packet := &Phase2Packet{}
	phase2Packet.packetHeader = *makeDefaultPacketHeader(Phase2)
	phase2Packet.globuleHashSignature = randomArray66()
	phase2Packet.SignatureHeaderSection1 = randomArray66()
	phase2Packet.SignatureHeaderSection2 = randomArray66()
	phase2Packet.bitSet, _ = NewBitSet(134)

	vote := &MissingNode{NodeIndex: 25}

	phase2Packet.AddVote(vote)

	return phase2Packet
}

func TestPhase2Packet_Deserialize(t *testing.T) {
	checkSerializationDeserialization(t, makePhase2Packet())
}

func TestPhase2Packet_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makePhase2Packet(), "unexpected EOF")

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

	success := packet.AddClaim(makeNodeJoinClaim(true))
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
	packet.globuleHashSignature = randomArray66()
	packet.SignatureHeaderSection1 = randomArray66()
	var err error
	packet.bitset, err = NewBitSet(100)
	assert.NoError(t, err)

	count := 70
	refs := initRefs(count)
	cells := initBitCells(refs)
	bitset, err := NewBitSet(len(cells))
	assert.NoError(t, err)

	packet.bitset = bitset

	return packet
}
