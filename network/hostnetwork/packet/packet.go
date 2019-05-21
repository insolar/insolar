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
	"encoding/gob"
	"io"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

// Packet is network packet object for requests and responses.
type Packet struct {
	Sender    *host.Host
	Receiver  *host.Host
	Type      types.PacketType
	RequestID types.RequestID

	TraceID    string
	Data       interface{}
	IsResponse bool
}

func (p *Packet) GetSender() insolar.Reference {
	return p.Sender.NodeID
}

func (p *Packet) GetSenderHost() *host.Host {
	return p.Sender
}

func (p *Packet) GetType() types.PacketType {
	return p.Type
}

func (p *Packet) GetData() interface{} {
	return p.Data
}

func (p *Packet) GetRequestID() types.RequestID {
	return p.RequestID
}

// SerializePacket converts packet to byte slice.
func SerializePacket(q *Packet) ([]byte, error) {
	var msgBuffer bytes.Buffer
	enc := gob.NewEncoder(&msgBuffer)
	err := enc.Encode(q)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize packet")
	}

	length := msgBuffer.Len()

	var lengthBytes [8]byte
	binary.PutUvarint(lengthBytes[:], uint64(length))

	var result []byte
	result = append(result, lengthBytes[:]...)
	result = append(result, msgBuffer.Bytes()...)

	return result, nil
}

// DeserializePacket reads packet from io.Reader.
func DeserializePacket(conn io.Reader) (*Packet, error) {

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

	msg := &Packet{}
	dec := gob.NewDecoder(bytes.NewReader(buf))

	err = dec.Decode(msg)
	if err != nil {
		log.Error("[ DeserializePacket ] couldn't decode packet: ", err)
		return nil, err
	}

	log.Debugf("[ DeserializePacket ] decoded packet to %#v", msg)

	return msg, nil
}

func init() {
	gob.Register(&RequestPulse{})
	gob.Register(&ResponsePulse{})
}
