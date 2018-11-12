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
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/stretchr/testify/require"
)

func checkGetAndSetPacketHeader(t *testing.T, packet PacketRoutable, header *RoutingHeader) {
	err := packet.SetPacketHeader(header)
	require.NoError(t, err)

	newHeader, err := packet.GetPacketHeader()
	require.NoError(t, err)

	require.Equal(t, header, newHeader)
}

func checkSetPacketHeader_BadType(t *testing.T, packet PacketRoutable, header *RoutingHeader) {
	err := packet.SetPacketHeader(header)
	require.Contains(t, err.Error(), "wrong packet type")
}

func makeHeader(packetType types.PacketType) *RoutingHeader {
	header := &RoutingHeader{}
	header.PacketType = packetType
	header.OriginID = core.ShortNodeID(23)
	header.TargetID = core.ShortNodeID(33)

	return header
}

func TestPhase1Packet_GetAndSetPacketHeader(t *testing.T) {
	checkGetAndSetPacketHeader(t, &Phase1Packet{}, makeHeader(types.Phase1))
}

func TestPhase2Packet_SetAndGetPacketHeader(t *testing.T) {
	checkGetAndSetPacketHeader(t, &Phase2Packet{}, makeHeader(types.Phase2))
}

func TestPhase1Packet_SetPacketHeader_BadType(t *testing.T) {
	checkSetPacketHeader_BadType(t, &Phase1Packet{}, makeHeader(types.Phase2))
}

func TestPhase2Packet_SetPacketHeader_BadType(t *testing.T) {
	checkSetPacketHeader_BadType(t, &Phase2Packet{}, makeHeader(types.Phase1))
}
