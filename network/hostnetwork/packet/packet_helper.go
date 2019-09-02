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
	"strconv"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/pkg/errors"
)

func (p *Packet) SetRequest(request interface{}) {
	var r isRequest_Request
	switch t := request.(type) {
	case *RPCRequest:
		r = &Request_RPC{t}
	case *PulseRequest:
		r = &Request_Pulse{t}
	case *BootstrapRequest:
		r = &Request_Bootstrap{t}
	case *AuthorizeRequest:
		r = &Request_Authorize{t}
	case *SignCertRequest:
		r = &Request_SignCert{t}
	case *UpdateScheduleRequest:
		r = &Request_UpdateSchedule{t}
	case *ReconnectRequest:
		r = &Request_Reconnect{t}
	default:
		panic("Request payload is not a valid protobuf struct!")
	}
	p.Payload = &Packet_Request{Request: &Request{Request: r}}
}

func (p *Packet) SetResponse(response interface{}) {
	var r isResponse_Response
	switch t := response.(type) {
	case *RPCResponse:
		r = &Response_RPC{t}
	case *BasicResponse:
		r = &Response_Basic{t}
	case *BootstrapResponse:
		r = &Response_Bootstrap{t}
	case *AuthorizeResponse:
		r = &Response_Authorize{t}
	case *SignCertResponse:
		r = &Response_SignCert{t}
	case *ErrorResponse:
		r = &Response_Error{t}
	case *UpdateScheduleResponse:
		r = &Response_UpdateSchedule{t}
	case *ReconnectResponse:
		r = &Response_Reconnect{t}
	default:
		panic("Response payload is not a valid protobuf struct!")
	}
	p.Payload = &Packet_Response{Response: &Response{Response: r}}
}

func (p *Packet) GetType() types.PacketType {
	// TODO: make p.Type of type PacketType instead of uint32
	return types.PacketType(p.Type)
}

func (p *Packet) GetSender() insolar.Reference {
	return p.Sender.NodeID
}

func (p *Packet) GetSenderHost() *host.Host {
	return p.Sender
}

func (p *Packet) GetRequestID() types.RequestID {
	return types.RequestID(p.RequestID)
}

func (p *Packet) IsResponse() bool {
	return p.GetResponse() != nil
}

// SerializePacket converts packet to byte slice.
func SerializePacket(p *Packet) ([]byte, error) {
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

func DeserializePacketRaw(conn io.Reader) (*ReceivedPacket, error) {
	reader := NewCapturingReader(conn)

	lengthBytes := make([]byte, 8)
	if _, err := io.ReadFull(reader, lengthBytes); err != nil {
		return nil, err
	}
	lengthReader := bytes.NewReader(lengthBytes)
	length, err := binary.ReadUvarint(lengthReader)
	if err != nil {
		return nil, io.ErrUnexpectedEOF
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, errors.Wrap(err, "failed to read packet")
	}

	msg := &Packet{}
	err = msg.Unmarshal(buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode packet")
	}

	receivedPacket := NewReceivedPacket(msg, reader.Captured())
	return receivedPacket, nil
}

// DeserializePacket reads packet from io.Reader.
func DeserializePacket(logger insolar.Logger, conn io.Reader) (*ReceivedPacket, error) {
	receivedPacket, err := DeserializePacketRaw(conn)
	if err != nil {
		return nil, err
	}
	logger.Debugf("[ DeserializePacket ] decoded packet to %s", receivedPacket.DebugString())
	return receivedPacket, nil
}

func (p *Packet) DebugString() string {
	if p == nil {
		return "nil"
	}
	return `&Packet{` +
		`Sender:` + p.Sender.String() + `,` +
		`Receiver:` + p.Receiver.String() + `,` +
		`RequestID:` + strconv.FormatUint(p.RequestID, 10) + `,` +
		`TraceID:` + p.TraceID + `,` +
		`Type:` + p.GetType().String() + `,` +
		`IsResponse:` + strconv.FormatBool(p.IsResponse()) + `,` +
		`}`
}

func NewPacket(sender, receiver *host.Host, packetType types.PacketType, id uint64) *Packet {
	return &Packet{
		// Polymorph field should be non-default so we have first byte 0x80 in serialized representation
		Polymorph: 1,
		Sender:    sender,
		Receiver:  receiver,
		Type:      uint32(packetType),
		RequestID: id,
	}
}

type CapturingReader struct {
	io.Reader
	buffer bytes.Buffer
}

func NewCapturingReader(reader io.Reader) *CapturingReader {
	return &CapturingReader{Reader: reader}
}

func (r *CapturingReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.buffer.Write(p)
	return n, err
}

func (r *CapturingReader) Captured() []byte {
	return r.buffer.Bytes()
}
