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
	"github.com/stretchr/testify/require"
)

func makeDefaultPacketHeader() *PacketHeader {
	packetHeader := &PacketHeader{}
	packetHeader.Routing = uint8(2)
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

func checkSerialization(t *testing.T, orig interface{}, target interface{}) {
	data := serializeData(t, orig.(Serializer))
	r := bytes.NewReader(data)
	err := target.(Serializer).Deserialize(r)
	require.NoError(t, err)
	require.Equal(t, orig, target)
}

func checkBadDataSerialization(t *testing.T, orig interface{}, target interface{}, msg string) {
	require.Equal(t, reflect.TypeOf(orig), reflect.TypeOf(target), "Types must be the same")
	data := serializeData(t, orig.(Serializer))
	r := bytes.NewReader(data[:len(data)-1])
	err := target.(Serializer).Deserialize(r)
	require.EqualError(t, err, msg)
}

func TestPacketHeaderReadWrite(t *testing.T) {
	checkSerialization(t, makeDefaultPacketHeader(), &PacketHeader{})
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
	checkSerialization(t, makeDefaultPulseDataExt(), &PulseDataExt{})
}

func TestPulseDataExtReadWrite_BadData(t *testing.T) {
	checkBadDataSerialization(t, makeDefaultPulseDataExt(), &PulseDataExt{},
		"[ PulseDataExt.Deserialize ] Can't read Entropy: unexpected EOF")
}

func TestPulseDataReadWrite(t *testing.T) {
	pulseData := &PulseData{}
	pulseData.PulseNumer = uint32(32)
	pulseData.Data = makeDefaultPulseDataExt()

	checkSerialization(t, pulseData, &PulseData{})
}

func TestPulseDataReadWrite_BadData(t *testing.T) {
	pulseData := &PulseData{}
	pulseData.PulseNumer = uint32(32)
	pulseData.Data = makeDefaultPulseDataExt()
	checkBadDataSerialization(t, pulseData, &PulseData{},
		"[ PulseData.Deserialize ] Can't read PulseDataExt: [ PulseDataExt.Deserialize ] Can't read Entropy: unexpected EOF")
}

func TestNodePulseProofReadWrite(t *testing.T) {
	nodePulseProof := &NodePulseProof{}
	nodePulseProof.NodeSignature = uint64(63)
	nodePulseProof.NodeStateHash = uint64(64)
	checkSerialization(t, nodePulseProof, &NodePulseProof{})
}

func TestNodePulseProofReadWrite_BadData(t *testing.T) {
	nodePulseProof := &NodePulseProof{}
	nodePulseProof.NodeSignature = uint64(63)
	nodePulseProof.NodeStateHash = uint64(64)
	checkBadDataSerialization(t, nodePulseProof, &NodePulseProof{},
		"[ PulseDataExt.Deserialize ] Can't read Entropy: unexpected EOF")
}
