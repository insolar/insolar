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

	"github.com/insolar/network/node"
)

type messageType int

const (
	// TypePing is message type for ping method
	TypePing = messageType(iota + 1)
	// TypeStore is message type for store method
	TypeStore
	// TypeFindNode is message type for FindNode method
	TypeFindNode
	// TypeFindValue is message type for FindValue method
	TypeFindValue
	// TypeRPC is message type for RPC method
	TypeRPC
	// TypeRelay is message type for request target to be a relay
	TypeRelay
)

// RequestID is 64 bit unsigned int request id
type RequestID uint64

// Message is DHT message object
type Message struct {
	Sender    *node.Node
	Receiver  *node.Node
	Type      messageType
	RequestID RequestID

	Data       interface{}
	Error      error
	IsResponse bool
}

// NewPingMessage can be used as a shortcut for creating ping messages instead of message Builder
func NewPingMessage(sender, receiver *node.Node) *Message {
	return &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypePing,
	}
}

// NewRelayMessage uses for send a command to target node to make it as relay
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

// IsValid checks if message data is a valid structure for current message type
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
	default:
		valid = false
	}

	return valid
}

// IsForMe checks if message is addressed to our node
func (m *Message) IsForMe(origin node.Origin) bool {
	return origin.Contains(m.Receiver) || m.Type == TypePing && origin.Address.Equal(*m.Receiver.Address)
}

// SerializeMessage converts message to byte slice
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

// DeserializeMessage reads message from io.Reader
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

	msgBytes := make([]byte, length)
	_, err = conn.Read(msgBytes)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(msgBytes)
	msg := &Message{}
	dec := gob.NewDecoder(reader)

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

	gob.Register(&ResponseDataFindNode{})
	gob.Register(&ResponseDataFindValue{})
	gob.Register(&ResponseDataStore{})
	gob.Register(&ResponseDataRPC{})
	gob.Register(&ResponseRelay{})
}
