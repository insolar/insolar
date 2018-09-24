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

// Type is a enum type of event
type Type byte

const (
	baseEventType            = Type(iota)
	CallMethodEventType      // CallMethodEvent - Simply call method and return result
	CallConstructorEventType // CallConstructorEvent is a event for calling constructor and obtain its reaction

	// Ledger
	TypeGetCode                // TypeGetCode - retrieve code from storage.
	TypeGetClass               // TypeGetClass - latest state of the class known to storage
	TypeGetObject              // TypeGetObject returns descriptor for latest state of the object known to storage.
	TypeGetDelegate            // TypeGetDelegate returns descriptor for latest state of the object known to storage.
	TypeDeclareType            // TypeGetDelegate creates new type.
	TypeDeployCode             // TypeDeployCode creates new code.
	TypeActivateClass          // TypeActivateClass activates class.
	TypeDeactivateClass        // TypeDeactivateClass deactivates class.
	TypeUpdateClass            // TypeUpdateClass amends class.
	TypeActivateObject         // TypeActivateObject activates object.
	TypeActivateObjectDelegate // TypeActivateObjectDelegate similar to ActivateObjType but it creates object as parent's delegate of provided class.
	TypeDeactivateObject       // TypeDeactivateObject deactivates object.
	TypeUpdateObject           // TypeUpdateObject amends object.
)

// GetEmptyMessage constructs specified event
func getEmptyEvent(mt Type) (core.Event, error) {
	switch mt {
	case baseEventType:
		return nil, errors.New("working with event type == 0 is prohibited")
	case CallMethodEventType:
		return &CallMethodEvent{}, nil
	case CallConstructorEventType:
		return &CallConstructorEvent{}, nil

	// Ledger
	case TypeGetCode:
		return &GetCode{}, nil
	case TypeGetClass:
		return &GetClass{}, nil
	case TypeGetObject:
		return &GetObject{}, nil
	case TypeGetDelegate:
		return &GetDelegate{}, nil
	case TypeDeclareType:
		return &DeclareType{}, nil
	case TypeDeployCode:
		return &DeployCode{}, nil
	case TypeActivateClass:
		return &ActivateClass{}, nil
	case TypeDeactivateClass:
		return &DeactivateClass{}, nil
	case TypeUpdateClass:
		return &UpdateClass{}, nil
	case TypeActivateObject:
		return &ActivateObject{}, nil
	case TypeActivateObjectDelegate:
		return &ActivateObjectDelegate{}, nil
	case TypeDeactivateObject:
		return &DeactivateObject{}, nil
	case TypeUpdateObject:
		return &UpdateObject{}, nil
	default:
		return nil, errors.Errorf("unimplemented event type %d", mt)
	}
}

func serialize(event core.Event, t Type) (io.Reader, error) {
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

	event, err := getEmptyEvent(Type(b[0]))
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
}
