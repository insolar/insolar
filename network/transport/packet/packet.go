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

// NewPingPacket can be used as a shortcut for creating ping packets instead of packet Builder.
func NewPingPacket(sender, receiver *host.Host) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     types.TypePing,
	}
}

// IsValid checks if packet data is a valid structure for current packet type.
func (m *Packet) IsValid() (valid bool) { // nolint: gocyclo
	switch m.Type {
	case types.TypePing:
		valid = true
	case types.TypeFindHost:
		_, valid = m.Data.(*RequestDataFindHost)
	case types.TypeFindValue:
		_, valid = m.Data.(*RequestDataFindValue)
	case types.TypeStore:
		_, valid = m.Data.(*RequestDataStore)
	case types.TypeRPC:
		_, valid = m.Data.(*RequestDataRPC)
	case types.TypeRelay:
		_, valid = m.Data.(*RequestRelay)
	case types.TypeAuthentication:
		_, valid = m.Data.(*RequestAuthentication)
	case types.TypeCheckOrigin:
		_, valid = m.Data.(*RequestCheckOrigin)
	case types.TypeObtainIP:
		_, valid = m.Data.(*RequestObtainIP)
	case types.TypeRelayOwnership:
		_, valid = m.Data.(*RequestRelayOwnership)
	case types.TypeKnownOuterHosts:
		_, valid = m.Data.(*RequestKnownOuterHosts)
	case types.TypeCheckNodePriv:
		_, valid = m.Data.(*RequestCheckNodePriv)
	case types.TypeCascadeSend:
		_, valid = m.Data.(*RequestCascadeSend)
	case types.TypePulse:
		_, valid = m.Data.(*RequestPulse)
	case types.TypeGetRandomHosts:
		_, valid = m.Data.(*RequestGetRandomHosts)
	case types.TypeGetNonce:
		_, valid = m.Data.(*RequestGetNonce)
	case types.TypeCheckSignedNonce:
		_, valid = m.Data.(*RequestCheckSignedNonce)
	case types.TypeExchangeUnsyncLists:
		_, valid = m.Data.(*RequestExchangeUnsyncLists)
	case types.TypeExchangeUnsyncHash:
		_, valid = m.Data.(*RequestExchangeUnsyncHash)
	case types.TypeDisconnect:
		_, valid = m.Data.(*RequestDisconnect)
	default:
		valid = false
	}

	return valid
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
	_, err := conn.Read(lengthBytes)
	if err != nil {
		return nil, err
	}

	lengthReader := bytes.NewBuffer(lengthBytes)
	length, err := binary.ReadUvarint(lengthReader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read variant")
	}

	var reader bytes.Buffer
	for uint64(reader.Len()) < length {
		n, err := reader.ReadFrom(conn)
		if err != nil && n == 0 {
			log.Debugln(err.Error())
		}
	}

	msg := &Packet{}
	dec := gob.NewDecoder(&reader)

	err = dec.Decode(msg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize packet")
	}

	return msg, nil
}

func init() {
	gob.Register(&RequestDataFindHost{})
	gob.Register(&RequestDataFindValue{})
	gob.Register(&RequestDataStore{})
	gob.Register(&RequestDataRPC{})
	gob.Register(&RequestRelay{})
	gob.Register(&RequestAuthentication{})
	gob.Register(&RequestCheckOrigin{})
	gob.Register(&RequestObtainIP{})
	gob.Register(&RequestRelayOwnership{})
	gob.Register(&RequestKnownOuterHosts{})
	gob.Register(&RequestCheckNodePriv{})
	gob.Register(&RequestCascadeSend{})
	gob.Register(&RequestPulse{})
	gob.Register(&RequestGetRandomHosts{})
	gob.Register(&RequestGetNonce{})
	gob.Register(&RequestCheckSignedNonce{})
	gob.Register(&RequestExchangeUnsyncLists{})
	gob.Register(&RequestExchangeUnsyncHash{})
	gob.Register(&RequestDisconnect{})

	gob.Register(&ResponseDataFindHost{})
	gob.Register(&ResponseDataFindValue{})
	gob.Register(&ResponseDataStore{})
	gob.Register(&ResponseDataRPC{})
	gob.Register(&ResponseRelay{})
	gob.Register(&ResponseAuthentication{})
	gob.Register(&ResponseCheckOrigin{})
	gob.Register(&ResponseObtainIP{})
	gob.Register(&ResponseRelayOwnership{})
	gob.Register(&ResponseKnownOuterHosts{})
	gob.Register(&ResponseCheckNodePriv{})
	gob.Register(&ResponseCascadeSend{})
	gob.Register(&ResponsePulse{})
	gob.Register(&ResponseGetRandomHosts{})
	gob.Register(&ResponseGetNonce{})
	gob.Register(&ResponseExchangeUnsyncLists{})
	gob.Register(&ResponseExchangeUnsyncHash{})
	gob.Register(&ResponseDisconnect{})
	gob.Register(&ResponseCheckSignedNonce{})
}
