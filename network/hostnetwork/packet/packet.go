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
	"log"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
)

type packetType int

const (
	// TypePing is packet type for ping method.
	TypePing = packetType(iota + 1)
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
	// TypeAuth is packet type for authentication between hosts.
	TypeAuth
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
	// TypeCascadeSend is the packet type for the cascade send message feature
	TypeCascadeSend
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

// NewCheckNodePrivPacket used to create message for check node privileges request.
func NewCheckNodePrivPacket(sender, receiver *host.Host, roleKey string) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeCheckNodePriv,
		Data: &RequestCheckNodePriv{
			RoleKey: roleKey,
		},
	}
}

// NewPingPacket can be used as a shortcut for creating ping packets instead of packet Builder.
func NewPingPacket(sender, receiver *host.Host) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypePing,
	}
}

// NewRelayPacket uses for send a command to target host to make it as relay.
func NewRelayPacket(command CommandType, sender, receiver *host.Host) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeRelay,
		Data: &RequestRelay{
			Command: command,
		},
	}
}

// NewAuthPacket uses for starting authentication.
func NewAuthPacket(command CommandType, sender, receiver *host.Host) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeAuth,
		Data:     &RequestAuth{Command: command},
	}
}

// NewCheckOriginPacket uses for check originality.
func NewCheckOriginPacket(sender, receiver *host.Host) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeCheckOrigin,
		Data:     &RequestCheckOrigin{},
	}
}

// NewObtainIPPacket uses for get self IP.
func NewObtainIPPacket(sender, receiver *host.Host) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeObtainIP,
		Data:     &RequestObtainIP{},
	}
}

// NewRelayOwnershipPacket uses for relay ownership request.
func NewRelayOwnershipPacket(sender, receiver *host.Host, ready bool) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeRelayOwnership,
		Data:     &RequestRelayOwnership{Ready: ready},
	}
}

// NewKnownOuterHostsPacket uses to notify all hosts in home subnet about known outer hosts.
func NewKnownOuterHostsPacket(sender, receiver *host.Host, hosts int) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeKnownOuterHosts,
		Data: &RequestKnownOuterHosts{
			ID:         sender.ID.String(),
			OuterHosts: hosts,
		},
	}
}

// IsValid checks if packet data is a valid structure for current packet type.
func (m *Packet) IsValid() (valid bool) {
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
	case TypeAuth:
		_, valid = m.Data.(*RequestAuth)
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
	default:
		valid = false
	}

	return valid
}

// IsForMe checks if packet is addressed to our host.
func (m *Packet) IsForMe(origin host.Origin) bool {
	return origin.Contains(m.Receiver) || m.Type == TypePing && origin.Address.Equal(*m.Receiver.Address)
}

// SerializePacket converts packet to byte slice.
func SerializePacket(q *Packet) ([]byte, error) {
	var msgBuffer bytes.Buffer
	enc := gob.NewEncoder(&msgBuffer)
	err := enc.Encode(q)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	var reader bytes.Buffer
	for uint64(reader.Len()) < length {
		n, _ := reader.ReadFrom(conn)
		if err != nil && n == 0 {
			log.Println(err.Error())
		}
	}

	msg := &Packet{}
	dec := gob.NewDecoder(&reader)

	err = dec.Decode(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func init() {
	gob.Register(&RequestDataFindHost{})
	gob.Register(&RequestDataFindValue{})
	gob.Register(&RequestDataStore{})
	gob.Register(&RequestDataRPC{})
	gob.Register(&RequestRelay{})
	gob.Register(&RequestAuth{})
	gob.Register(&RequestCheckOrigin{})
	gob.Register(&RequestObtainIP{})
	gob.Register(&RequestRelayOwnership{})
	gob.Register(&RequestKnownOuterHosts{})
	gob.Register(&RequestCheckNodePriv{})
	gob.Register(&RequestCascadeSend{})

	gob.Register(&ResponseDataFindHost{})
	gob.Register(&ResponseDataFindValue{})
	gob.Register(&ResponseDataStore{})
	gob.Register(&ResponseDataRPC{})
	gob.Register(&ResponseRelay{})
	gob.Register(&ResponseAuth{})
	gob.Register(&ResponseCheckOrigin{})
	gob.Register(&ResponseObtainIP{})
	gob.Register(&ResponseRelayOwnership{})
	gob.Register(&ResponseKnownOuterHosts{})
	gob.Register(&ResponseCheckNodePriv{})
	gob.Register(&ResponseCascadeSend{})

	gob.Register(&id.ID{})
}
