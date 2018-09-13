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

// Package message represents message that messagerouter can route
package message

import (
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// BaseMessage base of message class family, do not use it standalone
type baseMessage struct {
	Request core.RecordRef
	Domain  core.RecordRef
}

func (baseMessage) Serialize() (io.Reader, error) {
	panic("Do not use base")
}

func (baseMessage) GetReference() core.RecordRef {
	panic("Do not use base")
}

// MessageType is a enum type of message
type MessageType byte

const (
	baseMessageType            = MessageType(iota)
	CallMethodMessageType      // CallMethodMessage - Simply call method and return result
	CallConstructorMessageType // CallConstructorMessage is a message for calling constructor and obtain its response
	DelegateMessageType        // DelegateMessage is a message for injecting a delegate
)

// GetEmptyMessage constructs specified message
func getEmptyMessage(mt MessageType) (core.Message, error) {
	switch mt {
	case baseMessageType:
		return nil, errors.New("working with message type == 0 is prohibited")
	case CallMethodMessageType:
		return &CallMethodMessage{}, nil
	case CallConstructorMessageType:
		return &CallConstructorMessage{}, nil
	case DelegateMessageType:
		return &DelegateMessage{}, nil
	default:
		return nil, errors.Errorf("unimplemented messagetype %d", mt)
	}
}

// Deserialize returns a message
func Deserialize(buff io.Reader) (core.Message, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short slice for deserialize message")
	}

	m, err := getEmptyMessage(MessageType(b[0]))
	if err != nil {
		return nil, err
	}
	enc := gob.NewDecoder(buff)
	err = enc.Decode(m)
	return m, err
}
