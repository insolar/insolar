/*
 *    Copyright 2018 INS Ecosystem
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

package message

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"log"

	"github.com/insolar/insolar/network/host/node"
)

type messageType int

const (
	// TypePing is message type for ping method.
	TypePing = messageType(iota + 1)
	// TypeStore is message type for store method.
	TypeStore
	// TypeFindNode is message type for FindNode method.
	TypeFindNode
	// TypeFindValue is message type for FindValue method.
	TypeFindValue
	// TypeRPC is message type for RPC method.
	TypeRPC
	// TypeRelay is message type for request target to be a relay.
	TypeRelay
	// TypeAuth is message type for authentication between nodes.
	TypeAuth
	// TypeCheckOrigin is message to check originality of some node.
	TypeCheckOrigin
	// TypeObtainIP is message to get itself IP from another node.
	TypeObtainIP
	// TypeRelayOwnership is message to say all other nodes that current node have a static IP.
	TypeRelayOwnership
	// TypeKnownOuterNodes is message to say how much outer nodes current node know.
	TypeKnownOuterNodes
)

// RequestID is 64 bit unsigned int request id.
type RequestID uint64

// Message is DHT message object.
type Message struct {
	Sender        *node.Node
	Receiver      *node.Node
	Type          messageType
	RequestID     RequestID
	RemoteAddress string

	Data       interface{}
	Error      error
	IsResponse bool
}

// NewPingMessage can be used as a shortcut for creating ping messages instead of message Builder.
func NewPingMessage(sender, receiver *node.Node) *Message {
	return &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypePing,
	}
}

// NewRelayMessage uses for send a command to target node to make it as relay.
func NewRelayMessage(command CommandType, sender, receiver *node.Node) *Message {
	return &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeRelay,
		Data: &RequestRelay{
			Command: command,
		},
	}
}

// NewAuthMessage uses for starting authentication.
func NewAuthMessage(command CommandType, sender, receiver *node.Node) *Message {
	return &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeAuth,
		Data:     &RequestAuth{Command: command},
	}
}

// NewCheckOriginMessage uses for check originality.
func NewCheckOriginMessage(sender, receiver *node.Node) *Message {
	return &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeCheckOrigin,
		Data:     &RequestCheckOrigin{},
	}
}

// NewObtainIPMessage uses for get self IP.
func NewObtainIPMessage(sender, receiver *node.Node) *Message {
	return &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeObtainIP,
		Data:     &RequestObtainIP{},
	}
}

// NewRelayOwnershipMessage uses for relay ownership request.
func NewRelayOwnershipMessage(sender, receiver *node.Node, ready bool) *Message {
	return &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeRelayOwnership,
		Data:     &RequestRelayOwnership{Ready: ready},
	}
}

// NewKnownOuterNodesMessage uses to notify all nodes in home subnet about known outer nodes.
func NewKnownOuterNodesMessage(sender, receiver *node.Node, nodes int) *Message {
	return &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypeKnownOuterNodes,
		Data: &RequestKnownOuterNodes{
			ID:         sender.ID.String(),
			OuterNodes: nodes,
		},
	}
}

// IsValid checks if message data is a valid structure for current message type.
func (m *Message) IsValid() (valid bool) {
	switch m.Type {
	case TypePing:
		valid = true
	case TypeFindNode:
		_, valid = m.Data.(*RequestDataFindNode)
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
	case TypeKnownOuterNodes:
		_, valid = m.Data.(*RequestKnownOuterNodes)
	default:
		valid = false
	}

	return valid
}

// IsForMe checks if message is addressed to our node.
func (m *Message) IsForMe(origin node.Origin) bool {
	return origin.Contains(m.Receiver) || m.Type == TypePing && origin.Address.Equal(*m.Receiver.Address)
}

// SerializeMessage converts message to byte slice.
func SerializeMessage(q *Message) ([]byte, error) {
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

// DeserializeMessage reads message from io.Reader.
func DeserializeMessage(conn io.Reader) (*Message, error) {

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

	msg := &Message{}
	dec := gob.NewDecoder(&reader)

	err = dec.Decode(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func init() {
	gob.Register(&RequestDataFindNode{})
	gob.Register(&RequestDataFindValue{})
	gob.Register(&RequestDataStore{})
	gob.Register(&RequestDataRPC{})
	gob.Register(&RequestRelay{})
	gob.Register(&RequestAuth{})
	gob.Register(&RequestCheckOrigin{})
	gob.Register(&RequestObtainIP{})
	gob.Register(&RequestRelayOwnership{})
	gob.Register(&RequestKnownOuterNodes{})

	gob.Register(&ResponseDataFindNode{})
	gob.Register(&ResponseDataFindValue{})
	gob.Register(&ResponseDataStore{})
	gob.Register(&ResponseDataRPC{})
	gob.Register(&ResponseRelay{})
	gob.Register(&ResponseAuth{})
	gob.Register(&ResponseCheckOrigin{})
	gob.Register(&ResponseObtainIP{})
	gob.Register(&ResponseRelayOwnership{})
	gob.Register(&ResponseKnownOuterNodes{})
}
