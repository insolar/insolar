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

package packet

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/pkg/errors"
)

func (p *PacketBackend) SetRequest(request interface{}) {
	var r isRequest_Request
	switch t := request.(type) {
	case *Ping:
		r = &Request_Ping{t}
	case *RPCRequest:
		r = &Request_RPC{t}
	case *CascadeRequest:
		r = &Request_Cascade{t}
	case *PulseRequest:
		r = &Request_Pulse{t}
	case *BootstrapRequest:
		r = &Request_Bootstrap{t}
	case *AuthorizeRequest:
		r = &Request_Authorize{t}
	case *RegisterRequest:
		r = &Request_Register{t}
	case *GenesisRequest:
		r = &Request_Genesis{t}
	default:
		panic("Request payload is not a valid protobuf struct!")
	}
	p.Payload = &PacketBackend_Request{Request: &Request{Request: r}}
}

func (p *PacketBackend) SetResponse(response interface{}) {
	var r isResponse_Response
	switch t := response.(type) {
	case *Ping:
		r = &Response_Ping{t}
	case *RPCResponse:
		r = &Response_RPC{t}
	case *BasicResponse:
		r = &Response_Basic{t}
	case *BootstrapResponse:
		r = &Response_Bootstrap{t}
	case *AuthorizeResponse:
		r = &Response_Authorize{t}
	case *RegisterResponse:
		r = &Response_Register{t}
	case *GenesisResponse:
		r = &Response_Genesis{t}
	default:
		panic("Response payload is not a valid protobuf struct!")
	}
	p.Payload = &PacketBackend_Response{Response: &Response{Response: r}}
}

func (p *PacketBackend) GetType() types.PacketType {
	// TODO: make p.Type of type PacketType instead of uint32
	return types.PacketType(p.Type)
}

func (p *PacketBackend) GetSender() insolar.Reference {
	return p.Sender.NodeID
}

func (p *PacketBackend) GetSenderHost() *host.Host {
	return p.Sender
}

func (p *PacketBackend) GetRequestID() types.RequestID {
	return types.RequestID(p.RequestID)
}

func (p *PacketBackend) IsResponse() bool {
	return p.GetResponse() != nil
}

// SerializePacketBackend converts packet to byte slice.
func SerializePacketBackend(p *PacketBackend) ([]byte, error) {
	data, err := p.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize packet")
	}

	var lengthBytes [8]byte
	binary.PutUvarint(lengthBytes[:], uint64(p.Size()))

	var result []byte
	result = append(result, lengthBytes[:]...)
	result = append(result, data...)

	return result, nil
}

// DeserializePacketBackend reads packet from io.Reader.
func DeserializePacketBackend(conn io.Reader) (*PacketBackend, error) {
	lengthBytes := make([]byte, 8)
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}
	lengthReader := bytes.NewBuffer(lengthBytes)
	length, err := binary.ReadUvarint(lengthReader)
	if err != nil {
		return nil, io.ErrUnexpectedEOF
	}

	log.Debugf("[ DeserializePacket ] packet length %d", length)
	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		log.Error("[ DeserializePacket ] couldn't read packet: ", err)
		return nil, err
	}
	log.Debugf("[ DeserializePacket ] read packet")

	msg := &PacketBackend{}
	err = msg.Unmarshal(buf)
	if err != nil {
		log.Error("[ DeserializePacket ] couldn't decode packet: ", err)
		return nil, err
	}

	log.Debugf("[ DeserializePacket ] decoded packet to %#v", msg)

	return msg, nil
}
