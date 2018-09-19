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

// Package event represents event that eventbus can route
package event

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// BaseMessage base of event class family, do not use it standalone
type baseEvent struct {
	Request core.RecordRef
	Domain  core.RecordRef
}

func (baseEvent) Serialize() (io.Reader, error) {
	panic("Do not use base")
}

func (baseEvent) GetReference() core.RecordRef {
	panic("Do not use base")
}

// EventType is a enum type of event
type EventType byte

const (
	baseEventType            = EventType(iota)
	CallMethodEventType      // CallMethodEvent - Simply call method and return result
	CallConstructorEventType // CallConstructorEvent is a event for calling constructor and obtain its response
	DelegateEventType        // DelegateEvent is a event for injecting a delegate
	ChildEventType           // ChildEvent is a event for saving a child
	UpdateObjectEventType    // UpdateObjectEvent is a event for updating an object
	GetObjectEventType       // GetObjectEvent is a event for retrieving an object
)

// GetEmptyMessage constructs specified event
func getEmptyEvent(mt EventType) (core.Event, error) {
	switch mt {
	case baseEventType:
		return nil, errors.New("working with event type == 0 is prohibited")
	case CallMethodEventType:
		return &CallMethodEvent{}, nil
	case CallConstructorEventType:
		return &CallConstructorEvent{}, nil
	case DelegateEventType:
		return &DelegateMessage{}, nil
	case ChildEventType:
		return &ChildMessage{}, nil
	case UpdateObjectEventType:
		return &UpdateObjectMessage{}, nil
	case GetObjectEventType:
		return &GetObjectMessage{}, nil
	default:
		return nil, errors.Errorf("unimplemented event type %d", mt)
	}
}

func serialize(event core.Event, t EventType) (io.Reader, error) {
	buff := &bytes.Buffer{}
	_, err := buff.Write([]byte{byte(t)})
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(buff)
	err = enc.Encode(event)
	return buff, err
}

// Deserialize returns a event
func Deserialize(buff io.Reader) (core.Event, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short slice for deserialize event")
	}

	event, err := getEmptyEvent(EventType(b[0]))
	if err != nil {
		return nil, err
	}
	enc := gob.NewDecoder(buff)
	err = enc.Decode(event)
	return event, err
}

func init() {
	gob.Register(&CallConstructorEvent{})
	gob.Register(&CallMethodEvent{})
	gob.Register(&DelegateMessage{})
	gob.Register(&ChildMessage{})
	gob.Register(&UpdateObjectMessage{})
	gob.Register(&GetObjectMessage{})
}
