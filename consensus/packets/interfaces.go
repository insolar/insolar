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
	"io"
	"strconv"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/packet/types"
)

type RoutingHeader struct {
	OriginID   core.ShortNodeID
	TargetID   core.ShortNodeID
	PacketType types.PacketType
}

type PacketRoutable interface {
	// SetPacketHeader set routing information for transport level.
	SetPacketHeader(header *RoutingHeader) error
	// GetPacketHeader get routing information from transport level.
	GetPacketHeader() (*RoutingHeader, error)
}

type Serializer interface {
	Serialize() ([]byte, error)
	Deserialize(data io.Reader) error
}

type HeaderSkipDeserializer interface {
	DeserializeWithoutHeader(data io.Reader, header *PacketHeader) error
}

type ConsensusPacket interface {
	HeaderSkipDeserializer
	Serializer
	PacketRoutable
}

func ExtractPacket(reader io.Reader) (ConsensusPacket, error) {
	header := PacketHeader{}
	err := header.Deserialize(reader)
	if err != nil {
		return nil, errors.New("[ ExtractPacket ] Can't read packet header")
	}

	var packet ConsensusPacket
	switch header.PacketT {
	case Phase1:
		packet = &Phase1Packet{}
	case Phase2:
		packet = &Phase2Packet{}
	case Phase3:
		packet = &Phase3Packet{}
	default:
		return nil, errors.New("[ ExtractPacket ] Unknown extract packet type. " + strconv.Itoa(int(header.PacketT)))
	}

	err = packet.DeserializeWithoutHeader(reader, &header)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExtractPacket ] Can't DeserializeWithoutHeader packet")
	}

	return packet, nil
}
