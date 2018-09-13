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

package message

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
)

// CallConstructorMessage is a message for calling constructor and obtain its response
type CallConstructorMessage struct {
	baseMessage
	ClassRef  core.RecordRef
	Name      string
	Arguments core.Arguments
}

// GetOperatingRole returns operating jet role for given message type.
func (m *CallConstructorMessage) GetOperatingRole() core.JetRole {
	return core.RoleVirtualExecutor
}

// Get reference returns referenced object.
func (m *CallConstructorMessage) GetReference() core.RecordRef {
	return m.ClassRef
}

// Serialize serializes message.
func (m *CallConstructorMessage) Serialize() (io.Reader, error) {
	buff := &bytes.Buffer{}
	buff.Write([]byte{byte(CallConstructorMessageType)})
	enc := gob.NewEncoder(buff)
	err := enc.Encode(m)
	return buff, err
}
