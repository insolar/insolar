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
	"testing"

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

func serializePacketHeader(t *testing.T, packetHeader *PacketHeader) []byte {
	data, err := packetHeader.Serialize()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	return data
}

func TestPacketHeaderReadWrite(t *testing.T) {
	packetHeader := makeDefaultPacketHeader()
	data := serializePacketHeader(t, packetHeader)

	newPacketHeader := PacketHeader{}
	r := bytes.NewReader(data)
	err := newPacketHeader.Deserialize(r)
	require.NoError(t, err)

	require.Equal(t, packetHeader, newPacketHeader)
}

func TestPacketHeaderReadWrite_BadData(t *testing.T) {
	data := serializePacketHeader(t, makeDefaultPacketHeader())
	newPacketHeader := PacketHeader{}
	r := bytes.NewReader(data[:len(data)-2])
	err := newPacketHeader.Deserialize(r)
	require.EqualError(t, err, "[ Deserialize ] Can't read TargetNodeID: unexpected EOF")
}
