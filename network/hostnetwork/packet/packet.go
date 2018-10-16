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
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/pkg/errors"
)

//go:generate stringer -type=packetType
type packetType int

const (
	// TypePing is packet type for ping method.
	TypePing packetType = iota + 1
	// TypeStore is packet type for store method.
	TypeStore
	// TypeFindHost is packet type for FindHost method.
	TypeFindHost
	// TypeFindValue is packet type for FindValue method.
	TypeFindValue
	// TypeRPC is packet type for RPC method.
	TypeRPC
	// TypeRelay is packet type for request target to be a relay.
	TypeRelay
	// TypeAuthentication is packet type for authentication between hosts.
	TypeAuthentication
	// TypeCheckOrigin is packet to check originality of some host.
	TypeCheckOrigin
	// TypeObtainIP is packet to get itself IP from another host.
	TypeObtainIP
	// TypeRelayOwnership is packet to say all other hosts that current host have a static IP.
	TypeRelayOwnership
	// TypeKnownOuterHosts is packet to say how much outer hosts current host know.
	TypeKnownOuterHosts
	// TypeCheckNodePriv is packet to check preset node privileges.
	TypeCheckNodePriv
	// TypeCascadeSend is the packet type for the cascade send message feature.
	TypeCascadeSend
	// TypePulse is packet type for the messages received from pulsars.
	TypePulse
	// TypeGetRandomHosts is packet type for the call to get random hosts of the DHT network.
	TypeGetRandomHosts
	// TypeGetNonce is packet to request a nonce to sign it.
	TypeGetNonce
	// TypeCheckSignedNonce is packet to check a signed nonce.
	TypeCheckSignedNonce
	// TypeExchangeUnsyncLists is packet type to exchange unsync lists during consensus
	TypeExchangeUnsyncLists
	// TypeExchangeUnsyncHash is packet type to exchange hash of merged unsync lists during consensus
	TypeExchangeUnsyncHash
	// TypeDisconnect is packet to disconnect from active list.
	TypeDisconnect
)

// RequestID is 64 bit unsigned int request id.
type RequestID uint64

// Packet is DHT packet object.
type Packet struct {
	Sender        *host.Host
	Receiver      *host.Host
	Type          packetType
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
		Type:     TypePing,
	}
}

// IsValid checks if packet data is a valid structure for current packet type.
func (m *Packet) IsValid() (valid bool) { // nolint: gocyclo
	switch m.Type {
	case TypePing:
		valid = true
	case TypeFindHost:
		_, valid = m.Data.(*RequestDataFindHost)
	case TypeFindValue:
		_, valid = m.Data.(*RequestDataFindValue)
	case TypeStore:
		_, valid = m.Data.(*RequestDataStore)
	case TypeRPC:
		_, valid = m.Data.(*RequestDataRPC)
	case TypeRelay:
		_, valid = m.Data.(*RequestRelay)
	case TypeAuthentication:
		_, valid = m.Data.(*RequestAuthentication)
	case TypeCheckOrigin:
		_, valid = m.Data.(*RequestCheckOrigin)
	case TypeObtainIP:
		_, valid = m.Data.(*RequestObtainIP)
	case TypeRelayOwnership:
		_, valid = m.Data.(*RequestRelayOwnership)
	case TypeKnownOuterHosts:
		_, valid = m.Data.(*RequestKnownOuterHosts)
	case TypeCheckNodePriv:
		_, valid = m.Data.(*RequestCheckNodePriv)
	case TypeCascadeSend:
		_, valid = m.Data.(*RequestCascadeSend)
	case TypePulse:
		_, valid = m.Data.(*RequestPulse)
	case TypeGetRandomHosts:
		_, valid = m.Data.(*RequestGetRandomHosts)
	case TypeGetNonce:
		_, valid = m.Data.(*RequestGetNonce)
	case TypeCheckSignedNonce:
		_, valid = m.Data.(*RequestCheckSignedNonce)
	case TypeExchangeUnsyncLists:
		_, valid = m.Data.(*RequestExchangeUnsyncLists)
	case TypeExchangeUnsyncHash:
		_, valid = m.Data.(*RequestExchangeUnsyncHash)
	case TypeDisconnect:
		_, valid = m.Data.(*RequestDisconnect)
	default:
		valid = false
	}

	return valid
}

// IsForMe checks if packet is addressed to our host.
func (m *Packet) IsForMe(origin host.Origin) bool {
	return origin.Contains(m.Receiver) || m.Type == TypePing && origin.IDs[0].Equal(m.Receiver.ID.Bytes()) //origin.Address.Equal(*m.Receiver.Address)
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

	gob.Register(&id.ID{})
}
