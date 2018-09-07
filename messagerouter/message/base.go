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

// Package Message represents message that messagerouter can route
package message

import (
	"fmt"

	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// BaseMessage base of message class family, do not use it standalone
type BaseMessage struct {
	Request core.RecordRef
	Domain  core.RecordRef
}

func (m *BaseMessage) Serialize() (io.Reader, error) {
	panic("BaseMessage is not usable object")
}

// MessageType is a enum type of message
type MessageType byte

const (
	BaseMessageType MessageType = iota
	CallMethodMessageType
	CallConstructorMessageType
	MessageTypesCount
)

// GetEmptyMessage constructs specified message
func GetEmptyMessage(mt MessageType) core.Message {
	switch mt {
	case 0:
		panic("working with message type == 0 is prohibited")
	case CallMethodMessageType:
		return &CallMethodMessage{}
	case CallConstructorMessageType:
		return &CallConstructorMessage{}
	default:
		panic(fmt.Sprintf("unimplemented messagetype %d", mt))
	}
}

// Deserialize returns a message
func Deserialize(buff io.Reader) (core.Message, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short slice for deserialize message")
	}

	m := GetEmptyMessage(MessageType(b[0]))
	enc := gob.NewDecoder(buff)
	err = enc.Decode(m)
	return m, err
}
