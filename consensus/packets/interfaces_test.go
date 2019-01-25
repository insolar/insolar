/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
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
