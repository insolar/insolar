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

// Package message represents message that messagebus can route
package message

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

const (
	// Logicrunner

	TypeCallMethod      = core.MessageType(iota) // TypeCallMethod calls method and returns result
	TypeCallConstructor                          // TypeCallConstructor is a message for calling constructor and obtain its reply

	// Ledger

	TypeGetCode                // TypeGetCode retrieves code from storage.
	TypeGetClass               // TypeGetClass retrieves class from storage.
	TypeGetObject              // TypeGetObject retrieves object from storage.
	TypeGetDelegate            // TypeGetDelegate retrieves object represented as provided class.
	TypeDeclareType            // TypeDeclareType creates new type.
	TypeDeployCode             // TypeDeployCode creates new code.
	TypeActivateClass          // TypeActivateClass activates class.
	TypeDeactivateClass        // TypeDeactivateClass deactivates class.
	TypeUpdateClass            // TypeUpdateClass amends class.
	TypeActivateObject         // TypeActivateObject activates object.
	TypeActivateObjectDelegate // TypeActivateObjectDelegate similar to ActivateObjType but it creates object as parent's delegate of provided class.
	TypeDeactivateObject       // TypeDeactivateObject deactivates object.
	TypeUpdateObject           // TypeUpdateObject amends object.
)

// GetEmptyMessage constructs specified message
func getEmptyMessage(mt core.MessageType) (core.Message, error) {
	switch mt {
	// Logicrunner
	case TypeCallMethod:
		return &CallMethod{}, nil
	case TypeCallConstructor:
		return &CallConstructor{}, nil
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
		return nil, errors.Errorf("unimplemented message type %d", mt)
	}
}

// Serialize returns encoded message.
func Serialize(msg core.Message) (io.Reader, error) {
	buff := &bytes.Buffer{}
	_, err := buff.Write([]byte{byte(msg.Type())})
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(buff)
	err = enc.Encode(msg)
	return buff, err
}

// Deserialize returns decoded message.
func Deserialize(buff io.Reader) (core.Message, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short slice for deserialize message")
	}

	msg, err := getEmptyMessage(core.MessageType(b[0]))
	if err != nil {
		return nil, err
	}
	enc := gob.NewDecoder(buff)
	err = enc.Decode(msg)
	return msg, err
}

func init() {
	// Logicrunner
	gob.Register(&CallConstructor{})
	gob.Register(&CallMethod{})
	// Ledger
	gob.Register(&GetCode{})
	gob.Register(&GetClass{})
	gob.Register(&GetObject{})
	gob.Register(&GetDelegate{})
	gob.Register(&DeclareType{})
	gob.Register(&DeployCode{})
	gob.Register(&ActivateClass{})
	gob.Register(&DeactivateClass{})
	gob.Register(&UpdateClass{})
	gob.Register(&ActivateObject{})
	gob.Register(&ActivateObjectDelegate{})
	gob.Register(&DeactivateObject{})
	gob.Register(&UpdateObject{})
}
