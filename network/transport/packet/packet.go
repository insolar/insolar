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

package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

// RequestID is 64 bit unsigned int request id.
type RequestID uint64

// Packet is DHT packet object.
type Packet struct {
	Sender        *host.Host
	Receiver      *host.Host
	Type          types.PacketType
	RequestID     RequestID
	RemoteAddress string

	Data       interface{}
	Error      error
	IsResponse bool
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
	n, err := conn.Read(lengthBytes)
	if err != nil {
		return nil, err
	}

	log.Debugf("[ DeserializePacket ] read %d bytes", n)

	lengthReader := bytes.NewBuffer(lengthBytes)
	length, err := binary.ReadUvarint(lengthReader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read variant")
	}

	log.Debugf("[ DeserializePacket ] packet length %d", length)

	buf := make([]byte, length)

	var readLength int
	for readLength < int(length) {
		n, err = conn.Read(buf[readLength:])
		readLength = readLength + n
		log.Debugf("read %d bytes", n)
		if err != nil && n == 0 {
			log.Debugln(err.Error())
		}
	}

	msg := &Packet{}
	dec := gob.NewDecoder(bytes.NewReader(buf))

	err = dec.Decode(msg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize packet")
	}

	return msg, nil
}

func init() {
	gob.Register(&RequestPulse{})
	gob.Register(&RequestGetRandomHosts{})

	gob.Register(&ResponsePulse{})
	gob.Register(&ResponseGetRandomHosts{})
}
