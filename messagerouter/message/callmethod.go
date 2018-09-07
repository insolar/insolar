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
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
)

// CallMethodMessage - Simply call method and return result
type CallMethodMessage struct {
	BaseMessage
	ObjectRef core.RecordRef
	Request   core.RecordRef
	Method    string
	Arguments core.Arguments
}

// GetReference implements core.Message
func (m *CallMethodMessage) GetReference() core.RecordRef {
	return m.ObjectRef
}

// GetReference implements core.Message
func (m *CallMethodMessage) Serialize() (io.Reader, error) {
	buff := &bytes.Buffer{}
	buff.Write([]byte{byte(CallMethodMessageType)})
	enc := gob.NewEncoder(buff)
	err := enc.Encode(m)
	return buff, err
}
